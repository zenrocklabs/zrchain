package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/pkg/errors"
)

func (k Keeper) Mints(goCtx context.Context, req *types.QueryMintsRequest) (*types.QueryMintsResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}

	keys, pageRes, err := k.queryBridge(goCtx, k.mintStore, req.Pagination, req.Creator, req.Denom, req.Status, req.Id)
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

	keys, pageRes, err := k.queryBridge(goCtx, k.burnStore, req.Pagination, "", req.Denom, req.Status, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryBurnsResponse{
		Burns:      keys,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) queryBridge(goCtx context.Context, store collections.Map[uint64, types.Bridge], pagination *query.PageRequest, creator, denom string, status types.BridgeStatus, id uint64) ([]*types.Bridge, *query.PageResponse, error) {
	keys, pageRes, err := query.CollectionFilteredPaginate(
		goCtx,
		store,
		pagination,
		func(key uint64, value types.Bridge) (bool, error) {
			statusMatch := status == types.BridgeStatus_BRIDGE_STATUS_UNSPECIFIED ||
				status == value.State

			creatorMatch := creator == "" || creator == value.Creator

			denomMatch := denom == "" || denom == value.Denom

			idMatch := id == 0 || id == key

			return statusMatch && creatorMatch && denomMatch && idMatch, nil
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

func (k Keeper) Stats(goCtx context.Context, req *types.QueryStatsRequest) (*types.QueryStatsResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}

	mintKeys, _, err := k.queryBridge(goCtx, k.mintStore, nil, "", req.Denom, types.BridgeStatus_BRIDGE_STATUS_COMPLETED, 0)
	if err != nil {
		return nil, err
	}

	burnKeys, _, err := k.queryBridge(goCtx, k.burnStore, nil, "", req.Denom, types.BridgeStatus_BRIDGE_STATUS_COMPLETED, 0)
	if err != nil {
		return nil, err
	}

	var totalMints uint64
	for _, mint := range mintKeys {

		if req.Address == "" || mint.Creator == req.Address {
			totalMints += mint.Amount
		}
	}

	var totalBurns uint64
	for _, burn := range burnKeys {

		if req.Address == "" || burn.RecipientAddress == req.Address {
			totalBurns += burn.Amount
		}
	}

	return &types.QueryStatsResponse{
		TotalMints: totalMints,
		TotalBurns: totalBurns,
	}, nil
}
