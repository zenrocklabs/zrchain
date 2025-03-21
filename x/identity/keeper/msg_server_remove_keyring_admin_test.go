package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v6/x/identity/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_RemoveKeyringAdmin(t *testing.T) {

	defaultKrWithAdmins := types.Keyring{
		Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
		Creator:     "testCreator",
		Description: "testDescription",
		Admins:      []string{"testCreator", "admin1", "admin2", "admin3"},
		KeyReqFee:   0,
		SigReqFee:   0,
		IsActive:    true,
	}

	type args struct {
		keyring *types.Keyring
		msg     *types.MsgRemoveKeyringAdmin
	}

	tests := []struct {
		name        string
		args        args
		want        *types.MsgRemoveKeyringAdminResponse
		wantKeyring *types.Keyring
		wantErr     bool
	}{
		{
			name: "PASS: remove keyring admin",
			args: args{
				keyring: &defaultKrWithAdmins,
				msg:     types.NewMsgRemoveKeyringAdmin("testCreator", "keyring1pfnq7r04rept47gaf5cpdew2", "admin3"),
			},
			want: &types.MsgRemoveKeyringAdminResponse{},
			wantKeyring: &types.Keyring{
				Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
				Creator:     "testCreator",
				Description: "testDescription",
				Admins:      []string{"testCreator", "admin1", "admin2"},
				KeyReqFee:   0,
				SigReqFee:   0,
				IsActive:    true,
			},
		},
		{
			name: "FAIL: remove single admin",
			args: args{
				keyring: &types.Keyring{
					Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
					Creator:     "testCreator",
					Description: "testDescription",
					Admins:      []string{"testCreator"},
					KeyReqFee:   0,
					SigReqFee:   0,
					IsActive:    true,
				},
				msg: types.NewMsgRemoveKeyringAdmin("testCreator", "keyring1pfnq7r04rept47gaf5cpdew2", "testCreator"),
			},
			want:    &types.MsgRemoveKeyringAdminResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: keyring is nil or not found",
			args: args{
				keyring: &defaultKrWithAdmins,
				msg:     types.NewMsgRemoveKeyringAdmin("testCreator", "notAKeyring", "admin1"),
			},
			want:    &types.MsgRemoveKeyringAdminResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: removed admin is no admin",
			args: args{
				keyring: &defaultKrWithAdmins,
				msg:     types.NewMsgRemoveKeyringAdmin("testCreator", "keyring1pfnq7r04rept47gaf5cpdew2", "noAdmin"),
			},
			want:    &types.MsgRemoveKeyringAdminResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: creator is no admin owner",
			args: args{
				keyring: &defaultKrWithAdmins,
				msg:     types.NewMsgRemoveKeyringAdmin("noAdmin", "keyring1pfnq7r04rept47gaf5cpdew2", "admin1"),
			},
			want:    &types.MsgRemoveKeyringAdminResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: keyring is not active",
			args: args{
				keyring: &types.Keyring{
					Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
					Creator:     "testCreator",
					Description: "testDescription",
					Admins:      []string{"testCreator", "admin1", "admin2", "admin3"},
					KeyReqFee:   0,
					SigReqFee:   0,
					IsActive:    false,
				},
				msg: types.NewMsgRemoveKeyringAdmin("noAdmin", "keyring1pfnq7r04rept47gaf5cpdew2", "admin1"),
			},
			want:    &types.MsgRemoveKeyringAdminResponse{},
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

			got, err := msgSer.RemoveKeyringAdmin(ctx, tt.args.msg)
			require.Equal(t, tt.wantErr, err != nil)

			if !tt.wantErr {
				require.Equal(t, tt.want, got)

				gotKeyring, err := ik.KeyringStore.Get(ctx, tt.args.keyring.Address)
				require.NoError(t, err)
				require.Equal(t, tt.wantKeyring, &gotKeyring)
			}
		})
	}

}
