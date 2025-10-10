package v9_test

import (
	"testing"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"

	v9 "github.com/Zenrock-Foundation/zrchain/v6/x/validation/migrations/v9"
	validation "github.com/Zenrock-Foundation/zrchain/v6/x/validation/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"

	zenbtctypes "github.com/Zenrock-Foundation/zrchain/v6/x/zenbtc/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"
)

func TestClearEthereumNonceData_NonMainnet(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(validation.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey).WithChainID("gardia-9") // non-mainnet

	kvStoreService := runtime.NewKVStoreService(storeKey)
	sb := collections.NewSchemaBuilder(kvStoreService)

	ethNonceData := collections.NewMap(
		sb,
		types.LastUsedEthereumNonceKey,
		types.LastUsedEthereumNonceIndex,
		collections.Uint64Key,
		codec.CollValue[zenbtctypes.NonceData](cdc),
	)

	// Set some initial (non-zero) nonce data
	require.NoError(t, ethNonceData.Set(ctx, 1, zenbtctypes.NonceData{Nonce: 10, PrevNonce: 9, Counter: 3, Skip: false}))
	require.NoError(t, ethNonceData.Set(ctx, 2, zenbtctypes.NonceData{Nonce: 42, PrevNonce: 41, Counter: 7, Skip: false}))

	require.NoError(t, v9.ClearEthereumNonceData(ctx, ethNonceData))

	// Verify they were zeroed & Skip set to true
	err := ethNonceData.Walk(ctx, nil, func(key uint64, value zenbtctypes.NonceData) (bool, error) {
		require.Equal(t, uint64(0), value.Nonce, "nonce should be zeroed for key %d", key)
		require.Equal(t, uint64(0), value.PrevNonce, "prev nonce should be zeroed for key %d", key)
		require.Equal(t, uint64(0), value.Counter, "counter should be zeroed for key %d", key)
		require.True(t, value.Skip, "skip flag should be true for key %d", key)
		return false, nil
	})
	require.NoError(t, err)
}

func TestClearEthereumNonceData_MainnetNoOp(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(validation.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey).WithChainID("diamond-1") // mainnet

	kvStoreService := runtime.NewKVStoreService(storeKey)
	sb := collections.NewSchemaBuilder(kvStoreService)

	ethNonceData := collections.NewMap(
		sb,
		types.LastUsedEthereumNonceKey,
		types.LastUsedEthereumNonceIndex,
		collections.Uint64Key,
		codec.CollValue[zenbtctypes.NonceData](cdc),
	)

	original := zenbtctypes.NonceData{Nonce: 5, PrevNonce: 4, Counter: 2, Skip: false}
	require.NoError(t, ethNonceData.Set(ctx, 7, original))

	require.NoError(t, v9.ClearEthereumNonceData(ctx, ethNonceData))

	// Ensure value unchanged
	got, err := ethNonceData.Get(ctx, 7)
	require.NoError(t, err)
	require.Equal(t, original, got, "nonce data should be unchanged on mainnet")
}
