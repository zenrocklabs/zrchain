# Sidecar Generics Refactor Analysis

**Author**: Peyton-Spencer\
**Date**: October 23, 2025\
**Topics Covered**: Sidecar architecture, Go generics adoption, reflection performance costs

---

## Table of Contents

- [Overview](#overview)
- [Key Accomplishments](#key-accomplishments)
- [Technical Changes](#technical-changes)
- [Code References](#code-references)
- [Challenges & Solutions](#challenges--solutions)
- [Decisions Made](#decisions-made)
- [Testing & Iteration](#testing--iteration)
- [Benchmark Plan](#benchmark-plan)
- [Module & State Changes](#module--state-changes)
- [Security & Decentralization](#security--decentralization)
- [Next Steps](#next-steps)

---

## Executive Summary

Audited the sidecar’s Solana event pipeline and documented a generics-based refactor plan. The current implementation leans heavily on `[]any` slices and reflection, which hurts readability and performance. The proposed design introduces typed adapters and generic helpers so we can remove most dynamic typing while keeping event decoding flexible.

[🔝 back to top](#table-of-contents)

---

## Overview

- Identified every location where the sidecar stores Solana events as `[]any`.
- Highlighted the reflection hot paths inside `processMintTransaction` and `processBurnTransaction`.
- Outlined a type-safe alternative that uses Go generics and generated adapters to convert SPL binding structs into API events without `reflect`.

[🔝 back to top](#table-of-contents)

---

## Key Accomplishments

- Catalogued `any` usage across mint, burn, and pending-transaction flows.
- Quantified risks: runtime panics, extra allocations, maintenance overhead.
- Proposed a staged migration strategy to typed event envelopes and generic processing.

[🔝 back to top](#table-of-contents)

---

## Technical Changes

### Current Pain Points

- `processSignatures` and downstream helpers propagate `[]any` (`sidecar/oracle.go:3097-3492`).
- Event decoders rely on reflection to access struct fields (`sidecar/oracle.go:1740-1766`, `sidecar/oracle.go:3601-3629`).
- Pending transaction retries rehydrate events via type assertions (`sidecar/oracle.go:2213-2402`, `sidecar/oracle.go:2652-2755`).

### Proposed Refactor

- Define small adapters that convert generated SPL structs into the API types without reflection.

  ```go
  type SolanaEventAdapter[T any] interface {
      ToMintEvent(slot uint64, sig solana.Signature, idx int) (api.SolanaMintEvent, bool)
      ToBurnEvent(slot uint64, sig solana.Signature, idx int) (api.BurnEvent, bool)
  }
  ```

- Introduce generic helpers so the compiler enforces event types and we avoid `[]any`.

  ```go
  func processTransaction[T SolanaEventAdapter[T]](
      tx *solrpc.GetTransactionResult,
      program solana.PublicKey,
      sig solana.Signature,
      decode func(*solrpc.GetTransactionResult, solana.PublicKey) ([]T, error),
  ) ([]api.SolanaMintEvent, error) {
      events, err := decode(tx, program)
      if err != nil {
          return nil, err
      }
      out := make([]api.SolanaMintEvent, 0, len(events))
      for idx, ev := range events {
          if mint, ok := ev.ToMintEvent(uint64(tx.Slot), sig, idx); ok {
              out = append(out, mint)
          }
      }
      return out, nil
  }
  ```

- Replace shared `[]any` collections with concrete slices (e.g., `[]api.SolanaMintEvent`) and constrain generic helpers with small interfaces so downstream code remains type-safe.
- Generalize `processSignatures` with a type parameter so batching stays generic without `any`.

  ```go
  type SolanaEvent interface {
      TxSignature() string
  }

  func processSignatures[E SolanaEvent](
      ctx context.Context,
      sigs []*solrpc.TransactionSignature,
      fetch func(context.Context, *solrpc.TransactionSignature) (E, error),
  ) ([]E, solana.Signature, []string, error) {
      out := make([]E, 0, len(sigs))
      for _, sig := range sigs {
          evt, err := fetch(ctx, sig)
          if err != nil {
              // collect for retry queue
              continue
          }
          out = append(out, evt)
      }
      newest := solana.Signature{}
      if n := len(out); n > 0 {
          newest = sigs[n-1].Signature
      }
      return out, newest, nil, nil
  }
  ```

- Update pending transaction storage to keep typed payloads, removing reflective reprocessing.

[🔝 back to top](#table-of-contents)

---

## Code References

- [`sidecar/oracle.go`](../../sidecar/oracle.go) – `[]any` usage in mint/burn pipelines (`1595-1907`, `2652-2781`, `3097-3492`, `3601-3629`).
- [`sidecar/oracle_test.go`](../../sidecar/oracle_test.go) – Test helpers returning `[]any` (`235-236`).
- [`sidecar/utils.go`](../../sidecar/utils.go) – Event merge logic keyed by `any`-derived hashes (`314-356`).

[🔝 back to top](#table-of-contents)

---

## Challenges & Solutions

### Challenge 1: Dynamic Event Decoding

**Problem**: SPL bindings emit distinct struct types per program, so early implementations used `interface{}` plus reflection to stay generic.

**Solution**: Generate small adapters for each SPL contract and wire them into generic helpers. This keeps the call sites generic without sacrificing static typing.

[🔝 back to top](#table-of-contents)

---

## Decisions Made

- Favor generics over reflection to improve performance and safety in hot paths.
- Stage the migration: prototype zenBTC mint path first, then extend to burns and zenZEC once metrics confirm parity.

[🔝 back to top](#table-of-contents)

---

## Testing & Iteration

### Tests Performed

- No automated tests executed yet; refactor plan only.

### Iteration Cycles

| Iteration | Focus                     | Changes                         | Result              |
| --------- | ------------------------- | --------------------------------| ------------------- |
| 1         | Debt assessment           | Catalogued `any`/`reflect` usage| Refactor plan ready |

[🔝 back to top](#table-of-contents)

---

## Benchmark Plan

- Add Go benchmark cases alongside `sidecar/oracle_test.go` to quantify CPU cost before and after the generics migration. Representative ideas:
  - `BenchmarkProcessTransaction_Mint` exercising the current reflection-based pipeline.
  - `BenchmarkProcessSignatures_Batch10` to measure batching overhead and pending-queue churn.
- Use `testing.B` with canned `GetTransactionResult` fixtures so runs are deterministic and avoid RPC latency masking CPU differences.
- Capture baseline allocations (`b.ReportAllocs()`) to validate the generics path actually reduces heap pressure.
- Prioritize benchmark coverage for the desync remediation work first; generics optimization should not delay the gossip and watermark fixes already identified in [`docs/reports/2025-10-23-sidecar-desync-analysis.md`](../../docs/reports/2025-10-23-sidecar-desync-analysis.md).

[🔝 back to top](#table-of-contents)

---

## Module & State Changes

- No module or chain state changes; documentation-only assessment.

[🔝 back to top](#table-of-contents)

---

## Security & Decentralization

- Removing reflection reduces the chance of runtime panics that could desync validators’ sidecars, indirectly strengthening oracle availability.
- Typed event handling lowers the risk of malformed data leaking into vote extensions.

[🔝 back to top](#table-of-contents)

---

## Next Steps

- Implement the generic zenBTC mint adapter and benchmark against current reflection-based code.
- Extend adapters to zenZEC and burn flows once performance gains are confirmed.
- Update pending transaction persistence to store typed payloads and remove the remaining `any` helpers.

[🔝 back to top](#table-of-contents)
