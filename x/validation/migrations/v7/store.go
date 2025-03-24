package v7

import (
	"cosmossdk.io/collections"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func UpdateParams(ctx sdk.Context, params collections.Item[types.HVParams]) error {
	if err := params.Set(ctx, *types.DefaultHVParams(ctx)); err != nil {
		return err
	}
	return nil
}
