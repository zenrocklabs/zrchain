package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RemoveKeyringParty(goCtx context.Context, msg *types.MsgRemoveKeyringParty) (*types.MsgRemoveKeyringPartyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	kr, err := k.KeyringStore.Get(ctx, msg.KeyringAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "keyring %s not found", msg.KeyringAddr)
	}

	if !kr.IsActive {
		return nil, errorsmod.Wrapf(types.ErrInactive, "keyring %s is not active", msg.KeyringAddr)
	}

	if !kr.IsParty(msg.Party) {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "party %s is not a party of the keyring %s", msg.Party, msg.KeyringAddr)
	}

	if !kr.IsAdmin(msg.Creator) {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "tx creator %s is not keyring admin", msg.Creator)
	}

	kr.RemoveParty(msg.Party)

	if msg.DecreaseThreshold {
		kr.PartyThreshold--
	}

	if len(kr.Parties) < int(kr.PartyThreshold) {
		kr.PartyThreshold = uint32(len(kr.Parties))
	}

	if err := k.KeyringStore.Set(ctx, kr.Address, kr); err != nil {
		return nil, errorsmod.Wrapf(types.ErrInternal, "failed to set keyring %s", kr.Address)
	}

	return &types.MsgRemoveKeyringPartyResponse{}, nil
}
