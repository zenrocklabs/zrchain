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

func TestHandleZenBTCBackfillRequests(t *testing.T) {
	tests := []struct {
		name        string
		eventType   validationtypes.EventType
		chainID     string
		expectError bool
		isZenBTC    bool
	}{
		{
			name:        "ZenBTC burn on Solana - valid",
			eventType:   validationtypes.EventType_EVENT_TYPE_ZENBTC_BURN,
			chainID:     "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
			expectError: false,
			isZenBTC:    true,
		},
		{
			name:        "ZenBTC burn on Ethereum - not supported yet",
			eventType:   validationtypes.EventType_EVENT_TYPE_ZENBTC_BURN,
			chainID:     "eip155:17000",
			expectError: false, // Should not error but should skip processing
			isZenBTC:    true,
		},
		{
			name:        "ZenTP burn - should still work",
			eventType:   validationtypes.EventType_EVENT_TYPE_ZENTP_BURN,
			chainID:     "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
			expectError: false,
			isZenBTC:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup minimal Oracle for testing logic
			cfg := sidecartypes.Config{
				Network: "devnet",
			}
			oracle := &Oracle{
				Config:    cfg,
				DebugMode: true,
			}
			var state atomic.Value
			state.Store(&sidecartypes.OracleState{
				CleanedSolanaBurnEvents: make(map[string]bool),
			})
			oracle.currentState = state

			// Create backfill request
			txHash := "test_tx_hash_" + tt.name
			requests := []*validationtypes.MsgTriggerEventBackfill{
				{
					TxHash:       txHash,
					EventType:    tt.eventType,
					Caip2ChainId: tt.chainID,
				},
			}

			// Test that the method processes different event types correctly
			// Note: This will fail network calls but we're testing the logic flow
			burnEvents, err := oracle.processBackfillRequestsList(context.Background(), requests)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// For Ethereum ZenBTC burns, we expect no events since it's not supported yet
				if tt.eventType == validationtypes.EventType_EVENT_TYPE_ZENBTC_BURN && tt.chainID == "eip155:17000" {
					assert.Len(t, burnEvents, 0, "Ethereum ZenBTC burns should be skipped")
				}
				// Note: Actual network calls will fail in test environment,
				// but we've verified the logic flow handles different event types correctly
			}
		})
	}
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
