package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	neutrino "github.com/Zenrock-Foundation/zrchain/v5/sidecar/neutrino"
	aggregatorv3 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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
	contractInstance, err := middleware.NewContractZRServiceManager(common.HexToAddress(o.Config.EthOracle.ContractAddrs.ServiceManager), o.EthClient)
	if err != nil {
		return fmt.Errorf("failed to create contract instance: %w", err)
	}
	tempEthClient, priceFeed := o.initPriceFeed()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-o.mainLoopTicker.C:
			if err := o.fetchAndProcessState(contractInstance, tempEthClient, priceFeed); err != nil {
				log.Printf("Error fetching and processing state: %v", err)
			}
		}
	}
}

func (o *Oracle) fetchAndProcessState(contractInstance *middleware.ContractZRServiceManager, tempEthClient *ethclient.Client, priceFeed *aggregatorv3.AggregatorV3Interface) error {
	latestHeader, err := o.EthClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to fetch latest block: %w", err)
	}

	targetBlockNumber := new(big.Int).Sub(latestHeader.Number, BlocksBeforeFinality)

	delegations, err := o.getContractState(contractInstance, targetBlockNumber)
	if err != nil {
		return fmt.Errorf("failed to get contract state: %w", err)
	}

	header, err := o.EthClient.HeaderByNumber(context.Background(), targetBlockNumber)
	if err != nil {
		return fmt.Errorf("failed to fetch ETH block data: %w", err)
	}

	gasPrice, err := o.EthClient.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch gas price: %w", err)
	}

	mainnetLatestHeader, err := tempEthClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to fetch latest block: %w", err)
	}
	targetBlockNumberMainnet := new(big.Int).Sub(mainnetLatestHeader.Number, BlocksBeforeFinality)

	ETHUSDPrice, err := o.fetchEthPrice(priceFeed, targetBlockNumberMainnet)
	if err != nil {
		return fmt.Errorf("failed to fetch ETH price: %w", err)
	}

	o.updateChan <- OracleState{
		Delegations:    delegations,
		EthBlockHeight: header.Number.Uint64(),
		EthBlockHash:   header.Hash().Hex(),
		EthGasLimit:    header.GasLimit,
		EthGasPrice:    gasPrice.Uint64(),
		ETHUSDPrice:    ETHUSDPrice,
		ROCKUSDPrice:   0, // placeholder until we have a price feed for ROCK
	}

	return nil
}

func (o *Oracle) getContractState(contractInstance *middleware.ContractZRServiceManager, height *big.Int) (map[string]map[string]*big.Int, error) {
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

		// log.Printf("Operator: %s, Validator: %s, Amount: %s", operator.Hex(), validatorAddress, amount.String())

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
