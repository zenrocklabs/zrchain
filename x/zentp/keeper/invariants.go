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

	pendingMints, err := k.GetMintsWithStatus(ctx, types.BridgeStatus_BRIDGE_STATUS_PENDING)
	if err != nil {
		// It's ok if there are no pending mints, so we don't return an error.
	}

	var pendingAmount math.Int = math.ZeroInt()
	for _, mint := range pendingMints {
		pendingAmount = pendingAmount.Add(math.NewIntFromUint64(mint.Amount))
	}

	zrchainSupply := k.bankKeeper.GetSupply(ctx, params.BondDenom).Amount

	totalSupply := zrchainSupply.Add(solanaSupply).Add(pendingAmount)
	if newAmount.IsPositive() {
		totalSupply = totalSupply.Add(newAmount)
	}

	if totalSupply.GT(sdkmath.NewIntFromUint64(rockCap)) {
		return errors.Errorf("total ROCK supply including pending would exceed cap (%s), bridge disabled", totalSupply.String())
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
