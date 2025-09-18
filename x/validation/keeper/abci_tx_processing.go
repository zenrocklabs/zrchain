package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	solSystem "github.com/gagliardetto/solana-go/programs/system"
	zenbtctypes "github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

// EVMQueueArgs describes the parameters needed to process an EVM-based tx queue.
type EVMQueueArgs[T any] struct {
	KeyID               uint64
	RequestedNonce      uint64
	NonceRequestedStore collections.Map[uint64, bool]
	GetPendingTxs       func(ctx sdk.Context) ([]T, error)
	DispatchTx          func(tx T) error
	OnTxConfirmed       func(tx T) error // called when the head tx is confirmed (nonce advanced)
}

// SolanaQueueArgs describes the parameters needed to process a Solana-based tx queue.
type SolanaQueueArgs[T any] struct {
	NonceAccountKey        uint64
	NonceAccount           *solSystem.NonceAccount
	NonceRequestedStore    collections.Map[uint64, bool]
	GetPendingTxs          func(ctx sdk.Context) ([]T, error)
	DispatchTx             func(tx T) error
	UpdatePendingTxStatus  func(tx T) error // status/timeout checks for head each block
}

// processEVMQueue processes an EVM queue with clear nonce-advance and dispatch semantics.
func processEVMQueue[T any](k *Keeper, ctx sdk.Context, args EVMQueueArgs[T]) {
	isRequested, err := isNonceRequested(ctx, args.NonceRequestedStore, args.KeyID)
	if err != nil {
		k.Logger(ctx).Error("error checking nonce request state", "keyID", args.KeyID, "error", err)
		return
	}
	if !isRequested {
		return
	}

	pendingTxs, err := args.GetPendingTxs(ctx)
	if err != nil {
		k.Logger(ctx).Error("error getting pending transactions", "error", err)
		return
	}
	if len(pendingTxs) == 0 {
		if err := k.clearNonceRequest(ctx, args.NonceRequestedStore, args.KeyID); err != nil {
			k.Logger(ctx).Error("error clearing ethereum nonce request", "keyID", args.KeyID, "error", err)
		}
		return
	}

	nonceData, err := k.getNonceDataWithInit(ctx, args.KeyID)
	if err != nil {
		k.Logger(ctx).Error("error getting nonce data", "keyID", args.KeyID, "error", err)
		return
	}
	k.Logger(ctx).Info("Nonce info", "nonce", nonceData.Nonce, "prev", nonceData.PrevNonce, "counter", nonceData.Counter, "skip", nonceData.Skip, "requested", args.RequestedNonce)

	// Defensive: if we have a known non-zero nonce but requested is zero (no consensus yet), do nothing.
	if nonceData.Nonce != 0 && args.RequestedNonce == 0 {
		return
	}

	// If on-chain nonce advanced for head, run continuation callback and update prev nonce.
	nonceUpdated, err := handleNonceUpdate(k, ctx, args.KeyID, args.RequestedNonce, nonceData, pendingTxs[0], args.OnTxConfirmed)
	if err != nil {
		k.Logger(ctx).Error("error handling nonce update", "keyID", args.KeyID, "error", err)
		return
	}

	// If only one pending and it's now confirmed (nonce advanced), we can clear the request for this key.
	if len(pendingTxs) == 1 && nonceUpdated {
		if err := k.clearNonceRequest(ctx, args.NonceRequestedStore, args.KeyID); err != nil {
			k.Logger(ctx).Error("error clearing ethereum nonce request", "keyID", args.KeyID, "error", err)
		}
		return
	}

	if nonceData.Skip {
		return
	}

	// If head confirmed, try to dispatch the next item, else retry dispatching head (idempotent).
	idx := 0
	if nonceUpdated {
		idx = 1
	}
	if idx < len(pendingTxs) {
		if err := args.DispatchTx(pendingTxs[idx]); err != nil {
			k.Logger(ctx).Error("tx dispatch callback error", "keyID", args.KeyID, "error", err)
		}
	}
}

// processSolanaQueue processes a Solana queue with clear nonce and status/timeout semantics.
func processSolanaQueue[T any](k *Keeper, ctx sdk.Context, args SolanaQueueArgs[T]) {
	isRequested, err := isNonceRequested(ctx, args.NonceRequestedStore, args.NonceAccountKey)
	if err != nil {
		k.Logger(ctx).Error("error checking nonce request state", "nonce_account_key", args.NonceAccountKey, "error", err)
		return
	}
	if !isRequested {
		return
	}

	pendingTxs, err := args.GetPendingTxs(ctx)
	if err != nil {
		k.Logger(ctx).Error("error getting pending transactions", "error", err)
		return
	}
	if len(pendingTxs) == 0 {
		if err := k.clearNonceRequest(ctx, args.NonceRequestedStore, args.NonceAccountKey); err != nil {
			k.Logger(ctx).Error("error clearing solana nonce request", "nonce_account_key", args.NonceAccountKey, "error", err)
		}
		return
	}

	if args.NonceAccount == nil || args.NonceAccount.Nonce.IsZero() {
		k.Logger(ctx).Error("solana nonce is zero or missing", "nonce_account_key", args.NonceAccountKey)
		return
	}

	// Update head status/timeout and then attempt dispatch (both idempotent).
	if err := args.UpdatePendingTxStatus(pendingTxs[0]); err != nil {
		k.Logger(ctx).Error("error handling solana transaction status check", "nonce_account_key", args.NonceAccountKey, "error", err)
		return
	}
	if err := args.DispatchTx(pendingTxs[0]); err != nil {
		k.Logger(ctx).Error("tx dispatch callback error", "nonce_account_key", args.NonceAccountKey, "error", err)
	}
}

