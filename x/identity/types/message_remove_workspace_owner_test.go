package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v5/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgRemoveWorkspaceOwner_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgRemoveWorkspaceOwner
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgRemoveWorkspaceOwner{
				Creator:       "invalid_address",
				Owner:         sample.AccAddress(),
				WorkspaceAddr: sample.WorkspaceAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid owner",
			msg: MsgRemoveWorkspaceOwner{
				Creator:       sample.AccAddress(),
				Owner:         "invalid_address",
				WorkspaceAddr: sample.WorkspaceAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid workspace",
			msg: MsgRemoveWorkspaceOwner{
				Creator:       sample.AccAddress(),
				Owner:         sample.AccAddress(),
				WorkspaceAddr: "invalid_workspace",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid",
			msg: MsgRemoveWorkspaceOwner{
				Creator:       sample.AccAddress(),
				Owner:         sample.AccAddress(),
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
