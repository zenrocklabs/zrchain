package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"maps"
	"math/big"
	"net/http"
	"reflect"
	"sort"
	"sync"
	"time"

	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solrock/generated/rock_spl_token"
	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solzenbtc/generated/zenbtc_spl_token"
	"github.com/Zenrock-Foundation/zrchain/v6/go-client"
	neutrino "github.com/Zenrock-Foundation/zrchain/v6/sidecar/neutrino"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
	"github.com/beevik/ntp"
	sdkBech32 "github.com/cosmos/cosmos-sdk/types/bech32"
	aggregatorv3 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	zenbtc "github.com/zenrocklabs/zenbtc/bindings"
	zenbtctypes "github.com/zenrocklabs/zenbtc/x/zenbtc/types"
	middleware "github.com/zenrocklabs/zenrock-avs/contracts/bindings/ZrServiceManager"

	validationkeeper "github.com/Zenrock-Foundation/zrchain/v6/x/validation/keeper"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	zentptypes "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	solana "github.com/gagliardetto/solana-go"
	solanagoSystem "github.com/gagliardetto/solana-go/programs/system"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	jsonrpc "github.com/gagliardetto/solana-go/rpc/jsonrpc"
	// Added for bin.Marshal
)

func NewOracle(
	config sidecartypes.Config,
	ethClient *ethclient.Client,
	neutrinoServer *neutrino.NeutrinoServer,
	solanaClient *solrpc.Client,
	zrChainQueryClient *client.QueryClient,
	debugMode bool,
	skipInitialWait bool,
) *Oracle {
	o := &Oracle{
		stateCache:         make([]sidecartypes.OracleState, 0),
		Config:             config,
		EthClient:          ethClient,
		neutrinoServer:     neutrinoServer,
		solanaClient:       solanaClient,
		zrChainQueryClient: zrChainQueryClient,
		DebugMode:          debugMode,
		SkipInitialWait:    skipInitialWait,
	}

	// Load initial state from cache file
	latestDiskState, historicalStates, err := loadStateDataFromFile(o.Config.StateFile)
	if err != nil {
		slog.Error("Critical error loading state from file, initializing with empty state", "file", o.Config.StateFile, "error", err)
		o.currentState.Store(&EmptyOracleState)
		o.stateCache = []sidecartypes.OracleState{EmptyOracleState}
		// lastSol*SigStr fields will remain empty strings (zero value)
	} else {
		if latestDiskState != nil {
			o.currentState.Store(latestDiskState)
			o.stateCache = historicalStates
			o.lastSolRockMintSigStr = latestDiskState.LastSolRockMintSig
			o.lastSolZenBTCMintSigStr = latestDiskState.LastSolZenBTCMintSig
			o.lastSolZenBTCBurnSigStr = latestDiskState.LastSolZenBTCBurnSig
			o.lastSolRockBurnSigStr = latestDiskState.LastSolRockBurnSig
			slog.Info("Loaded state from file",
				"rockMintSig", o.lastSolRockMintSigStr,
				"zenBTCMintSig", o.lastSolZenBTCMintSigStr,
				"zenBTCBurnSig", o.lastSolZenBTCBurnSigStr,
				"rockBurnSig", o.lastSolRockBurnSigStr)
		} else {
			// File didn't exist, was empty, or had non-critical parse issues treated as fresh start
			slog.Info("State file not found or empty/invalid, initializing with empty state", "file", o.Config.StateFile)
			o.currentState.Store(&EmptyOracleState)
			o.stateCache = []sidecartypes.OracleState{EmptyOracleState}
			// lastSol*SigStr fields will remain empty strings
		}
	}

	// Initialize the function fields with the real implementations
	o.getSolanaZenBTCBurnEventsFn = o.getSolanaZenBTCBurnEvents
	o.getSolanaRockBurnEventsFn = o.getSolanaRockBurnEvents
	o.rpcCallBatchFn = o.solanaClient.RPCCallBatch
	o.getTransactionFn = o.solanaClient.GetTransaction
	o.getSignaturesForAddressFn = o.solanaClient.GetSignaturesForAddressWithOpts

	return o
}

func (o *Oracle) runOracleMainLoop(ctx context.Context) error {
	serviceManager, err := middleware.NewContractZrServiceManager(
		common.HexToAddress(sidecartypes.ServiceManagerAddresses[o.Config.Network]),
		o.EthClient,
	)
	if err != nil {
		return fmt.Errorf("failed to create contract instance: %w", err)
	}
	zenBTCControllerHolesky, err := zenbtc.NewZenBTController(
		common.HexToAddress(sidecartypes.ZenBTCControllerAddresses[o.Config.Network]),
		o.EthClient,
	)
	if err != nil {
		return fmt.Errorf("failed to create contract instance: %w", err)
	}
	mainnetEthClient, btcPriceFeed, ethPriceFeed := o.initPriceFeed()

	// Initial alignment: Fetch NTP time once at startup
	ntpTime, err := ntp.Time("time.google.com")
	if err != nil {
		// If NTP fails at startup, panic. Sidecars require time sync to establish consensus.
		slog.Error("Failed to fetch NTP time at startup. Cannot proceed.", "error", err)
		panic(fmt.Sprintf("FATAL: Failed to fetch NTP time at startup: %v. Cannot proceed.", err))
	}

	mainLoopTickerIntervalDuration := sidecartypes.MainLoopTickerInterval

	// Align the start time to the nearest MainLoopTickerInterval.
	// This runs only if NTP succeeded (checked by the panic above) and skipInitialWait is false
	if !o.SkipInitialWait {
		alignedStart := ntpTime.Truncate(mainLoopTickerIntervalDuration).Add(mainLoopTickerIntervalDuration)
		initialSleep := time.Until(alignedStart)
		if initialSleep > 0 {
			slog.Info("Initial alignment: Sleeping until start ticker.",
				"sleepDuration", initialSleep.Round(time.Millisecond),
				"alignedStart", alignedStart.Format("15:04:05.00"))
			time.Sleep(initialSleep)
		}
	} else {
		slog.Info("Skipping initial alignment wait due to --skip-initial-wait flag. Firing initial tick immediately.")
		go o.processOracleTick(serviceManager, zenBTCControllerHolesky, btcPriceFeed, ethPriceFeed, mainnetEthClient, time.Now(), mainLoopTickerIntervalDuration)
	}

	mainLoopTicker := time.NewTicker(mainLoopTickerIntervalDuration)
	defer mainLoopTicker.Stop()
	o.mainLoopTicker = mainLoopTicker
	slog.Info("Ticker synced, awaiting initial oracle data fetch", "interval", mainLoopTickerIntervalDuration)

	for {
		select {
		case <-ctx.Done():
			return nil
		case tickTime := <-o.mainLoopTicker.C:
			o.processOracleTick(serviceManager, zenBTCControllerHolesky, btcPriceFeed, ethPriceFeed, mainnetEthClient, tickTime, mainLoopTickerIntervalDuration)
		}
	}
}

func (o *Oracle) processOracleTick(
	serviceManager *middleware.ContractZrServiceManager,
	zenBTCControllerHolesky *zenbtc.ZenBTController,
	btcPriceFeed *aggregatorv3.AggregatorV3Interface,
	ethPriceFeed *aggregatorv3.AggregatorV3Interface,
	mainnetEthClient *ethclient.Client,
	tickTime time.Time,
	mainLoopTickerIntervalDuration time.Duration,
) {
	successfulFetch := true
	newState, err := o.fetchAndProcessState(serviceManager, zenBTCControllerHolesky, btcPriceFeed, ethPriceFeed, mainnetEthClient)
	if err != nil {
		slog.Error("Error fetching and processing state", "error", err)
		successfulFetch = false
	}

	// --- Intra-loop NTP check and wait (with fallback to ticker time) ---
	var sleepDuration time.Duration
	var nextIntervalMark time.Time
	alignmentSource := "NTP"

	// Attempt to fetch current NTP time *after* processing
	ntpTimeNow, err := ntp.Time("time.google.com")
	if err != nil {
		// NTP Failed: Fallback to using the captured ticker time
		slog.Warn("Error fetching NTP time for alignment. Falling back to ticker time.", "error", err)
		alignmentSource = "Local Ticker Fallback"
		// Calculate the next interval boundary based on when the ticker fired.
		nextIntervalMark = tickTime.Truncate(mainLoopTickerIntervalDuration).Add(mainLoopTickerIntervalDuration)
	} else {
		// NTP Succeeded: Calculate alignment based on NTP time.
		nextIntervalMark = ntpTimeNow.Truncate(mainLoopTickerIntervalDuration).Add(mainLoopTickerIntervalDuration)
	}

	// Calculate how long to sleep until the calculated mark
	sleepDuration = time.Until(nextIntervalMark)

	if sleepDuration > 0 {
		slog.Info("State fetched. Waiting until next aligned interval mark to apply update.",
			"sleepDuration", sleepDuration.Round(time.Millisecond),
			"alignmentSource", alignmentSource,
			"nextIntervalMark", nextIntervalMark.Format("15:04:05.00"))
		time.Sleep(sleepDuration)
	} else {
		// If fetching took longer than the interval OR NTP failed and ticker time also leads to negative sleep, log a warning.
		slog.Warn("State fetching took too long relative to alignment. Update applied immediately.", "alignmentSource", alignmentSource)
	}
	// --- End of intra-loop wait ---

	// Send the fetched state exactly at the interval mark (or immediately if delayed)
	if successfulFetch {
		slog.Info("Received AVS contract state for", "network", sidecartypes.NetworkNames[o.Config.Network], "block", newState.EthBlockHeight)
		slog.Info("Received prices", "ROCK/USD", newState.ROCKUSDPrice, "BTC/USD", newState.BTCUSDPrice, "ETH/USD", newState.ETHUSDPrice)
		o.currentState.Store(&newState)
		o.CacheState()
	}

	// Clean up burn events *after* sending state update
	o.cleanUpBurnEvents()
	// Clean up mint events *after* sending state update
	o.cleanUpMintEvents()
}

