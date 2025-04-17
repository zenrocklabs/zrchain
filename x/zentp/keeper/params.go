package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
)

func (k Keeper) GetSolanaParams(ctx context.Context) *types.Solana {
	params, err := k.ParamStore.Get(ctx)
	if err != nil ||
		params.Solana == nil ||
		params.Solana.SignerKeyId == 0 ||
		params.Solana.ProgramId == "" ||
		params.Solana.MintAddress == "" ||
		// params.Solana.MultisigKeyAddress == "" ||
		params.Solana.Btl == 0 ||
		params.Solana.FeeWallet == "" ||
		params.Solana.NonceAuthorityKey == 0 ||
		params.Solana.NonceAccountKey == 0 {

		return types.DefaultSolanaParams
	}
	return params.Solana
}
