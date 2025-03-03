package keeper

import (
	"context"
	"fmt"

	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	pol "github.com/Zenrock-Foundation/zrchain/v5/policy"
	policykeeper "github.com/Zenrock-Foundation/zrchain/v5/x/policy/keeper"
	policytypes "github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
)

func (k msgServer) UpdateKeyPolicy(goCtx context.Context, msg *types.MsgUpdateKeyPolicy) (*types.MsgUpdateKeyPolicyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	key, err := k.KeyStore.Get(ctx, msg.KeyId)
	if err != nil {
		return nil, fmt.Errorf("key %v not found", msg.KeyId)
	}

	ws, err := k.identityKeeper.GetWorkspace(ctx, key.WorkspaceAddr)
	if err != nil {
		return nil, fmt.Errorf("workspace %s not found", key.WorkspaceAddr)
	}

	if msg.SignPolicyId > 0 {
		if err := k.policyKeeper.PolicyMembersAreOwners(goCtx, msg.SignPolicyId, ws.Owners); err != nil {
			return nil, err
		}
	}

	act, err := k.policyKeeper.AddAction(ctx, msg.Creator, msg, key.SignPolicyId, 0, nil)
	if err != nil {
		return nil, err
	}

	return k.UpdateKeyPolicyActionHandler(ctx, act)
}

func (k msgServer) UpdateKeyPolicyPolicyGenerator(ctx sdk.Context, msg *types.MsgUpdateKeyPolicy) (pol.Policy, error) {
	key, err := k.KeyStore.Get(ctx, msg.KeyId)
	if err != nil {
		return nil, fmt.Errorf("key %v not found", msg.KeyId)
	}

	ws, err := k.identityKeeper.GetWorkspace(ctx, key.WorkspaceAddr)
	if err != nil {
		return nil, fmt.Errorf("workspace %s not found", key.WorkspaceAddr)
	}

	return ws.PolicyUpdateKeyPolicy(), nil
}

func (k msgServer) UpdateKeyPolicyActionHandler(ctx sdk.Context, act *policytypes.Action) (*types.MsgUpdateKeyPolicyResponse, error) {
	return policykeeper.TryExecuteAction(
		k.policyKeeper,
		k.cdc,
		ctx,
		act,
		k.updateKey,
	)
}

func (k msgServer) updateKey(ctx sdk.Context, msg *types.MsgUpdateKeyPolicy) (*types.MsgUpdateKeyPolicyResponse, error) {
	key, err := k.KeyStore.Get(ctx, msg.KeyId)
	if err != nil {
		return nil, fmt.Errorf("key %v not found", msg.KeyId)
	}

	key.SignPolicyId = msg.SignPolicyId
	k.KeyStore.Set(ctx, key.Id, key)

	return &types.MsgUpdateKeyPolicyResponse{}, nil
}