func (o *Oracle) fetchAndProcessState(
	serviceManager *middleware.ContractZrServiceManager,
	zenBTCControllerHolesky *zenbtc.ZenBTController,
	btcPriceFeed *aggregatorv3.AggregatorV3Interface,
	ethPriceFeed *aggregatorv3.AggregatorV3Interface,
	tempEthClient *ethclient.Client,
) (sidecartypes.OracleState, error) {
	ctx := context.Background()
	var wg sync.WaitGroup

	slog.Info("Retrieving latest header", "network", sidecartypes.NetworkNames[o.Config.Network], "time", time.Now().Format("15:04:05.00"))
	latestHeader, err := o.EthClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return sidecartypes.OracleState{}, fmt.Errorf("failed to fetch latest block: %w", err)
	}
	slog.Info("Retrieved latest header", "network", sidecartypes.NetworkNames[o.Config.Network], "block", latestHeader.Number.Uint64(), "time", time.Now().Format("15:04:05.00"))
	targetBlockNumber := new(big.Int).Sub(latestHeader.Number, big.NewInt(sidecartypes.EthBlocksBeforeFinality))

	// Check base fee availability
	if latestHeader.BaseFee == nil {
		return sidecartypes.OracleState{}, fmt.Errorf("base fee not available (pre-London fork?)")
	}

	update := o.initializeStateUpdate()
	var updateMutex sync.Mutex
	errChan := make(chan error, 16)

	// Fetch Ethereum contract data (AVS delegations and redemptions on EigenLayer)
	o.fetchEthereumContractData(&wg, serviceManager, zenBTCControllerHolesky, targetBlockNumber, update, &updateMutex, errChan)

	// Fetch network data (gas estimates, tips, Solana fees)
	o.fetchNetworkData(&wg, ctx, update, &updateMutex, errChan)

	// Fetch price data (ROCK, BTC, ETH)
	o.fetchPriceData(&wg, btcPriceFeed, ethPriceFeed, tempEthClient, ctx, update, &updateMutex, errChan)

	// Fetch zenBTC burn events from Ethereum
	o.fetchEthereumBurnEvents(&wg, latestHeader, update, &updateMutex, errChan)

	// Fetch Solana mint events for zenBTC
	o.processSolanaMintEvents(&wg, update, &updateMutex, errChan)

	// Fetch Solana burn events for zenBTC and ROCK
	o.fetchSolanaBurnEvents(&wg, update, &updateMutex, errChan)

	// Fetch and populate backfill requests from zrChain
	o.processBackfillRequests(&wg, update, &updateMutex)

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			slog.Error("Error during state fetch", "error", err)
			return sidecartypes.OracleState{}, err
		}
	}

	// Update signature strings and build final state
	return o.buildFinalState(update, latestHeader, targetBlockNumber)
}

func (o *Oracle) fetchEthereumContractData(
	wg *sync.WaitGroup,
	serviceManager *middleware.ContractZrServiceManager,
	zenBTCControllerHolesky *zenbtc.ZenBTController,
	targetBlockNumber *big.Int,
	update *oracleStateUpdate,
	updateMutex *sync.Mutex,
	errChan chan<- error,
) {
	// Fetches the state of AVS delegations from the service manager contract.
	wg.Add(1)
	go func() {
		defer wg.Done()
		delegations, err := o.getServiceManagerState(serviceManager, targetBlockNumber)
		if err != nil {
			errChan <- fmt.Errorf("failed to get contract state: %w", err)
			return
		}
		updateMutex.Lock()
		update.eigenDelegations = delegations
		updateMutex.Unlock()
	}()

	// Fetches pending zenBTC redemptions from the zenBTC controller contract.
	wg.Add(1)
	go func() {
		defer wg.Done()
		redemptions, err := o.getRedemptions(zenBTCControllerHolesky, targetBlockNumber)
		if err != nil {
			errChan <- fmt.Errorf("failed to get zenBTC contract state: %w", err)
			return
		}
		updateMutex.Lock()
		update.redemptions = redemptions
		updateMutex.Unlock()
	}()
}

func (o *Oracle) fetchNetworkData(
	wg *sync.WaitGroup,
	ctx context.Context,
	update *oracleStateUpdate,
	updateMutex *sync.Mutex,
	errChan chan<- error,
) {
	// Fetches the suggested gas tip cap (priority fee) for Ethereum transactions.
	wg.Add(1)
	go func() {
		defer wg.Done()
		suggestedTip, err := o.EthClient.SuggestGasTipCap(ctx)
		if err != nil {
			errChan <- fmt.Errorf("failed to get suggested priority fee: %w", err)
			return
		}
		updateMutex.Lock()
		update.suggestedTip = suggestedTip
		updateMutex.Unlock()
	}()

	// Fetches the current fee in lamports required per signature on Solana.
	wg.Add(1)
	go func() {
		defer wg.Done()
		lamportsPerSignature, err := o.getSolanaLamportsPerSignature(ctx)
		if err != nil {
			log.Printf("Warning: getSolanaLamportsPerSignature failed: %v. Using potentially stale/default value.", err)
		}
		updateMutex.Lock()
		update.solanaLamportsPerSignature = lamportsPerSignature
		updateMutex.Unlock()
	}()

	// Estimates the gas required for a zenBTC stake call on Ethereum.
	wg.Add(1)
	go func() {
		defer wg.Done()
		stakeCallData, err := validationkeeper.EncodeStakeCallData(big.NewInt(1000000000))
		if err != nil {
			errChan <- fmt.Errorf("failed to encode stake call data: %w", err)
			return
		}
		addr := common.HexToAddress(sidecartypes.ZenBTCControllerAddresses[o.Config.Network])
		estimatedGas, err := o.EthClient.EstimateGas(context.Background(), ethereum.CallMsg{
			From: common.HexToAddress(sidecartypes.WhitelistedRoleAddresses[o.Config.Network]),
			To:   &addr,
			Data: stakeCallData,
		})
		if err != nil {
			errChan <- fmt.Errorf("failed to estimate gas for stake call: %w", err)
			return
		}
		updateMutex.Lock()
		update.estimatedGas = (estimatedGas * sidecartypes.GasEstimationBuffer) / 100
		updateMutex.Unlock()
	}()
}

func (o *Oracle) fetchPriceData(
	wg *sync.WaitGroup,
	btcPriceFeed *aggregatorv3.AggregatorV3Interface,
	ethPriceFeed *aggregatorv3.AggregatorV3Interface,
	tempEthClient *ethclient.Client,
	ctx context.Context,
	update *oracleStateUpdate,
	updateMutex *sync.Mutex,
	errChan chan<- error,
) {
	httpTimeout := sidecartypes.DefaultHTTPTimeout

	// Fetches the latest ROCK/USD price from the specified public endpoint.
	wg.Add(1)
	go func() {
		defer wg.Done()
		client := &http.Client{
			Timeout: httpTimeout,
		}
		resp, err := client.Get(sidecartypes.ROCKUSDPriceURL)
		if err != nil {
			errChan <- fmt.Errorf("failed to retrieve ROCK price data: %w", err)
			return
		}
		defer resp.Body.Close()

		var priceData []PriceData
		if err := json.NewDecoder(resp.Body).Decode(&priceData); err != nil || len(priceData) == 0 {
			errChan <- fmt.Errorf("failed to decode ROCK price data or empty data: %w", err)
			return
		}
		priceDec, err := math.LegacyNewDecFromStr(priceData[0].Last)
		if err != nil {
			errChan <- fmt.Errorf("failed to parse ROCK price data: %w", err)
			return
		}
		updateMutex.Lock()
		update.ROCKUSDPrice = priceDec
		updateMutex.Unlock()
	}()

	// Fetches the latest BTC/USD and ETH/USD prices from Chainlink price feeds on Ethereum mainnet.
	wg.Add(1)
	go func() {
		defer wg.Done()
		mainnetLatestHeader, err := tempEthClient.HeaderByNumber(ctx, nil)
		if err != nil {
			errChan <- fmt.Errorf("failed to fetch latest mainnet block: %w", err)
			return
		}
		targetBlockNumberMainnet := new(big.Int).Sub(mainnetLatestHeader.Number, big.NewInt(sidecartypes.EthBlocksBeforeFinality))

		if btcPriceFeed == nil || ethPriceFeed == nil {
			errChan <- fmt.Errorf("BTC or ETH price feed not initialized")
			return
		}

		btcPrice, err := o.fetchPrice(btcPriceFeed, targetBlockNumberMainnet)
		if err != nil {
			errChan <- fmt.Errorf("failed to fetch BTC price: %w", err)
			return
		}

		ethPrice, err := o.fetchPrice(ethPriceFeed, targetBlockNumberMainnet)
		if err != nil {
			errChan <- fmt.Errorf("failed to fetch ETH price: %w", err)
			return
		}

		updateMutex.Lock()
		update.BTCUSDPrice = btcPrice
		update.ETHUSDPrice = ethPrice
		updateMutex.Unlock()
	}()
}

func (o *Oracle) fetchEthereumBurnEvents(
	wg *sync.WaitGroup,
	latestHeader *ethtypes.Header,
	update *oracleStateUpdate,
	updateMutex *sync.Mutex,
	errChan chan<- error,
) {
	// Fetches and processes recent zenBTC burn events from Ethereum within a defined block range.
	wg.Add(1)
	go func() {
		defer wg.Done()
		events, err := o.processEthereumBurnEvents(latestHeader)
		if err != nil {
			errChan <- fmt.Errorf("failed to process Ethereum burn events: %w", err)
			return
		}
		updateMutex.Lock()
		update.ethBurnEvents = events
		updateMutex.Unlock()
	}()
}

func (o *Oracle) processSolanaMintEvents(
	wg *sync.WaitGroup,
	update *oracleStateUpdate,
	updateMutex *sync.Mutex,
	errChan chan<- error,
) {
	// Fetches new ROCK and zenBTC mint events from Solana since the last processed signature,
	// and merges them with the existing cached events.
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Get new events using watermarking
		lastKnownRockSig := o.GetLastProcessedSolSignature(sidecartypes.SolRockMint)
		rockEvents, newRockSig, err := o.getSolROCKMints(sidecartypes.SolRockProgramID[o.Config.Network], lastKnownRockSig)
		if err != nil {
			errChan <- fmt.Errorf("failed to process Solana ROCK mint events: %w", err)
			return
		}

		lastKnownZenBTCSig := o.GetLastProcessedSolSignature(sidecartypes.SolZenBTCMint)
		zenbtcEvents, newZenBTCSig, err := o.getSolZenBTCMints(sidecartypes.ZenBTCSolanaProgramID[o.Config.Network], lastKnownZenBTCSig)
		if err != nil {
			errChan <- fmt.Errorf("failed to process Solana zenBTC mint events: %w", err)
			return
		}

		allNewEvents := append(rockEvents, zenbtcEvents...)

		// Get current state to merge with new mint events
		currentState := o.currentState.Load().(*sidecartypes.OracleState)

		// Use helper function to merge events
		mergedMintEvents := mergeNewMintEvents(currentState.SolanaMintEvents, currentState.CleanedSolanaMintEvents, allNewEvents, "Solana mint")

		updateMutex.Lock()
		update.SolanaMintEvents = mergedMintEvents
		if !newRockSig.IsZero() {
			update.latestSolanaSigs[sidecartypes.SolRockMint] = newRockSig
		}
		if !newZenBTCSig.IsZero() {
			update.latestSolanaSigs[sidecartypes.SolZenBTCMint] = newZenBTCSig
		}
		updateMutex.Unlock()
	}()
}

