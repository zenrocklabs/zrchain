package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sync"
	"time"

	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v5/go-client"
	neutrino "github.com/Zenrock-Foundation/zrchain/v5/sidecar/neutrino"
	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v5/sidecar/shared"
	aggregatorv3 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	zenbtc "github.com/zenrocklabs/zenbtc/bindings"
	middleware "github.com/zenrocklabs/zenrock-avs/contracts/bindings/ZrServiceManager"

	validationkeeper "github.com/Zenrock-Foundation/zrchain/v5/x/validation/keeper"
	solana "github.com/gagliardetto/solana-go/rpc"
)

func NewOracle(
	config sidecartypes.Config,
	ethClient *ethclient.Client,
	neutrinoServer *neutrino.NeutrinoServer,
	solanaClient *solana.Client,
	zrChainQueryClient *client.QueryClient,
	ticker *time.Ticker,
) *Oracle {
	o := &Oracle{
		stateCache:         make([]sidecartypes.OracleState, 0),
		Config:             config,
		EthClient:          ethClient,
		neutrinoServer:     neutrinoServer,
		solanaClient:       solanaClient,
		zrChainQueryClient: zrChainQueryClient,
		updateChan:         make(chan sidecartypes.OracleState, 32),
		mainLoopTicker:     ticker,
	}
	o.currentState.Store(&EmptyOracleState)

	// Load initial state from cache file
	if err := o.LoadFromFile(o.Config.StateFile); err != nil {
		log.Printf("Error loading state from file: %v", err)
	}

	return o
}

func (o *Oracle) runAVSContractOracleLoop(ctx context.Context) error {
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

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-o.mainLoopTicker.C:
			if err := o.fetchAndProcessState(serviceManager, zenBTCControllerHolesky, btcPriceFeed, ethPriceFeed, mainnetEthClient); err != nil {
				log.Printf("Error fetching and processing state: %v", err)
			}
			o.cleanUpEthBurnEvents()
		}
	}
}

