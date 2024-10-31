package keeper

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"

	errorsmod "cosmossdk.io/errors"

	"cosmossdk.io/collections"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CreateWorkspace(ctx sdk.Context, workspace *types.Workspace) (string, error) {
	count, err := k.WorkspaceCount.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return "", err
		}
		count = 1
	} else {
		count++
	}

	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, count)
	addrHash := sha256.Sum256(buf)
	workspace.Address = sdk.MustBech32ifyAddressBytes(types.PrefixWorkspaceAddress, sdk.AccAddress(addrHash[:types.WorkspaceAddressLength]))

	if err := k.WorkspaceStore.Set(ctx, workspace.Address, *workspace); err != nil {
		return "", err
	}

	if err := k.WorkspaceCount.Set(ctx, count); err != nil {
		return "", err
	}

	return workspace.Address, nil
}

func (k *Keeper) storeChildWorkspace(ctx sdk.Context, parent, child *types.Workspace) (*types.MsgNewChildWorkspaceResponse, error) {
	childAddr, err := k.CreateWorkspace(ctx, child)
	if err != nil {
		return nil, errorsmod.Wrap(err, "create child workspace failed")
	}

	parent.AddChild(child)

	if err = k.WorkspaceStore.Set(ctx, parent.Address, *parent); err != nil {
		return nil, errorsmod.Wrapf(types.ErrInternal, "failed to set parent workspace %s", parent.Address)
	}

	return &types.MsgNewChildWorkspaceResponse{Address: childAddr}, nil
}
