package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	pol "github.com/Zenrock-Foundation/zrchain/v4/policy"
	policykeeper "github.com/Zenrock-Foundation/zrchain/v4/x/policy/keeper"
	policytypes "github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"

	"github.com/Zenrock-Foundation/zrchain/v4/x/identity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RemoveWorkspaceOwner(goCtx context.Context, msg *types.MsgRemoveWorkspaceOwner) (*types.MsgRemoveWorkspaceOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ws, err := k.WorkspaceStore.Get(ctx, msg.WorkspaceAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "workspace %s not found", msg.WorkspaceAddr)
	}

	if !ws.IsOwner(msg.Creator) {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "creator %s is not a creator of the workspace", msg.Creator)
	}

	if !ws.IsOwner(msg.Owner) {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "owner %s is not an owner of the workspace %s", msg.Owner, msg.WorkspaceAddr)
	}

	if ws.SignPolicyId > 0 {
		participants, err := k.policyKeeper.GetPolicyParticipants(ctx, ws.SignPolicyId)
		if err != nil {
			return nil, err
		}
		if _, ok := participants[msg.Owner]; ok {
			return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "user %s is used in signpolicy of workspace %s", msg.Owner, msg.WorkspaceAddr)
		}
	}
	if ws.AdminPolicyId > 0 {
		participants, err := k.policyKeeper.GetPolicyParticipants(ctx, ws.AdminPolicyId)
		if err != nil {
			return nil, err
		}
		if _, ok := participants[msg.Owner]; ok {
			return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "user %s is used in signpolicy of workspace %s", msg.Owner, msg.WorkspaceAddr)
		}
	}

	act, err := k.policyKeeper.AddAction(ctx, msg.Creator, msg, ws.AdminPolicyId, msg.Btl, nil)
	if err != nil {
		return nil, err
	}
	return k.RemoveOwnerActionHandler(ctx, act)
}

func (k msgServer) RemoveOwnerPolicyGenerator(ctx sdk.Context, msg *types.MsgRemoveWorkspaceOwner) (pol.Policy, error) {
	ws, err := k.WorkspaceStore.Get(ctx, msg.WorkspaceAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "workspace %s not found", msg.WorkspaceAddr)
	}

	pol := ws.PolicyRemoveOwner()
	return pol, nil
}

func (k msgServer) RemoveOwnerActionHandler(ctx sdk.Context, act *policytypes.Action) (*types.MsgRemoveWorkspaceOwnerResponse, error) {
	return policykeeper.TryExecuteAction(
		&k.policyKeeper,
		k.cdc,
		ctx,
		act,
		func(ctx sdk.Context, msg *types.MsgRemoveWorkspaceOwner) (*types.MsgRemoveWorkspaceOwnerResponse, error) {
			ws, err := k.WorkspaceStore.Get(ctx, msg.WorkspaceAddr)
			if err != nil {
				return nil, errorsmod.Wrapf(types.ErrNotFound, "workspace %s not found", msg.WorkspaceAddr)
			}

			if !ws.IsOwner(msg.Owner) {
				return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "creator %s is not an owner of the child workspace", msg.Owner)
			}
			if len(ws.Owners) == 1 {
				return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "cannot remove last owner %s", msg.Owner)
			}
			ws.RemoveOwner(msg.Owner)

			if err := k.WorkspaceStore.Set(ctx, ws.Address, ws); err != nil {
				return nil, errorsmod.Wrapf(types.ErrInternal, "failed to set workspace %s", msg.WorkspaceAddr)
			}

			return &types.MsgRemoveWorkspaceOwnerResponse{}, nil
		},
	)
}
