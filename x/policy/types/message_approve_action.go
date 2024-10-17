package types

import (
	errorsmod "cosmossdk.io/errors"
	codec "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgApproveAction{}

func NewMsgApproveAction(creator string, actionType string, actionId uint64, additionalSignatures []*codec.Any) *MsgApproveAction {
	return &MsgApproveAction{
		Creator:              creator,
		ActionType:           actionType,
		ActionId:             actionId,
		AdditionalSignatures: additionalSignatures,
	}
}

func (msg *MsgApproveAction) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if len(msg.ActionType) > 512 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid actiontype, max len 512")
	}
	return nil
}
