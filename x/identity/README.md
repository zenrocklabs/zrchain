# `x/identity`

## Abstract

The following documents specify the identity module.

Zenrock's Identity Module is responsible for identity management on Zenrock. There are currently two objects that are stored and managed using this module: Keyrings & Workspaces.

## Contents

* [Concepts](#concepts)
    * [Keyrings](#keyrings)
    * [Workspaces](#workspaces)
* [State](#state)
    * [Keyring Store](#keyring-store)
    * [Workspace Store](#workspace-store)
* [Msg Service](#msg-service)
    * [Msg/NewWorkspace](#msgnewworkspace)
    * [Msg/AddWorkspaceOwner](#msgaddworkspaceowner)
    * [Msg/RemoveWorkspaceOwner](#msgremoveworkspaceowner)
    * [Msg/NewChildWorkspace](#msgnewchildworkspace)
    * [Msg/AppendChildWorkspace](#msgappendchildworkspace)
    * [Msg/UpdateWorkspace](#msgupdateworkspace)
    * [Msg/NewKeyring](#msgnewkeyring)
    * [Msg/AddKeyringParty](#msgaddkeyringparty)
    * [Msg/AddKeyringAdmin](#msgaddkeyringadmin)
    * [Msg/RemoveKeyringParty](#msgremovekeyringparty)
    * [Msg/RemoveKeyringAdmin](#msgremovekeyringadmin)
    * [Msg/DeactivateKeyring](#msgdeactivatekeyring)
    * [Msg/UpdateKeyring](#msgupdatekeyring)
* [Events](#events)
    * [EventNewWorkspace](#eventnewworkspace)
    * [EventAddOwnerToWorkspace](#eventaddownertoworkspace)
    * [EventOwnerAddedToWorkspace](#eventowneraddedtoworkspace)
* [Parameters](#parameters)
    * [KeyringCreationFee](#keyringcreationfee)
* [Client](#client)
    * [CLI](#cli)
    * [gRPC](#grpc)
    * [REST](#rest)

## Concepts

### Keyrings

Keyrings represent off-chain systems that provide keys and signatures. One example of this is the Zenrock MPC system.

Keyrings have a distinct address. This is used to select the keyring when performing key and signature requests.

Keyrings also define how many parties are involved. Only parties inside the keyring object are eligible to broadcast responses on-chain.

Lastly, a keyring provider can set the costs that occur for each key and signature request.

### Workspaces

Workspaces allow users (Workspace owners) to manage sets of keys. This allows key rotation and decouples the risk from managing wallets through one single account.

Apart from having associated keys, workspace settings also define policies for admin and signature tasks.

Policies define governance of the Workspace. They typically specify which combination of accounts must provide approval for a signature to be considered valid.

Policies also provide a governance structure for Workspace changes and can be used to prevent giving a single account full control over the Workspace.

...

## State

The `identity` module uses the `collections` package which provides collection storage.

Here's the list of collections stored as part of the `identity` module.

### Keyring store

The `KeyringStore` stores `Keyring`: `String(KeyringAddress) -> ProtocolBuffer(Keyring)`.

### Workspace store

The `WorkspaceStore` stores `Workspace`: `String(WorkspaceAddress) -> ProtocolBuffer(Workspace)`.

## Msg Service

### Msg/NewWorkspace

A new workspace can be created with the `MsgNewWorkspace`, which can specify the admin policy, sign policy and an optional list of additional owners to be added, additional to the creator of the tx.
By default you can specify 0 for admin or sign policy id in which case the creator will be the only member of a default policy.

```proto
message MsgNewWorkspace {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  uint64 admin_policy_id = 2;
  uint64 sign_policy_id = 3;
  repeated string additional_owners = 4; // Optional
}
```

It's expected to fail if

* if the the creator or additional owners are not a member of the specified admin or sign policy

### Msg/AddWorkspaceOwner

An additional member can be added with the `MsgAddWorkspaceOwner` message.

```go reference
message MsgAddWorkspaceOwner {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  string workspace_addr = 2;
  string new_owner = 3;
  uint64 btl = 4;
}
```

It's expected to fail if

* the workspace is not found
* the transaction creator is not an owner of the workspace
* the new owner is already an owner of the workspace

### Msg/RemoveWorkspaceOwner

An existing owner can be removed from a workspace with the `MsgRemoveWorkspaceOwner` message.

```proto
message MsgRemoveWorkspaceOwner {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  string workspace_addr = 2;
  string owner = 3;
  uint64 btl = 4;
}
```

It's expected to fail if

* the workspace is not found
* the transaction creator is not an owner of the workspace
* the existing owner is not an owner of the workspace
* the existing owner is a member of the assigned admin or sign policy

### Msg/NewChildWorkspace

A new child workspace can be created with the `MsgNewChildWorkspace` message.

```proto
message MsgNewChildWorkspace {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  string parent_workspace_addr = 2;
  uint64 btl = 3;
}
```

It's expected to fail if

* the parent workspace is not found
* the transaction creator is not an owner of the workspace

### Msg/AppendChildWorkspace

An existing workspace can be added to another existing workspace as a child workspace with the `MsgAppendChildWorkspace` message.

```proto
message MsgAppendChildWorkspace {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  string parent_workspace_addr = 2;
  string child_workspace_addr = 3;
  uint64 btl = 4;
}
```

It's expected to fail if

* the parent or child workspace is not found
* the transaction creator is not an owner of the parent or child workspace
* the child workspace is already a child of the parent workspace

### Msg/UpdateWorkspace

An existing workspace can be update with the `MsgUpdateWorkspace` message.

```proto
message MsgUpdateWorkspace {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  string workspace_addr = 2;
  uint64 admin_policy_id = 3;
  uint64 sign_policy_id = 4;
  uint64 btl = 5;
}
```

It's expected to fail if

* the workspace is not found
* the transaction creator is not an owner of the workspace
* there are no updates to the policies

### Msg/NewKeyring

A new keyring can be created with the `MsgNewKeyring` message.

```proto
message MsgNewKeyring {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  string description = 2;
  uint32 party_threshold = 3;
  uint64 key_req_fee = 4;
  uint64 sig_req_fee = 5;
}
```

It's expected to fail if

* the balance of the transaction sender is to low to pay the fee

### Msg/AddKeyringParty

A new party can be added to an existing keyring with the `MsgAddKeyringParty` message.

```proto
message MsgAddKeyringParty {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  string keyring_addr = 2;
  string party = 3;
  bool increase_threshold = 4; // Optional flag
}
```

It's expected to fail if

* the keyring is not found
* the keyring is not active
* the party is already added to the keyring
* the transaction creator is not an admin of the keyring

### Msg/AddKeyringAdmin

A new admin can be added to an existing keyring with the `MsgAddKeyringAdmin` message.

```proto
message MsgAddKeyringAdmin {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  string keyring_addr = 2;
  string admin = 3;
}
```

It's expected to fail if

* the keyring is not found
* the keyring is not active
* the admin is already an admin of the keyring
* the transaction creator is not an admin of the keyring

### Msg/RemoveKeyringParty

An party can be removed from a keyring with the `MsgRemoveKeyringParty` message.

```proto
message MsgRemoveKeyringParty {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  string keyring_addr = 2;
  string party = 3;
  bool decrease_threshold = 4; // Optional flag
}
```

It's expected to fail if

* the keyring is not found
* the keyring is not active
* the party is not a party of the keyring
* the transaction creator is not an admin of the keyring

### Msg/RemoveKeyringAdmin

An admin can be removed from a keyring with the `MsgRemoveKeyringAdmin` message.

```proto
message MsgRemoveKeyringAdmin {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  string keyring_addr = 2;
  string admin = 3;
}
```

It's expected to fail if

* the keyring is not found
* the keyring is not active
* the admin is not an admin of the keyring
* the transaction creator is not an admin of the keyring
* the admin to remove is the last admin in the keyring

### Msg/DeactivateKeyring

A keyring can be deactivated using the `MsgDeactivateKeyring` message.

```proto
message MsgDeactivateKeyring {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  string keyring_addr = 2;
}
```

It's expected to fail if

* the keyring is not found
* the transaction creator is not a keyring admin

### Msg/UpdateKeyring

A keyring can be updated with the `MsgUpdateKeyring` message.

```proto
message MsgUpdateKeyring {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  string keyring_addr = 2;
  uint32 party_threshold = 3;
  uint64 key_req_fee = 4;
  uint64 sig_req_fee = 5;
  string description = 6;
  bool is_active = 7;
}
```

It's expected to fail if

* keyring is not found
* the transaction creator is not a keyring admin

## Events

The identity module emits the following events:

### EventNewWorkspace

| Type                             | Attribute Key  | Attribute Value                   |
| -------------------------------- | -------------  | --------------------------------  |
| message                          | action         | /zrchain.identity.MsgNewWorkspace |
| message                          | module         | identity                          |
| new_workspace                    | workspace_addr | {workapce_addr}                   |


### EventAddOwnerToWorkspace

| Type                             | Attribute Key  | Attribute Value                        |
| -------------------------------- | -------------  | --------------------------------       |
| message                          | action         | /zrchain.identity.MsgAddWorkspaceOwner |
| message                          | module         | identity                               |
| add_owner_to_workspace           | action_id      | {action_id}                            |
 

### EventOwnerAddedToWorkspace

| Type                             | Attribute Key  | Attribute Value                        |
| -------------------------------- | -------------  | --------------------------------       |
| message                          | action         | /zrchain.identity.MsgAddWorkspaceOwner |
| message                          | module         | identity                               |
| owner_added_to_workspace         | workspace_addr | {workspace_addr}                       |
| owner_added_to_workspace         | owner_addr     | {owner_addr}                           |

## Parameters

### KeyringCreationFee

The KeyringCreationFee specifies the amount in urock that needs to be paid by the transaction sender to create a keyring.

## Client

### CLI

A user can query and interact with the `identity` module using the CLI.

#### Query

The `query` commands allow users to query `identity` state.

```bash
zenrockd query identity --help
```

##### keyrings

The `keyrings` command allows users to query the available keyrings.

```bash
zenrockd query identity keyrings 
```

Example:

```bash
zenrockd query identity keyrings
```

Example Output:

```bash
keyrings:
- address: keyring1k6vc6vhp6e6l3rxalue9v4ux
  admins:
  - zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  description: Keyring with Fees
  is_active: true
  key_req_fee: "2"
  parties:
  - zen10kmgv5gzygnecf46x092ecfe5xcvvv9rdaxmts
  sig_req_fee: "2"
- address: keyring1pfnq7r04rept47gaf5cpdew2
  admins:
  - zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  description: ZenrockKMS
  is_active: true
  parties:
  - zen10kmgv5gzygnecf46x092ecfe5xcvvv9rdaxmts
pagination:
  total: "2"
```

##### keyring-by-address

The `keyring-by-address` command allows users to query ...todo...

```bash
zenrockd query identity keyring-by-address [keyring-addr]
```

Example:

```bash
zenrockd query identity keyring-by-address keyring1k6vc6vhp6e6l3rxalue9v4ux
```

Example Output:

```bash
keyring:
  address: keyring1k6vc6vhp6e6l3rxalue9v4ux
  admins:
  - zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  description: Keyring with Fees
  is_active: true
  key_req_fee: "2"
  parties:
  - zen10kmgv5gzygnecf46x092ecfe5xcvvv9rdaxmts
  sig_req_fee: "2"
```

##### workspaces

The `workspaces` command allows users to query the available workspaces.

```bash
zenrockd query identity workspaces
```

Example:

```bash
zenrockd query identity workspaces
```

Example Output:

```bash
pagination:
  total: "1"
workspaces:
- address: workspace14a2hpadpsy9h4auve2z8lw
  creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  owners:
  - zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
```

##### workspace-by-address

The `workspace-by-address` command allows users to query a workspace by address.

```bash
zenrockd query identity workspace-by-address [workspace-addr]
```

Example:

```bash
zenrockd query identity workspace-by-address workspace14a2hpadpsy9h4auve2z8lw
```

Example Output:

```bash
workspace:
  address: workspace14a2hpadpsy9h4auve2z8lw
  creator: zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
  owners:
  - zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty
```

### Transactions

The `tx` commands allow users to interact with the `identity` module.
In the examples below, the btl flag is used to specify the amount of blocks within which the users of the admin policy have to approve the action that is created (see policy module), if not specified the btl from the policy or a default of 1000 will be used.

```bash
zenrockd tx identity --help
```

#### new-workspace

The `new-workspace` command allows users to create a new workspace.

```bash
zenrockd tx tx identity new-workspace --admin-policy-id [admin-policy-id] --sign-policy-id [sign-policy-id] --additional-owners [additional-owners]
```

Example:

```bash
zenrockd tx identity new-workspace --admin-policy-id 1 --sign-policy-id 1 --additional-owners zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq
```

#### add-workspace-owner

The `add-workspace-owner` command allows users to add owners to a workspace.

```bash
zenrockd tx identity add-workspace-owner [workspace-addr] [owner-address] --btl [btl]
```

Example:

```bash
zenrockd tx identity add-workspace-owner workspace14a2hpadpsy9h4auve2z8lw zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq 
```

#### remove-workspace-owner

The `remove-workspace-owner` command allows users to .....

```bash
zenrockd tx identity remove-workspace-owner [workspace-address] [owner-address] --btl [btl]
```

Example:

```bash
zenrockd tx identity remove-workspace-owner workspace14a2hpadpsy9h4auve2z8lw zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq 
```

#### new-child-workspace

The `new-child-workspace` command allows users to create and add a new child workspace to an existing workspace.

```bash
zenrockd tx identity new-child-workspace [parent-workspace-addr] --btl [btl]
```

Example:

```bash
zenrockd tx identity new-child-workspace workspace14a2hpadpsy9h4auve2z8lw
```

#### append-child-workspace

The `append-child-workspace` command allows users to add an existing workspace as a child to another workspace.

```bash
zenrockd tx identity append-child-workspace [parent-workspace-addr] [child-workspace-addr] --btl [btl]
```

Example:

```bash
zenrockd tx identity append-child-workspace workspace1mphgzyhncnzyggfxmv4nmh workspace1xklrytgff7w32j52v34w36
```

#### update-workspace

The `update-workspace` command allows users to update the policies on a workspace.

```bash
zenrockd tx identity update-workspace [workspace-address] [admin-policy-id] [sign-policy-id] --btl [btl]
```

Example:

```bash
zenrockd tx identity update-workspace workspace14a2hpadpsy9h4auve2z8lw 1 1
```

#### new-keyring

The `new-keyring` command allows users to create a new keyring.
The fees are specified as a uint64 in urock.

```bash
zenrockd tx identity new-keyring [description] [key-request-fee] [sign-request-fee]
```

Example:

```bash
zenrockd tx identity new-keyring keyring1 10000 10000
```

#### add-keyring-admin

The `add-keyring-admin` command allows users to add an admin to a keyring.

```bash
zenrockd tx identity add-keyring-admin [keyring-addr] [admin]
```

Example:

```bash
zenrockd tx identity add-keyring-admin keyring1k6vc6vhp6e6l3rxalue9v4ux zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq
```

#### add-keyring-party

The `add-keyring-party` command allows users to add a signing party to a keyring.
The increase threshold parameter specifies whether or not to increase the signature threshold for the keyring.

```bash
zenrockd tx identity add-keyring-party [keyring-addr] [party] --increase-threshold [true]
```

Example:

```bash
zenrockd tx identity add-keyring-party add-keyring-party keyring1k6vc6vhp6e6l3rxalue9v4ux zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq --increase-threshold true
```

#### remove-keyring-admin

The `remove-keyring-admin` command allows users to remove an admin from a keyring.

```bash
zenrockd tx identity remove-keyring-admin [keyring-addr] [admin]
```

Example:

```bash
zenrockd tx identity remove-keyring-admin keyring1k6vc6vhp6e6l3rxalue9v4ux zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq
```

#### remove-keyring-party

The `remove-keyring-party` command allows users to remove a signing party from a keyring.
The decrease threshold parameter specifies whether or not to decrease the signature threshold for the keyring.


```bash
zenrockd tx identity remove-keyring-party [keyring-addr] [party] --decrease-threshold [true] 
```

Example:

```bash
zenrockd tx identity remove-keyring-party keyring1k6vc6vhp6e6l3rxalue9v4ux zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq --decrease-threshold true
```

#### deactivate-keyring

The `deactivate-keyring` command allows users to deactivate a keyring, this will disable the keyring from creating new keys or signing new requests.

```bash
zenrockd tx identity deactivate-keyring [keyring-addr]
```

Example:

```bash
zenrockd tx identity deactivate-keyring keyring1k6vc6vhp6e6l3rxalue9v4ux
```

#### update-keyring

The `update-keyring` command allows users to .....

```bash
zenrockd tx identity update-keyring [keyring-addr] [is-active:true|false] [party-threshold] [key-req-fee] [sig-req-fee] [description]
```

Example:

```bash
zenrockd tx identity update-keyring keyring1k6vc6vhp6e6l3rxalue9v4ux true 1 10000 10000 keyring1
```

### gRPC

A user can query the `identity` module using gRPC endpoints.

#### Workspaces

The `Workspaces` endpoint allows users to query all workspaces.

```bash
zrchain.identity.Query/Workspaces
```

Example:

```bash
grpcurl -plaintext localhost:9090 zrchain.identity.Query/Workspaces
```

Example Output:

```json
{
  "workspaces": [
    {
      "address": "workspace14a2hpadpsy9h4auve2z8lw",
      "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
      "owners": [
        "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
      ]
    }
  ],
  "pagination": {
    "total": "1"
  }
}
```

#### WorkspaceByAddr

The `WorkspaceByAddr` endpoint allows users to query for a single workspace by address.

```bash
zrchain.identity.Query/WorkspaceByAddress
```

Example:

```bash
grpcurl -plaintext \
  -d '{"workspace_addr": "workspace14a2hpadpsy9h4auve2z8lw"}' localhost:9090 zrchain.identity.Query/WorkspaceByAddress
```

Example Output:

```bash
{
  "workspace": {
    "address": "workspace14a2hpadpsy9h4auve2z8lw",
    "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
    "owners": [
      "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
    ]
  }
}
```

#### Keyrings

The `todo` endpoint allows users to query all keyrings.

```bash
zrchain.identity.Query/Keyrings 
```

Example:

```bash
grpcurl -plaintext localhost:9090 zrchain.identity.Query/Keyrings 
```

Example Output:

```bash
{
  "keyrings": [
    {
      "address": "keyring1k6vc6vhp6e6l3rxalue9v4ux",
      "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
      "description": "Keyring with Fees",
      "admins": [
        "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
      ],
      "parties": [
        "zen10kmgv5gzygnecf46x092ecfe5xcvvv9rdaxmts"
      ],
      "keyReqFee": "2",
      "sigReqFee": "2",
      "isActive": true
    },
    {
      "address": "keyring1pfnq7r04rept47gaf5cpdew2",
      "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
      "description": "ZenrockKMS",
      "admins": [
        "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
      ],
      "parties": [
        "zen10kmgv5gzygnecf46x092ecfe5xcvvv9rdaxmts"
      ],
      "isActive": true
    }
  ],
  "pagination": {
    "total": "2"
  }
}
```

#### KeyringByAddress

The `KeyringByAddress` endpoint allows users to query for a single keyring by address.

```bash
zrchain.identity.Query/KeyringByAddress
```

Example:

```bash
grpcurl -plaintext \
    -d '{"keyring_addr": "keyring1pfnq7r04rept47gaf5cpdew2"}' localhost:9090 zrchain.identity.Query/KeyringByAddress
```

Example Output:

```bash
{
  "keyring": {
    "address": "keyring1pfnq7r04rept47gaf5cpdew2",
    "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
    "description": "ZenrockKMS",
    "admins": [
      "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
    ],
    "parties": [
      "zen10kmgv5gzygnecf46x092ecfe5xcvvv9rdaxmts"
    ],
    "isActive": true
  }
}
```

### REST

A user can query the `identity` module using REST endpoints.

#### workspaces

The `workspaces` endpoint allows users to query all workspaces.

```bash
/zrchain/identity/workspaces
```

Example:

```bash
curl localhost:1317/zrchain/identity/workspaces
```

Example Output:

```json
{
  "workspaces": [
    {
      "address": "workspace14a2hpadpsy9h4auve2z8lw",
      "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
      "owners": [
        "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
      ],
      "child_workspaces": [],
      "admin_policy_id": "0",
      "sign_policy_id": "0",
      "alias": ""
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "1"
  }
}
```

#### workspace_by_address

The `workspace_by_address` endpoint allows users to query a single workspace by address.

```bash
/zrchain/identity/workspace_by_address/{workspace_addr}
```

Example:

```bash
curl http://localhost:1317/zrchain/identity/workspace_by_address/workspace14a2hpadpsy9h4auve2z8lw 
```

Example Output:

```json
{
  "workspace": {
    "address": "workspace14a2hpadpsy9h4auve2z8lw",
    "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
    "owners": [
      "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
    ],
    "child_workspaces": [],
    "admin_policy_id": "0",
    "sign_policy_id": "0",
    "alias": ""
  }
}
```

#### keyrings

The `keyrings` endpoint allows users to query all keyrings.

```bash
/zrchain/identity/keyrings
```

Example:

```bash
curl localhost:1317/zrchain/identity/keyrings
```

Example Output:

```bash
{
  "keyrings": [
    {
      "address": "keyring1k6vc6vhp6e6l3rxalue9v4ux",
      "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
      "description": "Keyring with Fees",
      "admins": [
        "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
      ],
      "parties": [
        "zen10kmgv5gzygnecf46x092ecfe5xcvvv9rdaxmts"
      ],
      "party_threshold": 0,
      "key_req_fee": "2",
      "sig_req_fee": "2",
      "is_active": true
    },
    {
      "address": "keyring1pfnq7r04rept47gaf5cpdew2",
      "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
      "description": "ZenrockKMS",
      "admins": [
        "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
      ],
      "parties": [
        "zen10kmgv5gzygnecf46x092ecfe5xcvvv9rdaxmts"
      ],
      "party_threshold": 0,
      "key_req_fee": "0",
      "sig_req_fee": "0",
      "is_active": true
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "2"
  }
}
```

#### keyring_by_address

The `keyring_by_address` endpoint allows users to query a single keyring by address.

```bash
/zrchain/identity/keyring_by_address/{keyring_addr}
```

Example:

```bash
curl localhost:1317/zrchain/identity/keyring_by_address/keyring1k6vc6vhp6e6l3rxalue9v4ux
```

Example Output:

```bash
{
  "keyring": {
    "address": "keyring1k6vc6vhp6e6l3rxalue9v4ux",
    "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
    "description": "Keyring with Fees",
    "admins": [
      "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
    ],
    "parties": [
      "zen10kmgv5gzygnecf46x092ecfe5xcvvv9rdaxmts"
    ],
    "party_threshold": 0,
    "key_req_fee": "2",
    "sig_req_fee": "2",
    "is_active": true
  }
}
```
