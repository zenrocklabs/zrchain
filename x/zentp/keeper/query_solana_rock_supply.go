package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

func (k Keeper) QuerySolanaROCKSupply(goCtx context.Context, req *types.QuerySolanaROCKSupplyRequest) (*types.QuerySolanaROCKSupplyResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	supply, err := k.GetSolanaROCKSupply(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QuerySolanaROCKSupplyResponse{
		Amount: supply.Uint64(),
	}, nil
}
