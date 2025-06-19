package v3

import (
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
)

func UpdateMintStore(ctx sdk.Context, oldMintsCol collections.Map[uint64, types.Bridge], newMintsCol collections.Map[uint64, types.Bridge]) error {
	mintStore, err := oldMintsCol.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	mints, err := mintStore.Values()
	if err != nil {
		return err
	}

	for _, mint := range mints {
		if err := newMintsCol.Set(ctx, mint.Id, mint); err != nil {
			return err
		}
	}

	return nil
}

func UpdateBurnStore(ctx sdk.Context, oldBurnsCol collections.Map[uint64, types.Bridge], newBurnsCol collections.Map[uint64, types.Bridge]) error {
	burnStore, err := oldBurnsCol.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	burns, err := burnStore.Values()
	if err != nil {
		return err
	}

	for _, burn := range burns {
		if err := newBurnsCol.Set(ctx, burn.Id, burn); err != nil {
			return err
		}
	}

	return nil
}

// ...
