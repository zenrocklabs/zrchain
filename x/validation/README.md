# Zenrock Hybrid Validation Module - `x/validation`

This is a fork of the Cosmos SDK's `x/staking` module implementing our Hybrid Validation (HV) logic. This version differs from the original in a few key areas:

### Oracle
- `abci.go`, `abci_helpers.go`, and `abci_types.go` contain oracle-related logic - [the overall architecture is outlined in this doc on Notion.](https://www.notion.so/VE-Based-Oracle-System-f86bca07ce27425ca2f82817f345a1bb?pvs=4)

### Consensus Voting Power
- `/keeper/val_state_change.go` implements the logic to enable using off-chain assets as collateral for staking, ensuring correct voting power weighting by multiplying the staking token amounts by their prices.

### Staking Rewards
- The above `val_state_change.go` file also contains the current implementation for provisioning staking rewards for these assets (currently only ETH is supported), storing them in the `AVSRewardsPool` keeper store.
- The mechanism for withdrawing these rewards to Ethereum as ETH / ERC-20 tokens still needs to be designed & implemented. EigenLayer's implementation of this is also not fully finalised yet.

### Slashing
- `/keeper/slash.go` implements slashing of tokens delegated to validators through the AVS contract. In a future version this will be propagated to the contract as the tokens are stored there.

### Misc
- Custom validator types as well as other custom types are used in place of vanilla `x/staking` SDK ones across many files
- i.e. `ValidatorHV` also contains the `AVSTokens` field for the amount of restaked ETH delegated to the validator.