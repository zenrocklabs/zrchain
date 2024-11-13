package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v5/x/identity/module"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	"github.com/stretchr/testify/require"
)

var defaultWs = types.Workspace{
	Address: "workspace14a2hpadpsy9h4auve2z8lw",
	Creator: "testOwner",
	Owners:  []string{"testOwner"},
}

func Test_msgServer_AddWorkspaceOwner(t *testing.T) {
	type args struct {
		workspace *types.Workspace
		msg       *types.MsgAddWorkspaceOwner
	}
	tests := []struct {
		name          string
		args          args
		want          *types.MsgAddWorkspaceOwnerResponse
		wantWorkspace *types.Workspace
		wantErr       bool
	}{
		{
			name: "PASS: add workspace owner",
			args: args{
				workspace: &defaultWs,
				msg:       types.NewMsgAddWorkspaceOwner("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", "testOwner2", 1000),
			},
			want: &types.MsgAddWorkspaceOwnerResponse{},
			wantWorkspace: &types.Workspace{
				Address: "workspace14a2hpadpsy9h4auve2z8lw",
				Creator: "testOwner",
				Owners:  []string{"testOwner", "testOwner2"},
			},
			wantErr: false,
		},
		{
			name: "FAIL: workspace is nil or not found",
			args: args{
				workspace: &defaultWs,
				msg:       types.NewMsgAddWorkspaceOwner("testOwner", "notAWorkspace", "testOwner2", 1000),
			},
			want:    &types.MsgAddWorkspaceOwnerResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: owner is already owner",
			args: args{
				workspace: &defaultWs,
				msg:       types.NewMsgAddWorkspaceOwner("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", "testOwner", 1000),
			},
			want:    &types.MsgAddWorkspaceOwnerResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: creator is no admin owner",
			args: args{
				workspace: &defaultWs,
				msg:       types.NewMsgAddWorkspaceOwner("noOwner", "workspace14a2hpadpsy9h4auve2z8lw", "testOwner", 1000),
			},
			want:    &types.MsgAddWorkspaceOwnerResponse{},
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

			got, err := msgSer.AddWorkspaceOwner(ctx, tt.args.msg)
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
