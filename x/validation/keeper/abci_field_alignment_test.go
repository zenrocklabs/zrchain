package keeper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestFieldHandlerAlignment ensures that the field handlers and enum constants are properly aligned
// This test prevents nil map panics by catching misalignments at compile/test time
func TestFieldHandlerAlignment(t *testing.T) {
	handlers := initializeFieldHandlers()

	t.Run("all_handlers_have_unique_fields", func(t *testing.T) {
		fieldsSeen := make(map[VoteExtensionField]bool)
		for _, handler := range handlers {
			require.False(t, fieldsSeen[handler.Field],
				"Duplicate field in handlers: %s", handler.Field.String())
			fieldsSeen[handler.Field] = true
		}
	})

	t.Run("all_handlers_have_valid_string_representation", func(t *testing.T) {
		for _, handler := range handlers {
			str := handler.Field.String()
			require.NotEqual(t, "Unknown", str,
				"Handler field %d has no String() representation", handler.Field)
		}
	})

	t.Run("all_handlers_have_GetValue_function", func(t *testing.T) {
		for _, handler := range handlers {
			require.NotNil(t, handler.GetValue,
				"Handler for field %s is missing GetValue function", handler.Field.String())
		}
	})

	t.Run("all_handlers_have_SetValue_function", func(t *testing.T) {
		for _, handler := range handlers {
			require.NotNil(t, handler.SetValue,
				"Handler for field %s is missing SetValue function", handler.Field.String())
		}
	})

	t.Run("all_enum_fields_have_handlers", func(t *testing.T) {
		// Create a map of all fields that have handlers
		handlerFields := make(map[VoteExtensionField]bool)
		for _, handler := range handlers {
			handlerFields[handler.Field] = true
		}

		// Check all known enum constants (from VEFieldEigenDelegationsHash to VEFieldLatestZcashHeaderHash)
		// We iterate through the enum range and verify each has a handler
		enumFields := []VoteExtensionField{
			VEFieldEigenDelegationsHash,
			VEFieldEthBurnEventsHash,
			VEFieldSolanaBurnEventsHash,
			VEFieldRedemptionsHash,
			VEFieldRequestedBtcHeaderHash,
			VEFieldRequestedBtcBlockHeight,
			VEFieldEthBlockHeight,
			VEFieldEthGasLimit,
			VEFieldEthBaseFee,
			VEFieldEthTipCap,
			VEFieldRequestedStakerNonce,
			VEFieldRequestedEthMinterNonce,
			VEFieldRequestedUnstakerNonce,
			VEFieldRequestedCompleterNonce,
			VEFieldROCKUSDPrice,
			VEFieldBTCUSDPrice,
			VEFieldETHUSDPrice,
			VEFieldLatestBtcBlockHeight,
			VEFieldLatestBtcHeaderHash,
			VEFieldSolanaMintNoncesHash,
			VEFieldSolanaAccountsHash,
			VEFieldSolanaMintEventsHash,
			VEFieldRequestedZcashBlockHeight,
			VEFieldRequestedZcashHeaderHash,
			VEFieldLatestZcashBlockHeight,
			VEFieldLatestZcashHeaderHash,
		}

		for _, field := range enumFields {
			require.True(t, handlerFields[field],
				"Enum field %s (%d) does not have a corresponding handler", field.String(), field)
		}

		// Verify the count matches
		require.Equal(t, len(enumFields), len(handlers),
			"Number of enum fields (%d) does not match number of handlers (%d)",
			len(enumFields), len(handlers))
	})
}

