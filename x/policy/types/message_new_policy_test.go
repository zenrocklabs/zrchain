package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgNewPolicy_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgNewPolicy
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgNewPolicy{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid name",
			msg: MsgNewPolicy{
				Creator: sample.AccAddress(),
				Name:    sample.StringLen(256),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid",
			msg: MsgNewPolicy{
				Creator: sample.AccAddress(),
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
