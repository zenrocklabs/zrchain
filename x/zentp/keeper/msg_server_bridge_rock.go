package keeper

import (
	"context"

	"github.com/pkg/errors"

	sdkmath "cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v5/app/params"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v5/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) BridgeRock(goCtx context.Context, req *types.MsgBridgeRock) (*types.MsgBridgeRockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	//pa := k.GetParams(goCtx)
	//pa.Solana.NonceAuthorityPubKey = "success"
	//if err := k.SetParams(ctx, pa); err != nil {
	//	panic(err)
	//}
	_, err := treasurytypes.Caip2ToKeyType(req.DestinationChain)
	if err != nil {
		return nil, err
	}

	if treasurytypes.ValidateChainAddress(req.DestinationChain, req.RecipientAddress) != nil {
		return nil, errors.New("invalid recipient address: " + req.RecipientAddress)
	}
	p := k.GetParams(ctx)
	totalAmount := req.Amount + p.Solana.Fee // TODO: do this chain agnostic
	bal := k.bankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(req.Creator), params.BondDenom)
	if bal.IsLT(sdk.NewCoin("urock", sdkmath.NewIntFromUint64(totalAmount))) {
		return nil, errors.New("not enough balance")
	}

	mintsCount, err := k.MintCount.Get(ctx)
	if err != nil {
	}
	if err != nil {
		return nil, err
	}

	mintsCount++
	err = k.mintStore.Set(ctx, mintsCount, types.BridgeRock{
		Id:               mintsCount,
		SourceChain:      "cosmos:" + ctx.ChainID(),
		DestinationChain: req.DestinationChain,
		Amount:           req.Amount,
		RecipientAddress: req.RecipientAddress,
		State:            types.BridgeStatus_BRIDGE_STATUS_NEW,
	})
	if err != nil {
		return nil, err
	}

	if err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		sdk.MustAccAddressFromBech32(req.RecipientAddress),
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(totalAmount))),
	); err != nil {
		return nil, err
	}

	return &types.MsgBridgeRockResponse{
		Id: mintsCount,
	}, nil
}
