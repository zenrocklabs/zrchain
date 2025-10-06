package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAcknowledgePoolTransfer{}

func NewMsgAcknowledgePoolTransfer(creator string, swapId uint64, sourceTxHash string, status SwapStatus, rejectReason string) *MsgAcknowledgePoolTransfer {
	return &MsgAcknowledgePoolTransfer{
		Creator:      creator,
		SwapId:       swapId,
		SourceTxHash: sourceTxHash,
		Status:       status,
		RejectReason: rejectReason,
	}
}

func (msg *MsgAcknowledgePoolTransfer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Status != SwapStatus_SWAP_STATUS_COMPLETED && msg.Status != SwapStatus_SWAP_STATUS_REJECTED {
		return fmt.Errorf("msg status is not completed or rejected: %s", msg.Status)
	}
	return nil
}
