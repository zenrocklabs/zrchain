# `x/policy`

## Abstract

The following documents specify the policy module.

Zenrock's Policy Module handles permissions and conditions around managing workspaces, keyrings and assets on Zenrock.

## Contents

* [Concepts](#concepts)
    * [Policies](#policies)
    * [Passkeys](#passkeys)
* [State](#state)
    * [Action Store](#action-store)
    * [Policy Store](#policy-store)
    * [SignMethod Store](#signmethod-store)
* [Msg Service](#msg-service)
    * [Msg/NewPolicy](#msgnewpolicy)
    * [Msg/ApproveAction](#msgapproveaction)
    * [Msg/RevokeAction](#msgrevokeaction)
    * [Msg/AddMultiGrant](#msgaddmultigrant)
    * [Msg/RemoveMultiGrant](#msgremovemultigrant)
    * [Msg/AddSignMethod](#msgaddsignmethod)
    * [Msg/RemoveSignMethod](#msgremovesignmethod)
* [Events](#events)
    * [EventNewAction](#eventnewaction)
* [Client](#client)
    * [CLI](#cli)
    * [gRPC](#grpc)
    * [REST](#rest)

## Concepts

### Policies

Policies are integrated in some messages of the identity and treasury modules. Such messages are only valid if the conditions specified by the policy - for example, multiple approvals - are fulfilled.

When approvals from other owners or admins are required, an action to approve the request is issued to the respective accounts.

Actions are valid for a defined number of blocks - if the terms of the policy are not fulfilled by the specified block height, the message will be rejected.

This ensures that relevant accounts are notified of pending actions and prevents requests from getting stuck in a perpetually pending state.

### Passkeys

Passkeys enable keys derived from biometrics to be used as workspace approvers. We will add them to the workspace so that users could approve an action for a trusted device they own, adding another factor to their approval process.

When passkeys are implemented, approvals can only be sourced from a certain device or classic Web2 account. This will add a second layer of security to the workspace. Passkeys may be used to ease access and improve the usability of policies.

## State

The `policy` module uses the `collections` package which provides collection storage.

Here's the list of collections stored as part of the `policy` module.

### Action Store

The `ActionStore` stores `Action`: `BigEndian(ActionId) -> ProtocolBuffer(Action)`.

### Policy Store

The `PolicyStore` stores `Policy`: `BigEndian(PolicyId) -> ProtocolBuffer(Policy)`.

### SignMethod Store

The `SignMethodStore` stores `SignMethod`: `Pair(String(owner),String(id)) -> ProtocolBuffer(SignMethod)`.

## Msg Service

### Msg/NewPolicy

A new policy can be added using the `MsgNewPolicy` message.
The policy is specified as a BoolParsePolicy object.

```proto
message MsgNewPolicy {
  option (cosmos.msg.v1.signer) = "creator";
  string              creator = 1;
  string              name    = 2;
  google.protobuf.Any policy  = 3;
  uint64              btl     = 4; 
}

message BoolparserPolicy {
  // Definition of the policy, eg.
  // "t1 + t2 + t3 > 1"
  string definition = 1;
  repeated PolicyParticipant participants = 2;
}

message PolicyParticipant {
  string abbreviation = 1 [deprecated = true];
  string address = 2;
}
```

It's expected to fail if

* the specified policy definition is not valid 
* a participant has no address specified 
* a participant is not used in the expression
* a participant address is used more than once
* a participant address is not a valid Bech32 address

### Msg/ApproveAction

An action can be approved with the `MsgApproveAction` message.
The additional signature can currently only be of type `AdditionalSignaturePasskey`

```proto
message MsgApproveAction {
  option (cosmos.msg.v1.signer) = "creator";
           string              creator               = 1;
           string              action_type           = 2;
           uint64              action_id             = 3;
  repeated google.protobuf.Any additional_signatures = 4;
}

message AdditionalSignaturePasskey {
  bytes raw_id = 1;
  bytes authenticator_data = 2;
  bytes client_data_json = 3;
  bytes signature = 4;
}
```

It's expected to fail if

* the action is not found
* the action is no longer has pending state
* the action btl has been reached

### Msg/RevokeAction

An action can be revoked with the `MsgRevokeAction` message.

```proto
message MsgRevokeAction {
  option (cosmos.msg.v1.signer) = "creator";
  string creator   = 1;
  uint64 action_id = 2;
}
```

It's expected to fail if

* the action is not found
* the action is no longer has pending state
* the transaction creator is not the creator of the action

### Msg/AddMultiGrant

Multiple grants can be granted with the `MsgAddMultiGrant` message.

```proto
message MsgAddMultiGrant {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  string grantee = 2;
  repeated string msgs = 3;
}
```

It's expected to fail if

* the grantee is not a valid address
* a specified msg is not a valid message identifier

### Msg/RemoveMultiGrant

Multiple grants can be removed with the `MsgRemoveMultiGrant` message.

```proto
message MsgRemoveMultiGrant {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  string grantee = 2;
  repeated string msgs = 3;
}
```

It's expected to fail if

* the grantee is not a valid address
* a specified msg is not a valid message identifier


### Msg/AddSignMethod

An additional signmethod can be added with the `MsgAddSignMethod` message.

```proto
message MsgAddSignMethod {
  option (cosmos.msg.v1.signer) = "creator";
  string              creator = 1;
  google.protobuf.Any config  = 2;
}

message SignMethodPasskey {
  bytes raw_id = 1;
  bytes attestation_object = 2;
  bytes client_data_json = 3;
  bool active = 4;
}
```

It's expected to fail if

* the config is not of type SignMethodPasskey
* the config is not doenst contain the correct challenge or signature (see GUIDE.md)

### Msg/RemoveSignMethod

An additional signmethod can be removed with the `MsgRemoveSignMethod` message.

```proto
message MsgRemoveSignMethod {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  string id      = 2;
}
```

It's expected to fail if

* the signmethod is not found for the creator of the transaction

## Parameters

### DefaultBTL

The default amount of blocks that an action will be active for users to approve it. 
If no BTL (blocks to live) value is specified in either policy or cli this value will used.

### MinimumBTL

The minumum amount of blocks an action will be active.
If the BTL value specified in either policy or cli is smaller than this value, then this value is used.

## Events

The policy module emits the following events:

### EventNewAction

| Type                             | Attribute Key    | Attribute Value     |
| -------------------------------- | -------------    | ------------------- |
| new_action                       | action_id        | {action_id}         |
| new_action                       | participant_addr | {participant_addr}  |

## Client

### CLI

A user can query and interact with the `policy` module using the CLI.

#### Query

The `query` commands allow users to query `policy` state.

```bash
zenrockd query policy --help
```

##### actions

The `actions` command allows users to query all actions.

```bash
zenrockd query policy actions
```

Example:

```bash
zenrockd query policy actions
```

Example Output:

```yaml
- approvers:
  - zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  btl: "7744"
  creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  id: "1"
  msg:
    type: /zrchain.treasury.MsgNewKeyRequest
    value:
      btl: "1000"
      creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
      key_type: ecdsa
      keyring_addr: keyring1pfnq7r04rept47gaf5cpdew2
      workspace_addr: workspace1mphgzyhncnzyggfxmv4nmh
  policy_id: "1"
  status: ACTION_STATUS_PENDING
pagination:
  total: "1"
```

##### action-details-by-id

The `action-details-by-id` command allows users to query the details of an action by id.

```bash
zenrockd query policy action-details-by-id [id]
```

Example:

```bash
zenrockd query policy action-details-by-id 1
```

Example Output:

```yaml
action:
  approvers:
  - zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  btl: "9199"
  creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  id: "1"
  msg:
    type: /zrchain.treasury.MsgNewKeyRequest
    value:
      btl: "1000"
      creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
      key_type: ecdsa
      keyring_addr: keyring1pfnq7r04rept47gaf5cpdew2
      workspace_addr: workspace1mphgzyhncnzyggfxmv4nmh
  policy_id: "1"
  status: ACTION_STATUS_PENDING
approvers:
- zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
current_height: "9000"
id: "1"
pending_approvers:
- zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq
policy:
  btl: "1000"
  creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  id: "1"
  name: alice_faucet
  policy:
    type: /zrchain.policy.BoolparserPolicy
    value:
      definition: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq > 1
      participants:
      - address: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
      - address: zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq
```

##### policies

The `policies` command allows users to query all policies.

```bash
zenrockd query policy policies
```

Example:

```bash
zenrockd query policy policies
```

Example Output:

```yaml
pagination:
  total: "1"
policies:
- policy:
    btl: "1000"
    creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
    id: "1"
    name: policy2
    policy:
      type: /zrchain.policy.BoolparserPolicy
      value:
        definition: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq > 1
        participants:
        - address: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
        - address: zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq
```

##### policies-by-creator

The `policies-by-creator` command allows users to query polcies that were create by one or more users.

```bash
zenrockd query policy policies-by-creator [creators]
```

Example:

```bash
zenrockd query policy policies-by-creator zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty,zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq
```

Example Output:

```yaml
pagination:
  total: "2"
policies:
- creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  id: "1"
  name: policy1
  policy:
    type: /zrchain.policy.BoolparserPolicy
    value:
      definition: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq > 1
      participants:
      - address: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
      - address: zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq
- creator: zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq
  id: "2"
  name: policy2
  policy:
    type: /zrchain.policy.BoolparserPolicy
    value:
      definition: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq > 1
      participants:
      - address: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
      - address: zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq
```

##### policy-by-id

The `policy-by-id` command allows users to query a policy by id.

```bash
zenrockd query policy policy-by-id [id]
```

Example:

```bash
zenrockd query policy policy-by-id 1
```

Example Output:

```yaml
policy:
  policy:
    btl: "1000"
    creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
    id: "1"
    name: policy1
    policy:
      type: /zrchain.policy.BoolparserPolicy
      value:
        definition: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq > 1
        participants:
        - address: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
        - address: zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq
```

##### sign-methods-by-address

The `sign-methods-by-address` command allows users to query the signmethods for a user.

```bash
zenrockd query policy sign-methods-by-address [address]
```

Example:

```bash
zenrockd query policy sign-methods-by-address zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
```

Example Output:

```yaml
config:
- type: /zrchain.policy.SignMethodPasskey
  value:
    active: true
    attestation_object: o2Nm...
    client_data_json: eyJ0e...
    raw_id: KhlOB...
pagination:
  total: "1"
```

### Transactions

The `tx` commands allow users to interact with the `policy` module.

```bash
zenrockd tx policy --help
```

#### new-policy

The `new-policy` command allows users to create a new policy.

```bash
zenrockd tx policy new-policy [name] [policy] --btl [btl]
```

Example:

```bash
zenrockd tx policy new-policy policy1 '{"@type":"/zrchain.policy.BoolparserPolicy", "definition": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq > 1", "participants":[{ "address":"zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty" },{  "address":"zen17ekx844yl3ftmcl47ryc7fz5cd7uhxq4f5ma5q" }]}'
```

#### approve-action

The `approve-action` command allows users to approve an action.

```bash
zenrockd tx policy approve-action [action-type] [action-id]
```

Example:

```bash
zenrockd tx policy approve-action "/zrchain.treasury.MsgNewKeyRequest" 1
```

#### revoke-action

The `revoke-action` command allows users to revoke an action.

```bash
zenrockd tx policy revoke-action [action-id]
```

Example:

```bash
zenrockd tx policy revoke-action 1
```

#### add-multi-grant

The `add-multi-grant` command allows users to add multiple grants for the specified grantee.

```bash
zenrockd tx policy add-multi-grant [grantee] [msgs] --from [granter]
```

Example:

```bash
zenrockd tx policy add-multi-grant zen17ekx844yl3ftmcl47ryc7fz5cd7uhxq4f5ma5q /zrchain.policy.MsgApproveAction,/zrchain.policy.MsgRevokeAction
```

#### remove-multi-grant

The `remove-multi-grant` command allows users to remove multiple grants for the specified grantee.

```bash
zenrockd tx policy remove-multi-grant [grantee] [msgs] --from [granter]
```

Example:

```bash
zenrockd tx policy remove-multi-grant zen17ekx844yl3ftmcl47ryc7fz5cd7uhxq4f5ma5q /zrchain.policy.MsgApproveAction,/zrchain.policy.MsgRevokeAction
```

#### add-sign-method

The `add-sign-method` command allows users to add an additional sign method.

```bash
zenrockd tx policy add-sign-method [id] [config]
```

Example:

```bash
export json_payload=$(cat <<EOF
{
    "@type": "/zrchain.policy.SignMethodPasskey",
    "raw_id": "KhlO...",
    "client_data_json": "eyJ0...",
    "attestation_object": "o2Nm..."
}
EOF
)

zenrockd tx policy add-sign-method "$json_payload" --from alice 
```

#### remove-sign-method

The `remove-sign-method` command allows users to to remove a signmethod.

```bash
zenrockd tx policy remove-sign-method [id] --from [owner]
```

Example:

```bash
zenrockd tx policy remove-sign-method 1
```

### gRPC

A user can query the `policy` module using gRPC endpoints.

#### Actions

The `Actions` endpoint allows users to query for all actions.

```bash
zrchain.policy.Query/Actions
```

Example:

```bash
grpcurl -plaintext localhost:9090 zrchain.policy.Query/Actions
```

Example Output:

```json
{
  "pagination": {
    "total": "1"
  },
  "actions": [
    {
      "id": "1",
      "approvers": [
        "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
      ],
      "status": "ACTION_STATUS_PENDING",
      "policyId": "1",
      "msg": {
        "@type": "/zrchain.treasury.MsgNewKeyRequest",
        "btl": "1000",
        "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
        "keyType": "ecdsa",
        "keyringAddr": "keyring1pfnq7r04rept47gaf5cpdew2",
        "workspaceAddr": "workspace1mphgzyhncnzyggfxmv4nmh"
      },
      "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
      "btl": "9199"
    }
  ]
}
```

#### Policies

The `Policies` endpoint allows users to query for all policies.

```bash
zrchain.policy.Query/Policies
```

Example:

```bash
grpcurl -plaintext localhost:9090 zrchain.policy.Query/Policies
```

Example Output:

```json
{
  "pagination": {
    "total": "2"
  },
  "policies": [
    {
      "policy": {
        "creator": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
        "id": "1",
        "name": "policy1",
        "policy": {
          "@type": "/zrchain.policy.BoolparserPolicy",
          "definition": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq \u003e 1",
          "participants": [
            {
              "address": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
            },
            {
              "address": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq"
            }
          ]
        },
        "btl": "1000"
      }
    },
    {
      "policy": {
        "creator": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
        "id": "2",
        "name": "policy2",
        "policy": {
          "@type": "/zrchain.policy.BoolparserPolicy",
          "definition": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq \u003e 1",
          "participants": [
            {
              "address": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
            },
            {
              "address": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq"
            }
          ]
        },
        "btl": "1000"
      }
    }
  ]
}
```

#### PolicyById

The `PolicyById` endpoint allows users to query for a policy by id.

```bash
zrchain.policy.Query/PolicyById
```

Example:

```bash
grpcurl -plaintext \
  -d '{"id": 1}' localhost:9090 zrchain.policy.Query/PolicyById
```

Example Output:

```json
{
  "policy": {
    "policy": {
      "creator": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
      "id": "1",
      "name": "policy1",
      "policy": {
        "@type": "/zrchain.policy.BoolparserPolicy",
        "definition": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq \u003e 1",
        "participants": [
          {
            "address": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
          },
          {
            "address": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq"
          }
        ]
      },
      "btl": "1000"
    }
  }
}
```

#### PoliciesByCreator

The `PoliciesByCreator` endpoint allows users to query for policies created by one or more specified users.

```bash
zrchain.policy.Query/PoliciesByCreator
```

Example:

```bash
grpcurl -plaintext \
  -d '{"creators": ["zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq","zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq"]}' \
  localhost:9090 zrchain.policy.Query/PoliciesByCreator
```

Example Output:

```json
{
  "policies": [
    {
      "creator": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
      "id": "1",
      "name": "policy1",
      "policy": {
        "@type": "/zrchain.policy.BoolparserPolicy",
        "definition": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq \u003e 1",
        "participants": [
          {
            "address": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
          },
          {
            "address": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq"
          }
        ]
      }
    },
    {
      "creator": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
      "id": "2",
      "name": "policy2",
      "policy": {
        "@type": "/zrchain.policy.BoolparserPolicy",
        "definition": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq \u003e 1",
        "participants": [
          {
            "address": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
          },
          {
            "address": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq"
          }
        ]
      }
    }
  ],
  "pagination": {
    "total": "2"
  }
}
```

#### ActionDetailsById

The `ActionDetailsById` endpoint allows users to query for action details for a specific action id.

```bash
zrchain.policy.Query/ActionDetailsById
```

Example:

```bash
grpcurl -plaintext \
  -d '{"id":1}' localhost:9090 zrchain.policy.Query/ActionDetailsById
```

Example Output:

```json
{
  "id": "1",
  "action": {
    "id": "1",
    "approvers": [
      "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
    ],
    "status": "ACTION_STATUS_PENDING",
    "policyId": "1",
    "msg": {
      "@type": "/zrchain.treasury.MsgNewKeyRequest",
      "btl": "1000",
      "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
      "keyType": "ecdsa",
      "keyringAddr": "keyring1pfnq7r04rept47gaf5cpdew2",
      "workspaceAddr": "workspace1mphgzyhncnzyggfxmv4nmh"
    },
    "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
    "btl": "9199"
  },
  "policy": {
    "creator": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
    "id": "1",
    "name": "alice_faucet",
    "policy": {
      "@type": "/zrchain.policy.BoolparserPolicy",
      "definition": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq \u003e 1",
      "participants": [
        {
          "address": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
        },
        {
          "address": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq"
        }
      ]
    },
    "btl": "1000"
  },
  "approvers": [
    "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
  ],
  "pendingApprovers": [
    "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq"
  ],
  "currentHeight": "9000"
}
```

#### SignMethodsByAddress

The `SignMethodsByAddress` endpoint allows users to query for signmethods for a specific user.

```bash
zrchain.policy.Query/SignMethodsByAddress
```

Example:

```bash
grpcurl -plaintext \
  -d '{"address": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"}' localhost:9090 zrchain.policy.Query/SignMethodsByAddress
```

Example Output:

```json
{
  "pagination": {
    "total": "1"
  },
  "config": [
    {
      "@type": "/zrchain.policy.SignMethodPasskey",
      "active": true,
      "attestationObject": "o2Nm...",
      "clientDataJson": "eyJ0...",
      "rawId": "KhlO..."
    }
  ]
}
```

### REST

A user can query the `policy` module using REST endpoints.

#### actions

The `actions` endpoint allows users to query .....

```bash
/zrchain/policy/actions
```

Example:

```bash
curl localhost:1317/zrchain/policy/actions
```

Example Output:

```json
{
  "pagination": {
    "next_key": null,
    "total": "1"
  },
  "actions": [
    {
      "id": "1",
      "approvers": [
        "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
      ],
      "status": "ACTION_STATUS_PENDING",
      "policy_id": "1",
      "msg": {
        "@type": "/zrchain.treasury.MsgNewKeyRequest",
        "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
        "workspace_addr": "workspace1mphgzyhncnzyggfxmv4nmh",
        "keyring_addr": "keyring1pfnq7r04rept47gaf5cpdew2",
        "key_type": "ecdsa",
        "btl": "1000",
        "index": "0",
        "ext_requester": "",
        "ext_key_type": "0",
        "sign_policy_id": "0"
      },
      "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
      "btl": "9199",
      "policy_data": []
    }
  ]
}
```

#### policies

The `policies` endpoint allows users to query .....

```bash
/zrchain/policy/policies
```

Example:

```bash
curl localhost:1317/zrchain/policy/policies
```

Example Output:

```json
{
  "pagination": {
    "next_key": null,
    "total": "1"
  },
  "policies": [
    {
      "policy": {
        "creator": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
        "id": "1",
        "name": "policy1",
        "policy": {
          "@type": "/zrchain.policy.BoolparserPolicy",
          "definition": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq > 1",
          "participants": [
            {
              "address": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
            },
            {
              "address": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq"
            }
          ]
        },
        "btl": "1000"
      },
      "metadata": null
    },
  ]
}
```

#### policy_by_id

The `policy_by_id` endpoint allows users to query .....

```bash
/zrchain/policy/policy_by_id/{id}
```

Example:

```bash
curl localhost:1317/zrchain/policy/policy_by_id/1
```

Example Output:

```json
{
  "policy": {
    "policy": {
      "creator": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
      "id": "1",
      "name": "alice_faucet",
      "policy": {
        "@type": "/zrchain.policy.BoolparserPolicy",
        "definition": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq > 1",
        "participants": [
          {
            "address": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
          },
          {
            "address": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq"
          }
        ]
      },
      "btl": "1000"
    },
    "metadata": null
  }
}
```

#### policy_by_creator

The `policies_by_creator` endpoint allows users to query .....

```bash
/zrchain/policy/policies_by_creator/{creators}
```

Example:

```bash
curl localhost:1317/zrchain/policy/policies_by_creator/zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq,zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq
```

Example Output:

```json
{
  "policies": [
    {
      "creator": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
      "id": "1",
      "name": "policy1",
      "policy": {
        "@type": "/zrchain.policy.BoolparserPolicy",
        "definition": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq > 1",
        "participants": [
          {
            "address": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
          },
          {
            "address": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq"
          }
        ]
      },
      "btl": "0"
    },
    {
      "creator": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
      "id": "2",
      "name": "policy2",
      "policy": {
        "@type": "/zrchain.policy.BoolparserPolicy",
        "definition": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq > 1",
        "participants": [
          {
            "address": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
          },
          {
            "address": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq"
          }
        ]
      },
      "btl": "0"
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "2"
  }
}
```

#### action_details_by_id

The `action_details_by_id` endpoint allows users to query .....

```bash
/zrchain/policy/action_details_by_id/{id}
```

Example:

```bash
curl localhost:1317/zrchain/policy/action_details_by_id/1
```

Example Output:

```json
{
  "id": "1",
  "action": {
    "id": "1",
    "approvers": [
      "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
    ],
    "status": "ACTION_STATUS_PENDING",
    "policy_id": "1",
    "msg": {
      "@type": "/zrchain.treasury.MsgNewKeyRequest",
      "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
      "workspace_addr": "workspace1mphgzyhncnzyggfxmv4nmh",
      "keyring_addr": "keyring1pfnq7r04rept47gaf5cpdew2",
      "key_type": "ecdsa",
      "btl": "1000",
      "index": "0",
      "ext_requester": "",
      "ext_key_type": "0",
      "sign_policy_id": "0"
    },
    "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
    "btl": "9199",
    "policy_data": []
  },
  "policy": {
    "creator": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
    "id": "1",
    "name": "policy1",
    "policy": {
      "@type": "/zrchain.policy.BoolparserPolicy",
      "definition": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq > 1",
      "participants": [
        {
          "address": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
        },
        {
          "address": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq"
        }
      ]
    },
    "btl": "1000"
  },
  "approvers": [
    "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
  ],
  "pending_approvers": [
    "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq"
  ],
  "current_height": "9000"
}
```

#### sign_methods_by_address

The `sign_methods_by_address` endpoint allows users to query .....

```bash
/zrchain/policy/sign_methods_by_address/{address}
```

Example:

```bash
curl localhost:1317/zrchain/policy/sign_methods_by_address/zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
```

Example Output:

```json
{
  "pagination": {
    "next_key": null,
    "total": "1"
  },
  "config": [
    {
      "@type": "/zrchain.policy.SignMethodPasskey",
      "raw_id": "KhlO...",
      "attestation_object": "o2Nmb...",
      "client_data_json": "eyJ0...",
      "active": true
    }
  ]
}
```
