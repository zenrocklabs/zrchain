package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"regexp"
)

const SOLANA_TOKEN_ADDRESS = "CpvmayWu1wnZDbny8uiaojhBVtLjVwS8Zif2b4h6mhWt"
const SOLANA_ROCK_PROGRAM_ID = "9CNTbJY29vHPThkMXCVNozdhXtWrWHyxVy39EhpRtiXe"
const SOLANA_SIGNER_ADDRESS = "5V6vvNUuKvjgcb1faPCTKqWGwevgeENk63fnqbDtJTDB"
const SOLANA_FEE_WALLET = "6Uz2z3pfQgn5SpGp9THV6uWqAQdDVxF3df8XCiTgf8Qq"
const ZENROCK_ALICE_ADDRESS = "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"

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
