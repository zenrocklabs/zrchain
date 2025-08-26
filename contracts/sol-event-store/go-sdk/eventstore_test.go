package eventstore

import (
	"strings"
	"testing"

	"github.com/gagliardetto/solana-go"
)

func TestBitcoinAddressType(t *testing.T) {
	tests := []struct {
		address      string
		expectedType string
	}{
		// P2PKH (Legacy) addresses
		{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "P2PKH (Legacy)"},
		{"1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2", "P2PKH (Legacy)"},

		// P2SH addresses
		{"3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy", "P2SH"},
		{"3QJmV3qfvL9SuYo34YihAf3sRCW3qSinyC", "P2SH"},

		// Bech32 P2WPKH addresses (42 chars)
		{"bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4", "P2WPKH (Bech32)"},
		{"bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh", "P2WPKH (Bech32)"},

		// Bech32 P2WSH addresses (62 chars)
		{"bc1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3qccfmv3", "P2WSH (Bech32)"},

		// Bech32m P2TR (Taproot) addresses (62 chars)
		{"bc1p5d7rjq7g6rdk2yhzks9smlaqtedr4dekq08ge8ztwac72sfr9rusxg3297", "P2TR (Taproot/Bech32m)"},

		// Testnet addresses
		{"mzBc4XEFSdzCDcTxAgf6EZXgsZWpztRhef", "Testnet P2PKH"},
		{"n1ZCYg9YXtB5XrQjMq1YXy2PVJeKkQn3LJ", "Testnet P2PKH"},
		{"2N2JD6wb56AfK4tfmM6PwdVmoYk2dCKf4Br", "Testnet P2SH"},
		{"2MzQwSSnBHWHqSAqtTVQ6v47XtaisrJa1Vc", "Testnet P2SH"},
		{"tb1qw508d6qejxtdg4y5r3zarvary0c5xw7kxpjzsx", "Testnet P2WPKH (Bech32)"},
		{"tb1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3qccfmv3", "Testnet P2WSH (Bech32)"},
		{"tb1p5d7rjq7g6rdk2yhzks9smlaqtedr4dekq08ge8ztwac72sfr9rusxg3297", "Testnet P2TR (Taproot)"},

		// Regtest addresses (bcrt1) - these are often 44 and 64 chars respectively
		{"bcrt1qw508d6qejxtdg4y5r3zarvary0c5xw7kw508d6q", "Unknown (45 chars)"},                        // Invalid length for regtest
		{"bcrt1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3qrp33g0", "Regtest P2WSH (Bech32)"}, // Correct length
		{"bcrt1p5d7rjq7g6rdk2yhzks9smlaqtedr4dekq08ge8ztwac72sfr9rusxg329p0", "Unknown (65 chars)"},    // Wrong length

		// Edge cases
		{"", "empty"},
		{"invalid", "Unknown (7 chars)"},
		{"1", "Unknown (1 chars)"},
		{"bc1", "Unknown (3 chars)"},
		{"toolongaddressstringx", "Unknown (21 chars)"},
	}

	for _, test := range tests {
		result := GetBitcoinAddressType(test.address)
		if result != test.expectedType {
			t.Errorf("GetBitcoinAddressType(%q) = %q, want %q", test.address, result, test.expectedType)
		}
	}
}

func TestFlexibleAddressString(t *testing.T) {
	tests := []struct {
		name     string
		address  FlexibleAddress
		expected string
	}{
		{
			name: "Empty address",
			address: FlexibleAddress{
				Len:  0,
				Data: [63]uint8{},
			},
			expected: "",
		},
		{
			name: "Short address",
			address: FlexibleAddress{
				Len: 5,
				Data: func() [63]uint8 {
					var data [63]uint8
					copy(data[:], []byte("hello"))
					return data
				}(),
			},
			expected: "hello",
		},
		{
			name: "Bitcoin P2PKH address",
			address: FlexibleAddress{
				Len: 34,
				Data: func() [63]uint8 {
					var data [63]uint8
					copy(data[:], []byte("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"))
					return data
				}(),
			},
			expected: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
		},
		{
			name: "Bitcoin bech32 address",
			address: FlexibleAddress{
				Len: 42,
				Data: func() [63]uint8 {
					var data [63]uint8
					copy(data[:], []byte("bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4"))
					return data
				}(),
			},
			expected: "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.address.String()
			if result != test.expected {
				t.Errorf("FlexibleAddress.String() = %q, want %q", result, test.expected)
			}
		})
	}
}

