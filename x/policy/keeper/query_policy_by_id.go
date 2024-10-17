package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
)

func (k Keeper) PolicyById(goCtx context.Context, req *types.QueryPolicyByIdRequest) (*types.QueryPolicyByIdResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidArgument, "request is nil")
	}

	policyPb, err := k.PolicyStore.Get(goCtx, req.Id)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "policy %d not found", req.Id)
	}

	res, err := types.NewPolicyResponse(k.cdc, &policyPb)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInternal, "new policy response error %s", err.Error())
	}

	return &types.QueryPolicyByIdResponse{Policy: res}, nil
}
