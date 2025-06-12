package integration

import (
	"crypto/rand"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
)

func CreateECDSAPubkey() ([]byte, error) {
	seed := make([]byte, 64)
	rand.Read(seed)
	privateKey, _ := btcec.PrivKeyFromBytes(seed[:])
	ecdsaPriv := privateKey.ToECDSA()
	prvD := math.PaddedBigBytes(ecdsaPriv.D, 32)
	prv, err := crypto.ToECDSA(prvD)
	if err != nil {
		return nil, err
	}
	pubKeyBytes := crypto.CompressPubkey(&prv.PublicKey)
	return pubKeyBytes, nil
}

func AccAddress() string {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return cosmostypes.AccAddress(addr).String()
}
