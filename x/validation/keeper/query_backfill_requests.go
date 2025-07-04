package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
)

func (k Querier) QueryBackfillRequests(ctx context.Context, req *types.QueryBackfillRequestsRequest) (*types.QueryBackfillRequestsResponse, error) {
	backfillRequests, err := k.BackfillRequests.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// Return empty BackfillRequests when the collection is empty
			emptyRequests := types.BackfillRequests{}
			return &types.QueryBackfillRequestsResponse{
				BackfillRequests: &emptyRequests,
			}, nil
		}
		return nil, err
	}

	return &types.QueryBackfillRequestsResponse{
		BackfillRequests: &backfillRequests,
	}, nil
}
