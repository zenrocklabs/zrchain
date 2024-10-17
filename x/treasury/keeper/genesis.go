package keeper

import (
	"github.com/Zenrock-Foundation/zrchain/v4/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) ExportState(ctx sdk.Context, genState *types.GenesisState) error {
	keyReqStore, err := k.KeyRequestStore.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	genState.KeyRequests, err = keyReqStore.Values()
	if err != nil {
		return err
	}

	keyStore, err := k.KeyStore.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	genState.Keys, err = keyStore.Values()
	if err != nil {
		return err
	}

	sigReqStore, err := k.SignRequestStore.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	genState.SignRequests, err = sigReqStore.Values()
	if err != nil {
		return err
	}

	sigTxReqStore, err := k.SignTransactionRequestStore.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	genState.SignTxRequests, err = sigTxReqStore.Values()
	if err != nil {
		return err
	}

	icaTxReqStore, err := k.ICATransactionRequestStore.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	genState.IcaTxRequests, err = icaTxReqStore.Values()
	if err != nil {
		return err
	}

	return nil
}
