#!/usr/bin/env bash

K1="alice"
K2="bob"
CHAINID="zenrock"
KEYRING="test"
LOGLEVEL="info"
NON_VALIDATOR=false
ALTERNATE_HOME=false
ALTERNATE_PORTS=""
LOCALNET=""
HOME_DIR="$HOME/.zrchain"
ALTERNATE_HOME_DIR="$HOME/.zrchain_alt"
VALIDATOR_HOME="$HOME/.zrchain"
SIDECAR_ADDR=""
MNEMONIC1="strategy social surge orange pioneer tiger skill endless lend slide one jazz pipe expose detect soup fork cube trigger frown stick wreck ring tissue"
MNEMONIC2="fee buzz avocado dolphin syrup rule access cave close puppy lemon round able bronze fame give spoon company since fog error trip toast unable"

# Start the node with optional flags
zenrockd start --home $HOME_DIR --pruning=nothing --log_level $LOGLEVEL \
--minimum-gas-prices=0.0001urock --api.enable --api.enabled-unsafe-cors \
${NON_VALIDATOR:+--non-validator} \
${SIDECAR_ADDR:+--sidecar-addr "$SIDECAR_ADDR"}
