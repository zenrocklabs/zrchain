package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
)

func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "request is nil")
	}

	params, err := k.ParamStore.Get(goCtx)
	if err != nil {
		return nil, nil
	}

	return &types.QueryParamsResponse{Params: params}, nil
}
