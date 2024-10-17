package types

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (sig *AdditionalSignaturePasskey) GetConfigId() string {
	return base64.RawURLEncoding.EncodeToString(sig.RawId)
}

func (sig *AdditionalSignaturePasskey) Verify(_ sdk.Context, config SignMethod, act Action) string {
	cfg, ok := config.(*SignMethodPasskey)
	if !ok {
		return ""
	}

	addr := config.GetParticipantId()
	if !sig.verifyChallenge(addr, act) {
		return ""
	}

	pubkeyBytes, err := cfg.GetPublicKey()
	if err != nil {
		return ""
	}

	pubkey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     big.NewInt(0).SetBytes(pubkeyBytes[0:32]),
		Y:     big.NewInt(0).SetBytes(pubkeyBytes[32:64]),
	}

	clientDataHash := sha256.Sum256(sig.ClientDataJson)
	sigData := append(sig.AuthenticatorData, clientDataHash[:]...)

	type ECDSASignature struct {
		R, S *big.Int
	}

	e := &ECDSASignature{}
	f := crypto.SHA256.New
	h := f()

	h.Write(sigData)

	if _, err := asn1.Unmarshal(sig.Signature, e); err != nil {
		return ""
	}

	if ecdsa.Verify(pubkey, h.Sum(nil), e.R, e.S) {
		return addr
	}

	return ""
}

func (sig *AdditionalSignaturePasskey) verifyChallenge(addr string, act Action) bool {
	clientDataJson := CollectedClientData{}
	err := json.Unmarshal(sig.ClientDataJson, &clientDataJson)
	if err != nil {
		return false
	}

	var challenge []byte = nil
	for _, kv := range act.PolicyData {
		if kv.Key == "challenge-"+addr {
			challenge = kv.GetValue()
			break
		}
	}

	if challenge == nil {
		return false
	}

	incomingChallenge, err := base64.RawURLEncoding.DecodeString(clientDataJson.Challenge)
	if err != nil {
		return false
	}

	if !bytes.Equal(challenge, incomingChallenge) {
		return false
	}

	return true
}
