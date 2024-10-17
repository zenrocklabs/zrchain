package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v4/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgAddKeyringAdmin_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAddKeyringAdmin
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgAddKeyringAdmin{
				Creator:     "invalid_address",
				Admin:       sample.AccAddress(),
				KeyringAddr: sample.KeyringAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid admin address",
			msg: MsgAddKeyringAdmin{
				Creator:     sample.AccAddress(),
				Admin:       "invalid_address",
				KeyringAddr: sample.KeyringAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid keyring address",
			msg: MsgAddKeyringAdmin{
				Creator:     sample.AccAddress(),
				Admin:       sample.AccAddress(),
				KeyringAddr: "invalid_keyring",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: MsgAddKeyringAdmin{
				Creator:     sample.AccAddress(),
				Admin:       sample.AccAddress(),
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
