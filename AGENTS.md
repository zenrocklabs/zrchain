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

