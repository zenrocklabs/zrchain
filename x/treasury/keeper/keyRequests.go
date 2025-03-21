package keeper

import (
	"errors"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
)

func (k Keeper) AppendKeyRequest(ctx sdk.Context, keyRequest *types.KeyRequest) (uint64, error) {
	count, err := k.KeyRequestCount.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return 0, err
		} else {
			keyRequest.Id = 1
		}
	} else {
		keyRequest.Id = count + 1
	}

	if err := k.KeyRequestStore.Set(ctx, keyRequest.Id, *keyRequest); err != nil {
		return 0, err
	}

	if err := k.KeyRequestCount.Set(ctx, keyRequest.Id); err != nil {
		return 0, err
	}

	return keyRequest.Id, nil
}
