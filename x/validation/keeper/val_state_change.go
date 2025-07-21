package keeper

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	m "math"
	"math/big"
	"sort"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	gogotypes "github.com/cosmos/gogoproto/types"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
)

// BlockValidatorUpdates calculates the ValidatorUpdates for the current block
// Called in each EndBlock
func (k Keeper) BlockValidatorUpdates(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Calculate validator set changes.
	//
	// NOTE: ApplyAndReturnValidatorSetUpdates has to come before
	// UnbondAllMatureValidatorQueue.
	// This fixes a bug when the unbonding period is instant (is the case in
	// some of the tests). The test expected the validator to be completely
	// unbonded after the Endblocker (go from Bonded -> Unbonding during
	// ApplyAndReturnValidatorSetUpdates and then Unbonding -> Unbonded during
	// UnbondAllMatureValidatorQueue).
	validatorUpdates, err := k.ApplyAndReturnValidatorSetUpdates(ctx)
	if err != nil {
		return nil, err
	}

	// Provision the AVS staking rewards for the current block
	if err := k.provisionAVSValidatorRewards(sdkCtx); err != nil {
		return nil, err
	}
	if err := k.provisionAVSDelegatorRewardsAndCommissions(sdkCtx); err != nil {
		return nil, err
	}

	// unbond all mature validators from the unbonding queue
	if err = k.UnbondAllMatureValidators(ctx); err != nil {
		return nil, err
	}

	// Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds, err := k.DequeueAllMatureUBDQueue(ctx, sdkCtx.BlockHeader().Time)
	if err != nil {
		return nil, err
	}

	for _, dvPair := range matureUnbonds {
		addr, err := k.validatorAddressCodec.StringToBytes(dvPair.ValidatorAddress)
		if err != nil {
			return nil, err
		}
		delegatorAddress, err := k.authKeeper.AddressCodec().StringToBytes(dvPair.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		balances, err := k.CompleteUnbonding(ctx, delegatorAddress, addr)
		if err != nil {
			continue
		}

		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbonding,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, dvPair.ValidatorAddress),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvPair.DelegatorAddress),
			),
		)
	}

	// Remove all mature redelegations from the red queue.
	matureRedelegations, err := k.DequeueAllMatureRedelegationQueue(ctx, sdkCtx.BlockHeader().Time)
	if err != nil {
		return nil, err
	}

	for _, dvvTriplet := range matureRedelegations {
		valSrcAddr, err := k.validatorAddressCodec.StringToBytes(dvvTriplet.ValidatorSrcAddress)
		if err != nil {
			return nil, err
		}
		valDstAddr, err := k.validatorAddressCodec.StringToBytes(dvvTriplet.ValidatorDstAddress)
		if err != nil {
			return nil, err
		}
		delegatorAddress, err := k.authKeeper.AddressCodec().StringToBytes(dvvTriplet.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		balances, err := k.CompleteRedelegation(
			ctx,
			delegatorAddress,
			valSrcAddr,
			valDstAddr,
		)
		if err != nil {
			continue
		}

		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteRedelegation,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvvTriplet.DelegatorAddress),
				sdk.NewAttribute(types.AttributeKeySrcValidator, dvvTriplet.ValidatorSrcAddress),
				sdk.NewAttribute(types.AttributeKeyDstValidator, dvvTriplet.ValidatorDstAddress),
			),
		)
	}

	return validatorUpdates, nil
}

// provisionAVSValidatorRewards calculates the rewards for all the bonded Zenrock validators
// which have received AVS delegations and provisions them in the rewards pool for later withdrawal.
func (k Keeper) provisionAVSValidatorRewards(ctx sdk.Context) error {
	return k.IterateBondedZenrockValidatorsByPower(ctx, func(index int64, validator types.ValidatorHV) error {
		rewardsPerBlock := k.calculateRewardsPerBlock(ctx, validator.TokensAVS)
		amount := rewardsPerBlock.TruncateInt()
		if err := k.addRewards(ctx, validator.OperatorAddress, amount); err != nil {
			k.Logger(ctx).Error("couldn't set rewards for validator", "address", validator.OperatorAddress, "error", err)
			return err
		}
		return nil
	})
}

