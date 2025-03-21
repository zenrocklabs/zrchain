package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgAddMultiGrant_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAddMultiGrant
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgAddMultiGrant{
				Creator: "invalid_address",
				Grantee: sample.AccAddress(),
				Msgs:    []string{"/zrchain.policy.MsgApproveAction"},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid grantee",
			msg: MsgAddMultiGrant{
				Creator: sample.AccAddress(),
				Grantee: "invalid_address",
				Msgs:    []string{"/zrchain.policy.MsgApproveAction"},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid msg",
			msg: MsgAddMultiGrant{
				Creator: sample.AccAddress(),
				Grantee: sample.AccAddress(),
				Msgs:    []string{sample.StringLen(513)},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid address",
			msg: MsgAddMultiGrant{
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
