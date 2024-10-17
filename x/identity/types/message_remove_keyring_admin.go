package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgRemoveKeyringAdmin{}

func NewMsgRemoveKeyringAdmin(creator string, keyringAddr string, admin string) *MsgRemoveKeyringAdmin {
	return &MsgRemoveKeyringAdmin{
		Creator:     creator,
		KeyringAddr: keyringAddr,
		Admin:       admin,
	}
}

func (msg *MsgRemoveKeyringAdmin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid admin address (%s)", err)
	}

	krbz, err := sdk.GetFromBech32(msg.KeyringAddr, PrefixKeyringAddress)
	if err != nil || len(krbz) != KeyringAddressLength {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid keyring address (%s)", err)
	}

	return nil
}
