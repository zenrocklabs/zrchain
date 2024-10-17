package keeper_test

import (
	"reflect"
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v4/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/identity/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v4/x/identity/module"
	"github.com/Zenrock-Foundation/zrchain/v4/x/identity/types"
	pol "github.com/Zenrock-Foundation/zrchain/v4/x/policy/module"
	policytypes "github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_RemoveWorkspaceOwner(t *testing.T) {

	var defaultWsWithOwners = types.Workspace{
		Address:       "workspace14a2hpadpsy9h4auve2z8lw",
		Creator:       "testOwner",
		Owners:        []string{"testOwner", "testOwner2", "testOwner3"},
		AdminPolicyId: 0,
		SignPolicyId:  0,
	}

	var defaultWsWithOwnersAndPolicy = types.Workspace{
		Address:       "workspace14a2hpadpsy9h4auve2z8lw",
		Creator:       "testOwner",
		Owners:        []string{"testOwner", "testOwner2"},
		AdminPolicyId: 1,
		SignPolicyId:  0,
	}

	type args struct {
		workspace *types.Workspace
		msg       *types.MsgRemoveWorkspaceOwner
		policies  []policytypes.Policy
	}
	tests := []struct {
		name          string
		args          args
		want          *types.MsgRemoveWorkspaceOwnerResponse
		wantWorkspace *types.Workspace
		wantErr       bool
	}{
		{
			name: "PASS: remove workspace owner",
			args: args{
				workspace: &defaultWsWithOwners,
				msg:       types.NewMsgRemoveWorkspaceOwner("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", "testOwner2", 1000),
				policies:  []policytypes.Policy{},
			},
			want: &types.MsgRemoveWorkspaceOwnerResponse{},
			wantWorkspace: &types.Workspace{
				Address:       "workspace14a2hpadpsy9h4auve2z8lw",
				Creator:       "testOwner",
				Owners:        []string{"testOwner", "testOwner3"},
				AdminPolicyId: 0,
				SignPolicyId:  0,
			},
		},
		{
			name: "PASS: remove workspace creator",
			args: args{
				workspace: &defaultWsWithOwners,
				msg:       types.NewMsgRemoveWorkspaceOwner("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", "testOwner", 1000),
				policies:  []policytypes.Policy{},
			},
			want: &types.MsgRemoveWorkspaceOwnerResponse{},
			wantWorkspace: &types.Workspace{
				Address:       "workspace14a2hpadpsy9h4auve2z8lw",
				Creator:       "testOwner",
				Owners:        []string{"testOwner2", "testOwner3"},
				AdminPolicyId: 0,
				SignPolicyId:  0,
			},
		},
		{
			name: "FAIL: remove single owner",
			args: args{
				workspace: &types.Workspace{
					Address:       "workspace14a2hpadpsy9h4auve2z8lw",
					Creator:       "testOwner",
					Owners:        []string{"testOwner"},
					AdminPolicyId: 0,
					SignPolicyId:  0,
				},
				msg:      types.NewMsgRemoveWorkspaceOwner("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", "testOwner", 1000),
				policies: []policytypes.Policy{},
			},
			want:    &types.MsgRemoveWorkspaceOwnerResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: workspace is nil or not found",
			args: args{
				workspace: &defaultWsWithOwners,
				msg:       types.NewMsgRemoveWorkspaceOwner("testOwner", "notAWorkspace", "testOwner2", 1000),
				policies:  []policytypes.Policy{},
			},
			want:    &types.MsgRemoveWorkspaceOwnerResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: removed owner is no owner",
			args: args{
				workspace: &defaultWsWithOwners,
				msg:       types.NewMsgRemoveWorkspaceOwner("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", "noOwner", 1000),
				policies:  []policytypes.Policy{},
			},
			want:    &types.MsgRemoveWorkspaceOwnerResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: creator is no admin owner",
			args: args{
				workspace: &defaultWsWithOwners,
				msg:       types.NewMsgRemoveWorkspaceOwner("noOwner", "workspace14a2hpadpsy9h4auve2z8lw", "testOwner", 1000),
				policies:  []policytypes.Policy{},
			},
			want:    &types.MsgRemoveWorkspaceOwnerResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: removed owner is policy member",
			args: args{
				workspace: &defaultWsWithOwnersAndPolicy,
				msg:       types.NewMsgRemoveWorkspaceOwner("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", "testOwner2", 1000),
				policies: []policytypes.Policy{
					policy1,
				},
			},
			want:    &types.MsgRemoveWorkspaceOwnerResponse{},
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

			got, err := msgSer.RemoveWorkspaceOwner(ctx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RemoveWorkspaceOwner() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("RemoveWorkspaceOwner() got = %v, want %v", got, tt.want)
				}

				gotWorkspace, err := ik.WorkspaceStore.Get(ctx, tt.args.workspace.Address)
				require.NoError(t, err)

				if !reflect.DeepEqual(&gotWorkspace, tt.wantWorkspace) {
					t.Errorf("RemoveWorkspaceOwner() got = %v, want %v", gotWorkspace, tt.wantWorkspace)
				}

				act, err := pk.ActionStore.Get(ctx, 1)
				require.Nil(t, err)
				assert.Equal(t, uint64(1000), act.Btl)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}
