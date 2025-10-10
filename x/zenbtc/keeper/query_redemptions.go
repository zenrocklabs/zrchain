package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

func (k Keeper) GetRedemptions(goCtx context.Context, req *types.QueryRedemptionsRequest) (*types.QueryRedemptionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	redemptions := make([]types.Redemption, 0)
	var queryRange collections.Range[uint64]

	if err := k.Redemptions.Walk(ctx, queryRange.StartInclusive(req.StartIndex), func(key uint64, redemption types.Redemption) (bool, error) {
		switch req.Status {
		case types.RedemptionStatus_INITIATED:
			if redemption.Status == types.RedemptionStatus_INITIATED {
				redemptions = append(redemptions, redemption)
			}
		case types.RedemptionStatus_UNSTAKED:
			if redemption.Status == types.RedemptionStatus_UNSTAKED {
				redemptions = append(redemptions, redemption)
			}
		case types.RedemptionStatus_COMPLETED:
			if redemption.Status == types.RedemptionStatus_COMPLETED {
				redemptions = append(redemptions, redemption)
			}
		default: // don't filter
			redemptions = append(redemptions, redemption)
		}
		return false, nil
	}); err != nil {
		return nil, err
	}

	return &types.QueryRedemptionsResponse{Redemptions: redemptions}, nil
}
