package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	pol "github.com/Zenrock-Foundation/zrchain/v6/x/policy/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
)

func TestKeeper_PoliciesByCreators(t *testing.T) {
	policy, _ := codectypes.NewAnyWithValue(&types.BoolparserPolicy{
		Definition: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq > 1",
		Participants: []*types.PolicyParticipant{
			{
				Address: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
			},
			{
				Address: "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
			},
		},
	})

	policy1 := types.Policy{
		Creator: "testCreator",
		Id:      1,
		Name:    "testPolicy",
		Policy:  policy,
	}
	policy2 := types.Policy{
		Creator: "testCreator2",
		Id:      2,
		Name:    "testPolicy",
		Policy:  policy,
	}

	type args struct {
		policy []types.Policy
		req    *types.QueryPoliciesByCreatorRequest
	}
	tests := []struct {
		name    string
		args    args
		want    types.QueryPoliciesByCreatorResponse
		wantErr bool
	}{
		{
			name: "PASS: get policy for 1 creator",
			args: args{
				policy: []types.Policy{policy1, policy2},
				req: &types.QueryPoliciesByCreatorRequest{
					Creators: []string{"testCreator"},
				},
			},
			want: types.QueryPoliciesByCreatorResponse{
				Policies:   []*types.Policy{&policy1},
				Pagination: &query.PageResponse{Total: 1},
			},
		},
		{
			name: "PASS: get policy for multiple creators",
			args: args{
				policy: []types.Policy{policy1, policy2},
				req: &types.QueryPoliciesByCreatorRequest{
					Creators: []string{"testCreator", "testCreator2"},
				},
			},
			want: types.QueryPoliciesByCreatorResponse{
				Policies:   []*types.Policy{&policy1, &policy2},
				Pagination: &query.PageResponse{Total: 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			pk := keepers.PolicyKeeper
			ctx := keepers.Ctx

			genesis := types.GenesisState{
				PortId:   types.PortID,
				Policies: tt.args.policy,
			}
			pol.InitGenesis(ctx, *pk, genesis)

			got, err := pk.PoliciesByCreator(ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PoliciesByCreator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				require.Equal(t, len(tt.want.Policies), len(got.Policies), "Those values should be the same")
				for i := range got.Policies {
					require.Equal(t, tt.want.Policies[i].Creator, got.Policies[i].Creator)
				}
			}
		})
	}
}
