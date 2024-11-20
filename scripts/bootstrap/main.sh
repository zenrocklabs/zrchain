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
MINT_COINS="80000000000000"
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

# Set initial mint parameters
jq '.app_state.mint.params = {
    "mint_denom": "urock",
    "inflation_rate_change": "0.000000000000000000",
    "inflation_max": "0.000000000000000000",
    "inflation_min": "0.000000000000000000",
    "goal_bonded": "0.670000000000000000",
    "blocks_per_year": "6311520",
}' $ARTIFACTS_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $ARTIFACTS_DIR/config/genesis.json

# Create Zenrock account
# TOTAL_SUPPLY - Validators bonded tokens
REMAINING=$((TOTAL_SUPPLY - MINT_COINS - N_VALIDATORS * BOND_COINS))
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

# Add funds for mint module - not a validator
$NODE_BIN genesis add-genesis-account "zen1m3h30wlvsf8llruxtpukdvsy0km2kum8ju4et3" "80000000000000urock" --keyring-backend "$KEYRING_BACKEND" --home "$ARTIFACTS_DIR" --module-name mint

# Genesis params
tmpfile=$(mktemp)
# jq '.app_state["staking"]["params"]["bond_denom"]="urock"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"
# jq '.app_state["crisis"]["constant_fee"]["denom"]="urock"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"
# jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="urock"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"
# jq '.app_state["mint"]["params"]["mint_denom"]="urock"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"
jq '.consensus_params["block"]["max_gas"]="20000000"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"
# Enable AVS
jq '.consensus["params"]["abci"]["vote_extensions_enable_height"]="1"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"

# TODO: replace these values depending on the chain (devnet, gardia, mainnet)
# jq '.app_state["gov"]["params"]["voting_period"]="1800s"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"
# jq '.app_state["gov"]["params"]["expedited_voting_period"]="1200s"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"
# Unbonding time to 86400s --> 1 day (devnet and testnet)
jq '.app_state["validation"]["params"]["unbonding_time"]="86400s"' < "$GENESIS" > "$tmpfile" && mv "$tmpfile" "$GENESIS"

jq '.app_state.treasury.params = {
  "mpc_keyring": "",
  "zr_sign_address": "",
  "keyring_commission": 10,
  "keyring_commission_destination": "zen1m3h30wlvsf8llruxtpukdvsy0km2kum8ju4et3",
  "min_gas_fee": "2.5urock"
}' $ARTIFACTS_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $ARTIFACTS_DIR/config/genesis.json

# Update governance parameters
jq '.app_state.gov.params = {
    "min_deposit": [
        {
            "denom": "urock",
            "amount": "1000000000"
        }
    ],
    "max_deposit_period": "172800s",
    "voting_period": "345600s",
    "quorum": "0.334000000000000000",
    "threshold": "0.500000000000000000",
    "veto_threshold": "0.334000000000000000",
    "min_initial_deposit_ratio": "0.000000000000000000",
    "proposal_cancel_ratio": "0.500000000000000000",
    "proposal_cancel_dest": "",
    "expedited_voting_period": "86400s",
    "expedited_threshold": "0.667000000000000000",
    "expedited_min_deposit": [
        {
            "denom": "urock",
            "amount": "5000000000"
        }
    ],
    "burn_vote_quorum": false,
    "burn_proposal_deposit_prevote": false,
    "burn_vote_veto": true,
    "min_deposit_ratio": "0.010000000000000000"
}' $ARTIFACTS_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $ARTIFACTS_DIR/config/genesis.json

# Add distribution parameters
jq '.app_state.distribution.params = {
    "community_tax": "0.020000000000000000",
    "base_proposer_reward": "0.000000000000000000",
    "bonus_proposer_reward": "0.000000000000000000",
    "withdraw_addr_enabled": false
}' $ARTIFACTS_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $ARTIFACTS_DIR/config/genesis.json

# Add crisis parameters
jq '.app_state.crisis = {
    "constant_fee": {
        "denom": "urock",
        "amount": "1000000000"
    }
}' $ARTIFACTS_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $ARTIFACTS_DIR/config/genesis.json

# Add denom_metadata to the bank
jq '.app_state.bank.denom_metadata = [
    {
        "description": "The native staking token of the Zenrock blockchain",
        "denom_units": [
            {
                "denom": "urock",
                "exponent": 0,
                "aliases": ["microrock"]
            },
            {
                "denom": "rock",
                "exponent": 6,
                "aliases": []
            }
        ],
        "base": "urock",
        "display": "rock",
        "name": "ROCK",
        "symbol": "ROCK"
    }
]' $ARTIFACTS_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $ARTIFACTS_DIR/config/genesis.json

# Add identity parameters
jq '.app_state.identity.params = {
    "keyring_creation_fee": "10000000000"
}' $ARTIFACTS_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $ARTIFACTS_DIR/config/genesis.json

# Add policy parameters
jq '.app_state.policy.params = {
    "minimum_btl": "10",
    "default_btl": "1000"
}' $ARTIFACTS_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $ARTIFACTS_DIR/config/genesis.json

# Add slashing parameters
jq '.app_state.slashing.params = {
    "signed_blocks_window": "10000",
    "min_signed_per_window": "0.500000000000000000",
    "downtime_jail_duration": "600s",
    "slash_fraction_double_sign": "0.050000000000000000",
    "slash_fraction_downtime": "0.005000000000000000"
}' $ARTIFACTS_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $ARTIFACTS_DIR/config/genesis.json

# Add validation parameters
jq '.app_state.validation.params = {
    "unbonding_time": "1814400s",
    "max_validators": 75,
    "max_entries": 7,
    "historical_entries": 10000,
    "bond_denom": "urock",
    "min_commission_rate": "0.050000000000000000",
    "HVParams": {
        "AVSRewardsRate": "0.00",
        "BlockTime": 5,
        "ZenBTCParams": {
            "zenBTCEthContractAddr": "",
            "zenBTCDepositKeyringAddr": "",
            "zenBTCMinterKeyID": 0
        }
    }
}' $ARTIFACTS_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $ARTIFACTS_DIR/config/genesis.json

# Add wasm parameters
jq '.app_state.wasm.params = {
    "code_upload_access": {
        "permission": "Nobody",
        "addresses": []
    },
    "instantiate_default_permission": "Nobody"
}' $ARTIFACTS_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $ARTIFACTS_DIR/config/genesis.json

# Collect generated txs
$NODE_BIN genesis collect-gentxs --home "$ARTIFACTS_DIR" 1>/dev/null 2>&1
rm -rf ./data

echo "Bootstrap finished"
