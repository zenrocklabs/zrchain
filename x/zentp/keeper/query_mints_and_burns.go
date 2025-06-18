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

	keys, pageRes, err := k.queryBridge(goCtx, k.mintStore, req.Pagination, req.Creator, req.Denom, "", req.SourceTxHash, req.Status, req.Id, req.TxId)
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

	keys, pageRes, err := k.queryBridge(goCtx, k.burnStore, req.Pagination, "", req.Denom, req.RecipientAddress, req.SourceTxHash, req.Status, req.Id, req.TxId)
	if err != nil {
		return nil, err
	}

	return &types.QueryBurnsResponse{
		Burns:      keys,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) queryBridge(goCtx context.Context, store collections.Map[uint64, types.Bridge], pagination *query.PageRequest, creator, denom, recipientAddress, sourceTxHash string, status types.BridgeStatus, id, txId uint64) ([]*types.Bridge, *query.PageResponse, error) {
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

			txIdMatch := txId == 0 || txId == value.TxId

			recipientAddressMatch := recipientAddress == "" || recipientAddress == value.RecipientAddress

			sourceTxHashMatch := sourceTxHash == "" || sourceTxHash == value.TxHash

			return statusMatch && creatorMatch && denomMatch && idMatch && txIdMatch && recipientAddressMatch && sourceTxHashMatch, nil
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
