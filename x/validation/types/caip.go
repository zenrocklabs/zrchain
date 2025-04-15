package types

import (
	context "context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ZChainID string

func (z *ZChainID) Uint64() uint64 {
	res, _ := strconv.ParseUint(string(*z), 10, 64)
	return res
}

// IsValidCAIP2 checks if a string follows the CAIP-2 format: "<namespace>:<reference>"
func IsValidCAIP2(input string) bool {
	regex := `^[a-zA-Z0-9\-]+:[a-zA-Z0-9\-]+$`
	re := regexp.MustCompile(regex)
	return re.MatchString(input)
}

// ExtractCAIP2Parts extracts the namespace and reference from a CAIP-2 string.
// Returns an error if the string is not a valid CAIP-2.
func ExtractCAIP2Parts(input string) (string, string, error) {
	if !IsValidCAIP2(input) {
		return "", "", errors.New("invalid CAIP-2 string")
	}

	parts := strings.Split(input, ":")
	return parts[0], parts[1], nil
}

// ExtractEVMChainID Checks if a CAIP-2 string is EVM based and extracts the chain ID.
func ExtractEVMChainID(input string) (ZChainID, error) {
	namespace, id, err := ExtractCAIP2Parts(input)
	if err != nil {
		return "", err
	}

	if namespace != "eip155" {
		return "", errors.New("CAIP-2 is not of EVM type")
	}

	return ZChainID(id), nil
}

// TODO(Sasha): add support for Solana
// ValidateChainID validates a chain ID that can be either a CAIP-2 string or a uint64.
// For CAIP-2 strings, currently only EVM chains (eip155 namespace) are supported.
// Returns the validated uint64 chain ID and any error.
func ValidateChainID(ctx context.Context, chainID any) (ZChainID, error) {
	var chID ZChainID

	switch chainIDInput := chainID.(type) {
	case string:
		if IsValidCAIP2(chainIDInput) {
			var err error
			namespace, _, err := ExtractCAIP2Parts(chainIDInput)
			if err != nil {
				return "", fmt.Errorf("invalid CAIP-2 chain ID: %w", err)
			}
			if namespace == "eip155" {
				chID, err = ExtractEVMChainID(chainIDInput)
				if err != nil {
					return "", fmt.Errorf("invalid EVM chain ID: %w", err)
				}
			} else if namespace == "solana" {
				chID, err = ExtractSolanaNetwork(chainIDInput)
				if err != nil {
					return "", fmt.Errorf("invalid Solana chain ID: %w", err)
				}
			} else {
				return "", fmt.Errorf("unsupported chain type: %s (only EVM/eip155 chains are supported)", namespace)
			}
		} else {
			return "", fmt.Errorf("invalid chain ID format: %s (expected uint64 or CAIP-2 format with eip155 namespace)", chainIDInput)
		}
	case uint64:
		chID = ZChainID(strconv.FormatUint(chainIDInput, 10))
	default:
		return "", fmt.Errorf("unsupported chain ID type: %T (expected uint64 or string)", chainID)
	}

	allowedChainIDs := []ZChainID{"17000", "EtWTRABZaYq6iMfeYKouRu166VU2xqa1"}
	if strings.HasPrefix(sdk.UnwrapSDKContext(ctx).ChainID(), "diamond") {
		allowedChainIDs = append(allowedChainIDs, "1")
	}
	if !slices.Contains(allowedChainIDs, chID) {
		return "", fmt.Errorf("unsupported EVM chain ID: %s (allowed values: %v)", chainID, allowedChainIDs)
	}

	return chID, nil
}

// ExtractSolanaNetwork checks if a CAIP-2 string is Solana-based and extracts the network type.
// Returns the network type (mainnet, testnet, or devnet) and an error if not a valid Solana CAIP-2.
func ExtractSolanaNetwork(input string) (ZChainID, error) {
	namespace, reference, err := ExtractCAIP2Parts(input)
	if err != nil {
		return "", err
	}

	if namespace != "solana" {
		return "", errors.New("CAIP-2 is not of Solana type")
	}

	// Validate that the reference is one of the allowed Solana networks
	allowedNetworks := []string{
		"5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp", // mainnet
		"EtWTRABZaYq6iMfeYKouRu166VU2xqa1", // devnet
		"4uhcVJyU9pJkvQyS88uRDiswHXSCkY3z", // testnet
	}
	if !slices.Contains(allowedNetworks, reference) {
		return "", fmt.Errorf("invalid Solana network: %s (allowed values: %v)", reference, allowedNetworks)
	}

	return ZChainID(reference), nil
}

// IsSolanaCAIP2 checks if a CAIP-2 string represents a Solana network (mainnet, testnet, or devnet).
func IsSolanaCAIP2(input string) bool {
	network, err := ExtractSolanaNetwork(input)
	return err == nil && network != ""
}

// IsEthereumCAIP2 checks if a CAIP-2 string represents an Ethereum network (eip155 namespace).
func IsEthereumCAIP2(input string) bool {
	_, err := ExtractEVMChainID(input)
	return err == nil
}
