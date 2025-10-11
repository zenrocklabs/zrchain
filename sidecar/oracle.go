package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"maps"
	"math/big"
	"net/http"
	"os"
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

	validationkeeper "github.com/Zenrock-Foundation/zrchain/v6/x/validation/keeper"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	zenbtctypes "github.com/Zenrock-Foundation/zrchain/v6/x/zenbtc/types"
	zentptypes "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	zenbtc "github.com/Zenrock-Foundation/zrchain/v6/zenbtc/bindings"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	solana "github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	jsonrpc "github.com/gagliardetto/solana-go/rpc/jsonrpc"
	middleware "github.com/zenrocklabs/zenrock-avs/contracts/bindings/ZrServiceManager"
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
	forceTestReset bool,
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
		ForceTestReset:     forceTestReset,

		// Initialize performance optimization fields
		solanaRateLimiter: make(chan struct{}, sidecartypes.SolanaMaxConcurrentRPCCalls), // Configurable concurrent Solana RPC calls
		transactionCache:  make(map[string]*CachedTxResult),
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
			o.lastSolZenZECMintSigStr = latestDiskState.LastSolZenZECMintSig
			o.lastSolZenZECBurnSigStr = latestDiskState.LastSolZenZECBurnSig
			o.lastSolRockBurnSigStr = latestDiskState.LastSolRockBurnSig
			slog.Info("Loaded state from file",
				"rockMintSig", o.lastSolRockMintSigStr,
				"zenBTCMintSig", o.lastSolZenBTCMintSigStr,
				"zenBTCBurnSig", o.lastSolZenBTCBurnSigStr,
				"zenZECMintSig", o.lastSolZenZECMintSigStr,
				"zenZECBurnSig", o.lastSolZenZECBurnSigStr,
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
	o.getSolanaZenBTCBurnEventsFn = func(ctx context.Context, programID string, lastKnownSig solana.Signature) ([]api.BurnEvent, solana.Signature, error) {
		// This function is only used for backward compatibility in places where we don't have access to state update
		// Failed transactions will be lost here, but this is only used in non-critical paths
		dummyUpdate := &oracleStateUpdate{pendingTransactions: make(map[string]sidecartypes.PendingTxInfo)}
		dummyMutex := &sync.Mutex{}
		events, sig, err := o.getSolanaZenBTCBurnEvents(ctx, programID, lastKnownSig, dummyUpdate, dummyMutex)
		if len(dummyUpdate.pendingTransactions) > 0 {
			slog.Warn("Lost failed transactions in backward compatibility function", "count", len(dummyUpdate.pendingTransactions), "eventType", "Solana zenBTC burn")
		}
		return events, sig, err
	}
	o.getSolanaZenZECBurnEventsFn = func(ctx context.Context, programID string, lastKnownSig solana.Signature) ([]api.BurnEvent, solana.Signature, error) {
		dummyUpdate := &oracleStateUpdate{pendingTransactions: make(map[string]sidecartypes.PendingTxInfo)}
		dummyMutex := &sync.Mutex{}
		events, sig, err := o.getSolanaZenZECBurnEvents(ctx, programID, lastKnownSig, dummyUpdate, dummyMutex)
		if len(dummyUpdate.pendingTransactions) > 0 {
			slog.Warn("Lost failed transactions in backward compatibility function", "count", len(dummyUpdate.pendingTransactions), "eventType", "Solana zenZEC burn")
		}
		return events, sig, err
	}
	o.getSolanaRockBurnEventsFn = func(ctx context.Context, programID string, lastKnownSig solana.Signature) ([]api.BurnEvent, solana.Signature, error) {
		// This function is only used for backward compatibility in places where we don't have access to state update
		// Failed transactions will be lost here, but this is only used in non-critical paths
		dummyUpdate := &oracleStateUpdate{pendingTransactions: make(map[string]sidecartypes.PendingTxInfo)}
		dummyMutex := &sync.Mutex{}
		events, sig, err := o.getSolanaRockBurnEvents(ctx, programID, lastKnownSig, dummyUpdate, dummyMutex)
		if len(dummyUpdate.pendingTransactions) > 0 {
			slog.Warn("Lost failed transactions in backward compatibility function", "count", len(dummyUpdate.pendingTransactions), "eventType", "Solana ROCK burn")
		}
		return events, sig, err
	}

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
	// Initialize periodic reset scheduling (interval computed dynamically; may be overridden later by test flag)
	o.scheduleNextReset(time.Now().UTC(), time.Duration(sidecartypes.OracleStateResetIntervalHours)*time.Hour)

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
	zenBTCController, err := zenbtc.NewZenBTController(
		common.HexToAddress(sidecartypes.ZenBTCControllerAddresses[o.Config.Network]),
		o.EthClient,
	)
	if err != nil {
		return fmt.Errorf("failed to create contract instance: %w", err)
	}
	mainnetEthClient, btcPriceFeed, ethPriceFeed := o.initPriceFeed()

	// Allow customization of ticker interval in regnet network
	mainLoopTickerIntervalDuration := func() time.Duration {
		if o.Config.Network == "regnet" {
			configValue := time.Duration(o.Config.E2ETestsTickerInterval) * time.Second
			configValue = max(configValue, 15*time.Second)
			return configValue
		}
		return sidecartypes.MainLoopTickerInterval
	}()

	var tickCancel context.CancelFunc = func() {}
	defer tickCancel()

	// Align the start time to the nearest MainLoopTickerInterval.
	if !o.SkipInitialWait {
		ntpTime, err := ntp.Time(sidecartypes.NTPServer)
		if err != nil {
			slog.Error("Failed to fetch NTP time at startup. Cannot proceed.", "error", err)
			panic(fmt.Sprintf("FATAL: Failed to fetch NTP time at startup: %v. Cannot proceed.", err))
		}
		alignedStart := ntpTime.Truncate(mainLoopTickerIntervalDuration).Add(mainLoopTickerIntervalDuration)
		initialSleep := time.Until(alignedStart)
		if initialSleep > 0 {
			slog.Info("Initial alignment: Sleeping until start ticker.",
				"sleepDuration", initialSleep.Round(time.Millisecond),
				"alignedStart", alignedStart.Format(sidecartypes.TimeFormatPrecise))
			time.Sleep(initialSleep)
		}
	} else {
		slog.Info("Skipping initial alignment wait due to --skip-initial-wait flag. Firing initial tick immediately.")
		var initialTickCtx context.Context
		initialTickCtx, tickCancel = context.WithCancel(ctx)
		go o.processOracleTick(initialTickCtx, serviceManager, zenBTCController, btcPriceFeed, ethPriceFeed, mainnetEthClient, time.Now())
	}

	mainLoopTicker := time.NewTicker(mainLoopTickerIntervalDuration)
	defer mainLoopTicker.Stop()
	o.mainLoopTicker = mainLoopTicker
	slog.Info("Ticker synced, awaiting initial oracle data fetch", "interval", mainLoopTickerIntervalDuration)

	// Start persistent pending transaction processor that runs across all ticks
	if o.solanaClient != nil {
		go o.processPendingTransactionsPersistent(ctx)
	}

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
			go o.processOracleTick(tickCtx, serviceManager, zenBTCController, btcPriceFeed, ethPriceFeed, mainnetEthClient, tickTime)
		}
	}
}

func (o *Oracle) processOracleTick(
	tickCtx context.Context,
	serviceManager *middleware.ContractZrServiceManager,
	zenBTCController *zenbtc.ZenBTController,
	btcPriceFeed *aggregatorv3.AggregatorV3Interface,
	ethPriceFeed *aggregatorv3.AggregatorV3Interface,
	mainnetEthClient *ethclient.Client,
	tickTime time.Time,
) {
	// Perform scheduled reset if due (evaluated using the tick's aligned time)
	o.maybePerformScheduledReset(tickTime.UTC())
	newState, err := o.fetchAndProcessState(tickCtx, serviceManager, zenBTCController, btcPriceFeed, ethPriceFeed, mainnetEthClient)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			slog.Info("Data fetch time limit reached. Applying partially gathered state to meet tick deadline.", "tickTime", tickTime.Format(sidecartypes.TimeFormatPrecise))
		} else {
			slog.Error("Error fetching and processing state, applying partial update with fallbacks", "error", err)
		}
		// Continue to apply the partial state rather than aborting entirely
	}

	// Always apply the state update (even if partial) - the individual event fetching functions
	// have their own watermark protection to prevent event loss
	slog.Info("Applying state update for tick", "tickTime", tickTime.Format(sidecartypes.TimeFormatPrecise))
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
	oldZenZECMint := o.lastSolZenZECMintSigStr
	oldZenZECBurn := o.lastSolZenZECBurnSigStr
	oldRockBurn := o.lastSolRockBurnSigStr

	o.currentState.Store(&newState)

	// Log event counts in each state field every tick
	slog.Info("State event counts for this tick",
		"ethBurnEvents", len(newState.EthBurnEvents),
		"cleanedEthBurnEvents", len(newState.CleanedEthBurnEvents),
		"solanaBurnEvents", len(newState.SolanaBurnEvents),
		"cleanedSolanaBurnEvents", len(newState.CleanedSolanaBurnEvents),
		"solanaMintEvents", len(newState.SolanaMintEvents),
		"cleanedSolanaMintEvents", len(newState.CleanedSolanaMintEvents),
		"redemptions", len(newState.Redemptions),
		"pendingSolanaTxs", len(newState.PendingSolanaTxs))

	// Update the oracle's high-watermark fields from the newly applied state.
	// These are used as the starting point for the next fetch cycle.
	o.lastSolRockMintSigStr = newState.LastSolRockMintSig
	o.lastSolZenBTCMintSigStr = newState.LastSolZenBTCMintSig
	o.lastSolZenBTCBurnSigStr = newState.LastSolZenBTCBurnSig
	o.lastSolZenZECMintSigStr = newState.LastSolZenZECMintSig
	o.lastSolZenZECBurnSigStr = newState.LastSolZenZECBurnSig
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
	if oldZenZECMint != o.lastSolZenZECMintSigStr {
		slog.Info("Updated zenZEC mint watermark", "old", oldZenZECMint, "new", o.lastSolZenZECMintSigStr)
		watermarkChanged = true
	}
	if oldZenZECBurn != o.lastSolZenZECBurnSigStr {
		slog.Info("Updated zenZEC burn watermark", "old", oldZenZECBurn, "new", o.lastSolZenZECBurnSigStr)
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
	zenBTCController *zenbtc.ZenBTController,
	btcPriceFeed *aggregatorv3.AggregatorV3Interface,
	ethPriceFeed *aggregatorv3.AggregatorV3Interface,
	tempEthClient *ethclient.Client,
) (sidecartypes.OracleState, error) {
	var wg sync.WaitGroup

	// Log initial state at beginning of tick
	initialState := o.currentState.Load().(*sidecartypes.OracleState)
	slog.Info("TICK START STATE SNAPSHOT",
		"solanaMintEvents", len(initialState.SolanaMintEvents),
		"cleanedSolanaMintEvents", len(initialState.CleanedSolanaMintEvents),
		"solanaBurnEvents", len(initialState.SolanaBurnEvents),
		"cleanedSolanaBurnEvents", len(initialState.CleanedSolanaBurnEvents),
		"pendingSolanaTxs", len(initialState.PendingSolanaTxs))

	slog.Info("Retrieving latest header", "network", sidecartypes.NetworkNames[o.Config.Network], "time", time.Now().Format(sidecartypes.TimeFormatPrecise))
	latestHeader, err := o.EthClient.HeaderByNumber(tickCtx, nil)
	if err != nil {
		return sidecartypes.OracleState{}, fmt.Errorf("failed to fetch latest block: %w", err)
	}
	slog.Info("Retrieved latest header", "network", sidecartypes.NetworkNames[o.Config.Network], "block", latestHeader.Number.Uint64(), "time", time.Now().Format(sidecartypes.TimeFormatPrecise))
	targetBlockNumber := new(big.Int).Sub(latestHeader.Number, big.NewInt(sidecartypes.EthBlocksBeforeFinality))

	// Check base fee availability
	if latestHeader.BaseFee == nil {
		return sidecartypes.OracleState{}, fmt.Errorf("base fee not available (pre-London fork?)")
	}

	update := o.initializeStateUpdate()
	var updateMutex sync.Mutex
	errChan := make(chan error, sidecartypes.ErrorChannelBufferSize)

	// Use a separate context for the goroutines that can be canceled
	// if the main tick context is canceled.
	routinesCtx, cancelRoutines := context.WithCancel(tickCtx)
	defer cancelRoutines()

	// Pending transactions are now processed by a persistent background goroutine
	// Started in runOracleMainLoop, not per-tick

	// Fetch Ethereum contract data (AVS delegations and redemptions on EigenLayer)
	// o.fetchEthereumContractData(routinesCtx, &wg, serviceManager, zenBTCController, targetBlockNumber, update, &updateMutex, errChan)

	// Fetch network data (gas estimates, tips, Solana fees)
	o.fetchNetworkData(routinesCtx, &wg, update, &updateMutex, errChan)

	// Fetch price data (ROCK, BTC, ETH)
	o.fetchPriceData(routinesCtx, &wg, btcPriceFeed, ethPriceFeed, tempEthClient, update, &updateMutex, errChan)

	// Fetch zenBTC burn events from Ethereum
	o.fetchEthereumBurnEvents(routinesCtx, &wg, latestHeader, update, &updateMutex, errChan)

	// Fetch Solana mint events for zenBTC and ROCK (only if Solana is enabled)
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
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Panic recovered in waitgroup goroutine", "panic", r)
				close(waitChan)
			}
		}()
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

