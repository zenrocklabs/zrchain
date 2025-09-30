package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgAcknowledgePoolTransfer_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAcknowledgePoolTransfer
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgAcknowledgePoolTransfer{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgAcknowledgePoolTransfer{
				Creator: sample.AccAddress(),
				Status:  SwapStatus_SWAP_STATUS_COMPLETED,
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
