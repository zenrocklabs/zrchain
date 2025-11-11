# Sidecar Desync Root Cause Analysis: RPC Data Divergence

**Author**: Claude Code Investigation + Peyton Spencer
**Date**: November 11, 2025
**Severity**: CRITICAL - Blocks consensus on all Solana events
**Status**: Root cause identified; infrastructure fix pending (contract upgrade)

---

## Executive Summary

The persistent sidecar desync issue affecting validators is caused by **different validators fetching different subsets of Solana transactions** due to RPC failures and rate limiting. When transaction fetch operations fail on some validators but succeed on others, they end up with different event lists, resulting in different hash values and consensus failures.

**Root Cause**: Solana RPC rate limiting and transient failures cause `getSignaturesForAddress()` and batch transaction fetches to fail unpredictably. Different validators receive partial/incomplete transaction data, leading to:
- Validator A: Events [1, 2, 3, 4, 5] → Hash 0xAAAA...
- Validator B: Events [1, 2, 3, 5] (missing 4 due to RPC error) → Hash 0xBBBB...
- No 2/3+ consensus on hash → Events not processed

**Affected Components**:
- Solana event fetch functions with RPC-based pagination
- Rate limiting semaphore (20 concurrent calls limit)
- Batch request retry logic (falls back to individual requests)
- Pending transaction queue (stores failed transactions)

**Why It Happens**: Solana RPC endpoints have rate limits and occasional failures. When parallel validators query Solana simultaneously:
1. Some succeed in fetching full transaction list
2. Others hit rate limits, timeouts, or RPC errors
3. Each validator caches partial results
4. Next tick hashes include different transaction sets
5. Hashes diverge → consensus fails

**Real-World Trigger**: 10+ pending transactions + 3-4 validators querying in parallel = high probability at least one validator hits RPC limits

---

## Root Cause Deep Dive

### Event Fetching Architecture with RPC Risk

The sidecar fetches Solana events through signature-based pagination:

```
getSolanaEvents() flow:
  1. Get watermark signature (last processed)
  2. Call getSignaturesForAddress(limit=50)
     ├─ Validator A: Gets 50 signatures ✓
     ├─ Validator B: Gets 40 signatures (rate limited) ⚠️
     └─ Validator C: Gets 50 signatures ✓

  3. Fill signature gaps (up to 10 pages backfill)
     ├─ Validator A: Fills gaps, finds 50 total ✓
     ├─ Validator B: Hits rate limit again, stops at 40 ⚠️
     └─ Validator C: Fills gaps, finds 50 total ✓

  4. Batch fetch transactions
     ├─ Validator A: Batches of 10 × 5 = 50 txs ✓
     ├─ Validator B: Batches of 10 × 4 = 40 txs ⚠️
     └─ Validator C: Batches of 10 × 5 = 50 txs ✓

  5. Result: Different validators have different tx counts
```

### Key Constants Causing Divergence

From `sidecar/shared/config.go`:

```go
SolanaMaxConcurrentRPCCalls = 20          // Semaphore - shared across validators
SolanaEventScanTxLimit      = 50          // Initial signature batch
SolanaMaxBackfillPages      = 10          // Can be hit by slow validators
SolanaEventFetchBatchSize   = 10          // Transaction batch size
SolanaBatchTimeout          = 20 * time.Second  // Can expire on slow RPC
```

**Problem**: Each validator independently hits these limits at different times:
- Fast validators: Complete all fetches before timeout
- Slow validators: Hit timeout, get partial results
- Result: Different event lists

### Watermark Advancement Masks the Issue

The watermark advancement logic in `processSolanaMintEvents()` (lines 961-984) and `fetchSolanaBurnEvents()` (lines 1032-1113) uses:

```go
// Conservative advancement: advance if (new_sig != old_sig) OR (error == nil)
if !newRockSig.Equals(lastKnownRockSig) || rockErr == nil {
    update.latestSolanaSigs[sidecartypes.SolRockMint] = newRockSig
}
```

**This compounds the problem**: Validators with partial failures still advance their watermarks because:
- They got *some* new transactions (newSig != oldSig) ✓
- Or they had no errors (rockErr == nil) ✓

So each validator advances past different transaction sets, then next tick fetches more from that diverged point.

### Event List Sorting Happens But Too Late

Events ARE sorted in `buildFinalState()` (lines 1267-1273 for mint, 1257-1263 for burn):

```go
sort.Slice(update.SolanaMintEvents, func(i, j int) bool {
    if update.SolanaMintEvents[i].Height != update.SolanaMintEvents[j].Height {
        return update.SolanaMintEvents[i].Height < update.SolanaMintEvents[j].Height
    }
    return update.SolanaMintEvents[i].TxSig < update.SolanaMintEvents[j].TxSig
})
```

