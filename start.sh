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
MNEMONIC1="strategy social surge orange pioneer tiger skill endless lend slide one jazz pipe expose detect soup fork cube trigger frown stick wreck ring tissue"
MNEMONIC2="fee buzz avocado dolphin syrup rule access cave close puppy lemon round able bronze fame give spoon company since fog error trip toast unable"
MNEMONIC3="exclude try nephew main caught favorite tone degree lottery device tissue tent ugly mouse pelican gasp lava flush pen river noise remind balcony emerge"


# Start the node
if [ "$NON_VALIDATOR" = true ]; then
    zenrockd start --home $HOME_DIR --pruning=nothing --log_level $LOGLEVEL \
    --minimum-gas-prices=0.0001urock --api.enable --api.enabled-unsafe-cors --non-validator
else
    zenrockd start --home $HOME_DIR --pruning=nothing --log_level $LOGLEVEL \
    --minimum-gas-prices=0.0001urock --api.enable --api.enabled-unsafe-cors
fi
