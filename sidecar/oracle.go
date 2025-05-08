package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"math/big"
	"net/http"
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
	aggregatorv3 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	zenbtc "github.com/zenrocklabs/zenbtc/bindings"
	middleware "github.com/zenrocklabs/zenrock-avs/contracts/bindings/ZrServiceManager"

	validationkeeper "github.com/Zenrock-Foundation/zrchain/v6/x/validation/keeper"
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
) *Oracle {
	o := &Oracle{
		stateCache:         make([]sidecartypes.OracleState, 0),
		Config:             config,
		EthClient:          ethClient,
		neutrinoServer:     neutrinoServer,
		solanaClient:       solanaClient,
		zrChainQueryClient: zrChainQueryClient,
		updateChan:         make(chan sidecartypes.OracleState, 32),
	}
	o.currentState.Store(&EmptyOracleState)

	// Load initial state from cache file
	if err := o.LoadFromFile(o.Config.StateFile); err != nil {
		log.Printf("Error loading state from file: %v", err)
	} else {
		// If LoadFromFile successfully populates o.currentState
		if loadedState, ok := o.currentState.Load().(*sidecartypes.OracleState); ok && loadedState != nil {
			o.lastSolRockMintSigStr = loadedState.LastSolRockMintSig
			o.lastSolZenBTCMintSigStr = loadedState.LastSolZenBTCMintSig
			o.lastSolZenBTCBurnSigStr = loadedState.LastSolZenBTCBurnSig
			o.lastSolRockBurnSigStr = loadedState.LastSolRockBurnSig
			log.Printf("Loaded last Solana signatures from state: RockMint=%s, ZenBTCMint=%s, ZenBTCBurn=%s, RockBurn=%s",
				o.lastSolRockMintSigStr, o.lastSolZenBTCMintSigStr, o.lastSolZenBTCBurnSigStr, o.lastSolRockBurnSigStr)
		} else {
			log.Println("Loaded state was nil or not OracleState type, initializing Solana signatures as empty.")
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

	// Initial alignment: Fetch NTP time once at startup
	// ntpTime, err := ntp.Time("time.google.com")
	// if err != nil {
	// 	// If NTP fails at startup, panic. Sidecars require time sync to establish consensus.
	// 	log.Fatalf("FATAL: Failed to fetch NTP time at startup: %v. Cannot proceed.", err)
	// }

	mainLoopTickerIntervalDuration := time.Duration(sidecartypes.MainLoopTickerIntervalSeconds) * time.Second

	// Align the start time to the nearest MainLoopTickerInterval.
	// This runs only if NTP succeeded (checked by the panic above)
	// alignedStart := ntpTime.Truncate(mainLoopTickerIntervalDuration).Add(mainLoopTickerIntervalDuration)
	// initialSleep := time.Until(alignedStart)
	// if initialSleep > 0 {
	// 	log.Printf("Initial alignment: Sleeping %v until %v to start ticker.", initialSleep.Round(time.Millisecond), alignedStart.Format("15:04:05.00"))
	// 	time.Sleep(initialSleep)
	// }

	mainLoopTicker := time.NewTicker(mainLoopTickerIntervalDuration)
	defer mainLoopTicker.Stop()
	o.mainLoopTicker = mainLoopTicker
	log.Printf("Ticker synched, awaiting initial oracle data fetch (%ds interval)...", sidecartypes.MainLoopTickerIntervalSeconds)

	for {
		select {
		case <-ctx.Done():
			return nil
		case tickTime := <-o.mainLoopTicker.C:
			newState, err := o.fetchAndProcessState(serviceManager, zenBTCControllerHolesky, btcPriceFeed, ethPriceFeed, mainnetEthClient)
			if err != nil {
				log.Printf("Error fetching and processing state: %v", err)
				continue // Skip sending update on error
			}

			// --- Intra-loop NTP check and wait (with fallback to ticker time) ---
			var sleepDuration time.Duration
			var nextIntervalMark time.Time
			alignmentSource := "NTP"

			// Attempt to fetch current NTP time *after* processing
			ntpTimeNow, err := ntp.Time("time.google.com")
			if err != nil {
				// NTP Failed: Fallback to using the captured ticker time
				log.Printf("Warning: Error fetching NTP time for alignment: %v. Falling back to ticker time.", err)
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
				log.Printf("State fetched. Waiting %v until next %s-aligned interval mark (%v) to apply update.",
					sleepDuration.Round(time.Millisecond),
					alignmentSource,
					nextIntervalMark.Format("15:04:05.00"))
				time.Sleep(sleepDuration)
			} else {
				// If fetching took longer than the interval OR NTP failed and ticker time also leads to negative sleep, log a warning.
				log.Printf("Warning: State fetching took too long relative to %s alignment. Update applied immediately.", alignmentSource)
			}
			// --- End of intra-loop wait ---

			// Send the fetched state exactly at the interval mark (or immediately if delayed)
			o.updateChan <- newState

			// Clean up burn events *after* sending state update
			o.cleanUpBurnEvents()
		}
	}
}

func (o *Oracle) fetchAndProcessState(
	serviceManager *middleware.ContractZrServiceManager,
	zenBTCControllerHolesky *zenbtc.ZenBTController,
	btcPriceFeed *aggregatorv3.AggregatorV3Interface,
	ethPriceFeed *aggregatorv3.AggregatorV3Interface,
	tempEthClient *ethclient.Client,
) (sidecartypes.OracleState, error) {
	const httpTimeout = 10 * time.Second

	ctx := context.Background()
	var wg sync.WaitGroup

	log.Printf("Retrieving latest %s header at %v", sidecartypes.NetworkNames[o.Config.Network], time.Now().Format("15:04:05.00"))
	latestHeader, err := o.EthClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return sidecartypes.OracleState{}, fmt.Errorf("failed to fetch latest block: %w", err)
	}
	log.Printf("Retrieved latest %s header (block %d) at %v", sidecartypes.NetworkNames[o.Config.Network], latestHeader.Number.Uint64(), time.Now().Format("15:04:05.00"))
	targetBlockNumber := new(big.Int).Sub(latestHeader.Number, big.NewInt(sidecartypes.EthBlocksBeforeFinality))

	// Check base fee availability
	if latestHeader.BaseFee == nil {
		return sidecartypes.OracleState{}, fmt.Errorf("base fee not available (pre-London fork?)")
	}

	type oracleStateUpdate struct {
		eigenDelegations           map[string]map[string]*big.Int
		redemptions                []api.Redemption
		suggestedTip               *big.Int
		estimatedGas               uint64
		ethBurnEvents              []api.BurnEvent
		solanaBurnEvents           []api.BurnEvent
		ROCKUSDPrice               math.LegacyDec
		BTCUSDPrice                math.LegacyDec
		ETHUSDPrice                math.LegacyDec
		solanaLamportsPerSignature uint64
		SolanaMintEvents           []api.SolanaMintEvent
		latestSolanaSigs           map[string]solana.Signature // To store newest sigs from Solana event fetchers
	}

	update := &oracleStateUpdate{
		latestSolanaSigs: make(map[string]solana.Signature),
		// Initialize slices/maps to avoid nil issues if goroutines don't return anything
		SolanaMintEvents: []api.SolanaMintEvent{},
		solanaBurnEvents: []api.BurnEvent{},
		eigenDelegations: make(map[string]map[string]*big.Int),
		redemptions:      []api.Redemption{},
		ethBurnEvents:    []api.BurnEvent{},
	}
	var updateMutex sync.Mutex
	errChan := make(chan error, 16)

	// Fetch eigen delegations
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

	// Fetch redemptions
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

	// Get suggested priority fee
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

	// Fetch Solana lamports per signature
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Use the dedicated function which now handles potential errors and fallback
		lamportsPerSignature, err := o.getSolanaLamportsPerSignature(ctx)
		if err != nil {
			// Log the error but let the default value (or potentially last known value) be used
			log.Printf("Warning: getSolanaLamportsPerSignature failed: %v. Using potentially stale/default value.", err)
			// Don't push to errChan here, as getSolanaLamportsPerSignature provides a fallback
		}
		updateMutex.Lock()
		// Only update if we got a valid value (getSolanaLamportsPerSignature returns 5000 on error)
		// Or decide if we always want to update with the returned value (including fallback)
		update.solanaLamportsPerSignature = lamportsPerSignature
		updateMutex.Unlock()
	}()

	// Estimate gas for stake call
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
		// Apply buffer (original logic)
		update.estimatedGas = (estimatedGas * 110) / 100
		updateMutex.Unlock()
	}()

	// Fetch ROCK price
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

		var priceData []PriceData // Remove sidecartypes qualifier
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

	// Fetch BTC and ETH prices
	wg.Add(1)
	go func() {
		defer wg.Done()
		mainnetLatestHeader, err := tempEthClient.HeaderByNumber(ctx, nil)
		if err != nil {
			errChan <- fmt.Errorf("failed to fetch latest mainnet block: %w", err)
			return
		}
		targetBlockNumberMainnet := new(big.Int).Sub(mainnetLatestHeader.Number, big.NewInt(sidecartypes.EthBlocksBeforeFinality))

		// Check if price feeds are initialized
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

	// Process ETH burn events
	wg.Add(1)
	go func() {
		defer wg.Done()
		events, err := o.processEthBurnEvents(latestHeader)
		if err != nil {
			errChan <- fmt.Errorf("failed to process Ethereum burn events: %w", err)
			return
		}
		updateMutex.Lock()
		update.ethBurnEvents = events
		updateMutex.Unlock()
	}()

	// Fetch SolROCK Mints using watermarking
	wg.Add(1)
	go func() {
		defer wg.Done()
		lastKnownSig := o.GetLastProcessedSolSignature("solRockMint")
		events, newestSig, err := o.getSolROCKMints(sidecartypes.SolRockProgramID[o.Config.Network], lastKnownSig)
		if err != nil {
			errChan <- fmt.Errorf("failed to process SolROCK mint events: %w", err)
			return
		}
		updateMutex.Lock()
		if len(events) > 0 {
			update.SolanaMintEvents = append(update.SolanaMintEvents, events...)
		}
		if !newestSig.IsZero() { // Update latest sig even if no new events, to advance watermark
			update.latestSolanaSigs["solRockMint"] = newestSig
		}
		updateMutex.Unlock()
	}()

	// Fetch SolZenBTC Mints using watermarking
	wg.Add(1)
	go func() {
		defer wg.Done()
		lastKnownSig := o.GetLastProcessedSolSignature("solZenBTCMint")
		events, newestSig, err := o.getSolZenBTCMints(sidecartypes.ZenBTCSolanaProgramID[o.Config.Network], lastKnownSig)
		if err != nil {
			errChan <- fmt.Errorf("failed to process SolZenBTC mint events: %w", err)
			return
		}
		updateMutex.Lock()
		if len(events) > 0 {
			update.SolanaMintEvents = append(update.SolanaMintEvents, events...)
		}
		if !newestSig.IsZero() {
			update.latestSolanaSigs["solZenBTCMint"] = newestSig
		}
		updateMutex.Unlock()
	}()

	// Fetch Solana ZenBTC burn events using watermarking
	wg.Add(1)
	go func() {
		defer wg.Done()
		lastKnownSig := o.GetLastProcessedSolSignature("solZenBTCBurn")
		events, newestSig, err := o.getSolanaZenBTCBurnEvents(sidecartypes.ZenBTCSolanaProgramID[o.Config.Network], lastKnownSig)
		if err != nil {
			errChan <- fmt.Errorf("failed to process Solana ZenBTC burn events: %w", err)
			return
		}
		updateMutex.Lock()
		if len(events) > 0 {
			update.solanaBurnEvents = append(update.solanaBurnEvents, events...)
		}
		if !newestSig.IsZero() {
			update.latestSolanaSigs["solZenBTCBurn"] = newestSig
		}
		updateMutex.Unlock()
	}()

	// Fetch Solana ROCK burn events using watermarking
	wg.Add(1)
	go func() {
		defer wg.Done()
		lastKnownSig := o.GetLastProcessedSolSignature("solRockBurn")
		events, newestSig, err := o.getSolanaRockBurnEvents(sidecartypes.SolRockProgramID[o.Config.Network], lastKnownSig)
		if err != nil {
			errChan <- fmt.Errorf("failed to process Solana ROCK burn events: %w", err)
			return
		}
		updateMutex.Lock()
		if len(events) > 0 {
			update.solanaBurnEvents = append(update.solanaBurnEvents, events...)
		}
		if !newestSig.IsZero() {
			update.latestSolanaSigs["solRockBurn"] = newestSig
		}
		updateMutex.Unlock()
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			// Log the first error encountered and return
			log.Printf("Error during state fetch: %v", err)
			return sidecartypes.OracleState{}, err
		}
	}

	// Update the main Oracle's last signature strings *before* creating newState
	// This ensures that if CacheState is called with newState, it has the latest signatures.
	if len(update.latestSolanaSigs) > 0 {
		if sig, ok := update.latestSolanaSigs["solRockMint"]; ok && !sig.IsZero() {
			o.lastSolRockMintSigStr = sig.String()
		}
		if sig, ok := update.latestSolanaSigs["solZenBTCMint"]; ok && !sig.IsZero() {
			o.lastSolZenBTCMintSigStr = sig.String()
		}
		if sig, ok := update.latestSolanaSigs["solZenBTCBurn"]; ok && !sig.IsZero() {
			o.lastSolZenBTCBurnSigStr = sig.String()
		}
		if sig, ok := update.latestSolanaSigs["solRockBurn"]; ok && !sig.IsZero() {
			o.lastSolRockBurnSigStr = sig.String()
		}
		log.Printf("Updated latest Solana signatures: RockMint=%s, ZenBTCMint=%s, ZenBTCBurn=%s, RockBurn=%s",
			o.lastSolRockMintSigStr, o.lastSolZenBTCMintSigStr, o.lastSolZenBTCBurnSigStr, o.lastSolRockBurnSigStr)
	}

	currentState := o.currentState.Load().(*sidecartypes.OracleState)

	// Ensure update fields that might not have been populated are not nil
	if update.suggestedTip == nil {
		update.suggestedTip = big.NewInt(0) // Or use last known value from currentState?
		log.Println("Warning: suggestedTip was nil, using 0.")
	}
	if update.ROCKUSDPrice.IsNil() {
		update.ROCKUSDPrice = currentState.ROCKUSDPrice // Use last known good price
		log.Println("Warning: ROCKUSDPrice was nil, using last known state value.")
	}
	if update.BTCUSDPrice.IsNil() {
		update.BTCUSDPrice = currentState.BTCUSDPrice // Use last known good price
		log.Println("Warning: BTCUSDPrice was nil, using last known state value.")
	}
	if update.ETHUSDPrice.IsNil() {
		update.ETHUSDPrice = currentState.ETHUSDPrice // Use last known good price
		log.Println("Warning: ETHUSDPrice was nil, using last known state value.")
	}
	if update.solanaLamportsPerSignature == 0 {
		update.solanaLamportsPerSignature = currentState.SolanaLamportsPerSignature // Use last known good value
		log.Println("Warning: solanaLamportsPerSignature was 0, using last known state value.")
	}

	newState := sidecartypes.OracleState{
		EigenDelegations:           update.eigenDelegations,
		EthBlockHeight:             targetBlockNumber.Uint64(),
		EthGasLimit:                update.estimatedGas,
		EthBaseFee:                 latestHeader.BaseFee.Uint64(),
		EthTipCap:                  update.suggestedTip.Uint64(),
		SolanaLamportsPerSignature: update.solanaLamportsPerSignature,
		EthBurnEvents:              update.ethBurnEvents,
		CleanedEthBurnEvents:       currentState.CleanedEthBurnEvents, // Preserve cleaned events
		SolanaBurnEvents:           update.solanaBurnEvents,
		CleanedSolanaBurnEvents:    currentState.CleanedSolanaBurnEvents, // Preserve cleaned events
		Redemptions:                update.redemptions,
		SolanaMintEvents:           update.SolanaMintEvents,
		ROCKUSDPrice:               update.ROCKUSDPrice,
		BTCUSDPrice:                update.BTCUSDPrice,
		ETHUSDPrice:                update.ETHUSDPrice,
		// Persist latest Solana signatures
		LastSolRockMintSig:   o.lastSolRockMintSigStr,
		LastSolZenBTCMintSig: o.lastSolZenBTCMintSigStr,
		LastSolZenBTCBurnSig: o.lastSolZenBTCBurnSigStr,
		LastSolRockBurnSig:   o.lastSolRockBurnSigStr,
	}

	// Store the new state into the Oracle's current state immediately.
	o.currentState.Store(&newState)
	// Optional: Cache state immediately after updating it, if persistence per cycle is desired.
	// o.CacheState() // Uncomment if CacheState should run here instead of only in cleanUpBurnEvents

	if sidecartypes.DebugMode {
		log.Printf("\nState fetched (pre-update send): %+v\n", newState)
	}

	return newState, nil
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

	quorumNumber := uint8(0)

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
				log.Printf("Failed to get stake for operator %s: %v", operator.Hex(), err)
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

