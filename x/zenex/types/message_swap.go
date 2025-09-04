package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSwap{}

func NewMsgSwap(creator string, pair string, workspace string, amount string, yield bool, senderKey uint64, recipientKey uint64) *MsgSwap {

	amountIn, err := sdkmath.LegacyNewDecFromStr(amount)
	if err != nil {
		panic(err)
	}

	return &MsgSwap{
		Creator:      creator,
		Pair:         pair,
		Workspace:    workspace,
		AmountIn:     amountIn,
		Yield:        yield,
		SenderKey:    senderKey,
		RecipientKey: recipientKey,
	}
}

func (msg *MsgSwap) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
