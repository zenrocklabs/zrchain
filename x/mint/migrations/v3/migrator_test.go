package v3_test

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v5/x/mint/exported"
	"github.com/Zenrock-Foundation/zrchain/v5/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// "testing"
	// "github.com/stretchr/testify/require"
	// "cosmossdk.io/collections"
	// "cosmossdk.io/math"
	// storetypes "cosmossdk.io/store/types"
	// "github.com/Zenrock-Foundation/zrchain/v5/x/mint"
	// "github.com/Zenrock-Foundation/zrchain/v5/x/mint/exported"
	// v3 "github.com/Zenrock-Foundation/zrchain/v5/x/mint/migrations/v3"
	// "github.com/Zenrock-Foundation/zrchain/v5/x/mint/types"
	// "github.com/cosmos/cosmos-sdk/codec"
	// "github.com/cosmos/cosmos-sdk/runtime"
	// "github.com/cosmos/cosmos-sdk/testutil"
	// sdk "github.com/cosmos/cosmos-sdk/types"
	// moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
)

type mockSubspace struct {
	ps types.Params
}

func newMockSubspace(ps types.Params) mockSubspace {
	return mockSubspace{ps: ps}
}

func (ms mockSubspace) GetParamSet(ctx sdk.Context, ps exported.ParamSet) {
	*ps.(*types.Params) = ms.ps
}

func TestMigrate(t *testing.T) {
	t.SkipNow()
	// encCfg := moduletestutil.MakeTestEncodingConfig(mint.AppModuleBasic{})
	// cdc := encCfg.Codec

	// storeKey := storetypes.NewKVStoreKey(v3.ModuleName)
	// tKey := storetypes.NewTransientStoreKey("transient_test")
	// ctx := testutil.DefaultContext(storeKey, tKey)
	// kvStoreService := runtime.NewKVStoreService(storeKey)
	// store := kvStoreService.OpenKVStore(ctx)

	// // Initialize the parameters with expected values
	// params := types.Params{
	// 	MintDenom:           "urock",
	// 	InflationRateChange: math.LegacyNewDecWithPrec(0, 2),  // 0.13
	// 	InflationMax:        math.LegacyNewDecWithPrec(0, 2),  // 0.20
	// 	InflationMin:        math.LegacyNewDecWithPrec(0, 2),  // 0.07
	// 	GoalBonded:          math.LegacyNewDecWithPrec(67, 2), // 0.67
	// 	BlocksPerYear:       6311520,                          // for one year
	// }

	// // Create a schema builder and prefix for the item
	// schemaBuilder := collections.NewSchemaBuilder(kvStoreService)
	// prefix := collections.NewPrefix(v3.ParamsKey)

	// // Set the parameters in the store
	// err := store.Set(v3.ParamsKey, cdc.MustMarshal(&params)) // Ensure the key is set
	// require.NoError(t, err)

	// require.NoError(t, v3.UpdateParams(ctx, collections.NewItem(schemaBuilder, prefix, "params", codec.CollValue[types.Params](cdc))))

	// var res types.Params
	// bz, err := store.Get(v3.ParamsKey)
	// require.NoError(t, err)
	// require.NoError(t, cdc.Unmarshal(bz, &res))
	// require.Equal(t, params, res)
}
