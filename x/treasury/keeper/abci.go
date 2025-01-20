package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlock(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	fmt.Printf("block height: %d\n", ctx.BlockHeight())
	err := k.AddBlockTime(goCtx)
	if err != nil {
		return err
	}

	err = k.CheckForKeyMPCTimeouts(goCtx)
	if err != nil {
		return err
	}

	err = k.CheckForSignatureMPCTimeouts(goCtx)
	if err != nil {
		return err
	}

	return nil
}
