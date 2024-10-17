package shared

import "fmt"

func WalletTypeToKeyType(wt uint64) (string, error) {
	switch wt {
	case 0, 1:
		return "bitcoin", nil
	case 60: // ethereum
		return "ecdsa", nil
	case 501: // solana
		return "ed25519", nil

	}

	return "", fmt.Errorf("unknown wallet type")
}