func (o *Oracle) processEthBurnEvents(latestHeader *ethtypes.Header) ([]api.BurnEvent, error) {
	fromBlock := new(big.Int).Sub(latestHeader.Number, big.NewInt(int64(sidecartypes.EthBurnEventsBlockRange)))
	toBlock := latestHeader.Number
	newEthBurnEvents, err := o.getEthBurnEvents(fromBlock, toBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to get Ethereum burn events: %w", err)
	}

	// Get current state to merge with new burn events
	currentState := o.currentState.Load().(*sidecartypes.OracleState)

	// Create a map of existing events for quick lookup
	existingEthBurnEvents := make(map[string]bool)
	for _, event := range currentState.EthBurnEvents {
		key := fmt.Sprintf("%s-%s-%d", event.ChainID, event.TxID, event.LogIndex)
		existingEthBurnEvents[key] = true
	}

	// Only add new events that aren't already in our cache and haven't been cleaned up
	mergedEthBurnEvents := make([]api.BurnEvent, len(currentState.EthBurnEvents))
	copy(mergedEthBurnEvents, currentState.EthBurnEvents)
	for _, event := range newEthBurnEvents {
		key := fmt.Sprintf("%s-%s-%d", event.ChainID, event.TxID, event.LogIndex)
		if !existingEthBurnEvents[key] && !currentState.CleanedEthBurnEvents[key] {
			mergedEthBurnEvents = append(mergedEthBurnEvents, event)
		}
	}

	return mergedEthBurnEvents, nil
}

