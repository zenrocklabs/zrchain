package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgNewWorkspace{}

func NewMsgNewWorkspace(creator string, adminPolicyID uint64, signPolicyID uint64, additionalOwners ...string) *MsgNewWorkspace {
	return &MsgNewWorkspace{
		Creator:          creator,
		AdminPolicyId:    adminPolicyID,
		SignPolicyId:     signPolicyID,
		AdditionalOwners: additionalOwners,
	}
}

func (msg *MsgNewWorkspace) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	for _, addr := range msg.AdditionalOwners {
		_, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address for additional owner %s (%s)", addr, err)
		}
	}
	return nil
}
