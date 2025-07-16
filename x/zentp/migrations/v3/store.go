package v3

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	appparams "github.com/Zenrock-Foundation/zrchain/v6/app/params"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
)

func UpdateMintStore(ctx sdk.Context, oldMintsCol collections.Map[uint64, types.Bridge], newMintsCol collections.Map[uint64, types.Bridge]) error {
	oldMintsCol.Walk(ctx, nil, func(key uint64, mint types.Bridge) (bool, error) {
		if err := newMintsCol.Set(ctx, key, mint); err != nil {
			return true, err
		}
		return false, nil
	})

	return nil
}

func UpdateBurnStore(ctx sdk.Context, oldBurnsCol collections.Map[uint64, types.Bridge], newBurnsCol collections.Map[uint64, types.Bridge]) error {
	oldBurnsCol.Walk(ctx, nil, func(key uint64, burn types.Bridge) (bool, error) {
		if err := newBurnsCol.Set(ctx, key, burn); err != nil {
			return true, err
		}
		return false, nil
	})

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

	var pendingTransferAmounts sdk.Coins

	for _, mint := range pendingMints {
		pendingTransferAmounts = pendingTransferAmounts.Add(sdk.NewCoin(mint.Denom, math.NewIntFromUint64(mint.Amount)))
	}

	zentpAddr := accountKeeper.GetModuleAddress(types.ModuleName)
	zentpBalance := bankKeeper.GetBalance(ctx, zentpAddr, appparams.BondDenom)

	amountToSend := zentpBalance.Sub(sdk.NewCoin(appparams.BondDenom, pendingTransferAmounts.AmountOf(appparams.BondDenom)))

	if !amountToSend.IsZero() {
		err = bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.ZentpCollectorName, sdk.NewCoins(amountToSend))
		if err != nil {
			return err
		}
	}

	return nil
}
