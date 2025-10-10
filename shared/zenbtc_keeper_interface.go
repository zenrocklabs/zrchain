package shared

import (
	"context"

	"cosmossdk.io/math"
	zenbtctypes "github.com/Zenrock-Foundation/zrchain/v6/x/zenbtc/types"
)

type ZenBTCKeeper interface {
	GetStakerKeyID(ctx context.Context) uint64
	GetEthMinterKeyID(ctx context.Context) uint64
	GetParams(ctx context.Context) (zenbtctypes.Params, error)
	GetSolanaParams(ctx context.Context) *zenbtctypes.Solana
	GetUnstakerKeyID(ctx context.Context) uint64
	GetCompleterKeyID(ctx context.Context) uint64
	GetRewardsDepositKeyID(ctx context.Context) uint64
	GetChangeAddressKeyIDs(ctx context.Context) []uint64
	GetControllerAddr(ctx context.Context) string
	GetEthTokenAddr(ctx context.Context) string
	GetBitcoinProxyAddress(ctx context.Context) string
	SetPendingMintTransaction(ctx context.Context, pendingMintTransaction zenbtctypes.PendingMintTransaction) error
	WalkPendingMintTransactions(ctx context.Context, fn func(id uint64, pendingMintTransaction zenbtctypes.PendingMintTransaction) (stop bool, err error)) error
	GetPendingMintTransaction(ctx context.Context, id uint64) (zenbtctypes.PendingMintTransaction, error)
	HasPendingMintTransaction(ctx context.Context, id uint64) (bool, error)
	GetSupply(ctx context.Context) (zenbtctypes.Supply, error)
	SetSupply(ctx context.Context, supply zenbtctypes.Supply) error
	HasRedemption(ctx context.Context, id uint64) (bool, error)
	SetRedemption(ctx context.Context, id uint64, redemption zenbtctypes.Redemption) error
	GetRedemption(ctx context.Context, id uint64) (zenbtctypes.Redemption, error)
	WalkRedemptions(ctx context.Context, fn func(id uint64, redemption zenbtctypes.Redemption) (stop bool, err error)) error
	WalkRedemptionsDescending(ctx context.Context, fn func(id uint64, redemption zenbtctypes.Redemption) (stop bool, err error)) error
	GetExchangeRate(ctx context.Context) (math.LegacyDec, error)
	GetBurnEvent(ctx context.Context, id uint64) (zenbtctypes.BurnEvent, error)
	SetBurnEvent(ctx context.Context, id uint64, burnEvent zenbtctypes.BurnEvent) error
	CreateBurnEvent(ctx context.Context, burnEvent *zenbtctypes.BurnEvent) (uint64, error)
	WalkBurnEvents(ctx context.Context, fn func(id uint64, burnEvent zenbtctypes.BurnEvent) (stop bool, err error)) error
	GetFirstPendingEthMintTransaction(ctx context.Context) (uint64, error)
	SetFirstPendingEthMintTransaction(ctx context.Context, id uint64) error
	GetFirstPendingSolMintTransaction(ctx context.Context) (uint64, error)
	SetFirstPendingSolMintTransaction(ctx context.Context, id uint64) error
	GetFirstPendingBurnEvent(ctx context.Context) (uint64, error)
	SetFirstPendingBurnEvent(ctx context.Context, id uint64) error
	GetFirstPendingRedemption(ctx context.Context) (uint64, error)
	SetFirstPendingRedemption(ctx context.Context, id uint64) error
	GetFirstPendingStakeTransaction(ctx context.Context) (uint64, error)
	SetFirstPendingStakeTransaction(ctx context.Context, id uint64) error
	GetFirstRedemptionAwaitingSign(ctx context.Context) (uint64, error)
	SetFirstRedemptionAwaitingSign(ctx context.Context, id uint64) error
}
