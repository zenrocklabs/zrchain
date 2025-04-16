package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := keepertest.ZentpKeeper(t)
	params := types.DefaultParams()
	require.NoError(t, keeper.ParamStore.Set(ctx, params))

	response, err := keeper.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
