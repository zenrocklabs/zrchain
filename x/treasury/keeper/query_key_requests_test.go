package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	treasury "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
)

func TestKeeper_KeyRequests(t *testing.T) {

	type args struct {
		keyReqs []types.KeyRequest
		req     *types.QueryKeyRequestsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *types.QueryKeyRequestsResponse
		wantErr bool
	}{
		{
			name: "PASS: return key requests for a workspace and a keyring",
			args: args{
				keyReqs: []types.KeyRequest{defaultKeyRequest},
				req: &types.QueryKeyRequestsRequest{
					KeyringAddr: defaultKeyRequest.KeyringAddr,
				},
			},
			want: &types.QueryKeyRequestsResponse{
				KeyRequests: []*types.KeyReqResponse{&defaultKeyReqResponse},
				Pagination:  &query.PageResponse{Total: 1},
			},
		},
		{
			name: "PASS: return only the request that matches the req.KeyringAddr",
			args: args{
				keyReqs: []types.KeyRequest{keyRequestWithKeyring2, defaultKeyRequest},
				req: &types.QueryKeyRequestsRequest{
					KeyringAddr: defaultKeyRequest.KeyringAddr,
				},
			},
			want: &types.QueryKeyRequestsResponse{
				KeyRequests: []*types.KeyReqResponse{&defaultKeyReqResponse},
				Pagination:  &query.PageResponse{Total: 1},
			},
		},
		{
			name: "PASS: return only the request that matches the req.Status",
			args: args{
				keyReqs: []types.KeyRequest{keyRequestWithKeyring3, keyRequestWithKeyring2, defaultKeyRequest},
				req: &types.QueryKeyRequestsRequest{
					Status: types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
				},
			},
			want: &types.QueryKeyRequestsResponse{
				KeyRequests: []*types.KeyReqResponse{&defaultKeyReqResponse, &defaultKeyReqResponse2},
				Pagination:  &query.PageResponse{Total: 2},
			},
		},
		{
			name: "PASS: return only the request that matches the req.WorkspaceAddr",
			args: args{
				keyReqs: []types.KeyRequest{keyRequestWithKeyring3, keyRequestWithKeyring2, defaultKeyRequest},
				req: &types.QueryKeyRequestsRequest{
					WorkspaceAddr: keyRequestWithKeyring3.WorkspaceAddr,
				},
			},
			want: &types.QueryKeyRequestsResponse{
				KeyRequests: []*types.KeyReqResponse{&defaultKeyReqResponse3},
				Pagination:  &query.PageResponse{Total: 1},
			},
		},
		{
			name: "PASS: keyringAddr with no key requests",
			args: args{
				keyReqs: []types.KeyRequest{defaultKeyRequest},
				req: &types.QueryKeyRequestsRequest{
					KeyringAddr: "notAKeyringAddr",
				},
			},
			want: &types.QueryKeyRequestsResponse{
				KeyRequests: nil,
				Pagination:  &query.PageResponse{},
			},
			wantErr: false,
		},
		{
			name: "PASS: all requests (no flags provided)",
			args: args{
				keyReqs: []types.KeyRequest{defaultKeyRequest},
				req:     &types.QueryKeyRequestsRequest{},
			},
			want: &types.QueryKeyRequestsResponse{
				KeyRequests: []*types.KeyReqResponse{&defaultKeyReqResponse},
				Pagination:  &query.PageResponse{Total: 1},
			},
			wantErr: false,
		},
		{
			name: "PASS: no requests of matching type",
			args: args{
				keyReqs: []types.KeyRequest{defaultKeyRequest},
				req: &types.QueryKeyRequestsRequest{
					Status: types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
				},
			},
			want: &types.QueryKeyRequestsResponse{
				KeyRequests: nil,
				Pagination:  &query.PageResponse{},
			},
			wantErr: false,
		},
		{
			name: "FAIL: invalid request",
			args: args{
				keyReqs: []types.KeyRequest{defaultKeyRequest},
				req:     nil,
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
				KeyRequests: tt.args.keyReqs,
			}
			treasury.InitGenesis(ctx, *tk, genesis)

			got, err := tk.KeyRequests(ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
