package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v4/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/identity/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v4/x/identity/module"
	"github.com/Zenrock-Foundation/zrchain/v4/x/identity/types"
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
				if _, err := msgSer.NewKeyring(ctx, tt.args.msgKeyring); err != nil {
					t.Fatal(err)
				}
			}
			got, err := ik.Keyrings(ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keyrings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got.Keyrings) != tt.want {
				t.Errorf("Keyrings() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}
