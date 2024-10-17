package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgNewSignatureRequest{}

func NewMsgNewSignatureRequest(creator string, keyId uint64, dataForSigning string, btl uint64) *MsgNewSignatureRequest {
	return &MsgNewSignatureRequest{
		Creator:        creator,
		KeyId:          keyId,
		DataForSigning: dataForSigning,
		Btl:            btl,
	}
}

func (msg *MsgNewSignatureRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
