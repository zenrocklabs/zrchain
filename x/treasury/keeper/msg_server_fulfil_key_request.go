package keeper

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	identitytypes "github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) FulfilKeyRequest(goCtx context.Context, msg *types.MsgFulfilKeyRequest) (*types.MsgFulfilKeyRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	req, err := k.KeyRequestStore.Get(ctx, msg.RequestId)
	if err != nil {
		return nil, fmt.Errorf("request %v not found", msg.RequestId)
	}

	keyring, err := k.identityKeeper.GetKeyring(ctx, req.KeyringAddr)
	if err != nil || !keyring.IsActive {
		return k.rejectKeyRequest(ctx, &req, fmt.Sprintf("keyring %s is nil or is inactive", req.KeyringAddr))
	}

	if len(msg.KeyringPartySignature) != 64 {
		return k.rejectKeyRequest(ctx, &req, fmt.Sprintf("invalid mpc party signature, should be 64 bytes, is %d", len(msg.KeyringPartySignature)))
	}

	if err := k.validateKeyRequest(msg, &req, keyring); err != nil {
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

	return nil
}

func (k msgServer) handleKeyRequestFulfilment(ctx sdk.Context, msg *types.MsgFulfilKeyRequest, req *types.KeyRequest) (*types.MsgFulfilKeyRequestResponse, error) {
	// Reject if a party tries to sign more than once
	for _, sig := range req.KeyringPartySignatures {
		if sig.Creator == msg.Creator {
			req.Status = types.KeyRequestStatus_KEY_REQUEST_STATUS_REJECTED
			errMsg := fmt.Sprintf("party %v already sent a fulfilment", msg.Creator)
			req.RejectReason = errMsg
			if err := k.KeyRequestStore.Set(ctx, req.Id, *req); err != nil {
				return nil, err
			}
			return &types.MsgFulfilKeyRequestResponse{}, nil
		}
	}

	// Check against public key from other parties
	pubKey := (msg.Result.(*types.MsgFulfilKeyRequest_Key)).Key.PublicKey
	if len(req.PublicKey) > 0 {
		if !bytes.Equal(req.PublicKey, pubKey) {
			rejectReason := fmt.Sprintf("public key mismatch, expected %x, got %x", req.PublicKey, pubKey)
			return k.rejectKeyRequest(ctx, req, rejectReason)
		}
	}

	if res, err := k.validatePublicKey(ctx, req, req.KeyType, pubKey); err != nil || res != nil {
		return res, err
	}

	// Append party signature
	req.KeyringPartySignatures = append(req.KeyringPartySignatures, &types.PartySignature{
		Creator:   msg.Creator,
		Signature: msg.KeyringPartySignature,
	})

	keyring, err := k.identityKeeper.GetKeyring(ctx, req.KeyringAddr)
	if err != nil || !keyring.IsActive {
		return k.rejectKeyRequest(ctx, req, fmt.Sprintf("keyring %s is nil or is inactive", req.KeyringAddr))
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
		// Store public key from first party's response so we can check other parties respond with the same key
		req.PublicKey = pubKey
		req.Status = types.KeyRequestStatus_KEY_REQUEST_STATUS_PARTIAL
	}

	if err := k.KeyRequestStore.Set(ctx, req.Id, *req); err != nil {
		return nil, err
	}

	if req.Fee > 0 {
		feeRecipient := keyring.Address
		if keyring.DelegateFees {
			feeRecipient = types.KeyringCollectorName
		}
		err := k.SplitKeyringFee(ctx, msg.Creator, feeRecipient, req.Fee)
		if err != nil {
			return nil, err
		}
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

	if req.Fee > 0 {
		err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx,
			types.KeyringEscrowName,
			sdk.MustAccAddressFromBech32(req.Creator),
			sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(req.Fee))),
		)
		if err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventKeyRequestRejected,
			sdk.NewAttribute(types.AttributeRequestId, strconv.FormatUint(req.GetId(), 10)),
		),
	})

	return &types.MsgFulfilKeyRequestResponse{}, nil
}

// validateECDSAKeyDetails performs common validation for ECDSA-based public keys.
func (k msgServer) validateECDSAKeyDetails(ctx sdk.Context, req *types.KeyRequest, pubKey []byte, keyTypeStr string) (*types.MsgFulfilKeyRequestResponse, error) {
	keyLen := len(pubKey)
	if keyLen != 33 && keyLen != 65 {
		return k.rejectKeyRequest(ctx, req, fmt.Sprintf("invalid %s public key %x of length %v, expected 33 or 65 bytes", keyTypeStr, pubKey, keyLen))
	}

	// Constraint 1: PK must start with prefix 0x02 or 0x03 if ECDSA type.
	// This applies to all ECDSA types passed to this function, regardless of length.
	if pubKey[0] != 0x02 && pubKey[0] != 0x03 {
		return k.rejectKeyRequest(ctx, req, fmt.Sprintf("invalid %s public key %x: prefix must be 0x02 or 0x03, got %x", keyTypeStr, pubKey, pubKey[0]))
	}

	// Constraint 2: It should not contain more than 4 leading zeros after the prefix.
	// This means if pubKey[1] through pubKey[5] are all zero, it's a violation.
	// (i.e., 5 or more leading zeros after the prefix is not allowed).
	// This check is relevant as keyLen is guaranteed to be 33 or 65.
	numLeadingZerosAfterPrefix := 0
	for i := 1; i <= 5; i++ { // Check byte indices 1 through 5
		if pubKey[i] == 0 {
			numLeadingZerosAfterPrefix++
		} else {
			break // Stop counting at the first non-zero byte
		}
	}

	if numLeadingZerosAfterPrefix >= 5 {
		return k.rejectKeyRequest(ctx, req, fmt.Sprintf("invalid %s public key %x: contains %d leading zeros after prefix (max 4 allowed)", keyTypeStr, pubKey, numLeadingZerosAfterPrefix))
	}

	return nil, nil
}

func (k msgServer) validatePublicKey(ctx sdk.Context, req *types.KeyRequest, keyType types.KeyType, pubKey []byte) (*types.MsgFulfilKeyRequestResponse, error) {
	switch keyType {
	case types.KeyType_KEY_TYPE_ECDSA_SECP256K1:
		return k.validateECDSAKeyDetails(ctx, req, pubKey, "ecdsa_secp256k1")
	case types.KeyType_KEY_TYPE_EDDSA_ED25519:
		if l := len(pubKey); l != ed25519.PublicKeySize {
			return k.rejectKeyRequest(ctx, req, fmt.Sprintf("invalid eddsa_ed25519 public key %x of length %v, expected %d bytes", pubKey, l, ed25519.PublicKeySize))
		}
	case types.KeyType_KEY_TYPE_BITCOIN_SECP256K1:
		return k.validateECDSAKeyDetails(ctx, req, pubKey, "bitcoin_secp256k1")
	default:
		return k.rejectKeyRequest(ctx, req, fmt.Sprintf("invalid key type: %v", keyType.String()))
	}
	return nil, nil // Should be unreachable due to default case, but good practice.
}

func (k msgServer) rejectKeyRequest(ctx sdk.Context, req *types.KeyRequest, rejectReason string) (*types.MsgFulfilKeyRequestResponse, error) {
	req.Status = types.KeyRequestStatus_KEY_REQUEST_STATUS_REJECTED
	req.RejectReason = rejectReason

	if err := k.KeyRequestStore.Set(ctx, req.Id, *req); err != nil {
		return nil, err
	}

	return &types.MsgFulfilKeyRequestResponse{}, nil
}
