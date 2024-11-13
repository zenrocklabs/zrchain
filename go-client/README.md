# go-client

A comprehensive gRPC client library for querying and interacting with `zenrockd`, the Zenrock blockchain daemon.

While the Cosmos SDK provides auto-generated protobuf clients, this package aims to be more ergonomic and developer-friendly by providing:

- Strongly typed interfaces for all operations
- Simplified transaction handling
- Built-in retries and error handling
- Convenient helper methods for common operations

## Key Components

### Core Client Files

- `raw_tx_client.go`: Base client implementation that handles:
  - Transaction building and signing
  - Gas estimation and fee calculation  
  - Broadcasting and confirmation
  - Identity and key management

### Module-Specific Clients

- `tx_treasury.go`: Treasury module operations including:
  - Key Request Management
    - Creating new key requests (standard and ZrSign)
    - Fulfilling key requests with public keys
    - Tracking request status
  - Signature Request Operations  
    - Creating signature requests for transactions
    - Managing ZrSign signature requests
    - Fulfilling requests with signatures
    - Rejecting invalid requests
  - Authorization
    - Generating and verifying party signatures
    - Managing request permissions

- `tx_zenbtc.go`: Bitcoin bridge functionality:
  - Deposit Verification
    - Validating Bitcoin transaction inclusion
    - Merkle proof verification
    - Amount and address validation
  - Cross-chain Transaction Management
    - Initiating Bitcoin transactions
    - Tracking transaction status
    - Managing Bitcoin addresses

## Usage

The clients are designed to be used together, with the specialized clients (`TreasuryTxClient`, `ZenBTCTxClient`) wrapping the base `RawTxClient` to provide domain-specific functionality while handling common operations like transaction building, signing and broadcasting automatically.

See individual client documentation for detailed usage examples and API references.
