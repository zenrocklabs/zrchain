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
