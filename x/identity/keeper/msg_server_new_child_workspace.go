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

func (k msgServer) NewChildWorkspace(goCtx context.Context, msg *types.MsgNewChildWorkspace) (*types.MsgNewChildWorkspaceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	parent, err := k.WorkspaceStore.Get(ctx, msg.ParentWorkspaceAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "parent workspace %s not found", msg.ParentWorkspaceAddr)
	}

	if !parent.IsOwner(msg.Creator) {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "creator %s is not an owner of the workspace %s", msg.Creator, msg.ParentWorkspaceAddr)
	}

	act, err := k.policyKeeper.AddAction(ctx, msg.Creator, msg, parent.AdminPolicyId, msg.Btl, nil)
	if err != nil {
		return nil, err
	}
	return k.NewChildWorkspaceActionHandler(ctx, act)
}

func (k msgServer) NewChildWorkspacePolicyGenerator(ctx sdk.Context, msg *types.MsgNewChildWorkspace) (pol.Policy, error) {
	parent, err := k.WorkspaceStore.Get(ctx, msg.ParentWorkspaceAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "parent workspace %s not found", msg.ParentWorkspaceAddr)
	}

	pol := parent.PolicyAppendChild()
	return pol, nil
}

func (k msgServer) NewChildWorkspaceActionHandler(ctx sdk.Context, act *policytypes.Action) (*types.MsgNewChildWorkspaceResponse, error) {
	return policykeeper.TryExecuteAction(
		k.policyKeeper,
		k.cdc,
		ctx,
		act,
		func(ctx sdk.Context, msg *types.MsgNewChildWorkspace) (*types.MsgNewChildWorkspaceResponse, error) {
			parent, err := k.WorkspaceStore.Get(ctx, msg.ParentWorkspaceAddr)
			if err != nil {
				return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "parent workspace %s not found", msg.ParentWorkspaceAddr)
			}

			child := &types.Workspace{
				Creator:       msg.Creator,
				Owners:        []string{msg.Creator},
				AdminPolicyId: parent.AdminPolicyId,
				SignPolicyId:  parent.SignPolicyId,
			}

			return k.storeChildWorkspace(ctx, &parent, child)
		},
	)
}
