package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

func (k Keeper) GetLockTransactions(goCtx context.Context, req *types.QueryLockTransactionsRequest) (*types.QueryLockTransactionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	asset := req.Asset
	if asset == types.Asset_ASSET_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "asset must be specified")
	}

	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rangeForAsset := collections.NewPrefixedPairRange[string, string](assetKey)

	lockTransactions := []*types.LockTransaction{}
	if err := k.LockTransactions.Walk(ctx, rangeForAsset, func(key collections.Pair[string, string], value types.LockTransaction) (bool, error) {
		// ensure asset set (for legacy migrated entries)
		value.Asset = asset
		lockTransactions = append(lockTransactions, &value)
		return false, nil
	}); err != nil {
		return nil, err
	}

	return &types.QueryLockTransactionsResponse{LockTransactions: lockTransactions}, nil
}
