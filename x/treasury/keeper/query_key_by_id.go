package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
)

func (k Keeper) KeyByID(goCtx context.Context, req *types.QueryKeyByIDRequest) (*types.QueryKeyByIDResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidArgument, "request is nil")
	}

	key, err := k.KeyStore.Get(goCtx, req.Id)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "key %d not found", req.Id)
	}

	return &types.QueryKeyByIDResponse{
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
	}, nil
}
