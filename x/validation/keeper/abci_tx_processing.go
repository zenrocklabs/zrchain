package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	zenbtctypes "github.com/Zenrock-Foundation/zrchain/v6/x/zenbtc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	solSystem "github.com/gagliardetto/solana-go/programs/system"
)

// TxDispatchRequestHandlerInterface defines an explicit controller to check and clear
// whether dispatch has been requested for a given key.
type TxDispatchRequestHandlerInterface[K comparable] interface {
	IsTxDispatchRequested(ctx sdk.Context, key K) (bool, error)
	ClearTxDispatchRequest(ctx sdk.Context, key K) error
}

// TxDispatchRequestHandler adapts a collections.Map[K,bool] to the
// TxDispatchRequestHandler interface.
type TxDispatchRequestHandler[K comparable] struct {
	Store collections.Map[K, bool]
}

func (checker TxDispatchRequestHandler[K]) IsTxDispatchRequested(ctx sdk.Context, key K) (bool, error) {
	isTxDispatchRequested, err := checker.Store.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("error reading dispatch request flag: %w", err)
	}
	return isTxDispatchRequested, nil
}

func (checker TxDispatchRequestHandler[K]) ClearTxDispatchRequest(ctx sdk.Context, key K) error {
	// Remove the flag entry entirely (treats absence as not requested)
	if err := checker.Store.Remove(ctx, key); err != nil && !errors.Is(err, collections.ErrNotFound) {
		return fmt.Errorf("error clearing dispatch request flag: %w", err)
	}
	return nil
}

// TxQueueProcessor captures the business callbacks for pending tx handling.
type TxQueueProcessor[T any] struct {
	GetPendingTxs         func(ctx sdk.Context) ([]T, error)
	DispatchTx            func(item T) error
	OnTxConfirmed         func(item T) error // Ethereum
	UpdatePendingTxStatus func(item T) error // Solana
}

// EthereumTxProcessor encapsulates Ethereum queue control flow (dispatch request, nonce, confirm, dispatch).
type EthereumTxProcessor[T any] struct {
	KeyID                    uint64
	RequestedNonce           uint64
	TxDispatchRequestHandler TxDispatchRequestHandler[uint64]
	TxQueueProcessor         TxQueueProcessor[T]
	Keeper                   *Keeper
}

func (r EthereumTxProcessor[T]) ProcessTxs(ctx sdk.Context) {
	k := r.Keeper
	requested, err := r.TxDispatchRequestHandler.IsTxDispatchRequested(ctx, r.KeyID)
	if err != nil {
		k.Logger(ctx).Error("error checking dispatch request", "keyID", r.KeyID, "error", err)
		return
	}
	if !requested {
		return
	}

	items, err := r.TxQueueProcessor.GetPendingTxs(ctx)
	if err != nil {
		k.Logger(ctx).Error("error getting pending transactions", "error", err)
		return
	}
	if len(items) == 0 {
		if err := r.TxDispatchRequestHandler.ClearTxDispatchRequest(ctx, r.KeyID); err != nil {
			k.Logger(ctx).Error("error clearing dispatch request", "keyID", r.KeyID, "error", err)
		}
		return
	}

	nonceData, err := k.getNonceDataWithInit(ctx, r.KeyID)
	if err != nil {
		k.Logger(ctx).Error("error getting nonce data", "keyID", r.KeyID, "error", err)
		return
	}
	k.Logger(ctx).Info("Nonce info", "nonce", nonceData.Nonce, "prev", nonceData.PrevNonce, "counter", nonceData.Counter, "skip", nonceData.Skip, "requested", r.RequestedNonce)

	if nonceData.Nonce != 0 && r.RequestedNonce == 0 {
		return
	}

	nonceUpdated, err := handleNonceUpdate(k, ctx, r.KeyID, r.RequestedNonce, nonceData, items[0], r.TxQueueProcessor.OnTxConfirmed)
	if err != nil {
		k.Logger(ctx).Error("error handling nonce update", "keyID", r.KeyID, "error", err)
		return
	}

	if len(items) == 1 && nonceUpdated {
		if err := r.TxDispatchRequestHandler.ClearTxDispatchRequest(ctx, r.KeyID); err != nil {
			k.Logger(ctx).Error("error clearing dispatch request", "keyID", r.KeyID, "error", err)
		}
		return
	}

	if nonceData.Skip {
		return
	}

	idx := 0
	if nonceUpdated {
		idx = 1
	}
	if idx < len(items) {
		if err := r.TxQueueProcessor.DispatchTx(items[idx]); err != nil {
			k.Logger(ctx).Error("tx dispatch callback error", "keyID", r.KeyID, "error", err)
		}
	}
}

// SolanaTxProcessor encapsulates Solana queue control flow (dispatch request, status, dispatch).
type SolanaTxProcessor[T any] struct {
	NonceAccountKey          uint64
	NonceAccount             *solSystem.NonceAccount
	TxDispatchRequestHandler TxDispatchRequestHandler[uint64]
	TxQueueProcessor         TxQueueProcessor[T]
	Keeper                   *Keeper
}

