# zenBTC TVL Utility

A command-line utility to query the total value locked (TVL) for zenBTC on Solana.

## Overview

This utility:
1. Derives the zenBTC mint address from the program ID using the `"wrapped_mint"` seed
2. Queries the Solana RPC for the token supply
3. Displays comprehensive information including program ID, mint address, total supply, decimals, and raw amount

## Installation

Build the utility:

```bash
go build -o zenbtc-tvl cmd/tvl/zenbtc/main.go
```

Or run directly:

```bash
go run cmd/tvl/zenbtc/main.go
```

## Usage

### Basic Usage (Mainnet)

Query zenBTC TVL on Solana mainnet with default settings:

```bash
go run cmd/tvl/zenbtc/main.go
```

### Custom RPC Endpoint

Use a custom Solana RPC endpoint:

```bash
go run cmd/tvl/zenbtc/main.go -rpc https://my-custom-rpc.solana.com
```

### Custom Program ID

Query a different zenBTC program (e.g., testnet):

```bash
go run cmd/tvl/zenbtc/main.go \
  -program-id 9Gfr1YrMca5hyYRDP2nGxYkBWCSZtBm1oXBZyBdtYgNL \
  -rpc https://api.testnet.solana.com
```

## Command-Line Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-rpc` | `https://api.mainnet-beta.solana.com` | Solana RPC endpoint URL |
| `-program-id` | `9t9RfpterTs95eXbKQWeAriZqET13TbjwDa6VW6LJHFb` | zenBTC program ID (mainnet) |
| `-h` | - | Display help message |

## Output Example

```
═══════════════════════════════════════════════════════════
              zenBTC Solana Mainnet TVL
═══════════════════════════════════════════════════════════

RPC Endpoint:      https://api.mainnet-beta.solana.com
Program ID:        9t9RfpterTs95eXbKQWeAriZqET13TbjwDa6VW6LJHFb
Mint Address:      9hX59xHHnaZXLU6quvm5uGY2iDiT3jczaReHy6A6TYKw
Mint Bump Seed:    254

───────────────────────────────────────────────────────────
                    Token Supply
───────────────────────────────────────────────────────────

Total Minted:      74.02429085 zenBTC
Raw Amount:        7402429085 (smallest units)
Decimals:          8

Slot:              376141298
═══════════════════════════════════════════════════════════
```

## How It Works

### Mint Address Derivation

The zenBTC mint address is a Program Derived Address (PDA) calculated using:
- **Seed**: `"wrapped_mint"`
- **Program ID**: The zenBTC program address

This is the same derivation logic used in the Solana smart contract (see `contracts/solzenbtc/pdas.go`).

### Network-Specific Program IDs

| Network | Program ID |
|---------|------------|
| Mainnet | `9t9RfpterTs95eXbKQWeAriZqET13TbjwDa6VW6LJHFb` |
| Testnet/Regnet | `9Gfr1YrMca5hyYRDP2nGxYkBWCSZtBm1oXBZyBdtYgNL` |
| Devnet/Localnet | `886BBKmJ71jrqWhZBAzqhdNRKk761Mx9WVCPzFkY4uBb` |

## API Integration

This utility can be integrated into monitoring systems or dashboards:

```bash
# Get just the TVL amount
go run cmd/tvl/zenbtc/main.go | grep "Total Minted" | awk '{print $3}'
```

## Dependencies

- `github.com/gagliardetto/solana-go` - Solana Go SDK for RPC calls and PDA derivation

## Related Files

- Source: `cmd/tvl/zenbtc/main.go`
- PDA Logic: `contracts/solzenbtc/pdas.go`
- Program IDs: `sidecar/shared/types.go`