// sendError sends an error to the channel if the context is not done.
// This prevents panics from sending on a closed channel.
func sendError(ctx context.Context, errChan chan<- error, err error) {
	select {
	case <-ctx.Done():
		// Only log context cancellation errors that aren't actually about context cancellation
		if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) && !strings.Contains(err.Error(), "context canceled") {
			slog.Debug("Context canceled, dropping error", "err", err)
		}
	case errChan <- err:
	}
}

// handleTransactionFailure consolidates the repetitive failure handling pattern
func (o *Oracle) handleTransactionFailure(
	failedSignatures []string,
	signature solana.Signature,
	newestSigProcessed solana.Signature,
	eventTypeName string,
	message string,
	err error,
) ([]string, solana.Signature) {
	if err != nil {
		slog.Warn(message+", adding to failed signatures",
			"eventType", eventTypeName,
			"signature", signature,
			"error", err)
	} else {
		slog.Warn(message+", adding to failed signatures",
			"eventType", eventTypeName,
			"signature", signature)
	}

	failedSignatures = append(failedSignatures, signature.String())
	newestSigProcessed = signature
	return failedSignatures, newestSigProcessed
}

// fetchAndUpdateState is a generic helper function to reduce goroutine boilerplate for state updates
func fetchAndUpdateState[T any](
	ctx context.Context,
	wg *sync.WaitGroup,
	errChan chan<- error,
	updateMutex *sync.Mutex,
	fetchFunc func(ctx context.Context) (T, error),
	updateFunc func(result T, update *oracleStateUpdate),
	update *oracleStateUpdate,
	errorMsg string,
) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		result, err := fetchFunc(ctx)
		if err != nil {
			sendError(ctx, errChan, fmt.Errorf("%s: %w", errorMsg, err))
			return
		}
		updateMutex.Lock()
		updateFunc(result, update)
		updateMutex.Unlock()
	}()
}

func (o *Oracle) fetchEthereumContractData(
	ctx context.Context,
	wg *sync.WaitGroup,
	serviceManager *middleware.ContractZrServiceManager,
	zenBTCController *zenbtc.ZenBTController,
	targetBlockNumber *big.Int,
	update *oracleStateUpdate,
	updateMutex *sync.Mutex,
	errChan chan<- error,
) {
	// Fetches the state of AVS delegations from the service manager contract.
	fetchAndUpdateState(
		ctx, wg, errChan, updateMutex,
		func(ctx context.Context) (map[string]map[string]*big.Int, error) {
			return o.getServiceManagerState(ctx, serviceManager, targetBlockNumber)
		},
		func(result map[string]map[string]*big.Int, update *oracleStateUpdate) {
			update.eigenDelegations = result
		},
		update,
		"failed to get contract state",
	)

	// Fetches pending zenBTC redemptions from the zenBTC controller contract.
	fetchAndUpdateState(
		ctx, wg, errChan, updateMutex,
		func(ctx context.Context) ([]api.Redemption, error) {
			return o.getRedemptions(ctx, zenBTCController, targetBlockNumber)
		},
		func(result []api.Redemption, update *oracleStateUpdate) {
			update.redemptions = result
		},
		update,
		"failed to get zenBTC contract state",
	)
}

// applyGasBuffer applies the gas estimation buffer to a gas value
func applyGasBuffer(gas uint64) uint64 {
	return (gas * sidecartypes.GasEstimationBuffer) / 100
}

// handleGasEstimationFallback logs an error and returns a fallback gas estimate with buffer applied
func handleGasEstimationFallback(err error, operation string) uint64 {
	slog.Error(operation, "error", err)
	bufferedGas := applyGasBuffer(sidecartypes.WrapCallGasLimitFallback)
	slog.Info("Using fallback gas estimate", "gas", bufferedGas)
	return bufferedGas
}

