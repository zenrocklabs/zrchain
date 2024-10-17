package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgFulfilICATransactionRequest{}

func NewMsgFulfilICATransactionRequest(creator string, requestID uint64, status SignRequestStatus, keyringPartySig, signedData []byte, rejectReason string) *MsgFulfilICATransactionRequest {
	return &MsgFulfilICATransactionRequest{
		Creator:               creator,
		RequestId:             requestID,
		Status:                status,
		KeyringPartySignature: keyringPartySig,
		SignedData:            signedData,
		RejectReason:          rejectReason,
	}
}

func (msg *MsgFulfilICATransactionRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.RejectReason) > 512 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid reject reason, max len is 512")
	}
	return nil
}
