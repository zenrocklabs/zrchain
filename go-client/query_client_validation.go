package client

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v5/x/validation/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc"
)

// ValidationQueryClient is the client for the validation module.
type ValidationQueryClient struct {
	client types.QueryClient
}

func NewValidationQueryClient(c *grpc.ClientConn) *ValidationQueryClient {
	return &ValidationQueryClient{
		client: types.NewQueryClient(c),
	}
}

func (c *ValidationQueryClient) ActiveSetValidators(ctx context.Context, pageRequest *query.PageRequest) (*types.QueryValidatorsResponse, error) {
	return c.client.Validators(ctx, &types.QueryValidatorsRequest{Status: types.Bonded.String(), Pagination: pageRequest})
}
