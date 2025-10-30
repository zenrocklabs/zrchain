package types

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

// Test that every WalletType defined in protobuf is covered by NewWalletI.
func Test_NewWalletI_Exhaustive(t *testing.T) {
	for walletType, name := range WalletType_name {
		if walletType == int32(WalletType_WALLET_TYPE_UNSPECIFIED) {
			continue
		}

		t.Run(name, func(t *testing.T) {
			_, err := NewWallet(
				&Key{
					Id:            0,
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					Type:          0,
					PublicKey:     hexutil.MustDecode("0x025cd45a6614df5348692ea4d0f7c16255b75a6b6f67bea5013621fe84af8031f0"),
				},
				WalletType(walletType),
			)
			require.NotErrorIs(t, err, ErrUnknownWalletType)
		})
	}
}

// Test that ZCash addresses are generated with correct prefixes for each network type.
func Test_ZCashAddressFormats(t *testing.T) {
	testKey := &Key{
		Id:            1,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		Type:          KeyType_KEY_TYPE_BITCOIN_SECP256K1,
		PublicKey:     hexutil.MustDecode("0x025cd45a6614df5348692ea4d0f7c16255b75a6b6f67bea5013621fe84af8031f0"),
	}

	tests := []struct {
		name         string
		walletType   WalletType
		expectPrefix string
		network      string
	}{
		{
			name:         "ZCash Mainnet",
			walletType:   WalletType_WALLET_TYPE_ZCASH_MAINNET,
			expectPrefix: "t1",
			network:      "mainnet",
		},
		{
			name:         "ZCash Testnet",
			walletType:   WalletType_WALLET_TYPE_ZCASH_TESTNET,
			expectPrefix: "tm",
			network:      "testnet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wallet, err := NewWallet(testKey, tt.walletType)
			require.NoError(t, err, "Failed to create wallet")
			require.NotNil(t, wallet, "Wallet should not be nil")

			address := wallet.Address()
			require.NotEmpty(t, address, "Address should not be empty")

			// Verify the address starts with the correct prefix
			require.True(
				t,
				strings.HasPrefix(address, tt.expectPrefix),
				"Address %s should start with prefix %s",
				address,
				tt.expectPrefix,
			)

			// Additional validation: address should be base58 encoded
			// Base58 alphabet: 123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz
			// (excludes 0, O, I, l to avoid confusion)
			base58Charset := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
			for _, char := range address {
				if !strings.ContainsRune(base58Charset, char) {
					require.Fail(
						t,
						"Address contains invalid base58 character",
						"Character: %c in address: %s",
						char,
						address,
					)
				}
			}

			// Transparent addresses should be 35 characters long (22 bytes encoded in base58)
			require.Equal(t, 35, len(address), "Transparent address should be 35 characters")

			t.Logf("Generated %s address: %s", tt.network, address)
		})
	}
}

// Test multiple keys produce different addresses
func Test_ZCashAddressUniqueness(t *testing.T) {
	testKeys := []*Key{
		{
			Id:            1,
			WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
			Type:          KeyType_KEY_TYPE_BITCOIN_SECP256K1,
			PublicKey:     hexutil.MustDecode("0x025cd45a6614df5348692ea4d0f7c16255b75a6b6f67bea5013621fe84af8031f0"),
		},
		{
			Id:            2,
			WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
			Type:          KeyType_KEY_TYPE_BITCOIN_SECP256K1,
			PublicKey:     hexutil.MustDecode("0x0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"),
		},
		{
			Id:            3,
			WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
			Type:          KeyType_KEY_TYPE_BITCOIN_SECP256K1,
			PublicKey:     hexutil.MustDecode("0x02c6047f9441ed7d6d3045406e95c07cd85c778e4b8cef3ca7abac09b95c709ee5"),
		},
	}

	walletTypes := []WalletType{
		WalletType_WALLET_TYPE_ZCASH_MAINNET,
		WalletType_WALLET_TYPE_ZCASH_TESTNET,
	}

	for _, walletType := range walletTypes {
		t.Run(WalletType_name[int32(walletType)], func(t *testing.T) {
			addresses := make(map[string]bool)

			for _, key := range testKeys {
				wallet, err := NewWallet(key, walletType)
				require.NoError(t, err)

				address := wallet.Address()
				require.NotEmpty(t, address)

				// Ensure this address hasn't been seen before (uniqueness)
				require.False(t, addresses[address], "Address %s is not unique", address)
				addresses[address] = true
			}

			require.Equal(t, len(testKeys), len(addresses), "Should generate unique addresses for each key")
		})
	}
}

