package zentp_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/testutil/nullify"
	zentp "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:           types.DefaultParams(),
		Burns:            testutil.DefaultBurns,
		Mints:            testutil.DefaultMints,
		SolanaRockSupply: 1000000,
		ZentpFees:        types.DefaultGenesis().ZentpFees,

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ZentpKeeper(t)
	zentp.InitGenesis(ctx, k, genesisState)
	got := zentp.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
