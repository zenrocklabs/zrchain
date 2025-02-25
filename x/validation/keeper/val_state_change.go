package keeper

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"

	abci "github.com/cometbft/cometbft/abci/types"
	gogotypes "github.com/cosmos/gogoproto/types"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v5/x/validation/types"
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
	totalPower := math.ZeroInt()
	amtFromBondedToNotBonded, amtFromNotBondedToBonded := math.ZeroInt(), math.ZeroInt()

	// Retrieve the last validator set.
	// The persistent set is updated later in this function.
	// (see LastValidatorPowerKey).
	last, err := k.getLastValidatorsByAddr(ctx)
	if err != nil {
		return nil, err
	}

	// Iterate over validators, highest power to lowest.
	iterator, err := k.ValidatorsPowerStoreIterator(ctx)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {
		// everything that is iterated in this loop is becoming or already a
		// part of the bonded validator set
		valAddr := sdk.ValAddress(iterator.Value())
		validator := k.mustGetValidator(ctx, valAddr)

		if validator.Jailed {
			panic("should never retrieve a jailed validator from the power store")
		}

		// if we get to a zero-power validator (which we don't bond),
		// there are no more possible bonded validators
		if validator.PotentialConsensusPower(k.PowerReduction(ctx)) == 0 {
			break
		}

		// apply the appropriate state change if necessary
		switch {
		case validator.IsUnbonded():
			validator, err = k.unbondedToBonded(ctx, validator)
			if err != nil {
				return
			}
			amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(validator.GetTokens())
		case validator.IsUnbonding():
			validator, err = k.unbondingToBonded(ctx, validator)
			if err != nil {
				return
			}
			amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(validator.GetTokens())
		case validator.IsBonded():
			// no state change
		default:
			panic("unexpected validator status")
		}

		// fetch the old power bytes
		valAddrStr, err := k.validatorAddressCodec.BytesToString(valAddr)
		if err != nil {
			return nil, err
		}
		oldPowerBytes, found := last[valAddrStr]

		stakeableAssets, err := k.GetStakeableAssetPrices(ctx)
		if err != nil {
			return nil, err
		}

		nativePower := math.LegacyZeroDec()
		avsPower := math.LegacyZeroDec()

		pricesAreValid := false
		for _, asset := range stakeableAssets {
			if !asset.PriceUSD.IsZero() {
				pricesAreValid = true
				break
			}
		}

		for _, asset := range stakeableAssets {
			switch asset.Asset {
			case types.Asset_ROCK:
				if pricesAreValid {
					nativePower = asset.PriceUSD.MulInt64(validator.ConsensusPower(powerReduction))
				} else {
					nativePower = math.LegacyNewDec(validator.ConsensusPower(powerReduction))
				}
			case types.Asset_BTC:
				exchangeRate, err := k.zenBTCKeeper.GetExchangeRate(ctx)
				if err != nil {
					exchangeRate = math.LegacyOneDec()
				}
				avsPower = asset.PriceUSD.Mul(exchangeRate).MulInt64(adjustPowerToPrecision(validator.TokensAVS, asset.Precision).Int64())
			case types.Asset_ETH:
				continue
			}
		}

		newPower := nativePower.Add(avsPower).TruncateInt64()
		if newPower == 0 {
			newPower = 1
		}

		consAddr, err := validator.GetConsAddr()
		if err != nil {
			return nil, err
		}

		k.Logger(ctx).Debug(fmt.Sprintf(
			"\nvalidator: %s | %s\ntoken stake: native=%d, avs=%d, total=%d\nstake value: native=%s, avs=%s, total=%d",
			valAddrStr, sdk.ConsAddress(consAddr).String(),
			validator.ConsensusPower(powerReduction), validator.AVSConsensusPower(powerReduction),
			validator.ConsensusPower(powerReduction)+validator.AVSConsensusPower(powerReduction),
			nativePower.String(), avsPower.String(), newPower,
		))

		newPowerBytes := k.cdc.MustMarshal(&gogotypes.Int64Value{Value: newPower})

		// update the validator set if power has changed
		if !found || !bytes.Equal(oldPowerBytes, newPowerBytes) {
			updates = append(updates, validator.ABCIValidatorUpdate(powerReduction, newPower))

			if err = k.SetLastValidatorPower(ctx, valAddr, newPower); err != nil {
				return nil, err
			}
		}

		delete(last, valAddrStr)
		count++

		totalPower = totalPower.Add(math.NewInt(newPower))
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
		str, err := k.validatorAddressCodec.StringToBytes(validator.GetOperator())
		if err != nil {
			return nil, err
		}
		amtFromBondedToNotBonded = amtFromBondedToNotBonded.Add(validator.GetTokens())
		if err = k.DeleteLastValidatorPower(ctx, str); err != nil {
			return nil, err
		}

		updates = append(updates, validator.ABCIValidatorUpdateZero())
	}

	// Update the pools based on the recent updates in the validator set:
	// - The tokens from the non-bonded candidates that enter the new validator set need to be transferred
	// to the Bonded pool.
	// - The tokens from the bonded validators that are being kicked out from the validator set
	// need to be transferred to the NotBonded pool.
	switch {
	// Compare and subtract the respective amounts to only perform one transfer.
	// This is done in order to avoid doing multiple updates inside each iterator/loop.
	case amtFromNotBondedToBonded.GT(amtFromBondedToNotBonded):
		if err = k.notBondedTokensToBonded(ctx, amtFromNotBondedToBonded.Sub(amtFromBondedToNotBonded)); err != nil {
			return nil, err
		}
	case amtFromNotBondedToBonded.LT(amtFromBondedToNotBonded):
		if err = k.bondedTokensToNotBonded(ctx, amtFromBondedToNotBonded.Sub(amtFromNotBondedToBonded)); err != nil {
			return nil, err
		}
	default: // equal amounts of tokens; no update required
	}

	// set total power on lookup index if there are any updates
	if len(updates) > 0 {
		if err = k.SetLastTotalPower(ctx, totalPower); err != nil {
			return nil, err
		}
	}

	// set the list of validator updates
	if err = k.SetValidatorUpdates(ctx, updates); err != nil {
		return nil, err
	}

	return updates, err
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

func adjustPowerToPrecision(tokens math.Int, assetPrecision uint32) math.Int {
	precisionDiff := assetPrecision - 6 // relative to ROCK's 6 decimals of precision
	powerOf10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precisionDiff)), nil)
	return tokens.Quo(math.NewIntFromBigInt(powerOf10))
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