// TestHandlerGetSetSymmetry ensures that GetValue and SetValue are symmetric operations
func TestHandlerGetSetSymmetry(t *testing.T) {
	handlers := initializeFieldHandlers()

	testCases := map[VoteExtensionField]any{
		VEFieldEigenDelegationsHash:       []byte("test_hash_1"),
		VEFieldEthBurnEventsHash:          []byte("test_hash_2"),
		VEFieldSolanaBurnEventsHash:       []byte("test_hash_3"),
		VEFieldRedemptionsHash:            []byte("test_hash_4"),
		VEFieldRequestedBtcHeaderHash:     []byte("test_hash_5"),
		VEFieldLatestBtcHeaderHash:        []byte("test_hash_6"),
		VEFieldSolanaMintNoncesHash:       []byte("test_hash_7"),
		VEFieldSolanaAccountsHash:         []byte("test_hash_8"),
		VEFieldSolanaMintEventsHash:       []byte("test_hash_9"),
		VEFieldRequestedBtcBlockHeight:    int64(12345),
		VEFieldEthBlockHeight:             uint64(67890),
		VEFieldEthGasLimit:                uint64(21000),
		VEFieldEthBaseFee:                 uint64(20000000000),
		VEFieldEthTipCap:                  uint64(1000000000),
		VEFieldRequestedStakerNonce:       uint64(100),
		VEFieldRequestedEthMinterNonce:    uint64(200),
		VEFieldRequestedUnstakerNonce:     uint64(300),
		VEFieldRequestedCompleterNonce:    uint64(400),
		VEFieldLatestBtcBlockHeight:       int64(800000),
		VEFieldROCKUSDPrice:               "1.25",
		VEFieldBTCUSDPrice:                "45000.00",
		VEFieldETHUSDPrice:                "2800.00",
		VEFieldRequestedZcashBlockHeight:  int64(1000000),
		VEFieldRequestedZcashHeaderHash:   []byte("zcash_hash_1"),
		VEFieldLatestZcashBlockHeight:     int64(2000000),
		VEFieldLatestZcashHeaderHash:      []byte("zcash_hash_2"),
	}

	for _, handler := range handlers {
		testValue, ok := testCases[handler.Field]
		require.True(t, ok, "Missing test value for field %s", handler.Field.String())

		t.Run(handler.Field.String(), func(t *testing.T) {
			var ve VoteExtension

			// Set the value
			handler.SetValue(testValue, &ve)

			// Get the value back
			retrievedValue := handler.GetValue(ve)

			// Verify they match
			switch expected := testValue.(type) {
			case []byte:
				retrieved, ok := retrievedValue.([]byte)
				require.True(t, ok, "GetValue returned wrong type for %s", handler.Field.String())
				require.Equal(t, expected, retrieved, "GetValue/SetValue mismatch for %s", handler.Field.String())
			default:
				require.Equal(t, testValue, retrievedValue, "GetValue/SetValue mismatch for %s", handler.Field.String())
			}
		})
	}
}

// TestFieldVoteMapInitialization simulates the initialization pattern in GetConsensusAndPluralityVEData
// to ensure no nil map panics can occur
func TestFieldVoteMapInitialization(t *testing.T) {
	t.Run("field_handler_based_initialization_covers_all_access_patterns", func(t *testing.T) {
		// Get field handlers (this is what GetConsensusAndPluralityVEData does first)
		fieldHandlers := initializeFieldHandlers()

		// Initialize fieldVotes map using the same pattern as the fixed code
		fieldVotes := make(map[VoteExtensionField]map[string]fieldVote)
		for _, handler := range fieldHandlers {
			fieldVotes[handler.Field] = make(map[string]fieldVote)
		}

		// Simulate accessing the maps for each handler (what happens in the vote processing loop)
		for _, handler := range fieldHandlers {
			votes := fieldVotes[handler.Field]
			require.NotNil(t, votes,
				"Map for field %s was not initialized", handler.Field.String())

			// Simulate writing to the map (this would panic if the map was nil)
			key := "test_key"
			votes[key] = fieldVote{
				value:     "test_value",
				votePower: 1000,
			}

			// Verify the write succeeded
			require.Contains(t, votes, key)
			require.Equal(t, int64(1000), votes[key].votePower)
		}

		// Simulate the consensus processing loop (another access pattern)
		for _, handler := range fieldHandlers {
			votes := fieldVotes[handler.Field]
			require.NotNil(t, votes,
				"Map for field %s was not initialized in consensus loop", handler.Field.String())

			// This simulates the len(votes) check in the actual code
			_ = len(votes)
		}
	})
}

