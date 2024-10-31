package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v5/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgAddKeyringParty_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAddKeyringParty
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgAddKeyringParty{
				Creator:     "invalid_address",
				KeyringAddr: sample.KeyringAddress(),
				Party:       sample.AccAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid keyring address",
			msg: MsgAddKeyringParty{
				Creator:     sample.AccAddress(),
				KeyringAddr: "invalid_keyring",
				Party:       sample.AccAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid party address",
			msg: MsgAddKeyringParty{
				Creator:     sample.AccAddress(),
				KeyringAddr: sample.KeyringAddress(),
				Party:       "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: MsgAddKeyringParty{
				Creator:     sample.AccAddress(),
				KeyringAddr: sample.KeyringAddress(),
				Party:       sample.AccAddress(),
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
