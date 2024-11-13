package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/stretchr/testify/require"
)

func Test_ParseTX(t *testing.T) {
	unsignedPlusTX := "010000000001011cf81f651d4135576e226b9e3afd6d5bb3d204fc87258c732b90c7f49224c36a0100000000ffffffff02701700000000000016001496d1dc8156bf66e09b58582a8fe07fd2cf780cd4b3180000000000001600140b7ae56af7989b75aa086a194276063dc1d34568050877380000000000002103d6b23e342370ae51214d6e198f86ac0065a3dc15c21a2a348060acb7de8a0ad91600140b7ae56af7989b75aa086a194276063dc1d34568080100000000000000209b0c1c721ceea9a96c1e6dd99e37be2d1d0602ba564806ea531663d88cf2de9800000000"
	utx, _ := hex.DecodeString(unsignedPlusTX)
	wallet := btcWallet(t)
	transfer, err := wallet.ParseTx(utx, nil)
	require.Nil(t, err, "Error parsing tx")

	hexhash := hex.EncodeToString(transfer.SigHashes[0])
	require.True(t, hexhash == "9b0c1c721ceea9a96c1e6dd99e37be2d1d0602ba564806ea531663d88cf2de98", "Error parsing tx")
}

func btcWallet(t *testing.T) *BitcoinWallet {
	t.Helper()
	hashedSeed := sha256.Sum256([]byte("example seed"))

	// Generate secp256k1 private key from the hashed seed
	privateKey, err := crypto.ToECDSA(hashedSeed[:])

	if err != nil {
		log.Fatal("Failed to generate private key:", err)
	}

	// Serialize the public key in a compressed format
	publicKeyBytes := crypto.CompressPubkey(&privateKey.PublicKey)
	fmt.Println(hex.EncodeToString(publicKeyBytes))

	k := &Key{
		Id:            0,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		Type:          KeyType_KEY_TYPE_BITCOIN_SECP256K1,
		PublicKey:     publicKeyBytes,
	}

	wallet, err := NewBTCWallet(k, &chaincfg.TestNet3Params)
	require.NoError(t, err)
	return wallet
}
