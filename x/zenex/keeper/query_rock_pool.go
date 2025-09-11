package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) RockPool(goCtx context.Context, req *types.QueryRockPoolRequest) (*types.QueryRockPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	balance := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(types.ZenexCollectorName), params.BondDenom)

	redeemableAssets, err := k.getRedeemableAssets(ctx, balance.Amount.Uint64())
	if err != nil {
		return nil, err
	}

	return &types.QueryRockPoolResponse{
		RockBalance:      balance.Amount.Uint64(),
		RedeemableAssets: redeemableAssets,
	}, nil
}

func (k Keeper) getRedeemableAssets(ctx context.Context, rockBalance uint64) ([]types.RedeemableAsset, error) {
	assets, err := k.validationKeeper.GetAssets(ctx)
	if err != nil {
		return nil, err
	}

	// rockBtcPrice, err := k.GetPrice(sdk.UnwrapSDKContext(ctx), types.TradePair_TRADE_PAIR_ROCK_BTC)
	// if err != nil {
	// 	return nil, err
	// }

	btcRockPrice, err := k.GetPrice(sdk.UnwrapSDKContext(ctx), types.TradePair_TRADE_PAIR_BTC_ROCK)
	if err != nil {
		return nil, err
	}

	redeemableAssets := make([]types.RedeemableAsset, 0)
	var amountOut uint64
	for _, asset := range assets {
		switch asset {
		case validationtypes.Asset_BTC:
			amountOut, err = k.GetAmountOut(sdk.UnwrapSDKContext(ctx), types.TradePair_TRADE_PAIR_BTC_ROCK, rockBalance, btcRockPrice)
			if err != nil {
				return nil, err
			}
		default:
			continue
		}
		redeemableAssets = append(redeemableAssets, types.RedeemableAsset{
			Asset:  asset,
			Amount: amountOut,
		})
	}

	return redeemableAssets, nil
}
