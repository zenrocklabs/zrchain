# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**zrChain** is a Cosmos SDK-based blockchain providing "bedrock security infrastructure" for omnichain operations. It specializes in wrapped asset management (currently Bitcoin via zenBTC, expanding to ZCash via DCT framework) with oracle consensus, MPC signing, and cross-chain settlement.

**Tech Stack**: Go 1.25+, Cosmos SDK v0.50.13, CometBFT, Protocol Buffers (proto3), Sidecar oracle (Go)

## Build & Development Commands

### Installation
```bash
./scripts/dev-setup.sh   # Install Go/just/uv and proto toolchain; safe to rerun
make install             # Build and install zenrockd binary to $GOPATH/bin
```

### Building
```bash
make build           # Build to ./build/
make build-linux     # Cross-compile for Linux (amd64)
make build-sidecar   # Build sidecar oracle binary
```

### Testing
```bash
# Unit tests (skip e2e and proxy tests)
go test -race -v $(go list ./... | grep -v '/tests/e2e')

# Full test suite (includes e2e, requires running Docker daemon)
go test -race -v ./...

# Run single test
go test -race -v -run TestName ./path/to/package

# Run tests for a specific module (e.g., dct)
go test -v ./x/dct/...

# With coverage
go test -cover ./...
```

### Linting & Code Generation
```bash
make proto            # Generate all Protobuf files (runs proto-gen)
make proto-gen        # Generate Go code from .proto files
make proto-lint       # Lint Protobuf definitions
make proto-format     # Format .proto files with clang-format
make proto-all        # Format, lint, and generate in sequence
make web-gen          # Generate web client and hooks from Protobuf
```

### Running Locally
```bash
# Option 1: Simple start
zenrockd start

# Option 2: Init script (recommended for development with sample data)
./init.sh --localnet 1    # Terminal 1: First validator + sidecar
./init.sh --localnet 2    # Terminal 2: Second validator

# Sidecar (oracle) - required for wrapped asset operations
make build-sidecar
make run-sidecar           # Main sidecar on port 9090
make run-alt-sidecar       # Alternative sidecar on port 9393
```

### Complete Workflow
```bash
make proto-all    # Regenerate all proto files
make build        # Build everything
go test -v ./... # Run tests
make install      # Install to GOPATH
```

## Architecture Overview

### Module Structure (`x/` directory)

The blockchain consists of specialized modules:

- **`x/dct/`** (Digital Currency Tokens) - Generalized framework for wrapped assets (v1+)
  - Handles zenZEC and future non-BTC wrapped assets
  - Primary entry: `VerifyDepositBlockInclusion` message

- **`x/zenbtc/`** (Digital Bitcoin) - Bitcoin wrapped asset (v0, production)
  - Handles all Bitcoin deposits and minting
  - Kept separate from DCT to ensure stability
  - Primary entry: `VerifyDepositBlockInclusion` message

- **`x/validation/`** - ABCI hooks, vote extensions, oracle consensus
  - `keeper/abci.go` - Main PreBlocker flow
  - `keeper/abci_zenbtc.go` - zenBTC lifecycle hooks
  - `keeper/abci_dct.go` - DCT lifecycle hooks
  - Manages BTC and ZCash block header consensus

- **`x/treasury/`** - Key management and MPC signature requests
  - Stores keyrings and manages signing workflows

- **`x/identity/`** - Workspaces and keyring management
  - Manages user identities and multi-sig keyrings

- **`x/policy/`** - Governance and action approval policies

- **`x/mint/`**, **`x/zentp/`**, **`x/zenex/`** - Additional specialized modules

### Critical Collections (State Storage)

Located in `x/validation/keeper/keeper.go`:

```go
// Bitcoin headers used by zenBTC module
BtcBlockHeaders collections.Map[int64, sidecar.BTCBlockHeader]
LatestBtcHeaderHeight collections.Item[int64]

// ZCash headers used by DCT module
ZcashBlockHeaders collections.Map[int64, sidecar.BTCBlockHeader]
LatestZcashHeaderHeight collections.Item[int64]
```

### Sidecar Oracle Architecture

**Location**: `sidecar/` directory

The oracle sidecar provides blockchain data and serves gRPC endpoints:

- **Core**: `main.go`, `oracle.go`, `server.go`
- **Chain Clients**: `bitcoin_client.go`, `zcash_client.go`
- **Configuration**: `sidecar/config.yaml` (RPC endpoints for Bitcoin/ZCash networks)
- **gRPC API**: `sidecar/proto/api/` (service definitions)

**Vote Extension Fields** (in `x/validation/keeper/abci_types.go`):
- Bitcoin: `RequestedBtcBlockHeight`, `LatestBtcBlockHeight`, header hashes
- ZCash: `RequestedZcashBlockHeight`, `LatestZcashBlockHeight`, header hashes

## Wrapped Asset Flow Example (zenZEC)

### 1. Deposit & Verification
- User deposits ZEC to designated address
- ZCash Proxy detects deposit, generates Merkle proof
- Calls `x/dct/Msg/VerifyDepositBlockInclusion` with ASSET_ZENZEC
- Chain verifies proof using stored `ZcashBlockHeaders`
- Creates `PendingMintTransaction` (status: DEPOSITED)

