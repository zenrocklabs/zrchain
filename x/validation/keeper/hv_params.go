package keeper

import (
	"context"

	"cosmossdk.io/math"

	"github.com/Zenrock-Foundation/zrchain/v5/x/validation/types"
)

var (
	DefaultAVSRewardsRate, _       = math.LegacyNewDecFromStr("0.03") // 0.03 == 3% APR
	DefaultBlockTime         int64 = 1                                // seconds
	DefaultStakeableAssets         = []*types.AssetData{
		{Asset: types.Asset_ROCK, Precision: 6, PriceUSD: math.LegacyZeroDec()},
		{Asset: types.Asset_zenBTC, Precision: 8, PriceUSD: math.LegacyZeroDec()},
		{Asset: types.Asset_stETH, Precision: 18, PriceUSD: math.LegacyZeroDec()},
	}
	DefaultHVParamsAuthority = "zen1sd3fwcpw2mdw3pxexmlg34gsd78r0sxrk5weh3"
)

// NewParams creates a new Params instance
func NewHVParams(avsRewardsRate math.LegacyDec, blockTime int64, stakeableAssets []*types.AssetData, authority string) *types.HVParams {
	return &types.HVParams{
		AVSRewardsRate:  avsRewardsRate,
		BlockTime:       blockTime,
		StakeableAssets: stakeableAssets,
		Authority:       authority,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultHVParams() *types.HVParams {
	return NewHVParams(
		DefaultAVSRewardsRate,
		DefaultBlockTime,
		DefaultStakeableAssets,
		DefaultHVParamsAuthority,
	)
}

func (k Keeper) GetAVSRewardsRate(ctx context.Context) math.LegacyDec {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return DefaultAVSRewardsRate
	}
	return params.AVSRewardsRate
}

func (k Keeper) GetBlockTime(ctx context.Context) int64 {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return DefaultBlockTime
	}
	return params.BlockTime
}

func (k Keeper) GetStakeableAssets(ctx context.Context) []*types.AssetData {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return DefaultStakeableAssets
	}
	return params.StakeableAssets
}

func (k Keeper) GetHVParamsAuthority(ctx context.Context) string {
	params, err := k.HVParams.Get(ctx)
	if err != nil || params.Authority == "" {
		return DefaultHVParamsAuthority
	}
	return params.Authority
}
