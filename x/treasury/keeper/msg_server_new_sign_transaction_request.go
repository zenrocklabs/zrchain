package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	pol "github.com/Zenrock-Foundation/zrchain/v6/policy"
	policykeeper "github.com/Zenrock-Foundation/zrchain/v6/x/policy/keeper"
	policytypes "github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) NewSignTransactionRequest(goCtx context.Context, msg *types.MsgNewSignTransactionRequest) (*types.MsgNewSignTransactionRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if len(msg.KeyIds) == 0 {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "key ids cannot be empty")
	}

	if len(msg.UnsignedTransaction) == 0 {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "unsigned transaction cannot be empty")
	}

	var keys []types.Key
	for _, keyId := range msg.KeyIds {
		key, err := k.KeyStore.Get(ctx, keyId)
		if err != nil {
			return nil, fmt.Errorf("key %v not found", keyId)
		}
		keys = append(keys, key)
	}

	key := keys[0]
	signPolicyId := key.SignPolicyId

	ws, err := k.identityKeeper.GetWorkspace(ctx, key.WorkspaceAddr)
	if err != nil {
		return nil, fmt.Errorf("workspace %s not found", key.WorkspaceAddr)
	}
	if signPolicyId == 0 {
		signPolicyId = ws.SignPolicyId
	}

	keyring, err := k.identityKeeper.GetKeyring(ctx, key.KeyringAddr)
	if err != nil || !keyring.IsActive {
		return nil, fmt.Errorf("problem with keyring found: %v, IsActive: %v", err, keyring.IsActive)
	}

	// use wallet to parse unsigned transaction
	w, err := types.NewWallet(&key, msg.WalletType)
	if err != nil {
		return nil, err
	}

	parser, ok := w.(types.TxParser)
	if !ok {
		return nil, fmt.Errorf("wallet %v does not implement TxParser", msg.WalletType)
	}

	var meta types.Metadata
	if err := k.cdc.UnpackAny(msg.Metadata, &meta); err != nil {
		return nil, fmt.Errorf("failed to unpack metadata: %w", err)
	}
	tx, err := parser.ParseTx(msg.UnsignedTransaction, meta)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tx: %w", err)
	}

	ctx.Logger().Debug("parsed layer 1 tx", "wallet", w, "tx", tx)

	act, err := k.policyKeeper.AddAction(ctx, msg.Creator, msg, signPolicyId, msg.Btl, map[string][]byte{
		types.TxValueKey:        []byte(tx.Amount.String()),
		types.DataForSigningKey: tx.DataForSigning,
	}, ws.Owners)
	if err != nil {
		return nil, err
	}
	return k.NewSignTransactionRequestActionHandler(ctx, act)
}

func (k msgServer) NewSignTransactionRequestPolicyGenerator(ctx sdk.Context, msg *types.MsgNewSignTransactionRequest) (pol.Policy, error) {
	key, err := k.KeyStore.Get(ctx, msg.KeyIds[0])
	if err != nil {
		return nil, fmt.Errorf("key not found")
	}

	ws, err := k.identityKeeper.GetWorkspace(ctx, key.WorkspaceAddr)
	if err != nil {
		return nil, fmt.Errorf("workspace not found")
	}

	pol := ws.PolicyNewSignTransactionRequest()
	return pol, nil
}

func (k msgServer) NewSignTransactionRequestActionHandler(ctx sdk.Context, act *policytypes.Action) (*types.MsgNewSignTransactionRequestResponse, error) {
	return policykeeper.TryExecuteAction(
		k.policyKeeper,
		k.cdc,
		ctx,
		act,
		func(ctx sdk.Context, msg *types.MsgNewSignTransactionRequest) (*types.MsgNewSignTransactionRequestResponse, error) {
			dataForSigning := act.GetPolicyDataMap()[types.DataForSigningKey]
			return k.HandleSignTransactionRequest(ctx, msg, dataForSigning)
		},
	)
}
