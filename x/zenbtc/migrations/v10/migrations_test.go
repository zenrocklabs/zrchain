package v10_test

import (
	"testing"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"

	v10 "github.com/zenrocklabs/zenbtc/x/zenbtc/migrations/v10"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

func makeTestEncodingConfig() moduletestutil.TestEncodingConfig {
	return moduletestutil.MakeTestEncodingConfig()
}

func TestRemoveInvalidRedemptions(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	amberCtx := testutil.DefaultContext(storeKey, tKey).WithChainID(v10.AmberChainIDPrefix + "-1")
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)
	redemptions := collections.NewMap(schemaBuilder, collections.NewPrefix(0), "redemptions", collections.Uint64Key, codec.CollValue[types.Redemption](cdc))
	firstPendingRedemption := collections.NewItem(schemaBuilder, collections.NewPrefix(1), "first_pending_redemption", collections.Uint64Value)
	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// Sample redemption data with mixed statuses
	testData := []types.Redemption{
		{Data: types.RedemptionData{Id: 1}, Status: types.RedemptionStatus_COMPLETED},
		{Data: types.RedemptionData{Id: 2}, Status: types.RedemptionStatus_UNSTAKED},
		{Data: types.RedemptionData{Id: 3}, Status: types.RedemptionStatus_COMPLETED},
		{Data: types.RedemptionData{Id: 4}, Status: types.RedemptionStatus_AWAITING_SIGN},
		{Data: types.RedemptionData{Id: 5}, Status: types.RedemptionStatus_COMPLETED},
		{Data: types.RedemptionData{Id: 6}, Status: types.RedemptionStatus_UNSTAKED},
		{Data: types.RedemptionData{Id: 7}, Status: types.RedemptionStatus_INITIATED},
	}

	for _, item := range testData {
		err := redemptions.Set(amberCtx, item.Data.Id, item)
		require.NoError(t, err)
	}

	// Set initial first pending
	initialFirstPending := uint64(100)
	err = firstPendingRedemption.Set(amberCtx, initialFirstPending)
	require.NoError(t, err)

	// Run the migration
	err = v10.RemoveInvalidRedemptions(amberCtx, redemptions, firstPendingRedemption)
	require.NoError(t, err)

	// Verification - UNSTAKED and AWAITING_SIGN should be removed
	_, err = redemptions.Get(amberCtx, 2)
	require.Error(t, err, "Redemption 2 (UNSTAKED) should have been removed")
	_, err = redemptions.Get(amberCtx, 4)
	require.Error(t, err, "Redemption 4 (AWAITING_SIGN) should have been removed")
	_, err = redemptions.Get(amberCtx, 6)
	require.Error(t, err, "Redemption 6 (UNSTAKED) should have been removed")

	// COMPLETED and INITIATED redemptions should still exist
	_, err = redemptions.Get(amberCtx, 1)
	require.NoError(t, err, "Redemption 1 (COMPLETED) should still exist")
	_, err = redemptions.Get(amberCtx, 3)
	require.NoError(t, err, "Redemption 3 (COMPLETED) should still exist")
	_, err = redemptions.Get(amberCtx, 5)
	require.NoError(t, err, "Redemption 5 (COMPLETED) should still exist")
	_, err = redemptions.Get(amberCtx, 7)
	require.NoError(t, err, "Redemption 7 (INITIATED) should still exist")

	// Verify count of remaining redemptions
	iter, err := redemptions.Iterate(amberCtx, nil)
	require.NoError(t, err)
	keys, err := iter.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 4, "Should have 4 redemptions remaining")

	// Verify first pending is updated to the latest valid redemption (ID 7, which is INITIATED)
	newFirstPending, err := firstPendingRedemption.Get(amberCtx)
	require.NoError(t, err)
	require.Equal(t, uint64(7), newFirstPending, "First pending should be updated to 7 (latest valid redemption)")
}