func (o *Oracle) fetchSolanaBurnEvents(
	wg *sync.WaitGroup,
	update *oracleStateUpdate,
	updateMutex *sync.Mutex,
	errChan chan<- error,
) {
	var zenBtcEvents, rockEvents []api.BurnEvent
	var zenBtcErr, rockErr error
	var wgEvents sync.WaitGroup

	// Fetches new zenBTC burn events from Solana since the last processed signature.
	wgEvents.Add(1)
	go func() {
		defer wgEvents.Done()
		lastKnownSig := o.GetLastProcessedSolSignature(sidecartypes.SolZenBTCBurn)
		var newestSig solana.Signature
		zenBtcEvents, newestSig, zenBtcErr = o.getSolanaZenBTCBurnEventsFn(sidecartypes.ZenBTCSolanaProgramID[o.Config.Network], lastKnownSig)
		if zenBtcErr == nil && !newestSig.IsZero() {
			updateMutex.Lock()
			update.latestSolanaSigs[sidecartypes.SolZenBTCBurn] = newestSig
			updateMutex.Unlock()
		}
	}()

	// Fetches new ROCK burn events from Solana since the last processed signature.
	wgEvents.Add(1)
	go func() {
		defer wgEvents.Done()
		lastKnownSig := o.GetLastProcessedSolSignature(sidecartypes.SolRockBurn)
		var newestSig solana.Signature
		rockEvents, newestSig, rockErr = o.getSolanaRockBurnEventsFn(sidecartypes.SolRockProgramID[o.Config.Network], lastKnownSig)
		if rockErr == nil && !newestSig.IsZero() {
			updateMutex.Lock()
			update.latestSolanaSigs[sidecartypes.SolRockBurn] = newestSig
			updateMutex.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		wgEvents.Wait() // Wait for both Solana burn fetches to complete

		if zenBtcErr != nil {
			errChan <- fmt.Errorf("failed to process Solana zenBTC burn events: %w", zenBtcErr)
		}
		if rockErr != nil {
			errChan <- fmt.Errorf("failed to process Solana ROCK burn events: %w", rockErr)
		}

		// If either failed, we don't proceed with merging to avoid partial state.
		if zenBtcErr != nil || rockErr != nil {
			return
		}

		// Merge and sort
		allNewSolanaBurnEvents := append(zenBtcEvents, rockEvents...)
		sort.Slice(allNewSolanaBurnEvents, func(i, j int) bool {
			if allNewSolanaBurnEvents[i].Height != allNewSolanaBurnEvents[j].Height {
				return allNewSolanaBurnEvents[i].Height < allNewSolanaBurnEvents[j].Height
			}
			return allNewSolanaBurnEvents[i].LogIndex < allNewSolanaBurnEvents[j].LogIndex
		})

		// Get current state to merge with new burn events
		currentState := o.currentState.Load().(*sidecartypes.OracleState)

		// Use helper function to merge events
		mergedBurnEvents := mergeNewBurnEvents(currentState.SolanaBurnEvents, currentState.CleanedSolanaBurnEvents, allNewSolanaBurnEvents, "Solana")

		updateMutex.Lock()
		update.solanaBurnEvents = mergedBurnEvents
		updateMutex.Unlock()
	}()
}

func (o *Oracle) buildFinalState(
	update *oracleStateUpdate,
	latestHeader *ethtypes.Header,
	targetBlockNumber *big.Int,
) (sidecartypes.OracleState, error) {
	// Update the main Oracle's last signature strings
	if len(update.latestSolanaSigs) > 0 {
		if sig, ok := update.latestSolanaSigs[sidecartypes.SolRockMint]; ok && !sig.IsZero() {
			o.lastSolRockMintSigStr = sig.String()
		}
		if sig, ok := update.latestSolanaSigs[sidecartypes.SolZenBTCMint]; ok && !sig.IsZero() {
			o.lastSolZenBTCMintSigStr = sig.String()
		}
		if sig, ok := update.latestSolanaSigs[sidecartypes.SolZenBTCBurn]; ok && !sig.IsZero() {
			o.lastSolZenBTCBurnSigStr = sig.String()
		}
		if sig, ok := update.latestSolanaSigs[sidecartypes.SolRockBurn]; ok && !sig.IsZero() {
			o.lastSolRockBurnSigStr = sig.String()
		}
		slog.Info("Updated latest Solana signatures",
			"rockMint", o.lastSolRockMintSigStr,
			"zenBTCMint", o.lastSolZenBTCMintSigStr,
			"zenBTCBurn", o.lastSolZenBTCBurnSigStr,
			"rockBurn", o.lastSolRockBurnSigStr)
	}

	currentState := o.currentState.Load().(*sidecartypes.OracleState)

	// Apply fallbacks for nil values
	o.applyFallbacks(update, currentState)

	// EigenDelegations map keys are deterministically ordered by encoding/json, so no manual sorting needed.

	// Sort all event slices to ensure deterministic order
	sort.Slice(update.ethBurnEvents, func(i, j int) bool {
		if update.ethBurnEvents[i].Height != update.ethBurnEvents[j].Height {
			return update.ethBurnEvents[i].Height < update.ethBurnEvents[j].Height
		}
		return update.ethBurnEvents[i].LogIndex < update.ethBurnEvents[j].LogIndex
	})
	sort.Slice(update.solanaBurnEvents, func(i, j int) bool {
		if update.solanaBurnEvents[i].Height != update.solanaBurnEvents[j].Height {
			return update.solanaBurnEvents[i].Height < update.solanaBurnEvents[j].Height
		}
		return update.solanaBurnEvents[i].LogIndex < update.solanaBurnEvents[j].LogIndex
	})
	sort.Slice(update.redemptions, func(i, j int) bool {
		return update.redemptions[i].Id < update.redemptions[j].Id
	})
	sort.Slice(update.SolanaMintEvents, func(i, j int) bool {
		if update.SolanaMintEvents[i].Height != update.SolanaMintEvents[j].Height {
			return update.SolanaMintEvents[i].Height < update.SolanaMintEvents[j].Height
		}
		// Use TxSig as a secondary sort key for determinism if heights are identical
		return update.SolanaMintEvents[i].TxSig < update.SolanaMintEvents[j].TxSig
	})

	newState := sidecartypes.OracleState{
		EigenDelegations:           update.eigenDelegations,
		EthBlockHeight:             targetBlockNumber.Uint64(),
		EthGasLimit:                update.estimatedGas,
		EthBaseFee:                 latestHeader.BaseFee.Uint64(),
		EthTipCap:                  update.suggestedTip.Uint64(),
		SolanaLamportsPerSignature: update.solanaLamportsPerSignature,
		EthBurnEvents:              update.ethBurnEvents,
		CleanedEthBurnEvents:       currentState.CleanedEthBurnEvents,
		SolanaBurnEvents:           update.solanaBurnEvents,
		CleanedSolanaBurnEvents:    currentState.CleanedSolanaBurnEvents,
		Redemptions:                update.redemptions,
		SolanaMintEvents:           update.SolanaMintEvents,
		CleanedSolanaMintEvents:    currentState.CleanedSolanaMintEvents,
		ROCKUSDPrice:               update.ROCKUSDPrice,
		BTCUSDPrice:                update.BTCUSDPrice,
		ETHUSDPrice:                update.ETHUSDPrice,
		LastSolRockMintSig:         o.lastSolRockMintSigStr,
		LastSolZenBTCMintSig:       o.lastSolZenBTCMintSigStr,
		LastSolZenBTCBurnSig:       o.lastSolZenBTCBurnSigStr,
		LastSolRockBurnSig:         o.lastSolRockBurnSigStr,
	}

	if o.DebugMode {
		jsonData, err := json.MarshalIndent(newState, "", "  ")
		if err != nil {
			slog.Error("Error marshalling state to JSON for logging", "error", err)
			slog.Info("State fetched (pre-update send - fallback)", "state", newState)
		} else {
			slog.Info("State fetched (pre-update send)", "state", string(jsonData))
		}
	}

	return newState, nil
}

func (o *Oracle) applyFallbacks(update *oracleStateUpdate, currentState *sidecartypes.OracleState) {
	// Ensure update fields that might not have been populated are not nil
	if update.suggestedTip == nil {
		update.suggestedTip = big.NewInt(0)
		slog.Warn("suggestedTip was nil, using 0")
	}
	if update.ROCKUSDPrice.IsNil() {
		update.ROCKUSDPrice = currentState.ROCKUSDPrice
		slog.Warn("ROCKUSDPrice was nil, using last known state value")
	}
	if update.BTCUSDPrice.IsNil() {
		update.BTCUSDPrice = currentState.BTCUSDPrice
		slog.Warn("BTCUSDPrice was nil, using last known state value")
	}
	if update.ETHUSDPrice.IsNil() {
		update.ETHUSDPrice = currentState.ETHUSDPrice
		slog.Warn("ETHUSDPrice was nil, using last known state value")
	}
	if update.solanaLamportsPerSignature == 0 {
		update.solanaLamportsPerSignature = currentState.SolanaLamportsPerSignature
		slog.Warn("solanaLamportsPerSignature was 0, using last known state value")
	}
}

func (o *Oracle) getServiceManagerState(contractInstance *middleware.ContractZrServiceManager, height *big.Int) (map[string]map[string]*big.Int, error) {
	delegations := make(map[string]map[string]*big.Int)

	callOpts := &bind.CallOpts{
		BlockNumber: height,
	}

	// Retrieve all validators from the contract
	allValidators, err := contractInstance.GetAllValidator(callOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to get all validators: %w", err)
	}

	quorumNumber := sidecartypes.EigenLayerQuorumNumber

	// Iterate over all validators
	for _, validator := range allValidators {
		validatorAddr := validator.ValidatorAddr
		operators := validator.Operators

		// Initialize the map for this validator if not already
		if delegations[validatorAddr] == nil {
			delegations[validatorAddr] = make(map[string]*big.Int)
		}

		// Iterate over operators associated with the validator
		for _, operator := range operators {
			// Get the stake amount for the operator
			amount, err := contractInstance.GetEigenStake(callOpts, operator, quorumNumber)
			if err != nil {
				slog.Error("Failed to get stake for operator", "operator", operator.Hex(), "error", err)
				continue
			}

			// Only consider positive stake amounts
			if amount.Cmp(big.NewInt(0)) > 0 {
				operatorAddr := operator.Hex()
				// Sum up the stake if operator already exists under this validator
				if existingAmount, exists := delegations[validatorAddr][operatorAddr]; exists {
					delegations[validatorAddr][operatorAddr] = new(big.Int).Add(existingAmount, amount)
				} else {
					delegations[validatorAddr][operatorAddr] = amount
				}
			}
		}
	}

	return delegations, nil
}

func (o *Oracle) processEthereumBurnEvents(latestHeader *ethtypes.Header) ([]api.BurnEvent, error) {
	fromBlock := new(big.Int).Sub(latestHeader.Number, big.NewInt(int64(sidecartypes.EthBurnEventsBlockRange)))
	toBlock := latestHeader.Number
	newEthBurnEvents, err := o.getEthBurnEvents(fromBlock, toBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to get Ethereum burn events: %w", err)
	}

	// Get current state to merge with new burn events
	currentState := o.currentState.Load().(*sidecartypes.OracleState)

	// Use helper function to merge events
	mergedEthBurnEvents := mergeNewBurnEvents(currentState.EthBurnEvents, currentState.CleanedEthBurnEvents, newEthBurnEvents, "Ethereum")

	return mergedEthBurnEvents, nil
}

// reconcileBurnEventsWithZRChain checks a list of burn events against the chain and returns the events
// that should remain in the cache and an updated map of cleaned events.
func (o *Oracle) reconcileBurnEventsWithZRChain(
	ctx context.Context,
	eventsToClean []api.BurnEvent,
	cleanedEvents map[string]bool,
	chainTypeName string, // For logging purposes (e.g., "Ethereum", "Solana")
) ([]api.BurnEvent, map[string]bool) { // Removed error return for simplicity now

	remainingEvents := make([]api.BurnEvent, 0)
	updatedCleanedEvents := make(map[string]bool)
	// Copy existing cleaned events map to avoid modifying the original directly
	maps.Copy(updatedCleanedEvents, cleanedEvents)

	for _, event := range eventsToClean {
		// Check if this specific event was already cleaned in a previous run but is still in the eventsToClean list for some reason
		key := fmt.Sprintf("%s-%s-%d", event.ChainID, event.TxID, event.LogIndex)
		if _, alreadyCleaned := updatedCleanedEvents[key]; alreadyCleaned {
			slog.Info("Skipping already cleaned burn event", "chain", chainTypeName, "txID", event.TxID, "logIndex", event.LogIndex, "chainID", event.ChainID)
			continue // Skip to next event if already marked as cleaned
		}

		var foundOnChain bool

		// 1. Check ZenBTC keeper
		zenbtcResp, err := o.zrChainQueryClient.ZenBTCQueryClient.BurnEvents(ctx, 0, event.TxID, event.LogIndex, event.ChainID)
		if err != nil {
			slog.Error("Error querying ZenBTC for burn event", "chain", chainTypeName, "txID", event.TxID, "logIndex", event.LogIndex, "chainID", event.ChainID, "error", err)
			// Keep events that we failed to query, they might succeed next time. We'll let it continue to the ZenTP check.
		}

		if zenbtcResp != nil && len(zenbtcResp.BurnEvents) > 0 {
			foundOnChain = true
		}

		// 2. If not found and it's a Solana event, check ZenTP keeper as well
		if !foundOnChain && chainTypeName == "Solana" {
			// The destination address for Solana burns is a 32-byte key, but zrchain uses the first 20 bytes for the Cosmos address.
			if len(event.DestinationAddr) >= 20 {
				bech32Addr, err := sdkBech32.ConvertAndEncode("zen", event.DestinationAddr[:20])
				if err != nil {
					slog.Error("Error converting destination address to bech32 for ZenTP query", "txID", event.TxID, "error", err)
				} else {
					zentpResp, err := o.zrChainQueryClient.ZenTPQueryClient.Burns(ctx, bech32Addr, event.TxID)
					if err != nil {
						slog.Error("Error querying ZenTP for Solana burn event", "txID", event.TxID, "address", bech32Addr, "error", err)
					}
					// Check zentpResp and its content.
					if zentpResp != nil && len(zentpResp.Burns) > 0 {
						foundOnChain = true
					}
				}
			} else {
				slog.Warn("Skipping ZenTP check for Solana burn event due to short destination address", "txID", event.TxID, "addressLength", len(event.DestinationAddr))
			}
		}

		// 3. Update state based on whether it was found
		if !foundOnChain {
			remainingEvents = append(remainingEvents, event)
		} else {
			// Event found on chain, mark it as cleaned by adding to the map
			updatedCleanedEvents[key] = true
			slog.Info("Removing burn event from cache as it's now on chain", "chain", chainTypeName, "txID", event.TxID, "logIndex", event.LogIndex, "chainID", event.ChainID)
			// Do *not* add to remainingEvents
		}
	}

	return remainingEvents, updatedCleanedEvents
}

func (o *Oracle) cleanUpBurnEvents() {
	o.cleanupMutex.Lock()
	defer o.cleanupMutex.Unlock()

	currentState := o.currentState.Load().(*sidecartypes.OracleState)

	// Check if there are any events to clean up at all
	initialEthCount := len(currentState.EthBurnEvents)
	initialSolCount := len(currentState.SolanaBurnEvents)
	if initialEthCount == 0 && initialSolCount == 0 {
		return // Nothing to clean
	}

	ctx := context.Background()
	stateChanged := false

	// Clean up Ethereum events
	remainingEthEvents, updatedCleanedEthEvents := o.reconcileBurnEventsWithZRChain(ctx, currentState.EthBurnEvents, currentState.CleanedEthBurnEvents, "Ethereum")
	if len(remainingEthEvents) != initialEthCount {
		slog.Info("Removed Ethereum burn events from cache", "count", initialEthCount-len(remainingEthEvents))
		stateChanged = true
	}

	// Clean up Solana events
	remainingSolEvents, updatedCleanedSolEvents := o.reconcileBurnEventsWithZRChain(ctx, currentState.SolanaBurnEvents, currentState.CleanedSolanaBurnEvents, "Solana")
	if len(remainingSolEvents) != initialSolCount {
		slog.Info("Removed Solana burn events from cache", "count", initialSolCount-len(remainingSolEvents))
		stateChanged = true
	}

	// Update the current state only if changes were made to either list
	if stateChanged {
		newState := *currentState // Copy existing state
		newState.EthBurnEvents = remainingEthEvents
		newState.CleanedEthBurnEvents = updatedCleanedEthEvents
		newState.SolanaBurnEvents = remainingSolEvents
		newState.CleanedSolanaBurnEvents = updatedCleanedSolEvents
		newState.SolanaMintEvents = currentState.SolanaMintEvents
		newState.CleanedSolanaMintEvents = currentState.CleanedSolanaMintEvents

		o.currentState.Store(&newState)
		o.CacheState() // Persist the updated state
		slog.Info("Burn event cache state updated and saved.")
	} else {
		slog.Info("No burn events removed from cache during cleanup.")
	}
}

// getEthBurnEvents retrieves all ZenBTCTokenRedemption (burn) events from the specified block range,
// converts them into []api.BurnEvent with correctly populated fields, and formats the chainID in CAIP-2 format.
func (o *Oracle) getEthBurnEvents(fromBlock, toBlock *big.Int) ([]api.BurnEvent, error) {
	ctx := context.Background()
	tokenAddress := common.HexToAddress(sidecartypes.ZenBTCTokenAddresses.Ethereum[o.Config.Network])

	// Create a new instance of the ZenBTC token contract
	zenBTCInstance, err := zenbtc.NewZenBTCFilterer(tokenAddress, o.EthClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZenBTC token contract filterer: %w", err)
	}

	// Set up the filter options
	endBlock := toBlock.Uint64()
	filterOpts := &bind.FilterOpts{
		Start:   fromBlock.Uint64(),
		End:     &endBlock,
		Context: ctx,
	}

	// Use the generated FilterTokenRedemption method
	iterator, err := zenBTCInstance.FilterTokenRedemption(filterOpts, nil) // nil means no filter on redeemer address
	if err != nil {
		return nil, fmt.Errorf("failed to filter token redemption events: %w", err)
	}
	defer iterator.Close()

	chainID, err := o.EthClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	var burnEvents []api.BurnEvent
	for iterator.Next() {
		event := iterator.Event
		if event == nil {
			continue
		}

		// Use block number as deterministic ordering key
		height := uint64(event.Raw.BlockNumber)

		burnEvents = append(burnEvents, api.BurnEvent{
			TxID:            event.Raw.TxHash.Hex(),
			LogIndex:        uint64(event.Raw.Index),
			ChainID:         fmt.Sprintf("eip155:%s", chainID.String()),
			DestinationAddr: event.DestAddr,
			Amount:          event.Value,
			IsZenBTC:        true,
			Height:          height,
		})
	}

	if err := iterator.Error(); err != nil {
		return nil, fmt.Errorf("error iterating through token redemption events: %w", err)
	}

	return burnEvents, nil
}

func (o *Oracle) getRedemptions(contractInstance *zenbtc.ZenBTController, height *big.Int) ([]api.Redemption, error) {
	callOpts := &bind.CallOpts{
		BlockNumber: height,
	}

	redemptionData, err := contractInstance.GetReadyForComplete(callOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent redemptions: %w", err)
	}

	redemptions := make([]api.Redemption, 0)
	for _, redemption := range redemptionData {
		redemptions = append(redemptions, api.Redemption{
			Id:                 redemption.Nonce.Uint64(),
			DestinationAddress: redemption.DestinationAddress,
			Amount:             redemption.ZenBTCValue.Uint64(),
		})
	}

	return redemptions, nil
}

func (o *Oracle) reconcileMintEventsWithZRChain(
	ctx context.Context,
	eventsToClean []api.SolanaMintEvent,
	cleanedEvents map[string]bool,
) ([]api.SolanaMintEvent, map[string]bool) {
	remainingEvents := make([]api.SolanaMintEvent, 0)
	updatedCleanedEvents := make(map[string]bool)
	maps.Copy(updatedCleanedEvents, cleanedEvents)

	for _, event := range eventsToClean {
		key := base64.StdEncoding.EncodeToString(event.SigHash)
		if _, alreadyCleaned := updatedCleanedEvents[key]; alreadyCleaned {
			continue
		}

		var foundOnChain bool

		// Check ZenBTC keeper
		zenbtcResp, err := o.zrChainQueryClient.ZenBTCQueryClient.PendingMintTransaction(ctx, event.TxSig)
		if err != nil {
			slog.Debug("Error querying ZenBTC for mint event", "txSig", event.TxSig, "error", err)
		}

		if zenbtcResp != nil && zenbtcResp.PendingMintTransaction != nil &&
			zenbtcResp.PendingMintTransaction.Status == zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED {
			foundOnChain = true
		}

		// If not found, check ZenTP keeper as well
		if !foundOnChain {
			zentpResp, err := o.zrChainQueryClient.ZenTPQueryClient.Mints(ctx, "", event.TxSig, zentptypes.BridgeStatus_BRIDGE_STATUS_COMPLETED)
			if err != nil {
				slog.Debug("Error querying ZenTP for mint event", "txSig", event.TxSig, "error", err)
			}
			if zentpResp != nil && len(zentpResp.Mints) > 0 {
				foundOnChain = true
			}
		}

		if !foundOnChain {
			remainingEvents = append(remainingEvents, event)
		} else {
			updatedCleanedEvents[key] = true
			slog.Info("Removing Solana mint event from cache as it's now on chain", "txSig", event.TxSig, "sigHash", key)
		}
	}

	return remainingEvents, updatedCleanedEvents
}

func (o *Oracle) cleanUpMintEvents() {
	o.cleanupMutex.Lock()
	defer o.cleanupMutex.Unlock()

	currentState := o.currentState.Load().(*sidecartypes.OracleState)

	initialMintCount := len(currentState.SolanaMintEvents)
	if initialMintCount == 0 {
		return
	}

	ctx := context.Background()
	stateChanged := false

	remainingMintEvents, updatedCleanedMintEvents := o.reconcileMintEventsWithZRChain(ctx, currentState.SolanaMintEvents, currentState.CleanedSolanaMintEvents)
	if len(remainingMintEvents) != initialMintCount {
		slog.Info("Removed Solana mint events from cache", "count", initialMintCount-len(remainingMintEvents))
		stateChanged = true
	}

	if stateChanged {
		newState := *currentState
		newState.SolanaMintEvents = remainingMintEvents
		newState.CleanedSolanaMintEvents = updatedCleanedMintEvents
		o.currentState.Store(&newState)
		o.CacheState()
		slog.Info("Mint event cache state updated and saved.")
	} else {
		slog.Info("No mint events removed from cache during cleanup.")
	}
}

// processTransactionFunc defines the function signature for processing a single Solana transaction.
// It returns a slice of events (as any), and an error if processing fails.
type processTransactionFunc func(
	txResult *solrpc.GetTransactionResult,
	program solana.PublicKey,
	sig solana.Signature,
	debugMode bool,
) ([]any, error)

// getSolanaEvents is a generic function to fetch and process events from a Solana program.
// It handles the boilerplate of fetching signatures, batching transaction details, and managing
// a watermark to only process new transactions. It delegates the specific logic for parsing
// a transaction to the provided processTransaction function.
func (o *Oracle) getSolanaEvents(
	programIDStr string,
	lastKnownSig solana.Signature,
	eventTypeName string,
	// processTransaction is a callback function that contains the logic to parse a single, fetched
	// Solana transaction and extract the relevant events (e.g., mints or burns).
	processTransaction processTransactionFunc,
) ([]any, solana.Signature, error) {
	limit := sidecartypes.SolanaEventScanTxLimit
	program, err := solana.PublicKeyFromBase58(programIDStr)
	if err != nil {
		return nil, lastKnownSig, fmt.Errorf("failed to obtain program public key for %s: %w", eventTypeName, err)
	}

	// Fetch latest signatures for the program address
	allSignatures, err := o.getSignaturesForAddressFn(context.Background(), program, &solrpc.GetSignaturesForAddressOpts{
		Limit:      &limit,
		Commitment: solrpc.CommitmentConfirmed,
	})
	if err != nil {
		return nil, lastKnownSig, fmt.Errorf("failed to get %s signatures: %w", eventTypeName, err)
	}

	if len(allSignatures) == 0 {
		slog.Info("Retrieved 0 events (no signatures found)", "eventType", eventTypeName)
		return []any{}, lastKnownSig, nil
	}

	// The newest signature from the node's perspective for this program address.
	newestSigFromNode := allSignatures[0].Signature
	newSignaturesToFetchDetails := make([]*solrpc.TransactionSignature, 0)

	// Filter and prepare signatures for processing
	newSignaturesToFetchDetails, _ = o.filterNewSignatures(allSignatures, lastKnownSig, eventTypeName, newestSigFromNode, limit)
	if len(newSignaturesToFetchDetails) == 0 {
		return []any{}, newestSigFromNode, nil
	}

	var processedEvents []any
	lastSuccessfullyProcessedSig := lastKnownSig
	internalBatchSize := sidecartypes.SolanaEventFetchBatchSize
	maxTxVersion := sidecartypes.SolanaTransactionVersion0

	for i := 0; i < len(newSignaturesToFetchDetails); i += internalBatchSize {
		end := min(i+internalBatchSize, len(newSignaturesToFetchDetails))
		currentBatchSignatures := newSignaturesToFetchDetails[i:end]

		batchRequests := make(jsonrpc.RPCRequests, 0, len(currentBatchSignatures))
		for j, sigInfo := range currentBatchSignatures {
			batchRequests = append(batchRequests, &jsonrpc.RPCRequest{
				Method: "getTransaction",
				Params: []any{
					sigInfo.Signature.String(),
					map[string]any{
						"encoding":                       solana.EncodingBase64,
						"commitment":                     solrpc.CommitmentConfirmed,
						"maxSupportedTransactionVersion": &maxTxVersion,
					},
				},
				ID:      uint64(j),
				JSONRPC: "2.0",
			})
		}

		var batchResponses jsonrpc.RPCResponses
		var err error
		// Execute the batch request with retries
		for retry := 0; retry < sidecartypes.SolanaEventFetchMaxRetries; retry++ {
			batchResponses, err = o.rpcCallBatchFn(context.Background(), batchRequests)
			if err == nil {
				// Quick check for any errors inside the response. If so, we'll retry the whole batch.
				hasErrors := false
				for _, resp := range batchResponses {
					if resp.Error != nil {
						hasErrors = true
						break
					}
				}
				if !hasErrors {
					break // Success, exit retry loop.
				}
				err = fmt.Errorf("response contains errors") // Set err to non-nil to trigger retry
			}

			// If we are here, it means there was an error.
			slog.Warn("Sub-batch GetTransaction failed after retries. Retrying...",
				"eventType", eventTypeName,
				"startIndex", i,
				"endIndex", end-1,
				"error", err,
				"retryCount", retry+1,
				"maxRetries", sidecartypes.SolanaEventFetchMaxRetries)

			// Don't sleep on the last attempt
			if retry < sidecartypes.SolanaEventFetchMaxRetries-1 {
				time.Sleep(sidecartypes.SolanaEventFetchRetrySleep)
			}
		}

		// If the batch request failed after all retries, fall back to individual requests
		if err != nil {
			slog.Warn("Batch request failed after all retries. Falling back to individual requests.", "eventType", eventTypeName)
			for _, sigInfo := range currentBatchSignatures {
				var txResult *solrpc.GetTransactionResult
				var err error

				// Retry individual fallback requests
				for retry := 0; retry < sidecartypes.SolanaFallbackMaxRetries; retry++ {
					txResult, err = o.getTransactionFn(context.Background(), sigInfo.Signature, &solrpc.GetTransactionOpts{
						Encoding:                       solana.EncodingBase64,
						Commitment:                     solrpc.CommitmentConfirmed,
						MaxSupportedTransactionVersion: &maxTxVersion,
					})
					if err == nil && txResult != nil {
						break // Success, exit retry loop
					}

					// Log retry attempts
					if retry < sidecartypes.SolanaFallbackMaxRetries-1 {
						slog.Warn("Fallback GetTransaction failed.",
							"tx", sigInfo.Signature,
							"eventType", eventTypeName,
							"error", err,
							"retryCount", retry+1,
							"maxRetries", sidecartypes.SolanaFallbackMaxRetries)
						time.Sleep(sidecartypes.SolanaEventFetchRetrySleep)
					}
				}

				if err != nil {
					slog.Error("Error in fallback GetTransaction after all retries", "tx", sigInfo.Signature, "eventType", eventTypeName, "error", err)
					continue
				}
				if txResult == nil {
					slog.Error("Nil result in fallback GetTransaction after all retries", "tx", sigInfo.Signature, "eventType", eventTypeName)
					continue
				}

				events, err := processTransaction(txResult, program, sigInfo.Signature, o.DebugMode)
				if err != nil {
					slog.Error("Failed to process events in fallback, skipping. Event is likely of an unrelated type.", "tx", sigInfo.Signature, "eventType", eventTypeName, "error", err)
					continue
				}

				if len(events) > 0 {
					processedEvents = append(processedEvents, events...)
				}
				lastSuccessfullyProcessedSig = sigInfo.Signature
				time.Sleep(sidecartypes.SolanaFallbackSleepInterval) // Rate limit individual fallback requests
			}
			continue // Skip the rest of the loop for this batch
		}

		if end < len(newSignaturesToFetchDetails) {
			time.Sleep(sidecartypes.SolanaSleepInterval)
		}

		for _, resp := range batchResponses {
			requestIndex, ok := parseRPCResponseID(resp, eventTypeName)
			if !ok {
				continue
			}
			if !validateRequestIndex(requestIndex, len(currentBatchSignatures), eventTypeName) {
				continue
			}
			sig := currentBatchSignatures[requestIndex].Signature

			if resp.Error != nil {
				// This should ideally not be hit if the retry logic above is working, but kept as a safeguard.
				slog.Error("Error in sub-batch GetTransaction result. This transaction will be missed in this cycle.", "tx", sig, "eventType", eventTypeName, "error", resp.Error)
				continue
			}
			if resp.Result == nil {
				slog.Error("Nil result field in sub-batch response", "tx", sig, "eventType", eventTypeName)
				continue
			}

			var txResult solrpc.GetTransactionResult
			if err := json.Unmarshal(resp.Result, &txResult); err != nil {
				slog.Error("Failed to unmarshal GetTransactionResult", "tx", sig, "eventType", eventTypeName, "error", err)
				continue
			}

			// Call the processor function to handle the token-specific logic.
			events, err := processTransaction(&txResult, program, sig, o.DebugMode)
			if err != nil {
				slog.Warn("Failed to process events, skipping. Event is likely of an unrelated type.", "tx", sig, "eventType", eventTypeName, "error", err)
				continue // Skip this transaction
			}

			if len(events) > 0 {
				processedEvents = append(processedEvents, events...)
			}

			// This signature has been processed successfully, so we can advance the watermark.
			lastSuccessfullyProcessedSig = sig
		}
	}

	slog.Info("From inspected transactions, retrieved new events. Newest last processed signature:",
		"count", len(processedEvents),
		"eventType", eventTypeName,
		"newestLastProcessedSig", lastSuccessfullyProcessedSig.String())
	return processedEvents, lastSuccessfullyProcessedSig, nil
}

// processMintTransaction is a generic helper that processes a single Solana transaction to extract mint events.
// It's designed to be reusable for different SPL tokens (like ROCK and zenBTC) by accepting functions
// that handle token-specific decoding and data extraction.
func (o *Oracle) processMintTransaction(
	txResult *solrpc.GetTransactionResult,
	program solana.PublicKey,
	sig solana.Signature,
	debugMode bool,
	// decodeEvents is a function that knows how to decode all events for a specific SPL token program.
	decodeEvents func(*solrpc.GetTransactionResult, solana.PublicKey) ([]any, error),
	// getEventData is a function that knows how to extract mint details from a specific "TokensMintedWithFee" event type for that SPL token.
	getEventData func(any) (recipient solana.PublicKey, value, fee uint64, mint solana.PublicKey, ok bool),
	eventTypeName string,
) ([]any, error) {
	decodedEvents, err := decodeEvents(txResult, program)
	if err != nil {
		return nil, err
	}

	// Extract transaction details for SigHash calculation
	if txResult.Transaction == nil {
		slog.Debug("Transaction envelope is nil in GetTransactionResult", "tx", sig, "type", eventTypeName)
		return nil, nil // Not an error, just no data
	}
	solTX, err := txResult.Transaction.GetTransaction()
	if err != nil || solTX == nil {
		return nil, fmt.Errorf("failed to get solana.Transaction from GetTransactionResult for sig %s: %w", sig, err)
	}

	if len(solTX.Signatures) != 2 {
		slog.Debug("Transaction does not have exactly 2 signatures; skipping SigHash calculation", "tx", sig.String(), "type", eventTypeName, "signatures", len(solTX.Signatures))
		return nil, nil // Not an error, just not the transaction type we are looking for.
	}
	combined := append(solTX.Signatures[0][:], solTX.Signatures[1][:]...)
	sigHash := sha256.Sum256(combined)

	var mintEvents []any
	for _, event := range decodedEvents {
		// Use reflection to access fields of the event, which could be of type
		// *rock_spl_token.Event or *zenbtc_spl_token.Event.
		eventValue := reflect.ValueOf(event)
		if eventValue.Kind() == reflect.Ptr {
			eventValue = eventValue.Elem()
		}

		if eventValue.Kind() != reflect.Struct {
			continue // Should not happen with current implementation
		}

		eventNameField := eventValue.FieldByName("Name")
		eventDataField := eventValue.FieldByName("Data")

		if !eventNameField.IsValid() || !eventDataField.IsValid() {
			continue // Should not happen if event structs are as expected
		}

		if eventNameField.String() == "TokensMintedWithFee" {
			recipient, value, fee, mint, ok := getEventData(eventDataField.Interface())
			if !ok {
				slog.Warn("Type assertion failed for TokensMintedWithFeeEventData", "eventType", eventTypeName, "tx", sig)
				continue
			}
			mintEvent := api.SolanaMintEvent{
				SigHash:   sigHash[:],
				Height:    uint64(txResult.Slot),
				Recipient: recipient.Bytes(),
				Value:     value,
				Fee:       fee,
				Mint:      mint.Bytes(),
				TxSig:     sig.String(),
			}
			mintEvents = append(mintEvents, mintEvent)
			if debugMode {
				slog.Debug("Solana Mint Event",
					"eventType", eventTypeName,
					"txSig", sig.String(),
					"sigHash", fmt.Sprintf("%x", mintEvent.SigHash),
					"recipient", solana.PublicKeyFromBytes(mintEvent.Recipient).String(),
					"height", mintEvent.Height,
					"value", mintEvent.Value,
					"fee", mintEvent.Fee,
					"mint", solana.PublicKeyFromBytes(mintEvent.Mint).String())
			}
		}
	}
	return mintEvents, nil
}

func (o *Oracle) getSolROCKMints(programID string, lastKnownSig solana.Signature) ([]api.SolanaMintEvent, solana.Signature, error) {
	eventTypeName := "Solana ROCK mint"
	// processor defines how to extract ROCK mint events from a single Solana transaction.
	// It's passed to the generic getSolanaEvents function to handle the specific logic for this token type.
	processor := func(txResult *solrpc.GetTransactionResult, program solana.PublicKey, sig solana.Signature, debugMode bool) ([]any, error) {
		return o.processMintTransaction(txResult, program, sig, debugMode,
			// Decode all events for the ROCK SPL token from the given transaction.
			func(tx *solrpc.GetTransactionResult, prog solana.PublicKey) ([]any, error) {
				events, err := rock_spl_token.DecodeEvents(tx, prog)
				if err != nil {
					return nil, err
				}
				var interfaceEvents []any
				for _, event := range events {
					interfaceEvents = append(interfaceEvents, event)
				}
				return interfaceEvents, nil
			},
			// Extract the relevant details (recipient, value, fee, mint) from a specific
			// TokensMintedWithFee event for the ROCK SPL token.
			func(data any) (solana.PublicKey, uint64, uint64, solana.PublicKey, bool) {
				eventData, ok := data.(*rock_spl_token.TokensMintedWithFeeEventData)
				if !ok {
					return solana.PublicKey{}, 0, 0, solana.PublicKey{}, false
				}
				return eventData.Recipient, eventData.Value, eventData.Fee, eventData.Mint, true
			},
			eventTypeName,
		)
	}

	untypedEvents, newWatermark, err := o.getSolanaEvents(programID, lastKnownSig, eventTypeName, processor)
	if err != nil {
		return nil, lastKnownSig, err
	}

	mintEvents := make([]api.SolanaMintEvent, len(untypedEvents))
	for i, untypedEvent := range untypedEvents {
		mintEvents[i] = untypedEvent.(api.SolanaMintEvent)
	}

	return mintEvents, newWatermark, nil
}

func (o *Oracle) getSolZenBTCMints(programID string, lastKnownSig solana.Signature) ([]api.SolanaMintEvent, solana.Signature, error) {
	eventTypeName := "Solana zenBTC mint"
	// processor defines how to extract zenBTC mint events from a single Solana transaction.
	// It's passed to the generic getSolanaEvents function to handle the specific logic for this token type.
	processor := func(txResult *solrpc.GetTransactionResult, program solana.PublicKey, sig solana.Signature, debugMode bool) ([]any, error) {
		return o.processMintTransaction(txResult, program, sig, debugMode,
			// Decode all events for the zenBTC SPL token from the given transaction.
			func(tx *solrpc.GetTransactionResult, prog solana.PublicKey) ([]any, error) {
				events, err := zenbtc_spl_token.DecodeEvents(tx, prog)
				if err != nil {
					return nil, err
				}
				var interfaceEvents []any
				for _, event := range events {
					interfaceEvents = append(interfaceEvents, event)
				}
				return interfaceEvents, nil
			},
			// Extract the relevant details (recipient, value, fee, mint) from a specific
			// TokensMintedWithFee event for the zenBTC SPL token.
			func(data any) (solana.PublicKey, uint64, uint64, solana.PublicKey, bool) {
				eventData, ok := data.(*zenbtc_spl_token.TokensMintedWithFeeEventData)
				if !ok {
					return solana.PublicKey{}, 0, 0, solana.PublicKey{}, false
				}
				return eventData.Recipient, eventData.Value, eventData.Fee, eventData.Mint, true
			},
			eventTypeName,
		)
	}

	untypedEvents, newWatermark, err := o.getSolanaEvents(programID, lastKnownSig, eventTypeName, processor)
	if err != nil {
		return nil, lastKnownSig, err
	}

	mintEvents := make([]api.SolanaMintEvent, len(untypedEvents))
	for i, untypedEvent := range untypedEvents {
		mintEvents[i] = untypedEvent.(api.SolanaMintEvent)
	}

	return mintEvents, newWatermark, nil
}

// getSolanaLamportsPerSignature fetches the current lamports per signature from the Solana network
// Uses the same slot rounding logic as getSolanaRecentBlockhash for consistency
func (o *Oracle) getSolanaLamportsPerSignature(ctx context.Context) (uint64, error) {
	// Create a simple dummy transaction to estimate fees.
	// Using placeholder public keys. These don't need to exist or have funds
	// as the transaction is not actually sent, only used for fee calculation.
	dummySigner := solana.MustPublicKeyFromBase58("1nc1nerator11111111111111111111111111111111")   // Use Sol Incinerator or other valid non-program ID
	dummyReceiver := solana.MustPublicKeyFromBase58("Stake11111111111111111111111111111111111111") // Use StakeProgramID as another valid key

	// Get a recent blockhash
	recentBlockhashResult, err := o.solanaClient.GetLatestBlockhash(ctx, solrpc.CommitmentConfirmed)
	if err != nil {
		slog.Warn("Failed to GetLatestBlockhash for fee calculation. Returning default lamports/sig.", "error", err, "defaultLamports", sidecartypes.DefaultSolanaFeeReturned)
		return sidecartypes.DefaultSolanaFeeReturned, fmt.Errorf("GetLatestBlockhash RPC call failed: %w", err)
	}
	if recentBlockhashResult == nil || recentBlockhashResult.Value == nil {
		slog.Warn("Incomplete GetLatestBlockhash result for fee calculation. Returning default lamports/sig.", "defaultLamports", sidecartypes.DefaultSolanaFeeReturned)
		return sidecartypes.DefaultSolanaFeeReturned, fmt.Errorf("GetLatestBlockhash returned nil result or value")
	}
	recentBlockhash := recentBlockhashResult.Value.Blockhash

	// Create a new transaction builder
	txBuilder := solana.NewTransactionBuilder()

	// Add a simple transfer instruction (e.g., transfer 1 lamport)
	// The actual details of the instruction don't matter as much as its presence and size.
	transferIx := solanagoSystem.NewTransferInstruction(
		1,           // 1 lamport
		dummySigner, // Use dummySigner as the source of the transfer
		dummyReceiver,
	).Build()
	txBuilder.AddInstruction(transferIx)
	txBuilder.SetFeePayer(dummySigner) // dummySigner is also the fee payer
	txBuilder.SetRecentBlockHash(recentBlockhash)

	// The message needs to be compiled and serialized.
	// For `getFeeForMessage`, we typically don't need to sign it.
	// First, build the transaction.
	tx, err := txBuilder.Build()
	if err != nil {
		slog.Warn("Failed to build transaction for fee calculation. Returning default lamports/sig.", "error", err, "defaultLamports", sidecartypes.DefaultSolanaFeeReturned)
		return sidecartypes.DefaultSolanaFeeReturned, fmt.Errorf("failed to build transaction for fee calculation: %w", err)
	}
	messageData := tx.Message // tx.Message is of type solana.Message (a struct)

	// Get the serialized message bytes using the standard MarshalBinary interface:
	serializedMessage, err := messageData.MarshalBinary()
	if err != nil {
		slog.Warn("Failed to serialize message using messageData.MarshalBinary for fee calculation. Returning default lamports/sig.", "error", err, "defaultLamports", sidecartypes.DefaultSolanaFeeReturned)
		return sidecartypes.DefaultSolanaFeeReturned, fmt.Errorf("failed to serialize message using messageData.MarshalBinary: %w", err)
	}

	// Call GetFeeForMessage (expects base64 encoded message string)
	msgBase64 := base64.StdEncoding.EncodeToString(serializedMessage)
	resp, err := o.solanaClient.GetFeeForMessage(ctx, msgBase64, solrpc.CommitmentConfirmed)
	if err != nil {
		slog.Warn("Failed to get Solana fees via GetFeeForMessage. Returning default lamports/sig.", "error", err, "defaultLamports", sidecartypes.DefaultSolanaFeeReturned)
		return sidecartypes.DefaultSolanaFeeReturned, fmt.Errorf("GetFeeForMessage RPC call failed: %w", err)
	}

	if resp == nil || resp.Value == nil {
		slog.Warn("Incomplete fee data from Solana RPC (GetFeeForMessage response or value is nil). Returning default lamports/sig.", "defaultLamports", sidecartypes.DefaultSolanaFeeReturned)
		return sidecartypes.DefaultSolanaFeeReturned, fmt.Errorf("GetFeeForMessage returned nil response or value")
	}

	// The fee is returned in lamports for the entire message.
	// To get lamports per signature, we'd typically divide by the number of signatures
	// in a standard transaction. For now, let's assume the fee returned is representative enough
	// or that it implicitly means per signature in this context for typical transactions.
	// Solana's documentation on `getFeeForMessage` states it returns the fee for the message in lamports.
	// If a transaction has 1 signature, this value is effectively lamports per signature.
	// If it could have more, this logic might need refinement or be based on a typical tx signature count.
	// For now, we'll use the direct value, assuming it's the intended lamports-per-signature proxy.
	lamports := *resp.Value
	if lamports == 0 {
		// It's possible for fees to be 0 on devnet/testnet or if priority fees are not needed.
		// However, consistently returning 0 might indicate an issue or a need for a non-zero default
		// if the oracle relies on a non-zero fee for some calculations.
		slog.Warn("Solana GetFeeForMessage returned 0 for LamportsPerSignature, using default if required by downstream logic, otherwise using 0.", "defaultLamports", sidecartypes.DefaultSolanaFeeReturned)
		// Depending on requirements, you might return 0 here, or stick to a default like 5000.
		// Let's return 0 if the network says 0, but log it.
		return 0, nil // Or return 5000, nil if a non-zero value is strictly necessary downstream
	}
	return lamports, nil
}

// processBurnTransaction is a generic helper that processes a single Solana transaction to extract burn events.
// It is reusable for different SPL tokens by accepting functions that handle token-specific logic.
func (o *Oracle) processBurnTransaction(
	txResult *solrpc.GetTransactionResult,
	program solana.PublicKey,
	sig solana.Signature,
	debugMode bool,
	// decodeEvents is a function that knows how to decode all events for a specific SPL token program.
	decodeEvents func(*solrpc.GetTransactionResult, solana.PublicKey) ([]any, error),
	// getEventData is a function that knows how to extract burn details from a specific "TokenRedemption" event type.
	getEventData func(any) (destAddr []byte, value uint64, ok bool),
	eventTypeName string,
	chainID string,
	isZenBTC bool,
) ([]any, error) {
	decodedEvents, err := decodeEvents(txResult, program)
	if err != nil {
		return nil, err // The caller will log this as a "likely unrelated type"
	}

	if debugMode {
		slog.Debug("Processing tx: found events", "eventType", eventTypeName, "tx", sig, "eventCount", len(decodedEvents))
		for i, event := range decodedEvents {
			// Use reflection to see event details for debugging
			eventValue := reflect.ValueOf(event)
			if eventValue.Kind() == reflect.Ptr {
				eventValue = eventValue.Elem()
			}

			if eventValue.Kind() != reflect.Struct {
				continue // Should not happen with current implementation
			}

			eventNameField := eventValue.FieldByName("Name")
			if eventNameField.IsValid() {
				slog.Debug("Event details", "eventIndex", i, "eventName", eventNameField.String(), "eventType", fmt.Sprintf("%T", event))
			}
		}
	}

	var burnEvents []any
	for logIndex, event := range decodedEvents {
		// Use reflection to access fields of the event, which could be of type
		// *rock_spl_token.Event or *zenbtc_spl_token.Event.
		eventValue := reflect.ValueOf(event)
		if eventValue.Kind() == reflect.Ptr {
			eventValue = eventValue.Elem()
		}

		if eventValue.Kind() != reflect.Struct {
			continue // Should not happen
		}

		eventNameField := eventValue.FieldByName("Name")
		eventDataField := eventValue.FieldByName("Data")

		if !eventNameField.IsValid() || !eventDataField.IsValid() {
			continue // Should not happen
		}

		if eventNameField.String() == "TokenRedemption" {
			destAddr, value, ok := getEventData(eventDataField.Interface())
			if !ok {
				slog.Warn("Type assertion failed for TokenRedemptionEventData", "eventType", eventTypeName, "tx", sig)
				continue
			}
			burnEvent := api.BurnEvent{
				TxID:            sig.String(),
				LogIndex:        uint64(logIndex),
				ChainID:         chainID,
				DestinationAddr: destAddr,
				Amount:          value,
				IsZenBTC:        isZenBTC,
				Height:          uint64(txResult.Slot),
			}
			burnEvents = append(burnEvents, burnEvent)
			if debugMode {
				slog.Debug("Burn Event",
					"eventType", eventTypeName,
					"txID", burnEvent.TxID,
					"logIndex", burnEvent.LogIndex,
					"chainID", burnEvent.ChainID,
					"destinationAddr", fmt.Sprintf("%x", burnEvent.DestinationAddr),
					"amount", burnEvent.Amount)
			}
		}
	}
	return burnEvents, nil
}

// getSolanaZenBTCBurnEvents retrieves ZenBTC burn events from Solana.
func (o *Oracle) getSolanaZenBTCBurnEvents(programID string, lastKnownSig solana.Signature) ([]api.BurnEvent, solana.Signature, error) {
	eventTypeName := "Solana zenBTC burn"
	chainID := sidecartypes.SolanaCAIP2[o.Config.Network]

	// processor defines how to extract zenBTC burn events from a single Solana transaction.
	// It's passed to the generic getSolanaEvents function to handle the specific logic for this token type.
	processor := func(txResult *solrpc.GetTransactionResult, program solana.PublicKey, sig solana.Signature, debugMode bool) ([]any, error) {
		return o.processBurnTransaction(txResult, program, sig, debugMode,
			// Decode all events for the zenBTC SPL token from the given transaction.
			func(tx *solrpc.GetTransactionResult, prog solana.PublicKey) ([]any, error) {
				events, err := zenbtc_spl_token.DecodeEvents(tx, prog)
				if err != nil {
					return nil, err
				}
				var interfaceEvents []any
				for _, event := range events {
					interfaceEvents = append(interfaceEvents, event)
				}
				return interfaceEvents, nil
			},
			// Extract the relevant details (destination address, value) from a specific
			// TokenRedemption event for the zenBTC SPL token.
			func(data any) (destAddr []byte, value uint64, ok bool) {
				eventData, ok := data.(*zenbtc_spl_token.TokenRedemptionEventData)
				if !ok {
					return nil, 0, false
				}
				return eventData.DestAddr, eventData.Value, true
			},
			eventTypeName, chainID, true,
		)
	}

	untypedEvents, newWatermark, err := o.getSolanaEvents(programID, lastKnownSig, eventTypeName, processor)
	if err != nil {
		return nil, lastKnownSig, err
	}

	burnEvents := make([]api.BurnEvent, len(untypedEvents))
	for i, untypedEvent := range untypedEvents {
		burnEvents[i] = untypedEvent.(api.BurnEvent)
	}

	return burnEvents, newWatermark, nil
}

// getSolanaRockBurnEvents retrieves Rock burn events from Solana.
func (o *Oracle) getSolanaRockBurnEvents(programID string, lastKnownSig solana.Signature) ([]api.BurnEvent, solana.Signature, error) {
	eventTypeName := "Solana ROCK burn"
	chainID := sidecartypes.SolanaCAIP2[o.Config.Network]

	// processor defines how to extract ROCK burn events from a single Solana transaction.
	// It's passed to the generic getSolanaEvents function to handle the specific logic for this token type.
	processor := func(txResult *solrpc.GetTransactionResult, program solana.PublicKey, sig solana.Signature, debugMode bool) ([]any, error) {
		return o.processBurnTransaction(txResult, program, sig, debugMode,
			// Decode all events for the ROCK SPL token from the given transaction.
			func(tx *solrpc.GetTransactionResult, prog solana.PublicKey) ([]any, error) {
				events, err := rock_spl_token.DecodeEvents(tx, prog)
				if err != nil {
					return nil, err
				}
				var interfaceEvents []any
				for _, event := range events {
					interfaceEvents = append(interfaceEvents, event)
				}
				return interfaceEvents, nil
			},
			// Extract the relevant details (destination address, value) from a specific
			// TokenRedemption event for the ROCK SPL token.
			func(data any) (destAddr []byte, value uint64, ok bool) {
				eventData, ok := data.(*rock_spl_token.TokenRedemptionEventData)
				if !ok {
					return nil, 0, false
				}
				return eventData.DestAddr[:], eventData.Value, true
			},
			eventTypeName, chainID, false,
		)
	}

	untypedEvents, newWatermark, err := o.getSolanaEvents(programID, lastKnownSig, eventTypeName, processor)
	if err != nil {
		return nil, lastKnownSig, err
	}

	burnEvents := make([]api.BurnEvent, len(untypedEvents))
	for i, untypedEvent := range untypedEvents {
		burnEvents[i] = untypedEvent.(api.BurnEvent)
	}

	return burnEvents, newWatermark, nil
}

// getSolanaBurnEventFromSig fetches and decodes burn events from a single Solana transaction signature.
func (o *Oracle) getSolanaBurnEventFromSig(sigStr string, programID string) (*api.BurnEvent, error) {
	program, err := solana.PublicKeyFromBase58(programID)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain program public key for burn event backfill: %w", err)
	}

	sig, err := solana.SignatureFromBase58(sigStr)
	if err != nil {
		return nil, fmt.Errorf("invalid signature string for backfill: %w", err)
	}

	v0 := uint64(0)
	txResult, err := o.solanaClient.GetTransaction(
		context.Background(),
		sig,
		&solrpc.GetTransactionOpts{
			Encoding:                       solana.EncodingBase64,
			Commitment:                     solrpc.CommitmentConfirmed,
			MaxSupportedTransactionVersion: &v0,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction for backfill sig %s: %w", sig, err)
	}
	if txResult == nil {
		return nil, fmt.Errorf("nil transaction result for backfill sig %s", sig)
	}

	chainID := sidecartypes.SolanaCAIP2[o.Config.Network]

	// This is for zentp (ROCK) burns for now.
	events, err := rock_spl_token.DecodeEvents(txResult, program)
	if err != nil {
		slog.Warn("Failed to decode Solana ROCK burn events for backfill tx. Event is likely of an unrelated type.", "tx", sig, "error", err)
		return nil, err
	}

	for logIndex, event := range events {
		if event.Name == "TokenRedemption" {
			eventData, ok := event.Data.(*rock_spl_token.TokenRedemptionEventData)
			if !ok {
				slog.Warn("Type assertion failed for Solana ROCK TokenRedemptionEventData on backfill tx", "tx", sig)
				continue
			}
			burnEvent := &api.BurnEvent{
				TxID:            sig.String(),
				LogIndex:        uint64(logIndex),
				ChainID:         chainID,
				DestinationAddr: eventData.DestAddr[:],
				Amount:          eventData.Value,
				IsZenBTC:        false, // This is a ROCK burn
				Height:          uint64(txResult.Slot),
			}
			if o.DebugMode {
				slog.Debug("Backfilled Solana ROCK Burn Event",
					"txID", burnEvent.TxID,
					"logIndex", burnEvent.LogIndex,
					"chainID", burnEvent.ChainID,
					"destinationAddr", fmt.Sprintf("%x", burnEvent.DestinationAddr),
					"amount", burnEvent.Amount)
			}
			// Return the first matching event found.
			return burnEvent, nil
		}
	}

	// No matching event was found in the transaction.
	return nil, nil
}

// processBackfillRequests polls for backfill requests and processes them.
func (o *Oracle) processBackfillRequests(
	wg *sync.WaitGroup,
	update *oracleStateUpdate,
	updateMutex *sync.Mutex,
) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		backfillResp, err := o.zrChainQueryClient.ValidationQueryClient.BackfillRequests(context.Background())
		if err != nil {
			// Don't push to errChan, as this is not a critical failure. Just log it.
			slog.Error("Failed to query backfill requests", "error", err)
			return
		}

		if backfillResp == nil || backfillResp.BackfillRequests == nil || len(backfillResp.BackfillRequests.Requests) == 0 {
			return // No backfill requests
		}
		slog.Info("Found backfill requests to process", "count", len(backfillResp.BackfillRequests.Requests))
		o.handleBackfillRequests(backfillResp.BackfillRequests.Requests, update, updateMutex)
	}()
}

