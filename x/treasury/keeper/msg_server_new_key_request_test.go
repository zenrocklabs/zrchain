package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v5/x/identity/module"
	idTypes "github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/keeper"
	treasury "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_NewKeyRequest(t *testing.T) {

	type args struct {
		keyring   *idTypes.Keyring
		workspace *idTypes.Workspace
		msg       *types.MsgNewKeyRequest
	}
	tests := []struct {
		name           string
		args           args
		wantKeyRequest *types.KeyRequest
		want           *types.MsgNewKeyRequestResponse
		wantErr        bool
	}{
		{
			name: "PASS: request a new ecdsa key",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				msg:       types.NewMsgNewKeyRequest("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", "keyring1pfnq7r04rept47gaf5cpdew2", "ecdsa", 1000, 0),
			},
			wantKeyRequest: &types.KeyRequest{
				Id:            1,
				Creator:       "testOwner",
				WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
				KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
				KeyType:       types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
			},
			want: &types.MsgNewKeyRequestResponse{KeyReqId: 1},
		},
		{
			name: "PASS: request a new eddsa key",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				msg:       types.NewMsgNewKeyRequest("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", "keyring1pfnq7r04rept47gaf5cpdew2", "ed25519", 1000, 0),
			},
			wantKeyRequest: &types.KeyRequest{
				Id:            1,
				Creator:       "testOwner",
				WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
				KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
				KeyType:       types.KeyType_KEY_TYPE_EDDSA_ED25519,
				Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
			},
			want: &types.MsgNewKeyRequestResponse{KeyReqId: 1},
		},
		{
			name: "FAIL: workspace not found",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				msg:       types.NewMsgNewKeyRequest("testOwner", "notAWorkspace", "keyring1pfnq7r04rept47gaf5cpdew2", "ecdsa", 1000, 0),
			},
			want:    &types.MsgNewKeyRequestResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: policy not found",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				msg:       types.NewMsgNewKeyRequest("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", "keyring1pfnq7r04rept47gaf5cpdew2", "ecdsa", 1000, 1),
			},
			want:    &types.MsgNewKeyRequestResponse{},
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
					Parties:     []string{},
					KeyReqFee:   0,
					SigReqFee:   0,
					IsActive:    false,
				},
				workspace: &defaultWs,
				msg:       types.NewMsgNewKeyRequest("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", "keyring1pfnq7r04rept47gaf5cpdew2", "ecdsa", 1000, 0),
			},
			want:    &types.MsgNewKeyRequestResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: invalid workspace",
			args: args{
				keyring: &idTypes.Keyring{
					Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
					Creator:     "testCreator",
					Description: "testDescription",
					Admins:      []string{"testCreator"},
					Parties:     []string{},
					KeyReqFee:   0,
					SigReqFee:   0,
					IsActive:    false,
				},
				workspace: &defaultWs,
				msg:       types.NewMsgNewKeyRequest("testOwner", "invalid", "keyring1pfnq7r04rept47gaf5cpdew2", "ecdsa", 1000, 0),
			},
			want:    &types.MsgNewKeyRequestResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: invalid keyring",
			args: args{
				keyring: &idTypes.Keyring{
					Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
					Creator:     "testCreator",
					Description: "testDescription",
					Admins:      []string{"testCreator"},
					Parties:     []string{},
					KeyReqFee:   0,
					SigReqFee:   0,
					IsActive:    false,
				},
				workspace: &defaultWs,
				msg:       types.NewMsgNewKeyRequest("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", "invalid", "ecdsa", 1000, 0),
			},
			want:    &types.MsgNewKeyRequestResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: invalid key type",
			args: args{
				keyring: &idTypes.Keyring{
					Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
					Creator:     "testCreator",
					Description: "testDescription",
					Admins:      []string{"testCreator"},
					Parties:     []string{},
					KeyReqFee:   0,
					SigReqFee:   0,
					IsActive:    false,
				},
				workspace: &defaultWs,
				msg:       types.NewMsgNewKeyRequest("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", "keyring1pfnq7r04rept47gaf5cpdew2", "invalid", 1000, 0),
			},
			want:    &types.MsgNewKeyRequestResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			ik := keepers.IdentityKeeper
			tk := keepers.TreasuryKeeper
			pk := keepers.PolicyKeeper
			ctx := keepers.Ctx
			msgSer := keeper.NewMsgServerImpl(*tk)

			idGenesis := idTypes.GenesisState{
				PortId:     idTypes.PortID,
				Keyrings:   []idTypes.Keyring{*tt.args.keyring},
				Workspaces: []idTypes.Workspace{*tt.args.workspace},
			}
			identity.InitGenesis(ctx, *ik, idGenesis)

			trGenesis := types.GenesisState{
				PortId: types.PortID,
			}
			treasury.InitGenesis(ctx, *tk, trGenesis)

			got, err := msgSer.NewKeyRequest(ctx, tt.args.msg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)

				gotKeyReq, err := tk.KeyRequestStore.Get(ctx, got.KeyReqId)
				require.NoError(t, err)
				require.Equal(t, tt.wantKeyRequest, &gotKeyReq)

				act, err := pk.ActionStore.Get(ctx, 1)
				require.NoError(t, err)
				require.Equal(t, uint64(1000), act.Btl)
			}
		})
	}
}
