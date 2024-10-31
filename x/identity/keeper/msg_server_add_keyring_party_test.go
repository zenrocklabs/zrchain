package keeper_test

import (
	"reflect"
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v5/x/identity/module"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_AddKeyringParty(t *testing.T) {

	var defaultKr = types.Keyring{
		Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
		Creator:     "testCreator",
		Description: "testDescription",
		Admins:      []string{"testCreator"},
		Parties:     []string{},
		KeyReqFee:   0,
		SigReqFee:   0,
		IsActive:    true,
	}

	var wantKr = types.Keyring{
		Address:        "keyring1pfnq7r04rept47gaf5cpdew2",
		Creator:        "testCreator",
		Description:    "testDescription",
		Admins:         []string{"testCreator"},
		Parties:        []string{"testParty"},
		KeyReqFee:      0,
		SigReqFee:      0,
		IsActive:       true,
		PartyThreshold: 1,
	}

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
				keyring: &defaultKr,
				msg:     types.NewMsgAddKeyringParty("testCreator", "keyring1pfnq7r04rept47gaf5cpdew2", "testParty"),
			},
			want:        &types.MsgAddKeyringPartyResponse{},
			wantKeyring: &wantKr,
		},
		{
			name: "FAIL: keyring not found",
			args: args{
				keyring: &defaultKr,
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
				keyring: &defaultKr,
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
			if (err != nil) != tt.wantErr {
				t.Fatalf("AddKeyringParty() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("AddKeyringParty() got = %v, want %v", got, tt.want)
				}

				gotKeyring, err := ik.KeyringStore.Get(ctx, tt.args.keyring.Address)
				require.NoError(t, err)

				if !reflect.DeepEqual(&gotKeyring, tt.wantKeyring) {
					t.Fatalf("NewKeyring() got = %v, want %v", gotKeyring, tt.wantKeyring)
				}
			}
		})
	}
}
