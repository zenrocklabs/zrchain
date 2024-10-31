package keeper_test

import (
	"reflect"
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	pol "github.com/Zenrock-Foundation/zrchain/v5/x/policy/module"
	policytypes "github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_NewWorkspace(t *testing.T) {
	type args struct {
		msg      *types.MsgNewWorkspace
		policies []policytypes.Policy
	}
	tests := []struct {
		name        string
		args        args
		want        *types.MsgNewWorkspaceResponse
		wantCreated *types.Workspace
		wantErr     bool
	}{
		{
			name: "PASS: create a workspace",
			args: args{
				msg:      types.NewMsgNewWorkspace("testOwner", 0, 0),
				policies: []policytypes.Policy{},
			},
			want: &types.MsgNewWorkspaceResponse{
				Addr: "workspace10j06zdk5gyl6v9ekzwem0v",
			},
			wantCreated: &types.Workspace{
				Address:         "workspace10j06zdk5gyl6v9ekzwem0v",
				Creator:         "testOwner",
				Owners:          []string{"testOwner"},
				ChildWorkspaces: nil,
			},
		},
		{
			name: "PASS: create a workspace with additional owners",
			args: args{
				msg:      types.NewMsgNewWorkspace("testOwner", 0, 0, "owner1", "owner2"),
				policies: []policytypes.Policy{},
			},
			want: &types.MsgNewWorkspaceResponse{
				Addr: "workspace10j06zdk5gyl6v9ekzwem0v",
			},
			wantCreated: &types.Workspace{
				Address:         "workspace10j06zdk5gyl6v9ekzwem0v",
				Creator:         "testOwner",
				Owners:          []string{"testOwner", "owner1", "owner2"},
				ChildWorkspaces: nil,
			},
		},
		{
			name: "FAIL: add owner twice",
			args: args{
				msg: types.NewMsgNewWorkspace("testOwner", 0, 0, "testOwner", "owner2"),
			},
			wantErr: true,
		},
		{
			name: "FAIL: create workspace withpolicy and missing owner",
			args: args{
				msg: types.NewMsgNewWorkspace("testOwner", 1, 0),
				policies: []policytypes.Policy{
					policy1,
				},
			},
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

			polGenesis := policytypes.GenesisState{
				PortId:   policytypes.PortID,
				Policies: tt.args.policies,
			}

			pol.InitGenesis(ctx, *pk, polGenesis)

			got, err := msgSer.NewWorkspace(ctx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWorkspace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.want, got)

			if !tt.wantErr {
				gotWorkspace, err := ik.WorkspaceStore.Get(ctx, got.Addr)
				require.NoError(t, err)

				if !reflect.DeepEqual(&gotWorkspace, tt.wantCreated) {
					t.Errorf("NewWorkspace() got = %v, want %v", gotWorkspace, tt.wantCreated)
				}
			} else {
				require.NotNil(t, err)
			}
		})
	}
}
