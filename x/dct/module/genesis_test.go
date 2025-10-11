package dct_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/testutil/nullify"

	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/keeper"
	dct "github.com/Zenrock-Foundation/zrchain/v6/x/dct/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: *keeper.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.DctKeeper(t)
	dct.InitGenesis(ctx, k, genesisState)
	got := dct.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
