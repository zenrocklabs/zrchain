package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeactivateKeyring(goCtx context.Context, msg *types.MsgDeactivateKeyring) (*types.MsgDeactivateKeyringResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	kr, err := k.KeyringStore.Get(ctx, msg.KeyringAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "keyring %s not found", msg.KeyringAddr)
	}

	if !kr.IsAdmin(msg.Creator) {
		return nil, fmt.Errorf("keyring deactivations should be requested by admins")
	}

	kr.SetStatus(false)

	if err = k.KeyringStore.Set(ctx, kr.Address, kr); err != nil {
		return nil, errorsmod.Wrapf(types.ErrInternal, "failed to set keyring %s", kr.Address)
	}

	return &types.MsgDeactivateKeyringResponse{}, nil
}
