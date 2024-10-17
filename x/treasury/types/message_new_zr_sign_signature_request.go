package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgNewZrSignSignatureRequest{}

func NewMsgNewZrSignSignatureRequest(creator string, address string, keyType uint64, walletType WalletType, walletIndex uint64, cacheID []byte) *MsgNewZrSignSignatureRequest {
	return &MsgNewZrSignSignatureRequest{
		Creator:     creator,
		Address:     address,
		KeyType:     keyType,
		WalletIndex: walletIndex,
		WalletType:  walletType,
		CacheId:     cacheID,
	}
}

func (msg *MsgNewZrSignSignatureRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
