package v7_test

import (
	"testing"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"

	v7 "github.com/zenrocklabs/zenbtc/x/zenbtc/migrations/v7"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

func makeTestEncodingConfig() moduletestutil.TestEncodingConfig {
	return moduletestutil.MakeTestEncodingConfig()
}

func TestPurgeInvalidState(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)
	burnEvents := collections.NewMap(schemaBuilder, collections.NewPrefix(0), "burn_events", collections.Uint64Key, codec.CollValue[types.BurnEvent](cdc))
	firstPendingBurnEvent := collections.NewItem(schemaBuilder, collections.NewPrefix(2), "first_pending_burn_event", collections.Uint64Value)
	redemptions := collections.NewMap(schemaBuilder, collections.NewPrefix(1), "redemptions", collections.Uint64Key, codec.CollValue[types.Redemption](cdc))
	firstPendingRedemption := collections.NewItem(schemaBuilder, collections.NewPrefix(3), "first_pending_redemption", collections.Uint64Value)
	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// Sample data
	testData := []types.BurnEvent{
		{Id: 1, ChainID: v7.SolanaChainIDPrefix + "addr1", DestinationAddr: []byte("addr1"), Amount: 100},
		{Id: 2, ChainID: "another-chain", DestinationAddr: []byte("addr2"), Amount: 200},
		{Id: 3, ChainID: v7.SolanaChainIDPrefix + "addr3", DestinationAddr: []byte("addr3"), Amount: 300},
		{Id: 4, ChainID: "another-chain-2", DestinationAddr: []byte("addr4"), Amount: 400},
	}

	for _, item := range testData {
		err := burnEvents.Set(ctx, item.Id, item)
		require.NoError(t, err)
	}

	initialFirstPending := uint64(100)
	err = firstPendingBurnEvent.Set(ctx, initialFirstPending)
	require.NoError(t, err)

	// Run the migration
	err = v7.PurgeInvalidState(ctx, burnEvents, firstPendingBurnEvent, redemptions, firstPendingRedemption)
	require.NoError(t, err)

	// Verification
	_, err = burnEvents.Get(ctx, 1)
	require.Error(t, err, "Event 1 should have been removed")
	_, err = burnEvents.Get(ctx, 3)
	require.Error(t, err, "Event 3 should have been removed")

	event2, err := burnEvents.Get(ctx, 2)
	require.NoError(t, err, "Event 2 should still exist")
	require.Equal(t, testData[1], event2)

	event4, err := burnEvents.Get(ctx, 4)
	require.NoError(t, err, "Event 4 should still exist")
	require.Equal(t, testData[3], event4)

	iter, err := burnEvents.Iterate(ctx, nil)
	require.NoError(t, err)
	keys, err := iter.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 2)

	newFirstPending, err := firstPendingBurnEvent.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), newFirstPending, "First pending burn event should be updated to the first removed ID")

	// Test no solana events
	// Clear solana events
	err = burnEvents.Remove(ctx, uint64(2))
	require.NoError(t, err)
	err = burnEvents.Remove(ctx, uint64(4))
	require.NoError(t, err)

	// re-add them without solana prefix
	testData2 := []types.BurnEvent{
		{Id: 5, ChainID: "another-chain-3", DestinationAddr: []byte("addr5"), Amount: 500},
		{Id: 6, ChainID: "another-chain-4", DestinationAddr: []byte("addr6"), Amount: 600},
	}
	for _, item := range testData2 {
		err := burnEvents.Set(ctx, item.Id, item)
		require.NoError(t, err)
	}

	currentFirstPending, err := firstPendingBurnEvent.Get(ctx)
	require.NoError(t, err)

	err = v7.PurgeInvalidState(ctx, burnEvents, firstPendingBurnEvent, redemptions, firstPendingRedemption)
	require.NoError(t, err)
	newFirstPendingAfterNoop, err := firstPendingBurnEvent.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, currentFirstPending, newFirstPendingAfterNoop, "First pending burn event should not change if no events are removed")

}

func TestPurgeInvalidState_Redemptions(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	gardiaCtx := testutil.DefaultContext(storeKey, tKey).WithChainID(v7.GardiaChainIDPrefix + "-1")
	nonGardiaCtx := testutil.DefaultContext(storeKey, tKey).WithChainID("some-other-chain")
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)

	burnEvents := collections.NewMap(schemaBuilder, collections.NewPrefix(0), "burn_events", collections.Uint64Key, codec.CollValue[types.BurnEvent](cdc))
	firstPendingBurnEvent := collections.NewItem(schemaBuilder, collections.NewPrefix(2), "first_pending_burn_event", collections.Uint64Value)
	redemptions := collections.NewMap(schemaBuilder, collections.NewPrefix(1), "redemptions", collections.Uint64Key, codec.CollValue[types.Redemption](cdc))
	firstPendingRedemption := collections.NewItem(schemaBuilder, collections.NewPrefix(3), "first_pending_redemption", collections.Uint64Value)

	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// Sample redemption data
	redemptionData := []types.Redemption{
		{Data: types.RedemptionData{Id: 1}, Status: types.RedemptionStatus_INITIATED},
		{Data: types.RedemptionData{Id: 2}, Status: types.RedemptionStatus_UNSTAKED},
		{Data: types.RedemptionData{Id: 3}, Status: types.RedemptionStatus_COMPLETED},
		{Data: types.RedemptionData{Id: 4}, Status: types.RedemptionStatus_UNSTAKED},
	}

	for _, item := range redemptionData {
		err := redemptions.Set(gardiaCtx, item.Data.Id, item)
		require.NoError(t, err)
	}

	initialFirstPendingRedemption := uint64(999)
	err = firstPendingRedemption.Set(gardiaCtx, initialFirstPendingRedemption)
	require.NoError(t, err)

	// Run migration on a non-gardia chain first
	err = v7.PurgeInvalidState(nonGardiaCtx, burnEvents, firstPendingBurnEvent, redemptions, firstPendingRedemption)
	require.NoError(t, err)

	// Verify no redemptions were removed
	iter, err := redemptions.Iterate(gardiaCtx, nil)
	require.NoError(t, err)
	keys, err := iter.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 4, "No redemptions should be removed on a non-gardia chain")

	firstPendingRedemptionAfterNoop, err := firstPendingRedemption.Get(gardiaCtx)
	require.NoError(t, err)
	require.Equal(t, initialFirstPendingRedemption, firstPendingRedemptionAfterNoop, "First pending redemption should not change on a non-gardia chain")

	// Run migration on a gardia chain
	err = v7.PurgeInvalidState(gardiaCtx, burnEvents, firstPendingBurnEvent, redemptions, firstPendingRedemption)
	require.NoError(t, err)

	// Verification
	_, err = redemptions.Get(gardiaCtx, 2)
	require.Error(t, err, "Redemption 2 should have been removed")
	_, err = redemptions.Get(gardiaCtx, 4)
	require.Error(t, err, "Redemption 4 should have been removed")

	_, err = redemptions.Get(gardiaCtx, 1)
	require.NoError(t, err, "Redemption 1 should still exist")
	_, err = redemptions.Get(gardiaCtx, 3)
	require.NoError(t, err, "Redemption 3 should still exist")

	iter, err = redemptions.Iterate(gardiaCtx, nil)
	require.NoError(t, err)
	keys, err = iter.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 2, "There should be 2 redemptions remaining")

	newFirstPendingRedemption, err := firstPendingRedemption.Get(gardiaCtx)
	require.NoError(t, err)
	require.Equal(t, uint64(2), newFirstPendingRedemption, "First pending redemption should be updated to the first removed ID")
}
