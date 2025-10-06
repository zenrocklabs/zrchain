package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgZenexBitcoinTransfer_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgZenexTransferRequest
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgZenexTransferRequest{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgZenexTransferRequest{
				Creator:    sample.AccAddress(),
				WalletType: treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET,
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
