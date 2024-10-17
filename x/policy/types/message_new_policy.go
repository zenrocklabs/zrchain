package types

import (
	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgNewPolicy{}

const TypeMsgNewPolicy = "new_policy"

func NewMsgNewPolicy(creator string, name string, policy *codectypes.Any, btl uint64) *MsgNewPolicy {
	return &MsgNewPolicy{
		Creator: creator,
		Name:    name,
		Policy:  policy,
		Btl:     btl,
	}
}

func (msg *MsgNewPolicy) Route() string {
	return RouterKey
}

func (msg *MsgNewPolicy) Type() string {
	return TypeMsgNewPolicy
}

func (msg *MsgNewPolicy) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if len(msg.Name) > 255 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid name, max len 255")
	}

	return nil
}
