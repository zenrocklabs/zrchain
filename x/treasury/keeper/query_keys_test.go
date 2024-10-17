package keeper_test

import (
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"

	keepertest "github.com/Zenrock-Foundation/zrchain/v4/testutil/keeper"
	treasury "github.com/Zenrock-Foundation/zrchain/v4/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v4/x/treasury/types"
)

func TestKeeper_Keys(t *testing.T) {
	type args struct {
		keys []types.Key
		req  *types.QueryKeysRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *types.QueryKeysResponse
		wantErr bool
	}{
		{
			name: "PASS: ecdsa - return key requests for a workspace and a keyring",
			args: args{
				keys: []types.Key{defaultECDSAKey},
				req:  &types.QueryKeysRequest{},
			},
			want: &types.QueryKeysResponse{
				Keys: []*types.KeyAndWalletResponse{{Key: &defaultECDSAKeyResponse, Wallets: []*types.WalletResponse{
					{Address: "zen1tun7x3s2ywksa32nl38d3fuuv8nk5ang97v73r", Type: types.WalletType_WALLET_TYPE_NATIVE.String()},
					{Address: "0xdEa33aE3DA8f2EbA6efBB3EF5d143415438a6541", Type: types.WalletType_WALLET_TYPE_EVM.String()},
				}}},
				Pagination: &query.PageResponse{Total: 1},
			},
		},
		{
			name: "PASS: Bitcoin - return key requests for a workspace and a keyring",
			args: args{
				keys: []types.Key{defaultBitcoinKey},
				req:  &types.QueryKeysRequest{},
			},
			want: &types.QueryKeysResponse{
				Keys: []*types.KeyAndWalletResponse{{Key: &defaultBitcoinKeyResponse, Wallets: []*types.WalletResponse{
					{Address: "tb1qtun7x3s2ywksa32nl38d3fuuv8nk5angr0m8zv", Type: types.WalletType_WALLET_TYPE_BTC_TESTNET.String()},
					{Address: "bc1qtun7x3s2ywksa32nl38d3fuuv8nk5angffq5el", Type: types.WalletType_WALLET_TYPE_BTC_MAINNET.String()},
					{Address: "bcrt1qtun7x3s2ywksa32nl38d3fuuv8nk5angpxz249", Type: types.WalletType_WALLET_TYPE_BTC_REGNET.String()},
				}}},
				Pagination: &query.PageResponse{Total: 1},
			},
		},

		{
			name: "PASS: ecdsa - return keys for native addresses",
			args: args{
				keys: []types.Key{defaultECDSAKey},
				req: &types.QueryKeysRequest{
					WalletType: types.WalletType_WALLET_TYPE_NATIVE,
				},
			},
			want: &types.QueryKeysResponse{
				Keys: []*types.KeyAndWalletResponse{{Key: &defaultECDSAKeyResponse, Wallets: []*types.WalletResponse{
					{Address: "zen1tun7x3s2ywksa32nl38d3fuuv8nk5ang97v73r", Type: types.WalletType_WALLET_TYPE_NATIVE.String()},
				}}},
				Pagination: &query.PageResponse{Total: 1},
			},
		},
		{
			name: "PASS: ecdsa - return keys for native addresses (prefixed)",
			args: args{
				keys: []types.Key{defaultECDSAKey},
				req: &types.QueryKeysRequest{
					WalletType: types.WalletType_WALLET_TYPE_NATIVE,
					Prefixes:   []string{"celestia", "cosmos", "zen"},
				},
			},
			want: &types.QueryKeysResponse{
				Keys: []*types.KeyAndWalletResponse{{Key: &defaultECDSAKeyResponse, Wallets: []*types.WalletResponse{
					{Address: "celestia1tun7x3s2ywksa32nl38d3fuuv8nk5angwe008l", Type: types.WalletType_WALLET_TYPE_NATIVE.String()},
					{Address: "cosmos1tun7x3s2ywksa32nl38d3fuuv8nk5angln7laj", Type: types.WalletType_WALLET_TYPE_NATIVE.String()},
					{Address: "zen1tun7x3s2ywksa32nl38d3fuuv8nk5ang97v73r", Type: types.WalletType_WALLET_TYPE_NATIVE.String()},
				}}},
				Pagination: &query.PageResponse{Total: 1},
			},
		},
		{
			name: "PASS: ecdsa - return keys for eth addresses",
			args: args{
				keys: []types.Key{defaultECDSAKey},
				req: &types.QueryKeysRequest{
					WalletType: types.WalletType_WALLET_TYPE_EVM,
				},
			},
			want: &types.QueryKeysResponse{
				Keys: []*types.KeyAndWalletResponse{{Key: &defaultECDSAKeyResponse, Wallets: []*types.WalletResponse{
					{Address: "0xdEa33aE3DA8f2EbA6efBB3EF5d143415438a6541", Type: types.WalletType_WALLET_TYPE_EVM.String()},
				}}},
				Pagination: &query.PageResponse{Total: 1},
			},
		},
		{
			name: "PASS: ecdsa - return keys for celestia addresses",
			args: args{
				keys: []types.Key{defaultECDSAKey},
				req: &types.QueryKeysRequest{
					WalletType: types.WalletType_WALLET_TYPE_NATIVE,
					Prefixes:   []string{"celestia"},
				},
			},
			want: &types.QueryKeysResponse{
				Keys: []*types.KeyAndWalletResponse{{Key: &defaultECDSAKeyResponse, Wallets: []*types.WalletResponse{
					{Address: "celestia1tun7x3s2ywksa32nl38d3fuuv8nk5angwe008l", Type: types.WalletType_WALLET_TYPE_NATIVE.String()},
				}}},
				Pagination: &query.PageResponse{Total: 1},
			},
		},
		{
			name: "PASS: eddsa - return keys for eddsa addresses",
			args: args{
				keys: []types.Key{defaultEdDSAKey},
				req: &types.QueryKeysRequest{
					WalletType: types.WalletType_WALLET_TYPE_UNSPECIFIED,
				},
			},
			want: &types.QueryKeysResponse{
				Keys: []*types.KeyAndWalletResponse{{Key: &defaultEdDSAKeyResponse, Wallets: []*types.WalletResponse{
					{Address: "3pPzT5vhum6GN5pfKVAqB5MA4C3sZmZviyZTYpqRsz6R", Type: types.WalletType_WALLET_TYPE_SOLANA.String()},
				}}},
				Pagination: &query.PageResponse{Total: 1},
			},
		},
		{
			name: "PASS: eddsa - return keys for all addresses",
			args: args{
				keys: []types.Key{defaultEdDSAKey},
				req: &types.QueryKeysRequest{
					WalletType: types.WalletType_WALLET_TYPE_UNSPECIFIED,
				},
			},
			want: &types.QueryKeysResponse{
				Keys: []*types.KeyAndWalletResponse{{Key: &defaultEdDSAKeyResponse, Wallets: []*types.WalletResponse{
					{Address: "3pPzT5vhum6GN5pfKVAqB5MA4C3sZmZviyZTYpqRsz6R", Type: types.WalletType_WALLET_TYPE_SOLANA.String()},
				}}},
				Pagination: &query.PageResponse{Total: 1},
			},
		},
		{
			name: "PASS: ecdsa - return keys for all addresses from a specific workspace",
			args: args{
				keys: []types.Key{defaultECDSAKey},
				req: &types.QueryKeysRequest{
					WorkspaceAddr: defaultECDSAKey.WorkspaceAddr,
				},
			},
			want: &types.QueryKeysResponse{
				Keys: []*types.KeyAndWalletResponse{{Key: &defaultECDSAKeyResponse, Wallets: []*types.WalletResponse{
					{Address: "zen1tun7x3s2ywksa32nl38d3fuuv8nk5ang97v73r", Type: types.WalletType_WALLET_TYPE_NATIVE.String()},
					{Address: "0xdEa33aE3DA8f2EbA6efBB3EF5d143415438a6541", Type: types.WalletType_WALLET_TYPE_EVM.String()},
				}}},
				Pagination: &query.PageResponse{Total: 1},
			},
		},
		{
			name: "PASS: ecdsa - return keys for workspace with no keys (nothing returned)",
			args: args{
				keys: []types.Key{defaultECDSAKey},
				req: &types.QueryKeysRequest{
					WorkspaceAddr: "anotherWorkspace",
				},
			},
			want:    &types.QueryKeysResponse{Pagination: &query.PageResponse{}},
			wantErr: false,
		},
		{
			name: "FAIL: request is nil",
			args: args{
				keys: []types.Key{defaultECDSAKey},
				req:  nil,
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
				PortId: types.PortID,
				Keys:   tt.args.keys,
			}
			treasury.InitGenesis(ctx, *tk, genesis)

			got, err := tk.Keys(ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() got = %v, want %v", got, tt.want)
			}
		})
	}
}
