package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"

	"github.com/Zenrock-Foundation/zrchain/v4/x/treasury/types"
)

func (k Keeper) KeyByID(goCtx context.Context, req *types.QueryKeyByIDRequest) (*types.QueryKeyByIDResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidArgument, "request is nil")
	}

	fmt.Println("foo", "checkpoint", "3b")

	key, err := k.KeyStore.Get(goCtx, req.Id)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "key %d not found", req.Id)
	}

	fmt.Println("foo", "checkpoint", "3c")

	return &types.QueryKeyByIDResponse{
		Key: &types.KeyResponse{
			Id:            key.Id,
			WorkspaceAddr: key.WorkspaceAddr,
			KeyringAddr:   key.KeyringAddr,
			Type:          key.Type.String(),
			PublicKey:     key.PublicKey,
			Index:         key.Index,
			SignPolicyId:  key.SignPolicyId,
		},
		Wallets: processWallets(key, req.WalletType, req.Prefixes),
	}, nil
}
