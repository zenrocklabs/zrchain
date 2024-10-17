package keeper_test

import (
	"reflect"
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v4/testutil/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v4/x/identity/module"
	"github.com/Zenrock-Foundation/zrchain/v4/x/identity/types"
)

func TestKeeper_KeyringByAddress(t *testing.T) {

	var defaultKr = types.Keyring{
		Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
		Creator:     "testCreator",
		Description: "testDescription",
		Admins:      []string{"testCreator"},
		KeyReqFee:   0,
		SigReqFee:   0,
		IsActive:    true,
	}

	type args struct {
		keyring *types.Keyring
		req     *types.QueryKeyringByAddressRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *types.QueryKeyringByAddressResponse
		wantErr bool
	}{
		{
			name: "PASS: get a keyring by address",
			args: args{
				keyring: &defaultKr,
				req: &types.QueryKeyringByAddressRequest{
					KeyringAddr: "keyring1pfnq7r04rept47gaf5cpdew2",
				},
			},
			want: &types.QueryKeyringByAddressResponse{Keyring: &types.Keyring{
				Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
				Creator:     "testCreator",
				Description: "testDescription",
				Admins:      []string{"testCreator"},
				Parties:     nil,
				KeyReqFee:   0,
				SigReqFee:   0,
				IsActive:    true,
			}},
		},
		{
			name: "FAIL: keyring by address not found",
			args: args{
				keyring: &defaultKr,
				req: &types.QueryKeyringByAddressRequest{
					KeyringAddr: "noKeyringAddress",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "FAIL: invalid request",
			args: args{
				keyring: &defaultKr,
				req:     nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			ik := keepers.IdentityKeeper
			ctx := keepers.Ctx

			genesis := types.GenesisState{
				PortId:   types.PortID,
				Keyrings: []types.Keyring{*tt.args.keyring},
			}
			identity.InitGenesis(ctx, *ik, genesis)

			got, err := ik.KeyringByAddress(ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyringByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyringByAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}
