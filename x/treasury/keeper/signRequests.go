package keeper

import (
	"errors"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
)

func (k Keeper) CreateSignRequest(ctx sdk.Context, signRequest *types.SignRequest) (uint64, error) {
	count, err := k.SignRequestCount.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return 0, err
		}
		signRequest.Id = 1
	} else {
		signRequest.Id = count + 1
	}

	if err := k.SignRequestStore.Set(ctx, signRequest.Id, *signRequest); err != nil {
		return 0, err
	}

	if err := k.SignRequestCount.Set(ctx, signRequest.Id); err != nil {
		return 0, err
	}

	return signRequest.Id, nil
}
