package zentp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.ParamStore.Set(ctx, genState.Params); err != nil {
		panic(err)
	}
	if err := k.MintCount.Set(ctx, uint64(0)); err != nil {
		panic(err)
	}
	if err := k.BurnCount.Set(ctx, uint64(0)); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	params, err := k.ParamStore.Get(ctx)
	if err != nil {
		params = types.DefaultParams()
	}
	genesis.Params = params

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