// handleBackfillRequests processes a slice of backfill requests.
func (o *Oracle) handleBackfillRequests(requests []*validationtypes.MsgTriggerEventBackfill, update *oracleStateUpdate, updateMutex *sync.Mutex) {
	if len(requests) == 0 {
		return
	}

	var newBurnEvents []api.BurnEvent

	for i, req := range requests {
		// For now, only handle ZenTP burn events.
		if req.EventType == validationtypes.EventType_EVENT_TYPE_ZENTP_BURN {
			slog.Info("Processing zentp burn backfill request", "txHash", req.TxHash)
			programID := sidecartypes.SolRockProgramID[o.Config.Network]
			event, err := o.getSolanaBurnEventFromSig(req.TxHash, programID)
			if err != nil {
				slog.Error("Error processing backfill request", "txHash", req.TxHash, "error", err)
				continue
			}
			if event != nil {
				newBurnEvents = append(newBurnEvents, *event)
			}

			// Pause between requests to avoid rate-limiting, but not after the final one.
			if i < len(requests)-1 {
				time.Sleep(sidecartypes.SolanaFallbackSleepInterval)
			}
		}
	}

	if len(newBurnEvents) == 0 {
		return
	}

	updateMutex.Lock()
	defer updateMutex.Unlock()

	// Get cleaned events from the persisted state to check against duplicates
	currentState := o.currentState.Load().(*sidecartypes.OracleState)

	// Use helper function to merge backfilled events with existing ones
	update.solanaBurnEvents = mergeNewBurnEvents(update.solanaBurnEvents, currentState.CleanedSolanaBurnEvents, newBurnEvents, "backfilled Solana")
}

