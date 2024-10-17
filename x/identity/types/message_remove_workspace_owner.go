package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgRemoveWorkspaceOwner{}

func NewMsgRemoveWorkspaceOwner(creator string, workspaceAddr string, owner string, btl uint64) *MsgRemoveWorkspaceOwner {
	return &MsgRemoveWorkspaceOwner{
		Creator:       creator,
		WorkspaceAddr: workspaceAddr,
		Owner:         owner,
		Btl:           btl,
	}
}

func (msg *MsgRemoveWorkspaceOwner) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	wsbz, err := sdk.GetFromBech32(msg.WorkspaceAddr, PrefixWorkspaceAddress)
	if err != nil || len(wsbz) != WorkspaceAddressLength {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid workspace address (%s)", err)
	}

	return nil
}
