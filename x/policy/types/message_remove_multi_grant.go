package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgRemoveMultiGrant{}

func NewMsgRemoveMultiGrant(creator string, grantee string, msgs []string) *MsgAddMultiGrant {
	return &MsgAddMultiGrant{
		Creator: creator,
		Grantee: grantee,
		Msgs:    msgs,
	}
}

func (msg *MsgRemoveMultiGrant) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Grantee)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid grantee address (%s)", err)
	}

	if len(msg.Msgs) == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "at least 1 msg type is required")
	}

	for _, msg := range msg.Msgs {
		if len(msg) > 512 {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid msg, max len 512")
		}
	}
	return nil
}
