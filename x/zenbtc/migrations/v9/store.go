package v9

import (
	"strings"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

const (
	AmberChainIDPrefix = "amber"
)

// RemoveBurnedEvents removes all burn events with BURNED status
// and updates the first pending burn event pointer to the last non-BURNED event before the deleted ones.
// Only runs on chains with "amber" prefix (devnet only).
func RemoveBurnedEvents(
	ctx sdk.Context,
	burnEvents collections.Map[uint64, types.BurnEvent],
	firstPendingBurnEvent collections.Item[uint64],
) error {
	// Only run on amber (devnet) chains
	if !strings.HasPrefix(ctx.ChainID(), AmberChainIDPrefix) {
		return nil
	}

	var keysToRemove []uint64
	var minBurnedID *uint64
	var maxNonBurnedBeforeFirstBurned *uint64

	// First pass: find the minimum BURNED ID
	if err := burnEvents.Walk(ctx, nil, func(key uint64, value types.BurnEvent) (stop bool, err error) {
		if value.Status == types.BurnStatus_BURN_STATUS_BURNED {
			if minBurnedID == nil || key < *minBurnedID {
				minBurnedID = &key
			}
		}
		return false, nil
	}); err != nil {
		return err
	}

	// Second pass: collect BURNED events to remove and find the last non-BURNED event before the first BURNED one
	if err := burnEvents.Walk(ctx, nil, func(key uint64, value types.BurnEvent) (stop bool, err error) {
		if value.Status == types.BurnStatus_BURN_STATUS_BURNED {
			keysToRemove = append(keysToRemove, key)
		} else {
			// Track the maximum non-BURNED ID that comes before the first BURNED event
			if minBurnedID != nil && key < *minBurnedID {
				if maxNonBurnedBeforeFirstBurned == nil || key > *maxNonBurnedBeforeFirstBurned {
					maxNonBurnedBeforeFirstBurned = &key
				}
			}
		}
		return false, nil
	}); err != nil {
		return err
	}

	// Remove the identified burn events
	for _, key := range keysToRemove {
		if err := burnEvents.Remove(ctx, key); err != nil {
			return err
		}
	}

	// Update first pending burn event to the last non-BURNED event before the deleted ones
	if len(keysToRemove) > 0 && maxNonBurnedBeforeFirstBurned != nil {
		if err := firstPendingBurnEvent.Set(ctx, *maxNonBurnedBeforeFirstBurned); err != nil {
			return err
		}
	}

	return nil
}