// provisionAVSDelegatorRewardsAndCommissions calculates the rewards for all the AVS delegators
// and provisions them in the rewards pool for later withdrawal. It also calculates the commissions
// owed to the validators from the delegation rewards and provisions them in the rewards pool.
func (k Keeper) provisionAVSDelegatorRewardsAndCommissions(ctx sdk.Context) error {
	return k.AVSDelegations.Walk(ctx, nil, func(key collections.Pair[string, string], delegated math.Int) (bool, error) {
		validatorAddr, delegatorAddr := key.K1(), key.K2()
		validator, err := k.GetZenrockValidatorFromBech32(ctx, validatorAddr)
		if err != nil {
			// This can error if an AVS operator delegates to a validator address that doesn't exist.
			// Since it is possible for this to error in normal circumstances, we shouldn't return the error or it will halt the chain.
			// In this case we must simply log the error and continue.
			k.Logger(ctx).Debug("couldn't get validator", "address", validatorAddr, "error", err)
			return false, nil
		}

		rewardsPerBlock := k.calculateRewardsPerBlock(ctx, delegated)
		if rewardsPerBlock.IsZero() {
			return false, nil
		}
		commission := rewardsPerBlock.Mul(validator.Commission.Rate).TruncateInt()
		delegatorReward := rewardsPerBlock.Sub(commission.ToLegacyDec()).TruncateInt()

		if err := k.addRewards(ctx, delegatorAddr, delegatorReward); err != nil {
			return true, fmt.Errorf("couldn't set rewards for delegator %s: %w", delegatorAddr, err)
		}

		if err := k.addRewards(ctx, validatorAddr, commission); err != nil {
			return true, fmt.Errorf("couldn't set rewards for validator %s: %w", validatorAddr, err)
		}

		return false, nil
	})
}

// calculateRewardsPerBlock calculates the rewards per block for a given amount of staked AVS tokens.
// This uses the current AVS rewards rate set in the module parameters via governance.
func (k Keeper) calculateRewardsPerBlock(ctx sdk.Context, tokens math.Int) math.LegacyDec {
	secondsPerYear := math.LegacyNewDec(365 * 24 * 60 * 60)
	blockTime := math.LegacyNewDec(k.GetBlockTime(ctx))
	if blockTime.IsZero() {
		return math.LegacyZeroDec()
	}
	blocksPerYear := secondsPerYear.Quo(blockTime)
	rewardsPerYear := k.GetAVSRewardsRate(ctx).Mul(tokens.ToLegacyDec())
	return rewardsPerYear.Quo(blocksPerYear)
}

// Helper function to get current rewards or zero if not found
func (k Keeper) getCurrentRewards(ctx sdk.Context, addr string) (math.Int, error) {
	currentRewards, err := k.AVSRewardsPool.Get(ctx, addr)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// If no rewards are found, return zero (since returning an error here halts the chain)
			return math.ZeroInt(), nil
		}
		// Any other error is critical and should be returned
		return math.ZeroInt(), err
	}
	return currentRewards, nil
}

// Helper function to add rewards to the pool
func (k Keeper) addRewards(ctx sdk.Context, addr string, amount math.Int) error {
	currentRewards, err := k.getCurrentRewards(ctx, addr)
	if err != nil {
		return err
	}
	return k.AVSRewardsPool.Set(ctx, addr, currentRewards.Add(amount))
}

