package client

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v5/x/validation/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc"
)

// ValidationQueryClient provides a client interface for interacting with the validation module.
// It wraps the auto-generated validation QueryClient to provide a more ergonomic interface
// for validator-related queries.
type ValidationQueryClient struct {
	client types.QueryClient
}

// NewValidationQueryClient returns a new ValidationQueryClient with the supplied GRPC client connection.
// This is the main constructor for creating a validation query client instance.
func NewValidationQueryClient(c *grpc.ClientConn) *ValidationQueryClient {
	return &ValidationQueryClient{
		client: types.NewQueryClient(c),
	}
}

// BondedValidators retrieves a paginated list of currently bonded validators.
// These validators are actively participating in block production and consensus.
func (c *ValidationQueryClient) BondedValidators(ctx context.Context, pageRequest *query.PageRequest) (*types.QueryValidatorsResponse, error) {
	return c.client.Validators(ctx, &types.QueryValidatorsRequest{
		Status:     types.Bonded.String(),
		Pagination: pageRequest,
	})
}

// UnbondedValidators retrieves a paginated list of currently unbonded validators.
// These validators are not currently participating in block production or consensus.
func (c *ValidationQueryClient) UnbondedValidators(ctx context.Context, pageRequest *query.PageRequest) (*types.QueryValidatorsResponse, error) {
	return c.client.Validators(ctx, &types.QueryValidatorsRequest{
		Status:     types.Unbonded.String(),
		Pagination: pageRequest,
	})
}

// UnbondingValidators retrieves a paginated list of currently unbonding validators.
// These validators are in the process of transitioning from bonded to unbonded status.
func (c *ValidationQueryClient) UnbondingValidators(ctx context.Context, pageRequest *query.PageRequest) (*types.QueryValidatorsResponse, error) {
	return c.client.Validators(ctx, &types.QueryValidatorsRequest{
		Status:     types.Unbonding.String(),
		Pagination: pageRequest,
	})
}
