package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) QuerySupply(ctx context.Context, req *types.QuerySupplyRequest) (*types.QuerySupplyResponse, error) {
	asset := req.Asset
	var supplies []*types.AssetSupply

	if asset == types.Asset_ASSET_UNSPECIFIED {
		// Return all asset supplies
		_ = k.Supply.Walk(sdk.UnwrapSDKContext(ctx), nil, func(assetKey string, supply types.Supply) (bool, error) {
			exchangeRate, err := k.GetExchangeRate(sdk.UnwrapSDKContext(ctx), supply.Asset)
			if err != nil {
				if !errors.Is(err, collections.ErrNotFound) {
					return true, err
				}
			}
			supplies = append(supplies, &types.AssetSupply{
				Supply:       supply,
				ExchangeRate: exchangeRate.String(),
			})
			return false, nil
		})
	} else {
		// Return supply for specific asset
		supply, err := k.GetSupply(sdk.UnwrapSDKContext(ctx), asset)
		if err != nil {
			if errors.Is(err, collections.ErrNotFound) {
				return &types.QuerySupplyResponse{Supplies: []*types.AssetSupply{}}, nil
			}
			return nil, err
		}

		exchangeRate, err := k.GetExchangeRate(sdk.UnwrapSDKContext(ctx), asset)
		if err != nil {
			if !errors.Is(err, collections.ErrNotFound) {
				return nil, err
			}
		}

		supplies = append(supplies, &types.AssetSupply{
			Supply:       supply,
			ExchangeRate: exchangeRate.String(),
		})
	}

	return &types.QuerySupplyResponse{Supplies: supplies}, nil
}
