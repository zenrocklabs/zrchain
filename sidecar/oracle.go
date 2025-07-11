package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"math/big"
	"net/http"
	"reflect"
	"sort"
	"strings"
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
	solrpc "github.com/gagliardetto/solana-go/rpc"
	jsonrpc "github.com/gagliardetto/solana-go/rpc/jsonrpc"
	// Added for bin.Marshal
)

// sendError sends an error to the channel if the context is not done.
// This prevents panics from sending on a closed channel.
func sendError(ctx context.Context, errChan chan<- error, err error) {
	select {
	case <-ctx.Done():
		slog.Warn("Context canceled, dropping error", "err", err)
	case errChan <- err:
	}
}

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

		// Initialize performance optimization fields
		solanaRateLimiter: make(chan struct{}, sidecartypes.SolanaMaxConcurrentRPCCalls), // Configurable concurrent Solana RPC calls
		transactionCache:  make(map[string]*CachedTxResult),
		httpClientPool: &sync.Pool{
			New: func() interface{} {
				return &http.Client{
					Timeout: sidecartypes.SolanaRPCTimeout, // Use longer timeout for Solana operations
					Transport: &http.Transport{
						MaxIdleConns:        100,
						MaxIdleConnsPerHost: 10,
						IdleConnTimeout:     90 * time.Second,
					},
				}
			},
		},
		batchRequestPool: &sync.Pool{
			New: func() interface{} {
				return make(jsonrpc.RPCRequests, 0, sidecartypes.SolanaEventFetchBatchSize)
			},
		},
		eventProcessorPool: &sync.Pool{
			New: func() interface{} {
				return &EventProcessor{
					events: make([]any, 0, 100),
				}
			},
		},
	}

	// Load initial state from cache file
	latestDiskState, historicalStates, err := loadStateDataFromFile(o.Config.StateFile)
	if err != nil {
		slog.Warn("Unable to load oracle state cache from disk - starting with clean state", "file", o.Config.StateFile, "error", err)
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

	// Only initialize Solana-related functions if Solana client is available
	if o.solanaClient != nil {
		o.rpcCallBatchFn = o.solanaClient.RPCCallBatch
		o.getTransactionFn = o.solanaClient.GetTransaction
		o.getSignaturesForAddressFn = o.solanaClient.GetSignaturesForAddressWithOpts
	} else {
		// Set dummy functions that return empty results when Solana is disabled
		o.rpcCallBatchFn = func(ctx context.Context, rpcs jsonrpc.RPCRequests) (jsonrpc.RPCResponses, error) {
			return jsonrpc.RPCResponses{}, nil
		}
		o.getTransactionFn = func(ctx context.Context, signature solana.Signature, opts *solrpc.GetTransactionOpts) (*solrpc.GetTransactionResult, error) {
			return nil, fmt.Errorf("solana functionality disabled")
		}
		o.getSignaturesForAddressFn = func(ctx context.Context, account solana.PublicKey, opts *solrpc.GetSignaturesForAddressOpts) ([]*solrpc.TransactionSignature, error) {
			return []*solrpc.TransactionSignature{}, nil
		}
	}

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

	mainLoopTickerIntervalDuration := sidecartypes.MainLoopTickerInterval
	var tickCancel context.CancelFunc = func() {}
	defer tickCancel()

	// Align the start time to the nearest MainLoopTickerInterval.
	if !o.SkipInitialWait {
		ntpTime, err := ntp.Time("time.google.com")
		if err != nil {
			slog.Error("Failed to fetch NTP time at startup. Cannot proceed.", "error", err)
			panic(fmt.Sprintf("FATAL: Failed to fetch NTP time at startup: %v. Cannot proceed.", err))
		}
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
		var initialTickCtx context.Context
		initialTickCtx, tickCancel = context.WithCancel(ctx)
		go o.processOracleTick(initialTickCtx, serviceManager, zenBTCControllerHolesky, btcPriceFeed, ethPriceFeed, mainnetEthClient, time.Now())
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
			// Cancel the previous tick's processing context. This signals the previous
			// fetchAndProcessState to wrap up and apply its (potentially partial) state.
			tickCancel()

			// Create a new context for the new tick.
			var tickCtx context.Context
			tickCtx, tickCancel = context.WithCancel(ctx)

			// Start the new tick's processing in a goroutine.
			go o.processOracleTick(tickCtx, serviceManager, zenBTCControllerHolesky, btcPriceFeed, ethPriceFeed, mainnetEthClient, tickTime)
		}
	}
}

func (o *Oracle) processOracleTick(
	tickCtx context.Context,
	serviceManager *middleware.ContractZrServiceManager,
	zenBTCControllerHolesky *zenbtc.ZenBTController,
	btcPriceFeed *aggregatorv3.AggregatorV3Interface,
	ethPriceFeed *aggregatorv3.AggregatorV3Interface,
	mainnetEthClient *ethclient.Client,
	tickTime time.Time,
) {
	newState, err := o.fetchAndProcessState(tickCtx, serviceManager, zenBTCControllerHolesky, btcPriceFeed, ethPriceFeed, mainnetEthClient)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			slog.Info("Data fetch time limit reached. Applying partially gathered state to meet tick deadline.", "tickTime", tickTime.Format("15:04:05.00"))
		} else {
			slog.Error("Error fetching and processing state, applying partial update with fallbacks", "error", err)
		}
		// Continue to apply the partial state rather than aborting entirely
	}

	// Always apply the state update (even if partial) - the individual event fetching functions
	// have their own watermark protection to prevent event loss
	slog.Info("Applying state update for tick", "tickTime", tickTime.Format("15:04:05.00"))
	slog.Info("Received AVS contract state for", "network", sidecartypes.NetworkNames[o.Config.Network], "block", newState.EthBlockHeight)
	slog.Info("Received prices", "ROCK/USD", newState.ROCKUSDPrice, "BTC/USD", newState.BTCUSDPrice, "ETH/USD", newState.ETHUSDPrice)
	o.applyStateUpdate(newState)
}

// applyStateUpdate commits a new state to the oracle. It updates the current in-memory state,
// updates the high-watermark fields on the oracle object itself, and persists the new state to disk.
// This is the single, atomic point of truth for state transitions.
func (o *Oracle) applyStateUpdate(newState sidecartypes.OracleState) {
	// Log watermark changes for debugging
	oldRockMint := o.lastSolRockMintSigStr
	oldZenBTCMint := o.lastSolZenBTCMintSigStr
	oldZenBTCBurn := o.lastSolZenBTCBurnSigStr
	oldRockBurn := o.lastSolRockBurnSigStr

	o.currentState.Store(&newState)

	// Log event counts in each state field every tick
	slog.Info("State event counts per tick",
		"ethBurnEvents", len(newState.EthBurnEvents),
		"cleanedEthBurnEvents", len(newState.CleanedEthBurnEvents),
		"solanaBurnEvents", len(newState.SolanaBurnEvents),
		"cleanedSolanaBurnEvents", len(newState.CleanedSolanaBurnEvents),
		"solanaMintEvents", len(newState.SolanaMintEvents),
		"cleanedSolanaMintEvents", len(newState.CleanedSolanaMintEvents),
		"redemptions", len(newState.Redemptions))

	// Update the oracle's high-watermark fields from the newly applied state.
	// These are used as the starting point for the next fetch cycle.
	o.lastSolRockMintSigStr = newState.LastSolRockMintSig
	o.lastSolZenBTCMintSigStr = newState.LastSolZenBTCMintSig
	o.lastSolZenBTCBurnSigStr = newState.LastSolZenBTCBurnSig
	o.lastSolRockBurnSigStr = newState.LastSolRockBurnSig

	// Log any watermark changes
	watermarkChanged := false
	if oldRockMint != o.lastSolRockMintSigStr {
		slog.Info("Updated ROCK mint watermark", "old", oldRockMint, "new", o.lastSolRockMintSigStr)
		watermarkChanged = true
	}
	if oldZenBTCMint != o.lastSolZenBTCMintSigStr {
		slog.Info("Updated zenBTC mint watermark", "old", oldZenBTCMint, "new", o.lastSolZenBTCMintSigStr)
		watermarkChanged = true
	}
	if oldZenBTCBurn != o.lastSolZenBTCBurnSigStr {
		slog.Info("Updated zenBTC burn watermark", "old", oldZenBTCBurn, "new", o.lastSolZenBTCBurnSigStr)
		watermarkChanged = true
	}
	if oldRockBurn != o.lastSolRockBurnSigStr {
		slog.Info("Updated ROCK burn watermark", "old", oldRockBurn, "new", o.lastSolRockBurnSigStr)
		watermarkChanged = true
	}

	// Only log the comprehensive watermark summary if any watermarks actually changed
	if watermarkChanged {
		slog.Info("Applied new state and updated watermarks",
			"rockMint", o.lastSolRockMintSigStr,
			"zenBTCMint", o.lastSolZenBTCMintSigStr,
			"zenBTCBurn", o.lastSolZenBTCBurnSigStr,
			"rockBurn", o.lastSolRockBurnSigStr)
	}

	o.CacheState()
}

