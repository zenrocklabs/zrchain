package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	zentp "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
)

func Test_QueryMintsAndBurns(t *testing.T) {

	mints := testutil.DefaultMints

	tests := []struct {
		name    string
		req     *types.QueryMintsRequest
		want    *types.QueryMintsResponse
		wantErr bool
	}{
		{
			name: "PASS: default",
			req: &types.QueryMintsRequest{
				Pagination: &query.PageRequest{},
			},
			want: &types.QueryMintsResponse{
				Mints: []*types.Bridge{&mints[0], &mints[1], &mints[2]},
				Pagination: &query.PageResponse{
					Total: uint64(len(mints)),
				},
			},
		},
		{
			name: "PASS: filter by id",
			req: &types.QueryMintsRequest{
				Id:         1,
				Pagination: &query.PageRequest{},
			},
			want: &types.QueryMintsResponse{
				Mints: []*types.Bridge{&mints[0]},
				Pagination: &query.PageResponse{
					Total: 1,
				},
			},
		}, {
			name: "PASS: filter by status",
			req: &types.QueryMintsRequest{
				Status: types.BridgeStatus_BRIDGE_STATUS_PENDING,
			},
			want: &types.QueryMintsResponse{
				Mints: []*types.Bridge{&mints[2]},
				Pagination: &query.PageResponse{
					Total: 1,
				},
			},
		},
		{
			name: "PASS: filter by creator",
			req: &types.QueryMintsRequest{
				Creator:    "zenrock1vpg325p5a7v73jaqyk02muh6n22p852h60anxx",
				Pagination: &query.PageRequest{},
			},
			want: &types.QueryMintsResponse{
				Mints: []*types.Bridge{&mints[0], &mints[2]},
				Pagination: &query.PageResponse{
					Total: 2,
				},
			},
		},
		{
			name: "PASS: filter by denom",
			req: &types.QueryMintsRequest{
				Denom:      "urock",
				Pagination: &query.PageRequest{},
			},
			want: &types.QueryMintsResponse{
				Mints: []*types.Bridge{&mints[0], &mints[1], &mints[2]},
				Pagination: &query.PageResponse{
					Total: 3,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zk, ctx := keepertest.ZentpKeeper(t)

			genesis := types.GenesisState{
				Params: types.DefaultParams(),
				Mints:  mints,
			}
			zentp.InitGenesis(ctx, zk, genesis)

			got, err := zk.Mints(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mints() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
