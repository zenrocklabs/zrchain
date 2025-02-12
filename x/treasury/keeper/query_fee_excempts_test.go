package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
)

func TestFeeExcemptsQuery(t *testing.T) {
	keepers := keepertest.NewTest(t)
	tk := keepers.TreasuryKeeper
	ctx := keepers.Ctx
	noFeeMsgs := types.DefaultNoFeeMsgs()
	for _, msg := range noFeeMsgs {
		require.NoError(t, tk.NoFeeMsgsList.Set(ctx, msg))
	}

	response, err := tk.FeeExcempts(ctx, &types.QueryFeeExcemptsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryFeeExcemptsResponse{NoFeeMsgs: noFeeMsgs}, response)
}
