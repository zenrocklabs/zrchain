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
// For CAIP-2 strings, EVM (eip155 namespace) and Solana chains are supported.
// Returns the validated ZChainID and any error.
func ValidateChainID(ctx context.Context, chainID any) (ZChainID, error) {
	var chID ZChainID
	isEVM := false // Flag to determine if the context is EVM for final validation

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
				isEVM = true
			} else if namespace == "solana" {
				// ExtractSolanaNetwork already validates against its own allowed list of networks
				solanaChID, err := ExtractSolanaNetwork(ctx, chainIDInput)
				if err != nil {
					return "", fmt.Errorf("invalid Solana chain ID: %w", err)
				}
				return solanaChID, nil // Solana ID is validated, return directly
			} else {
				return "", fmt.Errorf("unsupported CAIP-2 namespace: %s (supported: eip155, solana)", namespace)
			}
		} else {
			// Updated error message for non-CAIP2 strings
			return "", fmt.Errorf("invalid chain ID format: %s (expected uint64 or CAIP-2, e.g., eip155:1, solana:network-id)", chainIDInput)
		}
	case uint64:
		chID = ZChainID(strconv.FormatUint(chainIDInput, 10))
		isEVM = true
	default:
		return "", fmt.Errorf("unsupported chain ID type: %T (expected uint64 or string)", chainID)
	}

	// This validation is now specifically for EVM chain IDs
	if isEVM {
		allowedEVMChainIDs := []ZChainID{"17000"}
		if strings.HasPrefix(sdk.UnwrapSDKContext(ctx).ChainID(), "diamond") {
			allowedEVMChainIDs = append(allowedEVMChainIDs, "1")
		}

		if !slices.Contains(allowedEVMChainIDs, chID) {
			// Use chID in the error message, as chainID is the raw input
			return "", fmt.Errorf("unsupported EVM chain ID: %s (allowed: %v)", chID, allowedEVMChainIDs)
		}
	}
	// If it was Solana, it returned early.
	// If it was an unsupported CAIP-2 namespace or invalid format, it returned an error early.
	// If it was EVM, it has now been validated.

	return chID, nil
}

// ExtractSolanaNetwork checks if a CAIP-2 string is Solana-based and extracts the network type.
// Returns the network type (mainnet, testnet, or devnet) and an error if not a valid Solana CAIP-2.
func ExtractSolanaNetwork(ctx context.Context, input string) (ZChainID, error) {
	namespace, reference, err := ExtractCAIP2Parts(input)
	if err != nil {
		return "", err
	}

	if namespace != "solana" {
		return "", errors.New("CAIP-2 is not of Solana type")
	}

	// Validate that the reference is one of the allowed Solana networks
	allowedNetworks := []string{
		"EtWTRABZaYq6iMfeYKouRu166VU2xqa1", // devnet
		"4uhcVJyU9pJkvQyS88uRDiswHXSCkY3z", // testnet
	}
	if strings.HasPrefix(sdk.UnwrapSDKContext(ctx).ChainID(), "diamond") {
		allowedNetworks = append(allowedNetworks, "5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp") // mainnet
	}
	if !slices.Contains(allowedNetworks, reference) {
		return "", fmt.Errorf("invalid Solana network: %s (allowed values: %v)", reference, allowedNetworks)
	}

	return ZChainID(reference), nil
}

// IsSolanaCAIP2 checks if a CAIP-2 string represents a Solana network (mainnet, testnet, or devnet).
func IsSolanaCAIP2(ctx context.Context, input string) bool {
	network, err := ExtractSolanaNetwork(ctx, input)
	return err == nil && network != ""
}

// IsEthereumCAIP2 checks if a CAIP-2 string represents an Ethereum network (eip155 namespace).
func IsEthereumCAIP2(input string) bool {
	_, err := ExtractEVMChainID(input)
	return err == nil
}
