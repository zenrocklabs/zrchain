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
zenrockd tx wasm store target/wasm32-unknown-unknown/release/parser_example.wasm --from alice $TXFLAG -y 

get code id:
zenrockd query wasm list-code 
export CODE_ID=1

instantiate the contract:
zenrockd tx wasm instantiate $CODE_ID "{}" --from alice --label "example contract" $TXFLAG -y --no-admin

get contract address:
zenrockd query wasm list-contract-by-code $CODE_ID 
export CONTRACT=zen14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s38wvxu

calling the contract from go code:

sudomsg := contractcall.SudoMessage{
    ParseInput: contractcall.ParseInput{
        Input: []byte{1, 2, 3, 4},
    },
}

addr, err := sdk.AccAddressFromBech32("zen14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s38wvxu")
if err == nil {
    res, err := contractcall.Sudo(k.wasmKeeperFn(), sdk.UnwrapSDKContext(goCtx), addr, sudomsg)
    if err != nil {
        k.Logger().Debug(err.Error())
    }
    k.Logger().Debug(fmt.Sprintf("%v", res))
}

```