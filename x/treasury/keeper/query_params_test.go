package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
)

func TestParamsQuery(t *testing.T) {
	keepers := keepertest.NewTest(t)
	tk := keepers.TreasuryKeeper
	ctx := keepers.Ctx
	params := types.DefaultParams()
	require.NoError(t, tk.ParamStore.Set(ctx, params))

	response, err := tk.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
