package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"

	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

// CreatePendingMintTransaction stores a new pending mint transaction for the specified asset.
func (k Keeper) CreatePendingMintTransaction(ctx context.Context, mintTx *dcttypes.PendingMintTransaction) (uint64, error) {
	if mintTx.Asset == dcttypes.Asset_ASSET_UNSPECIFIED {
		return 0, dcttypes.ErrUnknownAsset
	}

	assetKey, err := k.getAssetKey(mintTx.Asset)
	if err != nil {
		return 0, err
	}

	count, err := k.PendingMintTransactionCount.Get(ctx, assetKey)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return 0, err
		}
		count = 0
	}

	nextID := count + 1
	mintTx.Id = nextID
	mintTx.Status = dcttypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED

	key, err := k.pendingMintKey(mintTx.Asset, nextID)
	if err != nil {
		return 0, err
	}

	if err := k.PendingMintTransactions.Set(ctx, key, *mintTx); err != nil {
		return 0, err
	}

	if err := k.PendingMintTransactionCount.Set(ctx, assetKey, nextID); err != nil {
		return 0, err
	}

	return nextID, nil
}