// reconcileBurnEventsWithChain checks a list of burn events against the chain and returns the events
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
			log.Printf("Skipping already cleaned %s burn event (txID: %s, logIndex: %d, chainID: %s)", chainTypeName, event.TxID, event.LogIndex, event.ChainID)
			continue // Skip to next event if already marked as cleaned
		}

		resp, err := o.zrChainQueryClient.ZenBTCQueryClient.BurnEvents(ctx, 0, event.TxID, event.LogIndex, event.ChainID)
		if err != nil || resp == nil {
			// Log the specific chain type in the error
			log.Printf("Error querying %s burn event (txID: %s, logIndex: %d, chainID: %s): %v", chainTypeName, event.TxID, event.LogIndex, event.ChainID, err)
			// Keep events that we failed to query, they might succeed next time
			remainingEvents = append(remainingEvents, event)
			continue // Continue checking other events even if one query fails
		}

		// If the event is not found on chain, keep it in our cache
		if len(resp.BurnEvents) == 0 {
			remainingEvents = append(remainingEvents, event)
		} else {
			// Event found on chain, mark it as cleaned by adding to the map
			updatedCleanedEvents[key] = true
			log.Printf("Removing %s burn event from cache as it's now on chain (txID: %s, logIndex: %d, chainID: %s)", chainTypeName, event.TxID, event.LogIndex, event.ChainID)
			// Do *not* add to remainingEvents
		}
	}

	return remainingEvents, updatedCleanedEvents
}

