package keeper

import (
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"

	v2 "github.com/cosmos/cosmos-sdk/x/staking/migrations/v2"
	v3 "github.com/cosmos/cosmos-sdk/x/staking/migrations/v3"
	v4 "github.com/cosmos/cosmos-sdk/x/staking/migrations/v4"
	v5 "github.com/cosmos/cosmos-sdk/x/staking/migrations/v5"

	v6 "github.com/Zenrock-Foundation/zrchain/v6/x/validation/migrations/v6"
	v7 "github.com/Zenrock-Foundation/zrchain/v6/x/validation/migrations/v7"
	v8 "github.com/Zenrock-Foundation/zrchain/v6/x/validation/migrations/v8"
	v9 "github.com/Zenrock-Foundation/zrchain/v6/x/validation/migrations/v9"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper         *Keeper
	legacySubspace exported.Subspace
}

// NewMigrator returns a new Migrator instance.
func NewMigrator(keeper *Keeper, legacySubspace exported.Subspace) Migrator {
	return Migrator{
		keeper:         keeper,
		legacySubspace: legacySubspace,
	}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	store := runtime.KVStoreAdapter(m.keeper.storeService.OpenKVStore(ctx))
	return v2.MigrateStore(ctx, store)
}

// Migrate2to3 migrates x/staking state from consensus version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	store := runtime.KVStoreAdapter(m.keeper.storeService.OpenKVStore(ctx))
	return v3.MigrateStore(ctx, store, m.keeper.cdc, m.legacySubspace)
}

// Migrate3to4 migrates x/staking state from consensus version 3 to 4.
func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	store := runtime.KVStoreAdapter(m.keeper.storeService.OpenKVStore(ctx))
	return v4.MigrateStore(ctx, store, m.keeper.cdc, m.legacySubspace)
}

// Migrate4to5 migrates x/staking state from consensus version 4 to 5.
func (m Migrator) Migrate4to5(ctx sdk.Context) error {
	store := runtime.KVStoreAdapter(m.keeper.storeService.OpenKVStore(ctx))
	return v5.MigrateStore(ctx, store, m.keeper.cdc)
}

// Migrate5to6 migrates x/staking state from consensus version 5 to 6.
func (m Migrator) Migrate5to6(ctx sdk.Context) error {
	return v6.UpdateParams(ctx, m.keeper.HVParams)
}

// Migrate6to7 migrates x/staking state from consensus version 6 to 7.
func (m Migrator) Migrate6to7(ctx sdk.Context) error {
	return v7.UpdateParams(ctx, m.keeper.HVParams)
}

func (m Migrator) Migrate7to8(ctx sdk.Context) error {
	return v8.UpdateBtcBlockHeaders(ctx, m.keeper.BtcBlockHeaders, m.keeper.ValidationInfos)
}

func (m Migrator) Migrate8to9(ctx sdk.Context) error {
	return v9.ClearEthereumNonceData(ctx, m.keeper.LastUsedEthereumNonce)
}
