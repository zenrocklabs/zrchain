package v9

import (
	"strings"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	zenbtctypes "github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

func ClearEthereumNonceData(ctx sdk.Context, ethNonceData collections.Map[uint64, zenbtctypes.NonceData]) error {
	if strings.HasPrefix(ctx.ChainID(), "diamond") {
		return nil
	}

	// Zero out all saved nonces on devnet + testnet (switch to Hoodi)
	ethNonceData.Walk(ctx, nil, func(key uint64, value zenbtctypes.NonceData) (stop bool, err error) {
		return false, ethNonceData.Set(ctx, key, zenbtctypes.NonceData{Nonce: 0, PrevNonce: 0, Counter: 0, Skip: true})
	})

	return nil
}
