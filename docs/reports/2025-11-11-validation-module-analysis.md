# zrChain Validation Module: Vote Extension & Consensus Analysis

## Executive Summary

This document provides a detailed analysis of how vote extensions are constructed, how Solana event hashes are computed and validated, and how consensus is reached in the zrChain validation module.

---

## Sidecar Oracle Event Flow (Context for Vote Extensions)

### Event Fetching Architecture

The sidecar oracle constructs the data that gets hashed in vote extensions through a 60-second tick cycle:

**Main Tick Flow** (`runOracleMainLoop()` - line 211):
```
Tick fires every 60 seconds
    ↓
processSolanaMintEvents() - Parallel fetch:
    ├─ getSolROCKMints()      - ROCK mint events
    ├─ getSolZenBTCMints()    - zenBTC mint events
    └─ getSolZenZECMints()    - zenZEC mint events via EventStore
    ↓
fetchSolanaBurnEvents() - Parallel fetch:
    ├─ getSolanaZenBTCBurnEvents()
    ├─ getSolanaZenZECBurnEvents()
    └─ getSolanaRockBurnEvents()
    ↓
fetchEthereumBurnEvents() - Ethereum burn events
    ↓
buildFinalState() - Merge & sort all events
    ↓
OracleState created with:
    - SolanaMintEvents[]
    - SolanaBurnEvents[]
    - EthBurnEvents[]
    - All other data fields
    ↓
(Every 60 seconds) Vote extension construction uses this state
```

### Critical Detail: Event List Order Matters

When validators construct vote extensions, they hash the complete event lists:

```go
// x/validation/keeper/abci.go: ConstructVoteExtension
SolanaMintEventsHash = deriveHash(oracleData.SolanaMintEvents)   // HASH includes list order
SolanaBurnEventsHash = deriveHash(oracleData.SolanaBurnEvents)   // HASH includes list order
```

**This means**:
- Same events in different order = different hash
- Different hashes = no consensus
- No consensus = events don't get processed

---

## 1. Vote Extension Construction Flow

### 1.1 ConstructVoteExtension Function

**Location:** `x/validation/keeper/abci.go:158-255`

**Purpose:** Builds a vote extension from oracle data that each validator proposes during vote extension phase.

**Process Flow:**

```
Oracle Data Collection
    ↓
gatherOracleDataForVoteExtension() - Fetches:
    - Block headers (Bitcoin, ZCash)
    - Ethereum nonces and gas parameters
    - Solana nonces and accounts
    - Solana mint/burn events
    - Ethereum burn events
    - Redemptions
    - Price feeds
    ↓
Hash Derivation for Each Field:
    deriveHash(data) → SHA256(JSON(data))
    ↓
VoteExtension Structure Assembly
    (25 fields with type-safe identifiers)
```

### 1.2 Hash Derivation Function

**Location:** `x/validation/keeper/abci_utils.go:333-339`

```go
func deriveHash[T any](data T) ([32]byte, error) {
    dataBz, err := json.Marshal(data)
    if err != nil {
        return [32]byte{}, fmt.Errorf("error encoding data: %w", err)
    }
    return sha256.Sum256(dataBz), nil
}
```

**Critical Details:**
- Uses JSON marshaling for data serialization
- JSON ordering and field representation must be deterministic across validators
- Returns 32-byte SHA256 hash
- All validators MUST produce identical hashes for same data

**Hash Fields in Vote Extension:**
1. `EigenDelegationsHash` - AVS delegation state
2. `EthBurnEventsHash` - Ethereum burn events
3. `RedemptionsHash` - Redemption transactions
4. `SolanaMintNoncesHash` - Solana nonce accounts
5. `SolanaAccountsHash` - Associated token accounts
6. **`SolanaMintEventsHash`** - Solana mint events (CRITICAL)
7. `SolanaBurnEventsHash` - Solana burn events
8. `LatestBtcHeaderHash` - Latest Bitcoin header
9. `RequestedBtcHeaderHash` - Requested historical Bitcoin header
10. `LatestZcashHeaderHash` - Latest ZCash header
11. `RequestedZcashHeaderHash` - Requested ZCash header

