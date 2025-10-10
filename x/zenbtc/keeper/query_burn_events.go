package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueryBurnEvents(ctx context.Context, req *types.QueryBurnEventsRequest) (*types.QueryBurnEventsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var matchingBurnEvents []*types.BurnEvent
	var queryRange collections.Range[uint64]

	if err := k.BurnEvents.Walk(ctx, queryRange.StartInclusive(req.StartIndex), func(_ uint64, burnEvent types.BurnEvent) (bool, error) {
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

	return &types.QueryBurnEventsResponse{BurnEvents: matchingBurnEvents}, nil
}
