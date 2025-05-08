package v2_test

import (
	"testing"

	"cosmossdk.io/collections"
	validation "github.com/Zenrock-Foundation/zrchain/v6/x/validation/module"
	v2 "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/migrations/v2"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"

	storetypes "cosmossdk.io/store/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
)

var (
	params = &types.Solana{
		SignerKeyId:       10,
		ProgramId:         "DXREJumiQhNejXa1b5EFPUxtSYdyJXBdiHeu6uX1ribA",
		NonceAccountKey:   12,
		NonceAuthorityKey: 11,
		MintAddress:       "StVNdHNSFK3uVTL5apWHysgze4M8zrsqwjEAH1JM87i",
		FeeWallet:         "FzqGcRG98v1KhKxatX2Abb2z1aJ2rViQwBK5GHByKCAd",
		Fee:               0,
		Btl:               20,
	}
)

func TestMigrate(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(validation.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)

	kvStoreService := runtime.NewKVStoreService(storeKey)
	store := kvStoreService.OpenKVStore(ctx)
	sb := collections.NewSchemaBuilder(kvStoreService)
	paramsCol := collections.NewItem(sb, types.ParamsKey, types.ParamsIndex, codec.CollValue[types.Solana](cdc))
	err := paramsCol.Set(ctx, types.Solana{})
	require.NoError(t, err)
	require.NoError(t, v2.UpdateParams(ctx, paramsCol))

	// Get the value from the params collection
	expectedParams, err := paramsCol.Get(ctx)
	require.NoError(t, err)

	var res types.Params
	bz, err := store.Get(types.ParamsKey)
	require.NoError(t, err)
	require.NoError(t, cdc.Unmarshal(bz, &res))

	require.Equal(t, expectedParams, res)
}

func TestMigrateFail(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(validation.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)

	kvStoreService := runtime.NewKVStoreService(storeKey)
	sb := collections.NewSchemaBuilder(kvStoreService)
	params := collections.NewItem(sb, types.ParamsKey, types.ParamsIndex, codec.CollValue[types.Solana](cdc))

	require.NoError(t, v2.UpdateParams(ctx, params))
}
