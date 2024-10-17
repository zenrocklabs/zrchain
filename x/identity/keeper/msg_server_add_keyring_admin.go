package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/Zenrock-Foundation/zrchain/v4/x/identity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddKeyringAdmin(goCtx context.Context, msg *types.MsgAddKeyringAdmin) (*types.MsgAddKeyringAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	kr, err := k.KeyringStore.Get(ctx, msg.KeyringAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "keyring %s not found", msg.KeyringAddr)
	}

	if !kr.IsActive {
		return nil, errorsmod.Wrapf(types.ErrInactive, "keyring %s is not active", msg.KeyringAddr)
	}

	if kr.IsAdmin(msg.Admin) {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "admin %s is already an admin", msg.Creator)
	}

	if !kr.IsAdmin(msg.Creator) {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "tx creator %s is not keyring admin", msg.Creator)
	}

	kr.AddAdmin(msg.Admin)

	if err = k.KeyringStore.Set(ctx, kr.Address, kr); err != nil {
		return nil, errorsmod.Wrapf(types.ErrInternal, "failed to set keyring %s", kr.Address)
	}

	return &types.MsgAddKeyringAdminResponse{}, nil
}
