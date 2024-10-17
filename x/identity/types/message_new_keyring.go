package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgNewKeyring{}

func NewMsgNewKeyring(creator string, description string, keyReqFee, sigReqFee uint64) *MsgNewKeyring {
	return &MsgNewKeyring{
		Creator:     creator,
		Description: description,
		KeyReqFee:   keyReqFee,
		SigReqFee:   sigReqFee,
	}
}

func (msg *MsgNewKeyring) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Description) > 255 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid description, max len is 255")
	}
	return nil
}
