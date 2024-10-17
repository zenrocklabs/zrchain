package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v4/policy"
	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) NewPolicy(goCtx context.Context, msg *types.MsgNewPolicy) (*types.MsgNewPolicyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params, err := k.Keeper.ParamStore.Get(ctx)
	if err != nil {
		return nil, err
	}

	if msg.Btl == 0 {
		msg.Btl = params.DefaultBtl
	}
	if msg.Btl < params.MinimumBtl {
		msg.Btl = params.MinimumBtl
	}

	var p policy.Policy
	if err := k.cdc.UnpackAny(msg.Policy, &p); err != nil {
		return nil, err
	}
	if err := p.Validate(); err != nil {
		return nil, err
	}

	policyPb := &types.Policy{
		Creator: msg.Creator,
		Name:    msg.Name,
		Policy:  msg.Policy,
		Btl:     msg.Btl,
	}

	id, err := k.CreatePolicy(ctx, policyPb)
	if err != nil {
		return nil, err
	}

	return &types.MsgNewPolicyResponse{
		Id: id,
	}, nil
}
