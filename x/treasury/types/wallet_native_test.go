package types

import (
	"crypto/sha256"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func Test_NativeWallet_Address(t *testing.T) {
	wallet := nativeWallet(t)
	require.Equal(t, "zen1egz60et40xxzm5rhtlj7caskpvqmqujrgeauaa", wallet.Address())
}

func nativeWallet(t *testing.T) *NativeWallet {
	t.Helper()
	hashedSeed := sha256.Sum256([]byte("example seed"))

	// Generate secp256k1 private key from the hashed seed
	privateKey, err := crypto.ToECDSA(hashedSeed[:])
	if err != nil {
		log.Fatal("Failed to generate private key:", err)
	}

	// Serialize the public key in a compressed format
	publicKeyBytes := crypto.CompressPubkey(&privateKey.PublicKey)

	k := &Key{
		Id:            0,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		Type:          KeyType_KEY_TYPE_ECDSA_SECP256K1,
		PublicKey:     publicKeyBytes,
	}

	wallet, err := NewNativeWallet(k)
	require.NoError(t, err)
	return wallet
}
