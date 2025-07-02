# ZenBTC and ZenTP Flows

This document outlines the sequence of operations for the ZenBTC and ZenTP protocols within the zrchain ecosystem, illustrated with Mermaid sequence diagrams.

## ZenBTC Protocol

ZenBTC allows for the trust-minimized bridging of Bitcoin to and from other blockchains like Ethereum and Solana.

### Deposit and Mint

This flow describes how a user deposits BTC and mints zenBTC on a destination chain.

```mermaid
sequenceDiagram
    participant User
    participant zrchain as zrchain (zenbtc)
    participant zrchain_val as zrchain (validation)
    participant Sidecar
    participant Bitcoin
    participant EigenLayer
    participant Solana
    participant Ethereum

    User->>Bitcoin: Deposit BTC
    User->>zrchain: MsgVerifyDepositBlockInclusion(proof)

    zrchain->>zrchain_val: Get BTC Block Header (via Sidecar)
    zrchain->>zrchain: Verify proof, Create PendingMintTransaction (status: DEPOSITED)
    zrchain->>zrchain_val: Request Staker Nonce

    Note over zrchain_val,EigenLayer: Consensus: Staking on EigenLayer
    zrchain_val->>zrchain_val: PreBlocker: processZenBTCStaking()
    zrchain_val->>EigenLayer: constructStakeTx & submitEthereumTransaction
    EigenLayer-->>Sidecar: Stake Event
    Sidecar-->>zrchain_val: Oracle data (nonce update)

    zrchain_val->>zrchain: txContinuationCallback: Update PendingMintTransaction (status: STAKED)
    zrchain_val->>zrchain_val: Request Minter Nonce (ETH or SOL)

    alt Mint on Solana
        Note over zrchain_val,Solana: Consensus: Minting zenBTC on Solana
        zrchain_val->>zrchain_val: PreBlocker: processZenBTCMintsSolana()
        zrchain_val->>Solana: PrepareSolanaMintTx & submitSolanaTransaction
        Solana-->>Sidecar: Mint Event (TokensMintedWithFee)
        Sidecar-->>zrchain_val: Oracle data (SolanaMintEvents)
        zrchain_val->>zrchain: PreBlocker: processSolanaZenBTCMintEvents()
        zrchain->>zrchain: Match event, Update PendingMintTransaction (status: MINTED)
        zrchain-->>User: zenBTC minted on Solana
    else Mint on Ethereum
        Note over zrchain_val,Ethereum: Consensus: Minting zenBTC on Ethereum
        zrchain_val->>zrchain_val: PreBlocker: processZenBTCMintsEthereum()
        zrchain_val->>Ethereum: constructMintTx & submitEthereumTransaction
        Ethereum-->>Sidecar: Mint Event
        Sidecar-->>zrchain_val: Oracle data (nonce update)
        zrchain_val->>zrchain: txContinuationCallback: Update PendingMintTransaction (status: MINTED)
        zrchain-->>User: zenBTC minted on Ethereum
    end
```

### Redemption and Burn

This flow describes how a user burns zenBTC on a destination chain to redeem their original BTC.

```mermaid
sequenceDiagram
    participant User
    participant DestinationChain as Ethereum / Solana
    participant zrchain as zrchain (zenbtc)
    participant zrchain_val as zrchain (validation)
    participant Sidecar
    participant EigenLayer
    participant Bitcoin

    User->>DestinationChain: Burn zenBTC
    DestinationChain-->>Sidecar: Burn Event
    Sidecar-->>zrchain_val: Oracle data (burn event)

    Note over zrchain_val: Consensus on Burn Event
    zrchain_val->>zrchain: PreBlocker: storeNewZenBTCBurnEvents()
    zrchain->>zrchain: Create BurnEvent (status: BURNED)
    zrchain->>zrchain_val: Request Unstaker Nonce

    Note over zrchain_val,EigenLayer: Consensus: Unstaking from EigenLayer
    zrchain_val->>zrchain_val: PreBlocker: processZenBTCBurnEvents()
    zrchain_val->>EigenLayer: constructUnstakeTx & submitEthereumTransaction
    EigenLayer-->>Sidecar: Unstake Initiated Event
    Sidecar-->>zrchain_val: Oracle data (nonce update)
    zrchain_val->>zrchain: txContinuationCallback: Update BurnEvent (status: UNSTAKING)

    Note over Sidecar,EigenLayer: Sidecar monitors for unstake completion
    EigenLayer-->>Sidecar: Unstake Ready Event
    Sidecar-->>zrchain_val: Oracle data (redemption ready)

    zrchain_val->>zrchain: PreBlocker: storeNewZenBTCRedemptions()
    zrchain->>zrchain: Update Redemption (status: READY)
    zrchain_val->>zrchain_val: Request Completer Nonce

    Note over zrchain_val,EigenLayer: Consensus: Completing Unstake
    zrchain_val->>zrchain_val: PreBlocker: processZenBTCRedemptions()
    zrchain_val->>EigenLayer: constructCompleteTx & submitEthereumTransaction
    EigenLayer-->>Sidecar: Unstake Complete Event
    Sidecar-->>zrchain_val: Oracle data (nonce update)
    zrchain_val->>zrchain: txContinuationCallback: Update Redemption (status: COMPLETED)

    User->>zrchain: MsgSubmitUnsignedRedemptionTx(unsignedBtcTx)
    zrchain->>zrchain: Verify unsigned tx
    zrchain->>Bitcoin: Request signature for BTC tx (via Treasury)
    Bitcoin-->>User: Signed BTC Transaction
    User->>Bitcoin: Broadcast signed tx
```

## ZenTP Protocol

ZenTP (Zenrock Transport Protocol) is used for bridging native zrchain assets to other blockchains, such as Solana.

### Bridge to Solana (Mint solROCK)

This flow describes bridging a native asset from zrchain to Solana, resulting in the minting of a corresponding SPL token (e.g., solROCK).

```mermaid
sequenceDiagram
    participant User
    participant zrchain as zrchain (zentp)
    participant zrchain_val as zrchain (validation)
    participant Sidecar
    participant Solana

    User->>zrchain: MsgBridge(amount, solana_addr)
    zrchain->>zrchain: Lock User's native tokens
    zrchain->>zrchain: Create Bridge object (status: PENDING)
    zrchain->>zrchain_val: Request Solana Nonce & Account Info

    Note over zrchain_val,Solana: Consensus for Minting solROCK on Solana
    zrchain_val->>zrchain_val: PreBlocker: processSolanaROCKMints()
    zrchain_val->>Solana: PrepareSolanaMintTx & submitSolanaTransaction

    Solana-->>Sidecar: Mint Event (TokensMintedWithFee)
    Sidecar-->>zrchain_val: Oracle data (mint event)

    zrchain_val->>zrchain: PreBlocker: processSolanaROCKMintEvents()
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
    participant zrchain as zrchain (zentp)
    participant zrchain_val as zrchain (validation)
    participant Sidecar
    participant Solana

    User->>Solana: Burn solROCK (providing zrchain address)
    Solana-->>Sidecar: Burn Event (TokenRedemption)
    Sidecar-->>zrchain_val: Oracle data (burn event)

    Note over zrchain_val: Consensus on Burn Event
    zrchain_val->>zrchain_val: PreBlocker: processSolanaROCKBurnEvents()
    zrchain->>zrchain: Verify burn event is new
    zrchain->>zrchain: Mint native ROCK tokens
    zrchain->>User: Send native ROCK tokens to user's zrchain address
```
