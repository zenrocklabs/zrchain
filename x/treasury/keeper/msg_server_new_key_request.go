package keeper

import (
	"context"
	"fmt"
	"slices"

	errorsmod "cosmossdk.io/errors"
	pol "github.com/Zenrock-Foundation/zrchain/v5/policy"
	identitytypes "github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	policykeeper "github.com/Zenrock-Foundation/zrchain/v5/x/policy/keeper"
	policytypes "github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

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

	workspaceBytes, err := sdk.GetFromBech32(msg.WorkspaceAddr, identitytypes.PrefixWorkspaceAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid workspace address (%s)", err)
	}
	if len(workspaceBytes) != identitytypes.WorkspaceAddressLength {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "workspace address length %d is invalid for workspace %s, should be %d", len(workspaceBytes), msg.WorkspaceAddr, identitytypes.WorkspaceAddressLength)
	}

	keyringBytes, err := sdk.GetFromBech32(msg.KeyringAddr, identitytypes.PrefixKeyringAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid keyring address (%s)", err)
	}
	if len(keyringBytes) != identitytypes.KeyringAddressLength {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "keyring address length %d is invalid for keyring %s, should be %d", len(keyringBytes), msg.KeyringAddr, identitytypes.KeyringAddressLength)
	}

	if !slices.Contains(types.ValidKeyTypes, msg.KeyType) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid keytype %s, valid types %+v", msg.KeyType, types.ValidKeyTypes)
	}
	ws, err := k.identityKeeper.WorkspaceStore.Get(goCtx, msg.WorkspaceAddr)
	if err != nil {
		return nil, fmt.Errorf("workspace %s not found", msg.WorkspaceAddr)
	}

	if msg.SignPolicyId > 0 {
		if err = k.policyKeeper.PolicyMembersAreOwners(goCtx, msg.SignPolicyId, ws.Owners); err != nil {
			return nil, err
		}
	}

	// we have to check if the keyring is Active or not
	keyring, err := k.identityKeeper.KeyringStore.Get(goCtx, msg.KeyringAddr)
	if err != nil || !keyring.IsActive {
		return nil, fmt.Errorf("keyring %s is nil or is inactive", msg.KeyringAddr)
	}

	if metadata := msg.ZenbtcMetadata; metadata != nil {
		// TODO: add other chainIDs before zenBTC mainnet upgrade
		if metadata.Caip2ChainId != "eip155:17000" {
			return nil, fmt.Errorf("unsupported mint recipient chainID for zenBTC deposit key")
		}
		// TODO: add alternate check for Solana address format
		if !common.IsHexAddress(metadata.RecipientAddr) {
			return nil, fmt.Errorf("mint recipient address for zenBTC deposit key must be a valid Ethereum address")
		}
		if metadata.ChainType == types.WalletType_WALLET_TYPE_UNSPECIFIED {
			return nil, fmt.Errorf("unsupported chain type for zenBTC deposit key")
		}
	}

	// TODO: do we want to have this check below?
	// if _, err := k.policyKeeper.PolicyStore.Get(ctx, ws.SignPolicyId); err != nil {
	// 	return nil, err
	// }

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
