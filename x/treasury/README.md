# `x/treasury`

## Abstract

The following documents specify the treasury module.

Zenrock's Treasury Module is responsible for requesting and managing keys.

The treasury module facilitates a wide-range of applications - including cross-chain interaction - by fulfilling key and signature requests.

## Contents

* [Concepts](#concepts)
  * [Supported Keys](#supported-keys)
  * [Signatures](#signatures)
* [State](#state)
  * [KeyStore](#keystore)
  * [KeyRequestStore](#keyrequeststore)
  * [SignRequestStore](#signrequeststore)
  * [SignTransactionRequestStore](#signtransactionrequeststore)
  * [ICATransactionRequestStore](#icatransactionrequeststore)
* [Msg Service](#msg-service)
  * [Msg/NewKeyRequest](#msgnewkeyrequest)
  * [Msg/FulfilKeyRequest](#msgfulfilkeyrequest)
  * [Msg/NewSignatureRequest](#msgnewsignaturerequest)
  * [Msg/NewSignTransactionRequest](#msgnewsigntransactionrequest)
  * [Msg/FulfilSignatureRequest](#msgfulfilsignaturerequest)
  * [Msg/TransferFromKeyring](#msgtransferfromkeyring)
  * [Msg/NewICATransactionRequest](#msgnewicatransactionrequest)
  * [Msg/FulfilICATransactionRequest](#msgfulfilicatransactionrequest)
  * [Msg/UpdateKeyPolicy](#msgupdatekeypolicy)
  * [Msg/NewZrSignSignatureRequest](#msgnewzrsignsignaturerequest)
* [Events](#events)
  * [EventNewKeyRequest](#eventnewkeyrequest)
  * [EventKeyRequestFulfilled](#eventkeyrequestfulfilled)
  * [EventKeyRequestRejected](#eventkeyrequestrejected)
  * [EventNewSignRequest](#eventnewsignrequest)
  * [EventSignRequestFulfilled](#eventsignrequestfulfilled)
  * [EventSignRequestRejected](#eventsignrequestrejected)
  * [EventNewICATransactionRequest](#eventnewicatransactionrequest)
  * [EventICATransactionRequestFulfilled](#eventicatransactionrequestfulfilled)
* [Client](#client)
  * [CLI](#cli)
  * [gRPC](#grpc)
  * [REST](#rest)

## Concepts

### Supported Keys

Zenrock has the capacity to generate ecdsa secp256k1 and eddsa ed25519 keys. Key requests are processed off-chain by registered Keyrings, which subsequently store generated keys on-chain.

Keys can be used to derive valid EVM, Cosmos, Solana, and Bitcoin addresses.

Zenrock-generated addresses behave like standard self-hosted wallets and are able to interact with the relevant network natively. This exposes users to a broad range of the most popular networks.

### Signatures

Zenrock provides a signature request service that allows users to either request arbitrary signatures or request for transactions to be signed.

The transaction signing service may include broadcasting the transaction to the relevant layer one network as necessary.

The web application allows users to either manually define a transaction payload or use services like WalletConnect to submit unsigned transactions/messages to the Zenrock system.

## State

The `treasury` module uses the `collections` package which provides collection storage.

Here's the list of collections stored as part of the `treasury` module.

### KeyStore

The `KeyStore` stores `Key`: `BigEndian(KeyId) -> ProtocolBuffer(Key)`.

### KeyRequestStore

The `KeyRequestStore` stores `KeyRequest`: `BigEndian(KeyRequestId) -> ProtocolBuffer(KeyRequest)`.

### SignRequestStore

The `SignRequestStore` stores `SignRequest`: `BigEndian(SignRequestId) -> ProtocolBuffer(SignRequest)`.

### SignTransactionRequestStore

The `SignTransactionRequestStore` stores `BigEndian(SignTransactionRequestId) -> ProtocolBuffer(SignTransactionRequest)`.

### ICATransactionRequestStore

The `ICATransactionRequestStore` stores `BigEndian(ICATransactionRequestId) -> ProtocolBuffer(ICATransactionRequest)`.

## Msg Service

### Msg/NewKeyRequest

A new key can be requested with the `MsgNewKeyRequest` message.

```proto
message MsgNewKeyRequest {
  option (cosmos.msg.v1.signer) = "creator";
  string creator        = 1;
  string workspace_addr = 2;
  string keyring_addr   = 3;
  string key_type       = 4;
  uint64 btl            = 5;
  uint64 index          = 6;
  string ext_requester  = 7;
  uint64 ext_key_type   = 8;
  uint64 sign_policy_id = 9;
}
```

It's expected to fail if

* the workspace is not found
* the policy participants are not owners of the workspace
* the keyring is not active
* the transaction creators balance is to low to pay the keyring key creation fee
* the key type is not ecdsa, ed25519, bitcoin

### Msg/FulfilKeyRequest

A key request can be fulfilled with the `MsgFulfilKeyRequest` message, this message is used by the keyrings to reply to a keyrequest.

```proto
message MsgFulfilKeyRequest {
  option (cosmos.msg.v1.signer) = "creator";
  string           creator    = 1;
  uint64           request_id = 2;
  KeyRequestStatus status     = 3;

  // Holds the result of the request. If status is approved, the result will
  // contain the requested key's public key that can be used for signing
  // payloads.
  // If status is rejected, the result will contain the reason.
  oneof result {
    MsgNewKey key           = 4;
    string    reject_reason = 5;
  }
  bytes keyring_party_signature = 6;
}

message MsgNewKey {
  bytes public_key = 1;
}
```

It's expected to fail if

* the key request is not found
* the keyring is not active
* the transaction creator is not a keyring party
* the key request is already completed
* the keyring signature is not 64 bytes
* the public key has an invalid size for the type

### Msg/NewSignatureRequest

A new signature can be requested with the `MsgNewSignatureRequest` message.

```proto
message MsgNewSignatureRequest {
  option (cosmos.msg.v1.signer) = "creator";
  string              creator                     = 1;
  repeated uint64     key_ids                     = 2;
  string              data_for_signing            = 3;
  uint64              btl                         = 4;
  bytes               cache_id                    = 5;
  bytes               verify_signing_data         = 6;
  VerificationVersion verify_signing_data_version = 7;
  uint64              mpc_btl                     = 8;
  bytes               zenbtc_tx_bytes             = 9;  // Optional
}

enum VerificationVersion {
  UNKNOWN      = 0;
  BITCOIN_PLUS = 1;
}
```

It's expected to fail if

* the payload is not a 32 byte hash in hex format
* the key or its workspace can not be found
* the keyring to which the key belongs to is not found or not active
* in case of a bitcoin signature request, that the tx in verify_signing_data matches the hash in data_for_signing
* the transaction creators balance is too low to pay the keyring key creation fee
* the number of key IDs doesn't match the number of data elements
* the request is for a zenBTC deposit key but not from the Bitcoin proxy service

### Msg/NewSignTransactionRequest

A new sign transaction can be requested with the `MsgNewSignTransactionRequest` message.

```proto
message MsgNewSignTransactionRequest {
  option (cosmos.msg.v1.signer) = "creator";
  string     creator              = 1;
  uint64     key_id               = 2;
  WalletType wallet_type          = 3;
  bytes      unsigned_transaction = 4;
  
  // Additional metadata required when parsing the unsigned transaction.
  google.protobuf.Any metadata = 5;
  uint64              btl      = 6;
  bytes               cache_id = 7;
  bool                no_broadcast = 8;
  uint64              mpc_btl = 9;
}

message MetadataEthereum {
  uint64 chain_id = 1;
}

// Define an enum for Solana network types
enum SolanaNetworkType {
  UNDEFINED = 0;
  MAINNET = 1;
  DEVNET = 2;
  TESTNET = 3;
}

message MetadataSolana {
  SolanaNetworkType network = 1;
  string mintAddress = 2;
}
```

It's expected to fail if

* the key or its workspace can not be found
* the keyring to which the key belongs to is not found or not active
* the transaction cannot be parsed (depends on the wallet type)
* the metadata is invalid (depends on the wallet type)
* the transaction creators balance is too low to pay the keyring key creation fee
* the data for signing is empty
* the request is for a zenBTC deposit key but not from the Bitcoin proxy service

### Msg/FulfilSignatureRequest

A signature request can be fulfilled with the `MsgFulfilSignatureRequest` message, this message is used by the keyrings to reply to a signature request.

```proto
message MsgFulfilSignatureRequest {
  option (cosmos.msg.v1.signer) = "creator";
  string            creator                 = 1;
  uint64            request_id              = 2;
  SignRequestStatus status                  = 3;
  bytes             keyring_party_signature = 4;
  bytes             signed_data             = 5;
  string            reject_reason           = 6;
}
```

It's expected to fail if

* the request is not found
* the signature request is already completed
* signed data is missing or contains an invalid signature
* the keyring signature is not 64 bytes
* the keyring is not found or inactive
* the signed data contains an invalid signature for the data in the sign request
* the transaction creator is not a keyring party

### Msg/TransferFromKeyring

Funds that have accumulated in a keyring as a result of key and sign requests can be transfered to an account with the `MsgTransferFromKeyring` message.

```proto
message MsgTransferFromKeyring {
  option (cosmos.msg.v1.signer) = "creator";
  string creator   = 1;
  string keyring   = 2;
  string recipient = 3;
  uint64 amount    = 4;
  string denom     = 5;
}
```

It's expected to fail if

* the keyring is not found
* the transaction creator or recepient are not a keyring admin
* the recipient is not a valid address
* the requested amount is higher than the keyring's balance

### Msg/UpdateKeyPolicy

The sign policy on a key can be updated with the `MsgUpdateKeyPolicy` message.

```proto
message MsgUpdateKeyPolicy {
  option (cosmos.msg.v1.signer) = "creator";
  string creator        = 1;
  uint64 key_id         = 2;
  uint64 sign_policy_id = 3;
}
```

It's expected to fail if

* the key or its workspace is not found
* the new policy's participants are not owners of the workspace
* the transaction creator is not authorized to modify the key policy
* the specified policy ID does not exist

### Msg/NewZrSignSignatureRequest

A new ZrSign signature can be requested with the `MsgNewZrSignSignatureRequest` message.

```proto
message MsgNewZrSignSignatureRequest {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  string address = 2;
  uint64 key_type = 3;
  uint64 wallet_index = 4;
  bytes cache_id = 5;
  string data = 6;
  bytes verify_signing_data = 7;
  VerificationVersion verify_signing_data_version = 8;
  WalletType wallet_type = 9;
  google.protobuf.Any metadata = 10;
  bool no_broadcast = 11;
  uint64 btl = 12;
  bool tx = 13;
}
```

It's expected to fail if

* the ZrSign workspace cannot be found for the address
* the data for signing is invalid
* the wallet type is not supported
* the transaction verification fails
* the creator's balance is too low to pay the required fees

### Msg/NewICATransactionRequest

A new ICA transaction can be requested with the `MsgNewICATransactionRequest` message.

```proto
message MsgNewICATransactionRequest {
  option (cosmos.msg.v1.signer) = "creator";
  string creator                    = 1;
  uint64 key_id                     = 2;
  string input_payload              = 3;
  string connection_id              = 4;
  uint64 relative_timeout_timestamp = 5;
  uint64 btl                        = 6;
}
```

It's expected to fail if

* the key or its workspace cannot be found
* the keyring to which the key belongs is not found or not active
* the connection ID is invalid
* the transaction creator's balance is too low to pay the keyring fee
* the request is for a zenBTC deposit key but not from the Bitcoin proxy service

### Msg/FulfilICATransactionRequest

A ICA transaction request can be fulfilled with the `MSgFulfilICATransactionRequest` message.

```proto
message MsgFulfilICATransactionRequest {
  option (cosmos.msg.v1.signer) = "creator";
  string            creator                 = 1;
  uint64            request_id              = 2;
  SignRequestStatus status                  = 3;
  bytes             keyring_party_signature = 4;
  bytes             signed_data             = 5;
  string            reject_reason           = 6;
}
```

It's expected to fail if

* the request is not found
* the ICA transaction request is already completed
* signed data is missing or contains an invalid signature
* the keyring signature is not 64 bytes
* the keyring is not found or inactive
* the transaction creator is not a keyring party

## Events

The treasury module emits the following events:

### EventNewKeyRequest

| Type                             | Attribute Key | Attribute Value                  |
| -------------------------------- | ------------- | -------------------------------- |
| new_key_request                  | request_id    | {request_id}                     |

### EventKeyRequestFulfilled

| Type                             | Attribute Key | Attribute Value                  |
| -------------------------------- | ------------- | -------------------------------- |
| key_request_fulfilled            | request_id    | {request_id}                     |

### EventKeyRequestRejected

| Type                             | Attribute Key | Attribute Value                  |
| -------------------------------- | ------------- | -------------------------------- |
| key_request_rejected             | request_id    | {request_id}                     |

### EventNewSignRequest

| Type                             | Attribute Key | Attribute Value                  |
| -------------------------------- | ------------- | -------------------------------- |
| new_sign_request                 | request_id    | {request_id}                     |

### EventSignRequestFulfilled

| Type                             | Attribute Key | Attribute Value                  |
| -------------------------------- | ------------- | -------------------------------- |
| sign_request_fulfilled           | request_id    | {request_id}                     |

### EventSignRequestRejected

| Type                             | Attribute Key | Attribute Value                  |
| -------------------------------- | ------------- | -------------------------------- |
| sign_request_rejected            | request_id    | {request_id}                     |

### EventNewICATransactionRequest

| Type                             | Attribute Key | Attribute Value                  |
| -------------------------------- | ------------- | -------------------------------- |
| new_key_request                  | request_id    | {request_id}                     |

### EventICATransactionRequestFulfilled

| Type                              | Attribute Key | Attribute Value                  |
| --------------------------------- | ------------- | -------------------------------- |
| ica_transaction_request_fulfilled | request_id    | {request_id}                     |

## Parameters

### KeyringCommission

The percentage of the keyring sign commission that will be sent to the KeyringCommissionDestination.

### KeyringCommissionDestination

The destination address to receive the keyring sign commission percentage.

### MinGasFee

The minimum gas fee to be paid on the network.

### MPCKeyring

The MPC Keyring address that will be used by ZrSign.

### DefaultBTL

The default BTL (Blocks To Live) for MPC key and signature requests.

## Client

### CLI

A user can query and interact with the `treasury` module using the CLI.

#### Query

The `query` commands allow users to query `treasury` state.

```bash
zenrockd query treasury --help
```

##### key-requests

The `key-requests` command allows users to query all key requests.

```bash
zenrockd query treasury key-requests
```

Example:

```bash
zenrockd query treasury key-requests
```

Example Output:

```yaml
key_requests:
- creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  id: "1"
  key_type: KEY_TYPE_ECDSA_SECP256K1
  keyring_addr: keyring1pfnq7r04rept47gaf5cpdew2
  status: KEY_REQUEST_STATUS_PENDING
  workspace_addr: workspace1mphgzyhncnzyggfxmv4nmh
pagination:
  total: "1"
```

##### key-request-by-id

The `key-request-by-id` command allows users to query a key request by id.

```bash
zenrockd query treasury key-request-by-id --id [id]
```

Example:

```bash
zenrockd query treasury key-request-by-id --id 1
```

Example Output:

```yaml
key_request:
  creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  id: "1"
  key_type: KEY_TYPE_ECDSA_SECP256K1
  keyring_addr: keyring1pfnq7r04rept47gaf5cpdew2
  status: KEY_REQUEST_STATUS_PENDING
  workspace_addr: workspace1mphgzyhncnzyggfxmv4nmh
```

##### keys

The `keys` command allows users to query all keys.

```bash
zenrockd query treasury keys
```

Example:

```bash
zenrockd query treasury keys
```

Example Output:

```yaml
keys:
- key:
    id: "1"
    keyring_addr: keyring1pfnq7r04rept47gaf5cpdew2
    public_key: AjtNNCpqUnzKjC2Sas52NqGfCSUodOI2j+ejJdUO/+2z
    type: KEY_TYPE_ECDSA_SECP256K1
    workspace_addr: workspace1mphgzyhncnzyggfxmv4nmh
  wallets:
  - address: zen1qfgae6jr3xyawcz7uuh7ez06eteksxgfhdk7uf
    type: WALLET_TYPE_NATIVE
  - address: 0xCF4324b9FD2fBC83b98a1483E854D31D3E45944C
    type: WALLET_TYPE_EVM
pagination:
  total: "1"
```

##### key-by-id

The `key-by-id` command allows users to query a key by id.

```bash
zenrockd query treasury key-by-id [id]
```

Example:

```bash
zenrockd query treasury key-by-id 1
```

Example Output:

```yaml
key:
  id: "1"
  keyring_addr: keyring1pfnq7r04rept47gaf5cpdew2
  public_key: AjtNNCpqUnzKjC2Sas52NqGfCSUodOI2j+ejJdUO/+2z
  type: KEY_TYPE_ECDSA_SECP256K1
  workspace_addr: workspace1mphgzyhncnzyggfxmv4nmh
wallets:
- address: zen1qfgae6jr3xyawcz7uuh7ez06eteksxgfhdk7uf
  type: WALLET_TYPE_NATIVE
- address: 0xCF4324b9FD2fBC83b98a1483E854D31D3E45944C
  type: WALLET_TYPE_EVM
```

##### signature-requests

The `signature-requests` command allows users to query all signature requests.

```bash
zenrockd query treasury signature-requests
```

Example:

```bash
zenrockd query treasury signature-requests
```

Example Output:

```yaml
pagination:
  total: "1"
sign_requests:
- creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  data_for_signing:
  - FYST...
  id: "1"
  key_id: "1"
  key_type: KEY_TYPE_ECDSA_SECP256K1
  keyring_party_signatures:
  - C3JV...
  metadata:
    type: /zrchain.treasury.MetadataEthereum
    value:
      chain_id: "11155111"
  signed_data:
  - sign_request_id: "1"
    signed_data: +L8bs...
  status: SIGN_REQUEST_STATUS_FULFILLED
```

##### signature-request-by-id

The `signature-request-by-id` command allows users to query a signature request by id.

```bash
zenrockd query treasury signature-request-by-id [id]
```

Example:

```bash
zenrockd query treasury signature-request-by-id 1
```

Example Output:

```yaml
sign_request:
  creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  data_for_signing:
  - FYST...
  id: "1"
  key_id: "1"
  key_type: KEY_TYPE_ECDSA_SECP256K1
  keyring_party_signatures:
  - C3JV...
  metadata:
    type: /zrchain.treasury.MetadataEthereum
    value:
      chain_id: "11155111"
  signed_data:
  - sign_request_id: "1"
    signed_data: +L8bs...
  status: SIGN_REQUEST_STATUS_FULFILLED
```

##### sign-transaction-requests

The `sign-transaction-requests` command allows users to query all sign transaction requests.

```bash
zenrockd query treasury sign-transaction-requests 
```

Example:

```bash
zenrockd query treasury sign-transaction-requests
```

Example Output:

```yaml
pagination:
  total: "1"
sign_transaction_requests:
- sign_request:
    creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
    data_for_signing:
    - FYST...
    id: "2"
    key_id: "1"
    key_type: KEY_TYPE_ECDSA_SECP256K1
    keyring_party_signatures:
    - wnD4z...
    metadata:
      type: /zrchain.treasury.MetadataEthereum
      value:
        chain_id: "11155111"
    signed_data:
    - sign_request_id: "2"
      signed_data: +L8bs...
    status: SIGN_REQUEST_STATUS_FULFILLED
  sign_transaction_requests:
    creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
    id: "1"
    key_id: "1"
    sign_request_id: "2"
    unsigned_transaction: +MqAh...
    wallet_type: WALLET_TYPE_EVM
```

##### sign-transaction-request-by-id

The `sign-transaction-request-by-id` command allows users to query a single sign transaction request by id.

```bash
zenrockd query treasury sign-transaction-request-by-id [id]
```

Example:

```bash
zenrockd query treasury sign-transaction-request-by-id 1
```

Example Output:

```yaml
sign_transaction_request:
  creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  id: "1"
  key_id: "1"
  sign_request_id: "2"
  unsigned_transaction: +MqAh...
  wallet_type: WALLET_TYPE_EVM
```

### Transactions

The `tx` commands allow users to interact with the `treasury` module.

```bash
zenrockd tx treasury --help
```

#### new-key-request

The `new-key-request` command allows users to request a new key of a specific type.

```bash
zenrockd tx treasury new-key-request [workspace-addr] [keyring-addr] [secp256k1|ed25519|bitcoin] --btl [btl] --sign-policy-id [policy_id]
```

Example:

```bash
zenrockd tx treasury new-key-request workspace1mphgzyhncnzyggfxmv4nmh keyring1pfnq7r04rept47gaf5cpdew2 ecdsa
```

#### new-signature-request

The `new-signature-request` command allows users to request a new signature for a hash.

```bash
zenrockd tx treasury new-signature-request [key-id] [data-for-signing] --btl [btl] --cache-id [cache-id]
```

Example:

```bash
zenrockd tx treasury new-signature-request 1 824c611eff75402dcc6a3b7982006ad48b3239da2350e8a6337b0c737c6dc8d0
```

#### new-sign-transaction-request

The `new-sign-transaction-request` command allows users to request a new signature for a transaction.

```bash
zenrockd tx treasury new-sign-transaction-request [key-id] [wallet-type] [unsigned-tx] --metadata [metadata] --btl [btl] --cache-id [cache-id]
```

Example:

```bash
zenrockd tx treasury new-sign-transaction-request 1 evm 'f8ca...70c64' --metadata '{ "@type": "/zrchain.treasury.MetadataEthereum", "chain_id": "11155111" }'
```

#### update-key-policy

The `update-key-policy` command allows users to udpate the signing policy on a key.

```bash
zenrockd tx treasury update-key-policy [key-id] [sign_policy_id] 
```

Example:

```bash
zenrockd tx treasury update-key-policy 1 1
```

#### transfer-from-keyring

The `transfer-from-keyring` command allows users to transfer funds from the keyring to an admin.

```bash
zenrockd tx treasury transfer-from-keyring [keyring] [recipient] [amount] [denom]
```

Example:

```bash
zenrockd tx treasury transfer-from-keyring keyring1w887ucurq2nmnj5mq5uaju6a zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty 10 urock
```

### gRPC

A user can query the `treasury` module using gRPC endpoints.

#### Keys

The `Keys` endpoint allows users to query all keys.

```bash
zrchain.treasury.Query/Keys
```

Example:

```bash
grpcurl -plaintext localhost:9090 zrchain.treasury.Query/Keys
```

Example Output:

```json
{
  "keys": [
    {
      "key": {
        "id": "1",
        "workspaceAddr": "workspace1mphgzyhncnzyggfxmv4nmh",
        "keyringAddr": "keyring1pfnq7r04rept47gaf5cpdew2",
        "type": "KEY_TYPE_ECDSA_SECP256K1",
        "publicKey": "AjtNNCpqUnzKjC2Sas52NqGfCSUodOI2j+ejJdUO/+2z"
      },
      "wallets": [
        {
          "address": "zen1qfgae6jr3xyawcz7uuh7ez06eteksxgfhdk7uf",
          "type": "WALLET_TYPE_NATIVE"
        },
        {
          "address": "0xCF4324b9FD2fBC83b98a1483E854D31D3E45944C",
          "type": "WALLET_TYPE_EVM"
        }
      ]
    }
  ],
  "pagination": {
    "total": "1"
  }
}
```

#### KeyById

The `KeyById` endpoint allows users to query for a single key by id.

```bash
zrchain.treasury.Query/KeyByID
```

Example:

```bash
grpcurl -plaintext \
    -d '{"id":1}' localhost:9090 zrchain.treasury.Query/KeyByID
```

Example Output:

```json
{
  "key": {
    "id": "1",
    "workspaceAddr": "workspace1mphgzyhncnzyggfxmv4nmh",
    "keyringAddr": "keyring1pfnq7r04rept47gaf5cpdew2",
    "type": "KEY_TYPE_ECDSA_SECP256K1",
    "publicKey": "AjtNNCpqUnzKjC2Sas52NqGfCSUodOI2j+ejJdUO/+2z"
  },
  "wallets": [
    {
      "address": "zen1qfgae6jr3xyawcz7uuh7ez06eteksxgfhdk7uf",
      "type": "WALLET_TYPE_NATIVE"
    },
    {
      "address": "0xCF4324b9FD2fBC83b98a1483E854D31D3E45944C",
      "type": "WALLET_TYPE_EVM"
    }
  ]
}
```

#### KeyRequests

The `KeyRequests` endpoint allows users to query for all key requests.

```bash
zrchain.treasury.Query/KeyRequests
```

Example:

```bash
grpcurl -plaintext localhost:9090 zrchain.treasury.Query/KeyRequests
```

Example Output:

```json
{
  "keyRequests": [
    {
      "id": "1",
      "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
      "workspaceAddr": "workspace1mphgzyhncnzyggfxmv4nmh",
      "keyringAddr": "keyring1pfnq7r04rept47gaf5cpdew2",
      "keyType": "KEY_TYPE_ECDSA_SECP256K1",
      "status": "KEY_REQUEST_STATUS_FULFILLED",
      "keyringPartySignatures": [
        "C3JV..."
      ]
    }
  ],
  "pagination": {
    "total": "1"
  }
}
```

#### KeyRequestById

The `KeyRequestById` endpoint allows users to query for a single key request by id.

```bash
zrchain.treasury.Query/KeyRequestByID
```

Example:

```bash
grpcurl -plaintext \
    -d '{"id":1}' localhost:9090 zrchain.treasury.Query/KeyRequestByID
```

Example Output:

```json
{
  "keyRequest": {
    "id": "1",
    "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
    "workspaceAddr": "workspace1mphgzyhncnzyggfxmv4nmh",
    "keyringAddr": "keyring1pfnq7r04rept47gaf5cpdew2",
    "keyType": "KEY_TYPE_ECDSA_SECP256K1",
    "status": "KEY_REQUEST_STATUS_FULFILLED",
    "keyringPartySignatures": [
      "C3JV..."
    ]
  }
}
```

#### SignTransactionRequests

The `SignTransactionRequests` endpoint allows users to query for all sign transaction requests.

```bash
zrchain.treasury.Query/SignTransactionRequests
```

Example:

```bash
grpcurl -plaintext localhost:9090 zrchain.treasury.Query/SignTransactionRequests
```

Example Output:

```json
{
  "signTransactionRequests": [
    {
      "signTransactionRequests": {
        "id": "1",
        "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
        "keyId": "1",
        "walletType": "WALLET_TYPE_EVM",
        "unsignedTransaction": "+MqAh...",
        "signRequestId": "3"
      },
      "signRequest": {
        "id": "3",
        "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
        "keyId": "1",
        "keyType": "KEY_TYPE_ECDSA_SECP256K1",
        "dataForSigning": [
          "FYST..."
        ],
        "status": "SIGN_REQUEST_STATUS_PENDING",
        "metadata": {
          "@type": "/zrchain.treasury.MetadataEthereum",
          "chainId": "11155111"
        }
      }
    }
  ],
  "pagination": {
    "total": "1"
  }
}
```

#### SignTransactionRequestById

The `SignTransactionRequestById` endpoint allows users to query for a single sign transaction request by id.

```bash
zrchain.treasury.Query/SignTransactionRequestByID
```

Example:

```bash
grpcurl -plaintext \
    -d '{"id":1}' localhost:9090 zrchain.treasury.Query/SignTransactionRequestByID
```

Example Output:

```json
{
  "signTransactionRequest": {
    "id": "1",
    "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
    "keyId": "1",
    "walletType": "WALLET_TYPE_EVM",
    "unsignedTransaction": "+MqA...",
    "signRequestId": "3"
  }
}
```

#### SignatureRequests

The `SignatureRequests` endpoint allows users to query for all signature requests.

```bash
zrchain.treasury.Query/SignatureRequests
```

Example:

```bash
grpcurl -plaintext localhost:9090 zrchain.treasury.Query/SignatureRequests
```

Example Output:

```json
{
  "signRequests": [
    {
      "id": "1",
      "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
      "keyId": "1",
      "keyType": "KEY_TYPE_ECDSA_SECP256K1",
      "dataForSigning": [
        "FYST..."
      ],
      "status": "SIGN_REQUEST_STATUS_FULFILLED",
      "signedData": [
        {
          "signRequestId": "1",
          "signedData": "+L8b..."
        }
      ],
      "keyringPartySignatures": [
        "utCB..."
      ],
      "metadata": {
        "@type": "/zrchain.treasury.MetadataEthereum",
        "chainId": "11155111"
      }
    }
  ],
  "pagination": {
    "total": "1"
  }
}
```

#### SignatureRequestsById

The `SignatureRequestsById` endpoint allows users to query for a single signature request by id.

```bash
zrchain.treasury.Query/SignatureRequestByID
```

Example:

```bash
grpcurl -plaintext \
    -d '{"id":1}' localhost:9090 zrchain.treasury.Query/SignatureRequestByID
```

Example Output:

```json
{
  "signRequest": {
    "id": "1",
    "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
    "keyId": "1",
    "keyType": "KEY_TYPE_ECDSA_SECP256K1",
    "dataForSigning": [
      "FYST..."
    ],
    "status": "SIGN_REQUEST_STATUS_FULFILLED",
    "signedData": [
      {
        "signRequestId": "1",
        "signedData": "+L8bs..."
      }
    ],
    "keyringPartySignatures": [
      "C3JV..."
    ],
    "metadata": {
      "@type": "/zrchain.treasury.MetadataEthereum",
      "chainId": "11155111"
    }
  }
}
```

### REST

A user can query the `treasury` module using REST endpoints.

#### keys

The `keys` endpoint allows users to query all keys in a workspace.

```bash
/zrchain/treasury/keys/{workspace_addr}/{wallet_type}/{prefixes}
```

Example:

```bash
curl localhost:1317/zrchain/treasury/keys/workspace1mphgzyhncnzyggfxmv4nmh/WALLET_TYPE_UNSPECIFIED/zen,cosmos
```

Example Output:

```json
{
  "keys": [
    {
      "key": {
        "id": "1",
        "workspace_addr": "workspace1mphgzyhncnzyggfxmv4nmh",
        "keyring_addr": "keyring1pfnq7r04rept47gaf5cpdew2",
        "type": "KEY_TYPE_ECDSA_SECP256K1",
        "public_key": "AjtNNCpqUnzKjC2Sas52NqGfCSUodOI2j+ejJdUO/+2z",
        "index": "0",
        "sign_policy_id": "0"
      },
      "wallets": [
        {
          "address": "zen1qfgae6jr3xyawcz7uuh7ez06eteksxgfhdk7uf",
          "type": "WALLET_TYPE_NATIVE"
        },
        {
          "address": "0xCF4324b9FD2fBC83b98a1483E854D31D3E45944C",
          "type": "WALLET_TYPE_EVM"
        }
      ]
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "1"
  }
}
```

#### key_by_id

The `key_by_id` endpoint allows users to query a specific key by id.

```bash
/zrchain/treasury/key_by_id/{id}/{wallet_type}/{prefixes}
```

Example:

```bash
curl localhost:1317/zrchain/treasury/key_by_id/1/WALLET_TYPE_UNSPECIFIED/zen,cosmos
```

Example Output:

```json
{
  "key": {
    "id": "1",
    "workspace_addr": "workspace1mphgzyhncnzyggfxmv4nmh",
    "keyring_addr": "keyring1pfnq7r04rept47gaf5cpdew2",
    "type": "KEY_TYPE_ECDSA_SECP256K1",
    "public_key": "AjtNNCpqUnzKjC2Sas52NqGfCSUodOI2j+ejJdUO/+2z",
    "index": "0",
    "sign_policy_id": "0"
  },
  "wallets": [
    {
      "address": "zen1qfgae6jr3xyawcz7uuh7ez06eteksxgfhdk7uf",
      "type": "WALLET_TYPE_NATIVE"
    },
    {
      "address": "cosmos1qfgae6jr3xyawcz7uuh7ez06eteksxgfdqylsc",
      "type": "WALLET_TYPE_NATIVE"
    },
    {
      "address": "0xCF4324b9FD2fBC83b98a1483E854D31D3E45944C",
      "type": "WALLET_TYPE_EVM"
    }
  ]
}
```

#### key_requests

The `key_requests` endpoint allows users to query key requests for a specific keyring by status and workspace.

```bash
/zrchain/treasury/key_requests/{keyring_addr}/{status}/{workspace_addr}
```

Example:

```bash
curl localhost:1317/zrchain/treasury/key_requests/keyring1pfnq7r04rept47gaf5cpdew2/KEY_REQUEST_STATUS_UNSPECIFIED/
```

Example Output:

```json
{
  "key_requests": [
    {
      "id": "1",
      "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
      "workspace_addr": "workspace1mphgzyhncnzyggfxmv4nmh",
      "keyring_addr": "keyring1pfnq7r04rept47gaf5cpdew2",
      "key_type": "KEY_TYPE_ECDSA_SECP256K1",
      "status": "KEY_REQUEST_STATUS_FULFILLED",
      "keyring_party_signatures": [
        "C3JV..."
      ],
      "reject_reason": "",
      "index": "0",
      "sign_policy_id": "0"
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "1"
  }
}
```

#### key_request_by_id

The `key_request_by_id` endpoint allows users to query a specific key request by id.

```bash
/zrchain/treasury/key_request_by_id/{id}
```

Example:

```bash
curl localhost:1317/zrchain/treasury/key_request_by_id/1
```

Example Output:

```json
{
  "key_request": {
    "id": "1",
    "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
    "workspace_addr": "workspace1mphgzyhncnzyggfxmv4nmh",
    "keyring_addr": "keyring1pfnq7r04rept47gaf5cpdew2",
    "key_type": "KEY_TYPE_ECDSA_SECP256K1",
    "status": "KEY_REQUEST_STATUS_FULFILLED",
    "keyring_party_signatures": [
      "C3JV..."
    ],
    "reject_reason": "",
    "index": "0",
    "sign_policy_id": "0"
  }
}
```

#### signature_requests

The `signature_requests` endpoint allows users to query all signature requests for a specific keyring ans status.

```bash
/zrchain/treasury/signature_requests/{keyring_addr}/{status}
```

Example:

```bash
curl localhost:1317/zrchain/treasury/signature_requests/keyring1pfnq7r04rept47gaf5cpdew2/SIGN_REQUEST_STATUS_UNSPECIFIED 
```

Example Output:

```json
{
  "sign_requests": [
    {
      "id": "1",
      "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
      "key_id": "1",
      "key_type": "KEY_TYPE_ECDSA_SECP256K1",
      "data_for_signing": [
        "FYST..."
      ],
      "status": "SIGN_REQUEST_STATUS_FULFILLED",
      "signed_data": [
        {
          "sign_request_id": "1",
          "signed_data": "+L8bs..."
        }
      ],
      "keyring_party_signatures": [
        "C3JV..."
      ],
      "reject_reason": "",
      "metadata": {
        "@type": "/zrchain.treasury.MetadataEthereum",
        "chain_id": "11155111"
      },
      "parent_req_id": "0",
      "child_req_ids": [],
      "cache_id": null
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "1"
  }
}
```

#### signature_request_by_id

The `signature_request_by_id` endpoint allows users to query a specific signature request by id.

```bash
/zrchain/treasury/signature_request_by_id/{id}
```

Example:

```bash
curl localhost:1317/zrchain/treasury/signature_request_by_id/1
```

Example Output:

```json
{
  "sign_request": {
    "id": "1",
    "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
    "key_id": "1",
    "key_type": "KEY_TYPE_ECDSA_SECP256K1",
    "data_for_signing": [
      "FYST..."
    ],
    "status": "SIGN_REQUEST_STATUS_FULFILLED",
    "signed_data": [
      {
        "sign_request_id": "1",
        "signed_data": "+L8bs..."
      }
    ],
    "keyring_party_signatures": [
      "C3JV..."
    ],
    "reject_reason": "",
    "metadata": {
      "@type": "/zrchain.treasury.MetadataEthereum",
      "chain_id": "11155111"
    },
    "parent_req_id": "0",
    "child_req_ids": [],
    "cache_id": null
  }
}
```

#### sign_transaction_request

The `sign_transaction_request` endpoint allows users to query signature requests for a specific wallet, key and status.

```bash
/zrchain/treasury/sign_transaction_request/{wallet_type}/{key_id}/{status}
```

Example:

```bash
curl localhost:1317/zrchain/treasury/sign_transaction_request/WALLET_TYPE_UNSPECIFIED/1/
```

Example Output:

```json
{
  "sign_transaction_requests": [
    {
      "sign_transaction_requests": {
        "id": "1",
        "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
        "key_id": "1",
        "wallet_type": "WALLET_TYPE_EVM",
        "unsigned_transaction": "+MqAh...",
        "sign_request_id": "6",
        "no_broadcast": false
      },
      "sign_request": {
        "id": "6",
        "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
        "key_id": "1",
        "key_type": "KEY_TYPE_ECDSA_SECP256K1",
        "data_for_signing": [
          "FYST..."
        ],
        "status": "SIGN_REQUEST_STATUS_FULFILLED",
        "signed_data": [
          {
            "sign_request_id": "6",
            "signed_data": "+L8bs..."
          }
        ],
        "keyring_party_signatures": [
          "jprc..."
        ],
        "reject_reason": "",
        "metadata": {
          "@type": "/zrchain.treasury.MetadataEthereum",
          "chain_id": "11155111"
        },
        "parent_req_id": "0",
        "child_req_ids": [],
        "cache_id": null
      }
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "1"
  }
}
```

#### sign_transaction_request_by_id

The `sign_transaction_request_by_id` endpoint allows users to query a specific sign transaction request by id.

```bash
/zrchain/treasury/sign_transaction_request_by_id/{id}
```

Example:

```bash
curl localhost:1317/zrchain/treasury/sign_transaction_request_by_id/1
```

Example Output:

```json
{
  "sign_transaction_request": {
    "id": "1",
    "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
    "key_id": "1",
    "wallet_type": "WALLET_TYPE_EVM",
    "unsigned_transaction": "+MqAh...",
    "sign_request_id": "2",
    "no_broadcast": false
  }
}
```
