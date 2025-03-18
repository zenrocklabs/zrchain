package keeper

import (
	"context"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	pol "github.com/Zenrock-Foundation/zrchain/v5/policy"
	policykeeper "github.com/Zenrock-Foundation/zrchain/v5/x/policy/keeper"
	policytypes "github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"

	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddWorkspaceOwner(goCtx context.Context, msg *types.MsgAddWorkspaceOwner) (*types.MsgAddWorkspaceOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ws, err := k.WorkspaceStore.Get(ctx, msg.WorkspaceAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "workspace %s not found", msg.WorkspaceAddr)
	}

	if !ws.IsOwner(msg.Creator) {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "tx creator %s is not owner of the workspace", msg.Creator)
	}

	if ws.IsOwner(msg.NewOwner) {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "new owner %s is already an owner of the workspace", msg.Creator)
	}

	act, err := k.policyKeeper.AddAction(ctx, msg.Creator, msg, ws.AdminPolicyId, msg.Btl, nil, ws.Owners)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventAddOwnerToWorkspace,
			sdk.NewAttribute(types.AttributeActionId, strconv.FormatUint(act.GetId(), 10)),
		),
	})

	return k.AddOwnerActionHandler(ctx, act)
}

func (k msgServer) AddOwnerPolicyGenerator(ctx sdk.Context, msg *types.MsgAddWorkspaceOwner) (pol.Policy, error) {
	ws, err := k.WorkspaceStore.Get(ctx, msg.WorkspaceAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "workspace %s not found", msg.WorkspaceAddr)
	}

	pol := ws.PolicyAddOwner()
	return pol, nil
}

func (k msgServer) AddOwnerActionHandler(ctx sdk.Context, act *policytypes.Action) (*types.MsgAddWorkspaceOwnerResponse, error) {
	return policykeeper.TryExecuteAction(
		k.policyKeeper,
		k.cdc,
		ctx,
		act,
		func(ctx sdk.Context, msg *types.MsgAddWorkspaceOwner) (*types.MsgAddWorkspaceOwnerResponse, error) {
			ws, err := k.WorkspaceStore.Get(ctx, msg.WorkspaceAddr)
			if err != nil {
				return nil, errorsmod.Wrapf(types.ErrNotFound, "workspace %s not found", msg.WorkspaceAddr)
			}

			if err := ws.AddOwner(msg.NewOwner); err != nil {
				return nil, err
			}

			if err = k.WorkspaceStore.Set(ctx, ws.Address, ws); err != nil {
				return nil, errorsmod.Wrapf(types.ErrInternal, "failed to set workspace %s", msg.WorkspaceAddr)
			}

			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventOwnerAddedToWorkspace,
					sdk.NewAttribute(types.AttributeWorkspaceAddr, msg.WorkspaceAddr),
					sdk.NewAttribute(types.AttributeOwnerAddr, msg.NewOwner),
				),
			})

			return &types.MsgAddWorkspaceOwnerResponse{}, nil
		},
	)
}
