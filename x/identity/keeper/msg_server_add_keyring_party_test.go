package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v6/x/identity/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_AddKeyringParty(t *testing.T) {

	type args struct {
		keyring *types.Keyring
		msg     *types.MsgAddKeyringParty
	}
	tests := []struct {
		name        string
		args        args
		want        *types.MsgAddKeyringPartyResponse
		wantKeyring *types.Keyring
		wantErr     bool
	}{
		{
			name: "PASS: add a party to a keyring",
			args: args{
				keyring: &testutil.DefaultKr,
				msg:     types.NewMsgAddKeyringParty("testCreator", "keyring1pfnq7r04rept47gaf5cpdew2", "testParty"),
			},
			want: &types.MsgAddKeyringPartyResponse{},
			wantKeyring: &types.Keyring{
				Address:        testutil.DefaultKr.Address,
				Creator:        testutil.DefaultKr.Creator,
				Description:    testutil.DefaultKr.Description,
				Admins:         testutil.DefaultKr.Admins,
				Parties:        append(testutil.DefaultKr.Parties, "testParty"),
				KeyReqFee:      testutil.DefaultKr.KeyReqFee,
				SigReqFee:      testutil.DefaultKr.SigReqFee,
				IsActive:       testutil.DefaultKr.IsActive,
				PartyThreshold: 1,
			},
		},
		{
			name: "FAIL: keyring not found",
			args: args{
				keyring: &testutil.DefaultKr,
				msg:     types.NewMsgAddKeyringParty("testCreator", "invalidKeyring", "testParty"),
			},
			want:    &types.MsgAddKeyringPartyResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: party is already in the keyring",
			args: args{
				keyring: &types.Keyring{
					Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
					Creator:     "testCreator",
					Description: "testDescription",
					Admins:      []string{"testCreator"},
					Parties:     []string{"testParty"},
					KeyReqFee:   0,
					SigReqFee:   0,
					IsActive:    true,
				},
				msg: types.NewMsgAddKeyringParty("testCreator", "keyring1pfnq7r04rept47gaf5cpdew2", "testParty"),
			},
			want:    &types.MsgAddKeyringPartyResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: creator no keyring admin",
			args: args{
				keyring: &testutil.DefaultKr,
				msg:     types.NewMsgAddKeyringParty("notKeyringAdmin", "keyring1pfnq7r04rept47gaf5cpdew2", "testParty"),
			},
			want:    &types.MsgAddKeyringPartyResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: inactive keyring",
			args: args{
				keyring: &types.Keyring{
					Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
					Creator:     "testCreator",
					Description: "testDescription",
					Admins:      []string{"testCreator"},
					Parties:     []string{},
					KeyReqFee:   0,
					SigReqFee:   0,
					IsActive:    false,
				},
				msg: types.NewMsgAddKeyringParty("testCreator", "keyring1pfnq7r04rept47gaf5cpdew2", "testParty"),
			},
			want:    &types.MsgAddKeyringPartyResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			ik := keepers.IdentityKeeper
			ctx := keepers.Ctx
			msgSer := keeper.NewMsgServerImpl(*ik)

			genesis := types.GenesisState{
				PortId:   types.PortID,
				Keyrings: []types.Keyring{*tt.args.keyring},
			}
			identity.InitGenesis(ctx, *ik, genesis)

			got, err := msgSer.AddKeyringParty(ctx, tt.args.msg)
			require.Equal(t, tt.wantErr, err != nil)

			if !tt.wantErr {
				require.Equal(t, tt.want, got, "AddKeyringParty response does not match expected value")
				gotKeyring, err := ik.KeyringStore.Get(ctx, tt.args.keyring.Address)
				require.NoError(t, err)
				require.Equal(t, tt.wantKeyring, &gotKeyring, "keyring does not match expected value")
			}
		})
	}
}
