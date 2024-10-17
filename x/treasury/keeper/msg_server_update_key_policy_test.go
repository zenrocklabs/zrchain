package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v4/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/treasury/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/treasury/types"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"

	identity "github.com/Zenrock-Foundation/zrchain/v4/x/identity/module"
	idtypes "github.com/Zenrock-Foundation/zrchain/v4/x/identity/types"
	policymodule "github.com/Zenrock-Foundation/zrchain/v4/x/policy/module"
	policytypes "github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
	treasurymodule "github.com/Zenrock-Foundation/zrchain/v4/x/treasury/module"
)

func Test_msgServer_UpdateKeyPolicy(t *testing.T) {

	type args struct {
		wsOwners []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "PASS: update key signing policy",
			args: args{
				wsOwners: []string{"testOwner", "testOwner2"},
			},
			wantErr: false,
		},
		{
			name: "FAIL: policy members must be member or workspace",
			args: args{
				wsOwners: []string{"testOwner"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			keepers := keepertest.NewTest(t)
			tk := keepers.TreasuryKeeper
			pk := keepers.PolicyKeeper
			ik := keepers.IdentityKeeper
			ctx := keepers.Ctx
			msgSer := keeper.NewMsgServerImpl(*tk)

			// create workspace
			idGenesis := idtypes.GenesisState{
				Workspaces: []idtypes.Workspace{
					{
						Address: "workspace14a2hpadpsy9h4auve2z8lw",
						Creator: "testOwner",
						Owners:  tt.args.wsOwners,
					},
				},
			}
			identity.InitGenesis(ctx, *ik, idGenesis)

			// create policies
			policyGenesis := policytypes.GenesisState{
				Policies: []policytypes.Policy{
					policy1,
					policy2,
				},
			}
			policymodule.InitGenesis(ctx, *pk, policyGenesis)

			// create key
			treasuryGenesis := types.GenesisState{
				Keys: []types.Key{
					{
						Id:            1,
						WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
						KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
						Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
						PublicKey:     []byte{3, 49, 216, 1, 30, 80, 59, 231, 205, 4, 237, 36, 101, 79, 9, 119, 30, 114, 19, 70, 242, 44, 91, 95, 238, 20, 245, 197, 193, 86, 22, 120, 35},
						SignPolicyId:  1,
					},
				},
			}
			treasurymodule.InitGenesis(ctx, *tk, treasuryGenesis)

			// execute update action
			_, err := msgSer.UpdateKeyPolicy(ctx, &types.MsgUpdateKeyPolicy{
				Creator:      "testOwner",
				KeyId:        1,
				SignPolicyId: 2,
			})

			if !tt.wantErr {
				require.Nil(t, err)
				// action needs to be approved using the existing policy id
				act, err := pk.ActionStore.Get(ctx, 1)
				require.Nil(t, err)
				assert.Equal(t, uint64(1), act.PolicyId)

				// after update the key needs to have the new policy id
				key, err := tk.KeyStore.Get(ctx, 1)
				require.Nil(t, err)
				assert.Equal(t, uint64(2), key.SignPolicyId)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}
