package v5beta4

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/Zenrock-Foundation/zrchain/v5/app/upgrades"
)

const UpgradeName = "v5beta4"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
		Renamed: []storetypes.StoreRename{},
	},
}