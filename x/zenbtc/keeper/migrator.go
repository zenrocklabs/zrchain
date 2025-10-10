package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v2 "github.com/zenrocklabs/zenbtc/x/zenbtc/migrations/v2"
	v3 "github.com/zenrocklabs/zenbtc/x/zenbtc/migrations/v3"
	v4 "github.com/zenrocklabs/zenbtc/x/zenbtc/migrations/v4"
	v5 "github.com/zenrocklabs/zenbtc/x/zenbtc/migrations/v5"
	v6 "github.com/zenrocklabs/zenbtc/x/zenbtc/migrations/v6"
	v7 "github.com/zenrocklabs/zenbtc/x/zenbtc/migrations/v7"
	v8 "github.com/zenrocklabs/zenbtc/x/zenbtc/migrations/v8"
	v9 "github.com/zenrocklabs/zenbtc/x/zenbtc/migrations/v9"
	v10 "github.com/zenrocklabs/zenbtc/x/zenbtc/migrations/v10"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper *Keeper
}

// NewMigrator returns a new Migrator instance.
func NewMigrator(keeper *Keeper) Migrator {
	return Migrator{
		keeper: keeper,
	}
}

// Migrate1to2 migrates x/zenbtc params from consensus version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {

	if err := v2.UpdateParams(ctx, m.keeper.Params); err != nil {
		return err
	}

	return nil
}

// Migrate2to3 migrates x/zenbtc params from consensus version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {

	if err := v3.UpdateParams(ctx, m.keeper.Params); err != nil {
		return err
	}

	return nil
}

// Migrate3to4 migrates x/zenbtc params from consensus version 3 to 4.
func (m Migrator) Migrate3to4(ctx sdk.Context) error {

	if err := v4.UpdateParams(ctx, m.keeper.Params); err != nil {
		return err
	}

	return nil
}

// Migrate4to5 migrates x/zenbtc params from consensus version 4 to 5.
func (m Migrator) Migrate4to5(ctx sdk.Context) error {

	if err := v5.UpdateParams(ctx, m.keeper.Params); err != nil {
		return err
	}

	return nil
}

// Migrate5to6 migrates x/zenbtc params from consensus version 5 to 6.
func (m Migrator) Migrate5to6(ctx sdk.Context) error {

	if err := v6.MigrateLockTransactions(ctx, m.keeper.LockTransactionStore, m.keeper.LockTransactions, m.keeper.SetAuthority); err != nil {
		return err
	}

	return nil
}

// Migrate6to7 migrates x/zenbtc from consensus version 6 to 7.
func (m Migrator) Migrate6to7(ctx sdk.Context) error {
	return v7.PurgeInvalidState(ctx, m.keeper.BurnEvents, m.keeper.FirstPendingBurnEvent, m.keeper.Redemptions, m.keeper.FirstPendingRedemption)
}

// Migrate7to8 migrates x/zenbtc from consensus version 7 to 8.
func (m Migrator) Migrate7to8(ctx sdk.Context) error {
	return v8.RemoveStakedMints(ctx, m.keeper.PendingMintTransactionsMap, m.keeper.FirstPendingSolMintTransaction)
}

// Migrate8to9 migrates x/zenbtc from consensus version 8 to 9.
func (m Migrator) Migrate8to9(ctx sdk.Context) error {
	return v9.RemoveBurnedEvents(ctx, m.keeper.BurnEvents, m.keeper.FirstPendingBurnEvent)
}

// Migrate9to10 migrates x/zenbtc from consensus version 9 to 10.
func (m Migrator) Migrate9to10(ctx sdk.Context) error {
	return v10.RemoveInvalidRedemptions(ctx, m.keeper.Redemptions, m.keeper.FirstPendingRedemption)
}
