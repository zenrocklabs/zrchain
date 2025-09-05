package zenex_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/testutil/nullify"
	zenex "github.com/Zenrock-Foundation/zrchain/v6/x/zenex/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ZenexKeeper(t)
	zenex.InitGenesis(ctx, k, genesisState)
	got := zenex.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
