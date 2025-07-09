package keeper_test

import (
	"encoding/json"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/keeper"
	validationtestutil "github.com/Zenrock-Foundation/zrchain/v6/x/validation/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Helper function to create a reusable last commit
func createTestLastCommit() abci.CommitInfo {
	return abci.CommitInfo{
		Round: 1,
		Votes: []abci.VoteInfo{
			{
				Validator: abci.Validator{
					Address: []byte("test-validator"),
					Power:   1000,
				},
				BlockIdFlag: 1,
			},
		},
	}
}

// Helper function to create a reusable extended last commit
func createTestExtendedLastCommit() abci.ExtendedCommitInfo {
	return abci.ExtendedCommitInfo{
		Round: 1,
		Votes: []abci.ExtendedVoteInfo{
			{
				Validator: abci.Validator{
					Address: []byte("test-validator"),
					Power:   1000,
				},
				BlockIdFlag: 1,
			},
		},
	}
}

func TestBeginBlocker(t *testing.T) {
	suite := new(ValidationKeeperTestSuite)
	suite.SetT(&testing.T{})
	keeper, _ := suite.ValidationKeeperSetupTest()
	ctx := sdk.UnwrapSDKContext(suite.ctx)

	err := keeper.BeginBlocker(ctx)
	require.NoError(t, err)
}

func TestEndBlocker(t *testing.T) {
	suite := new(ValidationKeeperTestSuite)
	suite.SetT(&testing.T{})
	keeper, _ := suite.ValidationKeeperSetupTest()
	ctx := sdk.UnwrapSDKContext(suite.ctx)

	_, err := keeper.EndBlocker(ctx)
	require.NoError(t, err)
}

func TestExtendVoteHandler(t *testing.T) {

	type args struct {
		req *abci.RequestExtendVote
	}
	tests := []struct {
		name    string
		args    args
		want    *abci.ResponseExtendVote
		wantErr bool
	}{
		{
			name: "PASS: extend vote handler",
			args: args{
				req: &abci.RequestExtendVote{
					Hash:               []byte("test"),
					Height:             1,
					Time:               time.Now(),
					Txs:                [][]byte{[]byte(`{"EigenDelegationsHash":"uhVXdw9X1G/iRkwfVMBjUFFsCgsB33yWKu4h5ierVJI=","EthBaseFee":3732027422,"EthBlockHeight":22796583,"EthBurnEventsHash":"dCNOmK/nSY+12vHzasLXiswzlGT5UHA7jAGYkvmCuQs=","EthGasLimit":249091,"EthTipCap":72578,"LatestBtcBlockHeight":902951,"LatestBtcHeaderHash":"uPjzvaQD965jAViGFwf7CUtMrY7EwhHyvWpHDMeOU6Y=","ROCKUSDPrice":"0.047030000000000000","RedemptionsHash":"dCNOmK/nSY+12vHzasLXiswzlGT5UHA7jAGYkvmCuQs=","RequestedBtcBlockHeight":0,"RequestedBtcHeaderHash":null,"RequestedCompleterNonce":0,"RequestedEthMinterNonce":0,"RequestedStakerNonce":0,"RequestedUnstakerNonce":0,"SidecarVersionName":"rose_moon","SolanaAccountsHash":"RBNvo1WzZ4oRRq0W9+hknpT7T8If536DEMBg9hyq/4o=","SolanaBurnEventsHash":"dCNOmK/nSY+12vHzasLXiswzlGT5UHA7jAGYkvmCuQs=","SolanaLamportsPerSignature":0,"SolanaMintEventsHash":"Zp729xYaghztbJRLKnyJfwyGnIlbMvMeV2CNm9/5Li0=","SolanaMintNoncesHash":"RBNvo1WzZ4oRRq0W9+hknpT7T8If536DEMBg9hyq/4o=","ZRChainBlockHeight":3401684}`)},
					ProposedLastCommit: createTestLastCommit(),
					Misbehavior:        nil,
					NextValidatorsHash: []byte("test-next-validators-hash"),
					ProposerAddress:    []byte("test-proposer-address"),
				},
			},
			want: &abci.ResponseExtendVote{VoteExtension: validationtestutil.SampleVoteExtension},
		},
		{
			name: "PASS: extend vote handler",
			args: args{
				req: &abci.RequestExtendVote{
					Hash:               []byte("test"),
					Height:             2,
					Time:               time.Now(),
					Txs:                [][]byte{[]byte(`{"EigenDelegationsHash":"uhVXdw9X1G/iRkwfVMBjUFFsCgsB33yWKu4h5ierVJI=","EthBaseFee":3732027422,"EthBlockHeight":22796583,"EthBurnEventsHash":"dCNOmK/nSY+12vHzasLXiswzlGT5UHA7jAGYkvmCuQs=","EthGasLimit":249091,"EthTipCap":72578,"LatestBtcBlockHeight":902951,"LatestBtcHeaderHash":"uPjzvaQD965jAViGFwf7CUtMrY7EwhHyvWpHDMeOU6Y=","ROCKUSDPrice":"0.047030000000000000","RedemptionsHash":"dCNOmK/nSY+12vHzasLXiswzlGT5UHA7jAGYkvmCuQs=","RequestedBtcBlockHeight":0,"RequestedBtcHeaderHash":null,"RequestedCompleterNonce":0,"RequestedEthMinterNonce":0,"RequestedStakerNonce":0,"RequestedUnstakerNonce":0,"SidecarVersionName":"rose_moon","SolanaAccountsHash":"RBNvo1WzZ4oRRq0W9+hknpT7T8If536DEMBg9hyq/4o=","SolanaBurnEventsHash":"dCNOmK/nSY+12vHzasLXiswzlGT5UHA7jAGYkvmCuQs=","SolanaLamportsPerSignature":0,"SolanaMintEventsHash":"Zp729xYaghztbJRLKnyJfwyGnIlbMvMeV2CNm9/5Li0=","SolanaMintNoncesHash":"RBNvo1WzZ4oRRq0W9+hknpT7T8If536DEMBg9hyq/4o=","ZRChainBlockHeight":3401684}`)},
					ProposedLastCommit: createTestLastCommit(),
					Misbehavior:        nil,
					NextValidatorsHash: []byte("test-next-validators-hash"),
					ProposerAddress:    []byte("test-proposer-address"),
				},
			},
			want: &abci.ResponseExtendVote{VoteExtension: validationtestutil.SampleVoteExtensionHeight2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			suite := new(ValidationKeeperTestSuite)
			suite.SetT(&testing.T{})
			keeper, ctrl := suite.ValidationKeeperSetupTest()
			defer ctrl.Finish()

			ctx := sdk.UnwrapSDKContext(suite.ctx)

			got, err := keeper.ExtendVoteHandler(ctx, tt.args.req)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}

}

func TestVerifyVoteExtensionHandler(t *testing.T) {
	type args struct {
		req *abci.RequestVerifyVoteExtension
	}
	tests := []struct {
		name    string
		args    args
		want    *abci.ResponseVerifyVoteExtension
		wantErr bool
	}{
		{
			name: "PASS: verify vote extension handler",
			args: args{
				req: &abci.RequestVerifyVoteExtension{
					Hash:             []byte("test"),
					ValidatorAddress: []byte("test-validator"),
					Height:           1,
					VoteExtension:    validationtestutil.SampleVoteExtension,
				},
			},
			want: &abci.ResponseVerifyVoteExtension{
				Status: abci.ResponseVerifyVoteExtension_ACCEPT,
			},
		},
		{
			name: "FAIL: vote extension with height mismatch",
			args: args{
				req: &abci.RequestVerifyVoteExtension{
					Hash:             []byte("test"),
					ValidatorAddress: []byte("test-validator"),
					Height:           2, // Different from vote extension height
					VoteExtension:    validationtestutil.SampleVoteExtension,
				},
			},
			want: &abci.ResponseVerifyVoteExtension{
				Status: abci.ResponseVerifyVoteExtension_REJECT,
			},
		},
		{
			name: "FAIL: vote extension too large",
			args: args{
				req: &abci.RequestVerifyVoteExtension{
					Hash:             []byte("test"),
					ValidatorAddress: []byte("test-validator"),
					Height:           1,
					VoteExtension:    make([]byte, 1025), // Exceeds VoteExtBytesLimit of 1024
				},
			},
			want: &abci.ResponseVerifyVoteExtension{
				Status: abci.ResponseVerifyVoteExtension_REJECT,
			},
		},
		{
			name: "FAIL: invalid JSON vote extension",
			args: args{
				req: &abci.RequestVerifyVoteExtension{
					Hash:             []byte("test"),
					ValidatorAddress: []byte("test-validator"),
					Height:           1,
					VoteExtension:    []byte("invalid json"),
				},
			},
			want: &abci.ResponseVerifyVoteExtension{
				Status: abci.ResponseVerifyVoteExtension_REJECT,
			},
		},
		{
			name: "FAIL: invalid vote extension data",
			args: args{
				req: &abci.RequestVerifyVoteExtension{
					Hash:             []byte("test"),
					ValidatorAddress: []byte("test-validator"),
					Height:           1,
					VoteExtension:    []byte(`{"ZRChainBlockHeight":0,"EigenDelegationsHash":"","EthBlockHeight":0,"EthBaseFee":0,"EthTipCap":0,"EthGasLimit":0,"EthBurnEventsHash":"","RedemptionsHash":"","ROCKUSDPrice":"","BTCUSDPrice":"","ETHUSDPrice":"","LatestBtcBlockHeight":0,"LatestBtcHeaderHash":""}`),
				},
			},
			want: &abci.ResponseVerifyVoteExtension{
				Status: abci.ResponseVerifyVoteExtension_REJECT,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := new(ValidationKeeperTestSuite)
			suite.SetT(&testing.T{})
			keeper, ctrl := suite.ValidationKeeperSetupTest()
			defer ctrl.Finish()

			resp, err := keeper.VerifyVoteExtensionHandler(sdk.UnwrapSDKContext(suite.ctx), tt.args.req)
			require.NoError(t, err)
			require.Equal(t, tt.want, resp)
		})
	}
}

func TestPrepareProposal(t *testing.T) {

	// WIP - requires proper vote extension and oracle data to fully test

	voteExt := keeper.VoteExtension{
		EthBlockHeight: 12345,
		EthGasLimit:    21000,
		EthBaseFee:     20000000000,
		EthTipCap:      1000000000,
		BTCUSDPrice:    "50000.00",
		ETHUSDPrice:    "3000.00",
		ROCKUSDPrice:   "1.00",
	}
	voteExtBytes, err := json.Marshal(voteExt)
	require.NoError(t, err)

	type args struct {
		req *abci.RequestPrepareProposal
	}
	tests := []struct {
		name    string
		args    args
		want    *abci.ResponsePrepareProposal
		wantErr bool
	}{
		{
			name: "test with consensus",
			args: args{
				req: &abci.RequestPrepareProposal{
					MaxTxBytes: 4000,
					Txs:        [][]byte{},
					Height:     2,
					Time:       time.Now(),
					LocalLastCommit: abci.ExtendedCommitInfo{
						Round: 1,
						Votes: []abci.ExtendedVoteInfo{
							{
								Validator: abci.Validator{
									Address: []byte("QDagxuKQqu3HMpWLmNIgCEhR9b0="),
									Power:   1000000,
								},
								BlockIdFlag:   1,
								VoteExtension: voteExtBytes,
							},
						},
					},
					Misbehavior:        nil,
					NextValidatorsHash: []byte("test-next-validators-hash"),
					ProposerAddress:    []byte("test-proposer-address"),
				},
			},
			want: &abci.ResponsePrepareProposal{
				Txs: [][]byte{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := new(ValidationKeeperTestSuite)
			suite.SetT(&testing.T{})
			keeper, ctrl := suite.ValidationKeeperSetupTest()
			defer ctrl.Finish()

			ctx := sdk.UnwrapSDKContext(suite.ctx)
			ctx = ctx.WithBlockHeight(2)
			consensusParams := ctx.ConsensusParams()
			consensusParams.Abci = &cmtproto.ABCIParams{
				VoteExtensionsEnableHeight: 1,
			}
			ctx = ctx.WithConsensusParams(consensusParams)

			resp, err := keeper.PrepareProposal(ctx, tt.args.req)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotEmpty(t, resp)

			// Verify it's valid JSON
			var oracleData map[string]interface{}
			err = json.Unmarshal(resp, &oracleData)
			require.NoError(t, err)

			require.Contains(t, oracleData, "ConsensusData")
			require.Contains(t, oracleData, "FieldVotePowers")

			_, ok := oracleData["FieldVotePowers"].(map[string]interface{})
			require.True(t, ok)
		})
	}
}

func TestProcessProposal(t *testing.T) {

	type args struct {
		req *abci.RequestProcessProposal
	}
	tests := []struct {
		name    string
		args    args
		want    *abci.ResponseProcessProposal
		wantErr bool
	}{
		{
			name: "PASS: process proposal with vote extensions enabled, but with empty oracle data",
			args: args{
				req: &abci.RequestProcessProposal{
					Txs:    [][]byte{[]byte(`{"EigenDelegationsMap":null,"ValidatorDelegations":null,"RequestedBtcBlockHeight":0,"RequestedBtcBlockHeader":{},"LatestBtcBlockHeight":0,"LatestBtcBlockHeader":{},"EthBlockHeight":0,"EthGasLimit":0,"EthBaseFee":0,"EthTipCap":0,"RequestedStakerNonce":0,"RequestedEthMinterNonce":0,"RequestedUnstakerNonce":0,"RequestedCompleterNonce":0,"SolanaMintNonces":null,"SolanaAccounts":null,"SolanaLamportsPerSignature":0,"SolanaMintEvents":null,"SolanaZenBTCMintEvents":null,"EthBurnEvents":null,"SolanaBurnEvents":null,"Redemptions":null,"ROCKUSDPrice":"","BTCUSDPrice":"","ETHUSDPrice":"","ConsensusData":{"votes":[{"validator":{"address":"QDagxuKQqu3HMpWLmNIgCEhR9b0=","power":125000000},"extension_signature":"QB/lPpqzBJAW+iNF37X5PVrHpuHJ/ZmKWcFX6JdwTxYPAjomEHI9BqzF9EOSpp3CQ1/OikFMlITSR+eqIhgaCg==","block_id_flag":2}]},"FieldVotePowers":{}}`)},
					Height: 2,
					Hash:   []byte("test"),
					Time:   time.Now(),
					ProposedLastCommit: abci.CommitInfo{
						Round: 1,
						Votes: []abci.VoteInfo{
							{
								Validator: abci.Validator{
									Address: []byte("QDagxuKQqu3HMpWLmNIgCEhR9b0="),
									Power:   1000000,
								},
								BlockIdFlag: 1,
							},
						},
					},
					Misbehavior:        nil,
					NextValidatorsHash: []byte("test-next-validators-hash"),
					ProposerAddress:    []byte("test-proposer-address"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := new(ValidationKeeperTestSuite)
			suite.SetT(&testing.T{})
			keeper, ctrl := suite.ValidationKeeperSetupTest()
			defer ctrl.Finish()

			// Get the ubermock controller from the suite and finish it too
			if suite.zenBTCCtrl != nil {
				defer suite.zenBTCCtrl.Finish()
			}

			ctx := sdk.UnwrapSDKContext(suite.ctx)
			// Set block height to 2 (not 1, since vote extensions are disabled for height 1)
			ctx = ctx.WithBlockHeight(2)

			// Enable vote extensions by setting consensus parameters
			consensusParams := ctx.ConsensusParams()
			consensusParams.Abci = &cmtproto.ABCIParams{
				VoteExtensionsEnableHeight: 1, // Enable from height 1 onwards
			}
			ctx = ctx.WithConsensusParams(consensusParams)

			resp, err := keeper.ProcessProposal(ctx, tt.args.req)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotEmpty(t, resp)
		})
	}
}

func TestPreBlocker(t *testing.T) {

	type args struct {
		req         *abci.RequestFinalizeBlock
		blockHeight int64
	}
	tests := []struct {
		name    string
		args    args
		want    *abci.ResponseFinalizeBlock
		wantErr bool
	}{
		{
			name: "PASS: pre blocker with vote extensions enabled and consensus data",
			args: args{
				req: &abci.RequestFinalizeBlock{
					Txs:    [][]byte{[]byte(`{"BTCUSDPrice":"106603.530000000000000000","ETHUSDPrice":"2422.093500000000000000","EigenDelegationsHash":"uhVXdw9X1G/iRkwfVMBjUFFsCgsB33yWKu4h5ierVJI=","EthBaseFee":3732027422,"EthBlockHeight":22796583,"EthBurnEventsHash":"dCNOmK/nSY+12vHzasLXiswzlGT5UHA7jAGYkvmCuQs=","EthGasLimit":249091,"EthTipCap":72578,"LatestBtcBlockHeight":902951,"LatestBtcHeaderHash":"uPjzvaQD965jAViGFwf7CUtMrY7EwhHyvWpHDMeOU6Y=","ROCKUSDPrice":"0.047030000000000000","RedemptionsHash":"dCNOmK/nSY+12vHzasLXiswzlGT5UHA7jAGYkvmCuQs=","RequestedBtcBlockHeight":0,"RequestedBtcHeaderHash":null,"RequestedCompleterNonce":0,"RequestedEthMinterNonce":0,"RequestedStakerNonce":0,"RequestedUnstakerNonce":0,"SidecarVersionName":"rose_moon","SolanaAccountsHash":"RBNvo1WzZ4oRRq0W9+hknpT7T8If536DEMBg9hyq/4o=","SolanaBurnEventsHash":"dCNOmK/nSY+12vHzasLXiswzlGT5UHA7jAGYkvmCuQs=","SolanaLamportsPerSignature":0,"SolanaMintEventsHash":"Zp729xYaghztbJRLKnyJfwyGnIlbMvMeV2CNm9/5Li0=","SolanaMintNoncesHash":"RBNvo1WzZ4oRRq0W9+hknpT7T8If536DEMBg9hyq/4o=","ZRChainBlockHeight":3401684}`)},
					Height: 3,
					Time:   time.Now(),
					DecidedLastCommit: abci.CommitInfo{
						Round: 1,
						Votes: []abci.VoteInfo{
							{
								Validator: abci.Validator{
									Address: []byte("QDagxuKQqu3HMpWLmNIgCEhR9b0="),
									Power:   1000000,
								},
								BlockIdFlag: 1,
							},
						},
					},
					Misbehavior:        nil,
					NextValidatorsHash: []byte("test-next-validators-hash"),
					ProposerAddress:    []byte("test-proposer-address"),
				},
				blockHeight: 3,
			},
		},
		{
			name: "PASS: pre blocker with vote extensions enabled and consensus data",
			args: args{
				req: &abci.RequestFinalizeBlock{
					Txs:    [][]byte{[]byte(validationtestutil.VoteExt)},
					Height: 3,
					Time:   time.Now(),
					DecidedLastCommit: abci.CommitInfo{
						Round: 1,
						Votes: []abci.VoteInfo{
							{
								Validator: abci.Validator{
									Address: []byte("QDagxuKQqu3HMpWLmNIgCEhR9b0="),
									Power:   1000000,
								},
								BlockIdFlag: 1,
							},
						},
					},
					Misbehavior:        nil,
					NextValidatorsHash: []byte("test-next-validators-hash"),
					ProposerAddress:    []byte("test-proposer-address"),
				},
				blockHeight: 4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := new(ValidationKeeperTestSuite)
			suite.SetT(&testing.T{})
			keeper, ctrl := suite.ValidationKeeperSetupTest()
			defer ctrl.Finish()

			ctx := sdk.UnwrapSDKContext(suite.ctx)
			ctx = ctx.WithBlockHeight(tt.args.blockHeight)

			consensusParams := ctx.ConsensusParams()
			consensusParams.Abci = &cmtproto.ABCIParams{
				VoteExtensionsEnableHeight: 2,
			}
			ctx = ctx.WithConsensusParams(consensusParams)

			err := keeper.PreBlocker(ctx, tt.args.req)
			require.NoError(t, err)
		})
	}
}

func TestGetValidatedOracleData(t *testing.T) {

	type args struct {
		voteExt         keeper.VoteExtension
		fieldVotePowers map[keeper.VoteExtensionField]int64
	}
	tests := []struct {
		name    string
		args    args
		want    *keeper.OracleData
		wantErr bool
	}{
		{
			name: "PASS: get validated oracle data",
			args: args{
				voteExt: validationtestutil.SampleDecodedVoteExtension,
				fieldVotePowers: map[keeper.VoteExtensionField]int64{
					keeper.VEFieldZRChainBlockHeight:      1,
					keeper.VEFieldEigenDelegationsHash:    1,
					keeper.VEFieldRequestedBtcBlockHeight: 1,
					keeper.VEFieldRequestedBtcHeaderHash:  1,
					keeper.VEFieldEthBlockHeight:          1,
					keeper.VEFieldEthGasLimit:             1,
					keeper.VEFieldEthBaseFee:              1,
					keeper.VEFieldEthTipCap:               1,
					keeper.VEFieldRequestedStakerNonce:    1,
					keeper.VEFieldRequestedEthMinterNonce: 1,
					keeper.VEFieldRequestedUnstakerNonce:  1,
					keeper.VEFieldRequestedCompleterNonce: 1,
					keeper.VEFieldSolanaMintNoncesHash:    1,
					keeper.VEFieldSolanaAccountsHash:      1,
					keeper.VEFieldEthBurnEventsHash:       1,
					keeper.VEFieldSolanaBurnEventsHash:    1,
					keeper.VEFieldSolanaMintEventsHash:    1,
					keeper.VEFieldRedemptionsHash:         1,
					keeper.VEFieldROCKUSDPrice:            1,
					keeper.VEFieldBTCUSDPrice:             1,
					keeper.VEFieldETHUSDPrice:             1,
					keeper.VEFieldLatestBtcBlockHeight:    1,
					keeper.VEFieldLatestBtcHeaderHash:     1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := new(ValidationKeeperTestSuite)
			suite.SetT(&testing.T{})
			keeper, ctrl := suite.ValidationKeeperSetupTest()
			defer ctrl.Finish()

			ctx := sdk.UnwrapSDKContext(suite.ctx)
			ctx = ctx.WithBlockHeight(3)

			consensusParams := ctx.ConsensusParams()
			consensusParams.Abci = &cmtproto.ABCIParams{
				VoteExtensionsEnableHeight: 2,
			}
			ctx = ctx.WithConsensusParams(consensusParams)

			got, err := keeper.GetValidatedOracleData(ctx, tt.args.voteExt, tt.args.fieldVotePowers)
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}

func TestUpdateValidatorStakes(t *testing.T) {

	type args struct {
		oracleData keeper.OracleData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "PASS: update validator stakes",
			args: args{
				oracleData: validationtestutil.SampleOracleData,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := new(ValidationKeeperTestSuite)
			suite.SetT(&testing.T{})
			keeper, ctrl := suite.ValidationKeeperSetupTest()
			defer ctrl.Finish()

			ctx := sdk.UnwrapSDKContext(suite.ctx)
			keeper.UpdateValidatorStakes(ctx, tt.args.oracleData)
			// No error to check since the function doesn't return anything
		})
	}
}

func TestRemoveStaleValidatorDelegations(t *testing.T) {
	type args struct {
		validatorInAVSDelegationSet map[string]bool
		existingValidators          map[string]sdkmath.Int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "PASS: update validator stakes",
			args: args{
				validatorInAVSDelegationSet: map[string]bool{
					"zenvaloper1tnh2q55v8wyygtt9srz5safamzdengsns4jcd6": true,
					"zenvaloper1ghekyjucln7y67ntx7cf27m9dpuxxemn953g2g": true,
				},
				existingValidators: map[string]sdkmath.Int{
					"zenvaloper1tnh2q55v8wyygtt9srz5safamzdengsns4jcd6": sdkmath.NewInt(1000000),
					"zenvaloper1ghekyjucln7y67ntx7cf27m9dpuxxemn953g2g": sdkmath.NewInt(1000000),
					"zenvaloper1p8wcgrjr4pjju90xg6u9cgq55dxwq8j7ves9zy": sdkmath.NewInt(1000000),
				},
			},
			want: 2,
		},
		{
			name: "PASS: remove validator from validatorInAVSDelegationSet",
			args: args{
				validatorInAVSDelegationSet: map[string]bool{
					"zenvaloper1tnh2q55v8wyygtt9srz5safamzdengsns4jcd6": true,
				},
				existingValidators: map[string]sdkmath.Int{
					"zenvaloper1tnh2q55v8wyygtt9srz5safamzdengsns4jcd6": sdkmath.NewInt(1000000),
					"zenvaloper1ghekyjucln7y67ntx7cf27m9dpuxxemn953g2g": sdkmath.NewInt(1000000),
					"zenvaloper1p8wcgrjr4pjju90xg6u9cgq55dxwq8j7ves9zy": sdkmath.NewInt(1000000),
				},
			},
			want: 1,
		},
		{
			name: "PASS: no validators in validatorInAVSDelegationSet",
			args: args{
				validatorInAVSDelegationSet: map[string]bool{},
				existingValidators: map[string]sdkmath.Int{
					"zenvaloper1tnh2q55v8wyygtt9srz5safamzdengsns4jcd6": sdkmath.NewInt(1000000),
					"zenvaloper1ghekyjucln7y67ntx7cf27m9dpuxxemn953g2g": sdkmath.NewInt(1000000),
					"zenvaloper1p8wcgrjr4pjju90xg6u9cgq55dxwq8j7ves9zy": sdkmath.NewInt(1000000),
				},
			},
			want: 0,
		},
		{
			name: "PASS: three validators in validatorInAVSDelegationSet",
			args: args{
				validatorInAVSDelegationSet: map[string]bool{
					"zenvaloper1tnh2q55v8wyygtt9srz5safamzdengsns4jcd6": true,
					"zenvaloper1ghekyjucln7y67ntx7cf27m9dpuxxemn953g2g": true,
					"zenvaloper1p8wcgrjr4pjju90xg6u9cgq55dxwq8j7ves9zy": true,
				},
				existingValidators: map[string]sdkmath.Int{
					"zenvaloper1tnh2q55v8wyygtt9srz5safamzdengsns4jcd6": sdkmath.NewInt(1000000),
					"zenvaloper1ghekyjucln7y67ntx7cf27m9dpuxxemn953g2g": sdkmath.NewInt(1000000),
					"zenvaloper1p8wcgrjr4pjju90xg6u9cgq55dxwq8j7ves9zy": sdkmath.NewInt(1000000),
				},
			},
			want: 3,
		},
		{
			name: "PASS: remove validators gracefully when validator store is empty",
			args: args{
				validatorInAVSDelegationSet: map[string]bool{
					"zenvaloper1tnh2q55v8wyygtt9srz5safamzdengsns4jcd6": true,
				},
				existingValidators: map[string]sdkmath.Int{
					"zenvaloper1tnh2q55v8wyygtt9srz5safamzdengsns4jcd6": sdkmath.NewInt(1000000),
					"zenvaloper1ghekyjucln7y67ntx7cf27m9dpuxxemn953g2g": sdkmath.NewInt(1000000),
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := new(ValidationKeeperTestSuite)
			suite.SetT(&testing.T{})
			keeper, ctrl := suite.ValidationKeeperSetupTest()
			defer ctrl.Finish()

			ctx := sdk.UnwrapSDKContext(suite.ctx)

			for validator, amount := range tt.args.existingValidators {
				// Create validator in the validator store
				validatorHV := types.ValidatorHV{
					OperatorAddress:   validator,
					ConsensusPubkey:   nil, // Not needed for this test
					Jailed:            false,
					Status:            types.Bonded,
					TokensNative:      amount,
					DelegatorShares:   sdkmath.LegacyNewDecFromInt(amount),
					Description:       types.Description{},
					UnbondingHeight:   0,
					UnbondingTime:     time.Time{},
					Commission:        types.Commission{},
					MinSelfDelegation: sdkmath.ZeroInt(),
					TokensAVS:         sdkmath.ZeroInt(),
				}

				err := keeper.SetValidator(ctx, validatorHV)
				require.NoError(t, err)

				err = keeper.ValidatorDelegations.Set(ctx, validator, amount)
				require.NoError(t, err)
			}

			initialCount := 0
			keeper.ValidatorDelegations.Walk(ctx, nil, func(key string, value sdkmath.Int) (bool, error) {
				initialCount++
				return false, nil
			})

			keeper.RemoveStaleValidatorDelegations(ctx, tt.args.validatorInAVSDelegationSet)

			validatorDelegations := 0
			keeper.ValidatorDelegations.Walk(ctx, nil, func(key string, value sdkmath.Int) (bool, error) {
				validatorDelegations++
				return false, nil
			})

			require.Equal(t, tt.want, validatorDelegations)

		})
	}
}
