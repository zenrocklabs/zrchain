package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Dcts(goCtx context.Context, req *types.QueryDctsRequest) (*types.QueryDctsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	dcts, pageRes, err := query.CollectionFilteredPaginate(
		goCtx,
		k.DctStore,
		req.Pagination,
		func(key string, value types.Dct) (bool, error) {
			return true, nil
		},
		func(key string, value types.Dct) (*types.Dct, error) {
			return &value, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.QueryDctsResponse{
		Dct:        dcts,
		Pagination: pageRes,
	}, nil
}