// ApplyAndReturnValidatorSetUpdates applies and return accumulated updates to the bonded validator set. Also,
// * Updates the active valset as keyed by LastValidatorPowerKey.
// * Updates the total power as keyed by LastTotalPowerKey.
// * Updates validator status' according to updated powers.
// * Updates the fee pool bonded vs not-bonded tokens.
// * Updates relevant indices.
// It gets called once after genesis, another time maybe after genesis transactions,
// then once at every EndBlock.
//
// CONTRACT: Only validators with non-zero power or zero-power that were bonded
// at the previous block height or were removed from the validator set entirely
// are returned to CometBFT.
func (k Keeper) ApplyAndReturnValidatorSetUpdates(ctx context.Context) (updates []abci.ValidatorUpdate, err error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	maxValidators := params.MaxValidators
	powerReduction := k.PowerReduction(ctx)
	totalPowerCalculated := math.ZeroInt()
	amtFromBondedToNotBonded, amtFromNotBondedToBonded := math.ZeroInt(), math.ZeroInt()

	last, err := k.getLastValidatorsByAddr(ctx)
	if err != nil {
		return nil, err
	}

	// Check and jail validators with excessive mismatched vote extensions
	if err := k.checkAndJailValidatorsForMismatchedVoteExtensions(ctx); err != nil {
		return nil, fmt.Errorf("failed to check validators for vote extension mismatches: %w", err)
	}

	// Cleanup old mismatch count records every window size blocks to prevent storage bloat
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	windowSize := k.GetVEWindowSize(ctx)
	if sdkCtx.BlockHeight()%windowSize == 0 {
		if err := k.cleanupOldMismatchCounts(ctx, sdkCtx.BlockHeight()); err != nil {
			k.Logger(ctx).Error("Failed to cleanup old mismatch counts", "error", err)
			// Don't return error here as it's not critical for validator set updates
		}
	}

	iterator, err := k.ValidatorsPowerStoreIterator(ctx)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	// Fetch all stakeable asset prices and data once.
	stakeableAssetsWithPrices, err := k.GetStakeableAssetPrices(ctx)
	if err != nil {
		return nil, err
	}

	// Determine if any prices are valid (non-zero).
	pricesAreValid := false
	for _, asset := range stakeableAssetsWithPrices {
		if !asset.PriceUSD.IsZero() {
			pricesAreValid = true
			break
		}
	}

	// Identify native (ROCK) and primary AVS (BTC) asset data from the fetched list.
	// This assumes fixed asset types for native and primary AVS.
	// A more dynamic system might involve looking up asset types based on validator's specific holdings.
	var nativeAssetData *types.AssetData
	var primaryAVSAssetData *types.AssetData // e.g., for BTC

	for _, asset := range stakeableAssetsWithPrices {
		if asset.Asset == types.Asset_ROCK { // Assuming types.Asset_ROCK is your native asset enum/const
			nativeAssetData = asset
		}
		if asset.Asset == types.Asset_BTC { // Assuming types.Asset_BTC is your primary AVS asset enum/const
			primaryAVSAssetData = asset
		}
		// Add more AVS asset types if a validator can hold multiple, and adjust logic below.
	}

	for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {
		valAddr := sdk.ValAddress(iterator.Value())
		validator := k.mustGetValidator(ctx, valAddr)

		if validator.Jailed {
			panic("should never retrieve a jailed validator from the power store")
		}

		// If we get to a zero-power validator (which we don't bond),
		// there are no more possible bonded validators.
		// This check relies on PotentialConsensusPower using the *native* tokens primarily.
		// If AVS tokens could make a zero-native-token validator powerful, this check might need adjustment.
		if validator.PotentialConsensusPower(powerReduction) == 0 && validator.TokensAVS.IsZero() {
			break
		}

		switch {
		case validator.IsUnbonded():
			validator, err = k.unbondedToBonded(ctx, validator)
			if err != nil {
				return nil, err
			}
			amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(validator.GetTokens())
		case validator.IsUnbonding():
			validator, err = k.unbondingToBonded(ctx, validator)
			if err != nil {
				return nil, err
			}
			amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(validator.GetTokens())
		case validator.IsBonded():
			// no state change
		default:
			panic("unexpected validator status")
		}

		valAddrStr, err := k.validatorAddressCodec.BytesToString(valAddr)
		if err != nil {
			return nil, err
		}
		oldPowerBytes, found := last[valAddrStr]

		// Use the new helper function to calculate power
		// Pass validator.TokensAVS directly, assuming it corresponds to primaryAVSAssetData (e.g. BTC)
		finalConsensusPower, nativeValueDecimal, avsValueDecimal := k.calculateValidatorPowerComponents(
			validator,
			powerReduction,
			nativeAssetData,
			primaryAVSAssetData, // Pass data for the AVS tokens the validator holds
			validator.TokensAVS, // Pass the amount of those AVS tokens
			pricesAreValid,
		)

		consAddr, err := validator.GetConsAddr()
		if err != nil {
			return nil, err
		}

		configuredAVSAssetType := "AVS_STAKE_DISABLED" // Default if AVS asset isn't configured/found
		if primaryAVSAssetData != nil {
			configuredAVSAssetType = primaryAVSAssetData.GetAsset().String()
		}

		k.Logger(ctx).Debug(fmt.Sprintf(
			"\nvalidator: %s | %s\ntoken stake: native_units=%d, avs_raw_units=%s (for %s)\nstake value: native_contrib=%s, avs_contrib=%s, total_power_calc=%d",
			valAddrStr, sdk.ConsAddress(consAddr).String(),
			validator.ConsensusPower(powerReduction), validator.TokensAVS.String(), configuredAVSAssetType,
			nativeValueDecimal.String(), avsValueDecimal.String(), finalConsensusPower,
		))

		newPowerBytes := k.cdc.MustMarshal(&gogotypes.Int64Value{Value: finalConsensusPower})

		if !found || !bytes.Equal(oldPowerBytes, newPowerBytes) {
			// The ABCIValidatorUpdate's Power field should be finalConsensusPower.
			// The validator.ABCIValidatorUpdate method takes the final calculated power directly.
			updates = append(updates, validator.ABCIValidatorUpdate(powerReduction, finalConsensusPower))

			if err = k.SetLastValidatorPower(ctx, valAddr, finalConsensusPower); err != nil {
				return nil, err
			}
		}

		delete(last, valAddrStr)
		count++
		totalPowerCalculated = totalPowerCalculated.Add(math.NewInt(finalConsensusPower))
	}

	noLongerBonded, err := sortNoLongerBonded(last, k.validatorAddressCodec)
	if err != nil {
		return nil, err
	}

	for _, valAddrBytes := range noLongerBonded {
		validator := k.mustGetValidator(ctx, sdk.ValAddress(valAddrBytes))
		validator, err = k.bondedToUnbonding(ctx, validator)
		if err != nil {
			return nil, err
		}

		amtFromBondedToNotBonded = amtFromBondedToNotBonded.Add(validator.GetTokens())

		// Use sdk.ValAddress for DeleteLastValidatorPower if it expects that
		if err = k.DeleteLastValidatorPower(ctx, sdk.ValAddress(valAddrBytes)); err != nil { // Corrected to pass sdk.ValAddress if this was the intention from the comment. The original code used valAddrBytes (a []byte) which implicitly converts to sdk.ValAddress, but this is more explicit.
			return nil, err
		}
		updates = append(updates, validator.ABCIValidatorUpdateZero())
	}

	// Update the pools based on the recent updates in the validator set:
	switch {
	case amtFromNotBondedToBonded.GT(amtFromBondedToNotBonded):
		if err = k.notBondedTokensToBonded(sdk.UnwrapSDKContext(ctx), amtFromNotBondedToBonded.Sub(amtFromBondedToNotBonded)); err != nil {
			return nil, err
		}
	case amtFromNotBondedToBonded.LT(amtFromBondedToNotBonded):
		if err = k.bondedTokensToNotBonded(sdk.UnwrapSDKContext(ctx), amtFromBondedToNotBonded.Sub(amtFromNotBondedToBonded)); err != nil {
			return nil, err
		}
	}

	if len(updates) > 0 {
		if err = k.SetLastTotalPower(ctx, totalPowerCalculated); err != nil {
			return nil, err
		}
	}
	if err = k.SetValidatorUpdates(ctx, updates); err != nil { // This was removed in SDK v0.47+, might not be needed
		return nil, err
	}

	return updates, nil
}

