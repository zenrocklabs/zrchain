package keeper

import (
	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

const rockCap = 1_000_000_000_000_000 // 1bn ROCK in urock

// CheckROCKSupplyCap checks if the total ROCK supply (on-chain + pending) plus any new amount would exceed the cap.
// A positive newAmount is for a new bridge-to-Solana request.
// A zero newAmount is for checking existing state, e.g. before completing a mint or burn.
func (k Keeper) CheckROCKSupplyCap(ctx sdk.Context, newAmount math.Int) error {
	solanaSupply, err := k.GetSolanaROCKSupply(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get solana rock supply")
	}

		zrchainSupply := k.bankKeeper.GetSupply(ctx, params.BondDenom).Amount

	if newAmount.IsPositive() {
		// This check is for new bridge requests from zrchain to solana.
		// It ensures the bridge amount does not exceed the supply available for bridging.
		zentpModuleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
		zentpModuleBalance := k.bankKeeper.GetBalance(ctx, zentpModuleAddr, params.BondDenom).Amount
		availableForBridging := zrchainSupply.Sub(zentpModuleBalance)
		if newAmount.GT(availableForBridging) {
			return errors.Errorf("bridge amount %s exceeds available zrchain rock supply for bridging %s", newAmount.String(), availableForBridging.String())
		}
	}

	// Total supply is the sum of ROCK on zrchain and ROCK on Solana.
	// The zrchainSupply from the bank keeper includes all tokens, even those
	// held in the module account for pending bridges. This is the correct total.
	totalSupply := zrchainSupply.Add(solanaSupply)

	// A bridge operation does not change the total supply, so we do not add newAmount here.
	// We just check if the current total supply is already over the cap.
	if totalSupply.GT(sdkmath.NewIntFromUint64(rockCap)) {
		return errors.Errorf("total ROCK supply %s exceeds cap (%d), bridge disabled", totalSupply.String(), rockCap)
	}
	return nil
}

// CheckCanBurnFromSolana checks if a burn from Solana is valid.
func (k Keeper) CheckCanBurnFromSolana(ctx sdk.Context, burnAmount math.Int) error {
	solanaSupply, err := k.GetSolanaROCKSupply(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get solana rock supply")
	}
	if burnAmount.GT(solanaSupply) {
		return errors.Errorf("attempt to bridge from solana exceeds solana ROCK supply, amount: %s, supply: %s", burnAmount.String(), solanaSupply.String())
	}
	// Also check cap on total supply. After burn, zrchain supply will increase. Total supply stays same.
	return k.CheckROCKSupplyCap(ctx, math.ZeroInt())
}
