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

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --localnet)
            LOCALNET=$2
            shift 2
            ;;
        --alternate-home)
            ALTERNATE_HOME=true
            shift
            ;;
        --non-validator)
            NON_VALIDATOR=true
            NO_SIDECAR="--no-sidecar"
            shift
            ;;
        --no-sidecar)
            NO_SIDECAR="--no-sidecar"
            shift
            ;;
        --no-vote-extensions)
            NO_VOTE_EXTENSIONS=true
            NO_SIDECAR="--no-sidecar"
            shift
            ;;
        --alternate-ports)
            ALTERNATE_PORTS="--alternate-ports"
            shift
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Adjust settings based on LOCALNET flag
if [ -n "$LOCALNET" ]; then
    if [ "$LOCALNET" -eq 1 ]; then
        # Validator 1: Normal flow but waits for gentx
        WAIT_FOR_GENTX=true
    elif [ "$LOCALNET" -eq 2 ]; then
        # Validator 2: Uses alternate ports and home directory
        ALTERNATE_HOME=true
        ALTERNATE_PORTS="--alternate-ports"
    else
        echo "Invalid value for --localnet. Use 1 or 2."
        exit 1
    fi
fi

# Set the appropriate home directory
if [ "$ALTERNATE_HOME" = true ]; then
    HOME_DIR=$ALTERNATE_HOME_DIR
fi

# Set the moniker based on node type and home directory
if [ "$NON_VALIDATOR" = true ]; then
    MONIKER="zenrock_non_validator"
elif [ "$ALTERNATE_HOME" = true ]; then
    MONIKER="zenrock_alt"
else
    MONIKER="zenrock"
fi

# Function to automate key recovery using expect
function recover_key_with_mnemonic() {
    local KEY_NAME=$1
    local MNEMONIC=$2
    local KEYRING=$3
    local HOME_DIR=$4

    expect << EOF
spawn zenrockd keys add $KEY_NAME --recover --keyring-backend $KEYRING --home $HOME_DIR
expect "Enter your bip39 mnemonic"
send "$MNEMONIC\r"
expect eof
EOF
}

# Clean up existing data directories and files
if [ "$NON_VALIDATOR" = false ] && [ "$LOCALNET" != "2" ]; then
    # Only Validator 1 cleans up shared directories
    echo "Cleaning up old gentx files and genesis.json..."
    rm -rf ./gentx
    rm -f ./genesis.json
fi

# Remove existing daemon and client data
rm -rf $HOME_DIR

set -e

make install

# Only make sidecar if not using --localnet 2 or --no-sidecar/--no-vote-extensions
if [[ "$LOCALNET" != "2" && -z "$NO_SIDECAR" && -z "$NO_VOTE_EXTENSIONS" ]]; then
    make sidecar
fi

rm -rf sidecar/neutrino/neutrino_*/*.bin
rm -rf sidecar/neutrino/neutrino_*/*.db

if [ "$NON_VALIDATOR" = false ]; then
    # Add keys for Alice and Bob using their mnemonics
    recover_key_with_mnemonic $K1 "$MNEMONIC1" $KEYRING $HOME_DIR
    recover_key_with_mnemonic $K2 "$MNEMONIC2" $KEYRING $HOME_DIR
fi

# Initialize the node
zenrockd init $MONIKER --chain-id $CHAINID --home $HOME_DIR

function ssed {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        gsed "$@"
    else
        sed "$@"
    fi
}

if [ "$NON_VALIDATOR" = false ]; then
    if [ "$LOCALNET" != "2" ]; then
        # First validator node in localnet or default single-node flow: allocate genesis accounts
        zenrockd genesis add-genesis-account $K1 2000000000000000urock --keyring-backend $KEYRING --home $HOME_DIR
        zenrockd genesis add-genesis-account $K2 2000000000000000urock --keyring-backend $KEYRING --home $HOME_DIR
        zenrockd genesis add-genesis-account zen10kmgv5gzygnecf46x092ecfe5xcvvv9rdaxmts 1000000000000000urock --keyring-backend $KEYRING --home $HOME_DIR
        zenrockd genesis add-genesis-account zen1zpmqphp46nsn097ysltk4j5wmpjn9gd5gwyfnq 1000000000000000urock --keyring-backend $KEYRING --home $HOME_DIR
    else
        # Second validator node: copy genesis.json
        cp $VALIDATOR_HOME/config/genesis.json $HOME_DIR/config/genesis.json
    fi
