package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"

	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"

	"github.com/cosmos/cosmos-sdk/types/query"
)

func (k Keeper) Keyrings(goCtx context.Context, req *types.QueryKeyringsRequest) (*types.QueryKeyringsResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "invalid arguments: request is nil")
	}

	keyrings, pageRes, err := query.CollectionFilteredPaginate[string, types.Keyring, collections.Map[string, types.Keyring], *types.Keyring](
		goCtx,
		k.KeyringStore,
		req.Pagination,
		nil,
		func(key string, value types.Keyring) (*types.Keyring, error) {
			return &value, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.QueryKeyringsResponse{
		Keyrings:   keyrings,
		Pagination: pageRes,
	}, nil
}
