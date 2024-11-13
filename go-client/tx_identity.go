package client

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
)

// Package client provides functionality for managing addresses and signing transactions
// on the Zenrock blockchain. This file implements the Identity type and methods for
// creating new addresses from seed phrases.

// addrPrefix defines the Bech32 prefix for Zenrock addresses
// All addresses on the Zenrock network start with this prefix
var addrPrefix = "zen"

// init configures the SDK's address prefix settings.
// This ensures that all generated addresses use the correct Zenrock prefix.
func init() {
	// set up SDK config (singleton)
	config := sdktypes.GetConfig()
	config.SetBech32PrefixForAccount(addrPrefix, addrPrefix+"pub")
}

// Identity represents an account on Zenrock. It contains both the public address
// and private key necessary for signing transactions on the network.
//
// Fields:
//   - Address: The Bech32-encoded account address (starts with 'zen')
//   - PrivKey: The secp256k1 private key used for signing
type Identity struct {
	Address sdktypes.AccAddress
	PrivKey *secp256k1.PrivKey
}

// NewIdentityFromSeed creates a new Identity from a BIP39 seed phrase and derivation path.
// This method should primarily be used for testing purposes, as storing seed phrases
// in production environments is not recommended for security reasons.
//
// Parameters:
//   - derivationPath: The BIP32/44 derivation path (e.g., "m/44'/118'/0'/0/0")
//   - seedPhrase: A BIP39 mnemonic seed phrase
//
// Returns:
//   - Identity: A new identity containing the derived address and private key
//   - error: An error if seed phrase conversion or key derivation fails
//
// Example:
//
//	identity, err := NewIdentityFromSeed(
//	    "m/44'/118'/0'/0/0",
//	    "seed phrase words here...",
//	)
//	if err != nil {
//	    // Handle error
//	}
//
// Security Note:
//
// This method should NOT be used in production environments. Instead, use secure
// key management solutions like hardware security modules (HSMs) or key management
// services (KMS) to store and manage private keys.
func NewIdentityFromSeed(derivationPath, seedPhrase string) (Identity, error) {
	// Convert the seed phrase to a seed
	seedBytes, err := bip39.NewSeedWithErrorChecking(seedPhrase, "")
	if err != nil {
		return Identity{}, fmt.Errorf("failed to convert seed phrase to seed: %w", err)
	}

	// Create a master key and derive the desired key
	masterKey, ch := hd.ComputeMastersFromSeed(seedBytes)
	derivedKey, err := hd.DerivePrivateKeyForPath(masterKey, ch, derivationPath)
	if err != nil {
		return Identity{}, fmt.Errorf("failed to derive private key: %w", err)
	}

	// Generate a private key object from the bytes
	privKey, _ := btcec.PrivKeyFromBytes(derivedKey)

	// Convert the public key to a Cosmos secp256k1.PublicKey
	cosmosPrivKey := &secp256k1.PrivKey{
		Key: privKey.Serialize(),
	}

	// Get the address of the public key
	addr := sdktypes.AccAddress(cosmosPrivKey.PubKey().Address())

	return Identity{
		Address: addr,
		PrivKey: cosmosPrivKey,
	}, nil
}
