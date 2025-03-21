package v3_test

import (
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	v3 "github.com/Zenrock-Foundation/zrchain/v6/x/mint/migrations/v3"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"

	storetypes "cosmossdk.io/store/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/mint"
	"github.com/Zenrock-Foundation/zrchain/v6/x/mint/exported"
	"github.com/Zenrock-Foundation/zrchain/v6/x/mint/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
)

var (
	stakingYield             = math.LegacyNewDecWithPrec(7, 2)
	burnRate                 = math.LegacyNewDecWithPrec(10, 2)
	protocolWalletRate       = math.LegacyNewDecWithPrec(30, 2)
	protocolWalletAddress    = "zen1fhln2vnudxddpymqy82vzqhnlsfh4stjd683ze"
	additionalStakingRewards = math.LegacyNewDecWithPrec(30, 2)
	additionalMpcRewards     = math.LegacyNewDecWithPrec(5, 2)
	additionalBurnRate       = math.LegacyNewDecWithPrec(25, 2)
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
	encCfg := moduletestutil.MakeTestEncodingConfig(mint.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(v3.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)

	kvStoreService := runtime.NewKVStoreService(storeKey)
	store := kvStoreService.OpenKVStore(ctx)
	sb := collections.NewSchemaBuilder(kvStoreService)
	params := collections.NewItem(sb, types.ParamsKey, "minter", codec.CollValue[types.Params](cdc))
	err := params.Set(ctx, types.DefaultParams())
	require.NoError(t, err)
	require.NoError(t, v3.UpdateParams(ctx, params))

	var res types.Params
	bz, err := store.Get(v3.ParamsKey)
	require.NoError(t, err)
	require.NoError(t, cdc.Unmarshal(bz, &res))

	require.True(t, res.AdditionalBurnRate.Equal(additionalBurnRate))
	require.True(t, res.AdditionalMpcRewards.Equal(additionalMpcRewards))
	require.True(t, res.StakingYield.Equal(stakingYield))
	require.True(t, res.BurnRate.Equal(burnRate))
	require.True(t, res.ProtocolWalletRate.Equal(protocolWalletRate))
	require.Equal(t, res.ProtocolWalletAddress, protocolWalletAddress)
	require.True(t, res.AdditionalStakingRewards.Equal(additionalStakingRewards))
}

func TestMigrateFail(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(mint.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(v3.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)

	kvStoreService := runtime.NewKVStoreService(storeKey)
	sb := collections.NewSchemaBuilder(kvStoreService)
	params := collections.NewItem(sb, types.ParamsKey, "minter", codec.CollValue[types.Params](cdc))

	require.Error(t, v3.UpdateParams(ctx, params))
}
