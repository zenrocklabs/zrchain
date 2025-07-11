package v3

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
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

func SendZentpFeesToMintModule(
	ctx sdk.Context,
	getPendingMints func(context.Context) ([]*types.Bridge, error),
	getBridgeFeeParams func(context.Context) (sdk.AccAddress, math.LegacyDec, error),
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
) error {
	pendingMints, err := getPendingMints(ctx)
	if err != nil {
		return err
	}

	var pendingFees sdk.Coins
	_, bridgeFee, err := getBridgeFeeParams(ctx)
	if err != nil {
		return err
	}

	for _, mint := range pendingMints {
		amountInt := math.NewIntFromUint64(mint.Amount)
		bridgeFeeAmount := math.LegacyNewDecFromInt(amountInt).Mul(bridgeFee).TruncateInt()
		if bridgeFeeAmount.IsPositive() {
			pendingFees = pendingFees.Add(sdk.NewCoin(mint.Denom, bridgeFeeAmount))
		}
	}

	zentpAddr := accountKeeper.GetModuleAddress(types.ModuleName)
	zentpBalance := bankKeeper.SpendableCoins(ctx, zentpAddr)

	amountToSend := zentpBalance.Sub(pendingFees...)

	if amountToSend.IsAllPositive() {
		err = bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.ZentpCollectorName, amountToSend)
		if err != nil {
			return err
		}
	}

	return nil
}
