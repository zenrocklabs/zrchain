package types

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	solana "github.com/gagliardetto/solana-go"
)

func NativeAddress(key *Key, prefix string) (string, error) {
	if prefix == "" {
		prefix = "zen"
	}

	k, err := key.ToECDSASecp256k1()
	if err != nil {
		return "", err
	}
	var pubkey secp256k1.PubKey
	pubkey.Key = crypto.CompressPubkey(k)
	bech32Address := sdk.MustBech32ifyAddressBytes(prefix, pubkey.Address())

	return bech32Address, nil
}

func EthereumAddress(key *Key) (string, error) {
	k, err := key.ToECDSASecp256k1()
	if err != nil {
		return "", err
	}
	addr := crypto.PubkeyToAddress(*k)
	return addr.Hex(), nil
}

func BitcoinP2SH(key *Key, chain *chaincfg.Params) (string, error) {
	k, err := key.ToBitcoinSecp256k1()
	if err != nil {
		return "", err
	}
	var pubkey secp256k1.PubKey
	pubkey.Key = crypto.CompressPubkey(k)
	witnessProg := btcutil.Hash160(pubkey.Key)
	p2wpkhAddress, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, chain)
	if err != nil {
		return "", fmt.Errorf("error creating new P2SH address (NewAddressWitnessPubKeyHash): %w", err)
	}
	p2wpkhScript, err := txscript.PayToAddrScript(p2wpkhAddress)
	if err != nil {
		return "", fmt.Errorf("error creating new P2SH address (PayToAddrScript): %w", err)
	}
	p2shAddress, err := btcutil.NewAddressScriptHash(p2wpkhScript, chain)
	if err != nil {
		return "", fmt.Errorf("error creating new P2SH address (NewAddressScriptHash): %w", err)
	}
	return p2shAddress.EncodeAddress(), nil
}

func BitcoinP2WPKH(key *Key, chain *chaincfg.Params) (string, error) {
	k, err := key.ToBitcoinSecp256k1()
	if err != nil {
		return "", err
	}
	var pubkey secp256k1.PubKey
	pubkey.Key = crypto.CompressPubkey(k)
	publicKey, err := btcec.ParsePubKey(pubkey.Key)
	if err != nil {
		return "", err
	}
	// Generate P2WPKH address from the public key
	witnessProg := btcutil.Hash160(publicKey.SerializeCompressed())
	address, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, chain)
	if err != nil {
		return "", err
	}
	return address.EncodeAddress(), nil
}

func SolanaAddress(key *Key) (string, error) {
	_, err := key.ToEdDSAEd25519()
	if err != nil {
		return "", err
	}
	pk := solana.PublicKeyFromBytes(key.PublicKey)
	return pk.String(), nil
}
