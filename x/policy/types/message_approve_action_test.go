package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgApproveAction_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgApproveAction
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgApproveAction{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid action type",
			msg: MsgApproveAction{
				Creator:    sample.AccAddress(),
				ActionType: sample.StringLen(513),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid address",
			msg: MsgApproveAction{
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
