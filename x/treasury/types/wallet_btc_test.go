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
	unsignedPlusTX := "0100000000010115c3c881a51c5edc9d10018d5e63b7dfbcb319a723598cca6dacad440d25407b0000000000ffffffff02e80300000000000016001492ff4f363ec63adde62fbb232b916b2a3050d5a0d20f00000000000017a9141be8b5d7e39d0624c5d9452ba820ded1e89316d48706160014050431872c2326d2075f4eb03ef5b8efc4b8cab20810270000000000002103ac3f68dd29a1628283191069bcfb43a74fbbf752808ab8f3c86c14a68aeece6917a9141be8b5d7e39d0624c5d9452ba820ded1e89316d4870801000000000000002051fbd9df8f361bdc88fa1fe0b8148fa56c9d30ddbe507fa5edf1978d0499202000000000"
	utx, _ := hex.DecodeString(unsignedPlusTX)
	wallet := btcWallet(t)
	transfer, err := wallet.ParseTx(utx, nil)
	require.Nil(t, err, "Error parsing tx")

	hexhash := hex.EncodeToString(transfer.SigHashes[0])
	require.True(t, hexhash == "51fbd9df8f361bdc88fa1fe0b8148fa56c9d30ddbe507fa5edf1978d04992020", "Error parsing tx")
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
