package types

import (
	errorsmod "cosmossdk.io/errors"
	identitytypes "github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgTransferFromKeyring{}

func NewMsgTransferFromKeyring(creator string, keyring string, recipient string, amount uint64, denom string) *MsgTransferFromKeyring {
	return &MsgTransferFromKeyring{
		Creator:   creator,
		Keyring:   keyring,
		Recipient: recipient,
		Amount:    amount,
		Denom:     denom,
	}
}

func (msg *MsgTransferFromKeyring) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
	}

	krbz, err := sdk.GetFromBech32(msg.Keyring, identitytypes.PrefixKeyringAddress)
	if err != nil || len(krbz) != identitytypes.KeyringAddressLength {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid keyring address (%s)", err)
	}

	return nil
}