func (o *Oracle) fetchNetworkData(
	ctx context.Context,
	wg *sync.WaitGroup,
	update *oracleStateUpdate,
	updateMutex *sync.Mutex,
	errChan chan<- error,
) {
	// Fetches the suggested gas tip cap (priority fee) for Ethereum transactions.
	fetchAndUpdateState(
		ctx, wg, errChan, updateMutex,
		func(ctx context.Context) (*big.Int, error) {
			return o.EthClient.SuggestGasTipCap(ctx)
		},
		func(result *big.Int, update *oracleStateUpdate) {
			update.suggestedTip = result
		},
		update,
		"failed to get suggested priority fee",
	)

	// Estimates the gas required for a zenBTC stake call on Ethereum.
	fetchAndUpdateState(
		ctx, wg, errChan, updateMutex,
		func(ctx context.Context) (uint64, error) {
			network := o.Config.Network
			whitelistedRoleAddr := sidecartypes.WhitelistedRoleAddresses[network]
			controllerAddr := sidecartypes.ZenBTCControllerAddresses[network]

			wrapCallData, err := validationkeeper.EncodeWrapCallData(
				common.HexToAddress("0x0000000000000000000000000000000000000000"),
				big.NewInt(1000000000),
				1,
			)
			if err != nil {
				return handleGasEstimationFallback(err, "Failed to encode wrap call data"), nil
			}

			addr := common.HexToAddress(controllerAddr)
			estimatedGas, err := o.EthClient.EstimateGas(context.Background(), ethereum.CallMsg{
				From: common.HexToAddress(whitelistedRoleAddr),
				To:   &addr,
				Data: wrapCallData,
			})
			if err != nil {
				return handleGasEstimationFallback(err, "Gas estimation failed"), nil
			}

			finalGas := applyGasBuffer(estimatedGas)
			return finalGas, nil
		},
		func(result uint64, update *oracleStateUpdate) {
			update.estimatedGas = result
		},
		update,
		"failed to estimate gas for zenBTC stake call",
	)
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
	fetchAndUpdateState(
		ctx, wg, errChan, updateMutex,
		func(ctx context.Context) (math.LegacyDec, error) {
			client := &http.Client{
				Timeout: httpTimeout,
			}
			resp, err := client.Get(sidecartypes.ROCKUSDPriceURL)
			if err != nil {
				return math.LegacyDec{}, fmt.Errorf("failed to retrieve ROCK price data: %w", err)
			}
			defer resp.Body.Close()

			var priceData []PriceData
			if err := json.NewDecoder(resp.Body).Decode(&priceData); err != nil || len(priceData) == 0 {
				return math.LegacyDec{}, fmt.Errorf("failed to decode ROCK price data or empty data: %w", err)
			}
			priceDec, err := math.LegacyNewDecFromStr(priceData[0].Last)
			if err != nil {
				return math.LegacyDec{}, fmt.Errorf("failed to parse ROCK price data: %w", err)
			}
			return priceDec, nil
		},
		func(result math.LegacyDec, update *oracleStateUpdate) {
			update.ROCKUSDPrice = result
		},
		update,
		"failed to fetch ROCK price data",
	)

	// Fetches the latest BTC/USD and ETH/USD prices from Chainlink price feeds on Ethereum mainnet.
	fetchAndUpdateState(
		ctx, wg, errChan, updateMutex,
		func(ctx context.Context) (cryptoPrices, error) {
			mainnetLatestHeader, err := tempEthClient.HeaderByNumber(ctx, nil)
			if err != nil {
				return cryptoPrices{}, fmt.Errorf("failed to fetch latest mainnet block: %w", err)
			}
			targetBlockNumberMainnet := new(big.Int).Sub(mainnetLatestHeader.Number, big.NewInt(sidecartypes.EthBlocksBeforeFinality))

			if btcPriceFeed == nil || ethPriceFeed == nil {
				return cryptoPrices{}, fmt.Errorf("BTC or ETH price feed not initialized")
			}

			btcPrice, err := o.fetchPrice(btcPriceFeed, targetBlockNumberMainnet)
			if err != nil {
				return cryptoPrices{}, fmt.Errorf("failed to fetch BTC price: %w", err)
			}

			ethPrice, err := o.fetchPrice(ethPriceFeed, targetBlockNumberMainnet)
			if err != nil {
				return cryptoPrices{}, fmt.Errorf("failed to fetch ETH price: %w", err)
			}

			return cryptoPrices{
				BTCUSDPrice: btcPrice,
				ETHUSDPrice: ethPrice,
			}, nil
		},
		func(result cryptoPrices, update *oracleStateUpdate) {
			update.BTCUSDPrice = result.BTCUSDPrice
			update.ETHUSDPrice = result.ETHUSDPrice
		},
		update,
		"failed to fetch crypto prices",
	)
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
	fetchAndUpdateState(
		ctx, wg, errChan, updateMutex,
		func(ctx context.Context) (burnEventResult, error) {
			currentState := o.currentState.Load().(*sidecartypes.OracleState)

			fromBlock := new(big.Int).Sub(latestHeader.Number, big.NewInt(int64(sidecartypes.EthBurnEventsBlockRange)))
			toBlock := latestHeader.Number
			newEvents, err := o.getEthBurnEvents(ctx, fromBlock, toBlock)
			if err != nil {
				// Don't return error, just log it and continue with empty events
				slog.Warn("Failed to get Ethereum burn events, proceeding with reconciliation only", "error", err)
				newEvents = []api.BurnEvent{} // Ensure slice is not nil
			}

			// Reconcile and merge
			remainingEvents, cleanedEvents := o.reconcileBurnEventsWithZRChain(ctx, currentState.EthBurnEvents, currentState.CleanedEthBurnEvents, "Ethereum")
			mergedEvents := mergeNewBurnEvents(remainingEvents, cleanedEvents, newEvents, "Ethereum")

			return burnEventResult{
				ethBurnEvents:        mergedEvents,
				cleanedEthBurnEvents: cleanedEvents,
			}, nil
		},
		func(result burnEventResult, update *oracleStateUpdate) {
			update.ethBurnEvents = result.ethBurnEvents
			update.cleanedEthBurnEvents = result.cleanedEthBurnEvents
		},
		update,
		"failed to process Ethereum burn events",
	)
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

		// Parallel fetch of ROCK, zenBTC, and zenZEC mint events
		var rockEvents, zenbtcEvents, zenzecEvents []api.SolanaMintEvent
		var newRockSig, newZenBTCSig, newZenZECSig solana.Signature
		var rockErr, zenbtcErr, zenzecErr error
		var mintWg sync.WaitGroup

		// Fetch ROCK mint events in parallel
		mintWg.Add(1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("Panic recovered in ROCK mint goroutine", "panic", r)
					rockErr = fmt.Errorf("panic in ROCK mint processing: %v", r)
				}
				mintWg.Done()
			}()
			lastKnownRockSig := o.GetLastProcessedSolSignature(sidecartypes.SolRockMint)
			rockEvents, newRockSig, rockErr = o.getSolROCKMints(ctx, sidecartypes.SolRockProgramID[o.Config.Network], lastKnownRockSig, update, updateMutex)
		}()

		// Fetch zenBTC mint events in parallel
		mintWg.Add(1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("Panic recovered in zenBTC mint goroutine", "panic", r)
					zenbtcErr = fmt.Errorf("panic in zenBTC mint processing: %v", r)
				}
				mintWg.Done()
			}()
			lastKnownZenBTCSig := o.GetLastProcessedSolSignature(sidecartypes.SolZenBTCMint)
			zenbtcEvents, newZenBTCSig, zenbtcErr = o.getSolZenBTCMints(ctx, sidecartypes.ZenBTCSolanaProgramID[o.Config.Network], lastKnownZenBTCSig, update, updateMutex)
		}()

		// Fetch zenZEC mint events in parallel
		mintWg.Add(1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("Panic recovered in zenZEC mint goroutine", "panic", r)
					zenzecErr = fmt.Errorf("panic in zenZEC mint processing: %v", r)
				}
				mintWg.Done()
			}()
			lastKnownZenZECSig := o.GetLastProcessedSolSignature(sidecartypes.SolZenZECMint)
			zenzecEvents, newZenZECSig, zenzecErr = o.getSolZenZECMints(ctx, sidecartypes.ZenZECSolanaProgramID[o.Config.Network], lastKnownZenZECSig, update, updateMutex)
		}()

		mintWg.Wait()

		// Handle errors after parallel execution
		if rockErr != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to get Solana ROCK mint events, applying partial results: %w", rockErr))
		}
		if zenbtcErr != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to get Solana zenBTC mint events, applying partial results: %w", zenbtcErr))
		}
		if zenzecErr != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to get Solana zenZEC mint events, applying partial results: %w", zenzecErr))
		}

		allNewEvents := append(rockEvents, zenbtcEvents...)
		allNewEvents = append(allNewEvents, zenzecEvents...)

		// Reconcile and merge
		slog.Info("Before reconciliation",
			"pendingMintEvents", len(currentState.SolanaMintEvents),
			"cleanedMintEventsMap", len(currentState.CleanedSolanaMintEvents),
			"newlyFetchedEvents", len(allNewEvents))
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
		slog.Info("MINT EVENT MERGE INPUT",
			"remainingEventsFromReconcile", len(remainingEvents),
			"cleanedEventsFromReconcile", len(cleanedEvents),
			"newEventsToMerge", len(allNewEvents),
			"reconciliationSucceeded", reconcileErr == nil)
		mergedMintEvents := mergeNewMintEvents(remainingEvents, cleanedEvents, allNewEvents, "Solana mint")
		slog.Info("MINT EVENT MERGE OUTPUT",
			"mergedMintEventsCount", len(mergedMintEvents))

		updateMutex.Lock()

		initialUpdateCount := len(update.SolanaMintEvents)

		// Check if there are any existing events from pending transactions
		existingPendingEvents := make(map[string]api.SolanaMintEvent)
		for _, event := range update.SolanaMintEvents {
			existingPendingEvents[event.TxSig] = event
		}

		// Merge new events with existing pending events, avoiding duplicates
		// Use hash map for O(1) lookup instead of O(n²) linear search
		newEventSigs := make(map[string]bool, len(mergedMintEvents))
		finalEvents := make([]api.SolanaMintEvent, 0, len(mergedMintEvents)+len(existingPendingEvents))

		// First add all new events and build lookup map
		for _, event := range mergedMintEvents {
			finalEvents = append(finalEvents, event)
			newEventSigs[event.TxSig] = true
		}

		// Then add any pending events that weren't already included (O(1) lookup)
		for sig, event := range existingPendingEvents {
			if !newEventSigs[sig] {
				finalEvents = append(finalEvents, event)
			}
		}

		slog.Info("FINAL UPDATE COMPOSITION",
			"mergedMintEventsCount", len(mergedMintEvents),
			"existingPendingEventsCount", len(existingPendingEvents),
			"finalEventsCount", len(finalEvents),
			"initialUpdateCount", initialUpdateCount)

		update.SolanaMintEvents = finalEvents
		update.cleanedSolanaMintEvents = cleanedEvents

		// Calculate deduplication statistics
		pendingEventsPreserved := len(finalEvents) - len(mergedMintEvents)
		duplicatesRemoved := len(existingPendingEvents) - pendingEventsPreserved

		slog.Info("Final merge and state update completed",
			"finalMintEvents", len(update.SolanaMintEvents),
			"finalCleanedEvents", len(update.cleanedSolanaMintEvents),
			"addedToUpdate", len(update.SolanaMintEvents)-initialUpdateCount,
			"newEvents", len(mergedMintEvents),
			"existingPendingEvents", len(existingPendingEvents),
			"pendingEventsPreserved", pendingEventsPreserved,
			"duplicatesRemoved", duplicatesRemoved,
			"reconcileErr", reconcileErr != nil)

		// Get current watermarks for comparison
		lastKnownRockSig := o.GetLastProcessedSolSignature(sidecartypes.SolRockMint)
		lastKnownZenBTCSig := o.GetLastProcessedSolSignature(sidecartypes.SolZenBTCMint)
		lastKnownZenZECSig := o.GetLastProcessedSolSignature(sidecartypes.SolZenZECMint)

		// Update watermarks if we have valid new signatures and they represent progress
		// Allow advancement even on timeout errors as long as failed transactions are in pending queue
		// Fatal errors (like pending store failures) are indicated by unchanged watermark signatures
		if !newRockSig.IsZero() && (!newRockSig.Equals(lastKnownRockSig) || rockErr == nil) {
			update.latestSolanaSigs[sidecartypes.SolRockMint] = newRockSig
			if rockErr != nil {
				slog.Info("Advancing ROCK mint watermark despite partial processing error (failed transactions safely stored in pending queue)",
					"newWatermark", newRockSig,
					"partialError", rockErr)
			}
		} else if rockErr != nil && !newRockSig.IsZero() && newRockSig.Equals(lastKnownRockSig) {
			slog.Warn("Blocking ROCK mint watermark advancement due to fatal error",
				"unchangedWatermark", newRockSig,
				"fatalError", rockErr)
		}

		if !newZenBTCSig.IsZero() && (!newZenBTCSig.Equals(lastKnownZenBTCSig) || zenbtcErr == nil) {
			update.latestSolanaSigs[sidecartypes.SolZenBTCMint] = newZenBTCSig
			if zenbtcErr != nil {
				slog.Info("Advancing zenBTC mint watermark despite partial processing error (failed transactions safely stored in pending queue)",
					"newWatermark", newZenBTCSig,
					"partialError", zenbtcErr)
			}
		} else if zenbtcErr != nil && !newZenBTCSig.IsZero() && newZenBTCSig.Equals(lastKnownZenBTCSig) {
			slog.Warn("Blocking zenBTC mint watermark advancement due to fatal error",
				"unchangedWatermark", newZenBTCSig,
				"fatalError", zenbtcErr)
		}

		if !newZenZECSig.IsZero() && (!newZenZECSig.Equals(lastKnownZenZECSig) || zenzecErr == nil) {
			update.latestSolanaSigs[sidecartypes.SolZenZECMint] = newZenZECSig
			if zenzecErr != nil {
				slog.Info("Advancing zenZEC mint watermark despite partial processing error (failed transactions safely stored in pending queue)",
					"newWatermark", newZenZECSig,
					"partialError", zenzecErr)
			}
		} else if zenzecErr != nil && !newZenZECSig.IsZero() && newZenZECSig.Equals(lastKnownZenZECSig) {
			slog.Warn("Blocking zenZEC mint watermark advancement due to fatal error",
				"unchangedWatermark", newZenZECSig,
				"fatalError", zenzecErr)
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

		var zenBtcEvents, rockEvents, zenZECEvents []api.BurnEvent
		var zenBtcErr, rockErr, zenZECErr error
		var burnWg sync.WaitGroup

		// Fetches new zenBTC burn events from Solana in parallel
		burnWg.Add(1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("Panic recovered in zenBTC burn goroutine", "panic", r)
					zenBtcErr = fmt.Errorf("panic in zenBTC burn processing: %v", r)
				}
				burnWg.Done()
			}()
			lastKnownSig := o.GetLastProcessedSolSignature(sidecartypes.SolZenBTCBurn)
			var newestSig solana.Signature
			zenBtcEvents, newestSig, zenBtcErr = o.getSolanaZenBTCBurnEvents(ctx, sidecartypes.ZenBTCSolanaProgramID[o.Config.Network], lastKnownSig, update, updateMutex)
			if !newestSig.IsZero() {
				lastKnownSig := o.GetLastProcessedSolSignature(sidecartypes.SolZenBTCBurn)
				// Allow watermark advancement if we have progress or no errors
				if !newestSig.Equals(lastKnownSig) || zenBtcErr == nil {
					updateMutex.Lock()
					update.latestSolanaSigs[sidecartypes.SolZenBTCBurn] = newestSig
					updateMutex.Unlock()
					if zenBtcErr != nil {
						slog.Info("Advancing zenBTC burn watermark despite partial processing error (failed transactions safely stored in pending queue)",
							"newWatermark", newestSig,
							"partialError", zenBtcErr)
					}
				} else if zenBtcErr != nil && newestSig.Equals(lastKnownSig) {
					slog.Warn("Blocking zenBTC burn watermark advancement due to fatal error",
						"unchangedWatermark", newestSig,
						"fatalError", zenBtcErr)
				}
			}
		}()

		// Fetches new zenZEC burn events in parallel
		burnWg.Add(1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("Panic recovered in zenZEC burn goroutine", "panic", r)
					zenZECErr = fmt.Errorf("panic in zenZEC burn processing: %v", r)
				}
				burnWg.Done()
			}()
			lastKnownSig := o.GetLastProcessedSolSignature(sidecartypes.SolZenZECBurn)
			var newestSig solana.Signature
			zenZECEvents, newestSig, zenZECErr = o.getSolanaZenZECBurnEvents(ctx, sidecartypes.ZenZECSolanaProgramID[o.Config.Network], lastKnownSig, update, updateMutex)
			if !newestSig.IsZero() {
				lastKnownSig := o.GetLastProcessedSolSignature(sidecartypes.SolZenZECBurn)
				if !newestSig.Equals(lastKnownSig) || zenZECErr == nil {
					updateMutex.Lock()
					update.latestSolanaSigs[sidecartypes.SolZenZECBurn] = newestSig
					updateMutex.Unlock()
					if zenZECErr != nil {
						slog.Info("Advancing zenZEC burn watermark despite partial processing error (failed transactions safely stored in pending queue)",
							"newWatermark", newestSig,
							"partialError", zenZECErr)
					}
				} else if zenZECErr != nil && newestSig.Equals(lastKnownSig) {
					slog.Warn("Blocking zenZEC burn watermark advancement due to fatal error",
						"unchangedWatermark", newestSig,
						"fatalError", zenZECErr)
				}
			}
		}()

		// Fetches new ROCK burn events from Solana in parallel
		burnWg.Add(1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("Panic recovered in ROCK burn goroutine", "panic", r)
					rockErr = fmt.Errorf("panic in ROCK burn processing: %v", r)
				}
				burnWg.Done()
			}()
			lastKnownSig := o.GetLastProcessedSolSignature(sidecartypes.SolRockBurn)
			var newestSig solana.Signature
			rockEvents, newestSig, rockErr = o.getSolanaRockBurnEvents(ctx, sidecartypes.SolRockProgramID[o.Config.Network], lastKnownSig, update, updateMutex)
			if !newestSig.IsZero() {
				lastKnownSig := o.GetLastProcessedSolSignature(sidecartypes.SolRockBurn)
				// Allow watermark advancement if we have progress or no errors
				if !newestSig.Equals(lastKnownSig) || rockErr == nil {
					updateMutex.Lock()
					update.latestSolanaSigs[sidecartypes.SolRockBurn] = newestSig
					updateMutex.Unlock()
					if rockErr != nil {
						slog.Info("Advancing ROCK burn watermark despite partial processing error (failed transactions safely stored in pending queue)",
							"newWatermark", newestSig,
							"partialError", rockErr)
					}
				} else if rockErr != nil && newestSig.Equals(lastKnownSig) {
					slog.Warn("Blocking ROCK burn watermark advancement due to fatal error",
						"unchangedWatermark", newestSig,
						"fatalError", rockErr)
				}
			}
		}()

		burnWg.Wait() // Wait for all parallel burn fetches to complete

		if zenBtcErr != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to process Solana zenBTC burn events, applying partial results: %w", zenBtcErr))
		}
		if zenZECErr != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to process Solana zenZEC burn events, applying partial results: %w", zenZECErr))
		}
		if rockErr != nil {
			sendError(ctx, errChan, fmt.Errorf("failed to process Solana ROCK burn events, applying partial results: %w", rockErr))
		}

		// Merge and sort all new events (which will be empty if fetches failed)
		allNewSolanaBurnEvents := append(zenBtcEvents, zenZECEvents...)
		allNewSolanaBurnEvents = append(allNewSolanaBurnEvents, rockEvents...)
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

		// A map is used to efficiently de-duplicate events by their transaction ID.
		combinedEventsMap := make(map[string]api.BurnEvent)

		// First, preserve any existing pending events that were already added to the state update
		existingPendingEvents := make(map[string]api.BurnEvent)
		for _, event := range update.solanaBurnEvents {
			existingPendingEvents[event.TxID] = event
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

		// Merge the newly fetched Solana burn events with existing events
		mergedEvents := mergeNewBurnEvents(baseEvents, cleanedEvents, allNewSolanaBurnEvents, "Solana")

		// Ensure pending events are preserved in the final result
		// Use hash map for O(1) lookup instead of O(n²) linear search
		mergedEventTxIDs := make(map[string]bool, len(mergedEvents))
		finalEvents := make([]api.BurnEvent, 0, len(mergedEvents)+len(existingPendingEvents))

		// First add all merged events and build lookup map
		for _, event := range mergedEvents {
			finalEvents = append(finalEvents, event)
			mergedEventTxIDs[event.TxID] = true
		}

		// Then add any pending events that weren't already included (O(1) lookup)
		for txID, event := range existingPendingEvents {
			if !mergedEventTxIDs[txID] {
				finalEvents = append(finalEvents, event)
			}
		}

		update.solanaBurnEvents = finalEvents
		update.cleanedSolanaBurnEvents = cleanedEvents

		// Log deduplication results
		if len(existingPendingEvents) > 0 || len(mergedEvents) > 0 {
			// Calculate deduplication statistics only when logging
			pendingEventsPreserved := len(finalEvents) - len(mergedEvents)
			duplicatesRemoved := len(existingPendingEvents) - pendingEventsPreserved

			slog.Info("Burn events merge completed",
				"finalBurnEvents", len(update.solanaBurnEvents),
				"mergedEvents", len(mergedEvents),
				"existingPendingEvents", len(existingPendingEvents),
				"pendingEventsPreserved", pendingEventsPreserved,
				"duplicatesRemoved", duplicatesRemoved)
		}
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
	lastSolZenZECMintSig := o.lastSolZenZECMintSigStr
	lastSolZenZECBurnSig := o.lastSolZenZECBurnSigStr
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
	if sig, ok := update.latestSolanaSigs[sidecartypes.SolZenZECMint]; ok && !sig.IsZero() {
		lastSolZenZECMintSig = sig.String()
	}
	if sig, ok := update.latestSolanaSigs[sidecartypes.SolZenZECBurn]; ok && !sig.IsZero() {
		lastSolZenZECBurnSig = sig.String()
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
		LastSolZenZECMintSig:    lastSolZenZECMintSig,
		LastSolZenZECBurnSig:    lastSolZenZECBurnSig,
		LastSolRockBurnSig:      lastSolRockBurnSig,
		PendingSolanaTxs:        update.pendingTransactions,
	}

	slog.Info("FINAL STATE CONSTRUCTED",
		"finalSolanaMintEvents", len(newState.SolanaMintEvents),
		"finalCleanedSolanaMintEvents", len(newState.CleanedSolanaMintEvents),
		"finalSolanaBurnEvents", len(newState.SolanaBurnEvents),
		"finalCleanedSolanaBurnEvents", len(newState.CleanedSolanaBurnEvents),
		"finalPendingSolanaTxs", len(newState.PendingSolanaTxs))

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

	// Debug: Log the contract address being used
	contractAddr := sidecartypes.ZenBTCControllerAddresses[o.Config.Network]
	slog.Info("DEBUG: Attempting to get redemptions from contract",
		"contractAddress", contractAddr,
		"blockHeight", height.String(),
		"network", o.Config.Network,
		"networkName", sidecartypes.NetworkNames[o.Config.Network],
	)

	redemptionData, err := contractInstance.GetReadyForComplete(callOpts)
	if err != nil {
		slog.Error("DEBUG: Failed to get redemptions from contract",
			"contractAddress", contractAddr,
			"error", err,
		)
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
	remaining := make([]api.SolanaMintEvent, 0)
	updatedCleanedEvents := make(map[string]bool)
	maps.Copy(updatedCleanedEvents, cleanedEvents)

	slog.Info("MINT RECONCILIATION START",
		"inputEventsCount", len(eventsToClean),
		"inputCleanedEventsCount", len(cleanedEvents))

	var zenbtcQueryErrors, zentpQueryErrors int
	var lastZenbtcError, lastZentpError error
	var eventsKeptDueToQuery, eventsRemovedFromChain int

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
			remaining = append(remaining, event)
			eventsKeptDueToQuery++
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
				remaining = append(remaining, event)
				continue
			}
			if zentpResp != nil && len(zentpResp.Mints) > 0 {
				foundOnChain = true
			}
		}

		if !foundOnChain {
			remaining = append(remaining, event)
		} else {
			updatedCleanedEvents[key] = true
			eventsRemovedFromChain++
			slog.Info("Removing Solana mint event from cache as it's now on chain", "txSig", event.TxSig, "sigHash", key)
		}
	}

	// Log reconciliation summary
	slog.Info("MINT RECONCILIATION SUMMARY",
		"inputEventsCount", len(eventsToClean),
		"remainingEventsCount", len(remaining),
		"cleanedEventsCount", len(updatedCleanedEvents),
		"eventsRemovedFromChain", eventsRemovedFromChain,
		"eventsKeptDueToQueryFailure", eventsKeptDueToQuery,
		"zenbtcQueryErrors", zenbtcQueryErrors,
		"zentpQueryErrors", zentpQueryErrors,
		"hasErrors", zenbtcQueryErrors > 0 || zentpQueryErrors > 0)

	// Log summary of any query errors (but not for context cancellation)
	if zenbtcQueryErrors > 0 {
		if !errors.Is(lastZenbtcError, context.Canceled) && !errors.Is(lastZenbtcError, context.DeadlineExceeded) && !strings.Contains(lastZenbtcError.Error(), "context canceled") {
			slog.Warn("Failed to query zrChain for zenBTC mint events, keeping in cache",
				"failedCount", zenbtcQueryErrors,
				"totalEvents", len(eventsToClean),
				"lastError", lastZenbtcError)
		} else {
			slog.Debug("ZrChain ZenBTC query canceled due to context, keeping mint events in cache",
				"failedCount", zenbtcQueryErrors,
				"totalEvents", len(eventsToClean))
		}
	}
	if zentpQueryErrors > 0 {
		if !errors.Is(lastZentpError, context.Canceled) && !errors.Is(lastZentpError, context.DeadlineExceeded) && !strings.Contains(lastZentpError.Error(), "context canceled") {
			slog.Warn("Failed to query zrChain ZenTP for mint events, keeping in cache",
				"failedCount", zentpQueryErrors,
				"totalEvents", len(eventsToClean),
				"lastError", lastZentpError)
		} else {
			slog.Debug("ZrChain ZenTP query canceled due to context, keeping mint events in cache",
				"failedCount", zentpQueryErrors,
				"totalEvents", len(eventsToClean))
		}
	}

	return remaining, updatedCleanedEvents, nil
}

