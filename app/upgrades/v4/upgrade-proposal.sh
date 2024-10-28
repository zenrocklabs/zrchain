#!/bin/bash 

# Add a policy with old format
zenrockd tx policy new-policy alice_bob '{"@type":"/zrchain.policy.BoolparserPolicy", "definition": "u1 + u2 > 1", "participants":[{ "abbreviation":"u1", "address":"zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty" },{ "abbreviation":"u2", "address":"zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq" }]}' --yes --from alice --chain-id zenrock --gas-prices 0.0001urock
# Send SoftwareUpgrade proposal - Upgrade Name: v4
zenrockd tx upgrade software-upgrade v4 --title upgrade --summary upgrade --upgrade-height 100 --upgrade-info "{}" --no-validate --deposit 10000000urock --from alice --yes --chain-id zenrock --gas-prices 0.0001urock
# Deposit for the proposal - Proposal ID: 1
zenrockd tx gov deposit 1 10000000urock --from alice --yes --chain-id zenrock --gas-prices 0.0001urock
# Vote for the proposal
zenrockd tx gov vote 1 yes --from alice --yes --chain-id zenrock --gas-prices 0.0001urock
# Check gov vote status
zenrockd q gov proposal 1
# Check gov vote status
zenrockd q gov votes 1

