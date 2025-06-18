package shared

import (
	"math/big"

	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	solrpc "github.com/gagliardetto/solana-go/rpc"
)

// Network constants
const (
	NetworkDevnet  = "devnet"
	NetworkTestnet = "testnet"
	NetworkMainnet = "mainnet"
)

// PriceFeed struct with fields for different price feeds
type PriceFeed struct {
	BTC string
	ETH string
}

// ZenBTCToken struct to hold token addresses for different blockchains
type ZenBTCToken struct {
	Ethereum map[string]string
}

// Contract address constants and other network-specific configuration values
// NB: these constants should not be changed as they are important for synchronicity.
// Modifying them will exponentially increase the risk of your validator being slashed
var (
	// ServiceManagerAddresses maps network names to service manager contract addresses
	ServiceManagerAddresses = map[string]string{
		NetworkDevnet:  "0xe2Aaf5A9a04cac7f3D43b4Afb7463850E1caEfB3",
		NetworkTestnet: "0xa559CDb9e029fc4078170122eBf7A3e622a764E4",
		NetworkMainnet: "0x4ca852BD78D9B7295874A7D223023Bff011b7EB3",
	}

	// PriceFeedAddresses contains addresses for different price feed contracts
	PriceFeedAddresses = PriceFeed{
		BTC: "0xF4030086522a5bEEa4988F8cA5B36dbC97BeE88c",
		ETH: "0x5f4eC3Df9cbd43714FE2740f5E3616155c5b8419",
	}

	// ZenBTCControllerAddresses maps network names to ZenBTC controller contract addresses
	ZenBTCControllerAddresses = map[string]string{
		NetworkDevnet:  "0x2844bd31B68AE5a0335c672e6251e99324441B73",
		NetworkTestnet: "0xaCE3634AAd9bCC48ef6A194f360F7ACe51F7d9f1",
		NetworkMainnet: "0xa87bE298115bE701A12F34F9B4585586dF052008",
	}

	// ZenBTCTokenAddresses holds token addresses for different blockchains
	ZenBTCTokenAddresses = ZenBTCToken{
		Ethereum: map[string]string{
			NetworkDevnet:  "0x7692E9a796001FeE9023853f490A692bAB2E4834",
			NetworkTestnet: "0xfA32a2D7546f8C7c229F94E693422A786DaE5E18",
			NetworkMainnet: "0x2fE9754d5D28bac0ea8971C0Ca59428b8644C776",
		},
	}

	// WhitelistedRoleAddresses maps network names to whitelisted role addresses
	WhitelistedRoleAddresses = map[string]string{
		NetworkDevnet:  "0x697bc4CAC913792f3D5BFdfE7655881A3b73e7Fe",
		NetworkTestnet: "0x75F1068e904815398045878A41e4324317c93aE4",
		NetworkMainnet: "0xBc17325952D043cCe5Bf1e4F42E26aE531962ED0",
	}

	// NetworkNames maps network identifiers to their human-readable names
	NetworkNames = map[string]string{
		NetworkDevnet:  "Holešky Ethereum Testnet",
		NetworkTestnet: "Holešky Ethereum Testnet",
		NetworkMainnet: "Ethereum Mainnet",
	}

	ZenBTCSolanaProgramID = map[string]string{
		NetworkDevnet:  "2pbhSDGggjXdRxp6qYjyeWLhvv4Ptf2r7QG8tbiBAZHq",
		NetworkTestnet: "9Gfr1YrMca5hyYRDP2nGxYkBWCSZtBm1oXBZyBdtYgNL",
		NetworkMainnet: "9t9RfpterTs95eXbKQWeAriZqET13TbjwDa6VW6LJHFb",
	}
	SolRockProgramID = map[string]string{
		NetworkDevnet:  "DXREJumiQhNejXa1b5EFPUxtSYdyJXBdiHeu6uX1ribA",
		NetworkTestnet: "4qXvX1jzVH2deMQGLZ8DXyQNkPdnMNQxHudyZEZAEa4f",
		NetworkMainnet: "3WyacwnCNiz4Q1PedWyuwodYpLFu75jrhgRTZp69UcA9",
	}

	// Solana RPC endpoints
	SolanaRPCEndpoints = map[string]string{
		NetworkDevnet:  solrpc.DevNet_RPC,
		NetworkTestnet: solrpc.DevNet_RPC,
		NetworkMainnet: solrpc.MainNetBeta_RPC,
	}

	// Solana CAIP-2 Identifiers (Map network name to CAIP-2 string)
	// Solana devnet is used for both devnet and testnet environments
	SolanaCAIP2 = map[string]string{
		NetworkDevnet:  "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		NetworkTestnet: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		NetworkMainnet: "solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp",
	}

	// ROCK Price feed URL
	ROCKUSDPriceURL = "https://api.gateio.ws/api/v4/spot/tickers?currency_pair=ROCK_USDT"

	// Oracle tuning parameters
	MainLoopTickerIntervalSeconds = 60 // Seconds
	OracleCacheSize               = 20
	EthBurnEventsBlockRange       = 1000
	EthBlocksBeforeFinality       = int64(5) // TODO: should this be increased?
	SolanaEventScanTxLimit        = 1000
	SolanaEventFetchBatchSize     = 50
)

