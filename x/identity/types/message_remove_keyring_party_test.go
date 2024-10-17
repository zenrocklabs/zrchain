package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v4/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgRemoveKeyringParty_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgRemoveKeyringParty
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgRemoveKeyringParty{
				Creator:     "invalid_address",
				Party:       sample.AccAddress(),
				KeyringAddr: sample.KeyringAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid keyring address",
			msg: MsgRemoveKeyringParty{
				Creator:     sample.AccAddress(),
				Party:       sample.AccAddress(),
				KeyringAddr: "invalid_keyring",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid party address",
			msg: MsgRemoveKeyringParty{
				Creator:     sample.AccAddress(),
				Party:       "invalid_address",
				KeyringAddr: sample.KeyringAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid",
			msg: MsgRemoveKeyringParty{
				Creator:     sample.AccAddress(),
				Party:       sample.AccAddress(),
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
