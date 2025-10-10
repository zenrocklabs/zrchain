package zenbtc_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/testutil/nullify"

	"github.com/zenrocklabs/zenbtc/x/zenbtc/keeper"
	zenbtc "github.com/zenrocklabs/zenbtc/x/zenbtc/module"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: *keeper.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ZenbtcKeeper(t)
	zenbtc.InitGenesis(ctx, k, genesisState)
	got := zenbtc.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
