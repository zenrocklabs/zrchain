package keeper

import (
	"context"

	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DctWallets(
	goCtx context.Context,
	req *types.QueryDctWalletsRequest,
) (*types.QueryDctWalletsResponse, error) {
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
			case req.AssetType != dcttypes.Asset_ASSET_UNSPECIFIED && value.ZenbtcMetadata.Asset != req.AssetType:
				return false, nil
			case req.MintChainId != "" && value.ZenbtcMetadata.Caip2ChainId != req.MintChainId:
				return false, nil
			case req.WalletType != types.WalletType_WALLET_TYPE_UNSPECIFIED && value.ZenbtcMetadata.ChainType != req.WalletType:
				return false, nil
			case req.RecipientAddr != "" && value.ZenbtcMetadata.RecipientAddr != req.RecipientAddr:
				return false, nil
			case req.RecipientAddr != "" && value.ZenbtcMetadata.RecipientAddr != req.RecipientAddr:
				return false, nil
			}

			assetMatch := (value.ZenbtcMetadata.Asset == req.AssetType)
			recipientAddressMatch := (req.RecipientAddr == "" || value.ZenbtcMetadata.RecipientAddr == req.RecipientAddr)
			chainIdMatch := (req.MintChainId == "" || value.ZenbtcMetadata.Caip2ChainId == req.MintChainId)
			returnAddrMatch := (req.RecipientAddr == "" || value.ZenbtcMetadata.RecipientAddr == req.RecipientAddr)

			return assetMatch && recipientAddressMatch && chainIdMatch && returnAddrMatch, nil
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

	return &types.QueryDctWalletsResponse{
		DctWallets: keys,
		Pagination: pageRes,
	}, nil
}
