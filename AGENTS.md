# ZrChain Agent Guide

## Repository Structure

### Core Modules
```
x/
├── zenbtc/          # Bitcoin-specific wrapped asset (v0 - production)
│   ├── keeper/      # Business logic
│   └── types/       # Proto-generated types
├── dct/             # Digital Currency Tokens (v1+ - generalized framework)
│   ├── keeper/      # Handles zenZEC and future wrapped assets
│   └── types/       # Proto-generated types
├── validation/      # ABCI logic, vote extensions, oracle consensus
│   └── keeper/
│       ├── abci.go           # Main PreBlocker flow
│       ├── abci_zenbtc.go    # zenBTC mint/burn/redemption
│       ├── abci_dct.go       # DCT mint/burn/redemption
│       ├── abci_types.go     # Vote extension definitions
│       └── abci_utils.go     # Helper functions
├── treasury/        # Key management, signature requests
└── identity/        # Keyrings and workspaces
```

### Sidecar Oracle
```
sidecar/
├── main.go                   # Initialization
├── oracle.go                 # Core oracle logic
├── server.go                 # gRPC endpoints
├── types.go                  # Oracle struct
├── shared/types.go           # Config and OracleState
├── bitcoin_client.go         # Bitcoin RPC client
├── zcash_client.go           # ZCash RPC client (NEW)
└── proto/api/                # gRPC service definitions
```

### Protobuf Definitions
```
proto/zrchain/
├── dct/                      # DCT messages
├── zenbtc/                   # zenBTC messages
└── validation/               # Validation module messages
```

## Module Architecture (CRITICAL)

### zenBTC Module (`x/zenbtc`)
- **Purpose**: Handles ALL Bitcoin deposits (legacy v0 system)
- **Endpoint**: `zrchain.zenbtc.Msg/VerifyDepositBlockInclusion`
- **Headers**: Uses `BtcBlockHeaders` collection
- **Flow**: BTC deposit → Bitcoin headers → zenBTC mint → Solana/EVM
- **Status**: Production, DO NOT break

### DCT Module (`x/dct`)
- **Purpose**: Handles ALL non-BTC wrapped assets (v1+ extensible)
- **Endpoint**: `zrchain.dct.Msg/VerifyDepositBlockInclusion`
- **Assets**: zenZEC (ASSET_ZENZEC), future assets
- **REJECTS**: ASSET_ZENBTC deposits with clear error message
- **Headers**: Asset-specific (ZcashBlockHeaders for zenZEC)
- **Flow**: ZEC deposit → ZCash headers → zenZEC mint → Solana/EVM

### Why Separate Modules?
1. zenBTC is battle-tested production code
2. DCT is generalized framework for new assets
3. Keeps changes isolated to prevent breaking zenBTC
4. Clear separation of concerns

## Complete Wrapped Asset Flow (zenZEC Example)

### 1. Deposit & Verification
**Files**: `x/dct/keeper/msg_server_verify_deposit_block_inclusion.go`

1. User deposits ZEC to deposit address
2. ZCash Proxy detects deposit, generates Merkle proof
3. Proxy calls `MsgVerifyDepositBlockInclusion` with:
   - `asset = ASSET_ZENZEC`
   - `chain_name = "zcash-testnet"` (or mainnet/regtest)
   - `raw_tx`, `proof`, `block_height`
4. Chain fetches ZCash header: `k.validationKeeper.ZcashBlockHeaders.Get(ctx, block_height)`
5. Verifies Merkle proof using `bitcoin.VerifyBTCLockTransaction()`
6. Creates `PendingMintTransaction` with status DEPOSITED
7. Updates supply: `CustodiedAmount += deposit`, `PendingAmount += wrapped_amount`
8. Requests Solana nonce and account via validation keeper

### 2. Mint Processing
**Files**: `x/validation/keeper/abci_dct.go:processDCTMintsSolana()`

1. PreBlocker detects pending mint (status = DEPOSITED)
2. Validates consensus on required fields (nonce, BTC/USD price, accounts)
3. Calculates fee using `CalculateFlatZenBTCMintFee()`
4. Prepares Solana mint transaction via `PrepareSolanaMintTx()`
5. Submits for MPC signing via `submitSolanaTransaction()`
6. Updates status to STAKED, stores `ZrchainTxId` and `BlockHeight`

### 3. Mint Confirmation
**Files**: `x/validation/keeper/abci_dct.go:processSolanaDCTMintEvents()`

1. Sidecar detects Solana mint event via `SolanaMintEvents`
2. Validators include event hash in vote extensions
3. PreBlocker matches event signature hash to pending mint
4. Updates supply: `PendingAmount -= amount`, `MintedAmount += amount`
5. Updates transaction status to MINTED
6. User receives zenZEC on Solana

### 4. Burn Detection
**Files**: `x/validation/keeper/abci_dct.go:storeNewDCTBurnEvents()`

1. User burns zenZEC on Solana
2. Sidecar detects burn event via `SolanaBurnEvents`
3. PreBlocker creates `BurnEvent` (status = BURNED)
4. Creates `Redemption` (status = UNSTAKED) with destination address

