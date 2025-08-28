package types

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgInitDct_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgInitDct
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgInitDct{
				Creator:          "invalid_address",
				Amount:           sdk.NewCoin("test", math.NewInt(100)),
				DestinationChain: "test",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgInitDct{
				Creator:          sample.AccAddress(),
				Amount:           sdk.NewCoin("test", math.NewInt(100)),
				DestinationChain: "test",
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
