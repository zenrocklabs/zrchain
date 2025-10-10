package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"

	types "github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

func (k Keeper) CreatePendingMintTransaction(ctx context.Context, mintTransaction *types.PendingMintTransaction) (uint64, error) {
	count, err := k.PendingMintTransactionCount.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return 0, err
		} else {
			mintTransaction.Id = 1
		}
	} else {
		mintTransaction.Id = count + 1
	}

	mintTransaction.Status = types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED

	if err := k.PendingMintTransactionsMap.Set(ctx, mintTransaction.Id, *mintTransaction); err != nil {
		return 0, err
	}

	if err := k.PendingMintTransactionCount.Set(ctx, mintTransaction.Id); err != nil {
		return 0, err
	}

	return mintTransaction.Id, nil
}
