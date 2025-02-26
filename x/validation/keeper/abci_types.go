package keeper

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
	sidecar "github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	VoteExtBytesLimit = 1024
)

var (
	ACCEPT_VOTE     = &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}
	REJECT_VOTE     = &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_REJECT}
	ACCEPT_PROPOSAL = &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}
	REJECT_PROPOSAL = &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}

	ErrOracleSidecar = errors.New("oracle sidecar error")
)

type (
	VoteExtension struct {
		ZRChainBlockHeight         int64
		EigenDelegationsHash       []byte
		RequestedBtcBlockHeight    int64
		RequestedBtcHeaderHash     []byte
		EthBlockHeight             uint64
		EthGasLimit                uint64
		EthBaseFee                 uint64
		EthTipCap                  uint64
		RequestedStakerNonce       uint64
		RequestedEthMinterNonce    uint64
		RequestedUnstakerNonce     uint64
		RequestedCompleterNonce    uint64
		SolanaLamportsPerSignature uint64
		EthBurnEventsHash          []byte
		RedemptionsHash            []byte
		ROCKUSDPrice               math.LegacyDec
		BTCUSDPrice                math.LegacyDec
		ETHUSDPrice                math.LegacyDec
		LatestBtcBlockHeight       int64
		LatestBtcHeaderHash        []byte
	}

	VEWithVotePower struct {
		VoteExtension []byte
		VotePower     int64
	}

	OracleData struct {
		EigenDelegationsMap        map[string]map[string]*big.Int
		ValidatorDelegations       []ValidatorDelegations
		RequestedBtcBlockHeight    int64
		RequestedBtcBlockHeader    sidecar.BTCBlockHeader
		LatestBtcBlockHeight       int64
		LatestBtcBlockHeader       sidecar.BTCBlockHeader
		EthBlockHeight             uint64
		EthGasLimit                uint64
		EthBaseFee                 uint64
		EthTipCap                  uint64
		RequestedStakerNonce       uint64
		RequestedEthMinterNonce    uint64
		RequestedUnstakerNonce     uint64
		RequestedCompleterNonce    uint64
		SolanaLamportsPerSignature uint64
		EthBurnEvents              []api.BurnEvent
		Redemptions                []api.Redemption
		ROCKUSDPrice               math.LegacyDec
		BTCUSDPrice                math.LegacyDec
		ETHUSDPrice                math.LegacyDec
		ConsensusData              abci.ExtendedCommitInfo
		FieldVotePowers            map[VoteExtensionField]int64 // Track which fields reached consensus
	}

	ValidatorDelegations struct {
		Validator string
		Stake     math.Int
	}

	sidecarClient interface {
		GetSidecarState(ctx context.Context, _ *sidecar.SidecarStateRequest, opts ...grpc.CallOption) (*sidecar.SidecarStateResponse, error)
		GetSidecarStateByEthHeight(ctx context.Context, req *sidecar.SidecarStateByEthHeightRequest, opts ...grpc.CallOption) (*sidecar.SidecarStateResponse, error)
		GetBitcoinBlockHeaderByHeight(ctx context.Context, in *sidecar.BitcoinBlockHeaderByHeightRequest, opts ...grpc.CallOption) (*sidecar.BitcoinBlockHeaderResponse, error)
		GetLatestBitcoinBlockHeader(ctx context.Context, in *sidecar.LatestBitcoinBlockHeaderRequest, opts ...grpc.CallOption) (*sidecar.BitcoinBlockHeaderResponse, error)
		GetLatestEthereumNonceForAccount(ctx context.Context, in *sidecar.LatestEthereumNonceForAccountRequest, opts ...grpc.CallOption) (*sidecar.LatestEthereumNonceForAccountResponse, error)
	}
)

func NewSidecarClient(serverAddr string) (sidecarClient, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return sidecar.NewSidecarServiceClient(conn), nil
}

func ContainsVoteExtension(tx []byte, unmarshalTx sdk.TxDecoder) bool {
	if _, err := unmarshalTx(tx); err != nil {
		return true
	}
	return false
}

func VoteExtensionsEnabled(ctx sdk.Context) bool {
	// The genesis block should not have vote extensions enabled
	if ctx.BlockHeight() == 1 {
		return false
	}

	consParams := ctx.ConsensusParams()
	if consParams.Abci == nil || consParams.Abci.VoteExtensionsEnableHeight == 0 {
		return false
	}

	return ctx.BlockHeight() > consParams.Abci.VoteExtensionsEnableHeight
}

