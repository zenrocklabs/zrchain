package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v5/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgFulfilSignatureRequest_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgFulfilSignatureRequest
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgFulfilSignatureRequest{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid reject reason",
			msg: MsgFulfilSignatureRequest{
				Creator:      "invalid_address",
				RejectReason: sample.StringLen(513),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid",
			msg: MsgFulfilSignatureRequest{
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