func (o *Oracle) cleanUpBurnEvents() {
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
		log.Printf("Removed %d Ethereum burn events from cache", initialEthCount-len(remainingEthEvents))
		stateChanged = true
	}

	// Clean up Solana events
	remainingSolEvents, updatedCleanedSolEvents := o.reconcileBurnEventsWithZRChain(ctx, currentState.SolanaBurnEvents, currentState.CleanedSolanaBurnEvents, "Solana")
	if len(remainingSolEvents) != initialSolCount {
		log.Printf("Removed %d Solana burn events from cache", initialSolCount-len(remainingSolEvents))
		stateChanged = true
	}

	// Update the current state only if changes were made to either list
	if stateChanged {
		newState := *currentState // Copy existing state
		newState.EthBurnEvents = remainingEthEvents
		newState.CleanedEthBurnEvents = updatedCleanedEthEvents
		newState.SolanaBurnEvents = remainingSolEvents
		newState.CleanedSolanaBurnEvents = updatedCleanedSolEvents

		o.currentState.Store(&newState)
		o.CacheState() // Persist the updated state
		log.Println("Burn event cache state updated and saved.")
	} else {
		log.Println("No burn events removed from cache during cleanup.")
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

		burnEvents = append(burnEvents, api.BurnEvent{
			TxID:            event.Raw.TxHash.Hex(),
			LogIndex:        uint64(event.Raw.Index),
			ChainID:         fmt.Sprintf("eip155:%s", chainID.String()),
			DestinationAddr: event.DestAddr,
			Amount:          event.Value,
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

func (o *Oracle) getSolROCKMints(programID string, lastKnownSig solana.Signature) ([]api.SolanaMintEvent, solana.Signature, error) {
	limit := sidecartypes.SolanaEventScanTxLimit
	program, err := solana.PublicKeyFromBase58(programID)
	if err != nil {
		return nil, solana.Signature{}, fmt.Errorf("failed to obtain program public key for SolROCK: %w", err)
	}

	// Fetch latest signatures for the program address
	allSignatures, err := o.solanaClient.GetSignaturesForAddressWithOpts(context.Background(), program, &solrpc.GetSignaturesForAddressOpts{
		Limit:      &limit,
		Commitment: solrpc.CommitmentConfirmed, // Or Finalized depending on requirement
	})
	if err != nil {
		// Return existing watermark on error
		return nil, lastKnownSig, fmt.Errorf("failed to get SolROCK signatures: %w", err)
	}

	if len(allSignatures) == 0 {
		log.Printf("retrieved 0 rock solana mint events (no signatures found)")
		// No signatures found at all, return existing watermark
		return []api.SolanaMintEvent{}, lastKnownSig, nil
	}

	// The newest signature from the node's perspective for this program address
	newestSigFromNode := allSignatures[0].Signature
	// The type returned by GetSignaturesForAddressWithOpts is []*rpc.TransactionSignature
	newSignaturesToFetchDetails := make([]*solrpc.TransactionSignature, 0)

	// Filter signatures: find signatures newer than the last one we processed.
	// Signatures are returned newest first.
	for _, sigInfo := range allSignatures {
		if !lastKnownSig.IsZero() && sigInfo.Signature == lastKnownSig {
			break // Found the last processed signature, stop collecting.
		}
		newSignaturesToFetchDetails = append(newSignaturesToFetchDetails, sigInfo)
	}

	if len(newSignaturesToFetchDetails) == 0 {
		// No *new* signatures since the last check.
		log.Printf("No new SolROCK mint signatures since last check (%s). Newest from node: %s", lastKnownSig, newestSigFromNode)
		// Return the newest signature seen from the node to update the watermark
		return []api.SolanaMintEvent{}, newestSigFromNode, nil
	}

	log.Printf("Found %d new SolROCK mint signatures to process since %s.", len(newSignaturesToFetchDetails), lastKnownSig)

	// Reverse the slice so we process the oldest *new* signature first.
	for i, j := 0, len(newSignaturesToFetchDetails)-1; i < j; i, j = i+1, j-1 {
		newSignaturesToFetchDetails[i], newSignaturesToFetchDetails[j] = newSignaturesToFetchDetails[j], newSignaturesToFetchDetails[i]
	}

	var mintEvents []api.SolanaMintEvent
	const internalBatchSize = 20 // Define a smaller batch size for getTransaction calls
	v0 := uint64(0)              // Define v0 for pointer

	for i := 0; i < len(newSignaturesToFetchDetails); i += internalBatchSize {
		end := i + internalBatchSize
		if end > len(newSignaturesToFetchDetails) {
			end = len(newSignaturesToFetchDetails)
		}
		currentBatchSignatures := newSignaturesToFetchDetails[i:end]

		batchRequests := make(jsonrpc.RPCRequests, 0, len(currentBatchSignatures))
		// Build the batch request for GetTransaction
		for j, sigInfo := range currentBatchSignatures {
			batchRequests = append(batchRequests, &jsonrpc.RPCRequest{
				Method: "getTransaction",
				Params: []any{
					sigInfo.Signature.String(),
					map[string]any{ // Use map for options
						"encoding":                       solana.EncodingJSON,
						"commitment":                     solrpc.CommitmentConfirmed,
						"maxSupportedTransactionVersion": &v0, // Pass pointer to 0
					},
				},
				ID:      uint64(j), // Use index within the current sub-batch for ID
				JSONRPC: "2.0",
			})
		}

		// Execute the batch request
		batchResponses, err := o.solanaClient.RPCCallBatch(context.Background(), batchRequests)
		if err != nil {
			log.Printf("SolROCK mints sub-batch GetTransaction failed (signatures %d to %d): %v. Skipping this sub-batch.", i, end-1, err)
			// Potentially return lastKnownSig here if we want to stop processing further sub-batches on error
			// return nil, lastKnownSig, fmt.Errorf("SolROCK mints sub-batch GetTransaction failed: %w", err)
			continue // Skip to the next sub-batch
		}

		// Process the results
		for _, resp := range batchResponses { // Iterate over RPCResponses
			if resp == nil { // Check if response object itself is nil
				log.Printf("Nil RPCResponse object in sub-batch response (SolROCK mint)")
				continue
			}

			// Use resp.ID for mapping and getting original signature from currentBatchSignatures
			var requestIndex int // This index is relative to currentBatchSignatures
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
					log.Printf("Failed to convert json.Number ID to int64 (SolROCK mint): %v", err)
					continue
				}
				requestIndex = int(idInt64)
			default:
				log.Printf("Invalid response ID type %T received (SolROCK mint)", resp.ID)
				continue
			}

			if requestIndex < 0 || requestIndex >= len(currentBatchSignatures) {
				log.Printf("Invalid response ID %d received for sub-batch (SolROCK mint)", requestIndex)
				continue
			}
			sig := currentBatchSignatures[requestIndex].Signature // Get sig from the current sub-batch

			if resp.Error != nil { // Check for RPC error in the response
				log.Printf("Error in sub-batch GetTransaction result for tx %s (SolROCK mint): %v", sig, resp.Error)
				continue // Skip this transaction
			}
			if resp.Result == nil { // Check if the Result field is nil
				log.Printf("Nil result field in sub-batch response for tx %s (SolROCK mint)", sig)
				continue
			}

			// Unmarshal the json.RawMessage result into GetTransactionResult
			var txResult solrpc.GetTransactionResult
			if err := json.Unmarshal(resp.Result, &txResult); err != nil {
				log.Printf("Failed to unmarshal GetTransactionResult for tx %s (SolROCK mint): %v", sig, err)
				continue
			}

			// Decode events using the result
			events, err := rock_spl_token.DecodeEvents(&txResult, program)
			if err != nil {
				log.Printf("Failed to decode SolROCK mint events for tx %s: %v", sig, err)
				continue // Skip this transaction
			}

			// Extract transaction details for SigHash calculation
			if txResult.Transaction == nil {
				log.Printf("Transaction envelope is nil in GetTransactionResult for tx %s (SolROCK mint)", sig)
				continue
			}
			solTX, err := txResult.Transaction.GetTransaction()
			if err != nil || solTX == nil {
				log.Printf("Failed to get solana.Transaction from GetTransactionResult for SolROCK sig %s: %v", sig, err)
				continue // Skip this transaction
			}

			// Original SigHash logic requiring 2 signatures
			if len(solTX.Signatures) != 2 {
				log.Printf("Expected 2 signatures for sighash calculation for SolROCK mint, got %d for sig %s. Skipping event.", len(solTX.Signatures), sig)
				continue
			}
			combined := append(solTX.Signatures[0][:], solTX.Signatures[1][:]...)
			sigHash := sha256.Sum256(combined) // Define sigHash here for this transaction

			blockTimeUnix := int64(0) // Define blockTimeUnix here for this transaction
			if txResult.BlockTime != nil {
				blockTimeUnix = txResult.BlockTime.Time().Unix()
			}

			// Append valid events from this transaction
			for _, event := range events {
				if event.Name == "TokensMintedWithFee" {
					e, ok := event.Data.(*rock_spl_token.TokensMintedWithFeeEventData)
					if !ok {
						log.Printf("Type assertion failed for SolROCK TokensMintedWithFeeEventData on tx %s", sig)
						continue
					}
					mintEvents = append(mintEvents, api.SolanaMintEvent{
						SigHash:   sigHash[:],
						Date:      blockTimeUnix,
						Recipient: e.Recipient.Bytes(),
						Value:     e.Value,
						Fee:       e.Fee,
						Mint:      e.Mint.Bytes(),
					})
				}
			}
		}
	}

	log.Printf("Retrieved %d new SolROCK mint events. Newest signature from node: %s", len(mintEvents), newestSigFromNode)
	// Return the collected events and the newest signature seen from the node to update the watermark
	return mintEvents, newestSigFromNode, nil
}

func (o *Oracle) getSolZenBTCMints(programID string, lastKnownSig solana.Signature) ([]api.SolanaMintEvent, solana.Signature, error) {
	limit := sidecartypes.SolanaEventScanTxLimit // Use constant

	program, err := solana.PublicKeyFromBase58(programID)
	if err != nil {
		return nil, solana.Signature{}, fmt.Errorf("failed to obtain program public key for SolZenBTC: %w", err)
	}

	// Fetch latest signatures for the program address
	allSignatures, err := o.solanaClient.GetSignaturesForAddressWithOpts(context.Background(), program, &solrpc.GetSignaturesForAddressOpts{
		Limit:      &limit,
		Commitment: solrpc.CommitmentConfirmed,
	})
	if err != nil {
		// Return existing watermark on error
		return nil, lastKnownSig, fmt.Errorf("failed to get SolZenBTC mint signatures: %w", err)
	}

	if len(allSignatures) == 0 {
		log.Printf("retrieved 0 zenbtc solana mint events (no signatures found)")
		// No signatures found at all, return existing watermark
		return []api.SolanaMintEvent{}, lastKnownSig, nil
	}

	// The newest signature from the node's perspective for this program address
	newestSigFromNode := allSignatures[0].Signature
	newSignaturesToFetchDetails := make([]*solrpc.TransactionSignature, 0)

	// Filter signatures: find signatures newer than the last one we processed.
	// Signatures are returned newest first.
	for _, sigInfo := range allSignatures {
		if !lastKnownSig.IsZero() && sigInfo.Signature == lastKnownSig {
			break // Found the last processed signature, stop collecting.
		}
		newSignaturesToFetchDetails = append(newSignaturesToFetchDetails, sigInfo)
	}

	if len(newSignaturesToFetchDetails) == 0 {
		// No *new* signatures since the last check.
		log.Printf("No new SolZenBTC mint signatures since last check (%s). Newest from node: %s", lastKnownSig, newestSigFromNode)
		// Return the newest signature seen from the node to update the watermark
		return []api.SolanaMintEvent{}, newestSigFromNode, nil
	}

	log.Printf("Found %d new SolZenBTC mint signatures to process since %s.", len(newSignaturesToFetchDetails), lastKnownSig)

	// Reverse the slice so we process the oldest *new* signature first.
	for i, j := 0, len(newSignaturesToFetchDetails)-1; i < j; i, j = i+1, j-1 {
		newSignaturesToFetchDetails[i], newSignaturesToFetchDetails[j] = newSignaturesToFetchDetails[j], newSignaturesToFetchDetails[i]
	}

	var mintEvents []api.SolanaMintEvent
	const internalBatchSize = 20 // Define a smaller batch size for getTransaction calls
	v0 := uint64(0)              // Define v0 for pointer for maxSupportedTransactionVersion

	for i := 0; i < len(newSignaturesToFetchDetails); i += internalBatchSize {
		end := i + internalBatchSize
		if end > len(newSignaturesToFetchDetails) {
			end = len(newSignaturesToFetchDetails)
		}
		currentBatchSignatures := newSignaturesToFetchDetails[i:end]

		batchRequests := make(jsonrpc.RPCRequests, 0, len(currentBatchSignatures))
		// Build the batch request for GetTransaction
		for j, sigInfo := range currentBatchSignatures {
			batchRequests = append(batchRequests, &jsonrpc.RPCRequest{
				Method: "getTransaction",
				Params: []any{
					sigInfo.Signature.String(),
					map[string]any{
						"encoding":                       solana.EncodingJSON,
						"commitment":                     solrpc.CommitmentConfirmed,
						"maxSupportedTransactionVersion": &v0, // Pass pointer to 0
					},
				},
				ID:      uint64(j), // Use index within the current sub-batch for ID
				JSONRPC: "2.0",
			})
		}

		// Execute the batch request
		batchResponses, err := o.solanaClient.RPCCallBatch(context.Background(), batchRequests)
		if err != nil {
			log.Printf("SolZenBTC mints sub-batch GetTransaction failed (signatures %d to %d): %v. Skipping this sub-batch.", i, end-1, err)
			continue // Skip to the next sub-batch
		}

		// Process the results
		for _, resp := range batchResponses { // Iterate over RPCResponses
			if resp == nil { // Check if response object itself is nil
				log.Printf("Nil RPCResponse object in sub-batch response (SolZenBTC mint)")
				continue
			}

			var requestIndex int // This index is relative to currentBatchSignatures
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
					log.Printf("Failed to convert json.Number ID to int64 (SolZenBTC mint): %v", err)
					continue
				}
				requestIndex = int(idInt64)
			default:
				log.Printf("Invalid response ID type %T received (SolZenBTC mint)", resp.ID)
				continue
			}

			if requestIndex < 0 || requestIndex >= len(currentBatchSignatures) {
				log.Printf("Invalid response ID %d received for sub-batch (SolZenBTC mint)", requestIndex)
				continue
			}
			sig := currentBatchSignatures[requestIndex].Signature // Get sig from the current sub-batch

			if resp.Error != nil { // Check for RPC error in the response
				log.Printf("Error in sub-batch GetTransaction result for tx %s (SolZenBTC mint): %v", sig, resp.Error)
				continue // Skip this transaction
			}
			if resp.Result == nil { // Check if the Result field is nil
				log.Printf("Nil result field in sub-batch response for tx %s (SolZenBTC mint)", sig)
				continue
			}

			// Unmarshal the json.RawMessage result into GetTransactionResult
			var txResult solrpc.GetTransactionResult
			if err := json.Unmarshal(resp.Result, &txResult); err != nil {
				log.Printf("Failed to unmarshal GetTransactionResult for tx %s (SolZenBTC mint): %v", sig, err)
				continue
			}

			// Decode events using the result
			events, err := zenbtc_spl_token.DecodeEvents(&txResult, program)
			if err != nil {
				log.Printf("Failed to decode SolZenBTC mint events for tx %s: %v", sig, err)
				continue // Skip this transaction
			}

			// Extract transaction details for SigHash calculation
			if txResult.Transaction == nil {
				log.Printf("Transaction envelope is nil in GetTransactionResult for tx %s (SolZenBTC mint)", sig)
				continue
			}
			solTX, err := txResult.Transaction.GetTransaction()
			if err != nil || solTX == nil {
				log.Printf("Failed to get solana.Transaction from GetTransactionResult for SolZenBTC sig %s: %v", sig, err)
				continue // Skip this transaction
			}

			// Sighash calculation from original logic - requires 2 signatures
			if len(solTX.Signatures) != 2 {
				log.Printf("Expected 2 signatures for sighash calculation for SolZenBTC mint, got %d for sig %s. Skipping event.", len(solTX.Signatures), sig)
				continue
			}
			combined := append(solTX.Signatures[0][:], solTX.Signatures[1][:]...)
			sigHash := sha256.Sum256(combined) // Define sigHash for this transaction

			blockTimeUnix := int64(0) // Define blockTimeUnix for this transaction
			if txResult.BlockTime != nil {
				blockTimeUnix = txResult.BlockTime.Time().Unix()
			}

			for _, event := range events {
				if event.Name == "TokensMintedWithFee" {
					e, ok := event.Data.(*zenbtc_spl_token.TokensMintedWithFeeEventData)
					if !ok {
						log.Printf("Type assertion failed for SolZenBTC TokensMintedWithFeeEventData on tx %s", sig)
						continue
					}
					mintEvents = append(mintEvents, api.SolanaMintEvent{
						SigHash:   sigHash[:],
						Date:      blockTimeUnix,
						Recipient: e.Recipient.Bytes(),
						Value:     e.Value,
						Fee:       e.Fee,
						Mint:      e.Mint.Bytes(),
					})
				}
			}
		}
	}

	log.Printf("Retrieved %d new SolZenBTC mint events. Newest signature from node: %s", len(mintEvents), newestSigFromNode)
	// Return the collected events and the newest signature seen from the node to update the watermark
	return mintEvents, newestSigFromNode, nil
}

