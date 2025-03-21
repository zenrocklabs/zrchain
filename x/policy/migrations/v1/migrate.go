package v1

import (
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
)

func UpdateParams(ctx sdk.Context, params collections.Item[types.Params]) error {
	oldParams, err := params.Get(ctx)
	if err != nil {
		return err
	}

	// ...

	params.Set(ctx, oldParams)

	return nil
}

// ...
