package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"

	types "github.com/Zenrock-Foundation/zrchain/v6/x/zenbtc/types"
)

func (k Keeper) CreateBurnEvent(ctx context.Context, burnEvent *types.BurnEvent) (uint64, error) {
	count, err := k.BurnEventCount.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return 0, err
		}
		burnEvent.Id = 1
	} else {
		burnEvent.Id = count + 1
	}

	if burnEvent.Status == types.BurnStatus_BURN_STATUS_UNSPECIFIED {
		burnEvent.Status = types.BurnStatus_BURN_STATUS_UNSTAKING
	}

	if err := k.BurnEvents.Set(ctx, burnEvent.Id, *burnEvent); err != nil {
		return 0, err
	}

	if err := k.BurnEventCount.Set(ctx, burnEvent.Id); err != nil {
		return 0, err
	}

	return burnEvent.Id, nil
}
