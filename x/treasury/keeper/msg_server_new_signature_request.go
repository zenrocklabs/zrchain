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
	if err = validatePayload(payload, key.Type); err != nil {
		return nil, err
	}

	signPolicyID := key.SignPolicyId

	if signPolicyID == 0 {
		ws, getDataErr := k.identityKeeper.GetWorkspace(ctx, key.WorkspaceAddr)
		if getDataErr != nil {
			return nil, fmt.Errorf("workspace %s not found: %v", key.WorkspaceAddr, getDataErr)
		}
		signPolicyID = ws.SignPolicyId
	}

	keyring, err := k.identityKeeper.GetKeyring(ctx, key.KeyringAddr)
	if err != nil || !keyring.IsActive {
		return nil, fmt.Errorf("keyring %s is nil or is inactive", keyring.Address)
	}

	act, err := k.policyKeeper.AddAction(ctx, msg.Creator, msg, signPolicyID, msg.Btl, nil)
	if err != nil {
		return nil, err
	}

	var dataForSigning [][]byte
	for _, p := range payload {
		data, decodeErr := hex.DecodeString(p)
		if decodeErr != nil {
			return nil, decodeErr
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

	ws, err := k.identityKeeper.GetWorkspace(ctx, key.WorkspaceAddr)
	if err != nil {
		return nil, fmt.Errorf("workspace %s not found", key.WorkspaceAddr)
	}

	return ws.PolicyNewSignatureRequest(), nil
}

func (k msgServer) NewSignatureRequestActionHandler(ctx sdk.Context, act *policytypes.Action) (*types.MsgNewSignatureRequestResponse, error) {
	return policykeeper.TryExecuteAction(
		k.policyKeeper,
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

// validatePayload validates the payload according to the following rules:
// 1. The payload slice must not be empty.
// 2. The payload slice must contain at least one non-empty string.
// 3. If the payload has exactly one string:
//   - For non-EDDSA keys (e.g. ECDSA), the hex-encoded string must be exactly ecdsaHexEncodedLength characters long.
//   - For EDDSA keys, the hex-encoded string must be shorter than eddsaMaxHexEncodedLength.
func validatePayload(payload []string, keyType types.KeyType) error {
	const (
		ecdsaHexEncodedLength    = 64   // ECDSA keys require exactly 64 characters.
		eddsaMaxHexEncodedLength = 2000 // EDDSA keys require fewer than 2000 characters.
	)

	if len(payload) == 0 {
		return fmt.Errorf("payload is empty")
	}

	nonEmptyStringFound := false
	for _, item := range payload {
		if strings.TrimSpace(item) != "" {
			nonEmptyStringFound = true
			break
		}
	}
	if !nonEmptyStringFound {
		return fmt.Errorf("payload is full of empty strings")
	}

	for _, data := range payload {
		if keyType != types.KeyType_KEY_TYPE_EDDSA_ED25519 {
			if len(data) != ecdsaHexEncodedLength {
				return fmt.Errorf("data for signing for ecdsa key should have a hex-encoded length of %d, not: %d", ecdsaHexEncodedLength, len(data))
			}
		} else {
			if len(data) >= eddsaMaxHexEncodedLength {
				return fmt.Errorf("data for signing for eddsa key should have a hex-encoded length smaller than %d, not: %d", eddsaMaxHexEncodedLength, len(data))
			}
		}
	}

	return nil
}
