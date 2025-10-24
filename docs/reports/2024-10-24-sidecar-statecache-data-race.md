# Sidecar stateCache Data Race Discovery

**Author**: Peyton-Spencer\
**Date**: October 24, 2024\
**Topics Covered**: Concurrency bugs, Data races, Thread safety, Oracle state management

---

## Table of Contents

- [Executive Summary](#executive-summary)
- [Key Accomplishments](#key-accomplishments)
- [Technical Analysis](#technical-analysis)
- [Code References](#code-references)
- [Challenges & Solutions](#challenges--solutions)
- [Decisions Made](#decisions-made)
- [Security & Decentralization](#security--decentralization)
- [Next Steps](#next-steps)

---

## Executive Summary

During code review, identified a critical data race condition in the sidecar's `Oracle.stateCache` field. The slice is accessed concurrently from multiple goroutines without proper synchronization, creating potential for panics, corrupted data reads, and inconsistent state persistence. The issue exists because only `resetMutex` guards some operations (specifically the scheduled reset), but not all read/write operations on `stateCache`.

**Severity**: High - Can cause runtime panics and state corruption  
**Impact**: Production sidecar stability  
**Status**: Discovered, not yet fixed

[Back to top](#table-of-contents)

---

## Key Accomplishments

- ✅ Identified data race in `Oracle.stateCache` concurrent access
- ✅ Mapped all read and write locations across the codebase
- ✅ Documented specific race scenarios with detailed analysis
- ✅ Identified the root cause (inconsistent mutex usage)
- ✅ Scoped the required fix (extend `resetMutex` to all stateCache operations)

[Back to top](#table-of-contents)

---

## Technical Analysis

### The Data Race

`Oracle.stateCache` is a `[]sidecartypes.OracleState` slice that is accessed concurrently from multiple goroutines without consistent mutex protection. While `Oracle.currentState` is properly protected using `atomic.Value`, the `stateCache` slice backing it is not.

### Unprotected Write Locations

1. **`CacheState()`** ([`sidecar/utils.go:98-114`](../../sidecar/utils.go))
   - Appends new state to `stateCache`
   - Calls `SaveToFile()` to persist state
   - Called from `applyStateUpdate()` in main oracle loop
   - **No mutex protection**

2. **`SaveToFile()`** ([`sidecar/utils.go:68-96`](../../sidecar/utils.go))
   - Reads `stateCache` to encode as JSON
   - Creates temporary file and atomically renames
   - **No mutex protection**

3. **`SetStateCacheForTesting()`** ([`sidecar/utils.go:211-220`](../../sidecar/utils.go))
   - Completely replaces `stateCache` slice
   - Used in test setup
   - **No mutex protection**

### Unprotected Read Locations

1. **`getStateByEthHeight()`** ([`sidecar/utils.go:117-125`](../../sidecar/utils.go))
   - Iterates over `stateCache` in reverse
   - Even has a TODO comment acknowledging the race: `// TODO: possible data race with stateCache -- concurrent read and write to stateCache?`
   - **No mutex protection**

2. **`SaveToFile()`** (same function as above)
   - Reads `stateCache` to encode
   - **No mutex protection**

### Protected Operations

1. **`performFullStateResetLocked()`** ([`sidecar/oracle.go:3932-3955`](../../sidecar/oracle.go))
   - Clears `stateCache` to empty state
   - **Protected by `resetMutex`** (when called from `maybePerformScheduledReset`)
   - However, the mutex is only held during the reset, not during other operations

### Race Scenarios

#### Scenario 1: Concurrent Append and Reset

```
Thread 1 (Main Loop)          Thread 2 (Scheduled Reset)
--------------------          --------------------------
CacheState()
  stateCache = append(...)    
                              performFullStateResetLocked()
                                resetMutex.Lock()
                                stateCache = []...{EmptyState}
                                resetMutex.Unlock()
  SaveToFile(stateCache)      
```

**Problem**: Thread 1's append happens without mutex, while Thread 2 clears the slice. The append could:
- Operate on stale slice capacity/length
- Result in lost state update
- Cause slice corruption if append triggers reallocation

#### Scenario 2: Iterator During Reset

```
Thread 1 (gRPC Handler)       Thread 2 (Scheduled Reset)
-----------------------       --------------------------
getStateByEthHeight(height)
  for i := len(stateCache)-1...
                              performFullStateResetLocked()
                                resetMutex.Lock()
                                stateCache = []...{EmptyState}
                                resetMutex.Unlock()
    if stateCache[i].Eth...   <- PANIC: index out of range
```

**Problem**: Iterator relies on slice length, but slice is replaced mid-iteration. Classic concurrent modification panic.

#### Scenario 3: Concurrent SaveToFile Operations

```
Thread 1 (Main Loop)          Thread 2 (Main Loop - different tick)
--------------------          ------------------------------------
applyStateUpdate()
  CacheState()
    stateCache = append(...)
                              applyStateUpdate()
                                CacheState()
                                  stateCache = append(...)
                                  SaveToFile()
    SaveToFile()                    json.Encode(stateCache)
      json.Encode(stateCache)
```

**Problem**: Multiple concurrent appends to the same slice without synchronization. Slice internals (length, capacity, backing array pointer) could be corrupted.

#### Scenario 4: Encode During Modification

```
Thread 1 (Main Loop)          Thread 2 (Scheduled Reset)
--------------------          --------------------------
CacheState()
  SaveToFile()
    json.Encode(stateCache)
      for _, state := range stateCache
                              performFullStateResetLocked()
                                resetMutex.Lock()
                                stateCache = []...{EmptyState}
                                resetMutex.Unlock()
        encode state...       <- Encoding partially old, partially new state
```

**Problem**: JSON encoder reads slice while it's being replaced. Could result in:
- Partial state written to disk
- Encoding errors
- Corrupted state file

### Why This Wasn't Caught Earlier

1. **Race detector not run in production**: Data races are non-deterministic and may only surface under specific timing conditions
2. **Low frequency of scheduled resets**: `maybePerformScheduledReset()` only runs every 24 hours (or 2 minutes in test mode)
3. **Most operations are reads**: Read-read races don't cause corruption (though they may see inconsistent state)
4. **Atomic currentState masks the issue**: Since `currentState` is properly atomic, most code uses that instead of `stateCache` directly

### Correct Fix

Extend `resetMutex` to guard **all** operations on `stateCache`:

1. **Acquire `resetMutex` in `CacheState()`** before appending and calling `SaveToFile()`
2. **Acquire `resetMutex` in `getStateByEthHeight()`** before iterating
3. **Acquire `resetMutex` in `SetStateCacheForTesting()`** before replacing slice
4. **Keep existing `resetMutex` in `performFullStateResetLocked()`**

Alternative (more surgical): Consider whether `stateCache` needs to be a field at all, or if it could be refactored to only live within the critical section of `SaveToFile()`. However, `getStateByEthHeight()` needs historical access, so the field is necessary.

[Back to top](#table-of-contents)

---

## Code References

### Files Analyzed

- [`sidecar/oracle.go`](../../sidecar/oracle.go) - Oracle struct definition, scheduled reset logic
- [`sidecar/utils.go`](../../sidecar/utils.go) - State management functions (CacheState, SaveToFile, getStateByEthHeight)
- [`sidecar/types.go`](../../sidecar/types.go) - Type definitions
- [`sidecar/main.go`](../../sidecar/main.go) - Initialization and startup

### Files Modified

- [`sidecar/oracle2_test.go`](../../sidecar/oracle2_test.go) - Added race condition tests (210 lines added)

### Specific Functions Affected

- `CacheState()` - Needs mutex protection
- `SaveToFile()` - Needs mutex protection
- `getStateByEthHeight()` - Needs mutex protection (has TODO comment acknowledging this)
- `SetStateCacheForTesting()` - Needs mutex protection
- `performFullStateResetLocked()` - Already has mutex protection (caller holds lock)
- `maybePerformScheduledReset()` - Holds mutex correctly, but scope is too narrow

[Back to top](#table-of-contents)

---

## Challenges & Solutions

### Challenge 1: Identifying the Root Cause

**Problem**: Initial suspicion was vague ("concurrent read/write possibility"). Needed to trace all access patterns to confirm.

**Solution**: Systematic grep for `stateCache` across codebase, then manual analysis of each usage site to determine:
- Is this a read or write?
- Which goroutine context does this run in?
- Is there mutex protection?
- What are the timing relationships?

### Challenge 2: Understanding the Mutex Strategy

**Problem**: Code already has `resetMutex`, but it's not consistently applied. Why was it added only for resets?

**Solution**: Recognized that `resetMutex` was likely added later to fix scheduled reset races, but the original design oversight (unprotected stateCache) wasn't fully addressed. The fix is to extend the existing mutex to all operations, not add a new one.

### Challenge 3: Deciding Between Fix Approaches

**Problem**: Multiple potential fixes:
1. Add mutex to all operations
2. Replace slice with sync.Map or other concurrent data structure
3. Redesign to avoid shared mutable state

**Solution**: Option 1 (extend existing mutex) is most appropriate because:
- Minimal code change
- Consistent with existing pattern
- `stateCache` isn't a performance bottleneck (infrequent updates)
- Preserves existing API and structure

[Back to top](#table-of-contents)

---

## Decisions Made

### Decision 1: Report Issue, Don't Fix Immediately

**Rationale**: Creating a comprehensive report allows for:
- Team review and discussion of fix approach
- Consideration of whether other related issues exist
- Proper testing and validation before merge
- Documentation for future reference

**Trade-off**: Issue remains unfixed temporarily, but risk is relatively low given infrequent reset schedule.

### Decision 2: Use Existing resetMutex

**Rationale**: Adding a new mutex (e.g., `stateCacheMutex`) would:
- Increase cognitive complexity
- Create potential for lock ordering issues
- Duplicate functionality

Better to extend the scope of the existing `resetMutex` to protect all state cache operations.

### Decision 3: Mutex Over Refactoring

**Rationale**: While a redesign could eliminate the shared mutable state entirely, this would be a much larger change:
- Higher risk of introducing new bugs
- More extensive testing required
- Unclear if benefits outweigh costs

The mutex fix is low-risk, well-understood, and sufficient.

[Back to top](#table-of-contents)

---

## Security & Decentralization

### Security Implications

#### Impact of Data Race on Oracle Security

**Corrupted State Persistence**:
- If `SaveToFile()` encodes corrupted state, the oracle could restart with invalid historical data
- Could affect `getStateByEthHeight()` lookups used by gRPC endpoints
- **Mitigation after fix**: Proper mutex ensures atomic state transitions

**Panic in Production**:
- Index-out-of-bounds panic in `getStateByEthHeight()` could crash the sidecar
- Sidecar crash prevents validator from submitting vote extensions
- Missing vote extensions could affect consensus quorum (though not individually fatal)
- **Mitigation after fix**: Race-free iteration eliminates panic risk

**State Inconsistency**:
- Concurrent modification could lead to inconsistent view of historical states
- gRPC clients querying historical state could receive incorrect data
- **Mitigation after fix**: Mutex provides consistent view of stateCache

### Impact on Decentralization

**Validator Reliability**:
- Data races affecting individual validators reduce overall network reliability
- If multiple validators experience the same race condition simultaneously, consensus could be impacted
- **Importance**: Sidecar stability is critical to validator participation in oracle consensus

**No Direct Security Exploit**:
- This is a stability bug, not an exploitable vulnerability
- Does not affect cryptographic security or consensus rules
- Does not allow unauthorized state modification (beyond what race conditions could cause)

### Testing with Race Detector

**Critical**: This issue should be caught by Go's race detector:

```bash
go test -race ./sidecar/...
```

**Recommendation**: Add race detector to CI pipeline to prevent similar issues:
- Run integration tests with `-race` flag
- Accept performance overhead in CI (not production)
- Catch concurrency bugs before merge

[Back to top](#table-of-contents)

---

## Next Steps

### Immediate Actions

1. **Run Race Detector**: Execute `go test -race -run TestStateCacheDataRace ./sidecar/` to confirm the race
2. **Implement Fix**: Add mutex protection to all stateCache operations
3. **Verify Fix**: Re-run race tests after applying mutex fix to confirm resolution
4. **Code Review**: Get second set of eyes on the fix and broader concurrency patterns

### Testing Strategy

1. **Race Detector Tests**: 
   - **NEW**: Added `TestStateCacheDataRace` in `sidecar/oracle2_test.go` (4 scenarios)
   - **NEW**: Added `TestStateCacheRaceWithSetStateCacheForTesting` 
   - Run with: `go test -race -run TestStateCacheDataRace ./sidecar/`
   - Will detect races in: CacheState, getStateByEthHeight, SaveToFile, SetStateCacheForTesting
   - Verify no races detected after mutex fix is applied

2. **Integration Testing**:
   - Test with `--test-reset` flag (2-minute reset interval) to exercise reset path
   - Concurrent gRPC requests calling `getStateByEthHeight()` during reset
   - Load testing with multiple ticks appending to stateCache

3. **Regression Testing**:
   - Verify state persistence still works correctly
   - Confirm historical state queries return correct data
   - Check scheduled reset behavior unchanged

### Longer-Term Improvements

1. **Enable Race Detector in CI**:
   - Add `-race` flag to test commands in GitHub Actions
   - Consider running race detector nightly even if not on every PR (performance overhead)

2. **Concurrency Audit**:
   - Review other Oracle fields for similar issues
   - Document which fields require which locks
   - Consider adding lock-checking assertions in debug builds

3. **Observability**:
   - Add metrics for stateCache size
   - Monitor for any mutex contention after fix (shouldn't be significant)
   - Alert on sidecar crashes (existing, but ensure it's working)

4. **Documentation**:
   - Add comments to Oracle struct documenting which mutex protects which fields
   - Document goroutine safety in function comments
   - Update contribution guide with concurrency testing requirements

### Related Questions to Investigate

- Are there other slices or maps in Oracle that need protection?
- Is `transactionCache` properly protected? (It has its own mutex, good)
- Could we simplify the state management to reduce concurrency complexity?
- Should `getStateByEthHeight()` even exist, or can we redesign to avoid needing it?

[Back to top](#table-of-contents)

---

## Appendix A: Test Coverage

### Race Condition Tests Added

Two comprehensive test functions were added to `sidecar/oracle2_test.go`:

#### 1. TestStateCacheDataRace (Lines 1606-1780)

Tests 4 concurrent access scenarios:

**Scenario 1: Concurrent CacheState()**
- 10 goroutines, 100 iterations each
- All append to `stateCache` simultaneously
- Reproduces: Concurrent slice append race

**Scenario 2: Concurrent Read and Write**
- 10 reader goroutines calling `getStateByEthHeight()`
- 10 writer goroutines calling `CacheState()`
- 100 iterations each
- Reproduces: Read-write race during iteration

**Scenario 3: Concurrent SaveToFile()**
- 10 goroutines reading `stateCache` for JSON encoding
- 10 goroutines modifying `stateCache` via `CacheState()`
- Reproduces: Encode-during-modification race

**Scenario 4: Simulated Reset During Operations**
- 10 goroutines iterating over `stateCache` 
- 5 goroutines replacing `stateCache` (simulating reset)
- Reproduces: Iterator panic when slice is replaced mid-loop
- **Most dangerous scenario** - can cause production crashes

#### 2. TestStateCacheRaceWithSetStateCacheForTesting (Lines 1782-1815)

- 20 goroutines concurrently calling `SetStateCacheForTesting()`
- 50 iterations per goroutine
- Tests race in test helper function

### Running the Tests

```bash
# Run just the race tests
go test -race -run TestStateCacheDataRace ./sidecar/

# Run all sidecar tests with race detector (takes longer)
go test -race ./sidecar/...

# Expected output BEFORE fix:
# ==================
# WARNING: DATA RACE
# Write at 0x... by goroutine X:
#   main.(*Oracle).CacheState()
#       /path/to/sidecar/utils.go:105
# 
# Previous read at 0x... by goroutine Y:
#   main.(*Oracle).getStateByEthHeight()
#       /path/to/sidecar/utils.go:120
# ==================

# Expected output AFTER fix:
# PASS
# ok      github.com/Zenrock-Foundation/zrchain/v6/sidecar    2.456s
```

### Test Files Organization

After investigation, the distinction between test files:

- **`oracle_test.go`** (260 lines)
  - Integration-style tests
  - Some make real network calls (marked with `t.Skip()`)
  - Tests full oracle initialization and external dependencies
  - Example: `TestFetchSolanaBurnEvents_Integration`

- **`oracle2_test.go`** (1815 lines after additions)
  - Comprehensive unit tests
  - Extensive mocking and helper functions
  - Isolated component testing
  - Includes the new race condition tests
  - Better for testing specific concurrent behaviors

**Recommendation**: Consider merging these files in the future or establishing clear naming convention (e.g., `oracle_integration_test.go` vs `oracle_unit_test.go`).

[Back to top](#table-of-contents)

---

## Appendix B: Code Snippets

### Current (Buggy) Code

```go
// sidecar/utils.go:98-114
func (o *Oracle) CacheState() {
	currentState := o.currentState.Load().(*sidecartypes.OracleState)
	newState := *currentState // Create a copy of the current state

	// Cache the new state
	o.stateCache = append(o.stateCache, newState) // ⚠️ RACE: Unprotected append
	if len(o.stateCache) > sidecartypes.OracleCacheSize {
		o.stateCache = o.stateCache[1:] // ⚠️ RACE: Unprotected slice modification
	}

	if err := o.SaveToFile(o.Config.StateFile); err != nil {
		log.Printf("Error saving state to file: %v", err)
	}
}
```

```go
// sidecar/utils.go:117-125
func (o *Oracle) getStateByEthHeight(height uint64) (*sidecartypes.OracleState, error) {
	// TODO: possible data race with stateCache -- concurrent read and write to stateCache?
	// Search in reverse order to efficiently find the most recent state with matching height
	for i := len(o.stateCache) - 1; i >= 0; i-- { // ⚠️ RACE: Unprotected iteration
		if o.stateCache[i].EthBlockHeight == height {
			return &o.stateCache[i], nil
		}
	}
	return nil, fmt.Errorf("state with Ethereum block height %d not found", height)
}
```

### Proposed Fix (Conceptual)

```go
func (o *Oracle) CacheState() {
	currentState := o.currentState.Load().(*sidecartypes.OracleState)
	newState := *currentState

	o.resetMutex.Lock() // ✅ FIX: Acquire mutex before modifying stateCache
	o.stateCache = append(o.stateCache, newState)
	if len(o.stateCache) > sidecartypes.OracleCacheSize {
		o.stateCache = o.stateCache[1:]
	}
	
	if err := o.SaveToFile(o.Config.StateFile); err != nil {
		log.Printf("Error saving state to file: %v", err)
	}
	o.resetMutex.Unlock() // ✅ FIX: Release after SaveToFile completes
}
```

```go
func (o *Oracle) getStateByEthHeight(height uint64) (*sidecartypes.OracleState, error) {
	o.resetMutex.Lock() // ✅ FIX: Acquire mutex before reading stateCache
	defer o.resetMutex.Unlock()
	
	for i := len(o.stateCache) - 1; i >= 0; i-- {
		if o.stateCache[i].EthBlockHeight == height {
			// Make a copy to return (safe after mutex release)
			state := o.stateCache[i]
			return &state, nil
		}
	}
	return nil, fmt.Errorf("state with Ethereum block height %d not found", height)
}
```

[Back to top](#table-of-contents)