// getSolanaRecentBlockhashWithSlot fetches a recent Solana blockhash from the block with height divisible by SolanaSlotRoundingFactor
// (i.e., a block height that's a multiple of the rounding factor) and returns both the blockhash and slot
// func (o *Oracle) getSolanaRecentBlockhash(ctx context.Context) (string, uint64, error) {
// 	blockhash, slot, _, err := o.getSolanaBlockInfoAtRoundedSlot(ctx)
// 	if err != nil {
// 		return "", 0, err
// 	}
// 	return blockhash, slot, nil
// }

// getSolanaLamportsPerSignature fetches the current lamports per signature from the Solana network
// Uses the same slot rounding logic as getSolanaRecentBlockhash for consistency
func (o *Oracle) getSolanaLamportsPerSignature(ctx context.Context) (uint64, error) {
	// Create a simple dummy transaction to estimate fees.
	// Using placeholder public keys. These don't need to exist or have funds
	// as the transaction is not actually sent, only used for fee calculation.
	dummySender := solana.MustPublicKeyFromBase58("11111111111111111111111111111111")              // A known valid public key (System Program)
	dummyReceiver := solana.MustPublicKeyFromBase58("Stake11111111111111111111111111111111111111") // Use StakeProgramID as another valid key

	// Get a recent blockhash
	recentBlockhashResult, err := o.solanaClient.GetLatestBlockhash(ctx, solrpc.CommitmentConfirmed)
	if err != nil {
		log.Printf("Failed to GetLatestBlockhash for fee calculation: %v. Returning default 5000 lamports/sig.", err)
		return 5000, fmt.Errorf("GetLatestBlockhash RPC call failed: %w", err)
	}
	if recentBlockhashResult == nil || recentBlockhashResult.Value == nil {
		log.Printf("Incomplete GetLatestBlockhash result for fee calculation. Returning default 5000 lamports/sig.")
		return 5000, fmt.Errorf("GetLatestBlockhash returned nil result or value")
	}
	recentBlockhash := recentBlockhashResult.Value.Blockhash

	// Create a new transaction builder
	txBuilder := solana.NewTransactionBuilder()

	// Add a simple transfer instruction (e.g., transfer 1 lamport)
	// The actual details of the instruction don't matter as much as its presence and size.
	transferIx := solanagoSystem.NewTransferInstruction(
		1, // 1 lamport
		dummySender,
		dummyReceiver,
	).Build()
	txBuilder.AddInstruction(transferIx)
	txBuilder.SetFeePayer(dummySender)
	txBuilder.SetRecentBlockHash(recentBlockhash)

	// The message needs to be compiled and serialized.
	// For `getFeeForMessage`, we typically don't need to sign it.
	// First, build the transaction.
	tx, err := txBuilder.Build()
	if err != nil {
		log.Printf("Failed to build transaction for fee calculation: %v. Returning default 5000 lamports/sig.", err)
		return 5000, fmt.Errorf("failed to build transaction for fee calculation: %w", err)
	}
	messageData := tx.Message // tx.Message is of type solana.Message (a struct)

	// Get the serialized message bytes using the standard MarshalBinary interface:
	serializedMessage, err := messageData.MarshalBinary()
	if err != nil {
		log.Printf("Failed to serialize message using messageData.MarshalBinary for fee calculation: %v. Returning default 5000 lamports/sig.", err)
		return 5000, fmt.Errorf("failed to serialize message using messageData.MarshalBinary: %w", err)
	}

	// Call GetFeeForMessage (expects base64 encoded message string)
	msgBase64 := base64.StdEncoding.EncodeToString(serializedMessage)
	resp, err := o.solanaClient.GetFeeForMessage(ctx, msgBase64, solrpc.CommitmentConfirmed)
	if err != nil {
		log.Printf("Failed to get Solana fees via GetFeeForMessage: %v. Returning default 5000 lamports/sig.", err)
		return 5000, fmt.Errorf("GetFeeForMessage RPC call failed: %w", err)
	}

	if resp == nil || resp.Value == nil {
		log.Printf("Incomplete fee data from Solana RPC (GetFeeForMessage response or value is nil). Returning default 5000 lamports/sig.")
		return 5000, fmt.Errorf("GetFeeForMessage returned nil response or value")
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
		log.Printf("Warning: Solana GetFeeForMessage returned 0 for LamportsPerSignature, using default 5000 if required by downstream logic, otherwise using 0.")
		// Depending on requirements, you might return 0 here, or stick to a default like 5000.
		// Let's return 0 if the network says 0, but log it. If downstream MUST have non-zero, this is an issue.
		return 0, nil // Or return 5000, nil if a non-zero value is strictly necessary downstream
	}
	return lamports, nil
}

