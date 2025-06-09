package keeper

import (
	"context"

	"github.com/pkg/errors"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Bridge(goCtx context.Context, req *types.MsgBridge) (*types.MsgBridgeResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	const rockCap = 1_000_000_000_000_000 // 1bn ROCK in urock
	solanaSupply, err := k.GetSolanaROCKSupply(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get solana rock supply")
	}

	pendingMints, err := k.GetMintsWithStatus(ctx, types.BridgeStatus_BRIDGE_STATUS_PENDING)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pending mints")
	}

	var pendingAmount math.Int = math.ZeroInt()
	for _, mint := range pendingMints {
		pendingAmount = pendingAmount.Add(math.NewIntFromUint64(mint.Amount))
	}

	zrchainSupply := k.bankKeeper.GetSupply(ctx, params.BondDenom).Amount
	// Check if current supply + pending supply + new bridge amount > cap
	totalSupply := zrchainSupply.Add(solanaSupply).Add(pendingAmount).Add(math.NewIntFromUint64(req.Amount))
	if totalSupply.GT(sdkmath.NewIntFromUint64(rockCap)) {
		return nil, errors.Errorf("total ROCK supply including pending bridges would exceed cap (%s), bridge disabled", totalSupply.String())
	}

	if _, err := treasurytypes.Caip2ToKeyType(req.DestinationChain); err != nil {
		return nil, err
	}

	if treasurytypes.ValidateChainAddress(req.DestinationChain, req.RecipientAddress) != nil {
		return nil, errors.New("invalid recipient address: " + req.RecipientAddress)
	}

	if req.Denom != params.BondDenom {
		return nil, errors.New("invalid denomination")
	}

	baseAmount, err := k.AddFeeToBridgeAmount(ctx, req.Amount)
	if err != nil {
		return nil, err
	}

	// Use safe math to prevent overflow
	baseAmountInt := sdkmath.NewIntFromUint64(baseAmount)
	feeInt := sdkmath.NewIntFromUint64(k.GetSolanaParams(ctx).Fee)
	totalAmountInt := baseAmountInt.Add(feeInt)

	bal := k.bankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(req.Creator), params.BondDenom)
	if bal.Amount.LT(totalAmountInt) {
		return nil, errors.New("not enough balance")
	}

	mintsCount, err := k.MintCount.Get(ctx)
	if err != nil {
		return nil, err
	}

	mintsCount++
	if err = k.mintStore.Set(ctx, mintsCount, types.Bridge{
		Id:               mintsCount,
		Creator:          req.Creator,
		SourceAddress:    req.Creator,
		SourceChain:      "cosmos:" + ctx.ChainID(),
		DestinationChain: req.DestinationChain,
		Amount:           req.Amount,
		Denom:            req.Denom,
		RecipientAddress: req.RecipientAddress,
		State:            types.BridgeStatus_BRIDGE_STATUS_PENDING,
	}); err != nil {
		return nil, err
	}

	if err = k.MintCount.Set(ctx, mintsCount); err != nil {
		return nil, err
	}

	if err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		sdk.MustAccAddressFromBech32(req.Creator),
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, totalAmountInt)),
	); err != nil {
		return nil, err
	}

	if err = k.validationKeeper.SetSolanaRequestedNonce(goCtx, k.GetSolanaParams(ctx).NonceAccountKey, true); err != nil {
		return nil, err
	}

	if err = k.validationKeeper.SetSolanaZenTPRequestedAccount(goCtx, req.RecipientAddress, true); err != nil {
		return nil, err
	}

	return &types.MsgBridgeResponse{
		Id: mintsCount,
	}, nil
}
