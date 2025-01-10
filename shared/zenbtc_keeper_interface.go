package shared

import "context"

type ZenBTCKeeper interface {
	GetMinterKeyID(ctx context.Context) uint64
	GetUnstakerKeyID(ctx context.Context) uint64
	GetEthBatcherAddr(ctx context.Context) string
	GetBitcoinProxyAddress(ctx context.Context) string
}
