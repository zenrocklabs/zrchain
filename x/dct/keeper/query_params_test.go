package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

func TestParamsQuery(t *testing.T) {
	dctKeeper, ctx := DctKeeper(t)
	params := keeper.DefaultParams()
	require.NoError(t, dctKeeper.Params.Set(ctx, *params))

	response, err := dctKeeper.QueryParams(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: *params}, response)
}
