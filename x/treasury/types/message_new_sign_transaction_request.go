package types

import (
	errorsmod "cosmossdk.io/errors"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgNewSignTransactionRequest{}

func NewMsgNewSignTransactionRequest(creator string, keyID uint64, walletType WalletType, unsignedTx []byte, meta *cdctypes.Any, btl uint64) *MsgNewSignTransactionRequest {
	return &MsgNewSignTransactionRequest{
		Creator:             creator,
		KeyId:               keyID,
		UnsignedTransaction: unsignedTx,
		WalletType:          walletType,
		Metadata:            meta,
		Btl:                 btl,
	}
}

func (msg *MsgNewSignTransactionRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
