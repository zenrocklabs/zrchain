package main

import (
	"context"
	"crypto/sha256"
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
	solrpc "github.com/gagliardetto/solana-go/rpc"
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
	ntpTime, err := ntp.Time("time.google.com")
	if err != nil {
		// If NTP fails at startup, log warning and proceed without alignment.
		log.Printf("Warning: Failed to fetch NTP time at startup: %v. Initial ticker alignment skipped.", err)
		ntpTime = time.Now() // Use local time as fallback for duration calculation
	}

	// Define ticker interval duration
	mainLoopTickerIntervalDuration := time.Duration(sidecartypes.MainLoopTickerIntervalSeconds) * time.Second

	// Align the start time to the nearest MainLoopTickerInterval if NTP succeeded
	if err == nil { // Only align if NTP fetch was successful
		alignedStart := ntpTime.Truncate(mainLoopTickerIntervalDuration).Add(mainLoopTickerIntervalDuration)
		initialSleep := time.Until(alignedStart)
		if initialSleep > 0 {
			log.Printf("Initial alignment: Sleeping %v until %v to start ticker.", initialSleep.Round(time.Millisecond), alignedStart.Format("15:04:05.00"))
			time.Sleep(initialSleep)
		}
	}

	mainLoopTicker := time.NewTicker(mainLoopTickerIntervalDuration)
	defer mainLoopTicker.Stop()
	o.mainLoopTicker = mainLoopTicker

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
	}

	update := &oracleStateUpdate{}
	var updateMutex sync.Mutex
	errChan := make(chan error, 16) // needs to be larger than the number of goroutines we're spawning

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
		lamportsPerSignature, err := o.getSolanaLamportsPerSignature(ctx)
		if err != nil {
			errChan <- fmt.Errorf("failed to get Solana lamports per signature: %w", err)
			return
		}
		updateMutex.Lock()
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

		var priceData []PriceData
		if err := json.NewDecoder(resp.Body).Decode(&priceData); err != nil {
			errChan <- fmt.Errorf("failed to decode ROCK price data: %w", err)
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
			errChan <- fmt.Errorf("failed to fetch latest block: %w", err)
			return
		}
		targetBlockNumberMainnet := new(big.Int).Sub(mainnetLatestHeader.Number, big.NewInt(sidecartypes.EthBlocksBeforeFinality))

		// Fetch BTC price
		btcPrice, err := o.fetchPrice(btcPriceFeed, targetBlockNumberMainnet)
		if err != nil {
			errChan <- fmt.Errorf("failed to fetch BTC price: %w", err)
			return
		}

		// Fetch ETH price
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

	// Fetch SolROCK Mints
	wg.Add(1)
	go func() {
		defer wg.Done()
		events, err := o.getSolROCKMints(sidecartypes.SolRockProgramID[o.Config.Network])
		if err != nil {
			errChan <- fmt.Errorf("failed to process SolROCK mint events: %w", err)
			return
		}
		updateMutex.Lock()
		update.SolanaMintEvents = append(update.SolanaMintEvents, events...)
		updateMutex.Unlock()
	}()

	// Fetch SolZenBTC Mints
	wg.Add(1)
	go func() {
		defer wg.Done()
		events, err := o.getSolZenBTCMints(sidecartypes.ZenBTCSolanaProgramID[o.Config.Network])
		if err != nil {
			errChan <- fmt.Errorf("failed to process SolZenBTC mint events: %w", err)
			return
		}
		updateMutex.Lock()
		update.SolanaMintEvents = append(update.SolanaMintEvents, events...)
		updateMutex.Unlock()
	}()

	// Fetch Solana burn events
	wg.Add(1)
	go func() {
		defer wg.Done()
		solanaProgramID := sidecartypes.ZenBTCSolanaProgramID[o.Config.Network]
		events, err := o.getSolanaZenBTCBurnEvents(solanaProgramID)
		if err != nil {
			errChan <- fmt.Errorf("failed to process Solana burn events: %w", err)
			return
		}
		updateMutex.Lock()
		update.solanaBurnEvents = events
		updateMutex.Unlock()
	}()

	// Fetch Solana ROCK burn events
	wg.Add(1)
	go func() {
		defer wg.Done()
		solanaProgramID := sidecartypes.SolRockProgramID[o.Config.Network]
		events, err := o.getSolanaRockBurnEvents(solanaProgramID)
		if err != nil {
			errChan <- fmt.Errorf("failed to process Solana burn events: %w", err)
			return
		}
		updateMutex.Lock()
		update.solanaBurnEvents = append(update.solanaBurnEvents, events...)
		updateMutex.Unlock()
	}()

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		if err != nil {
			return sidecartypes.OracleState{}, err
		}
	}

	// Get current state to preserve cleaned events
	currentState := o.currentState.Load().(*sidecartypes.OracleState)

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
		ROCKUSDPrice:               update.ROCKUSDPrice,
		BTCUSDPrice:                update.BTCUSDPrice,
		ETHUSDPrice:                update.ETHUSDPrice,
	}

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

