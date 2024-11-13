package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ZrSignKeys(goCtx context.Context, req *types.QueryZrSignKeysRequest) (*types.QueryZrSignKeysResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	workspaces, err := k.identityKeeper.GetZrSignWorkspaces(goCtx, req.Address, req.WalletType)
	if err != nil {
		return &types.QueryZrSignKeysResponse{}, nil
	}

	result := &types.QueryZrSignKeysResponse{}
	keys, pageRes, err := query.CollectionFilteredPaginate[uint64, types.Key, collections.Map[uint64, types.Key], *types.ZrSignKeyEntry](
		goCtx,
		k.KeyStore,
		nil,
		func(key uint64, value types.Key) (bool, error) {
			for _, w := range workspaces {
				if w == value.WorkspaceAddr {
					return true, nil
				}
			}
			return false, nil
		},
		func(key uint64, value types.Key) (*types.ZrSignKeyEntry, error) {
			var keyType string
			for kt, w := range workspaces {
				if w == value.WorkspaceAddr {
					keyType = kt
					break
				}
			}
			address, err := zrSignWallet(&value, keyType)
			if err != nil {
				return nil, err
			}
			return &types.ZrSignKeyEntry{
				WalletType: keyType,
				Index:      value.Index,
				Address:    address,
				Id:         value.Id,
			}, nil
		})
	if err != nil {
		return nil, err
	}

	result.Keys = append(result.Keys, keys...)
	result.Pagination = pageRes

	return result, nil
}

func zrSignWallet(key *types.Key, keyType string) (string, error) {
	switch keyType {
	case "60":
		return types.EthereumAddress(key)
	case "1":
		return types.BitcoinP2WPKH(key, &chaincfg.TestNet3Params)
	case "0":
		return types.BitcoinP2WPKH(key, &chaincfg.MainNetParams)
	case "501":
		return types.SolanaAddress(key)
	default:
		return "", fmt.Errorf("unknown keyType: %s", keyType)
	}
}
