#!/bin/sh

key_types=(
	"ecdsa"    #1
	"ecdsa"    #2
	"bitcoin"  #3
	"ecdsa"    #4
	"ecdsa"    #5
	"ecdsa"    #6
	"ecdsa"    #7
	"ed25519"  #8
	"ed25519"  #9
	"ed25519"  #10
	"ed25519"  #11
        "ed25519"  #12
        "ed25519"  #13
) 

workspace="workspace14a2hpadpsy9h4auve2z8lw"
keyring="keyring1pfnq7r04rept47gaf5cpdew2"
from="alice"
chain_id="zenrock"
fees="10000urock"
backend="test"


for key_type in "${key_types[@]}"; do
    echo "Requesting key of type: $key_type"
    zenrockd tx treasury new-key-request "$workspace" "$keyring" "$key_type" \
        --from "$from" \
        --chain-id "$chain_id" \
        --fees "$fees" \
        --keyring-backend "$backend" \
        --yes

    echo "Sleeping for 15 seconds..."
    sleep 15
done

# Only runs if alice is the admin key
zenrockd tx zentp set-solana-rock-supply 100000000000000 --from alice --chain-id zenrock --fees 500000urock --yes

echo "mammal romance obtain swarm disorder snake apology debris daughter magnet column scrub crowd drift empty rebuild address first patch believe myself grow aware muffin" |zenrockd keys add btcproxy --recover --keyring-backend test
sleep 5
zenrockd tx bank send alice zen1trdxe6r48aqvhm026akay7tjnzuarf2rxuz0ah  10000000000urock --keyring-backend test --chain-id zenrock  --fees 200urock --yes
sleep 5
zenrockd tx identity add-workspace-owner workspace14a2hpadpsy9h4auve2z8lw zen1trdxe6r48aqvhm026akay7tjnzuarf2rxuz0ah --chain-id zenrock --btl 100 --fees 20urock --keyring-backend test --from alice --yes

