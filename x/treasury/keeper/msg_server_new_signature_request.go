package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	pol "github.com/Zenrock-Foundation/zrchain/v5/policy"
	policykeeper "github.com/Zenrock-Foundation/zrchain/v5/x/policy/keeper"
	policytypes "github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"

	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) NewSignatureRequest(goCtx context.Context, msg *types.MsgNewSignatureRequest) (*types.MsgNewSignatureRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	key, err := k.KeyStore.Get(ctx, msg.KeyIds[0])
	if err != nil {
		return nil, fmt.Errorf("key %v not found", msg.KeyIds[0])
	}

	payload := strings.Split(msg.DataForSigning, ",")
	if key.Type != types.KeyType_KEY_TYPE_EDDSA_ED25519 && len(payload) == 1 && len(payload[0]) != 64 {
		return nil, fmt.Errorf("data for signing for ecdsa key should be a have a hex-encoded length of 64, not: %d", len(payload[0]))
	}
	if key.Type == types.KeyType_KEY_TYPE_EDDSA_ED25519 && len(payload) == 1 && len(payload[0]) >= 2000 {
		return nil, fmt.Errorf("data for signing for eddsa key should have a hex-encoded length of smaller than 2000, not: %d", len(payload[0]))
	}

	signPolicyID := key.SignPolicyId

	if signPolicyID == 0 {
		ws, err := k.identityKeeper.WorkspaceStore.Get(ctx, key.WorkspaceAddr)
		if err != nil {
			return nil, fmt.Errorf("workspace %s not found: %v", key.WorkspaceAddr, err)
		}
		signPolicyID = ws.SignPolicyId
	}

	keyring, err := k.identityKeeper.KeyringStore.Get(ctx, key.KeyringAddr)
	if err != nil || !keyring.IsActive {
		return nil, fmt.Errorf("keyring %s is nil or is inactive", keyring.Address)
	}

	act, err := k.policyKeeper.AddAction(ctx, msg.Creator, msg, signPolicyID, msg.Btl, nil)
	if err != nil {
		return nil, err
	}

	var dataForSigning [][]byte
	for _, p := range payload {
		data, err := hex.DecodeString(p)
		if err != nil {
			return nil, err
		}
		dataForSigning = append(dataForSigning, data)
	}

	verified, err := VerifyDataForSigning(dataForSigning, msg.VerifySigningData, msg.VerifySigningDataVersion)
	if verified == types.Verification_Failed {
		return nil, fmt.Errorf("transaction & hash verfication transaction did not verify")
	}
	if err != nil {
		return nil, fmt.Errorf("error whilst verifying transaction & hashes %s", err.Error())
	}
	return k.NewSignatureRequestActionHandler(ctx, act)
}

func (k msgServer) NewSignatureRequestPolicyGenerator(ctx sdk.Context, msg *types.MsgNewSignatureRequest) (pol.Policy, error) {
	key, err := k.KeyStore.Get(ctx, msg.KeyIds[0])
	if err != nil {
		return nil, fmt.Errorf("key %v not found", msg.KeyIds[0])
	}

	ws, err := k.identityKeeper.WorkspaceStore.Get(ctx, key.WorkspaceAddr)
	if err != nil {
		return nil, fmt.Errorf("workspace %s not found", key.WorkspaceAddr)
	}

	return ws.PolicyNewSignatureRequest(), nil
}

func (k msgServer) NewSignatureRequestActionHandler(ctx sdk.Context, act *policytypes.Action) (*types.MsgNewSignatureRequestResponse, error) {
	return policykeeper.TryExecuteAction(
		&k.policyKeeper,
		k.cdc,
		ctx,
		act,
		k.HandleSignatureRequest,
	)
}

// VerifyDataForSigning - checks dataforSigning against any supplied verification data
func VerifyDataForSigning(dataForSigning [][]byte, verifySigningDataTx []byte, verificationVersion types.VerificationVersion) (status types.VerificationStatus, err error) {
	if dataForSigning == nil || verifySigningDataTx == nil {
		return types.Verification_NotVerified, nil
	}

	switch verificationVersion {
	case types.VerificationVersion_UNKNOWN:
		return types.Verification_NotVerified, nil
	case types.VerificationVersion_BITCOIN_PLUS:
		// Here we check that the Bitcoin Signatures required are properly derived from the Supplied Transaction
		return types.VerifyBitcoinSigHashes(dataForSigning, verifySigningDataTx)
	default:
		return types.Verification_NotVerified, nil
	}
}
