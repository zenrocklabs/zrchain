package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
)

func TestGetParams(t *testing.T) {
	keepers := keepertest.NewTest(t)
	ik := keepers.IdentityKeeper
	ctx := keepers.Ctx
	params := types.DefaultParams()

	require.NoError(t, ik.ParamStore.Set(ctx, params))
	got, err := ik.ParamStore.Get(ctx)
	require.NoError(t, err)
	require.EqualValues(t, params, got)
}