func (o *Oracle) fetchAndProcessState(
	tickCtx context.Context,
	serviceManager *middleware.ContractZrServiceManager,
	zenBTCControllerHolesky *zenbtc.ZenBTController,
	btcPriceFeed *aggregatorv3.AggregatorV3Interface,
	ethPriceFeed *aggregatorv3.AggregatorV3Interface,
	tempEthClient *ethclient.Client,
) (sidecartypes.OracleState, error) {
	var wg sync.WaitGroup

	slog.Info("Retrieving latest header", "network", sidecartypes.NetworkNames[o.Config.Network], "time", time.Now().Format("15:04:05.00"))
	latestHeader, err := o.EthClient.HeaderByNumber(tickCtx, nil)
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

	// Use a separate context for the goroutines that can be canceled
	// if the main tick context is canceled.
	routinesCtx, cancelRoutines := context.WithCancel(tickCtx)
	defer cancelRoutines()

	// Fetch Ethereum contract data (AVS delegations and redemptions on EigenLayer)
	o.fetchEthereumContractData(routinesCtx, &wg, serviceManager, zenBTCControllerHolesky, targetBlockNumber, update, &updateMutex, errChan)

	// Fetch network data (gas estimates, tips, Solana fees)
	o.fetchNetworkData(routinesCtx, &wg, update, &updateMutex, errChan)

	// Fetch price data (ROCK, BTC, ETH)
	o.fetchPriceData(routinesCtx, &wg, btcPriceFeed, ethPriceFeed, tempEthClient, update, &updateMutex, errChan)

	// Fetch zenBTC burn events from Ethereum
	o.fetchEthereumBurnEvents(routinesCtx, &wg, latestHeader, update, &updateMutex, errChan)

	// Fetch Solana mint events for zenBTC (only if Solana is enabled)
	if o.solanaClient != nil {
		o.processSolanaMintEvents(routinesCtx, &wg, update, &updateMutex, errChan)
	}

	// Fetch Solana burn events for zenBTC and ROCK (only if Solana is enabled)
	if o.solanaClient != nil {
		o.fetchSolanaBurnEvents(routinesCtx, &wg, update, &updateMutex, errChan)
	}

	// Fetch and populate backfill requests from zrChain
	o.processBackfillRequests(routinesCtx, &wg, update, &updateMutex)

	// Wait for all goroutines to complete, or for the tick to be canceled.
	waitChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitChan)
	}()

	select {
	case <-waitChan:
		// All tasks completed normally.
	case <-tickCtx.Done():
		// Tick was canceled. The goroutines have been notified via cancelRoutines().
		slog.Warn("Oracle tick deadline approached. Applying state updates with partial data to maintain synchronization timing.")
		// We must wait for the goroutines to finish before closing the error channel.
		<-waitChan
	}

	close(errChan)

	// Collect all errors but don't fail - log them and continue with partial state
	var collectedErrors []error
	for err := range errChan {
		if err != nil {
			collectedErrors = append(collectedErrors, err)
			// Don't log an error if it's just a context cancellation, as this is expected when the tick times out.
			if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				slog.Warn("Component error during state fetch (continuing with partial state)", "error", err)
			}
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

	// If the original tick context was canceled, return that error so the caller knows.
	if tickCtx.Err() != nil {
		return finalState, tickCtx.Err()
	}

	return finalState, nil
}

func (o *Oracle) fetchEthereumContractData(
	ctx context.Context,
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
		delegations, err := o.getServiceManagerState(ctx, serviceManager, targetBlockNumber)
		if err != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to get contract state: %w", err))
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
		redemptions, err := o.getRedemptions(ctx, zenBTCControllerHolesky, targetBlockNumber)
		if err != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to get zenBTC contract state: %w", err))
			return
		}
		updateMutex.Lock()
		update.redemptions = redemptions
		updateMutex.Unlock()
	}()
}

func (o *Oracle) fetchNetworkData(
	ctx context.Context,
	wg *sync.WaitGroup,
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
			sendError(ctx, errChan, fmt.Errorf("failed to get suggested priority fee: %w", err))
			return
		}
		updateMutex.Lock()
		update.suggestedTip = suggestedTip
		updateMutex.Unlock()
	}()

	// Estimates the gas required for a zenBTC stake call on Ethereum.
	wg.Add(1)
	go func() {
		defer wg.Done()
		stakeCallData, err := validationkeeper.EncodeStakeCallData(big.NewInt(1000000000))
		if err != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to encode stake call data: %w", err))
			return
		}
		addr := common.HexToAddress(sidecartypes.ZenBTCControllerAddresses[o.Config.Network])
		estimatedGas, err := o.EthClient.EstimateGas(context.Background(), ethereum.CallMsg{
			From: common.HexToAddress(sidecartypes.WhitelistedRoleAddresses[o.Config.Network]),
			To:   &addr,
			Data: stakeCallData,
		})
		if err != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to estimate gas for stake call: %w", err))
			return
		}
		updateMutex.Lock()
		update.estimatedGas = (estimatedGas * sidecartypes.GasEstimationBuffer) / 100
		updateMutex.Unlock()
	}()
}

func (o *Oracle) fetchPriceData(
	ctx context.Context,
	wg *sync.WaitGroup,
	btcPriceFeed *aggregatorv3.AggregatorV3Interface,
	ethPriceFeed *aggregatorv3.AggregatorV3Interface,
	tempEthClient *ethclient.Client,
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
			sendError(ctx, errChan, fmt.Errorf("failed to retrieve ROCK price data: %w", err))
			return
		}
		defer resp.Body.Close()

		var priceData []PriceData
		if err := json.NewDecoder(resp.Body).Decode(&priceData); err != nil || len(priceData) == 0 {
			sendError(ctx, errChan, fmt.Errorf("failed to decode ROCK price data or empty data: %w", err))
			return
		}
		priceDec, err := math.LegacyNewDecFromStr(priceData[0].Last)
		if err != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to parse ROCK price data: %w", err))
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
			sendError(ctx, errChan, fmt.Errorf("failed to fetch latest mainnet block: %w", err))
			return
		}
		targetBlockNumberMainnet := new(big.Int).Sub(mainnetLatestHeader.Number, big.NewInt(sidecartypes.EthBlocksBeforeFinality))

		if btcPriceFeed == nil || ethPriceFeed == nil {
			sendError(ctx, errChan, fmt.Errorf("BTC or ETH price feed not initialized"))
			return
		}

		btcPrice, err := o.fetchPrice(btcPriceFeed, targetBlockNumberMainnet)
		if err != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to fetch BTC price: %w", err))
			return
		}

		ethPrice, err := o.fetchPrice(ethPriceFeed, targetBlockNumberMainnet)
		if err != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to fetch ETH price: %w", err))
			return
		}

		updateMutex.Lock()
		update.BTCUSDPrice = btcPrice
		update.ETHUSDPrice = ethPrice
		updateMutex.Unlock()
	}()
}

