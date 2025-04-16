package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.ZentpKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.ParamStore.Set(ctx, params))
	p, err := k.ParamStore.Get(ctx)
	require.NoError(t, err)
	require.EqualValues(t, params, p)
}
