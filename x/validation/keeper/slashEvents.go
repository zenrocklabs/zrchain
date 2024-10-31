package keeper

import (
	"errors"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v5/x/validation/types"
)

func (k Keeper) CreateSlashEvent(ctx sdk.Context, slashEvent *types.SlashEvent) (uint64, error) {
	var id uint64

	count, err := k.SlashEventCount.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return 0, err
		} else {
			id = 1
		}
	} else {
		id = count + 1
	}

	if err := k.SlashEvents.Set(ctx, id, *slashEvent); err != nil {
		return 0, err
	}

	if err := k.SlashEventCount.Set(ctx, id); err != nil {
		return 0, err
	}

	return id, nil
}
