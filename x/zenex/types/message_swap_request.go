package types

import (
	"errors"
	"slices"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSwapRequest{}

func NewMsgSwap(creator, pair, workspace, destinationCaip2 string, amountIn uint64, rockKeyId uint64, btcKeyId uint64) *MsgSwapRequest {

	return &MsgSwapRequest{
		Creator:   creator,
		Pair:      pair,
		Workspace: workspace,
		AmountIn:  amountIn,
		RockKeyId: rockKeyId,
		BtcKeyId:  btcKeyId,
	}
}

func (msg *MsgSwapRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.AmountIn == 0 {
		return errors.New("amount in is 0")
	}

	if msg.Pair == "" {
		return errors.New("pair is empty")
	}

	if msg.Workspace == "" {
		return errors.New("workspace is empty")
	}

	if !slices.Contains(ValidPairTypes, msg.Pair) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid keytype %s, valid types %+v", msg.Pair, ValidPairTypes)
	}

	return nil
}
