#!/bin/sh

key_types=(
	"ed25519"  #1
	"ecdsa"    #2
	"ecdsa"    #3
	"ecdsa"    #4
	"ecdsa"    #5
	"ecdsa"    #6
	"ecdsa"    #7
	"ed25519"  #8
	"ed25519"  #9
	"ed25519"  #10
	"ed25519"  #11
        "ed25519"  #12
	"bitcoin"  #13
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

