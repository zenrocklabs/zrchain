package keeper

import (
	"context"
	"fmt"
	"slices"

	errorsmod "cosmossdk.io/errors"
	pol "github.com/Zenrock-Foundation/zrchain/v4/policy"
	identitytypes "github.com/Zenrock-Foundation/zrchain/v4/x/identity/types"
	policykeeper "github.com/Zenrock-Foundation/zrchain/v4/x/policy/keeper"
	policytypes "github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
	"github.com/Zenrock-Foundation/zrchain/v4/x/treasury/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) NewKeyRequest(goCtx context.Context, msg *types.MsgNewKeyRequest) (*types.MsgNewKeyRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	store, err := k.ParamStore.Get(ctx)
	if err != nil {
		return nil, err
	}
	if msg.Creator == store.ZrSignAddress {
		return k.zrSignKeyRequest(goCtx, msg)
	}

	wsbz, err := sdk.GetFromBech32(msg.WorkspaceAddr, identitytypes.PrefixWorkspaceAddress)
	if err != nil || len(wsbz) != identitytypes.WorkspaceAddressLength {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid workspace address (%s)", err)
	}

	krbz, err := sdk.GetFromBech32(msg.KeyringAddr, identitytypes.PrefixKeyringAddress)
	if err != nil || len(krbz) != identitytypes.KeyringAddressLength {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid keyring address (%s)", err)
	}

	if !slices.Contains(types.ValidKeyTypes, msg.KeyType) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid keytype %s, valid types %+v", msg.KeyType, types.ValidKeyTypes)
	}
	ws, err := k.identityKeeper.WorkspaceStore.Get(goCtx, msg.WorkspaceAddr)
	if err != nil {
		return nil, fmt.Errorf("workspace %s not found", msg.WorkspaceAddr)
	}

	if msg.SignPolicyId > 0 {
		if err := k.policyKeeper.PolicyMembersAreOwners(goCtx, msg.SignPolicyId, ws.Owners); err != nil {
			return nil, err
		}
	}

	// we have to check if the keyring is Active or not
	keyring, err := k.identityKeeper.KeyringStore.Get(goCtx, msg.KeyringAddr)
	if err != nil || !keyring.IsActive {
		return nil, fmt.Errorf("keyring %s is nil or is inactive", msg.KeyringAddr)
	}

	k.policyKeeper.PolicyStore.Get(ctx, ws.SignPolicyId)
	act, err := k.policyKeeper.AddAction(ctx, msg.Creator, msg, ws.SignPolicyId, msg.Btl, nil)
	if err != nil {
		return nil, err
	}

	return k.NewKeyRequestActionHandler(ctx, act)
}

func (k msgServer) NewKeyRequestPolicyGenerator(ctx sdk.Context, msg *types.MsgNewKeyRequest) (pol.Policy, error) {
	ws, err := k.identityKeeper.WorkspaceStore.Get(ctx, msg.WorkspaceAddr)
	if err != nil {
		return nil, fmt.Errorf("workspace not found")
	}

	pol := ws.PolicyNewKeyRequest()
	return pol, nil
}

func (k msgServer) NewKeyRequestActionHandler(ctx sdk.Context, act *policytypes.Action) (*types.MsgNewKeyRequestResponse, error) {
	return policykeeper.TryExecuteAction(
		&k.policyKeeper,
		k.cdc,
		ctx,
		act,
		k.newKeyRequest,
	)
}
