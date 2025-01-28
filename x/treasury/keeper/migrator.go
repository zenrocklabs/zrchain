package keeper

import (
	v2 "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/migrations/v2"
	v3 "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/migrations/v3"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
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

	err := v2.ChangeKeyIdtoKeyIds(ctx, m.keeper.SignRequestStore, m.keeper.cdc)
	if err != nil {
		ctx.Logger().With("error", err).Error("failed to migrate treasury module")
		return err
	}

	return nil
}

func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	ctx.Logger().With("module", types.ModuleName).Info("starting migration to v3")

	err := v3.ChangeZenBtcMetadataChainIdtoCaip2Id(ctx, m.keeper.KeyStore, m.keeper.KeyRequestStore, m.keeper.cdc)
	if err != nil {
		ctx.Logger().With("error", err).Error("failed to migrate treasury module")
		return err
	}

	return nil
}