---

## 2. Solana Mint Event Hash Computation & Consensus

### 2.1 How Solana Mint Events Reach Vote Extensions

**Data Flow:**

```
Sidecar Oracle
    ↓
GetSidecarState() → Fetches SolanaMintEvents
    ↓
oracleData.SolanaMintEvents = resp.SolanaMintEvents
    ↓
ConstructVoteExtension():
    solanaMintEventsHash = deriveHash(oracleData.SolanaMintEvents)
    ↓
VoteExtension.SolanaMintEventsHash = solanaMintEventsHash[:]
    ↓
Marshal to JSON → ExtendVote (sent to consensus)
```

### 2.2 Event Hash Mismatch Scenarios

**Cause 1: Sidecar Data Inconsistency**
- Different validators have sidecars at different block heights
- Sidecar not fully synced for all validators
- Different Solana RPC responses

**Cause 2: Event Ordering**
- JSON marshaling of slice order matters
- If validators fetch events in different order → different hash
- Sidecar implementation may not sort events consistently

**Cause 3: Event Filtering Logic**
- Validators may filter events differently
- Missing event detection logic differences
- Timeout handling variations

### 2.3 Event Structure

**SolanaMintEvent Fields:**
```go
type SolanaMintEvent struct {
    Coint           Coin      // ZENBTC, ZENZEC, ROCK, etc.
    TxSig           string    // Transaction signature
    Recipient       []byte    // Recipient public key
    Mint            []byte    // Mint address
    Value           uint64    // Token amount
    LogIndex        uint32    // Event index in transaction
    SigHash         []byte    // Event signature hash (used for deduplication)
}
```

---

## 3. Consensus Mechanism

### 3.1 GetConsensusAndPluralityVEData Function

**Location:** `x/validation/keeper/abci_utils.go:129-260`

**Key Innovation:** Field-by-field consensus instead of all-or-nothing

**Voting Algorithm:**

```
Vote Processing:
    For each validator in ExtendedCommitInfo.Votes:
        ↓
    Parse vote extension JSON
    ↓
    For each field (25 total):
        - Get field value from vote extension
        - Create key via genericGetKey(value) → deterministic string
        - Track field_name → {value: any, votePower: int64}
    ↓
    Calculate thresholds:
        - superMajority = ((totalVotePower * 2) / 3) + 1
        - simpleMajority = (totalVotePower / 2) + 1
    ↓
    For each field:
        - Find max vote power for field
        - Collect all values with max power (handle ties)
        - Tie-breaking: lexicographic sort on serialized representation
        - If max power >= supermajority: consensus reached
        - If gas field AND max power >= simpleMajority: consensus reached
```

### 3.2 Consensus Thresholds

**Supermajority (2/3+):**
- All hash fields
- All header fields
- Ethereum parameters
- Nonces
- Prices
- Redemptions

**Simple Majority (>50%):**
- Only gas-related fields:
  - `VEFieldEthGasLimit`
  - `VEFieldEthBaseFee`
  - `VEFieldEthTipCap`

**No Consensus Required:**
- `SidecarVersionName` (informational only)

### 3.3 Tie-Breaking Mechanism

**Location:** `abci_utils.go` lines 221-237

**Problem:** Multiple validators vote for different values with same voting power

**Solution:** Deterministic tie-breaking based on lexicographic ordering

```go
// Sort tied values by their string representation
slices.SortFunc(tiedValues, func(a, b struct {
    key   string
    value any
}) int {
    return strings.Compare(a.key, b.key)
})

// Select first (lowest lexicographically)
mostVotedValue := tiedValues[0].value
```

**Guarantee:** All validators will deterministically select the same value

### 3.4 Required Voting Validators

**Consensus Requirement:**
- Need supermajority (2/3+) of voting power
- Example: 3 validators with 100 power each
  - Total: 300
  - Required: ((300*2)/3)+1 = 201 power
  - Any 3 validators or 2 of largest validators = consensus

**No Single Point Failure:** As long as 2/3+ validators are online and agree

---

## 4. Hash Validation & Mismatch Detection

### 4.1 Validation Flow

