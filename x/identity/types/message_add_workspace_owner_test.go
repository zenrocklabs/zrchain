package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v4/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgAddWorkspaceOwner_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAddWorkspaceOwner
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgAddWorkspaceOwner{
				Creator:       "invalid_address",
				WorkspaceAddr: sample.WorkspaceAddress(),
				NewOwner:      sample.AccAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid workspace address",
			msg: MsgAddWorkspaceOwner{
				Creator:       sample.AccAddress(),
				WorkspaceAddr: "invalid_workspace",
				NewOwner:      sample.AccAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid owner address",
			msg: MsgAddWorkspaceOwner{
				Creator:       sample.AccAddress(),
				WorkspaceAddr: sample.WorkspaceAddress(),
				NewOwner:      "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: MsgAddWorkspaceOwner{
				Creator:       sample.AccAddress(),
				WorkspaceAddr: sample.WorkspaceAddress(),
				NewOwner:      sample.AccAddress(),
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
