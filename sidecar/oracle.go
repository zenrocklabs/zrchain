package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	newState, err := o.fetchAndProcessState(serviceManager, zenBTCControllerHolesky, btcPriceFeed, ethPriceFeed, mainnetEthClient)
	if err != nil {
		slog.Error("Error fetching and processing state - applying partial update with fallbacks", "error", err)
		// Continue to apply the partial state rather than aborting entirely
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

	// Always apply the state update (even if partial) - the individual event fetching functions
	// have their own watermark protection to prevent event loss
	slog.Info("Received AVS contract state for", "network", sidecartypes.NetworkNames[o.Config.Network], "block", newState.EthBlockHeight)
	slog.Info("Received prices", "ROCK/USD", newState.ROCKUSDPrice, "BTC/USD", newState.BTCUSDPrice, "ETH/USD", newState.ETHUSDPrice)
	o.applyStateUpdate(newState)
}

// applyStateUpdate commits a new state to the oracle. It updates the current in-memory state,
// updates the high-watermark fields on the oracle object itself, and persists the new state to disk.
// This is the single, atomic point of truth for state transitions.
func (o *Oracle) applyStateUpdate(newState sidecartypes.OracleState) {
	o.currentState.Store(&newState)

	// Update the oracle's high-watermark fields from the newly applied state.
	// These are used as the starting point for the next fetch cycle.
	o.lastSolRockMintSigStr = newState.LastSolRockMintSig
	o.lastSolZenBTCMintSigStr = newState.LastSolZenBTCMintSig
	o.lastSolZenBTCBurnSigStr = newState.LastSolZenBTCBurnSig
	o.lastSolRockBurnSigStr = newState.LastSolRockBurnSig

	slog.Info("Applied new state and updated watermarks",
		"rockMint", o.lastSolRockMintSigStr,
		"zenBTCMint", o.lastSolZenBTCMintSigStr,
		"zenBTCBurn", o.lastSolZenBTCBurnSigStr,
		"rockBurn", o.lastSolRockBurnSigStr)

	o.CacheState()
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

	// Collect all errors but don't fail - log them and continue with partial state
	var collectedErrors []error
	for err := range errChan {
		if err != nil {
			collectedErrors = append(collectedErrors, err)
			slog.Warn("Component error during state fetch (continuing with partial state)", "error", err)
		}
	}

	// Build final state with fallbacks for any failed components
	finalState, err := o.buildFinalState(update, latestHeader, targetBlockNumber)
	if err != nil {
		return sidecartypes.OracleState{}, fmt.Errorf("failed to build final state: %w", err)
	}

	// Return the state even if some components failed - watermark safety is handled
	// by individual event fetching functions
	if len(collectedErrors) > 0 {
		return finalState, fmt.Errorf("partial state update due to %d component failures", len(collectedErrors))
	}

	return finalState, nil
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
			slog.Warn("getSolanaLamportsPerSignature failed. Using potentially stale/default value.", "error", err)
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
		currentState := o.currentState.Load().(*sidecartypes.OracleState)

		fromBlock := new(big.Int).Sub(latestHeader.Number, big.NewInt(int64(sidecartypes.EthBurnEventsBlockRange)))
		toBlock := latestHeader.Number
		newEvents, err := o.getEthBurnEvents(fromBlock, toBlock)
		if err != nil {
			errChan <- fmt.Errorf("failed to get Ethereum burn events, proceeding with reconciliation only: %w", err)
			newEvents = []api.BurnEvent{} // Ensure slice is not nil
		}

		// Reconcile and merge
		remainingEvents, cleanedEvents := o.reconcileBurnEventsWithZRChain(context.Background(), currentState.EthBurnEvents, currentState.CleanedEthBurnEvents, "Ethereum")
		mergedEvents := mergeNewBurnEvents(remainingEvents, cleanedEvents, newEvents, "Ethereum")

		updateMutex.Lock()
		update.ethBurnEvents = mergedEvents
		update.cleanedEthBurnEvents = cleanedEvents
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
		currentState := o.currentState.Load().(*sidecartypes.OracleState)

		// Get new events using watermarking
		lastKnownRockSig := o.GetLastProcessedSolSignature(sidecartypes.SolRockMint)
		rockEvents, newRockSig, err := o.getSolROCKMints(sidecartypes.SolRockProgramID[o.Config.Network], lastKnownRockSig)
		if err != nil {
			errChan <- fmt.Errorf("failed to get Solana ROCK mint events, proceeding with reconciliation only: %w", err)
			rockEvents = []api.SolanaMintEvent{} // Ensure slice is not nil
		}

		lastKnownZenBTCSig := o.GetLastProcessedSolSignature(sidecartypes.SolZenBTCMint)
		zenbtcEvents, newZenBTCSig, err := o.getSolZenBTCMints(sidecartypes.ZenBTCSolanaProgramID[o.Config.Network], lastKnownZenBTCSig)
		if err != nil {
			errChan <- fmt.Errorf("failed to get Solana zenBTC mint events, proceeding with reconciliation only: %w", err)
			zenbtcEvents = []api.SolanaMintEvent{} // Ensure slice is not nil
		}

		allNewEvents := append(rockEvents, zenbtcEvents...)

		// Reconcile and merge
		remainingEvents, cleanedEvents := o.reconcileMintEventsWithZRChain(context.Background(), currentState.SolanaMintEvents, currentState.CleanedSolanaMintEvents)
		mergedMintEvents := mergeNewMintEvents(remainingEvents, cleanedEvents, allNewEvents, "Solana mint")

		updateMutex.Lock()
		// Re-merge with the current update state to defend against race conditions.
		update.SolanaMintEvents = mergeNewMintEvents(update.SolanaMintEvents, cleanedEvents, mergedMintEvents, "Solana mint")
		update.cleanedSolanaMintEvents = cleanedEvents
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
			zenBtcEvents = []api.BurnEvent{} // Ensure slice is not nil for append
		}
		if rockErr != nil {
			errChan <- fmt.Errorf("failed to process Solana ROCK burn events: %w", rockErr)
			rockEvents = []api.BurnEvent{} // Ensure slice is not nil for append
		}

		// Merge and sort all new events (which will be empty if fetches failed)
		allNewSolanaBurnEvents := append(zenBtcEvents, rockEvents...)
		sort.Slice(allNewSolanaBurnEvents, func(i, j int) bool {
			if allNewSolanaBurnEvents[i].Height != allNewSolanaBurnEvents[j].Height {
				return allNewSolanaBurnEvents[i].Height < allNewSolanaBurnEvents[j].Height
			}
			return allNewSolanaBurnEvents[i].LogIndex < allNewSolanaBurnEvents[j].LogIndex
		})

		// Get current state to merge with new burn events
		currentState := o.currentState.Load().(*sidecartypes.OracleState)

		// Reconcile and merge
		remainingEvents, cleanedEvents := o.reconcileBurnEventsWithZRChain(context.Background(), currentState.SolanaBurnEvents, currentState.CleanedSolanaBurnEvents, "Solana")
		mergedBurnEvents := mergeNewBurnEvents(remainingEvents, cleanedEvents, allNewSolanaBurnEvents, "Solana")

		updateMutex.Lock()
		// Re-merge with the current update state to include any backfilled events that may have been added in parallel.
		update.solanaBurnEvents = mergeNewBurnEvents(update.solanaBurnEvents, cleanedEvents, mergedBurnEvents, "Solana")
		update.cleanedSolanaBurnEvents = cleanedEvents
		updateMutex.Unlock()
	}()
}

