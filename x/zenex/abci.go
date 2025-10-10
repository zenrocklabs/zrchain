package zenex

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(goctx context.Context, k keeper.Keeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)

	ctx := sdk.UnwrapSDKContext(goctx)

	requiredRockBalance, err := k.GetRequiredRockBalance(ctx)
	if err != nil {
		k.Logger().Error("failed to get required rock balance", "error", err)
		return nil
	}

	rockFeePoolBalance := k.GetRockFeePoolBalance(ctx)
	if rockFeePoolBalance >= requiredRockBalance {
		err := k.CreateRockBtcSwap(ctx, rockFeePoolBalance)
		if err != nil {
			k.Logger().Error("failed to create rock btc swap", "error", err)
			return nil
		}
	}

	return nil
}
