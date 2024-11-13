package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ActionDetailsById(goCtx context.Context, req *types.QueryActionDetailsByIdRequest) (*types.QueryActionDetailsByIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	action, err := k.ActionStore.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var (
		policy types.Policy
	)
	approvers := []string{}
	pendingApprovers := []string{}

	if action.PolicyId > 0 {
		policy, err = k.PolicyStore.Get(ctx, action.PolicyId)
		if err != nil {
			return nil, err
		}

		pol, err := PolicyForAction(ctx, &k, &action)
		if err != nil {
			return nil, err
		}

		participantMap := map[string]string{}
		switch p := pol.(type) {
		case *types.BoolparserPolicy:
			for _, participant := range p.Participants {
				participantMap[participant.Address] = participant.Address
			}
		}

		for _, abbrev := range action.Approvers {
			if address, ok := participantMap[abbrev]; ok {
				approvers = append(approvers, address)
				delete(participantMap, abbrev)
			}
		}

		for _, address := range participantMap {
			pendingApprovers = append(pendingApprovers, address)
		}
	} else {
		approvers = action.Approvers
	}

	return &types.QueryActionDetailsByIdResponse{
		Id:               action.Id,
		Action:           &action,
		Policy:           &policy,
		Approvers:        approvers,
		PendingApprovers: pendingApprovers,
		CurrentHeight:    uint64(ctx.BlockHeight()),
	}, nil
}
