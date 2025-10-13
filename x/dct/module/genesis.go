package dct

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.InitGenesis(ctx, &genState); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := keeper.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export
	err := k.ExportState(ctx, genesis)
	if err != nil {
		panic(err)
	}

	return genesis
}
