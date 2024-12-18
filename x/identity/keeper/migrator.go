package keeper

import (
	v2 "github.com/Zenrock-Foundation/zrchain/v5/x/identity/migrations/v2"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrator is a struct for handling in-place state migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns Migrator instance for the state migration.
func NewMigrator(k Keeper) Migrator {
	return Migrator{
		keeper: k,
	}
}

// Migrate2to3 migrates the x/treasury module state from the consensus version 2 to
// version 3. Specifically, it adds a new field 'fees' to the keyring type
// and migrates 'key_req_fee' and 'sig_req_fee' to use the new  fees type
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v2.MigrateKeyrings(ctx, m.keeper.KeyringStore, m.keeper.cdc)
}
