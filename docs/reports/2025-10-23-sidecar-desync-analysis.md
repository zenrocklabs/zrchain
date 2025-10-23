# Sidecar Desync Analysis: Solana Events Hash Mismatch

**Author**: Peyton-Spencer  
**Date**: October 23, 2025  
**Topics Covered**: Sidecar Oracle, Solana RPC Rate Limiting, Vote Extension Consensus, Event Processing

---

## Table of Contents

- [Overview](#overview)
- [Executive Summary](#executive-summary)
- [Architecture Analysis](#architecture-analysis)
- [Primary Theory: RPC Rate Limiting](#primary-theory-rpc-rate-limiting)
- [Additional Theories](#additional-theories)
- [Code References](#code-references)
- [Root Cause Analysis](#root-cause-analysis)
- [Recommendations](#recommendations)
- [Next Steps](#next-steps)

---

## Executive Summary

Sidecars are desyncing in testnet, with validators failing to reach consensus on the `SolanaMintEventsHash` and `SolanaBurnEventsHash` fields in vote extensions. This indicates that different sidecar instances are fetching different sets of pending Solana events, resulting in different hash values.

The primary suspected cause is Solana RPC rate limiting, which could cause different sidecars to receive partial or inconsistent data when querying for events. This report analyzes the event fetching architecture and proposes multiple theories for why desync occurs.

[🔝 back to top](#table-of-contents)

---

## Overview

The sidecar oracle is responsible for fetching off-chain data (including Solana events) and submitting it to validators via vote extensions. During consensus in the PreBlocker phase, validators hash the Solana events and vote on the hash. A supermajority (2/3+) must agree on the hash for it to be accepted.

**Problem**: Validators are not reaching consensus on `SolanaMintEventsHash` and `SolanaBurnEventsHash`, suggesting different sidecars are seeing different pending events.

[🔝 back to top](#table-of-contents)

---

## Architecture Analysis

### Event Fetching Flow

1. **Main Loop Tick** (every 60 seconds by default)
   - Parallel goroutines fetch Solana mint events and burn events
   - Each event type (zenBTC mint, ROCK mint, zenZEC mint, zenBTC burn, ROCK burn, zenZEC burn) has its own watermark signature

2. **Event Fetching Process** (`getSolanaEvents()`)
   - Uses watermark-based pagination: fetches signatures since last known signature
   - Calls `getSignaturesForAddress()` with `limit=64` (configurable)
   - Detects gaps and backfills up to 10 pages if watermark not found
   - Processes signatures in batches with adaptive batch size (2-10 transactions per batch)
   - Uses rate limiting semaphore with max 20 concurrent RPC calls

3. **Batch Processing** (`processSignatures()`)
   - Groups signatures into batches (default: 10 per batch)
   - Calls `rpcCallBatch()` to fetch multiple transactions at once
   - Handles batch failures with exponential backoff (reduces batch size by half)
   - Falls back to individual transaction fetches if batch fails
   - Failed transactions are added to pending queue for retry

4. **Reconciliation** (`reconcileMintEventsWithZRChain()` and `reconcileBurnEventsWithZRChain()`)
   - Checks cached events against on-chain state
   - Removes events that have been confirmed on-chain
   - Preserves events that are still pending

5. **Event Merging** (`mergeNewMintEvents()` and `mergeNewBurnEvents()`)
   - Combines newly fetched events with existing pending events
   - Uses transaction signature as deduplication key
   - Maintains chronological order by block height and log index

6. **Vote Extension Construction** (`ConstructVoteExtension()`)
   - Hashes the final list of events using `deriveHash()`
   - Includes hash in vote extension for consensus

### Rate Limiting Mechanisms

```go
// From sidecar/shared/types.go
SolanaMaxConcurrentRPCCalls = 20          // Semaphore size
SolanaRPCTimeout            = 10 * time.Second
SolanaBatchTimeout          = 20 * time.Second
SolanaRateLimiterTimeout    = 10 * time.Second
SolanaEventScanTxLimit      = 64          // Signatures per page
SolanaMaxBackfillPages      = 10          // Max pages to backfill
SolanaEventFetchBatchSize   = 10          // Transactions per batch
```

[🔝 back to top](#table-of-contents)

---

## Primary Theory: RPC Rate Limiting

### Theory Statement

**Solana RPC nodes implement rate limiting that causes different sidecars to receive partial or inconsistent responses when making parallel requests. This leads to different watermark advancement and different sets of pending events.**

### Evidence Supporting This Theory

1. **Aggressive Batching**: The sidecar fetches multiple event types concurrently (6 event types × multiple batches), creating burst traffic patterns that may trigger rate limits.

2. **Watermark Advancement on Partial Success**: If a sidecar successfully processes some transactions but fails on others, it may still advance the watermark. The code at `sidecar/oracle.go:3466-3492` blocks watermark advancement if pending store operations fail, but transient RPC errors during batch processing could still cause inconsistency.

3. **Gap Detection Failures**: When the watermark signature is not found in the initial page of 64 signatures, the sidecar attempts to backfill up to 10 pages (`SolanaMaxBackfillPages`). If rate limiting prevents successful backfill, different sidecars may end up with different signature ranges.

4. **Batch Failure Handling**: When a batch RPC call fails, the code reduces the batch size and retries. However, the reduced batch size may fetch a different subset of transactions if the RPC node's view has changed between retries.

### Code Locations

- **Rate limiter acquisition**: `sidecar/oracle.go:3425-3433`
  ```go
  select {
  case o.solanaRateLimiter <- struct{}{}:
      defer func() { <-o.solanaRateLimiter }()
  case <-ctx.Done():
      return nil, lastKnownSig, ctx.Err()
  case <-time.After(sidecartypes.SolanaRateLimiterTimeout):
      return nil, lastKnownSig, fmt.Errorf("rate limiter timeout...")
  }
  ```

- **Batch request execution**: `sidecar/oracle.go:2880-2932`
  ```go
  batchResponses, batchErr := o.rpcCallBatchFn(batchCtx, batchRequests)
  if batchErr != nil {
      newBatchSize := max(currentBatchSize/2, minBatchSize)
      return nil, newBatchSize, batchErr
  }
  ```

- **Gap filling logic**: `sidecar/oracle.go:3528-3592`
  ```go
  for page := 0; page < sidecartypes.SolanaMaxBackfillPages; page++ {
      // ... backfill logic
  }
  slog.Error("Unable to locate starting watermark signature after scanning maximum pages")
  ```

### Why This Causes Desync

1. **Sidecar A** makes batch request → rate limited → receives partial response → advances watermark to last successful transaction
2. **Sidecar B** makes batch request → not rate limited (or different rate limit bucket) → receives complete response → advances watermark further
3. **Next tick**: Sidecar A and B start from different watermarks
4. **Result**: Sidecar A has more pending events than Sidecar B because its watermark is behind
5. **Consequence**: Different event lists → different hashes → no consensus

### Assessment

**Plausible.** The sidecar still advances Solana watermarks whenever it observes a newer signature, even if some transactions only make it into the local pending queue (`sidecar/oracle.go:965-986`, `sidecar/oracle.go:3160-3212`, `sidecar/oracle.go:3466-3492`). Because batch and fallback processing respond to RPC errors independently on every validator, the pending sets diverge while the shared watermark moves forward, reproducing the observed hash mismatches.

[🔝 back to top](#table-of-contents)

---

## Additional Theories

### Theory 2: Non-Deterministic Event Ordering

**Statement**: Solana events fetched in parallel (ROCK mint, zenBTC mint, zenZEC mint, etc.) are merged into a single list and sorted. If events have identical block heights and log indices, the sort may be non-deterministic across sidecars.

**Supporting Code**: `sidecar/oracle.go:1255-1261`
```go
sort.Slice(update.SolanaMintEvents, func(i, j int) bool {
    if update.SolanaMintEvents[i].Height != update.SolanaMintEvents[j].Height {
        return update.SolanaMintEvents[i].Height < update.SolanaMintEvents[j].Height
    }
    // Use TxSig as a secondary sort key for determinism if heights are identical
    return update.SolanaMintEvents[i].TxSig < update.SolanaMintEvents[j].TxSig
})
```

**Mitigation**: The code includes `TxSig` as a secondary sort key, so this should be deterministic. However, if TxSig values are ever empty or malformed, this could cause issues.

**Likelihood**: Low (secondary sort key should prevent this)

**Assessment**: Unlikely. `update.SolanaMintEvents` is sorted deterministically by height and signature (`sidecar/oracle.go:1255-1260`), and mint decoding always supplies a `TxSig`, so there is no obvious path to nondeterministic ordering.

---

### Theory 3: Reconciliation Timing Issues

**Statement**: The reconciliation step queries the zrChain state to remove confirmed events. If sidecars perform reconciliation at slightly different times (due to network latency or processing delays), they may see different on-chain states and thus retain different sets of pending events.

**Supporting Code**: `sidecar/oracle.go:1485-1578` (reconciliation logic)

**Why This Happens**:
- Sidecar A reconciles at block N → removes events confirmed in block N
- Sidecar B reconciles at block N+1 → removes events confirmed in blocks N and N+1
- Result: Sidecar A retains more pending events

**Likelihood**: Medium (timing-dependent)

**Assessment**: Plausible. Each validator queries pending mint status independently and keeps events when those RPC calls fail (`sidecar/oracle.go:1502-1548`), so timing or error differences can leave one sidecar with an event that another has already dropped.

---

### Theory 4: Transaction Cache Inconsistency

**Statement**: The sidecar caches transaction results with a 5-minute TTL. If a cached transaction is used by one sidecar but expired for another, they may process transactions differently.

**Supporting Code**: `sidecar/oracle.go:3495-3526` (cache implementation)

**Why This Could Cause Issues**:
- Cached transactions may contain different data than fresh fetches if the RPC node's view has changed
- Cache expiration is time-based, not deterministic across sidecars

**Likelihood**: Low (cache is only used for performance, not correctness)

**Assessment**: Unlikely. Cached `getTransaction` results are just reused within the same sidecar during retries and expire after five minutes (`sidecar/oracle.go:3383-3405`, `sidecar/oracle.go:3508-3512`), so once the RPC succeeds every validator will decode the same event payload.

---

### Theory 5: Pending Transaction Queue Divergence

**Statement**: Failed transactions are added to a pending queue (`PendingSolanaTxs`) for retry. If sidecars experience different failure patterns, their pending queues diverge, leading to different event lists.

**Supporting Code**: `sidecar/oracle.go:3466-3486`
```go
for _, failedSig := range failedSignatures {
    if pendingErr := o.addPendingTransaction(failedSig, eventTypeName, update, updateMutex); pendingErr != nil {
        // ...
    }
}
```

**Why This Happens**:
- Transient RPC failures are non-deterministic
- Different sidecars may fail on different transactions
- Pending transactions are retried in subsequent ticks, creating divergent states

**Likelihood**: High (directly related to RPC reliability)

**Assessment**: Highly plausible. Failed signatures are appended to a per-instance pending map while the watermark advances (`sidecar/oracle.go:3160-3212`, `sidecar/oracle.go:3466-3492`), and the asynchronous retry loop (`sidecar/oracle.go:2213-2242`) gives plenty of room for validators to disagree on which events remain outstanding.

---

### Theory 6: Signature Gap Handling Inconsistency

**Statement**: When a gap is detected (watermark not found in initial page), the backfill logic may behave differently across sidecars if the RPC node returns different results during pagination.

**Supporting Code**: `sidecar/oracle.go:3528-3592`

**Critical Issue**: If backfill fails or times out after scanning `SolanaMaxBackfillPages` (10 pages), the code proceeds with whatever data was collected:
```go
slog.Error("Unable to locate starting watermark signature after scanning maximum pages. Proceeding with collected data")
```

This means different sidecars may collect different amounts of data before giving up, leading to divergent watermarks.

**Likelihood**: High (especially if RPC nodes are rate-limited or slow)

**Assessment**: Plausible. After ten backfill pages the helper returns whatever it has and logs an error (`sidecar/oracle.go:3528-3591`), so validators connected to different RPC histories can assemble different signature windows.

---

### Theory 7: Batch Response Error Handling Asymmetry

**Statement**: When a batch RPC call returns partial errors (some responses succeed, some fail), the error handling may cause different sidecars to process different subsets of transactions.

**Supporting Code**: `sidecar/oracle.go:2902-2923`
```go
if batchErr == nil {
    errorCount := 0
    for i, resp := range batchResponses {
        if resp.Error != nil {
            errorCount++
            batchErr = fmt.Errorf("response contains errors: %v", resp.Error)
            break  // <-- Stops at first error
        }
    }
}
```

**Issue**: The loop breaks on the first error, meaning transactions processed before the error are handled, but subsequent ones are not. If the error occurs at different positions in the batch for different sidecars, they process different subsets.

**Likelihood**: Medium (depends on RPC response consistency)

**Assessment**: Plausible. The batch executor aborts on the first error (`sidecar/oracle.go:2880-2932`) and pushes any remaining failures into the fallback path, which marks those signatures pending (`sidecar/oracle.go:3160-3212`); the exact failure position varies per validator, feeding the divergence described in Theories 1 and 5.

[🔝 back to top](#table-of-contents)

---

## Why `getTransaction` Calls Are Required

`getSignaturesForAddress` gives us ordered signatures and slots, but it omits the data required by the bridge pipeline:

- **Per-Event Data**: Sidecars must know recipient, amount, fee, mint, and log index to populate `api.SolanaMintEvent`/`api.BurnEvent`. These fields are decoded from the transaction logs in `processMintTransaction` and `processBurnTransaction` (`sidecar/oracle.go:3594-3619`, `sidecar/oracle.go:2652-2721`).
- **Multiple Events per Signature**: One Solana transaction can emit several mint/burn events. Hashing the signature alone would collapse distinct `SigHash` values and break the dedup logic in `mergeNewMintEvents` (`sidecar/utils.go:314-356`).
- **Bridge Reconciliation**: The validation module later consumes the full event payload to credit mint amounts and complete redemptions (`x/validation/keeper/abci_dct.go:processDCTMintsSolana`, `x/validation/keeper/abci_dct.go:processSolanaDCTMintEvents`). Signatures alone cannot drive those state updates.
- **Failure Classification**: `getTransaction` tells us whether the transaction actually emitted the expected program logs. Without it, pending retries (`sidecar/oracle.go:2213-2242`) could never confirm completion.

Therefore the sidecar hashes full event objects for consensus *and* retains them for downstream processing; fetching only signatures would lose critical information.

[🔝 back to top](#table-of-contents)

---

## Code References

### Sidecar Oracle Files

- [`sidecar/oracle.go`](../../sidecar/oracle.go) - Core oracle loop and event fetching logic
  - Lines 783-992: `processSolanaMintEvents()` - Mint event fetching
  - Lines 994-1197: `fetchSolanaBurnEvents()` - Burn event fetching
  - Lines 3410-3493: `getSolanaEvents()` - Generic event fetching with watermarks
  - Lines 3528-3592: `fetchAndFillSignatureGap()` - Gap detection and backfill
  - Lines 3089-3288: `processSignatures()` - Batch processing with adaptive sizing
  - Lines 2880-2932: `executeBatchRequest()` - Batch RPC execution
  - Lines 1485-1578: `reconcileMintEventsWithZRChain()` - Event reconciliation

- [`sidecar/types.go`](../../sidecar/types.go) - Oracle struct and type definitions
  - Lines 44-86: Oracle struct with rate limiter and cache fields

- [`sidecar/shared/types.go`](../../sidecar/shared/types.go) - Constants and configuration
  - Lines 119-173: Tuning parameters (rate limits, batch sizes, timeouts)

### Validation Module Files

- [`x/validation/keeper/abci.go`](../../x/validation/keeper/abci.go) - Vote extension construction
  - Lines 158-255: `ConstructVoteExtension()` - Hashes events for consensus
  - Lines 214-221: Solana event hash derivation

- [`x/validation/keeper/abci_types.go`](../../x/validation/keeper/abci_types.go) - Vote extension types
  - Lines 289-318: Vote extension field enums
  - Lines 430-463: Field handlers for Solana events

- [`x/validation/keeper/abci_utils.go`](../../x/validation/keeper/abci_utils.go) - Validation logic
  - Lines 1566-1575: Hash validation for Solana events

[🔝 back to top](#table-of-contents)

---

## Root Cause Analysis

### Primary Root Cause: Race Conditions During RPC Failures

The core issue is that **watermark advancement and event accumulation are not atomic operations**. When RPC failures occur (due to rate limiting or other issues), different sidecars can end up with:

1. **Different watermarks** (last successfully processed signature)
2. **Different pending event lists** (failed transactions added to pending queue)
3. **Different reconciliation states** (timing-dependent on-chain state queries)

These three sources of divergence compound over time, causing hash mismatches.

### Cascade Effect

```
Tick 1:
├── Sidecar A: RPC rate limited → processes 40/64 txs → watermark advances to tx_40
└── Sidecar B: RPC succeeds    → processes 64/64 txs → watermark advances to tx_64

Tick 2:
├── Sidecar A: Fetches from tx_40 → gets 24 "duplicate" txs + 40 new txs
└── Sidecar B: Fetches from tx_64 → gets 64 new txs

Result:
├── Sidecar A: 64 pending events (40 + 40 - 16 confirmed)
└── Sidecar B: 50 pending events (64 - 14 confirmed)

Hash mismatch!
```

[🔝 back to top](#table-of-contents)

---

## Recommendations

### Immediate Actions (High Priority)

1. **Add Deterministic Watermark Checkpoints**
   - After each successful fetch, store watermark + event count in state file
   - Compare watermarks across sidecars before vote extension construction
   - If watermarks differ, log warning and sync to lowest common watermark

2. **Implement Backoff and Retry with Jitter**
   - Add exponential backoff with jitter to RPC calls to reduce thundering herd
   - Spread out parallel event fetches with small delays

3. **Add Diagnostic Logging**
   - Log watermark positions for each event type after every tick
   - Log pending event count and hash values
   - Log RPC call latencies and failure rates
   - Add metrics for rate limit detection (429 responses, timeout frequency)

4. **Improve Batch Error Handling**
   - Instead of breaking on first error, process all successful responses
   - Only mark failed transactions as pending, don't discard entire batch

### Medium-Term Fixes (Medium Priority)

5. **Introduce Gossip-Based Consensus Round for Events**
   - Stand up a lightweight libp2p-style mesh where each sidecar publishes its `latestSolanaSigs` and pending signature digest (`sidecar/oracle.go:889-939`, `sidecar/oracle.go:3466-3492`)
   - Require signed gossip frames so peers can reject tampered data and converge on the minimum watermark before vote extension construction
   - Use the mesh for early divergence alerts while zrChain remains the canonical decision point in PreBlocker

6. **Implement Event Snapshot Synchronization**
   - Periodically (e.g., every 10 minutes), sidecars exchange pending event lists
   - Detect divergence early and resync before hash mismatch occurs

7. **Add RPC Health Monitoring**
   - Track success rate per RPC endpoint
   - Dynamically switch to backup RPC nodes if primary is rate-limited
   - Use multiple RPC nodes and cross-validate responses

### Long-Term Solutions (Lower Priority)

8. **Move to Pull-Based Event Model**
   - Instead of pushing events via vote extensions, validators pull from sidecar
   - Sidecars maintain deterministic event streams
   - Validators query stream at specific block heights for consistency

9. **Implement Merkle Proof-Based Events**
   - Hash individual events into a Merkle tree
   - Vote extensions include only Merkle root
   - Validators can fetch missing events and verify against root

10. **Redesign Watermark Management**
    - Store watermarks on-chain as part of validation module state
    - All validators advance watermark atomically during PreBlocker
    - Sidecars fetch events relative to on-chain watermark

[🔝 back to top](#table-of-contents)

---

## Proposed Sidecar Gossip Network

- **Purpose**: Give validators rapid, off-chain visibility into each peer’s Solana watermarks and pending signature queue so the group can clamp to the lowest shared position before hashing events. This supplements the existing vote extension without requiring an immediate zrChain upgrade.
- **Protocol Shape**: Reuse the existing oracle state update machinery (`sidecar/oracle.go:352-942`) to emit compact frames containing: event type, watermark signature, count of pending retries, and a hash of `SolanaMintEvents`/`solanaBurnEvents`. Peers sign each frame with their validator key so recipients can drop unauthenticated data.
- **Network Layer**: A libp2p (or quic-go) mesh keyed by validator identities. Bootstrapping peers from zrChain validator metadata keeps discovery simple, while gossipsub topics per asset (`sol-zenbtc-mint`, `sol-zenzec-burn`, etc.) cap noise.
- **Divergence Handling**: If a peer advertises a higher watermark or conflicting hash, we back off and rescan from the minimum advertised signature, preventing the optimistic watermark advancement noted in `sidecar/oracle.go:965-986`.
- **Integration Path**: Start as an optional “warn-only” daemon that logs mismatches. Once stable, feed the lowest observed watermark back into `getSolanaEvents` to ensure every sidecar processes the same signature range before constructing the vote extension.

[🔝 back to top](#table-of-contents)

---

## Next Steps

### Investigation Phase

1. **Deploy Enhanced Logging** (Today)
   - Add watermark and event count logging to all testnet sidecars
   - Monitor for 24 hours to confirm desync patterns

2. **Analyze Testnet Logs** (Within 48 hours)
   - Identify which event types desync most frequently
   - Correlate with RPC error rates and rate limit indicators
   - Determine if desync occurs during specific times (high load periods)

3. **Reproduce Locally** (Within 1 week)
   - Set up local testnet with multiple validators
   - Introduce artificial RPC rate limiting to trigger desync
   - Validate that rate limiting is the root cause

### Implementation Phase

4. **Implement Diagnostic Logging** (Week 1)
   - PR with enhanced logging and metrics

5. **Fix Batch Error Handling** (Week 1)
   - PR to process all successful batch responses instead of failing fast

6. **Add Watermark Checkpoints** (Week 2)
   - PR to store and compare watermarks across sidecars

7. **Implement Retry with Jitter** (Week 2)
   - PR to add exponential backoff and jitter to RPC calls

8. **Test in Testnet** (Week 3)
   - Deploy fixes to testnet
   - Monitor for 1 week to confirm desync is resolved

9. **Mainnet Deployment** (Week 4)
   - If testnet is stable, coordinate mainnet upgrade

### Validation Criteria

- **Success Metric**: 99% of vote extensions reach consensus on Solana event hashes
- **Monitoring**: No desync incidents for 7 consecutive days in testnet
- **Rollback Plan**: If desync persists, revert to simpler watermark management (single event type per tick)

[🔝 back to top](#table-of-contents)

---

## Security & Decentralization

### Security Implications

1. **Oracle Consensus Security**: The desync issue weakens oracle consensus by preventing validators from agreeing on off-chain state. This could be exploited by malicious actors to stall the chain or manipulate event processing.

2. **Watermark Rollback Risk**: If watermarks diverge significantly, a validator's sidecar may need to roll back and refetch events, potentially missing time-sensitive operations (e.g., redemptions).

3. **Pending Event Queue Growth**: Divergent pending queues could grow unbounded if reconciliation fails repeatedly, consuming memory and degrading performance.

### Decentralization Posture

- **Validator Independence**: Each validator runs its own sidecar and RPC node, maintaining decentralization. However, if all validators use the same RPC provider, rate limiting becomes a centralized failure point.

- **Recommendation**: Encourage validators to use diverse RPC providers (Helius, QuickNode, self-hosted) to reduce correlated failures.

- **No Trusted Bridge**: Event verification relies on cryptographic proofs and consensus, not a single oracle. Desync does not compromise this property, but delays event processing.

[🔝 back to top](#table-of-contents)

---

## Deployment Considerations

### Testnet Deployment

1. **Coordinate with Validators**: Announce logging upgrade and request validators to capture logs during desync events
2. **Non-Breaking Change**: Logging and metrics enhancements are backward compatible
3. **Monitoring Setup**: Ensure Prometheus/Grafana dashboards capture new metrics

### Mainnet Deployment

1. **Phased Rollout**: Deploy fixes to 1-2 validators first, monitor for 48 hours
2. **Governance Proposal**: If changes affect consensus logic, submit on-chain proposal
3. **Emergency Rollback**: Keep previous sidecar binary available for quick rollback

### Validator Requirements

- **RPC Node Capacity**: Validators may need to upgrade Solana RPC nodes to handle increased traffic
- **Logging Storage**: Enhanced logging may increase disk I/O and storage requirements
- **Alert Configuration**: Set up alerts for watermark divergence and hash mismatch events

[🔝 back to top](#table-of-contents)
