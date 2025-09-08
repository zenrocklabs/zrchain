package keeper

import (
	"context"
	"fmt"

	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ZenexTransferRequest(goCtx context.Context, msg *types.MsgZenexTransferRequest) (*types.MsgZenexTransferRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Creator != k.GetParams(ctx).BtcProxyAddress {
		return nil, fmt.Errorf("message creator is not the btc proxy address")
	}

	swap, err := k.SwapsStore.Get(ctx, msg.SwapId)
	if err != nil {
		return nil, err
	}

	if swap.Status != types.SwapStatus_SWAP_STATUS_REQUESTED {
		return nil, fmt.Errorf("swap status is not requested")
	}

	var senderKeyId uint64
	switch swap.Pair {
	case "rockbtc":
		senderKeyId = swap.BtcKeyId
	case "btcrock":
		senderKeyId = swap.BtcChangeKeyId
	default:
		return nil, fmt.Errorf("invalid pair: %s", swap.Pair)
	}

	bitcoinTx := treasurytypes.NewMsgNewSignTransactionRequest(
		swap.Creator,
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

	return &types.MsgZenexTransferRequestResponse{
		SignTxId: signTxResponse.Id,
	}, nil
}
