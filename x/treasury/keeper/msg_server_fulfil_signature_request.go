package keeper

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/btcsuite/btcd/btcec/v2"
	bitcoinecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"

	bitcoinutils "github.com/Zenrock-Foundation/zrchain/v6/bitcoin"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
)

func (k msgServer) FulfilSignatureRequest(goCtx context.Context, msg *types.MsgFulfilSignatureRequest) (*types.MsgFulfilSignatureRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	req, key, err := k.fulfilRequestSetup(ctx, msg.RequestId)
	if err != nil {
		return nil, err
	}

	if req.Status != types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING &&
		req.Status != types.SignRequestStatus_SIGN_REQUEST_STATUS_PARTIAL {
		return nil, fmt.Errorf("request is not pending/partial, can't be updated")
	}

	if err := k.validateZenBTCSignRequest(ctx, *req, *key); err != nil {
		req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED
		req.RejectReason = err.Error()
	} else {
		switch msg.Status {
		case types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED, types.SignRequestStatus_SIGN_REQUEST_STATUS_PARTIAL:
			if err := k.handleSignatureRequest(ctx, msg, req, key); err != nil {
				return nil, err
			}
		case types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED:
			if err := k.handleSignatureRequestRejection(ctx, msg, req); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("invalid status field %s, should be either fulfilled/partial/rejected", msg.Status)
		}
	}

	if err := k.SignRequestStore.Set(ctx, req.Id, *req); err != nil {
		return nil, fmt.Errorf("failed to set sign request: %w", err)
	}

	keyring, err := k.identityKeeper.GetKeyring(ctx, key.KeyringAddr)
	if err != nil || !keyring.IsActive {
		req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED
		req.RejectReason = fmt.Sprintf("keyring %s is inactive or not found", key.KeyringAddr)
		return nil, nil
	}
	if !keyring.IsParty(msg.Creator) && req.Status != types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED {
		return nil, fmt.Errorf("only one party of the keyring can fulfil signature request")
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

	eventType := types.EventSignRequestFulfilled
	if req.Status == types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED {
		eventType = types.EventSignRequestRejected
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute(types.AttributeRequestId, strconv.FormatUint(req.GetId(), 10)),
		),
	})

	return &types.MsgFulfilSignatureRequestResponse{}, nil
}

func (k msgServer) fulfilRequestSetup(ctx sdk.Context, requestID uint64) (*types.SignRequest, *types.Key, error) {
	req, err := k.SignRequestStore.Get(ctx, requestID)
	if err != nil {
		return nil, nil, fmt.Errorf("request not found")
	}

	key, err := k.KeyStore.Get(ctx, req.KeyIds[0])
	if err != nil {
		return nil, nil, fmt.Errorf("key not found")
	}

	return &req, &key, nil
}

func (k msgServer) handleSignatureRequest(ctx sdk.Context, msg *types.MsgFulfilSignatureRequest, req *types.SignRequest, key *types.Key) error {
	sigData := msg.GetSignedData()
	if len(sigData) == 0 {
		req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED
		req.RejectReason = "missing signature data"
		return nil
	}

	// Reject if a party tries to sign more than once
	for _, sig := range req.KeyringPartySigs {
		if sig.Creator == msg.Creator {
			req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED
			errMsg := fmt.Sprintf("party %v already sent a fulfilment", msg.Creator)
			req.RejectReason = errMsg
			return nil
		}
	}

	// Reject invalid signature
	if msg.KeyringPartySignature == nil || len(msg.KeyringPartySignature) != 64 {
		req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED
		req.RejectReason = "invalid length of mpc party signature"
		return nil
	}

	keyring, err := k.identityKeeper.GetKeyring(ctx, key.KeyringAddr)
	if err != nil || !keyring.IsActive {
		req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED
		req.RejectReason = fmt.Sprintf("keyring %s is nil or is inactive", key.KeyringAddr)
		return nil
	}

	if req.KeyType == types.KeyType_KEY_TYPE_BITCOIN_SECP256K1 {
		sigDataBitcoin, err := bitcoinutils.ConvertECDSASigtoBitcoinSig(hex.EncodeToString(sigData))
		if err != nil {
			req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED
			req.RejectReason = fmt.Sprintf("failed to convert ecdsa signature to bitcoin signature: %v", err)
			return nil
		}
		sigData, err = hex.DecodeString(sigDataBitcoin)
		if err != nil {
			req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED
			req.RejectReason = fmt.Sprintf("failed to decode bitcoin signature: %v", err)
			return nil
		}
	}

	if len(req.DataForSigning) == 1 {
		if err := k.verifySignature(ctx, req, key, sigData); err != nil {
			req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED
			req.RejectReason = fmt.Sprintf("signature verification failed: %v", err)
			return nil
		}
	}

	// Append party signature
	req.KeyringPartySigs = append(req.KeyringPartySigs, &types.PartySignature{
		Creator:   msg.Creator,
		Signature: msg.KeyringPartySignature,
	})

	// Check against signed data from other parties
	if len(req.SignedData) == 0 {
		// Only append if this is the first party to respond
		req.SignedData = append(req.SignedData, &types.SignedDataWithID{
			SignRequestId: msg.RequestId,
			SignedData:    sigData,
		})
	} else {
		// Store signed data from first party's response so we can check other parties respond with the same data
		if !bytes.Equal(req.SignedData[0].SignedData, sigData) {
			req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED
			req.RejectReason = fmt.Sprintf("signed data mismatch, expected %x, got %x", req.SignedData[0].SignedData, sigData)
			return nil
		}
	}

	if len(req.KeyringPartySigs) >= int(keyring.PartyThreshold) {
		req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED

		if req.ParentReqId != 0 {
			if err := k.updateParentRequest(ctx, req, sigData, msg.RequestId); err != nil {
				return err
			}
		}
	} else {
		req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_PARTIAL
	}

	return nil
}

