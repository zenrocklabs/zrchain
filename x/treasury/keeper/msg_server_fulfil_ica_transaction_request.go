package keeper

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) FulfilICATransactionRequest(goCtx context.Context, msg *types.MsgFulfilICATransactionRequest) (*types.MsgFulfilICATransactionRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	req, key, err := k.fulfilRequestSetup(ctx, msg.RequestId)
	if err != nil {
		return nil, err
	}

	if err = k.processICATransactionRequest(ctx, msg.Status, req, key, msg.GetSignedData(), msg.KeyringPartySignature, msg.Creator, msg.GetRejectReason(), msg.RequestId); err != nil {
		return nil, err
	}

	if err := k.SignRequestStore.Set(ctx, req.Id, *req); err != nil {
		return nil, fmt.Errorf("failed to set sign request: %w", err)
	}

	eventType := types.EventICATransactionRequestFulfilled
	if req.Status == types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED {
		eventType = types.EventSignRequestRejected
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute(types.AttributeRequestId, strconv.FormatUint(req.GetId(), 10)),
		),
	})

	return &types.MsgFulfilICATransactionRequestResponse{}, nil
}

func (k msgServer) processICATransactionRequest(ctx sdk.Context, status types.SignRequestStatus, req *types.SignRequest, key *types.Key, sigData []byte, keyringPartySignature []byte, creator string, rejectReason string, requestId uint64) error {
	if err := k.validateICATransactionRequest(ctx, req, key, sigData, keyringPartySignature, status); err != nil {
		return err
	}

	switch status {
	case types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED, types.SignRequestStatus_SIGN_REQUEST_STATUS_PARTIAL:
		if err := k.handleICATransactionRequestFulfilment(ctx, req, key, sigData, keyringPartySignature, creator, requestId); err != nil {
			return err
		}
	case types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED:
		req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED
		req.RejectReason = rejectReason
	default:
		return fmt.Errorf("invalid status field %s, should be either fulfilled/partial/rejected", status)
	}

	return nil
}

func (k msgServer) validateICATransactionRequest(ctx sdk.Context, req *types.SignRequest, key *types.Key, sigData []byte, keyringPartySignature []byte, status types.SignRequestStatus) error {
	if status == types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED {
		return nil
	}

	if len(sigData) == 0 {
		return fmt.Errorf("missing signature data")
	}

	if keyringPartySignature == nil || len(keyringPartySignature) != 64 {
		return fmt.Errorf("invalid mpc party signature")
	}

	if len(req.DataForSigning) != 1 || len(req.DataForSigning[0]) != 32 {
		return fmt.Errorf("data for signing should be a single 32-byte hash for ICA transaction requests")
	}

	return k.verifySignature(ctx, req, key, sigData)
}

func (k msgServer) handleICATransactionRequestFulfilment(ctx sdk.Context, req *types.SignRequest, key *types.Key, sigData []byte, keyringPartySignature []byte, creator string, requestId uint64) error {
	sigExists := false
	for _, sig := range req.KeyringPartySigs {
		if bytes.Equal(sig.Signature, keyringPartySignature) {
			sigExists = true
		}
	}
	if !sigExists {
		req.KeyringPartySigs = append(req.KeyringPartySigs, &types.PartySignature{
			Creator:   creator,
			Signature: keyringPartySignature,
		})
	}

	keyring, err := k.identityKeeper.GetKeyring(ctx, key.KeyringAddr)
	if err != nil || !keyring.IsActive {
		return fmt.Errorf("keyring %s is nil or is inactive", key.KeyringAddr)
	}

	if len(req.KeyringPartySignatures) >= int(keyring.PartyThreshold) {
		req.SignedData = append(req.SignedData, &types.SignedDataWithID{
			SignRequestId: requestId,
			SignedData:    sigData,
		})
		req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED
	} else {
		req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_PARTIAL
	}

	return nil
}
