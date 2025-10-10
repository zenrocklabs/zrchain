package v9_test

import (
	"testing"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"

	v9 "github.com/zenrocklabs/zenbtc/x/zenbtc/migrations/v9"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

func makeTestEncodingConfig() moduletestutil.TestEncodingConfig {
	return moduletestutil.MakeTestEncodingConfig()
}

func TestRemoveBurnedEvents(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	amberCtx := testutil.DefaultContext(storeKey, tKey).WithChainID(v9.AmberChainIDPrefix + "-1")
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)
	burnEvents := collections.NewMap(schemaBuilder, collections.NewPrefix(0), "burn_events", collections.Uint64Key, codec.CollValue[types.BurnEvent](cdc))
	firstPendingBurnEvent := collections.NewItem(schemaBuilder, collections.NewPrefix(1), "first_pending_burn_event", collections.Uint64Value)
	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// Sample burn event data with mixed statuses
	testData := []types.BurnEvent{
		{Id: 1, ChainID: "chain-1", Status: types.BurnStatus_BURN_STATUS_UNSTAKING, Amount: 100},
		{Id: 2, ChainID: "chain-1", Status: types.BurnStatus_BURN_STATUS_BURNED, Amount: 200},
		{Id: 3, ChainID: "chain-1", Status: types.BurnStatus_BURN_STATUS_BURNED, Amount: 300},
		{Id: 4, ChainID: "chain-1", Status: types.BurnStatus_BURN_STATUS_UNSTAKING, Amount: 400},
		{Id: 5, ChainID: "chain-1", Status: types.BurnStatus_BURN_STATUS_BURNED, Amount: 500},
	}

	for _, item := range testData {
		err := burnEvents.Set(amberCtx, item.Id, item)
		require.NoError(t, err)
	}

	// Set first pending to point to a BURNED event
	initialFirstPending := uint64(2)
	err = firstPendingBurnEvent.Set(amberCtx, initialFirstPending)
	require.NoError(t, err)

	// Run the migration
	err = v9.RemoveBurnedEvents(amberCtx, burnEvents, firstPendingBurnEvent)
	require.NoError(t, err)

	// Verification - BURNED events should be removed
	_, err = burnEvents.Get(amberCtx, 2)
	require.Error(t, err, "Event 2 (BURNED) should have been removed")
	_, err = burnEvents.Get(amberCtx, 3)
	require.Error(t, err, "Event 3 (BURNED) should have been removed")
	_, err = burnEvents.Get(amberCtx, 5)
	require.Error(t, err, "Event 5 (BURNED) should have been removed")

	// Non-BURNED events should still exist
	event1, err := burnEvents.Get(amberCtx, 1)
	require.NoError(t, err, "Event 1 should still exist")
	require.Equal(t, testData[0], event1)

	event4, err := burnEvents.Get(amberCtx, 4)
	require.NoError(t, err, "Event 4 should still exist")
	require.Equal(t, testData[3], event4)

	// Verify count of remaining events
	iter, err := burnEvents.Iterate(amberCtx, nil)
	require.NoError(t, err)
	keys, err := iter.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 2, "Should have 2 events remaining")

	// Verify first pending is updated to the last non-burned event before the deleted ones
	newFirstPending, err := firstPendingBurnEvent.Get(amberCtx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), newFirstPending, "First pending should be updated to 1 (last non-burned event before first burned)")
}

func TestRemoveBurnedEvents_PointerAfterBurnedEvents(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	amberCtx := testutil.DefaultContext(storeKey, tKey).WithChainID(v9.AmberChainIDPrefix + "-1")
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)
	burnEvents := collections.NewMap(schemaBuilder, collections.NewPrefix(0), "burn_events", collections.Uint64Key, codec.CollValue[types.BurnEvent](cdc))
	firstPendingBurnEvent := collections.NewItem(schemaBuilder, collections.NewPrefix(1), "first_pending_burn_event", collections.Uint64Value)
	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// Events with BURNED at the beginning
	testData := []types.BurnEvent{
		{Id: 1, ChainID: "chain-1", Status: types.BurnStatus_BURN_STATUS_BURNED, Amount: 100},
		{Id: 2, ChainID: "chain-1", Status: types.BurnStatus_BURN_STATUS_BURNED, Amount: 200},
		{Id: 3, ChainID: "chain-1", Status: types.BurnStatus_BURN_STATUS_UNSTAKING, Amount: 300},
		{Id: 4, ChainID: "chain-1", Status: types.BurnStatus_BURN_STATUS_UNSTAKING, Amount: 400},
	}

	for _, item := range testData {
		err := burnEvents.Set(amberCtx, item.Id, item)
		require.NoError(t, err)
	}

	// Set first pending to point after all BURNED events
	initialFirstPending := uint64(3)
	err = firstPendingBurnEvent.Set(amberCtx, initialFirstPending)
	require.NoError(t, err)

	// Run the migration
	err = v9.RemoveBurnedEvents(amberCtx, burnEvents, firstPendingBurnEvent)
	require.NoError(t, err)

	// Verify BURNED events removed
	_, err = burnEvents.Get(amberCtx, 1)
	require.Error(t, err, "Event 1 (BURNED) should have been removed")
	_, err = burnEvents.Get(amberCtx, 2)
	require.Error(t, err, "Event 2 (BURNED) should have been removed")

	// Non-BURNED events should still exist
	_, err = burnEvents.Get(amberCtx, 3)
	require.NoError(t, err, "Event 3 should still exist")
	_, err = burnEvents.Get(amberCtx, 4)
	require.NoError(t, err, "Event 4 should still exist")

	// First pending should be unchanged because there were no non-burned events before the first burned one
	newFirstPending, err := firstPendingBurnEvent.Get(amberCtx)
	require.NoError(t, err)
	require.Equal(t, initialFirstPending, newFirstPending, "First pending should remain unchanged (no non-burned events before burned ones)")
}