// getSolanaZenBTCBurnEvents retrieves ZenBTC burn events from Solana.
func (o *Oracle) getSolanaZenBTCBurnEvents(programID string, lastKnownSig solana.Signature) ([]api.BurnEvent, solana.Signature, error) {
	limit := sidecartypes.SolanaEventScanTxLimit

	program, err := solana.PublicKeyFromBase58(programID)
	if err != nil {
		return nil, solana.Signature{}, fmt.Errorf("failed to obtain program public key for SolZenBTC burn: %w", err)
	}

	// Fetch latest signatures for the program address
	allSignatures, err := o.solanaClient.GetSignaturesForAddressWithOpts(context.Background(), program, &solrpc.GetSignaturesForAddressOpts{
		Limit:      &limit,
		Commitment: solrpc.CommitmentConfirmed,
	})
	if err != nil {
		// Return existing watermark on error
		return nil, lastKnownSig, fmt.Errorf("failed to get Solana ZenBTC burn signatures: %w", err)
	}

	if len(allSignatures) == 0 {
		log.Printf("retrieved 0 zenbtc solana burn events (no signatures found)")
		// No signatures found at all, return existing watermark
		return []api.BurnEvent{}, lastKnownSig, nil
	}

	// The newest signature from the node's perspective for this program address
	newestSigFromNode := allSignatures[0].Signature
	newSignaturesToFetchDetails := make([]*solrpc.TransactionSignature, 0)

	// Filter signatures: find signatures newer than the last one we processed.
	// Signatures are returned newest first.
	for _, sigInfo := range allSignatures {
		if !lastKnownSig.IsZero() && sigInfo.Signature == lastKnownSig {
			break // Found the last processed signature, stop collecting.
		}
		newSignaturesToFetchDetails = append(newSignaturesToFetchDetails, sigInfo)
	}

	if len(newSignaturesToFetchDetails) == 0 {
		// No *new* signatures since the last check.
		log.Printf("No new SolZenBTC burn signatures since last check (%s). Newest from node: %s", lastKnownSig, newestSigFromNode)
		// Return the newest signature seen from the node to update the watermark
		return []api.BurnEvent{}, newestSigFromNode, nil
	}

	log.Printf("Found %d new SolZenBTC burn signatures to process since %s.", len(newSignaturesToFetchDetails), lastKnownSig)

	// Reverse the slice so we process the oldest *new* signature first.
	for i, j := 0, len(newSignaturesToFetchDetails)-1; i < j; i, j = i+1, j-1 {
		newSignaturesToFetchDetails[i], newSignaturesToFetchDetails[j] = newSignaturesToFetchDetails[j], newSignaturesToFetchDetails[i]
	}

	var burnEvents []api.BurnEvent
	const internalBatchSize = 20 // Define a smaller batch size for getTransaction calls
	v0 := uint64(0)              // Define v0 for pointer for maxSupportedTransactionVersion

	for i := 0; i < len(newSignaturesToFetchDetails); i += internalBatchSize {
		end := i + internalBatchSize
		if end > len(newSignaturesToFetchDetails) {
			end = len(newSignaturesToFetchDetails)
		}
		currentBatchSignatures := newSignaturesToFetchDetails[i:end]

		batchRequests := make(jsonrpc.RPCRequests, 0, len(currentBatchSignatures))
		// Build the batch request for GetTransaction
		for j, sigInfo := range currentBatchSignatures {
			batchRequests = append(batchRequests, &jsonrpc.RPCRequest{
				Method: "getTransaction",
				Params: []any{
					sigInfo.Signature.String(),
					map[string]any{
						"encoding":                       solana.EncodingJSON,
						"commitment":                     solrpc.CommitmentConfirmed,
						"maxSupportedTransactionVersion": &v0, // Pass pointer to 0
					},
				},
				ID:      uint64(j), // Use index within the current sub-batch for ID
				JSONRPC: "2.0",
			})
		}

		// Execute the batch request
		batchResponses, err := o.solanaClient.RPCCallBatch(context.Background(), batchRequests)
		if err != nil {
			log.Printf("SolZenBTC burn sub-batch GetTransaction failed (signatures %d to %d): %v. Skipping this sub-batch.", i, end-1, err)
			continue // Skip to the next sub-batch
		}

		chainID := sidecartypes.SolanaCAIP2[o.Config.Network] // Moved here as it's needed per batch processing if successful

		// Process the results
		for _, resp := range batchResponses { // Iterate over RPCResponses
			if resp == nil { // Check if response object itself is nil
				log.Printf("Nil RPCResponse object in sub-batch response (SolZenBTC burn)")
				continue
			}

			var requestIndex int // This index is relative to currentBatchSignatures
			switch id := resp.ID.(type) {
			case float64:
				requestIndex = int(id)
			case int:
				requestIndex = id
			case uint64:
				requestIndex = int(id)
			case json.Number:
				idInt64, err := id.Int64()
				if err != nil {
					log.Printf("Failed to convert json.Number ID to int64 (SolZenBTC burn): %v", err)
					continue
				}
				requestIndex = int(idInt64)
			default:
				log.Printf("Invalid response ID type %T received (SolZenBTC burn)", resp.ID)
				continue
			}

			if requestIndex < 0 || requestIndex >= len(currentBatchSignatures) {
				log.Printf("Invalid response ID %d received for sub-batch (SolZenBTC burn)", requestIndex)
				continue
			}
			originalSignature := currentBatchSignatures[requestIndex].Signature // Get sig from the current sub-batch

			if resp.Error != nil { // Check for RPC error in the response
				log.Printf("Error in sub-batch GetTransaction result for tx %s (SolZenBTC burn): %v", originalSignature, resp.Error)
				continue // Skip this transaction
			}
			if resp.Result == nil { // Check if the Result field is nil
				log.Printf("Nil result field in sub-batch response for tx %s (SolZenBTC burn)", originalSignature)
				continue
			}

			// Unmarshal the json.RawMessage result into GetTransactionResult
			var txResult solrpc.GetTransactionResult
			if err := json.Unmarshal(resp.Result, &txResult); err != nil {
				log.Printf("Failed to unmarshal GetTransactionResult for tx %s (SolZenBTC burn): %v", originalSignature, err)
				continue
			}

			// Decode events using the result
			events, err := zenbtc_spl_token.DecodeEvents(&txResult, program)
			if err != nil {
				log.Printf("Failed to decode Solana ZenBTC burn events for tx %s: %v", originalSignature, err)
				continue // Skip this transaction
			}

			// Process burn events for this transaction
			for logIndex, event := range events {
				if event.Name == "TokenRedemption" {
					e, ok := event.Data.(*zenbtc_spl_token.TokenRedemptionEventData)
					if !ok {
						log.Printf("Type assertion failed for SolZenBTC TokenRedemptionEventData on tx %s", originalSignature)
						continue
					}
					burnEvents = append(burnEvents, api.BurnEvent{
						TxID:            originalSignature.String(),
						LogIndex:        uint64(logIndex),
						ChainID:         chainID,
						DestinationAddr: e.DestAddr[:],
						Amount:          e.Value,
					})
				}
			}
		}
	}

	log.Printf("Retrieved %d new SolZenBTC burn events. Newest signature from node: %s", len(burnEvents), newestSigFromNode)
	// Return the collected events and the newest signature seen from the node to update the watermark
	return burnEvents, newestSigFromNode, nil
}

