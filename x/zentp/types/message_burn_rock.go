package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgBurnRock{}

func NewMsgBurnRock(creator string, chainId string, keyId int32, amount int32, recipient string) *MsgBurnRock {
	return &MsgBurnRock{
		Creator:   creator,
		ChainId:   chainId,
		KeyId:     keyId,
		Amount:    amount,
		Recipient: recipient,
	}
}

func (msg *MsgBurnRock) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
