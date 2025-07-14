package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	treasury "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"github.com/zenrocklabs/goem/ethereum"
)

var defaultECDSAKeyWithZenBTCMetadata = types.Key{
	Id:            1,
	WorkspaceAddr: "testWorkspace",
	KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
	Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
	PublicKey:     []byte{0x03, 0xca, 0x27, 0xea, 0x7b, 0x06, 0x41, 0x49, 0x7b, 0x19, 0xa7, 0x23, 0xe3, 0xb9, 0x25, 0x90, 0x80, 0x1c, 0x7e, 0x79, 0xb1, 0x14, 0x25, 0x3f, 0xc1, 0xe9, 0x9d, 0xf1, 0xfd, 0x97, 0x30, 0x52, 0x6b},
	ZenbtcMetadata: &types.ZenBTCMetadata{
		RecipientAddr: "0x9D450478FDB879C2900Ad54A0A407B0607b20478",
		ChainType:     types.WalletType_WALLET_TYPE_EVM,
		Caip2ChainId:  ethereum.HoodiCAIP2,
		ReturnAddress: "tb1qypwjx7yj5jz0gw0vh76348ypa2ns7tfwsnhlh9",
	},
}

var defaultECDSAKeyResponseWithZenBTCMetadata = types.KeyResponse{
	Id:            1,
	WorkspaceAddr: "testWorkspace",
	KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
	Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1.String(),
	PublicKey:     []byte{0x03, 0xca, 0x27, 0xea, 0x7b, 0x06, 0x41, 0x49, 0x7b, 0x19, 0xa7, 0x23, 0xe3, 0xb9, 0x25, 0x90, 0x80, 0x1c, 0x7e, 0x79, 0xb1, 0x14, 0x25, 0x3f, 0xc1, 0xe9, 0x9d, 0xf1, 0xfd, 0x97, 0x30, 0x52, 0x6b},
	ZenbtcMetadata: &types.ZenBTCMetadata{
		RecipientAddr: "0x9D450478FDB879C2900Ad54A0A407B0607b20478",
		ChainType:     types.WalletType_WALLET_TYPE_EVM,
		Caip2ChainId:  ethereum.HoodiCAIP2,
		ReturnAddress: "tb1qypwjx7yj5jz0gw0vh76348ypa2ns7tfwsnhlh9",
	},
}

