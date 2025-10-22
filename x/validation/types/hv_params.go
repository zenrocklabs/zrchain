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
	DefaultPriceRetentionBlockRange int64 = 320                              // blocks
	DefaultVEJailingEnabled         bool  = false                            // enable VE jailing by default
	DefaultVEJailDurationMinutes    int64 = 60                               // 60 minutes jail duration
	DefaultVEWindowSize             int64 = 320                              // 320 blocks window for VE mismatch tracking
	DefaultVEJailThreshold          int64 = 160                              // 160 mismatches before jailing

	DefaultTestnetStakeableAssets = []*AssetData{
		{Asset: Asset_ROCK, Precision: 6, PriceUSD: math.LegacyZeroDec()},
		{Asset: Asset_BTC, Precision: 8, PriceUSD: math.LegacyZeroDec()},
		{Asset: Asset_ZEC, Precision: 8, PriceUSD: math.LegacyZeroDec()},
	}
	DefaultMainnetStakeableAssets = []*AssetData{
		{Asset: Asset_ROCK, Precision: 6, PriceUSD: math.LegacyZeroDec()},
		{Asset: Asset_BTC, Precision: 8, PriceUSD: math.LegacyZeroDec()},
		{Asset: Asset_ZEC, Precision: 18, PriceUSD: math.LegacyZeroDec()},
	}
)

// NewParams creates a new Params instance
func NewHVParams(avsRewardsRate math.LegacyDec, blockTime int64, stakeableAssets []*AssetData, priceRetentionBlockRange int64, veJailingEnabled bool, veJailDurationMinutes int64, veWindowSize int64, veJailThreshold int64) *HVParams {
	return &HVParams{
		AVSRewardsRate:           avsRewardsRate,
		BlockTime:                blockTime,
		StakeableAssets:          stakeableAssets,
		PriceRetentionBlockRange: priceRetentionBlockRange,
		VEJailingEnabled:         veJailingEnabled,
		VEJailDurationMinutes:    veJailDurationMinutes,
		VEWindowSize:             veWindowSize,
		VEJailThreshold:          veJailThreshold,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultHVParams(ctx context.Context) *HVParams {
	return NewHVParams(
		DefaultAVSRewardsRate,
		DefaultBlockTime,
		GetDefaultStakeableAssets(ctx),
		DefaultPriceRetentionBlockRange,
		DefaultVEJailingEnabled,
		DefaultVEJailDurationMinutes,
		DefaultVEWindowSize,
		DefaultVEJailThreshold,
	)
}

func GetDefaultStakeableAssets(ctx context.Context) []*AssetData {
	if strings.HasPrefix(sdk.UnwrapSDKContext(ctx).ChainID(), "diamond") {
		return DefaultMainnetStakeableAssets
	}
	return DefaultTestnetStakeableAssets
}
