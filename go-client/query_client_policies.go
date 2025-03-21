package client

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
	"google.golang.org/grpc"
)

// Package client provides a gRPC client implementation for interacting with the Zenrock policy module.
// It offers a simplified interface for querying policy-related information and action details from zrchain.

// PolicyQueryClient provides a client interface for interacting with the policy module.
// It wraps the auto-generated policy QueryClient to provide a more ergonomic interface
// for policy-related queries.
type PolicyQueryClient struct {
	client types.QueryClient
}

// NewPolicyQueryClient returns a new PolicyQueryClient with the supplied GRPC client connection.
// This is the main constructor for creating a policy query client instance.
//
// Parameters:
//   - c: A gRPC client connection to a Zenrock node
//
// Returns:
//   - *PolicyQueryClient: A new policy query client instance
func NewPolicyQueryClient(c *grpc.ClientConn) *PolicyQueryClient {
	return &PolicyQueryClient{
		client: types.NewQueryClient(c),
	}
}

// GetActionDetailsById retrieves the details of a specific action by its ID.
// This method allows querying information about policy actions that have been
// submitted to the blockchain.
//
// Parameters:
//   - ctx: Context for the request, can be used for timeouts and cancellation
//   - actionId: The unique identifier of the action to query
//
// Returns:
//   - *types.QueryActionDetailsByIdResponse: Contains the details of the requested action
//   - error: An error if the query fails or if the action is not found
//
// Example:
//
//	details, err := client.GetActionDetailsById(context.Background(), 123)
//	if err != nil {
//	    // Handle error
//	}
func (c *PolicyQueryClient) GetActionDetailsById(ctx context.Context, actionId uint64) (*types.QueryActionDetailsByIdResponse, error) {
	return c.client.ActionDetailsById(ctx, &types.QueryActionDetailsByIdRequest{
		Id: actionId,
	})
}
