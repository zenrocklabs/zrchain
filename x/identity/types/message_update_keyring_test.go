package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateKeyring_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateKeyring
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateKeyring{
				Creator:     "invalid_address",
				KeyringAddr: sample.KeyringAddress(),
				Description: "valid description",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid keyring address",
			msg: MsgUpdateKeyring{
				Creator:     sample.AccAddress(),
				KeyringAddr: "invalid keyring",
				Description: "valid description",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid description",
			msg: MsgUpdateKeyring{
				Creator:     sample.AccAddress(),
				KeyringAddr: sample.KeyringAddress(),
				Description: sample.StringLen(256),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid address",
			msg: MsgUpdateKeyring{
				Creator:     sample.AccAddress(),
				KeyringAddr: sample.KeyringAddress(),
				Description: "valid description",
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
