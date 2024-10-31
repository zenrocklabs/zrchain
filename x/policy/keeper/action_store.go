package keeper

import (
	"errors"

	"cosmossdk.io/collections"
	"github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CreateAction(ctx sdk.Context, action *types.Action) (uint64, error) {
	count, err := k.ActionCount.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return 0, err
		}
		action.Id = 1
	} else {
		action.Id = count + 1
	}

	if err := k.ActionStore.Set(ctx, action.Id, *action); err != nil {
		return 0, err
	}

	if err := k.ActionCount.Set(ctx, action.Id); err != nil {
		return 0, err
	}

	return action.Id, nil
}