func TestTokensMintedWithFeeGetID(t *testing.T) {
	tests := []struct {
		name     string
		event    TokensMintedWithFee
		expected uint64
	}{
		{
			name: "Zero ID",
			event: TokensMintedWithFee{
				ID: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			expected: 0,
		},
		{
			name: "Small ID",
			event: TokensMintedWithFee{
				ID: [16]uint8{42, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			expected: 42,
		},
		{
			name: "Large ID",
			event: TokensMintedWithFee{
				ID: [16]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			expected: 0xFFFFFFFFFFFFFFFF, // Maximum uint64
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.event.GetID()
			if result != test.expected {
				t.Errorf("TokensMintedWithFee.GetID() = %d, want %d", result, test.expected)
			}
		})
	}
}

func TestZenbtcTokenRedemptionGetID(t *testing.T) {
	tests := []struct {
		name     string
		event    ZenbtcTokenRedemption
		expected uint64
	}{
		{
			name: "Zero ID",
			event: ZenbtcTokenRedemption{
				ID: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			expected: 0,
		},
		{
			name: "Sequential ID",
			event: ZenbtcTokenRedemption{
				ID: [16]uint8{100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			expected: 100,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.event.GetID()
			if result != test.expected {
				t.Errorf("ZenbtcTokenRedemption.GetID() = %d, want %d", result, test.expected)
			}
		})
	}
}

func TestZenbtcTokenRedemptionGetBitcoinAddress(t *testing.T) {
	tests := []struct {
		name     string
		event    ZenbtcTokenRedemption
		expected string
	}{
		{
			name: "Empty address",
			event: ZenbtcTokenRedemption{
				DestAddr: FlexibleAddress{Len: 0, Data: [63]uint8{}},
			},
			expected: "",
		},
		{
			name: "P2PKH address",
			event: ZenbtcTokenRedemption{
				DestAddr: FlexibleAddress{
					Len: 34,
					Data: func() [63]uint8 {
						var data [63]uint8
						copy(data[:], []byte("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"))
						return data
					}(),
				},
			},
			expected: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
		},
		{
			name: "Bech32 address",
			event: ZenbtcTokenRedemption{
				DestAddr: FlexibleAddress{
					Len: 42,
					Data: func() [63]uint8 {
						var data [63]uint8
						copy(data[:], []byte("bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4"))
						return data
					}(),
				},
			},
			expected: "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.event.GetBitcoinAddress()
			if result != test.expected {
				t.Errorf("ZenbtcTokenRedemption.GetBitcoinAddress() = %q, want %q", result, test.expected)
			}
		})
	}
}

func TestRockTokenRedemptionGetID(t *testing.T) {
	tests := []struct {
		name     string
		event    RockTokenRedemption
		expected uint64
	}{
		{
			name: "Zero ID",
			event: RockTokenRedemption{
				ID: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			expected: 0,
		},
		{
			name: "Random ID",
			event: RockTokenRedemption{
				ID: [16]uint8{200, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			expected: 456, // 200 + (1 << 8) = 200 + 256 = 456
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.event.GetID()
			if result != test.expected {
				t.Errorf("RockTokenRedemption.GetID() = %d, want %d", result, test.expected)
			}
		})
	}
}

func TestRockTokenRedemptionGetBitcoinAddress(t *testing.T) {
	tests := []struct {
		name     string
		event    RockTokenRedemption
		expected string
	}{
		{
			name: "P2SH address",
			event: RockTokenRedemption{
				DestAddr: FlexibleAddress{
					Len: 34,
					Data: func() [63]uint8 {
						var data [63]uint8
						copy(data[:], []byte("3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy"))
						return data
					}(),
				},
			},
			expected: "3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy",
		},
		{
			name: "Taproot address",
			event: RockTokenRedemption{
				DestAddr: FlexibleAddress{
					Len: 62,
					Data: func() [63]uint8 {
						var data [63]uint8
						copy(data[:], []byte("bc1p5d7rjq7g6rdk2yhzks9smlaqtedr4dekq08ge8ztwac72sfr9rusxg3297"))
						return data
					}(),
				},
			},
			expected: "bc1p5d7rjq7g6rdk2yhzks9smlaqtedr4dekq08ge8ztwac72sfr9rusxg3297",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.event.GetBitcoinAddress()
			if result != test.expected {
				t.Errorf("RockTokenRedemption.GetBitcoinAddress() = %q, want %q", result, test.expected)
			}
		})
	}
}

func TestGetShardPDA(t *testing.T) {
	// Mock client for testing PDA generation
	client := &Client{
		programID: mustParsePublicKey("4KFjSTnBjbJbWXAiwpWjBCCfAKhjqMp3yfXYpoR3eVis"),
	}

	tests := []struct {
		name       string
		seed       string
		shardIndex uint16
		shouldErr  bool
	}{
		{
			name:       "Valid zenbtc wrap shard",
			seed:       ZENBTC_WRAP_SHARD_SEED,
			shardIndex: 0,
			shouldErr:  false,
		},
		{
			name:       "Valid zenbtc wrap shard max",
			seed:       ZENBTC_WRAP_SHARD_SEED,
			shardIndex: ZENBTC_WRAP_SHARD_COUNT - 1,
			shouldErr:  false,
		},
		{
			name:       "Valid rock unwrap shard",
			seed:       ROCK_UNWRAP_SHARD_SEED,
			shardIndex: 5,
			shouldErr:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pda, err := client.getShardPDA(test.seed, test.shardIndex)

			if test.shouldErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify PDA is not the default (zero) public key
			if pda.IsZero() {
				t.Error("Generated PDA should not be zero")
			}
		})
	}
}

func TestGetAllShardAddresses(t *testing.T) {
	client := &Client{
		programID: mustParsePublicKey("4KFjSTnBjbJbWXAiwpWjBCCfAKhjqMp3yfXYpoR3eVis"),
	}

	tests := []struct {
		name       string
		seed       string
		shardCount uint16
	}{
		{
			name:       "ZenBTC wrap shards",
			seed:       ZENBTC_WRAP_SHARD_SEED,
			shardCount: ZENBTC_WRAP_SHARD_COUNT,
		},
		{
			name:       "Rock unwrap shards",
			seed:       ROCK_UNWRAP_SHARD_SEED,
			shardCount: ROCK_UNWRAP_SHARD_COUNT,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			addresses, err := client.getAllShardAddresses(test.seed, test.shardCount)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(addresses) != int(test.shardCount) {
				t.Errorf("Expected %d addresses, got %d", test.shardCount, len(addresses))
				return
			}

			// Verify all addresses are unique and non-zero
			seen := make(map[string]bool)
			for i, addr := range addresses {
				if addr.IsZero() {
					t.Errorf("Address at index %d is zero", i)
				}

				addrStr := addr.String()
				if seen[addrStr] {
					t.Errorf("Duplicate address found: %s", addrStr)
				}
				seen[addrStr] = true
			}
		})
	}
}

func TestConstants(t *testing.T) {
	// Test that constants match expected values
	if TARGET_WRAP_EVENTS != 1000 {
		t.Errorf("TARGET_WRAP_EVENTS = %d, want 1000", TARGET_WRAP_EVENTS)
	}
	if TARGET_UNWRAP_EVENTS != 1020 {
		t.Errorf("TARGET_UNWRAP_EVENTS = %d, want 1020", TARGET_UNWRAP_EVENTS)
	}

	if SHARD_SIZE_WRAP != 100 {
		t.Errorf("SHARD_SIZE_WRAP = %d, want 100", SHARD_SIZE_WRAP)
	}

	if SHARD_SIZE_UNWRAP != 60 {
		t.Errorf("SHARD_SIZE_UNWRAP = %d, want 60", SHARD_SIZE_UNWRAP)
	}

	if ZENBTC_WRAP_SHARD_COUNT != 10 {
		t.Errorf("ZENBTC_WRAP_SHARD_COUNT = %d, want 10", ZENBTC_WRAP_SHARD_COUNT)
	}

	if ZENBTC_UNWRAP_SHARD_COUNT != 17 {
		t.Errorf("ZENBTC_UNWRAP_SHARD_COUNT = %d, want 17", ZENBTC_UNWRAP_SHARD_COUNT)
	}

	if ROCK_WRAP_SHARD_COUNT != 10 {
		t.Errorf("ROCK_WRAP_SHARD_COUNT = %d, want 10", ROCK_WRAP_SHARD_COUNT)
	}

	if ROCK_UNWRAP_SHARD_COUNT != 17 {
		t.Errorf("ROCK_UNWRAP_SHARD_COUNT = %d, want 17", ROCK_UNWRAP_SHARD_COUNT)
	}

	// Test capacity calculations
	zenbtcWrapCapacity := ZENBTC_WRAP_SHARD_COUNT * SHARD_SIZE_WRAP
	zenbtcUnwrapCapacity := ZENBTC_UNWRAP_SHARD_COUNT * SHARD_SIZE_UNWRAP
	rockWrapCapacity := ROCK_WRAP_SHARD_COUNT * SHARD_SIZE_WRAP
	rockUnwrapCapacity := ROCK_UNWRAP_SHARD_COUNT * SHARD_SIZE_UNWRAP

	if zenbtcWrapCapacity != 1000 {
		t.Errorf("ZenBTC wrap capacity = %d, want 1000", zenbtcWrapCapacity)
	}

	if zenbtcUnwrapCapacity != 1020 {
		t.Errorf("ZenBTC unwrap capacity = %d, want 1020", zenbtcUnwrapCapacity)
	}

	if rockWrapCapacity != 1000 {
		t.Errorf("Rock wrap capacity = %d, want 1000", rockWrapCapacity)
	}

	if rockUnwrapCapacity != 1020 {
		t.Errorf("Rock unwrap capacity = %d, want 1020", rockUnwrapCapacity)
	}

	totalCapacity := zenbtcWrapCapacity + zenbtcUnwrapCapacity + rockWrapCapacity + rockUnwrapCapacity
	if totalCapacity != 4040 {
		t.Errorf("Total capacity = %d, want 4040", totalCapacity)
	}
}

func TestNewClient(t *testing.T) {
	// Test with nil program ID (should use default)
	client1 := NewClient(nil, nil)
	if client1 == nil {
		t.Error("NewClient returned nil")
	}

	expectedProgramID := mustParsePublicKey(DEFAULT_PROGRAM_ID)
	if !client1.programID.Equals(expectedProgramID) {
		t.Errorf("Default program ID = %s, want %s", client1.programID, expectedProgramID)
	}

	// Test with custom program ID
	customProgramID := mustParsePublicKey("11111111111111111111111111111112")
	client2 := NewClient(nil, &customProgramID)
	if !client2.programID.Equals(customProgramID) {
		t.Errorf("Custom program ID = %s, want %s", client2.programID, customProgramID)
	}
}

// Helper function to parse public keys for testing
func mustParsePublicKey(s string) solana.PublicKey {
	pk, err := solana.PublicKeyFromBase58(s)
	if err != nil {
		panic(err)
	}
	return pk
}

// Benchmark tests
func BenchmarkGetBitcoinAddressType(b *testing.B) {
	addresses := []string{
		"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
		"3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy",
		"bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4",
		"bc1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3qccfmv3",
		"bc1p5d7rjq7g6rdk2yhzks9smlaqtedr4dekq08ge8ztwac72sfr9rusxg3297",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		addr := addresses[i%len(addresses)]
		GetBitcoinAddressType(addr)
	}
}

func BenchmarkFlexibleAddressString(b *testing.B) {
	addr := FlexibleAddress{
		Len: 42,
		Data: func() [63]uint8 {
			var data [63]uint8
			copy(data[:], []byte("bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4"))
			return data
		}(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = addr.String()
	}
}

// Test ring buffer behavior
func TestRingBufferLogic(t *testing.T) {
	// Test that the constants match expected ring buffer capacities
	zenbtcWrapCapacity := ZENBTC_WRAP_SHARD_COUNT * SHARD_SIZE_WRAP
	zenbtcUnwrapCapacity := ZENBTC_UNWRAP_SHARD_COUNT * SHARD_SIZE_UNWRAP
	rockWrapCapacity := ROCK_WRAP_SHARD_COUNT * SHARD_SIZE_WRAP
	rockUnwrapCapacity := ROCK_UNWRAP_SHARD_COUNT * SHARD_SIZE_UNWRAP

	if zenbtcWrapCapacity != 1000 {
		t.Errorf("ZenBTC wrap capacity = %d, want 1000", zenbtcWrapCapacity)
	}
	if zenbtcUnwrapCapacity != 1020 {
		t.Errorf("ZenBTC unwrap capacity = %d, want 1020", zenbtcUnwrapCapacity)
	}
	if rockWrapCapacity != 1000 {
		t.Errorf("Rock wrap capacity = %d, want 1000", rockWrapCapacity)
	}
	if rockUnwrapCapacity != 1020 {
		t.Errorf("Rock unwrap capacity = %d, want 1020", rockUnwrapCapacity)
	}

	totalCapacity := zenbtcWrapCapacity + zenbtcUnwrapCapacity + rockWrapCapacity + rockUnwrapCapacity
	if totalCapacity != 4040 {
		t.Errorf("Total capacity = %d, want 4040", totalCapacity)
	}
}

// Test event ID wraparound behavior
func TestEventIDWraparound(t *testing.T) {
	// Test that large event IDs don't cause issues
	largeEventIDs := []uint64{
		1000000,
		uint64(1) << 32,
		uint64(1) << 48,
		^uint64(0), // Maximum uint64
	}

	for _, eventID := range largeEventIDs {
		event := TokensMintedWithFee{
			ID: func() [16]uint8 {
				var id [16]uint8
				// Set the lower 64 bits
				for i := 0; i < 8; i++ {
					id[i] = uint8(eventID >> (i * 8))
				}
				return id
			}(),
		}

		result := event.GetID()
		if result != eventID {
			t.Errorf("GetID() = %d, want %d for large event ID", result, eventID)
		}
	}
}

// Test Bitcoin address edge cases
func TestBitcoinAddressEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		address  string
		expected string
	}{
		// Length boundary tests
		{"Short legacy", "1", "Unknown (1 chars)"},
		{"Long legacy", "1" + strings.Repeat("A", 40), "Unknown (41 chars)"},
		{"Short P2SH", "3", "Unknown (1 chars)"},
		{"Long P2SH", "3" + strings.Repeat("A", 40), "Unknown (41 chars)"},

		// Bech32 length variations
		{"Short bech32", "bc1q", "Unknown (4 chars)"},
		{"Medium bech32", "bc1q" + strings.Repeat("a", 30), "Unknown (34 chars)"},
		{"Long bech32", "bc1q" + strings.Repeat("a", 70), "Unknown (74 chars)"},

		// Case sensitivity
		{"Uppercase legacy", "1A1ZP1EP5QGEFI2DMPTFTL5SLMV7DIVFNA", "P2PKH (Legacy)"},
		{"Mixed case bech32", "BC1QW508D6QEJXTDG4Y5R3ZARVARY0C5XW7KV8F3T4", "Unknown (42 chars)"},

		// Invalid characters
		{"Invalid chars", "bc1q!@#$%^&*()", "Unknown (14 chars)"},
		{"Numbers only", "123456789", "Unknown (9 chars)"},

		// Boundary conditions for regtest
		{"Regtest too short", "bcrt1q", "Unknown (6 chars)"},
		{"Regtest wrong length", "bcrt1q" + strings.Repeat("a", 30), "Unknown (36 chars)"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := GetBitcoinAddressType(test.address)
			if result != test.expected {
				t.Errorf("GetBitcoinAddressType(%q) = %q, want %q", test.address, result, test.expected)
			}
		})
	}
}

// Test FlexibleAddress boundary conditions
func TestFlexibleAddressBoundary(t *testing.T) {
	tests := []struct {
		name string
		len  uint8
		data string
		want string
	}{
		{"Zero length", 0, "", ""},
		{"Single char", 1, "a", "a"},
		{"Max length", 63, strings.Repeat("x", 63), strings.Repeat("x", 63)},
		{"Length mismatch", 5, strings.Repeat("y", 10), "yyyyy"}, // Should only use first 5 chars
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			addr := FlexibleAddress{
				Len: test.len,
				Data: func() [63]uint8 {
					var data [63]uint8
					copy(data[:], []byte(test.data))
					return data
				}(),
			}

			result := addr.String()
			if result != test.want {
				t.Errorf("FlexibleAddress.String() = %q, want %q", result, test.want)
			}
		})
	}
}
