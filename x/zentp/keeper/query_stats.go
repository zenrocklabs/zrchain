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

	// Mints
	totalMinted := math.ZeroInt()
	var mintsCount uint64
	mintPagination := &query.PageRequest{
		CountTotal: true,
	}
	for {
		mintKeys, pageResMint, err := k.queryBridge(goCtx, k.MintStore, mintPagination, req.Address, req.Denom, "", "", types.BridgeStatus_BRIDGE_STATUS_COMPLETED, 0, 0)
		if err != nil {
			return nil, err
		}

		if mintPagination.CountTotal {
			mintsCount = pageResMint.Total
		}

		for _, mint := range mintKeys {
			if req.Address == "" || mint.Creator == req.Address {
				totalMinted = totalMinted.Add(math.NewIntFromUint64(mint.Amount))
			}
		}

		if pageResMint.NextKey == nil {
			break
		}
		mintPagination.Key = pageResMint.NextKey
		mintPagination.CountTotal = false
	}

	// Burns
	totalBurned := math.ZeroInt()
	var burnsCount uint64
	burnPagination := &query.PageRequest{
		CountTotal: true,
	}
	for {
		burnKeys, pageResBurn, err := k.queryBridge(goCtx, k.BurnStore, burnPagination, "", req.Denom, req.Address, "", types.BridgeStatus_BRIDGE_STATUS_COMPLETED, 0, 0)
		if err != nil {
			return nil, err
		}

		if burnPagination.CountTotal {
			burnsCount = pageResBurn.Total
		}

		for _, burn := range burnKeys {
			if req.Address == "" || burn.RecipientAddress == req.Address {
				totalBurned = totalBurned.Add(math.NewIntFromUint64(burn.Amount))
			}
		}

		if pageResBurn.NextKey == nil {
			break
		}
		burnPagination.Key = pageResBurn.NextKey
		burnPagination.CountTotal = false
	}

	if req.ShowFees {
		zentpFees, err := k.ZentpFees.Get(goCtx)
		if err != nil {
			return nil, err
		}

		return &types.QueryStatsResponse{
			TotalMinted: totalMinted.Uint64(),
			TotalBurned: totalBurned.Uint64(),
			MintsCount:  mintsCount,
			BurnsCount:  burnsCount,
			ZentpFees:   &zentpFees,
		}, nil
	}

	return &types.QueryStatsResponse{
		TotalMinted: totalMinted.Uint64(),
		TotalBurned: totalBurned.Uint64(),
		MintsCount:  mintsCount,
		BurnsCount:  burnsCount,
	}, nil
}
