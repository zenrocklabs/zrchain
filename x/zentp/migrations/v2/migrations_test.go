package v2_test

import (
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
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
	params = &types.Params{
		Solana: &types.Solana{
			SignerKeyId:       10,
			ProgramId:         "AgoRvPWg2R7nkKhxvipvms79FmxQr75r2GwNSpPtxcLg",
			NonceAccountKey:   12,
			NonceAuthorityKey: 11,
			MintAddress:       "4oUDGAy46CmemmozTt6kWT5E3rqkLp2rCvAumpMWqR5T",
			FeeWallet:         "5aLz81F9uugwKBmvUY3DcXB1B7G2Yf7tB9zacdJBhZbh",
			Fee:               0,
			Btl:               20,
		},
		BridgeFee: math.LegacyNewDecWithPrec(1, 2),
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
