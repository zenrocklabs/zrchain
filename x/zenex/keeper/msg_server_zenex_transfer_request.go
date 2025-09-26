package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
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

	if swap.Status != types.SwapStatus_SWAP_STATUS_INITIATED {
		return nil, fmt.Errorf("swap status is not requested")
	}

	if msg.RejectReason != "" {
		swap.Status = types.SwapStatus_SWAP_STATUS_REJECTED
		swap.RejectReason = msg.RejectReason
		err = k.SwapsStore.Set(ctx, swap.SwapId, swap)
		if err != nil {
			return nil, err
		}
		// Release previously pending escrowed funds
		if swap.Pair == types.TradePair_TRADE_PAIR_ROCK_BTC {
			rockAddress, err := k.GetRockAddress(ctx, swap.RockKeyId)
			if err != nil {
				return nil, err
			}
			// Sending Swap.AmountIn to the mpc rock address
			err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ZenexCollectorName, sdk.MustAccAddressFromBech32(rockAddress), sdk.NewCoins(sdk.NewCoin(params.BondDenom, math.NewIntFromUint64(swap.Data.AmountIn))))
			if err != nil {
				return nil, err
			}
		}
		return &types.MsgZenexTransferRequestResponse{}, nil
	}

	keyIDs := make([]uint64, len(msg.DataForSigning))
	hashes := make([]string, len(msg.DataForSigning))
	for i, input := range msg.DataForSigning {
		keyIDs[i] = input.KeyId
		hashes[i] = input.Hash
	}

	signReq := &treasurytypes.MsgNewSignatureRequest{
		Creator:        msg.Creator,
		KeyIds:         keyIDs,
		DataForSigning: strings.Join(hashes, ","), // hex string, each unsigned utxo is separated by comma
		CacheId:        msg.CacheId,
	}

	signReqResponse, err := k.treasuryKeeper.HandleSignatureRequest(ctx, signReq)
	if err != nil {
		return nil, err
	}

	swap.SignReqId = signReqResponse.SigReqId
	swap.Status = types.SwapStatus_SWAP_STATUS_REQUESTED
	swap.UnsignedPlusTx = hex.EncodeToString(msg.UnsignedPlusTx)
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
		SignReqId: signReqResponse.SigReqId,
	}, nil
}
