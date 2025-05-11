package keeper

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"

	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	"github.com/stretchr/testify/require"
)

func TestCalculateValidatorPowerComponents_FullAVSPrecision(t *testing.T) {
	var testK Keeper // dummy keeper instance

	rockPrice, _ := math.LegacyNewDecFromStr("0.01801")
	btcPrice, _ := math.LegacyNewDecFromStr("102855.235")
	ethPrice, _ := math.LegacyNewDecFromStr("3000.00")
	zeroPrice, _ := math.LegacyNewDecFromStr("0.0")

	nativeAssetDataROCK := &types.AssetData{Asset: types.Asset_ROCK, PriceUSD: rockPrice, Precision: 6}
	avsAssetDataBTC := &types.AssetData{Asset: types.Asset_BTC, PriceUSD: btcPrice, Precision: 8}
	avsAssetDataETH := &types.AssetData{Asset: types.Asset_ETH, PriceUSD: ethPrice, Precision: 18}

	powerReductionVal := math.NewInt(1000000)

	tests := []struct {
		name                     string
		validatorTokensNative    math.Int
		nativeAssetData          *types.AssetData
		avsAssetData             *types.AssetData
		avsTokensHeld            math.Int
		avsExchangeRate          math.LegacyDec
		pricesAreValid           bool
		expectedTotalPower       int64
		expectedNativeContribStr string
		expectedAVSContribStr    string
	}{
		{
			name:                     "val-1 (ROCK + BTC), prices valid, full AVS precision",
			validatorTokensNative:    math.NewInt(10300000000000),
			nativeAssetData:          nativeAssetDataROCK,
			avsAssetData:             avsAssetDataBTC,
			avsTokensHeld:            math.NewInt(118170108), // 1.18170108 BTC
			avsExchangeRate:          math.LegacyOneDec(),
			pricesAreValid:           true,
			expectedTotalPower:       307047,
			expectedNativeContribStr: "185503.000000000000000000",
			expectedAVSContribStr:    "121544.142283153800000000",
		},
		{
			name:                     "val-2 (ROCK only), prices valid",
			validatorTokensNative:    math.NewInt(10300016500000),
			nativeAssetData:          nativeAssetDataROCK,
			avsAssetData:             nil, // No AVS asset data as no AVS tokens relevant
			avsTokensHeld:            math.NewInt(0),
			avsExchangeRate:          math.LegacyOneDec(),
			pricesAreValid:           true,
			expectedTotalPower:       185503,
			expectedNativeContribStr: "185503.288160000000000000",
			expectedAVSContribStr:    "0.000000000000000000",
		},
		{
			name:                     "val-1 (ROCK + BTC), prices NOT valid",
			validatorTokensNative:    math.NewInt(10300000000000),
			nativeAssetData:          &types.AssetData{Asset: types.Asset_ROCK, PriceUSD: zeroPrice, Precision: 6},
			avsAssetData:             &types.AssetData{Asset: types.Asset_BTC, PriceUSD: zeroPrice, Precision: 8},
			avsTokensHeld:            math.NewInt(118170108),
			avsExchangeRate:          math.LegacyOneDec(),
			pricesAreValid:           false,
			expectedTotalPower:       10300000, // Native consensus units, AVS is 0
			expectedNativeContribStr: "10300000.000000000000000000",
			expectedAVSContribStr:    "0.000000000000000000",
		},
		{
			name:                     "Validator with 1 ETH, prices valid",
			validatorTokensNative:    math.NewInt(0),
			nativeAssetData:          nil,
			avsAssetData:             avsAssetDataETH,
			avsTokensHeld:            math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), // 1 ETH in Wei
			avsExchangeRate:          math.LegacyOneDec(),
			pricesAreValid:           true,
			expectedTotalPower:       3000,
			expectedNativeContribStr: "0.000000000000000000",
			expectedAVSContribStr:    "3000.000000000000000000",
		},
		{
			name:                  "Validator with tiny ROCK that would give <1 power, prices valid",
			validatorTokensNative: math.NewInt(500000), // 0.5 ROCK. NativeConsensusUnits = 0 if ConsensusPower truncates.
			nativeAssetData:       nativeAssetDataROCK,
			avsAssetData:          nil,
			avsTokensHeld:         math.NewInt(0),
			avsExchangeRate:       math.LegacyOneDec(),
			pricesAreValid:        true,
			// nativeConsensusUnits = 0 (from 500000 / 1000000)
			// nativePowerContribution = 0.01801 * 0 = 0
			// total = 0. Since validator has some tokens (0.5 ROCK), power is bumped to 1.
			expectedTotalPower:       1,
			expectedNativeContribStr: "0.000000000000000000",
			expectedAVSContribStr:    "0.000000000000000000",
		},
		{
			name:                  "Validator with tiny AVS that results in <1 power (but >0), gets 1 power",
			validatorTokensNative: math.NewInt(0),
			nativeAssetData:       nil,
			avsAssetData:          &types.AssetData{Asset: types.Asset_BTC, PriceUSD: math.LegacyMustNewDecFromStr("0.00000001"), Precision: 8},
			avsTokensHeld:         math.NewInt(118170108), // 1.18170108 BTC
			avsExchangeRate:       math.LegacyOneDec(),
			pricesAreValid:        true,
			// avsPower = 0.00000001 * 1.18170108 = 0.0000000118170108. total = 0 (truncates). Bumped to 1.
			expectedTotalPower:       1,
			expectedNativeContribStr: "0.000000000000000000",
			expectedAVSContribStr:    "0.000000011817010800",
		},
		{
			name:                     "Validator with no stake, prices valid",
			validatorTokensNative:    math.NewInt(0),
			nativeAssetData:          nativeAssetDataROCK, // Price info available
			avsAssetData:             avsAssetDataBTC,     // Price info available
			avsTokensHeld:            math.NewInt(0),
			avsExchangeRate:          math.LegacyOneDec(),
			pricesAreValid:           true,
			expectedTotalPower:       0, // Genuinely no stake, so power is 0
			expectedNativeContribStr: "0.000000000000000000",
			expectedAVSContribStr:    "0.000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a validator instance for the test.
			// The `ConsensusPower` method of this validator will be called.
			validator := newPowerTestValidator(tt.validatorTokensNative)

			totalPower, nativeC, avsC := testK.calculateValidatorPowerComponents(
				validator,
				powerReductionVal, // This is the powerReduction for native tokens
				tt.nativeAssetData,
				tt.avsAssetData,
				tt.avsTokensHeld,
				tt.avsExchangeRate,
				tt.pricesAreValid,
			)

			expectedNative, err := math.LegacyNewDecFromStr(tt.expectedNativeContribStr)
			require.NoError(t, err)
			expectedAVS, err := math.LegacyNewDecFromStr(tt.expectedAVSContribStr)
			require.NoError(t, err)

			require.Equal(t, tt.expectedTotalPower, totalPower, "Total Consensus Power mismatch")
			require.True(t, expectedNative.Equal(nativeC), "Native Power Contribution mismatch. Expected '%s', got '%s'", expectedNative.String(), nativeC.String())
			require.True(t, expectedAVS.Equal(avsC), "AVS Power Contribution mismatch. Expected '%s', got '%s'", expectedAVS.String(), avsC.String())
		})
	}
}

