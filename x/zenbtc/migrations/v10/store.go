package v10

import (
	"strings"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

const (
	AmberChainIDPrefix = "amber"
)

// RemoveInvalidRedemptions removes all redemptions with UNSTAKED or AWAITING_SIGN status
// and updates the first pending redemption pointer to the latest non-UNSTAKED/non-AWAITING_SIGN redemption.
// Only runs on chains with "amber" prefix (devnet only).
func RemoveInvalidRedemptions(
	ctx sdk.Context,
	redemptions collections.Map[uint64, types.Redemption],
	firstPendingRedemption collections.Item[uint64],
) error {
	// Only run on amber (devnet) chains
	if !strings.HasPrefix(ctx.ChainID(), AmberChainIDPrefix) {
		return nil
	}

	var keysToRemove []uint64
	var maxValidID *uint64

	// Find all redemptions with UNSTAKED or AWAITING_SIGN status
	// and track the maximum valid (non-UNSTAKED/non-AWAITING_SIGN) redemption ID
	if err := redemptions.Walk(ctx, nil, func(key uint64, value types.Redemption) (stop bool, err error) {
		if value.Status == types.RedemptionStatus_UNSTAKED || value.Status == types.RedemptionStatus_AWAITING_SIGN {
			keysToRemove = append(keysToRemove, key)
		} else {
			// Track the maximum ID of any non-invalid redemption
			if maxValidID == nil || key > *maxValidID {
				maxValidID = &key
			}
		}
		return false, nil
	}); err != nil {
		return err
	}

	// Remove the identified redemptions
	for _, key := range keysToRemove {
		if err := redemptions.Remove(ctx, key); err != nil {
			return err
		}
	}

	// Update first pending redemption to the latest valid redemption
	if len(keysToRemove) > 0 && maxValidID != nil {
		if err := firstPendingRedemption.Set(ctx, *maxValidID); err != nil {
			return err
		}
	}

	return nil
}
