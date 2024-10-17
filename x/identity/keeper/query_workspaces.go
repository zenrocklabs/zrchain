package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"

	"github.com/Zenrock-Foundation/zrchain/v4/x/identity/types"

	"github.com/cosmos/cosmos-sdk/types/query"
)

func (k Keeper) Workspaces(goCtx context.Context, req *types.QueryWorkspacesRequest) (*types.QueryWorkspacesResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "invalid arguments: request is nil")
	}

	workspaces, pageRes, err := query.CollectionFilteredPaginate[string, types.Workspace, collections.Map[string, types.Workspace], *types.Workspace](
		goCtx,
		k.WorkspaceStore,
		req.Pagination,
		func(key string, value types.Workspace) (bool, error) {
			return (req.Owner != "" && value.IsOwner(req.Owner)) ||
				(req.Creator != "" && value.Creator == req.Creator) ||
				(req.Creator == "" && req.Owner == ""), nil
		},
		func(key string, value types.Workspace) (*types.Workspace, error) {
			return &value, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.QueryWorkspacesResponse{
		Workspaces: workspaces,
		Pagination: pageRes,
	}, nil
}