func (o *Oracle) fetchEthereumBurnEvents(
	ctx context.Context,
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
		newEvents, err := o.getEthBurnEvents(ctx, fromBlock, toBlock)
		if err != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to get Ethereum burn events, proceeding with reconciliation only: %w", err))
			newEvents = []api.BurnEvent{} // Ensure slice is not nil
		}

		// Reconcile and merge
		remainingEvents, cleanedEvents := o.reconcileBurnEventsWithZRChain(ctx, currentState.EthBurnEvents, currentState.CleanedEthBurnEvents, "Ethereum")
		mergedEvents := mergeNewBurnEvents(remainingEvents, cleanedEvents, newEvents, "Ethereum")

		updateMutex.Lock()
		update.ethBurnEvents = mergedEvents
		update.cleanedEthBurnEvents = cleanedEvents
		updateMutex.Unlock()
	}()
}

func (o *Oracle) processSolanaMintEvents(
	ctx context.Context,
	wg *sync.WaitGroup,
	update *oracleStateUpdate,
	updateMutex *sync.Mutex,
	errChan chan<- error,
) {
	// Fetches new ROCK and zenBTC mint events from Solana in parallel since the last processed signature,
	// and merges them with the existing cached events.
	wg.Add(1)
	go func() {
		defer wg.Done()
		currentState := o.currentState.Load().(*sidecartypes.OracleState)

		// Parallel fetch of ROCK and zenBTC mint events
		var rockEvents, zenbtcEvents []api.SolanaMintEvent
		var newRockSig, newZenBTCSig solana.Signature
		var rockErr, zenbtcErr error
		var mintWg sync.WaitGroup

		// Fetch ROCK mint events in parallel
		mintWg.Add(1)
		go func() {
			defer mintWg.Done()
			lastKnownRockSig := o.GetLastProcessedSolSignature(sidecartypes.SolRockMint)
			rockEvents, newRockSig, rockErr = o.getSolROCKMints(ctx, sidecartypes.SolRockProgramID[o.Config.Network], lastKnownRockSig)
		}()

		// Fetch zenBTC mint events in parallel
		mintWg.Add(1)
		go func() {
			defer mintWg.Done()
			lastKnownZenBTCSig := o.GetLastProcessedSolSignature(sidecartypes.SolZenBTCMint)
			zenbtcEvents, newZenBTCSig, zenbtcErr = o.getSolZenBTCMints(ctx, sidecartypes.ZenBTCSolanaProgramID[o.Config.Network], lastKnownZenBTCSig)
		}()

		mintWg.Wait()

		// Handle errors after parallel execution
		if rockErr != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to get Solana ROCK mint events, applying partial results: %w", rockErr))
		}
		if zenbtcErr != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to get Solana zenBTC mint events, applying partial results: %w", zenbtcErr))
		}

		allNewEvents := append(rockEvents, zenbtcEvents...)

		// Reconcile and merge
		slog.Info("Before reconciliation",
			"pendingMintEvents", len(currentState.SolanaMintEvents),
			"cleanedMintEventsMap", len(currentState.CleanedSolanaMintEvents))
		remainingEvents, cleanedEvents, reconcileErr := o.reconcileMintEventsWithZRChain(ctx, currentState.SolanaMintEvents, currentState.CleanedSolanaMintEvents)

		// Only log reconciliation results if there was activity or errors
		if reconcileErr != nil || len(remainingEvents) != len(currentState.SolanaMintEvents) || len(cleanedEvents) != len(currentState.CleanedSolanaMintEvents) {
			slog.Info("After reconciliation",
				"remainingEvents", len(remainingEvents),
				"cleanedEventsMap", len(cleanedEvents),
				"reconcileErr", reconcileErr)
		}

		if reconcileErr != nil {
			// Reconciliation failed, so 'remainingEvents' contains all the previously pending events.
			slog.Warn("Failed to reconcile Solana mint events with zrchain - retaining unconfirmed Solana mint events for next cycle", "error", reconcileErr)
		}

		// `remainingEvents` now contains all events that are still pending (either all of them if reconcile failed, or a subset if it succeeded).
		// Merge the newly fetched events (`allNewEvents`) into this list.
		// This call will produce a non-confusing log because `remainingEvents` is passed as the existing set.
		mergedMintEvents := mergeNewMintEvents(remainingEvents, cleanedEvents, allNewEvents, "Solana mint")

		updateMutex.Lock()
		// Since there is no backfill for mint events, there is no race condition with other goroutines
		// modifying update.SolanaMintEvents. We can assign the final list directly.
		// This avoids a second, confusingly-logged merge.
		initialUpdateCount := len(update.SolanaMintEvents)
		update.SolanaMintEvents = mergedMintEvents
		update.cleanedSolanaMintEvents = cleanedEvents

		// Only log if there were actual changes in the final merge
		if len(update.SolanaMintEvents) != initialUpdateCount || len(mergedMintEvents) > 0 {
			slog.Info("Final merge and state update completed",
				"finalMintEvents", len(update.SolanaMintEvents),
				"finalCleanedEvents", len(update.cleanedSolanaMintEvents),
				"addedToUpdate", len(update.SolanaMintEvents)-initialUpdateCount)
		}
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
	ctx context.Context,
	wg *sync.WaitGroup,
	update *oracleStateUpdate,
	updateMutex *sync.Mutex,
	errChan chan<- error,
) {
	// Fetch zenBTC and ROCK burn events from Solana in parallel
	wg.Add(1)
	go func() {
		defer wg.Done()

		var zenBtcEvents, rockEvents []api.BurnEvent
		var zenBtcErr, rockErr error
		var burnWg sync.WaitGroup

		// Fetches new zenBTC burn events from Solana in parallel
		burnWg.Add(1)
		go func() {
			defer burnWg.Done()
			lastKnownSig := o.GetLastProcessedSolSignature(sidecartypes.SolZenBTCBurn)
			var newestSig solana.Signature
			zenBtcEvents, newestSig, zenBtcErr = o.getSolanaZenBTCBurnEventsFn(ctx, sidecartypes.ZenBTCSolanaProgramID[o.Config.Network], lastKnownSig)
			if !newestSig.IsZero() {
				updateMutex.Lock()
				update.latestSolanaSigs[sidecartypes.SolZenBTCBurn] = newestSig
				updateMutex.Unlock()
			}
		}()

		// Fetches new ROCK burn events from Solana in parallel
		burnWg.Add(1)
		go func() {
			defer burnWg.Done()
			lastKnownSig := o.GetLastProcessedSolSignature(sidecartypes.SolRockBurn)
			var newestSig solana.Signature
			rockEvents, newestSig, rockErr = o.getSolanaRockBurnEventsFn(ctx, sidecartypes.SolRockProgramID[o.Config.Network], lastKnownSig)
			if !newestSig.IsZero() {
				updateMutex.Lock()
				update.latestSolanaSigs[sidecartypes.SolRockBurn] = newestSig
				updateMutex.Unlock()
			}
		}()

		burnWg.Wait() // Wait for both parallel burn fetches to complete

		if zenBtcErr != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to process Solana zenBTC burn events, applying partial results: %w", zenBtcErr))
		}
		if rockErr != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to process Solana ROCK burn events, applying partial results: %w", rockErr))
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

		// Reconcile with zrChain to see which of the pending events have been processed.
		remainingEvents, cleanedEvents := o.reconcileBurnEventsWithZRChain(ctx, currentState.SolanaBurnEvents, currentState.CleanedSolanaBurnEvents, "Solana")

		updateMutex.Lock()
		defer updateMutex.Unlock()

		// Manually construct the base list of events to avoid confusing logs from `mergeNewBurnEvents`.
		// The base list consists of events that were added by the backfill process (`update.solanaBurnEvents`)
		// and events from the previous state that have not yet been reconciled (`remainingEvents`).
		// A map is used to efficiently de-duplicate events by their transaction ID.
		combinedEventsMap := make(map[string]api.BurnEvent)

		// Add events that might have been added by the backfill process first.
		for _, event := range update.solanaBurnEvents {
			combinedEventsMap[event.TxID] = event
		}
		// Add the remaining unreconciled events from the previous state.
		for _, event := range remainingEvents {
			combinedEventsMap[event.TxID] = event
		}

		// Convert the map back to a slice. The order is not important here as the final state is sorted later.
		baseEvents := make([]api.BurnEvent, 0, len(combinedEventsMap))
		for _, event := range combinedEventsMap {
			baseEvents = append(baseEvents, event)
		}

		// Now, merge the newly fetched Solana burn events (`allNewSolanaBurnEvents`) into the
		// de-duplicated list of existing events (`baseEvents`). The log message from this
		// call will now correctly report only the truly new events as "added".
		update.solanaBurnEvents = mergeNewBurnEvents(baseEvents, cleanedEvents, allNewSolanaBurnEvents, "Solana")
		update.cleanedSolanaBurnEvents = cleanedEvents
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
		EigenDelegations:        update.eigenDelegations,
		EthBlockHeight:          targetBlockNumber.Uint64(),
		EthGasLimit:             update.estimatedGas,
		EthBaseFee:              latestHeader.BaseFee.Uint64(),
		EthTipCap:               update.suggestedTip.Uint64(),
		EthBurnEvents:           update.ethBurnEvents,
		CleanedEthBurnEvents:    update.cleanedEthBurnEvents,
		SolanaBurnEvents:        update.solanaBurnEvents,
		CleanedSolanaBurnEvents: update.cleanedSolanaBurnEvents,
		Redemptions:             update.redemptions,
		SolanaMintEvents:        update.SolanaMintEvents,
		CleanedSolanaMintEvents: update.cleanedSolanaMintEvents,
		ROCKUSDPrice:            update.ROCKUSDPrice,
		BTCUSDPrice:             update.BTCUSDPrice,
		ETHUSDPrice:             update.ETHUSDPrice,
		LastSolRockMintSig:      lastSolRockMintSig,
		LastSolZenBTCMintSig:    lastSolZenBTCMintSig,
		LastSolZenBTCBurnSig:    lastSolZenBTCBurnSig,
		LastSolRockBurnSig:      lastSolRockBurnSig,
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
	if update.estimatedGas == 0 {
		update.estimatedGas = currentState.EthGasLimit
		slog.Warn("estimatedGas was 0, using last known state value")
	}
}

