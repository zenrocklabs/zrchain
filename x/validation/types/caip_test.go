package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidCAIP2(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid Ethereum Mainnet", "eth:mainnet", true},
		{"Valid Bitcoin Mainnet", "bitcoin:mainnet", true},
		{"Valid Solana Devnet", "solana:devnet", true},
		{"Invalid Missing Colon", "invalid_format", false},
		{"Invalid Hyphen Instead of Colon", "eth-mainnet", false},
		{"Invalid Extra Colon", "eth:mainnet:extra", false},
		{"Empty String", "", false},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, IsValidCAIP2(tc.input))
		})
	}
}

func TestExtractCAIP2Parts(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectedNS  string
		expectedRef string
		expectError bool
	}{
		{"Valid Ethereum Mainnet", "eth:mainnet", "eth", "mainnet", false},
		{"Valid Bitcoin Mainnet", "bitcoin:mainnet", "bitcoin", "mainnet", false},
		{"Valid Solana Devnet", "solana:devnet", "solana", "devnet", false},
		{"Invalid Missing Colon", "invalid_format", "", "", true},
		{"Invalid Hyphen Instead of Colon", "eth-mainnet", "", "", true},
		{"Invalid Extra Colon", "eth:mainnet:extra", "", "", true},
		{"Empty String", "", "", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ns, ref, err := ExtractCAIP2Parts(tc.input)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedNS, ns)
				require.Equal(t, tc.expectedRef, ref)
			}
		})
	}
}