func (k Keeper) GetStakeableAssetPrices(ctx context.Context) ([]*types.AssetData, error) {
	stakeableAssets := k.GetStakeableAssets(ctx)

	for _, asset := range stakeableAssets {
		assetPrice, err := k.AssetPrices.Get(ctx, asset.Asset)
		if err != nil {
			if !errors.Is(err, collections.ErrNotFound) {
				return nil, err
			}
			asset.PriceUSD = math.LegacyZeroDec()
			if err := k.AssetPrices.Set(ctx, asset.Asset, asset.PriceUSD); err != nil {
				return nil, err
			}
		} else {
			asset.PriceUSD = assetPrice
		}
	}

	return stakeableAssets, nil
}

// calculateValidatorPowerComponents computes the native and AVS power contributions
// and the total truncated consensus power for a validator, using full decimal precision for AVS token values.
func (k Keeper) calculateValidatorPowerComponents(
	validator types.ValidatorHV, // Used for validator.ConsensusPower
	powerReduction math.Int,
	nativeAssetData *types.AssetData,
	avsAssetData *types.AssetData, // Contains PriceUSD and Precision for the AVS token
	avsTokens math.Int, // Raw amount of AVS tokens (e.g., satoshis)
	pricesAreValid bool,
) (totalConsensusPower int64, nativePowerContribution math.LegacyDec, avsPowerContribution math.LegacyDec) {

	nativePowerContribution = math.LegacyZeroDec()
	avsPowerContribution = math.LegacyZeroDec()

	// Calculate Native Power Contribution
	nativeConsensusPowerUnits := validator.ConsensusPower(powerReduction)
	if nativeAssetData != nil {
		if pricesAreValid {
			nativePowerContribution = nativeAssetData.PriceUSD.MulInt64(nativeConsensusPowerUnits)
		} else {
			nativePowerContribution = math.LegacyNewDec(nativeConsensusPowerUnits)
		}
	} else if !pricesAreValid {
		nativePowerContribution = math.LegacyNewDec(nativeConsensusPowerUnits)
	}

	// Calculate AVS Power Contribution with full decimal precision
	if avsAssetData != nil && avsTokens.IsPositive() {
		if pricesAreValid && !avsAssetData.PriceUSD.IsZero() {
			// Convert raw AVS tokens (math.Int) to math.LegacyDec
			avsTokensDec := math.LegacyNewDecFromInt(avsTokens)

			// Calculate the divisor (10^assetPrecision) as math.LegacyDec
			if avsAssetData.Precision > 0 { // Avoid division by zero if precision is 0 (10^0 = 1 handled by next lines)
				powerOf10DenominatorBigInt := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(avsAssetData.Precision)), nil)
				divisorDec := math.LegacyNewDecFromBigInt(powerOf10DenominatorBigInt)

				if !divisorDec.IsZero() {
					// Calculate whole AVS units as math.LegacyDec (e.g., 1.18170108 for BTC)
					wholeAVSUnitsDec := avsTokensDec.Quo(divisorDec)
					avsPowerContribution = avsAssetData.PriceUSD.Mul(wholeAVSUnitsDec)
				} else {
					// This case should ideally not be reached if precision > 0
					avsPowerContribution = math.LegacyZeroDec()
				}
			} else { // asset.Precision is 0, so 10^0 = 1. Treat tokens as whole units.
				avsPowerContribution = avsAssetData.PriceUSD.Mul(avsTokensDec)
			}
		}
		// If pricesAreValid is false, or avsAssetData.PriceUSD is zero, avsPowerContribution remains ZeroDec
	}

	calculatedTotalDecimal := nativePowerContribution.Add(avsPowerContribution)

	// Prevent panic from TruncateInt64 and ensure power is not negative.
	// m.MaxInt64 comes from the standard "math" package (aliased as 'm').
	// math.LegacyNewDec and math.LegacyZeroDec come from "cosmossdk.io/math".
	maxInt64AsDec := math.LegacyNewDec(m.MaxInt64)

	if calculatedTotalDecimal.GT(maxInt64AsDec) {
		totalConsensusPower = m.MaxInt64
	} else if calculatedTotalDecimal.IsNegative() {
		totalConsensusPower = 0
	} else {
		// It's safe to truncate now, and it's not negative.
		totalConsensusPower = calculatedTotalDecimal.TruncateInt64()
	}

	// The calling function `ApplyAndReturnValidatorSetUpdates` filters out validators
	// where (validator.PotentialConsensusPower(powerReduction) == 0 && validator.TokensAVS.IsZero()).
	// If this function is called, the validator is considered to have some basis for staking.
	// Therefore, if the calculated power truncates to 0 (and was not negative), it should be set to 1,
	// unless it represents a truly valueless validator.
	// This logic applies if calculatedTotalDecimal was >= 0.
	if totalConsensusPower == 0 && !calculatedTotalDecimal.IsNegative() {
		if calculatedTotalDecimal.IsZero() && validator.TokensNative.IsZero() && !avsTokens.IsPositive() {
			// Truly valueless or original was exactly zero, power remains 0.
			// totalConsensusPower is already 0.
		} else {
			// Had some small positive value that truncated to 0, or was zero but had some tokens, gets power 1.
			totalConsensusPower = 1
		}
	}
	// If calculatedTotalDecimal was negative, totalConsensusPower was set to 0 and remains 0 due to the
	// '!calculatedTotalDecimal.IsNegative()' condition above.

	return totalConsensusPower, nativePowerContribution, avsPowerContribution
}

