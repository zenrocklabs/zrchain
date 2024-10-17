package identity_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v4/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/testutil/nullify"
	identity "github.com/Zenrock-Foundation/zrchain/v4/x/identity/module"
	"github.com/Zenrock-Foundation/zrchain/v4/x/identity/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # genesis/test/state
	}

	keepers := keepertest.NewTest(t)
	ik := keepers.IdentityKeeper
	ctx := keepers.Ctx
	identity.InitGenesis(ctx, *ik, genesisState)
	got := identity.ExportGenesis(ctx, *ik)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	// this line is used by starport scaffolding # genesis/test/assert
}
