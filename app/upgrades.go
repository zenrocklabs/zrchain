package app

import (
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/Zenrock-Foundation/zrchain/v4/app/upgrades"
	v2 "github.com/Zenrock-Foundation/zrchain/v4/app/upgrades/v2"
	v3 "github.com/Zenrock-Foundation/zrchain/v4/app/upgrades/v3"
)

var Upgrades = []upgrades.Upgrade{
	v2.Upgrade,
	v3.Upgrade,
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

		if upgradeInfo.Name == upgrade.UpgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
			// configure store loader that checks if version == upgradeHeight and applies store upgrades
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &upgrade.StoreUpgrades))
		}
	}
}
