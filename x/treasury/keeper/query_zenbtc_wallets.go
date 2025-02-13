package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ZenbtcWallets(
	goCtx context.Context,
	req *types.QueryZenbtcWalletsRequest,
) (*types.QueryZenbtcWalletsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	keys, pageRes, err := query.CollectionFilteredPaginate(
		goCtx,
		k.KeyStore,
		req.Pagination,
		func(key uint64, value types.Key) (bool, error) {

			if value.ZenbtcMetadata == nil {
				return false, nil
			}

			switch {
			case req.MintChainId != "" && value.ZenbtcMetadata.Caip2ChainId != req.MintChainId:
				return false, nil
			case req.ChainType != types.WalletType_WALLET_TYPE_UNSPECIFIED && value.ZenbtcMetadata.ChainType != req.ChainType:
				return false, nil
			case req.RecipientAddr != "" && value.ZenbtcMetadata.RecipientAddr != req.RecipientAddr:
				return false, nil
			case req.ReturnAddr != "" && value.ZenbtcMetadata.ReturnAddress != req.ReturnAddr:
				return false, nil
			}

			recipientAddressMatch := (req.RecipientAddr == "" || value.ZenbtcMetadata.RecipientAddr == req.RecipientAddr)
			chainIdMatch := (req.MintChainId == "" || value.ZenbtcMetadata.Caip2ChainId == req.MintChainId)
			returnAddrMatch := (req.ReturnAddr == "" || value.ZenbtcMetadata.ReturnAddress == req.ReturnAddr)

			return recipientAddressMatch && chainIdMatch && returnAddrMatch, nil
		},
		func(key uint64, value types.Key) (*types.KeyAndWalletResponse, error) {
			return &types.KeyAndWalletResponse{
				Key: &types.KeyResponse{
					Id:             value.Id,
					WorkspaceAddr:  value.WorkspaceAddr,
					KeyringAddr:    value.KeyringAddr,
					Type:           value.Type.String(),
					PublicKey:      value.PublicKey,
					Index:          value.Index,
					SignPolicyId:   value.SignPolicyId,
					ZenbtcMetadata: value.ZenbtcMetadata,
				},
				Wallets: processWallets(value, types.WalletType_WALLET_TYPE_UNSPECIFIED, nil),
			}, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.QueryZenbtcWalletsResponse{
		ZenbtcWallets: keys,
		Pagination:    pageRes,
	}, nil
}
