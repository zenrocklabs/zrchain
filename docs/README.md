# zenBTC and zenTP Flows

This document outlines the sequence of operations for the zenBTC and zenTP protocols within the zrchain ecosystem, illustrated with Mermaid sequence diagrams.

## zenBTC Protocol

zenBTC allows for the trust-minimized bridging of Bitcoin to and from other blockchains like Ethereum and Solana.

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

    Sidecar->>Bitcoin: Polls for new blocks
    Sidecar-->>zrchain: Report BTC Block Header (via vote extension)
    
    zrchain->>zrchain: Verify proof, Create PendingMintTransaction (status: DEPOSITED)
    zrchain->>zrchain: Request Staker Nonce for EigenLayer

    Note over zrchain,EigenLayer: Consensus: Staking on EigenLayer
    zrchain->>zrchain: PreBlocker: processZenBTCStaking()
    zrchain->>MPC Cluster: constructStakeTx() -> SignTransactionRequest
    MPC Cluster-->>zrchain: Fulfill SignTransactionRequest
    Relayer->>zrchain: Poll for fulfilled requests
    zrchain-->>Relayer: Signed Stake Tx
    Relayer->>EigenLayer: Broadcast Stake Tx

    Sidecar->>EigenLayer: Polls for nonce update after tx broadcast
    Sidecar-->>zrchain: Reports new nonce (via vote extension)
    Note over zrchain: Data is verified via Vote Extension consensus
    zrchain->>zrchain: PreBlocker confirms tx, updates status to STAKED
    zrchain->>zrchain: Request Minter Nonce (ETH or SOL)

    alt Mint on Solana
        Note over zrchain,Solana: Consensus: Minting zenBTC on Solana
        zrchain->>zrchain: PreBlocker: processZenBTCMintsSolana()
        zrchain->>MPC Cluster: PrepareSolanaMintTx() -> SignTransactionRequest
        MPC Cluster-->>zrchain: Fulfill SignTransactionRequest
        Relayer->>zrchain: Poll for fulfilled requests
        zrchain-->>Relayer: Signed Mint Tx
        Relayer->>Solana: Broadcast Mint Tx
        
        Sidecar->>Solana: Scans for Mint Events
        Sidecar-->>zrchain: Reports new Mint Events (via vote extension)
        Note over zrchain: Data is verified via Vote Extension consensus
        zrchain->>zrchain: PreBlocker: processSolanaZenBTCMintEvents()
        zrchain->>zrchain: Match event, Update PendingMintTransaction (status: MINTED)
        zrchain-->>User: zenBTC minted on Solana
    else Mint on Ethereum
        Note over zrchain,Ethereum: Consensus: Minting zenBTC on Ethereum
        zrchain->>zrchain: PreBlocker: processZenBTCMintsEthereum()
        zrchain->>MPC Cluster: constructMintTx() -> SignTransactionRequest
        MPC Cluster-->>zrchain: Fulfill SignTransactionRequest
        Relayer->>zrchain: Poll for fulfilled requests
        zrchain-->>Relayer: Signed Mint Tx
        Relayer->>Ethereum: Broadcast Mint Tx

        Sidecar->>Ethereum: Polls for nonce update after tx broadcast
        Sidecar-->>zrchain: Reports new nonce (via vote extension)
        Note over zrchain: Data is verified via Vote Extension consensus
        zrchain->>zrchain: PreBlocker confirms tx, updates status to MINTED
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
    Note over zrchain: Data is verified via Vote Extension consensus

    zrchain->>zrchain: PreBlocker: storeNewZenBTCBurnEvents()
    zrchain->>zrchain: Create BurnEvent (status: BURNED)
    zrchain->>zrchain: Request Unstaker Nonce for EigenLayer

    Note over zrchain,EigenLayer: Consensus: Unstaking from EigenLayer
    zrchain->>zrchain: PreBlocker: processZenBTCBurnEvents()
    zrchain->>MPC Cluster: constructUnstakeTx() -> SignTransactionRequest
    MPC Cluster-->>zrchain: Fulfill SignTransactionRequest
    Relayer->>zrchain: Poll for fulfilled requests
    zrchain-->>Relayer: Signed Unstake Tx
    Relayer->>EigenLayer: Broadcast Unstake Tx
    
    Sidecar->>EigenLayer: Polls for nonce update after tx broadcast
    Sidecar-->>zrchain: Reports new nonce (via vote extension)
    Note over zrchain: Data is verified via Vote Extension consensus
    zrchain->>zrchain: PreBlocker confirms tx, updates status to UNSTAKING

    Note over Sidecar,EigenLayer: Sidecar polls EigenLayer for unstake completion
    Sidecar->>EigenLayer: Polls for unstake completion
    Sidecar-->>zrchain: Reports unstake ready (via vote extension)
    Note over zrchain: Data is verified via Vote Extension consensus
    zrchain->>zrchain: PreBlocker: storeNewZenBTCRedemptions()
    zrchain->>zrchain: Update Redemption (status: UNSTAKED)

    Note over zrchain,EigenLayer: Consensus: Completing withdrawal from EigenLayer
    zrchain->>zrchain: PreBlocker: processZenBTCRedemptions()
    zrchain->>MPC Cluster: constructCompleteWithdrawalTx() -> SignTransactionRequest
    MPC Cluster-->>zrchain: Fulfill SignTransactionRequest
    Relayer->>zrchain: Poll for fulfilled requests
    zrchain-->>Relayer: Signed CompleteWithdrawal Tx
    Relayer->>EigenLayer: Broadcast Tx
    
    Sidecar->>EigenLayer: Polls for nonce update after tx broadcast
    Sidecar-->>zrchain: Reports new nonce (via vote extension)
    Note over zrchain: Data is verified via Vote Extension consensus
    zrchain->>zrchain: PreBlocker confirms tx, updates status to READY_FOR_BTC_RELEASE

    Note over Bitcoin Proxy, zrchain: Proxy polls for redemptions
    Bitcoin Proxy->>zrchain: Poll for READY_FOR_BTC_RELEASE redemptions
    zrchain-->>Bitcoin Proxy: Redemption Info (UTXOs)
    Bitcoin Proxy->>zrchain: MsgSubmitUnsignedRedemptionTx(UTXOs)
    zrchain->>MPC Cluster: Request signature for BTC tx
    MPC Cluster-->>zrchain: Fulfill SignTransactionRequest
    Bitcoin Proxy->>zrchain: Poll for fulfilled BTC tx
    zrchain-->>Bitcoin Proxy: Signed BTC Transaction
    Bitcoin Proxy->>Bitcoin: Broadcast signed tx
    Bitcoin-->>User: Receives redeemed BTC
    
    Note over zrchain: (Post-broadcast) Redemption marked as COMPLETED
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

    Note over zrchain,Solana: Consensus for Minting solROCK on Solana
    zrchain->>zrchain: PreBlocker: processSolanaROCKMints()
    zrchain->>MPC Cluster: PrepareSolanaMintTx() -> SignTransactionRequest
    MPC Cluster-->>zrchain: Fulfill SignTransactionRequest
    Relayer->>zrchain: Poll for fulfilled requests
    zrchain-->>Relayer: Signed Mint Tx
    Relayer->>Solana: Broadcast Mint Tx

    Sidecar->>Solana: Scans for Mint Events
    Sidecar-->>zrchain: Reports new Mint Events (via vote extension)
    Note over zrchain: Data is verified via Vote Extension consensus
    zrchain->>zrchain: PreBlocker: processSolanaROCKMintEvents()
    zrchain->>zrchain: Match event to Bridge object
    zrchain->>zrchain: Burn locked native tokens
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
    Note over zrchain: Data is verified via Vote Extension consensus

    zrchain->>zrchain: PreBlocker: processSolanaROCKBurnEvents()
    zrchain->>zrchain: Verify burn event is new
    zrchain->>zrchain: Mint native ROCK tokens
    zrchain->>User: Send native ROCK tokens to user's zrchain address
```
""
