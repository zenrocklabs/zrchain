package shared

import (
	"context"

	"cosmossdk.io/math"
	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

// DCTKeeper defines the subset of keeper functionality required by other modules
// to interact with decentralised custody token state.
type DCTKeeper interface {
	GetAuthority() string

	GetParams(ctx context.Context) (dcttypes.Params, error)
	GetAssetParams(ctx context.Context, asset dcttypes.Asset) (dcttypes.AssetParams, error)
	ListSupportedAssets(ctx context.Context) ([]dcttypes.Asset, error)
	GetSolanaParams(ctx context.Context, asset dcttypes.Asset) (*dcttypes.Solana, error)

	GetDepositKeyringAddr(ctx context.Context, asset dcttypes.Asset) (string, error)
	GetStakerKeyID(ctx context.Context, asset dcttypes.Asset) (uint64, error)
	GetEthMinterKeyID(ctx context.Context, asset dcttypes.Asset) (uint64, error)
	GetUnstakerKeyID(ctx context.Context, asset dcttypes.Asset) (uint64, error)
	GetCompleterKeyID(ctx context.Context, asset dcttypes.Asset) (uint64, error)
	GetRewardsDepositKeyID(ctx context.Context, asset dcttypes.Asset) (uint64, error)
	GetChangeAddressKeyIDs(ctx context.Context, asset dcttypes.Asset) ([]uint64, error)
	GetProxyAddress(ctx context.Context, asset dcttypes.Asset) (string, error)
	GetBitcoinProxyAddress(ctx context.Context, asset dcttypes.Asset) (string, error)

	SetPendingMintTransaction(ctx context.Context, pendingMintTransaction dcttypes.PendingMintTransaction) error
	WalkPendingMintTransactions(ctx context.Context, asset dcttypes.Asset, fn func(id uint64, pendingMintTransaction dcttypes.PendingMintTransaction) (stop bool, err error)) error
	GetPendingMintTransaction(ctx context.Context, asset dcttypes.Asset, id uint64) (dcttypes.PendingMintTransaction, error)
	HasPendingMintTransaction(ctx context.Context, asset dcttypes.Asset, id uint64) (bool, error)
	GetFirstPendingSolMintTransaction(ctx context.Context, asset dcttypes.Asset) (uint64, error)
	SetFirstPendingSolMintTransaction(ctx context.Context, asset dcttypes.Asset, id uint64) error
	GetFirstPendingEthMintTransaction(ctx context.Context, asset dcttypes.Asset) (uint64, error)
	SetFirstPendingEthMintTransaction(ctx context.Context, asset dcttypes.Asset, id uint64) error
	GetFirstPendingStakeTransaction(ctx context.Context, asset dcttypes.Asset) (uint64, error)
	SetFirstPendingStakeTransaction(ctx context.Context, asset dcttypes.Asset, id uint64) error

	GetSupply(ctx context.Context, asset dcttypes.Asset) (dcttypes.Supply, error)
	SetSupply(ctx context.Context, supply dcttypes.Supply) error
	GetExchangeRate(ctx context.Context, asset dcttypes.Asset) (math.LegacyDec, error)

	// Burn events
	CreateBurnEvent(ctx context.Context, asset dcttypes.Asset, burnEvent *dcttypes.BurnEvent) (uint64, error)
	WalkBurnEvents(ctx context.Context, asset dcttypes.Asset, fn func(id uint64, burnEvent dcttypes.BurnEvent) (stop bool, err error)) error

	// Redemptions
	SetRedemption(ctx context.Context, asset dcttypes.Asset, id uint64, redemption dcttypes.Redemption) error
	GetRedemption(ctx context.Context, asset dcttypes.Asset, id uint64) (dcttypes.Redemption, error)
	WalkRedemptions(ctx context.Context, asset dcttypes.Asset, fn func(id uint64, redemption dcttypes.Redemption) (stop bool, err error)) error
	WalkRedemptionsDescending(ctx context.Context, asset dcttypes.Asset, fn func(id uint64, redemption dcttypes.Redemption) (stop bool, err error)) error
	GetFirstRedemptionAwaitingSign(ctx context.Context, asset dcttypes.Asset) (uint64, error)
	SetFirstRedemptionAwaitingSign(ctx context.Context, asset dcttypes.Asset, id uint64) error
}
