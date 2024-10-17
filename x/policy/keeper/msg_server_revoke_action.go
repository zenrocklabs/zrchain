package keeper

import (
	"context"
	"fmt"

	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RevokeAction(goCtx context.Context, msg *types.MsgRevokeAction) (*types.MsgRevokeActionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	act, err := k.ActionStore.Get(ctx, msg.ActionId)
	if err != nil {
		return nil, fmt.Errorf("action not found")
	}

	if act.Creator != msg.Creator {
		return nil, fmt.Errorf("action creator does not match")
	}

	if act.Status != types.ActionStatus_ACTION_STATUS_PENDING {
		return nil, fmt.Errorf("action status is not pending")
	}

	act.Status = types.ActionStatus_ACTION_STATUS_REVOKED

	if err := k.ActionStore.Set(ctx, act.Id, act); err != nil {
		return nil, err
	}

	return &types.MsgRevokeActionResponse{}, nil
}
