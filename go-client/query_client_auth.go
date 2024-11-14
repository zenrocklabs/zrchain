package client

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc"
)

// Package client provides a gRPC client implementation for interacting with the Cosmos SDK auth module.
// It offers a simplified and opinionated interface for querying account information from a Cosmos-based
// blockchain node.

// AuthQueryClient stores a query client for the zenrock auth module.
// It wraps the auto-generated Cosmos SDK auth QueryClient to provide
// a more ergonomic interface for common auth-related queries.
type AuthQueryClient struct {
	client authtypes.QueryClient
}

// NewAuthQueryClient returns a new AuthQueryClient with the supplied GRPC client connection.
// This is the main constructor for creating an auth query client instance.
//
// Parameters:
//   - c: A gRPC client connection to a Cosmos SDK node
//
// Returns:
//   - *AuthQueryClient: A new auth query client instance
func NewAuthQueryClient(c *grpc.ClientConn) *AuthQueryClient {
	return &AuthQueryClient{
		client: authtypes.NewQueryClient(c),
	}
}

// Account retrieves the auth account information for the supplied address.
// It handles the unmarshaling of the account data and ensures the account
// type is a standard Cosmos SDK BaseAccount.
//
// Parameters:
//   - ctx: Context for the request, can be used for timeouts and cancellation
//   - addr: The bech32 address of the account to query
//
// Returns:
//   - types.AccountI: The account interface containing account details
//   - error: An error if the query fails or if the account type is unsupported
//
// Example:
//
//	account, err := client.Account(context.Background(), "zen1...")
//	if err != nil {
//	    // Handle error
//	}
func (c *AuthQueryClient) Account(ctx context.Context, addr string) (types.AccountI, error) {
	res, err := c.client.Account(ctx, &authtypes.QueryAccountRequest{
		Address: addr,
	})
	if err != nil {
		return nil, err
	}

	// Check account type
	if res.Account.TypeUrl != "/cosmos.auth.v1beta1.BaseAccount" {
		return nil, fmt.Errorf("unknown account type: %s", res.Account.TypeUrl)
	}

	baseAccount := &authtypes.BaseAccount{}
	if err := baseAccount.Unmarshal(res.Account.Value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account: %w", err)
	}
	return baseAccount, nil
}
