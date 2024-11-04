package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v5/x/identity/module"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_Keyrings(t *testing.T) {

	type args struct {
		req          *types.QueryKeyringsRequest
		msgKeyring   *types.MsgNewKeyring
		keyringCount int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "PASS: create 100 keyrings",
			args: args{
				req: &types.QueryKeyringsRequest{
					Pagination: nil,
				},
				msgKeyring:   types.NewMsgNewKeyring("testCreator", "testDescription", 0, 0),
				keyringCount: 100,
			},
			want: 100,
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

			for i := 0; i < tt.args.keyringCount; i++ {
				_, err := msgSer.NewKeyring(ctx, tt.args.msgKeyring)
				require.NoError(t, err)
			}

			got, err := ik.Keyrings(ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, len(got.Keyrings))
		})
	}
}
