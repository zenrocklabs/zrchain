package keeper

import (
	v2 "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/migrations/v2"
	v3 "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/migrations/v3"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migrator struct {
	keeper Keeper
}

func NewMigrator(keeper Keeper) *Migrator {
	return &Migrator{
		keeper: keeper,
	}
}

func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	ctx.Logger().With("module", types.ModuleName).Info("starting migration to v2")

	if err := v2.UpdateParams(ctx, m.keeper.ParamStore); err != nil {
		ctx.Logger().With("error", err).Error("failed to migrate zentp module")
		return err
	}

	return nil
}

func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	ctx.Logger().With("module", types.ModuleName).Info("starting zentp migration to v3")

	if err := v3.UpdateMintStore(ctx, m.keeper.mintStore, m.keeper.MintStore); err != nil {
		ctx.Logger().With("error", err).Error("failed to migrate zentp module")
		return err
	}

	if err := v3.UpdateBurnStore(ctx, m.keeper.burnStore, m.keeper.BurnStore); err != nil {
		ctx.Logger().With("error", err).Error("failed to migrate zentp module")
		return err
	}

	return nil
}
