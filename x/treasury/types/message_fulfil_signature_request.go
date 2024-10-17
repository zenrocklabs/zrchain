package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgFulfilSignatureRequest{}

func NewMsgFulfilSignatureRequest(creator string, requestID uint64, status SignRequestStatus, keyringPartySignature, signedData []byte, rejectReason string) *MsgFulfilSignatureRequest {
	return &MsgFulfilSignatureRequest{
		Creator:               creator,
		RequestId:             requestID,
		Status:                status,
		KeyringPartySignature: keyringPartySignature,
		SignedData:            signedData,
		RejectReason:          rejectReason,
	}
}

func (msg *MsgFulfilSignatureRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.RejectReason) > 512 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid reject reason, max len is 512")
	}
	return nil
}
