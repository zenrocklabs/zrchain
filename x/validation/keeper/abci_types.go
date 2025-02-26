package keeper

import (
	"context"
	"errors"
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
		BtcBlockHeight             int64
		BtcHeaderHash              []byte
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
	}

	VEWithVotePower struct {
		VoteExtension []byte
		VotePower     int64
	}

	OracleData struct {
		EigenDelegationsMap        map[string]map[string]*big.Int
		ValidatorDelegations       []ValidatorDelegations
		BtcBlockHeight             int64
		BtcBlockHeader             sidecar.BTCBlockHeader
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
	if ve.BtcBlockHeight == 0 {
		logger.Error("invalid vote extension: BtcBlockHeight is 0")
		invalid = true
	}
	if len(ve.BtcHeaderHash) == 0 {
		logger.Error("invalid vote extension: BtcHeaderHash is empty")
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

	return invalid
}

// HasAnyOracleData returns true if this OracleData contains any meaningful data
// beyond the ConsensusData (which is always present)
func (o OracleData) HasAnyOracleData() bool {
	// Check if the oracle data has any ethereum or bitcoin data
	if o.EthBlockHeight > 0 || o.BtcBlockHeight > 0 {
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
	VEFieldZRChainBlockHeight   VoteExtensionField = iota // Height is always essential
	VEFieldEigenDelegationsHash                           // AVS delegations are essential for validator updates
	VEFieldEthBurnEventsHash                              // Events from Ethereum
	VEFieldRedemptionsHash                                // Redemption data
	VEFieldBtcHeaderHash                                  // Bitcoin header data

	// Standard fields
	VEFieldBtcBlockHeight
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
)

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
	case VEFieldBtcHeaderHash:
		return "BtcHeaderHash"
	case VEFieldBtcBlockHeight:
		return "BtcBlockHeight"
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

// fieldVote represents a voted value with its accumulated voting power
type fieldVote struct {
	value     any
	votePower int64
}
