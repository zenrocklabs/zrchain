package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
)

func (k Querier) QueryBackfillRequests(ctx context.Context, req *types.QueryBackfillRequestsRequest) (*types.QueryBackfillRequestsResponse, error) {
	backfillRequests, err := k.BackfillRequests.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryBackfillRequestsResponse{
		BackfillRequests: &backfillRequests,
	}, nil
}
