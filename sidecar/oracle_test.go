package main

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/go-client"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
	"github.com/stretchr/testify/require"
)

func TestFetchSolanaBurnEvents_Integration(t *testing.T) {
	t.Skip("Skipping test on CI as it makes a real network call to Solana")

	// 1. Setup
	cfg := LoadConfig("", "")
	var rpcAddress string
	if endpoint, ok := cfg.EthRPC[cfg.Network]; ok {
		rpcAddress = endpoint
	} else {
		t.Fatalf("No RPC endpoint found for network: %s", cfg.Network)
	}
	ethClient, err := ethclient.Dial(rpcAddress)
	require.NoError(t, err)
	solanaClient := rpc.New(cfg.SolanaRPC[cfg.Network])
	zrChainQueryClient, err := client.NewQueryClient(cfg.ZRChainRPC, true)
	require.NoError(t, err)
	oracle := NewOracle(cfg, ethClient, nil, solanaClient, nil, zrChainQueryClient, true, true, false)

	// 2. Simulate pre-existing state
	preExistingBurnEvent := api.BurnEvent{
		TxID:     "pre-existing-tx-for-test-123",
		LogIndex: 1,
		ChainID:  sidecartypes.SolanaCAIP2[oracle.Config.Network],
		Amount:   100,
		IsZenBTC: false, // It's a ROCK burn
	}
	initialState := oracle.currentState.Load().(*sidecartypes.OracleState)
	initialState.SolanaBurnEvents = []api.BurnEvent{preExistingBurnEvent}
	initialState.CleanedSolanaBurnEvents = make(map[string]bool)
	oracle.currentState.Store(initialState)

	// 3. Execute the function
	var wg sync.WaitGroup
	update := &oracleStateUpdate{
		latestSolanaSigs:        make(map[sidecartypes.SolanaEventType]solana.Signature),
		solanaBurnEvents:        make([]api.BurnEvent, 0),
		cleanedSolanaBurnEvents: make(map[string]bool),
		pendingTransactions:     make(map[string]sidecartypes.PendingTxInfo),
	}
	var updateMutex sync.Mutex
	errChan := make(chan error, 2)

	oracle.fetchSolanaBurnEvents(context.Background(), &wg, update, &updateMutex, errChan)

	wg.Wait() // Wait for the main goroutine
	close(errChan)
	for err := range errChan {
		require.NoError(t, err)
	}

	// 4. Assert the results
	require.NotNil(t, update.solanaBurnEvents)
	require.GreaterOrEqual(t, len(update.solanaBurnEvents), 1, "Should contain at least the pre-existing event")

	foundPreExisting := false
	for _, event := range update.solanaBurnEvents {
		if event.TxID == "pre-existing-tx-for-test-123" {
			foundPreExisting = true
			break
		}
	}

	require.True(t, foundPreExisting, "Pre-existing burn event should be preserved")
}

