package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueryPendingMintTransactions(ctx context.Context, req *types.QueryPendingMintTransactionsRequest) (*types.QueryPendingMintTransactionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var pendingMintResponses []*types.PendingMintTransaction
	var queryRange collections.Range[uint64]

	if err := k.PendingMintTransactionsMap.Walk(ctx, queryRange.StartInclusive(req.StartIndex), func(_ uint64, mint types.PendingMintTransaction) (bool, error) {
		switch req.Status {
		case types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED:
			if mint.Status == types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED {
				pendingMintResponses = append(pendingMintResponses, &mint)
			}
		case types.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED:
			if mint.Status == types.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED {
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

	return &types.QueryPendingMintTransactionsResponse{
		PendingMintTransactions: pendingMintResponses,
	}, nil
}
