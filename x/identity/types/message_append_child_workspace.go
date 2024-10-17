package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAppendChildWorkspace{}

func NewMsgAppendChildWorkspace(creator string, parentWorkspaceAddr string, childWorkspaceAddr string, btl uint64) *MsgAppendChildWorkspace {
	return &MsgAppendChildWorkspace{
		Creator:             creator,
		ParentWorkspaceAddr: parentWorkspaceAddr,
		ChildWorkspaceAddr:  childWorkspaceAddr,
		Btl:                 btl,
	}
}

func (msg *MsgAppendChildWorkspace) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	wsbz, err := sdk.GetFromBech32(msg.ParentWorkspaceAddr, PrefixWorkspaceAddress)
	if err != nil || len(wsbz) != WorkspaceAddressLength {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid parent workspace address (%s)", err)
	}

	wsbz, err = sdk.GetFromBech32(msg.ChildWorkspaceAddr, PrefixWorkspaceAddress)
	if err != nil || len(wsbz) != WorkspaceAddressLength {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid child workspace address (%s)", err)
	}

	return nil
}
