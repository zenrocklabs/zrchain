package client

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	"google.golang.org/grpc"
)

// DCTQueryClient provides a client interface for interacting with the dct module.
// It wraps the auto-generated dct QueryClient to provide a more ergonomic interface
// for querying DCT-related information, such as lock transactions across multiple assets.
type DCTQueryClient struct {
	client types.QueryClient
}

// NewDCTQueryClient returns a new DCTQueryClient with the supplied GRPC client connection.
// This is the main constructor for creating a DCT query client instance.
//
// Parameters:
//   - c: A gRPC client connection to a Zenrock node
//
// Returns:
//   - *DCTQueryClient: A new DCT query client instance
func NewDCTQueryClient(c *grpc.ClientConn) *DCTQueryClient {
	return &DCTQueryClient{
		client: types.NewQueryClient(c),
	}
}

// LockTransactions retrieves all lock transactions in the DCT module for a specific asset.
// These transactions represent digital assets that have been locked in the system.
//
// Parameters:
//   - ctx: Context for the request, can be used for timeouts and cancellation
//   - asset: The asset type (e.g., ASSET_ZENBTC, ASSET_ZENZEC)
//
// Returns:
//   - *types.QueryLockTransactionsResponse: Contains the list of lock transactions
//   - error: An error if the query fails
//
// Example:
//
//	locks, err := client.LockTransactions(context.Background(), types.ASSET_ZENBTC)
//	if err != nil {
//	    // Handle error
//	}
func (c *DCTQueryClient) LockTransactions(ctx context.Context, asset types.Asset) (*types.QueryLockTransactionsResponse, error) {
	return c.client.GetLockTransactions(ctx, &types.QueryLockTransactionsRequest{Asset: asset})
}

// Redemptions retrieves redemptions for a specific asset.
//
// Parameters:
//   - ctx: Context for the request
//   - asset: The asset type (e.g., ASSET_ZENBTC, ASSET_ZENZEC)
//   - startIndex: Starting index for pagination
//   - status: Filter by redemption status
//
// Returns:
//   - *types.QueryRedemptionsResponse: Contains the list of redemptions
//   - error: An error if the query fails
func (c *DCTQueryClient) Redemptions(ctx context.Context, asset types.Asset, startIndex uint64, status types.RedemptionStatus) (*types.QueryRedemptionsResponse, error) {
	return c.client.GetRedemptions(ctx, &types.QueryRedemptionsRequest{Asset: asset, StartIndex: startIndex, Status: status})
}

// Params retrieves the current parameters for the DCT module.
//
// Parameters:
//   - ctx: Context for the request
//
// Returns:
//   - *types.QueryParamsResponse: Contains the DCT module parameters
//   - error: An error if the query fails
func (c *DCTQueryClient) Params(ctx context.Context) (*types.QueryParamsResponse, error) {
	return c.client.QueryParams(ctx, &types.QueryParamsRequest{})
}

// BurnEvents retrieves burn events for a specific asset.
//
// Parameters:
//   - ctx: Context for the request
//   - asset: The asset type (e.g., ASSET_ZENBTC, ASSET_ZENZEC)
//   - startIndex: Starting index for pagination
//   - txID: Filter by transaction ID (empty string for no filter)
//   - logIndex: Filter by log index
//   - chainID: Filter by CAIP-2 chain ID (empty string for no filter)
//   - status: Filter by burn status
//
// Returns:
//   - *types.QueryBurnEventsResponse: Contains the list of burn events
//   - error: An error if the query fails
func (c *DCTQueryClient) BurnEvents(ctx context.Context, asset types.Asset, startIndex uint64, txID string, logIndex uint64, chainID string) (*types.QueryBurnEventsResponse, error) {
	return c.client.QueryBurnEvents(ctx, &types.QueryBurnEventsRequest{
		Asset:        asset,
		StartIndex:   startIndex,
		TxID:         txID,
		LogIndex:     logIndex,
		Caip2ChainID: chainID,
	})
}

// PendingMintTransactions retrieves pending mint transactions for a specific asset.
//
// Parameters:
//   - ctx: Context for the request
//   - asset: The asset type (e.g., ASSET_ZENBTC, ASSET_ZENZEC)
//   - startIndex: Starting index for pagination
//   - status: Filter by mint transaction status
//
// Returns:
//   - *types.QueryPendingMintTransactionsResponse: Contains the list of pending mint transactions
//   - error: An error if the query fails
func (c *DCTQueryClient) PendingMintTransactions(ctx context.Context, asset types.Asset, startIndex uint64, status types.MintTransactionStatus) (*types.QueryPendingMintTransactionsResponse, error) {
	return c.client.QueryPendingMintTransactions(ctx, &types.QueryPendingMintTransactionsRequest{
		Asset:      asset,
		StartIndex: startIndex,
		Status:     status,
	})
}

// PendingMintTransaction retrieves a specific pending mint transaction by hash.
//
// Parameters:
//   - ctx: Context for the request
//   - asset: The asset type (e.g., ASSET_ZENBTC, ASSET_ZENZEC)
//   - txHash: The transaction hash to look up
//
// Returns:
//   - *types.QueryPendingMintTransactionResponse: Contains the pending mint transaction details
//   - error: An error if the query fails
func (c *DCTQueryClient) PendingMintTransaction(ctx context.Context, asset types.Asset, txHash string) (*types.QueryPendingMintTransactionResponse, error) {
	return c.client.QueryPendingMintTransaction(ctx, &types.QueryPendingMintTransactionRequest{
		Asset:  asset,
		TxHash: txHash,
	})
}

// Supply retrieves the current supply information for DCT assets.
//
// Parameters:
//   - ctx: Context for the request
//   - asset: The asset type (e.g., ASSET_ZENBTC, ASSET_ZENZEC)
//
// Returns:
//   - *types.QuerySupplyResponse: Contains the supply information
//   - error: An error if the query fails
func (c *DCTQueryClient) Supply(ctx context.Context, asset types.Asset) (*types.QuerySupplyResponse, error) {
	return c.client.QuerySupply(ctx, &types.QuerySupplyRequest{Asset: asset})
}
