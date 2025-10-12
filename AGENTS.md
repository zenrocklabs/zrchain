# ZenZEC (ZCash) Integration Agent Plan

## Current Status
The DCT zenZEC minting flow is implemented in the chain code but **NOT fully wired up**. The critical missing piece is ZCash block header tracking in the sidecar oracle.

## Critical Gap
When a user deposits ZEC to mint zenZEC, the chain needs to verify the deposit transaction is included in a ZCash block. Currently:
- ✅ Chain code expects block headers via `k.validationKeeper.BtcBlockHeaders.Get(ctx, msg.BlockHeight)`
- ✅ Chain code can verify ZCash transactions (same format as Bitcoin)
- ❌ **Sidecar is NOT fetching ZCash block headers**
- ❌ **Vote extensions are NOT including ZCash headers**

## Architecture Decision
**NOT using SPV node** - Instead using RPC endpoint to fetch ZCash headers on-demand.

## Implementation Plan for ZCash Block Header Tracking

### Phase 1: Sidecar Configuration & RPC Client ✓ COMPLETE

**Goal**: Add ZCash RPC client to sidecar for fetching block headers

**Tasks**:
1. ✅ Create `sidecar/zcash_client.go` with RPC methods:
   - `NewZcashClient(rpcURL string, enabled bool) *ZcashClient`
   - `GetBlockCount(ctx) (int64, error)`
   - `GetBlockHash(ctx, height) (string, error)`
   - `GetBlockHeader(ctx, hash) (*api.BTCBlockHeader, error)`
   - `GetBlockHeaderByHeight(ctx, height) (*api.BTCBlockHeader, error)`
   - `GetLatestBlockHeader(ctx) (*api.BTCBlockHeader, int64, error)`

2. ✅ Update `sidecar/types.go`:
   - Add `zcashClient *ZcashClient` field to Oracle struct
   - Add `lastZcashHeaderHeight int64` tracking field

3. ✅ Update `sidecar/shared/types.go`:
   - Add ZCash header fields to `OracleState`:
     - `LatestZcashBlockHeight int64`
     - `LatestZcashBlockHeader *api.BTCBlockHeader`
     - `RequestedZcashBlockHeight int64`
     - `RequestedZcashBlockHeader *api.BTCBlockHeader`
   - Add `ZcashRPC map[string]string` to Config struct

4. ✅ Update `sidecar/config.yaml.example`:
   - Add `zcash_rpc` section with devnet/testnet/mainnet endpoints

5. ✅ Update `sidecar/main.go`:
   - Add `validateZcashClient()` function (similar to validateSolanaClient)
   - Initialize ZCash client and pass to NewOracle()

6. ✅ Update `sidecar/oracle.go`:
   - Add `zcashClient *ZcashClient` parameter to NewOracle()
   - Initialize oracle.zcashClient field

7. ✅ Update `sidecar/server.go`:
   - Add gRPC methods: `GetZcashBlockHeaderByHeight()`, `GetLatestZcashBlockHeader()`

**Success Criteria**: ✅
- Sidecar successfully connects to ZCash RPC endpoint
- Can fetch latest ZCash block header via RPC
- ZCash client properly initialized in Oracle struct

### Phase 2: Vote Extension Integration ✓ COMPLETE

**Goal**: Include ZCash headers in vote extensions for validator consensus

**Tasks**:
1. ✅ Update `x/validation/keeper/abci_types.go`:
   - Add ZCash fields to `VoteExtension` struct (4 new fields)
   - Add ZCash fields to `OracleData` struct (4 new fields)
   - Add ZCash field enum constants
   - Add ZCash field handlers to `initializeFieldHandlers()`
   - Add ZCash methods to `sidecarClient` interface

2. ✅ Update `proto/zrchain/dct/mint.proto`:
   - Add `RequestedZcashHeaders` message

3. ✅ Update `x/validation/keeper/keeper.go`:
   - Add `RequestedHistoricalZcashHeaders` collection

4. ✅ Update `x/validation/types/keys.go`:
   - Add `RequestedHistoricalZcashHeadersKey` and index

5. ✅ Update `x/validation/keeper/abci_utils.go`:
   - Add `retrieveZcashHeaders()` function

6. ✅ Update `x/validation/keeper/abci.go`:
   - Update `gatherOracleDataForVoteExtension()` to fetch ZCash headers
   - Update `ConstructVoteExtension()` to hash and include ZCash headers

7. ✅ Update `sidecar/proto/api/sidecar_service.proto`:
   - Add ZCash gRPC methods

8. ✅ Regenerate protobuf files (chain + sidecar)

**Success Criteria**: ✅
- Vote extensions include ZCash header hashes
- Validators can reach consensus on ZCash block headers
- ZCash headers transmitted efficiently in vote extensions

### Phase 3: Chain State Storage (IN PROGRESS)

**Goal**: Store ZCash headers that reach consensus in chain state

**Tasks**:
1. Add `ZcashBlockHeaders` collection to validation keeper (similar to `BtcBlockHeaders`)
2. Add `LatestZcashHeaderHeight` tracking
3. Update PreBlocker to store ZCash headers when they reach consensus
4. Add `storeZcashBlockHeaders()` function similar to `storeBitcoinHeaders()`
5. Update GetValidatedOracleData to populate ZCash headers from consensus

**Success Criteria**:
- ZCash headers stored in chain state after reaching consensus
- Can query stored ZCash headers by height
- Latest ZCash header height tracked correctly

### Phase 4: Deposit Verification Updates

**Goal**: Use ZCash headers for ZCash deposit verification instead of Bitcoin headers

**Tasks**:
1. Update `x/dct/keeper/msg_server_verify_deposit_block_inclusion.go`:
   - Detect ZCash chains (check CAIP-2 chain ID or asset type)
   - Use `k.validationKeeper.ZcashBlockHeaders.Get()` instead of `BtcBlockHeaders` for ZCash
   - Keep Bitcoin logic for Bitcoin chains

2. Ensure proper chain detection:
   - ZCash chains should route to ZcashBlockHeaders
   - Bitcoin chains should route to BtcBlockHeaders

**Success Criteria**:
- ZCash deposits verified using ZCash block headers
- Bitcoin deposits continue using Bitcoin block headers
- No cross-chain header confusion

### Phase 5: Testing & Validation

**Goal**: Ensure end-to-end ZCash minting flow works correctly

**Tasks**:
1. Test ZCash header fetching from sidecar RPC
2. Verify vote extensions include ZCash headers
3. Test ZCash header storage in chain state
4. Test zenZEC deposit verification with ZCash headers
5. Perform full zenZEC minting flow test

**Success Criteria**:
- ZCash headers fetched and stored successfully
- zenZEC minting works end-to-end with ZCash deposits
- No errors in deposit verification

## Current Implementation Status

✅ **Phase 1 COMPLETE**: Sidecar RPC client and configuration
✅ **Phase 2 COMPLETE**: Vote extension integration
⚠️ **Phase 3 IN PROGRESS**: Chain state storage
⏳ **Phase 4 PENDING**: Deposit verification updates
⏳ **Phase 5 PENDING**: Testing & validation
