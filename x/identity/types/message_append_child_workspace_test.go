package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v5/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgAppendChildWorkspace_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAppendChildWorkspace
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgAppendChildWorkspace{
				Creator:             "invalid_address",
				ParentWorkspaceAddr: sample.WorkspaceAddress(),
				ChildWorkspaceAddr:  sample.WorkspaceAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid parent workspace address",
			msg: MsgAppendChildWorkspace{
				Creator:             sample.AccAddress(),
				ParentWorkspaceAddr: "invalid_workspace",
				ChildWorkspaceAddr:  sample.WorkspaceAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid child workspace address",
			msg: MsgAppendChildWorkspace{
				Creator:             sample.AccAddress(),
				ParentWorkspaceAddr: sample.WorkspaceAddress(),
				ChildWorkspaceAddr:  "invalid_workspace",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: MsgAppendChildWorkspace{
				Creator:             sample.AccAddress(),
				ParentWorkspaceAddr: sample.WorkspaceAddress(),
				ChildWorkspaceAddr:  sample.WorkspaceAddress(),
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