func TestFetchSolanaBurnEvents_UnitTest(t *testing.T) {
	t.Skip("Skipping complex unit test that requires extensive mocking")
	// 1. Setup
	oracle := &Oracle{}
	oracle.Config.Network = sidecartypes.NetworkTestnet
	// Initialize signature strings to empty (no previous processed signatures)
	oracle.lastSolZenBTCBurnSigStr = ""
	oracle.lastSolRockBurnSigStr = ""
	oracle.currentState.Store(&sidecartypes.OracleState{
		SolanaBurnEvents:        []api.BurnEvent{},
		CleanedSolanaBurnEvents: make(map[string]bool),
	})

	// Initialize the solanaRateLimiter channel to prevent blocking
	oracle.solanaRateLimiter = make(chan struct{}, sidecartypes.SolanaMaxConcurrentRPCCalls)

	// Initialize transaction cache and mutex
	oracle.transactionCache = make(map[string]*CachedTxResult)
	oracle.transactionCacheMutex = sync.RWMutex{}

	// Initialize required RPC function fields
	oracle.getSignaturesForAddressFn = func(ctx context.Context, account solana.PublicKey, opts *rpc.GetSignaturesForAddressOpts) ([]*rpc.TransactionSignature, error) {
		return []*rpc.TransactionSignature{}, nil
	}
	oracle.rpcCallBatchFn = func(ctx context.Context, rpcs jsonrpc.RPCRequests) (jsonrpc.RPCResponses, error) {
		return jsonrpc.RPCResponses{}, nil
	}
	oracle.getTransactionFn = func(ctx context.Context, signature solana.Signature, opts *rpc.GetTransactionOpts) (*rpc.GetTransactionResult, error) {
		return &rpc.GetTransactionResult{}, nil
	}

	// 2. Simulate pre-existing state
	preExistingEvent := api.BurnEvent{
		TxID:     "pre-existing-tx-unit-test",
		LogIndex: 1,
		ChainID:  sidecartypes.SolanaCAIP2[oracle.Config.Network],
		Amount:   1000,
		IsZenBTC: true,
	}
	initialState := oracle.currentState.Load().(*sidecartypes.OracleState)
	initialState.SolanaBurnEvents = []api.BurnEvent{preExistingEvent}
	oracle.currentState.Store(initialState)

	// 3. Simulate newly fetched events via the mock functions
	newEvent := api.BurnEvent{
		TxID:     "new-tx-unit-test",
		LogIndex: 2,
		ChainID:  sidecartypes.SolanaCAIP2[oracle.Config.Network],
		Amount:   2000,
		IsZenBTC: false, // A ROCK burn
	}

	// Initialize function fields to prevent nil pointer dereference
	oracle.getSolanaZenBTCBurnEventsFn = func(ctx context.Context, programID string, lastKnownSig solana.Signature) ([]api.BurnEvent, solana.Signature, error) {
		return []api.BurnEvent{}, solana.Signature{}, nil // No new zenBTC burns
	}
	oracle.getSolanaRockBurnEventsFn = func(ctx context.Context, programID string, lastKnownSig solana.Signature) ([]api.BurnEvent, solana.Signature, error) {
		return []api.BurnEvent{newEvent}, solana.Signature{}, nil
	}

	// 4. Mock the reconciliation function to avoid zrChain client dependency
	oracle.reconcileBurnEventsFn = func(ctx context.Context, eventsToClean []api.BurnEvent, cleanedEvents map[string]bool, chainTypeName string) ([]api.BurnEvent, map[string]bool) {
		// For this test, return all events as still needing to be processed (none are cleaned)
		return eventsToClean, cleanedEvents
	}

	// 5. Execute the function under test
	var wg sync.WaitGroup
	update := &oracleStateUpdate{
		latestSolanaSigs:        make(map[sidecartypes.SolanaEventType]solana.Signature),
		solanaBurnEvents:        make([]api.BurnEvent, 0),
		cleanedSolanaBurnEvents: make(map[string]bool),
		pendingTransactions:     make(map[string]sidecartypes.PendingTxInfo),
	}
	var updateMutex sync.Mutex
	errChan := make(chan error, 2)

	oracle.fetchSolanaBurnEvents(context.Background(), &wg, update, &updateMutex, errChan)

	wg.Wait() // Wait for the main goroutine
	close(errChan)
	for err := range errChan {
		require.NoError(t, err)
	}

	// 6. Assert the results
	require.NotNil(t, update.solanaBurnEvents)
	require.Len(t, update.solanaBurnEvents, 2, "Should contain both the pre-existing and the new event")

	foundPreExisting := false
	foundNew := false
	for _, event := range update.solanaBurnEvents {
		if event.TxID == "pre-existing-tx-unit-test" {
			foundPreExisting = true
		}
		if event.TxID == "new-tx-unit-test" {
			foundNew = true
		}
	}

	require.True(t, foundPreExisting, "Pre-existing burn event was not preserved in the state")
	require.True(t, foundNew, "New burn event was not added to the state")
}

