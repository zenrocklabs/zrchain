package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Swaps(goCtx context.Context, req *types.QuerySwapsRequest) (*types.QuerySwapsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	swaps, pageRes, err := k.querySwap(goCtx, k.SwapsStore, req.Pagination, req.Creator, req.Pair, req.Workspace, req.Status, req.SwapId)
	if err != nil {
		return nil, err
	}

	return &types.QuerySwapsResponse{
		Swaps:      swaps,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) querySwap(goCtx context.Context, store collections.Map[uint64, types.Swap], pagination *query.PageRequest, creator, pair, workspace string, status types.SwapStatus, swapId uint64) ([]*types.Swap, *query.PageResponse, error) {
	swaps, pageRes, err := query.CollectionFilteredPaginate(
		goCtx,
		store,
		pagination,
		func(key uint64, value types.Swap) (bool, error) {
			statusMatch := status == types.SwapStatus_SWAP_STATUS_UNSPECIFIED || status == value.Status

			creatorMatch := creator == "" || creator == value.Creator

			pairMatch := pair == "" || pair == value.Pair

			workspaceMatch := workspace == "" || workspace == value.Workspace

			swapIdMatch := swapId == 0 || swapId == value.SwapId

			return statusMatch && creatorMatch && pairMatch && workspaceMatch && swapIdMatch, nil
		},
		func(key uint64, value types.Swap) (*types.Swap, error) {
			return &value, nil
		},
	)

	return swaps, pageRes, err
}
