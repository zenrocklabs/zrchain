package client

import (
	"context"

	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
	"google.golang.org/grpc"
)

// ZenBTCQueryClient is the client for the zenbtc module.
type ZenBTCQueryClient struct {
	client types.QueryClient
}

func NewZenBTCQueryClient(c *grpc.ClientConn) *ZenBTCQueryClient {
	return &ZenBTCQueryClient{
		client: types.NewQueryClient(c),
	}
}

func (c *ZenBTCQueryClient) LockTransactions(ctx context.Context) (*types.QueryLockTransactionsResponse, error) {
	return c.client.LockTransactions(ctx, &types.QueryLockTransactionsRequest{})
}
