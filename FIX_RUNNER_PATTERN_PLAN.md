# Queue Processor Refactor Plan (zenBTC/zenTP ABCI)

Owner: Validation module (x/validation/keeper)
Status: Design/Implementation plan

## Goals

- Make transaction queue processing explicit, uniform, and easy to read.
- Reduce boilerplate and error-prone conditionals (nonce gates, empty queues, confirmation vs. retry semantics).
- Adopt intuitive names the team agreed on:
  - GetPendingTxs(ctx) ([]T, error)
  - DispatchTx(item T) error
  - OnTxConfirmed(item T) error (EVM)
  - UpdatePendingTxStatus(item T) error (Solana; replaces ambiguous "OnTick")
- Replace confusing NonceRequestedStore with an “idiot-proof” dispatch request checker name.

## Naming finalization

- Dispatch request checker interface
  - Name: IsDispatchRequestedChecker[K comparable]
  - Methods:
    - IsDispatchRequested(ctx sdk.Context, key K) (bool, error)
    - ClearDispatchRequest(ctx sdk.Context, key K) error
  - Map adapter name: MapBoolDispatchRequested[K comparable]

- Queue processor:
  ```go
  type QueueProcessor[T any] struct {
    GetPendingTxs         func(ctx sdk.Context) ([]T, error)
    DispatchTx            func(item T) error
    OnTxConfirmed         func(item T) error            // EVM-specific
    UpdatePendingTxStatus func(item T) error            // Solana-specific
  }
  ```

- Runners (control-flow wrappers):
  - EVMRunner[T]
    - Fields: KeyID, RequestedNonce, Checker IsDispatchRequestedChecker[uint64], Processor QueueProcessor[T], Keeper *Keeper
    - Method: Run(ctx sdk.Context)
  - SolanaRunner[T]
    - Fields: NonceAccountKey, NonceAccount, Checker IsDispatchRequestedChecker[uint64], Processor QueueProcessor[T], Keeper *Keeper
    - Method: Run(ctx sdk.Context)

## Files to change/add

- Add: x/validation/keeper/queues.go
  - IsDispatchRequestedChecker interface
  - MapBoolDispatchRequested adapter
  - QueueProcessor[T]
  - EVMRunner[T], SolanaRunner[T]

- Edit: x/validation/keeper/abci_tx_processing.go
  - Keep processEVMQueue/processSolanaQueue for compatibility but re-implement them by delegating to EVMRunner/SolanaRunner with the new names.
  - Update EVMQueueArgs and SolanaQueueArgs field names to:
    - GetPendingTxs, DispatchTx, OnTxConfirmed, UpdatePendingTxStatus
    - Replace NonceRequestedStore with Checker IsDispatchRequestedChecker[uint64]
  - For this refactor pass, we can deprecate old names with a comment, but we will update all internal usages to new names.

- Edit: x/validation/keeper/abci_zenbtc.go
  - Replace queue calls with runner calls safely (see roll-out steps), or call the updated processEVMQueue/processSolanaQueue with new fields.
  - Ensure no leftovers of old fields appear in composite literals.

- Edit: x/validation/keeper/abci_zentp.go
  - Replace queue calls with runner calls safely (see roll-out steps), or call the updated processSolanaQueue with new fields.

- No changes to business logic in:
  - abci_utils.go, abci.go — except import/use if necessary.

## Roll-out steps (incremental, safe)

1) Introduce queues.go with:
   - IsDispatchRequestedChecker interface
   - MapBoolDispatchRequested adapter for collections.Map[K,bool]
   - QueueProcessor[T]
   - EVMRunner[T], SolanaRunner[T]
   - These compile in isolation and do not alter behavior.

2) Update EVMQueueArgs/SolanaQueueArgs in abci_tx_processing.go to the new names and introduce the Checker field:
   - EVMQueueArgs[T] {
       KeyID, RequestedNonce, Checker IsDispatchRequestedChecker[uint64],
       GetPendingTxs, DispatchTx, OnTxConfirmed
     }
   - SolanaQueueArgs[T] {
       NonceAccountKey, NonceAccount, Checker IsDispatchRequestedChecker[uint64],
       GetPendingTxs, DispatchTx, UpdatePendingTxStatus
     }

3) Reimplement processEVMQueue and processSolanaQueue to delegate to EVMRunner/SolanaRunner respectively:
   - Build a runner from args and call runner.Run(ctx).
   - This keeps call sites unchanged while using the new control flow internally.

4) Build and verify.