// TestInitializeFieldHandlersDeterministic ensures the function always returns the same handlers
func TestInitializeFieldHandlersDeterministic(t *testing.T) {
	// Call the function multiple times
	handlers1 := initializeFieldHandlers()
	handlers2 := initializeFieldHandlers()
	handlers3 := initializeFieldHandlers()

	// Verify they all have the same length
	require.Equal(t, len(handlers1), len(handlers2),
		"initializeFieldHandlers returned different lengths on subsequent calls")
	require.Equal(t, len(handlers1), len(handlers3),
		"initializeFieldHandlers returned different lengths on subsequent calls")

	// Verify they all have the same fields in the same order
	for i := range handlers1 {
		require.Equal(t, handlers1[i].Field, handlers2[i].Field,
			"Field at index %d differs between calls", i)
		require.Equal(t, handlers1[i].Field, handlers3[i].Field,
			"Field at index %d differs between calls", i)
	}
}

// TestVoteExtensionStructCompleteness ensures the VoteExtension struct has all necessary fields
func TestVoteExtensionStructCompleteness(t *testing.T) {
	handlers := initializeFieldHandlers()

	t.Run("all_handlers_can_get_and_set_values", func(t *testing.T) {
		var ve VoteExtension
		veType := reflect.TypeOf(ve)

		for _, handler := range handlers {
			// Try to set a value (this will panic if the struct field doesn't exist)
			require.NotPanics(t, func() {
				switch handler.Field {
				case VEFieldEigenDelegationsHash, VEFieldEthBurnEventsHash, VEFieldSolanaBurnEventsHash,
					VEFieldRedemptionsHash, VEFieldRequestedBtcHeaderHash, VEFieldLatestBtcHeaderHash,
					VEFieldSolanaMintNoncesHash, VEFieldSolanaAccountsHash, VEFieldSolanaMintEventsHash,
					VEFieldRequestedZcashHeaderHash, VEFieldLatestZcashHeaderHash:
					handler.SetValue([]byte("test"), &ve)
				case VEFieldRequestedBtcBlockHeight, VEFieldLatestBtcBlockHeight,
					VEFieldRequestedZcashBlockHeight, VEFieldLatestZcashBlockHeight:
					handler.SetValue(int64(123), &ve)
				case VEFieldEthBlockHeight, VEFieldEthGasLimit, VEFieldEthBaseFee, VEFieldEthTipCap,
					VEFieldRequestedStakerNonce, VEFieldRequestedEthMinterNonce,
					VEFieldRequestedUnstakerNonce, VEFieldRequestedCompleterNonce:
					handler.SetValue(uint64(456), &ve)
				case VEFieldROCKUSDPrice, VEFieldBTCUSDPrice, VEFieldETHUSDPrice:
					handler.SetValue("1.23", &ve)
				}
			}, "SetValue panicked for field %s", handler.Field.String())

			// Try to get a value (this will panic if the struct field doesn't exist)
			require.NotPanics(t, func() {
				_ = handler.GetValue(ve)
			}, "GetValue panicked for field %s", handler.Field.String())
		}

		// Log the number of fields for documentation
		t.Logf("VoteExtension struct has %d fields", veType.NumField())
		t.Logf("initializeFieldHandlers returns %d handlers", len(handlers))
	})

	t.Run("all_struct_fields_are_handled_or_explicitly_excluded", func(t *testing.T) {
		var ve VoteExtension
		veType := reflect.TypeOf(ve)

		// Fields that are informational/metadata and don't require consensus voting
		excludedFields := map[string]string{
			"SidecarVersionName": "Informational field that doesn't require consensus",
		}

		// Build a map of struct field names that have handlers
		handlerFieldNames := make(map[string]bool)
		for _, handler := range handlers {
			// Get the actual struct field name by using GetValue on a zero VoteExtension
			// and checking which field it accesses
			fieldName := getStructFieldNameForHandler(handler)
			if fieldName != "" {
				handlerFieldNames[fieldName] = true
			}
		}

		// Check every struct field
		var uncoveredFields []string
		for i := 0; i < veType.NumField(); i++ {
			field := veType.Field(i)
			fieldName := field.Name

			// Skip if it's explicitly excluded
			if reason, excluded := excludedFields[fieldName]; excluded {
				t.Logf("Field %s is excluded: %s", fieldName, reason)
				continue
			}

			// Check if it has a handler
			if !handlerFieldNames[fieldName] {
				uncoveredFields = append(uncoveredFields, fieldName)
			}
		}

		require.Empty(t, uncoveredFields,
			"The following struct fields are neither handled nor explicitly excluded: %v. "+
				"Either add handlers for them or add them to the excludedFields map with a reason.",
			uncoveredFields)

		// Calculate expected counts
		totalStructFields := veType.NumField()
		excludedCount := len(excludedFields)
		expectedHandlers := totalStructFields - excludedCount

		t.Logf("VoteExtension struct: %d total fields", totalStructFields)
		t.Logf("Excluded fields: %d", excludedCount)
		t.Logf("Expected handlers: %d", expectedHandlers)
		t.Logf("Actual handlers: %d", len(handlers))

		require.Equal(t, expectedHandlers, len(handlers),
			"Number of handlers (%d) should equal total struct fields (%d) minus excluded fields (%d) = %d",
			len(handlers), totalStructFields, excludedCount, expectedHandlers)
	})
}