func (o *Oracle) getSolROCKMints(programID string) ([]api.SolanaMintEvent, error) {
	limit := sidecartypes.SolanaEventScanTxLimit

	program, err := solana.PublicKeyFromBase58(programID)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain program public key: %w", err)
	}
	signatures, err := o.solanaClient.GetSignaturesForAddressWithOpts(context.Background(), program, &solrpc.GetSignaturesForAddressOpts{
		Limit:      &limit,
		Commitment: solrpc.CommitmentConfirmed,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get SolROCK redemptions: %w", err)
	}

	var mintEvents []api.SolanaMintEvent

	for _, signature := range signatures {
		tx, err := o.solanaClient.GetTransaction(context.Background(), signature.Signature, &solrpc.GetTransactionOpts{
			Commitment: solrpc.CommitmentConfirmed,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get SolROCK redemption transaction: %w", err)
		}

		events, err := rock_spl_token.DecodeEvents(tx, program)
		if err != nil {
			return nil, fmt.Errorf("failed to decode SolROCK redemption events: %w", err)
		}

		solTX, err := tx.Transaction.GetTransaction()
		if err != nil || len(solTX.Signatures) != 2 {
			continue
		}
		combined := append(solTX.Signatures[0][:], solTX.Signatures[1][:]...)
		sigHash := sha256.Sum256(combined)
		for _, event := range events {
			if event.Name == "TokensMintedWithFee" {
				e := event.Data.(*rock_spl_token.TokensMintedWithFeeEventData)
				mintEvents = append(mintEvents, api.SolanaMintEvent{
					SigHash:   sigHash[:],
					Date:      tx.BlockTime.Time().Unix(),
					Recipient: e.Recipient.Bytes(),
					Value:     e.Value,
					Fee:       e.Fee,
					Mint:      e.Mint.Bytes(),
				})
			}

		}
	}
	log.Printf("retrieved %d rock solana mint events", len(mintEvents))
	return mintEvents, nil
}

func (o *Oracle) getSolZenBTCMints(programID string) ([]api.SolanaMintEvent, error) {
	limit := 1000

	program, err := solana.PublicKeyFromBase58(programID)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain program public key: %w", err)
	}
	signatures, err := o.solanaClient.GetSignaturesForAddressWithOpts(context.Background(), program, &solrpc.GetSignaturesForAddressOpts{
		Limit:      &limit,
		Commitment: solrpc.CommitmentConfirmed,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get SolROCK redemptions: %w", err)
	}

	var mintEvents []api.SolanaMintEvent

	for _, signature := range signatures {
		tx, err := o.solanaClient.GetTransaction(context.Background(), signature.Signature, &solrpc.GetTransactionOpts{
			Commitment: solrpc.CommitmentConfirmed,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get SolROCK redemption transaction: %w", err)
		}

		events, err := zenbtc_spl_token.DecodeEvents(tx, program)
		if err != nil {
			return nil, fmt.Errorf("failed to decode SolROCK redemption events: %w", err)
		}

		solTX, err := tx.Transaction.GetTransaction()
		if err != nil || len(solTX.Signatures) != 2 {
			continue
		}
		combined := append(solTX.Signatures[0][:], solTX.Signatures[1][:]...)
		sigHash := sha256.Sum256(combined)
		for _, event := range events {
			if event.Name == "TokensMintedWithFee" {
				e := event.Data.(*zenbtc_spl_token.TokensMintedWithFeeEventData)
				mintEvents = append(mintEvents, api.SolanaMintEvent{
					SigHash:   sigHash[:],
					Date:      tx.BlockTime.Time().Unix(),
					Recipient: e.Recipient.Bytes(),
					Value:     e.Value,
					Fee:       e.Fee,
					Mint:      e.Mint.Bytes(),
				})
			}

		}
	}
	log.Printf("retrieved %d zenbtc solana mint events", len(mintEvents))
	return mintEvents, nil
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
	// _, _, feeCalculator, err := o.getSolanaBlockInfoAtRoundedSlot(ctx)
	// if err != nil {
	// 	return 0, err
	// }
	// return feeCalculator, nil
	return 5000, nil
}

// getSolanaZenBTCBurnEvents retrieves ZenBTC burn events from Solana.
func (o *Oracle) getSolanaZenBTCBurnEvents(programID string) ([]api.BurnEvent, error) {
	limit := sidecartypes.SolanaEventScanTxLimit

	program, err := solana.PublicKeyFromBase58(programID)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain program public key: %w", err)
	}
	signatures, err := o.solanaClient.GetSignaturesForAddressWithOpts(context.Background(), program, &solrpc.GetSignaturesForAddressOpts{
		Limit:      &limit,
		Commitment: solrpc.CommitmentConfirmed,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get Solana ZenBTC burn signatures: %w", err)
	}

	var burnEvents []api.BurnEvent

	for _, signature := range signatures {
		tx, err := o.solanaClient.GetTransaction(context.Background(), signature.Signature, &solrpc.GetTransactionOpts{
			Commitment: solrpc.CommitmentConfirmed,
		})
		if err != nil {
			// Log error and continue to next signature
			log.Printf("Failed to get Solana ZenBTC burn transaction %s: %v", signature.Signature, err)
			continue
		}

		events, err := zenbtc_spl_token.DecodeEvents(tx, program)
		if err != nil {
			// Log error and continue to next signature
			log.Printf("Failed to decode Solana ZenBTC burn events for tx %s: %v", signature.Signature, err)
			continue
		}

		// Solana CAIP-2 Identifier
		chainID := sidecartypes.SolanaCAIP2[o.Config.Network]

		for logIndex, event := range events {
			if event.Name == "TokenRedemption" {
				e := event.Data.(*zenbtc_spl_token.TokenRedemptionEventData)

				burnEvents = append(burnEvents, api.BurnEvent{
					TxID:            signature.Signature.String(), // Use transaction signature as TxID
					LogIndex:        uint64(logIndex),             // Use log index within the transaction
					ChainID:         chainID,
					DestinationAddr: e.DestAddr[:],
					Amount:          e.Value,
				})
			}
		}
	}
	return burnEvents, nil
}

// getSolanaRockBurnEvents retrieves ZenBTC burn events from Solana.
func (o *Oracle) getSolanaRockBurnEvents(programID string) ([]api.BurnEvent, error) {
	limit := sidecartypes.SolanaEventScanTxLimit

	program, err := solana.PublicKeyFromBase58(programID)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain program public key: %w", err)
	}
	signatures, err := o.solanaClient.GetSignaturesForAddressWithOpts(context.Background(), program, &solrpc.GetSignaturesForAddressOpts{
		Limit:      &limit,
		Commitment: solrpc.CommitmentConfirmed,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get Solana ZenBTC burn signatures: %w", err)
	}

	var burnEvents []api.BurnEvent

	for _, signature := range signatures {
		tx, err := o.solanaClient.GetTransaction(context.Background(), signature.Signature, &solrpc.GetTransactionOpts{
			Commitment: solrpc.CommitmentConfirmed,
		})
		if err != nil {
			// Log error and continue to next signature
			log.Printf("Failed to get Solana ZenBTC burn transaction %s: %v", signature.Signature, err)
			continue
		}

		events, err := rock_spl_token.DecodeEvents(tx, program)
		if err != nil {
			// Log error and continue to next signature
			log.Printf("Failed to decode Solana ZenBTC burn events for tx %s: %v", signature.Signature, err)
			continue
		}

		// Solana CAIP-2 Identifier
		chainID := sidecartypes.SolanaCAIP2[o.Config.Network]

		for logIndex, event := range events {
			if event.Name == "TokenRedemption" {
				e := event.Data.(*rock_spl_token.TokenRedemptionEventData)

				burnEvents = append(burnEvents, api.BurnEvent{
					TxID:            signature.Signature.String(), // Use transaction signature as TxID
					LogIndex:        uint64(logIndex),             // Use log index within the transaction
					ChainID:         chainID,
					DestinationAddr: e.DestAddr[:],
					Amount:          e.Value,
				})
			}
		}
	}
	return burnEvents, nil
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
