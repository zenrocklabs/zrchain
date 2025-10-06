# zenBTC and zenTP Flows

This document outlines the sequence of operations for the zenBTC and zenTP protocols within the zrChain ecosystem, illustrated with Mermaid sequence diagrams.

## Overview of Consensus Mechanism

The zrChain network uses a **Vote Extension** based consensus mechanism where validators run sidecar processes that monitor external blockchains (Bitcoin, Solana, Ethereum) and report their state. Each validator submits their observed data as a vote extension, and only data that reaches **supermajority consensus** (>67% of voting power) is accepted and processed on-chain (except for less critical fields that only require a simple majority).

### Vote Extension Lifecycle

1. **ExtendVoteHandler**: Each validator's sidecar collects oracle data and creates vote extensions containing hashes
2. **VerifyVoteExtensionHandler**: Validators verify each other's vote extensions for basic validity
3. **PrepareProposal**: The proposer aggregates vote extensions and determines consensus fields
4. **ProcessProposal**: All validators verify the proposed oracle data matches consensus
5. **PreBlocker**: Oracle data with consensus is applied to on-chain state and triggers transaction processing

### Vote Extension Process

1. **Sidecar Data Collection**: Each validator's sidecar continuously monitors external chains and collects oracle data
2. **Vote Extension Creation**: During `ExtendVoteHandler`, validators query their sidecar and create vote extensions containing hashes of the observed data
3. **Consensus Verification**: In `PrepareProposal`, the proposer queries sidecar state to validate it against vote extensions from the previous block, then determines which fields have supermajority consensus
4. **State Application**: In `PreBlocker`, only fields with consensus are processed and applied to on-chain state

This ensures that external blockchain state is only acted upon when there is strong validator agreement, providing security against oracle manipulation.

### Consensus Details

**Supermajority Threshold**: Most fields require >2/3 of total voting power to reach consensus
**Simple Majority Fields**: Gas-related fields (ETH gas prices, Solana fees) use >1/2 threshold for faster updates
**Deterministic Tie-Breaking**: When multiple values have equal vote power, lexicographic ordering ensures all validators select the same result
**Field-Level Consensus**: Each data field (prices, nonces, events) reaches consensus independently, allowing partial state updates

This granular consensus approach maximizes system uptime by allowing critical operations to proceed even when some oracle data is unavailable.

## Key Components

- **Sidecar**: Synchronised oracle system, polled by zrChain validators and enshrined by ROCK stake
- **Vote Extensions**: CometBFT mechanism to extend consensus over arbitrary non-transaction data
- **MPC Stack**: Monitors zrChain for cryptographic requests and generates keys/signatures using threshold cryptography (GG21 MPC algo)
- **Relayer**: Service that broadcasts signed transactions to external blockchains
- **Bitcoin Proxy**: Specialized trustless service for Bitcoin transaction monitoring and construction

## zenBTC Protocol

zenBTC allows for the trust-minimized bridging of Bitcoin to and from other blockchains like Solana and Ethereum.

### Deposit and Mint

This flow shows the complete process from depositing BTC to receiving zenBTC tokens on the destination chain.
Deposited BTC is custodied and tracked on zrChain and zenBTC is minted directly on the destination chain. The system no longer routes through EigenLayer staking for deposits; instead, zenBTC supply, custodied BTC, and validator bedrock accounting are updated upon mint confirmation.

