package keeper

import (
	"context"
	"errors"
	"strings"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
)

func (q Querier) SolanaCounters(ctx context.Context, req *types.QuerySolanaCountersRequest) (*types.QuerySolanaCountersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	resp := &types.QuerySolanaCountersResponse{
		Counters: make(map[string]*types.SolanaCounters),
	}

	store := q.Keeper.SolanaCounters
	asset := strings.TrimSpace(req.Asset)
	if asset != "" {
		counters, err := store.Get(sdkCtx, asset)
		if err != nil {
			if errors.Is(err, collections.ErrNotFound) {
				return resp, nil
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
		copy := counters
		resp.Counters[asset] = &copy
		return resp, nil
	}

	if err := store.Walk(sdkCtx, nil, func(assetKey string, counters types.SolanaCounters) (bool, error) {
		copy := counters
		resp.Counters[assetKey] = &copy
		return false, nil
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
}