// Validator state transitions

func (k Keeper) bondedToUnbonding(ctx context.Context, validator types.ValidatorHV) (types.ValidatorHV, error) {
	if !validator.IsBonded() {
		panic(fmt.Sprintf("bad state transition bondedToUnbonding, validator: %v\n", validator))
	}

	return k.BeginUnbondingValidator(ctx, validator)
}

func (k Keeper) unbondingToBonded(ctx context.Context, validator types.ValidatorHV) (types.ValidatorHV, error) {
	if !validator.IsUnbonding() {
		panic(fmt.Sprintf("bad state transition unbondingToBonded, validator: %v\n", validator))
	}

	return k.bondValidator(ctx, validator)
}

func (k Keeper) unbondedToBonded(ctx context.Context, validator types.ValidatorHV) (types.ValidatorHV, error) {
	if !validator.IsUnbonded() {
		panic(fmt.Sprintf("bad state transition unbondedToBonded, validator: %v\n", validator))
	}

	return k.bondValidator(ctx, validator)
}

// UnbondingToUnbonded switches a validator from unbonding state to unbonded state
func (k Keeper) UnbondingToUnbonded(ctx context.Context, validator types.ValidatorHV) (types.ValidatorHV, error) {
	if !validator.IsUnbonding() {
		return types.ValidatorHV{}, fmt.Errorf("bad state transition unbondingToUnbonded, validator: %v", validator)
	}

	return k.completeUnbondingValidator(ctx, validator)
}

// Vote extension mismatch detection configuration