func (o *Oracle) getSolROCKMints(ctx context.Context, programID string, lastKnownSig solana.Signature, update *oracleStateUpdate, updateMutex *sync.Mutex) ([]api.SolanaMintEvent, solana.Signature, error) {
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
			api.Coin_ROCK,
		)
	}

	untypedEvents, newWatermark, err := o.getSolanaEvents(ctx, programID, lastKnownSig, eventTypeName, processor, update, updateMutex)

	// Always process partial results, even if an error occurred.
	mintEvents := make([]api.SolanaMintEvent, len(untypedEvents))
	for i, untypedEvent := range untypedEvents {
		mintEvents[i] = untypedEvent.(api.SolanaMintEvent)
	}

	// Return the (potentially partial) events and the updated watermark along with the error.
	return mintEvents, newWatermark, err
}

func (o *Oracle) getSolZenBTCMints(ctx context.Context, programID string, lastKnownSig solana.Signature, update *oracleStateUpdate, updateMutex *sync.Mutex) ([]api.SolanaMintEvent, solana.Signature, error) {
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
			api.Coin_ZENBTC,
		)
	}

	untypedEvents, newWatermark, err := o.getSolanaEvents(ctx, programID, lastKnownSig, eventTypeName, processor, update, updateMutex)

	// Always process partial results, even if an error occurred.
	mintEvents := make([]api.SolanaMintEvent, len(untypedEvents))
	for i, untypedEvent := range untypedEvents {
		mintEvents[i] = untypedEvent.(api.SolanaMintEvent)
	}

	// Return the (potentially partial) events and the updated watermark along with the error.
	return mintEvents, newWatermark, err
}

func (o *Oracle) getSolZenZECMints(ctx context.Context, programID string, lastKnownSig solana.Signature, update *oracleStateUpdate, updateMutex *sync.Mutex) ([]api.SolanaMintEvent, solana.Signature, error) {
	eventTypeName := "Solana zenZEC mint"
	processor := func(txResult *solrpc.GetTransactionResult, program solana.PublicKey, sig solana.Signature, debugMode bool) ([]any, error) {
		return o.processMintTransaction(txResult, program, sig, debugMode,
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
			func(data any) (solana.PublicKey, uint64, uint64, solana.PublicKey, bool) {
				eventData, ok := data.(*zenbtc_spl_token.TokensMintedWithFeeEventData)
				if !ok {
					return solana.PublicKey{}, 0, 0, solana.PublicKey{}, false
				}
				return eventData.Recipient, eventData.Value, eventData.Fee, eventData.Mint, true
			},
			eventTypeName,
			api.Coin_ZENZEC,
		)
	}

	untypedEvents, newWatermark, err := o.getSolanaEvents(ctx, programID, lastKnownSig, eventTypeName, processor, update, updateMutex)

	mintEvents := make([]api.SolanaMintEvent, len(untypedEvents))
	for i, untypedEvent := range untypedEvents {
		mintEvents[i] = untypedEvent.(api.SolanaMintEvent)
	}

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
	coin api.Coin,
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
				Coin:            coin,
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
func (o *Oracle) getSolanaZenBTCBurnEvents(ctx context.Context, programID string, lastKnownSig solana.Signature, update *oracleStateUpdate, updateMutex *sync.Mutex) ([]api.BurnEvent, solana.Signature, error) {
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
			eventTypeName, chainID, true, api.Coin_ZENBTC,
		)
	}

	untypedEvents, newWatermark, err := o.getSolanaEvents(ctx, programID, lastKnownSig, eventTypeName, processor, update, updateMutex)

	// Always process partial results, even if an error occurred.
	burnEvents := make([]api.BurnEvent, len(untypedEvents))
	for i, untypedEvent := range untypedEvents {
		burnEvents[i] = untypedEvent.(api.BurnEvent)
	}

	// Return the (potentially partial) events and the updated watermark along with the error.
	return burnEvents, newWatermark, err
}

func (o *Oracle) getSolanaZenZECBurnEvents(ctx context.Context, programID string, lastKnownSig solana.Signature, update *oracleStateUpdate, updateMutex *sync.Mutex) ([]api.BurnEvent, solana.Signature, error) {
	eventTypeName := "Solana zenZEC burn"
	chainID := sidecartypes.SolanaCAIP2[o.Config.Network]

	processor := func(txResult *solrpc.GetTransactionResult, program solana.PublicKey, sig solana.Signature, debugMode bool) ([]any, error) {
		return o.processBurnTransaction(txResult, program, sig, debugMode,
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
			func(data any) (destAddr []byte, value uint64, ok bool) {
				eventData, ok := data.(*zenbtc_spl_token.TokenRedemptionEventData)
				if !ok {
					return nil, 0, false
				}
				return eventData.DestAddr, eventData.Value, true
			},
			eventTypeName, chainID, false, api.Coin_ZENZEC,
		)
	}

	untypedEvents, newWatermark, err := o.getSolanaEvents(ctx, programID, lastKnownSig, eventTypeName, processor, update, updateMutex)

	burnEvents := make([]api.BurnEvent, len(untypedEvents))
	for i, untypedEvent := range untypedEvents {
		burnEvents[i] = untypedEvent.(api.BurnEvent)
	}

	return burnEvents, newWatermark, err
}

