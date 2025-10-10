package v8

import (
	"strings"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

const (
	AmberChainIDPrefix = "amber"
)

// RemoveStakedMints removes all Solana pending mint transactions with STAKED status
// and updates the first pending Solana mint transaction index.
// Only runs on chains with "amber" prefix (devnet only).
func RemoveStakedMints(
	ctx sdk.Context,
	pendingMintTransactionsMap collections.Map[uint64, types.PendingMintTransaction],
	firstPendingSolMintTransaction collections.Item[uint64],
) error {
	// Only run on amber (devnet) chains
	if !strings.HasPrefix(ctx.ChainID(), AmberChainIDPrefix) {
		return nil
	}

	var keysToRemove []uint64
	var minSolID *uint64

	// Find all Solana pending mint transactions with STAKED status
	if err := pendingMintTransactionsMap.Walk(ctx, nil, func(key uint64, value types.PendingMintTransaction) (stop bool, err error) {
		if value.ChainType == types.WalletType_WALLET_TYPE_SOLANA {
			if value.Status == types.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED {
				keysToRemove = append(keysToRemove, key)
			} else {
				// Track the minimum ID for non-STAKED Solana transactions
				if minSolID == nil || key < *minSolID {
					minSolID = &key
				}
			}
		}
		return false, nil
	}); err != nil {
		return err
	}

	// Remove the identified pending mint transactions
	for _, key := range keysToRemove {
		if err := pendingMintTransactionsMap.Remove(ctx, key); err != nil {
			return err
		}
	}

	// Update first pending Solana mint transaction if we have a new minimum
	if minSolID != nil {
		if err := firstPendingSolMintTransaction.Set(ctx, *minSolID); err != nil {
			return err
		}
	}

	return nil
}
