package v8_test

import (
	"testing"

	"cosmossdk.io/collections"
	v8 "github.com/Zenrock-Foundation/zrchain/v6/x/validation/migrations/v8"
	validation "github.com/Zenrock-Foundation/zrchain/v6/x/validation/module"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"

	storetypes "cosmossdk.io/store/types"

	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
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
	sb := collections.NewSchemaBuilder(kvStoreService)
	btcBlockHeaders := collections.NewMap(sb, types.BtcBlockHeadersKey, types.BtcBlockHeadersIndex, collections.Int64Key, codec.CollValue[api.BTCBlockHeader](cdc))
	err := btcBlockHeaders.Set(ctx, 1, api.BTCBlockHeader{
		Version:    1,
		PrevBlock:  "0000000000000000000000000000000000000000000000000000000000000000",
		MerkleRoot: "0000000000000000000000000000000000000000000000000000000000000000",
		TimeStamp:  1,
		Bits:       1,
		Nonce:      1,
		BlockHash:  "0000000000000000000000000000000000000000000000000000000000000001",
	})
	require.NoError(t, err)
	err = btcBlockHeaders.Set(ctx, 2, api.BTCBlockHeader{
		Version:    2,
		PrevBlock:  "0000000000000000000000000000000000000000000000000000000000000001",
		MerkleRoot: "0000000000000000000000000000000000000000000000000000000000000001",
		TimeStamp:  2,
		Bits:       2,
		Nonce:      2,
		BlockHash:  "0000000000000000000000000000000000000000000000000000000000000002",
	})
	require.NoError(t, err)
	require.NoError(t, v8.UpdateBtcBlockHeaders(ctx, btcBlockHeaders))

	btcBlockHeaders.Walk(ctx, nil, func(key int64, value api.BTCBlockHeader) (stop bool, err error) {
		require.Equal(t, key, value.BlockHeight)
		return false, nil
	})
}