// getSolanaRockBurnEvents retrieves Rock burn events from Solana.
func (o *Oracle) getSolanaRockBurnEvents(ctx context.Context, programID string, lastKnownSig solana.Signature, update *oracleStateUpdate, updateMutex *sync.Mutex) ([]api.BurnEvent, solana.Signature, error) {
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
			eventTypeName, chainID, false, api.Coin_ROCK,
		)
	}

	untypedEvents, newWatermark, err := o.getSolanaEvents(ctx, programID, lastKnownSig, eventTypeName, processor, update, updateMutex)

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
		newBurnEvents, err := o.fetchAndProcessBackfillRequests(ctx)
		if err != nil {
			// Don't fail the whole operation, just log the error
			// But don't log context cancellation errors as they are expected
			if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) && !strings.Contains(err.Error(), "context canceled") {
				slog.Error("Failed to process backfill requests", "error", err)
			} else {
				slog.Debug("Backfill request processing canceled due to context", "error", err)
			}
			return
		}

		if len(newBurnEvents) == 0 {
			return // No events to process
		}

		updateMutex.Lock()
		defer updateMutex.Unlock()

		// Get cleaned events from the persisted state to check against duplicates
		currentState := o.currentState.Load().(*sidecartypes.OracleState)

		// Use helper function to merge backfilled events with existing ones
		update.solanaBurnEvents = mergeNewBurnEvents(update.solanaBurnEvents, currentState.CleanedSolanaBurnEvents, newBurnEvents, "backfilled Solana")
	}()
}

// fetchAndProcessBackfillRequests queries for backfill requests and processes them
func (o *Oracle) fetchAndProcessBackfillRequests(ctx context.Context) ([]api.BurnEvent, error) {
	backfillResp, err := o.zrChainQueryClient.ValidationQueryClient.BackfillRequests(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query backfill requests: %w", err)
	}

	if backfillResp == nil || backfillResp.BackfillRequests == nil || len(backfillResp.BackfillRequests.Requests) == 0 {
		return []api.BurnEvent{}, nil // No backfill requests
	}

	slog.Info("Found backfill requests to process", "count", len(backfillResp.BackfillRequests.Requests))
	return o.processBackfillRequestsList(ctx, backfillResp.BackfillRequests.Requests)
}

// processBackfillRequestsList processes a slice of backfill requests and returns new burn events
func (o *Oracle) processBackfillRequestsList(ctx context.Context, requests []*validationtypes.MsgTriggerEventBackfill) ([]api.BurnEvent, error) {
	if len(requests) == 0 {
		return []api.BurnEvent{}, nil
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
					return newBurnEvents, ctx.Err()
				}
			}
		}
	}

	return newBurnEvents, nil
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
	case sidecartypes.SolZenZECMint:
		sigStr = o.lastSolZenZECMintSigStr
	case sidecartypes.SolZenZECBurn:
		sigStr = o.lastSolZenZECBurnSigStr
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
		return "none"
	}
	return sig.String()
}

// addPendingTransaction adds a failed transaction to the pending queue in the state update
func (o *Oracle) addPendingTransaction(signature string, eventType string, update *oracleStateUpdate, updateMutex *sync.Mutex) error {
	updateMutex.Lock()
	defer updateMutex.Unlock()

	if update.pendingTransactions == nil {
		update.pendingTransactions = make(map[string]sidecartypes.PendingTxInfo)
	}

	now := time.Now()
	if existing, exists := update.pendingTransactions[signature]; exists {

		updated := sidecartypes.PendingTxInfo{
			Signature:    existing.Signature,
			EventType:    existing.EventType,
			RetryCount:   existing.RetryCount + 1,
			FirstAttempt: existing.FirstAttempt,
			LastAttempt:  now,
		}
		update.pendingTransactions[signature] = updated

		// Verify the update succeeded
		if stored, exists := update.pendingTransactions[signature]; !exists || stored.RetryCount != updated.RetryCount {
			return fmt.Errorf("failed to update pending transaction: %s", signature)
		}

		slog.Debug("Updated pending transaction retry count",
			"signature", signature,
			"eventType", eventType,
			"retryCount", updated.RetryCount)
	} else {
		// Add new pending transaction
		newTx := sidecartypes.PendingTxInfo{
			Signature:    signature,
			EventType:    eventType,
			RetryCount:   1,
			FirstAttempt: now,
			LastAttempt:  now,
		}
		update.pendingTransactions[signature] = newTx

		// Verify the addition succeeded
		if _, exists := update.pendingTransactions[signature]; !exists {
			return fmt.Errorf("failed to add pending transaction: %s", signature)
		}

		slog.Debug("Added new pending transaction",
			"signature", signature,
			"eventType", eventType)
	}

	return nil
}

// shouldRetryTransaction checks if a pending transaction should be retried
func (o *Oracle) shouldRetryTransaction(info sidecartypes.PendingTxInfo) bool {
	// Basic retry limit (can be made configurable later)
	maxRetries := sidecartypes.SolanaPendingTxMaxRetries
	if info.RetryCount >= maxRetries {
		return false
	}

	// Simple time-based retry interval (can be made exponential later)
	retryInterval := sidecartypes.SolanaPendingTxAllowRetryAfter
	return time.Since(info.LastAttempt) >= retryInterval
}

// processPendingTransactionsPersistent continuously retries all pending transactions until context cancellation
// This runs as a persistent background goroutine across all oracle ticks
func (o *Oracle) processPendingTransactionsPersistent(ctx context.Context) {
	slog.Info("Starting persistent pending transaction processing goroutine")

	// Retry loop - continue until context cancellation
	for {
		select {
		case <-ctx.Done():
			slog.Info("Persistent pending transaction processing cancelled")
			return
		default:
			// Check if there are any pending transactions to process
			currentState := o.currentState.Load().(*sidecartypes.OracleState)
			pendingCount := len(currentState.PendingSolanaTxs)

			if pendingCount == 0 {
				// No pending transactions, sleep briefly and check again
				select {
				case <-ctx.Done():
					return
				case <-time.After(sidecartypes.PendingTransactionCheckInterval):
					continue
				}
			}

			slog.Debug("Processing pending transactions round", "count", pendingCount)

			// Process one round of pending transactions
			o.processPendingTransactionsRoundPersistent(ctx)

			// Adaptive sleep based on activity
			sleepDuration := 2 * time.Second
			if pendingCount == 0 {
				sleepDuration = sidecartypes.PendingTransactionCheckInterval // Sleep longer if no transactions were processed
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(sleepDuration):
				// Continue to next round
			}
		}
	}
}

// processPendingTransactions continuously retries all pending transactions until context cancellation
func (o *Oracle) processPendingTransactions(ctx context.Context, wg *sync.WaitGroup, update *oracleStateUpdate, updateMutex *sync.Mutex) {
	defer wg.Done()

	slog.Info("Starting continuous pending transaction processing goroutine")

	// Retry loop - continue until context cancellation
	for {
		select {
		case <-ctx.Done():
			slog.Info("Pending transaction processing cancelled")
			return
		default:
			// Check if there are any pending transactions to process
			updateMutex.Lock()
			pendingCount := len(update.pendingTransactions)
			updateMutex.Unlock()

			if pendingCount == 0 {
				// No pending transactions, sleep briefly and check again
				select {
				case <-ctx.Done():
					return
				case <-time.After(sidecartypes.PendingTransactionCheckInterval):
					continue
				}
			}

			slog.Debug("Processing pending transactions round", "count", pendingCount)

			// Process one round of pending transactions
			pendingStats := o.processPendingTransactionsRound(ctx, update, updateMutex)

			// Log summary if any transactions were processed
			if pendingStats.totalProcessed > 0 {
				slog.Info("Pending transactions processed successfully",
					"successfulCount", len(pendingStats.successfulTxs),
					"totalProcessed", pendingStats.totalProcessed,
					"stillPending", pendingStats.totalProcessed-pendingStats.successCount,
					"successRate", fmt.Sprintf("%.1f%%", float64(pendingStats.successCount)/float64(pendingStats.totalProcessed)*100))
			}

			// Adaptive sleep based on activity
			sleepDuration := 2 * time.Second
			if pendingCount == 0 {
				sleepDuration = sidecartypes.PendingTransactionCheckInterval // Sleep longer if no transactions were processed
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(sleepDuration):
				// Continue to next round
			}
		}
	}
}

// processPendingTransactionsRound processes one round of pending transactions
// Returns (processedCount, successCount) where processedCount is attempts and successCount is completed
func (o *Oracle) processPendingTransactionsRound(ctx context.Context, update *oracleStateUpdate, updateMutex *sync.Mutex) pendingTransactionStats {
	// Create a copy to iterate over to avoid modifying map while iterating
	pendingCopy := make(map[string]sidecartypes.PendingTxInfo)
	updateMutex.Lock()
	maps.Copy(pendingCopy, update.pendingTransactions)
	updateMutex.Unlock()

	stats := pendingTransactionStats{
		successfulTxs: make([]string, 0),
	}

	for signature := range pendingCopy {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			slog.Debug("Pending transaction processing round cancelled",
				"processedInRound", stats.totalProcessed,
				"successfulInRound", stats.successCount)
			return stats
		default:
		}

		// Use live data for retry decision instead of stale snapshot
		updateMutex.Lock()
		current, exists := update.pendingTransactions[signature]
		if !exists {
			updateMutex.Unlock()
			continue
		}

		if !o.shouldRetryTransaction(current) {
			// Check if we should remove transactions that exceeded max retries
			if current.RetryCount >= sidecartypes.PendingTransactionMaxRetries {
				delete(update.pendingTransactions, signature)
				slog.Info("Removed pending transaction after max retries",
					"signature", signature,
					"eventType", current.EventType,
					"retryCount", current.RetryCount)
			}
			updateMutex.Unlock()
			continue
		}
		updateMutex.Unlock()

		// Attempt to retry the transaction
		sig, err := solana.SignatureFromBase58(signature)
		if err != nil {
			slog.Warn("Invalid signature in pending transactions", "signature", signature, "error", err)

			updateMutex.Lock()
			delete(update.pendingTransactions, signature)
			updateMutex.Unlock()
			continue
		}

		// Try to get and process the transaction
		txResult, err := o.retryIndividualTransaction(ctx, sig, current.EventType)
		stats.totalProcessed++ // Count this as a processing attempt

		if err != nil {
			updateMutex.Lock()
			if existing, exists := update.pendingTransactions[signature]; exists {
				updated := sidecartypes.PendingTxInfo{
					Signature:    existing.Signature,
					EventType:    existing.EventType,
					RetryCount:   existing.RetryCount + 1,
					FirstAttempt: existing.FirstAttempt,
					LastAttempt:  time.Now(),
				}
				update.pendingTransactions[signature] = updated
			}
			updateMutex.Unlock()
			slog.Debug("Pending transaction retry failed",
				"signature", signature,
				"eventType", current.EventType,
				"error", err)
			continue
		}

		if txResult != nil {
			// Transaction retrieved successfully, now try to process it
			events, err := o.processTransactionByEventType(txResult, sig, current.EventType)
			if err != nil {
				updateMutex.Lock()
				if existing, exists := update.pendingTransactions[signature]; exists {
					updated := sidecartypes.PendingTxInfo{
						Signature:    existing.Signature,
						EventType:    existing.EventType,
						RetryCount:   existing.RetryCount + 1,
						FirstAttempt: existing.FirstAttempt,
						LastAttempt:  time.Now(),
					}
					update.pendingTransactions[signature] = updated
				}
				updateMutex.Unlock()
				slog.Debug("Pending transaction processing failed",
					"signature", signature,
					"eventType", current.EventType,
					"error", err)
				continue
			}

			// Successfully processed - add events to the state update and remove from pending queue atomically
			updateMutex.Lock()
			if len(events) > 0 {
				// Add events to state update
				switch current.EventType {
				case "Solana ROCK mint", "Solana zenBTC mint":
					for _, event := range events {
						if mintEvent, ok := event.(api.SolanaMintEvent); ok {
							update.SolanaMintEvents = append(update.SolanaMintEvents, mintEvent)
						}
					}
				case "Solana zenBTC burn", "Solana ROCK burn":
					for _, event := range events {
						if burnEvent, ok := event.(api.BurnEvent); ok {
							update.solanaBurnEvents = append(update.solanaBurnEvents, burnEvent)
						}
					}
				}
				stats.successfulTxs = append(stats.successfulTxs, signature)
			}
			// Remove from pending queue in same atomic operation
			if update.pendingTransactions != nil {
				if _, exists := update.pendingTransactions[signature]; exists {
					delete(update.pendingTransactions, signature)
					stats.successCount++
					slog.Debug("Removed pending transaction after successful processing",
						"signature", signature)
				}
			}
			updateMutex.Unlock()
		} else {
			updateMutex.Lock()
			if existing, exists := update.pendingTransactions[signature]; exists {
				updated := sidecartypes.PendingTxInfo{
					Signature:    existing.Signature,
					EventType:    existing.EventType,
					RetryCount:   existing.RetryCount + 1,
					FirstAttempt: existing.FirstAttempt,
					LastAttempt:  time.Now(),
				}
				update.pendingTransactions[signature] = updated
			}
			updateMutex.Unlock()
		}
	}

	// Log processing statistics if any transactions were processed
	if stats.totalProcessed > 0 {
		roundSuccessRate := 0.0
		if stats.totalProcessed > 0 {
			roundSuccessRate = float64(stats.successCount) / float64(stats.totalProcessed) * 100
		}

		slog.Debug("Pending transaction processing round completed",
			"totalProcessed", stats.totalProcessed,
			"successfullyCompleted", stats.successCount,
			"stillPending", len(pendingCopy)-stats.successCount,
			"successRate", fmt.Sprintf("%.1f%%", roundSuccessRate))
	}

	return stats
}

