package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v5/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgDeactivateKeyring_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeactivateKeyring
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeactivateKeyring{
				Creator:     "invalid_address",
				KeyringAddr: sample.KeyringAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid keyring address",
			msg: MsgDeactivateKeyring{
				Creator:     sample.AccAddress(),
				KeyringAddr: "invalid_keyring",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: MsgDeactivateKeyring{
				Creator:     sample.AccAddress(),
				KeyringAddr: sample.KeyringAddress(),
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
