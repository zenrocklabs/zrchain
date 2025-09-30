package client

import (
	"context"

	identitytypes "github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	"google.golang.org/grpc"
)

// IdentityQueryClient provides a client interface for interacting with the treasury module.
// It wraps the auto-generated treasury QueryClient to provide a more ergonomic interface
// for treasury-related queries.
type IdentityQueryClient struct {
	client identitytypes.QueryClient
}

// NewIdentityQueryClient returns a new IdentityQueryClient with the supplied GRPC client connection.
// This is the main constructor for creating a treasury query client instance.
//
// Parameters:
//   - c: A gRPC client connection to a Zenrock node
//
// Returns:
//   - *IdentityQueryClient: A new identity query client instance
func NewIdentityQueryClient(c *grpc.ClientConn) *IdentityQueryClient {
	return &IdentityQueryClient{
		client: identitytypes.NewQueryClient(c),
	}
}

// Workspaces retrieves a paginated list of workspaces.
//
// Parameters:
//   - ctx: Context for the request, can be used for timeouts and cancellation
//   - page: Pagination parameters for the request
//
// Returns:
//   - []*identitytypes.Workspace: Slice of workspaces
//   - error: An error if the query fails
func (t *IdentityQueryClient) Workspaces(ctx context.Context, page *PageRequest) ([]*identitytypes.Workspace, error) {
	res, err := t.client.Workspaces(ctx, &identitytypes.QueryWorkspacesRequest{
		Pagination: page,
	})
	if err != nil {
		return nil, err
	}

	return res.Workspaces, nil
}

func (t *IdentityQueryClient) WorkspaceByAddress(ctx context.Context, workspaceAddr string) (*identitytypes.Workspace, error) {
	res, err := t.client.WorkspaceByAddress(ctx, &identitytypes.QueryWorkspaceByAddressRequest{
		WorkspaceAddr: workspaceAddr,
	})
	if err != nil {
		return nil, err
	}
	return res.Workspace, nil
}

func (t *IdentityQueryClient) Keyrings(ctx context.Context, page *PageRequest) ([]*identitytypes.Keyring, error) {
	res, err := t.client.Keyrings(ctx, &identitytypes.QueryKeyringsRequest{
		Pagination: page,
	})
	if err != nil {
		return nil, err
	}
	return res.Keyrings, nil
}

func (t *IdentityQueryClient) KeyringByAddress(ctx context.Context, keyringAddr string) (*identitytypes.Keyring, error) {
	res, err := t.client.KeyringByAddress(ctx, &identitytypes.QueryKeyringByAddressRequest{
		KeyringAddr: keyringAddr,
	})
	if err != nil {
		return nil, err
	}
	return res.Keyring, nil
}