### 2. Mint Processing
- PreBlocker (in `x/validation/keeper/abci_dct.go:processDCTMintsSolana()`) detects pending mint
- Validates oracle consensus on nonce, exchange rate, accounts
- Fee calculation is currently disabled (fee parameter is set to zero while DCT mint economics are finalized)
- Submits Solana mint transaction to MPC for signing
- Updates status to STAKED

### 3. Mint Confirmation
- Sidecar detects Solana mint event
- PreBlocker matches event signature to pending mint
- Updates status to MINTED, user receives zenZEC on Solana

### 4. Burn & Redemption
- User burns zenZEC on Solana
- Sidecar detects burn event, creates `Redemption` (UNSTAKED)
- ZCash Proxy submits unsigned redemption transaction
- MPC signs, proxy broadcasts to ZCash network

## Key Files & Dependencies

### Protobuf Definitions
- `proto/zrchain/dct/` - DCT messages
- `proto/zrchain/zenbtc/` - zenBTC messages
- `proto/zrchain/validation/` - Vote extensions and validation types

### Common Dependencies
- `cosmossdk.io/collections` - State storage
- `github.com/btcsuite/btcd` - Bitcoin library
- `github.com/cometbft/cometbft` - Consensus
- `github.com/cosmos/cosmos-sdk` - Framework (custom fork at v0.50.13-zenrock2)

### Important Utility Functions
- Bitcoin verification: `bitcoin.VerifyBTCLockTransaction()`
- Fee calculation: `CalculateFlatZenBTCMintFee()`
- Solana transaction prep: `PrepareSolanaMintTx()`
- Chain name mapping: `WalletTypeFromChainName()` (maps network names to wallet types)

## Common Development Patterns

### Adding a New Wrapped Asset

1. Define asset enum in `proto/zrchain/dct/params.proto`
2. Implement RPC client in sidecar (if new blockchain)
3. Add vote extension fields in `x/validation/keeper/abci_types.go`
4. Add header storage collection in `x/validation/keeper/keeper.go`
5. Handle deposit in `x/dct/keeper/msg_server_verify_deposit_block_inclusion.go`
6. Add asset-specific chain name mapping in `WalletTypeFromChainName()`
7. Configure asset parameters via governance

### Testing Block Header Verification

Headers are fetched via vote extensions during consensus and stored for later verification. To test deposit verification:

1. Ensure sidecar is running and configured with correct RPC endpoints
2. Trigger header request (happens automatically in consensus)
3. Inspect stored headers by exporting state (`zenrockd export > state.json` and check `app_state.validation.btc_block_headers` / `zcash_block_headers`), or by calling `keeper.GetBtcBlockHeaders()` / `GetZcashBlockHeaders()` in an integration test harness—there is no direct CLI query today
4. Submit deposit with matching block height

### Panics & Error Handling

⚠️ **Security Note**: Panics in module code halt the chain. Avoid panics for user-triggerable conditions. Return errors instead, especially for:
- Malformed user inputs
- External data validation failures
- Race conditions on concurrent access

## Git & PR Conventions

- **Branch naming**: Feature branches should be descriptive (e.g., `feature/zenbtc-fix`, `sasha/v6rev41p6`)
- **Commits**: Use conventional commits (`type(optional-scope)!: description`). `semantic-release` reads these: `feat` bumps the minor version, `fix`/`perf` bump the patch, and commits with `BREAKING CHANGE` notes (or `type!`) bump the major. Other conventional types (`docs`, `chore`, `refactor`, `test`, etc.) do not publish a release.
- **PRs**: One major change per PR, keep atomic and tested
- **Merging**: Use "Squash and merge" strategy
- **Testing**: All tests must pass locally before pushing: `go test ./...`

## Critical Constraints

### Do NOT Break zenBTC
- zenBTC is production code in `x/zenbtc/`
- Any changes affecting `VerifyDepositBlockInclusion` or block header storage must be regression tested
- Changes to validation ABCI must test both zenBTC and DCT flows

### Module Separation
- zenBTC endpoint: `zrchain.zenbtc.Msg/VerifyDepositBlockInclusion`
- DCT endpoint: `zrchain.dct.Msg/VerifyDepositBlockInclusion` (REJECTS ASSET_ZENBTC)
- DCT explicitly rejects Bitcoin deposits with clear error messages
- BTC headers stored in `BtcBlockHeaders`, ZCash headers in `ZcashBlockHeaders` (never mix)

### Chain Name Mappings
- Bitcoin: `"mainnet"`, `"testnet"`, `"regtest"`
- ZCash: `"zcash-mainnet"`, `"zcash-testnet"`, `"zcash-regtest"`
- Always validate chain names match expected format

## Debugging Tips

- **Block header sync issues**: Check sidecar logs and RPC endpoint connectivity
- **Verification failures**: Export state or query via test helpers to confirm header heights/hashes (no direct CLI yet)
- **Signature submission fails**: Check MPC server logs and signature request status (`zenrockd q treasury`)
- **Consensus not reached**: Verify all validators are running with same chain config
- **Local testing**: Use `./init.sh` to spin up multiple validators with pre-seeded genesis data

## Related Resources

- README.md - Project overview and setup
- CONTRIBUTING.md - Contribution guidelines and common security considerations
- AGENTS.md - Detailed wrapped asset flow architecture (refer to for complex flows)
- `.github/workflows/test.yml` - CI pipeline configuration
