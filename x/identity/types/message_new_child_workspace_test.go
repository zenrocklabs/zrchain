package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgNewChildWorkspace_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgNewChildWorkspace
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgNewChildWorkspace{
				Creator:             "invalid_address",
				ParentWorkspaceAddr: sample.WorkspaceAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid workspace address",
			msg: MsgNewChildWorkspace{
				Creator:             sample.AccAddress(),
				ParentWorkspaceAddr: "invalid_workspace",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: MsgNewChildWorkspace{
				Creator:             sample.AccAddress(),
				ParentWorkspaceAddr: sample.WorkspaceAddress(),
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
