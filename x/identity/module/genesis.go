package identity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetPort(ctx, genState.PortId)
	// Only try to bind to port if it is not already bound, since we may already own
	// port capability from capability InitGenesis
	if k.ShouldBound(ctx, genState.PortId) {
		// module binds to the port on InitChain
		// and claims the returned capability
		err := k.BindPort(ctx, genState.PortId)
		if err != nil {
			panic("could not claim port capability: " + err.Error())
		}
	}

	if err := k.ParamStore.Set(ctx, genState.Params); err != nil {
		panic(err)
	}

	for i := range genState.Keyrings {
		if err := k.KeyringStore.Set(ctx, genState.Keyrings[i].Address, genState.Keyrings[i]); err != nil {
			panic(err)
		}
	}

	if err := k.KeyringCount.Set(ctx, uint64(len(genState.Keyrings))); err != nil {
		panic(err)
	}

	for i := range genState.Workspaces {
		if err := k.WorkspaceStore.Set(ctx, genState.Workspaces[i].Address, genState.Workspaces[i]); err != nil {
			panic(err)
		}
	}

	if err := k.WorkspaceCount.Set(ctx, uint64(len(genState.Workspaces))); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	params, err := k.ParamStore.Get(ctx)
	if err != nil {
		panic(err)
	}
	genesis.Params = params
	genesis.PortId = k.GetPort(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	err = k.ExportState(ctx, genesis)
	if err != nil {
		panic(err)
	}

	return genesis
}
