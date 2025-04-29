package keeper

type Migrator struct {
	keeper Keeper
}

func NewMigrator(keeper Keeper) *Migrator {
	return &Migrator{
		keeper: keeper,
	}
}

// func (m Migrator) Migrate1to2(ctx sdk.Context) error {
// 	ctx.Logger().With("module", types.ModuleName).Info("starting migration to v2")

// if err := v2.UpdateParams(ctx, m.keeper.ParamStore); err != nil {
// 	ctx.Logger().With("error", err).Error("failed to migrate zentp module")
// 	return err
// }

// 	return nil
// }
