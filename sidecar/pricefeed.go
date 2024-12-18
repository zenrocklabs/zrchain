package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	aggregatorv3 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"
)

func (o *Oracle) initPriceFeed() (*ethclient.Client, *aggregatorv3.AggregatorV3Interface, *aggregatorv3.AggregatorV3Interface) {
	// once we switch the client in the main func to use mainnet, use it instead of this
	tempEthClient, err := ethclient.Dial(o.Config.EthOracle.RPC["mainnet"])
	if err != nil {
		log.Fatalf("failed to connect to the Ethereum client: %v", err)
	}

	btcPriceFeedAddr := common.HexToAddress(o.Config.EthOracle.ContractAddrs.PriceFeeds.BTC)
	btcPriceFeed, err := aggregatorv3.NewAggregatorV3Interface(btcPriceFeedAddr, tempEthClient) // use tempEthClient for now
	if err != nil {
		log.Fatalf("Failed to create Chainlink price feed instance: %v", err)
	}

	ethPriceFeedAddr := common.HexToAddress(o.Config.EthOracle.ContractAddrs.PriceFeeds.ETH)
	// ethPriceFeed, err := aggregatorv3.NewAggregatorV3Interface(priceFeedAddr, o.ethClient)
	ethPriceFeed, err := aggregatorv3.NewAggregatorV3Interface(ethPriceFeedAddr, tempEthClient) // use tempEthClient for now
	if err != nil {
		log.Fatalf("Failed to create Chainlink price feed instance: %v", err)
	}

	return tempEthClient, btcPriceFeed, ethPriceFeed
}

func (o *Oracle) fetchPrice(priceFeed *aggregatorv3.AggregatorV3Interface, blockNumber *big.Int) (float64, error) {
	callOpts := &bind.CallOpts{BlockNumber: blockNumber}

	roundData, err := priceFeed.LatestRoundData(callOpts)
	if err != nil {
		return 0, fmt.Errorf("error fetching latest round data: %v", err)
	}

	decimals, err := priceFeed.Decimals(callOpts)
	if err != nil {
		return 0, fmt.Errorf("error fetching decimals: %v", err)
	}

	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	price := new(big.Float).SetInt(roundData.Answer)
	divisorFloat := new(big.Float).SetInt(divisor)
	result := new(big.Float).Quo(price, divisorFloat)
	priceFloat, _ := result.Float64()

	return priceFloat, nil
}
