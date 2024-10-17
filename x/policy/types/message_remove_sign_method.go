package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgRemoveSignMethod{}

func NewMsgRemoveSignMethod(creator string, id string) *MsgRemoveSignMethod {
	return &MsgRemoveSignMethod{
		Creator: creator,
		Id:      id,
	}
}

func (msg *MsgRemoveSignMethod) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if len(msg.Id) > 512 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid id, max len 512")
	}

	return nil
}