**Location:** `x/validation/keeper/abci_utils.go:1487-1649`

**Three-Stage Validation:**

```
Stage 1: Gather Oracle Data (PreBlocker)
    ↓
    Proposer calls GetConsensusAndPluralityVEData()
    → Identifies which fields have supermajority consensus
    → Returns: consensusVE, pluralityVE, fieldVotePowers map
    
Stage 2: Validate Oracle Data Against Consensus
    ↓
    validateOracleData(consensusVE, oracleData, fieldVotePowers)
    → For each field with consensus:
        → If hash field: deriveHash(data) == voteExt.hash?
        → If scalar: voteExt.value == oracleData.value?
    → Collect mismatches
    
Stage 3: Handle Mismatches
    ↓
    handleValidationMismatches(mismatches, fieldVotePowers)
    → Delete field from fieldVotePowers (revoke consensus)
    → Log warning with expected vs actual
    → Continue block processing with remaining consensus fields
```

### 4.2 Hash Mismatch Scenarios

**Scenario 1: SolanaMintEventsHash Mismatch**

```
Consensus reached: SolanaMintEventsHash = 0xABCD...
Oracle data fetched: oracleData.SolanaMintEvents = [...]

Validation:
    expectedHash := voteExt.SolanaMintEventsHash  // 0xABCD...
    actualHash := deriveHash(oracleData.SolanaMintEvents)  // 0xEF01...
    
    if !bytes.Equal(expectedHash, actualHash):
        recordMismatch(VEFieldSolanaMintEventsHash, expectedHash, actualHash)
        delete(fieldVotePowers[VEFieldSolanaMintEventsHash])
        
Result: SolanaMintEventsHash consensus is REVOKED
        → processSolanaMintEvents() checks fieldHasConsensus()
        → If no consensus: events are NOT processed
```

### 4.3 Detection Mechanism

**Location:** `abci_utils.go` lines 1539-1542

```go
if fieldHasConsensus(fieldVotePowers, VEFieldSolanaMintEventsHash) {
    if err := validateHashField(
        VEFieldSolanaMintEventsHash.String(),
        voteExt.SolanaMintEventsHash,
        oracleData.SolanaMintEvents) {
        
        recordMismatch(VEFieldSolanaMintEventsHash, ...)
    }
}
```

---

## 5. Event Processing & Consensus Dependency

### 5.1 Solana Mint Event Processing

**Location:** `x/validation/keeper/abci.go:414-418` (PreBlocker)

**Consensus Dependency:**

```go
if fieldHasConsensus(oracleData.FieldVotePowers, VEFieldSolanaMintEventsHash) {
    k.processSolanaZenBTCMintEvents(ctx, oracleData)
    k.processSolanaROCKMintEvents(ctx, oracleData)
    k.processSolanaDCTMintEvents(ctx, oracleData)
}
```

**If Hash Mismatch Occurs:**
- Consensus revoked
- Events NOT processed that block
- Pending mint transactions remain in DEPOSITED status
- Next block: retry with latest oracle data

### 5.2 Event Matching Logic

**Event Matching Criteria:**
```go
// From abci_dct.go:264-418
for _, event := range oracleData.SolanaMintEvents {
    // Check 1: Coin type matches
    if event.Coint != coin {
        continue
    }
    
    // Check 2: Recipient matches
    if event.Recipient != pendingMint.RecipientAddress:
        continue
    }
    
    // Check 3: Amount matches
    if event.Value != pendingMint.Amount {
        continue
    }
    
    // Check 4: Event ID matches expected counter
    if eventID != counters.MintCounter + 1 {
        continue
    }
    
    matchedEvent = &event
    break
}
```

### 5.3 Event Deduplication

**Using Event Hash:**

```go
eventHash := base64.StdEncoding.EncodeToString(event.SigHash)
eventKey := collections.Join(asset.String(), eventHash)

if alreadyProcessed, err := k.ProcessedSolanaMintEvents.Get(ctx, eventKey); 
   err == nil && alreadyProcessed {
    // Already processed, skip
    continue
}
```

---

## 6. Known Issues & Race Conditions

### 6.1 JSON Marshaling Consistency

