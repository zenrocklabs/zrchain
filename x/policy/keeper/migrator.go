package keeper

import (
	v2 "github.com/Zenrock-Foundation/zrchain/v4/x/policy/migrations/v2"
	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
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
	//v1.UpdateParams(ctx, m.keeper.ParamStore)

	err := v2.UpdatePolicies(ctx, m.keeper.PolicyStore, m.keeper.cdc)
	ctx.Logger().Error("failed to migrate policies module", "error", err)

	return nil
}
