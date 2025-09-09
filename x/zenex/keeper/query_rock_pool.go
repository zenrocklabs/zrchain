package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
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

	return &types.QueryRockPoolResponse{
		RockBalance: balance.Amount.Uint64(),
	}, nil
}
