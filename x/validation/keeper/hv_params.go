package keeper

import (
	"context"

	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v5/x/validation/types"
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
		return types.DefaultStakeableAssets
	}
	return params.StakeableAssets
}

func (k Keeper) GetHVParamsAuthority(ctx context.Context) string {
	params, err := k.HVParams.Get(ctx)
	if err != nil || params.Authority == "" {
		return types.DefaultHVParamsAuthority
	}
	return params.Authority
}