else
    echo "Setting up as non-validator node"
    # Copy genesis file from validator node
    cp $VALIDATOR_HOME/config/genesis.json $HOME_DIR/config/genesis.json
fi

# Set block time to 1s
ssed -i 's|timeout_commit = "5s"|timeout_commit = "1s"|g' $HOME_DIR/config/config.toml

# Adjust ports if alternate ports are specified
if [[ -n "$ALTERNATE_PORTS" ]]; then
    # Change ports to avoid conflicts with other chains running locally
    ssed -i 's|26656|27656|g' $HOME_DIR/config/config.toml
    ssed -i 's|26657|27657|g' $HOME_DIR/config/config.toml
    ssed -i 's|6060|6760|g' $HOME_DIR/config/config.toml
    ssed -i 's|1317|1717|g' $HOME_DIR/config/app.toml
    ssed -i 's|9090|9790|g' $HOME_DIR/config/app.toml
    ssed -i 's|9091|9791|g' $HOME_DIR/config/app.toml
fi

# Configure persistent peers
if [ "$LOCALNET" = "2" ] || [ "$NON_VALIDATOR" = true ]; then
    # Get the first node's ID and P2P address
    VALIDATOR_NODE_ID=$(zenrockd tendermint show-node-id --home $VALIDATOR_HOME)
    VALIDATOR_LISTEN_ADDR=$(grep -A 3 "\[p2p\]" $VALIDATOR_HOME/config/config.toml | grep "laddr = " | cut -d '"' -f 2)
    # Extract IP and port from VALIDATOR_LISTEN_ADDR
    VALIDATOR_IP=$(echo $VALIDATOR_LISTEN_ADDR | cut -d ':' -f 2 | tr -d '/')
    VALIDATOR_PORT=$(echo $VALIDATOR_LISTEN_ADDR | cut -d ':' -f 3)
    # Construct the correct persistent_peers string
    PERSISTENT_PEERS="${VALIDATOR_NODE_ID}@127.0.0.1:${VALIDATOR_PORT}"
    # Update the config.toml file
    ssed -i "s|^persistent_peers *=.*|persistent_peers = \"$PERSISTENT_PEERS\"|" $HOME_DIR/config/config.toml
fi

