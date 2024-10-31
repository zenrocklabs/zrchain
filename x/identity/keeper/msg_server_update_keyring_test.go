package keeper_test

import (
	"reflect"
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v5/x/identity/module"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_UpdateKeyring(t *testing.T) {

	var defaultKr = types.Keyring{
		Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
		Creator:     "testCreator",
		Description: "testDescription",
		Admins:      []string{"testCreator"},
		IsActive:    true,
	}

	var wantKr = types.Keyring{
		Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
		Creator:     "testCreator",
		Description: "testDescription",
		Admins:      []string{"testCreator"},
		IsActive:    true,
	}

	type args struct {
		keyring *types.Keyring
		msg     *types.MsgUpdateKeyring
	}
	tests := []struct {
		name        string
		args        args
		want        *types.MsgUpdateKeyringResponse
		wantKeyring *types.Keyring
		wantErr     bool
	}{
		{
			name: "PASS: change keyring description",
			args: args{
				keyring: &defaultKr,
				msg:     types.NewMsgUpdateKeyring("testCreator", "keyring1pfnq7r04rept47gaf5cpdew2", "newDescription", true),
			},
			want: &types.MsgUpdateKeyringResponse{},
			wantKeyring: &types.Keyring{
				Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
				Creator:     "testCreator",
				Description: "newDescription",
				Admins:      []string{"testCreator"},
				IsActive:    true,
			},
		},
		{
			name: "FAIL: keyring not found",
			args: args{
				keyring: &defaultKr,
				msg:     types.NewMsgUpdateKeyring("testCreator", "invalidKeyring", "newDescription", true),
			},
			want:    &types.MsgUpdateKeyringResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: creator no keyring admin",
			args: args{
				keyring: &defaultKr,
				msg:     types.NewMsgUpdateKeyring("noAdmin", "keyring1pfnq7r04rept47gaf5cpdew2", "newDescription", true),
			},
			want:    &types.MsgUpdateKeyringResponse{},
			wantErr: true,
		},
		{
			name: "PASS: change keyring status to false",
			args: args{
				keyring: &types.Keyring{
					Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
					Creator:     "testCreator",
					Description: "testDescription",
					Admins:      []string{"testCreator"},
					IsActive:    true,
				},
				msg: types.NewMsgUpdateKeyring("testCreator", "keyring1pfnq7r04rept47gaf5cpdew2", "testDescription", false),
			},
			want: &types.MsgUpdateKeyringResponse{},
			wantKeyring: &types.Keyring{
				Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
				Creator:     "testCreator",
				Description: "testDescription",
				Admins:      []string{"testCreator"},
				IsActive:    false,
			},
		},
		{
			name: "PASS: change keyring status to true",
			args: args{
				keyring: &types.Keyring{
					Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
					Creator:     "testCreator",
					Description: "testDescription",
					Admins:      []string{"testCreator"},
					IsActive:    false,
				},
				msg: types.NewMsgUpdateKeyring("testCreator", "keyring1pfnq7r04rept47gaf5cpdew2", "testDescription", true),
			},
			want:        &types.MsgUpdateKeyringResponse{},
			wantKeyring: &wantKr,
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

			got, err := msgSer.UpdateKeyring(ctx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UpdateKeyring() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("UpdateKeyring() got = %v, want %v", got, tt.want)
				}

				gotKeyring, err := ik.KeyringStore.Get(ctx, tt.args.keyring.Address)
				require.NoError(t, err)

				if !reflect.DeepEqual(&gotKeyring, tt.wantKeyring) {
					t.Fatalf("UpdateKeyring() got = %v, want %v", gotKeyring, tt.wantKeyring)
				}
			}
		})
	}
}
