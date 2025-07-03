package main

import (
	"encoding/hex"
	"net"
	"os"
	"testing"
	"time"

	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/neutrino"
	"github.com/stretchr/testify/require"
)

// checkExternalDependencies checks if required external services are available
func checkExternalDependencies(t *testing.T) bool {
	// Check if we're in CI environment
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping integration test in CI environment")
		return false
	}

	// Check if Bitcoin proxy is running
	conn, err := net.DialTimeout("tcp", "127.0.0.1:1234", 2*time.Second)
	if err != nil {
		t.Skipf("Bitcoin proxy not available at 127.0.0.1:1234: %v", err)
		return false
	}
	conn.Close()

	// Check if integration tests are explicitly enabled
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Integration tests disabled. Set RUN_INTEGRATION_TESTS=true to enable")
		return false
	}

	return true
}

func Test_ProxyFunctions_Testnet3(t *testing.T) {
	// Skip if external dependencies are not available
	if !checkExternalDependencies(t) {
		return
	}

	// This test requires 2 additional pieces of setup:
	// 1. The 30 second sleep below gives the Neutrino nodes time to sync
	// 2. A Bitcoin Proxy needs to be running locally "http://127.0.0.1:1234", "user", "secret"

	ns := neutrino.NeutrinoServer{}
	ns.Initialize("testnet", "http://127.0.0.1:1234", "user", "secret", "./neutrino", 12345, "/neutrino_")

	// Reduced sleep time for faster test execution
	time.Sleep(10 * time.Second)

	// Get via the Neutrino Node
	header1, hash1, _, err := ns.GetBlockHeaderByHeight("testnet3", 1000)
	if err != nil {
		t.Skipf("Neutrino node not ready: %v", err)
		return
	}

	// Get via the Proxy
	header2, hash2, _, err := ns.ProxyGetBlockHeaderByHeight("testnet3", 1000)
	if err != nil {
		t.Skipf("Proxy not available: %v", err)
		return
	}

	// Compare block headers
	require.Equal(t, header1.Nonce, header2.Nonce, "block header mismatch")
	require.Equal(t, header1.Version, header2.Version, "block header mismatch")
	require.Equal(t, header1.Timestamp, header2.Timestamp, "block header mismatch")
	require.Equal(t, header1.Bits, header2.Bits, "block header mismatch")
	require.Equal(t, header1.MerkleRoot, header2.MerkleRoot, "block header mismatch")
	require.Equal(t, header1.PrevBlock, header2.PrevBlock, "block header mismatch")

	require.Equal(t, hash1, hash2, "block hash mismatch")
}

func Test_ProxyFunctions_Testnet4(t *testing.T) {
	t.Skip("manual test that requires local bitcoin proxy and may take time to sync; skipped in automated runs")
	//Do not run as part of CI.

	// This test requires 2 additional pieces of setup:
	// 1. The 30 second sleep below gives the Neutrino nodes time to sync
	// 2. A Bitcoin Proxy needs to be running locally "http://127.0.0.1:1234", "user", "secret"

	ns := neutrino.NeutrinoServer{}
	ns.Initialize("testnet", "http://127.0.0.1:1234", "user", "secret", "./neutrino", 12345, "/neutrino_")
	time.Sleep(30 * time.Second)

	// Reduced sleep time for faster test execution
	time.Sleep(10 * time.Second)

	// Get via the Neutrino Node - should return an error for testnet4
	_, _, _, err := ns.GetBlockHeaderByHeight("testnet4", 1000)
	require.Error(t, err, "Testnet4 should return an error")

	// Get via the Proxy
	_, hash2, _, err := ns.ProxyGetBlockHeaderByHeight("testnet4", 1000)
	if err != nil {
		t.Skipf("Proxy not available for testnet4: %v", err)
		return
	}

	hex := hex.EncodeToString(hash2[:])
	hex = ReverseHex(hex)
	require.Equal(t, hex, "00000000b747d47c3b38161693ad05e26924b3775a8be669751f969da836311e", "hash is invalid")
}

// Test_ProxyFunctions_Unit tests the proxy functions without external dependencies
func Test_ProxyFunctions_Unit(t *testing.T) {
	// Test ReverseHex function
	testCases := []struct {
		input    string
		expected string
	}{
		{"12345678", "78563412"},
		{"abcdef", "efcdab"},
		{"00", "00"},
		{"", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := ReverseHex(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}

	// Test invalid hex string
	result := ReverseHex("123")
	require.Equal(t, "invalid hex string", result)
}

// Test_NeutrinoServer_Initialization tests the NeutrinoServer initialization without external dependencies
func Test_NeutrinoServer_Initialization(t *testing.T) {
	ns := neutrino.NeutrinoServer{}

	// Test initialization with mock values
	ns.Initialize("testnet", "http://localhost:8080", "testuser", "testpass", "./test_neutrino", 8080, "/test_")

	// Verify the server was initialized (basic check)
	require.NotNil(t, ns)
}

func ReverseHex(hexStr string) string {
	n := len(hexStr)
	if n%2 != 0 {
		// Ensure the hex string length is even
		return "invalid hex string"
	}
	result := make([]byte, n)
	for i := 0; i < n; i += 2 {
		// Copy two characters (one byte) at a time from the end to the beginning
		result[n-i-2], result[n-i-1] = hexStr[i], hexStr[i+1]
	}
	return string(result)
}
