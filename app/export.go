package app

import (
	"encoding/json"
	"fmt"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	storetypes "cosmossdk.io/store/types"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	validation "github.com/Zenrock-Foundation/zrchain/v6/x/validation/module"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
)

// ExportAppStateAndValidators exports the state of the application for a genesis file.
func (app *ZenrockApp) ExportAppStateAndValidators(forZeroHeight bool, jailAllowedAddrs, modulesToExport []string) (servertypes.ExportedApp, error) {
	ctx := app.NewContextLegacy(true, cmtproto.Header{Height: app.LastBlockHeight()})

	// We export at last height + 1, because that's the height at which
	// CometBFT will start InitChain.
	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0
		if err := app.prepForZeroHeightGenesis(ctx, jailAllowedAddrs); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	genState, err := app.ModuleManager.ExportGenesisForModules(ctx, app.appCodec, modulesToExport)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	appState, err := json.MarshalIndent(genState, "", "  ")
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	validators, err := validation.WriteValidators(ctx, app.ValidationKeeper)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	return servertypes.ExportedApp{
		AppState:        appState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
	}, nil
}

// prepForZeroHeightGenesis prepares the blockchain state for genesis at height zero.
// This feature will be replaced by exports at specific block heights in the future.
func (app *ZenrockApp) prepForZeroHeightGenesis(ctx sdk.Context, jailAllowedAddrs []string) error {
	applyAllowedAddrs := len(jailAllowedAddrs) > 0
	allowedAddrsMap := make(map[string]bool)
	for _, addr := range jailAllowedAddrs {
		if _, err := sdk.ValAddressFromBech32(addr); err != nil {
			return fmt.Errorf("invalid address %s: %w", addr, err)
		}
		allowedAddrsMap[addr] = true
	}

	// Apply invariant checks on current state
	app.CrisisKeeper.AssertInvariants(ctx)

	if err := app.withdrawValidatorCommissions(ctx); err != nil {
		return err
	}

	if err := app.withdrawDelegatorRewards(ctx); err != nil {
		return err
	}

	app.DistrKeeper.DeleteAllValidatorSlashEvents(ctx)
	app.DistrKeeper.DeleteAllValidatorHistoricalRewards(ctx)

	if err := app.resetValidators(ctx, applyAllowedAddrs, allowedAddrsMap); err != nil {
		return err
	}

	if err := app.resetSlashingState(ctx); err != nil {
		return err
	}

	return nil
}

func (app *ZenrockApp) withdrawValidatorCommissions(ctx sdk.Context) error {
	return app.ValidationKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) bool {
		valBz, err := app.ValidationKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
		if err != nil {
			app.Logger().Error("Failed to convert address: %v", err)
			return true
		}

		_, err = app.DistrKeeper.WithdrawValidatorCommission(ctx, valBz)
		if err != nil {
			app.Logger().Error("Failed to withdraw commission: %v", err)
			return true
		}

		return false
	})
}

func (app *ZenrockApp) withdrawDelegatorRewards(ctx sdk.Context) error {
	dels, err := app.ValidationKeeper.GetAllDelegations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all delegations: %w", err)
	}

	for _, delegation := range dels {
		valAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			return fmt.Errorf("invalid validator address %s: %w", delegation.ValidatorAddress, err)
		}

		delAddr := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)
		_, err = app.DistrKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
		if err != nil {
			return fmt.Errorf("failed to withdraw delegation rewards for %s: %w", delegation.DelegatorAddress, err)
		}
	}

	return nil
}

func (app *ZenrockApp) resetValidators(ctx sdk.Context, applyAllowedAddrs bool, allowedAddrsMap map[string]bool) error {
	store := ctx.KVStore(app.GetKey(validationtypes.StoreKey))
	iter := storetypes.KVStoreReversePrefixIterator(store, stakingtypes.ValidatorsKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(stakingtypes.AddressFromValidatorsKey(iter.Key()))
		validator, err := app.ValidationKeeper.GetZenrockValidator(ctx, addr)
		if err != nil {
			return fmt.Errorf("expected validator not found: %w", err)
		}

		validator.UnbondingHeight = 0
		if applyAllowedAddrs && !allowedAddrsMap[addr.String()] {
			validator.Jailed = true
		}

		if err := app.ValidationKeeper.SetValidator(ctx, validator); err != nil {
			return fmt.Errorf("error setting validator: %w", err)
		}
	}

	if _, err := app.ValidationKeeper.ApplyAndReturnValidatorSetUpdates(ctx); err != nil {
		return fmt.Errorf("error applying validator set updates: %w", err)
	}

	return nil
}

func (app *ZenrockApp) resetSlashingState(ctx sdk.Context) error {
	return app.SlashingKeeper.IterateValidatorSigningInfos(ctx, func(addr sdk.ConsAddress, info slashingtypes.ValidatorSigningInfo) bool {
		// Reset start height on signing infos
		info.StartHeight = 0
		if err := app.SlashingKeeper.SetValidatorSigningInfo(ctx, addr, info); err != nil {
			app.Logger().Error("Failed to set signing info: %v", err)
			return true
		}
		return false
	})
}
