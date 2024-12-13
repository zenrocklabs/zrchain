package mint

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v5/x/mint/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/mint/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx context.Context, k keeper.Keeper, ic types.InflationCalculationFn) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)

	// fetch stored minter & params
	minter, err := k.Minter.Get(ctx)
	if err != nil {
		return err
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	// recalculate inflation rate
	totalStakingSupply, err := k.StakingTokenSupply(ctx)
	if err != nil {
		return err
	}

	bondedRatio, err := k.BondedRatio(ctx)
	if err != nil {
		return err
	}

	mintModuleBalance, err := k.GetMintModuleBalance(ctx)
	if err != nil {
		return err
	}

	totalRewards, err := k.ClaimTotalRewards(ctx)
	if err != nil {
		return err
	}

	totalRewardsRest, err := k.BaseDistribution(ctx, totalRewards)
	if err != nil {
		return err
	}

	totalBondedTokens, err := k.TotalBondedTokens(ctx)
	if err != nil {
		return err
	}

	totalBlockStakingReward, err := k.NextStakingReward(ctx, totalBondedTokens)
	if err != nil {
		return err
	}

	if totalBlockStakingReward.Amount.GT(totalRewardsRest.Amount) {
		topUpAmount, err := k.CalculateTopUp(ctx, totalBlockStakingReward, totalRewardsRest)
		if err != nil {
			return err
		}

		// if totalRewardsRest enough - top up from mint module
		if !topUpAmount.IsZero() || !mintModuleBalance.IsZero() {
			totalRewardsRest = totalBlockStakingReward
		}

		err = k.CheckModuleBalance(ctx, totalBlockStakingReward)
		if err != nil {
			return err
		}
	} else {
		excess, err := k.CalculateExcess(ctx, totalBlockStakingReward, totalRewardsRest)
		if err != nil {
			return err
		}

		err = k.ExcessDistribution(ctx, excess)
		if err != nil {
			return err
		}
	}

	err = k.AddCollectedFees(ctx, sdk.NewCoins(totalBlockStakingReward))
	if err != nil {
		return err
	}

	minter.Inflation = ic(ctx, minter, params, bondedRatio)
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalStakingSupply)
	if err = k.Minter.Set(ctx, minter); err != nil {
		return err
	}

	// mint coins, update supply
	mintedCoin := minter.BlockProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)

	err = k.MintCoins(ctx, mintedCoins)
	if err != nil {
		return err
	}

	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, mintedCoins)
	if err != nil {
		return err
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
