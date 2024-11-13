package keeper

import (
	"errors"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
)

func (k Keeper) CreateICATransactionRequest(ctx sdk.Context, icaTransactionRequest *types.ICATransactionRequest) (uint64, error) {
	count, err := k.ICATransactionRequestCount.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return 0, err
		}
		icaTransactionRequest.Id = 1
	} else {
		icaTransactionRequest.Id = count + 1
	}

	if err := k.ICATransactionRequestStore.Set(ctx, icaTransactionRequest.Id, *icaTransactionRequest); err != nil {
		return 0, err
	}

	if err := k.ICATransactionRequestCount.Set(ctx, icaTransactionRequest.Id); err != nil {
		return 0, err
	}

	return icaTransactionRequest.Id, nil
}
