package client

import (
	"context"

	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
	"google.golang.org/grpc"
)

// ZenBTCQueryClient provides a client interface for interacting with the zenbtc module.
// It wraps the auto-generated zenbtc QueryClient to provide a more ergonomic interface
// for querying ZenBTC-related information, such as lock transactions.
type ZenBTCQueryClient struct {
	client types.QueryClient
}

// NewZenBTCQueryClient returns a new ZenBTCQueryClient with the supplied GRPC client connection.
// This is the main constructor for creating a ZenBTC query client instance.
//
// Parameters:
//   - c: A gRPC client connection to a Zenrock node
//
// Returns:
//   - *ZenBTCQueryClient: A new ZenBTC query client instance
func NewZenBTCQueryClient(c *grpc.ClientConn) *ZenBTCQueryClient {
	return &ZenBTCQueryClient{
		client: types.NewQueryClient(c),
	}
}

// LockTransactions retrieves all lock transactions in the ZenBTC module.
// These transactions represent Bitcoin that has been locked in the system.
//
// Parameters:
//   - ctx: Context for the request, can be used for timeouts and cancellation
//
// Returns:
//   - *types.QueryLockTransactionsResponse: Contains the list of lock transactions
//   - error: An error if the query fails
//
// Example:
//
//	locks, err := client.LockTransactions(context.Background())
//	if err != nil {
//	    // Handle error
//	}
func (c *ZenBTCQueryClient) LockTransactions(ctx context.Context) (*types.QueryLockTransactionsResponse, error) {
	return c.client.GetLockTransactions(ctx, &types.QueryLockTransactionsRequest{})
}

func (c *ZenBTCQueryClient) Redemptions(ctx context.Context, startIndex uint64, status types.RedemptionStatus) (*types.QueryRedemptionsResponse, error) {
	return c.client.GetRedemptions(ctx, &types.QueryRedemptionsRequest{StartIndex: startIndex, Status: status})
}

func (c *ZenBTCQueryClient) Params(ctx context.Context) (*types.QueryParamsResponse, error) {
	return c.client.QueryParams(ctx, &types.QueryParamsRequest{})
}

func (c *ZenBTCQueryClient) BurnEvents(ctx context.Context, startIndex uint64, txID string, logIndex uint64, chainID string) (*types.QueryBurnEventsResponse, error) {
	return c.client.QueryBurnEvents(ctx, &types.QueryBurnEventsRequest{StartIndex: startIndex, TxID: txID, LogIndex: logIndex, ChainID: chainID})
}
