package v5_test

import (
	"testing"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"
	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	v5 "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/migrations/v5"
	treasury "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"
)

func TestMigrate(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(treasury.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)

	kvStoreService := runtime.NewKVStoreService(storeKey)
	sb := collections.NewSchemaBuilder(kvStoreService)

	keyStore := collections.NewMap(sb, types.KeysKey, types.KeysIndex, collections.Uint64Key, codec.CollValue[types.Key](cdc))

	zenzecKey := types.Key{
		Id:            1,
		KeyringAddr:   "keyring1k6vc6vhp6e6l3rxalue9v4ux",
		PublicKey:     []byte("A0m1rS0iR/3yLSL9kX5KrLBUO3hPKKpY7LhKLnOU9u0p"),
		Type:          types.KeyType_KEY_TYPE_BITCOIN_SECP256K1,
		WorkspaceAddr: "workspace1mphgzyhncnzyggfxmv4nmh",
		ZenbtcMetadata: &types.ZenBTCMetadata{
			Asset:         dcttypes.Asset_ASSET_ZENZEC,
			Caip2ChainId:  "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
			ChainType:     types.WalletType_WALLET_TYPE_SOLANA,
			RecipientAddr: "4LXTk6h63VNt2Pbe1ARucfuvM6qKD6pecGjFiCi7rbpL",
		},
	}

	zenbtcKey := types.Key{
		Id:            2,
		KeyringAddr:   "keyring1k6vc6vhp6e6l3rxalue9v4ux",
		PublicKey:     []byte("A0m1rS0iR/3yLSL9kX5KrLBUO3hPKKpY7LhKLnOU9u0p"),
		Type:          types.KeyType_KEY_TYPE_BITCOIN_SECP256K1,
		WorkspaceAddr: "workspace1mphgzyhncnzyggfxmv4nmh",
		ZenbtcMetadata: &types.ZenBTCMetadata{
			Caip2ChainId:  "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
			ChainType:     types.WalletType_WALLET_TYPE_SOLANA,
			RecipientAddr: "4LXTk6h63VNt2Pbe1ARucfuvM6qKD6pecGjFiCi7rbpL",
		},
	}

	zenbtcKey2 := types.Key{
		Id:            3,
		KeyringAddr:   "keyring1k6vc6vhp6e6l3rxalue9v4ux",
		PublicKey:     []byte("A0m1rS0iR/3yLSL9kX5KrLBUO3hPKKpY7LhKLnOU9u0p"),
		Type:          types.KeyType_KEY_TYPE_BITCOIN_SECP256K1,
		WorkspaceAddr: "workspace1mphgzyhncnzyggfxmv4nmh",
		ZenbtcMetadata: &types.ZenBTCMetadata{
			Asset:         dcttypes.Asset_ASSET_ZENBTC,
			Caip2ChainId:  "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
			ChainType:     types.WalletType_WALLET_TYPE_SOLANA,
			RecipientAddr: "4LXTk6h63VNt2Pbe1ARucfuvM6qKD6pecGjFiCi7rbpL",
		},
	}

	zenbtcKey3 := types.Key{
		Id:            4,
		KeyringAddr:   "keyring1k6vc6vhp6e6l3rxalue9v4ux",
		PublicKey:     []byte("A0m1rS0iR/3yLSL9kX5KrLBUO3hPKKpY7LhKLnOU9u0p"),
		Type:          types.KeyType_KEY_TYPE_BITCOIN_SECP256K1,
		WorkspaceAddr: "workspace1mphgzyhncnzyggfxmv4nmh",
		ZenbtcMetadata: &types.ZenBTCMetadata{
			Asset:         dcttypes.Asset_ASSET_UNSPECIFIED,
			Caip2ChainId:  "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
			ChainType:     types.WalletType_WALLET_TYPE_SOLANA,
			RecipientAddr: "4LXTk6h63VNt2Pbe1ARucfuvM6qKD6pecGjFiCi7rbpL",
		},
	}

	err := keyStore.Set(ctx, zenzecKey.Id, zenzecKey)
	require.NoError(t, err)
	err = keyStore.Set(ctx, zenbtcKey.Id, zenbtcKey)
	require.NoError(t, err)
	err = keyStore.Set(ctx, zenbtcKey2.Id, zenbtcKey2)
	require.NoError(t, err)
	err = keyStore.Set(ctx, zenbtcKey3.Id, zenbtcKey3)
	require.NoError(t, err)

	require.NoError(t, v5.UpdateZenbtcKeys(ctx, keyStore, cdc))

	// Check keys after migration
	migratedZenzecKey, err := keyStore.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, dcttypes.Asset_ASSET_ZENZEC, migratedZenzecKey.ZenbtcMetadata.Asset)

	migratedZenbtcKey2, err := keyStore.Get(ctx, 2)
	require.NoError(t, err)
	require.Equal(t, dcttypes.Asset_ASSET_ZENBTC, migratedZenbtcKey2.ZenbtcMetadata.Asset)

	migratedZenbtcKey3, err := keyStore.Get(ctx, 3)
	require.NoError(t, err)
	require.Equal(t, dcttypes.Asset_ASSET_ZENBTC, migratedZenbtcKey3.ZenbtcMetadata.Asset)

	migratedZenbtcKey4, err := keyStore.Get(ctx, 4)
	require.NoError(t, err)
	require.Equal(t, dcttypes.Asset_ASSET_ZENBTC, migratedZenbtcKey4.ZenbtcMetadata.Asset)
}
