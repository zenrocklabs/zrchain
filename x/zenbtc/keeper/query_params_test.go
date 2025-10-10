package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"

	"github.com/zenrocklabs/zenbtc/x/zenbtc/keeper"

	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

func TestParamsQuery(t *testing.T) {
	zenBTCKeeper, ctx := keepertest.ZenbtcKeeper(t)
	params := keeper.DefaultParams()
	require.NoError(t, zenBTCKeeper.Params.Set(ctx, *params))

	response, err := zenBTCKeeper.QueryParams(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: *params}, response)
}
