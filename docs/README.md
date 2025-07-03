# zenBTC and zenTP Flows

This document outlines the sequence of operations for the zenBTC and zenTP protocols within the zrchain ecosystem, illustrated with Mermaid sequence diagrams.

## Overview of Consensus Mechanism

The zrChain network uses a **Vote Extension** based consensus mechanism where validators run sidecar processes that monitor external blockchains (Bitcoin, Solana, Ethereum) and report their state. Each validator submits their observed data as a vote extension, and only data that reaches **supermajority consensus** (>67% of voting power) is accepted and processed on-chain.

### Vote Extension Lifecycle

1. **ExtendVoteHandler**: Each validator's sidecar collects oracle data and creates vote extensions containing hashes
2. **VerifyVoteExtensionHandler**: Validators verify each other's vote extensions for basic validity  
3. **PrepareProposal**: The proposer aggregates vote extensions and determines consensus fields
4. **ProcessProposal**: All validators verify the proposed oracle data matches consensus
5. **PreBlocker**: Oracle data with consensus is applied to on-chain state and triggers transaction processing

### Vote Extension Process

1. **Sidecar Data Collection**: Each validator's sidecar monitors external chains and collects oracle data
2. **Vote Extension Creation**: During `ExtendVoteHandler`, validators create vote extensions containing hashes of their observed data
3. **Consensus Verification**: In `PrepareProposal`/`ProcessProposal`, the network determines which fields have supermajority consensus
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
- **Vote Extensions**: CometBFT mechanism to extend consensus over arbitrary non-tx data
- **MPC Cluster**: Multi-party computation system for generating cryptographic signatures
- **Relayer**: Service that broadcasts signed transactions to external blockchains
- **Bitcoin Proxy**: Specialized service for Bitcoin transaction monitoring and construction

## zenBTC Protocol

zenBTC allows for the trust-minimized bridging of Bitcoin to and from other blockchains like Solana and Ethereum.

### Deposit and Mint

This flow describes how a user deposits BTC and how it is relayed to mint zenBTC on a destination chain.

```mermaid
sequenceDiagram
    participant User
    participant Web Frontend
    participant zrchain
    participant MPC Cluster
    participant Bitcoin
    participant Bitcoin Proxy
    participant Sidecar
    participant Relayer
    participant EigenLayer
    participant Solana
    participant Ethereum

    User->>Web Frontend: Request BTC Deposit Address
    Web Frontend->>zrchain: Request new deposit address
    zrchain->>MPC Cluster: Request new key provisioning
    MPC Cluster-->>zrchain: New Bitcoin Address
    zrchain-->>Web Frontend: Return deposit address
    Web Frontend-->>User: Display deposit address

    User->>Bitcoin: Deposit BTC to provided address
    
    Bitcoin Proxy->>Bitcoin: Detects deposit
    Bitcoin Proxy->>zrchain: MsgVerifyDepositBlockInclusion(proof)

    Sidecar->>Bitcoin: Polls for new block headers
    Sidecar->>Ethereum: Polls for ETH/BTC price feeds (Chainlink)
    Sidecar->>Ethereum: Polls for gas prices and estimates
    Sidecar->>Solana: Polls for lamports per signature fee
    Sidecar-->>zrchain: Report BTC Block Header, prices & network fees (via vote extension)
    Note over zrchain: Vote Extensions reach supermajority consensus on external chain data

    zrchain->>zrchain: Verify proof, Create PendingMintTransaction (status: DEPOSITED)
    zrchain->>zrchain: Request Staker Nonce for EigenLayer

    Sidecar->>EigenLayer: Polls for nonce values
    Sidecar-->>zrchain: Report nonce data via Vote Extension
    Note over zrchain: Vote Extensions reach supermajority consensus on nonce data
    Note over zrchain: Validates consensus on required fields (nonce, gas, prices) before transaction
    zrchain->>zrchain: PreBlocker: processZenBTCStaking()
    zrchain->>MPC Cluster: constructStakeTx() -> SignTransactionRequest
    MPC Cluster-->>zrchain: Fulfill SignTransactionRequest
    Relayer->>zrchain: Poll for fulfilled requests
    zrchain-->>Relayer: Signed Stake Tx
    Relayer->>EigenLayer: Broadcast Stake Tx

    Sidecar->>EigenLayer: Polls for nonce update after tx broadcast
    Sidecar-->>zrchain: Reports updated nonce (via vote extension)
    Note over zrchain: Vote Extensions reach supermajority consensus on updated nonce
    zrchain->>zrchain: PreBlocker confirms tx, updates status to STAKED
    zrchain->>zrchain: Request Minter Nonce (ETH or SOL)

    alt Mint on Solana
        Sidecar->>Solana: Polls for nonce and account data
        Sidecar-->>zrchain: Report Solana nonce/account data via Vote Extension
        Note over zrchain: Vote Extensions reach supermajority consensus on Solana data
        zrchain->>zrchain: PreBlocker: processZenBTCMintsSolana()
        zrchain->>MPC Cluster: PrepareSolanaMintTx() -> SignTransactionRequest
        MPC Cluster-->>zrchain: Fulfill SignTransactionRequest
        Relayer->>zrchain: Poll for fulfilled requests
        zrchain-->>Relayer: Signed Mint Tx
        Relayer->>Solana: Broadcast Mint Tx
        Note over zrchain: Solana transactions have BTL (Blocks To Live) timeout with retry logic
        
        Sidecar->>Solana: Scans for Mint Events
        Sidecar-->>zrchain: Reports new Mint Events (via vote extension)
        Note over zrchain: Vote Extensions reach supermajority consensus on mint events
        zrchain->>zrchain: PreBlocker: processSolanaZenBTCMintEvents()
        zrchain->>zrchain: Match event, Update PendingMintTransaction (status: MINTED)
        Note over zrchain: Updates zenBTC supply tracking (PendingZenBTC → MintedZenBTC)
        zrchain-->>User: zenBTC minted on Solana
    else Mint on Ethereum
        Sidecar->>Ethereum: Polls for nonce values 
        Sidecar-->>zrchain: Report Ethereum nonce data via Vote Extension
        Note over zrchain: Vote Extensions reach supermajority consensus on nonce data
        zrchain->>zrchain: PreBlocker: processZenBTCMintsEthereum()
        zrchain->>MPC Cluster: constructMintTx() -> SignTransactionRequest
        MPC Cluster-->>zrchain: Fulfill SignTransactionRequest
        Relayer->>zrchain: Poll for fulfilled requests
        zrchain-->>Relayer: Signed Mint Tx
        Relayer->>Ethereum: Broadcast Mint Tx

        Sidecar->>Ethereum: Polls for nonce update after tx broadcast
        Sidecar-->>zrchain: Reports updated nonce (via vote extension)
        Note over zrchain: Vote Extensions reach supermajority consensus on updated nonce
        zrchain->>zrchain: PreBlocker confirms tx, updates status to MINTED
        Note over zrchain: Updates zenBTC supply tracking (PendingZenBTC → MintedZenBTC)
        zrchain-->>User: zenBTC minted on Ethereum
    end
```

