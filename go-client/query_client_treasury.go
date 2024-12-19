package client

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc"
)

// PageRequest is an alias for the Cosmos SDK query.PageRequest type
// Used for pagination in queries that return multiple items
type PageRequest = query.PageRequest

// TreasuryQueryClient provides a client interface for interacting with the treasury module.
// It wraps the auto-generated treasury QueryClient to provide a more ergonomic interface
// for treasury-related queries.
type TreasuryQueryClient struct {
	client types.QueryClient
}

// NewTreasuryQueryClient returns a new TreasuryQueryClient with the supplied GRPC client connection.
// This is the main constructor for creating a treasury query client instance.
//
// Parameters:
//   - c: A gRPC client connection to a Zenrock node
//
// Returns:
//   - *TreasuryQueryClient: A new treasury query client instance
func NewTreasuryQueryClient(c *grpc.ClientConn) *TreasuryQueryClient {
	return &TreasuryQueryClient{
		client: types.NewQueryClient(c),
	}
}

// PendingKeyRequests retrieves a paginated list of pending key requests for a specific keyring address.
//
// Parameters:
//   - ctx: Context for the request, can be used for timeouts and cancellation
//   - page: Pagination parameters for the request
//   - keyringAddr: The address of the keyring to query pending requests for
//
// Returns:
//   - []*types.KeyReqResponse: Slice of pending key requests
//   - error: An error if the query fails
func (t *TreasuryQueryClient) PendingKeyRequests(ctx context.Context, page *PageRequest, keyringAddr string) ([]*types.KeyReqResponse, error) {
	res, err := t.client.KeyRequests(ctx, &types.QueryKeyRequestsRequest{
		Pagination:  page,
		KeyringAddr: keyringAddr,
		Status:      types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
	})
	if err != nil {
		return nil, err
	}

	return res.KeyRequests, nil
}

