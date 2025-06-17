package keeper

import (
	"context"
	"fmt"
	"slices"

	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
)

func (k msgServer) TriggerEventBackfill(ctx context.Context, msg *types.MsgTriggerEventBackfill) (*types.MsgTriggerEventBackfillResponse, error) {
	if k.authority != msg.Authority {
		return nil, fmt.Errorf("invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	backfillRequests, err := k.BackfillRequests.Get(ctx)
	if err != nil {
		return nil, err
	}

	if slices.Contains(backfillRequests.Requests, msg) {
		return nil, fmt.Errorf("backfill request already exists")
	}

	backfillRequests.Requests = append(backfillRequests.Requests, msg)

	if err = k.BackfillRequests.Set(ctx, backfillRequests); err != nil {
		return nil, err
	}

	return &types.MsgTriggerEventBackfillResponse{}, nil
}
