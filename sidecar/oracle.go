package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	neutrino "github.com/Zenrock-Foundation/zrchain/v5/sidecar/neutrino"
	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v5/sidecar/shared"
	aggregatorv3 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	zenbtc "github.com/zenrocklabs/zenbtc/bindings"
	middleware "github.com/zenrocklabs/zenrock-avs/contracts/bindings/ZrServiceManager"

	solana "github.com/gagliardetto/solana-go/rpc"
)

func NewOracle(config Config, ethClient *ethclient.Client, neutrinoServer *neutrino.NeutrinoServer, solanaClient *solana.Client, ticker *time.Ticker) *Oracle {
	o := &Oracle{
		stateCache:     make([]sidecartypes.OracleState, 0),
		Config:         config,
		EthClient:      ethClient,
		neutrinoServer: neutrinoServer,
		solanaClient:   solanaClient,
		updateChan:     make(chan sidecartypes.OracleState, 32),
		mainLoopTicker: ticker,
	}
	o.currentState.Store(&EmptyOracleState)

	return o
}

func (o *Oracle) runAVSContractOracleLoop(ctx context.Context) error {
	serviceManager, err := middleware.NewContractZrServiceManager(common.HexToAddress(o.Config.EthOracle.ContractAddrs.ServiceManager), o.EthClient)
	if err != nil {
		return fmt.Errorf("failed to create contract instance: %w", err)
	}
	zenBTCContractHolesky, err := zenbtc.NewZenBTC(common.HexToAddress(o.Config.EthOracle.ContractAddrs.ZenBTC.EthHolesky), o.EthClient)
	if err != nil {
		return fmt.Errorf("failed to create contract instance: %w", err)
	}
	tempEthClient, btcPriceFeed, ethPriceFeed := o.initPriceFeed()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-o.mainLoopTicker.C:
			if err := o.fetchAndProcessState(serviceManager, zenBTCContractHolesky, btcPriceFeed, ethPriceFeed, tempEthClient); err != nil {
				log.Printf("Error fetching and processing state: %v", err)
			}
		}
	}
}

func (o *Oracle) fetchAndProcessState(
	serviceManager *middleware.ContractZrServiceManager,
	zenBTCContractHolesky *zenbtc.ZenBTC,
	btcPriceFeed *aggregatorv3.AggregatorV3Interface,
	ethPriceFeed *aggregatorv3.AggregatorV3Interface,
	tempEthClient *ethclient.Client,
) error {
	ctx := context.Background()

	latestHeader, err := o.EthClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch latest block: %w", err)
	}

	targetBlockNumber := new(big.Int).Sub(latestHeader.Number, EthBlocksBeforeFinality)

	eigenDelegations, err := o.getServiceManagerState(serviceManager, targetBlockNumber)
	if err != nil {
		return fmt.Errorf("failed to get contract state: %w", err)
	}

	redemptionsEthereum, err := o.getRedemptionTrackerState(zenBTCContractHolesky, targetBlockNumber)
	if err != nil {
		return fmt.Errorf("failed to get zenBTC contract state: %w", err)
	}
	// log.Printf("redemptionsEthereum: %+v\n", redemptionsEthereum)

	// TODO: get redemptions on Solana + get BTC price

	// Get base fee from latest block
	if latestHeader.BaseFee == nil {
		return fmt.Errorf("base fee not available (pre-London fork?)")
	}

	// Get suggested priority fee from client
	suggestedTip, err := o.EthClient.SuggestGasTipCap(ctx)
	if err != nil {
		return fmt.Errorf("failed to get suggested priority fee: %w", err)
	}

	// We only need 1 signature for minting, so we can use an empty message
	// Message should contain your tx setup
	// solanaFee, err := o.solanaClient.GetFeeForMessage(ctx, sol.Message{
	// 	AccountKeys:         []sol.PublicKey{},
	// 	Header:              sol.MessageHeader{},
	// 	RecentBlockhash:     sol.Hash{},
	// 	Instructions:        []sol.CompiledInstruction{},
	// 	AddressTableLookups: sol.MessageAddressTableLookupSlice{},
	// }.ToBase64(), solana.CommitmentFinalized)
	// if err != nil {
	// 	return fmt.Errorf("failed to get solana fee: %w", err)
	// }

	mainnetLatestHeader, err := tempEthClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch latest block: %w", err)
	}
	targetBlockNumberMainnet := new(big.Int).Sub(mainnetLatestHeader.Number, EthBlocksBeforeFinality)

	BTCUSDPrice, err := o.fetchPrice(btcPriceFeed, targetBlockNumberMainnet)
	if err != nil {
		return fmt.Errorf("failed to fetch BTC price: %w", err)
	}

	// ETHUSDPrice, err := o.fetchPrice(ethPriceFeed, targetBlockNumberMainnet)
	// if err != nil {
	// 	return fmt.Errorf("failed to fetch ETH price: %w", err)
	// }

	o.updateChan <- sidecartypes.OracleState{
		EigenDelegations: eigenDelegations,
		EthBlockHeight:   targetBlockNumber.Uint64(),
		EthGasLimit:      latestHeader.GasLimit,
		EthBaseFee:       latestHeader.BaseFee.Uint64(),
		EthTipCap:        suggestedTip.Uint64(),
		// SolanaLamportsPerSignature: *solanaFee.Value,
		SolanaLamportsPerSignature: 5000, // TODO: update me
		RedemptionsEthereum:        redemptionsEthereum,
		RedemptionsSolana:          nil,         // TODO: update me
		ROCKUSDPrice:               0,           // TODO: add ROCKUSDPrice after TGE
		BTCUSDPrice:                BTCUSDPrice, // TODO: update me
		ETHUSDPrice:                4000.0,      // TODO: if we need ETH price let's uncomment above block
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

func (o *Oracle) getRedemptionTrackerState(contractInstance *zenbtc.ZenBTC, height *big.Int) ([]api.Redemption, error) {
	callOpts := &bind.CallOpts{
		BlockNumber: height,
	}

	recentRedemptions, err := contractInstance.GetRecentRedemptionData(callOpts, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent redemption data: %w", err)
	}

	// convert to []Redemptions - id[i] corresponds to destinationAddr[i] and amount[i]
	redemptions := make([]api.Redemption, 0)
	for i := 0; i < len(recentRedemptions.Ids); i++ {
		redemptions = append(redemptions, api.Redemption{
			Id:                 recentRedemptions.Ids[i],
			DestinationAddress: recentRedemptions.Destination[i],
			Amount:             recentRedemptions.Amounts[i],
		})
	}

	return redemptions, nil
}
