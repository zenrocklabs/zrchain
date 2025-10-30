package v6rev40

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/Zenrock-Foundation/zrchain/v6/app/upgrades"
	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

const UpgradeName = "v6rev40"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{dcttypes.ModuleName},
		Deleted: []string{},
		Renamed: []storetypes.StoreRename{},
	},
}
