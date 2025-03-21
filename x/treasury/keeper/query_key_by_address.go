package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) KeyByAddress(goCtx context.Context, req *types.QueryKeyByAddressRequest) (*types.QueryKeyByAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	keys := []types.Key{}
	if err := k.KeyStore.Walk(ctx, nil, func(key uint64, value types.Key) (bool, error) {
		if req.KeyringAddr != "" && value.KeyringAddr != req.KeyringAddr {
			return false, nil
		}

		if req.KeyType != types.KeyType_KEY_TYPE_UNSPECIFIED && value.Type != req.KeyType {
			return false, nil
		}

		keys = append(keys, value)
		return false, nil
	}); err != nil {
		return nil, err
	}

	for _, key := range keys {
		wallet, err := types.NewWallet(&key, req.WalletType)
		if err != nil {
			return nil, err
		}

		if req.Address == wallet.Address() {
			return &types.QueryKeyByAddressResponse{Response: &types.KeyAndWalletResponse{
				Key: &types.KeyResponse{
					Id:             key.Id,
					WorkspaceAddr:  key.WorkspaceAddr,
					KeyringAddr:    key.KeyringAddr,
					Type:           key.Type.String(),
					PublicKey:      key.PublicKey,
					Index:          key.Index,
					SignPolicyId:   key.SignPolicyId,
					ZenbtcMetadata: key.ZenbtcMetadata,
				},
				Wallets: processWallets(key, req.WalletType, req.Prefixes),
			}}, nil
		}
	}

	return &types.QueryKeyByAddressResponse{}, nil
}
