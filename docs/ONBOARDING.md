# Developer Onboarding Guide

Welcome to Zenrock! This guide will help you get up to speed with the zrChain codebase.

## Table of Contents

1. [First 30 Minutes](#first-30-minutes)
2. [Understanding the Architecture](#understanding-the-architecture)
3. [Project Structure](#project-structure)
4. [Common Development Tasks](#common-development-tasks)
5. [First Week Checklist](#first-week-checklist)
6. [Git & Release Workflow Basics](#git--release-workflow-basics)
7. [Worktrees & Graphite](#worktrees--graphite)
8. [Useful Commands](#useful-commands)
9. [Debugging Tips](#debugging-tips)
10. [Where to Get Help](#where-to-get-help)

## First 30 Minutes

### Step 1: Clone & Bootstrap Tooling (5 min)

```bash
git clone git@github.com:Zenrock-Foundation/zrchain.git
cd zrchain
./scripts/dev-setup.sh   # Installs Go, just, uv, proto toolchain, and checks Docker
make install
zenrockd version  # Should output version info
```

If `zenrockd` command not found:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

> `scripts/dev-setup.sh` is idempotent—you can rerun it any time to install missing dependencies. It auto-detects macOS/Linux package managers, installs Go 1.25+, [`uv`](https://docs.astral.sh/uv/), [`just`](https://github.com/casey/just), buf/protoc plugins, and walks you through Docker setup when it can’t automate it.

### Step 2: Start the Local Network (5 min)

Open **two terminal tabs** in the zrchain directory:

**Tab 1:**
```bash
./init.sh --localnet 1
```

**Tab 2:**
```bash
./init.sh --localnet 2
```

This starts 2 validators with pre-seeded data and a sidecar oracle. Wait ~30 seconds for the network to stabilize.

### Step 3: Verify It's Running (5 min)

```bash
# Check node status
zenrockd status

# Query some data
zenrockd query identity workspaces
zenrockd query zenbtc pending-mint-transactions
```

**You now have a working local environment!** ✅

### Step 4: Explore the Docs (15 min)

Read these in order:
1. This file (you're reading it!)
2. [`../CLAUDE.md`](../CLAUDE.md) - Developer reference
3. [`../AGENTS.md`](../AGENTS.md) - Architecture deep-dive

## Understanding the Architecture

### What is Zenrock?

Zenrock is a Cosmos blockchain that specializes in managing **wrapped assets** across multiple chains.

**Simple flow:**
```
User deposits Bitcoin
    ↓
Sidecar detects it
    ↓
Chain validates via oracle
    ↓
zenBTC minted on Solana
    ↓
User has bridged asset
```

### Core Components

#### 1. **The Blockchain (zrChain)**

Located in `/app/` and `/x/` directories.

Modules handle:
- **Identity** - Who you are (keyrings, workspaces)
- **Treasury** - What assets you control
- **Validation** - Validator consensus
- **ZenBTC** - Bitcoin wrapped asset (production)
- **DCT** - Digital Currency Tokens (new assets like zenZEC)
- **Policy** - Approval workflows

#### 2. **The Sidecar Oracle**

Located in `/sidecar/` directory.

This is a **separate service** that:
- Watches Bitcoin, ZCash, Solana, Ethereum networks
- Detects deposits and burn events
- Runs EigenLayer operator loops and maintains local keystores (bridging signatures live in the external MPC stack)
- Reports data to validators via vote extensions

#### 3. **Communication**

They talk via gRPC:
- Blockchain calls sidecar to get block headers, prices, etc.
- Sidecar submits transactions to blockchain for signature requests

### Architecture Diagram

```
┌────────────────────────────────────────────┐
│         zrChain Blockchain                 │
│  (Modules: identity, treasury, zenbtc...)  │
│                                            │
│  • Stores wrapped assets                   │
│  • Validates cross-chain operations        │
│  • Manages signatures                      │
└────────────────────────────────────────────┘
                    ↕ (gRPC)
┌────────────────────────────────────────────┐
│      Sidecar Oracle (Separate Process)     │
│                                            │
│  • Watches external blockchains            │
│  • Manages signing keys                    │
│  • Reports data via vote extensions        │
└────────────────────────────────────────────┘
```

## Project Structure

### Top-Level Directories

| Directory | Purpose |
|-----------|---------|
| `app/` | Application setup, module wiring, blockchain hooks |
| `cmd/` | CLI binary entry points |
| `x/` | **Custom blockchain modules** (where logic lives) |
| `sidecar/` | Oracle service for key management |
| `proto/` | Protocol Buffer definitions (messages, types) |
| `tests/` | End-to-end tests |
| `testutil/` | Testing utilities |
| `docs/` | Documentation |

### Understanding the `x/` Directory

Each module is self-contained:

```
x/{module}/
├── keeper/              # Business logic & state management
│   ├── keeper.go       # Core keeper setup
│   ├── msg_server.go   # Message handlers
│   ├── query.go        # Query handlers
│   └── *_test.go       # Tests
├── types/              # Data structures & messages
│   ├── messages.go     # Message definitions
│   ├── codec.go        # Encoding/decoding
│   └── params.proto    # Protobuf definitions
├── module/             # Module registration
└── simulation/         # Randomized testing
```

**Key modules:**

| Module | What it does |
|--------|------------|
| `identity` | Manages keyrings (groups of keys) and workspaces (user containers) |
| `treasury` | Manages assets and signature requests |
| `validation` | Replaces standard staking, manages validator set and oracle consensus |
| `zenbtc` | Bitcoin deposits become zenBTC tokens (production v0) |
| `dct` | Framework for new wrapped assets like zenZEC (v1+) |
| `policy` | Approval policies and governance |
| `mint` | Mints tokens |
| `zentp` | Zenrock token protocol |

### Key Files to Know

| File | Purpose |
|------|---------|
| `app/app.go` | Wires all modules together (start here to understand initialization) |
| `cmd/zenrockd/cmd/root.go` | CLI entry point |
| `app/abci.go` | Block lifecycle hooks (PreBlocker, BeginBlocker, EndBlocker) |
| `sidecar/main.go` | Sidecar startup |
| `sidecar/oracle.go` | Core oracle logic |
| `CLAUDE.md` | Developer reference guide |
| `AGENTS.md` | Detailed architecture documentation |

## Common Development Tasks

### Task 1: Understanding a Feature

**Goal:** Figure out how Bitcoin deposits work

1. Search for the entry point:
   ```bash
   grep -r "VerifyDepositBlockInclusion" x/*/keeper/
   ```

2. This finds it in `x/zenbtc/keeper/msg_server_verify_deposit_block_inclusion.go`

3. Read the file to understand the flow

4. Check related files:
   - `x/zenbtc/types/messages.go` - Message definition
   - `x/validation/keeper/abci.go` - Validation logic
   - Tests in `x/zenbtc/keeper/*_test.go`

### Task 2: Adding a New Message

**Goal:** Add a new transaction type

1. Define message in `/proto/zrchain/{module}/tx.proto`:
   ```protobuf
   message MsgMyNewAction {
     string creator = 1;
     string param1 = 2;
   }
   ```

2. Generate code:
   ```bash
   make proto-gen
   ```

3. Implement handler in `x/{module}/keeper/msg_server_my_new_action.go`:
   ```go
   func (k msgServer) MyNewAction(ctx context.Context, msg *types.MsgMyNewAction) (*types.MsgMyNewActionResponse, error) {
     // Your logic here
   }
   ```

4. Add CLI command in `x/{module}/module/cli/tx.go`

5. Write tests in `x/{module}/keeper/msg_server_test.go`

### Task 3: Running Tests

```bash
# Test specific module
go test -v ./x/identity/...

# Test specific test
go test -v -run TestNewKeyring ./x/identity/keeper/

# All tests (excludes e2e)
go test -race -v $(go list ./... | grep -v '/tests/e2e')

# With coverage (includes e2e; requires Docker/sidecar deps)
go test -cover ./...
```

### Task 4: Querying the Blockchain

```bash
# List all keyrings
zenrockd query identity keyrings

# Get specific keyring by address
zenrockd query identity keyring-by-address zen1...

# List workspaces
zenrockd query identity workspaces

# Check pending Bitcoin mints
zenrockd query zenbtc pending-mint-transactions

# Check module parameters
zenrockd query {module} params
```

## First Week Checklist

- [ ] **Day 1**
  - [ ] Clone repo
  - [ ] Run `make install`
  - [ ] Start local network with `./init.sh --localnet 1` and `./init.sh --localnet 2`
  - [ ] Verify network is running with `zenrockd status`

- [ ] **Day 2-3**
  - [ ] Read `CLAUDE.md` for developer reference
  - [ ] Read `AGENTS.md` for architecture
  - [ ] Pick one module (e.g., `x/identity`) and trace through a message handler
  - [ ] Run tests: `go test -v ./x/identity/...`

- [ ] **Day 4-5**
  - [ ] Understand the sidecar: read `sidecar/main.go` and `sidecar/oracle.go`
  - [ ] Understand ABCI hooks: read `app/abci.go`
  - [ ] Understand vote extensions: read `x/validation/keeper/abci_types.go`
  - [ ] Try making a small code change and run tests

- [ ] **Week 2+**
  - [ ] Work on first issue/task
  - [ ] Ask questions when stuck
  - [ ] Review PRs from teammates
  - [ ] Contribute your first feature

## Git & Release Workflow Basics

- Commit messages must follow the Conventional Commits pattern `type(optional-scope)!: description`.
- Our release pipeline uses `semantic-release`: `feat` bumps the **minor** version, `fix` and `perf` bump the **patch**, and commits marked with `BREAKING CHANGE` (or `type!`) bump the **major**.
- Other conventional types (`docs`, `chore`, `refactor`, `test`, etc.) keep the current version but must still follow the format.
- Pull request titles should mirror the same format because they are squashed on merge.
- Read [CONTRIBUTING.md](../CONTRIBUTING.md#git-conventions) for full examples and nuance.

## Worktrees & Graphite

We run parallel “agent” work in separate git worktrees. Everything lives in `scripts/git/` and is surfaced via `just` (installed by `scripts/dev-setup.sh`). The scripts detect Graphite automatically and keep stacks in sync when you opt in:

- `just worktree-new` (`scripts/git/worktree-new.sh`) – prompts for a branch name and base branch, then creates a new worktree. If the [Graphite CLI](https://graphite.dev/) (`gt`) is installed, you’ll be asked whether to use it; opting in calls `gt create -m …` on the parent branch so the new worktree sits in the correct stack.
- `just worktree-switch` – fuzzy-find an existing branch/worktree. With Graphite enabled it uses `gt get`/`gt checkout` to sync metadata, otherwise it falls back to plain git and auto-runs `gt track` afterward when possible.
- `just worktree-cleanup` – clean up stale worktrees via multi-select.

The scripts are Graphite-native when the CLI is present, but degrade gracefully to vanilla git so you can keep working even if you’re not using stacks. This is the backbone of the “Graphite + worktrees + agentic dev” workflow—fire up multiple worktrees, drop an AI agent into each, and let them chase tasks in parallel.

## Useful Commands

### Development

```bash
# Build
make install                # Install binary globally
make build                  # Build to ./build/
make build-linux           # Cross-compile for Linux

# Code Generation
make proto                  # Generate all Protobufs
make proto-gen             # Generate from .proto files
make proto-lint            # Lint .proto files
make proto-format          # Format .proto files
make web-gen               # Generate web client

# Testing
go test -v ./x/{module}/...              # Test module
go test -v -run TestName ./path/package # Test specific test
go test -race -v $(go list ./... | grep -v '/tests/e2e') # All unit tests (skip e2e)
go test -cover ./...                     # With coverage (includes e2e; requires Docker)

# Running
make install                             # Install zenrockd
./init.sh --localnet 1                  # Start validator
zenrockd start                          # Manual start
```

### Sidecar

```bash
make build-sidecar         # Build sidecar
make run-sidecar          # Run sidecar
make run-alt-sidecar      # Run on alternate port
```

### Querying

```bash
zenrockd query {module} {query-type}    # Generic query
zenrockd query identity workspaces       # List workspaces
zenrockd query identity keyrings        # List keyrings
zenrockd query zenbtc pending-mint-transactions  # Pending mints
zenrockd tx {module} {msg-type} --help  # Message help
```

## Debugging Tips

### Node Won't Start

```bash
# Check logs
tail -f ~/.zrchain/chain.log

# Reset local state
rm -rf ~/.zrchain

# Rebuild from scratch
make clean && make install
./init.sh --localnet 1
```

### Tests Failing

```bash
# Run with verbose output
go test -v -race ./x/{module}/...

# Run single test
go test -v -run TestName ./path

# Check for race conditions
go test -race ./...

# With coverage
go test -cover ./...
```

### Transaction Won't Execute

```bash
# Check balance
zenrockd query bank balances zen1{address}

# Check account exists
zenrockd query auth account zen1{address}

# Check module params (fees, etc.)
zenrockd query {module} params
```

### Sidecar Issues

```bash
# Check if running
lsof -i :9191  # Default sidecar port

# Check logs (look in terminal where you ran init.sh)
# Should show Bitcoin/ZCash header requests and price updates

# Restart
Ctrl+C in sidecar terminal and rerun
```

### Understanding Vote Extensions

Vote extensions are how the sidecar reports data to validators:

1. Validator asks sidecar for data (block headers, prices)
2. Sidecar responds
3. Validator includes in vote
4. 2/3+ consensus reached on that data
5. Data stored on-chain for later use

See `x/validation/keeper/abci_types.go` for what data is included.

## Where to Get Help

### Documentation

- **`CLAUDE.md`** - Developer reference (commands, architecture, common patterns)
- **`AGENTS.md`** - Detailed wrapped asset flow documentation
- **`CONTRIBUTING.md`** - Contribution guidelines, git conventions, security considerations
- **This file** - Getting started guide

### Code Search

```bash
# Find where something is used
grep -r "SearchTerm" x/

# Find message handlers
grep -r "func (k msgServer)" x/*/keeper/msg_server*.go

# Find query handlers
grep -r "func (k Keeper) " x/*/keeper/query*.go

# Find keeper initialization
grep -r "type Keeper struct" x/*/keeper/
```

### Testing & Experimentation

```bash
# Write a simple test file to understand a module
touch x/{module}/keeper/explore_test.go

# Run just that test
go test -v x/{module}/keeper/explore_test.go
```

### Team Communication

- **Discord** - [zenrockfoundation](https://discord.com/invite/zenrockfoundation)
- **Telegram** - [officialZenrock](https://t.me/officialZenrock)
- **GitHub Issues** - [Report bugs](https://github.com/Zenrock-Foundation/zrchain/issues)
- **PR Comments** - Ask questions on open PRs
- **Code Review** - Ask teammates to review your code early and often

## Key Concepts to Learn

### Modules

A module is a self-contained piece of blockchain logic. Each module has:
- **Messages** - Transactions users send
- **Queries** - Data users can read
- **State** - Data stored on-chain
- **Keeper** - Logic that manages the state

### Keepers

The `Keeper` is the main interface for a module. Think of it as the "brain" of the module. Other modules call keeper methods to read/write state.

### Collections

Modern Cosmos SDK uses `collections` for state management instead of old KVStore. Collections are like typed maps:

```go
MyData collections.Map[string, MyType]  // key → value
MyItem collections.Item[MyType]         // single item
```

### Protocol Buffers

Messages and types are defined in `.proto` files and auto-generated to Go:

```bash
make proto-gen  # Generates Go code from .proto files
```

### Vote Extensions

Validators can include extra data in their votes (used for oracle consensus):
1. PreBlocker calls sidecar for data
2. Data included in vote
3. Consensus on data reached
4. Data stored on-chain

See `x/validation/keeper/abci_types.go`.

### Wrapped Assets Flow (zenBTC Example)

```
1. User deposits BTC
   ↓
2. Sidecar detects deposit
   ↓
3. Sidecar calls chain: "Verify this deposit"
   ↓
4. Chain checks: "Is this BTC header valid?" (using data from vote extensions)
   ↓
5. If valid: Create PendingMintTransaction
   ↓
6. PreBlocker processes pending mints
   ↓
7. Chain submits to MPC: "Sign Solana mint transaction"
   ↓
8. MPC signs and returns signature
   ↓
9. Chain broadcasts to Solana
   ↓
10. User receives zenBTC on Solana
```

## Next Steps

1. **Immediate (Today)**
   - Run `make install` and start local network
   - Verify with `zenrockd status`

2. **This Week**
   - Read `CLAUDE.md` and `AGENTS.md`
   - Pick a module and understand it
   - Run tests
   - Make a small change and see tests pass

3. **Next Week**
   - Start working on your first task
   - Ask questions when stuck
   - Review teammates' code

Good luck, and welcome to the team! 🚀
