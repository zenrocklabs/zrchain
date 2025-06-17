package keeper

import (
	"context"
	"fmt"

	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
)

func (k msgServer) TriggerEventBackfill(ctx context.Context, msg *types.MsgTriggerEventBackfill) (*types.MsgTriggerEventBackfillResponse, error) {
	if k.authority != msg.Authority {
		return nil, fmt.Errorf("invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	return &types.MsgTriggerEventBackfillResponse{}, nil
}
