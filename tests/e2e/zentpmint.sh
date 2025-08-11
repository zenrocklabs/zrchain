#!/bin/sh

fees="500000urock"
amount=1000000
solana_caip="solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1"
solana_rpc="https://api.devnet.solana.com"


echo "Checking if zenrockd is installed"

# Check if zenrockd is available in PATH
if ! command -v zenrockd >/dev/null 2>&1; then
    echo "Error: zenrockd is not installed or not found in PATH"
    echo "Please install zenrockd first or ensure it's available in your PATH"
    exit 1
fi

echo "zenrockd found, proceeding with execution"

echo "Please enter environment: local, devnet or testnet"
read -r environment

# Validate the input
case $environment in
    local|devnet|testnet)
        echo "Selected environment: $environment"#
        if [ "$environment" = "local" ]; then
            rpc_url="http://localhost:26657"
            chain_id="zenrock"
            mint_address=$(zenrockd q zentp params --node $rpc_url | grep "mint_address:" | awk '{print $2}')
        elif [ "$environment" = "devnet" ]; then
            rpc_url="https://rpc.dev.zenrock.tech"
            chain_id="amber-1"
            mint_address=$(zenrockd q zentp params --node $rpc_url | grep "mint_address:" | awk '{print $2}')
        else
            rpc_url="https://rpc.gardia.zenrocklabs.io"
            chain_id="gardia-9"
            mint_address=$(zenrockd q zentp params --node $rpc_url | grep "mint_address:" | awk '{print $2}')
        fi
        ;;
    *)
        echo "Error: Invalid environment. Please enter 'local', 'devnet' or 'testnet'"
        exit 1
        ;;
esac

echo "Please enter the sender name as listed in your zenrockd (default: alice)"
read -r from

# If nothing is entered, default to alice
if [ -z "$from" ]; then
    from="alice"
fi

echo "Using sender address: $from"
echo "Current balance: $(zenrockd query bank balances $from --node $rpc_url)"

echo "Please enter the recipient solana address"
read -r solana_address

echo "Selected receiver: $solana_address"

# Check if curl is available for RPC calls
if ! command -v curl >/dev/null 2>&1; then
    echo "Warning: curl not found, cannot check Solana balance"
else
    
    # Get Solana balance using RPC
    balance_response=$(curl -s -X POST "$solana_rpc" \
        -H "Content-Type: application/json" \
        -d "{
            \"jsonrpc\": \"2.0\",
            \"id\": 1,
            \"method\": \"getBalance\",
            \"params\": [\"$solana_address\"]
        }")
    
    # Extract balance from response (in lamports)
    balance_lamports=$(echo "$balance_response" | grep -o '"value":[0-9]*' | cut -d':' -f2)
    
    if [ -n "$balance_lamports" ] && [ "$balance_lamports" != "null" ]; then
        # Convert lamports to SOL (1 SOL = 1,000,000,000 lamports)
        balance_sol=$(echo "scale=9; $balance_lamports / 1000000000" | bc 2>/dev/null || echo "calculation_failed")
        
        if [ "$balance_sol" != "calculation_failed" ]; then
            echo "Solana balance: $balance_sol SOL"
        else
            echo "Solana balance: $balance_lamports lamports"
        fi
    else
        echo "Could not retrieve Solana balance"
        echo "RPC Response: $balance_response"
    fi
fi

# Check token balance for the receiver using the mint_address
if [ -n "$mint_address" ]; then
    
    # Check if curl is available for RPC calls
    if command -v curl >/dev/null 2>&1; then        
        # Get token account info using RPC
        token_response=$(curl -s -X POST "$solana_rpc" \
            -H "Content-Type: application/json" \
            -d "{
                \"jsonrpc\": \"2.0\",
                \"id\": 1,
                \"method\": \"getTokenAccountsByOwner\",
                \"params\": [
                    \"$solana_address\",
                    {
                        \"mint\": \"$mint_address\"
                    },
                    {
                        \"encoding\": \"jsonParsed\"
                    }
                ]
            }")
        
        # Extract token balance from response
        token_balance=$(echo "$token_response" | grep -o '"uiAmount":[0-9.]*' | cut -d':' -f2 | head -1)
        
        if [ -n "$token_balance" ] && [ "$token_balance" != "null" ]; then
            echo "urock balance: $token_balance urock"
        else
            echo "urock token balance: 0 tokens (or account not found)"
            echo "RPC Response: $token_response"
        fi
    else
        echo "Warning: curl not found, cannot check token balance"
    fi
else
    echo "Warning: Could not extract mint address from zentp params"
fi

echo "Type number of iterations:"
read -r iterations

# Validate that iterations is a positive number
if ! echo "$iterations" | grep -E '^[0-9]+$' >/dev/null 2>&1; then
    echo "Error: Iterations must be a positive number"
    exit 1
fi

if [ "$iterations" -le 0 ]; then
    echo "Error: Iterations must be greater than 0"
    exit 1
fi

echo "Number of iterations: $iterations"

echo "Type amount of urock to mint:"
read -r amount

echo "Starting minting..."

for i in $(seq 1 $iterations); do
    total_sent=$((total_sent + amount))
    echo "Iteration $i of $iterations"
    echo "Minting $amount urock to $solana_address"
    echo "Starting transaction..."
    
    # Execute transaction and capture response
    tx_response=$(zenrockd tx zentp bridge $amount urock $solana_caip $solana_address \
        --from "$from" \
        --chain-id "$chain_id" \
        --fees "$fees" \
        --node "$rpc_url" \
        --yes)
    
    # Check if transaction was successful
    if echo "$tx_response" | grep -q "code: 0"; then
        # Extract txhash for successful transactions
        txhash=$(echo "$tx_response" | grep "txhash:" | awk '{print $2}')
        echo "Transaction $i successful - Hash: $txhash"
    else
        # Extract raw_log for failed transactions
        raw_log=$(echo "$tx_response" | grep "raw_log:" | cut -d"'" -f2)
        echo "Transaction $i failed: $raw_log"
    fi
    
    sleep 7
    
done

echo "Total amount sent: $total_sent urock in $iterations iterations"