func (ve VoteExtension) IsInvalid(logger log.Logger) bool {
	invalid := false

	if ve.ZRChainBlockHeight == 0 {
		logger.Error("invalid vote extension: ZRChainBlockHeight is 0")
		invalid = true
	}
	if len(ve.EigenDelegationsHash) == 0 {
		logger.Error("invalid vote extension: EigenDelegationsHash is empty")
		invalid = true
	}
	if ve.EthBlockHeight == 0 {
		logger.Error("invalid vote extension: EthBlockHeight is 0")
		invalid = true
	}
	if ve.EthBaseFee == 0 {
		logger.Error("invalid vote extension: EthBaseFee is 0")
		invalid = true
	}
	if ve.EthTipCap == 0 {
		logger.Error("invalid vote extension: EthTipCap is 0")
		invalid = true
	}
	if ve.EthGasLimit == 0 {
		logger.Error("invalid vote extension: EthGasLimit is 0")
		invalid = true
	}
	if ve.RequestedBtcBlockHeight == 0 {
		logger.Error("invalid vote extension: RequestedBtcBlockHeight is 0")
		invalid = true
	}
	if len(ve.RequestedBtcHeaderHash) == 0 {
		logger.Error("invalid vote extension: RequestedBtcHeaderHash is empty")
		invalid = true
	}
	if ve.SolanaLamportsPerSignature == 0 {
		logger.Error("invalid vote extension: SolanaLamportsPerSignature is 0")
		invalid = true
	}
	if len(ve.EthBurnEventsHash) == 0 {
		logger.Error("invalid vote extension: EthBurnEventsHash is empty")
		invalid = true
	}
	if len(ve.RedemptionsHash) == 0 {
		logger.Error("invalid vote extension: RedemptionsHash is empty")
		invalid = true
	}
	if ve.ROCKUSDPrice.IsNil() || ve.ROCKUSDPrice.IsZero() {
		logger.Error("invalid vote extension: ROCKUSDPrice is nil or zero")
		invalid = true
	}
	if ve.BTCUSDPrice.IsNil() || ve.BTCUSDPrice.IsZero() {
		logger.Error("invalid vote extension: BTCUSDPrice is nil or zero")
		invalid = true
	}
	if ve.ETHUSDPrice.IsZero() {
		logger.Error("invalid vote extension: ETHUSDPrice is zero")
		invalid = true
	}
	if ve.LatestBtcBlockHeight == 0 {
		logger.Error("invalid vote extension: LatestBtcBlockHeight is 0")
		invalid = true
	}
	if len(ve.LatestBtcHeaderHash) == 0 {
		logger.Error("invalid vote extension: LatestBtcHeaderHash is empty")
		invalid = true
	}

	return invalid
}

// HasAnyOracleData returns true if this OracleData contains any meaningful data
// beyond the ConsensusData (which is always present)
func (o OracleData) HasAnyOracleData() bool {
	// Check if the oracle data has any ethereum or bitcoin data
	if o.EthBlockHeight > 0 || o.RequestedBtcBlockHeight > 0 {
		return true
	}

	// Check for ethereum burn events
	if len(o.EthBurnEvents) > 0 {
		return true
	}

	// Check for redemptions
	if len(o.Redemptions) > 0 {
		return true
	}

	// Check for eigen delegations
	if len(o.EigenDelegationsMap) > 0 {
		return true
	}

	return false
}

// VoteExtensionField defines a type-safe identifier for vote extension fields
type VoteExtensionField int

const (
	// Essential fields - absence of these can cause data loss
	VEFieldZRChainBlockHeight     VoteExtensionField = iota // Height is always essential
	VEFieldEigenDelegationsHash                             // AVS delegations are essential for validator updates
	VEFieldEthBurnEventsHash                                // Events from Ethereum
	VEFieldRedemptionsHash                                  // Redemption data
	VEFieldRequestedBtcHeaderHash                           // Bitcoin header data

	// Standard fields
	VEFieldRequestedBtcBlockHeight
	VEFieldEthBlockHeight
	VEFieldEthGasLimit
	VEFieldEthBaseFee
	VEFieldEthTipCap
	VEFieldSolanaLamportsPerSignature
	VEFieldRequestedStakerNonce
	VEFieldRequestedEthMinterNonce
	VEFieldRequestedUnstakerNonce
	VEFieldRequestedCompleterNonce
	VEFieldROCKUSDPrice
	VEFieldBTCUSDPrice
	VEFieldETHUSDPrice
	VEFieldLatestBtcBlockHeight
	VEFieldLatestBtcHeaderHash
)