5) Replace call sites in abci_zenbtc.go and abci_zentp.go to use runner calls directly (optional in this pass).
   - If changed now, do it carefully:
     - Use a parenthesized composite literal when invoking methods in the same expression: `(EVMRunner[Type]{...}).Run(ctx)`.
     - Ensure the inner QueueProcessor composite only uses the new keys and contains commas at every field-end.
   - Alternatively, stick to the processEVMQueue/processSolanaQueue wrappers for this pass to reduce risk.

6) Build and re-verify.

7) Grep for any remaining references to: Pending:, Dispatch:, OnHeadConfirmed:, OnTick: and update to the new names.

8) Rename NonceRequestedStore usages at call sites to Checker with MapBoolDispatchRequested adapters:
   - EthereumNonceRequested -> MapBoolDispatchRequested[uint64]{M: k.EthereumNonceRequested}
   - SolanaNonceRequested   -> MapBoolDispatchRequested[uint64]{M: k.SolanaNonceRequested}

9) Logging/metrics (optional):
   - EVMRunner/SolanaRunner log:
     - dispatch requested? true/false
     - items queued, head state
     - confirmation/nonce status transitions
     - dispatch success/failure counts

10) Documentation:
   - Update docs/README.md queue segment: “All queues follow: GetPendingTxs → UpdatePendingTxStatus (Solana) → DispatchTx → OnTxConfirmed (EVM).”

## Example patterns (call sites)

- EVM mints (option A: wrapper)
```go
processEVMQueue(k, ctx, EVMQueueArgs[zenbtctypes.PendingMintTransaction]{
  KeyID:          k.zenBTCKeeper.GetEthMinterKeyID(ctx),
  RequestedNonce: oracleData.RequestedEthMinterNonce,
  Checker:        MapBoolDispatchRequested[uint64]{M: k.EthereumNonceRequested},
  GetPendingTxs:  func(ctx sdk.Context) ([]zenbtctypes.PendingMintTransaction, error) { ... },
  DispatchTx:     func(tx zenbtctypes.PendingMintTransaction) error { ... },
  OnTxConfirmed:  func(tx zenbtctypes.PendingMintTransaction) error { ... },
})
```

- EVM mints (option B: runner)
```go
(EVMRunner[zenbtctypes.PendingMintTransaction]{
  KeyID:          k.zenBTCKeeper.GetEthMinterKeyID(ctx),
  RequestedNonce: oracleData.RequestedEthMinterNonce,
  Checker:        MapBoolDispatchRequested[uint64]{M: k.EthereumNonceRequested},
  Processor: QueueProcessor[zenbtctypes.PendingMintTransaction]{
    GetPendingTxs: func(ctx sdk.Context) ([]zenbtctypes.PendingMintTransaction, error) { ... },
    DispatchTx:    func(tx zenbtctypes.PendingMintTransaction) error { ... },
    OnTxConfirmed: func(tx zenbtctypes.PendingMintTransaction) error { ... },
  },
  Keeper: k,
}).Run(ctx)
```

- Solana mints (runner or wrapper, analogous to above), using UpdatePendingTxStatus.

## Gotchas avoided

- Go composite literal invocation must be fully and correctly parenthesized when followed by `.Run(ctx)`.
- All struct fields must end with commas; no leftover fields from old structs.
- Do not mix old field names (Pending/Dispatch/OnHeadConfirmed/OnTick) with the new ones in the same literal.
- Keep wrapper functions during transition to reduce call-site surface changes.

## Testing strategy

- Build after Step 1, Step 3, Step 4.
- Add a thin unit test that injects a fake IsDispatchRequestedChecker and asserts:
  - When not requested, nothing happens.
  - When requested and queue empty, ClearDispatchRequest is called.
  - When requested and head confirmed, OnTxConfirmed is called, then next item is dispatched.
  - Solana: UpdatePendingTxStatus executed for head before dispatch.

## Timeline

- Day 0.5: Implement queues.go and wrapper delegation. Compile.
- Day 0.5: Update call sites (or keep wrappers), compile, smoke test.
- Day 0.5: Polish logs/metrics/doc updates.
- Day 0.5–1: Optional unit tests + follow-up cleanups.

## Decision points

- Adopt runners at call sites now vs. keep wrappers:
  - Safer: keep wrappers first (minimal diff), then adopt runners per function in small PRs.
  - If we must do it all now: follow the exact composite literal rules outlined above.

## Deliverables

- queues.go with the unified system
- Updated abci_tx_processing.go, abci_zenbtc.go, abci_zentp.go
- docs/README.md refresh (queue section) and this plan file
