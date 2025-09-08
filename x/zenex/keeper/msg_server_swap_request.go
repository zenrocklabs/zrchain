package keeper

import (
	"context"
	"errors"
	"fmt"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SwapRequest(goCtx context.Context, msg *types.MsgSwapRequest) (*types.MsgSwapRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	zenexPoolKeyId := k.GetParams(ctx).ZenexPoolKeyId

	workspace, err := k.identityKeeper.GetWorkspace(ctx, msg.Workspace)
	if err != nil {
		return nil, err
	}

	if !workspace.IsOwner(msg.Creator) && msg.Creator != k.GetParams(ctx).BtcProxyAddress {
		return nil, errors.New("message creator is not the owner of the workspace or the btc proxy address")
	}

	rockKey, err := k.treasuryKeeper.GetKey(ctx, msg.RockKeyId)
	if err != nil {
		return nil, err
	}

	btcKey, err := k.treasuryKeeper.GetKey(ctx, msg.BtcKeyId)
	if err != nil {
		return nil, err
	}

	if rockKey.WorkspaceAddr != msg.Workspace || btcKey.WorkspaceAddr != msg.Workspace {
		return nil, fmt.Errorf("rock key %d or btc key %d is not in the workspace %s", msg.RockKeyId, msg.BtcKeyId, msg.Workspace)
	}

	pair, price, err := k.GetPair(ctx, msg.Pair)
	if err != nil {
		return nil, err
	}

	// either returns BTC or ROCK amount to transfer
	// checks if the amount in is greater than the minimum satoshis
	amountOut, err := k.GetAmountOut(ctx, msg.Pair, msg.AmountIn, price)
	if err != nil {
		return nil, err
	}

	switch msg.Pair {
	case "rockbtc":
		err = k.EscrowRock(ctx, *rockKey, msg.AmountIn)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid pair: %s", msg.Pair)
	}

	swapCount, err := k.SwapsCount.Get(ctx)
	if err != nil {
		return nil, err
	}

	swapCount++
	swap := types.Swap{
		Creator: msg.Creator,
		SwapId:  swapCount,
		Status:  types.SwapStatus_SWAP_STATUS_REQUESTED,
		Pair:    msg.Pair,
		Data: &types.SwapData{
			BaseToken:  pair.BaseToken,
			QuoteToken: pair.QuoteToken,
			Price:      price,
			AmountIn:   msg.AmountIn,
			AmountOut:  amountOut,
		},
		RockKeyId:      msg.RockKeyId,
		BtcKeyId:       msg.BtcKeyId,
		ZenexPoolKeyId: zenexPoolKeyId,
		Workspace:      msg.Workspace,
	}

	err = k.SwapsStore.Set(ctx, swapCount, swap)
	if err != nil {
		return nil, err
	}

	err = k.SwapsCount.Set(ctx, swapCount)
	if err != nil {
		return nil, err
	}

	return &types.MsgSwapRequestResponse{SwapId: swap.SwapId}, nil
}
