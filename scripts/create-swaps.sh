#!/usr/bin/env sh

# Make a rock key id
zenrockd tx treasury new-key-request workspace14a2hpadpsy9h4auve2z8lw keyring1pfnq7r04rept47gaf5cpdew2 ecdsa --from alice --chain-id zenrock --fees 500000urock --yes

sleep 15

# Fund this key with id 14
zenrockd tx bank send alice zen12cmmca4m4fgaeuu2qw0y8whppdlztehxyfv8sx 10000000000000urock --from alice --yes --chain-id zenrock --fees 500000urock 

sleep 15

# Make a btc key id for the workspace
zenrockd tx treasury new-key-request workspace14a2hpadpsy9h4auve2z8lw keyring1pfnq7r04rept47gaf5cpdew2 bitcoin --from alice --chain-id zenrock --fees 500000urock --yes

sleep 15

# Make a btc key id for the zenex pool
zenrockd tx treasury new-key-request workspace14a2hpadpsy9h4auve2z8lw keyring1pfnq7r04rept47gaf5cpdew2 bitcoin --from alice --chain-id zenrock --fees 500000urock --yes

sleep 15

# Make a first rock-btc swap request to fill the zenex module account
zenrockd tx zenex swap workspace14a2hpadpsy9h4auve2z8lw rock-btc 60000000000 14 15 --from alice --chain-id zenrock --yes --fees 500000urock 

sleep 15

# Make a btc-rock swap request over the minimum satoshis (5500)
zenrockd tx zenex swap workspace14a2hpadpsy9h4auve2z8lw btc-rock 6001 14 15 --from alice --chain-id zenrock --yes --fees 500000urock

sleep 15

# Get both swaps
zenrockd q zenex swaps

