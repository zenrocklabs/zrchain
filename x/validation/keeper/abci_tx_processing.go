package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	solSystem "github.com/gagliardetto/solana-go/programs/system"
	zenbtctypes "github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

// checkForUpdateAndDispatchTx processes nonce updates and transaction dispatch.
// It contains separate logic for Ethereum and Solana based transactions due to their
// different nonce mechanisms and transaction lifecycles.
func checkForUpdateAndDispatchTx[T any](
	k *Keeper,
	ctx sdk.Context,
	keyID uint64,
	requestedEthNonce *uint64,
	requestedSolNonce *solSystem.NonceAccount,
	nonceReqStore collections.Map[uint64, bool],
	pendingTxs []T,
	txDispatchCallback func(tx T) error,
	txContinuationCallback func(tx T) error,
) {
	if len(pendingTxs) == 0 {
		return
	}

	// Ethereum transaction processing flow.
	if requestedEthNonce != nil {
		nonceData, err := k.getNonceDataWithInit(ctx, keyID)
		if err != nil {
			k.Logger(ctx).Error("error getting nonce data", "keyID", keyID, "error", err)
			return
		}
		k.Logger(ctx).Info("Nonce info",
			"nonce", nonceData.Nonce,
			"prev", nonceData.PrevNonce,
			"counter", nonceData.Counter,
			"skip", nonceData.Skip,
			"requested", requestedEthNonce,
		)
		if nonceData.Nonce != 0 && *requestedEthNonce == 0 {
			return
		}

		nonceUpdated, err := handleNonceUpdate(k, ctx, keyID, *requestedEthNonce, nonceData, pendingTxs[0], txContinuationCallback)
		if err != nil {
			k.Logger(ctx).Error("error handling nonce update", "keyID", keyID, "error", err)
			return
		}

		if len(pendingTxs) == 1 && nonceUpdated {
			if err := k.clearNonceRequest(ctx, nonceReqStore, keyID); err != nil {
				k.Logger(ctx).Error("error clearing ethereum nonce request", "keyID", keyID, "error", err)
			}
			return
		}

		if nonceData.Skip {
			return
		}

		// If tx[0] confirmed on-chain via nonce increment, dispatch tx[1]. If not then retry dispatching tx[0].
		txIndex := 0
		if nonceUpdated {
			txIndex = 1
		}

		if len(pendingTxs) <= txIndex {
			return
		}

		if err := txDispatchCallback(pendingTxs[txIndex]); err != nil {
			k.Logger(ctx).Error("tx dispatch callback error", "keyID", keyID, "error", err)
		}
		return
	}

	// Solana transaction processing flow.
	if requestedSolNonce != nil {
		k.Logger(ctx).Info("processing solana transaction with nonce", "nonce", requestedSolNonce.Nonce)

		if requestedSolNonce.Nonce.IsZero() {
			k.Logger(ctx).Error("solana nonce is zero")
			return
		}

		// For Solana, `txContinuationCallback` is a misnomer. It's a status/timeout checker for the head of the queue.
		// We call it, and then attempt to dispatch the same transaction. The dispatch is idempotent.
		if err := txContinuationCallback(pendingTxs[0]); err != nil {
			k.Logger(ctx).Error("error handling solana transaction status check", "keyID", keyID, "error", err)
			return
		}

		// If tx[0] is still pending, dispatch it. The dispatch callback is idempotent.
		if err := txDispatchCallback(pendingTxs[0]); err != nil {
			k.Logger(ctx).Error("tx dispatch callback error", "keyID", keyID, "error", err)
		}
	}
}

// processTransaction is a generic helper that encapsulates the common logic for nonce update and tx dispatch.
func processTransaction[T any](
	k *Keeper,
	ctx sdk.Context,
	keyID uint64,
	requestedEthNonce *uint64,
	requestedSolNonce *solSystem.NonceAccount,
	pendingGetter func(ctx sdk.Context) ([]T, error),
	txDispatchCallback func(tx T) error,
	txContinuationCallback func(tx T) error,
) {
	nonceReqStore := k.EthereumNonceRequested
	if requestedEthNonce == nil {
		nonceReqStore = k.SolanaNonceRequested
	}

	isRequested, err := isNonceRequested(ctx, nonceReqStore, keyID)
	if err != nil {
		k.Logger(ctx).Error("error checking nonce request state", "keyID", keyID, "error", err)
		return
	}
	if !isRequested {
		return
	}

	pendingTxs, err := pendingGetter(ctx)
	if err != nil {
		k.Logger(ctx).Error("error getting pending transactions", "error", err)
		return
	}

	if len(pendingTxs) == 0 {
		if err := k.clearNonceRequest(ctx, nonceReqStore, keyID); err != nil {
			k.Logger(ctx).Error("error clearing ethereum nonce request", "keyID", keyID, "error", err)
		}
		return
	}
	checkForUpdateAndDispatchTx(k, ctx, keyID, requestedEthNonce, requestedSolNonce, nonceReqStore, pendingTxs, txDispatchCallback, txContinuationCallback)
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