if [ "$NON_VALIDATOR" = false ]; then
    # Create gentx for the validator
    if [ "$LOCALNET" = "1" ] || [ -z "$LOCALNET" ]; then
        # First validator node or default
        zenrockd genesis gentx $K2 1000000000000000urock --keyring-backend $KEYRING --chain-id $CHAINID --home $HOME_DIR
    elif [ "$LOCALNET" = "2" ]; then
        # Second validator node
        zenrockd genesis gentx $K1 1000000000000000urock --keyring-backend $KEYRING --chain-id $CHAINID --home $HOME_DIR
    fi

    if [ -n "$LOCALNET" ]; then
        # Create the gentx directory if it doesn't exist
        mkdir -p ./gentx
        # Copy gentx files to the shared directory
        cp $HOME_DIR/config/gentx/*.json ./gentx/

        # Wait for all gentx files to be present
        if [ "$LOCALNET" = "1" ]; then
            echo "Waiting for other gentx files..."
            while [ $(ls ./gentx/*.json 2>/dev/null | wc -l) -lt 2 ]; do
                sleep 1
            done
        fi
    fi
fi

# Only the first validator collects gentxs and finalizes genesis.json
if [ "$NON_VALIDATOR" = false ] && ( [ "$LOCALNET" = "1" ] || [ -z "$LOCALNET" ] ); then
    # Copy all gentx files to the first validator's gentx directory
    if [ -n "$LOCALNET" ]; then
        cp ./gentx/*.json $HOME_DIR/config/gentx/
    fi

    # Collect genesis transactions
    zenrockd genesis collect-gentxs --home $HOME_DIR

    # Apply necessary modifications to genesis.json here
    # Set vote_extensions_enable_height
    if [ -z "$NO_VOTE_EXTENSIONS" ]; then
        jq '.consensus.params.abci.vote_extensions_enable_height = "1"' $HOME_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $HOME_DIR/config/genesis.json
    fi

    # Apply other necessary modifications from your original script
    jq '.app_state.identity.keyrings = [
      {
        "address": "keyring1pfnq7r04rept47gaf5cpdew2",
        "admins": ["zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"],
        "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
        "description": "ZenrockKMS",
        "is_active": true,
        "key_req_fee": 0,
        "parties": ["zen10kmgv5gzygnecf46x092ecfe5xcvvv9rdaxmts"],
        "sig_req_fee": 0
      },
      {
        "address": "keyring1k6vc6vhp6e6l3rxalue9v4ux",
        "admins": ["zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"],
        "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
        "description": "Keyring with Fees",
        "is_active": true,
        "key_req_fee": 2,
        "parties": ["zen10kmgv5gzygnecf46x092ecfe5xcvvv9rdaxmts"],
        "sig_req_fee": 2
      },
      {
        "address": "keyring1k6vc6vhp6e6l3rxaard6fd",
        "admins": ["zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"],
        "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
        "description": "TSS one",
        "is_active": true,
        "key_req_fee": 0,
        "parties": ["zen1qwnafe2s9eawhah5x6v4593v3tljdntl9zcqpn"],
        "sig_req_fee": 0
      }
    ]' $HOME_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $HOME_DIR/config/genesis.json

    jq '.app_state.identity.workspaces = [
      {
        "address": "workspace14a2hpadpsy9h4auve2z8lw",
        "creator": "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
        "owners": ["zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"]
      }
    ]' $HOME_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $HOME_DIR/config/genesis.json

    jq '.app_state.staking.params.bond_denom = "urock"' $HOME_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $HOME_DIR/config/genesis.json

    jq '.app_state.treasury.params = {
      "mpc_keyring": "keyring1pfnq7r04rept47gaf5cpdew2",
      "zr_sign_address": "zen1zpmqphp46nsn097ysltk4j5wmpjn9gd5gwyfnq",
      "keyring_commission": 10,
      "keyring_commission_destination": "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
      "min_gas_fee": "0.0001urock"
    }' $HOME_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $HOME_DIR/config/genesis.json

    jq '.app_state.gov.params.voting_period = "60s"' $HOME_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $HOME_DIR/config/genesis.json
    jq '.app_state.gov.params.expedited_voting_period = "30s"' $HOME_DIR/config/genesis.json > tmp_genesis.json && mv tmp_genesis.json $HOME_DIR/config/genesis.json

    # Validate the genesis file
    zenrockd genesis validate --home $HOME_DIR

    # Distribute the updated genesis.json to all nodes
    cp $HOME_DIR/config/genesis.json ./genesis.json
fi

# Wait for the genesis.json file to be finalized
if [ -n "$LOCALNET" ]; then
    while [ ! -f ./genesis.json ]; do
        sleep 1
    done

    # Copy the finalized genesis.json to all nodes
    cp ./genesis.json $HOME_DIR/config/genesis.json

    # Validate the genesis file
    zenrockd genesis validate --home $HOME_DIR
fi

# Start the Oracle Sidecar only for the first validator node
if [[ -z "$NO_SIDECAR" && "$NON_VALIDATOR" = false && ( "$LOCALNET" != "2" || -z "$LOCALNET" ) ]]; then
    (cd sidecar && ./sidecar > sidecar.log 2>&1 &)
    echo -e "\nStarting Oracle Sidecar, sleeping for 25 seconds to allow Sidecar to start pulling data...\n"
    sleep 25
fi

# Start the node
if [ "$NON_VALIDATOR" = true ]; then
    zenrockd start --home $HOME_DIR --pruning=nothing --log_level $LOGLEVEL \
    --minimum-gas-prices=0.0001urock --api.enable --api.enabled-unsafe-cors --non-validator
else
    zenrockd start --home $HOME_DIR --pruning=nothing --log_level $LOGLEVEL \
    --minimum-gas-prices=0.0001urock --api.enable --api.enabled-unsafe-cors
fi