```mermaid
sequenceDiagram
    participant User
    participant Web Frontend
    participant zrChain
    participant MPC Stack
    participant Bitcoin
    participant Bitcoin Proxy
    participant Sidecar
    participant Relayer
    participant Ethereum
    participant Solana

    User->>Web Frontend: Request BTC Deposit Address
    Web Frontend->>zrChain: Request new deposit address
    Note over zrChain: Create new Bitcoin key request
    MPC Stack->>zrChain: Poll for key requests
    zrChain-->>MPC Stack: Bitcoin key request found
    MPC Stack->>MPC Stack: Generate new Bitcoin key
    MPC Stack->>zrChain: Submit key request fulfillment transaction
    zrChain-->>Web Frontend: Return deposit address
    Web Frontend-->>User: Display deposit address

    User->>Bitcoin: Deposit BTC to provided address

    Note over Bitcoin: Block mined (deposit tx included)

    Sidecar->>Bitcoin: Polls for new block headers
    zrChain->>Sidecar: ExtendVoteHandler: Validators query sidecars to populate vote extensions with BTC headers/prices/fees
    zrChain->>Sidecar: PrepareProposal: Proposer validates sidecar state against vote extensions
    Note over zrChain: Vote Extensions reach supermajority consensus on external chain data

    Bitcoin Proxy->>Bitcoin: Detects deposit
    Bitcoin Proxy->>Bitcoin Proxy: Generate Merkle proof of BTC deposit transaction
    Bitcoin Proxy->>zrChain: MsgVerifyDepositBlockInclusion(proof)

    zrChain->>zrChain: Verify Merkle proof against Bitcoin block headers
    zrChain->>zrChain: System monitors for Bitcoin reorgs by requesting 6 previous headers
    zrChain->>zrChain: Compare + store historical Bitcoin headers for reorg protection
    zrChain->>zrChain: Validate transaction outputs and amounts
    zrChain->>zrChain: Create PendingMintTransaction (status: DEPOSITED)
    alt Destination = Solana
        zrChain->>zrChain: Request Solana Nonce & mark recipient account requested
    else Destination = EVM
        zrChain->>zrChain: Request Ethereum Minter Nonce
    end

    Sidecar->>Ethereum: Polls ETH/BTC prices and gas data
    Sidecar->>Solana: Polls for fees, nonce and account data
    zrChain->>Sidecar: ExtendVoteHandler: Validators populate vote extensions with nonces/accounts/prices
    zrChain->>Sidecar: PrepareProposal: Proposer validates nonces/accounts/prices against vote extensions
    Note over zrChain: Vote Extensions reach supermajority consensus on required fields

    alt Mint on Solana
        Note over zrChain: PreBlocker: processZenBTCMintsSolana()
        zrChain->>zrChain: Calculate zenBTC mint fee (flat) from BTC price & exchange rate
        zrChain->>zrChain: Determine mint amount from deposited BTC using current exchange rate
        zrChain->>zrChain: Create SignTransactionRequest for Solana mint
        MPC Stack->>zrChain: Poll for signature requests
        zrChain-->>MPC Stack: Solana mint transaction request found
        MPC Stack->>MPC Stack: Generate signature(s)
        MPC Stack->>zrChain: Submit signature request fulfillment transaction
        Relayer->>zrChain: Poll for fulfilled requests
        zrChain-->>Relayer: Signed Mint Tx picked up
        Relayer->>Solana: Broadcast Mint Tx
        Note over zrChain: Timeout management with BTL & AwaitingEventSince

        Sidecar->>Solana: Scans for zenBTC mint events
        zrChain->>Sidecar: ExtendVoteHandler: Validators populate vote extensions with mint events
        zrChain->>Sidecar: PrepareProposal: Proposer validates mint events against vote extensions
        Note over zrChain: Vote Extensions reach supermajority consensus on mint events
        Note over zrChain: PreBlocker: processSolanaZenBTCMintEvents()
        zrChain->>zrChain: Match event to PendingMintTransaction
        zrChain->>zrChain: Update PendingMintTransaction (status: MINTED)
        zrChain->>zrChain: Update supply (PendingZenBTC → MintedZenBTC) and adjust validator bedrock BTC
        zrChain-->>User: zenBTC minted on Solana
    else Mint on Ethereum
        Note over zrChain: PreBlocker: processZenBTCMintsEthereum()
        zrChain->>zrChain: Calculate zenBTC mint fee (flat) from BTC price & exchange rate
        zrChain->>zrChain: Determine mint amount from deposited BTC using current exchange rate
        zrChain->>zrChain: Create SignTransactionRequest for Ethereum mint
        MPC Stack->>zrChain: Poll for signature requests
        zrChain-->>MPC Stack: Ethereum mint transaction request found
        MPC Stack->>MPC Stack: Generate signature
        MPC Stack->>zrChain: Submit signature request fulfillment transaction
        Relayer->>zrChain: Poll for fulfilled requests
        zrChain-->>Relayer: Signed Mint Tx picked up
        Relayer->>Ethereum: Broadcast Mint Tx

        Sidecar->>Ethereum: Polls for nonce update after tx broadcast
        zrChain->>Sidecar: ExtendVoteHandler: Validators populate vote extensions with updated nonces
        Note over zrChain: Vote Extensions reach supermajority consensus on updated nonce
        zrChain->>zrChain: Update status to MINTED
        zrChain->>zrChain: Update supply (PendingZenBTC → MintedZenBTC) and adjust validator bedrock BTC
        zrChain-->>User: zenBTC minted on Ethereum
    end
```

