package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v5/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgTransferFromKeyring_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgTransferFromKeyring
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgTransferFromKeyring{
				Creator:   "invalid_address",
				Keyring:   sample.KeyringAddress(),
				Recipient: sample.AccAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid recipient",
			msg: MsgTransferFromKeyring{
				Creator:   sample.AccAddress(),
				Keyring:   sample.KeyringAddress(),
				Recipient: "invalid_recipient",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid keyring",
			msg: MsgTransferFromKeyring{
				Creator:   sample.AccAddress(),
				Keyring:   "invalid_keyring",
				Recipient: sample.AccAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid",
			msg: MsgTransferFromKeyring{
				Creator:   sample.AccAddress(),
				Recipient: sample.AccAddress(),
				Keyring:   sample.KeyringAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
