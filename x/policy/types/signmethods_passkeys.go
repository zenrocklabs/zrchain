package types

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/fxamacker/cbor/v2"
)

func (m *SignMethodPasskey) IsActive() bool {
	return m.Active
}

func (m *SignMethodPasskey) SetActive(active bool) {
	m.Active = active
}

func (m *SignMethodPasskey) GetConfigId() string {
	return base64.RawURLEncoding.EncodeToString(m.GetRawId())
}

func (m *SignMethodPasskey) GetParticipantId() string {
	return fmt.Sprintf("passkey{%s}", m.GetConfigId())
}

func (m *SignMethodPasskey) GetPublicKey() ([]byte, error) {
	attObj, err := m.getAttestationObject()
	if err != nil {
		return nil, err
	}

	return attObj.AuthData.AttData.getPublicKey()
}

func (m *SignMethodPasskey) VerifyConfig(ctx sdk.Context) error {
	var incomingTx tx.Tx
	if err := incomingTx.Unmarshal(ctx.TxBytes()); err != nil {
		return err
	}

	clientDataJson := CollectedClientData{}
	err := json.Unmarshal(m.ClientDataJson, &clientDataJson)
	if err != nil {
		return err
	}

	if incomingTx.AuthInfo == nil || len(incomingTx.AuthInfo.SignerInfos) < 1 {
		return fmt.Errorf("no signers in tx")
	}

	challenge := incomingTx.AuthInfo.SignerInfos[0].PublicKey.Value
	incomingChallenge, err := base64.RawURLEncoding.DecodeString(clientDataJson.Challenge)
	if err != nil {
		return err
	}

	if !bytes.Equal(challenge, incomingChallenge) {
		return fmt.Errorf("invalid challenge")
	}

	attObj, err := m.getAttestationObject()
	if err != nil {
		return err
	}

	if attObj.Format == "none" {
		if len(attObj.AttStatement) != 0 {
			return fmt.Errorf("attestation format none with attestation statement present")
		}
	} else {
		// TODO, see go-webauthn/protocol/attestation.go
		ctx.Logger().Warn(fmt.Sprintf("unimplemented attestation format: %s", attObj.Format))
	}

	return nil
}

func (m *SignMethodPasskey) getAttestationObject() (*AttestationObject, error) {
	var attObj AttestationObject
	if _, err := ctap2CBORDecMode.UnmarshalFirst(m.AttestationObject, &attObj); err != nil {
		return nil, err
	}

	if err := attObj.AuthData.unmarshal(attObj.RawAuthData); err != nil {
		return nil, err
	}

	return &attObj, nil
}

const nestedLevelsAllowed = 4

var ctap2CBORDecMode, _ = cbor.DecOptions{
	DupMapKey:       cbor.DupMapKeyEnforcedAPF,
	MaxNestedLevels: nestedLevelsAllowed,
	IndefLength:     cbor.IndefLengthForbidden,
	TagsMd:          cbor.TagsForbidden,
}.DecMode()
var ctap2CBOREncMode, _ = cbor.CTAP2EncOptions().EncMode()

type AuthenticatorFlags byte

const (
	FlagUserPresent AuthenticatorFlags = 1 << iota
	FlagRFU1
	FlagUserVerified
	FlagBackupEligible
	FlagBackupState
	FlagRFU2
	FlagAttestedCredentialData
	FlagHasExtensions
)
const (
	minAuthDataLength     = 37
	minAttestedAuthLength = 55
	maxCredentialIDLength = 1023
)

