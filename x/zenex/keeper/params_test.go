package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.ZenexKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
