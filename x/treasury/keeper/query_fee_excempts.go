package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
)

func (k Keeper) FeeExcempts(goCtx context.Context, req *types.QueryFeeExcemptsRequest) (*types.QueryFeeExcemptsResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "invalid arguments: request is nil")
	}

	noFeeMsgs := []string{}
	err := k.NoFeeMsgsList.Walk(goCtx, nil, func(msgUrl string) (stop bool, err error) {
		noFeeMsgs = append(noFeeMsgs, msgUrl)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryFeeExcemptsResponse{NoFeeMsgs: noFeeMsgs}, nil
}
