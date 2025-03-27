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
func ExtractEVMChainID(input string) (uint64, error) {
	namespace, id, err := ExtractCAIP2Parts(input)
	if err != nil {
		return 0, err
	}

	if namespace != "eip155" {
		return 0, errors.New("CAIP-2 is not of EVM type")
	}

	chainID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, err
	}

	return chainID, nil
}

// TODO(Sasha): add support for Solana
// ValidateChainID validates a chain ID that can be either a CAIP-2 string or a uint64.
// For CAIP-2 strings, currently only EVM chains (eip155 namespace) are supported.
// Returns the validated uint64 chain ID and any error.
func ValidateChainID(ctx context.Context, chainID any) (uint64, error) {
	var evmChainID uint64

	switch chainIDInput := chainID.(type) {
	case string:
		if IsValidCAIP2(chainIDInput) {
			var err error
			namespace, _, err := ExtractCAIP2Parts(chainIDInput)
			if err != nil {
				return 0, fmt.Errorf("invalid CAIP-2 chain ID: %w", err)
			}
			if namespace != "eip155" {
				return 0, fmt.Errorf("unsupported chain type: %s (only EVM/eip155 chains are supported)", namespace)
			}
			evmChainID, err = ExtractEVMChainID(chainIDInput)
			if err != nil {
				return 0, fmt.Errorf("invalid EVM chain ID: %w", err)
			}
		} else {
			return 0, fmt.Errorf("invalid chain ID format: %s (expected uint64 or CAIP-2 format with eip155 namespace)", chainIDInput)
		}
	case uint64:
		evmChainID = chainIDInput
	default:
		return 0, fmt.Errorf("unsupported chain ID type: %T (expected uint64 or string)", chainID)
	}

	allowedChainIDs := []uint64{17000}
	if strings.HasPrefix(sdk.UnwrapSDKContext(ctx).ChainID(), "diamond") {
		allowedChainIDs = append(allowedChainIDs, 1)
	}
	if !slices.Contains(allowedChainIDs, evmChainID) {
		return 0, fmt.Errorf("unsupported EVM chain ID: %d (allowed values: %v)", evmChainID, allowedChainIDs)
	}

	return evmChainID, nil
}

// ExtractSolanaNetwork checks if a CAIP-2 string is Solana-based and extracts the network type.
// Returns the network type (mainnet, testnet, or devnet) and an error if not a valid Solana CAIP-2.
func ExtractSolanaNetwork(input string) (string, error) {
	namespace, reference, err := ExtractCAIP2Parts(input)
	if err != nil {
		return "", err
	}

	if namespace != "solana" {
		return "", errors.New("CAIP-2 is not of Solana type")
	}

	// Validate that the reference is one of the allowed Solana networks
	allowedNetworks := []string{"mainnet", "testnet", "devnet"}
	if !slices.Contains(allowedNetworks, reference) {
		return "", fmt.Errorf("invalid Solana network: %s (allowed values: %v)", reference, allowedNetworks)
	}

	return reference, nil
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
