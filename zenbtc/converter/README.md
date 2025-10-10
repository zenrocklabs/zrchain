## AVS Rewards Swap

In order to use a single wallet to claim rewards from multiple operators the earner must be set to the address of the wallet that is going to be used in the config.

To get the address that is configured, you can execute the `balance` command.

```
go run cmd/rewards_swap/main.go balance
INF Balances for 0x935BFE680993dC00891B05b5fCD4B9CAFCF54fE7: ...
````

Next you will need the eigenlayer cli in order to add the earner for a specific operator,
you can find the cli in the following repo https://github.com/Layr-Labs/eigenlayer-cli.

once the cli is build, execute the following command:
```
./bin/eigenlayer rewards set-claimer -b -a <address_from_previous_command> -ea <operator_address> -n hoodi -k <path_to_ecdsa_key_file_from_sidecar> -v --eth-rpc-url https://ethereum-hoodi-rpc.publicnode.com 
```

for mainnet, use 
```
-n mainnet
--eth-rpc-url https://ethereum-rpc.publicnode.com
```

## CLI

The cli provides a set of separate tasks to test individual steps and a command to run all steps in sequence.

The full sequence is claim -> unwrap -> swap.

All commands support a broadcast flag that only when set to true, the cli will actually broadcast the tx.  Setting it to false can be used to test the tx.

### claim

claims the rewards of the specified operator into the eth account that is set in config.

claim [earner] [broadcast]

### unwrap

unwraps the specified amount of weth to eth, amount in native format 1e18

unwrap [amount] [broadcast]

### swap

swaps the specified amount of eth to btc using thorchain, amount in thorchain format 1e8

swap [destination] [amount] [broadcast]

### balance

gets the eth and weth balance of the configured eth account.

### full

executes a full cycle:
- claims all rewards for all configured operators into the configured eth account
- unwraps all weth to eth
- swaps all eth on the account, minus a margin for gas, into btc, sending it to the configured btc address

full [broadcast]
