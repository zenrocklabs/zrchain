package v4

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/Zenrock-Foundation/zrchain/v4/app/upgrades"
)

const UpgradeName = "v4"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{"zenbtc"},
		Deleted: []string{},
		Renamed: []storetypes.StoreRename{},
	},
}
