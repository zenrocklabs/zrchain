package keeper

import (
	"context"
	"fmt"
	"strconv"

	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ZenexTransferRequest(goCtx context.Context, msg *types.MsgZenexTransferRequest) (*types.MsgZenexTransferRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	swap, err := k.SwapsStore.Get(ctx, msg.SwapId)
	if err != nil {
		return nil, err
	}

	if msg.Creator != k.GetParams(ctx).BtcProxyAddress && msg.Creator != swap.Creator {
		return nil, fmt.Errorf("message creator is not the btc proxy address")
	}

	if swap.Status != types.SwapStatus_SWAP_STATUS_REQUESTED {
		return nil, fmt.Errorf("swap status is not requested")
	}

	var senderKeyId uint64
	var txCreator string
	switch swap.Pair {
	case types.TradePair_TRADE_PAIR_ROCK_BTC:
		senderKeyId = swap.BtcKeyId
		txCreator = swap.Creator
	case types.TradePair_TRADE_PAIR_BTC_ROCK:
		senderKeyId = swap.ZenexPoolKeyId
		txCreator = k.GetParams(ctx).BtcProxyAddress
	default:
		return nil, fmt.Errorf("invalid pair: %s", swap.Pair)
	}

	bitcoinTx := treasurytypes.NewMsgNewSignTransactionRequest(
		txCreator,
		[]uint64{senderKeyId},
		msg.WalletType,
		msg.UnsignedTx,
		nil,
		treasurytypes.DefaultParams().DefaultBtl,
	)

	signTxResponse, err := k.treasuryKeeper.MakeSignTransactionRequest(ctx, bitcoinTx)
	if err != nil {
		return nil, err
	}

	swap.SignTxId = signTxResponse.Id
	swap.Status = types.SwapStatus_SWAP_STATUS_SWAP_TRANSFER_REQUESTED
	err = k.SwapsStore.Set(ctx, swap.SwapId, swap)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventZenexTransferRequest,
			sdk.NewAttribute(types.AttributeSwapId, strconv.FormatUint(swap.SwapId, 10)),
			sdk.NewAttribute(types.AttributeNewSwapStatus, swap.Status.String()),
			sdk.NewAttribute(types.AttributePair, swap.Pair.String()),
		),
	})

	return &types.MsgZenexTransferRequestResponse{
		SignTxId: signTxResponse.Id,
	}, nil
}
