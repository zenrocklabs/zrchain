package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgNewKeyRequest{}

func NewMsgNewKeyRequest(creator string, workspaceAddr string, keyringAddr string, keyType string, btl uint64, signPolicyId uint64) *MsgNewKeyRequest {
	return &MsgNewKeyRequest{
		Creator:       creator,
		WorkspaceAddr: workspaceAddr,
		KeyringAddr:   keyringAddr,
		KeyType:       keyType,
		Btl:           btl,
		SignPolicyId:  signPolicyId,
	}
}

func (msg *MsgNewKeyRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}
