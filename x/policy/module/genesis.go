package policy

import (
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
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

	for _, p := range genState.Policies {
		if err := k.PolicyStore.Set(ctx, p.Id, p); err != nil {
			panic(err)
		}
	}

	if err := k.PolicyCount.Set(ctx, uint64(len(genState.Policies))); err != nil {
		panic(err)
	}

	for _, a := range genState.Actions {
		if err := k.ActionStore.Set(ctx, a.Id, a); err != nil {
			panic(err)
		}
	}

	if err := k.ActionCount.Set(ctx, uint64(len(genState.Actions))); err != nil {
		panic(err)
	}

	for _, a := range genState.SignMethods {
		if err := k.SignMethodStore.Set(ctx, collections.Join(a.Owner, a.Id), *a.Config); err != nil {
			panic(err)
		}
	}

	k.ParamStore.Set(ctx, genState.Params)
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
