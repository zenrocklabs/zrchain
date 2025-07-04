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

func Test_QueryMints(t *testing.T) {

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

func Test_QueryBurns(t *testing.T) {

	burns := testutil.DefaultBurns

	tests := []struct {
		name    string
		req     *types.QueryBurnsRequest
		want    *types.QueryBurnsResponse
		wantErr bool
	}{
		{
			name: "PASS: default",
			req: &types.QueryBurnsRequest{
				Pagination: &query.PageRequest{},
			},
			want: &types.QueryBurnsResponse{
				Burns: []*types.Bridge{&burns[0], &burns[1], &burns[2], &burns[3], &burns[4], &burns[5]},
				Pagination: &query.PageResponse{
					Total: uint64(len(burns)),
				},
			},
		},
		{
			name: "PASS: filter source tx hash",
			req: &types.QueryBurnsRequest{
				SourceTxHash: burns[0].TxHash,
				Pagination:   &query.PageRequest{},
			},
			want: &types.QueryBurnsResponse{
				Burns: []*types.Bridge{&burns[0]},
				Pagination: &query.PageResponse{
					Total: 1,
				},
			},
		},
		{
			name: "PASS: no tx hash found",
			req: &types.QueryBurnsRequest{
				SourceTxHash: "not-found",
				Pagination:   &query.PageRequest{},
			},
			want: &types.QueryBurnsResponse{
				Burns: nil,
				Pagination: &query.PageResponse{
					Total: 0,
				},
			},
		},
		{
			name: "PASS: tx id",
			req: &types.QueryBurnsRequest{
				Id:         4,
				Pagination: &query.PageRequest{},
			},
			want: &types.QueryBurnsResponse{
				Burns: []*types.Bridge{&burns[0]},
				Pagination: &query.PageResponse{
					Total: 1,
				},
			},
		},
		{
			name: "PASS: filter by recipient address",
			req: &types.QueryBurnsRequest{
				RecipientAddress: burns[0].RecipientAddress,
				Pagination:       &query.PageRequest{},
			},
			want: &types.QueryBurnsResponse{
				Burns: []*types.Bridge{&burns[0], &burns[3]},
				Pagination: &query.PageResponse{
					Total: 2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zk, ctx := keepertest.ZentpKeeper(t)

			genesis := types.GenesisState{
				Params: types.DefaultParams(),
				Burns:  burns,
			}
			zentp.InitGenesis(ctx, zk, genesis)

			got, err := zk.Burns(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Burns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
