package zenex

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	if err := k.SwapsCount.Set(ctx, uint64(len(genState.Swaps))); err != nil {
		panic(err)
	}

	for _, swap := range genState.Swaps {
		if err := k.SwapsStore.Set(ctx, swap.SwapId, swap); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	swaps, err := k.GetSwaps(ctx)
	if err != nil {
		panic(err)
	}
	genesis.Swaps = swaps

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