type AttestedCredentialData struct {
	AAGUID              []byte `json:"aaguid"`
	CredentialID        []byte `json:"credential_id"`
	CredentialPublicKey []byte `json:"public_key"`
}
type AuthenticatorData struct {
	RPIDHash []byte                 `json:"rpid"`
	Flags    AuthenticatorFlags     `json:"flags"`
	Counter  uint32                 `json:"sign_count"`
	AttData  AttestedCredentialData `json:"att_data"`
	ExtData  []byte                 `json:"ext_data"`
}
type AttestationObject struct {
	AuthData     AuthenticatorData
	RawAuthData  []byte                 `json:"authData"`
	Format       string                 `json:"fmt"`
	AttStatement map[string]interface{} `json:"attStmt,omitempty"`
}
type CollectedClientData struct {
	Challenge string `json:"challenge"`
	Origin    string `json:"origin"`
}
type PublicKeyData struct {
	_         bool  `cbor:",keyasint"`
	KeyType   int64 `cbor:"1,keyasint" json:"kty"`
	Algorithm int64 `cbor:"3,keyasint" json:"alg"`
}
type EC2PublicKeyData struct {
	PublicKeyData
	Curve  int64  `cbor:"-1,keyasint,omitempty" json:"crv"`
	XCoord []byte `cbor:"-2,keyasint,omitempty" json:"x"`
	YCoord []byte `cbor:"-3,keyasint,omitempty" json:"y"`
}

func (flag AuthenticatorFlags) HasAttestedCredentialData() bool {
	return (flag & FlagAttestedCredentialData) == FlagAttestedCredentialData
}

func (flag AuthenticatorFlags) HasExtensions() bool {
	return (flag & FlagHasExtensions) == FlagHasExtensions
}

func (a *AuthenticatorData) unmarshal(rawAuthData []byte) error {
	if len(rawAuthData) < minAuthDataLength {
		return fmt.Errorf("invalid auth data, len < 37")
	}

	a.RPIDHash = rawAuthData[:32]
	a.Flags = AuthenticatorFlags(rawAuthData[32])
	a.Counter = binary.BigEndian.Uint32(rawAuthData[33:37])
	remaining := len(rawAuthData) - minAuthDataLength

	if a.Flags.HasAttestedCredentialData() {
		if len(rawAuthData) > minAttestedAuthLength {
			if err := a.unmarshalAttestedData(rawAuthData); err != nil {
				return err
			}
			attDataLen := len(a.AttData.AAGUID) + 2 + len(a.AttData.CredentialID) + len(a.AttData.CredentialPublicKey)
			remaining = remaining - attDataLen
		} else {
			return fmt.Errorf("attested credential flag set but data is missing")
		}
	} else {
		if !a.Flags.HasExtensions() && len(rawAuthData) != 37 {
			return fmt.Errorf("attested credential flag not set")
		}
	}

	if a.Flags.HasExtensions() {
		if remaining != 0 {
			a.ExtData = rawAuthData[len(rawAuthData)-remaining:]
			remaining -= len(a.ExtData)
		} else {
			return fmt.Errorf("extensions flag set but extensions data is missing")
		}
	}

	if remaining != 0 {
		return fmt.Errorf("leftover bytes decoding AuthenticatorData")
	}

	return nil
}

func (a *AuthenticatorData) unmarshalAttestedData(rawAuthData []byte) (err error) {
	a.AttData.AAGUID = rawAuthData[37:53]

	idLength := binary.BigEndian.Uint16(rawAuthData[53:55])
	if len(rawAuthData) < int(55+idLength) {
		return fmt.Errorf("authenticator attestation data length too short")
	}

	if idLength > maxCredentialIDLength {
		return fmt.Errorf("authenticator attestation data credential id length too long")
	}

	a.AttData.CredentialID = rawAuthData[55 : 55+idLength]

	a.AttData.CredentialPublicKey, err = a.unmarshalCredentialPublicKey(rawAuthData[55+idLength:])
	if err != nil {
		return fmt.Errorf("could not unmarshal Credential Public Key: %v", err)
	}

	return nil
}

func (a *AuthenticatorData) unmarshalCredentialPublicKey(keyBytes []byte) (rawBytes []byte, err error) {
	var m interface{}

	if _, err = ctap2CBORDecMode.UnmarshalFirst(keyBytes, &m); err != nil {
		return nil, err
	}

	if rawBytes, err = ctap2CBOREncMode.Marshal(m); err != nil {
		return nil, err
	}

	return rawBytes, nil
}

func (a AttestedCredentialData) getPublicKey() ([]byte, error) {
	rawPubKey := a.CredentialPublicKey
	pubkeyData := EC2PublicKeyData{}

	if _, err := ctap2CBORDecMode.UnmarshalFirst(rawPubKey, &pubkeyData); err != nil {
		return nil, err
	}

	return append(pubkeyData.XCoord, pubkeyData.YCoord...), nil
}
