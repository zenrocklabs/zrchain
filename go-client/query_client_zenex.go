package client

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	"google.golang.org/grpc"
)

type ZenexQueryClient struct {
	client types.QueryClient
}

func NewZenexQueryClient(c *grpc.ClientConn) *ZenexQueryClient {
	return &ZenexQueryClient{
		client: types.NewQueryClient(c),
	}
}

func (c *ZenexQueryClient) Params(ctx context.Context) (*types.QueryParamsResponse, error) {
	return c.client.Params(ctx, &types.QueryParamsRequest{})
}

func (c *ZenexQueryClient) Swaps(ctx context.Context, creator, pair, workspace, sourceTxHash string, swapId uint64, status types.SwapStatus) (*types.QuerySwapsResponse, error) {
	return c.client.Swaps(ctx, &types.QuerySwapsRequest{
		Creator:      creator,
		SwapId:       swapId,
		Status:       status,
		Pair:         pair,
		Workspace:    workspace,
		SourceTxHash: sourceTxHash,
	})
}

func (c *ZenexQueryClient) RockPool(ctx context.Context) (*types.QueryRockPoolResponse, error) {
	return c.client.RockPool(ctx, &types.QueryRockPoolRequest{})
}