func (k msgServer) handleSignatureRequestRejection(ctx sdk.Context, msg *types.MsgFulfilSignatureRequest, req *types.SignRequest) error {
	req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED
	rejectReason := msg.GetRejectReason()
	if rejectReason == "" {
		rejectReason = "rejected with no reason provided"
	}
	req.RejectReason = rejectReason

	if req.ParentReqId != 0 {
		parentReq, err := k.SignRequestStore.Get(ctx, req.ParentReqId)
		if err != nil {
			return fmt.Errorf("parent request not found: %w", err)
		}
		parentReq.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED
		parentReq.RejectReason = "Child request " + strconv.FormatUint(req.Id, 10) + " rejected with reason: " + rejectReason

		if err := k.SignRequestStore.Set(ctx, parentReq.Id, parentReq); err != nil {
			return fmt.Errorf("failed to set parent sign request: %w", err)
		}
	}

	return nil
}

// calculateV calculates the value of v based on the given signature, public key, and hash.
// It iterates through two possible values (0,1) for the byte at index 64 of the signature (the last byte)
// and checks if the resulting public key matches the decompressed public key. If a match is
// found, it returns the calculated value of v. If no match is found, it returns an error.
func (k msgServer) calculateV(sigBytes, pubkeyBytes, hashBytes []byte) (v byte, err error) {
	pk, err := crypto.DecompressPubkey(pubkeyBytes)
	if err != nil {
		return 0, err
	}
	sigBytesToCheck := make([]byte, len(sigBytes))
	copy(sigBytesToCheck, sigBytes)
	for _, i := range []byte{0, 1} {
		sigBytesToCheck[64] = i
		pkey, err := crypto.SigToPub(hashBytes, sigBytesToCheck)
		if err == nil && pk.X.Cmp(pkey.X) == 0 && pk.Y.Cmp(pkey.Y) == 0 {
			return i, nil
		}
	}

	return 0, errors.New("unable to calculate v")
}

func (k msgServer) verifySignature(ctx sdk.Context, req *types.SignRequest, key *types.Key, sigData []byte) error {
	valid := false
	switch req.KeyType {
	case types.KeyType_KEY_TYPE_ECDSA_SECP256K1:
		if len(sigData) != 64 && len(sigData) != 65 {
			return fmt.Errorf("verifySignature- invalid ecdsa signature %x of length %v", sigData, len(sigData))
		}
		valid = crypto.VerifySignature(key.PublicKey, req.DataForSigning[0], sigData[:64])

		if valid && len(sigData) == 65 {
			v, err := k.calculateV(sigData, key.PublicKey, req.DataForSigning[0])
			if err != nil {
				ctx.Logger().Warn(err.Error())
			} else {
				sigData[64] = v
			}
		}
	case types.KeyType_KEY_TYPE_EDDSA_ED25519:
		if len(sigData) != ed25519.SignatureSize {
			return fmt.Errorf("verifySignature- invalid eddsa signature %x of length %v", sigData, len(sigData))
		}
		valid = ed25519.Verify(key.PublicKey, req.DataForSigning[0], sigData)
	case types.KeyType_KEY_TYPE_BITCOIN_SECP256K1:
		sig, err := bitcoinecdsa.ParseDERSignature(sigData)
		if err != nil {
			return fmt.Errorf("verifySignature - invalid Bitcoin signature %x - fail to parse ", sigData)
		}
		pubKey, err := btcec.ParsePubKey(key.PublicKey)
		if err != nil {
			return fmt.Errorf("verifySignature- invalid Bitcoin Public Key %x ", key.PublicKey)
		}
		valid = sig.Verify(req.DataForSigning[0], pubKey)
	default:
		return fmt.Errorf("verifySignature- invalid key type: %v", req.KeyType.String())
	}
	if !valid {
		return fmt.Errorf("verifySignature- invalid signature %x from keyring %s", sigData, key.KeyringAddr)
	}
	return nil
}

func (k msgServer) updateParentRequest(ctx sdk.Context, req *types.SignRequest, sigData []byte, requestId uint64) error {
	parentReq, err := k.SignRequestStore.Get(ctx, req.ParentReqId)
	if err != nil {
		return fmt.Errorf("parent request not found: %w", err)
	}

	sigExists := false
	for _, sd := range parentReq.SignedData {
		if bytes.Equal(sd.SignedData, sigData) {
			sigExists = true
			break
		}
	}
	if !sigExists {
		parentReq.SignedData = append(parentReq.SignedData, &types.SignedDataWithID{
			SignRequestId: requestId,
			SignedData:    sigData,
		})
		sort.Slice(parentReq.SignedData, func(i, j int) bool {
			return parentReq.SignedData[i].SignRequestId < parentReq.SignedData[j].SignRequestId
		})
	}

	if len(parentReq.SignedData) >= len(parentReq.DataForSigning) {
		parentReq.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED
	}

	if err := k.SignRequestStore.Set(ctx, parentReq.Id, parentReq); err != nil {
		return fmt.Errorf("failed to set parent sign request: %w", err)
	}

	return nil
}
