package keeper

import "context"

const defaultBTL = 20

func (k Keeper) GetDefaultBTL(ctx context.Context) (uint64, error) {
	params, err := k.ParamStore.Get(ctx)
	if err != nil {
		return defaultBTL, err
	}

	if params.DefaultBtl == 0 {
		return defaultBTL, nil
	}

	return params.DefaultBtl, nil
}
