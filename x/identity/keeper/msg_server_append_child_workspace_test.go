package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v6/x/identity/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	"github.com/stretchr/testify/require"
)

var childWs = types.Workspace{
	Address: "childWs",
	Creator: "testOwner",
	Owners:  []string{"testOwner"},
}

var invalidChildWs = types.Workspace{
	Address: "invalidChildWs",
	Creator: "testOwner2",
	Owners:  []string{"testOwner2"},
}

var wsWithChild = types.Workspace{
	Address:         "workspace14a2hpadpsy9h4auve2z8lw",
	Creator:         "testOwner",
	Owners:          []string{"testOwner"},
	ChildWorkspaces: []string{"childWs"},
}

func Test_msgServer_AppendChildWorkspace(t *testing.T) {
	type args struct {
		workspace *types.Workspace
		childWs   *types.Workspace
		msg       *types.MsgAppendChildWorkspace
	}
	tests := []struct {
		name          string
		args          args
		want          *types.MsgAppendChildWorkspaceResponse
		wantWorkspace *types.Workspace
		wantErr       bool
	}{
		{
			name: "PASS: add child workspace",
			args: args{
				workspace: &defaultWs,
				childWs:   &childWs,
				msg:       types.NewMsgAppendChildWorkspace("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", childWs.Address, 1000),
			},
			want: &types.MsgAppendChildWorkspaceResponse{},
			wantWorkspace: &types.Workspace{
				Address:         "workspace14a2hpadpsy9h4auve2z8lw",
				Creator:         "testOwner",
				Owners:          []string{"testOwner"},
				ChildWorkspaces: []string{"childWs"},
			},
		},
		{
			name: "FAIL: workspace is nil or not found",
			args: args{
				workspace: &defaultWs,
				childWs:   &childWs,
				msg:       types.NewMsgAppendChildWorkspace("testOwner", "notAWorkspace", childWs.Address, 1000),
			},
			want:    &types.MsgAppendChildWorkspaceResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: creator is not an owner of parent",
			args: args{
				workspace: &defaultWs,
				childWs:   &childWs,
				msg:       types.NewMsgAppendChildWorkspace("notAnOwner", "workspace14a2hpadpsy9h4auve2z8lw", childWs.Address, 1000),
			},
			want:    &types.MsgAppendChildWorkspaceResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: creator is not an owner of child",
			args: args{
				workspace: &defaultWs,
				childWs:   &invalidChildWs,
				msg:       types.NewMsgAppendChildWorkspace("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", invalidChildWs.Address, 1000),
			},
			want:    &types.MsgAppendChildWorkspaceResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: new child is already a child",
			args: args{
				workspace: &wsWithChild,
				childWs:   &childWs,
				msg:       types.NewMsgAppendChildWorkspace("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", "childWs", 1000),
			},
			want:    &types.MsgAppendChildWorkspaceResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: new child is nil",
			args: args{
				workspace: &wsWithChild,
				childWs:   &childWs,
				msg:       types.NewMsgAppendChildWorkspace("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", "noChild", 1000),
			},
			want:    &types.MsgAppendChildWorkspaceResponse{},
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
				Workspaces: []types.Workspace{*tt.args.workspace, *tt.args.childWs},
			}
			identity.InitGenesis(ctx, *ik, genesis)

			got, err := msgSer.AppendChildWorkspace(ctx, tt.args.msg)
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
