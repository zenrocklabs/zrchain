#!/bin/sh

solana_rpc="https://api.devnet.solana.com"

# Check if Solana CLI is installed
if ! command -v solana >/dev/null 2>&1; then
    echo "Error: Solana CLI is not installed. Please install it first."
    echo "Visit: https://docs.solana.com/cli/install-solana-cli-tools"
    exit 1
fi

# Check if RPC URL is set to devnet
RPC_URL=$(solana config get | grep "RPC URL" | awk '{print $3}')
if [ "$RPC_URL" != "https://api.devnet.solana.com" ]; then
    echo "Error: Solana config should be for devnet"
    echo "Current RPC URL: $RPC_URL"
    echo "Expected RPC URL: https://api.devnet.solana.com"
    echo "Run: solana config set --url https://api.devnet.solana.com"
    exit 1
fi

echo "Solana devnet address: $(solana address)"

# Check if zenrockd is available in PATH
if ! command -v zenrockd >/dev/null 2>&1; then
    echo "Error: zenrockd is not installed or not found in PATH"
    echo "Please install zenrockd first or ensure it's available in your PATH"
    exit 1
fi

echo "Please enter zrchain environment: local, devnet or testnet"
read -r environment

# Validate the input
case $environment in
    local|devnet|testnet)
        echo "Selected environment: $environment"
        if [ "$environment" = "local" ]; then
            rpc_url="http://localhost:26657"
            mint_address=$(zenrockd q zentp params --node $rpc_url | grep "mint_address:" | awk -F': ' '{print $2}' | tr -d ' ')
            program_id=$(zenrockd q zentp params --node $rpc_url | grep "program_id:" | awk -F': ' '{print $2}' | tr -d ' ')
            fee_wallet=$(zenrockd q zentp params --node $rpc_url | grep "fee_wallet:" | awk -F': ' '{print $2}' | tr -d ' ')
        elif [ "$environment" = "devnet" ]; then
            rpc_url="https://rpc.dev.zenrock.tech"
            mint_address=$(zenrockd q zentp params --node $rpc_url | grep "mint_address:" | awk -F': ' '{print $2}' | tr -d ' ')
            program_id=$(zenrockd q zentp params --node $rpc_url | grep "program_id:" | awk -F': ' '{print $2}' | tr -d ' ')
            fee_wallet=$(zenrockd q zentp params --node $rpc_url | grep "fee_wallet:" | awk -F': ' '{print $2}' | tr -d ' ')
        else
            rpc_url="https://rpc.gardia.zenrocklabs.io"
            mint_address=$(zenrockd q zentp params --node $rpc_url | grep "mint_address:" | awk -F': ' '{print $2}' | tr -d ' ')
            program_id=$(zenrockd q zentp params --node $rpc_url | grep "program_id:" | awk -F': ' '{print $2}' | tr -d ' ')
            fee_wallet=$(zenrockd q zentp params --node $rpc_url | grep "fee_wallet:" | awk -F': ' '{print $2}' | tr -d ' ')
        fi
        ;;
    *)
        echo "Error: Invalid environment. Please enter 'local', 'devnet' or 'testnet'"
        exit 1
        ;;
esac

# Get the keypair path from Solana config
keypair_path=$(solana config get | grep "Keypair Path" | awk -F': ' '{print $2}' | tr -d ' ')


# Get the public key using the keypair path
SIGNER_PUBKEY=$(solana-keygen pubkey "$keypair_path")
if [ $? -ne 0 ]; then
    echo "Failed to get signer public key"
    exit 1
fi

echo "Solana balance: $(solana balance)"
echo "urock balance: $(spl-token balance "$mint_address" --owner "$SIGNER_PUBKEY")"

echo "making unwrap transaction"

echo "enter destination address on zrchain:"
read -r dest_address

# Prompt for unwrap amount
echo "Please enter the amount of urock to unwrap:"
read -r unwrap_amount

# Validate the amount is a positive number
if ! [[ "$unwrap_amount" =~ ^[0-9]+(\.[0-9]+)?$ ]] || [ "$(echo "$unwrap_amount <= 0" | bc -l)" -eq 1 ]; then
    echo "Error: Please enter a valid positive number for the unwrap amount"
    exit 1
fi

echo "Unwrapping $unwrap_amount tokens from mint address: $mint_address"

# Capture the output and exit code
output=$(go run ../helper/create-unwrap.go \
    --signer "$keypair_path" \
    --mint "$mint_address" \
    --program "$program_id" \
    --rpc "$solana_rpc" \
    --amount "$unwrap_amount" \
    --dest "$dest_address" \
    --fee-wallet "$fee_wallet" \
)
exit_code=$?

echo "Go program exit code: $exit_code"
echo "Go program output: $output"

if [ $exit_code -ne 0 ]; then
    echo "Error: Go program failed with exit code $exit_code"
    exit 1
fi

echo "done"