package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgNewChildWorkspace{}

func NewMsgNewChildWorkspace(creator string, parentWorkspaceAddr string, btl uint64) *MsgNewChildWorkspace {
	return &MsgNewChildWorkspace{
		Creator:             creator,
		ParentWorkspaceAddr: parentWorkspaceAddr,
		Btl:                 btl,
	}
}

func (msg *MsgNewChildWorkspace) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	wsbz, err := sdk.GetFromBech32(msg.ParentWorkspaceAddr, PrefixWorkspaceAddress)
	if err != nil || len(wsbz) != WorkspaceAddressLength {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid parent workspace address (%s)", err)
	}

	return nil
}
