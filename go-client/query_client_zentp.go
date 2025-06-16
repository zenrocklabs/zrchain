package client

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	"google.golang.org/grpc"
)

type ZenTPQueryClient struct {
	client types.QueryClient
}

func NewZenTPQueryClient(c *grpc.ClientConn) *ZenTPQueryClient {
	return &ZenTPQueryClient{
		client: types.NewQueryClient(c),
	}
}

func (c *ZenTPQueryClient) Params(ctx context.Context) (*types.QueryParamsResponse, error) {
	return c.client.Params(ctx, &types.QueryParamsRequest{})
}

func (c *ZenTPQueryClient) Burns(ctx context.Context, recipientAddress string, sourceTxHash string) (*types.QueryBurnsResponse, error) {
	return c.client.Burns(ctx, &types.QueryBurnsRequest{
		RecipientAddress: recipientAddress,
		SourceTxHash:     sourceTxHash,
	})
}