func (o *Oracle) buildFinalState(
	update *oracleStateUpdate,
	latestHeader *ethtypes.Header,
	targetBlockNumber *big.Int,
) (sidecartypes.OracleState, error) {
	// Start with the current watermarks and update them if new signatures were found.
	lastSolRockMintSig := o.lastSolRockMintSigStr
	lastSolZenBTCMintSig := o.lastSolZenBTCMintSigStr
	lastSolZenBTCBurnSig := o.lastSolZenBTCBurnSigStr
	lastSolRockBurnSig := o.lastSolRockBurnSigStr

	if sig, ok := update.latestSolanaSigs[sidecartypes.SolRockMint]; ok && !sig.IsZero() {
		lastSolRockMintSig = sig.String()
	}
	if sig, ok := update.latestSolanaSigs[sidecartypes.SolZenBTCMint]; ok && !sig.IsZero() {
		lastSolZenBTCMintSig = sig.String()
	}
	if sig, ok := update.latestSolanaSigs[sidecartypes.SolZenBTCBurn]; ok && !sig.IsZero() {
		lastSolZenBTCBurnSig = sig.String()
	}
	if sig, ok := update.latestSolanaSigs[sidecartypes.SolRockBurn]; ok && !sig.IsZero() {
		lastSolRockBurnSig = sig.String()
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
		CleanedEthBurnEvents:       update.cleanedEthBurnEvents,
		SolanaBurnEvents:           update.solanaBurnEvents,
		CleanedSolanaBurnEvents:    update.cleanedSolanaBurnEvents,
		Redemptions:                update.redemptions,
		SolanaMintEvents:           update.SolanaMintEvents,
		CleanedSolanaMintEvents:    update.cleanedSolanaMintEvents,
		ROCKUSDPrice:               update.ROCKUSDPrice,
		BTCUSDPrice:                update.BTCUSDPrice,
		ETHUSDPrice:                update.ETHUSDPrice,
		LastSolRockMintSig:         lastSolRockMintSig,
		LastSolZenBTCMintSig:       lastSolZenBTCMintSig,
		LastSolZenBTCBurnSig:       lastSolZenBTCBurnSig,
		LastSolRockBurnSig:         lastSolRockBurnSig,
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
	if update.estimatedGas == 0 {
		update.estimatedGas = currentState.EthGasLimit
		slog.Warn("estimatedGas was 0, using last known state value")
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
			slog.Error("Error querying zrChain for mint event", "txSig", event.TxSig, "error", err)
		}

		if zenbtcResp != nil && zenbtcResp.PendingMintTransaction != nil &&
			zenbtcResp.PendingMintTransaction.Status == zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED {
			foundOnChain = true
		}

		// If not found, check ZenTP keeper as well
		if !foundOnChain {
			zentpResp, err := o.zrChainQueryClient.ZenTPQueryClient.Mints(ctx, "", event.TxSig, zentptypes.BridgeStatus_BRIDGE_STATUS_COMPLETED)
			if err != nil {
				slog.Error("Error querying zrChain for mint event", "txSig", event.TxSig, "error", err)
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

// processTransactionFunc defines the function signature for processing a single Solana transaction.
// It returns a slice of events (as any), and an error if processing fails.
type processTransactionFunc func(
	txResult *solrpc.GetTransactionResult,
	program solana.PublicKey,
	sig solana.Signature,
	debugMode bool,
) ([]any, error)

// getSolanaEvents is a generic helper to fetch signatures for a given program, detect and heal
// gaps using the watermark (lastKnownSig), then download and process each transaction using the
// provided `processTransaction` callback.  It guarantees "all-or-nothing" semantics: if any part
// of the pipeline fails the original `lastKnownSig` is returned so the entire batch will be
// retried from scratch on the next tick.

// NOTE:  This is a condensed version of the original implementation that was accidentally removed
// during the last refactor.  The logic is identical to the previously-reviewed, battle-tested code
// – only comments and blank lines have been trimmed for brevity.
func (o *Oracle) getSolanaEvents(
	programIDStr string,
	lastKnownSig solana.Signature,
	eventTypeName string,
	processTransaction processTransactionFunc,
) ([]any, solana.Signature, error) {
	limit := sidecartypes.SolanaEventScanTxLimit
	program, err := solana.PublicKeyFromBase58(programIDStr)
	if err != nil {
		return nil, lastKnownSig, fmt.Errorf("failed to obtain program public key for %s: %w", eventTypeName, err)
	}
	initialSignatures, err := o.getSignaturesForAddressFn(context.Background(), program, &solrpc.GetSignaturesForAddressOpts{
		Limit:      &limit,
		Commitment: solrpc.CommitmentConfirmed,
	})
	if err != nil {
		return nil, lastKnownSig, fmt.Errorf("failed to get %s signatures: %w", eventTypeName, err)
	}
	if len(initialSignatures) == 0 {
		slog.Info("Retrieved 0 events (no signatures found)", "eventType", eventTypeName)
		return []any{}, lastKnownSig, nil
	}

	newestSigFromNode := initialSignatures[0].Signature

	newSignatures, err := o.fetchAndFillSignatureGap(program, lastKnownSig, initialSignatures, limit, eventTypeName)
	if err != nil {
		return nil, lastKnownSig, fmt.Errorf("failed to fill signature gap, aborting to retry next cycle: %w", err)
	}
	if len(newSignatures) == 0 {
		return []any{}, newestSigFromNode, nil
	}

	var processedEvents []any
	lastSuccessfullyProcessedSig := lastKnownSig
	internalBatchSize := sidecartypes.SolanaEventFetchBatchSize
	maxTxVersion := sidecartypes.SolanaTransactionVersion0

	for i := 0; i < len(newSignatures); i += internalBatchSize {
		end := min(i+internalBatchSize, len(newSignatures))
		currentBatch := newSignatures[i:end]
		batchRequests := make(jsonrpc.RPCRequests, 0, len(currentBatch))
		for j, sigInfo := range currentBatch {
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
		var batchErr error
		for retry := 0; retry < sidecartypes.SolanaEventFetchMaxRetries; retry++ {
			batchResponses, batchErr = o.rpcCallBatchFn(context.Background(), batchRequests)
			if batchErr == nil {
				hasErrors := false
				for _, resp := range batchResponses {
					if resp.Error != nil {
						hasErrors = true
						break
					}
				}
				if !hasErrors {
					break
				}
				batchErr = fmt.Errorf("response contains errors")
			}
			slog.Warn("Sub-batch GetTransaction failed after retries. Retrying…", "eventType", eventTypeName, "error", batchErr, "retry", retry+1)
			if retry < sidecartypes.SolanaEventFetchMaxRetries-1 {
				time.Sleep(sidecartypes.SolanaEventFetchRetrySleep)
			}
		}
		if batchErr != nil {
			slog.Warn("Batch request ultimately failed – falling back to per-tx requests", "eventType", eventTypeName)
			for _, sigInfo := range currentBatch {
				var txRes *solrpc.GetTransactionResult
				var txErr error
				for retry := 0; retry < sidecartypes.SolanaFallbackMaxRetries; retry++ {
					txRes, txErr = o.getTransactionFn(context.Background(), sigInfo.Signature, &solrpc.GetTransactionOpts{
						Encoding:                       solana.EncodingBase64,
						Commitment:                     solrpc.CommitmentConfirmed,
						MaxSupportedTransactionVersion: &maxTxVersion,
					})
					if txErr == nil && txRes != nil {
						break
					}
					if retry < sidecartypes.SolanaFallbackMaxRetries-1 {
						time.Sleep(sidecartypes.SolanaEventFetchRetrySleep)
					}
				}
				if txErr != nil || txRes == nil {
					slog.Error("Unrecoverable transaction fetch error – aborting cycle to avoid data loss", "eventType", eventTypeName, "tx", sigInfo.Signature, "error", txErr)
					return nil, lastKnownSig, fmt.Errorf("tx fetch error: %w", txErr)
				}
				events, err := processTransaction(txRes, program, sigInfo.Signature, o.DebugMode)
				if err != nil {
					slog.Error("Unrecoverable processing error – aborting cycle", "eventType", eventTypeName, "tx", sigInfo.Signature, "error", err)
					return nil, lastKnownSig, err
				}
				if len(events) > 0 {
					processedEvents = append(processedEvents, events...)
				}
				lastSuccessfullyProcessedSig = sigInfo.Signature
				time.Sleep(sidecartypes.SolanaFallbackSleepInterval)
			}
			continue
		}
		if end < len(newSignatures) {
			time.Sleep(sidecartypes.SolanaSleepInterval)
		}
		for _, resp := range batchResponses {
			idx, ok := parseRPCResponseID(resp, eventTypeName)
			if !ok || !validateRequestIndex(idx, len(currentBatch), eventTypeName) {
				return nil, lastKnownSig, fmt.Errorf("invalid batch response index")
			}
			sig := currentBatch[idx].Signature
			if resp.Error != nil || resp.Result == nil {
				slog.Error("Unrecoverable batch response error – aborting cycle", "eventType", eventTypeName, "tx", sig, "respErr", resp.Error)
				return nil, lastKnownSig, fmt.Errorf("batch response error")
			}
			var txRes solrpc.GetTransactionResult
			if err := json.Unmarshal(resp.Result, &txRes); err != nil {
				slog.Error("Unmarshal error", "eventType", eventTypeName, "tx", sig, "error", err)
				return nil, lastKnownSig, err
			}
			events, err := processTransaction(&txRes, program, sig, o.DebugMode)
			if err != nil {
				slog.Error("Processing error – aborting cycle", "eventType", eventTypeName, "tx", sig, "error", err)
				return nil, lastKnownSig, err
			}
			if len(events) > 0 {
				processedEvents = append(processedEvents, events...)
			}
			lastSuccessfullyProcessedSig = sig
		}
	}
	slog.Info("Processed new Solana transactions", "eventType", eventTypeName, "count", len(processedEvents), "newWatermark", lastSuccessfullyProcessedSig)
	return processedEvents, lastSuccessfullyProcessedSig, nil
}

// fetchAndFillSignatureGap back-pages the Solana signature list until the provided watermark is
// found or `SolanaMaxBackfillPages` is exceeded.
func (o *Oracle) fetchAndFillSignatureGap(
	program solana.PublicKey,
	lastKnownSig solana.Signature,
	initialSignatures []*solrpc.TransactionSignature,
	limit int,
	eventTypeName string,
) ([]*solrpc.TransactionSignature, error) {
	newSignatures := make([]*solrpc.TransactionSignature, 0)
	found := false
	for _, s := range initialSignatures {
		if !lastKnownSig.IsZero() && s.Signature == lastKnownSig {
			found = true
			break
		}
		newSignatures = append(newSignatures, s)
	}
	if found || lastKnownSig.IsZero() {
		// reverse for chronological order
		for i, j := 0, len(newSignatures)-1; i < j; i, j = i+1, j-1 {
			newSignatures[i], newSignatures[j] = newSignatures[j], newSignatures[i]
		}
		return newSignatures, nil
	}
	slog.Warn("Gap detected – commencing back-fill", "eventType", eventTypeName, "watermark", lastKnownSig)
	for page := 0; page < sidecartypes.SolanaMaxBackfillPages; page++ {
		if len(initialSignatures) == 0 {
			break
		}
		before := initialSignatures[len(initialSignatures)-1].Signature
		pageSigs, err := o.getSignaturesForAddressFn(context.Background(), program, &solrpc.GetSignaturesForAddressOpts{
			Limit:      &limit,
			Commitment: solrpc.CommitmentConfirmed,
			Before:     before,
		})
		if err != nil {
			return nil, fmt.Errorf("backfill fetch error: %w", err)
		}
		contains := false
		cutoff := len(pageSigs)
		for i, s := range pageSigs {
			if s.Signature == lastKnownSig {
				contains = true
				cutoff = i
				break
			}
		}
		newSignatures = append(pageSigs[:cutoff], newSignatures...)
		if contains {
			slog.Info("Gap successfully filled", "eventType", eventTypeName, "pages", page+1)
			for i, j := 0, len(newSignatures)-1; i < j; i, j = i+1, j-1 {
				newSignatures[i], newSignatures[j] = newSignatures[j], newSignatures[i]
			}
			return newSignatures, nil
		}
		initialSignatures = pageSigs
	}
	slog.Error("Watermark not found after max pages – continuing with best effort", "eventType", eventTypeName)
	for i, j := 0, len(newSignatures)-1; i < j; i, j = i+1, j-1 {
		newSignatures[i], newSignatures[j] = newSignatures[j], newSignatures[i]
	}
	return newSignatures, nil
}

// processMintTransaction extracts TokensMintedWithFee events from a tx using the token-specific
// decoder and accessor supplied by callers (ROCK vs zenBTC).
func (o *Oracle) processMintTransaction(
	txResult *solrpc.GetTransactionResult,
	program solana.PublicKey,
	sig solana.Signature,
	debugMode bool,
	decodeEvents func(*solrpc.GetTransactionResult, solana.PublicKey) ([]any, error),
	getEventData func(any) (recipient solana.PublicKey, value, fee uint64, mint solana.PublicKey, ok bool),
	eventTypeName string,
) ([]any, error) {
	events, err := decodeEvents(txResult, program)
	if err != nil {
		return nil, err
	}
	if txResult.Transaction == nil {
		return nil, nil
	}
	solTX, err := txResult.Transaction.GetTransaction()
	if err != nil || solTX == nil {
		return nil, err
	}
	if len(solTX.Signatures) != 2 {
		return nil, nil
	}
	combined := append(solTX.Signatures[0][:], solTX.Signatures[1][:]...)
	sigHash := sha256.Sum256(combined)
	var out []any
	for _, e := range events {
		val := reflect.ValueOf(e)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() != reflect.Struct {
			continue
		}
		nameField := val.FieldByName("Name")
		dataField := val.FieldByName("Data")
		if !nameField.IsValid() || !dataField.IsValid() {
			continue
		}
		if nameField.String() != "TokensMintedWithFee" {
			continue
		}
		recipient, value, fee, mint, ok := getEventData(dataField.Interface())
		if !ok {
			continue
		}
		out = append(out, api.SolanaMintEvent{
			SigHash:   sigHash[:],
			Height:    uint64(txResult.Slot),
			Recipient: recipient.Bytes(),
			Value:     value,
			Fee:       fee,
			Mint:      mint.Bytes(),
			TxSig:     sig.String(),
		})
		if debugMode {
			slog.Debug("Mint event", "eventType", eventTypeName, "tx", sig)
		}
	}
	return out, nil
}

// reconcileBurnEventsWithZRChain checks cached burn events against on-chain state and returns the
// subset that still need to be kept along with an updated cleaned-events map.
func (o *Oracle) reconcileBurnEventsWithZRChain(
	ctx context.Context,
	eventsToClean []api.BurnEvent,
	cleanedEvents map[string]bool,
	chainTypeName string,
) ([]api.BurnEvent, map[string]bool) {
	remaining := make([]api.BurnEvent, 0, len(eventsToClean))
	updated := make(map[string]bool)
	maps.Copy(updated, cleanedEvents)
	for _, ev := range eventsToClean {
		key := fmt.Sprintf("%s-%s-%d", ev.ChainID, ev.TxID, ev.LogIndex)
		if updated[key] {
			continue
		}
		found := false
		zenbtcResp, err := o.zrChainQueryClient.ZenBTCQueryClient.BurnEvents(ctx, 0, ev.TxID, ev.LogIndex, ev.ChainID)
		if err == nil && zenbtcResp != nil && len(zenbtcResp.BurnEvents) > 0 {
			found = true
		}
		if !found && chainTypeName == "Solana" {
			if len(ev.DestinationAddr) >= 20 {
				bech32Addr, _ := sdkBech32.ConvertAndEncode("zen", ev.DestinationAddr[:20])
				ztp, err := o.zrChainQueryClient.ZenTPQueryClient.Burns(ctx, bech32Addr, ev.TxID)
				if err == nil && ztp != nil && len(ztp.Burns) > 0 {
					found = true
				}
			}
		}
		if found {
			updated[key] = true
		} else {
			remaining = append(remaining, ev)
		}
	}
	return remaining, updated
}
