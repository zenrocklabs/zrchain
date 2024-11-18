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

	// fetch stored minter & params
	minter, err := k.Minter.Get(ctx)
	if err != nil {
		return err
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	// TODO - create tests
	// TODO -
	// TODO - adjust events
	// TODO - remove legacy minting mechanism

	// DONE - add TotalBondedTokens
	// DONE - add NextStakingReward based on annual staking yield
	// DONE - add amount to get distributed to stakers
	// DONE - set new params for protocol wallet and distribution amounts
	// DONE - receive keyring fees as rewards
	// DONE - burn percentage
	// DONE - send to protocol wallet
	// DONE - handle if rewards are enough to cover the staking rewards
	// DONE - implement excess reward mechanism
	// DONE - handle top-up from protocol wallet to cover staking rewards

	// recalculate inflation rate
	totalStakingSupply, err := k.StakingTokenSupply(ctx)
	if err != nil {
		return err
	}

	bondedRatio, err := k.BondedRatio(ctx)
	if err != nil {
		return err
	}

	keyringRewards, err := k.ClaimKeyringFees(ctx)
	if err != nil {
		return err
	}

	// TODO - remove
	fmt.Println("keyringRewards: ", keyringRewards)

	keyringRewardsRest, err := k.BaseDistribution(ctx, keyringRewards)
	if err != nil {
		return err
	}

	// TODO - remove
	fmt.Println("keyringRewardsRest: ", keyringRewardsRest)

	totalBondedTokens, err := k.TotalBondedTokens(ctx)
	if err != nil {
		return err
	}

	// TODO - remove
	fmt.Println("totalBondedTokens: ", totalBondedTokens)

	totalBlockStakingReward, err := k.NextStakingReward(ctx, totalBondedTokens)
	if err != nil {
		return err
	}

	// TODO - remove
	fmt.Println("totalBlockStakingReward: ", totalBlockStakingReward)

	if totalBlockStakingReward.Amount.GT(keyringRewardsRest.Amount) {
		topUpAmount, err := k.CalculateTopUp(ctx, totalBlockStakingReward, keyringRewardsRest)
		if err != nil {
			return err
		}

		// TODO - remove
		fmt.Println("topUpAmount: ", topUpAmount)

		err = k.TopUpKeyringRewards(ctx, topUpAmount)
		if err != nil {
			return err
		}

		err = k.CheckModuleBalance(ctx, totalBlockStakingReward)
		if err != nil {
			return err
		}
	} else {
		excess, err := k.CalculateExcess(ctx, totalBlockStakingReward, keyringRewardsRest)
		if err != nil {
			return err
		}

		// TODO - remove
		fmt.Println("excess: ", excess)

		err = k.AdditionalBurn(ctx, excess)
		if err != nil {
			return err
		}

		// TODO - remove
		fmt.Println("excess after burn: ", excess)

		err = k.AdditionalMpcRewards(ctx, excess)
		if err != nil {
			return err
		}

		// TODO - remove
		fmt.Println("excess after mpc rewards: ", excess)

		err = k.AdditionalStakingRewards(ctx, excess)
		if err != nil {
			return err
		}

		// TODO - remove
		fmt.Println("excess after staking rewards: ", excess)
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
		),
	)

	return nil
}