func (o *Oracle) fetchAndProcessState(
	serviceManager *middleware.ContractZrServiceManager,
	zenBTCControllerHolesky *zenbtc.ZenBTController,
	btcPriceFeed *aggregatorv3.AggregatorV3Interface,
	ethPriceFeed *aggregatorv3.AggregatorV3Interface,
	tempEthClient *ethclient.Client,
) error {
	const httpTimeout = 10 * time.Second

	ctx := context.Background()
	var wg sync.WaitGroup

	latestHeader, err := o.EthClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch latest block: %w", err)
	}

	targetBlockNumber := new(big.Int).Sub(latestHeader.Number, EthBlocksBeforeFinality)

	// Check base fee availability
	if latestHeader.BaseFee == nil {
		return fmt.Errorf("base fee not available (pre-London fork?)")
	}

	type oracleStateUpdate struct {
		eigenDelegations map[string]map[string]*big.Int
		redemptions      []api.Redemption
		suggestedTip     *big.Int
		estimatedGas     uint64
		ethBurnEvents    []api.BurnEvent
		ROCKUSDPrice     math.LegacyDec
		BTCUSDPrice      math.LegacyDec
		ETHUSDPrice      math.LegacyDec
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
		resp, err := client.Get(ROCKUSDPriceURL)
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
		targetBlockNumberMainnet := new(big.Int).Sub(mainnetLatestHeader.Number, EthBlocksBeforeFinality)

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

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		if err != nil {
			return err
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
		SolanaLamportsPerSignature: 5000, // TODO: update me
		EthBurnEvents:              update.ethBurnEvents,
		CleanedEthBurnEvents:       currentState.CleanedEthBurnEvents,
		Redemptions:                update.redemptions,
		ROCKUSDPrice:               update.ROCKUSDPrice,
		BTCUSDPrice:                update.BTCUSDPrice,
		ETHUSDPrice:                update.ETHUSDPrice,
	}

	log.Printf("\nState update: %+v\n", newState)

	o.updateChan <- newState

	return nil
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
	fromBlock := new(big.Int).Sub(latestHeader.Number, big.NewInt(EthBurnEventsBlockRange))
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

func (o *Oracle) cleanUpEthBurnEvents() {
	currentState := o.currentState.Load().(*sidecartypes.OracleState)
	if len(currentState.EthBurnEvents) == 0 {
		return
	}

	ctx := context.Background()
	remainingEthBurnEvents := make([]api.BurnEvent, 0)

	// Check each Ethereum burn event against the chain
	for _, event := range currentState.EthBurnEvents {
		resp, err := o.zrChainQueryClient.ZenBTCQueryClient.BurnEvents(ctx, 0, event.TxID, event.LogIndex, event.ChainID)
		if err != nil {
			log.Printf("Error querying Ethereum burn event (txID: %s, logIndex: %d): %v", event.TxID, event.LogIndex, err)
			// Keep events that we failed to query
			remainingEthBurnEvents = append(remainingEthBurnEvents, event)
			continue
		}

		// If the event is not found on chain, keep it in our cache
		if len(resp.BurnEvents) == 0 {
			remainingEthBurnEvents = append(remainingEthBurnEvents, event)
		} else {
			key := fmt.Sprintf("%s-%s-%d", event.ChainID, event.TxID, event.LogIndex)
			if currentState.CleanedEthBurnEvents == nil {
				currentState.CleanedEthBurnEvents = make(map[string]bool)
			}
			currentState.CleanedEthBurnEvents[key] = true
			log.Printf("Removing Ethereum burn event from cache as it's now on chain (txID: %s, logIndex: %d, chainID: %s)", event.TxID, event.LogIndex, event.ChainID)
		}
	}

	// Update the current state with remaining events if any were removed
	if len(remainingEthBurnEvents) != len(currentState.EthBurnEvents) {
		log.Printf("Removed %d Ethereum burn events from cache", len(currentState.EthBurnEvents)-len(remainingEthBurnEvents))
		newState := *currentState
		newState.EthBurnEvents = remainingEthBurnEvents
		o.currentState.Store(&newState)
		o.CacheState()
	}
}

// getBurnEvents retrieves all ZenBTCTokenRedemption (burn) events from the specified block range,
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

// getSolanaRecentBlockhash is a wrapper around getSolanaRecentBlockhashWithSlot that returns only the blockhash
func (o *Oracle) getSolanaRecentBlockhash(ctx context.Context) (string, error) {
	blockhash, _, err := o.getSolanaRecentBlockhashWithSlot(ctx)
	return blockhash, err
}

// getSolanaRecentBlockhashWithSlot fetches a recent Solana blockhash from the block with height divisible by 50
// (i.e., a block height ending in 00 or 50) and returns both the blockhash and slot
func (o *Oracle) getSolanaRecentBlockhashWithSlot(ctx context.Context) (string, uint64, error) {
	// Get the latest block height
	resp, err := o.solanaClient.GetRecentBlockhash(ctx, solana.CommitmentFinalized)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get recent blockhash: %w", err)
	}

	// Get the slot for the recent blockhash
	slot, err := o.solanaClient.GetSlot(ctx, solana.CommitmentFinalized)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get current slot: %w", err)
	}

	// Calculate the nearest slot that is divisible by 50
	targetSlot := slot - (slot % 50)

	// If we're at slot 0, use the current slot's blockhash
	if targetSlot == 0 {
		return resp.Value.Blockhash.String(), slot, nil
	}

	// Get the blockhash for the target slot
	blockInfo, err := o.solanaClient.GetBlock(ctx, targetSlot)
	if err != nil {
		// Fallback to the recent blockhash if we can't get the target block
		log.Printf("Failed to get block at slot %d, using recent blockhash: %v", targetSlot, err)
		return resp.Value.Blockhash.String(), slot, nil
	}

	return blockInfo.Blockhash.String(), targetSlot, nil
}
