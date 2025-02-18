package types

import (
	"cosmossdk.io/math"
)

var (
	DefaultAVSRewardsRate, _       = math.LegacyNewDecFromStr("0.03") // 0.03 == 3% APR
	DefaultBlockTime         int64 = 5                                // seconds
	DefaultStakeableAssets         = []*AssetData{
		{Asset: Asset_ROCK, Precision: 6, PriceUSD: math.LegacyZeroDec()},
		{Asset: Asset_BTC, Precision: 8, PriceUSD: math.LegacyZeroDec()},
		{Asset: Asset_ETH, Precision: 18, PriceUSD: math.LegacyZeroDec()},
	}
	DefaultHVParamsAuthority              = "zen1sd3fwcpw2mdw3pxexmlg34gsd78r0sxrk5weh3"
	DefaultPriceRetentionBlockRange int64 = 100
)

// NewParams creates a new Params instance
func NewHVParams(avsRewardsRate math.LegacyDec, blockTime int64, stakeableAssets []*AssetData, authority string, priceRetentionBlockRange int64) *HVParams {
	return &HVParams{
		AVSRewardsRate:           avsRewardsRate,
		BlockTime:                blockTime,
		StakeableAssets:          stakeableAssets,
		PriceRetentionBlockRange: priceRetentionBlockRange,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultHVParams() *HVParams {
	return NewHVParams(
		DefaultAVSRewardsRate,
		DefaultBlockTime,
		DefaultStakeableAssets,
		DefaultHVParamsAuthority,
		DefaultPriceRetentionBlockRange,
	)
}
