package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueryPendingMintTransactions(ctx context.Context, req *types.QueryPendingMintTransactionsRequest) (*types.QueryPendingMintTransactionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var pendingMintResponses []*types.PendingMintTransaction

	// Walk through all assets if asset is unspecified, otherwise filter by asset
	asset := req.Asset
	if asset == types.Asset_ASSET_UNSPECIFIED {
		// Query all pending mint transactions across all assets
		if err := k.PendingMintTransactions.Walk(ctx, nil, func(key collections.Pair[string, uint64], mint types.PendingMintTransaction) (bool, error) {
			if key.K2() < req.StartIndex {
				return false, nil
			}
			switch req.Status {
			case types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED:
				if mint.Status == types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED {
					pendingMintResponses = append(pendingMintResponses, &mint)
				}
			case types.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED:
				if mint.Status == types.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED {
					pendingMintResponses = append(pendingMintResponses, &mint)
				}
			default: // don't filter by status
				pendingMintResponses = append(pendingMintResponses, &mint)
			}
			return false, nil
		}); err != nil {
			return nil, err
		}
	} else {
		// Query pending mints for a specific asset
		if err := k.WalkPendingMintTransactions(ctx, asset, func(id uint64, mint types.PendingMintTransaction) (bool, error) {
			if id < req.StartIndex {
				return false, nil
			}
			switch req.Status {
			case types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED:
				if mint.Status == types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED {
					pendingMintResponses = append(pendingMintResponses, &mint)
				}
			case types.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED:
				if mint.Status == types.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED {
					pendingMintResponses = append(pendingMintResponses, &mint)
				}
			default: // don't filter by status
				pendingMintResponses = append(pendingMintResponses, &mint)
			}
			return false, nil
		}); err != nil {
			return nil, err
		}
	}

	return &types.QueryPendingMintTransactionsResponse{
		PendingMintTransactions: pendingMintResponses,
	}, nil
}
