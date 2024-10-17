package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateKeyring{}

func NewMsgUpdateKeyring(creator string, keyringAddr string, description string, isActive bool) *MsgUpdateKeyring {
	return &MsgUpdateKeyring{
		Creator:     creator,
		KeyringAddr: keyringAddr,
		Description: description,
		IsActive:    isActive,
	}
}

func (msg *MsgUpdateKeyring) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	krbz, err := sdk.GetFromBech32(msg.KeyringAddr, PrefixKeyringAddress)
	if err != nil || len(krbz) != KeyringAddressLength {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid keyring address (%s)", err)
	}

	if len(msg.Description) > 255 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid description, max len is 255")
	}

	return nil
}
