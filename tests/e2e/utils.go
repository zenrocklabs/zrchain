package e2e

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gagliardetto/solana-go"
)

const ZENROCK_ALICE_ADDRESS = "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"

// Solana ROCK constants
const SOLANA_ROCK_TOKEN_ADDRESS = "CpvmayWu1wnZDbny8uiaojhBVtLjVwS8Zif2b4h6mhWt"
const SOLANA_ROCK_PROGRAM_ID = "9CNTbJY29vHPThkMXCVNozdhXtWrWHyxVy39EhpRtiXe"
const SOLANA_SIGNER_ADDRESS = "5V6vvNUuKvjgcb1faPCTKqWGwevgeENk63fnqbDtJTDB"

// Solana ZENBTC constants
const SOLANA_ZENBTC_TOKEN_ADDRESS = "3Ed6taAtCmDuw9wJnYLKJUcCtGq3YqXz98Cf4oX1b11p"
const SOLANA_ZENBTC_PROGRAM_ID = "9Gfr1YrMca5hyYRDP2nGxYkBWCSZtBm1oXBZyBdtYgNL"
const SOLANA_ZENBTC_MULTISIG = "6eKCiGUhCbi8KgHUEeoTXgSDBZEV1LipBQJVL2EqsABu"

const SOLANA_FEE_WALLET = "6Uz2z3pfQgn5SpGp9THV6uWqAQdDVxF3df8XCiTgf8Qq"

// Ethereum ZENBTC constants
const ETHEREUM_ZENBTC_CONTRACT = "0x745Aa06072bf149117C457C68b0531cF7567C4e1"
const ETHEREUM_CHAIN_ID = 560048

func extractSolanaPubkey(output string) (string, error) {
	re := regexp.MustCompile(`pubkey:\s+([A-Za-z0-9]+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return "", fmt.Errorf("pubkey not found")
	}
	return matches[1], nil
}

// LoadSolanaKeypair loads a solana keypair from a given path
func LoadSolanaPrivateKeyFromJSON(rawJSON string) (*solana.Wallet, error) {
	var keyBytes []byte
	err := json.Unmarshal([]byte(rawJSON), &keyBytes)
	if err != nil {
		return nil, fmt.Errorf("unmarshal key: %w", err)
	}
	priv := solana.PrivateKey(keyBytes)
	wallet := &solana.Wallet{
		PrivateKey: priv,
	}
	return wallet, nil
}

// ethereumWalletFromOutput reads the output of the cast command to create a wallet
func ethereumWalletFromOutput(output string) (*ecdsa.PrivateKey, error) {
	addrRe := regexp.MustCompile(`(?m)^Address:\s*(0x[a-fA-F0-9]{40})`)
	privRe := regexp.MustCompile(`(?m)^Private key:\s*(0x[a-fA-F0-9]{64})`)

	addrMatch := addrRe.FindStringSubmatch(output)
	privMatch := privRe.FindStringSubmatch(output)

	if len(addrMatch) < 2 || len(privMatch) < 2 {
		return nil, errors.New("failed to extract address or private key")
	}

	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(privMatch[1], "0x"))
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func extractBTCBalance(jsonStr string) (float64, error) {
	// Define a struct to match the JSON field we want
	var result struct {
		TotalAmount float64 `json:"total_amount"`
	}

	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return 0, err
	}

	return result.TotalAmount, nil
}

func randomBTCRegnetAddress() (string, error) {
	// Generate a new private key
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return "", err
	}

	// Get the compressed public key
	pubKey := privKey.PubKey()
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())

	// Generate a native SegWit address (bech32 starting with bcrt1)
	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, &chaincfg.RegressionNetParams)
	if err != nil {
		return "", err
	}

	return addr.EncodeAddress(), nil
}
