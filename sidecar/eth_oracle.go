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
	zenbtc "github.com/zenrocklabs/zenbtc/bindings"
	middleware "github.com/zenrocklabs/zenrock-avs/contracts/bindings/ZRServiceManager"

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
	serviceManager, err := middleware.NewContractZRServiceManager(common.HexToAddress(o.Config.EthOracle.ContractAddrs.ServiceManager), o.EthClient)
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

func (o *Oracle) fetchAndProcessState(serviceManager *middleware.ContractZRServiceManager, redemptionTrackerHolesky *zenbtc.RedemptionTracker, priceFeed *aggregatorv3.AggregatorV3Interface, tempEthClient *ethclient.Client) error {
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
	log.Printf("Delegations: %v\n", delegations)

	RedemptionsEthereum, err := o.getRedemptionTrackerState(redemptionTrackerHolesky, targetBlockNumber)
	if err != nil {
		// return fmt.Errorf("failed to get redemption tracker state: %w", err) // TODO: uncomment when RedemptionTracker is deployed
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

func (o *Oracle) getServiceManagerState(contractInstance *middleware.ContractZRServiceManager, height *big.Int) (map[string]map[string]*big.Int, error) {
	delegations := make(map[string]map[string]*big.Int)

	callOpts := &bind.CallOpts{
		BlockNumber: height,
	}

	// Retrieve all operators from the contract
	allOperators, err := contractInstance.GetAllOperators(callOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to get all operators: %w", err)
	}

	// Iterate over all operators
	for _, operator := range allOperators {
		// Fetch delegation details for the operator
		validatorAddress, amount, err := contractInstance.GetDelegationsForOperator(callOpts, operator)
		if err != nil {
			log.Printf("Failed to get delegation for operator %s: %v", operator.Hex(), err)
			continue
		}

		// Only consider positive delegation amounts
		if amount.Cmp(big.NewInt(0)) > 0 {
			if delegations[validatorAddress] == nil {
				delegations[validatorAddress] = make(map[string]*big.Int)
			}
			delegations[validatorAddress][operator.Hex()] = amount
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
