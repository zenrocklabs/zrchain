package main

import (
	"context"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/ethclient"
	sol "github.com/gagliardetto/solana-go"
	solana "github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"

	"github.com/Zenrock-Foundation/zrchain/v6/go-client"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/neutrino"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
)

var (
	EmptyOracleState = sidecartypes.OracleState{
		EthBlockHeight:          0,
		EthGasLimit:             0,
		EthBaseFee:              0,
		EthTipCap:               0,
		EthBurnEvents:           []api.BurnEvent{},
		CleanedEthBurnEvents:    make(map[string]bool),
		SolanaBurnEvents:        []api.BurnEvent{},
		CleanedSolanaBurnEvents: make(map[string]bool),
		Redemptions:             []api.Redemption{},
		SolanaMintEvents:        []api.SolanaMintEvent{},
		CleanedSolanaMintEvents: make(map[string]bool),
		ROCKUSDPrice:            math.LegacyNewDec(0),
		BTCUSDPrice:             math.LegacyNewDec(0),
		ETHUSDPrice:             math.LegacyNewDec(0),
		ZECUSDPrice:             math.LegacyNewDec(0),
		PendingSolanaTxs:        make(map[string]sidecartypes.PendingTxInfo),
	}
)

type Oracle struct {
	// currentState stores the oracle's current state atomically for lock-free reads.
	// Uses atomic.Pointer[T] for type-safe access without type assertions.
	//
	// IMMUTABILITY CONTRACT: Use copy-on-write pattern for all updates:
	//
	//   old := o.currentState.Load()  // Type-safe: returns *OracleState
	//   new := *old                   // Copy the struct
	//   new.Field = newValue          // Modify the copy
	//   o.currentState.Store(&new)    // Store new pointer
	//
	// NEVER mutate in-place (causes data race):
	//   state := o.currentState.Load()
	//   state.Field = newValue          // ❌ Mutates shared state!
	//   o.currentState.Store(state)     // ❌ Same pointer
	//
	// See docs/reports/2024-10-24-atomic-value-pointer-vs-value-storage-analysis.md
	// for detailed analysis of pointer vs value storage and concurrency patterns.
	currentState atomic.Pointer[sidecartypes.OracleState]

	stateCache         []sidecartypes.OracleState
	Config             sidecartypes.Config
	EthClient          *ethclient.Client
	neutrinoServer     *neutrino.NeutrinoServer
	zcashClient        *ZcashClient
	solanaClient       *solana.Client
	zrChainQueryClient *client.QueryClient
	mainLoopTicker     *time.Ticker
	DebugMode          bool
	SkipInitialWait    bool

	// Periodic reset control (scheduled UTC boundary resets). Interval derived on-the-fly from sidecartypes.OracleStateResetIntervalHours (or test flag).
	nextScheduledReset time.Time
	ForceTestReset     bool         // when true (set via test flag) use a 2-minute interval for rapid testing
	stateCacheMutex    sync.RWMutex // guards stateCache access and scheduled resets (RLock for reads, Lock for writes)

	// Last processed Solana signatures (managed as strings for persistence)
	lastSolRockMintSigStr   string
	lastSolZenBTCMintSigStr string
	lastSolZenBTCBurnSigStr string
	lastSolZenZECMintSigStr string
	lastSolZenZECBurnSigStr string
	lastSolRockBurnSigStr   string

	// ZCash header tracking
	lastZcashHeaderHeight int64

	// Performance optimization fields
	solanaRateLimiter     chan struct{}              // Semaphore for Solana RPC rate limiting
	transactionCache      map[string]*CachedTxResult // Cache for frequently accessed transactions
	transactionCacheMutex sync.RWMutex               // Protects transaction cache

	// Function fields for mocking
	getSolanaZenBTCBurnEventsFn func(ctx context.Context, programID string, lastKnownSig sol.Signature) ([]api.BurnEvent, sol.Signature, error)
	getSolanaZenZECBurnEventsFn func(ctx context.Context, programID string, lastKnownSig sol.Signature) ([]api.BurnEvent, sol.Signature, error)
	getSolanaRockBurnEventsFn   func(ctx context.Context, programID string, lastKnownSig sol.Signature) ([]api.BurnEvent, sol.Signature, error)
	rpcCallBatchFn              func(ctx context.Context, rpcs jsonrpc.RPCRequests) (jsonrpc.RPCResponses, error)
	getTransactionFn            func(ctx context.Context, signature sol.Signature, opts *solana.GetTransactionOpts) (out *solana.GetTransactionResult, err error)
	getSignaturesForAddressFn   func(ctx context.Context, account sol.PublicKey, opts *solana.GetSignaturesForAddressOpts) ([]*solana.TransactionSignature, error)
	reconcileBurnEventsFn       func(ctx context.Context, eventsToClean []api.BurnEvent, cleanedEvents map[string]bool, chainTypeName string) ([]api.BurnEvent, map[string]bool)
}

// CachedTxResult represents a cached transaction result with TTL
type CachedTxResult struct {
	Result    *solana.GetTransactionResult
	ExpiresAt time.Time
}

type oracleStateUpdate struct {
	redemptions             []api.Redemption
	suggestedTip            *big.Int
	estimatedGas            uint64
	ethBurnEvents           []api.BurnEvent
	cleanedEthBurnEvents    map[string]bool
	solanaBurnEvents        []api.BurnEvent
	cleanedSolanaBurnEvents map[string]bool
	ROCKUSDPrice            math.LegacyDec
	BTCUSDPrice             math.LegacyDec
	ETHUSDPrice             math.LegacyDec
	ZECUSDPrice             math.LegacyDec
	SolanaMintEvents        []api.SolanaMintEvent
	cleanedSolanaMintEvents map[string]bool
	latestSolanaSigs        map[sidecartypes.SolanaEventType]sol.Signature
	latestEventStoreCursors map[sidecartypes.SolanaEventType]string
	pendingTransactions     map[string]sidecartypes.PendingTxInfo
}

// PriceData represents Gate.io price response (deprecated, kept for reference)
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

// JupiterPriceResponse represents the response from Jupiter Price API v3
// The response is a map where keys are token addresses
type JupiterPriceResponse map[string]JupiterPriceData

// JupiterPriceData represents individual price data from Jupiter API
type JupiterPriceData struct {
	USDPrice       float64 `json:"usdPrice"`
	BlockID        int64   `json:"blockId"`
	Decimals       int     `json:"decimals"`
	PriceChange24h float64 `json:"priceChange24h"`
}
