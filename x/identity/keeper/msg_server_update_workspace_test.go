package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v6/x/identity/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	pol "github.com/Zenrock-Foundation/zrchain/v6/x/policy/module"
	policytypes "github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_UpdateWorkspace(t *testing.T) {

	type args struct {
		policies  []policytypes.Policy
		workspace *types.Workspace
		msg       *types.MsgUpdateWorkspace
	}
	tests := []struct {
		name          string
		args          args
		want          *types.MsgUpdateWorkspaceResponse
		wantWorkspace *types.Workspace
		wantErr       bool
	}{
		{
			name: "PASS: change sign and admin policy",
			args: args{
				policies:  []policytypes.Policy{testutil.Policy1, testutil.Policy2},
				workspace: &testutil.DefaultWsWithOwners,
				msg:       types.NewMsgUpdateWorkspace("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", 1, 2, 1000),
			},
			want: &types.MsgUpdateWorkspaceResponse{},
			wantWorkspace: &types.Workspace{
				Address:       "workspace14a2hpadpsy9h4auve2z8lw",
				Creator:       "testOwner",
				Owners:        []string{"testOwner", "testOwner2"},
				AdminPolicyId: 1,
				SignPolicyId:  2,
			},
		},
		{
			name: "FAIL: admin policy does not exist",
			args: args{
				policies:  []policytypes.Policy{},
				workspace: &testutil.DefaultWsWithOwners,
				msg:       types.NewMsgUpdateWorkspace("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", 1, 2, 1000),
			},
			want:    &types.MsgUpdateWorkspaceResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: sign policy does not exist",
			args: args{
				policies:  []policytypes.Policy{testutil.Policy1},
				workspace: &testutil.DefaultWsWithOwners,
				msg:       types.NewMsgUpdateWorkspace("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", 1, 2, 1000),
			},
			want:    &types.MsgUpdateWorkspaceResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: workspace does not exist",
			args: args{
				policies:  []policytypes.Policy{testutil.Policy1, testutil.Policy2},
				workspace: &testutil.DefaultWsWithOwners,
				msg:       types.NewMsgUpdateWorkspace("testOwner", "noWorkspace", 1, 2, 1000),
			},
			want:    &types.MsgUpdateWorkspaceResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: creator is no owner",
			args: args{
				policies:  []policytypes.Policy{testutil.Policy1, testutil.Policy2},
				workspace: &testutil.DefaultWsWithOwners,
				msg:       types.NewMsgUpdateWorkspace("noOwner", "workspace14a2hpadpsy9h4auve2z8lw", 1, 2, 1000),
			},
			want:    &types.MsgUpdateWorkspaceResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: no policy updates ",
			args: args{
				workspace: &testutil.DefaultWsWithOwners,
				msg:       types.NewMsgUpdateWorkspace("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", 0, 0, 1000),
				policies:  []policytypes.Policy{},
			},
			want:    &types.MsgUpdateWorkspaceResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: policy participant not part of workspace ",
			args: args{
				workspace: &testutil.DefaultWs,
				msg:       types.NewMsgUpdateWorkspace("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", 1, 0, 1000),
				policies: []policytypes.Policy{
					testutil.Policy1,
				},
			},
			want:    &types.MsgUpdateWorkspaceResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			ik := keepers.IdentityKeeper
			pk := keepers.PolicyKeeper
			ctx := keepers.Ctx
			msgSer := keeper.NewMsgServerImpl(*ik)

			genesis := types.GenesisState{
				PortId:     types.PortID,
				Workspaces: []types.Workspace{*tt.args.workspace},
			}
			identity.InitGenesis(ctx, *ik, genesis)

			polGenesis := policytypes.GenesisState{
				PortId:   policytypes.PortID,
				Policies: tt.args.policies,
			}

			pol.InitGenesis(ctx, *pk, polGenesis)

			got, err := msgSer.UpdateWorkspace(ctx, tt.args.msg)
			require.Equal(t, tt.wantErr, err != nil)

			if !tt.wantErr {
				require.Equal(t, tt.want, got)

				gotWorkspace, err := ik.WorkspaceStore.Get(ctx, tt.args.workspace.Address)
				require.NoError(t, err)
				require.Equal(t, tt.wantWorkspace, &gotWorkspace)

				act, err := pk.ActionStore.Get(ctx, 1)
				require.NoError(t, err)
				require.Equal(t, uint64(1000), act.Btl)
			}
		})
	}
}
