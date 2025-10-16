package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgZenexTransferRequest{}

func NewMsgZenexTransferRequest(creator string, swapId uint64, unsignedPlusTx []byte, walletType treasurytypes.WalletType, cacheId []byte, dataForSigning []*InputHashes, rejectReason string) *MsgZenexTransferRequest {
	return &MsgZenexTransferRequest{
		Creator:        creator,
		SwapId:         swapId,
		UnsignedPlusTx: unsignedPlusTx,
		WalletType:     walletType,
		CacheId:        cacheId,
		DataForSigning: dataForSigning,
		RejectReason:   rejectReason,
	}
}

func (msg *MsgZenexTransferRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.WalletType != treasurytypes.WalletType_WALLET_TYPE_BTC_MAINNET &&
		msg.WalletType != treasurytypes.WalletType_WALLET_TYPE_BTC_TESTNET &&
		msg.WalletType != treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET &&
		msg.WalletType != treasurytypes.WalletType_WALLET_TYPE_ZCASH_MAINNET &&
		msg.WalletType != treasurytypes.WalletType_WALLET_TYPE_ZCASH_TESTNET {
		return fmt.Errorf("invalid wallet type: %s", msg.WalletType.String())
	}

	if msg.RejectReason != "" && (msg.CacheId != nil || msg.UnsignedPlusTx != nil || msg.DataForSigning != nil) {
		return ErrInvalidRejectMsg
	}

	return nil
}