func TestRemoveInvalidRedemptions_NoInvalidRedemptions(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	amberCtx := testutil.DefaultContext(storeKey, tKey).WithChainID(v10.AmberChainIDPrefix + "-1")
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)
	redemptions := collections.NewMap(schemaBuilder, collections.NewPrefix(0), "redemptions", collections.Uint64Key, codec.CollValue[types.Redemption](cdc))
	firstPendingRedemption := collections.NewItem(schemaBuilder, collections.NewPrefix(1), "first_pending_redemption", collections.Uint64Value)
	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// Only valid redemptions
	testData := []types.Redemption{
		{Data: types.RedemptionData{Id: 1}, Status: types.RedemptionStatus_COMPLETED},
		{Data: types.RedemptionData{Id: 2}, Status: types.RedemptionStatus_INITIATED},
	}

	for _, item := range testData {
		err := redemptions.Set(amberCtx, item.Data.Id, item)
		require.NoError(t, err)
	}

	initialFirstPending := uint64(1)
	err = firstPendingRedemption.Set(amberCtx, initialFirstPending)
	require.NoError(t, err)

	// Run the migration
	err = v10.RemoveInvalidRedemptions(amberCtx, redemptions, firstPendingRedemption)
	require.NoError(t, err)

	// All redemptions should still exist
	iter, err := redemptions.Iterate(amberCtx, nil)
	require.NoError(t, err)
	keys, err := iter.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 2, "Should have 2 redemptions remaining")

	// First pending should remain unchanged (no invalid redemptions removed)
	newFirstPending, err := firstPendingRedemption.Get(amberCtx)
	require.NoError(t, err)
	require.Equal(t, initialFirstPending, newFirstPending, "First pending should remain unchanged")
}

func TestRemoveInvalidRedemptions_AllInvalid(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	amberCtx := testutil.DefaultContext(storeKey, tKey).WithChainID(v10.AmberChainIDPrefix + "-1")
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)
	redemptions := collections.NewMap(schemaBuilder, collections.NewPrefix(0), "redemptions", collections.Uint64Key, codec.CollValue[types.Redemption](cdc))
	firstPendingRedemption := collections.NewItem(schemaBuilder, collections.NewPrefix(1), "first_pending_redemption", collections.Uint64Value)
	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// All redemptions are UNSTAKED or AWAITING_SIGN
	testData := []types.Redemption{
		{Data: types.RedemptionData{Id: 1}, Status: types.RedemptionStatus_UNSTAKED},
		{Data: types.RedemptionData{Id: 2}, Status: types.RedemptionStatus_AWAITING_SIGN},
	}

	for _, item := range testData {
		err := redemptions.Set(amberCtx, item.Data.Id, item)
		require.NoError(t, err)
	}

	initialFirstPending := uint64(1)
	err = firstPendingRedemption.Set(amberCtx, initialFirstPending)
	require.NoError(t, err)

	// Run the migration
	err = v10.RemoveInvalidRedemptions(amberCtx, redemptions, firstPendingRedemption)
	require.NoError(t, err)

	// All redemptions should be removed
	iter, err := redemptions.Iterate(amberCtx, nil)
	require.NoError(t, err)
	keys, err := iter.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 0, "All redemptions should be removed")

	// First pending should remain unchanged (no valid redemptions to point to)
	newFirstPending, err := firstPendingRedemption.Get(amberCtx)
	require.NoError(t, err)
	require.Equal(t, initialFirstPending, newFirstPending, "First pending should remain unchanged (no valid redemptions)")
}

