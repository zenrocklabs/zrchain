package keeper

import (
	"github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) ExportState(ctx sdk.Context, genState *types.GenesisState) error {
	pStore, err := k.PolicyStore.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	genState.Policies, err = pStore.Values()
	if err != nil {
		return err
	}

	krStore, err := k.ActionStore.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	genState.Actions, err = krStore.Values()
	if err != nil {
		return err
	}

	smStore, err := k.SignMethodStore.Iterate(ctx, nil)
	if err != nil {
		return err
	}

	genState.SignMethods = []types.GenesisSignMethod{}
	kvs, err := smStore.KeyValues()
	if err != nil {
		return err
	}
	for _, kv := range kvs {
		genState.SignMethods = append(genState.SignMethods, types.GenesisSignMethod{
			Owner:  kv.Key.K1(),
			Id:     kv.Key.K2(),
			Config: &kv.Value,
		})
	}

	return nil
}
