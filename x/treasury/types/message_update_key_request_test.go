package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v4/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateKeyPolicy_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateKeyPolicy
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateKeyPolicy{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateKeyPolicy{
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