// GetKeys retrieves a paginated list of all keys in the treasury.
//
// Parameters:
//   - ctx: Context for the request
//   - offset: Starting position in the list
//   - pagesize: Number of items to return
//
// Returns:
//   - *types.QueryKeysResponse: Contains the list of keys and pagination info
//   - error: An error if the query fails
func (t *TreasuryQueryClient) GetKeys(ctx context.Context, offset uint64, pagesize uint64) (*types.QueryKeysResponse, error) {
	res, err := t.client.Keys(ctx, &types.QueryKeysRequest{
		Pagination: &query.PageRequest{Offset: offset, Limit: pagesize},
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetKey retrieves details for a specific key by its ID.
//
// Parameters:
//   - ctx: Context for the request
//   - keyId: The unique identifier of the key to query
//
// Returns:
//   - *types.QueryKeyByIDResponse: Contains the key details
//   - error: An error if the query fails or key is not found
func (t *TreasuryQueryClient) GetKey(ctx context.Context, keyId uint64) (*types.QueryKeyByIDResponse, error) {
	res, err := t.client.KeyByID(ctx, &types.QueryKeyByIDRequest{
		Id: keyId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetKey retrieves the last Key ID
// Parameters:
//   - ctx: Context for the request
//
// Returns:
//   - *types.QueryKeyByIDResponse: Contains the key details
//   - error: An error if the query fails or key is not found
func (t *TreasuryQueryClient) GetLastKey(ctx context.Context) (*types.QueryKeysResponse, error) {
	res, err := t.client.Keys(ctx, &types.QueryKeysRequest{
		Pagination: &query.PageRequest{Limit: 1, Reverse: true},
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetKeyRequest retrieves details for a specific key request by its ID.
//
// Parameters:
//   - ctx: Context for the request
//   - requestID: The unique identifier of the key request
//
// Returns:
//   - *types.KeyReqResponse: Contains the key request details
//   - error: An error if the query fails or request is not found
func (t *TreasuryQueryClient) GetKeyRequest(ctx context.Context, requestID uint64) (*types.KeyReqResponse, error) {
	res, err := t.client.KeyRequestByID(ctx, &types.QueryKeyRequestByIDRequest{
		Id: requestID,
	})
	if err != nil {
		return nil, err
	}

	return res.KeyRequest, nil
}

// PendingSignatureRequests retrieves a paginated list of pending signature requests for a specific keyring address.
//
// Parameters:
//   - ctx: Context for the request
//   - page: Pagination parameters
//   - keyringAddr: The address of the keyring to query pending requests for
//
// Returns:
//   - []*types.SignReqResponse: Slice of pending signature requests
//   - error: An error if the query fails
func (t *TreasuryQueryClient) PendingSignatureRequests(ctx context.Context, page *PageRequest, keyringAddr string) ([]*types.SignReqResponse, error) {
	res, err := t.client.SignatureRequests(ctx, &types.QuerySignatureRequestsRequest{
		Pagination:  page,
		KeyringAddr: keyringAddr,
		Status:      types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
	})
	if err != nil {
		return nil, err
	}

	return res.SignRequests, nil
}

// FulfilledSignatureRequests retrieves a paginated list of fulfilled signature requests.
//
// Parameters:
//   - ctx: Context for the request
//   - offset: Starting position in the list
//   - pagesize: Number of items to return
//
// Returns:
//   - *types.QuerySignatureRequestsResponse: Contains the list of fulfilled requests
//   - error: An error if the query fails
func (t *TreasuryQueryClient) FulfilledSignatureRequests(ctx context.Context, offset uint64, pagesize uint64) (*types.QuerySignatureRequestsResponse, error) {
	res, err := t.client.SignatureRequests(ctx, &types.QuerySignatureRequestsRequest{
		Pagination: &query.PageRequest{Offset: offset, Limit: pagesize},
		Status:     types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FulfilledSignatureRequests retrieves the latesat fulfilled signature request.
//
// Parameters:
//   - ctx: Context for the request
//
// Returns:
//   - *types.QuerySignatureRequestsResponse: Contains the list of fulfilled requests
//   - error: An error if the query fails
func (t *TreasuryQueryClient) LastFulfilledSignatureRequest(ctx context.Context) (*types.QuerySignatureRequestsResponse, error) {
	res, err := t.client.SignatureRequests(ctx, &types.QuerySignatureRequestsRequest{
		Pagination: &query.PageRequest{Limit: 1, Reverse: true},
		Status:     types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetSignatureRequest retrieves details for a specific signature request by its ID.
//
// Parameters:
//   - ctx: Context for the request
//   - requestID: The unique identifier of the signature request
//
// Returns:
//   - *types.SignReqResponse: Contains the signature request details
//   - error: An error if the query fails or request is not found
func (t *TreasuryQueryClient) GetSignatureRequest(ctx context.Context, requestID uint64) (*types.SignReqResponse, error) {
	res, err := t.client.SignatureRequestByID(ctx, &types.QuerySignatureRequestByIDRequest{
		Id: requestID,
	})
	if err != nil {
		return nil, err
	}

	return res.SignRequest, nil
}

// SignedTransactions retrieves a paginated list of fulfilled signature requests for a specific wallet type.
//
// Parameters:
//   - ctx: Context for the request
//   - page: Pagination parameters
//   - walletType: Type of wallet to filter transactions
//
// Returns:
//   - *types.QuerySignTransactionRequestsResponse: Contains the list of signed transactions
//   - error: An error if the query fails
func (t *TreasuryQueryClient) SignedTransactions(ctx context.Context, page *PageRequest, walletType types.WalletType) (*types.QuerySignTransactionRequestsResponse, error) {
	res, err := t.client.SignTransactionRequests(ctx, &types.QuerySignTransactionRequestsRequest{
		Pagination: page,
		WalletType: walletType,
		Status:     types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ZrSignKeys retrieves ZrSign keys for a specific address and wallet type.
//
// Parameters:
//   - ctx: Context for the request
//   - page: Pagination parameters
//   - address: The address to query keys for
//   - walletType: Type of wallet to filter keys
//
// Returns:
//   - *types.QueryZrSignKeysResponse: Contains the ZrSign keys
//   - error: An error if the query fails
func (t *TreasuryQueryClient) ZrSignKeys(ctx context.Context, page *PageRequest, address, walletType string) (*types.QueryZrSignKeysResponse, error) {
	res, err := t.client.ZrSignKeys(ctx, &types.QueryZrSignKeysRequest{
		Address:    address,
		WalletType: walletType,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetZenbtcWallets retrieves Zenbtc wallets for a specific address and wallet type.
//
// Parameters:
//   - ctx: Context for the request
//   - page: Pagination parameters
//   - mintChainId: The chain ID to filter wallets
//   - chainType: The type of wallet to filter
//   - recipientAddr: The recipient address to filter wallets
//   - returnAddr: The return address to filter wallets
//
// Returns:
//   - *types.QueryZenbtcWalletsResponse: Contains the Zenbtc wallets
//   - error: An error if the query fails
func (t *TreasuryQueryClient) GetZenbtcWallets(ctx context.Context, page *PageRequest, recipientAddr, returnAddr string, mintChainId uint64, chainType types.WalletType) (*types.QueryZenbtcWalletsResponse, error) {
	res, err := t.client.ZenbtcWallets(ctx, &types.QueryZenbtcWalletsRequest{
		RecipientAddr: recipientAddr,
		ChainType:     chainType,
		ReturnAddr:    returnAddr,
		MintChainId:   mintChainId,
		Pagination:    page,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
