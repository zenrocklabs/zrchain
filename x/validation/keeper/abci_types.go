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
	// if ve.ETHUSDPrice.IsZero() {
	// 	logger.Error("invalid vote extension: ETHUSDPrice is zero")
	// 	invalid = true
	// }

	return invalid
}
