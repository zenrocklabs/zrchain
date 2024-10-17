package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateWorkspace{}

func NewMsgUpdateWorkspace(creator string, workspaceAddr string, adminPolicyId uint64, signPolicyId uint64, btl uint64) *MsgUpdateWorkspace {
	return &MsgUpdateWorkspace{
		Creator:       creator,
		WorkspaceAddr: workspaceAddr,
		AdminPolicyId: adminPolicyId,
		SignPolicyId:  signPolicyId,
		Btl:           btl,
	}
}

func (msg *MsgUpdateWorkspace) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	wsbz, err := sdk.GetFromBech32(msg.WorkspaceAddr, PrefixWorkspaceAddress)
	if err != nil || len(wsbz) != WorkspaceAddressLength {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid workspace address (%s)", err)
	}

	return nil
}
