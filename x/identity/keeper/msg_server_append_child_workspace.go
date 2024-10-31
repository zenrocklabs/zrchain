package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	pol "github.com/Zenrock-Foundation/zrchain/v5/policy"
	policykeeper "github.com/Zenrock-Foundation/zrchain/v5/x/policy/keeper"
	policytypes "github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"

	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AppendChildWorkspace(goCtx context.Context, msg *types.MsgAppendChildWorkspace) (*types.MsgAppendChildWorkspaceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	child, err := k.WorkspaceStore.Get(ctx, msg.ChildWorkspaceAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "child workspace %s not found", msg.ChildWorkspaceAddr)
	}
	parent, err := k.WorkspaceStore.Get(ctx, msg.ParentWorkspaceAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "parent workspace %s not found", msg.ParentWorkspaceAddr)
	}

	if !parent.IsOwner(msg.Creator) {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "creator %s is not an owner of the parent workspace %s", msg.Creator, msg.ParentWorkspaceAddr)
	}

	if !child.IsOwner(msg.Creator) {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "creator %s is not an owner of the child workspace %s", msg.Creator, msg.ChildWorkspaceAddr)
	}

	if parent.IsChild(msg.ChildWorkspaceAddr) {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "new child is already a child workspace %s", msg.ChildWorkspaceAddr)
	}

	act, err := k.policyKeeper.AddAction(ctx, msg.Creator, msg, parent.AdminPolicyId, msg.Btl, nil)
	if err != nil {
		return nil, err
	}
	return k.AppendChildWorkspaceActionHandler(ctx, act)
}

func (k msgServer) AppendChildWorkspacePolicyGenerator(ctx sdk.Context, msg *types.MsgAppendChildWorkspace) (pol.Policy, error) {
	parent, err := k.WorkspaceStore.Get(ctx, msg.ParentWorkspaceAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "parent workspace %s not found", msg.ParentWorkspaceAddr)
	}

	pol := parent.PolicyAppendChild()
	return pol, nil
}

func (k msgServer) AppendChildWorkspaceActionHandler(ctx sdk.Context, act *policytypes.Action) (*types.MsgAppendChildWorkspaceResponse, error) {
	return policykeeper.TryExecuteAction(
		&k.policyKeeper,
		k.cdc,
		ctx,
		act,
		func(ctx sdk.Context, msg *types.MsgAppendChildWorkspace) (*types.MsgAppendChildWorkspaceResponse, error) {
			child, err := k.WorkspaceStore.Get(ctx, msg.ChildWorkspaceAddr)
			if err != nil {
				return nil, errorsmod.Wrapf(types.ErrNotFound, "child workspace %s not found", msg.ChildWorkspaceAddr)
			}
			parent, err := k.WorkspaceStore.Get(ctx, msg.ParentWorkspaceAddr)
			if err != nil {
				return nil, errorsmod.Wrapf(types.ErrNotFound, "parent workspace %s not found", msg.ParentWorkspaceAddr)
			}

			parent.AddChild(&child)

			if err = k.WorkspaceStore.Set(ctx, parent.Address, parent); err != nil {
				return nil, errorsmod.Wrapf(types.ErrInternal, "failed to set parent workspace %s", msg.ParentWorkspaceAddr)
			}

			return &types.MsgAppendChildWorkspaceResponse{}, nil
		},
	)
}
