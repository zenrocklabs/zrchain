package mint

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/mint/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/mint/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx context.Context, k keeper.Keeper, ic types.InflationCalculationFn) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)

	// fetch stored minter & params
	minter, err := k.Minter.Get(ctx)
	if err != nil {
		k.Logger(ctx).Error("failed to get minter", "error", err)
		return nil
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		k.Logger(ctx).Error("failed to get params", "error", err)
		return nil
	}

	// recalculate inflation rate
	totalStakingSupply, err := k.StakingTokenSupply(ctx)
	if err != nil {
		k.Logger(ctx).Error("failed to get staking token supply", "error", err)
		return nil
	}

	bondedRatio, err := k.BondedRatio(ctx)
	if err != nil {
		k.Logger(ctx).Error("failed to get bonded ratio", "error", err)
		return nil
	}

	mintModuleBalance, err := k.GetMintModuleBalance(ctx)
	if err != nil {
		k.Logger(ctx).Error("failed to get mint module balance", "error", err)
		return nil
	}

	totalRewards, err := k.ClaimTotalRewards(ctx)
	if err != nil {
		k.Logger(ctx).Error("failed to claim total rewards", "error", err)
		return nil
	}

	err = k.DistributeZentpFees(ctx)
	if err != nil {
		k.Logger(ctx).Error("failed to distribute zentp fees", "error", err)
		return nil
	}

	totalRewardsRest, err := k.BaseDistribution(ctx, totalRewards)
	if err != nil {
		k.Logger(ctx).Error("failed to calculate base distribution", "error", err)
		return nil
	}

	totalBondedTokens, err := k.TotalBondedTokens(ctx)
	if err != nil {
		k.Logger(ctx).Error("failed to get total bonded tokens", "error", err)
		return nil
	}

	totalBlockStakingReward, err := k.NextStakingReward(ctx, totalBondedTokens)
	if err != nil {
		k.Logger(ctx).Error("failed to get next staking reward", "error", err)
		return nil
	}

	if totalBlockStakingReward.Amount.GT(totalRewardsRest.Amount) {
		topUpAmount, err := k.CalculateTopUp(ctx, totalBlockStakingReward, totalRewardsRest)
		if err != nil {
			k.Logger(ctx).Error("failed to calculate top up amount", "error", err)
			return nil
		}

		// if totalRewardsRest enough - top up from mint module
		if !topUpAmount.IsZero() || !mintModuleBalance.IsZero() {
			totalRewardsRest = totalBlockStakingReward
		}

		if err := k.CheckModuleBalance(ctx, totalBlockStakingReward); err != nil {
			k.Logger(ctx).Error("failed to check module balance", "error", err)
			return nil
		}
	} else {
		excess, err := k.CalculateExcess(ctx, totalBlockStakingReward, totalRewardsRest)
		if err != nil {
			k.Logger(ctx).Error("failed to calculate excess", "error", err)
			return nil
		}

		if err := k.ExcessDistribution(ctx, excess); err != nil {
			k.Logger(ctx).Error("failed during excess distribution", "error", err)
			return nil
		}
	}

	if err := k.AddCollectedFees(ctx, sdk.NewCoins(totalBlockStakingReward)); err != nil {
		k.Logger(ctx).Error("failed to add collected fees (staking reward)", "error", err)
		return nil
	}

	minter.Inflation = ic(ctx, minter, params, bondedRatio)
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalStakingSupply)
	if err := k.Minter.Set(ctx, minter); err != nil {
		k.Logger(ctx).Error("failed to set minter", "error", err)
		return nil
	}

	// mint coins, update supply
	mintedCoin := minter.BlockProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)

	if err := k.MintCoins(ctx, mintedCoins); err != nil {
		k.Logger(ctx).Error("failed to mint coins", "error", err)
		return nil
	}

	// send the minted coins to the fee collector account
	if err := k.AddCollectedFees(ctx, mintedCoins); err != nil {
		k.Logger(ctx).Error("failed to add collected fees (minted coins)", "error", err)
		return nil
	}

	if mintedCoin.Amount.IsInt64() {
		defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyBondedRatio, bondedRatio.String()),
			sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
			sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
			sdk.NewAttribute(types.AttributeBlockStakingRewards, totalBlockStakingReward.Amount.String()),
			sdk.NewAttribute(types.AttributeTotalFees, totalRewards.Amount.String()),
		),
	)

	return nil
}
