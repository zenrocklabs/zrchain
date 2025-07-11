package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zenrocklabs/goem/ethereum"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

func TestIsSolanaCAIP2(t *testing.T) {
	mainnetCtx := sdk.Context{}.WithChainID("diamond-1")
	devnetCtx := sdk.Context{}.WithChainID("amber-1")

	testCases := []struct {
		name     string
		ctx      sdk.Context
		input    string
		expected bool
	}{
		// Mainnet context tests
		{"Mainnet Context - Valid Mainnet Chain", mainnetCtx, "solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp", true},
		{"Mainnet Context - Invalid Devnet Chain", mainnetCtx, "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1", false},
		{"Mainnet Context - Invalid Testnet Chain", mainnetCtx, "solana:4uhcVJyU9pJkvQyS88uRDiswHXSCkY3z", false},
		{"Mainnet Context - Invalid Solana Network", mainnetCtx, "solana:invalidnetwork", false},
		{"Mainnet Context - Invalid Namespace", mainnetCtx, "eth:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp", false},
		{"Mainnet Context - Invalid Format", mainnetCtx, "solana-5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp", false},
		{"Mainnet Context - Empty String", mainnetCtx, "", false},

		// Devnet context tests
		{"Devnet Context - Valid Devnet Chain", devnetCtx, "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1", true},
		{"Devnet Context - Invalid Mainnet Chain", devnetCtx, "solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp", false},
		{"Devnet Context - Invalid Solana Network", devnetCtx, "solana:invalidnetwork", false},
		{"Devnet Context - Invalid Namespace", devnetCtx, "eth:EtWTRABZaYq6iMfeYKouRu166VU2xqa1", false},
		{"Devnet Context - Invalid Format", devnetCtx, "solana-EtWTRABZaYq6iMfeYKouRu166VU2xqa1", false},
		{"Devnet Context - Empty String", devnetCtx, "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, IsSolanaCAIP2(tc.ctx, tc.input))
		})
	}
}

func TestValidateEVMChainID(t *testing.T) {
	mainnetCtx := sdk.Context{}.WithChainID("diamond-1")
	devnetCtx := sdk.Context{}.WithChainID("zenrock-1")

	testCases := []struct {
		name        string
		ctx         sdk.Context
		input       string
		expectedID  uint64
		expectError bool
	}{
		// Mainnet context tests
		{"Mainnet Ctx - Valid Mainnet EVM", mainnetCtx, "eip155:1", 1, false},
		{"Mainnet Ctx - Valid Holesky EVM", mainnetCtx, ethereum.HoleskyCAIP2, ethereum.HoleskyChainId.Uint64(), false},
		{"Mainnet Ctx - Valid Hoodi EVM", mainnetCtx, ethereum.HoodiCAIP2, ethereum.HoodiChainId.Uint64(), false},
		{"Mainnet Ctx - Invalid EVM Chain", mainnetCtx, "eip155:137", 0, true},
		{"Mainnet Ctx - Invalid Namespace", mainnetCtx, "eth:1", 0, true},
		{"Mainnet Ctx - Not CAIP-2", mainnetCtx, "not-caip2", 0, true},

		// Devnet context tests
		{"Devnet Ctx - Valid Holesky EVM", devnetCtx, ethereum.HoleskyCAIP2, ethereum.HoleskyChainId.Uint64(), false},
		{"Devnet Ctx - Valid Hoodi EVM", devnetCtx, ethereum.HoodiCAIP2, ethereum.HoodiChainId.Uint64(), false},
		{"Devnet Ctx - Invalid Mainnet EVM", devnetCtx, "eip155:1", 0, true},
		{"Devnet Ctx - Invalid EVM Chain", devnetCtx, "eip155:137", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := ValidateEVMChainID(tc.ctx, tc.input)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedID, id)
			}
		})
	}
}

func TestValidateSolanaChainID(t *testing.T) {
	mainnetCtx := sdk.Context{}.WithChainID("diamond-1")
	devnetCtx := sdk.Context{}.WithChainID("zenrock-1")

	testCases := []struct {
		name        string
		ctx         sdk.Context
		input       string
		expectedRef string
		expectError bool
	}{
		// Mainnet context tests
		{"Mainnet Ctx - Valid Mainnet Chain", mainnetCtx, "solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp", "5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp", false},
		{"Mainnet Ctx - Invalid Devnet Chain", mainnetCtx, "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1", "", true},
		{"Mainnet Ctx - Invalid Namespace", mainnetCtx, "sol:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp", "", true},
		{"Mainnet Ctx - Not CAIP-2", mainnetCtx, "not-caip2", "", true},

		// Devnet context tests
		{"Devnet Ctx - Valid Devnet Chain", devnetCtx, "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1", "EtWTRABZaYq6iMfeYKouRu166VU2xqa1", false},
		{"Devnet Ctx - Invalid Mainnet Chain", devnetCtx, "solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ref, err := ValidateSolanaChainID(tc.ctx, tc.input)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedRef, ref)
			}
		})
	}
}