func (r SolanaTxProcessor[T]) ProcessTxs(ctx sdk.Context) {
	k := r.Keeper
	requested, err := r.TxDispatchRequestHandler.IsTxDispatchRequested(ctx, r.NonceAccountKey)
	if err != nil {
		k.Logger(ctx).Error("error checking dispatch request", "nonce_account_key", r.NonceAccountKey, "error", err)
		return
	}
	if !requested {
		k.Logger(ctx).Info("SolanaTxProcessor: dispatch not requested", "nonce_account_key", r.NonceAccountKey)
		return
	}

	items, err := r.TxQueueProcessor.GetPendingTxs(ctx)
	if err != nil {
		k.Logger(ctx).Error("error getting pending transactions", "error", err)
		return
	}
	k.Logger(ctx).Info("SolanaTxProcessor: fetched pending transactions", "nonce_account_key", r.NonceAccountKey, "count", len(items))
	if len(items) == 0 {
		if err := r.TxDispatchRequestHandler.ClearTxDispatchRequest(ctx, r.NonceAccountKey); err != nil {
			k.Logger(ctx).Error("error clearing solana dispatch request", "nonce_account_key", r.NonceAccountKey, "error", err)
		}
		k.Logger(ctx).Info("SolanaTxProcessor: no pending transactions, cleared dispatch request", "nonce_account_key", r.NonceAccountKey)
		return
	}

	if r.NonceAccount == nil || r.NonceAccount.Nonce.IsZero() {
		k.Logger(ctx).Error("solana nonce is zero or missing", "nonce_account_key", r.NonceAccountKey)
		return
	}

	if r.TxQueueProcessor.UpdatePendingTxStatus != nil {
		if err := r.TxQueueProcessor.UpdatePendingTxStatus(items[0]); err != nil {
			k.Logger(ctx).Error("error updating solana tx status", "nonce_account_key", r.NonceAccountKey, "error", err)
			return
		}
	}
	k.Logger(ctx).Info("SolanaTxProcessor: dispatching transaction", "nonce_account_key", r.NonceAccountKey)
	if err := r.TxQueueProcessor.DispatchTx(items[0]); err != nil {
		k.Logger(ctx).Error("tx dispatch callback error", "nonce_account_key", r.NonceAccountKey, "error", err)
	}
}

// EthereumTxQueueArgs describes the parameters needed to process an Ethereum-based tx queue.
type EthereumTxQueueArgs[T any] struct {
	KeyID                    uint64
	RequestedNonce           uint64
	DispatchRequestedChecker TxDispatchRequestHandler[uint64]
	GetPendingTxs            func(ctx sdk.Context) ([]T, error)
	DispatchTx               func(tx T) error
	OnTxConfirmed            func(tx T) error // called when the head tx is confirmed (nonce advanced)
}

// SolanaTxQueueArgs describes the parameters needed to process a Solana-based tx queue.
type SolanaTxQueueArgs[T any] struct {
	NonceAccountKey          uint64
	NonceAccount             *solSystem.NonceAccount
	DispatchRequestedChecker TxDispatchRequestHandler[uint64]
	GetPendingTxs            func(ctx sdk.Context) ([]T, error)
	DispatchTx               func(tx T) error
	UpdatePendingTxStatus    func(tx T) error // status/timeout checks for head each block
}

// processEthereumTxQueue remains for backward-compat call sites; it delegates to EthereumTxProcessor.
func processEthereumTxQueue[T any](k *Keeper, ctx sdk.Context, args EthereumTxQueueArgs[T]) {
	(EthereumTxProcessor[T]{
		KeyID:                    args.KeyID,
		RequestedNonce:           args.RequestedNonce,
		TxDispatchRequestHandler: args.DispatchRequestedChecker,
		TxQueueProcessor: TxQueueProcessor[T]{
			GetPendingTxs: args.GetPendingTxs,
			DispatchTx:    args.DispatchTx,
			OnTxConfirmed: args.OnTxConfirmed,
		},
		Keeper: k,
	}).ProcessTxs(ctx)
}

// processSolanaTxQueue processes a Solana queue with clear nonce and status/timeout semantics.
func processSolanaTxQueue[T any](k *Keeper, ctx sdk.Context, args SolanaTxQueueArgs[T]) {
	(SolanaTxProcessor[T]{
		NonceAccountKey:          args.NonceAccountKey,
		NonceAccount:             args.NonceAccount,
		TxDispatchRequestHandler: args.DispatchRequestedChecker,
		TxQueueProcessor: TxQueueProcessor[T]{
			GetPendingTxs:         args.GetPendingTxs,
			DispatchTx:            args.DispatchTx,
			UpdatePendingTxStatus: args.UpdatePendingTxStatus,
		},
		Keeper: k,
	}).ProcessTxs(ctx)
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
