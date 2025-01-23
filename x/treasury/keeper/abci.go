package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlock(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	fmt.Println(ctx.BlockHeight())

	err := k.CheckForKeyMPCTimeouts(goCtx)
	if err != nil {
		return err
	}

	err = k.CheckForSignatureMPCTimeouts(goCtx)
	if err != nil {
		return err
	}

	return nil
}
