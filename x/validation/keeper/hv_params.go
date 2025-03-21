package keeper

import (
	"context"

	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
)

func (k Keeper) GetAVSRewardsRate(ctx context.Context) math.LegacyDec {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return types.DefaultAVSRewardsRate
	}
	return params.AVSRewardsRate
}

func (k Keeper) GetBlockTime(ctx context.Context) int64 {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return types.DefaultBlockTime
	}
	return params.BlockTime
}

func (k Keeper) GetStakeableAssets(ctx context.Context) []*types.AssetData {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return types.GetDefaultStakeableAssets(ctx)
	}
	return params.StakeableAssets
}

func (k Keeper) GetHVParamsAuthority(ctx context.Context) string {
	return k.authority
}

// GetPriceRetentionBlockRange returns the price retention block range
func (k Keeper) GetPriceRetentionBlockRange(ctx context.Context) int64 {
	params, err := k.HVParams.Get(ctx)
	if err != nil || params.PriceRetentionBlockRange <= 0 {
		return types.DefaultPriceRetentionBlockRange
	}
	return params.PriceRetentionBlockRange
}
