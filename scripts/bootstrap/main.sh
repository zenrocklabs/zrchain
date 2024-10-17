#!/usr/bin/env bash

# Set DIR
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

# Global default variables
NODE_BIN="$DIR/../../cmd/zenrockd/zenrockd"
SCRIPT=$(basename "${BASH_SOURCE[0]}")
N_VALIDATORS=1
KEYRING_BACKEND="test"
CHAIN_ID="zenrock"
TOTAL_SUPPLY="1000000000000000"
BOND_COINS="1000000000000"
ARTIFACTS_DIR="$DIR/artifacts"
GENESIS="$ARTIFACTS_DIR/config/genesis.json"

usage() {
    echo "Usage: $SCRIPT"
    echo ""
    echo "-n <number>  -- number of validators"
    echo "-c <string>  -- chain ID"
    echo "-b <string>  -- amount of coins to bond"
    exit 1
}

while getopts "h?n:c:b:" args; do
    case $args in
        h|\?) usage;;
        n ) N_VALIDATORS=${OPTARG};;
        c ) CHAIN_ID=${OPTARG};;
        b ) BOND_COINS=${OPTARG};;
    esac
done

set -eo pipefail

# Check zenrockd binary
if [ ! -f "$NODE_BIN" ]; then
    echo "Missing zenrockd binary"
    exit 1
fi

# Run init to get default config files
$NODE_BIN init bootstrap --chain-id "$CHAIN_ID" --home "$ARTIFACTS_DIR" 1>/dev/null 2>&1

# Configure chain-id
$NODE_BIN config set client chain-id "$CHAIN_ID" --home "$ARTIFACTS_DIR"

# Configure keyring backend
$NODE_BIN config set client keyring-backend "$KEYRING_BACKEND" --home "$ARTIFACTS_DIR"

# Configure
$NODE_BIN config set app minimum-gas-prices "0urock" --home "$ARTIFACTS_DIR"

# Create Zenrock account
# TOTAL_SUPPLY - Validators bonded tokens
REMAINING=$((TOTAL_SUPPLY - N_VALIDATORS * BOND_COINS))
echo $REMAINING
$NODE_BIN keys add "zenrock_account" --home "$ARTIFACTS_DIR" 2>&1 | \
    tail -n1 > "$ARTIFACTS_DIR/mnemonic_zenrock_account"
$NODE_BIN genesis add-genesis-account "zenrock_account" "${REMAINING}urock" \
    --home "$ARTIFACTS_DIR"

# Bootstrap validator/s
for i in $(seq "$N_VALIDATORS"); do
    echo "Bootstraping validator $i"
    $NODE_BIN keys add "validator_operator_$i" --home "$ARTIFACTS_DIR" 2>&1 | \
        tail -n1 > "$ARTIFACTS_DIR/mnemonic_$i"
    $NODE_BIN genesis add-genesis-account "validator_operator_$i" "${BOND_COINS}urock" \
        --home "$ARTIFACTS_DIR"
    $NODE_BIN genesis gentx "validator_operator_$i" "${BOND_COINS}urock" \
        --keyring-backend "$KEYRING_BACKEND" --home "$ARTIFACTS_DIR" \
        --moniker "validator-$i" --note "validator-$i" --chain-id "${CHAIN_ID}"
    mv "$ARTIFACTS_DIR/config/node_key.json" "$ARTIFACTS_DIR/validator_${i}_node_key.json"
    mv "$ARTIFACTS_DIR/config/priv_validator_key.json" "$ARTIFACTS_DIR/validator_${i}_priv_validator_key.json"
done

# Genesis params
tmpfile=$(mktemp)
jq '.app_state["staking"]["params"]["bond_denom"]="urock"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"
jq '.app_state["crisis"]["constant_fee"]["denom"]="urock"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"
jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="urock"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"
jq '.app_state["mint"]["params"]["mint_denom"]="urock"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"
jq '.consensus_params["block"]["max_gas"]="20000000"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"
# Enable AVS
jq '.consensus["params"]["abci"]["vote_extensions_enable_height"]="1"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"

# TODO: replace these values depending on the chain (devnet, gardia, mainnet)
jq '.app_state["gov"]["params"]["voting_period"]="1800s"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"
jq '.app_state["gov"]["params"]["expedited_voting_period"]="1200s"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"
# Unbonding time to 86400s --> 1 day (devnet and testnet)
jq '.app_state["validation"]["params"]["unbonding_time"]="86400s"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"

# Collect generated txs
$NODE_BIN genesis collect-gentxs --home "$ARTIFACTS_DIR" 1>/dev/null 2>&1
rm -rf ./data

echo "Bootstrap finished"
