package v3_test

import (
	"testing"

	"cosmossdk.io/collections"
	validation "github.com/Zenrock-Foundation/zrchain/v6/x/validation/module"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"

	storetypes "cosmossdk.io/store/types"

	v3 "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/migrations/v3"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
)

var (
	mint1 = &types.Bridge{
		Id:                 1,
		Denom:              "urock",
		Creator:            "test-creator",
		SourceAddress:      "test-creator",
		SourceChain:        "test-source-chain",
		DestinationChain:   "test-destination-chain",
		Amount:             1000,
		RecipientAddress:   "test-recipient-address",
		TxId:               123,
		TxHash:             "test-tx-hash",
		State:              types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		BlockHeight:        100,
		AwaitingEventSince: 0,
	}

	mint2 = &types.Bridge{
		Id:                 2,
		Denom:              "urock",
		Creator:            "test-creator",
		SourceAddress:      "test-creator",
		SourceChain:        "test-source-chain",
		DestinationChain:   "test-destination-chain",
		Amount:             2000,
		RecipientAddress:   "test-recipient-address",
		TxId:               456,
		TxHash:             "test-tx-hash-2",
		State:              types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		BlockHeight:        100,
		AwaitingEventSince: 0,
	}

	burn1 = &types.Bridge{
		Id:                 1,
		Denom:              "urock",
		Creator:            "test-creator",
		SourceAddress:      "test-creator",
		SourceChain:        "test-source-chain",
		DestinationChain:   "test-destination-chain",
		Amount:             1000,
		RecipientAddress:   "test-recipient-address",
		TxId:               123,
		TxHash:             "test-tx-hash",
		State:              types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		BlockHeight:        100,
		AwaitingEventSince: 0,
	}

	burn2 = &types.Bridge{
		Id:                 2,
		Denom:              "urock",
		Creator:            "test-creator",
		SourceAddress:      "test-creator",
		SourceChain:        "test-source-chain",
		DestinationChain:   "test-destination-chain",
		Amount:             1000,
		RecipientAddress:   "test-recipient-address",
		TxId:               123,
		TxHash:             "test-tx-hash",
		State:              types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		BlockHeight:        100,
		AwaitingEventSince: 0,
	}
)

func TestUpdateMintStore(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(validation.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)

	kvStoreService := runtime.NewKVStoreService(storeKey)
	sb := collections.NewSchemaBuilder(kvStoreService)
	oldMintsCol := collections.NewMap(sb, types.MintsKey, types.MintsIndex, collections.Uint64Key, codec.CollValue[types.Bridge](cdc))
	err := oldMintsCol.Set(ctx, 1, *mint1)
	require.NoError(t, err)
	err = oldMintsCol.Set(ctx, 2, *mint2)
	require.NoError(t, err)
	newMintsCol := collections.NewMap(sb, types.MintsKey, types.MintsIndex, collections.Uint64Key, codec.CollValue[types.Bridge](cdc))

	require.NoError(t, err)
	require.NoError(t, v3.UpdateMintStore(ctx, oldMintsCol, newMintsCol))

	mint, err := newMintsCol.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, mint1.Id, mint.Id)

	mintStore, err := newMintsCol.Iterate(ctx, nil)
	require.NoError(t, err)
	mints, err := mintStore.Values()
	require.NoError(t, err)
	require.Equal(t, 2, len(mints))

	mint, err = newMintsCol.Get(ctx, 2)
	require.NoError(t, err)
	require.Equal(t, mint2.Id, mint.Id)

	require.Equal(t, 2, len(mints))
}

func TestUpdateBurnStore(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(validation.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)

	kvStoreService := runtime.NewKVStoreService(storeKey)
	sb := collections.NewSchemaBuilder(kvStoreService)
	oldBurnsCol := collections.NewMap(sb, types.BurnsKey, types.BurnsIndex, collections.Uint64Key, codec.CollValue[types.Bridge](cdc))
	err := oldBurnsCol.Set(ctx, 1, *burn1)
	newBurnsCol := collections.NewMap(sb, types.BurnsKey, types.BurnsIndex, collections.Uint64Key, codec.CollValue[types.Bridge](cdc))

	require.NoError(t, err)
	require.NoError(t, v3.UpdateBurnStore(ctx, oldBurnsCol, newBurnsCol))

	burn, err := newBurnsCol.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, burn1.Id, burn.Id)
}