### 5. Redemption (Withdrawal)
**Files**: `x/dct/keeper/msg_server_submit_unsigned_redemption_tx.go`

1. ZCash Proxy polls for UNSTAKED redemptions
2. Proxy constructs unsigned ZCash transaction
3. Calls `MsgSubmitUnsignedRedemptionTx` with unsigned tx
4. Chain verifies outputs in `VerifyUnsignedRedemptionTX()`
5. Updates redemption status to AWAITING_SIGN
6. Creates signature request for MPC

### 6. Withdrawal Confirmation
**Files**: `x/validation/keeper/abci_dct.go:checkForDCTRedemptionFulfilment()`

1. MPC signs the ZCash transaction
2. PreBlocker detects fulfilled signature request
3. Calculates native amount to release using exchange rate
4. Updates supply: `MintedAmount -= amount`, `CustodiedAmount -= native_amount`
5. Updates redemption status to COMPLETED
6. Proxy broadcasts signed transaction to ZCash network

## Vote Extensions & Oracle Consensus

### ZCash Header Tracking
**Files**: 
- `x/validation/keeper/abci_types.go` - Vote extension fields
- `x/validation/keeper/abci.go` - Vote extension construction
- `x/validation/keeper/abci_utils.go` - Header retrieval

**Vote Extension Fields for ZCash**:
- `RequestedZcashBlockHeight` - Specific header height requested
- `RequestedZcashHeaderHash` - Hash of requested header
- `LatestZcashBlockHeight` - Latest known ZCash height
- `LatestZcashHeaderHash` - Hash of latest header

**Consensus Flow**:
1. `gatherOracleDataForVoteExtension()` calls sidecar for ZCash headers
2. `ConstructVoteExtension()` hashes headers and includes in vote
3. Validators reach 2/3+ consensus on header hashes
4. `PreBlocker` calls `storeZcashBlockHeaders()` to persist to chain state
5. Headers stored in `ZcashBlockHeaders` collection for deposit verification

### Collections (State Storage)
**File**: `x/validation/keeper/keeper.go`

```go
// Bitcoin headers (zenBTC module)
BtcBlockHeaders collections.Map[int64, sidecar.BTCBlockHeader]
LatestBtcHeaderHeight collections.Item[int64]

// ZCash headers (DCT module)
ZcashBlockHeaders collections.Map[int64, sidecar.BTCBlockHeader]
LatestZcashHeaderHeight collections.Item[int64]
```

## Key Configuration

### Sidecar Config
**File**: `sidecar/config.yaml`

```yaml
zcash_rpc:
  devnet: "http://zcash-node:8232"
  testnet: "http://zcash-testnet:8232"
  mainnet: "http://zcash-mainnet:8232"
```

### Chain Params (per asset)
**File**: `x/dct/types/params.proto`

```protobuf
message AssetParams {
  Asset asset = 1;                              // ASSET_ZENZEC
  string deposit_keyring_addr = 2;              // Keyring for deposits
  repeated uint64 change_address_key_ids = 8;   // Change addresses
  string proxy_address = 9;                     // Authorized proxy
  Solana solana = 12;                           // Solana-specific params
}
```

## Testing End-to-End

### Phase 5 Checklist (PENDING)
1. Configure ZCash RPC in `sidecar/config.yaml`
2. Start sidecar and verify ZCash client initialization
3. Trigger ZCash header request via vote extension
4. Verify headers stored: query `ZcashBlockHeaders` collection
5. Test deposit: Call DCT `VerifyDepositBlockInclusion` with ASSET_ZENZEC
6. Verify mint: Check Solana for zenZEC receipt
7. Test burn: Burn zenZEC on Solana
8. Verify redemption: Check ZCash transaction broadcast
9. Regression: Ensure zenBTC flow still works (separate module)

## Common Pitfalls

### Module Confusion
❌ **Wrong**: Calling zenBTC endpoint for zenZEC deposits
✅ **Right**: zenBTC module for BTC, DCT module for zenZEC

### Header Collections
❌ **Wrong**: Using `BtcBlockHeaders` for zenZEC verification
✅ **Right**: Using `ZcashBlockHeaders` for zenZEC, `BtcBlockHeaders` for zenBTC

### Chain Names
- Bitcoin: `"mainnet"`, `"testnet"`, `"regtest"`
- ZCash: `"zcash-mainnet"`, `"zcash-testnet"`, `"zcash-regtest"`

### Wallet Types
**File**: `x/dct/keeper/msg_server_verify_deposit_block_inclusion.go:WalletTypeFromChainName()`
- Maps chain names to wallet types
- Critical for treasury key lookups

## Adding New Wrapped Assets

1. **Add to Asset enum**: `proto/zrchain/dct/params.proto`
2. **Implement header fetching**: Add RPC client in sidecar (if new chain)
3. **Add vote extension fields**: `x/validation/keeper/abci_types.go`
4. **Add header storage**: `x/validation/keeper/keeper.go` (new collection)
5. **Update deposit verification**: `x/dct/keeper/msg_server_verify_deposit_block_inclusion.go`
6. **Add chain name mapping**: `WalletTypeFromChainName()`
7. **Configure asset params**: Chain governance or genesis
8. **Test end-to-end**: Follow Phase 5 checklist
