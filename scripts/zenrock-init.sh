#!/usr/bin/env sh

# install jq
apk add jq

# Clean up previous config
rm -rf ~/.zrchain

# Main config
/zenrockd init bootstrap --chain-id zenrock --overwrite
/zenrockd config set app minimum-gas-prices "2.5urock"
/zenrockd config set client chain-id zenrock
/zenrockd config set client keyring-backend test

# Change creation of blocks to 2s
sed -i 's/timeout_commit = "5s"/timeout_commit = "2s"/' ~/.zrchain/config/config.toml

# Set min_gas_fee
jq '.app_state.treasury.params.min_gas_fee = "2.5urock"' ~/.zrchain/config/genesis.json > /tmp/genesis.json && mv /tmp/genesis.json ~/.zrchain/config/genesis.json

# Set inflation to 0
jq '.app_state.mint.params = {
        "mint_denom": "urock",
        "inflation_rate_change": "0.000000000000000000",
        "inflation_max": "0.000000000000000000",
        "inflation_min": "0.000000000000000000",
        "goal_bonded": "0.670000000000000000",
        "blocks_per_year": "6311520",
        "staking_yield": "0.070000000000000000",
        "burn_rate": "0.100000000000000000",
        "protocol_wallet_rate": "0.300000000000000000",
        "protocol_wallet_address": "zen1vh2gdma746t88y7745qawy32m0qxx60gjw27jj",
        "additional_staking_rewards": "0.300000000000000000",
        "additional_mpc_rewards": "0.050000000000000000",
        "additional_burn_rate": "0.250000000000000000"
    }' ~/.zrchain/config/genesis.json > /tmp/genesis.json && mv /tmp/genesis.json ~/.zrchain/config/genesis.json


# Add validator account (zen14mh64dy233ar0z0z40fh84fz25k3ckmtpnc0pz)
echo "direct category any favorite alert symbol consider name always term patrol initial join profit mule arena glare problem whale critic choice world zebra inherit" | /zenrockd keys add "validator" --recover
/zenrockd genesis add-genesis-account "validator" 100000000000000000000000000urock
/zenrockd genesis add-genesis-account zen1m3h30wlvsf8llruxtpukdvsy0km2kum8ju4et3 100000000000000urock --module-name mint
/zenrockd genesis add-genesis-account zen1234wz2aaavp089ttnrj9jwjqraaqxkkadq0k03 0urock --module-name zenex_collector
/zenrockd genesis add-genesis-account zen1fpq2t9ygrst5lp5hl9d7fylppljp3xhhu37n4c 0urock --module-name zenex_fee_collector
/zenrockd genesis add-genesis-account zen14l8vvehfy0af0djxjx0uug0ladm57r6plfntx6 0urock --module-name zenex_btc_rewards_collector
/zenrockd genesis gentx "validator" 1000000000000urock --moniker "validator"
/zenrockd genesis collect-gentxs


# Start blockchain in background
/zenrockd start --grpc.address 0.0.0.0:9790 --rpc.laddr tcp://0.0.0.0:26657 &
zenrock_pid=$!

# Zenrock chain startup time
sleep 10

# Create connectors accounts
# connector-1: zen1qwnafe2s9eawhah5x6v4593v3tljdntl9zcqpn
# connector-2: zen1pyt0d94842803uyvhnt2aapzrf99mkje3ykuxg
# connector-3: zen1lm70sewfh30xsa3glpj33k2phvzdur43uztzqm
echo "top decade spare horn skin actor balcony swim prefer hood divert run sick save excess siege market proud enforce wood lecture drive near odor" | /zenrockd keys add "connector-1" --recover
echo "sausage diesel never robot balcony tube typical clap scrap little few try shock charge plunge creek quiz advance hub bomb border tape ecology scale" | /zenrockd keys add "connector-2" --recover
echo "peace gym gown lab hand lens grain tide faint actor artist desk guess length million clarify walnut foam satoshi alarm title elevator stairs fetch" | /zenrockd keys add "connector-3" --recover

# Create test user account (zen1tmwedgn3t3rn3ryzdjw78tev8khfhrhfemggw7)
echo "deal liberty toilet artefact during key home calm vanish shock paddle sustain still rotate tonight spoon insane isolate pistol prosper critic kidney diesel crack" | /zenrockd keys add "user-tests" --recover

# Send tokens to accounts
/zenrockd tx bank send validator zen1qwnafe2s9eawhah5x6v4593v3tljdntl9zcqpn 100000000000urock --fees 500000urock --yes
sleep 2
/zenrockd tx bank send validator zen1pyt0d94842803uyvhnt2aapzrf99mkje3ykuxg 100000000000urock --fees 500000urock --yes
sleep 2
/zenrockd tx bank send validator zen1lm70sewfh30xsa3glpj33k2phvzdur43uztzqm 100000000000urock --fees 500000urock --yes
sleep 2
/zenrockd tx bank send validator zen1tmwedgn3t3rn3ryzdjw78tev8khfhrhfemggw7 100000000000urock --fees 500000urock --yes
sleep 2

# Create keyring
# TSS v2: keyring1k6vc6vhp6e6l3rxalue9v4ux
/zenrockd tx identity new-keyring "TSS v2" 0 0 0 --from validator --fees 500000urock --party-threshold 3 --yes
sleep 2

# Add connector-1 and connector-2 to keyring TSS v2
/zenrockd tx identity add-keyring-party keyring1k6vc6vhp6e6l3rxalue9v4ux zen1qwnafe2s9eawhah5x6v4593v3tljdntl9zcqpn --from validator --fees 500000urock --yes
sleep 2
/zenrockd tx identity add-keyring-party keyring1k6vc6vhp6e6l3rxalue9v4ux zen1pyt0d94842803uyvhnt2aapzrf99mkje3ykuxg --from validator --fees 500000urock --yes
sleep 2
/zenrockd tx identity add-keyring-party keyring1k6vc6vhp6e6l3rxalue9v4ux zen1lm70sewfh30xsa3glpj33k2phvzdur43uztzqm --from validator --fees 500000urock --yes
sleep 2

# Create workspace (workspace10j06zdk5gyl6v9ekzwem0v)
/zenrockd tx identity new-workspace --from user-tests --fees 500000urock --yes

# Wait for zenrock process
wait $zenrock_pid


# new key keyring TSS v2
# /zenrockd tx treasury new-key-request workspace10j06zdk5gyl6v9ekzwem0v keyring1k6vc6vhp6e6l3rxalue9v4ux ecdsa --from user-tests --fees 500000urock --yes

# sign message (hello world)
# /zenrockd tx treasury new-signature-request 1 b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9 --from user-tests --fees 500000urock --yes
