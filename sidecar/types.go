package main

import (
	"math/big"
	"sync/atomic"
	"time"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/ethclient"
	solana "github.com/gagliardetto/solana-go/rpc"

	"github.com/Zenrock-Foundation/zrchain/v5/go-client"
	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/neutrino"
	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v5/sidecar/shared"
)

// NB: these constants should not be changed as they are important for synchronicity.
// Modifying them will exponentially increase the risk of your validator being slashed
const (
	MainLoopTickerInterval  = 10 * time.Second
	CacheSize               = 20
	EthBurnEventsBlockRange = 1000
	ROCKUSDPriceURL         = "https://api.gateio.ws/api/v4/spot/tickers?currency_pair=ROCK_USDT"
)

var (
	EmptyOracleState = sidecartypes.OracleState{
		EigenDelegations:           make(map[string]map[string]*big.Int),
		EthBlockHeight:             0,
		EthGasLimit:                0,
		EthBaseFee:                 0,
		EthTipCap:                  0,
		SolanaLamportsPerSignature: 0,
		EthBurnEvents:              []api.BurnEvent{},
		CleanedEthBurnEvents:       make(map[string]bool),
		Redemptions:                []api.Redemption{},
		ROCKUSDPrice:               math.LegacyNewDec(0),
		BTCUSDPrice:                math.LegacyNewDec(0),
		ETHUSDPrice:                math.LegacyNewDec(0),
	}
	EthBlocksBeforeFinality = big.NewInt(5) // TODO: should this be increased?
)

type Oracle struct {
	currentState       atomic.Value // *types.OracleState
	stateCache         []sidecartypes.OracleState
	Config             sidecartypes.Config
	EthClient          *ethclient.Client
	neutrinoServer     *neutrino.NeutrinoServer
	solanaClient       *solana.Client
	zrChainQueryClient *client.QueryClient
	updateChan         chan sidecartypes.OracleState
	mainLoopTicker     *time.Ticker
}

type PriceData struct {
	CurrencyPair     string `json:"currency_pair"`
	Last             string `json:"last"`
	LowestAsk        string `json:"lowest_ask"`
	LowestSize       string `json:"lowest_size"`
	HighestBid       string `json:"highest_bid"`
	HighestSize      string `json:"highest_size"`
	ChangePercentage string `json:"change_percentage"`
	BaseVolume       string `json:"base_volume"`
	QuoteVolume      string `json:"quote_volume"`
	High24h          string `json:"high_24h"`
	Low24h           string `json:"low_24h"`
}