**Issue:** JSON field ordering may differ across validators if:
- Go version differs
- Protobuf generation differs
- External library versions differ

**Mitigation:** Golang maps in JSON maintain insertion order in Go 1.18+

**Risk Level:** MEDIUM
- All validators use Go 1.25+ (per CLAUDE.md)
- If any dependencies changed: all validators must rebuild

### 6.2 Sidecar Sync Lag

**Issue:** Solana RPC nodes may not have latest events
```
Block N:
  Validator A: Fetches events up to Solana block 200,000
  Validator B: Fetches events up to Solana block 199,999
  → Different SolanaMintEventsHash
  → No consensus or mismatch detection
```

**Current Handling:**
- Relying on sidecar RPC reliability
- No explicit retry mechanism in consensus layer
- PreBlocker doesn't retry failed hash validations

**Risk Level:** HIGH
- Only mitigated by validator operational discipline

### 6.3 Event Ordering Non-Determinism

**Issue:** If sidecar returns events in non-deterministic order:
```
Validator A: [EventX, EventY, EventZ] → hash1
Validator B: [EventY, EventZ, EventX] → hash2
```

**Likelihood:** LOW
- Sidecar likely returns events by timestamp/height

**Risk Level:** MEDIUM

### 6.4 Hash Validation Idempotency

**Design:** Events use SigHash for deduplication, preventing double-processing

**Risk Level:** LOW (defensive programming)

---

## 7. Vote Extension Field Consensus Summary

| Field | Hash? | Threshold | Comment |
|-------|-------|-----------|---------|
| EigenDelegationsHash | Yes | 2/3+ | AVS state |
| EthBurnEventsHash | Yes | 2/3+ | Ethereum burns |
| SolanaBurnEventsHash | Yes | 2/3+ | Solana burns |
| **SolanaMintEventsHash** | Yes | 2/3+ | **CRITICAL** |
| SolanaMintNoncesHash | Yes | 2/3+ | Nonce accounts |
| SolanaAccountsHash | Yes | 2/3+ | Token accounts |
| LatestBtcHeaderHash | Yes | 2/3+ | Bitcoin header |
| LatestZcashHeaderHash | Yes | 2/3+ | ZCash header |
| EthBlockHeight | No | 2/3+ | Ethereum height |
| EthGasLimit | No | >50% | Gas parameter |
| EthBaseFee | No | >50% | Gas parameter |
| EthTipCap | No | >50% | Gas parameter |
| ROCKUSDPrice | No | 2/3+ | Price feed |
| BTCUSDPrice | No | 2/3+ | Price feed |
| ZECUSDPrice | No | 2/3+ | Price feed |
| Various Nonces | No | 2/3+ | Nonce values |

---

## 8. How Many Validators Need to Agree

**For Solana Mint Events Processing to Execute:**

Required: **2/3+ (supermajority)** of total validator voting power

**Example with 3 equal validators (100 power each):**
- Total: 300 power
- Required: ((300 × 2) ÷ 3) + 1 = 201 power
- This means: any 2 validators (200 power) = NOT enough
- All 3 validators = enough (300 power)

**Example with 4 validators (25, 25, 25, 25 power):**
- Total: 100 power
- Required: ((100 × 2) ÷ 3) + 1 = 67 power
- This means: 3 validators (75 power) = enough
- But: 2 validators (50 power) = NOT enough

**Key Insight:** Consensus is vote-power-weighted, not validator-count-based

---

## 9. Recent Code Changes

**Key Recent Commits:**
1. **429d1871** - Use event store for all flows (v6rev57)
2. **6d63d545** - ZCash redemptions support
3. **295616d9** - Enable ZenZEC redemptions
4. **c119a57c** - Increase VE size limit to 1600 bytes

---

## 10. Files Referenced

- `x/validation/keeper/abci.go` - Vote extension construction and processing
- `x/validation/keeper/abci_types.go` - Type definitions
- `x/validation/keeper/abci_utils.go` - Consensus and validation logic
- `x/validation/keeper/abci_dct.go` - DCT (ZenZEC) event processing
- `x/validation/keeper/abci_zenbtc.go` - ZenBTC event processing