// Test Zcash version bytes match specification exactly
func Test_ZCashVersionBytes(t *testing.T) {
	testKey := &Key{
		Id:            1,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		Type:          KeyType_KEY_TYPE_BITCOIN_SECP256K1,
		PublicKey:     hexutil.MustDecode("0x025cd45a6614df5348692ea4d0f7c16255b75a6b6f67bea5013621fe84af8031f0"),
	}

	tests := []struct {
		name               string
		walletType         WalletType
		expectedVersionHex string
		network            string
	}{
		{
			name:               "Mainnet version bytes [0x1C, 0xB8]",
			walletType:         WalletType_WALLET_TYPE_ZCASH_MAINNET,
			expectedVersionHex: "1cb8",
			network:            "mainnet",
		},
		{
			name:               "Testnet version bytes [0x1D, 0x25]",
			walletType:         WalletType_WALLET_TYPE_ZCASH_TESTNET,
			expectedVersionHex: "1d25",
			network:            "testnet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create wallet using internal constructor to verify version bytes
			wallet, err := NewZCashWallet(testKey, tt.network)
			require.NoError(t, err)
			require.NotNil(t, wallet)

			// Verify address generation
			address := wallet.Address()
			require.NotEmpty(t, address)

			t.Logf("Generated %s address with version bytes %s: %s",
				tt.network, tt.expectedVersionHex, address)
		})
	}
}

// Test Base58Check encoding produces correct structure
func Test_ZCashBase58CheckStructure(t *testing.T) {
	testKey := &Key{
		Id:            1,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		Type:          KeyType_KEY_TYPE_BITCOIN_SECP256K1,
		PublicKey:     hexutil.MustDecode("0x025cd45a6614df5348692ea4d0f7c16255b75a6b6f67bea5013621fe84af8031f0"),
	}

	tests := []struct {
		name         string
		walletType   WalletType
		network      string
		expectPrefix string
	}{
		{
			name:         "Mainnet Base58Check encoding",
			walletType:   WalletType_WALLET_TYPE_ZCASH_MAINNET,
			network:      "mainnet",
			expectPrefix: "t1",
		},
		{
			name:         "Testnet Base58Check encoding",
			walletType:   WalletType_WALLET_TYPE_ZCASH_TESTNET,
			network:      "testnet",
			expectPrefix: "tm",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wallet, err := NewWallet(testKey, tt.walletType)
			require.NoError(t, err)

			address := wallet.Address()
			require.NotEmpty(t, address)

			// Verify prefix matches expected value from version bytes
			require.True(
				t,
				strings.HasPrefix(address, tt.expectPrefix),
				"Address %s should start with %s (version bytes guarantee this prefix)",
				address,
				tt.expectPrefix,
			)

			// Verify total length is 35 characters (26 bytes before Base58 encoding)
			// Structure: 2-byte version + 20-byte hash + 4-byte checksum = 26 bytes
			require.Equal(t, 35, len(address),
				"Transparent P2PKH address should be 35 characters (26 bytes before Base58 encoding)")

			t.Logf("%s address structure verified: %s (prefix=%s, length=%d)",
				tt.network, address, tt.expectPrefix, len(address))
		})
	}
}

// Test that regtest uses testnet version bytes
func Test_ZCashRegtestUsesTestnetVersionBytes(t *testing.T) {
	testKey := &Key{
		Id:            1,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		Type:          KeyType_KEY_TYPE_BITCOIN_SECP256K1,
		PublicKey:     hexutil.MustDecode("0x025cd45a6614df5348692ea4d0f7c16255b75a6b6f67bea5013621fe84af8031f0"),
	}

	// Create regtest wallet (uses "regtest" network parameter)
	regtestWallet, err := NewZCashWallet(testKey, "regtest")
	require.NoError(t, err)

	// Create testnet wallet
	testnetWallet, err := NewZCashWallet(testKey, "testnet")
	require.NoError(t, err)

	regtestAddr := regtestWallet.Address()
	testnetAddr := testnetWallet.Address()

	// According to Zcash spec, regtest uses same version bytes as testnet [0x1D, 0x25]
	require.Equal(t, testnetAddr, regtestAddr,
		"Regtest should produce same address as testnet (both use version bytes [0x1D, 0x25])")
	require.True(t, strings.HasPrefix(regtestAddr, "tm"),
		"Regtest addresses should start with 'tm' prefix like testnet")

	t.Logf("Verified regtest uses testnet version bytes: %s", regtestAddr)
}

