package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_WorkspaceByAddress(t *testing.T) {
	type args struct {
		req          *types.QueryWorkspaceByAddressRequest
		msgWorkspace *types.MsgNewWorkspace
	}
	tests := []struct {
		name          string
		args          args
		want          *types.QueryWorkspaceByAddressResponse
		wantWorkspace *types.Workspace
		wantErr       bool
	}{
		{
			name: "PASS: happy path",
			args: args{
				req: &types.QueryWorkspaceByAddressRequest{
					WorkspaceAddr: "workspace10j06zdk5gyl6v9ekzwem0v",
				},
				msgWorkspace: types.NewMsgNewWorkspace("testOwner", 0, 0),
			},
			want: &types.QueryWorkspaceByAddressResponse{
				Workspace: &types.Workspace{
					Address:         "workspace10j06zdk5gyl6v9ekzwem0v",
					Creator:         "testOwner",
					Owners:          []string{"testOwner"},
					ChildWorkspaces: nil,
				},
			},
		},
		{
			name: "FAIL: req is nil",
			args: args{
				req:          nil,
				msgWorkspace: types.NewMsgNewWorkspace("testOwner", 0, 0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "FAIL: wrong workspace address",
			args: args{
				req: &types.QueryWorkspaceByAddressRequest{
					WorkspaceAddr: "wrongAddress",
				},
				msgWorkspace: types.NewMsgNewWorkspace("testOwner", 0, 0),
			},
			want: &types.QueryWorkspaceByAddressResponse{
				Workspace: &types.Workspace{
					Address:         "workspace10j06zdk5gyl6v9ekzwem0v",
					Creator:         "testOwner",
					Owners:          []string{"testOwner"},
					ChildWorkspaces: nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			ik := keepers.IdentityKeeper
			ctx := keepers.Ctx
			msgSer := keeper.NewMsgServerImpl(*ik)

			_, err := msgSer.NewWorkspace(ctx, tt.args.msgWorkspace)
			require.NoError(t, err)

			got, err := ik.WorkspaceByAddress(ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
