package shared

import (
	"math/big"
	"time"

	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	solrpc "github.com/gagliardetto/solana-go/rpc"
)

// Network constants
const (
	NetworkDevnet  = "devnet"
	NetworkRegnet  = "regnet"
	NetworkTestnet = "testnet"
	NetworkMainnet = "mainnet"
)

// Contract address constants and other network-specific configuration values
// NB: these constants should not be changed as they are important for synchronicity.
// Modifying them will exponentially increase the risk of your validator being slashed
var (
	// ServiceManagerAddresses maps network names to service manager contract addresses
	ServiceManagerAddresses = map[string]string{
		NetworkDevnet:  "0xe2Aaf5A9a04cac7f3D43b4Afb7463850E1caEfB3",
		NetworkRegnet:  "0xa6c639cC8506B13d7cb37bFa143318908050Fb70",
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
		NetworkRegnet:  "0x2419A36682f329d0B5e1834068C8C63046865504",
		NetworkTestnet: "0xaCE3634AAd9bCC48ef6A194f360F7ACe51F7d9f1",
		NetworkMainnet: "0xa87bE298115bE701A12F34F9B4585586dF052008",
	}

	// ZenBTCTokenAddresses holds token addresses for different blockchains
	ZenBTCTokenAddresses = ZenBTCToken{
		Ethereum: map[string]string{
			NetworkDevnet:  "0x7692E9a796001FeE9023853f490A692bAB2E4834",
			NetworkRegnet:  "0x745Aa06072bf149117C457C68b0531cF7567C4e1",
			NetworkTestnet: "0xfA32a2D7546f8C7c229F94E693422A786DaE5E18",
			NetworkMainnet: "0x2fE9754d5D28bac0ea8971C0Ca59428b8644C776",
		},
	}

	// WhitelistedRoleAddresses maps network names to whitelisted role addresses
	WhitelistedRoleAddresses = map[string]string{
		NetworkDevnet:  "0x697bc4CAC913792f3D5BFdfE7655881A3b73e7Fe",
		NetworkRegnet:  "0x75F1068e904815398045878A41e4324317c93aE4",
		NetworkTestnet: "0x75F1068e904815398045878A41e4324317c93aE4",
		NetworkMainnet: "0xBc17325952D043cCe5Bf1e4F42E26aE531962ED0",
	}

	// NetworkNames maps network identifiers to their human-readable names
	NetworkNames = map[string]string{
		NetworkDevnet:  "Hoodi Ethereum Testnet",
		NetworkRegnet:  "Hoodi Ethereum Testnet",
		NetworkTestnet: "Hoodi Ethereum Testnet",
		NetworkMainnet: "Ethereum Mainnet",
	}

	ZenBTCSolanaProgramID = map[string]string{
		NetworkDevnet:  "2pbhSDGggjXdRxp6qYjyeWLhvv4Ptf2r7QG8tbiBAZHq",
		NetworkTestnet: "9Gfr1YrMca5hyYRDP2nGxYkBWCSZtBm1oXBZyBdtYgNL",
		NetworkRegnet:  "9Gfr1YrMca5hyYRDP2nGxYkBWCSZtBm1oXBZyBdtYgNL",
		NetworkMainnet: "9t9RfpterTs95eXbKQWeAriZqET13TbjwDa6VW6LJHFb",
	}
	SolRockProgramID = map[string]string{
		NetworkDevnet:  "AgoRvPWg2R7nkKhxvipvms79FmxQr75r2GwNSpPtxcLg",
		NetworkRegnet:  "9CNTbJY29vHPThkMXCVNozdhXtWrWHyxVy39EhpRtiXe",
		NetworkTestnet: "4qXvX1jzVH2deMQGLZ8DXyQNkPdnMNQxHudyZEZAEa4f",
		NetworkMainnet: "3WyacwnCNiz4Q1PedWyuwodYpLFu75jrhgRTZp69UcA9",
	}

	// Solana RPC endpoints
	SolanaRPCEndpoints = map[string]string{
		NetworkDevnet:  solrpc.DevNet_RPC,
		NetworkRegnet:  solrpc.DevNet_RPC,
		NetworkTestnet: solrpc.DevNet_RPC,
		NetworkMainnet: solrpc.MainNetBeta_RPC,
	}

	// Solana CAIP-2 Identifiers (Map network name to CAIP-2 string)
	// Solana devnet is used for both devnet and testnet environments
	SolanaCAIP2 = map[string]string{
		NetworkDevnet:  "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		NetworkRegnet:  "solana:HK8b7Skns2TX3FvXQxm2mPQbY2nVY8GD",
		NetworkTestnet: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		NetworkMainnet: "solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp",
	}

	// ROCK Price feed URL - RISK OF SLASHING IF CHANGED
	ROCKUSDPriceURL = "https://api.gateio.ws/api/v4/spot/tickers?currency_pair=ROCK_USDT"

	// Oracle tuning parameters - RISK OF SLASHING IF CHANGED

	MainLoopTickerInterval         = 60 * time.Second
	OracleCacheSize                = 10
	EthBurnEventsBlockRange        = 1000
	EthBlocksBeforeFinality        = int64(8) // TODO: should this be increased?
	SolanaEventScanTxLimit         = 200
	SolanaMaxBackfillPages         = 10 // Max pages to fetch when filling a signature gap.
	SolanaEventFetchBatchSize      = 10
	SolanaEventFetchMinBatchSize   = 2
	SolanaSleepInterval            = 50 * time.Millisecond
	SolanaFallbackSleepInterval    = 10 * time.Millisecond // Sleep between individual fallback requests
	SolanaEventFetchMaxRetries     = 10
	SolanaFallbackMaxRetries       = 3 // Retries for individual fallback requests
	SolanaEventFetchRetrySleep     = 100 * time.Millisecond
	SolanaPendingTxMaxRetries      = 100
	SolanaPendingTxAllowRetryAfter = 5 * time.Second

	// RPC connection retry parameters
	RPCConnectionMaxRetries = 20
	RPCConnectionRetryDelay = 3 * time.Second

	// HTTP and RPC constants
	DefaultHTTPTimeout          = 10 * time.Second
	SolanaRPCTimeout            = 10 * time.Second // Longer timeout for Solana RPC operations
	SolanaBatchTimeout          = 20 * time.Second // Even longer for batch operations
	SolanaRateLimiterTimeout    = 10 * time.Second
	SolanaMaxConcurrentRPCCalls = 20          // Maximum concurrent Solana RPC calls (semaphore size)
	MaxSupportedSolanaTxVersion = uint64(0)   // Solana transaction version 0
	EigenLayerQuorumNumber      = uint8(0)    // EigenLayer quorum number for service manager
	GasEstimationBuffer         = uint64(110) // 110% buffer for gas estimation (10% extra)

	SidecarVersionName = "salmon_moon_r6"

	// VersionsRequiringCacheReset lists sidecar versions that need a one-time cache wipe.
	// This protects against subtle state incompatibilities after major upgrades.
	VersionsRequiringCacheReset = []string{"salmon_moon_r3", "salmon_moon_r4", "salmon_moon_r5", "salmon_moon_r6"}

	// Oracle processing constants
	ErrorChannelBufferSize              = 16                // Buffer size for error channels in goroutines
	InitialEventsSliceCapacity          = 100               // Initial capacity for events slice to reduce allocations
	StakeCallDataAmount                 = int64(1000000000) // Amount used for zenBTC stake call gas estimation
	PendingTransactionCheckInterval     = 5 * time.Second   // How often to check for pending transactions when queue is empty
	PendingTransactionStatusLogInterval = 15 * time.Second  // How often to log pending transaction processing status
	PendingTransactionMaxRetries        = 100               // Maximum retry attempts before removing from pending queue
	TransactionCacheTTL                 = 5 * time.Minute   // Time-to-live for cached transaction results
	NTPServer                           = "time.google.com" // NTP server for time synchronization
	TimeFormatPrecise                   = "15:04:05.00"     // Time format for precise logging (HH:MM:SS.ms)
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

// SolanaEventType defines a type for Solana event keys for type safety.
type SolanaEventType string

// Constants for Solana event types to avoid magic strings
const (
	SolRockMint   SolanaEventType = "solRockMint"
	SolZenBTCMint SolanaEventType = "solZenBTCMint"
	SolZenBTCBurn SolanaEventType = "solZenBTCBurn"
	SolRockBurn   SolanaEventType = "solRockBurn"
)

// PendingTxInfo represents a failed transaction that needs to be retried
type PendingTxInfo struct {
	Signature    string    `json:"signature"`
	EventType    string    `json:"eventType"`
	RetryCount   int       `json:"retryCount"`
	LastAttempt  time.Time `json:"lastAttempt"`
	FirstAttempt time.Time `json:"firstAttempt"`
}

type OracleState struct {
	EigenDelegations        map[string]map[string]*big.Int `json:"eigenDelegations"`
	EthBlockHeight          uint64                         `json:"ethBlockHeight"`
	EthGasLimit             uint64                         `json:"ethGasLimit"`
	EthBaseFee              uint64                         `json:"ethBaseFee"`
	EthTipCap               uint64                         `json:"ethTipCap"`
	EthBurnEvents           []api.BurnEvent                `json:"ethBurnEvents"`
	CleanedEthBurnEvents    map[string]bool                `json:"cleanedEthBurnEvents"`
	SolanaBurnEvents        []api.BurnEvent                `json:"solanaBurnEvents"`
	CleanedSolanaBurnEvents map[string]bool                `json:"cleanedSolanaBurnEvents"`
	Redemptions             []api.Redemption               `json:"redemptions"`
	ROCKUSDPrice            math.LegacyDec                 `json:"rockUSDPrice"`
	BTCUSDPrice             math.LegacyDec                 `json:"btcUSDPrice"`
	ETHUSDPrice             math.LegacyDec                 `json:"ethUSDPrice"`
	SolanaMintEvents        []api.SolanaMintEvent          `json:"solanaMintEvents"`
	CleanedSolanaMintEvents map[string]bool                `json:"cleanedSolanaMintEvents"`
	// Fields for watermarking Solana events
	LastSolRockMintSig   string `json:"lastSolRockMintSig,omitempty"`
	LastSolZenBTCMintSig string `json:"lastSolZenBTCMintSig,omitempty"`
	LastSolZenBTCBurnSig string `json:"lastSolZenBTCBurnSig,omitempty"`
	LastSolRockBurnSig   string `json:"lastSolRockBurnSig,omitempty"`
	// Pending transactions that failed processing and need to be retried
	PendingSolanaTxs map[string]PendingTxInfo `json:"pendingSolanaTxs,omitempty"`
}

type Config struct {
	Enabled                bool              `yaml:"enabled"`
	GRPCPort               int               `yaml:"grpc_port"`
	StateFile              string            `yaml:"state_file"`
	ZRChainRPC             string            `yaml:"zrchain_rpc"`
	OperatorConfig         string            `yaml:"operator_config"`
	Network                string            `yaml:"network"`
	EthRPC                 map[string]string `yaml:"eth_rpc"`
	SolanaRPC              map[string]string `yaml:"solana_rpc"`
	ProxyRPC               ProxyRPCConfig    `yaml:"proxy_rpc"`
	Neutrino               NeutrinoConfig    `yaml:"neutrino"`
	E2ETestsTickerInterval int               `yaml:"e2e_tests_ticker_interval"`
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
