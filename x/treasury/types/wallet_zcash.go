package types

import (
	"crypto/ecdsa"
	"crypto/sha256"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/ethereum/go-ethereum/crypto"
)

type ZCashWallet struct {
	key     *ecdsa.PublicKey
	network string // "mainnet", "testnet", or "regtest"
}

var _ Wallet = &ZCashWallet{}

// Zcash transparent address P2PKH version bytes (2 bytes, unlike Bitcoin's 1 byte)
var (
	// Mainnet P2PKH version bytes [0x1C, 0xB8] - produces "t1" prefix
	zcashMainnetP2PKHVersionBytes = []byte{0x1C, 0xB8}
	// Testnet P2PKH version bytes [0x1D, 0x25] - produces "tm" prefix (also used for Regtest)
	zcashTestnetP2PKHVersionBytes = []byte{0x1D, 0x25}
)

func NewZCashWallet(k *Key, network string) (*ZCashWallet, error) {
	pubkey, err := k.ToBitcoinSecp256k1()
	if err != nil {
		return nil, err
	}
	return &ZCashWallet{key: pubkey, network: network}, nil
}

func (w *ZCashWallet) Address() string {
	var pubkey secp256k1.PubKey
	pubkey.Key = crypto.CompressPubkey(w.key)
	publicKey, err := btcec.ParsePubKey(pubkey.Key)
	if err != nil {
		return ""
	}

	// Generate P2PKH transparent address from the public key
	pubKeyHash := btcutil.Hash160(publicKey.SerializeCompressed())

	// Select the appropriate version bytes based on network
	var versionBytes []byte
	switch w.network {
	case "mainnet":
		versionBytes = zcashMainnetP2PKHVersionBytes
	case "testnet", "regtest":
		versionBytes = zcashTestnetP2PKHVersionBytes
	default:
		return ""
	}

	// Encode the address using Base58Check with Zcash's 2-byte version prefix
	address := encodeZCashBase58Check(versionBytes, pubKeyHash)
	return address
}

// encodeZCashBase58Check encodes a Zcash transparent address using Base58Check
// with a 2-byte version prefix (unlike Bitcoin which uses 1 byte)
func encodeZCashBase58Check(versionBytes []byte, payload []byte) string {
	// Combine version bytes and payload
	b := make([]byte, 0, len(versionBytes)+len(payload)+4)
	b = append(b, versionBytes...)
	b = append(b, payload...)

	// Calculate checksum: first 4 bytes of double SHA256
	checksum := doubleSHA256(b)[:4]

	// Append checksum
	b = append(b, checksum...)

	// Base58 encode
	return base58.Encode(b)
}

// doubleSHA256 calculates SHA256(SHA256(data))
func doubleSHA256(data []byte) []byte {
	first := sha256.Sum256(data)
	second := sha256.Sum256(first[:])
	return second[:]
}
