package keeper

import (
	"context"

	"github.com/pkg/errors"

	sdkmath "cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Bridge(goCtx context.Context, req *types.MsgBridge) (*types.MsgBridgeResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := treasurytypes.Caip2ToKeyType(req.DestinationChain)
	if err != nil {
		return nil, err
	}

	if treasurytypes.ValidateChainAddress(req.DestinationChain, req.RecipientAddress) != nil {
		return nil, errors.New("invalid recipient address: " + req.RecipientAddress)
	}

	totalAmount, err := k.AddFeeToBridgeAmount(ctx, req.Amount)
	if err != nil {
		return nil, err
	}

	p := k.GetSolanaParams(ctx)
	totalAmount = totalAmount + p.Fee // TODO: do this chain agnostic
	bal := k.bankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(req.Creator), req.Denom)
	if bal.IsLT(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(totalAmount))) {
		return nil, errors.New("not enough balance")
	}

	mintsCount, err := k.MintCount.Get(ctx)
	if err != nil {
		return nil, err
	}

	mintsCount++
	err = k.mintStore.Set(ctx, mintsCount, types.Bridge{
		Id:               mintsCount,
		Creator:          req.Creator,
		SourceAddress:    req.SourceAddress,
		SourceChain:      "cosmos:" + ctx.ChainID(),
		DestinationChain: req.DestinationChain,
		Amount:           req.Amount,
		Denom:            req.Denom,
		RecipientAddress: req.RecipientAddress,
		State:            types.BridgeStatus_BRIDGE_STATUS_PENDING,
	})
	if err != nil {
		return nil, err
	}

	// TODO check err
	k.MintCount.Set(ctx, mintsCount)

	if err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		sdk.MustAccAddressFromBech32(req.SourceAddress),
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(totalAmount))),
	); err != nil {
		return nil, err
	}

	if err := k.validationKeeper.SetSolanaRequestedNonce(goCtx, p.NonceAccountKey, true); err != nil {
		return nil, err
	}

	if err := k.validationKeeper.SetSolanaZenTPRequestedAccount(goCtx, req.RecipientAddress, true); err != nil {
		return nil, err
	}

	return &types.MsgBridgeResponse{
		Id: mintsCount,
	}, nil
}
