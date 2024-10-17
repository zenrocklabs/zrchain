package types

import (
	"crypto/ecdsa"
	"crypto/sha256"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type NativeWallet struct {
	key *ecdsa.PublicKey
}

var _ Wallet = &NativeWallet{}
var _ TxParser = &NativeWallet{}

func NewNativeWallet(k *Key) (*NativeWallet, error) {
	pubkey, err := k.ToECDSASecp256k1()
	if err != nil {
		return nil, err
	}
	return &NativeWallet{key: pubkey}, nil
}

func (w *NativeWallet) Address() string {
	var pubkey secp256k1.PubKey
	pubkey.Key = crypto.CompressPubkey(w.key)
	bech32Address := sdk.MustBech32ifyAddressBytes("zen", pubkey.Address())
	return bech32Address
}

// ParseTx implements TxParser.
func (w *NativeWallet) ParseTx(b []byte, _ Metadata) (Transfer, error) {
	dataHash := sha256.Sum256(b)
	return Transfer{
		DataForSigning: dataHash[:],
	}, nil
}
