package keeper

import (
	"context"
	"fmt"

	"github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddSignMethod(goCtx context.Context, msg *types.MsgAddSignMethod) (*types.MsgAddSignMethodResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var signMethod types.SignMethod
	err := k.cdc.UnpackAny(msg.GetConfig(), &signMethod)
	if err != nil {
		return nil, err
	}

	if signMethod == nil {
		return nil, fmt.Errorf("empty config")
	}

	if err := signMethod.VerifyConfig(ctx); err != nil {
		return nil, fmt.Errorf("invalid config: %s", err)
	}

	signMethod.SetActive(true)

	if err := k.Keeper.SetSignMethod(ctx, msg.Creator, signMethod.GetConfigId(), signMethod); err != nil {
		return nil, err
	}

	return &types.MsgAddSignMethodResponse{}, nil
}