### Redemption and Burn

This flow describes how a user burns zenBTC on a destination chain to redeem their original BTC.

```mermaid
sequenceDiagram
    participant User
    participant DestinationChain as Ethereum / Solana
    participant Sidecar
    participant zrchain
    participant MPC Cluster
    participant Relayer
    participant EigenLayer
    participant Bitcoin Proxy
    participant Bitcoin

    User->>DestinationChain: Burn zenBTC
    Sidecar->>DestinationChain: Scans for Burn Events
    Sidecar-->>zrchain: Reports new Burn Events (via vote extension)
    Note over zrchain: Vote Extensions reach supermajority consensus on burn events

    zrchain->>zrchain: PreBlocker: storeNewZenBTCBurnEvents()
    zrchain->>zrchain: Create BurnEvent (status: BURNED)
    zrchain->>zrchain: Request Unstaker Nonce for EigenLayer

    Sidecar->>EigenLayer: Polls for nonce values
    Sidecar-->>zrchain: Report nonce data via Vote Extension
    Note over zrchain: Vote Extensions reach supermajority consensus on nonce data
    zrchain->>zrchain: PreBlocker: processZenBTCBurnEvents()
    zrchain->>MPC Cluster: constructUnstakeTx() -> SignTransactionRequest
    MPC Cluster-->>zrchain: Fulfill SignTransactionRequest
    Relayer->>zrchain: Poll for fulfilled requests
    zrchain-->>Relayer: Signed Unstake Tx
    Relayer->>EigenLayer: Broadcast Unstake Tx
    
    Sidecar->>EigenLayer: Polls for nonce update after tx broadcast
    Sidecar-->>zrchain: Reports updated nonce (via vote extension)
    Note over zrchain: Vote Extensions reach supermajority consensus on updated nonce
    zrchain->>zrchain: PreBlocker confirms tx, updates status to UNSTAKING

    Sidecar->>EigenLayer: Polls for unstake completion (redemption availability)
    Sidecar-->>zrchain: Reports redemption data when ready (via vote extension)
    Note over zrchain: Vote Extensions reach supermajority consensus on redemption data
    zrchain->>zrchain: PreBlocker: storeNewZenBTCRedemptions()
    zrchain->>zrchain: Update Redemption (status: UNSTAKED)
    zrchain->>zrchain: Wait for EigenLayer withdrawal delay period

    Note over zrchain: After withdrawal delay, redemption becomes available for completion

    Sidecar->>EigenLayer: Polls for completer nonce values
    Sidecar-->>zrchain: Report completer nonce data via Vote Extension
    Note over zrchain: Vote Extensions reach supermajority consensus on completer nonce
    Note over zrchain: Validates consensus on required fields (nonce, gas, prices) before transaction
    zrchain->>zrchain: PreBlocker: processZenBTCRedemptions()
    zrchain->>MPC Cluster: constructCompleteWithdrawalTx() -> SignTransactionRequest
    MPC Cluster-->>zrchain: Fulfill SignTransactionRequest
    Relayer->>zrchain: Poll for fulfilled requests
    zrchain-->>Relayer: Signed CompleteWithdrawal Tx
    Relayer->>EigenLayer: Broadcast Tx
    
    Sidecar->>EigenLayer: Polls for nonce update after tx broadcast
    Sidecar-->>zrchain: Reports updated nonce (via vote extension)
    Note over zrchain: Vote Extensions reach supermajority consensus on updated nonce
    zrchain->>zrchain: PreBlocker confirms tx, updates status to READY_FOR_BTC_RELEASE

    Bitcoin Proxy->>zrchain: Poll for READY_FOR_BTC_RELEASE redemptions
    zrchain-->>Bitcoin Proxy: Redemption Info (UTXOs)
    Bitcoin Proxy->>zrchain: MsgSubmitUnsignedRedemptionTx(UTXOs)
    zrchain->>MPC Cluster: Request signature for BTC tx
    Note over zrchain: Performs invariant checks (sufficient minted zenBTC and custodied BTC)
    Note over zrchain: Calculates BTC amount using current exchange rate
    MPC Cluster-->>zrchain: Fulfill SignTransactionRequest
    Bitcoin Proxy->>zrchain: Poll for fulfilled BTC tx
    zrchain-->>Bitcoin Proxy: Signed BTC Transaction
    Bitcoin Proxy->>Bitcoin: Broadcast signed tx
    Bitcoin-->>User: Receives redeemed BTC
    
    Sidecar->>Bitcoin: Monitors for transaction confirmation
    Sidecar-->>zrchain: Reports transaction inclusion (via vote extension)
    Note over zrchain: Vote Extensions reach consensus on Bitcoin transaction confirmation
    Note over zrchain: System monitors for Bitcoin reorgs by requesting 20 previous headers
    zrchain->>zrchain: Mark redemption as COMPLETED
```