func TestRemoveInvalidRedemptions_CompletedBeforeInvalid(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	amberCtx := testutil.DefaultContext(storeKey, tKey).WithChainID(v10.AmberChainIDPrefix + "-1")
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)
	redemptions := collections.NewMap(schemaBuilder, collections.NewPrefix(0), "redemptions", collections.Uint64Key, codec.CollValue[types.Redemption](cdc))
	firstPendingRedemption := collections.NewItem(schemaBuilder, collections.NewPrefix(1), "first_pending_redemption", collections.Uint64Value)
	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// COMPLETED redemptions followed by invalid ones
	testData := []types.Redemption{
		{Data: types.RedemptionData{Id: 1}, Status: types.RedemptionStatus_COMPLETED},
		{Data: types.RedemptionData{Id: 2}, Status: types.RedemptionStatus_COMPLETED},
		{Data: types.RedemptionData{Id: 3}, Status: types.RedemptionStatus_COMPLETED},
		{Data: types.RedemptionData{Id: 4}, Status: types.RedemptionStatus_UNSTAKED},
		{Data: types.RedemptionData{Id: 5}, Status: types.RedemptionStatus_AWAITING_SIGN},
	}

	for _, item := range testData {
		err := redemptions.Set(amberCtx, item.Data.Id, item)
		require.NoError(t, err)
	}

	initialFirstPending := uint64(1)
	err = firstPendingRedemption.Set(amberCtx, initialFirstPending)
	require.NoError(t, err)

	// Run the migration
	err = v10.RemoveInvalidRedemptions(amberCtx, redemptions, firstPendingRedemption)
	require.NoError(t, err)

	// Invalid redemptions should be removed
	_, err = redemptions.Get(amberCtx, 4)
	require.Error(t, err, "Redemption 4 (UNSTAKED) should have been removed")
	_, err = redemptions.Get(amberCtx, 5)
	require.Error(t, err, "Redemption 5 (AWAITING_SIGN) should have been removed")

	// COMPLETED redemptions should still exist
	_, err = redemptions.Get(amberCtx, 1)
	require.NoError(t, err, "Redemption 1 (COMPLETED) should still exist")
	_, err = redemptions.Get(amberCtx, 2)
	require.NoError(t, err, "Redemption 2 (COMPLETED) should still exist")
	_, err = redemptions.Get(amberCtx, 3)
	require.NoError(t, err, "Redemption 3 (COMPLETED) should still exist")

	// Verify first pending is updated to the latest valid redemption (ID 3)
	newFirstPending, err := firstPendingRedemption.Get(amberCtx)
	require.NoError(t, err)
	require.Equal(t, uint64(3), newFirstPending, "First pending should be updated to 3 (latest valid redemption)")
}

func TestRemoveInvalidRedemptions_NonAmberChain(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	// Use a non-amber chain ID
	nonAmberCtx := testutil.DefaultContext(storeKey, tKey).WithChainID("gardia-1")
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)
	redemptions := collections.NewMap(schemaBuilder, collections.NewPrefix(0), "redemptions", collections.Uint64Key, codec.CollValue[types.Redemption](cdc))
	firstPendingRedemption := collections.NewItem(schemaBuilder, collections.NewPrefix(1), "first_pending_redemption", collections.Uint64Value)
	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// Sample invalid redemptions
	testData := []types.Redemption{
		{Data: types.RedemptionData{Id: 1}, Status: types.RedemptionStatus_UNSTAKED},
		{Data: types.RedemptionData{Id: 2}, Status: types.RedemptionStatus_AWAITING_SIGN},
	}

	for _, item := range testData {
		err := redemptions.Set(nonAmberCtx, item.Data.Id, item)
		require.NoError(t, err)
	}

	initialFirstPending := uint64(1)
	err = firstPendingRedemption.Set(nonAmberCtx, initialFirstPending)
	require.NoError(t, err)

	// Run the migration on non-amber chain
	err = v10.RemoveInvalidRedemptions(nonAmberCtx, redemptions, firstPendingRedemption)
	require.NoError(t, err)

	// All redemptions should still exist (migration should be skipped)
	iter, err := redemptions.Iterate(nonAmberCtx, nil)
	require.NoError(t, err)
	keys, err := iter.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 2, "No redemptions should be removed on non-amber chain")

	// First pending should remain unchanged
	newFirstPending, err := firstPendingRedemption.Get(nonAmberCtx)
	require.NoError(t, err)
	require.Equal(t, initialFirstPending, newFirstPending, "First pending should remain unchanged on non-amber chain")
}
