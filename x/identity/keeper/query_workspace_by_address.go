package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) WorkspaceByAddress(goCtx context.Context, req *types.QueryWorkspaceByAddressRequest) (*types.QueryWorkspaceByAddressResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "invalid arguments: request is nil")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	workspace, err := k.WorkspaceStore.Get(ctx, req.WorkspaceAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "workspace %s is nil or not found", req.WorkspaceAddr)
	}

	return &types.QueryWorkspaceByAddressResponse{
		Workspace: &workspace,
	}, nil
}
