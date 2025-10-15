package types

import (
	"crypto/ecdsa"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/ethereum/go-ethereum/crypto"
)

type ZCashWallet struct {
	key     *ecdsa.PublicKey
	network string // "mainnet", "testnet", or "regtest"
}

var _ Wallet = &ZCashWallet{}

// ZCash unified address HRPs (Human-Readable Parts) for Bech32m encoding
var (
	// Mainnet unified address prefix "u1"
	zcashMainnetHRP = "u"
	// Testnet unified address prefix "utest1"
	zcashTestnetHRP = "utest"
	// Regtest unified address prefix "uregtest1"
	zcashRegtestHRP = "uregtest"
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

	// Generate P2PKH address from the public key
	pubKeyHash := btcutil.Hash160(publicKey.SerializeCompressed())

	// Select the appropriate HRP based on network
	var hrp string
	switch w.network {
	case "mainnet":
		hrp = zcashMainnetHRP
	case "testnet":
		hrp = zcashTestnetHRP
	case "regtest":
		hrp = zcashRegtestHRP
	default:
		return ""
	}

	// Encode the address using Bech32m encoding (used for unified addresses)
	address, err := encodeZCashUnifiedAddress(hrp, pubKeyHash)
	if err != nil {
		return ""
	}
	return address
}

// encodeZCashUnifiedAddress creates a ZCash unified address using Bech32m encoding
func encodeZCashUnifiedAddress(hrp string, payload []byte) (string, error) {
	// Convert payload to 5-bit groups for Bech32m encoding
	converted, err := bech32.ConvertBits(payload, 8, 5, true)
	if err != nil {
		return "", err
	}

	// Encode using Bech32m (version 1) which is used for unified addresses
	encoded, err := bech32.EncodeM(hrp, converted)
	if err != nil {
		return "", err
	}

	return encoded, nil
}
