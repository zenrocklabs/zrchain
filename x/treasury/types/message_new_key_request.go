package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgNewKeyRequest{}

func NewMsgNewKeyRequest(creator, workspaceAddr, keyringAddr, keyType string, btl, signPolicyId, mpcBtl uint64) *MsgNewKeyRequest {
	return &MsgNewKeyRequest{
		Creator:       creator,
		WorkspaceAddr: workspaceAddr,
		KeyringAddr:   keyringAddr,
		KeyType:       keyType,
		Btl:           btl,
		SignPolicyId:  signPolicyId,
		MpcBtl:        mpcBtl,
	}
}

func (msg *MsgNewKeyRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}
