package shared

import (
	"time"

	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
)

// Network constants
const (
	NetworkLocalnet = "localnet"
	NetworkRegnet   = "regnet"
	NetworkDevnet   = "devnet"
	NetworkTestnet  = "testnet"
	NetworkMainnet  = "mainnet"
)

// Contract address constants and other network-specific configuration values
// NB: these constants should not be changed as they are important for synchronicity.
// Modifying them will exponentially increase the risk of your validator being slashed
var (
	// PriceFeedAddresses contains addresses for different price feed contracts on Ethereum mainnet
	PriceFeedAddresses = PriceFeed{
		BTC: "0xF4030086522a5bEEa4988F8cA5B36dbC97BeE88c", // BTC/USD Chainlink feed
		ETH: "0x5f4eC3Df9cbd43714FE2740f5E3616155c5b8419", // ETH/USD Chainlink feed
	}

	// ZEC Price URL from Binance - RISK OF SLASHING IF CHANGED
	ZECUSDPriceURL = "https://api.binance.com/api/v3/ticker/price?symbol=ZECUSDT"

	// ROCK Token ID on Solana - RISK OF SLASHING IF CHANGED
	ROCKTokenID = "5VsPJ2EG7jjo3k2LPzQVriENKKQkNUTzujEzuaj4Aisf"

	// ROCK Price feed URL from Jupiter API - RISK OF SLASHING IF CHANGED
	ROCKUSDPriceURL = "https://lite-api.jup.ag/price/v3?ids=" + ROCKTokenID

	// ZenBTCControllerAddresses maps network names to ZenBTC controller contract addresses
	ZenBTCControllerAddresses = map[string]string{
		NetworkLocalnet: "0x2844bd31B68AE5a0335c672e6251e99324441B73",
		NetworkRegnet:   "0x2419A36682f329d0B5e1834068C8C63046865504",
		NetworkDevnet:   "0x2844bd31B68AE5a0335c672e6251e99324441B73",
		NetworkTestnet:  "0xaCE3634AAd9bCC48ef6A194f360F7ACe51F7d9f1",
		NetworkMainnet:  "0xa87bE298115bE701A12F34F9B4585586dF052008",
	}

	// ZenBTCTokenAddresses holds token addresses for different blockchains
	ZenBTCTokenAddresses = ZenBTCToken{
		Ethereum: map[string]string{
			NetworkLocalnet: "0x7692E9a796001FeE9023853f490A692bAB2E4834",
			NetworkRegnet:   "0x745Aa06072bf149117C457C68b0531cF7567C4e1",
			NetworkDevnet:   "0x7692E9a796001FeE9023853f490A692bAB2E4834",
			NetworkTestnet:  "0xfA32a2D7546f8C7c229F94E693422A786DaE5E18",
			NetworkMainnet:  "0x2fE9754d5D28bac0ea8971C0Ca59428b8644C776",
		},
	}

	// WhitelistedRoleAddresses maps network names to whitelisted role addresses
	WhitelistedRoleAddresses = map[string]string{
		NetworkLocalnet: "0x697bc4CAC913792f3D5BFdfE7655881A3b73e7Fe",
		NetworkRegnet:   "0x75F1068e904815398045878A41e4324317c93aE4",
		NetworkDevnet:   "0x697bc4CAC913792f3D5BFdfE7655881A3b73e7Fe",
		NetworkTestnet:  "0x75F1068e904815398045878A41e4324317c93aE4",
		NetworkMainnet:  "0xBc17325952D043cCe5Bf1e4F42E26aE531962ED0",
	}

	// NetworkNames maps network identifiers to their human-readable names
	NetworkNames = map[string]string{
		NetworkLocalnet: "Hoodi Ethereum Testnet",
		NetworkRegnet:   "Hoodi Ethereum Testnet",
		NetworkDevnet:   "Hoodi Ethereum Testnet",
		NetworkTestnet:  "Hoodi Ethereum Testnet",
		NetworkMainnet:  "Ethereum Mainnet",
	}

	ZenBTCSolanaProgramID = map[string]string{
		NetworkLocalnet: "886BBKmJ71jrqWhZBAzqhdNRKk761Mx9WVCPzFkY4uBb",
		NetworkRegnet:   "9Gfr1YrMca5hyYRDP2nGxYkBWCSZtBm1oXBZyBdtYgNL",
		NetworkDevnet:   "886BBKmJ71jrqWhZBAzqhdNRKk761Mx9WVCPzFkY4uBb",
		NetworkTestnet:  "9Gfr1YrMca5hyYRDP2nGxYkBWCSZtBm1oXBZyBdtYgNL",
		NetworkMainnet:  "9t9RfpterTs95eXbKQWeAriZqET13TbjwDa6VW6LJHFb",
	}

	ZenZECSolanaProgramID = map[string]string{
		NetworkLocalnet: "CfA179tLmrg9HeADSmGdQxZkRe85gxdNxjiARx5PXjiD",
		NetworkRegnet:   "",
		NetworkDevnet:   "CfA179tLmrg9HeADSmGdQxZkRe85gxdNxjiARx5PXjiD",
		NetworkTestnet:  "7y6tokt7ua3t8BaPVbpzS6wKS9PzKupQUBmXFrpAh83T",
		NetworkMainnet:  "BTzxmuLgNUfBeNCFxsoSVpEKjRTma75eWcaPjyDCUF88",
	}

	ZenZECEventStoreProgramID = map[string]string{
		NetworkLocalnet: "8jroAuLzFMzb1an5Rptf9Lx7JyTWFHMF1tA2ZW9Pymys",
		NetworkRegnet:   "",
		NetworkDevnet:   "8jroAuLzFMzb1an5Rptf9Lx7JyTWFHMF1tA2ZW9Pymys",
		NetworkTestnet:  "9gzutLMGsSPboaaXuyK87Uu3VsM1mRgFTddsCEFWwYVV",
		NetworkMainnet:  "7nT8EmELLWf4gSiWxT2xH5oR6GsWRn3X6dv2AY2bpa53",
	}

	ZenZECMintAddress = map[string]string{
		NetworkLocalnet: "3YMe7Bbus2rZiDR7ijRBhT6hNFvNwFnBgGEwxkw3L71g",
		NetworkRegnet:   "",
		NetworkDevnet:   "3YMe7Bbus2rZiDR7ijRBhT6hNFvNwFnBgGEwxkw3L71g",
		NetworkTestnet:  "H8X7ogzmLEuU36tHPFGFrMtFhHPNwbjuZW6FSGDgS9Jt",
		NetworkMainnet:  "JDt9rRGaieF6aN1cJkXFeUmsy7ZE4yY3CZb8tVMXVroS",
	}

	SolRockProgramID = map[string]string{
		NetworkLocalnet: "AgoRvPWg2R7nkKhxvipvms79FmxQr75r2GwNSpPtxcLg",
		NetworkRegnet:   "9CNTbJY29vHPThkMXCVNozdhXtWrWHyxVy39EhpRtiXe",
		NetworkDevnet:   "AgoRvPWg2R7nkKhxvipvms79FmxQr75r2GwNSpPtxcLg",
		NetworkTestnet:  "4qXvX1jzVH2deMQGLZ8DXyQNkPdnMNQxHudyZEZAEa4f",
		NetworkMainnet:  "3WyacwnCNiz4Q1PedWyuwodYpLFu75jrhgRTZp69UcA9",
	}

	// Solana CAIP-2 Identifiers (Map network name to CAIP-2 string)
	// Solana devnet is used for both devnet and testnet environments
	SolanaCAIP2 = map[string]string{
		NetworkLocalnet: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		NetworkRegnet:   "solana:HK8b7Skns2TX3FvXQxm2mPQbY2nVY8GD",
		NetworkDevnet:   "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		NetworkTestnet:  "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		NetworkMainnet:  "solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp",
	}

	// VersionsRequiringCacheReset lists sidecar versions that need a one-time cache wipe.
	// This protects against subtle state incompatibilities after major upgrades.
	VersionsRequiringCacheReset = []string{"sturgeon_moon_r4"}
	MaxSupportedSolanaTxVersion = func() *uint64 { u := uint64(0); return &u }() // Solana transaction version 0
)