// getStructFieldNameForHandler attempts to determine the struct field name for a handler
func getStructFieldNameForHandler(handler FieldHandler) string {
	// Map of VoteExtensionField enum to struct field name
	fieldMap := map[VoteExtensionField]string{
		VEFieldEigenDelegationsHash:       "EigenDelegationsHash",
		VEFieldEthBurnEventsHash:          "EthBurnEventsHash",
		VEFieldSolanaBurnEventsHash:       "SolanaBurnEventsHash",
		VEFieldRedemptionsHash:            "RedemptionsHash",
		VEFieldRequestedBtcHeaderHash:     "RequestedBtcHeaderHash",
		VEFieldRequestedBtcBlockHeight:    "RequestedBtcBlockHeight",
		VEFieldEthBlockHeight:             "EthBlockHeight",
		VEFieldEthGasLimit:                "EthGasLimit",
		VEFieldEthBaseFee:                 "EthBaseFee",
		VEFieldEthTipCap:                  "EthTipCap",
		VEFieldRequestedStakerNonce:       "RequestedStakerNonce",
		VEFieldRequestedEthMinterNonce:    "RequestedEthMinterNonce",
		VEFieldRequestedUnstakerNonce:     "RequestedUnstakerNonce",
		VEFieldRequestedCompleterNonce:    "RequestedCompleterNonce",
		VEFieldSolanaMintNoncesHash:       "SolanaMintNoncesHash",
		VEFieldSolanaAccountsHash:         "SolanaAccountsHash",
		VEFieldSolanaMintEventsHash:       "SolanaMintEventsHash",
		VEFieldROCKUSDPrice:               "ROCKUSDPrice",
		VEFieldBTCUSDPrice:                "BTCUSDPrice",
		VEFieldETHUSDPrice:                "ETHUSDPrice",
		VEFieldLatestBtcBlockHeight:       "LatestBtcBlockHeight",
		VEFieldLatestBtcHeaderHash:        "LatestBtcHeaderHash",
		VEFieldRequestedZcashBlockHeight:  "RequestedZcashBlockHeight",
		VEFieldRequestedZcashHeaderHash:   "RequestedZcashHeaderHash",
		VEFieldLatestZcashBlockHeight:     "LatestZcashBlockHeight",
		VEFieldLatestZcashHeaderHash:      "LatestZcashHeaderHash",
	}

	return fieldMap[handler.Field]
}

