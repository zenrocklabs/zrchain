package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateKeyPolicy{}

func NewMsgUpdateKeyPolicy(creator string, keyId uint64, signPolicyId uint64) *MsgUpdateKeyPolicy {
	return &MsgUpdateKeyPolicy{
		Creator:      creator,
		KeyId:        keyId,
		SignPolicyId: signPolicyId,
	}
}

func (msg *MsgUpdateKeyPolicy) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
