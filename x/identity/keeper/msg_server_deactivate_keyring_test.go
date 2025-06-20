package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v6/x/identity/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_DeactivateKeyring(t *testing.T) {

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
				keyring: &testutil.Keyring,
				msg:     types.NewMsgDeactivateKeyring("testCreator", "invalidKeyring"),
			},
			want:    &types.MsgDeactivateKeyringResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: creator no keyring admin",
			args: args{
				keyring: &testutil.Keyring,
				msg:     types.NewMsgDeactivateKeyring("noAdmin", "keyring1pfnq7r04rept47gaf5cpdew2"),
			},
			want:    &types.MsgDeactivateKeyringResponse{},
			wantErr: true,
		},
		{
			name: "PASS: change keyring status to false",
			args: args{
				keyring: &testutil.Keyring,
				msg:     types.NewMsgDeactivateKeyring("testCreator", "keyring1pfnq7r04rept47gaf5cpdew2"),
			},
			want:        &types.MsgDeactivateKeyringResponse{},
			wantKeyring: &testutil.WantKeyring,
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
			require.Equal(t, tt.wantErr, err != nil)

			if !tt.wantErr {
				require.Equal(t, tt.want, got, "response does not match expected value")
				gotKeyring, err := ik.KeyringStore.Get(ctx, tt.args.keyring.Address)
				require.NoError(t, err)
				require.Equal(t, tt.wantKeyring, &gotKeyring, "keyring does not match expected value")
			}
		})
	}
}
