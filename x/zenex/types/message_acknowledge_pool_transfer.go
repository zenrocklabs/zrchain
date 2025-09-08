package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAcknowledgePoolTransfer{}

func NewMsgAcknowledgePoolTransfer(creator string, swapId uint64, sourceTxHash string) *MsgAcknowledgePoolTransfer {
  return &MsgAcknowledgePoolTransfer{
		Creator: creator,
    SwapId: swapId,
    SourceTxHash: sourceTxHash,
	}
}

func (msg *MsgAcknowledgePoolTransfer) ValidateBasic() error {
  _, err := sdk.AccAddressFromBech32(msg.Creator)
  	if err != nil {
  		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
  	}
  return nil
}

