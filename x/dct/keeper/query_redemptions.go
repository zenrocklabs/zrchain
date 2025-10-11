package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

func (k Keeper) GetRedemptions(goCtx context.Context, req *types.QueryRedemptionsRequest) (*types.QueryRedemptionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	redemptions := make([]types.Redemption, 0)

	// Walk through all assets if asset is unspecified, otherwise filter by asset
	asset := req.Asset
	if asset == types.Asset_ASSET_UNSPECIFIED {
		// Query all redemptions across all assets
		if err := k.Redemptions.Walk(ctx, nil, func(key collections.Pair[string, uint64], redemption types.Redemption) (bool, error) {
			if key.K2() < req.StartIndex {
				return false, nil
			}
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
	} else {
		// Query redemptions for a specific asset
		if err := k.WalkRedemptions(ctx, asset, func(id uint64, redemption types.Redemption) (bool, error) {
			if id < req.StartIndex {
				return false, nil
			}
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
	}

	return &types.QueryRedemptionsResponse{Redemptions: redemptions}, nil
}
