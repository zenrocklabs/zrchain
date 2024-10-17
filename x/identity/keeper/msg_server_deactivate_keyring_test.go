package keeper_test

import (
	"reflect"
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v4/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/identity/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v4/x/identity/module"
	"github.com/Zenrock-Foundation/zrchain/v4/x/identity/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_DeactivateKeyring(t *testing.T) {
	var keyring = types.Keyring{
		Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
		Creator:     "testCreator",
		Description: "testDescription",
		Admins:      []string{"testCreator"},
		IsActive:    true,
	}

	var wantKeyring = types.Keyring{
		Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
		Creator:     "testCreator",
		Description: "testDescription",
		Admins:      []string{"testCreator"},
		IsActive:    false,
	}

	type args struct {
		keyring *types.Keyring
		msg     *types.MsgDeactivateKeyring
	}

	tests := []struct {
		name        string
		args        args
		want        *types.MsgDeactivateKeyringResponse
		wantKeyring *types.Keyring
		wantErr     bool
	}{
		{
			name: "FAIL: keyring not found",
			args: args{
				keyring: &keyring,
				msg:     types.NewMsgDeactivateKeyring("testCreator", "invalidKeyring"),
			},
			want:    &types.MsgDeactivateKeyringResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: creator no keyring admin",
			args: args{
				keyring: &keyring,
				msg:     types.NewMsgDeactivateKeyring("noAdmin", "keyring1pfnq7r04rept47gaf5cpdew2"),
			},
			want:    &types.MsgDeactivateKeyringResponse{},
			wantErr: true,
		},
		{
			name: "PASS: change keyring status to false",
			args: args{
				keyring: &keyring,
				msg:     types.NewMsgDeactivateKeyring("testCreator", "keyring1pfnq7r04rept47gaf5cpdew2"),
			},
			want:        &types.MsgDeactivateKeyringResponse{},
			wantKeyring: &wantKeyring,
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

			got, err := msgSer.DeactivateKeyring(ctx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("DeactivateKeyring() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("DeactivateKeyring() got = %v, want %v", got, tt.want)
				}

				gotKeyring, err := ik.KeyringStore.Get(ctx, tt.args.keyring.Address)
				require.NoError(t, err)

				if !reflect.DeepEqual(&gotKeyring, tt.wantKeyring) {
					t.Fatalf("DeactivateKeyring() got = %v, want %v", gotKeyring, tt.wantKeyring)
				}
			}
		})
	}
}
