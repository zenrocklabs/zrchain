package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/pkg/errors"
)

func (k Keeper) Mints(goCtx context.Context, req *types.QueryMintsRequest) (*types.QueryMintsResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}

	keys, pageRes, err := k.queryBridge(goCtx, req.Pagination, req.Creator, req.Denom, req.Status)
	if err != nil {
		return nil, err
	}

	return &types.QueryMintsResponse{
		Mints:      keys,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) Burns(goCtx context.Context, req *types.QueryBurnsRequest) (*types.QueryBurnsResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}

	keys, pageRes, err := k.queryBridge(goCtx, req.Pagination, "", req.Denom, req.Status)
	if err != nil {
		return nil, err
	}

	return &types.QueryBurnsResponse{
		Burns:      keys,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) queryBridge(goCtx context.Context, pagination *query.PageRequest, creator, denom string, status types.BridgeStatus) ([]*types.Bridge, *query.PageResponse, error) {
	keys, pageRes, err := query.CollectionFilteredPaginate(
		goCtx,
		k.mintStore,
		pagination,
		func(key uint64, value types.Bridge) (bool, error) {
			statusMatch := status == types.BridgeStatus_BRIDGE_STATUS_UNSPECIFIED ||
				status == value.State

			creatorMatch := creator == "" || creator == value.Creator

			denomMatch := denom == "" || denom == value.Denom

			return statusMatch && creatorMatch && denomMatch, nil
		},
		func(key uint64, value types.Bridge) (*types.Bridge, error) {
			return &value, nil
		},
	)
	if err != nil {
		return nil, nil, err
	}

	return keys, pageRes, err
}
