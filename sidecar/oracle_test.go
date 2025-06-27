package main

import (
	"sync"
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/go-client"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gagliardetto/solana-go"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
)

func TestFetchSolanaBurnEvents_PersistsUnprocessedEvents(t *testing.T) {
	

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
	solanaClient := solanarpc.New(cfg.SolanaRPC[cfg.Network])
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