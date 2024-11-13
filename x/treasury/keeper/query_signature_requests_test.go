package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	treasury "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
)

func TestKeeper_SignRequest(t *testing.T) {

	var defaultSignReq = types.SignRequest{
		Id:             1,
		Creator:        "testCreator",
		KeyId:          1,
		DataForSigning: [][]byte{[]byte("test")},
		Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
		KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
	}

	var defaultSignReqResponse = types.SignReqResponse{
		Id:             1,
		Creator:        "testCreator",
		KeyId:          1,
		DataForSigning: [][]byte{[]byte("test")},
		Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING.String(),
		KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1.String(),
	}

	var signReqExisting = types.SignRequest{
		Id:             2,
		Creator:        "testCreator",
		KeyId:          2,
		DataForSigning: [][]byte{[]byte("test")},
		Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
		KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
	}

	var emptyKeyringinKey = types.Key{
		Id:            1,
		WorkspaceAddr: "testWorkspace",
		KeyringAddr:   "",
		Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
		PublicKey:     []byte{0x03, 0xca, 0x27, 0xea, 0x7b, 0x06, 0x41, 0x49, 0x7b, 0x19, 0xa7, 0x23, 0xe3, 0xb9, 0x25, 0x90, 0x80, 0x1c, 0x7e, 0x79, 0xb1, 0x14, 0x25, 0x3f, 0xc1, 0xe9, 0x9d, 0xf1, 0xfd, 0x97, 0x30, 0x52, 0x6b},
	}

	type args struct {
		keys     []types.Key
		signReqs []types.SignRequest
		req      *types.QuerySignatureRequestsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *types.QuerySignatureRequestsResponse
		wantErr bool
	}{
		{
			name: "PASS: get all signature requests from keyring",
			args: args{
				keys:     []types.Key{defaultECDSAKey},
				signReqs: []types.SignRequest{defaultSignReq},
				req: &types.QuerySignatureRequestsRequest{
					KeyringAddr: defaultECDSAKey.KeyringAddr,
				},
			},
			want: &types.QuerySignatureRequestsResponse{
				SignRequests: []*types.SignReqResponse{&defaultSignReqResponse},
				Pagination:   &query.PageResponse{Total: 1},
			},
		},
		{
			name: "PASS: ignore requests to different keyring",
			args: args{
				keys:     []types.Key{secondECDSAKey, defaultECDSAKey},
				signReqs: []types.SignRequest{signReqExisting, defaultSignReq},
				req: &types.QuerySignatureRequestsRequest{
					KeyringAddr: defaultECDSAKey.KeyringAddr,
				},
			},
			want: &types.QuerySignatureRequestsResponse{
				SignRequests: []*types.SignReqResponse{&defaultSignReqResponse},
				Pagination:   &query.PageResponse{Total: 1},
			},
		},
		{
			name: "PASS: keyring not found (nothing returned)",
			args: args{
				keys:     []types.Key{defaultECDSAKey},
				signReqs: []types.SignRequest{defaultSignReq},
				req: &types.QuerySignatureRequestsRequest{
					KeyringAddr: "notAKeyring",
				},
			},
			want: &types.QuerySignatureRequestsResponse{
				SignRequests: nil,
				Pagination:   &query.PageResponse{},
			},
		},
		{
			name: "PASS: key not found (nothing returned)",
			args: args{
				keys: []types.Key{defaultECDSAKey},
				signReqs: []types.SignRequest{{
					Id:      1,
					Creator: "testCreator",
					KeyId:   50756,
				}},
				req: &types.QuerySignatureRequestsRequest{
					KeyringAddr: defaultKr.Address,
				},
			},
			want: &types.QuerySignatureRequestsResponse{
				SignRequests: nil,
				Pagination:   &query.PageResponse{},
			},
		},
		{
			name: "PASS: keyringAddr field is empty (no filter applied)",
			args: args{
				keys:     []types.Key{emptyKeyringinKey},
				signReqs: []types.SignRequest{signReqExisting, defaultSignReq},
				req: &types.QuerySignatureRequestsRequest{
					KeyringAddr: "",
				},
			},
			want: &types.QuerySignatureRequestsResponse{
				SignRequests: []*types.SignReqResponse{&defaultSignReqResponse},
				Pagination:   &query.PageResponse{Total: 1},
			},
		},
		{
			name: "PASS: get all signature requests with a specific status",
			args: args{
				keys: []types.Key{defaultECDSAKey},
				signReqs: []types.SignRequest{
					{
						Id:      1,
						Creator: "testCreator",
						KeyId:   1,
						KeyType: types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
						Status:  types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
					},
					{
						Id:      2,
						Creator: "testCreator",
						KeyId:   1,
						KeyType: types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
						Status:  types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
					},
					{
						Id:      3,
						Creator: "testCreator",
						KeyId:   1,
						KeyType: types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
						Status:  types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED,
					},
				},
				req: &types.QuerySignatureRequestsRequest{
					KeyringAddr: "keyring1pfnq7r04rept47gaf5cpdew2",
					Status:      types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
				},
			},
			want: &types.QuerySignatureRequestsResponse{
				SignRequests: []*types.SignReqResponse{
					{
						Id:      2,
						Creator: "testCreator",
						KeyId:   1,
						KeyType: types.KeyType_KEY_TYPE_ECDSA_SECP256K1.String(),
						Status:  types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING.String(),
					},
				},
				Pagination: &query.PageResponse{Total: 1},
			},
		},
		{
			name: "FAIL: invalid request",
			args: args{
				keys:     []types.Key{defaultECDSAKey},
				signReqs: []types.SignRequest{defaultSignReq},
				req:      nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			tk := keepers.TreasuryKeeper
			ctx := keepers.Ctx

			genesis := types.GenesisState{
				PortId:       types.PortID,
				Keys:         tt.args.keys,
				SignRequests: tt.args.signReqs,
			}
			treasury.InitGenesis(ctx, *tk, genesis)

			got, err := tk.SignatureRequests(ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