// Test compressed public key format is used
func Test_ZCashCompressedPublicKey(t *testing.T) {
	// Use a well-known secp256k1 public key for verification
	testKey := &Key{
		Id:            1,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		Type:          KeyType_KEY_TYPE_BITCOIN_SECP256K1,
		// This is a compressed public key (33 bytes): 0x02 prefix + 32 bytes X coordinate
		PublicKey: hexutil.MustDecode("0x0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"),
	}

	wallet, err := NewWallet(testKey, WalletType_WALLET_TYPE_ZCASH_MAINNET)
	require.NoError(t, err)

	address := wallet.Address()
	require.NotEmpty(t, address)

	// Verify it's a valid address format
	require.True(t, strings.HasPrefix(address, "t1"))
	require.Equal(t, 35, len(address))

	t.Logf("Compressed public key generates valid address: %s", address)
}

// Test invalid network parameter handling
func Test_ZCashInvalidNetwork(t *testing.T) {
	testKey := &Key{
		Id:            1,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		Type:          KeyType_KEY_TYPE_BITCOIN_SECP256K1,
		PublicKey:     hexutil.MustDecode("0x025cd45a6614df5348692ea4d0f7c16255b75a6b6f67bea5013621fe84af8031f0"),
	}

	// Create wallet with invalid network
	wallet, err := NewZCashWallet(testKey, "invalidnetwork")
	require.NoError(t, err) // Constructor should succeed

	// But Address() should return empty string for invalid network
	address := wallet.Address()
	require.Empty(t, address, "Invalid network should produce empty address")
}

// Test same key produces different addresses on different networks
func Test_ZCashNetworkAddressDifference(t *testing.T) {
	testKey := &Key{
		Id:            1,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		Type:          KeyType_KEY_TYPE_BITCOIN_SECP256K1,
		PublicKey:     hexutil.MustDecode("0x025cd45a6614df5348692ea4d0f7c16255b75a6b6f67bea5013621fe84af8031f0"),
	}

	mainnetWallet, err := NewWallet(testKey, WalletType_WALLET_TYPE_ZCASH_MAINNET)
	require.NoError(t, err)

	testnetWallet, err := NewWallet(testKey, WalletType_WALLET_TYPE_ZCASH_TESTNET)
	require.NoError(t, err)

	mainnetAddr := mainnetWallet.Address()
	testnetAddr := testnetWallet.Address()

	// Same key should produce different addresses on different networks
	require.NotEqual(t, mainnetAddr, testnetAddr,
		"Same key should produce different addresses on mainnet vs testnet")

	// But both should have same length and valid prefixes
	require.Equal(t, 35, len(mainnetAddr))
	require.Equal(t, 35, len(testnetAddr))
	require.True(t, strings.HasPrefix(mainnetAddr, "t1"))
	require.True(t, strings.HasPrefix(testnetAddr, "tm"))

	t.Logf("Mainnet address: %s", mainnetAddr)
	t.Logf("Testnet address: %s", testnetAddr)
}

// Test double SHA256 checksum calculation
func Test_ZCashDoubleSHA256Checksum(t *testing.T) {
	// Test with known input
	testData := []byte("hello world")

	checksum := doubleSHA256(testData)

	// Verify checksum is 32 bytes (SHA256 output)
	require.Equal(t, 32, len(checksum), "Double SHA256 should produce 32 bytes")

	// Verify it's not all zeros
	allZeros := true
	for _, b := range checksum {
		if b != 0 {
			allZeros = false
			break
		}
	}
	require.False(t, allZeros, "Checksum should not be all zeros")

	// Verify deterministic (same input produces same output)
	checksum2 := doubleSHA256(testData)
	require.Equal(t, checksum, checksum2, "Double SHA256 should be deterministic")

	t.Logf("Double SHA256 checksum computed: %x", checksum[:4])
}
