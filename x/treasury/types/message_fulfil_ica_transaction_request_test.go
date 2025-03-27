package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgFulfilICATransactionRequest_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgFulfilICATransactionRequest
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgFulfilICATransactionRequest{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid reject reason",
			msg: MsgFulfilICATransactionRequest{
				Creator:      sample.AccAddress(),
				RejectReason: sample.StringLen(513),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid",
			msg: MsgFulfilICATransactionRequest{
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
