package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v4/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/identity/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v4/x/identity/module"
	"github.com/Zenrock-Foundation/zrchain/v4/x/identity/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_NewKeyring(t *testing.T) {
	type args struct {
		msg *types.MsgNewKeyring
	}
	tests := []struct {
		name        string
		args        args
		want        *types.MsgNewKeyringResponse
		wantCreated *types.Keyring
		wantErr     bool
	}{
		{
			name: "PASS: create a keyring",
			args: args{
				msg: types.NewMsgNewKeyring("testCreator", "testDescription", 0, 0),
			},
			want: &types.MsgNewKeyringResponse{Addr: "keyring1k6vc6vhp6e6l3rxalue9v4ux"},
			wantCreated: &types.Keyring{
				Address:     "keyring1k6vc6vhp6e6l3rxalue9v4ux",
				Creator:     "testCreator",
				Description: "testDescription",
				Admins:      []string{"testCreator"},
				Parties:     nil,
				KeyReqFee:   0,
				SigReqFee:   0,
				IsActive:    true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			ik := keepers.IdentityKeeper
			ctx := keepers.Ctx
			msgSer := keeper.NewMsgServerImpl(*ik)

			genesis := types.GenesisState{
				PortId: types.PortID,
			}
			identity.InitGenesis(ctx, *ik, genesis)

			got, err := msgSer.NewKeyring(ctx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKeyring() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)

			gotKeyring, err := ik.KeyringStore.Get(ctx, got.Addr)
			require.NoError(t, err)

			require.Equal(t, tt.wantCreated, &gotKeyring)
		})
	}
}
