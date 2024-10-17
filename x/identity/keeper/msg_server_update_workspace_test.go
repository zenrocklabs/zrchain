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
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var defaultWsWithOwners = types.Workspace{
	Address: "workspace14a2hpadpsy9h4auve2z8lw",
	Creator: "testOwner",
	Owners:  []string{"testOwner", "testOwner2"},
}

var policy, _ = codectypes.NewAnyWithValue(&policytypes.BoolparserPolicy{
	Definition: "u1 + u2 > 1",
	Participants: []*policytypes.PolicyParticipant{
		{
			Abbreviation: "u1",
			Address:      "testOwner",
		},
		{
			Abbreviation: "u2",
			Address:      "testOwner2",
		},
	},
})

var policy1 = policytypes.Policy{
	Id:     1,
	Name:   "Policy1",
	Policy: policy,
}

var policy2 = policytypes.Policy{
	Id:     2,
	Name:   "Policy2",
	Policy: policy,
}

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
				policies:  []policytypes.Policy{policy1, policy2},
				workspace: &defaultWsWithOwners,
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
				workspace: &defaultWsWithOwners,
				msg:       types.NewMsgUpdateWorkspace("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", 1, 2, 1000),
			},
			want:    &types.MsgUpdateWorkspaceResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: sign policy does not exist",
			args: args{
				policies:  []policytypes.Policy{policy1},
				workspace: &defaultWsWithOwners,
				msg:       types.NewMsgUpdateWorkspace("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", 1, 2, 1000),
			},
			want:    &types.MsgUpdateWorkspaceResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: workspace does not exist",
			args: args{
				policies:  []policytypes.Policy{policy1, policy2},
				workspace: &defaultWsWithOwners,
				msg:       types.NewMsgUpdateWorkspace("testOwner", "noWorkspace", 1, 2, 1000),
			},
			want:    &types.MsgUpdateWorkspaceResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: creator is no owner",
			args: args{
				policies:  []policytypes.Policy{policy1, policy2},
				workspace: &defaultWsWithOwners,
				msg:       types.NewMsgUpdateWorkspace("noOwner", "workspace14a2hpadpsy9h4auve2z8lw", 1, 2, 1000),
			},
			want:    &types.MsgUpdateWorkspaceResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: no policy updates ",
			args: args{
				workspace: &defaultWsWithOwners,
				msg:       types.NewMsgUpdateWorkspace("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", 0, 0, 1000),
				policies:  []policytypes.Policy{},
			},
			want:    &types.MsgUpdateWorkspaceResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: policy participant not part of workspace ",
			args: args{
				workspace: &defaultWs,
				msg:       types.NewMsgUpdateWorkspace("testOwner", "workspace14a2hpadpsy9h4auve2z8lw", 1, 0, 1000),
				policies: []policytypes.Policy{
					policy1,
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
			if (err != nil) != tt.wantErr {
				t.Fatalf("UpdateWorkspace() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("UpdateWorkspace() got = %v, want %v", got, tt.want)
				}

				gotWorkspace, err := ik.WorkspaceStore.Get(ctx, tt.args.workspace.Address)
				require.NoError(t, err)

				if !reflect.DeepEqual(&gotWorkspace, tt.wantWorkspace) {
					t.Fatalf("UpdateWorkspace() got = %v, want %v", gotWorkspace, tt.wantWorkspace)
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
