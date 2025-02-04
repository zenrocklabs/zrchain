package types

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
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
	ns, id, err := ExtractCAIP2Parts(input)
	if err != nil {
		return 0, err
	}

	if ns != "eip155" {
		return 0, errors.New("CAIP-2 is not of EVM type")
	}

	chainId, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return 0, err
	}

	return chainId, nil
}
