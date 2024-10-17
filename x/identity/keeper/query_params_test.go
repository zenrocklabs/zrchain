package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/Zenrock-Foundation/zrchain/v4/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/identity/types"
)

func TestParamsQuery(t *testing.T) {
	keepers := keepertest.NewTest(t)
	ik := keepers.IdentityKeeper
	ctx := keepers.Ctx
	params := types.DefaultParams()
	require.NoError(t, ik.ParamStore.Set(ctx, params))

	response, err := ik.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
