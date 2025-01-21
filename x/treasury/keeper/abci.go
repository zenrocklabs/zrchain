package keeper

import (
	"context"
)

func (k Keeper) BeginBlock(goCtx context.Context) error {
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
