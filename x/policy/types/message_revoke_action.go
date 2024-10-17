package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgRevokeAction{}

func NewMsgRevokeAction(creator string, actionType string, actionId uint64) *MsgRevokeAction {
	return &MsgRevokeAction{
		Creator:  creator,
		ActionId: actionId,
	}
}

func (msg *MsgRevokeAction) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
