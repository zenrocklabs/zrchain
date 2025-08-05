package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgInitDct{}

func NewMsgInitDct(creator string, asset sdk.Coin, destinationChain string) *MsgInitDct {
	return &MsgInitDct{
		Creator:          creator,
		Asset:            &asset,
		DestinationChain: destinationChain,
	}
}

func (msg *MsgInitDct) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Asset == nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "asset is required")
	}
	if msg.DestinationChain == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "destination chain is required")
	}
	return nil
}
