package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgBurn_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgBurn
		err  error
	}{
		{
			name: "invalid authority",
			msg: MsgBurn{
				Authority:     "invalid_address",
				ModuleAccount: "zentp",
				Denom:         "urock",
				Amount:        100,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid message",
			msg: MsgBurn{
				Authority:     sample.GetAuthority(),
				ModuleAccount: "zentp",
				Denom:         "urock",
				Amount:        100,
			},
			err: nil,
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
