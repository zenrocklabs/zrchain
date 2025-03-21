package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v6/x/identity/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_NewChildWorkspace(t *testing.T) {
	type args struct {
		workspace *types.Workspace
		msg       *types.MsgNewChildWorkspace
	}
	tests := []struct {
		name          string
		args          args
		want          *types.MsgNewChildWorkspaceResponse
		wantWorkspace *types.Workspace
		wantErr       bool
	}{
		{
			name: "PASS: create new child workspace",
			args: args{
				workspace: &defaultWs,
				msg:       types.NewMsgNewChildWorkspace("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", 1000),
			},
			want: &types.MsgNewChildWorkspaceResponse{
				Address: "workspace1mphgzyhncnzyggfxmv4nmh",
			},
			wantWorkspace: &types.Workspace{
				Address:         "workspace14a2hpadpsy9h4auve2z8lw",
				Creator:         "testOwner",
				Owners:          []string{"testOwner"},
				ChildWorkspaces: []string{"workspace1mphgzyhncnzyggfxmv4nmh"},
			},
		},
		{
			name: "FAIL: workspace is nil or not found",
			args: args{
				workspace: &defaultWs,
				msg:       types.NewMsgNewChildWorkspace("testOwner", "notAWorkspace", 1000),
			},
			want:    &types.MsgNewChildWorkspaceResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: creator is not an owner",
			args: args{
				workspace: &defaultWs,
				msg:       types.NewMsgNewChildWorkspace("notAnOwner", "workspace14a2hpadpsy9h4auve2z8lw", 1000),
			},
			want:    &types.MsgNewChildWorkspaceResponse{},
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

			got, err := msgSer.NewChildWorkspace(ctx, tt.args.msg)
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
