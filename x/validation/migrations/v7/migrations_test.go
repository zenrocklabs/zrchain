package v7_test

import (
	"testing"

	"cosmossdk.io/collections"
	v6 "github.com/Zenrock-Foundation/zrchain/v5/x/validation/migrations/v6"
	validation "github.com/Zenrock-Foundation/zrchain/v5/x/validation/module"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"

	storetypes "cosmossdk.io/store/types"

	"github.com/Zenrock-Foundation/zrchain/v5/x/validation/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
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
	params := collections.NewItem(sb, types.HVParamsKey, types.HVParamsIndex, codec.CollValue[types.HVParams](cdc))
	err := params.Set(ctx, types.HVParams{})
	require.NoError(t, err)
	require.NoError(t, v6.UpdateParams(ctx, params))

	var res types.HVParams
	bz, err := store.Get(types.HVParamsKey)
	require.NoError(t, err)
	require.NoError(t, cdc.Unmarshal(bz, &res))

	p, err := params.Get(ctx)
	require.NoError(t, err)

	require.Equal(t, p, res)
}
