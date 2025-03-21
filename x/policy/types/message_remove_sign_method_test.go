package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgRemoveSignMethod_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgRemoveSignMethod
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgRemoveSignMethod{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid id",
			msg: MsgRemoveSignMethod{
				Creator: sample.AccAddress(),
				Id:      sample.StringLen(513),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid address",
			msg: MsgRemoveSignMethod{
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