### Redemption and Burn (Direct mode without EigenLayer)

This flow describes how a user burns zenBTC on a destination chain to redeem their original BTC.

```mermaid
sequenceDiagram
    participant User
    participant DestinationChain as Ethereum / Solana
    participant Sidecar
    participant zrChain
    participant MPC Stack
    participant Relayer
    participant Ethereum
    participant Bitcoin Proxy
    participant Bitcoin

    User->>DestinationChain: Burn zenBTC
    Sidecar->>DestinationChain: Scans for Burn Events
    zrChain->>Sidecar: ExtendVoteHandler: Validators query sidecars to populate vote extensions with zenBTC burn events
    zrChain->>Sidecar: PrepareProposal: Proposer validates burn events against vote extensions
    Note over zrChain: Vote Extensions reach supermajority consensus on burn events

    Note over zrChain: PreBlocker: storeNewZenBTCBurnEvents()
    zrChain->>zrChain: Check burn event not already processed (prevent duplicates)
    zrChain->>zrChain: Create BurnEvent (status: BURNED)
    zrChain->>zrChain: Create Redemption object (status: UNSTAKED)

    Bitcoin Proxy->>zrChain: Poll for UNSTAKED redemptions
    zrChain-->>Bitcoin Proxy: Redemption Info (amount, address)
    Bitcoin Proxy->>Bitcoin: Query UTXOs for available funds
    Bitcoin Proxy->>Bitcoin Proxy: Construct unsigned Bitcoin redemption transaction
    Bitcoin Proxy->>zrChain: MsgSubmitUnsignedRedemptionTx(unsigned_tx)
    zrChain->>zrChain: Parse and verify unsigned BTC transaction outputs
    zrChain->>zrChain: Validate invariants (minted zenBTC ≥ redemption amount)
    zrChain->>zrChain: Calculate BTC redemption amount: Convert burned zenBTC to BTC using current exchange rate
    zrChain->>zrChain: Flag redemptions as processed to prevent double-spending
    zrChain->>zrChain: Update redemption status to AWAITING_SIGN
    zrChain->>zrChain: Create SignTransactionRequest for BTC redemption
    MPC Stack->>zrChain: Poll for signature requests
    zrChain-->>MPC Stack: BTC redemption transaction request found
    MPC Stack->>MPC Stack: Generate signature
    MPC Stack->>zrChain: Submit signature request fulfillment transaction

    Note over zrChain: checkForRedemptionFulfilment() monitors AWAITING_SIGN redemptions
    Note over zrChain: When MPC signature is fulfilled, update status to COMPLETED

    Bitcoin Proxy->>zrChain: Poll for fulfilled BTC tx
    zrChain-->>Bitcoin Proxy: Signed BTC Transaction
    Bitcoin Proxy->>Bitcoin: Broadcast signed tx
    Bitcoin-->>User: Receives redeemed BTC
    zrChain->>zrChain: Adjust bedrock BTC accounting
    zrChain->>zrChain: Mark redemption as COMPLETED
    Note over zrChain: When MPC signature is fulfilled, update status to COMPLETED and adjust bedrock BTC accounting
```

## zenTP (Zenrock Transport Protocol)

zenTP is used for bridging native zrChain assets (currently only ROCK) to other blockchains (currently only Solana).
The protocol is a stripped-back iteration of zenBTC's bridging system, omitting restaking features amongst others.

### Bridge to Solana (Mint solROCK)

This flow describes bridging a native asset from zrChain to Solana, resulting in the minting of a corresponding SPL token (e.g., solROCK).

