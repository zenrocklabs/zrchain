package keeper

import (
	"context"
	"errors"
	"math/big"

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
		BtcMerkleRoot              string // TODO: change to hash of entire header?
		EthBlockHeight             uint64
		EthGasLimit                uint64
		EthBaseFee                 uint64
		EthTipCap                  uint64
		RequestedEthNonce          uint64
		SolanaLamportsPerSignature uint64
		EthereumRedemptionsHash    []byte
		SolanaRedemptionsHash      []byte
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
		RequestedEthNonce          uint64
		SolanaLamportsPerSignature uint64
		EthereumRedemptions        []api.Redemption
		SolanaRedemptions          []api.Redemption
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

func (ve VoteExtension) IsInvalid() bool { // Sasha: Should bitcoin fields be checked here? They're not as critical so maybe not
	return ve.ZRChainBlockHeight == 0 ||
		// TODO: uncomment this after TGE
		// voteExt.ROCKUSDPrice.IsZero() ||
		ve.ETHUSDPrice.IsZero() ||
		len(ve.AVSDelegationsHash) == 0 ||
		ve.EthBlockHeight == 0 ||
		len(ve.EthBlockHash) == 0 ||
		ve.EthBaseFee == 0 ||
		ve.EthTipCap == 0 ||
		ve.EthGasLimit == 0
	// ve.BtcBlockHeight == 0 ||
	// len(ve.BtcMerkleRoot) == 0
}
