package keeper_test

import (
	"reflect"
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	_ "github.com/Zenrock-Foundation/zrchain/v4/policy"
	keepertest "github.com/Zenrock-Foundation/zrchain/v4/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/keeper"
	policyModule "github.com/Zenrock-Foundation/zrchain/v4/x/policy/module"
	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
)

func Test_msgServer_NewPolicy(t *testing.T) {
	policy, err := codectypes.NewAnyWithValue(&types.BoolparserPolicy{
		Definition: "u1 + u2 > 1",
		Participants: []*types.PolicyParticipant{
			{
				Abbreviation: "u1",
				Address:      "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
			},
			{
				Abbreviation: "u2",
				Address:      "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
			},
		},
	})
	//invalidPolicy, err := codectypes.NewAnyWithValue(&types.BoolparserPolicy{
	//	Definition: "u1 + u2 > 1",
	//})
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name       string
		msg        *types.MsgNewPolicy
		want       *types.MsgNewPolicyResponse
		wantPolicy *types.QueryPolicyByIdResponse
		wantErr    bool
	}{
		{
			name:    "PASS new policy, minimum btl",
			msg:     types.NewMsgNewPolicy("testCreator", "testPolicy", policy, 0),
			want:    &types.MsgNewPolicyResponse{Id: 1},
			wantErr: false,
			wantPolicy: &types.QueryPolicyByIdResponse{
				Policy: &types.PolicyResponse{
					Policy: &types.Policy{
						Creator: "testCreator",
						Id:      1,
						Name:    "testPolicy",
						Policy:  policy,
						Btl:     1000, // default btl
					},
				},
			},
		},
		{
			name:    "PASS new policy, minimum btl",
			msg:     types.NewMsgNewPolicy("testCreator", "testPolicy", policy, 1),
			want:    &types.MsgNewPolicyResponse{Id: 1},
			wantErr: false,
			wantPolicy: &types.QueryPolicyByIdResponse{
				Policy: &types.PolicyResponse{
					Policy: &types.Policy{
						Creator: "testCreator",
						Id:      1,
						Name:    "testPolicy",
						Policy:  policy,
						Btl:     10, // minimum btl
					},
				},
			},
		},
		{
			name:    "PASS new policy, specified btl",
			msg:     types.NewMsgNewPolicy("testCreator", "testPolicy", policy, 100),
			want:    &types.MsgNewPolicyResponse{Id: 1},
			wantErr: false,
			wantPolicy: &types.QueryPolicyByIdResponse{
				Policy: &types.PolicyResponse{
					Policy: &types.Policy{
						Creator: "testCreator",
						Id:      1,
						Name:    "testPolicy",
						Policy:  policy,
						Btl:     100, // set btl
					},
				},
			},
		},

		// TODO: uncomment when BoolparsePolicy Validate() is implemented
		//{
		//	name:    "FAIL new policy - policy nil",
		//	msg:     types.NewMsgNewPolicy("testCreator", "testPolicy", invalidPolicy),
		//	wantErr: true,
		//	wantPolicy: &types.QueryPolicyByIdResponse{
		//		Policy: &types.PolicyResponse{
		//			Policy: &types.Policy{
		//				Id:     1,
		//				Name:   "testPolicy",
		//				Policy: policy,
		//			},
		//		},
		//	},
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			pk := keepers.PolicyKeeper
			msgSer := keeper.NewMsgServerImpl(*pk)

			plGenesis := types.GenesisState{
				Params: types.DefaultParams(),
			}
			policyModule.InitGenesis(keepers.Ctx, *pk, plGenesis)

			got, err := msgSer.NewPolicy(keepers.Ctx, tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPolicy() got = %v, want %v", got, tt.want)
			}

			gotPolicy, err := pk.PolicyById(keepers.Ctx, &types.QueryPolicyByIdRequest{Id: got.Id})
			if err != nil {
				t.Errorf("PolicyById() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(gotPolicy, tt.wantPolicy) {
				t.Errorf("NewPolicy() got = %v, want %v", gotPolicy, tt.wantPolicy)
				return
			}
		})
	}
}
