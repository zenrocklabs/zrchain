package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"cosmossdk.io/math"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	aggregatorv3 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"
)

func (o *Oracle) initPriceFeed() (*ethclient.Client, *aggregatorv3.AggregatorV3Interface, *aggregatorv3.AggregatorV3Interface) {
	mainnetEthClient, err := ethclient.Dial(o.Config.EthRPC[sidecartypes.NetworkMainnet])
	if err != nil {
		log.Fatalf("failed to connect to the Ethereum client: %v", err)
	}

	btcPriceFeedAddr := common.HexToAddress(sidecartypes.PriceFeedAddresses.BTC)
	btcPriceFeed, err := aggregatorv3.NewAggregatorV3Interface(btcPriceFeedAddr, mainnetEthClient)
	if err != nil {
		log.Fatalf("Failed to create BTC Chainlink price feed instance: %v", err)
	}

	ethPriceFeedAddr := common.HexToAddress(sidecartypes.PriceFeedAddresses.ETH)
	ethPriceFeed, err := aggregatorv3.NewAggregatorV3Interface(ethPriceFeedAddr, mainnetEthClient)
	if err != nil {
		log.Fatalf("Failed to create ETH Chainlink price feed instance: %v", err)
	}

	return mainnetEthClient, btcPriceFeed, ethPriceFeed
}

func (o *Oracle) fetchPrice(priceFeed *aggregatorv3.AggregatorV3Interface, blockNumber *big.Int) (math.LegacyDec, error) {
	callOpts := &bind.CallOpts{BlockNumber: blockNumber}

	roundData, err := priceFeed.LatestRoundData(callOpts)
	if err != nil {
		return math.LegacyNewDec(0), fmt.Errorf("error fetching latest round data: %v", err)
	}

	decimals, err := priceFeed.Decimals(callOpts)
	if err != nil {
		return math.LegacyNewDec(0), fmt.Errorf("error fetching decimals: %v", err)
	}

	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	price := math.LegacyNewDecFromBigInt(roundData.Answer)
	divisorDec := math.LegacyNewDecFromBigInt(divisor)
	result := price.Quo(divisorDec)

	return result, nil
}

// CoinGeckoResponse represents the response from CoinGecko API
type CoinGeckoResponse struct {
	Zcash struct {
		USD float64 `json:"usd"`
	} `json:"zcash"`
}

// fetchZECPrice fetches the ZEC/USD price from CoinGecko API
func (o *Oracle) fetchZECPrice(ctx context.Context) (math.LegacyDec, error) {
	httpClient := &http.Client{
		Timeout: sidecartypes.DefaultHTTPTimeout,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", sidecartypes.ZECUSDPriceURL, nil)
	if err != nil {
		return math.LegacyZeroDec(), fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return math.LegacyZeroDec(), fmt.Errorf("failed to fetch ZEC price from CoinGecko: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return math.LegacyZeroDec(), fmt.Errorf("CoinGecko API returned status %d", resp.StatusCode)
	}

	var result CoinGeckoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return math.LegacyZeroDec(), fmt.Errorf("failed to decode CoinGecko response: %w", err)
	}

	if result.Zcash.USD == 0 {
		return math.LegacyZeroDec(), fmt.Errorf("ZEC price is zero or missing from CoinGecko response")
	}

	// Convert float64 to LegacyDec
	priceStr := fmt.Sprintf("%.2f", result.Zcash.USD)
	price, err := math.LegacyNewDecFromStr(priceStr)
	if err != nil {
		return math.LegacyZeroDec(), fmt.Errorf("failed to convert ZEC price to decimal: %w", err)
	}

	return price, nil
}
