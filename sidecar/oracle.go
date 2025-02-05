package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	neutrino "github.com/Zenrock-Foundation/zrchain/v5/sidecar/neutrino"
	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v5/sidecar/shared"
	aggregatorv3 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	zenbtc "github.com/zenrocklabs/zenbtc/bindings"
	middleware "github.com/zenrocklabs/zenrock-avs/contracts/bindings/ZrServiceManager"

	validationkeeper "github.com/Zenrock-Foundation/zrchain/v5/x/validation/keeper"
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
	zenBTCControllerHolesky, err := zenbtc.NewZenBTController(
		common.HexToAddress(o.Config.EthOracle.ContractAddrs.ZenBTC.Controller[o.Config.Network]), o.EthClient,
	)
	if err != nil {
		return fmt.Errorf("failed to create contract instance: %w", err)
	}
	tempEthClient, btcPriceFeed, ethPriceFeed := o.initPriceFeed()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-o.mainLoopTicker.C:
			if err := o.fetchAndProcessState(serviceManager, zenBTCControllerHolesky, btcPriceFeed, ethPriceFeed, tempEthClient); err != nil {
				log.Printf("Error fetching and processing state: %v", err)
			}
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

	redemptions, err := o.getRedemptions(zenBTCControllerHolesky, targetBlockNumber)
	if err != nil {
		return fmt.Errorf("failed to get zenBTC contract state: %w", err)
	}

	// Get base fee from latest block
	if latestHeader.BaseFee == nil {
		return fmt.Errorf("base fee not available (pre-London fork?)")
	}

	// Get suggested priority fee from client
	suggestedTip, err := o.EthClient.SuggestGasTipCap(ctx)
	if err != nil {
		return fmt.Errorf("failed to get suggested priority fee: %w", err)
	}

	stakeCallData, err := validationkeeper.EncodeStakeCallData(big.NewInt(1000000000))
	if err != nil {
		return fmt.Errorf("failed to encode stake call data: %w", err)
	}
	addr := common.HexToAddress(o.Config.EthOracle.ContractAddrs.ZenBTC.Controller[o.Config.Network])
	estimatedGas, err := o.EthClient.EstimateGas(context.Background(), ethereum.CallMsg{
		From: common.HexToAddress("0xE1ca337e0a0839717ef86cdA53C51b08FE681e9c"),
		To:   &addr,
		Data: stakeCallData,
	})
	if err != nil {
		return fmt.Errorf("failed to estimate gas: %w", err)
	}
	incrementedGasLimit := (estimatedGas * 110) / 100

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

	resp, err := http.Get(ROCKUSDPriceURL)
	if err != nil {
		return fmt.Errorf("failed to retrieve ROCK price data: %w", err)
	}
	defer resp.Body.Close()

	var priceData []PriceData
	if err := json.NewDecoder(resp.Body).Decode(&priceData); err != nil {
		return fmt.Errorf("failed to decode ROCK price data: %w", err)
	}
	ROCKUSDPrice, err := strconv.ParseFloat(priceData[0].Last, 64)
	if err != nil {
		return fmt.Errorf("failed to parse ROCK price data: %w", err)
	}

	mainnetLatestHeader, err := tempEthClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch latest block: %w", err)
	}
	targetBlockNumberMainnet := new(big.Int).Sub(mainnetLatestHeader.Number, EthBlocksBeforeFinality)

	BTCUSDPrice, err := o.fetchPrice(btcPriceFeed, targetBlockNumberMainnet)
	if err != nil {
		return fmt.Errorf("failed to fetch BTC price: %w", err)
	}

	ETHUSDPrice, err := o.fetchPrice(ethPriceFeed, targetBlockNumberMainnet)
	if err != nil {
		return fmt.Errorf("failed to fetch ETH price: %w", err)
	}

	// Fetch burn events from the last 100 blocks and convert them to the desired format.
	fromBlock := new(big.Int).Sub(latestHeader.Number, big.NewInt(100))
	toBlock := latestHeader.Number
	ethBurnEvents, err := o.getBurnEvents(fromBlock, toBlock)
	if err != nil {
		return fmt.Errorf("failed to get burn events: %w", err)
	}

	o.updateChan <- sidecartypes.OracleState{
		EigenDelegations: eigenDelegations,
		EthBlockHeight:   targetBlockNumber.Uint64(),
		EthGasLimit:      incrementedGasLimit, // TODO: rename to EthStakeGasLimit
		EthBaseFee:       latestHeader.BaseFee.Uint64(),
		EthTipCap:        suggestedTip.Uint64(),
		// SolanaLamportsPerSignature: *solanaFee.Value,
		SolanaLamportsPerSignature: 5000, // TODO: update me
		EthBurnEvents:              ethBurnEvents,
		Redemptions:                redemptions,
		ROCKUSDPrice:               ROCKUSDPrice,
		BTCUSDPrice:                BTCUSDPrice,
		ETHUSDPrice:                ETHUSDPrice,
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

// getBurnEvents retrieves all ZenBTCTokenRedemption (burn) events from the specified block range,
// converts them into []api.BurnEvent with correctly populated fields, and formats the chainID in CAIP-2 format.
func (o *Oracle) getBurnEvents(fromBlock, toBlock *big.Int) ([]api.BurnEvent, error) {
	ctx := context.Background()
	tokenAddress := common.HexToAddress(o.Config.EthOracle.ContractAddrs.ZenBTC.Token.Ethereum[o.Config.Network])

	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []common.Address{tokenAddress},
		Topics: [][]common.Hash{
			{crypto.Keccak256Hash([]byte("ZenBTCTokenRedemption(address,uint256,bytes,uint256)"))},
		},
	}

	logs, err := o.EthClient.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to filter logs: %w", err)
	}

	// Create a new instance of the ZenBTC token contract to parse logs
	zenBTCInstance, err := zenbtc.NewZenBTC(tokenAddress, o.EthClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZenBTC token contract instance: %w", err)
	}

	chainID, err := o.EthClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	var burnEvents []api.BurnEvent
	for _, vLog := range logs {
		event, err := zenBTCInstance.ParseTokenRedemption(vLog)
		if err != nil {
			log.Printf("failed to parse burn event log: %v", err)
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

	return burnEvents, nil
}