// processPendingTransactionsRoundPersistent processes one round of pending transactions from the current state
// This version works with the persistent background processor and modifies the live oracle state
func (o *Oracle) processPendingTransactionsRoundPersistent(ctx context.Context) {
	currentState := o.currentState.Load().(*sidecartypes.OracleState)

	// Create a copy to iterate over to avoid modifying map while iterating
	pendingCopy := make(map[string]sidecartypes.PendingTxInfo)
	for k, v := range currentState.PendingSolanaTxs {
		pendingCopy[k] = v
	}

	processedCount := 0
	successCount := 0
	var successfulTransactions []string

	for signature := range pendingCopy {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			slog.Debug("Pending transaction processing round cancelled",
				"processedInRound", processedCount,
				"successfulInRound", successCount)
			return
		default:
		}

		// Get current transaction info (may have been updated by other goroutines)
		currentState = o.currentState.Load().(*sidecartypes.OracleState)
		current, exists := currentState.PendingSolanaTxs[signature]
		if !exists {
			continue
		}

		if !o.shouldRetryTransaction(current) {
			// Check if we should remove transactions that exceeded max retries
			if current.RetryCount >= sidecartypes.PendingTransactionMaxRetries {
				o.removePendingTransactionFromState(signature)
				slog.Info("Removed pending transaction after max retries",
					"signature", signature,
					"eventType", current.EventType,
					"retryCount", current.RetryCount)
			}
			continue
		}

		// Attempt to retry the transaction
		sig, err := solana.SignatureFromBase58(signature)
		if err != nil {
			slog.Warn("Invalid signature in pending transactions", "signature", signature, "error", err)
			o.removePendingTransactionFromState(signature)
			continue
		}

		// Try to get and process the transaction
		txResult, err := o.retryIndividualTransaction(ctx, sig, current.EventType)
		processedCount++ // Count this as a processing attempt

		if err != nil {
			o.updatePendingTransactionInState(signature, current)
			slog.Debug("Pending transaction retry failed",
				"signature", signature,
				"eventType", current.EventType,
				"error", err)
			continue
		}

		if txResult != nil {
			// Transaction retrieved successfully, now try to process it
			events, err := o.processTransactionByEventType(txResult, sig, current.EventType)
			if err != nil {
				o.updatePendingTransactionInState(signature, current)
				slog.Debug("Pending transaction processing failed",
					"signature", signature,
					"eventType", current.EventType,
					"error", err)
				continue
			}

			// Successfully processed - add events to the current state and remove from pending queue
			if len(events) > 0 {
				o.addEventsToCurrentState(events, current.EventType)
				successfulTransactions = append(successfulTransactions, signature)
			}

			// Remove from pending queue
			o.removePendingTransactionFromState(signature)
			successCount++
			slog.Debug("Removed pending transaction after successful processing",
				"signature", signature)
		} else {
			o.updatePendingTransactionInState(signature, current)
		}
	}

	// Log processing statistics if any transactions were processed
	if processedCount > 0 {
		roundSuccessRate := 0.0
		if processedCount > 0 {
			roundSuccessRate = float64(successCount) / float64(processedCount) * 100
		}

		slog.Debug("Pending transaction processing round completed",
			"totalProcessed", processedCount,
			"successfullyCompleted", successCount,
			"stillPending", len(pendingCopy)-successCount,
			"successRate", fmt.Sprintf("%.1f%%", roundSuccessRate))
	}

	// Log summary if any transactions were successfully processed
	if len(successfulTransactions) > 0 {
		slog.Info("Pending transactions processed successfully",
			"successfulCount", len(successfulTransactions),
			"totalProcessed", processedCount,
			"stillPending", len(pendingCopy)-successCount,
			"successRate", fmt.Sprintf("%.1f%%", float64(successCount)/float64(processedCount)*100))
	}
}

// Helper functions for managing pending transactions in the current state
func (o *Oracle) removePendingTransactionFromState(signature string) {
	for {
		currentState := o.currentState.Load().(*sidecartypes.OracleState)
		newPendingTxs := make(map[string]sidecartypes.PendingTxInfo)
		for k, v := range currentState.PendingSolanaTxs {
			if k != signature {
				newPendingTxs[k] = v
			}
		}

		newState := *currentState
		newState.PendingSolanaTxs = newPendingTxs

		if o.currentState.CompareAndSwap(currentState, &newState) {
			break
		}
		// Retry if CAS failed due to concurrent modification
	}
}

func (o *Oracle) updatePendingTransactionInState(signature string, txInfo sidecartypes.PendingTxInfo) {
	for {
		currentState := o.currentState.Load().(*sidecartypes.OracleState)
		newPendingTxs := make(map[string]sidecartypes.PendingTxInfo)
		for k, v := range currentState.PendingSolanaTxs {
			newPendingTxs[k] = v
		}

		// Update the transaction with incremented retry count
		updated := sidecartypes.PendingTxInfo{
			Signature:    txInfo.Signature,
			EventType:    txInfo.EventType,
			RetryCount:   txInfo.RetryCount + 1,
			FirstAttempt: txInfo.FirstAttempt,
			LastAttempt:  time.Now(),
		}
		newPendingTxs[signature] = updated

		newState := *currentState
		newState.PendingSolanaTxs = newPendingTxs

		if o.currentState.CompareAndSwap(currentState, &newState) {
			break
		}
		// Retry if CAS failed due to concurrent modification
	}
}

func (o *Oracle) addEventsToCurrentState(events []any, eventType string) {
	for {
		currentState := o.currentState.Load().(*sidecartypes.OracleState)
		newState := *currentState

		// Add events based on event type
		switch eventType {
		case "Solana ROCK mint", "Solana zenBTC mint":
			newMintEvents := make([]api.SolanaMintEvent, len(currentState.SolanaMintEvents))
			copy(newMintEvents, currentState.SolanaMintEvents)
			for _, event := range events {
				if mintEvent, ok := event.(api.SolanaMintEvent); ok {
					newMintEvents = append(newMintEvents, mintEvent)
				}
			}
			newState.SolanaMintEvents = newMintEvents
		case "Solana zenBTC burn", "Solana ROCK burn":
			newBurnEvents := make([]api.BurnEvent, len(currentState.SolanaBurnEvents))
			copy(newBurnEvents, currentState.SolanaBurnEvents)
			for _, event := range events {
				if burnEvent, ok := event.(api.BurnEvent); ok {
					newBurnEvents = append(newBurnEvents, burnEvent)
				}
			}
			newState.SolanaBurnEvents = newBurnEvents
		}

		if o.currentState.CompareAndSwap(currentState, &newState) {
			break
		}
		// Retry if CAS failed due to concurrent modification
	}
}

// processTransactionByEventType processes a transaction based on its event type
func (o *Oracle) processTransactionByEventType(txResult *solrpc.GetTransactionResult, sig solana.Signature, eventType string) ([]any, error) {
	// Determine program ID and processor function based on event type
	var programID string
	var processor processTransactionFunc

	switch eventType {
	case "Solana ROCK mint":
		programID = sidecartypes.SolRockProgramID[o.Config.Network]
		processor = func(txResult *solrpc.GetTransactionResult, program solana.PublicKey, sig solana.Signature, debugMode bool) ([]any, error) {
			return o.processMintTransaction(txResult, program, sig, debugMode,
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
				func(data any) (solana.PublicKey, uint64, uint64, solana.PublicKey, bool) {
					eventData, ok := data.(*rock_spl_token.TokensMintedWithFeeEventData)
					if !ok {
						return solana.PublicKey{}, 0, 0, solana.PublicKey{}, false
					}
					return eventData.Recipient, eventData.Value, eventData.Fee, eventData.Mint, true
				},
				eventType,
				api.Coin_ROCK,
			)
		}
	case "Solana zenBTC mint":
		programID = sidecartypes.ZenBTCSolanaProgramID[o.Config.Network]
		processor = func(txResult *solrpc.GetTransactionResult, program solana.PublicKey, sig solana.Signature, debugMode bool) ([]any, error) {
			return o.processMintTransaction(txResult, program, sig, debugMode,
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
				func(data any) (solana.PublicKey, uint64, uint64, solana.PublicKey, bool) {
					eventData, ok := data.(*zenbtc_spl_token.TokensMintedWithFeeEventData)
					if !ok {
						return solana.PublicKey{}, 0, 0, solana.PublicKey{}, false
					}
					return eventData.Recipient, eventData.Value, eventData.Fee, eventData.Mint, true
				},
				eventType,
				api.Coin_ZENBTC,
			)
		}
	case "Solana zenBTC burn":
		programID = sidecartypes.ZenBTCSolanaProgramID[o.Config.Network]
		processor = func(txResult *solrpc.GetTransactionResult, program solana.PublicKey, sig solana.Signature, debugMode bool) ([]any, error) {
			chainID := sidecartypes.SolanaCAIP2[o.Config.Network]
			return o.processBurnTransaction(txResult, program, sig, debugMode,
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
				func(data any) (destAddr []byte, value uint64, ok bool) {
					eventData, ok := data.(*zenbtc_spl_token.TokenRedemptionEventData)
					if !ok {
						return nil, 0, false
					}
					return eventData.DestAddr[:], eventData.Value, true
				},
				eventType, chainID, true, api.Coin_ZENBTC,
			)
		}
	case "Solana ROCK burn":
		programID = sidecartypes.SolRockProgramID[o.Config.Network]
		processor = func(txResult *solrpc.GetTransactionResult, program solana.PublicKey, sig solana.Signature, debugMode bool) ([]any, error) {
			chainID := sidecartypes.SolanaCAIP2[o.Config.Network]
			return o.processBurnTransaction(txResult, program, sig, debugMode,
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
				func(data any) (destAddr []byte, value uint64, ok bool) {
					eventData, ok := data.(*rock_spl_token.TokenRedemptionEventData)
					if !ok {
						return nil, 0, false
					}
					return eventData.DestAddr[:], eventData.Value, true
				},
				eventType, chainID, false, api.Coin_ROCK,
			)
		}
	default:
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}

	// Parse program ID
	program, err := solana.PublicKeyFromBase58(programID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse program ID for %s: %w", eventType, err)
	}

	// Process the transaction
	return processor(txResult, program, sig, o.DebugMode)
}

