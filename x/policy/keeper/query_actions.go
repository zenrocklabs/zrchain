package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
)

func (k Keeper) Actions(goCtx context.Context, req *types.QueryActionsRequest) (*types.QueryActionsResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "request is nil")
	}

	actions, pageRes, err := query.CollectionFilteredPaginate[uint64, types.Action, collections.Map[uint64, types.Action], types.ActionResponse](
		goCtx,
		k.ActionStore,
		req.Pagination,
		func(key uint64, value types.Action) (bool, error) {
			if req.Address != "" {
				pol, err := PolicyForAction(sdk.UnwrapSDKContext(goCtx), &k, &value)
				if err != nil {
					return false, nil
				}
				if _, err = pol.AddressToParticipant(req.Address); err != nil {
					return false, nil
				}
			}
			if req.Status != types.ActionStatus_ACTION_STATUS_UNSPECIFIED && value.Status != req.Status {
				return false, nil
			}
			return true, nil
		},
		func(key uint64, value types.Action) (types.ActionResponse, error) {
			return types.ActionResponse{
				Id:         value.Id,
				Approvers:  value.Approvers,
				Status:     value.Status.String(),
				PolicyId:   value.PolicyId,
				Msg:        value.Msg,
				Creator:    value.Creator,
				Btl:        value.Btl,
				PolicyData: value.PolicyData,
			}, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.QueryActionsResponse{
		Actions:    actions,
		Pagination: pageRes,
	}, nil
}