// checkAndJailValidatorsForMismatchedVoteExtensions checks all bonded validators
// and jails those who have submitted mismatched vote extensions for at least 100
// out of the last 200 blocks. This optimized version uses a sliding window counter
// to achieve O(V) time complexity instead of O(V × B × M).
func (k Keeper) checkAndJailValidatorsForMismatchedVoteExtensions(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentHeight := sdkCtx.BlockHeight()

	// Only check if we have at least the full window size of blocks of data
	if currentHeight < k.GetVEWindowSize(ctx) {
		return nil
	}

	// Get all bonded validators - O(V)
	bonded := make([]types.ValidatorHV, 0)
	err := k.IterateBondedZenrockValidatorsByPower(ctx, func(_ int64, validator types.ValidatorHV) error {
		if !validator.Jailed {
			bonded = append(bonded, validator)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// For each bonded validator, check their mismatch count - O(V)
	for _, validator := range bonded {
		consAddr, err := validator.GetConsAddr()
		if err != nil {
			k.Logger(ctx).Error("Failed to get consensus address for validator", "validator", validator.OperatorAddress, "error", err)
			continue
		}

		validatorHexAddr := hex.EncodeToString(consAddr)

		// Get the current mismatch count for this validator - O(1)
		mismatchCount, err := k.ValidatorMismatchCounts.Get(ctx, validatorHexAddr)
		if err != nil {
			// If no record exists, validator has no mismatches
			continue
		}

		// If the validator has reached the jail threshold or more mismatches, check if jailing is enabled
		jailThreshold := k.GetVEJailThreshold(ctx)
		if mismatchCount.TotalCount >= uint32(jailThreshold) {
			// Check if VE jailing is enabled
			if !k.GetVEJailingEnabled(ctx) {
				k.Logger(ctx).Info(
					"validator sidecar desynced (jailing disabled)",
					"validator", validator.OperatorAddress,
					"consensus_addr", sdk.ConsAddress(consAddr).String(),
					"mismatch_count", mismatchCount.TotalCount,
					"blocks_checked", k.GetVEWindowSize(ctx),
				)
				continue
			}

			k.Logger(ctx).Info(
				"validator sidecar desynced - jailing validator",
				"validator", validator.OperatorAddress,
				"consensus_addr", sdk.ConsAddress(consAddr).String(),
				"mismatch_count", mismatchCount.TotalCount,
				"blocks_checked", k.GetVEWindowSize(ctx),
			)

			if err := k.jailValidator(ctx, validator); err != nil {
				k.Logger(ctx).Error(
					"Failed to jail validator for mismatched vote extensions",
					"validator", validator.OperatorAddress,
					"error", err,
				)
				continue
			}

			// Get and update signing info to set jail duration based on parameter
			signInfo, err := k.slashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
			if err != nil {
				k.Logger(ctx).Error(
					"Failed to get validator signing info",
					"validator", validator.OperatorAddress,
					"error", err,
				)
				continue
			}

			// Calculate jail duration from parameter (convert minutes to duration)
			jailDurationMinutes := k.GetVEJailDurationMinutes(ctx)
			jailDuration := time.Duration(jailDurationMinutes) * time.Minute
			signInfo.JailedUntil = sdkCtx.BlockHeader().Time.Add(jailDuration)

			// Debug logging for JailedUntil troubleshooting
			k.Logger(ctx).Info(
				"Setting JailedUntil for validator",
				"validator", validator.OperatorAddress,
				"consensus_addr", sdk.ConsAddress(consAddr).String(),
				"current_time", sdkCtx.BlockHeader().Time.String(),
				"jail_duration_minutes", jailDurationMinutes,
				"jail_duration", jailDuration.String(),
				"jailed_until", signInfo.JailedUntil.String(),
				"start_height", signInfo.StartHeight,
			)

			if err := k.slashingKeeper.SetValidatorSigningInfo(ctx, consAddr, signInfo); err != nil {
				k.Logger(ctx).Error(
					"Failed to set validator signing info",
					"validator", validator.OperatorAddress,
					"error", err,
				)
				continue
			}

			// Clear the mismatch store for this validator
			if err := k.ValidatorMismatchCounts.Remove(ctx, validatorHexAddr); err != nil {
				k.Logger(ctx).Error(
					"Failed to remove mismatch count record for jailed validator",
					"validator", validator.OperatorAddress,
					"error", err,
				)
				continue
			}

			// Emit an event for the jailing
			sdkCtx.EventManager().EmitEvent(
				sdk.NewEvent(
					"validator_jailed_vote_extension_mismatch",
					sdk.NewAttribute("validator", validator.OperatorAddress),
					sdk.NewAttribute("consensus_address", sdk.ConsAddress(consAddr).String()),
					sdk.NewAttribute("mismatch_count", fmt.Sprintf("%d", mismatchCount.TotalCount)),
					sdk.NewAttribute("blocks_checked", fmt.Sprintf("%d", k.GetVEWindowSize(ctx))),
					sdk.NewAttribute("jail_duration_minutes", fmt.Sprintf("%d", jailDurationMinutes)),
					sdk.NewAttribute("jailed_until", signInfo.JailedUntil.String()),
				),
			)
		}
	}

	return nil
}

// cleanupOldMismatchCounts removes mismatch count records for validators that have
// no mismatches in the current window. This should be called periodically to prevent
// storage bloat.
func (k Keeper) cleanupOldMismatchCounts(ctx context.Context, currentHeight int64) error {
	windowStart := currentHeight - k.GetVEWindowSize(ctx) + 1

	// Iterate through all mismatch count records
	var toDelete []string
	err := k.ValidatorMismatchCounts.Walk(ctx, nil, func(validatorHexAddr string, mismatchCount types.ValidatorMismatchCount) (bool, error) {
		// Remove blocks outside the window
		validBlocks := make([]int64, 0, len(mismatchCount.MismatchBlocks))
		for _, block := range mismatchCount.MismatchBlocks {
			if block >= windowStart {
				validBlocks = append(validBlocks, block)
			}
		}

		if len(validBlocks) == 0 {
			// No mismatches in the current window, mark for deletion
			toDelete = append(toDelete, validatorHexAddr)
		} else if len(validBlocks) != len(mismatchCount.MismatchBlocks) {
			// Update the record with cleaned blocks
			mismatchCount.MismatchBlocks = validBlocks
			mismatchCount.TotalCount = uint32(len(validBlocks))
			if err := k.ValidatorMismatchCounts.Set(ctx, validatorHexAddr, mismatchCount); err != nil {
				return false, err
			}
		}

		return false, nil // continue iteration
	})

	if err != nil {
		return err
	}

	// Delete records with no recent mismatches
	for _, validatorHexAddr := range toDelete {
		if err := k.ValidatorMismatchCounts.Remove(ctx, validatorHexAddr); err != nil {
			k.Logger(ctx).Error("Failed to remove old mismatch count record", "validator", validatorHexAddr, "error", err)
		}
	}

	return nil
}

// send a validator to jail
func (k Keeper) jailValidator(ctx context.Context, validator types.ValidatorHV) error {
	if validator.Jailed {
		return types.ErrValidatorJailed.Wrapf("cannot jail already jailed validator, validator: %v", validator)
	}

	validator.Jailed = true
	if err := k.SetValidator(ctx, validator); err != nil {
		return err
	}

	return k.DeleteValidatorByPowerIndex(ctx, validator)
}

// remove a validator from jail
func (k Keeper) unjailValidator(ctx context.Context, validator types.ValidatorHV) error {
	if !validator.Jailed {
		return fmt.Errorf("cannot unjail already unjailed validator, validator: %v", validator)
	}

	validator.Jailed = false
	if err := k.SetValidator(ctx, validator); err != nil {
		return err
	}

	return k.SetValidatorByPowerIndex(ctx, validator)
}

// perform all the store operations for when a validator status becomes bonded
func (k Keeper) bondValidator(ctx context.Context, validator types.ValidatorHV) (types.ValidatorHV, error) {
	// delete the validator by power index, as the key will change
	if err := k.DeleteValidatorByPowerIndex(ctx, validator); err != nil {
		return validator, err
	}

	validator = validator.UpdateStatus(types.Bonded)

	// save the now bonded validator record to the two referenced stores
	if err := k.SetValidator(ctx, validator); err != nil {
		return validator, err
	}

	if err := k.SetValidatorByPowerIndex(ctx, validator); err != nil {
		return validator, err
	}

	// delete from queue if present
	if err := k.DeleteValidatorQueue(ctx, validator); err != nil {
		return validator, err
	}

	// trigger hook
	consAddr, err := validator.GetConsAddr()
	if err != nil {
		return validator, err
	}

	str, err := k.validatorAddressCodec.StringToBytes(validator.GetOperator())
	if err != nil {
		return validator, err
	}

	if err := k.Hooks().AfterValidatorBonded(ctx, consAddr, str); err != nil {
		return validator, err
	}

	return validator, err
}

// BeginUnbondingValidator performs all the store operations for when a validator begins unbonding
func (k Keeper) BeginUnbondingValidator(ctx context.Context, validator types.ValidatorHV) (types.ValidatorHV, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return validator, err
	}

	// delete the validator by power index, as the key will change
	if err = k.DeleteValidatorByPowerIndex(ctx, validator); err != nil {
		return validator, err
	}

	// sanity check
	if validator.Status != types.Bonded {
		panic(fmt.Sprintf("should not already be unbonded or unbonding, validator: %v\n", validator))
	}

	id, err := k.IncrementUnbondingID(ctx)
	if err != nil {
		return validator, err
	}

	validator = validator.UpdateStatus(types.Unbonding)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// set the unbonding completion time and completion height appropriately
	validator.UnbondingTime = sdkCtx.BlockHeader().Time.Add(params.UnbondingTime)
	validator.UnbondingHeight = sdkCtx.BlockHeader().Height

	validator.UnbondingIds = append(validator.UnbondingIds, id)

	// save the now unbonded validator record and power index
	if err = k.SetValidator(ctx, validator); err != nil {
		return validator, err
	}

	if err = k.SetValidatorByPowerIndex(ctx, validator); err != nil {
		return validator, err
	}

	// Adds to unbonding validator queue
	if err = k.InsertUnbondingValidatorQueue(ctx, validator); err != nil {
		return validator, err
	}

	// trigger hook
	consAddr, err := validator.GetConsAddr()
	if err != nil {
		return validator, err
	}

	str, err := k.validatorAddressCodec.StringToBytes(validator.GetOperator())
	if err != nil {
		return validator, err
	}

	if err := k.Hooks().AfterValidatorBeginUnbonding(ctx, consAddr, str); err != nil {
		return validator, err
	}

	if err := k.SetValidatorByUnbondingID(ctx, validator, id); err != nil {
		return validator, err
	}

	if err := k.Hooks().AfterUnbondingInitiated(ctx, id); err != nil {
		return validator, err
	}

	return validator, nil
}

// perform all the store operations for when a validator status becomes unbonded
func (k Keeper) completeUnbondingValidator(ctx context.Context, validator types.ValidatorHV) (types.ValidatorHV, error) {
	validator = validator.UpdateStatus(types.Unbonded)
	if err := k.SetValidator(ctx, validator); err != nil {
		return validator, err
	}

	return validator, nil
}

// map of operator bech32-addresses to serialized power
// We use bech32 strings here, because we can't have slices as keys: map[[]byte][]byte
type validatorsByAddr map[string][]byte

// get the last validator set
func (k Keeper) getLastValidatorsByAddr(ctx context.Context) (validatorsByAddr, error) {
	last := make(validatorsByAddr)

	iterator, err := k.LastValidatorsIterator(ctx)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// extract the validator address from the key (prefix is 1-byte, addrLen is 1-byte)
		valAddr := types.AddressFromLastValidatorPowerKey(iterator.Key())
		valAddrStr, err := k.validatorAddressCodec.BytesToString(valAddr)
		if err != nil {
			return nil, err
		}

		powerBytes := iterator.Value()
		last[valAddrStr] = make([]byte, len(powerBytes))
		copy(last[valAddrStr], powerBytes)
	}

	return last, nil
}

// given a map of remaining validators to previous bonded power
// returns the list of validators to be unbonded, sorted by operator address
func sortNoLongerBonded(last validatorsByAddr, ac address.Codec) ([][]byte, error) {
	// sort the map keys for determinism
	noLongerBonded := make([][]byte, len(last))
	index := 0

	for valAddrStr := range last {
		valAddrBytes, err := ac.StringToBytes(valAddrStr)
		if err != nil {
			return nil, err
		}
		noLongerBonded[index] = valAddrBytes
		index++
	}
	// sorted by address - order doesn't matter
	sort.SliceStable(noLongerBonded, func(i, j int) bool {
		// -1 means strictly less than
		return bytes.Compare(noLongerBonded[i], noLongerBonded[j]) == -1
	})

	return noLongerBonded, nil
}

// DebugValidatorSigningInfo logs the current signing info for a validator for debugging purposes
func (k Keeper) DebugValidatorSigningInfo(ctx context.Context, consAddr sdk.ConsAddress) {
	signInfo, err := k.slashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
	if err != nil {
		k.Logger(ctx).Error(
			"Failed to get signing info for debugging",
			"consensus_addr", consAddr.String(),
			"error", err,
		)
		return
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentTime := sdkCtx.BlockHeader().Time

	k.Logger(ctx).Info(
		"DEBUG: Validator signing info",
		"consensus_addr", consAddr.String(),
		"current_time", currentTime.String(),
		"jailed_until", signInfo.JailedUntil.String(),
		"is_jailed_until_past", currentTime.After(signInfo.JailedUntil),
		"time_remaining", signInfo.JailedUntil.Sub(currentTime).String(),
		"start_height", signInfo.StartHeight,
		"index_offset", signInfo.IndexOffset,
		"missed_blocks_counter", signInfo.MissedBlocksCounter,
		"tombstoned", signInfo.Tombstoned,
	)
}