func TestRemoveBurnedEvents_NoBurnedEvents(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	amberCtx := testutil.DefaultContext(storeKey, tKey).WithChainID(v9.AmberChainIDPrefix + "-1")
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)
	burnEvents := collections.NewMap(schemaBuilder, collections.NewPrefix(0), "burn_events", collections.Uint64Key, codec.CollValue[types.BurnEvent](cdc))
	firstPendingBurnEvent := collections.NewItem(schemaBuilder, collections.NewPrefix(1), "first_pending_burn_event", collections.Uint64Value)
	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// No BURNED events
	testData := []types.BurnEvent{
		{Id: 1, ChainID: "chain-1", Status: types.BurnStatus_BURN_STATUS_UNSTAKING, Amount: 100},
		{Id: 2, ChainID: "chain-1", Status: types.BurnStatus_BURN_STATUS_UNSTAKING, Amount: 200},
	}

	for _, item := range testData {
		err := burnEvents.Set(amberCtx, item.Id, item)
		require.NoError(t, err)
	}

	initialFirstPending := uint64(1)
	err = firstPendingBurnEvent.Set(amberCtx, initialFirstPending)
	require.NoError(t, err)

	// Run the migration
	err = v9.RemoveBurnedEvents(amberCtx, burnEvents, firstPendingBurnEvent)
	require.NoError(t, err)

	// All events should still exist
	iter, err := burnEvents.Iterate(amberCtx, nil)
	require.NoError(t, err)
	keys, err := iter.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 2, "Should have 2 events remaining")

	// First pending should remain unchanged
	newFirstPending, err := firstPendingBurnEvent.Get(amberCtx)
	require.NoError(t, err)
	require.Equal(t, initialFirstPending, newFirstPending, "First pending should remain unchanged")
}

func TestRemoveBurnedEvents_AllBurned(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	amberCtx := testutil.DefaultContext(storeKey, tKey).WithChainID(v9.AmberChainIDPrefix + "-1")
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)
	burnEvents := collections.NewMap(schemaBuilder, collections.NewPrefix(0), "burn_events", collections.Uint64Key, codec.CollValue[types.BurnEvent](cdc))
	firstPendingBurnEvent := collections.NewItem(schemaBuilder, collections.NewPrefix(1), "first_pending_burn_event", collections.Uint64Value)
	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// All events are BURNED
	testData := []types.BurnEvent{
		{Id: 1, ChainID: "chain-1", Status: types.BurnStatus_BURN_STATUS_BURNED, Amount: 100},
		{Id: 2, ChainID: "chain-1", Status: types.BurnStatus_BURN_STATUS_BURNED, Amount: 200},
	}

	for _, item := range testData {
		err := burnEvents.Set(amberCtx, item.Id, item)
		require.NoError(t, err)
	}

	initialFirstPending := uint64(1)
	err = firstPendingBurnEvent.Set(amberCtx, initialFirstPending)
	require.NoError(t, err)

	// Run the migration
	err = v9.RemoveBurnedEvents(amberCtx, burnEvents, firstPendingBurnEvent)
	require.NoError(t, err)

	// All events should be removed
	iter, err := burnEvents.Iterate(amberCtx, nil)
	require.NoError(t, err)
	keys, err := iter.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 0, "All events should be removed")

	// Note: First pending pointer behavior when all events are removed
	// The pointer will be updated if there are no non-burned events
	// In this case, the pointer remains as set since there's no minNonBurnedID
}

func TestRemoveBurnedEvents_NonAmberChain(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	// Use a non-amber chain ID
	nonAmberCtx := testutil.DefaultContext(storeKey, tKey).WithChainID("gardia-1")
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)
	burnEvents := collections.NewMap(schemaBuilder, collections.NewPrefix(0), "burn_events", collections.Uint64Key, codec.CollValue[types.BurnEvent](cdc))
	firstPendingBurnEvent := collections.NewItem(schemaBuilder, collections.NewPrefix(1), "first_pending_burn_event", collections.Uint64Value)
	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// Sample BURNED events
	testData := []types.BurnEvent{
		{Id: 1, ChainID: "chain-1", Status: types.BurnStatus_BURN_STATUS_BURNED, Amount: 100},
		{Id: 2, ChainID: "chain-1", Status: types.BurnStatus_BURN_STATUS_BURNED, Amount: 200},
	}

	for _, item := range testData {
		err := burnEvents.Set(nonAmberCtx, item.Id, item)
		require.NoError(t, err)
	}

	initialFirstPending := uint64(1)
	err = firstPendingBurnEvent.Set(nonAmberCtx, initialFirstPending)
	require.NoError(t, err)

	// Run the migration on non-amber chain
	err = v9.RemoveBurnedEvents(nonAmberCtx, burnEvents, firstPendingBurnEvent)
	require.NoError(t, err)

	// All events should still exist (migration should be skipped)
	iter, err := burnEvents.Iterate(nonAmberCtx, nil)
	require.NoError(t, err)
	keys, err := iter.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 2, "No events should be removed on non-amber chain")

	// First pending should remain unchanged
	newFirstPending, err := firstPendingBurnEvent.Get(nonAmberCtx)
	require.NoError(t, err)
	require.Equal(t, initialFirstPending, newFirstPending, "First pending should remain unchanged on non-amber chain")
}
