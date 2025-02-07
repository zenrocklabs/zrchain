package shared

import (
	"context"

	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

type ZenBTCKeeper interface {
	GetStakerKeyID(ctx context.Context) uint64
	GetEthMinterKeyID(ctx context.Context) uint64
	GetUnstakerKeyID(ctx context.Context) uint64
	GetCompleterKeyID(ctx context.Context) uint64
	GetEthBatcherAddr(ctx context.Context) string
	GetEthTokenAddr(ctx context.Context) string
	GetBitcoinProxyAddress(ctx context.Context) string
	SetPendingMintTransaction(ctx context.Context, pendingMintTransaction types.PendingMintTransaction) error
	WalkPendingMintTransactions(ctx context.Context, fn func(id uint64, pendingMintTransaction types.PendingMintTransaction) (stop bool, err error)) error
	GetSupply(ctx context.Context) (types.Supply, error)
	SetSupply(ctx context.Context, supply types.Supply) error
	HasRedemption(ctx context.Context, id uint64) (bool, error)
	SetRedemption(ctx context.Context, id uint64, redemption types.Redemption) error
	WalkRedemptions(ctx context.Context, fn func(id uint64, redemption types.Redemption) (stop bool, err error)) error
	GetExchangeRate(ctx context.Context) (float64, error)
	GetBurnEvents(ctx context.Context) (types.BurnEvents, error)
	SetBurnEvents(ctx context.Context, burnEvents types.BurnEvents) error
}
