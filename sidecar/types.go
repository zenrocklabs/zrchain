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

// / These constants should not be changed as they are important for synchronicity
const (
	MainLoopTickerInterval = 15 * time.Second
	CacheSize              = 20
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

type CoinMarketCapResponse struct {
	Status struct {
		Timestamp    string `json:"timestamp"`
		ErrorCode    int    `json:"error_code"`
		ErrorMessage string `json:"error_message"`
		Elapsed      int    `json:"elapsed"`
		CreditCount  int    `json:"credit_count"`
		Notice       string `json:"notice"`
	} `json:"status"`
	Data struct {
		ETH struct {
			Quote struct {
				USD struct {
					Price float64 `json:"price"`
				} `json:"USD"`
			} `json:"quote"`
		} `json:"ETH"`
	} `json:"data"`
}

type Config struct {
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
	ServiceManager string `yaml:"service_manager"`
	PriceFeed      string `yaml:"price_feed"`
	ZenBTC         ZenBTC `yaml:"zen_btc"`
}

type ZenBTC struct {
	EthMainnet string `yaml:"eth_mainnet"`
	EthHolesky string `yaml:"eth_holesky"`
}
