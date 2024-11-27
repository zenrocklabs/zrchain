package mint

import (
	"context"
	"fmt"

	"github.com/Zenrock-Foundation/zrchain/v5/x/mint/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/mint/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx context.Context, k keeper.Keeper, ic types.InflationCalculationFn) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)

	perms := k.GetModuleAccountPerms(ctx)
	fmt.Println("Mint Module Permissions at BeginBlocker:", perms)

	// fetch stored minter & params
	minter, err := k.Minter.Get(ctx)
	if err != nil {
		return err
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	// TODO remove
	fmt.Println("---")

	// TODO - create tests
	// TODO - adjust events
	// TODO - remove legacy minting mechanism

	// recalculate inflation rate
	totalStakingSupply, err := k.StakingTokenSupply(ctx)
	if err != nil {
		return err
	}

	bondedRatio, err := k.BondedRatio(ctx)
	if err != nil {
		return err
	}

	mintModuleBalance, err := k.CheckMintModuleBalance(ctx)
	if err != nil {
		return err
	}

	// TODO - remove
	fmt.Printf("mint module balance:\t\t%v\n", mintModuleBalance)

	totalRewards, err := k.ClaimTotalRewards(ctx)
	if err != nil {
		return err
	}

	// TODO - remove
	fmt.Printf("total rewards:\t\t\t%v\n", totalRewards)

	totalRewardsRest, err := k.BaseDistribution(ctx, totalRewards)
	if err != nil {
		return err
	}

	// TODO - remove
	fmt.Printf("rewards after burn&pw:\t\t%v\n", totalRewardsRest)

	totalBondedTokens, err := k.TotalBondedTokens(ctx)
	if err != nil {
		return err
	}

	// TODO - remove
	fmt.Printf("total bonded tokens:\t\t%v\n", totalBondedTokens)

	totalBlockStakingReward, err := k.NextStakingReward(ctx, totalBondedTokens)
	if err != nil {
		return err
	}

	// TODO - remove
	fmt.Printf("total staking rewards:\t\t%v\n", totalBlockStakingReward)

	if totalBlockStakingReward.Amount.GT(totalRewardsRest.Amount) {
		topUpAmount, err := k.CalculateTopUp(ctx, totalBlockStakingReward, totalRewardsRest)
		if err != nil {
			return err
		}

		// TODO - remove
		fmt.Printf("top-up amount:\t\t\t%v\n", topUpAmount)

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
		),
	)

	return nil
}
