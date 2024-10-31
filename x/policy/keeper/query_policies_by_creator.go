package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PoliciesByCreator(goCtx context.Context, req *types.QueryPoliciesByCreatorRequest) (*types.QueryPoliciesByCreatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	policies, pageRes, err := query.CollectionFilteredPaginate[uint64, types.Policy, collections.Map[uint64, types.Policy], *types.Policy](
		ctx,
		k.PolicyStore,
		req.Pagination,
		func(key uint64, value types.Policy) (include bool, err error) {
			isCreator := false
			for _, c := range req.Creators {
				if value.Creator == c {
					isCreator = true
					break
				}
			}
			return isCreator, nil
		},
		func(key uint64, value types.Policy) (*types.Policy, error) {
			return &types.Policy{
				Creator: value.Creator,
				Id:      value.Id,
				Name:    value.Name,
				Policy:  value.Policy,
			}, nil
		},
	)

	if err != nil {
		return nil, err
	}

	return &types.QueryPoliciesByCreatorResponse{
		Policies:   policies,
		Pagination: pageRes,
	}, nil
}
