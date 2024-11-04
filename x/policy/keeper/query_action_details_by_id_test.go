package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	pol "github.com/Zenrock-Foundation/zrchain/v5/x/policy/module"
	"github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_QueryActionDetailsById(t *testing.T) {
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

	newKeyReqMsg, _ := codectypes.NewAnyWithValue(&treasurytypes.MsgNewKeyRequest{
		Creator:       "some-address",
		WorkspaceAddr: "some-workspace",
		KeyringAddr:   "some-keyring",
		KeyType:       "some-key-type",
	})

	action := types.Action{
		Id:        1,
		Approvers: []string{"zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"},
		Status:    types.ActionStatus_ACTION_STATUS_PENDING,
		PolicyId:  1,
		Msg:       newKeyReqMsg,
		Creator:   "some-creator",
		Btl:       123,
	}

	type args struct {
		approvers []string
		actionId  uint64
		policyId  uint64
	}

	tests := []struct {
		name                 string
		args                 args
		wantApprovers        []string
		wantPendingApprovers []string
		wantErr              bool
	}{
		{
			name: "PASS: Get Action details",
			args: args{
				approvers: []string{"zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"},
				actionId:  1,
				policyId:  1,
			},
			wantErr:              false,
			wantApprovers:        []string{"zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"},
			wantPendingApprovers: []string{"zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq"},
		},
		{
			name: "PASS: policy id zero",
			args: args{
				approvers: []string{"some-address"},
				actionId:  1,
				policyId:  0,
			},
			wantErr:              false,
			wantApprovers:        []string{"some-address"},
			wantPendingApprovers: []string{},
		},
		{
			name: "FAIL: invalid action id",
			args: args{
				approvers: []string{"zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"},
				actionId:  2,
				policyId:  1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		keepers := keepertest.NewTest(t)
		pk := keepers.PolicyKeeper
		ctx := keepers.Ctx

		action.Approvers = tt.args.approvers
		action.PolicyId = tt.args.policyId
		genesis := types.GenesisState{
			PortId:   types.PortID,
			Policies: []types.Policy{policy1},
			Actions:  []types.Action{action},
		}

		pol.InitGenesis(ctx, *pk, genesis)

		details, err := pk.ActionDetailsById(ctx, &types.QueryActionDetailsByIdRequest{
			Id: tt.args.actionId,
		})

		if tt.wantErr {
			require.Error(t, err)
			require.Nil(t, details)
			return
		}

		require.NoError(t, err)
		require.Equal(t, tt.wantApprovers, details.Approvers)
		require.Equal(t, tt.wantPendingApprovers, details.PendingApprovers)
	}
}