// FieldHandler defines operations for processing a specific vote extension field
type FieldHandler struct {
	Field    VoteExtensionField
	GetValue func(ve VoteExtension) any
	SetValue func(v any, ve *VoteExtension)
}

// fieldVote represents a voted value with its accumulated voting power
type fieldVote struct {
	value     any
	votePower int64
}

// genericGetKey marshals a value to JSON for use as a map key
func genericGetKey(value any) string {
	if value == nil {
		return ""
	}

	// Handle byte slices specially since they need consistent representation
	if byteSlice, ok := value.([]byte); ok {
		return hex.EncodeToString(byteSlice)
	}

	// For other types, use JSON
	bytes, err := json.Marshal(value)
	if err != nil {
		// Fall back to string representation if marshaling fails
		return fmt.Sprintf("%v", value)
	}
	return string(bytes)
}

// String returns the string representation of a VoteExtensionField
func (f VoteExtensionField) String() string {
	switch f {
	case VEFieldZRChainBlockHeight:
		return "ZRChainBlockHeight"
	case VEFieldEigenDelegationsHash:
		return "EigenDelegationsHash"
	case VEFieldEthBurnEventsHash:
		return "EthBurnEventsHash"
	case VEFieldRedemptionsHash:
		return "RedemptionsHash"
	case VEFieldRequestedBtcHeaderHash:
		return "RequestedBtcHeaderHash"
	case VEFieldRequestedBtcBlockHeight:
		return "RequestedBtcBlockHeight"
	case VEFieldEthBlockHeight:
		return "EthBlockHeight"
	case VEFieldEthGasLimit:
		return "EthGasLimit"
	case VEFieldEthBaseFee:
		return "EthBaseFee"
	case VEFieldEthTipCap:
		return "EthTipCap"
	case VEFieldSolanaLamportsPerSignature:
		return "SolanaLamportsPerSignature"
	case VEFieldRequestedStakerNonce:
		return "RequestedStakerNonce"
	case VEFieldRequestedEthMinterNonce:
		return "RequestedEthMinterNonce"
	case VEFieldRequestedUnstakerNonce:
		return "RequestedUnstakerNonce"
	case VEFieldRequestedCompleterNonce:
		return "RequestedCompleterNonce"
	case VEFieldROCKUSDPrice:
		return "ROCKUSDPrice"
	case VEFieldBTCUSDPrice:
		return "BTCUSDPrice"
	case VEFieldETHUSDPrice:
		return "ETHUSDPrice"
	case VEFieldLatestBtcBlockHeight:
		return "LatestBtcBlockHeight"
	case VEFieldLatestBtcHeaderHash:
		return "LatestBtcHeaderHash"
	default:
		return "Unknown"
	}
}

// EssentialVoteExtensionFields defines fields where absence of consensus
// would cause data loss or chain issues
var EssentialVoteExtensionFields = []VoteExtensionField{
	VEFieldZRChainBlockHeight,   // Height is always essential
	VEFieldEigenDelegationsHash, // AVS delegations are essential for validator updates
	// Other fields can be made essential as needed in the future
}

// IsEssentialField returns true if the given field is considered essential
func IsEssentialField(field VoteExtensionField) bool {
	return slices.Contains(EssentialVoteExtensionFields, field)
}

// HasAllEssentialFields checks if all essential fields have reached consensus
func HasAllEssentialFields(fieldVotePowers map[VoteExtensionField]int64) bool {
	for _, field := range EssentialVoteExtensionFields {
		if _, ok := fieldVotePowers[field]; !ok {
			return false
		}
	}
	return true
}

