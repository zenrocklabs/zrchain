package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	neutrino "github.com/Zenrock-Foundation/zrchain/v5/sidecar/neutrino"
	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
	aggregatorv3 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	middleware "github.com/zenrocklabs/zenrock-avs/contracts/bindings/ZrServiceManager"
	zenbtc "github.com/zenrocklabs/zenbtc/bindings"

	solana "github.com/gagliardetto/solana-go/rpc"
)

func NewOracle(config Config, ethClient *ethclient.Client, neutrinoServer *neutrino.NeutrinoServer, solanaClient *solana.Client, ticker *time.Ticker) *Oracle {
	o := &Oracle{
		stateCache:     make([]OracleState, 0),
		Config:         config,
		EthClient:      ethClient,
		neutrinoServer: neutrinoServer,
		solanaClient:   solanaClient,
		updateChan:     make(chan OracleState, 32),
		mainLoopTicker: ticker,
	}

	// Initialize the current state
	initialState := &OracleState{
		Delegations: make(map[string]map[string]*big.Int),
	}
	o.currentState.Store(initialState)

	return o
}

func (o *Oracle) runAVSContractOracleLoop(ctx context.Context) error {
	contractInstance, err := middleware.NewContractZrServiceManager(common.HexToAddress(o.Config.EthOracle.ContractAddrs.ServiceManager), o.EthClient)
	if err != nil {
		return fmt.Errorf("failed to create contract instance: %w", err)
	}
	redemptionTrackerHolesky, err := zenbtc.NewRedemptionTracker(common.HexToAddress(o.Config.EthOracle.ContractAddrs.RedemptionTrackers.EthHolesky), o.EthClient)
	if err != nil {
		return fmt.Errorf("failed to create contract instance: %w", err)
	}
	tempEthClient, priceFeed := o.initPriceFeed()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-o.mainLoopTicker.C:
			if err := o.fetchAndProcessState(serviceManager, redemptionTrackerHolesky, priceFeed, tempEthClient); err != nil {
				log.Printf("Error fetching and processing state: %v", err)
			}
		}
	}
}

func (o *Oracle) fetchAndProcessState(contractInstance *middleware.ContractZrServiceManager, tempEthClient *ethclient.Client, priceFeed *aggregatorv3.AggregatorV3Interface) error {
	ctx := context.Background()

	latestHeader, err := o.EthClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch latest block: %w", err)
	}

	targetBlockNumber := new(big.Int).Sub(latestHeader.Number, BlocksBeforeFinality)

	delegations, err := o.getServiceManagerState(serviceManager, targetBlockNumber)
	if err != nil {
		return fmt.Errorf("failed to get contract state: %w", err)
	}

	RedemptionsEthereum, err := o.getRedemptionTrackerState(redemptionTrackerHolesky, targetBlockNumber)
	if err != nil {
		return fmt.Errorf("failed to get redemption tracker state: %w", err)
	}

	header, err := o.EthClient.HeaderByNumber(ctx, targetBlockNumber)
	if err != nil {
		return fmt.Errorf("failed to fetch ETH block data: %w", err)
	}

	// Get base fee from latest block
	if header.BaseFee == nil {
		return fmt.Errorf("base fee not available (pre-London fork?)")
	}

	// Get suggested priority fee from client
	suggestedTip, err := o.EthClient.SuggestGasTipCap(ctx)
	if err != nil {
		return fmt.Errorf("failed to get suggested priority fee: %w", err)
	}

	mainnetLatestHeader, err := tempEthClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch latest block: %w", err)
	}
	targetBlockNumberMainnet := new(big.Int).Sub(mainnetLatestHeader.Number, BlocksBeforeFinality)

	ETHUSDPrice, err := o.fetchEthPrice(priceFeed, targetBlockNumberMainnet)
	if err != nil {
		return fmt.Errorf("failed to fetch ETH price: %w", err)
	}

	o.updateChan <- OracleState{
		Delegations:         delegations,
		EthBlockHeight:      header.Number.Uint64(),
		EthBlockHash:        header.Hash().Hex(),
		EthGasLimit:         header.GasLimit,
		EthBaseFee:          header.BaseFee.Uint64(),
		EthTipCap:           suggestedTip.Uint64(),
		ETHUSDPrice:         ETHUSDPrice,
		ROCKUSDPrice:        0, // TODO: add ROCKUSDPrice after TGE
		RedemptionsEthereum: RedemptionsEthereum,
	}

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

func (o *Oracle) getRedemptionTrackerState(contractInstance *zenbtc.RedemptionTracker, height *big.Int) ([]api.Redemption, error) {
	callOpts := &bind.CallOpts{
		BlockNumber: height,
	}

	recentRedemptions, err := contractInstance.GetRecentRedemptions(callOpts, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get redemptions: %w", err)
	}

	// convert to []Redemptions - id[i] corresponds to destinationAddr[i] and amount[i]
	redemptions := make([]api.Redemption, 0)
	for i := 0; i < len(recentRedemptions.Ids); i++ {
		redemptions = append(redemptions, api.Redemption{
			Id:                 recentRedemptions.Ids[i],
			DestinationAddress: recentRedemptions.DestinationAddrs[i],
			Amount:             recentRedemptions.Amounts[i],
		})
	}

	return redemptions, nil
}
