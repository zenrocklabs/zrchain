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
func NewValidationQueryClient(c *grpc.ClientConn) *ValidationQueryClient {
	return &ValidationQueryClient{
		client: types.NewQueryClient(c),
	}
}

// ActiveSetValidators retrieves a paginated list of currently bonded (active) validators.
// These validators are participating in block production and consensus.
//
// Parameters:
//   - ctx: Context for the request, can be used for timeouts and cancellation
//   - pageRequest: Pagination parameters for the request
//
// Returns:
//   - *types.QueryValidatorsResponse: Contains the list of active validators and pagination info
//   - error: An error if the query fails
//
// Example:
//
//	validators, err := client.ActiveSetValidators(context.Background(), &query.PageRequest{
//	    Limit: 10,
//	    Offset: 0,
//	})
//	if err != nil {
//	    // Handle error
//	}
func (c *ValidationQueryClient) ActiveSetValidators(ctx context.Context, pageRequest *query.PageRequest) (*types.QueryValidatorsResponse, error) {
	return c.client.Validators(ctx, &types.QueryValidatorsRequest{
		Status:     types.Bonded.String(),
		Pagination: pageRequest,
	})
}
