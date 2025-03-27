package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v6/x/identity/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_AddKeyringAdmin(t *testing.T) {
	var defaultKr = types.Keyring{
		Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
		Creator:     "testCreator",
		Description: "testDescription",
		Admins:      []string{"testCreator"},
		KeyReqFee:   0,
		SigReqFee:   0,
		IsActive:    true,
	}

	var wantKr = types.Keyring{
		Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
		Creator:     "testCreator",
		Description: "testDescription",
		Admins:      []string{"testCreator", "testAdmin"},
		KeyReqFee:   0,
		SigReqFee:   0,
		IsActive:    true,
	}

	type args struct {
		keyring *types.Keyring
		msg     *types.MsgAddKeyringAdmin
	}
	tests := []struct {
		name        string
		args        args
		want        *types.MsgAddKeyringAdminResponse
		wantKeyring *types.Keyring
		wantErr     bool
	}{
		{
			name: "PASS: add a admin to a keyring",
			args: args{
				keyring: &defaultKr,
				msg:     types.NewMsgAddKeyringAdmin("testCreator", "keyring1pfnq7r04rept47gaf5cpdew2", "testAdmin"),
			},
			want:        &types.MsgAddKeyringAdminResponse{},
			wantKeyring: &wantKr,
		},
		{
			name: "FAIL: keyring not found",
			args: args{
				keyring: &defaultKr,
				msg:     types.NewMsgAddKeyringAdmin("testCreator", "invalidKeyring", "testAdmin"),
			},
			want:    &types.MsgAddKeyringAdminResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: admin is already in the keyring",
			args: args{
				keyring: &types.Keyring{
					Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
					Creator:     "testCreator",
					Description: "testDescription",
					Admins:      []string{"testCreator", "testAdmin"},
					KeyReqFee:   0,
					SigReqFee:   0,
					IsActive:    true,
				},
				msg: types.NewMsgAddKeyringAdmin("testCreator", "keyring1pfnq7r04rept47gaf5cpdew2", "testAdmin"),
			},
			want:    &types.MsgAddKeyringAdminResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: creator no keyring admin",
			args: args{
				keyring: &defaultKr,
				msg:     types.NewMsgAddKeyringAdmin("notKeyringAdmin", "keyring1pfnq7r04rept47gaf5cpdew2", "testAdmin"),
			},
			want:    &types.MsgAddKeyringAdminResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: inactive keyring",
			args: args{
				keyring: &types.Keyring{
					Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
					Creator:     "testCreator",
					Description: "testDescription",
					Admins:      []string{"testCreator"},
					Parties:     []string{},
					KeyReqFee:   0,
					SigReqFee:   0,
					IsActive:    false,
				},
				msg: types.NewMsgAddKeyringAdmin("testCreator", "keyring1pfnq7r04rept47gaf5cpdew2", "testAdmin"),
			},
			want:    &types.MsgAddKeyringAdminResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			ik := keepers.IdentityKeeper
			ctx := keepers.Ctx
			msgSer := keeper.NewMsgServerImpl(*ik)

			genesis := types.GenesisState{
				PortId:   types.PortID,
				Keyrings: []types.Keyring{*tt.args.keyring},
			}
			identity.InitGenesis(ctx, *ik, genesis)

			got, err := msgSer.AddKeyringAdmin(ctx, tt.args.msg)
			require.Equal(t, tt.wantErr, err != nil)

			if !tt.wantErr {
				require.Equal(t, tt.want, got, "AddKeyringAdmin response does not match expected value")
				gotKeyring, err := ik.KeyringStore.Get(ctx, tt.args.keyring.Address)
				require.NoError(t, err)

				require.Equal(t, tt.wantKeyring, &gotKeyring, "keyring does not match expected value")
			}
		})
	}
}
