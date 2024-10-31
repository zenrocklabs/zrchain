package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SignMethodsByAddress(goCtx context.Context, req *types.QuerySignMethodsByAddressRequest) (*types.QuerySignMethodsByAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	signMethods, pageRes, err := query.CollectionFilteredPaginate[collections.Pair[string, string], codectypes.Any, collections.Map[collections.Pair[string, string], codectypes.Any], *codectypes.Any](
		goCtx,
		k.SignMethodStore,
		req.Pagination,
		func(key collections.Pair[string, string], value codectypes.Any) (bool, error) {
			return key.K1() == req.Address, nil
		},
		func(key collections.Pair[string, string], value codectypes.Any) (*codectypes.Any, error) {
			return &value, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.QuerySignMethodsByAddressResponse{
		Pagination: pageRes,
		Config:     signMethods,
	}, nil
}
