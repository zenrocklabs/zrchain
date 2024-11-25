package app

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/Zenrock-Foundation/zrchain/v5/app/upgrades"
	v6 "github.com/Zenrock-Foundation/zrchain/v5/app/upgrades/v6"
	v3 "github.com/Zenrock-Foundation/zrchain/v5/x/mint/migrations/v3"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var Upgrades = []upgrades.Upgrade{
	v6.Upgrade,
}

func (app ZenrockApp) RegisterUpgradeHandlers() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	for _, upgrade := range Upgrades {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(app.ModuleManager, app.Configurator()),
		)

		ctx := context.Background()

		mintAccBase := authtypes.NewEmptyModuleAccount(v3.ModuleName, authtypes.Minter, authtypes.Burner)

		// Update the permissions
		macc := authtypes.NewModuleAccount(mintAccBase.BaseAccount, authtypes.Minter, authtypes.Burner)

		// Save the updated account
		app.AccountKeeper.SetModuleAccount(ctx, macc)

		if upgradeInfo.Name == upgrade.UpgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
			// configure store loader that checks if version == upgradeHeight and applies store upgrades
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &upgrade.StoreUpgrades))
		}
	}
}