// getSolanaRockBurnEvents retrieves Rock burn events from Solana.
func (o *Oracle) getSolanaRockBurnEvents(programID string, lastKnownSig solana.Signature) ([]api.BurnEvent, solana.Signature, error) {
	limit := sidecartypes.SolanaEventScanTxLimit
	program, err := solana.PublicKeyFromBase58(programID)
	if err != nil {
		return nil, solana.Signature{}, fmt.Errorf("failed to obtain program public key for SolRock burn: %w", err)
	}

	// Fetch latest signatures
	allSignatures, err := o.solanaClient.GetSignaturesForAddressWithOpts(context.Background(), program, &solrpc.GetSignaturesForAddressOpts{
		Limit:      &limit,
		Commitment: solrpc.CommitmentConfirmed,
	})
	if err != nil {
		return nil, lastKnownSig, fmt.Errorf("failed to get SolRock burn signatures: %w", err)
	}

	if len(allSignatures) == 0 {
		return []api.BurnEvent{}, lastKnownSig, nil
	}

	newestSigFromNode := allSignatures[0].Signature
	newSignaturesToFetchDetails := make([]*solrpc.TransactionSignature, 0)

	// Filter for new signatures
	for _, sigInfo := range allSignatures {
		if !lastKnownSig.IsZero() && sigInfo.Signature == lastKnownSig {
			break
		}
		// Skip errored signatures if needed (removed based on linter feedback in mint function)
		// if sigInfo.Error != nil { ... }
		newSignaturesToFetchDetails = append(newSignaturesToFetchDetails, sigInfo)
	}

	if len(newSignaturesToFetchDetails) == 0 {
		log.Printf("No new SolRock burn signatures since last check (%s). Newest from node: %s", lastKnownSig, newestSigFromNode)
		return []api.BurnEvent{}, newestSigFromNode, nil
	}

	log.Printf("Found %d new SolRock burn signatures to process since %s.", len(newSignaturesToFetchDetails), lastKnownSig)

	// Reverse to process oldest first
	for i, j := 0, len(newSignaturesToFetchDetails)-1; i < j; i, j = i+1, j-1 {
		newSignaturesToFetchDetails[i], newSignaturesToFetchDetails[j] = newSignaturesToFetchDetails[j], newSignaturesToFetchDetails[i]
	}

	var burnEvents []api.BurnEvent
	const internalBatchSize = 20 // Define a smaller batch size for getTransaction calls
	v0 := uint64(0)              // Define v0 for pointer

	for i := 0; i < len(newSignaturesToFetchDetails); i += internalBatchSize {
		end := i + internalBatchSize
		if end > len(newSignaturesToFetchDetails) {
			end = len(newSignaturesToFetchDetails)
		}
		currentBatchSignatures := newSignaturesToFetchDetails[i:end]

		batchRequests := make(jsonrpc.RPCRequests, 0, len(currentBatchSignatures))
		// Build batch request
		for j, sigInfo := range currentBatchSignatures {
			batchRequests = append(batchRequests, &jsonrpc.RPCRequest{
				Method: "getTransaction",
				Params: []any{
					sigInfo.Signature.String(),
					map[string]any{ // Use map for options
						"encoding":                       solana.EncodingJSON,
						"commitment":                     solrpc.CommitmentConfirmed,
						"maxSupportedTransactionVersion": &v0, // Pass pointer to 0
					},
				},
				ID:      uint64(j), // Use index within the current sub-batch for ID
				JSONRPC: "2.0",
			})
		}

		// Execute batch request
		batchResponses, err := o.solanaClient.RPCCallBatch(context.Background(), batchRequests)
		if err != nil {
			log.Printf("SolRock burn sub-batch GetTransaction failed (signatures %d to %d): %v. Skipping this sub-batch.", i, end-1, err)
			continue // Skip to the next sub-batch
		}

		chainID := sidecartypes.SolanaCAIP2[o.Config.Network] // Moved here

		// Process results
		for _, resp := range batchResponses {
			if resp == nil {
				log.Printf("Nil RPCResponse object in sub-batch response (SolRock burn)")
				continue
			}

			var requestIndex int // This index is relative to currentBatchSignatures
			switch id := resp.ID.(type) {
			case float64:
				requestIndex = int(id)
			case int:
				requestIndex = id
			case uint64:
				requestIndex = int(id)
			case json.Number:
				idInt64, err := id.Int64()
				if err != nil {
					log.Printf("Failed to convert json.Number ID to int64 (SolRock burn): %v", err)
					continue
				}
				requestIndex = int(idInt64)
			default:
				log.Printf("Invalid response ID type %T received (SolRock burn)", resp.ID)
				continue
			}

			if requestIndex < 0 || requestIndex >= len(currentBatchSignatures) {
				log.Printf("Invalid response ID %d received for sub-batch (SolRock burn)", requestIndex)
				continue
			}
			originalSignature := currentBatchSignatures[requestIndex].Signature // Get sig from the current sub-batch

			if resp.Error != nil {
				log.Printf("Error in sub-batch GetTransaction result for tx %s (SolRock burn): %v", originalSignature, resp.Error)
				continue
			}
			if resp.Result == nil {
				log.Printf("Nil result field in sub-batch response for tx %s (SolRock burn)", originalSignature)
				continue
			}

			var txResult solrpc.GetTransactionResult
			if err := json.Unmarshal(resp.Result, &txResult); err != nil {
				log.Printf("Failed to unmarshal GetTransactionResult for tx %s (SolRock burn): %v", originalSignature, err)
				continue
			}

			// Decode events
			events, err := rock_spl_token.DecodeEvents(&txResult, program)
			if err != nil {
				log.Printf("Failed to decode Solana Rock burn events for tx %s: %v", originalSignature, err)
				continue
			}

			// Process burn events
			for logIndex, event := range events {
				if event.Name == "TokenRedemption" {
					e, ok := event.Data.(*rock_spl_token.TokenRedemptionEventData)
					if !ok {
						log.Printf("Type assertion failed for SolRock TokenRedemptionEventData on tx %s", originalSignature)
						continue
					}
					burnEvents = append(burnEvents, api.BurnEvent{
						TxID:            originalSignature.String(),
						LogIndex:        uint64(logIndex),
						ChainID:         chainID,
						DestinationAddr: e.DestAddr[:],
						Amount:          e.Value,
					})
				}
			}
		}
	}

	log.Printf("Retrieved %d new SolRock burn events. Newest signature from node: %s", len(burnEvents), newestSigFromNode)
	return burnEvents, newestSigFromNode, nil
}

