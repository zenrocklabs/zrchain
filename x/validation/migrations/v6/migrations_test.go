package v6_test

import (
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
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

var (
	ZenBTCParams = &types.ZenBTCParams{
		ZenBTCEthBatcherAddr:      "0x912D79F8d489d0d007aBE0E26fD5d2f06BA4A2AA",
		ZenBTCDepositKeyringAddr:  "keyring1hpyh7xqr2w7h4eas5y8twnsg",
		ZenBTCWithdrawerKeyID:     1,
		ZenBTCMinterKeyID:         2,
		ZenBTCChangeAddressKeyIDs: []uint64{3},
		ZenBTCUnstakerKeyID:       4,
		ZenBTCRewardsDepositKeyID: 5,
		BitcoinProxyCreatorID:     "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		StakeableAssets: []*types.AssetData{
			{Asset: types.Asset_ROCK, Precision: 6, PriceUSD: math.LegacyZeroDec()},
			{Asset: types.Asset_zenBTC, Precision: 8, PriceUSD: math.LegacyZeroDec()},
			{Asset: types.Asset_stETH, Precision: 18, PriceUSD: math.LegacyZeroDec()},
		},
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
	params := collections.NewItem(sb, types.HVParamsKey, types.HVParamsIndex, codec.CollValue[types.HVParams](cdc))
	err := params.Set(ctx, types.HVParams{})
	require.NoError(t, err)
	require.NoError(t, v6.UpdateParams(ctx, params))

	var res types.HVParams
	bz, err := store.Get(types.HVParamsKey)
	require.NoError(t, err)
	require.NoError(t, cdc.Unmarshal(bz, &res))

	require.Equal(t, ZenBTCParams, res.ZenBTCParams)
}

func TestMigrateFail(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(validation.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)

	kvStoreService := runtime.NewKVStoreService(storeKey)
	sb := collections.NewSchemaBuilder(kvStoreService)
	params := collections.NewItem(sb, types.HVParamsKey, types.HVParamsIndex, codec.CollValue[types.HVParams](cdc))

	require.Error(t, v6.UpdateParams(ctx, params))
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
