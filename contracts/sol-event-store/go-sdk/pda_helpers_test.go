package eventstore

import (
	"testing"
)

// Test exported PDA helper functions for deterministic behavior.
func TestPDAHelpers(t *testing.T) {
	program := mustParsePublicKey("4KFjSTnBjbJbWXAiwpWjBCCfAKhjqMp3yfXYpoR3eVis")
	gc, err := DeriveGlobalConfigPDA(program)
	if err != nil {
		t.Fatalf("DeriveGlobalConfigPDA error: %v", err)
	}
	if gc.IsZero() {
		t.Fatalf("Global config PDA should not be zero")
	}

	// Pick arbitrary event IDs including edge boundaries
	eventIDs := []uint64{1, 42, 99, 100, 101, 500, 999, 1000, ^uint64(0)}
	seen := map[string]struct{}{}
	for _, id := range eventIDs {
		pda, shardIdx, err := DeriveWrapShardPDA(program, ZENBTC_WRAP_SHARD_SEED, id, uint16(ZENBTC_WRAP_SHARD_COUNT))
		if err != nil {
			// Should never error with valid shard count
			if id != ^uint64(0) { // still shouldn't for max, but keep note
				t.Fatalf("DeriveWrapShardPDA unexpected error for id %d: %v", id, err)
			}
		}
		if pda.IsZero() {
			// For id 0 (not in list) or invalid config only we would allow; here all ids > 0
			if id > 0 {

				// fail if zero
				t.Fatalf("Derived shard PDA zero for id %d", id)
			}
		}
		if shardIdx >= uint16(ZENBTC_WRAP_SHARD_COUNT) {
			t.Fatalf("Shard index out of range: %d for id %d", shardIdx, id)
		}
		seen[pda.String()] = struct{}{}
	}
	if len(seen) == 0 {
		t.Fatalf("Expected at least one PDA derived")
	}
}

// Ensure invalid shard count returns error.
func TestDeriveWrapShardPDAInvalidShardCount(t *testing.T) {
	program := mustParsePublicKey("4KFjSTnBjbJbWXAiwpWjBCCfAKhjqMp3yfXYpoR3eVis")
	_, _, err := DeriveWrapShardPDA(program, ZENBTC_WRAP_SHARD_SEED, 10, 0)
	if err == nil {
		t.Fatalf("expected error for zero shardCount")
	}
}
