package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgFulfilKeyRequest{}

func NewMsgFulfilKeyRequest(creator string, requestId uint64, status KeyRequestStatus, result isMsgFulfilKeyRequest_Result, KeyringPartySignature []byte) *MsgFulfilKeyRequest {
	return &MsgFulfilKeyRequest{
		Creator:               creator,
		RequestId:             requestId,
		Status:                status,
		Result:                result,
		KeyringPartySignature: KeyringPartySignature,
	}
}

func (msg *MsgFulfilKeyRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
