package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) KeyringByAddress(goCtx context.Context, req *types.QueryKeyringByAddressRequest) (*types.QueryKeyringByAddressResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "invalid arguments: request is nil")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	keyring, err := k.KeyringStore.Get(ctx, req.KeyringAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "keyring %s not found", req.KeyringAddr)
	}

	return &types.QueryKeyringByAddressResponse{
		Keyring: &keyring,
	}, nil
}
