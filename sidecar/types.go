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
		EigenDelegations:        make(map[string]map[string]*big.Int),
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
		PendingSolanaTxs:        make(map[string]sidecartypes.PendingTxInfo),
	}
)

type Oracle struct {
	currentState       atomic.Value // *types.OracleState
	stateCache         []sidecartypes.OracleState
	Config             sidecartypes.Config
	EthClient          *ethclient.Client
	neutrinoServer     *neutrino.NeutrinoServer
	solanaClient       *solana.Client
	zrChainQueryClient *client.QueryClient
	mainLoopTicker     *time.Ticker
	DebugMode          bool
	SkipInitialWait    bool

	// Last processed Solana signatures (managed as strings for persistence)
	lastSolRockMintSigStr   string
	lastSolZenBTCMintSigStr string
	lastSolZenBTCBurnSigStr string
	lastSolRockBurnSigStr   string

	// Performance optimization fields
	solanaRateLimiter     chan struct{}              // Semaphore for Solana RPC rate limiting
	transactionCache      map[string]*CachedTxResult // Cache for frequently accessed transactions
	transactionCacheMutex sync.RWMutex               // Protects transaction cache

	// Function fields for mocking
	getSolanaZenBTCBurnEventsFn func(ctx context.Context, programID string, lastKnownSig sol.Signature) ([]api.BurnEvent, sol.Signature, error)
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
