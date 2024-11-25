package keeper

import (
	"fmt"

	"github.com/Zenrock-Foundation/zrchain/v5/x/mint/exported"
	v2 "github.com/Zenrock-Foundation/zrchain/v5/x/mint/migrations/v2"
	v3 "github.com/Zenrock-Foundation/zrchain/v5/x/mint/migrations/v3"

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
	moduleAccount, perms := authKeeper.GetModuleAccountAndPermissions(ctx, v3.ModuleName)
	// mintAcc.GetAddress()
	// address := authKeeper.GetModuleAddress(v3.ModuleName)
	// address := mintAcc.GetAddress()
	// fmt.Println("Mint Module Address:", address)

	// account := authKeeper.GetAccount(ctx, address)

	baseAccount := authtypes.NewBaseAccount(
		authKeeper.GetModuleAddress(v3.ModuleName),
		nil,
		moduleAccount.GetAccountNumber(),
		moduleAccount.GetSequence(),
	)
	// baseAccount, ok := account.(*authtypes.BaseAccount)
	if !ok {
		return fmt.Errorf("account is not of type *authtypes.BaseAccount: %v", baseAccount)
	}

	macc := authtypes.NewModuleAccount(baseAccount, authtypes.Minter, authtypes.Burner)

	fmt.Println("Module Account GetModuleAccountAndPermissions:", macc)
	// fmt.Println("Module Account m.keeper.accountKeeper.Getmoduleaccount: ", macc)
	perms = macc.GetPermissions()
	fmt.Println("Mint Module Permissions BEFORE:", perms)

	m.keeper.accountKeeper.SetModuleAccount(ctx, macc)

	moduleAccount2 := authKeeper.GetModuleAccount(ctx, v3.ModuleName)
	perms = moduleAccount2.GetPermissions()
	fmt.Println("Mint Module Permissions AFTER:", perms)
	err := authKeeper.ValidatePermissions(moduleAccount)
	if err != nil {
		return err
	}

	// authtypes.NewPermissionsForAddress(v3.ModuleName, perms)

	// authKeeper.SetModuleAccount(ctx, moduleAccount)

	return v3.UpdateParams(ctx, m.keeper.Params, authKeeper)
}
