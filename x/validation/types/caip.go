package types

import (
	context "context"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zenrocklabs/goem/ethereum"
	"github.com/zenrocklabs/goem/solana"
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
		return "", "", fmt.Errorf("invalid CAIP-2 string: %s", input)
	}

	parts := strings.Split(input, ":")
	return parts[0], parts[1], nil
}

// ValidateEVMChainID validates an EVM chain ID from a CAIP-2 string (e.g., "eip155:1")
// and returns the chain ID as a uint64.
func ValidateEVMChainID(ctx context.Context, caip2 string) (uint64, error) {
	namespace, reference, err := ExtractCAIP2Parts(caip2)
	if err != nil {
		return 0, err
	}

	if namespace != ethereum.CAIP2Namespace {
		return 0, fmt.Errorf("CAIP-2 is not of EVM type (%s): %s", ethereum.CAIP2Namespace, caip2)
	}

	allowedEVMChainIDs := []string{ethereum.HoleskyChainId.String(), ethereum.HoodiChainId.String()}
	if strings.HasPrefix(sdk.UnwrapSDKContext(ctx).ChainID(), "diamond") {
		allowedEVMChainIDs = append(allowedEVMChainIDs, "1")
	}

	if !slices.Contains(allowedEVMChainIDs, reference) {
		return 0, fmt.Errorf("unsupported EVM chain ID: %s (allowed: %v)", reference, allowedEVMChainIDs)
	}

	chainID, err := strconv.ParseUint(reference, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid EVM chain ID reference: %s", reference)
	}

	return chainID, nil
}

// ValidateSolanaChainID validates a Solana chain ID from a CAIP-2 string.
// It enforces a strict check against the allowed Solana network for the current environment.
func ValidateSolanaChainID(ctx context.Context, caip2 string) (string, error) {
	namespace, reference, err := ExtractCAIP2Parts(caip2)
	if err != nil {
		return "", err
	}

	if namespace != solana.CAIP2Namespace {
		return "", fmt.Errorf("CAIP-2 is not of Solana type: %s", caip2)
	}

	var allowedChain string
	if strings.HasPrefix(sdk.UnwrapSDKContext(ctx).ChainID(), "diamond") {
		allowedChain = "solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp" // mainnet
	} else {
		if strings.HasPrefix(sdk.UnwrapSDKContext(ctx).ChainID(), "docker") {
			allowedChain = "solana:HK8b7Skns2TX3FvXQxm2mPQbY2nVY8GD" // regnet
		} else {
			allowedChain = "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1" // devnet
		}
	}

	if caip2 != allowedChain {
		return "", fmt.Errorf("invalid Solana network: %s (allowed value: %s)", caip2, allowedChain)
	}

	return reference, nil
}

// ValidateChainID validates a chain ID from a CAIP-2 string.
// Deprecated: Use ValidateEVMChainID or ValidateSolanaChainID directly for type safety
// and to avoid unnecessary logic.
func ValidateChainID(ctx context.Context, caip2 string) (string, error) {
	namespace, reference, err := ExtractCAIP2Parts(caip2)
	if err != nil {
		return "", err
	}

	switch namespace {
	case ethereum.CAIP2Namespace:
		_, err := ValidateEVMChainID(ctx, caip2)
		if err != nil {
			return "", err
		}
		return reference, nil
	case solana.CAIP2Namespace:
		return ValidateSolanaChainID(ctx, caip2)
	default:
		return "", fmt.Errorf("unsupported CAIP-2 namespace: %s (supported: %s, %s)", namespace, ethereum.CAIP2Namespace, solana.CAIP2Namespace)
	}
}

// IsSolanaCAIP2 checks if a CAIP-2 string represents a valid Solana network.
func IsSolanaCAIP2(ctx context.Context, input string) bool {
	_, err := ValidateSolanaChainID(ctx, input)
	return err == nil
}

// IsEthereumCAIP2 checks if a CAIP-2 string represents a valid Ethereum network.
func IsEthereumCAIP2(ctx context.Context, input string) bool {
	_, err := ValidateEVMChainID(ctx, input)
	return err == nil
}
