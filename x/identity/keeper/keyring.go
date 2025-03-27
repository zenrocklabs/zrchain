package keeper

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
)

func (k Keeper) CreateKeyring(ctx sdk.Context, keyring *types.Keyring) (string, error) {
	count, err := k.KeyringCount.Get(ctx)
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
	keyring.Address = sdk.MustBech32ifyAddressBytes(types.PrefixKeyringAddress, sdk.AccAddress(addrHash[13:13+types.KeyringAddressLength]))

	if err := k.KeyringStore.Set(ctx, keyring.Address, *keyring); err != nil {
		return "", errorsmod.Wrapf(types.ErrInternal, "failed to set keyring %s", keyring.Address)
	}

	if err := k.KeyringCount.Set(ctx, count); err != nil {
		return "", errorsmod.Wrapf(types.ErrInternal, "failed to set keyring count")
	}

	return keyring.Address, nil
}
