package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueryBurnEvents(ctx context.Context, req *types.QueryBurnEventsRequest) (*types.QueryBurnEventsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var matchingBurnEvents []*types.BurnEvent

	// Walk through all assets if asset is unspecified, otherwise filter by asset
	asset := req.Asset
	if asset == types.Asset_ASSET_UNSPECIFIED {
		// Query all burn events across all assets
		if err := k.BurnEvents.Walk(ctx, nil, func(key collections.Pair[string, uint64], burnEvent types.BurnEvent) (bool, error) {
			if key.K2() < req.StartIndex {
				return false, nil
			}
			if (req.TxID == "" || burnEvent.TxID == req.TxID) &&
				(req.LogIndex == 0 || burnEvent.LogIndex == req.LogIndex) &&
				(req.Caip2ChainID == "" || burnEvent.ChainID == req.Caip2ChainID) &&
				(req.Status == types.BurnStatus_BURN_STATUS_UNSPECIFIED || burnEvent.Status == req.Status) {

				matchingBurnEvents = append(matchingBurnEvents, &burnEvent)
			}
			return false, nil
		}); err != nil {
			return nil, err
		}
	} else {
		// Query burn events for a specific asset
		if err := k.WalkBurnEvents(ctx, asset, func(id uint64, burnEvent types.BurnEvent) (bool, error) {
			if id < req.StartIndex {
				return false, nil
			}
			if (req.TxID == "" || burnEvent.TxID == req.TxID) &&
				(req.LogIndex == 0 || burnEvent.LogIndex == req.LogIndex) &&
				(req.Caip2ChainID == "" || burnEvent.ChainID == req.Caip2ChainID) &&
				(req.Status == types.BurnStatus_BURN_STATUS_UNSPECIFIED || burnEvent.Status == req.Status) {

				matchingBurnEvents = append(matchingBurnEvents, &burnEvent)
			}
			return false, nil
		}); err != nil {
			return nil, err
		}
	}

	return &types.QueryBurnEventsResponse{BurnEvents: matchingBurnEvents}, nil
}
