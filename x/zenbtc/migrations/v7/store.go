package v7

import (
	"strings"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

const (
	GardiaChainIDPrefix = "gardia"
	SolanaChainIDPrefix = "solana:"
)

// PurgeInvalidState removes all burn events with a "solana:" prefix, and unstaked redemptions on gardia chains.
func PurgeInvalidState(
	ctx sdk.Context,
	burnEvents collections.Map[uint64, types.BurnEvent],
	firstPendingBurnEvent collections.Item[uint64],
	redemptions collections.Map[uint64, types.Redemption],
	firstPendingRedemption collections.Item[uint64],
) error {
	var keysToRemove []uint64
	var minID *uint64

	// Find all burn events with the target chain ID
	if err := burnEvents.Walk(ctx, nil, func(key uint64, value types.BurnEvent) (stop bool, err error) {
		if strings.HasPrefix(value.ChainID, SolanaChainIDPrefix) {
			keysToRemove = append(keysToRemove, key)
			if minID == nil || key < *minID {
				minID = &key
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

	// Set the new first pending burn event if any were removed
	if minID != nil {
		if err := firstPendingBurnEvent.Set(ctx, *minID); err != nil {
			return err
		}
	}

	if strings.HasPrefix(ctx.ChainID(), GardiaChainIDPrefix) {
		var redemptionsToRemove []uint64
		var minRedemptionID *uint64
		if err := redemptions.Walk(ctx, nil, func(key uint64, value types.Redemption) (stop bool, err error) {
			if value.Status == types.RedemptionStatus_UNSTAKED {
				redemptionsToRemove = append(redemptionsToRemove, key)
				if minRedemptionID == nil || key < *minRedemptionID {
					minRedemptionID = &key
				}
			}
			return false, nil
		}); err != nil {
			return err
		}

		for _, key := range redemptionsToRemove {
			if err := redemptions.Remove(ctx, key); err != nil {
				return err
			}
		}

		if minRedemptionID != nil {
			if err := firstPendingRedemption.Set(ctx, *minRedemptionID); err != nil {
				return err
			}
		}
	}

	return nil
}
