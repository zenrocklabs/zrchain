package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAddWorkspaceOwner{}

func NewMsgAddWorkspaceOwner(creator string, workspaceAddr string, newOwner string, btl uint64) *MsgAddWorkspaceOwner {
	return &MsgAddWorkspaceOwner{
		Creator:       creator,
		WorkspaceAddr: workspaceAddr,
		NewOwner:      newOwner,
		Btl:           btl,
	}
}

func (msg *MsgAddWorkspaceOwner) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	wsbz, err := sdk.GetFromBech32(msg.WorkspaceAddr, PrefixWorkspaceAddress)
	if err != nil || len(wsbz) != WorkspaceAddressLength {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid workspace address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.NewOwner)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	return nil
}
