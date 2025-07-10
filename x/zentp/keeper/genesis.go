package keeper

import (
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) ExportState(ctx sdk.Context, genState *types.GenesisState) error {
	mintStore, err := k.MintStore.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	genState.Mints, err = mintStore.Values()
	if err != nil {
		return err
	}

	burnStore, err := k.BurnStore.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	genState.Burns, err = burnStore.Values()
	if err != nil {
		return err
	}

	solanaRockSupply, err := k.SolanaROCKSupply.Get(ctx)
	if err != nil {
		return err
	}
	genState.SolanaRockSupply = solanaRockSupply.Uint64()

	zentpFees, err := k.ZentpFees.Get(ctx)
	if err != nil {
		return err
	}
	genState.ZentpFees = zentpFees

	return nil
}
