package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/Zenrock-Foundation/zrchain/v4/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/treasury/types"
)

func TestGetParams(t *testing.T) {
	keepers := keepertest.NewTest(t)
	tk := keepers.TreasuryKeeper
	ctx := keepers.Ctx
	params := types.DefaultParams()

	require.NoError(t, tk.ParamStore.Set(ctx, params))
	got, err := tk.ParamStore.Get(ctx)
	require.NoError(t, err)
	require.EqualValues(t, params, got)
}
