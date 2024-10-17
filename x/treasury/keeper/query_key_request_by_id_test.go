package keeper_test

import (
	"reflect"
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v4/testutil/keeper"
	treasury "github.com/Zenrock-Foundation/zrchain/v4/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v4/x/treasury/types"
)

func TestKeeper_KeyRequestByID(t *testing.T) {

	type args struct {
		keyReq *types.KeyRequest
		req    *types.QueryKeyRequestByIDRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *types.QueryKeyRequestByIDResponse
		wantErr bool
	}{
		{
			name: "PASS: get a key request by id",
			args: args{
				keyReq: &defaultKeyRequest,
				req: &types.QueryKeyRequestByIDRequest{
					Id: 1,
				},
			},
			want: &types.QueryKeyRequestByIDResponse{
				KeyRequest: &defaultKeyReqResponse,
			},
		},
		{
			name: "FAIL: key request not found",
			args: args{
				keyReq: &defaultKeyRequest,
				req: &types.QueryKeyRequestByIDRequest{
					Id: 10000,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "FAIL: invalid request",
			args: args{
				keyReq: &defaultKeyRequest,
				req:    nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			tk := keepers.TreasuryKeeper
			ctx := keepers.Ctx

			genesis := types.GenesisState{
				PortId:      types.PortID,
				KeyRequests: []types.KeyRequest{*tt.args.keyReq},
			}
			treasury.InitGenesis(ctx, *tk, genesis)

			got, err := tk.KeyRequestByID(ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyRequestByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyRequestByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