// Helper to get typed last processed Solana signature
func (o *Oracle) GetLastProcessedSolSignature(eventType sidecartypes.SolanaEventType) solana.Signature {
	var sigStr string
	switch eventType {
	case sidecartypes.SolRockMint:
		sigStr = o.lastSolRockMintSigStr
	case sidecartypes.SolZenBTCMint:
		sigStr = o.lastSolZenBTCMintSigStr
	case sidecartypes.SolZenBTCBurn:
		sigStr = o.lastSolZenBTCBurnSigStr
	case sidecartypes.SolRockBurn:
		sigStr = o.lastSolRockBurnSigStr
	default:
		slog.Warn("Unknown Solana event type for GetLastProcessedSolSignature", "eventType", eventType)
		return solana.Signature{} // Return zero signature for unknown types
	}

	if sigStr == "" {
		return solana.Signature{} // No signature stored yet, so zero value
	}
	sig, err := solana.SignatureFromBase58(sigStr)
	if err != nil {
		// Log the error but return a zero signature to proceed as if no prior sig exists
		slog.Warn("Could not parse stored signature string for event type. Treating as no prior signature.", "sigString", sigStr, "eventType", eventType, "error", err)
		return solana.Signature{}
	}
	return sig
}

// Helper function to parse RPC response ID into request index
func parseRPCResponseID(resp *jsonrpc.RPCResponse, eventType string) (int, bool) {
	if resp == nil {
		slog.Error("Nil RPCResponse object in sub-batch response", "eventType", eventType)
		return 0, false
	}

	var requestIndex int
	switch id := resp.ID.(type) {
	case float64: // JSON numbers often decode to float64
		requestIndex = int(id)
	case int:
		requestIndex = id
	case uint64: // Match the type we put in the request
		requestIndex = int(id)
	case json.Number:
		idInt64, err := id.Int64()
		if err != nil {
			slog.Error("Failed to convert json.Number ID to int64", "eventType", eventType, "error", err)
			return 0, false
		}
		requestIndex = int(idInt64)
	default:
		slog.Error("Invalid response ID type received", "idType", fmt.Sprintf("%T", resp.ID), "eventType", eventType)
		return 0, false
	}
	return requestIndex, true
}