## zenTP Protocol

zenTP (Zenrock Transport Protocol) is used for bridging native zrchain assets to other blockchains, such as Solana.

### Bridge to Solana (Mint solROCK)

This flow describes bridging a native asset from zrchain to Solana, resulting in the minting of a corresponding SPL token (e.g., solROCK).

```mermaid
sequenceDiagram
    participant User
    participant zrchain
    participant MPC Cluster
    participant Relayer
    participant Sidecar
    participant Solana

    User->>zrchain: MsgBridge(amount, solana_addr)
    zrchain->>zrchain: Lock User's native tokens
    zrchain->>zrchain: Create Bridge object (status: PENDING)
    zrchain->>zrchain: Request Solana Nonce & Account Info

    Sidecar->>Solana: Polls for nonce and account data
    Sidecar-->>zrchain: Report Solana nonce/account data via Vote Extension
    Note over zrchain: Vote Extensions reach supermajority consensus on Solana data
    zrchain->>zrchain: PreBlocker: processSolanaROCKMints()
    zrchain->>MPC Cluster: PrepareSolanaMintTx() -> SignTransactionRequest
    MPC Cluster-->>zrchain: Fulfill SignTransactionRequest
    Relayer->>zrchain: Poll for fulfilled requests
    zrchain-->>Relayer: Signed Mint Tx
    Relayer->>Solana: Broadcast Mint Tx
    Note over zrchain: Solana transactions have BTL (Blocks To Live) timeout with retry logic

    Sidecar->>Solana: Scans for Mint Events
    Sidecar-->>zrchain: Reports new Mint Events (via vote extension)
    Note over zrchain: Vote Extensions reach supermajority consensus on mint events
    zrchain->>zrchain: PreBlocker: processSolanaROCKMintEvents()
    zrchain->>zrchain: Match event to Bridge object
    zrchain->>zrchain: Burn locked native tokens
    Note over zrchain: Enforces 1bn token supply cap with invariant checks
    zrchain->>zrchain: Update solanaROCKSupply
    zrchain->>zrchain: Update Bridge object (status: COMPLETED)
    User->>Solana: Receives solROCK
```

### Burn from Solana (Redeem native ROCK)

This flow describes burning an SPL token on Solana to redeem the original native asset back on zrchain.

```mermaid
sequenceDiagram
    participant User
    participant zrchain
    participant Sidecar
    participant Solana

    User->>Solana: Burn solROCK (providing zrchain address)
    Sidecar->>Solana: Scans for Burn Events
    Sidecar-->>zrchain: Reports new Burn Events (via vote extension)
    Note over zrchain: Vote Extensions reach supermajority consensus on burn events

    zrchain->>zrchain: PreBlocker: processSolanaROCKBurnEvents()
    zrchain->>zrchain: Verify burn event is new
    zrchain->>zrchain: Mint native ROCK tokens
    zrchain->>User: Send native ROCK tokens to user's zrchain address
```
