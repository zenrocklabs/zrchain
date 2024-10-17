package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/Zenrock-Foundation/zrchain/v4/testutil/keeper"

	"github.com/Zenrock-Foundation/zrchain/v4/x/treasury/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/treasury/types"
)

func setupMsgServer(t testing.TB) (keeper.Keeper, types.MsgServer, context.Context) {
	keepers := keepertest.NewTest(t)
	tk := keepers.TreasuryKeeper
	ctx := keepers.Ctx
	return *tk, keeper.NewMsgServerImpl(*tk), ctx
}

func TestMsgServer(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)
}
