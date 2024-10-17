# Empty Example Contract

example derived from:
https://github.com/CosmWasm/cosmwasm/tree/main/contracts/empty

more info: 
https://docs.cosmwasm.com/docs/getting-started/interact-with-contract

```
install rust:
brew install rustup-init
rustup-init

compile the contract:
cargo wasm

setup env:
export RPC="http://localhost:26657"
export NODE=(--node $RPC)
export CHAIN_ID="zenrock"
export FEE_DENOM="urock"
export TXFLAG=($NODE --chain-id $CHAIN_ID --gas-prices 10$FEE_DENOM --gas auto --gas-adjustment 1.3)
export ZEN_KEYRING_ID=keyring1pfnq7r04rept47gaf5cpdew2

uploading the wasm binary:
zenrockd tx wasm store target/wasm32-unknown-unknown/release/contract_example.wasm --from alice $TXFLAG -y 

get code id:
zenrockd query wasm list-code 
export CODE_ID=1

instantiate the contract:
zenrockd tx wasm instantiate $CODE_ID "{}" --from alice --label "example contract" $TXFLAG -y --no-admin

get contract address:
zenrockd query wasm list-contract-by-code $CODE_ID 
export CONTRACT=zen14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s38wvxu

get the state of the contract:
zenrockd query wasm contract-state all $CONTRACT

create a workspace:
zenrockd tx wasm execute $CONTRACT '{"new_workspace_request":{ "admin_policy_id":"0", "sign_policy_id": "0" }}' --from alice $TXFLAG -y

create a key:
zenrockd tx wasm execute $CONTRACT '{"new_key_request":{ "workspace_addr": "workspace14a2hpadpsy9h4auve2z8lw", "keyring_addr": "keyring1pfnq7r04rept47gaf5cpdew2", "key_type": "ecdsa" }}' --from alice $TXFLAG -y

Make sure zenrockkms is running

zenrockd q treasury key-requests --keyring-addr $ZEN_KEYRING_ID

request a signature for data:
zenrockd tx wasm execute $CONTRACT '{"new_sign_data_request":{ "key_id":"1", "data_for_signing": "3132333435363738393031323334353637383930313233343536373839303132" }}' --from alice $TXFLAG -y

check the signature request:
zenrockd q treasury signature-requests --keyring-addr $ZEN_KEYRING_ID

request a signature for a transaction:
zenrockd tx wasm execute $CONTRACT '{"new_sign_transaction_request":{ "key_id":"1", "wallet_type": "eth", "unsigned_transaction": "6QKFAVlvlfGCUgiUX/E31LD9zUncowx89X5XigJtJ4mF6NSlEACAgICA", "metadata": { "type_url": "/zrchain.treasury.MetadataEthereum", "value": "CKftqAU=" } }}' --from alice $TXFLAG -y

add a workspace owner:
zenrockd tx wasm execute $CONTRACT '{"add_workspace_owner_request":{ "workspace_addr":"workspace14a2hpadpsy9h4auve2z8lw", "new_owner": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq" } }' --from alice $TXFLAG -y

check the workspace added owner:
zenrockd q identity workspace-by-address workspace14a2hpadpsy9h4auve2z8lw

query a keyring by address:
zenrockd query wasm contract-state smart $CONTRACT '{"keyring_by_address_query": { "keyring_addr": "keyring1pfnq7r04rept47gaf5cpdew2" }}' $NODE --output json

```