func (o *Oracle) getServiceManagerState(ctx context.Context, contractInstance *middleware.ContractZrServiceManager, height *big.Int) (map[string]map[string]*big.Int, error) {
	delegations := make(map[string]map[string]*big.Int)

	callOpts := &bind.CallOpts{
		BlockNumber: height,
		Context:     ctx,
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
			// TODO: remove this optimisation when we support multiple EigenLayer operators
			if operator != common.HexToAddress("0x4B2D2fE4DFa633C8a43FcECC05eAE4f4A84EF9f7") {
				continue
			}
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

func (o *Oracle) getEthBurnEvents(ctx context.Context, fromBlock, toBlock *big.Int) ([]api.BurnEvent, error) {
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

func (o *Oracle) getRedemptions(ctx context.Context, contractInstance *zenbtc.ZenBTController, height *big.Int) ([]api.Redemption, error) {
	callOpts := &bind.CallOpts{
		BlockNumber: height,
		Context:     ctx,
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
) ([]api.SolanaMintEvent, map[string]bool, error) {
	remainingEvents := make([]api.SolanaMintEvent, 0)
	updatedCleanedEvents := make(map[string]bool)
	maps.Copy(updatedCleanedEvents, cleanedEvents)

	var zenbtcQueryErrors, zentpQueryErrors int
	var lastZenbtcError, lastZentpError error

	for _, event := range eventsToClean {
		key := base64.StdEncoding.EncodeToString(event.SigHash)
		if _, alreadyCleaned := updatedCleanedEvents[key]; alreadyCleaned {
			continue
		}

		var foundOnChain bool

		// Check ZenBTC keeper
		zenbtcResp, err := o.zrChainQueryClient.ZenBTCQueryClient.PendingMintTransaction(ctx, event.TxSig)
		if err != nil {
			zenbtcQueryErrors++
			lastZenbtcError = err
			// If we fail to query zrChain for this specific event, we keep it in the cache
			// to retry later, but continue processing other events
			remainingEvents = append(remainingEvents, event)
			continue
		}

		if zenbtcResp != nil && zenbtcResp.PendingMintTransaction != nil &&
			zenbtcResp.PendingMintTransaction.Status == zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED {
			foundOnChain = true
		}

		// If not found, check ZenTP keeper as well
		if !foundOnChain {
			zentpResp, err := o.zrChainQueryClient.ZenTPQueryClient.Mints(ctx, "", event.TxSig, zentptypes.BridgeStatus_BRIDGE_STATUS_COMPLETED)
			if err != nil {
				zentpQueryErrors++
				lastZentpError = err
				// If we fail to query ZenTP for this specific event, keep it in cache to retry later
				remainingEvents = append(remainingEvents, event)
				continue
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

	// Log summary of any query errors
	if zenbtcQueryErrors > 0 {
		slog.Warn("Failed to query zrChain ZenBTC for mint events, keeping in cache",
			"failedCount", zenbtcQueryErrors,
			"totalEvents", len(eventsToClean),
			"lastError", lastZenbtcError)
	}
	if zentpQueryErrors > 0 {
		slog.Warn("Failed to query zrChain ZenTP for mint events, keeping in cache",
			"failedCount", zentpQueryErrors,
			"totalEvents", len(eventsToClean),
			"lastError", lastZentpError)
	}

	return remainingEvents, updatedCleanedEvents, nil
}

func (o *Oracle) getSolROCKMints(ctx context.Context, programID string, lastKnownSig solana.Signature) ([]api.SolanaMintEvent, solana.Signature, error) {
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

	untypedEvents, newWatermark, err := o.getSolanaEvents(ctx, programID, lastKnownSig, eventTypeName, processor)

	// Always process partial results, even if an error occurred.
	mintEvents := make([]api.SolanaMintEvent, len(untypedEvents))
	for i, untypedEvent := range untypedEvents {
		mintEvents[i] = untypedEvent.(api.SolanaMintEvent)
	}

	// Return the (potentially partial) events and the updated watermark along with the error.
	return mintEvents, newWatermark, err
}

func (o *Oracle) getSolZenBTCMints(ctx context.Context, programID string, lastKnownSig solana.Signature) ([]api.SolanaMintEvent, solana.Signature, error) {
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

	untypedEvents, newWatermark, err := o.getSolanaEvents(ctx, programID, lastKnownSig, eventTypeName, processor)

	// Always process partial results, even if an error occurred.
	mintEvents := make([]api.SolanaMintEvent, len(untypedEvents))
	for i, untypedEvent := range untypedEvents {
		mintEvents[i] = untypedEvent.(api.SolanaMintEvent)
	}

	// Return the (potentially partial) events and the updated watermark along with the error.
	return mintEvents, newWatermark, err
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
		slog.Info("Processing tx: found events", "eventType", eventTypeName, "tx", sig, "eventCount", len(decodedEvents))
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
				slog.Info("Event details", "eventIndex", i, "eventName", eventNameField.String(), "eventType", fmt.Sprintf("%T", event))
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
				slog.Info("Burn Event",
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
func (o *Oracle) getSolanaZenBTCBurnEvents(ctx context.Context, programID string, lastKnownSig solana.Signature) ([]api.BurnEvent, solana.Signature, error) {
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

	untypedEvents, newWatermark, err := o.getSolanaEvents(ctx, programID, lastKnownSig, eventTypeName, processor)

	// Always process partial results, even if an error occurred.
	burnEvents := make([]api.BurnEvent, len(untypedEvents))
	for i, untypedEvent := range untypedEvents {
		burnEvents[i] = untypedEvent.(api.BurnEvent)
	}

	// Return the (potentially partial) events and the updated watermark along with the error.
	return burnEvents, newWatermark, err
}

// getSolanaRockBurnEvents retrieves Rock burn events from Solana.
func (o *Oracle) getSolanaRockBurnEvents(ctx context.Context, programID string, lastKnownSig solana.Signature) ([]api.BurnEvent, solana.Signature, error) {
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

	untypedEvents, newWatermark, err := o.getSolanaEvents(ctx, programID, lastKnownSig, eventTypeName, processor)

	// Always process partial results, even if an error occurred.
	burnEvents := make([]api.BurnEvent, len(untypedEvents))
	for i, untypedEvent := range untypedEvents {
		burnEvents[i] = untypedEvent.(api.BurnEvent)
	}

	// Return the (potentially partial) events and the updated watermark along with the error.
	return burnEvents, newWatermark, err
}

// getSolanaBurnEventFromSig fetches and decodes burn events from a single Solana transaction signature.
func (o *Oracle) getSolanaBurnEventFromSig(ctx context.Context, sigStr string, programID string) (*api.BurnEvent, error) {
	program, err := solana.PublicKeyFromBase58(programID)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain program public key for burn event backfill: %w", err)
	}

	sig, err := solana.SignatureFromBase58(sigStr)
	if err != nil {
		return nil, fmt.Errorf("invalid signature string for backfill: %w", err)
	}

	if o.solanaClient == nil {
		return nil, fmt.Errorf("solana functionality is disabled")
	}

	txResult, err := o.solanaClient.GetTransaction(
		ctx,
		sig,
		&solrpc.GetTransactionOpts{
			Encoding:                       solana.EncodingBase64,
			Commitment:                     solrpc.CommitmentConfirmed,
			MaxSupportedTransactionVersion: &sidecartypes.MaxSupportedSolanaTxVersion,
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
				slog.Info("Backfilled Solana ROCK Burn Event",
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
	ctx context.Context,
	wg *sync.WaitGroup,
	update *oracleStateUpdate,
	updateMutex *sync.Mutex,
) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		backfillResp, err := o.zrChainQueryClient.ValidationQueryClient.BackfillRequests(ctx)
		if err != nil {
			// Don't push to errChan, as this is not a critical failure. Just log it.
			slog.Error("Failed to query backfill requests", "error", err)
			return
		}

		if backfillResp == nil || backfillResp.BackfillRequests == nil || len(backfillResp.BackfillRequests.Requests) == 0 {
			return // No backfill requests
		}
		slog.Info("Found backfill requests to process", "count", len(backfillResp.BackfillRequests.Requests))
		o.handleBackfillRequests(ctx, backfillResp.BackfillRequests.Requests, update, updateMutex)
	}()
}

// handleBackfillRequests processes a slice of backfill requests.
func (o *Oracle) handleBackfillRequests(ctx context.Context, requests []*validationtypes.MsgTriggerEventBackfill, update *oracleStateUpdate, updateMutex *sync.Mutex) {
	if len(requests) == 0 {
		return
	}

	var newBurnEvents []api.BurnEvent

	for i, req := range requests {
		// For now, only handle ZenTP burn events.
		if req.EventType == validationtypes.EventType_EVENT_TYPE_ZENTP_BURN {
			slog.Info("Processing zentp burn backfill request", "txHash", req.TxHash)
			programID := sidecartypes.SolRockProgramID[o.Config.Network]
			event, err := o.getSolanaBurnEventFromSig(ctx, req.TxHash, programID)
			if err != nil {
				slog.Error("Error processing backfill request", "txHash", req.TxHash, "error", err)
				continue
			}
			if event != nil {
				newBurnEvents = append(newBurnEvents, *event)
			}

			// Pause between requests to avoid rate-limiting, but not after the final one.
			if i < len(requests)-1 {
				timer := time.NewTimer(sidecartypes.SolanaFallbackSleepInterval)
				select {
				case <-timer.C:
				case <-ctx.Done():
					timer.Stop()
					return // Stop processing if the context is canceled
				}
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
		slog.Info("No watermark signature stored yet", "eventType", eventType)
		return solana.Signature{} // No signature stored yet, so zero value
	}
	sig, err := solana.SignatureFromBase58(sigStr)
	if err != nil {
		// Log the error but return a zero signature to proceed as if no prior sig exists
		slog.Warn("Could not parse stored signature string for event type. Treating as no prior signature.", "sigString", sigStr, "eventType", eventType, "error", err)
		return solana.Signature{}
	}

	slog.Info("Retrieved watermark signature", "eventType", eventType, "sig", sigStr)
	return sig
}

// formatWatermarkForLogging returns a user-friendly representation of a watermark signature
func formatWatermarkForLogging(sig solana.Signature) string {
	if sig.IsZero() {
		return "<none>"
	}
	return sig.String()
}

// processTransactionFunc defines the function signature for processing a single Solana transaction.
// It returns a slice of events (as any), and an error if processing fails.
type processTransactionFunc func(
	txResult *solrpc.GetTransactionResult,
	program solana.PublicKey,
	sig solana.Signature,
	debugMode bool,
) ([]any, error)

// getSolanaEvents is an optimized generic helper to fetch signatures for a given program, detect and heal
// gaps using the watermark (lastKnownSig), then download and process each transaction using the
// provided `processTransaction` callback. If any part of the transaction processing pipeline fails,
// it returns any partially processed events along with the watermark of the last successfully
// processed transaction. This allows the oracle to make incremental progress.

// PERFORMANCE OPTIMIZATIONS:
// - Rate limiting with semaphore to prevent RPC overload
// - Object pooling for batch requests and event processors
// - Transaction caching with TTL
// - Parallel batch processing with improved error handling
// - Exponential backoff retry strategy
// - Memory-efficient slice pre-allocation
// processSignatures takes a list of transaction signatures and processes them.
func (o *Oracle) processSignatures(
	ctx context.Context,
	signatures []*solrpc.TransactionSignature,
	program solana.PublicKey,
	eventTypeName string,
	processTransaction processTransactionFunc,
) ([]any, solana.Signature, error) {
	// Get event processor from pool for memory efficiency
	ep := o.eventProcessorPool.Get().(*EventProcessor)
	defer func() {
		ep.Reset()
		o.eventProcessorPool.Put(ep)
	}()

	var lastSuccessfullyProcessedSig solana.Signature
	// Adaptive batching parameters
	currentBatchSize := sidecartypes.SolanaEventFetchBatchSize
	minBatchSize := sidecartypes.SolanaEventFetchMinBatchSize

	// Track processing results for debugging
	totalNilResults := 0
	successfulTransactions := 0
	notFoundTransactions := 0
	processingErrors := 0
	emptyTransactions := 0

	// Pre-allocate with estimated capacity to reduce allocations
	estimatedEvents := len(signatures) * 2 // Estimate 2 events per signature on average
	if cap(ep.events) < estimatedEvents {
		ep.events = make([]any, 0, estimatedEvents)
	}

	// Process signatures with adaptive batching
	for i := 0; i < len(signatures); {
		if ctx.Err() != nil {
			return ep.GetEvents(), lastSuccessfullyProcessedSig, ctx.Err()
		}

		end := min(i+currentBatchSize, len(signatures))
		currentBatch := signatures[i:end]

		// Get batch request slice from pool
		batchRequests := o.batchRequestPool.Get().(jsonrpc.RPCRequests)
		batchRequests = (batchRequests)[:0] // Reset slice but keep capacity

		// Build batch requests
		for j, sigInfo := range currentBatch {
			batchRequests = append(batchRequests, &jsonrpc.RPCRequest{
				Method: "getTransaction",
				Params: []any{
					sigInfo.Signature.String(),
					map[string]any{
						"encoding":                       solana.EncodingBase64,
						"commitment":                     solrpc.CommitmentConfirmed,
						"maxSupportedTransactionVersion": &sidecartypes.MaxSupportedSolanaTxVersion,
					},
				},
				ID:      uint64(j),
				JSONRPC: "2.0",
			})
		}

		// Execute batch request
		batchCtx, batchCancel := context.WithTimeout(ctx, sidecartypes.SolanaBatchTimeout)
		batchResponses, batchErr := o.rpcCallBatchFn(batchCtx, batchRequests)
		batchCancel()

		// Debug logging for batch response
		slog.Debug("Batch RPC call completed",
			"eventType", eventTypeName,
			"batchSize", len(currentBatch),
			"requestCount", len(batchRequests),
			"responseCount", len(batchResponses),
			"batchError", batchErr)

		// Check for errors in the response itself
		if batchErr == nil {
			errorCount := 0
			for i, resp := range batchResponses {
				if resp.Error != nil {
					errorCount++
					slog.Debug("Individual response error",
						"eventType", eventTypeName,
						"responseIndex", i,
						"responseID", resp.ID,
						"error", resp.Error)
					batchErr = fmt.Errorf("response contains errors: %v", resp.Error)
					break
				}
			}
			if errorCount > 0 {
				slog.Debug("Batch response error summary",
					"eventType", eventTypeName,
					"errorCount", errorCount,
					"totalResponses", len(batchResponses))
			}
		}

		// If the batch failed, reduce the batch size and retry the same segment.
		if batchErr != nil {
			newBatchSize := max(currentBatchSize/2, minBatchSize)
			if currentBatchSize > minBatchSize {
				slog.Warn("Batch GetTransaction failed, reducing batch size and retrying",
					"eventType", eventTypeName, "error", batchErr, "oldSize", currentBatchSize, "newSize", newBatchSize)
				currentBatchSize = newBatchSize
			} else {
				// If we're already at the minimum batch size, switch to fallback for this batch
				slog.Warn("Batch transaction fetch failed at minimum batch size, switching to individual fallback", "eventType", eventTypeName, "size", len(currentBatch))
				if err := o.processFallbackTransactionsWithCaching(ctx, currentBatch, program, ep, &lastSuccessfullyProcessedSig, eventTypeName, processTransaction); err != nil {
					slog.Warn("Fallback processing failed - stopping to prevent watermark gaps", "eventType", eventTypeName, "processedCount", len(ep.GetEvents()), "newWatermark", lastSuccessfullyProcessedSig, "reason", err)
					o.batchRequestPool.Put(batchRequests)
					return ep.GetEvents(), lastSuccessfullyProcessedSig, err
				}
				i += len(currentBatch) // Advance past the successfully processed fallback segment
			}
			time.Sleep(sidecartypes.SolanaEventFetchRetrySleep) // Pause before retrying
			o.batchRequestPool.Put(batchRequests)
			continue // Retry the same segment `i`
		}

		// Success: Process the batch responses
		responseMap := make(map[int]*jsonrpc.RPCResponse, len(batchResponses))
		for _, resp := range batchResponses {
			if idx, ok := parseRPCResponseID(resp, eventTypeName); ok {
				responseMap[idx] = resp
			}
		}

		// Process responses in order to maintain watermark correctness
		for idx, sigInfo := range currentBatch {
			resp, exists := responseMap[idx]
			if !exists {
				err := fmt.Errorf("missing batch response for index %d", idx)
				slog.Warn("Incomplete batch response - stopping to prevent watermark gaps", "eventType", eventTypeName, "reason", err, "consecutiveSuccesses", idx)
				o.batchRequestPool.Put(batchRequests)
				return ep.GetEvents(), lastSuccessfullyProcessedSig, err
			}

			// Handle nil results from RPC (transaction not found/retrievable)
			if resp.Result == nil {
				totalNilResults++
				slog.Warn("Transaction returned nil result, attempting individual retry",
					"eventType", eventTypeName,
					"signature", sigInfo.Signature,
					"responseID", resp.ID,
					"responseError", resp.Error,
					"nilResultCount", totalNilResults)

				// Retry individual transaction
				if retryResult, err := o.retryIndividualTransaction(ctx, sigInfo.Signature, eventTypeName); err != nil {
					// ANY failure means we must stop to prevent gaps in watermark
					processingErrors++
					slog.Warn("Individual transaction retry failed - stopping to prevent watermark gaps",
						"eventType", eventTypeName,
						"signature", sigInfo.Signature,
						"error", err,
						"consecutiveSuccesses", idx,
						"strategy", "watermark_only_advances_through_consecutive_successes",
						"nextCycle", "will_retry_from_this_signature")
					o.batchRequestPool.Put(batchRequests)
					return ep.GetEvents(), lastSuccessfullyProcessedSig, err
				} else if retryResult != nil {
					// Process the successfully retried transaction
					events, err := processTransaction(retryResult, program, sigInfo.Signature, o.DebugMode)
					if err != nil {
						processingErrors++
						slog.Warn("Failed to process retried transaction - stopping to prevent watermark gaps",
							"eventType", eventTypeName,
							"signature", sigInfo.Signature,
							"error", err,
							"consecutiveSuccesses", idx)
						o.batchRequestPool.Put(batchRequests)
						return ep.GetEvents(), lastSuccessfullyProcessedSig, err
					}

					if len(events) > 0 {
						successfulTransactions++
						slog.Debug("Successfully processed retried transaction",
							"eventType", eventTypeName,
							"signature", sigInfo.Signature,
							"eventCount", len(events))
						for _, event := range events {
							ep.AddEvent(event)
						}
					} else {
						emptyTransactions++
						slog.Debug("Retried transaction processed but contained no events",
							"eventType", eventTypeName,
							"signature", sigInfo.Signature)
					}
					lastSuccessfullyProcessedSig = sigInfo.Signature
					continue
				} else {
					// Still nil after retry - stop to prevent gaps
					slog.Warn("Transaction still nil after retry - stopping to prevent watermark gaps",
						"eventType", eventTypeName,
						"signature", sigInfo.Signature,
						"consecutiveSuccesses", idx,
						"strategy", "zero_gaps_policy",
						"nextCycle", "will_retry_from_this_signature")
					o.batchRequestPool.Put(batchRequests)
					return ep.GetEvents(), lastSuccessfullyProcessedSig, fmt.Errorf("transaction still nil after retry")
				}
			}

			// Debug logging for response inspection
			if len(resp.Result) == 0 {
				slog.Warn("Response result is empty", "eventType", eventTypeName, "signature", sigInfo.Signature, "responseID", resp.ID)
			} else if len(resp.Result) < 10 {
				slog.Warn("Response result is very short", "eventType", eventTypeName, "signature", sigInfo.Signature, "rawResult", string(resp.Result), "resultLength", len(resp.Result))
			}

			var txRes solrpc.GetTransactionResult
			if err := json.Unmarshal(resp.Result, &txRes); err != nil {
				err = fmt.Errorf("unmarshal error: %w", err)
				resultStr := string(resp.Result)
				if len(resultStr) > 500 {
					resultStr = resultStr[:500] + "...(truncated)"
				}
				slog.Warn("Unmarshal error - stopping to prevent watermark gaps",
					"eventType", eventTypeName,
					"reason", err,
					"signature", sigInfo.Signature,
					"rawResult", resultStr,
					"resultLength", len(resp.Result),
					"responseID", resp.ID,
					"responseError", resp.Error,
					"consecutiveSuccesses", idx)
				o.batchRequestPool.Put(batchRequests)
				return ep.GetEvents(), lastSuccessfullyProcessedSig, err
			}

			o.cacheTransactionResult(sigInfo.Signature.String(), &txRes)

			events, err := processTransaction(&txRes, program, sigInfo.Signature, o.DebugMode)
			if err != nil {
				slog.Warn("Processing error - stopping to prevent watermark gaps",
					"eventType", eventTypeName,
					"reason", err,
					"signature", sigInfo.Signature,
					"consecutiveSuccesses", idx)
				o.batchRequestPool.Put(batchRequests)
				return ep.GetEvents(), lastSuccessfullyProcessedSig, err
			}

			if len(events) > 0 {
				successfulTransactions++
				slog.Debug("Successfully processed transaction",
					"eventType", eventTypeName,
					"signature", sigInfo.Signature,
					"eventCount", len(events))
				for _, event := range events {
					ep.AddEvent(event)
				}
			} else {
				emptyTransactions++
				slog.Debug("Transaction processed but contained no events",
					"eventType", eventTypeName,
					"signature", sigInfo.Signature)
			}
			lastSuccessfullyProcessedSig = sigInfo.Signature
		}

		// Advance to the next segment
		i += len(currentBatch)
		o.batchRequestPool.Put(batchRequests) // Return the batch request slice to the pool

		// Optional: slowly increase batch size on success
		if currentBatchSize < sidecartypes.SolanaEventFetchBatchSize {
			currentBatchSize = min(currentBatchSize+minBatchSize, sidecartypes.SolanaEventFetchBatchSize)
		}
		time.Sleep(sidecartypes.SolanaSleepInterval)
	}

	// Calculate processing statistics
	totalProcessed := len(signatures)
	successRate := float64(successfulTransactions) / float64(totalProcessed) * 100

	// Summary log with comprehensive batch processing results
	slog.Info("Batch processing summary",
		"eventType", eventTypeName,
		"totalSignatures", totalProcessed,
		"successfulTransactions", successfulTransactions,
		"emptyTransactions", emptyTransactions,
		"nilResults", totalNilResults,
		"notFoundTransactions", notFoundTransactions,
		"processingErrors", processingErrors,
		"successRate", fmt.Sprintf("%.1f%%", successRate),
		"extractedEvents", len(ep.GetEvents()),
		"newWatermark", lastSuccessfullyProcessedSig,
	)

	// Explain gap prevention strategy
	if processingErrors > 0 || totalNilResults > 0 {
		slog.Info("Gap prevention strategy: Watermark only advances through consecutive successes",
			"eventType", eventTypeName,
			"strategy", "stop_on_first_failure",
			"reason", "prevents_permanent_event_loss",
			"nextCycleWillRetry", "failed_transactions")
	}

	if totalNilResults > 0 {
		slog.Warn("Encountered nil results from Solana RPC", "eventType", eventTypeName, "nilResultCount", totalNilResults, "totalProcessed", len(signatures), "successfulEvents", len(ep.GetEvents()))
	}
	if len(ep.GetEvents()) > 0 {
		slog.Debug("Successfully extracted events from Solana transactions", "eventType", eventTypeName, "extractedEvents", len(ep.GetEvents()))
	}
	return ep.GetEvents(), lastSuccessfullyProcessedSig, nil
}

// getSolanaEvents is an optimized generic helper to fetch signatures for a given program, detect and heal
// gaps using the watermark (lastKnownSig), then download and process each transaction using the
// provided `processTransaction` callback. If any part of the transaction processing pipeline fails,
// it returns any partially processed events along with the watermark of the last successfully
// processed transaction. This allows the oracle to make incremental progress.

// PERFORMANCE OPTIMIZATIONS:
// - Rate limiting with semaphore to prevent RPC overload
// - Object pooling for batch requests and event processors
// - Transaction caching with TTL
// - Parallel batch processing with improved error handling
// - Exponential backoff retry strategy
// - Memory-efficient slice pre-allocation
// retryIndividualTransaction attempts to fetch a single transaction that returned nil in batch processing
func (o *Oracle) retryIndividualTransaction(ctx context.Context, sig solana.Signature, eventTypeName string) (*solrpc.GetTransactionResult, error) {
	// Check cache first
	if cached, found := o.getCachedTransactionResult(sig.String()); found {
		slog.Debug("Individual transaction retry found in cache", "eventType", eventTypeName, "signature", sig)
		return cached, nil
	}

	// Rate limiting: acquire semaphore slot
	select {
	case o.solanaRateLimiter <- struct{}{}:
		defer func() { <-o.solanaRateLimiter }()
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(sidecartypes.SolanaRateLimiterTimeout):
		return nil, fmt.Errorf("rate limiter timeout for individual retry after %v", sidecartypes.SolanaRateLimiterTimeout)
	}

	// Create timeout context for the individual request
	retryCtx, cancel := context.WithTimeout(ctx, sidecartypes.SolanaRPCTimeout)
	defer cancel()

	// Make individual RPC call with detailed error handling
	txResult, err := o.getTransactionFn(retryCtx, sig, &solrpc.GetTransactionOpts{
		Encoding:                       solana.EncodingBase64,
		Commitment:                     solrpc.CommitmentConfirmed,
		MaxSupportedTransactionVersion: &sidecartypes.MaxSupportedSolanaTxVersion,
	})

	if err != nil {
		// Categorize the error for better handling
		errStr := err.Error()
		if strings.Contains(errStr, "not found") {
			slog.Debug("Individual transaction retry confirmed not found", "eventType", eventTypeName, "signature", sig, "error", err)
			return nil, fmt.Errorf("not found: %w", err)
		} else if strings.Contains(errStr, "context canceled") || strings.Contains(errStr, "timeout") {
			slog.Debug("Individual transaction retry timed out", "eventType", eventTypeName, "signature", sig, "error", err)
			return nil, fmt.Errorf("timeout/canceled: %w", err)
		} else {
			slog.Debug("Individual transaction retry failed with RPC error", "eventType", eventTypeName, "signature", sig, "error", err)
			return nil, fmt.Errorf("rpc error: %w", err)
		}
	}

	// Handle nil result (transaction exists but no data returned)
	if txResult == nil {
		slog.Debug("Individual transaction retry returned nil result", "eventType", eventTypeName, "signature", sig)
		return nil, nil
	}

	// Cache the successful result
	o.cacheTransactionResult(sig.String(), txResult)
	slog.Debug("Individual transaction retry successful", "eventType", eventTypeName, "signature", sig)

	return txResult, nil
}

func (o *Oracle) getSolanaEvents(
	ctx context.Context,
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

	// Rate limiting: acquire semaphore slot
	select {
	case o.solanaRateLimiter <- struct{}{}:
		defer func() { <-o.solanaRateLimiter }()
	case <-ctx.Done():
		return nil, lastKnownSig, ctx.Err()
	case <-time.After(sidecartypes.SolanaRateLimiterTimeout):
		return nil, lastKnownSig, fmt.Errorf("rate limiter timeout for %s after %v", eventTypeName, sidecartypes.SolanaRateLimiterTimeout)
	}

	// Fetch initial signatures with timeout context
	fetchCtx, cancel := context.WithTimeout(ctx, sidecartypes.SolanaRPCTimeout)
	defer cancel()

	initialSignatures, err := o.getSignaturesForAddressFn(fetchCtx, program, &solrpc.GetSignaturesForAddressOpts{
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

	newSignatures, err := o.fetchAndFillSignatureGap(ctx, program, lastKnownSig, initialSignatures, limit, eventTypeName)
	if err != nil {
		return nil, lastKnownSig, fmt.Errorf("failed to fill signature gap, aborting to retry next cycle: %w", err)
	}
	if len(newSignatures) == 0 {
		slog.Info("No new signatures found", "eventType", eventTypeName, "watermark", formatWatermarkForLogging(lastKnownSig))
		return []any{}, newestSigFromNode, nil
	}

	slog.Info("Found new signatures", "eventType", eventTypeName, "count", len(newSignatures), "watermark", formatWatermarkForLogging(lastKnownSig), "newest", newestSigFromNode)

	events, lastSig, err := o.processSignatures(ctx, newSignatures, program, eventTypeName, processTransaction)
	if err != nil {
		return events, lastSig, err
	}

	return events, newestSigFromNode, nil
}

// cacheTransactionResult caches a transaction result with TTL
func (o *Oracle) cacheTransactionResult(sigStr string, txRes *solrpc.GetTransactionResult) {
	o.transactionCacheMutex.Lock()
	defer o.transactionCacheMutex.Unlock()

	// Clean expired entries (simple cleanup)
	now := time.Now()
	for key, cached := range o.transactionCache {
		if now.After(cached.ExpiresAt) {
			delete(o.transactionCache, key)
		}
	}

	// Cache the new result with 5-minute TTL
	o.transactionCache[sigStr] = &CachedTxResult{
		Result:    txRes,
		ExpiresAt: now.Add(5 * time.Minute),
	}
}

// getCachedTransactionResult retrieves a cached transaction result if available and not expired
func (o *Oracle) getCachedTransactionResult(sigStr string) (*solrpc.GetTransactionResult, bool) {
	o.transactionCacheMutex.RLock()
	defer o.transactionCacheMutex.RUnlock()

	cached, exists := o.transactionCache[sigStr]
	if !exists || time.Now().After(cached.ExpiresAt) {
		return nil, false
	}

	return cached.Result, true
}

// processFallbackTransactionsWithCaching handles individual transaction fetching with caching when batch processing fails
func (o *Oracle) processFallbackTransactionsWithCaching(
	ctx context.Context,
	currentBatch []*solrpc.TransactionSignature,
	program solana.PublicKey,
	ep *EventProcessor,
	lastSuccessfullyProcessedSig *solana.Signature,
	eventTypeName string,
	processTransaction processTransactionFunc,
) error {
	for _, sigInfo := range currentBatch {
		sigStr := sigInfo.Signature.String()

		// Check cache first
		var txRes *solrpc.GetTransactionResult
		if cached, found := o.getCachedTransactionResult(sigStr); found {
			txRes = cached
		} else {
			// Fetch from network with retry
			var txErr error
			retryDelay := sidecartypes.SolanaEventFetchRetrySleep

			for retry := 0; retry < sidecartypes.SolanaFallbackMaxRetries; retry++ {
				fallbackCtx, fallbackCancel := context.WithTimeout(ctx, sidecartypes.SolanaRPCTimeout)
				txRes, txErr = o.getTransactionFn(fallbackCtx, sigInfo.Signature, &solrpc.GetTransactionOpts{
					Encoding:                       solana.EncodingBase64,
					Commitment:                     solrpc.CommitmentConfirmed,
					MaxSupportedTransactionVersion: &sidecartypes.MaxSupportedSolanaTxVersion,
				})
				fallbackCancel()

				if txErr == nil && txRes != nil {
					// Cache the successful result
					o.cacheTransactionResult(sigStr, txRes)
					break
				}
				if retry < sidecartypes.SolanaFallbackMaxRetries-1 {
					time.Sleep(retryDelay)
					retryDelay = min(retryDelay*2, time.Second) // Exponential backoff for fallback
				}
			}
			if txErr != nil || txRes == nil {
				err := fmt.Errorf("unrecoverable tx fetch error: %w", txErr)
				slog.Error("Failed to fetch transaction after exhausting all retry attempts", "eventType", eventTypeName, "tx", sigInfo.Signature, "reason", err)
				return err
			}
		}

		events, err := processTransaction(txRes, program, sigInfo.Signature, o.DebugMode)
		if err != nil {
			slog.Error("Unrecoverable processing error in fallback", "eventType", eventTypeName, "tx", sigInfo.Signature, "error", err)
			return err
		}

		if len(events) > 0 {
			for _, event := range events {
				ep.AddEvent(event)
			}
		}
		*lastSuccessfullyProcessedSig = sigInfo.Signature
		time.Sleep(sidecartypes.SolanaFallbackSleepInterval)
	}
	return nil
}

// fetchAndFillSignatureGap back-pages the Solana signature list until the provided watermark is
// found or `SolanaMaxBackfillPages` is exceeded.
func (o *Oracle) fetchAndFillSignatureGap(
	ctx context.Context,
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
		pageSigs, err := o.getSignaturesForAddressFn(ctx, program, &solrpc.GetSignaturesForAddressOpts{
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
	slog.Error("Unable to locate starting watermark signature after scanning maximum pages. Proceeding with collected data", "eventType", eventTypeName)
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
			slog.Info("Mint event", "eventType", eventTypeName, "tx", sig)
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
	// Use mock function if available (for testing)
	if o.reconcileBurnEventsFn != nil {
		return o.reconcileBurnEventsFn(ctx, eventsToClean, cleanedEvents, chainTypeName)
	}

	remaining := make([]api.BurnEvent, 0, len(eventsToClean))
	updated := make(map[string]bool)
	maps.Copy(updated, cleanedEvents)

	var zenbtcQueryErrors, zentpQueryErrors, bech32EncodingErrors int
	var lastZenbtcError, lastZentpError, lastBech32Error error

	for _, ev := range eventsToClean {
		key := fmt.Sprintf("%s-%s-%d", ev.ChainID, ev.TxID, ev.LogIndex)
		if updated[key] {
			continue
		}
		found := false
		zenbtcResp, err := o.zrChainQueryClient.ZenBTCQueryClient.BurnEvents(ctx, 0, ev.TxID, ev.LogIndex, ev.ChainID)
		if err != nil {
			zenbtcQueryErrors++
			lastZenbtcError = err
		} else if zenbtcResp != nil && len(zenbtcResp.BurnEvents) > 0 {
			found = true
		}

		if !found && chainTypeName == "Solana" {
			if len(ev.DestinationAddr) >= 20 {
				bech32Addr, err := sdkBech32.ConvertAndEncode("zen", ev.DestinationAddr[:20])
				if err != nil {
					bech32EncodingErrors++
					lastBech32Error = err
				} else {
					ztp, err := o.zrChainQueryClient.ZenTPQueryClient.Burns(ctx, bech32Addr, ev.TxID)
					if err != nil {
						zentpQueryErrors++
						lastZentpError = err
					} else if ztp != nil && len(ztp.Burns) > 0 {
						found = true
					}
				}
			}
		}
		if found {
			updated[key] = true
		} else {
			remaining = append(remaining, ev)
		}
	}

	// Log summary of any query errors
	if zenbtcQueryErrors > 0 {
		slog.Error("Failed to query zrChain for zenBTC burn events",
			"failedCount", zenbtcQueryErrors,
			"totalEvents", len(eventsToClean),
			"chainType", chainTypeName,
			"lastError", lastZenbtcError)
	}

	if zentpQueryErrors > 0 {
		slog.Error("Failed to query zrChain for ZenTP burn events",
			"failedCount", zentpQueryErrors,
			"totalEvents", len(eventsToClean),
			"chainType", chainTypeName,
			"lastError", lastZentpError)
	}

	if bech32EncodingErrors > 0 {
		slog.Error("Failed to encode destination addresses for ZenTP burn checks",
			"failedCount", bech32EncodingErrors,
			"totalEvents", len(eventsToClean),
			"chainType", chainTypeName,
			"lastError", lastBech32Error)
	}

	return remaining, updated
}
