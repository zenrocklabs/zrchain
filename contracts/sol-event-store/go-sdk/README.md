# EventStore Go SDK

A high-performance Go client library for the Solana EventStore program that provides efficient access to Bitcoin bridge events across multiple shards.

## Features

- **🚀 Single RPC Call**: Fetch all 1000+ events across multiple shards in one efficient call
- **🪙 Bitcoin Address Support**: Automatic decoding of all Bitcoin address formats (Legacy, P2SH, Bech32, Taproot)
- **📊 Event Types**: Support for ZenBTC and Rock wrap/unwrap events
- **⚡ High Performance**: Optimized batch fetching with parallel shard processing
- **🎯 Type Safety**: Strongly typed Go structs with proper deserialization
- **🔍 Address Analysis**: Built-in Bitcoin address type detection and validation

## Installation

```bash
go get github.com/yourorg/eventstore-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/gagliardetto/solana-go/rpc"
    eventstore "github.com/yourorg/eventstore-sdk"
)

func main() {
    // Create RPC client
    client := rpc.New(rpc.DevNet_RPC)

    // Create EventStore client
    esClient := eventstore.NewClient(client, nil)

    ctx := context.Background()

    // Get all events in a single call
    allEvents, err := esClient.GetAllEvents(ctx)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("ZenBTC Wrap Events: %d\n", len(allEvents.ZenbtcWrapEvents))
    fmt.Printf("ZenBTC Unwrap Events: %d\n", len(allEvents.ZenbtcUnwrapEvents))
    fmt.Printf("Rock Wrap Events: %d\n", len(allEvents.RockWrapEvents))
    fmt.Printf("Rock Unwrap Events: %d\n", len(allEvents.RockUnwrapEvents))

    // Analyze Bitcoin addresses
    for _, event := range allEvents.ZenbtcUnwrapEvents {
        addr := event.GetBitcoinAddress()
        addrType := eventstore.GetBitcoinAddressType(addr)
        fmt.Printf("Bitcoin Address: %s (%s)\n", addr, addrType)
    }
}
```

## API Reference

### Client Creation

```go
// Use default EventStore program ID
client := eventstore.NewClient(rpcClient, nil)

// Use custom program ID
programID := solana.MustPublicKeyFromBase58("your-program-id")
client := eventstore.NewClient(rpcClient, &programID)
```

### Fetching Events

#### Get All Events (Recommended)
```go
// Fetches all events from all shards in a single RPC call
allEvents, err := client.GetAllEvents(ctx)
```

#### Get Specific Event Types
```go
// Get only ZenBTC wrap events
zenbtcWraps, err := client.GetZenbtcWrapEvents(ctx)

// Get only ZenBTC unwrap events (with Bitcoin addresses)
zenbtcUnwraps, err := client.GetZenbtcUnwrapEvents(ctx)

// Get only Rock wrap events
rockWraps, err := client.GetRockWrapEvents(ctx)

// Get only Rock unwrap events (with Bitcoin addresses)
rockUnwraps, err := client.GetRockUnwrapEvents(ctx)
```

#### Get Event Statistics
```go
counts, err := client.GetEventCounts(ctx)
// Returns: map[string]int{
//     "zenbtc_wrap": 245,
//     "zenbtc_unwrap": 187,
//     "rock_wrap": 156,
//     "rock_unwrap": 203,
//     "total": 791
// }
```

### Event Structures

#### Wrap Events (ZenBTC & Rock)
```go
type TokensMintedWithFee struct {
    Recipient solana.PublicKey
    Value     uint64  // Amount in smallest unit
    Fee       uint64  // Fee in smallest unit
    Mint      solana.PublicKey
    ID        [16]uint8  // 128-bit event ID
}

// Helper method
eventID := event.GetID()  // Returns uint64
```

#### Unwrap Events (ZenBTC & Rock)
```go
type ZenbtcTokenRedemption struct {
    Redeemer solana.PublicKey
    Value    uint64
    DestAddr FlexibleAddress  // Bitcoin address
    Fee      uint64
    Mint     solana.PublicKey
    ID       [16]uint8
}

// Helper methods
eventID := event.GetID()
bitcoinAddr := event.GetBitcoinAddress()  // Returns decoded string
```

### Bitcoin Address Utilities

```go
// Get human-readable address type
addrType := eventstore.GetBitcoinAddressType("bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh")
// Returns: "P2WPKH (Bech32)"

// Supported address types:
// - P2PKH (Legacy)           - starts with 1
// - P2SH                     - starts with 3
// - P2WPKH (Bech32)          - bc1q, 42 chars
// - P2WSH (Bech32)           - bc1q, 62 chars
// - P2TR (Taproot/Bech32m)   - bc1p, 62 chars
// - Testnet variants         - m/n, 2, tb1q, tb1p
```

## Architecture

The EventStore program uses deterministic sharding to store up to 1000+ events per type:

- **ZenBTC Wrap**: 10 shards × 100 events = 1000 events
- **ZenBTC Unwrap**: 17 shards × 60 events = 1020 events
- **Rock Wrap**: 10 shards × 100 events = 1000 events
- **Rock Unwrap**: 17 shards × 60 events = 1020 events

The SDK automatically:
1. Calculates all shard addresses
2. Fetches all shards in a single `getMultipleAccounts` RPC call
3. Deserializes and aggregates events across shards
4. Decodes Bitcoin addresses from the FlexibleAddress format

## Performance

- **Single RPC Call**: All ~4000 events fetched in one request
- **Efficient Deserialization**: Borsh binary format with zero-copy where possible
- **Minimal Memory**: Events processed incrementally during deserialization
- **Parallel Safe**: Client is safe for concurrent use

## Error Handling

The SDK gracefully handles:
- Non-existent shard accounts (when program is newly deployed)
- Corrupted account data (skips individual shards)
- Network timeouts (standard RPC client behavior)
- Invalid Bitcoin addresses (returns "Unknown" type)

## Examples

See the `example/` directory for complete usage examples including:
- Basic event fetching
- Bitcoin address analysis
- Real-time monitoring
- JSON export
- Event filtering and statistics

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

For questions and support:
- Create an issue on GitHub
- Check the example code in `example/main.go`
- Review the EventStore program documentation

---

**EventStore Go SDK** - Efficient access to Bitcoin-Solana bridge events with automatic Bitcoin address decoding and type detection.
