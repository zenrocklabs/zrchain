package v2_test

import (
	"testing"

	"cosmossdk.io/collections"
	validation "github.com/Zenrock-Foundation/zrchain/v6/x/validation/module"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"
	v2 "github.com/zenrocklabs/zenbtc/x/zenbtc/migrations/v2"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

var (
	params = &types.Params{
		DepositKeyringAddr:  "keyring1hpyh7xqr2w7h4eas5y8twnsg",
		StakerKeyID:         1,
		EthMinterKeyID:      2,
		UnstakerKeyID:       3,
		CompleterKeyID:      4,
		RewardsDepositKeyID: 5,
		ChangeAddressKeyIDs: []uint64{6},
		BitcoinProxyAddress: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		EthTokenAddr:        "0xC8CdeDd20cCb4c06884ac4C2fF952A0B7cC230a3",
		ControllerAddr:      "0x5b9Ea8d5486D388a158F026c337DF950866dA5e9",
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
	paramsCol := collections.NewItem(sb, types.ParamsKey, types.ParamsIndex, codec.CollValue[types.Params](cdc))
	err := paramsCol.Set(ctx, types.Params{})
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
	params := collections.NewItem(sb, types.ParamsKey, types.ParamsIndex, codec.CollValue[types.Params](cdc))

	require.NoError(t, v2.UpdateParams(ctx, params))
}

// type mockSubspace struct {
// 	ps types.HVParams
// }

// func newMockSubspace(ps types.HVParams) mockSubspace {
// 	return mockSubspace{ps: ps}
// }

// func (ms mockSubspace) GetParamSet(ctx sdk.Context, ps exported.ParamSet) {
// 	*ps.(*types.Params) = ms.ps
// }
