package main

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleBackfillRequests(t *testing.T) {
	t.Skip("Skipping test on CI as it makes a real network call to Solana")

	// 1. Setup a real Solana client
	cfg := LoadConfig("")
	solanaClient := solanarpc.New(cfg.SolanaRPC[cfg.Network])

	// 2. Setup a minimal Oracle with the real client
	oracle := &Oracle{
		Config:       cfg,
		solanaClient: solanaClient,
		DebugMode:    true,
	}
	var state atomic.Value
	state.Store(&sidecartypes.OracleState{
		CleanedSolanaBurnEvents: make(map[string]bool),
	})
	oracle.currentState = state

	// 3. Create the backfill request struct directly
	txHash := "fjBQLA8qFKtreJEpcTf5a8YdgaQbipmBACiWpaV4a6Ubx1caN6SzphJQvue4159VVKugZ6EUUHVNJbsqPpUc81B"
	requests := []*validationtypes.MsgTriggerEventBackfill{
		{
			TxHash:       txHash,
			EventType:    validationtypes.EventType_EVENT_TYPE_ZENTP_BURN,
			Caip2ChainId: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		},
	}

	// 4. Call handleBackfillRequests
	update := &oracleStateUpdate{
		solanaBurnEvents: make([]api.BurnEvent, 0),
	}
	var updateMutex sync.Mutex

	oracle.handleBackfillRequests(context.Background(), requests, update, &updateMutex)

	// 5. Assertions
	require.Len(t, update.solanaBurnEvents, 1, "should have one burn event")
	event := update.solanaBurnEvents[0]

	assert.Equal(t, txHash, event.TxID)
	assert.Equal(t, "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1", event.ChainID)
	// This is the actual amount from the on-chain transaction
	assert.Equal(t, uint64(300000), event.Amount)
	assert.False(t, event.IsZenBTC)
}
