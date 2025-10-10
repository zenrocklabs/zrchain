package v6_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"github.com/stretchr/testify/require"

	v6 "github.com/zenrocklabs/zenbtc/x/zenbtc/migrations/v6"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

var (
	KVBytes   = collections.NewPrefix(0)
	ItemBytes = collections.NewPrefix(1)
)

// Helper function to create a minimal encoding config
func makeTestEncodingConfig() moduletestutil.TestEncodingConfig {
	return moduletestutil.MakeTestEncodingConfig()
}

func TestMigrateLockTransactions(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)

	oldStoreKey := collections.NewPrefix(100) // Use distinct prefixes for testing
	newStoreKey := collections.NewPrefix(101)

	// Use the real codec
	collCdc := codec.CollValue[types.LockTransaction](cdc)

	oldStore := collections.NewMap(schemaBuilder, oldStoreKey, "old_lock", collections.PairKeyCodec(collections.StringKey, collections.Uint64Key), collCdc)
	newStore := collections.NewMap(schemaBuilder, newStoreKey, "new_lock", collections.StringKey, collCdc)

	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// Sample data
	testData := []struct {
		rawTx string
		vout  uint64
		data  types.LockTransaction
	}{
		{
			rawTx: "tx1",
			vout:  0,
			data:  types.LockTransaction{RawTx: "tx1", Vout: 0, Sender: "sender1", MintRecipient: "recipient1", Amount: 100, BlockHeight: 1},
		},
		{
			rawTx: "tx2",
			vout:  1,
			data:  types.LockTransaction{RawTx: "tx2", Vout: 1, Sender: "sender2", MintRecipient: "recipient2", Amount: 200, BlockHeight: 2},
		},
		{
			rawTx: "tx3withlongstringdata",
			vout:  999,
			data:  types.LockTransaction{RawTx: "tx3withlongstringdata", Vout: 999, Sender: "sender3", MintRecipient: "recipient3", Amount: 300, BlockHeight: 3},
		},
	}

	// Populate old store
	for _, item := range testData {
		oldKey := collections.Join(item.rawTx, item.vout)
		err := oldStore.Set(ctx, oldKey, item.data)
		require.NoError(t, err)
	}

	// Dummy authority setter for testing
	dummyAuthoritySetter := func(string) {}

	// Run the migration
	err = v6.MigrateLockTransactions(ctx, oldStore, newStore, dummyAuthoritySetter)
	require.NoError(t, err)

	// Verify new store
	count := 0
	for _, item := range testData {
		toBeHashed := fmt.Sprintf("%s:%d", item.rawTx, item.vout)
		hash := sha256.Sum256([]byte(toBeHashed))
		newKey := hex.EncodeToString(hash[:])

		migratedData, err := newStore.Get(ctx, newKey)
		require.NoError(t, err, "Expected item with key %s not found in new store", newKey)
		require.Equal(t, item.data, migratedData, "Data mismatch for key %s", newKey)
		count++
	}

	// Verify the count matches
	iter, err := newStore.Iterate(ctx, nil)
	require.NoError(t, err)
	allNewKeys, err := iter.Keys()
	require.NoError(t, err)
	require.Len(t, allNewKeys, len(testData), "Number of items in new store does not match expected")

	// Optional: Verify old store still exists (migration doesn't delete)
	iterOld, err := oldStore.Iterate(ctx, nil)
	require.NoError(t, err)
	allOldKeys, err := iterOld.Keys()
	require.NoError(t, err)
	require.Len(t, allOldKeys, len(testData), "Number of items in old store changed")
}

// Test with empty old store
func TestMigrateLockTransactionsEmpty(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName + "_empty") // Use unique key
	tKey := storetypes.NewTransientStoreKey("transient_test_empty")
	ctx := testutil.DefaultContext(storeKey, tKey)
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)

	oldStoreKey := collections.NewPrefix(102)
	newStoreKey := collections.NewPrefix(103)
	collCdc := codec.CollValue[types.LockTransaction](cdc)

	oldStore := collections.NewMap(schemaBuilder, oldStoreKey, "old_lock_empty", collections.PairKeyCodec(collections.StringKey, collections.Uint64Key), collCdc)
	newStore := collections.NewMap(schemaBuilder, newStoreKey, "new_lock_empty", collections.StringKey, collCdc)

	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// Dummy authority setter for testing
	dummyAuthoritySetter := func(string) {}

	// Run the migration on empty store
	err = v6.MigrateLockTransactions(ctx, oldStore, newStore, dummyAuthoritySetter)
	require.NoError(t, err)

	// Verify new store is empty
	iter, err := newStore.Iterate(ctx, nil)
	require.NoError(t, err)
	keys, err := iter.Keys()
	require.NoError(t, err)
	require.Empty(t, keys, "New store should be empty")
}
