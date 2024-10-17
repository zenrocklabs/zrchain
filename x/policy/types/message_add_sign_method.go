package types

import (
	errorsmod "cosmossdk.io/errors"
	codec "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAddSignMethod{}

func NewMsgAddSignMethod(creator string, config *codec.Any) *MsgAddSignMethod {
	return &MsgAddSignMethod{
		Creator: creator,
		Config:  config,
	}
}

func (msg *MsgAddSignMethod) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
