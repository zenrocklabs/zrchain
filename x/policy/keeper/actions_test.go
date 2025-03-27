package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	policy "github.com/Zenrock-Foundation/zrchain/v6/x/policy/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_AddAction_Btl(t *testing.T) {
	tests := []struct {
		name    string
		msgBtl  uint64
		polBtl  uint64
		wantBtl uint64
	}{
		{
			name:    "PASS: msg btl has precedence",
			msgBtl:  123,
			polBtl:  1234,
			wantBtl: 123,
		},
		{
			name:    "PASS: if no msg btl then policy btl is set",
			msgBtl:  0,
			polBtl:  1234,
			wantBtl: 1234,
		},
		{
			name:    "PASS: if no msg or pol btl then default btl is set",
			msgBtl:  0,
			polBtl:  0,
			wantBtl: 1000,
		},
		{
			name:    "PASS: if msg btl too small then minimum btl is set",
			msgBtl:  1,
			polBtl:  1234,
			wantBtl: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			pk := keepers.PolicyKeeper
			ctx := keepers.Ctx

			pol, _ := cdctypes.NewAnyWithValue(&types.BoolparserPolicy{
				Participants: []*types.PolicyParticipant{
					{
						Address: "some-creator",
					},
				},
			})

			polGenesis := types.GenesisState{
				Params: types.Params{
					MinimumBtl: 10,
					DefaultBtl: 1000,
				},
				Policies: []types.Policy{
					{
						Id:     1,
						Btl:    tt.polBtl,
						Policy: pol,
					},
				},
			}
			policy.InitGenesis(ctx, *pk, polGenesis)

			res, err := pk.AddAction(ctx, "some-creator", &types.MsgNewPolicy{}, 1, tt.msgBtl, nil, []string{"some-creator"})

			require.Nil(t, err)
			require.NotNil(t, res)

			require.Equal(t, tt.wantBtl, res.Btl)
		})
	}
}
