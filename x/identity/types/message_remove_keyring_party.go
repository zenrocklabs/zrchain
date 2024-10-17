package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgRemoveKeyringParty{}

func NewMsgRemoveKeyringParty(creator string, keyringAddr string, party string) *MsgRemoveKeyringParty {
	return &MsgRemoveKeyringParty{
		Creator:     creator,
		KeyringAddr: keyringAddr,
		Party:       party,
	}
}

func (msg *MsgRemoveKeyringParty) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Party)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid party address (%s)", err)
	}

	krbz, err := sdk.GetFromBech32(msg.KeyringAddr, PrefixKeyringAddress)
	if err != nil || len(krbz) != KeyringAddressLength {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid keyring address (%s)", err)
	}

	return nil
}
