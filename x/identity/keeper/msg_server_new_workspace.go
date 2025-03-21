package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) NewWorkspace(goCtx context.Context, msg *types.MsgNewWorkspace) (*types.MsgNewWorkspaceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	workspace := &types.Workspace{
		Creator:         msg.Creator,
		ChildWorkspaces: nil,
		AdminPolicyId:   msg.AdminPolicyId,
		SignPolicyId:    msg.SignPolicyId,
	}

	if err := workspace.AddOwner(msg.Creator); err != nil {
		return nil, err
	}

	for _, owner := range msg.AdditionalOwners {
		if err := workspace.AddOwner(owner); err != nil {
			return nil, err
		}
	}

	if err := k.policyKeeper.PolicyMembersAreOwners(ctx, msg.AdminPolicyId, workspace.Owners); err != nil {
		return nil, err
	}
	if err := k.policyKeeper.PolicyMembersAreOwners(ctx, msg.SignPolicyId, workspace.Owners); err != nil {
		return nil, err
	}

	addr, err := k.CreateWorkspace(ctx, workspace)
	if err != nil {
		return nil, err
	}

	res := &types.MsgNewWorkspaceResponse{
		Addr: addr,
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventNewWorkspace,
			sdk.NewAttribute(types.AttributeWorkspaceAddr, res.GetAddr()),
		),
	})

	return res, nil
}