// InitializeFieldHandlers creates handlers for all vote extension fields
func initializeFieldHandlers() []FieldHandler {
	return []FieldHandler{
		// Hash fields
		{
			Field:    VEFieldEigenDelegationsHash,
			GetValue: func(ve VoteExtension) any { return ve.EigenDelegationsHash },
			SetValue: func(v any, ve *VoteExtension) { ve.EigenDelegationsHash = v.([]byte) },
		},
		{
			Field:    VEFieldEthBurnEventsHash,
			GetValue: func(ve VoteExtension) any { return ve.EthBurnEventsHash },
			SetValue: func(v any, ve *VoteExtension) { ve.EthBurnEventsHash = v.([]byte) },
		},
		{
			Field:    VEFieldRedemptionsHash,
			GetValue: func(ve VoteExtension) any { return ve.RedemptionsHash },
			SetValue: func(v any, ve *VoteExtension) { ve.RedemptionsHash = v.([]byte) },
		},
		{
			Field:    VEFieldRequestedBtcHeaderHash,
			GetValue: func(ve VoteExtension) any { return ve.RequestedBtcHeaderHash },
			SetValue: func(v any, ve *VoteExtension) { ve.RequestedBtcHeaderHash = v.([]byte) },
		},
		{
			Field:    VEFieldLatestBtcHeaderHash,
			GetValue: func(ve VoteExtension) any { return ve.LatestBtcHeaderHash },
			SetValue: func(v any, ve *VoteExtension) { ve.LatestBtcHeaderHash = v.([]byte) },
		},

		// Integer fields
		{
			Field:    VEFieldRequestedBtcBlockHeight,
			GetValue: func(ve VoteExtension) any { return ve.RequestedBtcBlockHeight },
			SetValue: func(v any, ve *VoteExtension) { ve.RequestedBtcBlockHeight = v.(int64) },
		},
		{
			Field:    VEFieldEthBlockHeight,
			GetValue: func(ve VoteExtension) any { return ve.EthBlockHeight },
			SetValue: func(v any, ve *VoteExtension) { ve.EthBlockHeight = v.(uint64) },
		},
		{
			Field:    VEFieldEthGasLimit,
			GetValue: func(ve VoteExtension) any { return ve.EthGasLimit },
			SetValue: func(v any, ve *VoteExtension) { ve.EthGasLimit = v.(uint64) },
		},
		{
			Field:    VEFieldEthBaseFee,
			GetValue: func(ve VoteExtension) any { return ve.EthBaseFee },
			SetValue: func(v any, ve *VoteExtension) { ve.EthBaseFee = v.(uint64) },
		},
		{
			Field:    VEFieldEthTipCap,
			GetValue: func(ve VoteExtension) any { return ve.EthTipCap },
			SetValue: func(v any, ve *VoteExtension) { ve.EthTipCap = v.(uint64) },
		},
		{
			Field:    VEFieldSolanaLamportsPerSignature,
			GetValue: func(ve VoteExtension) any { return ve.SolanaLamportsPerSignature },
			SetValue: func(v any, ve *VoteExtension) { ve.SolanaLamportsPerSignature = v.(uint64) },
		},
		{
			Field:    VEFieldRequestedStakerNonce,
			GetValue: func(ve VoteExtension) any { return ve.RequestedStakerNonce },
			SetValue: func(v any, ve *VoteExtension) { ve.RequestedStakerNonce = v.(uint64) },
		},
		{
			Field:    VEFieldRequestedEthMinterNonce,
			GetValue: func(ve VoteExtension) any { return ve.RequestedEthMinterNonce },
			SetValue: func(v any, ve *VoteExtension) { ve.RequestedEthMinterNonce = v.(uint64) },
		},
		{
			Field:    VEFieldRequestedUnstakerNonce,
			GetValue: func(ve VoteExtension) any { return ve.RequestedUnstakerNonce },
			SetValue: func(v any, ve *VoteExtension) { ve.RequestedUnstakerNonce = v.(uint64) },
		},
		{
			Field:    VEFieldRequestedCompleterNonce,
			GetValue: func(ve VoteExtension) any { return ve.RequestedCompleterNonce },
			SetValue: func(v any, ve *VoteExtension) { ve.RequestedCompleterNonce = v.(uint64) },
		},
		{
			Field:    VEFieldLatestBtcBlockHeight,
			GetValue: func(ve VoteExtension) any { return ve.LatestBtcBlockHeight },
			SetValue: func(v any, ve *VoteExtension) { ve.LatestBtcBlockHeight = v.(int64) },
		},

		// Decimal fields
		{
			Field:    VEFieldROCKUSDPrice,
			GetValue: func(ve VoteExtension) any { return ve.ROCKUSDPrice },
			SetValue: func(v any, ve *VoteExtension) { ve.ROCKUSDPrice = v.(math.LegacyDec) },
		},
		{
			Field:    VEFieldBTCUSDPrice,
			GetValue: func(ve VoteExtension) any { return ve.BTCUSDPrice },
			SetValue: func(v any, ve *VoteExtension) { ve.BTCUSDPrice = v.(math.LegacyDec) },
		},
		{
			Field:    VEFieldETHUSDPrice,
			GetValue: func(ve VoteExtension) any { return ve.ETHUSDPrice },
			SetValue: func(v any, ve *VoteExtension) { ve.ETHUSDPrice = v.(math.LegacyDec) },
		},
	}
}
