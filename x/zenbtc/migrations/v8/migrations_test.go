package v8_test

import (
	"testing"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"

	v8 "github.com/zenrocklabs/zenbtc/x/zenbtc/migrations/v8"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

func makeTestEncodingConfig() moduletestutil.TestEncodingConfig {
	return moduletestutil.MakeTestEncodingConfig()
}

func TestRemoveStakedMints(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	amberCtx := testutil.DefaultContext(storeKey, tKey).WithChainID(v8.AmberChainIDPrefix + "-1")
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)
	pendingMintTransactionsMap := collections.NewMap(schemaBuilder, collections.NewPrefix(0), "pending_mint_transactions", collections.Uint64Key, codec.CollValue[types.PendingMintTransaction](cdc))
	firstPendingSolMintTransaction := collections.NewItem(schemaBuilder, collections.NewPrefix(1), "first_pending_sol_mint", collections.Uint64Value)
	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// Sample pending mint transaction data - all Solana
	testData := []types.PendingMintTransaction{
		{Id: 1, ChainType: types.WalletType_WALLET_TYPE_SOLANA, Status: types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED, Amount: 100},
		{Id: 2, ChainType: types.WalletType_WALLET_TYPE_SOLANA, Status: types.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED, Amount: 200},
		{Id: 3, ChainType: types.WalletType_WALLET_TYPE_SOLANA, Status: types.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED, Amount: 300},
		{Id: 4, ChainType: types.WalletType_WALLET_TYPE_SOLANA, Status: types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED, Amount: 400},
		{Id: 5, ChainType: types.WalletType_WALLET_TYPE_EVM, Status: types.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED, Amount: 500}, // EVM should be ignored
	}

	for _, item := range testData {
		err := pendingMintTransactionsMap.Set(amberCtx, item.Id, item)
		require.NoError(t, err)
	}

	// Set initial first pending value
	initialFirstPendingSol := uint64(200)
	err = firstPendingSolMintTransaction.Set(amberCtx, initialFirstPendingSol)
	require.NoError(t, err)

	// Run the migration
	err = v8.RemoveStakedMints(amberCtx, pendingMintTransactionsMap, firstPendingSolMintTransaction)
	require.NoError(t, err)

	// Verification - STAKED Solana transactions should be removed
	_, err = pendingMintTransactionsMap.Get(amberCtx, 2)
	require.Error(t, err, "Transaction 2 (Solana STAKED) should have been removed")
	_, err = pendingMintTransactionsMap.Get(amberCtx, 3)
	require.Error(t, err, "Transaction 3 (Solana STAKED) should have been removed")

	// Non-STAKED Solana transactions should still exist
	tx1, err := pendingMintTransactionsMap.Get(amberCtx, 1)
	require.NoError(t, err, "Transaction 1 should still exist")
	require.Equal(t, testData[0], tx1)

	tx4, err := pendingMintTransactionsMap.Get(amberCtx, 4)
	require.NoError(t, err, "Transaction 4 should still exist")
	require.Equal(t, testData[3], tx4)

	// EVM STAKED transaction should still exist (not removed)
	tx5, err := pendingMintTransactionsMap.Get(amberCtx, 5)
	require.NoError(t, err, "Transaction 5 (EVM) should still exist")
	require.Equal(t, testData[4], tx5)

	// Verify count of remaining transactions
	iter, err := pendingMintTransactionsMap.Iterate(amberCtx, nil)
	require.NoError(t, err)
	keys, err := iter.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 3, "Should have 3 transactions remaining")

	// Verify first pending Solana index is updated correctly
	newFirstPendingSol, err := firstPendingSolMintTransaction.Get(amberCtx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), newFirstPendingSol, "First pending Sol transaction should be 1")
}

func TestRemoveStakedMints_NoStakedTransactions(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	amberCtx := testutil.DefaultContext(storeKey, tKey).WithChainID(v8.AmberChainIDPrefix + "-1")
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)
	pendingMintTransactionsMap := collections.NewMap(schemaBuilder, collections.NewPrefix(0), "pending_mint_transactions", collections.Uint64Key, codec.CollValue[types.PendingMintTransaction](cdc))
	firstPendingSolMintTransaction := collections.NewItem(schemaBuilder, collections.NewPrefix(1), "first_pending_sol_mint", collections.Uint64Value)
	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// Sample data with no STAKED Solana transactions
	testData := []types.PendingMintTransaction{
		{Id: 1, ChainType: types.WalletType_WALLET_TYPE_SOLANA, Status: types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED, Amount: 100},
		{Id: 2, ChainType: types.WalletType_WALLET_TYPE_SOLANA, Status: types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED, Amount: 200},
	}

	for _, item := range testData {
		err := pendingMintTransactionsMap.Set(amberCtx, item.Id, item)
		require.NoError(t, err)
	}

	initialFirstPendingSol := uint64(1)
	err = firstPendingSolMintTransaction.Set(amberCtx, initialFirstPendingSol)
	require.NoError(t, err)

	// Run the migration
	err = v8.RemoveStakedMints(amberCtx, pendingMintTransactionsMap, firstPendingSolMintTransaction)
	require.NoError(t, err)

	// All transactions should still exist
	iter, err := pendingMintTransactionsMap.Iterate(amberCtx, nil)
	require.NoError(t, err)
	keys, err := iter.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 2, "Should have 2 transactions remaining")

	// First pending index should be updated to the minimum ID
	newFirstPendingSol, err := firstPendingSolMintTransaction.Get(amberCtx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), newFirstPendingSol, "First pending Sol should be 1")
}

