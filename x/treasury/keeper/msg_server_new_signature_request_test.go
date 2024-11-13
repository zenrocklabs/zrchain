package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v5/x/identity/module"
	idTypes "github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	policymodule "github.com/Zenrock-Foundation/zrchain/v5/x/policy/module"
	policytypes "github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/keeper"
	treasurymodule "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_NewSignatureRequest(t *testing.T) {

	type args struct {
		keyring   *idTypes.Keyring
		workspace *idTypes.Workspace
		key       *types.Key
		msg       *types.MsgNewSignatureRequest
	}
	tests := []struct {
		name            string
		args            args
		wantSignRequest *types.SignRequest
		want            *types.MsgNewSignatureRequestResponse
		wantErr         bool
	}{
		{
			name: "PASS: valid signature request",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key: &types.Key{
					Id:            1,
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
					PublicKey:     []byte{3, 49, 216, 1, 30, 80, 59, 231, 205, 4, 237, 36, 101, 79, 9, 119, 30, 114, 19, 70, 242, 44, 91, 95, 238, 20, 245, 197, 193, 86, 22, 120, 35},
					SignPolicyId:  1,
				},
				msg: types.NewMsgNewSignatureRequest("testOwner", 1, "778f572f33afab831365d52e563a0ddd9969a535f62262125bb06d43681c59aa", 1000),
			},
			wantSignRequest: &types.SignRequest{
				Id:             1,
				Creator:        "testOwner",
				KeyId:          1,
				DataForSigning: [][]byte{{119, 143, 87, 47, 51, 175, 171, 131, 19, 101, 213, 46, 86, 58, 13, 221, 153, 105, 165, 53, 246, 34, 98, 18, 91, 176, 109, 67, 104, 28, 89, 170}},
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
				KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
			},
			want: &types.MsgNewSignatureRequestResponse{SigReqId: 1},
		},
		{
			name: "FAIL: key not found",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key:       &defaultKey,
				msg:       types.NewMsgNewSignatureRequest("testOwner", 5, "778f572f33afab831365d52e563a0ddd9969a535f62262125bb06d43681c59aa", 1000),
			},
			wantErr: true,
		},
		{
			name: "FAIL: workspace not found",
			args: args{
				keyring: &defaultKr,
				workspace: &idTypes.Workspace{
					Address: "otherWorkspace",
					Creator: "testOwner",
					Owners:  []string{"testOwner"},
				},
				key: &defaultKey,
				msg: types.NewMsgNewSignatureRequest("testOwner", 1, "778f572f33afab831365d52e563a0ddd9969a535f62262125bb06d43681c59aa", 1000),
			},
			wantErr: true,
		},
		{
			name: "FAIL: keyring is inactive",
			args: args{
				keyring: &idTypes.Keyring{
					Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
					Creator:     "testCreator",
					Description: "testDescription",
					Admins:      []string{"testCreator"},
					Parties:     []string{"testCreator"},
					KeyReqFee:   0,
					SigReqFee:   0,
					IsActive:    false,
				},
				workspace: &defaultWs,
				key:       &defaultKey,
				msg:       types.NewMsgNewSignatureRequest("testOwner", 1, "778f572f33afab831365d52e563a0ddd9969a535f62262125bb06d43681c59aa", 1000),
			},
			wantErr: true,
		},
		{
			name: "FAIL: invalid length for data for signing",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key:       &defaultKey,
				msg:       types.NewMsgNewSignatureRequest("testOwner", 1, "778f572f", 1000),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			ik := keepers.IdentityKeeper
			pk := keepers.PolicyKeeper
			tk := keepers.TreasuryKeeper
			ctx := keepers.Ctx

			idGenesis := idTypes.GenesisState{
				PortId:     idTypes.PortID,
				Keyrings:   []idTypes.Keyring{*tt.args.keyring},
				Workspaces: []idTypes.Workspace{*tt.args.workspace},
			}
			identity.InitGenesis(ctx, *ik, idGenesis)

			tGenesis := types.GenesisState{
				PortId: types.PortID,
				Keys:   []types.Key{*tt.args.key},
			}
			treasurymodule.InitGenesis(ctx, *tk, tGenesis)

			pGenesis := policytypes.GenesisState{
				Policies: []policytypes.Policy{
					policy1,
				},
			}
			policymodule.InitGenesis(ctx, *pk, pGenesis)

			msgSer := keeper.NewMsgServerImpl(*tk)
			got, err := msgSer.NewSignatureRequest(ctx, tt.args.msg)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)

				gotSigReq, err := tk.SignRequestStore.Get(ctx, got.SigReqId)
				require.NoError(t, err)
				require.Equal(t, tt.wantSignRequest, &gotSigReq)

				act, err := pk.ActionStore.Get(ctx, 1)
				require.NoError(t, err)
				require.Equal(t, uint64(1000), act.Btl)
				require.Equal(t, tt.args.key.SignPolicyId, act.PolicyId)
			}
		})
	}
}
