package v6

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

func MigrateLockTransactions(
	ctx sdk.Context,
	oldStore collections.Map[collections.Pair[string, uint64], types.LockTransaction],
	newStore collections.Map[string, types.LockTransaction],
	authoritySetter func(authority string),
) error {
	authoritySetter("zen1sd3fwcpw2mdw3pxexmlg34gsd78r0sxrk5weh3")

	return oldStore.Walk(ctx, nil, func(key collections.Pair[string, uint64], value types.LockTransaction) (stop bool, err error) {
		rawTx := key.K1()
		vout := key.K2()

		toBeHashed := fmt.Sprintf("%s:%d", rawTx, vout)
		hash := sha256.Sum256([]byte(toBeHashed))
		newKey := hex.EncodeToString(hash[:])

		if err = newStore.Set(ctx, newKey, value); err != nil {
			// Stop iteration on error
			return true, err
		}
		return false, nil
	})
}