// TestEnumConstantCountMatchesHandlers ensures the number of VoteExtensionField enum constants
// matches the number of handlers exactly. This prevents someone from adding an enum constant
// but forgetting to add the corresponding handler.
func TestEnumConstantCountMatchesHandlers(t *testing.T) {
	handlers := initializeFieldHandlers()

	// All VoteExtensionField enum constants - this list MUST be updated when new fields are added
	allEnumConstants := []VoteExtensionField{
		VEFieldEigenDelegationsHash,
		VEFieldEthBurnEventsHash,
		VEFieldSolanaBurnEventsHash,
		VEFieldRedemptionsHash,
		VEFieldRequestedBtcHeaderHash,
		VEFieldRequestedBtcBlockHeight,
		VEFieldEthBlockHeight,
		VEFieldEthGasLimit,
		VEFieldEthBaseFee,
		VEFieldEthTipCap,
		VEFieldRequestedStakerNonce,
		VEFieldRequestedEthMinterNonce,
		VEFieldRequestedUnstakerNonce,
		VEFieldRequestedCompleterNonce,
		VEFieldROCKUSDPrice,
		VEFieldBTCUSDPrice,
		VEFieldETHUSDPrice,
		VEFieldLatestBtcBlockHeight,
		VEFieldLatestBtcHeaderHash,
		VEFieldSolanaMintNoncesHash,
		VEFieldSolanaAccountsHash,
		VEFieldSolanaMintEventsHash,
		VEFieldRequestedZcashBlockHeight,
		VEFieldRequestedZcashHeaderHash,
		VEFieldLatestZcashBlockHeight,
		VEFieldLatestZcashHeaderHash,
	}

	t.Run("enum_count_equals_handler_count", func(t *testing.T) {
		require.Equal(t, len(allEnumConstants), len(handlers),
			"Number of enum constants (%d) must equal number of handlers (%d). "+
				"If you added a new enum constant, you must also add it to initializeFieldHandlers() "+
				"and update the allEnumConstants list in this test.",
			len(allEnumConstants), len(handlers))
	})

	t.Run("every_enum_constant_has_handler", func(t *testing.T) {
		handlerFields := make(map[VoteExtensionField]bool)
		for _, handler := range handlers {
			handlerFields[handler.Field] = true
		}

		var missingHandlers []VoteExtensionField
		for _, enumConstant := range allEnumConstants {
			if !handlerFields[enumConstant] {
				missingHandlers = append(missingHandlers, enumConstant)
			}
		}

		require.Empty(t, missingHandlers,
			"The following enum constants do not have handlers: %v. "+
				"Add them to initializeFieldHandlers().",
			missingHandlers)
	})

	t.Run("every_handler_has_enum_constant", func(t *testing.T) {
		enumConstants := make(map[VoteExtensionField]bool)
		for _, constant := range allEnumConstants {
			enumConstants[constant] = true
		}

		var unknownHandlers []VoteExtensionField
		for _, handler := range handlers {
			if !enumConstants[handler.Field] {
				unknownHandlers = append(unknownHandlers, handler.Field)
			}
		}

		require.Empty(t, unknownHandlers,
			"The following handlers reference enum constants not in allEnumConstants list: %v. "+
				"Add them to the allEnumConstants list in this test.",
			unknownHandlers)
	})

	t.Run("all_enum_constants_are_unique", func(t *testing.T) {
		seen := make(map[VoteExtensionField]bool)
		var duplicates []VoteExtensionField

		for _, constant := range allEnumConstants {
			if seen[constant] {
				duplicates = append(duplicates, constant)
			}
			seen[constant] = true
		}

		require.Empty(t, duplicates,
			"Duplicate enum constants found in allEnumConstants list: %v",
			duplicates)
	})

	t.Logf("✅ Validation complete:")
	t.Logf("  - Enum constants: %d", len(allEnumConstants))
	t.Logf("  - Handlers: %d", len(handlers))
	t.Logf("  - All counts match and align perfectly")
}

