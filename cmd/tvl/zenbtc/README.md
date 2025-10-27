# zenBTC Cross-Chain TVL Utility

A command-line utility to query the total value locked (TVL) for zenBTC across Solana and Ethereum.

## Overview

This utility:
1. Derives the zenBTC mint address from the Solana program ID using the `"wrapped_mint"` seed
2. Queries Solana RPC for the SPL token supply
3. Queries Ethereum RPC for the ERC-20 token supply
4. Queries Zenrock API for the amount of custodied (locked) Bitcoin
5. Calculates combined TVL across both chains
6. Shows distribution percentages between Solana and Ethereum
7. Displays the actual Bitcoin backing the wrapped tokens

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

Query zenBTC TVL across Solana and Ethereum mainnets with default settings:

```bash
go run cmd/tvl/zenbtc/main.go
```

### Custom RPC Endpoints

Use custom RPC endpoints for either chain:

```bash
go run cmd/tvl/zenbtc/main.go \
  -solana-rpc https://my-solana-rpc.com \
  -ethereum-rpc https://my-eth-rpc.com
```

### Custom Program ID

Query a different zenBTC program (e.g., testnet):

```bash
go run cmd/tvl/zenbtc/main.go \
  -program-id 9Gfr1YrMca5hyYRDP2nGxYkBWCSZtBm1oXBZyBdtYgNL \
  -solana-rpc https://api.testnet.solana.com
```

## Command-Line Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-solana-rpc` | `https://api.mainnet-beta.solana.com` | Solana RPC endpoint URL |
| `-ethereum-rpc` | `https://mainnet.gateway.tenderly.co/` | Ethereum RPC endpoint URL |
| `-program-id` | `9t9RfpterTs95eXbKQWeAriZqET13TbjwDa6VW6LJHFb` | zenBTC program ID on Solana (mainnet) |
| `-h` | - | Display help message |

## Output Example

```
═══════════════════════════════════════════════════════════
          zenBTC Cross-Chain Total Value Locked
═══════════════════════════════════════════════════════════

┌─────────────────────────────────────────────────────────┐
│                    SOLANA MAINNET                       │
└─────────────────────────────────────────────────────────┘

RPC Endpoint:      https://api.mainnet-beta.solana.com
Program ID:        9t9RfpterTs95eXbKQWeAriZqET13TbjwDa6VW6LJHFb
Mint Address:      9hX59xHHnaZXLU6quvm5uGY2iDiT3jczaReHy6A6TYKw
Mint Bump Seed:    254

Total Supply:      74.02429085 zenBTC
Raw Amount:        7402429085
Percentage:        99.96%
Slot:              376159763

┌─────────────────────────────────────────────────────────┐
│                   ETHEREUM MAINNET                      │
└─────────────────────────────────────────────────────────┘

RPC Endpoint:      https://mainnet.gateway.tenderly.co/
Contract Address:  0x2fE9754d5D28bac0ea8971C0Ca59428b8644C776

Total Supply:      0.02900092 zenBTC
Raw Amount:        2900092
Percentage:        0.04%

═══════════════════════════════════════════════════════════
                   TOTAL ACROSS CHAINS
═══════════════════════════════════════════════════════════

Combined TVL:      74.05329177 zenBTC
Raw Amount:        7405329177

Distribution:
  Solana:          99.96%
  Ethereum:        0.04%

═══════════════════════════════════════════════════════════
                     BACKING ASSETS
═══════════════════════════════════════════════════════════

Locked BTC:        74.05395066 BTC
Raw Satoshis:      7405395066

Backing Ratio:     1.000089 (100.0089%)

═══════════════════════════════════════════════════════════
```

## How It Works

### Mint Address Derivation

The zenBTC mint address is a Program Derived Address (PDA) calculated using:
- **Seed**: `"wrapped_mint"`
- **Program ID**: The zenBTC program address

This is the same derivation logic used in the Solana smart contract (see `contracts/solzenbtc/pdas.go`).

### Custodied Bitcoin

The utility queries the Zenrock API at `https://api.diamond.zenrocklabs.io/zenbtc/supply` to retrieve the actual amount of Bitcoin that is custodied (locked) to back the zenBTC tokens. This provides transparency into the 1:1 backing of the wrapped asset.

The `custodiedBTC` field from the API response represents the total Bitcoin held in custody, measured in satoshis (the smallest Bitcoin unit). The utility automatically converts this to BTC for display.

### Backing Ratio

The backing ratio is calculated as:

```
Backing Ratio = Locked BTC / Total zenBTC Supply
```

This ratio shows the collateralization level:
- **1.0 (100%)** = Fully backed 1:1
- **> 1.0 (> 100%)** = Over-collateralized (more BTC locked than zenBTC in circulation)
- **< 1.0 (< 100%)** = Under-collateralized (less BTC than zenBTC - should never happen)

A ratio slightly above 1.0 is healthy and accounts for rounding differences in the system.

### Network-Specific Solana Program IDs

| Network | Program ID |
|---------|------------|
| Mainnet | `9t9RfpterTs95eXbKQWeAriZqET13TbjwDa6VW6LJHFb` |
| Testnet/Regnet | `9Gfr1YrMca5hyYRDP2nGxYkBWCSZtBm1oXBZyBdtYgNL` |
| Devnet/Localnet | `886BBKmJ71jrqWhZBAzqhdNRKk761Mx9WVCPzFkY4uBb` |

## API Integration

This utility can be integrated into monitoring systems or dashboards:

```bash
# Get combined TVL
go run cmd/tvl/zenbtc/main.go 2>/dev/null | grep "Combined TVL" | awk '{print $3}'

# Get locked BTC backing
go run cmd/tvl/zenbtc/main.go 2>/dev/null | grep "Locked BTC" | awk '{print $3}'

# Get backing ratio
go run cmd/tvl/zenbtc/main.go 2>/dev/null | grep "Backing Ratio" | awk '{print $3}'

# Get Solana supply
go run cmd/tvl/zenbtc/main.go 2>/dev/null | grep -A 10 "SOLANA MAINNET" | grep "Total Supply" | awk '{print $3}'

# Get Ethereum supply
go run cmd/tvl/zenbtc/main.go 2>/dev/null | grep -A 10 "ETHEREUM MAINNET" | grep "Total Supply" | awk '{print $3}'

# Get distribution percentages
go run cmd/tvl/zenbtc/main.go 2>/dev/null | grep -A 2 "Distribution:" | tail -2
```

## Data Sources

- **Solana Supply**: Queried directly from Solana mainnet via RPC (`GetTokenSupply`)
- **Ethereum Supply**: Queried from Ethereum mainnet via contract call (`TotalSupply`)
- **Custodied BTC**: Retrieved from Zenrock API endpoint (`https://api.diamond.zenrocklabs.io/zenbtc/supply`)

## Dependencies

- `github.com/gagliardetto/solana-go` - Solana Go SDK for RPC calls and PDA derivation
- `github.com/ethereum/go-ethereum` - Ethereum Go client for contract interactions
- `github.com/Zenrock-Foundation/zrchain/v6/zenbtc/bindings` - zenBTC ERC-20 contract bindings

## Related Files

- Source: `cmd/tvl/zenbtc/main.go`
- Solana PDA Logic: `contracts/solzenbtc/pdas.go`
- Contract Addresses: `sidecar/shared/types.go`
- Ethereum Bindings: `zenbtc/bindings/`