func TestRemoveStakedMints_AllStaked(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	amberCtx := testutil.DefaultContext(storeKey, tKey).WithChainID(v8.AmberChainIDPrefix + "-1")
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)
	pendingMintTransactionsMap := collections.NewMap(schemaBuilder, collections.NewPrefix(0), "pending_mint_transactions", collections.Uint64Key, codec.CollValue[types.PendingMintTransaction](cdc))
	firstPendingSolMintTransaction := collections.NewItem(schemaBuilder, collections.NewPrefix(1), "first_pending_sol_mint", collections.Uint64Value)
	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// All Solana transactions are STAKED
	testData := []types.PendingMintTransaction{
		{Id: 1, ChainType: types.WalletType_WALLET_TYPE_SOLANA, Status: types.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED, Amount: 100},
		{Id: 2, ChainType: types.WalletType_WALLET_TYPE_SOLANA, Status: types.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED, Amount: 200},
	}

	for _, item := range testData {
		err := pendingMintTransactionsMap.Set(amberCtx, item.Id, item)
		require.NoError(t, err)
	}

	initialFirstPendingSol := uint64(1)
	err = firstPendingSolMintTransaction.Set(amberCtx, initialFirstPendingSol)
	require.NoError(t, err)

	// Run the migration
	err = v8.RemoveStakedMints(amberCtx, pendingMintTransactionsMap, firstPendingSolMintTransaction)
	require.NoError(t, err)

	// All Solana transactions should be removed
	iter, err := pendingMintTransactionsMap.Iterate(amberCtx, nil)
	require.NoError(t, err)
	keys, err := iter.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 0, "All Solana STAKED transactions should be removed")

	// First pending index should remain unchanged since no valid Solana transactions exist
	newFirstPendingSol, err := firstPendingSolMintTransaction.Get(amberCtx)
	require.NoError(t, err)
	require.Equal(t, initialFirstPendingSol, newFirstPendingSol, "First pending Sol should remain unchanged")
}

func TestRemoveStakedMints_NonAmberChain(t *testing.T) {
	encCfg := makeTestEncodingConfig()
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	// Use a non-amber chain ID
	nonAmberCtx := testutil.DefaultContext(storeKey, tKey).WithChainID("gardia-1")
	kvStoreService := runtime.NewKVStoreService(storeKey)

	schemaBuilder := collections.NewSchemaBuilder(kvStoreService)
	pendingMintTransactionsMap := collections.NewMap(schemaBuilder, collections.NewPrefix(0), "pending_mint_transactions", collections.Uint64Key, codec.CollValue[types.PendingMintTransaction](cdc))
	firstPendingSolMintTransaction := collections.NewItem(schemaBuilder, collections.NewPrefix(1), "first_pending_sol_mint", collections.Uint64Value)
	_, err := schemaBuilder.Build()
	require.NoError(t, err)

	// Sample Solana STAKED transactions
	testData := []types.PendingMintTransaction{
		{Id: 1, ChainType: types.WalletType_WALLET_TYPE_SOLANA, Status: types.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED, Amount: 100},
		{Id: 2, ChainType: types.WalletType_WALLET_TYPE_SOLANA, Status: types.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED, Amount: 200},
	}

	for _, item := range testData {
		err := pendingMintTransactionsMap.Set(nonAmberCtx, item.Id, item)
		require.NoError(t, err)
	}

	initialFirstPendingSol := uint64(1)
	err = firstPendingSolMintTransaction.Set(nonAmberCtx, initialFirstPendingSol)
	require.NoError(t, err)

	// Run the migration on non-amber chain
	err = v8.RemoveStakedMints(nonAmberCtx, pendingMintTransactionsMap, firstPendingSolMintTransaction)
	require.NoError(t, err)

	// All transactions should still exist (migration should be skipped)
	iter, err := pendingMintTransactionsMap.Iterate(nonAmberCtx, nil)
	require.NoError(t, err)
	keys, err := iter.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 2, "No transactions should be removed on non-amber chain")

	// First pending index should remain unchanged
	newFirstPendingSol, err := firstPendingSolMintTransaction.Get(nonAmberCtx)
	require.NoError(t, err)
	require.Equal(t, initialFirstPendingSol, newFirstPendingSol, "First pending Sol should remain unchanged on non-amber chain")
}
