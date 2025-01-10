package shared

import (
	"context"

	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

type ZenBTCKeeper interface {
	GetMinterKeyID(ctx context.Context) uint64
	GetUnstakerKeyID(ctx context.Context) uint64
	GetEthBatcherAddr(ctx context.Context) string
	GetBitcoinProxyAddress(ctx context.Context) string
	GetPendingMintTransactions(ctx context.Context) (types.PendingMintTransactions, error)
	GetZenBTCExchangeRate(ctx context.Context) (float64, error)
}
