package types

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

// nolint:stylecheck,st1003
// revive:disable-next-line var-naming
func (k *Key) SetId(id uint64) { k.Id = id }

// nolint:stylecheck,st1003
// revive:disable-next-line var-naming
func (kr *KeyRequest) SetId(id uint64) { kr.Id = id }

// NewMsgFulfilKeyRequestKey is a utility function to generate a new successful
// FulfilKeyRequest result.
func NewMsgFulfilKeyRequestKey(publicKey []byte) isMsgFulfilKeyRequest_Result {
	return &MsgFulfilKeyRequest_Key{
		Key: &MsgNewKey{
			PublicKey: publicKey,
		},
	}
}

// NewMsgFulfilKeyRequestReject is a utility function to generate a new errored
// FulfilKeyRequest result.
func NewMsgFulfilKeyRequestReject(reason string) isMsgFulfilKeyRequest_Result {
	return &MsgFulfilKeyRequest_RejectReason{
		RejectReason: reason,
	}
}

// ToECDSASecp256k1 returns the key parsed as a ECDSA secp256k1 public key.
// It can be parssed as a compressed or uncompressed key.
func (k *Key) ToECDSASecp256k1() (*ecdsa.PublicKey, error) {
	if k.Type != KeyType_KEY_TYPE_ECDSA_SECP256K1 {
		return nil, fmt.Errorf("invalid key type, expected %s, got %s", KeyType_KEY_TYPE_ECDSA_SECP256K1, k.Type)
	}

	var pk *ecdsa.PublicKey

	if len(k.PublicKey) == 33 {
		// Compressed form
		var err error
		pk, err = crypto.DecompressPubkey(k.PublicKey)
		if err != nil {
			return nil, err
		}
	} else {
		// Uncompressed form
		var err error

		pk, err = crypto.UnmarshalPubkey(k.PublicKey)
		if err != nil {
			return nil, err
		}
	}

	return pk, nil
}

func (k *Key) ToBitcoinSecp256k1() (*ecdsa.PublicKey, error) {
	if k.Type != KeyType_KEY_TYPE_BITCOIN_SECP256K1 {
		return nil, fmt.Errorf("invalid key type, expected %s, got %s", KeyType_KEY_TYPE_BITCOIN_SECP256K1, k.Type)
	}

	var pk *ecdsa.PublicKey

	if len(k.PublicKey) == 33 {
		// Compressed form
		var err error
		pk, err = crypto.DecompressPubkey(k.PublicKey)
		if err != nil {
			return nil, err
		}
	} else {
		// Uncompressed form
		var err error

		pk, err = crypto.UnmarshalPubkey(k.PublicKey)
		if err != nil {
			return nil, err
		}
	}

	return pk, nil
}

// ToEdDSAEd25519 returns the key parsed as a EdDSA Ed25519 public key.
func (k *Key) ToEdDSAEd25519() (*ed25519.PublicKey, error) {
	if k.Type != KeyType_KEY_TYPE_EDDSA_ED25519 {
		return nil, fmt.Errorf("invalid key type, expected %s, got %s", KeyType_KEY_TYPE_EDDSA_ED25519, k.Type)
	}

	var pk *ed25519.PublicKey

	if len(k.PublicKey) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid key length, expect 32, got %d and key %v", len(k.PublicKey), k)
	}

	pubKey := ed25519.PublicKey(k.PublicKey)
	pk = &pubKey
	return pk, nil
}

func Caip2ToKeyType(caip string) (KeyType, error) {
	switch caip {
	case "solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp:7S3P4HxJpyyigGzodYwHtCxZyUQe9JiBMHyRWXArAaKv", // solana mainnet
		"solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1:DYw8jCTfwHNRJhhmFcbXvVDTqWMEVFBX6ZKUmG5CNSKK", // solana devnet
		"solana:4uhcVJyU9pJkvQyS88uRDiswHXSCkY3z:6LmSRCiu3z6NCSpF19oz1pHXkYkN4jWbj9K1nVELpDkT": // solana testnet
		return KeyType_KEY_TYPE_EDDSA_ED25519, nil
	case "eip155:1", // eth mainnet
		"eip155:11155111",  // sepolia
		"eip155:137",       // polygon main
		" eip155:80002",    // polygon amoy
		"eip155:56",        // bnb smartchan main
		"eip155:97",        // bnb smartchain test
		"eip155:43114",     // avalanche main
		"eip155:43113",     // avalanche fuji test
		"eip155:168587773": // blast sepolia
		return KeyType_KEY_TYPE_ECDSA_SECP256K1, nil
	}
}
