package keeper

import (
	"fmt"

	"github.com/Zenrock-Foundation/zrchain/v6/x/mint/exported"
	v2 "github.com/Zenrock-Foundation/zrchain/v6/x/mint/migrations/v2"
	v3 "github.com/Zenrock-Foundation/zrchain/v6/x/mint/migrations/v3"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Migrator is a struct for handling in-place state migrations.
type Migrator struct {
	keeper         Keeper
	legacySubspace exported.Subspace
}

// NewMigrator returns Migrator instance for the state migration.
func NewMigrator(k Keeper, ss exported.Subspace) Migrator {
	return Migrator{
		keeper:         k,
		legacySubspace: ss,
	}
}

// Migrate1to2 migrates the x/mint module state from the consensus version 1 to
// version 2. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/mint
// module state.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v2.Migrate(ctx, m.keeper.storeService.OpenKVStore(ctx), m.legacySubspace, m.keeper.cdc)
}

// Migrate migrates the x/mint module state from the consensus version 2 to
// version 3. Specifically, it adds several new parameters to the mint module
// and removes the legacy minter logic and replaces it with a deflationary
// model.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	authKeeper, ok := m.keeper.accountKeeper.(authkeeper.AccountKeeper)
	if !ok {
		return fmt.Errorf("accountKeeper is not of type authkeeper.AccountKeeper")
	}
	// Get the current mint module account
	moduleAccount, _ := authKeeper.GetModuleAccountAndPermissions(ctx, v3.ModuleName)

	// Create a new base account for the mint module with using the current attributes
	baseAccount := authtypes.NewBaseAccount(
		authKeeper.GetModuleAddress(v3.ModuleName),
		nil,
		moduleAccount.GetAccountNumber(),
		moduleAccount.GetSequence(),
	)

	// Create a new module account with the updated permissions
	macc := authtypes.NewModuleAccount(baseAccount, v3.ModuleName, authtypes.Minter, authtypes.Burner)

	// Set the new module account
	m.keeper.accountKeeper.SetModuleAccount(ctx, macc)

	err := authKeeper.ValidatePermissions(moduleAccount)
	if err != nil {
		return err
	}

	return v3.UpdateParams(ctx, m.keeper.Params)
}
