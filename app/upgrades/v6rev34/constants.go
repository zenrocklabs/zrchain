package v6rev34

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/Zenrock-Foundation/zrchain/v6/app/upgrades"
	zenextypes "github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
)

const UpgradeName = "v6rev34"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{zenextypes.ModuleName},
		Deleted: []string{},
		Renamed: []storetypes.StoreRename{},
	},
}
