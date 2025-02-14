package shared

import (
	"math/big"

	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
)

type OracleState struct {
	EigenDelegations           map[string]map[string]*big.Int `json:"eigenDelegations"`
	EthBlockHeight             uint64                         `json:"ethBlockHeight"`
	EthGasLimit                uint64                         `json:"ethGasLimit"`
	EthBaseFee                 uint64                         `json:"ethBaseFee"`
	EthTipCap                  uint64                         `json:"ethTipCap"`
	SolanaLamportsPerSignature uint64                         `json:"solanaLamportsPerSignature"`
	EthBurnEvents              []api.BurnEvent                `json:"ethBurnEvents"`
	Redemptions                []api.Redemption               `json:"redemptions"`
	ROCKUSDPrice               float64                        `json:"rockUSDPrice"`
	BTCUSDPrice                float64                        `json:"btcUSDPrice"`
	ETHUSDPrice                float64                        `json:"ethUSDPrice"` // TODO: remove field if we won't use ETH stake?
}

type Config struct {
	Enabled        bool              `yaml:"enabled"`
	GRPCPort       int               `yaml:"grpc_port"`
	StateFile      string            `yaml:"state_file"`
	ZRChainRPC     string            `yaml:"zrchain_rpc"`
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
