package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	pol "github.com/Zenrock-Foundation/zrchain/v5/policy"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	policy "github.com/Zenrock-Foundation/zrchain/v5/x/policy/keeper"
	policytypes "github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateWorkspace(goCtx context.Context, msg *types.MsgUpdateWorkspace) (*types.MsgUpdateWorkspaceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ws, err := k.WorkspaceStore.Get(ctx, msg.WorkspaceAddr)
	if err != nil {
		return nil, fmt.Errorf("workspace %s not found", msg.WorkspaceAddr)
	}

	if !ws.IsOwner(msg.Creator) {
		return nil, fmt.Errorf("creator %s is not an owner of the workspace %s", msg.Creator, msg.WorkspaceAddr)
	}

	if !ws.IsDifferent(msg.AdminPolicyId, msg.SignPolicyId) {
		return nil, fmt.Errorf("no updates to the policies")
	}

	if err := k.policyKeeper.PolicyMembersAreOwners(ctx, msg.AdminPolicyId, ws.Owners); err != nil {
		return nil, err
	}
	if err := k.policyKeeper.PolicyMembersAreOwners(ctx, msg.SignPolicyId, ws.Owners); err != nil {
		return nil, err
	}

	act, err := k.policyKeeper.AddAction(ctx, msg.Creator, msg, ws.AdminPolicyId, msg.Btl, nil)
	if err != nil {
		return nil, err
	}

	return k.UpdateWorkspaceActionHandler(ctx, act)
}

func (k msgServer) UpdateWorkspacePolicyGenerator(ctx sdk.Context, msg *types.MsgUpdateWorkspace) (pol.Policy, error) {
	ws, err := k.WorkspaceStore.Get(ctx, msg.WorkspaceAddr)
	if err != nil {
		return nil, fmt.Errorf("workspace not found")
	}

	pol := ws.PolicyUpdateWorkspace()
	return pol, nil
}

func (k msgServer) UpdateWorkspaceActionHandler(ctx sdk.Context, act *policytypes.Action) (*types.MsgUpdateWorkspaceResponse, error) {
	return policy.TryExecuteAction(
		&k.policyKeeper,
		k.cdc,
		ctx,
		act,
		func(ctx sdk.Context, msg *types.MsgUpdateWorkspace) (*types.MsgUpdateWorkspaceResponse, error) {
			ws, err := k.WorkspaceStore.Get(ctx, msg.WorkspaceAddr)
			if err != nil {
				return nil, fmt.Errorf("workspace not found")
			}

			if msg.AdminPolicyId != ws.AdminPolicyId {
				if msg.AdminPolicyId != 0 {
					if _, err := k.policyKeeper.PolicyStore.Get(ctx, msg.AdminPolicyId); err != nil {
						return nil, fmt.Errorf("admin policy %v not found", msg.AdminPolicyId)
					}
				}
				ws.AdminPolicyId = msg.AdminPolicyId
			}

			if msg.SignPolicyId != ws.SignPolicyId {
				if msg.SignPolicyId != 0 {
					if _, err := k.policyKeeper.PolicyStore.Get(ctx, msg.SignPolicyId); err != nil {
						return nil, fmt.Errorf("sign policy not found")
					}
				}
				ws.SignPolicyId = msg.SignPolicyId
			}

			if err = k.WorkspaceStore.Set(ctx, ws.Address, ws); err != nil {
				return nil, errorsmod.Wrapf(types.ErrInternal, "failed to set workspace %s", msg.WorkspaceAddr)
			}

			return &types.MsgUpdateWorkspaceResponse{}, nil
		},
	)
}
