package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgInitDct{}

func NewMsgInitDct(creator string, amount sdk.Coin, destinationChain string) *MsgInitDct {
	return &MsgInitDct{
		Creator:          creator,
		Amount:           amount,
		DestinationChain: destinationChain,
	}
}

func (msg *MsgInitDct) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Amount.Denom == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "amount is required")
	}
	if msg.DestinationChain == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "destination chain is required")
	}
	return nil
}