**But this doesn't help** because:
- Sorting [1, 2, 3, 4, 5] → [1, 2, 3, 4, 5]
- Sorting [1, 2, 3, 5] → [1, 2, 3, 5]
- Both sorted but DIFFERENT event counts
- Different event lists → different hashes → no consensus

---

## Why The Initial Analysis Was Incomplete

Initial investigation (in `2025-10-23-sidecar-desync-analysis.md`) correctly identified:
- ✅ RPC rate limiting as a theory
- ✅ Watermark advancement mechanism
- ✅ Batch processing error handling
- ❌ But concluded non-deterministic ordering was the root cause

**The real issue**: Not ordering, but **data divergence at the RPC fetch layer**. The ordering would only matter if validators had the same data in different orders—but they don't have the same data at all.

---

## Impact Assessment

### What's Affected

- ❌ All Solana mint events (ROCK, zenBTC, zenZEC)
- ❌ All Solana burn events (zenBTC, zenZEC, ROCK)
- ❌ Any reliance on `SolanaMintEventsHash` consensus
- ❌ Any reliance on `SolanaBurnEventsHash` consensus

### What Happens When Desync Occurs

In `PreBlocker` (`x/validation/keeper/abci.go:414-424`):

```go
if fieldHasConsensus(oracleData.FieldVotePowers, VEFieldSolanaMintEventsHash) {
    k.processSolanaZenBTCMintEvents(ctx, oracleData)      // Only if consensus!
    k.processSolanaROCKMintEvents(ctx, oracleData)
    k.processSolanaDCTMintEvents(ctx, oracleData)
}
```

**Result**: Events silently skipped when no consensus, pending transactions remain in DEPOSITED state indefinitely.

---

## Solution: Contract Upgrade (Infrastructure Fix)

The permanent fix is **upgrading Solana contracts** to provide a more reliable event fetching mechanism, likely:

1. **Event Store Program**: Central event log instead of signature-based pagination
2. **Deterministic cursors**: Track events by sequence number instead of signature
3. **Reduced RPC dependency**: Less reliant on rate-limited API calls
4. **Guaranteed consistency**: All validators can fetch same event set

**Timeline**: Planned for mainnet upgrade in the coming days

---

## Defensive Code Changes (Not Merged)

A PR was created with defensive sorting using `slices.Sorted(maps.Keys())` in mint and burn event loops (commit `b5c4ee6e`). This adds an extra layer of determinism but doesn't address the root cause. **The PR was closed** because:

1. Events are already sorted in `buildFinalState()` before hashing
2. The real issue is RPC data divergence, not ordering
3. Waiting for contract upgrade (infrastructure fix) is the proper solution

The defensive code is still a good practice pattern for future reference.

## Temporary Mitigations (Current)

While awaiting contract upgrade:

1. **Reduced reset period**: Restart consensus loop more frequently to force event reprocessing
2. **Pending transaction retry queue**: Store failed transactions, retry next tick
3. **Watermark advancement safety**: Conservative logic prevents watermark blocking
4. **Manual validator restarts**: Clear cached partial state when divergence detected

---

## Code References

**Event Fetching Components**:
- `sidecar/oracle.go:1776+` - `getSolanaZenBTCBurnEvents()`
- `sidecar/oracle.go:2000+` - `getSolZenBTCMints()`
- `sidecar/oracle.go:785+` - `processSolanaMintEvents()` (parallel fetching)
- `sidecar/oracle.go:1003+` - `fetchSolanaBurnEvents()` (parallel fetching)

**Watermark Management**:
- `sidecar/oracle.go:961-984` - Mint watermark advancement
- `sidecar/oracle.go:1032-1113` - Burn watermark advancement

**Rate Limiting**:
- `sidecar/shared/config.go:140-148` - Constants
- `sidecar/oracle.go:2950+` - `executeBatchRequest()` with semaphore

**Validation Module**:
- `x/validation/keeper/abci.go:414-424` - Consensus checks in PreBlocker
- `x/validation/keeper/abci.go:158-255` - Vote extension construction

---

## Conclusion

The sidecar desync issue is fundamentally caused by **Solana RPC unreliability** making different validators fetch different event subsets, not by code bugs in event ordering. The watermark advancement mechanism ensures the chain doesn't stall, but different validators can't agree on hashes when they have different data.

The real fix is **infrastructure-level**: upgrading the Solana contracts to provide event stores with deterministic, rate-limit-independent fetching. Short-term mitigations (reset periods, retries) keep the system functioning but don't solve the root cause.

**Expected Resolution**: Post-contract upgrade, validators will be able to fetch event data deterministically and consensus should reach consistently.
