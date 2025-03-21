package keeper

import (
	"cosmossdk.io/collections"
	policytypes "github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetSignMethod(ctx sdk.Context, owner string, id string, signMethod policytypes.SignMethod) error {
	config, err := codectypes.NewAnyWithValue(signMethod)
	if err != nil {
		return err
	}

	if err := k.SignMethodStore.Set(ctx, collections.Join(owner, id), *config); err != nil {
		return err
	}
	return nil
}

func (k Keeper) GetSignMethod(ctx sdk.Context, owner string, id string) (policytypes.SignMethod, error) {
	data, err := k.SignMethodStore.Get(ctx, collections.Join(owner, id))
	if err != nil {
		return nil, err
	}

	var signMethod policytypes.SignMethod
	if err := k.cdc.UnpackAny(&data, &signMethod); err != nil {
		return nil, err
	}
	return signMethod, nil
}
