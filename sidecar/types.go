package main

import (
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	solana "github.com/gagliardetto/solana-go/rpc"

	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/neutrino"
	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v5/sidecar/shared"
)

// NB: these constants should not be changed as they are important for synchronicity.
// Modifying them will exponentially increase the risk of your validator being slashed
const (
	MainLoopTickerInterval = 30 * time.Second
	CacheSize              = 20
	ROCKUSDPriceURL        = "https://api.gateio.ws/api/v4/spot/tickers?currency_pair=ROCK_USDT"
)

var (
	EmptyOracleState = sidecartypes.OracleState{
		EigenDelegations:           make(map[string]map[string]*big.Int),
		EthBlockHeight:             0,
		EthGasLimit:                0,
		EthBaseFee:                 0,
		EthTipCap:                  0,
		SolanaLamportsPerSignature: 0,
		RedemptionsEthereum:        []api.Redemption{},
		RedemptionsSolana:          []api.Redemption{},
		ROCKUSDPrice:               0,
		BTCUSDPrice:                0,
		ETHUSDPrice:                0,
	}
	// EthBlocksBeforeFinality   = big.NewInt(72)
	EthBlocksBeforeFinality = big.NewInt(0) // TODO: uncomment above and remove this line before mainnet
)

type Oracle struct {
	currentState   atomic.Value // *types.OracleState
	stateCache     []sidecartypes.OracleState
	Config         Config
	EthClient      *ethclient.Client
	neutrinoServer *neutrino.NeutrinoServer
	solanaClient   *solana.Client
	updateChan     chan sidecartypes.OracleState
	mainLoopTicker *time.Ticker
}

type Config struct {
	Enabled        bool              `yaml:"enabled"`
	GRPCPort       int               `yaml:"grpc_port"`
	StateFile      string            `yaml:"state_file"`
	OperatorConfig string            `yaml:"operator_config"`
	Network        string            `yaml:"network"`
	EthOracle      EthOracleConfig   `yaml:"eth_oracle"`
	SolanaRPC      map[string]string `yaml:"solana_rpc"`
	ProxyRPC       ProxyRPCConfig    `yaml:"proxy_rpc"`
	Neutrino       NeutrinoConfig    `yaml:"neutrino"`
}

type ProxyRPCConfig struct {
	URL      string `yaml:"url"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type NeutrinoConfig struct {
	Path string `yaml:"path"`
}

type EthOracleConfig struct {
	RPC           map[string]string `yaml:"rpc"`
	ContractAddrs ContractAddrs     `yaml:"contract_addrs"`
	NetworkName   map[string]string `yaml:"network_name"`
}

type ContractAddrs struct {
	ServiceManager string     `yaml:"service_manager"`
	PriceFeeds     PriceFeeds `yaml:"price_feeds"`
	ZenBTC         ZenBTC     `yaml:"zenbtc"`
}

type ZenBTC struct {
	Controller map[string]string `yaml:"controller"`
	Token      Networks          `yaml:"token"`
}

type Networks struct {
	Ethereum map[string]string `yaml:"ethereum"`
}

type PriceFeeds struct {
	BTC string `yaml:"btc"`
	ETH string `yaml:"eth"`
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