// Helper function to validate request index bounds
func validateRequestIndex(requestIndex int, batchSize int, eventType string) bool {
	if requestIndex < 0 || requestIndex >= batchSize {
		slog.Error("Invalid response ID received for sub-batch", "requestIndex", requestIndex, "eventType", eventType)
		return false
	}
	return true
}

// Helper function to generate event keys for deduplication
func generateBurnEventKey(event api.BurnEvent) string {
	return fmt.Sprintf("%s-%s-%d", event.ChainID, event.TxID, event.LogIndex)
}

func generateMintEventKey(event api.SolanaMintEvent) string {
	return base64.StdEncoding.EncodeToString(event.SigHash)
}

// Helper function to check if events already exist and merge new ones
func mergeNewBurnEvents(existingEvents []api.BurnEvent, cleanedEvents map[string]bool, newEvents []api.BurnEvent, eventTypeName string) []api.BurnEvent {
	// Create a map of existing events for quick lookup
	existingEventKeys := make(map[string]bool)
	for _, event := range existingEvents {
		key := generateBurnEventKey(event)
		existingEventKeys[key] = true
	}

	// Also check against already cleaned events
	for key := range cleanedEvents {
		existingEventKeys[key] = true
	}

	// Start with existing events
	mergedEvents := make([]api.BurnEvent, len(existingEvents))
	copy(mergedEvents, existingEvents)

	// Add new events if they don't already exist
	for _, event := range newEvents {
		key := generateBurnEventKey(event)
		if !existingEventKeys[key] {
			mergedEvents = append(mergedEvents, event)
			slog.Info("Added burn event to state", "type", eventTypeName, "txID", event.TxID)
		} else {
			slog.Debug("Skipping already present burn event", "type", eventTypeName, "txID", event.TxID)
		}
	}

	return mergedEvents
}