// SolanaEventType defines a type for Solana event keys for type safety.
type SolanaEventType string

// Constants for Solana event types to avoid magic strings
const (
	SolRockMint   SolanaEventType = "solRockMint"
	SolZenBTCMint SolanaEventType = "solZenBTCMint"
	SolZenBTCBurn SolanaEventType = "solZenBTCBurn"
	SolRockBurn   SolanaEventType = "solRockBurn"
)

type OracleState struct {
	EigenDelegations           map[string]map[string]*big.Int `json:"eigenDelegations"`
	EthBlockHeight             uint64                         `json:"ethBlockHeight"`
	EthGasLimit                uint64                         `json:"ethGasLimit"`
	EthBaseFee                 uint64                         `json:"ethBaseFee"`
	EthTipCap                  uint64                         `json:"ethTipCap"`
	SolanaLamportsPerSignature uint64                         `json:"solanaLamportsPerSignature"`
	EthBurnEvents              []api.BurnEvent                `json:"ethBurnEvents"`
	CleanedEthBurnEvents       map[string]bool                `json:"cleanedEthBurnEvents"`
	SolanaBurnEvents           []api.BurnEvent                `json:"solanaBurnEvents"`
	CleanedSolanaBurnEvents    map[string]bool                `json:"cleanedSolanaBurnEvents"`
	Redemptions                []api.Redemption               `json:"redemptions"`
	ROCKUSDPrice               math.LegacyDec                 `json:"rockUSDPrice"`
	BTCUSDPrice                math.LegacyDec                 `json:"btcUSDPrice"`
	ETHUSDPrice                math.LegacyDec                 `json:"ethUSDPrice"`
	SolanaMintEvents           []api.SolanaMintEvent          `json:"solanaMintEvents"`
	CleanedSolanaMintEvents    map[string]bool                `json:"cleanedSolanaMintEvents"`
	// Fields for watermarking Solana events
	LastSolRockMintSig   string `json:"lastSolRockMintSig,omitempty"`
	LastSolZenBTCMintSig string `json:"lastSolZenBTCMintSig,omitempty"`
	LastSolZenBTCBurnSig string `json:"lastSolZenBTCBurnSig,omitempty"`
	LastSolRockBurnSig   string `json:"lastSolRockBurnSig,omitempty"`
}

type Config struct {
	Enabled                  bool              `yaml:"enabled"`
	GRPCPort                 int               `yaml:"grpc_port"`
	StateFile                string            `yaml:"state_file"`
	ZRChainRPC               string            `yaml:"zrchain_rpc"`
	OperatorConfig           string            `yaml:"operator_config"`
	Network                  string            `yaml:"network"`
	EthRPC                   map[string]string `yaml:"eth_rpc"`
	SolanaRPC                map[string]string `yaml:"solana_rpc"`
	ProxyRPC                 ProxyRPCConfig    `yaml:"proxy_rpc"`
	Neutrino                 NeutrinoConfig    `yaml:"neutrino"`
	FilterEventsAfterStartup bool              `yaml:"filter_events_after_startup"`
}

type ProxyRPCConfig struct {
	URL      string `yaml:"url"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type NeutrinoConfig struct {
	Path string `yaml:"path"`
	Port int    `yaml:"port"`
}
