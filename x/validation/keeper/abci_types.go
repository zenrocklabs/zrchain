package keeper

import (
	"context"
	"errors"
	"math/big"

	"cosmossdk.io/math"
	sidecar "github.com/Zenrock-Foundation/zrchain/v4/sidecar/proto/api"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
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
		ZRChainBlockHeight int64
		ROCKUSDPrice       math.LegacyDec
		ETHUSDPrice        math.LegacyDec
		AVSDelegationsHash []byte
		BtcBlockHeight     int64
		BtcMerkleRoot      string
		EthBlockHeight     uint64
		EthBlockHash       common.Hash
		EthGasPrice        uint64
		EthGasLimit        uint64
		RequestedEthNonce  uint64
		EthTxHeight        uint64
		SolanaTxSlot       uint64
	}

	VEWithVotePower struct {
		VoteExtension []byte
		VotePower     int64
	}

	OracleData struct {
		ROCKUSDPrice         math.LegacyDec
		ETHUSDPrice          math.LegacyDec
		AVSDelegationsMap    map[string]map[string]*big.Int
		ValidatorDelegations []ValidatorDelegations
		BtcBlockHeight       int64
		BtcBlockHeader       sidecar.BTCBlockHeader
		EthBlockHeight       uint64
		EthBlockHash         common.Hash
		EthGasPrice          uint64
		EthGasLimit          uint64
		RequestedEthNonce    uint64
		EthTxHeight          uint64
		SolanaTxSlot         uint64
		ConsensusData        abci.ExtendedCommitInfo
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
		GetSolanaTransaction(ctx context.Context, in *sidecar.SolanaTransactionRequest, opts ...grpc.CallOption) (*sidecar.SolanaTransactionResponse, error)
		GetEthereumTransaction(ctx context.Context, in *sidecar.EthereumTransactionRequest, opts ...grpc.CallOption) (*sidecar.EthereumTransactionResponse, error)
		GetEthereumNonceAtHeight(ctx context.Context, in *sidecar.EthereumNonceAtHeightRequest, opts ...grpc.CallOption) (*sidecar.EthereumNonceAtHeightResponse, error)
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
		ve.EthGasPrice == 0 ||
		ve.EthGasLimit == 0 ||
		ve.BtcBlockHeight == 0 ||
		len(ve.BtcMerkleRoot) == 0
}
