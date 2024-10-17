package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
)

func (k Keeper) Policies(goCtx context.Context, req *types.QueryPoliciesRequest) (*types.QueryPoliciesResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "request is nil")
	}

	policies, pageRes, err := query.CollectionFilteredPaginate[uint64, types.Policy, collections.Map[uint64, types.Policy], types.PolicyResponse](
		goCtx,
		k.PolicyStore,
		req.Pagination,
		func(key uint64, value types.Policy) (bool, error) {
			return true, nil
		},
		func(key uint64, value types.Policy) (types.PolicyResponse, error) {
			// Transform each Policy to a PolicyResponse
			response, err := types.NewPolicyResponse(k.cdc, &value)
			if err != nil {
				return types.PolicyResponse{}, err
			}
			return *response, nil
		},
	)

	if err != nil {
		return nil, err
	}

	return &types.QueryPoliciesResponse{
		Policies:   policies,
		Pagination: pageRes,
	}, nil
}
