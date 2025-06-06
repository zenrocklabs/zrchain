package types

import (
	context "context"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	DefaultAVSRewardsRate, _              = math.LegacyNewDecFromStr("0.03") // 0.03 == 3% APR
	DefaultBlockTime                int64 = 5                                // seconds
	DefaultPriceRetentionBlockRange int64 = 100                              // blocks

	DefaultTestnetStakeableAssets = []*AssetData{
		{Asset: Asset_ROCK, Precision: 6, PriceUSD: math.LegacyZeroDec()},
	}
	DefaultMainnetStakeableAssets = []*AssetData{
		{Asset: Asset_ROCK, Precision: 6, PriceUSD: math.LegacyZeroDec()},
		{Asset: Asset_BTC, Precision: 8, PriceUSD: math.LegacyZeroDec()},
	}
)

// NewParams creates a new Params instance
func NewHVParams(avsRewardsRate math.LegacyDec, blockTime int64, stakeableAssets []*AssetData, priceRetentionBlockRange int64) *HVParams {
	return &HVParams{
		AVSRewardsRate:           avsRewardsRate,
		BlockTime:                blockTime,
		StakeableAssets:          stakeableAssets,
		PriceRetentionBlockRange: priceRetentionBlockRange,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultHVParams(ctx context.Context) *HVParams {
	return NewHVParams(
		DefaultAVSRewardsRate,
		DefaultBlockTime,
		GetDefaultStakeableAssets(ctx),
		DefaultPriceRetentionBlockRange,
	)
}

func GetDefaultStakeableAssets(ctx context.Context) []*AssetData {
	if strings.HasPrefix(sdk.UnwrapSDKContext(ctx).ChainID(), "diamond") {
		return DefaultMainnetStakeableAssets
	}
	return DefaultTestnetStakeableAssets
}