// getPendingTransactions is a generic helper that walks a store with key type uint64
// and returns a slice of items of type T that satisfy the provided predicate, up to a given limit.
// If limit is 0, all matching items will be returned.
func getPendingTransactions[T any](ctx sdk.Context, store interface {
	Walk(ctx sdk.Context, rng *collections.Range[uint64], fn func(key uint64, value T) (bool, error)) error
}, predicate func(T) bool, firstPendingID uint64, limit int) ([]T, error) {
	var results []T
	queryRange := &collections.Range[uint64]{}
	err := store.Walk(ctx, queryRange.StartInclusive(firstPendingID), func(key uint64, value T) (bool, error) {
		if predicate(value) {
			results = append(results, value)
			if limit > 0 && len(results) >= limit {
				return true, nil
			}
		}
		return false, nil
	})
	return results, err
}

// getNonceDataWithInit gets the nonce data for a key, initializing it if it doesn't exist
func (k *Keeper) getNonceDataWithInit(ctx sdk.Context, keyID uint64) (zenbtctypes.NonceData, error) {
	nonceData, err := k.LastUsedEthereumNonce.Get(ctx, keyID)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return zenbtctypes.NonceData{}, fmt.Errorf("error getting last used ethereum nonce: %w", err)
		}
		nonceData = zenbtctypes.NonceData{Nonce: 0, PrevNonce: 0, Counter: 0, Skip: true}
		if err := k.LastUsedEthereumNonce.Set(ctx, keyID, nonceData); err != nil {
			return zenbtctypes.NonceData{}, fmt.Errorf("error setting last used ethereum nonce: %w", err)
		}
	}
	return nonceData, nil
}

// isNonceRequested checks if a nonce has been requested for the given key
func isNonceRequested(ctx sdk.Context, store collections.Map[uint64, bool], keyID uint64) (bool, error) {
	requested, err := store.Get(ctx, keyID)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("error getting nonce request state: %w", err)
	}
	return requested, nil
}

// handleNonceUpdate handles the nonce update logic and returns whether an update occurred
func handleNonceUpdate[T any](
	k *Keeper,
	ctx sdk.Context,
	keyID uint64,
	requestedNonce uint64,
	nonceData zenbtctypes.NonceData,
	tx T,
	txContinuationCallback func(tx T) error,
) (bool, error) {
	if requestedNonce != nonceData.PrevNonce {
		if err := txContinuationCallback(tx); err != nil {
			return false, fmt.Errorf("tx continuation callback error: %w", err)
		}
		k.Logger(ctx).Warn("nonce updated for key",
			"keyID", keyID,
			"requestedNonce", requestedNonce,
			"prevNonce", nonceData.PrevNonce,
			"currentNonce", nonceData.Nonce,
		)
		nonceData.PrevNonce = nonceData.Nonce
		if err := k.LastUsedEthereumNonce.Set(ctx, keyID, nonceData); err != nil {
			return false, fmt.Errorf("error setting last used Ethereum nonce: %w", err)
		}
		return true, nil
	}
	return false, nil
}

// updateNonces handles updating nonce state for keys used for minting and unstaking.
func (k *Keeper) updateNonces(ctx sdk.Context, oracleData OracleData) {
	for _, key := range k.getZenBTCKeyIDs(ctx) {
		isRequested, err := isNonceRequested(ctx, k.EthereumNonceRequested, key)
		if err != nil {
			k.Logger(ctx).Error("error checking nonce request state", "keyID", key, "error", err)
			continue
		}
		if !isRequested {
			continue
		}

		var currentNonce uint64
		switch key {
		case k.zenBTCKeeper.GetStakerKeyID(ctx):
			currentNonce = oracleData.RequestedStakerNonce
		case k.zenBTCKeeper.GetEthMinterKeyID(ctx):
			currentNonce = oracleData.RequestedEthMinterNonce
		case k.zenBTCKeeper.GetUnstakerKeyID(ctx):
			currentNonce = oracleData.RequestedUnstakerNonce
		case k.zenBTCKeeper.GetCompleterKeyID(ctx):
			currentNonce = oracleData.RequestedCompleterNonce
		default:
			k.Logger(ctx).Error("invalid key ID", "keyID", key)
			continue
		}

		// Avoid erroneously setting nonce to zero if a non-zero nonce exists i.e. blocks with no consensus on VEs.
		nonceData, err := k.getNonceDataWithInit(ctx, key)
		if err != nil {
			k.Logger(ctx).Error("error getting nonce data", "keyID", key, "error", err)
			continue
		}
		if nonceData.Nonce != 0 && currentNonce == 0 {
			continue
		}

		if err := k.updateNonceState(ctx, key, currentNonce); err != nil {
			k.Logger(ctx).Error("error updating nonce state", "keyID", key, "error", err)
		}
	}
}

// clearNonceRequest resets the nonce-request flag for a given key.
func (k *Keeper) clearNonceRequest(ctx sdk.Context, store collections.Map[uint64, bool], keyID uint64) error {
	k.Logger(ctx).Warn("set requested nonce state to false", "keyID", keyID)
	return store.Set(ctx, keyID, false)
}