```mermaid
sequenceDiagram
    participant User
    participant zrChain
    participant MPC Stack
    participant Relayer
    participant Sidecar
    participant Solana

    User->>zrChain: MsgBridge(amount, solana_addr)
    zrChain->>zrChain: Validate amount against 1bn supply cap
    zrChain->>zrChain: Calculate bridge fee as percentage of amount and add to total cost
    zrChain->>zrChain: Lock User's native tokens in module
    zrChain->>zrChain: Create Bridge object (status: PENDING)
    zrChain->>zrChain: Request Solana Nonce & Account Info

    Sidecar->>Solana: Polls for nonce and account data
    zrChain->>Sidecar: ExtendVoteHandler: Validators query sidecars to populate vote extensions with SOL nonces/accounts
    zrChain->>Sidecar: PrepareProposal: Proposer validates Solana data against vote extensions
    Note over zrChain: Vote Extensions reach supermajority consensus on Solana data
    Note over zrChain: PreBlocker: processSolanaROCKMints()
    zrChain->>zrChain: Validates consensus on required fields (nonce, accounts) before transaction
    zrChain->>zrChain: Check if transaction already processed
    zrChain->>zrChain: Create SignTransactionRequest for Solana ROCK mint
    MPC Stack->>zrChain: Poll for signature requests
    zrChain-->>MPC Stack: Solana ROCK mint transaction request found
    MPC Stack->>MPC Stack: Generate signature
    MPC Stack->>zrChain: Submit signature request fulfillment transaction
    Relayer->>zrChain: Poll for fulfilled requests
    zrChain-->>Relayer: Signed Mint Tx picked up
    Relayer->>Solana: Broadcast Mint Tx
    Note over zrChain: Timeout management (when sidecars have consensus): BTL expiry check and AwaitingEventSince tracking
    zrChain->>zrChain: Check transaction BTL (blocks-to-live) expiry for retry timeout management
    zrChain->>zrChain: Set AwaitingEventSince timestamp to track confirmation wait time for secondary timeout logic

    Sidecar->>Solana: Scans for Mint Events
    zrChain->>Sidecar: ExtendVoteHandler: Validators query sidecars to populate vote extensions with solROCK mint events
    zrChain->>Sidecar: PrepareProposal: Proposer validates mint events against vote extensions
    Note over zrChain: Vote Extensions reach supermajority consensus on mint events
    Note over zrChain: PreBlocker: processSolanaROCKMintEvents()
    zrChain->>zrChain: Match event to Bridge request
    zrChain->>zrChain: Verify event not already processed
    zrChain->>zrChain: Burn locked native tokens from zenTP module
    zrChain->>zrChain: Enforces 1bn token supply cap with invariant checks
    zrChain->>zrChain: Update solanaROCKSupply
    zrChain->>zrChain: Update Bridge object (status: COMPLETED)
    User->>Solana: Receives solROCK
```

### Burn from Solana (Redeem native ROCK)

This flow describes burning an SPL token on Solana to redeem the original native asset back on zrChain.

```mermaid
sequenceDiagram
    participant User
    participant zrChain
    participant Sidecar
    participant Solana

    User->>Solana: Burn solROCK (embedding zrChain destination address into event)
    Sidecar->>Solana: Scans for Burn Events from bridge contract
    zrChain->>Sidecar: ExtendVoteHandler: Validators query sidecars to populate vote extensions with solROCK burn events
    zrChain->>Sidecar: PrepareProposal: Proposer validates burn events against vote extensions
    Note over zrChain: Vote Extensions reach supermajority consensus on burn events

    Note over zrChain: PreBlocker: processSolanaROCKBurnEvents()
    zrChain->>zrChain: Check burn not already processed (primary key = TxID + ChainID)
    zrChain->>zrChain: Check sufficient Solana ROCK supply exists
    zrChain->>zrChain: Calculate bridge fee as percentage of burned amount
    zrChain->>zrChain: Mint total burn amount to zentp module
    zrChain->>zrChain: Deduct from Solana ROCK supply tracking
    zrChain->>zrChain: Send (burn_amount - fee) to user's address
    zrChain->>zrChain: Retain fee in module account
    zrChain->>zrChain: Create Bridge record (status: COMPLETED)
    zrChain->>zrChain: Clear any corresponding backfill requests for processed events
    zrChain->>User: Native ROCK tokens received on zrChain
```
