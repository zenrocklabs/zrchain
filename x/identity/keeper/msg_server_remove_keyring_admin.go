package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RemoveKeyringAdmin(goCtx context.Context, msg *types.MsgRemoveKeyringAdmin) (*types.MsgRemoveKeyringAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	kr, err := k.KeyringStore.Get(ctx, msg.KeyringAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "keyring %s not found", msg.KeyringAddr)
	}

	if !kr.IsActive {
		return nil, errorsmod.Wrapf(types.ErrInactive, "keyring %s is not active", msg.KeyringAddr)
	}

	if !kr.IsAdmin(msg.Admin) {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "admin %s is not an admin of the keyring", msg.Admin)
	}

	if !kr.IsAdmin(msg.Creator) {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "tx creator %s is not keyring admin", msg.Creator)
	}

	if len(kr.Admins) == 1 {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "tx creator %s is the last admin in keyring", msg.Creator)
	}

	kr.RemoveAdmin(msg.Admin)

	if err := k.KeyringStore.Set(ctx, kr.Address, kr); err != nil {
		return nil, errorsmod.Wrapf(types.ErrInternal, "failed to set keyring %s", kr.Address)
	}

	return &types.MsgRemoveKeyringAdminResponse{}, nil
}