func TestKeeper_ZenBTCMetadata(t *testing.T) {
	type args struct {
		keys []types.Key
		req  *types.QueryZenbtcWalletsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *types.QueryZenbtcWalletsResponse
		wantErr bool
	}{
		{
			name: "PASS: return zenbtc metadata for empty request",
			args: args{
				keys: []types.Key{defaultECDSAKeyWithZenBTCMetadata},
				req:  &types.QueryZenbtcWalletsRequest{},
			},
			want: &types.QueryZenbtcWalletsResponse{
				ZenbtcWallets: []*types.KeyAndWalletResponse{
					{
						Key: &defaultECDSAKeyResponseWithZenBTCMetadata,
						Wallets: []*types.WalletResponse{
							{Address: "zen1tun7x3s2ywksa32nl38d3fuuv8nk5ang97v73r", Type: types.WalletType_WALLET_TYPE_NATIVE.String()},
							{Address: "0xdEa33aE3DA8f2EbA6efBB3EF5d143415438a6541", Type: types.WalletType_WALLET_TYPE_EVM.String()},
						},
					},
				},
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   1,
				},
			},
		},
		{
			name: "PASS: return zenbtc metadata for a recipient address request",
			args: args{
				keys: []types.Key{defaultECDSAKeyWithZenBTCMetadata},
				req: &types.QueryZenbtcWalletsRequest{
					RecipientAddr: "0x9D450478FDB879C2900Ad54A0A407B0607b20478",
				},
			},
			want: &types.QueryZenbtcWalletsResponse{
				ZenbtcWallets: []*types.KeyAndWalletResponse{
					{
						Key: &defaultECDSAKeyResponseWithZenBTCMetadata,
						Wallets: []*types.WalletResponse{
							{Address: "zen1tun7x3s2ywksa32nl38d3fuuv8nk5ang97v73r", Type: types.WalletType_WALLET_TYPE_NATIVE.String()},
							{Address: "0xdEa33aE3DA8f2EbA6efBB3EF5d143415438a6541", Type: types.WalletType_WALLET_TYPE_EVM.String()},
						},
					},
				},
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   1,
				},
			},
		},
		{
			name: "PASS: return zenbtc metadata for a chain type request",
			args: args{
				keys: []types.Key{defaultECDSAKeyWithZenBTCMetadata},
				req: &types.QueryZenbtcWalletsRequest{
					ChainType: types.WalletType_WALLET_TYPE_EVM,
				},
			},
			want: &types.QueryZenbtcWalletsResponse{
				ZenbtcWallets: []*types.KeyAndWalletResponse{
					{
						Key: &defaultECDSAKeyResponseWithZenBTCMetadata,
						Wallets: []*types.WalletResponse{
							{Address: "zen1tun7x3s2ywksa32nl38d3fuuv8nk5ang97v73r", Type: types.WalletType_WALLET_TYPE_NATIVE.String()},
							{Address: "0xdEa33aE3DA8f2EbA6efBB3EF5d143415438a6541", Type: types.WalletType_WALLET_TYPE_EVM.String()},
						},
					},
				},
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   1,
				},
			},
		},
		{
			name: "PASS: request only for EVM chain type",
			args: args{
				keys: []types.Key{
					{
						Id:            1,
						WorkspaceAddr: "testWorkspace",
						KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
						Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
						PublicKey:     []byte{0x03, 0xca, 0x27, 0xea, 0x7b, 0x06, 0x41, 0x49, 0x7b, 0x19, 0xa7, 0x23, 0xe3, 0xb9, 0x25, 0x90, 0x80, 0x1c, 0x7e, 0x79, 0xb1, 0x14, 0x25, 0x3f, 0xc1, 0xe9, 0x9d, 0xf1, 0xfd, 0x97, 0x30, 0x52, 0x6b},
						ZenbtcMetadata: &types.ZenBTCMetadata{
							RecipientAddr: "dummySolAddress",
							ChainType:     types.WalletType_WALLET_TYPE_SOLANA,
							Caip2ChainId:  ethereum.HoodiCAIP2,
							ReturnAddress: "tb1qypwjx7yj5jz0gw0vh76348ypa2ns7tfwsnhlh9",
						},
					},
					defaultECDSAKeyWithZenBTCMetadata,
				},
				req: &types.QueryZenbtcWalletsRequest{
					ChainType: types.WalletType_WALLET_TYPE_EVM,
				},
			},
			want: &types.QueryZenbtcWalletsResponse{
				ZenbtcWallets: []*types.KeyAndWalletResponse{
					{
						Key: &defaultECDSAKeyResponseWithZenBTCMetadata,
						Wallets: []*types.WalletResponse{
							{Address: "zen1tun7x3s2ywksa32nl38d3fuuv8nk5ang97v73r", Type: types.WalletType_WALLET_TYPE_NATIVE.String()},
							{Address: "0xdEa33aE3DA8f2EbA6efBB3EF5d143415438a6541", Type: types.WalletType_WALLET_TYPE_EVM.String()},
						},
					},
				},
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   1,
				},
			},
		},
		{
			name: "PASS: request Bitcoin wallet types for EVM chain type (metadata)",
			args: args{
				keys: []types.Key{
					{
						Id:            1,
						WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
						KeyringAddr:   "keyring1k6vc6vhp6e6l3rxalue9v4ux",
						Type:          types.KeyType_KEY_TYPE_BITCOIN_SECP256K1,
						PublicKey:     []byte{0x03, 0xca, 0x27, 0xea, 0x7b, 0x06, 0x41, 0x49, 0x7b, 0x19, 0xa7, 0x23, 0xe3, 0xb9, 0x25, 0x90, 0x80, 0x1c, 0x7e, 0x79, 0xb1, 0x14, 0x25, 0x3f, 0xc1, 0xe9, 0x9d, 0xf1, 0xfd, 0x97, 0x30, 0x52, 0x6b},
						ZenbtcMetadata: &types.ZenBTCMetadata{
							RecipientAddr: "0x9704Bc96D57180B3Cf4154fEf9Ba3A7aDFfDA9Ac",
							ChainType:     types.WalletType_WALLET_TYPE_EVM,
							Caip2ChainId:  "eip155:11555111",
							ReturnAddress: "tb1qypwjx7yj5jz0gw0vh76348ypa2ns7tfwsnhlh9",
						},
					},
				},
				req: &types.QueryZenbtcWalletsRequest{
					ChainType: types.WalletType_WALLET_TYPE_EVM,
				},
			},
			want: &types.QueryZenbtcWalletsResponse{
				ZenbtcWallets: []*types.KeyAndWalletResponse{
					{
						Key: &types.KeyResponse{
							Id:            1,
							WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
							KeyringAddr:   "keyring1k6vc6vhp6e6l3rxalue9v4ux",
							Type:          types.KeyType_KEY_TYPE_BITCOIN_SECP256K1.String(),
							PublicKey:     []byte{0x03, 0xca, 0x27, 0xea, 0x7b, 0x06, 0x41, 0x49, 0x7b, 0x19, 0xa7, 0x23, 0xe3, 0xb9, 0x25, 0x90, 0x80, 0x1c, 0x7e, 0x79, 0xb1, 0x14, 0x25, 0x3f, 0xc1, 0xe9, 0x9d, 0xf1, 0xfd, 0x97, 0x30, 0x52, 0x6b},
							ZenbtcMetadata: &types.ZenBTCMetadata{
								RecipientAddr: "0x9704Bc96D57180B3Cf4154fEf9Ba3A7aDFfDA9Ac",
								ChainType:     types.WalletType_WALLET_TYPE_EVM,
								Caip2ChainId:  "eip155:11555111",
								ReturnAddress: "tb1qypwjx7yj5jz0gw0vh76348ypa2ns7tfwsnhlh9",
							},
						},
						Wallets: []*types.WalletResponse{
							{Address: "tb1qtun7x3s2ywksa32nl38d3fuuv8nk5angr0m8zv", Type: types.WalletType_WALLET_TYPE_BTC_TESTNET.String()},
							{Address: "bc1qtun7x3s2ywksa32nl38d3fuuv8nk5angffq5el", Type: types.WalletType_WALLET_TYPE_BTC_MAINNET.String()},
							{Address: "bcrt1qtun7x3s2ywksa32nl38d3fuuv8nk5angpxz249", Type: types.WalletType_WALLET_TYPE_BTC_REGNET.String()},
						},
					},
				},
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   1,
				},
			},
		},
		{
			name: "PASS: return two keys with zenbtc metadata",
			args: args{
				keys: []types.Key{
					{
						Id:            2,
						WorkspaceAddr: "testWorkspace",
						KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
						Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
						PublicKey:     []byte{0x03, 0xca, 0x27, 0xea, 0x7b, 0x06, 0x41, 0x49, 0x7b, 0x19, 0xa7, 0x23, 0xe3, 0xb9, 0x25, 0x90, 0x80, 0x1c, 0x7e, 0x79, 0xb1, 0x14, 0x25, 0x3f, 0xc1, 0xe9, 0x9d, 0xf1, 0xfd, 0x97, 0x30, 0x52, 0x6b},
						ZenbtcMetadata: &types.ZenBTCMetadata{
							RecipientAddr: "anotherETHAddress",
							ChainType:     types.WalletType_WALLET_TYPE_EVM,
							Caip2ChainId:  ethereum.HoodiCAIP2,
							ReturnAddress: "tb1qypwjx7yj5jz0gw0vh76348ypa2ns7tfwsnhlh9",
						},
					},
					defaultECDSAKeyWithZenBTCMetadata,
				},
				req: &types.QueryZenbtcWalletsRequest{
					ChainType: types.WalletType_WALLET_TYPE_EVM,
				},
			},
			want: &types.QueryZenbtcWalletsResponse{
				ZenbtcWallets: []*types.KeyAndWalletResponse{
					{
						Key: &defaultECDSAKeyResponseWithZenBTCMetadata,
						Wallets: []*types.WalletResponse{
							{Address: "zen1tun7x3s2ywksa32nl38d3fuuv8nk5ang97v73r", Type: types.WalletType_WALLET_TYPE_NATIVE.String()},
							{Address: "0xdEa33aE3DA8f2EbA6efBB3EF5d143415438a6541", Type: types.WalletType_WALLET_TYPE_EVM.String()},
						},
					},
					{
						Key: &types.KeyResponse{
							Id:            2,
							WorkspaceAddr: "testWorkspace",
							KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
							Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1.String(),
							PublicKey:     []byte{0x03, 0xca, 0x27, 0xea, 0x7b, 0x06, 0x41, 0x49, 0x7b, 0x19, 0xa7, 0x23, 0xe3, 0xb9, 0x25, 0x90, 0x80, 0x1c, 0x7e, 0x79, 0xb1, 0x14, 0x25, 0x3f, 0xc1, 0xe9, 0x9d, 0xf1, 0xfd, 0x97, 0x30, 0x52, 0x6b},
							ZenbtcMetadata: &types.ZenBTCMetadata{
								RecipientAddr: "anotherETHAddress",
								ChainType:     types.WalletType_WALLET_TYPE_EVM,
								Caip2ChainId:  ethereum.HoodiCAIP2,
								ReturnAddress: "tb1qypwjx7yj5jz0gw0vh76348ypa2ns7tfwsnhlh9",
							},
						},
						Wallets: []*types.WalletResponse{
							{Address: "zen1tun7x3s2ywksa32nl38d3fuuv8nk5ang97v73r", Type: types.WalletType_WALLET_TYPE_NATIVE.String()},
							{Address: "0xdEa33aE3DA8f2EbA6efBB3EF5d143415438a6541", Type: types.WalletType_WALLET_TYPE_EVM.String()},
						},
					},
				},
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   2,
				},
			},
		},
		{
			name: "PASS: two filters applied",
			args: args{
				keys: []types.Key{
					{
						Id:            2,
						WorkspaceAddr: "testWorkspace",
						KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
						Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
						PublicKey:     []byte{0x03, 0xca, 0x27, 0xea, 0x7b, 0x06, 0x41, 0x49, 0x7b, 0x19, 0xa7, 0x23, 0xe3, 0xb9, 0x25, 0x90, 0x80, 0x1c, 0x7e, 0x79, 0xb1, 0x14, 0x25, 0x3f, 0xc1, 0xe9, 0x9d, 0xf1, 0xfd, 0x97, 0x30, 0x52, 0x6b},
						ZenbtcMetadata: &types.ZenBTCMetadata{
							RecipientAddr: "anotherETHAddress",
							ChainType:     types.WalletType_WALLET_TYPE_EVM,
							Caip2ChainId:  ethereum.HoodiCAIP2,
							ReturnAddress: "tb1qypwjx7yj5jz0gw0vh76348ypa2ns7tfwsnhlh9",
						},
					},
					defaultECDSAKeyWithZenBTCMetadata,
				},
				req: &types.QueryZenbtcWalletsRequest{
					RecipientAddr: "0x9D450478FDB879C2900Ad54A0A407B0607b20478",
					ChainType:     types.WalletType_WALLET_TYPE_EVM,
				},
			},
			want: &types.QueryZenbtcWalletsResponse{
				ZenbtcWallets: []*types.KeyAndWalletResponse{
					{
						Key: &defaultECDSAKeyResponseWithZenBTCMetadata,
						Wallets: []*types.WalletResponse{
							{Address: "zen1tun7x3s2ywksa32nl38d3fuuv8nk5ang97v73r", Type: types.WalletType_WALLET_TYPE_NATIVE.String()},
							{Address: "0xdEa33aE3DA8f2EbA6efBB3EF5d143415438a6541", Type: types.WalletType_WALLET_TYPE_EVM.String()},
						},
					},
				},
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   1,
				},
			},
		},
		{
			name: "PASS: key has no zenbtc metadata",
			args: args{
				keys: []types.Key{
					{
						Id:            1,
						WorkspaceAddr: "testWorkspace",
						KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
						Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
						PublicKey:     []byte{0x03, 0xca, 0x27, 0xea, 0x7b, 0x06, 0x41, 0x49, 0x7b, 0x19, 0xa7, 0x23, 0xe3, 0xb9, 0x25, 0x90, 0x80, 0x1c, 0x7e, 0x79, 0xb1, 0x14, 0x25, 0x3f, 0xc1, 0xe9, 0x9d, 0xf1, 0xfd, 0x97, 0x30, 0x52, 0x6b},
					},
				},
				req: &types.QueryZenbtcWalletsRequest{},
			},
			want: &types.QueryZenbtcWalletsResponse{
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   0,
				},
			},
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
		{
			name: "FAIL: ChainId does not match",
			args: args{
				keys: []types.Key{
					{
						Id:            1,
						WorkspaceAddr: "testWorkspace1",
						KeyringAddr:   "keyring1",
						Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
						ZenbtcMetadata: &types.ZenBTCMetadata{
							RecipientAddr: "recipient1",
							ChainType:     types.WalletType_WALLET_TYPE_EVM,
							Caip2ChainId:  "eip155:17001",
							ReturnAddress: "return1",
						},
					},
				},
				req: &types.QueryZenbtcWalletsRequest{
					MintChainId: ethereum.HoodiCAIP2,
					ChainType:   types.WalletType_WALLET_TYPE_EVM,
				},
			},
			want: &types.QueryZenbtcWalletsResponse{
				ZenbtcWallets: nil,
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   0,
				},
			},
		},
		{
			name: "FAIL: ChainType does not match",
			args: args{
				keys: []types.Key{
					{
						Id:            2,
						WorkspaceAddr: "testWorkspace2",
						KeyringAddr:   "keyring2",
						Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
						ZenbtcMetadata: &types.ZenBTCMetadata{
							RecipientAddr: "recipient2",
							ChainType:     types.WalletType_WALLET_TYPE_BTC_REGNET,
							Caip2ChainId:  ethereum.HoodiCAIP2,
							ReturnAddress: "return2",
						},
					},
				},
				req: &types.QueryZenbtcWalletsRequest{
					ChainType: types.WalletType_WALLET_TYPE_EVM,
				},
			},
			want: &types.QueryZenbtcWalletsResponse{
				ZenbtcWallets: nil,
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   0,
				},
			},
		},
		{
			name: "FAIL: RecipientAddr does not match",
			args: args{
				keys: []types.Key{
					{
						Id:            3,
						WorkspaceAddr: "testWorkspace3",
						KeyringAddr:   "keyring3",
						Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
						ZenbtcMetadata: &types.ZenBTCMetadata{
							RecipientAddr: "differentRecipient",
							ChainType:     types.WalletType_WALLET_TYPE_EVM,
							Caip2ChainId:  ethereum.HoodiCAIP2,
							ReturnAddress: "return3",
						},
					},
				},
				req: &types.QueryZenbtcWalletsRequest{
					RecipientAddr: "recipient2",
				},
			},
			want: &types.QueryZenbtcWalletsResponse{
				ZenbtcWallets: nil,
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   0,
				},
			},
		},
		{
			name: "FAIL: ReturnAddr does not match",
			args: args{
				keys: []types.Key{
					{
						Id:            4,
						WorkspaceAddr: "testWorkspace4",
						KeyringAddr:   "keyring4",
						Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
						ZenbtcMetadata: &types.ZenBTCMetadata{
							RecipientAddr: "recipient4",
							ChainType:     types.WalletType_WALLET_TYPE_EVM,
							Caip2ChainId:  ethereum.HoodiCAIP2,
							ReturnAddress: "differentReturn",
						},
					},
				},
				req: &types.QueryZenbtcWalletsRequest{
					ReturnAddr: "return3",
				},
			},
			want: &types.QueryZenbtcWalletsResponse{
				ZenbtcWallets: nil,
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   0,
				},
			},
		},
		{
			name: "FAIL: invalid chain type",
			args: args{
				keys: []types.Key{defaultECDSAKeyWithZenBTCMetadata},
				req: &types.QueryZenbtcWalletsRequest{
					ChainType: 10,
				},
			},
			want: &types.QueryZenbtcWalletsResponse{
				ZenbtcWallets: nil,
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   0,
				},
			},
		},
		{
			name: "FAIL: multiple filters applied with no matches",
			args: args{
				keys: []types.Key{
					defaultECDSAKeyWithZenBTCMetadata,
				},
				req: &types.QueryZenbtcWalletsRequest{
					RecipientAddr: "nonexistentRecipient",
					ChainType:     types.WalletType_WALLET_TYPE_EVM,
				},
			},
			want: &types.QueryZenbtcWalletsResponse{
				ZenbtcWallets: nil,
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   0,
				},
			},
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

			got, err := tk.ZenbtcWallets(ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
