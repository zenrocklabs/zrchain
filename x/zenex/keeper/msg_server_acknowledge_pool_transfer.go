package keeper

import (
	"context"
	"fmt"
	"strconv"

	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AcknowledgePoolTransfer(goCtx context.Context, msg *types.MsgAcknowledgePoolTransfer) (*types.MsgAcknowledgePoolTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	if msg.Creator != k.GetParams(ctx).BtcProxyAddress {
		return nil, fmt.Errorf("message creator is not the btc proxy address: %s", msg.Creator)
	}

	swap, err := k.SwapsStore.Get(ctx, msg.SwapId)
	if err != nil {
		return nil, err
	}

	if swap.Status != types.SwapStatus_SWAP_STATUS_REQUESTED {
		return nil, fmt.Errorf("swap status is not swap transfer requested: %s", swap.Status)
	}

	if msg.Status == types.SwapStatus_SWAP_STATUS_REJECTED {
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
		return &types.MsgAcknowledgePoolTransferResponse{}, nil
	}

	if swap.Pair == types.TradePair_TRADE_PAIR_BTC_ROCK {

		rockAddress, err := k.GetRockAddress(ctx, swap.RockKeyId)
		if err != nil {
			return nil, err
		}

		// Sending Swap.AmountOut to the rock address
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ZenexCollectorName, sdk.MustAccAddressFromBech32(rockAddress), sdk.NewCoins(sdk.NewCoin(params.BondDenom, math.NewIntFromUint64(swap.Data.AmountOut))))
		if err != nil {
			return nil, err
		}
	}

	swap.Status = msg.Status
	swap.SourceTxHash = msg.SourceTxHash
	err = k.SwapsStore.Set(ctx, swap.SwapId, swap)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventAcknowledgePoolTransfer,
			sdk.NewAttribute(types.AttributeSwapId, strconv.FormatUint(swap.SwapId, 10)),
			sdk.NewAttribute(types.AttributeNewSwapStatus, swap.Status.String()),
			sdk.NewAttribute(types.AttributePair, swap.Pair.String()),
		),
	})

	return &types.MsgAcknowledgePoolTransferResponse{}, nil
}
