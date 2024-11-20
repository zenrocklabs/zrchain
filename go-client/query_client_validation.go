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
//
// Parameters:
//   - c: A gRPC client connection to a Zenrock node
//
// Returns:
//   - *ValidationQueryClient: A new validation query client instance
//
// Example:
//
//	client := NewValidationQueryClient(grpc.DialContext(context.Background(), "localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials())))
func NewValidationQueryClient(c *grpc.ClientConn) *ValidationQueryClient {
	return &ValidationQueryClient{
		client: types.NewQueryClient(c),
	}
}

// UnbondedValidators retrieves a paginated list of currently unbonded validators.
// These validators are not currently participating in block production or consensus.
//
// Parameters:
//   - ctx: Context for the request, can be used for timeouts and cancellation
//   - pageRequest: Pagination parameters for the request
//
// Returns:
//   - *types.QueryValidatorsResponse: Contains the list of unbonded validators and pagination info
//   - error: An error if the query fails
//
// Example:
//
//	validators, err := client.UnbondedValidators(context.Background(), &query.PageRequest{
//	    Limit: 10,
//	    Offset: 0,
//	})
//	if err != nil {
//	    // Handle error
//	}
func (c *ValidationQueryClient) UnbondedValidators(ctx context.Context, pageRequest *query.PageRequest) (*types.QueryValidatorsResponse, error) {
	return c.client.Validators(ctx, &types.QueryValidatorsRequest{
		Status:     types.Unbonded.String(),
		Pagination: pageRequest,
	})
}

// UnbondingValidators retrieves a paginated list of currently unbonding validators.
// These validators are in the process of transitioning from bonded to unbonded status.
//
// Parameters:
//   - ctx: Context for the request, can be used for timeouts and cancellation
//   - pageRequest: Pagination parameters for the request
//
// Returns:
//   - *types.QueryValidatorsResponse: Contains the list of unbonding validators and pagination info
//   - error: An error if the query fails
//
// Example:
//
//	validators, err := client.UnbondingValidators(context.Background(), &query.PageRequest{
//	    Limit: 10,
//	    Offset: 0,
//	})
//	if err != nil {
//	    // Handle error
//	}
func (c *ValidationQueryClient) UnbondingValidators(ctx context.Context, pageRequest *query.PageRequest) (*types.QueryValidatorsResponse, error) {
	return c.client.Validators(ctx, &types.QueryValidatorsRequest{
		Status:     types.Unbonding.String(),
		Pagination: pageRequest,
	})
}
