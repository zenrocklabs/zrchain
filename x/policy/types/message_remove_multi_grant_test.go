package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgRemoveMultiGrant_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgRemoveMultiGrant
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgRemoveMultiGrant{
				Creator: "invalid_address",
				Grantee: sample.AccAddress(),
				Msgs:    []string{"/zrchain.policy.MsgApproveAction"},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid grantee",
			msg: MsgRemoveMultiGrant{
				Creator: sample.AccAddress(),
				Grantee: "invalid_address",
				Msgs:    []string{"/zrchain.policy.MsgApproveAction"},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid msg",
			msg: MsgRemoveMultiGrant{
				Creator: sample.AccAddress(),
				Grantee: sample.AccAddress(),
				Msgs:    []string{sample.StringLen(513)},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid address",
			msg: MsgRemoveMultiGrant{
				Creator: sample.AccAddress(),
				Grantee: sample.AccAddress(),
				Msgs:    []string{"/zrchain.policy.MsgApproveAction"},
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
