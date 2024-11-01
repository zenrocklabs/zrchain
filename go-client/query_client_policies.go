package client

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	"google.golang.org/grpc"
)

// TreasuryQueryClient is the client for the treasury module.
type PolicyQueryClient struct {
	client types.QueryClient
}

func NewPolicyQueryClient(c *grpc.ClientConn) *PolicyQueryClient {
	return &PolicyQueryClient{
		client: types.NewQueryClient(c),
	}
}

func (c *PolicyQueryClient) GetActionDetailsById(ctx context.Context, actionId uint64) (*types.QueryActionDetailsByIdResponse, error) {
	return c.client.ActionDetailsById(ctx, &types.QueryActionDetailsByIdRequest{
		Id: actionId,
	})
}
