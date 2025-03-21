package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddKeyringParty(goCtx context.Context, msg *types.MsgAddKeyringParty) (*types.MsgAddKeyringPartyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	kr, err := k.KeyringStore.Get(ctx, msg.KeyringAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "keyring %s not found", msg.KeyringAddr)
	}

	if !kr.IsActive {
		return nil, errorsmod.Wrapf(types.ErrInactive, "keyring %s is not active", msg.KeyringAddr)
	}

	if kr.IsParty(msg.Party) {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "party %s is already a party of the keyring %s", msg.Party, msg.KeyringAddr)
	}

	if !kr.IsAdmin(msg.Creator) {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "tx creator %s is not keyring admin", msg.Creator)
	}

	kr.AddParty(msg.Party)

	if msg.IncreaseThreshold {
		kr.PartyThreshold++
	}

	// When adding a party make sure the threshold is at least 1
	if kr.PartyThreshold == 0 {
		kr.PartyThreshold = 1
	}

	if err = k.KeyringStore.Set(ctx, kr.Address, kr); err != nil {
		return nil, errorsmod.Wrapf(types.ErrInternal, "failed to set keyring %s", kr.Address)
	}

	return &types.MsgAddKeyringPartyResponse{}, nil
}
