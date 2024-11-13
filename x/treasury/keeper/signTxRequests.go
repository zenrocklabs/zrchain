package keeper

import (
	"errors"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
)

func (k Keeper) CreateSignTransactionRequest(ctx sdk.Context, signTransactionRequest *types.SignTransactionRequest) (uint64, error) {
	count, err := k.SignTransactionRequestCount.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return 0, err
		}
		signTransactionRequest.Id = 1
	} else {
		signTransactionRequest.Id = count + 1
	}

	if err := k.SignTransactionRequestStore.Set(ctx, signTransactionRequest.Id, *signTransactionRequest); err != nil {
		return 0, err
	}

	if err := k.SignTransactionRequestCount.Set(ctx, signTransactionRequest.Id); err != nil {
		return 0, err
	}

	return signTransactionRequest.Id, nil
}
