package keeper

import (
	"context"

	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/pkg/errors"
)

func (k Keeper) Stats(goCtx context.Context, req *types.QueryStatsRequest) (*types.QueryStatsResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}

	pagination := &query.PageRequest{
		CountTotal: true,
	}

	mintKeys, pageResMint, err := k.queryBridge(goCtx, k.mintStore, pagination, req.Address, req.Denom, "", "", types.BridgeStatus_BRIDGE_STATUS_COMPLETED, 0, 0)
	if err != nil {
		return nil, err
	}

	burnKeys, pageResBurn, err := k.queryBridge(goCtx, k.burnStore, pagination, "", req.Denom, req.Address, "", types.BridgeStatus_BRIDGE_STATUS_COMPLETED, 0, 0)
	if err != nil {
		return nil, err
	}

	totalMinted := math.ZeroInt()
	mintsCount := pageResMint.Total
	for _, mint := range mintKeys {
		if req.Address == "" || mint.Creator == req.Address {
			totalMinted = totalMinted.Add(math.NewIntFromUint64(mint.Amount))
		}
	}

	totalBurned := math.ZeroInt()
	burnsCount := pageResBurn.Total
	for _, burn := range burnKeys {
		if req.Address == "" || burn.RecipientAddress == req.Address {
			totalBurned = totalBurned.Add(math.NewIntFromUint64(burn.Amount))
		}
	}

	return &types.QueryStatsResponse{
		TotalMinted: totalMinted.Uint64(),
		TotalBurned: totalBurned.Uint64(),
		MintsCount:  mintsCount,
		BurnsCount:  burnsCount,
	}, nil
}
