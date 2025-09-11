package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgSwap_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSwapRequest
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSwapRequest{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgSwapRequest{
				Creator:   sample.AccAddress(),
				Pair:      TradePair_TRADE_PAIR_ROCK_BTC,
				Workspace: sample.WorkspaceAddress(),
				AmountIn:  100000,
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
