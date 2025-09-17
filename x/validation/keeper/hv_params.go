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

// GetBedrockDefaultValOperAddr returns the HV param for default bedrock validator operator address
// Fallbacks to types.DefaultBedrockValidatorOperAddr if not set.
func (k Keeper) GetBedrockDefaultValOperAddr(ctx context.Context) string {
	// NOTE: We return the default for now if params retrieval fails or field is empty.
	// The proto contains the field; once protobufs are regenerated the HVParams struct will expose it.
	_ , err := k.HVParams.Get(ctx)
	if err != nil {
		return types.DefaultBedrockValidatorOperAddr
	}
	// TODO: read from params.BedrockDefaultValOperAddr when available in generated code
	return types.DefaultBedrockValidatorOperAddr
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

// GetVEJailingEnabled returns whether VE jailing is enabled
func (k Keeper) GetVEJailingEnabled(ctx context.Context) bool {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return types.DefaultVEJailingEnabled
	}
	return params.VEJailingEnabled
}

// GetVEJailDurationMinutes returns the VE jail duration in minutes
func (k Keeper) GetVEJailDurationMinutes(ctx context.Context) int64 {
	params, err := k.HVParams.Get(ctx)
	if err != nil || params.VEJailDurationMinutes <= 0 {
		return types.DefaultVEJailDurationMinutes
	}
	return params.VEJailDurationMinutes
}

// GetVEWindowSize returns the VE window size for mismatch tracking
func (k Keeper) GetVEWindowSize(ctx context.Context) int64 {
	params, err := k.HVParams.Get(ctx)
	if err != nil || params.VEWindowSize <= 0 {
		return types.DefaultVEWindowSize
	}
	return params.VEWindowSize
}

// GetVEJailThreshold returns the VE jail threshold for number of mismatches
func (k Keeper) GetVEJailThreshold(ctx context.Context) int64 {
	params, err := k.HVParams.Get(ctx)
	if err != nil || params.VEJailThreshold <= 0 {
		return types.DefaultVEJailThreshold
	}
	return params.VEJailThreshold
}
