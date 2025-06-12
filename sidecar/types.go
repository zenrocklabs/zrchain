package main

import (
	"math/big"
	"sync/atomic"
	"time"

	"cosmossdk.io/math"
	client "github.com/Zenrock-Foundation/zrchain/v6/go-client"
	neutrino "github.com/Zenrock-Foundation/zrchain/v6/sidecar/neutrino"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
)

// TODO: These event structs are temporary. They should be defined in `sidecar/proto/api/types.proto`
// and the generated Go files should be used instead (e.g., `api.EthStakeEvent`).
type EthStakeEvent struct {
	UnsignedTxHash []byte
}
type EthMintEvent struct {
	UnsignedTxHash []byte
}
type EthUnstakeEvent struct {
	UnsignedTxHash []byte
}
type EthCompletionEvent struct {
	UnsignedTxHash []byte
}

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
		EthStakeEvents:             []sidecartypes.EthStakeEvent{},
		EthMintEvents:              []sidecartypes.EthMintEvent{},
		EthUnstakeEvents:           []sidecartypes.EthUnstakeEvent{},
		EthCompletionEvents:        []sidecartypes.EthCompletionEvent{},
	}
)

type Oracle struct {
	currentState       atomic.Value // *types.OracleState
	stateCache         []sidecartypes.OracleState
	Config             sidecartypes.Config
	EthClient          *ethclient.Client
	neutrinoServer     *neutrino.NeutrinoServer
	solanaClient       *solrpc.Client
	zrChainQueryClient *client.QueryClient
	updateChan         chan sidecartypes.OracleState
	mainLoopTicker     *time.Ticker
	DebugMode          bool

	// Last processed Solana signatures (managed as strings for persistence)
	lastSolRockMintSigStr   string
	lastSolZenBTCMintSigStr string
	lastSolZenBTCBurnSigStr string
	lastSolRockBurnSigStr   string

	// Event caches to prevent re-processing
	cleanedEthBurnEvents map[string]bool
}

type oracleStateUpdate struct {
	eigenDelegations           map[string]map[string]*big.Int
	redemptions                []api.Redemption
	suggestedTip               *big.Int
	estimatedGas               uint64
	ethBurnEvents              []api.BurnEvent
	solanaBurnEvents           []api.BurnEvent
	ROCKUSDPrice               math.LegacyDec
	BTCUSDPrice                math.LegacyDec
	ETHUSDPrice                math.LegacyDec
	solanaLamportsPerSignature uint64
	SolanaMintEvents           []api.SolanaMintEvent
	latestSolanaSigs           map[sidecartypes.SolanaEventType]solana.Signature
	ethStakeEvents             []sidecartypes.EthStakeEvent
	ethMintEvents              []sidecartypes.EthMintEvent
	ethUnstakeEvents           []sidecartypes.EthUnstakeEvent
	ethCompletionEvents        []sidecartypes.EthCompletionEvent
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
