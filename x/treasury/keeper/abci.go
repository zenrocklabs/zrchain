package keeper

import (
	"context"
)

func (k Keeper) BeginBlock(goCtx context.Context) error {
	if err := k.CheckForKeyMPCTimeouts(goCtx); err != nil {
		k.Logger().Error("error in BeginBlock: CheckForKeyMPCTimeouts", "error", err)
	}

	if err := k.CheckForSignatureMPCTimeouts(goCtx); err != nil {
		k.Logger().Error("error in BeginBlock: CheckForSignatureMPCTimeouts", "error", err)
	}

	return nil
}