const (

	// Oracle tuning parameters - RISK OF SLASHING IF CHANGED

	MainLoopTickerInterval         = 60 * time.Second
	OracleCacheSize                = 5
	EthBurnEventsBlockRange        = 1000
	EthBlocksBeforeFinality        = int64(8) // TODO: should this be increased?
	SolanaEventScanTxLimit         = 50
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
	GasEstimationBuffer         = uint64(110) // 110% buffer for gas estimation (10% extra)

	// Oracle processing constants
	ErrorChannelBufferSize              = 16                // Buffer size for error channels in goroutines
	InitialEventsSliceCapacity          = 100               // Initial capacity for events slice to reduce allocations
	WrapCallGasLimitFallback            = uint64(100_000)   // Fallback gas limit for zenBTC wrap call
	PendingTransactionCheckInterval     = 5 * time.Second   // How often to check for pending transactions when queue is empty
	PendingTransactionStatusLogInterval = 15 * time.Second  // How often to log pending transaction processing status
	PendingTransactionMaxRetries        = 100               // Maximum retry attempts before removing from pending queue
	TransactionCacheTTL                 = 5 * time.Minute   // Time-to-live for cached transaction results
	NTPServer                           = "time.google.com" // NTP server for time synchronization
	TimeFormatPrecise                   = "15:04:05.00"     // Time format for precise logging (HH:MM:SS.ms)

	SidecarVersionName = "sturgeon_moon_r4"

	// OracleStateResetInterval controls how often (in UTC hours) the oracle
	// should perform a full in-memory + cache reset cycle.
	// Default: 24 (reset at the first tick after each UTC midnight).
	// Example values:
	//   24 => once per day (midnight UTC)
	//   12 => twice per day (00:00 UTC, 12:00 UTC)
	//    6 => every 6 hours
	//    3 => every 3 hours, etc.
	OracleStateResetInterval = 3 * time.Hour
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
	SolZenZECMint SolanaEventType = "solZenZECMint"
	SolZenZECBurn SolanaEventType = "solZenZECBurn"
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
	EthBlockHeight          uint64                `json:"ethBlockHeight"`
	EthGasLimit             uint64                `json:"ethGasLimit"`
	EthBaseFee              uint64                `json:"ethBaseFee"`
	EthTipCap               uint64                `json:"ethTipCap"`
	EthBurnEvents           []api.BurnEvent       `json:"ethBurnEvents"`
	CleanedEthBurnEvents    map[string]bool       `json:"cleanedEthBurnEvents"`
	SolanaBurnEvents        []api.BurnEvent       `json:"solanaBurnEvents"`
	CleanedSolanaBurnEvents map[string]bool       `json:"cleanedSolanaBurnEvents"`
	Redemptions             []api.Redemption      `json:"redemptions"`
	ROCKUSDPrice            math.LegacyDec        `json:"rockUSDPrice"`
	BTCUSDPrice             math.LegacyDec        `json:"btcUSDPrice"`
	ETHUSDPrice             math.LegacyDec        `json:"ethUSDPrice"`
	ZECUSDPrice             math.LegacyDec        `json:"zecUSDPrice"`
	SolanaMintEvents        []api.SolanaMintEvent `json:"solanaMintEvents"`
	CleanedSolanaMintEvents map[string]bool       `json:"cleanedSolanaMintEvents"`
	// ZCash block headers
	LatestZcashBlockHeight    int64               `json:"latestZcashBlockHeight,omitempty"`
	LatestZcashBlockHeader    *api.BTCBlockHeader `json:"latestZcashBlockHeader,omitempty"`
	RequestedZcashBlockHeight int64               `json:"requestedZcashBlockHeight,omitempty"`
	RequestedZcashBlockHeader *api.BTCBlockHeader `json:"requestedZcashBlockHeader,omitempty"`
	// Fields for watermarking Solana events
	LastSolRockMintSig   string `json:"lastSolRockMintSig,omitempty"`
	LastSolZenBTCMintSig string `json:"lastSolZenBTCMintSig,omitempty"`
	LastSolZenBTCBurnSig string `json:"lastSolZenBTCBurnSig,omitempty"`
	LastSolZenZECMintSig string `json:"lastSolZenZECMintSig,omitempty"`
	LastSolZenZECBurnSig string `json:"lastSolZenZECBurnSig,omitempty"`
	LastSolRockBurnSig   string `json:"lastSolRockBurnSig,omitempty"`
	// Pending transactions that failed processing and need to be retried
	PendingSolanaTxs map[string]PendingTxInfo `json:"pendingSolanaTxs,omitempty"`
}

type Config struct {
	Enabled                bool              `yaml:"enabled"`
	GRPCPort               int               `yaml:"grpc_port"`
	StateFile              string            `yaml:"state_file"`
	ZRChainRPC             string            `yaml:"zrchain_rpc"`
	Network                string            `yaml:"network"`
	EthRPC                 map[string]string `yaml:"eth_rpc"`
	SolanaRPC              map[string]string `yaml:"solana_rpc"`
	ZcashRPC               map[string]string `yaml:"zcash_rpc"`
	ProxyRPC               ProxyRPCConfig    `yaml:"proxy_rpc"`
	Neutrino               NeutrinoConfig    `yaml:"neutrino"`
	E2ETestsTickerInterval int               `yaml:"e2e_tests_ticker_interval"`
	DebugMode              bool              `yaml:"debug_mode"`
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