// Helper function to filter new signatures for processing
func (o *Oracle) filterNewSignatures(allSignatures []*solrpc.TransactionSignature, lastKnownSig solana.Signature, eventTypeName string, newestSigFromNode solana.Signature, limit int) ([]*solrpc.TransactionSignature, int) {
	newSignaturesToFetchDetails := make([]*solrpc.TransactionSignature, 0)

	// Filter signatures: find signatures newer than the last one we processed.
	var signaturesInspected int
	for _, sigInfo := range allSignatures {
		signaturesInspected++
		if !lastKnownSig.IsZero() && sigInfo.Signature == lastKnownSig {
			break // Found the last processed signature, stop collecting.
		}
		newSignaturesToFetchDetails = append(newSignaturesToFetchDetails, sigInfo)
	}

	if len(newSignaturesToFetchDetails) == 0 {
		if !lastKnownSig.IsZero() {
			slog.Info("No new signatures found since last processed",
				"eventType", eventTypeName,
				"lastSig", lastKnownSig.String(),
				"inspected", signaturesInspected,
				"total", limit,
				"newest", newestSigFromNode.String())
		} else {
			slog.Info("No signatures found in recent transactions",
				"eventType", eventTypeName,
				"recent", limit)
		}
		return newSignaturesToFetchDetails, signaturesInspected
	}

	if !lastKnownSig.IsZero() {
		if len(newSignaturesToFetchDetails) == len(allSignatures) {
			slog.Info("Last processed signature not found in latest transactions, processing full batch",
				"eventType", eventTypeName,
				"lastSig", lastKnownSig.String(),
				"batchSize", len(allSignatures))
		} else {
			slog.Info("Found new potential transactions to inspect",
				"eventType", eventTypeName,
				"newTxCount", len(newSignaturesToFetchDetails),
				"lastSig", lastKnownSig.String(),
				"inspected", signaturesInspected,
				"total", limit)
		}
	} else {
		slog.Info("No previous signature stored, processing recent transactions",
			"eventType", eventTypeName,
			"txCount", len(newSignaturesToFetchDetails),
			"recent", limit)
	}

	// Reverse the slice so we process the oldest *new* signature first.
	for i, j := 0, len(newSignaturesToFetchDetails)-1; i < j; i, j = i+1, j-1 {
		newSignaturesToFetchDetails[i], newSignaturesToFetchDetails[j] = newSignaturesToFetchDetails[j], newSignaturesToFetchDetails[i]
	}

	return newSignaturesToFetchDetails, signaturesInspected
}

// Helper function to merge new mint events
func mergeNewMintEvents(existingEvents []api.SolanaMintEvent, cleanedEvents map[string]bool, newEvents []api.SolanaMintEvent, eventTypeName string) []api.SolanaMintEvent {
	// Create a map of existing events for quick lookup
	existingEventKeys := make(map[string]bool)
	for _, event := range existingEvents {
		key := generateMintEventKey(event)
		existingEventKeys[key] = true
	}

	// Start with existing events
	mergedEvents := make([]api.SolanaMintEvent, len(existingEvents))
	copy(mergedEvents, existingEvents)

	// Add new events if they don't already exist
	for _, event := range newEvents {
		key := generateMintEventKey(event)
		if !existingEventKeys[key] {
			if cleaned, exists := cleanedEvents[key]; !exists || !cleaned {
				mergedEvents = append(mergedEvents, event)
				slog.Info("Added mint event to state", "type", eventTypeName, "txSig", event.TxSig)
			} else {
				slog.Debug("Skipping already present mint event", "type", eventTypeName, "txSig", event.TxSig)
			}
		}
	}

	return mergedEvents
}