// TestGenericGetKeyConsistency ensures the genericGetKey function is consistent
func TestGenericGetKeyConsistency(t *testing.T) {
	testCases := []struct {
		name     string
		value    any
		expected string
	}{
		{
			name:     "nil_value",
			value:    nil,
			expected: "",
		},
		{
			name:     "byte_slice",
			value:    []byte{0x01, 0x02, 0x03},
			expected: "010203",
		},
		{
			name:     "empty_byte_slice",
			value:    []byte{},
			expected: "",
		},
		{
			name:     "int64",
			value:    int64(12345),
			expected: "12345",
		},
		{
			name:     "uint64",
			value:    uint64(67890),
			expected: "67890",
		},
		{
			name:     "string",
			value:    "test_string",
			expected: `"test_string"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call genericGetKey multiple times to ensure consistency
			key1 := genericGetKey(tc.value)
			key2 := genericGetKey(tc.value)
			key3 := genericGetKey(tc.value)

			require.Equal(t, tc.expected, key1, "genericGetKey returned unexpected value")
			require.Equal(t, key1, key2, "genericGetKey is not consistent")
			require.Equal(t, key1, key3, "genericGetKey is not consistent")
		})
	}
}

// TestCompleteThreeWayAlignment is the ultimate validation that ensures:
// 1. Every VoteExtension struct field (except excluded ones) has an enum constant
// 2. Every enum constant has a handler
// 3. Every handler references a valid struct field
// This test ensures perfect alignment across all three components.
func TestCompleteThreeWayAlignment(t *testing.T) {
	var ve VoteExtension
	veType := reflect.TypeOf(ve)
	handlers := initializeFieldHandlers()

	// Fields that don't participate in consensus voting
	excludedStructFields := map[string]string{
		"SidecarVersionName": "Informational field - no consensus required",
	}

	// All enum constants
	allEnumConstants := []VoteExtensionField{
		VEFieldEigenDelegationsHash,
		VEFieldEthBurnEventsHash,
		VEFieldSolanaBurnEventsHash,
		VEFieldRedemptionsHash,
		VEFieldRequestedBtcHeaderHash,
		VEFieldRequestedBtcBlockHeight,
		VEFieldEthBlockHeight,
		VEFieldEthGasLimit,
		VEFieldEthBaseFee,
		VEFieldEthTipCap,
		VEFieldRequestedStakerNonce,
		VEFieldRequestedEthMinterNonce,
		VEFieldRequestedUnstakerNonce,
		VEFieldRequestedCompleterNonce,
		VEFieldROCKUSDPrice,
		VEFieldBTCUSDPrice,
		VEFieldETHUSDPrice,
		VEFieldLatestBtcBlockHeight,
		VEFieldLatestBtcHeaderHash,
		VEFieldSolanaMintNoncesHash,
		VEFieldSolanaAccountsHash,
		VEFieldSolanaMintEventsHash,
		VEFieldRequestedZcashBlockHeight,
		VEFieldRequestedZcashHeaderHash,
		VEFieldLatestZcashBlockHeight,
		VEFieldLatestZcashHeaderHash,
	}

	// Enum to struct field mapping
	enumToStructField := map[VoteExtensionField]string{
		VEFieldEigenDelegationsHash:       "EigenDelegationsHash",
		VEFieldEthBurnEventsHash:          "EthBurnEventsHash",
		VEFieldSolanaBurnEventsHash:       "SolanaBurnEventsHash",
		VEFieldRedemptionsHash:            "RedemptionsHash",
		VEFieldRequestedBtcHeaderHash:     "RequestedBtcHeaderHash",
		VEFieldRequestedBtcBlockHeight:    "RequestedBtcBlockHeight",
		VEFieldEthBlockHeight:             "EthBlockHeight",
		VEFieldEthGasLimit:                "EthGasLimit",
		VEFieldEthBaseFee:                 "EthBaseFee",
		VEFieldEthTipCap:                  "EthTipCap",
		VEFieldRequestedStakerNonce:       "RequestedStakerNonce",
		VEFieldRequestedEthMinterNonce:    "RequestedEthMinterNonce",
		VEFieldRequestedUnstakerNonce:     "RequestedUnstakerNonce",
		VEFieldRequestedCompleterNonce:    "RequestedCompleterNonce",
		VEFieldSolanaMintNoncesHash:       "SolanaMintNoncesHash",
		VEFieldSolanaAccountsHash:         "SolanaAccountsHash",
		VEFieldSolanaMintEventsHash:       "SolanaMintEventsHash",
		VEFieldROCKUSDPrice:               "ROCKUSDPrice",
		VEFieldBTCUSDPrice:                "BTCUSDPrice",
		VEFieldETHUSDPrice:                "ETHUSDPrice",
		VEFieldLatestBtcBlockHeight:       "LatestBtcBlockHeight",
		VEFieldLatestBtcHeaderHash:        "LatestBtcHeaderHash",
		VEFieldRequestedZcashBlockHeight:  "RequestedZcashBlockHeight",
		VEFieldRequestedZcashHeaderHash:   "RequestedZcashHeaderHash",
		VEFieldLatestZcashBlockHeight:     "LatestZcashBlockHeight",
		VEFieldLatestZcashHeaderHash:      "LatestZcashHeaderHash",
	}

	t.Run("counts_align_perfectly", func(t *testing.T) {
		totalStructFields := veType.NumField()
		excludedCount := len(excludedStructFields)
		enumCount := len(allEnumConstants)
		handlerCount := len(handlers)
		expectedCount := totalStructFields - excludedCount

		t.Logf("📊 Alignment Analysis:")
		t.Logf("  Total VoteExtension struct fields: %d", totalStructFields)
		t.Logf("  Excluded fields (no consensus): %d", excludedCount)
		t.Logf("  Fields requiring handlers: %d", expectedCount)
		t.Logf("  Enum constants defined: %d", enumCount)
		t.Logf("  Handlers initialized: %d", handlerCount)

		require.Equal(t, expectedCount, enumCount,
			"Enum count (%d) should equal struct fields (%d) minus excluded (%d) = %d",
			enumCount, totalStructFields, excludedCount, expectedCount)

		require.Equal(t, expectedCount, handlerCount,
			"Handler count (%d) should equal struct fields (%d) minus excluded (%d) = %d",
			handlerCount, totalStructFields, excludedCount, expectedCount)

		require.Equal(t, enumCount, handlerCount,
			"Enum count (%d) must equal handler count (%d)",
			enumCount, handlerCount)
	})

	t.Run("every_non_excluded_struct_field_has_enum_and_handler", func(t *testing.T) {
		// Build reverse map: struct field name -> enum constant
		structFieldToEnum := make(map[string]VoteExtensionField)
		for enum, fieldName := range enumToStructField {
			structFieldToEnum[fieldName] = enum
		}

		// Build map: enum constant -> has handler
		enumHasHandler := make(map[VoteExtensionField]bool)
		for _, handler := range handlers {
			enumHasHandler[handler.Field] = true
		}

		var missingEnum []string
		var missingHandler []string

		for i := 0; i < veType.NumField(); i++ {
			field := veType.Field(i)
			fieldName := field.Name

			// Skip excluded fields
			if _, excluded := excludedStructFields[fieldName]; excluded {
				continue
			}

			// Check if struct field has enum constant
			enumConstant, hasEnum := structFieldToEnum[fieldName]
			if !hasEnum {
				missingEnum = append(missingEnum, fieldName)
				continue
			}

			// Check if enum constant has handler
			if !enumHasHandler[enumConstant] {
				missingHandler = append(missingHandler, fieldName)
			}
		}

		require.Empty(t, missingEnum,
			"Struct fields without enum constants: %v. Add enum constants in VoteExtensionField.",
			missingEnum)

		require.Empty(t, missingHandler,
			"Struct fields with enum but no handler: %v. Add handlers in initializeFieldHandlers().",
			missingHandler)
	})

	t.Run("every_enum_has_struct_field_and_handler", func(t *testing.T) {
		// Build map: enum -> has handler
		enumHasHandler := make(map[VoteExtensionField]bool)
		for _, handler := range handlers {
			enumHasHandler[handler.Field] = true
		}

		var missingStructField []VoteExtensionField
		var missingHandler []VoteExtensionField

		for _, enumConstant := range allEnumConstants {
			// Check if enum has struct field mapping
			structFieldName, hasStructField := enumToStructField[enumConstant]
			if !hasStructField {
				missingStructField = append(missingStructField, enumConstant)
				continue
			}

			// Verify struct field actually exists
			_, found := veType.FieldByName(structFieldName)
			if !found {
				missingStructField = append(missingStructField, enumConstant)
				continue
			}

			// Check if enum has handler
			if !enumHasHandler[enumConstant] {
				missingHandler = append(missingHandler, enumConstant)
			}
		}

		require.Empty(t, missingStructField,
			"Enum constants without valid struct fields: %v. "+
				"Add struct fields to VoteExtension or fix enumToStructField mapping.",
			missingStructField)

		require.Empty(t, missingHandler,
			"Enum constants without handlers: %v. Add handlers in initializeFieldHandlers().",
			missingHandler)
	})

	t.Run("every_handler_has_enum_and_struct_field", func(t *testing.T) {
		var invalidHandlers []VoteExtensionField

		for _, handler := range handlers {
			// Check if handler's field is in enum constants
			enumFound := false
			for _, enumConstant := range allEnumConstants {
				if handler.Field == enumConstant {
					enumFound = true
					break
				}
			}

			if !enumFound {
				invalidHandlers = append(invalidHandlers, handler.Field)
				continue
			}

			// Check if enum maps to a struct field
			structFieldName, hasMaping := enumToStructField[handler.Field]
			if !hasMaping {
				invalidHandlers = append(invalidHandlers, handler.Field)
				continue
			}

			// Verify struct field exists
			_, found := veType.FieldByName(structFieldName)
			if !found {
				invalidHandlers = append(invalidHandlers, handler.Field)
			}
		}

		require.Empty(t, invalidHandlers,
			"Handlers with invalid enum or struct field mappings: %v",
			invalidHandlers)
	})

	t.Logf("✅ PERFECT ALIGNMENT VERIFIED:")
	t.Logf("  ✓ All non-excluded struct fields have enum constants and handlers")
	t.Logf("  ✓ All enum constants have struct fields and handlers")
	t.Logf("  ✓ All handlers reference valid enum constants and struct fields")
	t.Logf("  ✓ All counts match: %d fields = %d enums = %d handlers",
		veType.NumField()-len(excludedStructFields), len(allEnumConstants), len(handlers))
}

// TestFieldVoteMapNoNilPanic is a regression test specifically for the nil map panic
func TestFieldVoteMapNoNilPanic(t *testing.T) {
	t.Run("no_panic_when_accessing_initialized_maps", func(t *testing.T) {
		fieldHandlers := initializeFieldHandlers()
		fieldVotes := make(map[VoteExtensionField]map[string]fieldVote)

		// Initialize using the fixed pattern
		for _, handler := range fieldHandlers {
			fieldVotes[handler.Field] = make(map[string]fieldVote)
		}

		// This should never panic
		require.NotPanics(t, func() {
			for _, handler := range fieldHandlers {
				votes := fieldVotes[handler.Field]
				// Simulate the write that was causing the panic
				votes["test_key"] = fieldVote{
					value:     "test_value",
					votePower: 1000,
				}
			}
		}, "Accessing fieldVotes map caused a panic")
	})

	t.Run("panic_detection_if_using_old_pattern", func(t *testing.T) {
		// This test documents the OLD buggy behavior for reference
		// DO NOT change the initialization in GetConsensusAndPluralityVEData to match this!

		// Simulate the old buggy initialization pattern
		fieldVotes := make(map[VoteExtensionField]map[string]fieldVote)

		// OLD BUGGY CODE: Initialize based on enum range (DO NOT USE)
		// for i := VEFieldEigenDelegationsHash; i <= VEFieldLatestZcashHeaderHash; i++ {
		//     fieldVotes[i] = make(map[string]fieldVote)
		// }

		// Instead, we'll intentionally NOT initialize to show what happens
		fieldHandlers := initializeFieldHandlers()

		// Pick a field that might not be initialized
		handler := fieldHandlers[0]

		// This would panic if the map wasn't initialized
		require.Panics(t, func() {
			votes := fieldVotes[handler.Field] // This gets nil
			votes["test_key"] = fieldVote{     // This panics: assignment to entry in nil map
				value:     "test_value",
				votePower: 1000,
			}
		}, "Expected panic when accessing uninitialized map, but got none")
	})
}
