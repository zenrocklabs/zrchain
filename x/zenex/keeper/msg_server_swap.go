package keeper

import (
	"context"
	"errors"
	"slices"

	errorsmod "cosmossdk.io/errors"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) Swap(goCtx context.Context, msg *types.MsgSwap) (*types.MsgSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	workspace, err := k.identityKeeper.GetWorkspace(ctx, msg.Workspace)
	if err != nil {
		return nil, err
	}

	if !workspace.IsOwner(msg.Creator) && msg.Creator != k.GetParams(ctx).Btcproxyaddress {
		return nil, errors.New("sender key is not the owner of the workspace")
	}

	senderKey, err := k.treasuryKeeper.GetKey(ctx, msg.SenderKey)
	if err != nil {
		return nil, err
	}

	recipientKey, err := k.treasuryKeeper.GetKey(ctx, msg.RecipientKey)
	if err != nil {
		return nil, err
	}

	if senderKey.WorkspaceAddr != msg.Workspace || recipientKey.WorkspaceAddr != msg.Workspace {
		return nil, errors.New("sender key is not in the workspace")
	}

	if !slices.Contains(types.ValidPairTypes, msg.Pair) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid keytype %s, valid types %+v", msg.Pair, types.ValidPairTypes)
	}

	pair, err := k.GetPair(ctx, msg.Pair)
	if err != nil {
		return nil, err
	}

	rockBtcPrice, err := k.validationKeeper.GetRockBtcPrice(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: this calc is a placeholder
	amountOut := msg.AmountIn.Mul(rockBtcPrice).Abs()

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
			Price:      rockBtcPrice,
			AmountIn:   msg.AmountIn,
			AmountOut:  amountOut,
		},
		SenderKeyId:    msg.SenderKey,
		RecipientKeyId: msg.RecipientKey,
		Workspace:      msg.Workspace,
		ZenbtcYield:    msg.Yield,
	}

	err = k.SwapsStore.Set(ctx, swapCount, swap)
	if err != nil {
		return nil, err
	}

	err = k.SwapsCount.Set(ctx, swapCount)
	if err != nil {
		return nil, err
	}

	return &types.MsgSwapResponse{Id: swap.SwapId}, nil
}
