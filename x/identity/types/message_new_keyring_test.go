package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v5/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgNewKeyring_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgNewKeyring
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgNewKeyring{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid description",
			msg: MsgNewKeyring{
				Creator:     sample.AccAddress(),
				Description: sample.StringLen(256),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid",
			msg: MsgNewKeyring{
				Creator:     sample.AccAddress(),
				Description: sample.StringLen(255),
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
