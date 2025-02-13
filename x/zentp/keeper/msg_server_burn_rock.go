package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v5/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) BurnRock(goCtx context.Context, msg *types.MsgBurnRock) (*types.MsgBurnRockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgBurnRockResponse{}, nil
}
