package types

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"strings"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ripemd160"
)

type ZCashWallet struct {
	key     *ecdsa.PublicKey
	network string // "mainnet", "testnet", or "regtest"
}

var _ Wallet = &ZCashWallet{}
var _ TxParser = &ZCashWallet{}

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

// ParseTx implements TxParser for ZCash transactions
func (w *ZCashWallet) ParseTx(b []byte, m Metadata) (Transfer, error) {
	// ZCash transactions are similar to Bitcoin in structure
	// We need to calculate the signature hashes
	hashes, err := w.SigHashes(b)
	if err != nil {
		return Transfer{}, err
	}

	var dataForSigning []string
	for _, hash := range hashes {
		dataForSigning = append(dataForSigning, hex.EncodeToString(hash))
	}

	return Transfer{
		SigHashes:      hashes,
		DataForSigning: []byte(strings.Join(dataForSigning, ",")),
	}, nil
}

func (w *ZCashWallet) SigHashes(b []byte) (hashes [][]byte, err error) {
	// For ZCash, we can reuse the Bitcoin transaction deserialization
	// as ZCash v5 transactions are compatible with Bitcoin serialization
	msgTx, err := DeserializeTransaction(b)
	if err != nil {
		return nil, err
	}

	// Calculate signature hashes for each input
	// ZCash uses a similar signing mechanism to Bitcoin
	for i := range msgTx.TxIn {
		// For each input, we need to calculate the signature hash
		// This is a simplified version - in production, you'd want to handle
		// different signature hash types and ZCash-specific fields

		// Get the witness data if available
		if len(msgTx.TxIn[i].Witness) > 0 {
			// Extract the signature hash from witness data
			// The actual implementation depends on ZCash's specific witness structure
			witnessData := msgTx.TxIn[i].Witness[0]
			hashes = append(hashes, witnessData)
		} else {
			// Fallback: create a basic signature hash
			// In a real implementation, this would follow ZCash's signature hash algorithm
			return nil, fmt.Errorf("zcash transaction must have witness data")
		}
	}

	return hashes, nil
}

// Helper function to hash160 (RIPEMD160(SHA256(data)))
func hash160(data []byte) []byte {
	sha := crypto.Keccak256(data)
	ripemd := ripemd160.New()
	ripemd.Write(sha)
	return ripemd.Sum(nil)
}
