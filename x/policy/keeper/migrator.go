package keeper

import (
	v2 "github.com/Zenrock-Foundation/zrchain/v5/x/policy/migrations/v2"
	"github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
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

func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	ctx.Logger().With("module", types.ModuleName).Info("starting migration to v3")

	err := v2.UpdatePolicies(ctx, m.keeper.PolicyStore, m.keeper.cdc)
	if err != nil {
		ctx.Logger().With("error", err).Error("failed to migrate policies module")
		return err
	}

	return nil
}
