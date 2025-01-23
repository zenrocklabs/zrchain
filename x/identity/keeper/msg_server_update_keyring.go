package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateKeyring(goCtx context.Context, msg *types.MsgUpdateKeyring) (*types.MsgUpdateKeyringResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	kr, err := k.KeyringStore.Get(ctx, msg.KeyringAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "keyring %s not found", msg.KeyringAddr)
	}

	if !kr.IsAdmin(msg.Creator) {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "keyring updates must be performed by an admin - %s is not a keyring admin", msg.Creator)
	}

	kr.SetKeyReqFee(msg.KeyReqFee)
	kr.SetSigReqFee(msg.SigReqFee)
	kr.SetStatus(msg.IsActive)
	kr.SetDescription(msg.Description)
	kr.SetMpcMinimumBtl(msg.MpcMinimumBtl)
	kr.SetMpcDefaultBtl(msg.MpcDefaultBtl)

	if msg.PartyThreshold > 0 && msg.PartyThreshold <= uint32(len(kr.Parties)) {
		kr.SetPartyThreshold(msg.PartyThreshold)
	}

	if err := k.KeyringStore.Set(ctx, kr.Address, kr); err != nil {
		return nil, errorsmod.Wrapf(types.ErrInternal, "failed to set keyring %s", kr.Address)
	}

	return &types.MsgUpdateKeyringResponse{}, nil
}