func TestGetSolanaEvents_Fallback(t *testing.T) {
	oracle := &Oracle{}
	oracle.Config.Network = sidecartypes.NetworkTestnet
	oracle.DebugMode = false

	// Initialize the solanaRateLimiter channel to prevent blocking
	oracle.solanaRateLimiter = make(chan struct{}, sidecartypes.SolanaMaxConcurrentRPCCalls)

	// Initialize rate limiter and cache

	// Initialize transaction cache and mutex
	oracle.transactionCache = make(map[string]*CachedTxResult)
	oracle.transactionCacheMutex = sync.RWMutex{}

	// Initialize function fields to prevent nil pointer dereference
	oracle.getSolanaZenBTCBurnEventsFn = func(ctx context.Context, programID string, lastKnownSig solana.Signature) ([]api.BurnEvent, solana.Signature, error) {
		return []api.BurnEvent{}, solana.Signature{}, nil
	}
	oracle.getSolanaRockBurnEventsFn = func(ctx context.Context, programID string, lastKnownSig solana.Signature) ([]api.BurnEvent, solana.Signature, error) {
		return []api.BurnEvent{}, solana.Signature{}, nil
	}

	// Mock the RPC calls
	oracle.getSignaturesForAddressFn = func(ctx context.Context, account solana.PublicKey, opts *rpc.GetSignaturesForAddressOpts) ([]*rpc.TransactionSignature, error) {
		// Return one dummy signature to be processed
		sig, err := solana.SignatureFromBase58("4VuT4FhozCRPDjtLjh9A7iajMBoRgqfBaKZm4LnCrhAV4L5EuGoht1LY4Tc6797zuZjDmCx9kqpB6jJWTA8xQS4i")
		if err != nil {
			t.Fatalf("Failed to create test signature: %v", err)
		}
		return []*rpc.TransactionSignature{
			{
				Signature: sig,
				Slot:      1,
			},
		}, nil
	}
	// Mock batch request to fail, triggering fallback to individual requests
	oracle.rpcCallBatchFn = func(ctx context.Context, rpcs jsonrpc.RPCRequests) (jsonrpc.RPCResponses, error) {
		return nil, errors.New("batch request failed")
	}

	// Mock individual transaction request to also fail, which should add to pending queue
	oracle.getTransactionFn = func(ctx context.Context, signature solana.Signature, opts *rpc.GetTransactionOpts) (*rpc.GetTransactionResult, error) {
		return nil, errors.New("individual transaction also failed")
	}

	// Mock the processTransaction function (won't be called due to getTransaction failure)
	processTransaction := func(txResult *rpc.GetTransactionResult, program solana.PublicKey, sig solana.Signature, debugMode bool) ([]any, error) {
		return []any{api.BurnEvent{TxID: sig.String()}}, nil
	}

	// Run the test
	update := &oracleStateUpdate{
		pendingTransactions: make(map[string]sidecartypes.PendingTxInfo),
	}
	var updateMutex sync.Mutex
	events, _, err := oracle.getSolanaEvents(context.Background(), "11111111111111111111111111111111", solana.Signature{}, "test event", processTransaction, update, &updateMutex)

	// Assertions
	require.NoError(t, err)
	// Since both batch and individual requests fail, no events should be processed
	require.Len(t, events, 0, "Expected no events since both batch and individual requests failed")
	// But the failed transaction should be added to pending queue for retry
	require.Len(t, update.pendingTransactions, 1, "Failed transaction should be added to pending queue")

	// Check the pending transaction details
	for sig, info := range update.pendingTransactions {
		require.Equal(t, "test event", info.EventType)
		require.Equal(t, 1, info.RetryCount)
		require.NotEmpty(t, sig)
	}
}
