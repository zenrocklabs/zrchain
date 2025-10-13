package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueryPendingMintTransaction(ctx context.Context, req *types.QueryPendingMintTransactionRequest) (*types.QueryPendingMintTransactionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var pendingMintResponse *types.PendingMintTransaction

	// Search across all assets for the pending mint transaction with matching tx hash
	if err := k.WalkPendingMintTransactions(ctx, req.Asset, func(_ uint64, mint types.PendingMintTransaction) (bool, error) {
		if mint.TxHash == req.TxHash {
			pendingMintResponse = &mint
			return true, nil
		}
		return false, nil
	}); err != nil {
		return nil, err
	}

	if pendingMintResponse == nil {
		return nil, status.Error(codes.NotFound, "pending mint transaction not found")
	}

	return &types.QueryPendingMintTransactionResponse{
		PendingMintTransaction: pendingMintResponse,
	}, nil
}