// getSolanaBlockInfoAtRoundedSlot gets Solana block information from a slot divisible by SolanaSlotRoundingFactor
// Returns blockhash string, slot number, and lamports per signature
// func (o *Oracle) getSolanaBlockInfoAtRoundedSlot(ctx context.Context) (string, uint64, uint64, error) {
//// Get the latest block height
//resp, err := o.solanaClient.GetLatestBlockhash(ctx, solrpc.CommitmentFinalized)
//if err != nil {
//	return "", 0, 0, fmt.Errorf("failed to GetLatestBlockhash: %w", err)
//}
//
//dummySender := solana.MustPublicKeyFromBase58("11111111111111111111111111111111") // System Program ID (valid placeholder)
//dummyReceiver := solana.MustPublicKeyFromBase58("11111111111111111111111111111111")
//recentBlockhash := resp.Value.Blockhash
//// Create a transaction
//tx := solana.NewTransactionBuilder()
//
//// Add a transfer instruction
//transferIx := solprogram.Transfer{
//	FromPubkey: fromPubKey,
//	ToPubkey:   toPubKey,
//	Lamports:   1000, // Transfer 1000 lamports
//}
//
//tx.AddInstruction(transferIx)
//
//respFee, err := o.solanaClient.GetFeeForMessage(ctx, serialized, solrpc.CommitmentFinalized)
//if err != nil {
//	return "", 0, 0, fmt.Errorf("failed to GetFeeForMessage: %w", err)
//}
//// Default values from the recent block
//lamportsPerSignature := *respFee.Value
//
//// Get the slot for the recent blockhash
//slot, err := o.solanaClient.GetSlot(ctx, solrpc.CommitmentFinalized)
//if err != nil {
//	return recentBlockhash.String(), slot, lamportsPerSignature, fmt.Errorf("failed to get current slot: %w", err)
//}
//
//// Calculate the nearest slot that is divisible by the rounding factor
//targetSlot := slot - (slot % sidecartypes.SolanaSlotRoundingFactor)
//
//// If we're at slot 0, use the current slot's blockhash
//if targetSlot == 0 {
//	return recentBlockhash.String(), slot, lamportsPerSignature, nil
//}
//
//// Get the blockhash for the target slot
//blockInfo, err := o.solanaClient.GetBlock(ctx, targetSlot)
//if err != nil {
//	// Fallback to the recent blockhash if we can't get the target block
//	log.Printf("Failed to get block at slot %d, using recent blockhash: %v", targetSlot, err)
//	return recentBlockhash.String(), slot, lamportsPerSignature, nil
//}
//
//return blockInfo.Blockhash.String(), targetSlot, lamportsPerSignature, nil

// 	return "", 0, 0, nil
// }

// Helper to get typed last processed Solana signature
func (o *Oracle) GetLastProcessedSolSignature(eventType string) solana.Signature {
	var sigStr string
	switch eventType {
	case "solRockMint":
		sigStr = o.lastSolRockMintSigStr
	case "solZenBTCMint":
		sigStr = o.lastSolZenBTCMintSigStr
	case "solZenBTCBurn":
		sigStr = o.lastSolZenBTCBurnSigStr
	case "solRockBurn":
		sigStr = o.lastSolRockBurnSigStr
	default:
		log.Printf("Warning: Unknown Solana event type for GetLastProcessedSolSignature: %s", eventType)
		return solana.Signature{} // Return zero signature for unknown types
	}

	if sigStr == "" {
		return solana.Signature{} // No signature stored yet, so zero value
	}
	sig, err := solana.SignatureFromBase58(sigStr)
	if err != nil {
		// Log the error but return a zero signature to proceed as if no prior sig exists
		log.Printf("Warning: could not parse stored signature string '%s' for event type %s: %v. Treating as no prior signature.", sigStr, eventType, err)
		return solana.Signature{}
	}
	return sig
}
