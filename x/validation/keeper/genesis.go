package keeper

import (
	"context"
	"fmt"
	"strings"

	abci "github.com/cometbft/cometbft/abci/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
)

// InitGenesis sets the pool and parameters for the provided keeper.  For each
// validator in data, it sets that validator in the keeper along with manually
// setting the indexes. In addition, it also sets any delegations found in
// data. Finally, it updates the bonded validators.
// Returns final validator set after applying all declaration and delegations
func (k Keeper) InitGenesis(ctx context.Context, data *types.GenesisState) (res []abci.ValidatorUpdate) {
	bondedTokens := math.ZeroInt()
	notBondedTokens := math.ZeroInt()

	// We need to pretend to be "n blocks before genesis", where "n" is the
	// validator update delay, so that e.g. slashing periods are correctly
	// initialized for the validator set e.g. with a one-block offset - the
	// first TM block is at height 1, so state updates applied from
	// genesis.json are in block 0.
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx = sdkCtx.WithBlockHeight(1 - sdk.ValidatorUpdateDelay)
	ctx = sdkCtx

	if err := k.SetParams(ctx, types.Params(data.Params)); err != nil {
		panic(err)
	}

	hvParams := data.HVParams
	if err := k.HVParams.Set(ctx, hvParams); err != nil {
		panic(err)
	}

	if err := k.SetLastTotalPower(ctx, data.LastTotalPower); err != nil {
		panic(err)
	}

	for _, validator := range data.Validators {
		if err := k.SetValidator(ctx, validator); err != nil {
			panic(err)
		}

		// Manually set indices for the first time
		if err := k.SetValidatorByConsAddr(ctx, validator); err != nil {
			panic(err)
		}

		if err := k.SetValidatorByPowerIndex(ctx, validator); err != nil {
			panic(err)
		}

		// Call the creation hook if not exported
		if !data.Exported {
			valbz, err := k.ValidatorAddressCodec().StringToBytes(validator.GetOperator())
			if err != nil {
				panic(err)
			}
			if err := k.Hooks().AfterValidatorCreated(ctx, valbz); err != nil {
				panic(err)
			}
		}

		// update timeslice if necessary
		if validator.IsUnbonding() {
			if err := k.InsertUnbondingValidatorQueue(ctx, validator); err != nil {
				panic(err)
			}
		}

		switch validator.GetStatus() {
		case types.Bonded:
			bondedTokens = bondedTokens.Add(validator.GetTokens())

		case types.Unbonding, types.Unbonded:
			notBondedTokens = notBondedTokens.Add(validator.GetTokens())

		default:
			panic("invalid validator status")
		}
	}

	for _, delegation := range data.Delegations {
		delegatorAddress, err := k.authKeeper.AddressCodec().StringToBytes(delegation.DelegatorAddress)
		if err != nil {
			panic(fmt.Errorf("invalid delegator address: %s", err))
		}

		valAddr, err := k.validatorAddressCodec.StringToBytes(delegation.GetValidatorAddr())
		if err != nil {
			panic(err)
		}

		// Call the before-creation hook if not exported
		if !data.Exported {
			if err := k.Hooks().BeforeDelegationCreated(ctx, delegatorAddress, valAddr); err != nil {
				panic(err)
			}
		}

		if err := k.SetDelegation(ctx, delegation); err != nil {
			panic(err)
		}

		// Call the after-modification hook if not exported
		if !data.Exported {
			if err := k.Hooks().AfterDelegationModified(ctx, delegatorAddress, valAddr); err != nil {
				panic(err)
			}
		}
	}

	for _, ubd := range data.UnbondingDelegations {
		if err := k.SetUnbondingDelegation(ctx, ubd); err != nil {
			panic(err)
		}

		for _, entry := range ubd.Entries {
			if err := k.InsertUBDQueue(ctx, ubd, entry.CompletionTime); err != nil {
				panic(err)
			}
			notBondedTokens = notBondedTokens.Add(entry.Balance)
		}
	}

	for _, red := range data.Redelegations {
		if err := k.SetRedelegation(ctx, red); err != nil {
			panic(err)
		}

		for _, entry := range red.Entries {
			if err := k.InsertRedelegationQueue(ctx, red, entry.CompletionTime); err != nil {
				panic(err)
			}
		}
	}

	bondedCoins := sdk.NewCoins(sdk.NewCoin(data.Params.BondDenom, bondedTokens))
	notBondedCoins := sdk.NewCoins(sdk.NewCoin(data.Params.BondDenom, notBondedTokens))

	// check if the unbonded and bonded pools accounts exists
	bondedPool := k.GetBondedPool(ctx)
	if bondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	// TODO: remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862

	bondedBalance := k.bankKeeper.GetAllBalances(ctx, bondedPool.GetAddress())
	if bondedBalance.IsZero() {
		k.authKeeper.SetModuleAccount(ctx, bondedPool)
	}

	// if balance is different from bonded coins panic because genesis is most likely malformed
	if !bondedBalance.Equal(bondedCoins) {
		panic(fmt.Sprintf("bonded pool balance is different from bonded coins: %s <-> %s", bondedBalance, bondedCoins))
	}

	notBondedPool := k.GetNotBondedPool(ctx)
	if notBondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	notBondedBalance := k.bankKeeper.GetAllBalances(ctx, notBondedPool.GetAddress())
	if notBondedBalance.IsZero() {
		k.authKeeper.SetModuleAccount(ctx, notBondedPool)
	}

	// If balance is different from non bonded coins panic because genesis is most
	// likely malformed.
	if !notBondedBalance.Equal(notBondedCoins) {
		panic(fmt.Sprintf("not bonded pool balance is different from not bonded coins: %s <-> %s", notBondedBalance, notBondedCoins))
	}

	for asset, price := range data.AssetPrices {
		if err := k.AssetPrices.Set(ctx, types.Asset(asset), price); err != nil {
			panic(err)
		}
	}

	for _, btcBlockHeader := range data.BtcBlockHeaders {
		if err := k.BtcBlockHeaders.Set(ctx, btcBlockHeader.BlockHeight, btcBlockHeader); err != nil {
			panic(err)
		}
	}

	// TODO: check if this is correct
	k.SetSolanaRequestedNonce(ctx, k.zentpKeeper.GetSolanaParams(ctx).NonceAccountKey, true)

	for account, requested := range data.SolanaZentpAccountsRequested {
		k.SetSolanaZenTPRequestedAccount(ctx, account, requested)
	}

	for account, requested := range data.SolanaAccountsRequested {
		k.SolanaAccountsRequested.Set(ctx, account, requested)
	}

	for compositeKey, requested := range data.SolanaDctAccountsRequested {
		parts := strings.SplitN(compositeKey, ":", 2)
		if len(parts) != 2 {
			k.Logger(ctx).Error("invalid DCT Solana account key format", "key", compositeKey)
			continue
		}
		assetValue, ok := dcttypes.Asset_value[parts[0]]
		if !ok {
			k.Logger(ctx).Error("unknown DCT asset in Solana account request", "asset", parts[0])
			continue
		}
		asset := dcttypes.Asset(assetValue)
		key, err := k.dctAccountKey(asset, parts[1])
		if err != nil {
			k.Logger(ctx).Error("error constructing DCT account key", "asset", asset.String(), "owner", parts[1], "error", err)
			continue
		}
		if err := k.SolanaDCTAccountsRequested.Set(ctx, key, requested); err != nil {
			k.Logger(ctx).Error("error setting DCT Solana account request", "asset", asset.String(), "owner", parts[1], "error", err)
		}
	}

	k.BackfillRequests.Set(ctx, data.BackfillRequest)

	for _, ethereumNonce := range data.LastUsedEthereumNonce {
		k.EthereumNonceRequested.Set(ctx, ethereumNonce.Nonce, true)
	}

	k.RequestedHistoricalBitcoinHeaders.Set(ctx, data.RequestedHistoricalBitcoinHeaders)

	for address, amount := range data.AvsRewardsPool {
		k.AVSRewardsPool.Set(ctx, address, amount)
	}

	for solanaNonce, requested := range data.SolanaNonceRequested {
		k.SolanaNonceRequested.Set(ctx, solanaNonce, requested)
	}

	for ethereumNonce, requested := range data.EthereumNonceRequested {
		k.EthereumNonceRequested.Set(ctx, ethereumNonce, requested)
	}

	for account, requested := range data.SolanaZentpAccountsRequested {
		k.SetSolanaZenTPRequestedAccount(ctx, account, requested)
	}

	for solanaAccount, requested := range data.SolanaAccountsRequested {
		k.SolanaAccountsRequested.Set(ctx, solanaAccount, requested)
	}

	for _, validatorMismatchCount := range data.ValidatorMismatchCounts {
		k.ValidatorMismatchCounts.Set(ctx, validatorMismatchCount.ValidatorAddress, validatorMismatchCount)
	}

	for blockHeight, slashEvent := range data.SlashEvents {
		k.SlashEvents.Set(ctx, blockHeight, slashEvent)
	}

	err := k.SlashEventCount.Set(ctx, uint64(len(data.SlashEvents)))
	if err != nil {
		panic(err)
	}

	// TODO: check if this is correct
	for blockHeight, validationInfo := range data.ValidationInfos {
		k.ValidationInfos.Set(ctx, blockHeight, validationInfo)
	}

	if data.LastValidVeHeight > 0 {
		err = k.LastValidVEHeight.Set(ctx, data.LastValidVeHeight)
		if err != nil {
			panic(err)
		}
	}

	if data.LastCompletedZentpMintId > 0 {
		err = k.LastCompletedZentpMintID.Set(ctx, data.LastCompletedZentpMintId)
		if err != nil {
			panic(err)
		}
	}

	// don't need to run CometBFT updates if we exported
	if data.Exported {
		for _, lv := range data.LastValidatorPowers {
			valAddr, err := k.validatorAddressCodec.StringToBytes(lv.Address)
			if err != nil {
				panic(err)
			}

			err = k.SetLastValidatorPower(ctx, valAddr, lv.Power)
			if err != nil {
				panic(err)
			}

			validator, err := k.GetZenrockValidator(ctx, valAddr)
			if err != nil {
				panic(fmt.Sprintf("validator %s not found", lv.Address))
			}

			// keep the next-val-set offset, use the last power for the first block
			update := validator.ABCIValidatorUpdate(k.PowerReduction(ctx), lv.Power)
			res = append(res, update)
		}
	} else {
		var err error

		res, err = k.ApplyAndReturnValidatorSetUpdates(ctx)
		if err != nil {
			panic(err)
		}
	}

	return res
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, params, validators, and bonds found in
// the keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	var unbondingDelegations []types.UnbondingDelegation

	err := k.IterateUnbondingDelegations(ctx, func(_ int64, ubd types.UnbondingDelegation) (stop bool) {
		unbondingDelegations = append(unbondingDelegations, ubd)
		return false
	})
	if err != nil {
		panic(err)
	}

	var redelegations []types.Redelegation

	err = k.IterateRedelegations(ctx, func(_ int64, red types.Redelegation) (stop bool) {
		redelegations = append(redelegations, red)
		return false
	})
	if err != nil {
		panic(err)
	}

	var lastValidatorPowers []types.LastValidatorPower

	err = k.IterateLastValidatorPowers(ctx, func(addr sdk.ValAddress, power int64) (stop bool) {
		addrStr, err := k.validatorAddressCodec.BytesToString(addr)
		if err != nil {
			panic(err)
		}
		lastValidatorPowers = append(lastValidatorPowers, types.LastValidatorPower{Address: addrStr, Power: power})
		return false
	})
	if err != nil {
		panic(err)
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	hvParams, err := k.HVParams.Get(ctx)
	if err != nil {
		panic(err)
	}

	totalPower, err := k.GetLastTotalPower(ctx)
	if err != nil {
		panic(err)
	}

	allDelegations, err := k.GetAllDelegations(ctx)
	if err != nil {
		panic(err)
	}
	var delegations []types.Delegation
	for _, delegation := range allDelegations {
		delegations = append(delegations, types.Delegation(delegation))
	}

	allValidators, err := k.GetAllZenrockValidators(ctx)
	if err != nil {
		panic(err)
	}

	assetPricesMap, err := k.GetAssetPrices(ctx)
	if err != nil {
		panic(err)
	}

	// Convert Asset enum keys to int32 for proto compatibility
	assetPrices := make(map[int32]math.LegacyDec)
	for asset, price := range assetPricesMap {
		assetPrices[int32(asset)] = price
	}

	lastValidVeHeight, err := k.GetLastValidVeHeight(ctx)
	if err != nil {
		panic(err)
	}

	slashEvents, err := k.GetSlashEvents(ctx)
	if err != nil {
		panic(err)
	}

	slashEventCount := uint64(len(slashEvents))

	validationInfos, err := k.GetValidationInfos(ctx)
	if err != nil {
		panic(err)
	}

	btcBlockHeaders, err := k.GetBtcBlockHeaders(ctx)
	if err != nil {
		panic(err)
	}

	lastUsedSolanaNonce, err := k.GetLastUsedSolanaNonce(ctx)
	if err != nil {
		panic(err)
	}

	backfillRequest, err := k.GetBackfillRequests(ctx)
	if err != nil {
		panic(err)
	}

	lastUsedEthereumNonce, err := k.GetLastUsedEthereumNonce(ctx)
	if err != nil {
		panic(err)
	}

	requestedHistoricalBitcoinHeaders, err := k.GetRequestedHistoricalBitcoinHeaders(ctx)
	if err != nil {
		panic(err)
	}

	avsRewardsPool, err := k.GetAvsRewardsPool(ctx)
	if err != nil {
		panic(err)
	}

	ethereumNonceRequested, err := k.GetEthereumNonceRequested(ctx)
	if err != nil {
		panic(err)
	}

	solanaNonceRequested, err := k.GetSolanaNonceRequested(ctx)
	if err != nil {
		panic(err)
	}

	solanaAccountsRequested, err := k.GetSolanaAccountsRequested(ctx)
	if err != nil {
		panic(err)
	}

	solanaZenTPAccountsRequested, err := k.GetSolanaZenTPAccountsRequested(ctx)
	if err != nil {
		panic(err)
	}

	solanaDCTAccountsRequested, err := k.GetSolanaDCTAccountsRequested(ctx)
	if err != nil {
		panic(err)
	}

	validatorMismatchCounts, err := k.GetValidatorMismatchCounts(ctx)
	if err != nil {
		panic(err)
	}

	lastCompletedZentpMintID, err := k.GetLastCompletedZentpMintID(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Params:                            types.Params(params),
		LastTotalPower:                    totalPower,
		LastValidatorPowers:               lastValidatorPowers,
		Validators:                        allValidators,
		Delegations:                       delegations,
		UnbondingDelegations:              unbondingDelegations,
		Redelegations:                     redelegations,
		Exported:                          true,
		HVParams:                          hvParams,
		AssetPrices:                       assetPrices,
		LastValidVeHeight:                 lastValidVeHeight,
		SlashEvents:                       slashEvents,
		SlashEventCount:                   slashEventCount,
		ValidationInfos:                   validationInfos,
		BtcBlockHeaders:                   btcBlockHeaders,
		LastUsedSolanaNonce:               lastUsedSolanaNonce,
		BackfillRequest:                   backfillRequest,
		LastUsedEthereumNonce:             lastUsedEthereumNonce,
		RequestedHistoricalBitcoinHeaders: requestedHistoricalBitcoinHeaders,
		AvsRewardsPool:                    avsRewardsPool,
		EthereumNonceRequested:            ethereumNonceRequested,
		SolanaNonceRequested:              solanaNonceRequested,
		SolanaZentpAccountsRequested:      solanaZenTPAccountsRequested,
		SolanaAccountsRequested:           solanaAccountsRequested,
		SolanaDctAccountsRequested:        solanaDCTAccountsRequested,
		ValidatorMismatchCounts:           validatorMismatchCounts,
		LastCompletedZentpMintId:          lastCompletedZentpMintID,
	}
}
