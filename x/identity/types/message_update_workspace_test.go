package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v4/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateWorkspace_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateWorkspace
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateWorkspace{
				Creator:       "invalid_address",
				WorkspaceAddr: sample.WorkspaceAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid workspace address",
			msg: MsgUpdateWorkspace{
				Creator:       sample.AccAddress(),
				WorkspaceAddr: "invalid_workspace",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid",
			msg: MsgUpdateWorkspace{
				Creator:       sample.AccAddress(),
				WorkspaceAddr: sample.WorkspaceAddress(),
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
