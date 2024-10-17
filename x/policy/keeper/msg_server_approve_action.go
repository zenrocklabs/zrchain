package keeper

import (
	"context"
	"fmt"

	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ApproveAction(goCtx context.Context, msg *types.MsgApproveAction) (*types.MsgApproveActionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	act, err := k.ActionStore.Get(ctx, msg.ActionId)
	if err != nil {
		return nil, fmt.Errorf("action not found")
	}

	if act.Status != types.ActionStatus_ACTION_STATUS_PENDING {
		return nil, fmt.Errorf("action not pending %s", act.Status.String())
	}

	if act.Btl > 0 && act.Btl < uint64(ctx.BlockHeight()) {
		act.Status = types.ActionStatus_ACTION_STATUS_TIMEOUT
		if err := k.ActionStore.Set(ctx, act.Id, act); err != nil {
			return nil, err
		}
		return &types.MsgApproveActionResponse{
			Status: act.Status.String(),
		}, nil
	}

	policy, err := PolicyForAction(ctx, &k.Keeper, &act)
	if err != nil {
		return nil, err
	}

	participant, err := policy.AddressToParticipant(msg.Creator)
	if err != nil {
		return nil, err
	}

	if err := act.AddApprover(participant); err != nil {
		return nil, err
	}

	// process additional sigs here and add approver if valid
	for _, anySig := range msg.GetAdditionalSignatures() {
		var sig types.AdditionalSignature
		if err := k.cdc.UnpackAny(anySig, &sig); err != nil {
			return nil, err
		}

		signMethod, err := k.GetSignMethod(ctx, msg.Creator, sig.GetConfigId())
		if err != nil {
			return nil, err
		}

		if sig != nil && signMethod != nil && signMethod.IsActive() {
			addr := sig.Verify(ctx, signMethod, act)
			if addr != "" {
				participant, err := policy.AddressToParticipant(addr)
				if err != nil {
					continue // skip if address is not a member of the policy
				}
				if err := act.AddApprover(participant); err != nil {
					return nil, err
				}
			}
		}
	}

	if err := k.ActionStore.Set(ctx, act.Id, act); err != nil {
		return nil, err
	}

	h, ok := k.actionHandlers[msg.ActionType]
	if !ok {
		return nil, fmt.Errorf("action handler not found for %s", msg.ActionType)
	}

	if _, err := h(ctx, &act); err != nil {
		return nil, err
	}

	return &types.MsgApproveActionResponse{Status: act.Status.String()}, nil
}
