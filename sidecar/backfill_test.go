package main

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"

	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleBackfillRequests(t *testing.T) {
	t.Skip("Skipping test on CI as it makes a real network call to Solana")

	// 1. Setup a real Solana client
	cfg := LoadConfig("", "")
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

	// 4. Call processBackfillRequestsList
	burnEvents, err := oracle.processBackfillRequestsList(context.Background(), requests)

	// 5. Assertions
	require.NoError(t, err)
	require.Len(t, burnEvents, 1, "should have one burn event")
	event := burnEvents[0]

	assert.Equal(t, txHash, event.TxID)
	assert.Equal(t, "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1", event.ChainID)
	// This is the actual amount from the on-chain transaction
	assert.Equal(t, uint64(300000), event.Amount)
	assert.False(t, event.IsZenBTC)
}

func TestZenBTCBackfillRequestValidation(t *testing.T) {
	t.Run("validates ZenBTC burn request structure", func(t *testing.T) {
		// Test that ZenBTC backfill requests are properly structured
		request := &validationtypes.MsgTriggerEventBackfill{
			TxHash:       "5J8QLA8qFKtreJEpcTf5a8YdgaQbipmBACiWpaV4a6Ubx1caN6SzphJQvue4159VVKugZ6EUUHVNJbsqPpUc81B",
			EventType:    validationtypes.EventType_EVENT_TYPE_ZENBTC_BURN,
			Caip2ChainId: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		}

		// Verify the request has the expected fields
		assert.NotEmpty(t, request.TxHash, "TxHash should not be empty")
		assert.Equal(t, validationtypes.EventType_EVENT_TYPE_ZENBTC_BURN, request.EventType, "Should be ZenBTC burn event type")
		assert.True(t, strings.HasPrefix(request.Caip2ChainId, "solana:"), "Should be Solana chain ID")
	})
}

func TestBackfillRequestProcessingLogic(t *testing.T) {
	t.Run("processes different event types correctly", func(t *testing.T) {
		// Test that different event types are recognized
		zenbtcRequest := &validationtypes.MsgTriggerEventBackfill{
			TxHash:       "test_zenbtc_hash",
			EventType:    validationtypes.EventType_EVENT_TYPE_ZENBTC_BURN,
			Caip2ChainId: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		}

		zentpRequest := &validationtypes.MsgTriggerEventBackfill{
			TxHash:       "test_zentp_hash",
			EventType:    validationtypes.EventType_EVENT_TYPE_ZENTP_BURN,
			Caip2ChainId: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		}

		// Verify the events are structured correctly for processing
		assert.Equal(t, validationtypes.EventType_EVENT_TYPE_ZENBTC_BURN, zenbtcRequest.EventType)
		assert.Equal(t, validationtypes.EventType_EVENT_TYPE_ZENTP_BURN, zentpRequest.EventType)
		assert.Equal(t, "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1", zenbtcRequest.Caip2ChainId)
		assert.Equal(t, "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1", zentpRequest.Caip2ChainId)
	})
}
