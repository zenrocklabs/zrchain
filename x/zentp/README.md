# `x/zentp`

## Abstract

The following documents specify the zentp module.

The zentp module is a Cosmos SDK module in zrchain that handles token bridging operations between chains, specifically managing the minting and burning of tokens (particularly "Rock" tokens) across different blockchain networks. It provides functionality for bridging tokens from zrchain to destination chains, burning tokens from module accounts, and tracking mint and burn operations through dedicated storage.

## Concepts

The zentp module implements a token bridging system with the following key concepts:

- **Bridge Operations**: Each bridge operation (mint or burn) is represented by a `Bridge` type that tracks the source and destination chains, addresses, amounts, and operation status.
- **Operation States**: Bridge operations can be in one of several states: NEW, PENDING, COMPLETED, or FAILED, allowing for proper tracking of cross-chain operations.
- **Cross-Chain Bridging**: The module supports bridging tokens from zrchain to other chains (like Solana) using CAIP-2 chain identifiers and chain-specific address validation.
- **Governance Control**: Critical operations like burning tokens are gated through governance, ensuring proper control over token supply.

### Sidecar Integration

The zentp module relies on sidecar services to handle cross-chain operations, particularly for Solana integration. It continuously provides events form the solana ROCK contract and if zrchain sees an unhandled event it initiates the action. The sidecar services are responsible for:

The sidecar services provide the following functionality:

- **Mint Event Processing**: Provides mint events from Solana, enabling zrchain to initiate corresponding burn operations
- **Burn Event Processing**: Provides burn events from Solana, enabling zrchain to initiate corresponding mint operations
- **Account Information**: Supplies Solana durable nonce accounts and associated token accounts for recipients during transaction signing

This sidecar architecture enables a clean separation between on-chain state management and off-chain cross-chain operations, ensuring reliable and secure token bridging between zrchain and Solana.

## State

The `zentp` module uses the `collections` package which provides collection storage.

Here's the list of collections stored as part of the `zentp` module:

### Mint store

The `mintStore` stores `Bridge` operations: `Uint64(ID) -> ProtocolBuffer(Bridge)`.

### Burn store

The `burnStore` stores `Bridge` operations: `Uint64(ID) -> ProtocolBuffer(Bridge)`.

### Count stores

The module maintains two counter stores:

- `MintCount`: `Uint64` tracking the total number of mint operations
- `BurnCount`: `Uint64` tracking the total number of burn operations

## Msg Service

### Msg/UpdateParams

The module parameters can be updated with the `MsgUpdateParams` message, which is gated through governance.

```proto
message MsgUpdateParams {
  option (cosmos.msg.v1.signer) = "authority";
  string authority = 1;
  Params params = 2;
}
```

It's expected to fail if:

- the authority is not the module's authority address
- the parameters are invalid

### Msg/Bridge

A new bridge operation can be initiated with the `MsgBridge` message, which creates a mint request for Rock tokens on a destination chain.

```proto
message MsgBridge {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  string source_address = 2;
  uint64 amount = 3;
  string denom = 4;
  string destination_chain = 5;
  string recipient_address = 6;
}
```

It's expected to fail if:

- the creator is not authorized to perform the operation
- the source address is invalid
- the destination chain is not supported
- the recipient address is invalid for the destination chain
- the amount is zero or exceeds available balance

### Msg/Burn

Tokens can be burned from a module account with the `MsgBurn` message, which is gated through governance.

```proto
message MsgBurn {
  option (cosmos.msg.v1.signer) = "authority";
  string authority = 1;
  string module_account = 2;
  string denom = 3;
  uint64 amount = 4;
}
```

It's expected to fail if:

- the authority is not the module's authority address
- the module account is invalid
- the denomination is not supported
- the amount is zero or exceeds available balance

## Events

The zentp module emits the following events:

### EventBridge

