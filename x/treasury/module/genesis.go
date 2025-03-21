package treasury

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetPort(ctx, genState.PortId)

	if k.ShouldBound(ctx, genState.PortId) {
		if err := k.BindPort(ctx, genState.PortId); err != nil {
			panic("could not claim port capability: " + err.Error())
		}
	}

	if err := k.ParamStore.Set(ctx, genState.Params); err != nil {
		panic(err)
	}

	for _, keyRequest := range genState.KeyRequests {
		if err := k.KeyRequestStore.Set(ctx, keyRequest.Id, keyRequest); err != nil {
			panic(err)
		}
	}

	for _, key := range genState.Keys {
		if err := k.KeyStore.Set(ctx, key.Id, key); err != nil {
			panic(err)
		}
	}

	largerCount := max(len(genState.Keys), len(genState.KeyRequests))

	if err := k.KeyRequestCount.Set(ctx, uint64(largerCount)); err != nil {
		panic(err)
	}

	for _, signRequest := range genState.SignRequests {
		if err := k.SignRequestStore.Set(ctx, signRequest.Id, signRequest); err != nil {
			panic(err)
		}
	}

	if err := k.SignRequestCount.Set(ctx, uint64(len(genState.SignRequests))); err != nil {
		panic(err)
	}

	for _, signTxRequest := range genState.SignTxRequests {
		if err := k.SignTransactionRequestStore.Set(ctx, signTxRequest.Id, signTxRequest); err != nil {
			panic(err)
		}
	}

	if err := k.SignTransactionRequestCount.Set(ctx, uint64(len(genState.SignTxRequests))); err != nil {
		panic(err)
	}

	for _, icaTxRequest := range genState.IcaTxRequests {
		if err := k.ICATransactionRequestStore.Set(ctx, icaTxRequest.Id, icaTxRequest); err != nil {
			panic(err)
		}
	}

	if err := k.ICATransactionRequestCount.Set(ctx, uint64(len(genState.IcaTxRequests))); err != nil {
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
