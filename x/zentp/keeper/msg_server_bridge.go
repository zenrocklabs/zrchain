package keeper

import (
	"context"

	"github.com/pkg/errors"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Bridge(goCtx context.Context, req *types.MsgBridge) (*types.MsgBridgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := req.ValidateBasic(); err != nil {
		return nil, err
	}

	if err := k.Keeper.CheckROCKSupplyCap(ctx, math.NewIntFromUint64(req.Amount)); err != nil {
		return nil, err
	}

	if _, err := validationtypes.ValidateSolanaChainID(goCtx, req.DestinationChain); err != nil {
		return nil, err
	}

	if treasurytypes.ValidateChainAddress(req.DestinationChain, req.RecipientAddress) != nil {
		return nil, errors.New("invalid recipient address: " + req.RecipientAddress)
	}

	if req.Denom != params.BondDenom {
		return nil, errors.New("invalid denomination")
	}

	baseAmountInt := sdkmath.NewIntFromUint64(req.Amount)

	totalAmountInt, totalFeeInt, err := k.CalculateZentpMintFee(ctx, req.Amount)
	if err != nil {
		return nil, err
	}

	bal := k.bankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(req.Creator), params.BondDenom)
	if bal.Amount.LT(totalAmountInt) {
		return nil, errors.New("not enough balance")
	}

	mintsCount, err := k.MintCount.Get(ctx)
	if err != nil {
		return nil, err
	}

	mintsCount++
	if err = k.MintStore.Set(ctx, mintsCount, types.Bridge{
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
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, baseAmountInt)),
	); err != nil {
		return nil, err
	}

	if err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		sdk.MustAccAddressFromBech32(req.Creator),
		types.ZentpCollectorName,
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, totalFeeInt)),
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