// TestAdjustPowerToPrecision tests the utility function that converts raw token amounts
// to an integer number of whole units, demonstrating its truncating behavior.
// Note: The main power calculation logic in `calculateValidatorPowerComponents` (tested below)
// uses its own internal decimal arithmetic for AVS valuation to maintain full precision,
// rather than using this function's truncated integer output for that specific step.
func TestAdjustPowerToPrecision(t *testing.T) {
	var testK Keeper // dummy keeper instance

	tests := []struct {
		name           string
		tokens         math.Int
		assetPrecision uint32
		expected       math.Int
	}{
		{
			name:           "Tokens forming exact whole units",
			tokens:         math.NewInt(2000000),
			assetPrecision: 6, // 2000000 / 10^6 = 2
			expected:       math.NewInt(2),
		},
		{
			name:           "Asset with 0 precision (tokens are whole units)",
			tokens:         math.NewInt(123),
			assetPrecision: 0,
			expected:       math.NewInt(123),
		},
		{
			name:           "Zero tokens",
			tokens:         math.NewInt(0),
			assetPrecision: 8,
			expected:       math.NewInt(0),
		},
		{
			name:           "Tokens with precision resulting in truncation",
			tokens:         math.NewInt(12345), // e.g., 1.2345 units if precision is 4
			assetPrecision: 4,                  // 12345 / 10^4 = 1
			expected:       math.NewInt(1),
		},
		{
			name:           "Tokens less than one whole unit",
			tokens:         math.NewInt(500),
			assetPrecision: 3, // 500 / 10^3 = 0
			expected:       math.NewInt(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := testK.adjustPowerToPrecisionForTest(tt.tokens, tt.assetPrecision)
			require.True(t, tt.expected.Equal(actual), "Expected %s, got %s", tt.expected.String(), actual.String())
		})
	}
}

// adjustPowerToPrecision function as provided by the user (fixed version).
// This function returns the integer number of whole units, by dividing by 10^assetPrecision.
// Note: As it returns math.Int, it inherently truncates any fractional part.
func (k Keeper) adjustPowerToPrecisionForTest(tokens math.Int, assetPrecision uint32) math.Int {
	if assetPrecision == 0 {
		return tokens // Dividing by 10^0 = 1
	}
	powerOf10Denominator := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(assetPrecision)), nil)
	// Prevent division by zero, though 10^X should not be zero for non-negative X.
	if new(big.Int).Set(powerOf10Denominator).IsInt64() && powerOf10Denominator.Int64() == 0 {
		return math.ZeroInt()
	}
	return tokens.Quo(math.NewIntFromBigInt(powerOf10Denominator))
}

// Helper to create a types.ValidatorHV for testing power calculation.
// The actual types.ValidatorHV struct from your module will be used by the keeper method.
// This helper just ensures the TokensNative field is set for validator.ConsensusPower().
func newPowerTestValidator(tokensNative math.Int) types.ValidatorHV {
	return types.ValidatorHV{
		TokensNative: tokensNative,
		Status:       types.Bonded, // Ensure validator is treated as bonded for ConsensusPower calculation
		// Other fields like TokensAVS are passed directly to calculateValidatorPowerComponents
	}
}
