package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	idtypes "github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/policy/keeper"
	policy "github.com/Zenrock-Foundation/zrchain/v6/x/policy/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_RevokeAction(t *testing.T) {
	policyData := types.BoolparserPolicy{
		Definition: "creator + testApprover > 1",
		Participants: []*types.PolicyParticipant{
			{Address: "creator"},
			{Address: "testApprover"},
		},
	}

	addToWorkspaceMsg := idtypes.NewMsgAddWorkspaceOwner(
		"testApprover",
		"workspaceaddr",
		"newOwner",
		1000)
	addToWorkspaceMsgAny, _ := cdctypes.NewAnyWithValue(addToWorkspaceMsg)

	policyDataAny, err := cdctypes.NewAnyWithValue(&policyData)
	require.NoError(t, err)

	var defaultPolicy = types.Policy{
		Id:     1,
		Name:   "boolpolicy",
		Policy: policyDataAny,
	}

	var defaultAction = types.Action{
		Id:         1,
		Approvers:  []string{},
		Status:     types.ActionStatus_ACTION_STATUS_PENDING,
		PolicyId:   1,
		Msg:        addToWorkspaceMsgAny,
		Creator:    "creator",
		Btl:        1000,
		PolicyData: nil,
	}

	type args struct {
		action *types.Action
		policy *types.Policy
	}

	tests := []struct {
		name    string
		msg     *types.MsgRevokeAction
		wantErr bool
		args    *args
	}{
		{
			name: "PASS revoke action",
			msg: &types.MsgRevokeAction{
				Creator:  "creator",
				ActionId: 1,
			},
			wantErr: false,
			args: &args{
				action: &defaultAction,
				policy: &defaultPolicy,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			pk := keepers.PolicyKeeper
			msgSer := keeper.NewMsgServerImpl(*pk)

			polGenesis := types.GenesisState{
				Params:   types.Params{},
				PortId:   "42",
				Policies: []types.Policy{*tt.args.policy},
				Actions:  []types.Action{*tt.args.action},
			}
			policy.InitGenesis(keepers.Ctx, *pk, polGenesis)
			msgSer.RevokeAction(keepers.Ctx, tt.msg)
			_, err := pk.ActionStore.Get(keepers.Ctx, 1)
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}
