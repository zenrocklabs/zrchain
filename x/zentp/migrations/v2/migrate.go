package v2

import (
	"cosmossdk.io/collections"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func UpdateParams(
	ctx sdk.Context,
	params collections.Item[types.Params],
) error {
	p, err := params.Get(ctx)
	if err != nil {
		return err
	}

	p.Solana.NonceAccountKey = 10
	p.Solana.NonceAuthorityKey = 11
	p.Solana.SignerKeyId = 12

	return params.Set(ctx, p)
}
