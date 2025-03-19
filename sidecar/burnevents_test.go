package main

import (
	"math/big"
	"sync/atomic"
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v5/sidecar/shared"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

// Test that events older than EthBurnEventsBlockRange are removed
func TestRemoveStaleEvents(t *testing.T) {
	// Setup test oracle with minimal dependencies
	oracle := &Oracle{
		currentState: atomic.Value{},
	}

	// Current block is 10000
	currentBlockHeight := uint64(10000)
	cutoffHeight := currentBlockHeight - EthBurnEventsBlockRange

	// Setup test events
	testEvents := []api.BurnEvent{
		// Recent event that should be kept
		{
			TxID:        "0x1",
			LogIndex:    0,
			ChainID:     "eip155:1",
			Amount:      1000,
			BlockHeight: cutoffHeight + 1, // Just above cutoff
		},
		// Old event that should be removed
		{
			TxID:        "0x2",
			LogIndex:    0,
			ChainID:     "eip155:1",
			Amount:      2000,
			BlockHeight: cutoffHeight - 1, // Just below cutoff
		},
		// Another old event that should be removed
		{
			TxID:        "0x3",
			LogIndex:    0,
			ChainID:     "eip155:1",
			Amount:      3000,
			BlockHeight: cutoffHeight - 500, // Well below cutoff
		},
		// Edge case - exactly at cutoff (should be removed)
		{
			TxID:        "0x4",
			LogIndex:    0,
			ChainID:     "eip155:1",
			Amount:      4000,
			BlockHeight: cutoffHeight, // Exactly at cutoff
		},
	}

	// Create initial state
	initialState := &sidecartypes.OracleState{
		EthBurnEvents:        testEvents,
		CleanedEthBurnEvents: make(map[string]bool),
	}
	oracle.currentState.Store(initialState)

	// Create a mock header with our fixed block height
	header := &types.Header{
		Number: big.NewInt(int64(currentBlockHeight)),
	}

	// Call the function to remove stale events
	remainingEvents := oracle.removeStaleEvents(initialState, header)

	// Should only have 1 event left (the recent one)
	assert.Equal(t, 1, len(remainingEvents), "Should only have one event remaining")

	// Check it's the right event
	assert.Equal(t, "0x1", remainingEvents[0].TxID, "Only the recent event should remain")

	// Check that the cleaned events map contains the old events
	assert.True(t, initialState.CleanedEthBurnEvents["eip155:1-0x2-0"], "Old event 0x2 should be in cleaned events")
	assert.True(t, initialState.CleanedEthBurnEvents["eip155:1-0x3-0"], "Old event 0x3 should be in cleaned events")
	assert.True(t, initialState.CleanedEthBurnEvents["eip155:1-0x4-0"], "Edge case event 0x4 should be in cleaned events")
}