| Type                             | Attribute Key  | Attribute Value                   |
| -------------------------------- | -------------  | --------------------------------  |
| message                          | action         | /zrchain.zentp.MsgBridge         |
| message                          | module         | zentp                            |
| bridge                           | id             | {bridge_id}                      |
| bridge                           | source_address | {source_address}                 |
| bridge                           | amount         | {amount}                         |
| bridge                           | denom          | {denom}                          |
| bridge                           | dest_chain     | {destination_chain}              |
| bridge                           | recipient      | {recipient_address}              |

### EventBurn

| Type                             | Attribute Key  | Attribute Value                   |
| -------------------------------- | -------------  | --------------------------------  |
| message                          | action         | /zrchain.zentp.MsgBurn           |
| message                          | module         | zentp                            |
| burn                             | authority      | {authority}                      |
| burn                             | module_account | {module_account}                 |
| burn                             | denom          | {denom}                          |
| burn                             | amount         | {amount}                         |

### EventUpdateParams

| Type                             | Attribute Key  | Attribute Value                   |
| -------------------------------- | -------------  | --------------------------------  |
| message                          | action         | /zrchain.zentp.MsgUpdateParams   |
| message                          | module         | zentp                            |
| update_params                    | authority      | {authority}                      |

## Parameters

The zentp module has the following parameters:

### Solana Configuration

The module includes Solana-specific parameters for token bridging operations:

| Key | Type | Description |
| --- | ---- | ----------- |
| `signer_key_id` | `uint64` | The key ID of the signer for Solana transactions |
| `program_id` | `string` | The Solana program ID for token operations |
| `nonce_account_key` | `uint64` | The key ID for the nonce account |
| `nonce_authority_key` | `uint64` | The key ID for the nonce authority |
| `mint_address` | `string` | The address of the token mint on Solana |
| `fee_wallet` | `string` | The wallet address for collecting fees |
| `fee` | `uint64` | The fee amount for operations |

## Client

### CLI

A user can query and interact with the `zentp` module using the CLI.

#### Query

The `query` commands allow users to query `zentp` state.

```bash
zenrockd query zentp --help
```

##### Params

The `params` command allows users to query the module parameters.

```bash
zenrockd query zentp params
```

Example Output:

```bash
params:
  solana:
    signer_key_id: "1"
    program_id: "program_id"
    nonce_account_key: "2"
    nonce_authority_key: "3"
    mint_address: "mint_address"
    fee_wallet: "fee_wallet"
    fee: "1000"
```

##### Burns

The `burns` command allows users to query burn operations.

```bash
zenrockd query zentp burns [id] [denom]
```

Example:

```bash
zenrockd query zentp burns 1 rock
```

#### Transactions

The `tx` commands allow users to interact with the `zentp` module.

```bash
zenrockd tx zentp --help
```

##### Bridge

The `bridge` command allows users to initiate a bridge operation to mint tokens on a destination chain.

```bash
zenrockd tx zentp bridge [amount] [denom] [source-address] [destination-chain] [recipient-address]
```

Example:

```bash
zenrockd tx zentp bridge 1000 rock zen1... solana recipient_address
```

### gRPC

A user can query the `zentp` module using gRPC endpoints.

#### Params

The `Params` endpoint allows users to query the module parameters.

```bash
zrchain.zentp.Query/Params
```

Example:

```bash
grpcurl -plaintext localhost:9090 zrchain.zentp.Query/Params
```

#### Burns

The `Burns` endpoint allows users to query burn operations.

```bash
zrchain.zentp.Query/Burns
```

Example:

```bash
grpcurl -plaintext -d '{"id": 1, "denom": "rock"}' localhost:9090 zrchain.zentp.Query/Burns
```

### REST

A user can query the `zentp` module using REST endpoints.

#### params

The `params` endpoint allows users to query the module parameters.

```bash
/zrchain/zentp/params
```

Example:

```bash
curl localhost:1317/zrchain/zentp/params
```

#### burns

The `burns` endpoint allows users to query burn operations.

```bash
/zrchain/zentp/burns/{id}/{denom}
```

Example:

```bash
curl localhost:1317/zrchain/zentp/burns/1/rock
```