// addEventsToStateUpdate adds processed events to the state update
func (o *Oracle) addEventsToStateUpdate(events []any, eventType string, update *oracleStateUpdate, updateMutex *sync.Mutex) {
	updateMutex.Lock()
	defer updateMutex.Unlock()

	switch eventType {
	case "Solana ROCK mint", "Solana zenBTC mint":
		// Convert to mint events and add to state update
		for _, event := range events {
			if mintEvent, ok := event.(api.SolanaMintEvent); ok {
				update.SolanaMintEvents = append(update.SolanaMintEvents, mintEvent)
			}
		}
	case "Solana zenBTC burn", "Solana ROCK burn":
		// Convert to burn events and add to state update
		for _, event := range events {
			if burnEvent, ok := event.(api.BurnEvent); ok {
				update.solanaBurnEvents = append(update.solanaBurnEvents, burnEvent)
			}
		}
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

// getSolanaEvents is an optimized generic helper to fetch signatures for a given program, detect and heal
// gaps using the watermark (lastKnownSig), then download and process each transaction using the
// provided `processTransaction` callback. If any part of the transaction processing pipeline fails,
// it returns any partially processed events along with the watermark of the last successfully
// processed transaction. This allows the oracle to make incremental progress.

// PERFORMANCE OPTIMIZATIONS:
// - Rate limiting with semaphore to prevent RPC overload
// - Transaction caching with TTL
// - Parallel batch processing with improved error handling
// - Exponential backoff retry strategy
// - Memory-efficient slice pre-allocation
// processSignatures takes a list of transaction signatures and processes them.
// buildBatchRequests creates RPC batch requests for a set of transaction signatures
func (o *Oracle) buildBatchRequests(currentBatch []*solrpc.TransactionSignature) jsonrpc.RPCRequests {
	batchRequests := make(jsonrpc.RPCRequests, 0, sidecartypes.SolanaEventFetchBatchSize)

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

	return batchRequests
}

// executeBatchRequest executes a batch of RPC requests and handles batch-level errors
func (o *Oracle) executeBatchRequest(
	ctx context.Context,
	currentBatch []*solrpc.TransactionSignature,
	eventTypeName string,
	currentBatchSize int,
	minBatchSize int,
) (jsonrpc.RPCResponses, int, error) {
	batchRequests := o.buildBatchRequests(currentBatch)

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

	// If the batch failed, reduce the batch size for retry
	if batchErr != nil {
		newBatchSize := max(currentBatchSize/2, minBatchSize)
		return nil, newBatchSize, batchErr
	}

	return batchResponses, currentBatchSize, nil
}

// processIndividualTransaction handles the processing of a single transaction within a batch
func (o *Oracle) processIndividualTransaction(
	ctx context.Context,
	sigInfo *solrpc.TransactionSignature,
	resp *jsonrpc.RPCResponse,
	program solana.PublicKey,
	eventTypeName string,
	processTransaction processTransactionFunc,
	failedSignatures []string,
	newestSigProcessed solana.Signature,
	lastSuccessfullyProcessedSig solana.Signature,
	allEvents []any,
	stats *transactionProcessingStats,
) ([]string, solana.Signature, solana.Signature, []any) {
	// Handle nil results from RPC (transaction not found/retrievable)
	if resp.Result == nil {
		stats.totalNilResults++
		slog.Debug("Transaction returned nil result, attempting individual retry",
			"eventType", eventTypeName,
			"signature", sigInfo.Signature,
			"responseID", resp.ID,
			"responseError", resp.Error,
			"nilResultCount", stats.totalNilResults)

		// Retry individual transaction
		if retryResult, err := o.retryIndividualTransaction(ctx, sigInfo.Signature, eventTypeName); err != nil {
			// Add to pending queue instead of stopping
			stats.individualRetryFailures++
			stats.processingErrors++
			failedSignatures, newestSigProcessed = o.handleTransactionFailure(
				failedSignatures, sigInfo.Signature, newestSigProcessed, eventTypeName,
				"Individual transaction retry failed", err)
			return failedSignatures, newestSigProcessed, lastSuccessfullyProcessedSig, allEvents
		} else if retryResult != nil {
			// Process the successfully retried transaction
			events, err := processTransaction(retryResult, program, sigInfo.Signature, o.DebugMode)
			if err != nil {
				stats.processingErrors++
				failedSignatures, newestSigProcessed = o.handleTransactionFailure(
					failedSignatures, sigInfo.Signature, newestSigProcessed, eventTypeName,
					"Failed to process retried transaction", err)
				return failedSignatures, newestSigProcessed, lastSuccessfullyProcessedSig, allEvents
			}

			if len(events) > 0 {
				stats.successfulTransactions++
				stats.successfulWithEvents++
				stats.newTransactionsProcessed++
				slog.Debug("Retried transaction processed with events",
					"eventType", eventTypeName,
					"signature", sigInfo.Signature,
					"eventCount", len(events))
				allEvents = append(allEvents, events...)
			} else {
				stats.emptyTransactions++
				stats.successfulWithoutEvents++
				stats.newTransactionsProcessed++
				slog.Debug("Retried transaction processed but contained no events",
					"eventType", eventTypeName,
					"signature", sigInfo.Signature)
			}
			lastSuccessfullyProcessedSig = sigInfo.Signature
			newestSigProcessed = sigInfo.Signature
			return failedSignatures, newestSigProcessed, lastSuccessfullyProcessedSig, allEvents
		} else {
			// Still nil after retry - add to pending queue
			stats.individualRetryFailures++
			failedSignatures, newestSigProcessed = o.handleTransactionFailure(
				failedSignatures, sigInfo.Signature, newestSigProcessed, eventTypeName,
				"Transaction still nil after retry", nil)
			return failedSignatures, newestSigProcessed, lastSuccessfullyProcessedSig, allEvents
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
		failedSignatures, newestSigProcessed = o.handleTransactionFailure(
			failedSignatures, sigInfo.Signature, newestSigProcessed, eventTypeName,
			"Unmarshal error", err)
		return failedSignatures, newestSigProcessed, lastSuccessfullyProcessedSig, allEvents
	}

	o.cacheTransactionResult(sigInfo.Signature.String(), &txRes)

	events, err := processTransaction(&txRes, program, sigInfo.Signature, o.DebugMode)
	if err != nil {
		failedSignatures, newestSigProcessed = o.handleTransactionFailure(
			failedSignatures, sigInfo.Signature, newestSigProcessed, eventTypeName,
			"Processing error", err)
		return failedSignatures, newestSigProcessed, lastSuccessfullyProcessedSig, allEvents
	}

	if len(events) > 0 {
		stats.successfulTransactions++
		stats.successfulWithEvents++
		stats.newTransactionsProcessed++
		slog.Debug("Successfully processed transaction",
			"eventType", eventTypeName,
			"signature", sigInfo.Signature,
			"eventCount", len(events))
		allEvents = append(allEvents, events...)
	} else {
		stats.emptyTransactions++
		stats.successfulWithoutEvents++
		stats.newTransactionsProcessed++
		slog.Debug("Transaction processed but contained no events",
			"eventType", eventTypeName,
			"signature", sigInfo.Signature)
	}
	lastSuccessfullyProcessedSig = sigInfo.Signature
	newestSigProcessed = sigInfo.Signature
	return failedSignatures, newestSigProcessed, lastSuccessfullyProcessedSig, allEvents
}

// transactionProcessingStats holds statistics for transaction processing
type transactionProcessingStats struct {
	totalNilResults              int
	individualRetryFailures      int
	successfulTransactions       int
	notFoundTransactions         int
	processingErrors             int
	emptyTransactions            int
	successfulWithEvents         int
	successfulWithoutEvents      int
	pendingTransactionsProcessed int
	newTransactionsProcessed     int
}

// pendingTransactionStats holds statistics for pending transaction processing
type pendingTransactionStats struct {
	totalProcessed int
	successCount   int
	successfulTxs  []string
}

// cryptoPrices holds both BTC and ETH price data
type cryptoPrices struct {
	BTCUSDPrice math.LegacyDec
	ETHUSDPrice math.LegacyDec
}

// burnEventResult holds the result of processing burn events
type burnEventResult struct {
	ethBurnEvents        []api.BurnEvent
	cleanedEthBurnEvents map[string]bool
}

func (o *Oracle) processSignatures(
	ctx context.Context,
	signatures []*solrpc.TransactionSignature,
	program solana.PublicKey,
	eventTypeName string,
	processTransaction processTransactionFunc,
) ([]any, solana.Signature, []string, error) {
	// Create events slice directly
	allEvents := make([]any, 0, sidecartypes.InitialEventsSliceCapacity)
	// Collect failed transaction signatures
	failedSignatures := make([]string, 0)

	var lastSuccessfullyProcessedSig solana.Signature
	var newestSigProcessed solana.Signature
	// Adaptive batching parameters
	currentBatchSize := sidecartypes.SolanaEventFetchBatchSize
	minBatchSize := sidecartypes.SolanaEventFetchMinBatchSize

	// Track processing results for debugging
	stats := &transactionProcessingStats{}

	// Pre-allocate with estimated capacity to reduce allocations
	estimatedEvents := len(signatures) * 2 // Estimate 2 events per signature on average
	if cap(allEvents) < estimatedEvents {
		allEvents = make([]any, 0, estimatedEvents)
	}

	// Process signatures with adaptive batching
	for i := 0; i < len(signatures); {
		if ctx.Err() != nil {
			// Add ALL remaining unprocessed signatures to pending queue for optimistic watermark advancement
			for j := i; j < len(signatures); j++ {
				failedSignatures = append(failedSignatures, signatures[j].Signature.String())
				newestSigProcessed = signatures[j].Signature
			}
			slog.Debug("Context canceled during signature processing, added remaining signatures to pending queue",
				"eventType", eventTypeName,
				"processedSignatures", i,
				"totalSignatures", len(signatures),
				"addedToPending", len(signatures)-i)
			return allEvents, newestSigProcessed, failedSignatures, ctx.Err()
		}

		end := min(i+currentBatchSize, len(signatures))
		currentBatch := signatures[i:end]

		// Execute batch request using helper function
		batchResponses, newBatchSize, batchErr := o.executeBatchRequest(ctx, currentBatch, eventTypeName, currentBatchSize, minBatchSize)

		// If the batch failed, handle retry logic
		if batchErr != nil {
			// Check for context cancellation first - if cancelled, stop immediately
			if errors.Is(batchErr, context.Canceled) || errors.Is(batchErr, context.DeadlineExceeded) || strings.Contains(batchErr.Error(), "context canceled") {
				// Context cancelled - add remaining signatures to pending queue and return
				for j := i; j < len(signatures); j++ {
					failedSignatures = append(failedSignatures, signatures[j].Signature.String())
					newestSigProcessed = signatures[j].Signature
				}
				return allEvents, newestSigProcessed, failedSignatures, ctx.Err()
			}

			if currentBatchSize > minBatchSize {
				slog.Warn("Batch GetTransaction failed, reducing batch size and retrying",
					"eventType", eventTypeName, "error", batchErr, "oldSize", currentBatchSize, "newSize", newBatchSize)
				currentBatchSize = newBatchSize
			} else {
				// If we're already at the minimum batch size, try individual requests as fallback
				slog.Warn("Batch transaction fetch failed at minimum batch size, attempting individual fallback", "eventType", eventTypeName, "size", len(currentBatch))

				// Try individual requests for each transaction in the failed batch
				for idx, sigInfo := range currentBatch {
					// Check for context cancellation before each individual request
					if ctx.Err() != nil {
						// Add ALL remaining unprocessed transactions to failed signatures and return
						for j := idx; j < len(currentBatch); j++ {
							failedSignatures = append(failedSignatures, currentBatch[j].Signature.String())
							newestSigProcessed = currentBatch[j].Signature
						}
						// Return newestSigProcessed to advance watermark to cover failed transactions
						return allEvents, newestSigProcessed, failedSignatures, ctx.Err()
					}

					// Attempt individual retry
					retryResult, err := o.retryIndividualTransaction(ctx, sigInfo.Signature, eventTypeName)
					if err != nil {
						// Individual retry failed - add to pending queue
						stats.individualRetryFailures++
						stats.processingErrors++
						failedSignatures, newestSigProcessed = o.handleTransactionFailure(
							failedSignatures, sigInfo.Signature, newestSigProcessed, eventTypeName,
							"Individual fallback retry failed", err)
						continue
					} else if retryResult != nil {
						// Individual retry succeeded - process the transaction
						events, err := processTransaction(retryResult, program, sigInfo.Signature, o.DebugMode)
						if err != nil {
							stats.processingErrors++
							failedSignatures, newestSigProcessed = o.handleTransactionFailure(
								failedSignatures, sigInfo.Signature, newestSigProcessed, eventTypeName,
								"Failed to process individually retried transaction", err)
							continue
						}

						if len(events) > 0 {
							stats.successfulTransactions++
							slog.Debug("Successfully processed individual fallback transaction",
								"eventType", eventTypeName,
								"signature", sigInfo.Signature,
								"eventCount", len(events))
							allEvents = append(allEvents, events...)
						} else {
							stats.emptyTransactions++
							slog.Debug("Individual fallback transaction processed but contained no events",
								"eventType", eventTypeName,
								"signature", sigInfo.Signature)
						}
						lastSuccessfullyProcessedSig = sigInfo.Signature
						newestSigProcessed = sigInfo.Signature
					} else {
						// Individual retry returned nil - add to pending queue
						stats.individualRetryFailures++
						failedSignatures, newestSigProcessed = o.handleTransactionFailure(
							failedSignatures, sigInfo.Signature, newestSigProcessed, eventTypeName,
							"Individual fallback returned nil", nil)
					}
				}

				// Advance to the next segment since this batch is handled
				i += len(currentBatch)
				// Reset batch size for next attempt
				currentBatchSize = sidecartypes.SolanaEventFetchBatchSize
			}
			time.Sleep(sidecartypes.SolanaEventFetchRetrySleep) // Pause before retrying
			continue                                            // Retry the same segment `i`
		}

		// Success: Process the batch responses
		responseMap := make(map[int]*jsonrpc.RPCResponse, len(batchResponses))
		for _, resp := range batchResponses {
			if idx, ok := parseRPCResponseID(resp, eventTypeName); ok {
				responseMap[idx] = resp
			}
		}

		// Process responses in order
		for idx, sigInfo := range currentBatch {
			resp, exists := responseMap[idx]
			if !exists {
				failedSignatures, newestSigProcessed = o.handleTransactionFailure(
					failedSignatures, sigInfo.Signature, newestSigProcessed, eventTypeName,
					"Missing batch response", nil)
				continue
			}

			failedSignatures, newestSigProcessed, lastSuccessfullyProcessedSig, allEvents = o.processIndividualTransaction(
				ctx, sigInfo, resp, program, eventTypeName, processTransaction,
				failedSignatures, newestSigProcessed, lastSuccessfullyProcessedSig, allEvents, stats)
		}

		// Advance to the next segment
		i += len(currentBatch)

		// Optional: slowly increase batch size on success
		if currentBatchSize < sidecartypes.SolanaEventFetchBatchSize {
			currentBatchSize = min(currentBatchSize+minBatchSize, sidecartypes.SolanaEventFetchBatchSize)
		}
		time.Sleep(sidecartypes.SolanaSleepInterval)
	}

	// Calculate processing statistics
	totalProcessed := len(signatures)
	successRate := float64(stats.successfulTransactions) / float64(totalProcessed) * 100

	// Transaction accounting verification - ensure no transactions are lost
	accountedTransactions := stats.successfulTransactions + stats.emptyTransactions + len(failedSignatures)
	if accountedTransactions != totalProcessed {
		slog.Error("Transaction accounting mismatch - blocking watermark advancement",
			"eventType", eventTypeName,
			"totalSignatures", totalProcessed,
			"successful", stats.successfulTransactions,
			"empty", stats.emptyTransactions,
			"failed", len(failedSignatures),
			"accounted", accountedTransactions,
			"missing", totalProcessed-accountedTransactions)

		// Don't advance watermark - return last known good watermark
		return allEvents, lastSuccessfullyProcessedSig, failedSignatures,
			fmt.Errorf("transaction accounting mismatch: %d processed, %d accounted",
				totalProcessed, accountedTransactions)
	}

	// Summary log with comprehensive batch processing results
	slog.Info("Batch processing summary",
		"eventType", eventTypeName,
		"totalSignatures", totalProcessed,
		"successfulTransactions", stats.successfulTransactions,
		"successfulWithEvents", stats.successfulWithEvents,
		"successfulWithoutEvents", stats.successfulWithoutEvents,
		"nilResults", stats.totalNilResults,
		"notFoundTransactions", stats.notFoundTransactions,
		"processingErrors", stats.processingErrors,
		"individualRetryFailures", stats.individualRetryFailures,
		"successRate", fmt.Sprintf("%.1f%%", successRate),
		"extractedEvents", len(allEvents),
		"newWatermark", newestSigProcessed,
	)

	// Log pending transaction strategy
	if stats.processingErrors > 0 || stats.totalNilResults > 0 {
		slog.Info("Optimistic watermark advancement with pending queue",
			"eventType", eventTypeName,
			"strategy", "continue_with_pending_queue",
			"reason", "prevents_system_stalls",
			"pendingTransactions", "will_retry_failed_transactions")
	}

	// Summary logging for warning-level issues
	if stats.totalNilResults > 0 {
		slog.Warn("Transaction processing summary - Nil results encountered",
			"eventType", eventTypeName,
			"nilResultCount", stats.totalNilResults,
			"totalProcessed", len(signatures),
			"successfulEvents", len(allEvents))
	}

	if stats.individualRetryFailures > 0 {
		slog.Warn("Transaction processing summary - Individual retry failures",
			"eventType", eventTypeName,
			"retryFailureCount", stats.individualRetryFailures,
			"totalProcessed", len(signatures),
			"addedToPendingQueue", stats.individualRetryFailures)
	}

	if len(failedSignatures) > 0 {
		slog.Warn("Transaction processing summary - Failed signatures added to pending queue",
			"eventType", eventTypeName,
			"failedSignatureCount", len(failedSignatures),
			"totalProcessed", len(signatures),
			"willRetryInNextCycle", true)
	}
	if len(allEvents) > 0 {
		slog.Debug("Successfully extracted events from Solana transactions", "eventType", eventTypeName, "extractedEvents", len(allEvents))
	}
	// Return the newest signature processed for watermark advancement
	watermarkSig := newestSigProcessed
	if watermarkSig.IsZero() {
		watermarkSig = lastSuccessfullyProcessedSig
	}

	return allEvents, watermarkSig, failedSignatures, nil
}

// getSolanaEvents is an optimized generic helper to fetch signatures for a given program, detect and heal
// gaps using the watermark (lastKnownSig), then download and process each transaction using the
// provided `processTransaction` callback. If any part of the transaction processing pipeline fails,
// it returns any partially processed events along with the watermark of the last successfully
// processed transaction. This allows the oracle to make incremental progress.

// PERFORMANCE OPTIMIZATIONS:
// - Rate limiting with semaphore to prevent RPC overload
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
	update *oracleStateUpdate,
	updateMutex *sync.Mutex,
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
		return []any{}, lastKnownSig, nil
	}

	slog.Info("Found new signatures", "eventType", eventTypeName, "count", len(newSignatures), "watermark", formatWatermarkForLogging(lastKnownSig), "newest", newestSigFromNode)

	events, lastSig, failedSignatures, err := o.processSignatures(ctx, newSignatures, program, eventTypeName, processTransaction)

	// Handle failed signatures by adding them to pending queue
	// Only advance watermark if ALL pending store operations succeed
	var pendingStoreErrors []error
	for _, failedSig := range failedSignatures {
		if pendingErr := o.addPendingTransaction(failedSig, eventTypeName, update, updateMutex); pendingErr != nil {
			pendingStoreErrors = append(pendingStoreErrors, pendingErr)
			slog.Error("Failed to add transaction to pending store",
				"signature", failedSig,
				"eventType", eventTypeName,
				"error", pendingErr)
		}
	}

	// If any pending store operations failed, don't advance watermark
	if len(pendingStoreErrors) > 0 {
		slog.Error("Pending store operations failed - blocking watermark advancement",
			"eventType", eventTypeName,
			"failedOperations", len(pendingStoreErrors),
			"totalFailed", len(failedSignatures))
		return events, lastKnownSig, fmt.Errorf("pending store operations failed: %d errors", len(pendingStoreErrors))
	}

	if err != nil {
		return events, lastSig, err
	}

	return events, lastSig, nil
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
		ExpiresAt: now.Add(sidecartypes.TransactionCacheTTL),
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
	coin api.Coin,
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
			Coint:     coin,
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
		if !errors.Is(lastZenbtcError, context.Canceled) && !errors.Is(lastZenbtcError, context.DeadlineExceeded) && !strings.Contains(lastZenbtcError.Error(), "context canceled") {
			slog.Error("Failed to query zrChain for zenBTC burn events",
				"failedCount", zenbtcQueryErrors,
				"totalEvents", len(eventsToClean),
				"chainType", chainTypeName,
				"lastError", lastZenbtcError)
		} else {
			slog.Debug("ZrChain zenBTC burn query canceled due to context",
				"failedCount", zenbtcQueryErrors,
				"totalEvents", len(eventsToClean),
				"chainType", chainTypeName)
		}
	}
	if zentpQueryErrors > 0 {
		if !errors.Is(lastZentpError, context.Canceled) && !errors.Is(lastZentpError, context.DeadlineExceeded) && !strings.Contains(lastZentpError.Error(), "context canceled") {
			slog.Error("Failed to query zrChain for zenTP burn events",
				"failedCount", zentpQueryErrors,
				"totalEvents", len(eventsToClean),
				"chainType", chainTypeName,
				"lastError", lastZentpError)
		} else {
			slog.Debug("ZrChain ZenTP burn query canceled due to context",
				"failedCount", zentpQueryErrors,
				"totalEvents", len(eventsToClean),
				"chainType", chainTypeName)
		}
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

// --- Periodic Reset Logic ---

// maybePerformScheduledReset checks if a full state reset is due at the provided UTC time.
// It is safe to call on every tick; internal locking ensures idempotence per boundary.
func (o *Oracle) maybePerformScheduledReset(nowUTC time.Time) {
	o.resetMutex.Lock()
	defer o.resetMutex.Unlock()

	// Derive interval dynamically (test flag overrides).
	interval := time.Duration(sidecartypes.OracleStateResetIntervalHours) * time.Hour
	if o.ForceTestReset {
		interval = 2 * time.Minute
	}
	if interval <= 0 {
		interval = 24 * time.Hour
	}

	// If nextScheduledReset not set OR we switched into a shorter test interval, (re)compute it.
	if o.nextScheduledReset.IsZero() || (interval < time.Hour && o.nextScheduledReset.Sub(nowUTC) > interval) {
		o.scheduleNextReset(nowUTC, interval)
	}

	if !o.nextScheduledReset.IsZero() &&
		(nowUTC.Equal(o.nextScheduledReset) || nowUTC.After(o.nextScheduledReset)) {
		log.Printf("Performing scheduled oracle state reset at %s (nextScheduledReset=%s, interval=%s)",
			nowUTC.Format(time.RFC3339), o.nextScheduledReset.Format(time.RFC3339), interval)

		o.performFullStateResetLocked()

		// Schedule the next reset strictly after 'nowUTC'
		o.scheduleNextReset(nowUTC.Add(time.Second), interval)
		log.Printf("Next scheduled oracle state reset at %s", o.nextScheduledReset.Format(time.RFC3339))
	}
}

// scheduleNextReset computes the next UTC boundary time > now (or == now if exactly aligned)
// based on the provided interval (hours or test mode minutes).
func (o *Oracle) scheduleNextReset(nowUTC time.Time, interval time.Duration) {
	if interval <= 0 {
		interval = 24 * time.Hour
	}

	// For sub-hour (test) intervals, align to that minute boundary.
	if interval < time.Hour {
		trunc := nowUTC.Truncate(interval)
		if nowUTC.Equal(trunc) {
			o.nextScheduledReset = trunc
		} else {
			o.nextScheduledReset = trunc.Add(interval)
		}
		return
	}

	// Normal hour-based scheduling: align to multiples of interval since UTC midnight.
	midnight := nowUTC.Truncate(24 * time.Hour)
	elapsed := nowUTC.Sub(midnight)
	remainder := elapsed % interval
	if remainder == 0 {
		o.nextScheduledReset = nowUTC
	} else {
		o.nextScheduledReset = nowUTC.Add(interval - remainder)
	}
}

// performFullStateResetLocked clears in-memory caches and deletes the on-disk state file.
// Caller must hold o.resetMutex.
func (o *Oracle) performFullStateResetLocked() {
	// Clear oracle runtime state
	o.stateCache = []sidecartypes.OracleState{EmptyOracleState}
	o.currentState.Store(&EmptyOracleState)

	// Clear Solana signature watermarks
	o.lastSolRockMintSigStr = ""
	o.lastSolZenBTCMintSigStr = ""
	o.lastSolZenBTCBurnSigStr = ""
	o.lastSolZenZECMintSigStr = ""
	o.lastSolZenZECBurnSigStr = ""
	o.lastSolRockBurnSigStr = ""

	// Remove state file (ignore if it doesn't exist)
	if err := os.Remove(o.Config.StateFile); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: failed to remove state file during scheduled reset (%s): %v", o.Config.StateFile, err)
	}

	// Persist fresh empty state (recreates file)
	if err := o.SaveToFile(o.Config.StateFile); err != nil {
		log.Printf("Warning: failed to write fresh state file after reset (%s): %v", o.Config.StateFile, err)
	} else {
		log.Printf("State file reinitialized after scheduled reset: %s", o.Config.StateFile)
	}
}
