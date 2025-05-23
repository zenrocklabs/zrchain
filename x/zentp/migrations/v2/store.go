package v2

import (
	"fmt"
	"strings"

	"cosmossdk.io/collections"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func UpdateParams(ctx sdk.Context, params collections.Item[types.Params]) error {
	paramsMap := map[string]types.Solana{
		"zenrock": { // local
			SignerKeyId:       10,
			ProgramId:         "DXREJumiQhNejXa1b5EFPUxtSYdyJXBdiHeu6uX1ribA",
			NonceAccountKey:   12,
			NonceAuthorityKey: 11,
			MintAddress:       "StVNdHNSFK3uVTL5apWHysgze4M8zrsqwjEAH1JM87i",
			FeeWallet:         "FzqGcRG98v1KhKxatX2Abb2z1aJ2rViQwBK5GHByKCAd",
			Fee:               0,
			Btl:               20,
		},
		"amber": { // devnet
			// TODO: Configure devnet environment parameters
			SignerKeyId:       10,
			ProgramId:         "DXREJumiQhNejXa1b5EFPUxtSYdyJXBdiHeu6uX1ribA",
			NonceAccountKey:   12,
			NonceAuthorityKey: 11,
			MintAddress:       "StVNdHNSFK3uVTL5apWHysgze4M8zrsqwjEAH1JM87i",
			FeeWallet:         "FzqGcRG98v1KhKxatX2Abb2z1aJ2rViQwBK5GHByKCAd",
			Fee:               0,
			Btl:               20,
		},
		"gardia": { // testnet
			// TODO: Configure testnet environment parameters
			SignerKeyId:       10,
			ProgramId:         "DXREJumiQhNejXa1b5EFPUxtSYdyJXBdiHeu6uX1ribA",
			NonceAccountKey:   12,
			NonceAuthorityKey: 11,
			MintAddress:       "StVNdHNSFK3uVTL5apWHysgze4M8zrsqwjEAH1JM87i",
			FeeWallet:         "FzqGcRG98v1KhKxatX2Abb2z1aJ2rViQwBK5GHByKCAd",
			Fee:               0,
			Btl:               20,
		},
		"diamond": { // mainnet
			// TODO: Configure mainnet environment parameters
			SignerKeyId:       0,
			ProgramId:         "3WyacwnCNiz4Q1PedWyuwodYpLFu75jrhgRTZp69UcA9",
			NonceAccountKey:   0,
			NonceAuthorityKey: 0,
			MintAddress:       "5VsPJ2EG7jjo3k2LPzQVriENKKQkNUTzujEzuaj4Aisf",
			FeeWallet:         "7AnbfuYgwXXKo2Jn8HBc9HBrrNezEPvYA55NW2PWmSHQ",
			Fee:               0,
			Btl:               20,
		},
	}

	chainID := ctx.ChainID()
	if chainID == "" {
		chainID = "zenrock"
	}

	newParams := types.Params{
		Solana:    &types.Solana{},
		BridgeFee: types.DefaultParams().BridgeFee,
	}

	for prefix, paramSet := range paramsMap {
		if strings.HasPrefix(chainID, prefix) {
			newParams.Solana = &paramSet
			break
		}
	}

	if newParams.Solana == nil || newParams.BridgeFee.IsNil() || newParams.BridgeFee.IsNegative() {
		return fmt.Errorf("failed to update params for chain %s", chainID)
	}

	if err := params.Set(ctx, newParams); err != nil {
		return err
	}

	return nil
}
