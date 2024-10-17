package keeper

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"fmt"
	"strconv"

	identitytypes "github.com/Zenrock-Foundation/zrchain/v4/x/identity/types"
	"github.com/Zenrock-Foundation/zrchain/v4/x/treasury/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) FulfilKeyRequest(goCtx context.Context, msg *types.MsgFulfilKeyRequest) (*types.MsgFulfilKeyRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	req, err := k.KeyRequestStore.Get(ctx, msg.RequestId)
	if err != nil {
		return nil, fmt.Errorf("request %v not found", msg.RequestId)
	}

	keyring, err := k.identityKeeper.KeyringStore.Get(ctx, req.KeyringAddr)
	if err != nil || !keyring.IsActive {
		return nil, fmt.Errorf("keyring %s is nil or is inactive", req.KeyringAddr)
	}

	if err := k.validateKeyRequest(msg, &req, &keyring); err != nil {
		return nil, err
	}

	switch msg.Status {
	case types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED:
		return k.handleKeyRequestFulfilment(ctx, msg, &req)
	case types.KeyRequestStatus_KEY_REQUEST_STATUS_REJECTED:
		return k.handleKeyRequestRejection(ctx, msg, &req)
	default:
		return nil, fmt.Errorf("invalid status field, should be either fulfilled/rejected")
	}
}

func (k msgServer) validateKeyRequest(msg *types.MsgFulfilKeyRequest, req *types.KeyRequest, keyring *identitytypes.Keyring) error {
	if !keyring.IsParty(msg.Creator) {
		return fmt.Errorf("only one party of the keyring can fulfil key request")
	}

	if req.Status != types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING && req.Status != types.KeyRequestStatus_KEY_REQUEST_STATUS_PARTIAL {
		return fmt.Errorf("request %v is not pending/partial, can't be updated", req.Status)
	}

	if len(msg.KeyringPartySignature) != 64 {
		return fmt.Errorf("invalid mpc party signature, should be 64 bytes, is %d", len(msg.KeyringPartySignature))
	}

	return nil
}

func (k msgServer) handleKeyRequestFulfilment(ctx sdk.Context, msg *types.MsgFulfilKeyRequest, req *types.KeyRequest) (*types.MsgFulfilKeyRequestResponse, error) {
	pubKey := (msg.Result.(*types.MsgFulfilKeyRequest_Key)).Key.PublicKey

	if err := k.validatePublicKey(req.KeyType, pubKey); err != nil {
		return nil, err
	}

	sigExists := false
	for _, sig := range req.KeyringPartySignatures {
		if bytes.Equal(sig, msg.KeyringPartySignature) {
			sigExists = true
			break
		}
	}

	if !sigExists {
		req.KeyringPartySignatures = append(req.KeyringPartySignatures, msg.KeyringPartySignature)
	}

	keyring, err := k.identityKeeper.KeyringStore.Get(ctx, req.KeyringAddr)
	if err != nil || !keyring.IsActive {
		return nil, fmt.Errorf("keyring %s is nil or is inactive", req.KeyringAddr)
	}

	if len(req.KeyringPartySignatures) >= int(keyring.PartyThreshold) {
		req.Status = types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED

		key := &types.Key{
			Id:             req.Id,
			WorkspaceAddr:  req.WorkspaceAddr,
			KeyringAddr:    req.KeyringAddr,
			Type:           req.KeyType,
			PublicKey:      pubKey,
			Index:          req.Index,
			SignPolicyId:   req.SignPolicyId,
			ZenbtcMetadata: req.ZenbtcMetadata,
		}

		if err := k.KeyStore.Set(ctx, key.Id, *key); err != nil {
			return nil, err
		}
	} else {
		req.Status = types.KeyRequestStatus_KEY_REQUEST_STATUS_PARTIAL
	}

	if err := k.KeyRequestStore.Set(ctx, req.Id, *req); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventKeyRequestFulfilled,
			sdk.NewAttribute(types.AttributeRequestId, strconv.FormatUint(req.GetId(), 10)),
		),
	})

	return &types.MsgFulfilKeyRequestResponse{}, nil
}

func (k msgServer) handleKeyRequestRejection(ctx sdk.Context, msg *types.MsgFulfilKeyRequest, req *types.KeyRequest) (*types.MsgFulfilKeyRequestResponse, error) {
	req.Status = types.KeyRequestStatus_KEY_REQUEST_STATUS_REJECTED
	req.RejectReason = msg.Result.(*types.MsgFulfilKeyRequest_RejectReason).RejectReason

	if err := k.KeyRequestStore.Set(ctx, req.Id, *req); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventKeyRequestRejected,
			sdk.NewAttribute(types.AttributeRequestId, strconv.FormatUint(req.GetId(), 10)),
		),
	})

	return &types.MsgFulfilKeyRequestResponse{}, nil
}

func (k msgServer) validatePublicKey(keyType types.KeyType, pubKey []byte) error {
	switch keyType {
	case types.KeyType_KEY_TYPE_ECDSA_SECP256K1:
		if l := len(pubKey); l != 33 && l != 65 {
			return fmt.Errorf("invalid ecdsa public key %x of length %v", pubKey, l)
		}
	case types.KeyType_KEY_TYPE_EDDSA_ED25519:
		if l := len(pubKey); l != ed25519.PublicKeySize {
			return fmt.Errorf("invalid eddsa public key %x of length %v", pubKey, l)
		}
	case types.KeyType_KEY_TYPE_BITCOIN_SECP256K1:
		if l := len(pubKey); l != 33 && l != 65 {
			return fmt.Errorf("invalid bitcoin public key %x of length %v", pubKey, l)
		}
	default:
		return fmt.Errorf("invalid key type: %v", keyType.String())
	}
	return nil
}
