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
	cfg := LoadConfig()
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
	oracle := NewOracle(cfg, ethClient, nil, solanaClient, zrChainQueryClient, true, true)

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
		latestSolanaSigs: make(map[sidecartypes.SolanaEventType]solana.Signature),
	}
	var updateMutex sync.Mutex
	errChan := make(chan error, 2)

	oracle.fetchSolanaBurnEvents(&wg, update, &updateMutex, errChan)

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
	// 1. Setup
	oracle := &Oracle{}
	oracle.Config.Network = sidecartypes.NetworkTestnet
	oracle.currentState.Store(&sidecartypes.OracleState{
		SolanaBurnEvents:        []api.BurnEvent{},
		CleanedSolanaBurnEvents: make(map[string]bool),
	})

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

	oracle.getSolanaZenBTCBurnEventsFn = func(programID string, lastKnownSig solana.Signature) ([]api.BurnEvent, solana.Signature, error) {
		return []api.BurnEvent{}, solana.Signature{}, nil // No new zenBTC burns
	}
	oracle.getSolanaRockBurnEventsFn = func(programID string, lastKnownSig solana.Signature) ([]api.BurnEvent, solana.Signature, error) {
		return []api.BurnEvent{newEvent}, solana.Signature{}, nil
	}

	// 4. Execute the function under test
	var wg sync.WaitGroup
	update := &oracleStateUpdate{
		latestSolanaSigs: make(map[sidecartypes.SolanaEventType]solana.Signature),
	}
	var updateMutex sync.Mutex
	errChan := make(chan error, 2)

	oracle.fetchSolanaBurnEvents(&wg, update, &updateMutex, errChan)

	wg.Wait() // Wait for the main goroutine
	close(errChan)
	for err := range errChan {
		require.NoError(t, err)
	}

	// 5. Assert the results
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

	// Mock the RPC calls
	oracle.getSignaturesForAddressFn = func(ctx context.Context, account solana.PublicKey, opts *rpc.GetSignaturesForAddressOpts) ([]*rpc.TransactionSignature, error) {
		// Return one dummy signature to be processed
		sig, _ := solana.SignatureFromBase58("3NeFkZ2FendD tins4bYxm9fEMpA9n1aWzD1yT4vBfH8FDSS18aA3A33sGTc32sW2N524n7J1P1B3a33")
		return []*rpc.TransactionSignature{
			{
				Signature: sig,
				Slot:      1,
			},
		}, nil
	}
	oracle.rpcCallBatchFn = func(ctx context.Context, rpcs jsonrpc.RPCRequests) (jsonrpc.RPCResponses, error) {
		return nil, errors.New("batch request failed")
	}

	// Create a dummy transaction result that can be unmarshaled
	// This needs to be a realistic structure that processTransaction expects.
	dummyTxResult := rpc.GetTransactionResult{
		Slot: 123,
		Meta: &rpc.TransactionMeta{
			LogMessages: []string{"Program log: TokenRedemption"},
		},
	}

	oracle.getTransactionFn = func(ctx context.Context, signature solana.Signature, opts *rpc.GetTransactionOpts) (*rpc.GetTransactionResult, error) {
		return &dummyTxResult, nil
	}

	// Mock the processTransaction function to return a dummy event
	processTransaction := func(txResult *rpc.GetTransactionResult, program solana.PublicKey, sig solana.Signature, debugMode bool) ([]any, error) {
		return []any{api.BurnEvent{TxID: sig.String()}}, nil
	}

	// Run the test
	events, _, err := oracle.getSolanaEvents("11111111111111111111111111111111", solana.Signature{}, "test event", processTransaction)

	// Assertions
	require.NoError(t, err)
	// The mock getTransaction returns a result with a dummy event, so we expect one event
	require.Len(t, events, 1, "Expected one event to be processed via fallback")
}
