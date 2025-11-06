package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJupiterPriceResponse_Parsing(t *testing.T) {
	tests := []struct {
		name          string
		jsonResponse  string
		expectedPrice float64
		shouldError   bool
		errorContains string
	}{
		{
			name: "Valid Jupiter response",
			jsonResponse: `{
				"5VsPJ2EG7jjo3k2LPzQVriENKKQkNUTzujEzuaj4Aisf": {
					"usdPrice": 0.123456,
					"blockId": 378315808,
					"decimals": 6,
					"priceChange24h": 9.84
				}
			}`,
			expectedPrice: 0.123456,
			shouldError:   false,
		},
		{
			name: "Valid Jupiter response with different price",
			jsonResponse: `{
				"5VsPJ2EG7jjo3k2LPzQVriENKKQkNUTzujEzuaj4Aisf": {
					"usdPrice": 1.234567890123456,
					"blockId": 378315808,
					"decimals": 6,
					"priceChange24h": 5.23
				}
			}`,
			expectedPrice: 1.234567890123456,
			shouldError:   false,
		},
		{
			name: "Valid Jupiter response with small price",
			jsonResponse: `{
				"5VsPJ2EG7jjo3k2LPzQVriENKKQkNUTzujEzuaj4Aisf": {
					"usdPrice": 0.000123,
					"blockId": 378315808,
					"decimals": 6,
					"priceChange24h": -2.5
				}
			}`,
			expectedPrice: 0.000123,
			shouldError:   false,
		},
		{
			name: "Missing ROCK token in response",
			jsonResponse: `{
				"SomeOtherToken123": {
					"usdPrice": 1.23,
					"blockId": 123,
					"decimals": 9,
					"priceChange24h": 0
				}
			}`,
			shouldError:   true,
			errorContains: "not found",
		},
		{
			name:          "Empty data object",
			jsonResponse:  `{}`,
			shouldError:   true,
			errorContains: "not found",
		},
		{
			name: "Zero price",
			jsonResponse: `{
				"5VsPJ2EG7jjo3k2LPzQVriENKKQkNUTzujEzuaj4Aisf": {
					"usdPrice": 0,
					"blockId": 378315808,
					"decimals": 6,
					"priceChange24h": 0
				}
			}`,
			shouldError:   true,
			errorContains: "invalid",
		},
		{
			name: "Negative price",
			jsonResponse: `{
				"5VsPJ2EG7jjo3k2LPzQVriENKKQkNUTzujEzuaj4Aisf": {
					"usdPrice": -1.23,
					"blockId": 378315808,
					"decimals": 6,
					"priceChange24h": 0
				}
			}`,
			shouldError:   true,
			errorContains: "invalid",
		},
		{
			name:          "Invalid JSON",
			jsonResponse:  `{"invalid json`,
			shouldError:   true,
			errorContains: "decode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jupiterResp JupiterPriceResponse
			err := json.Unmarshal([]byte(tt.jsonResponse), &jupiterResp)

			if tt.shouldError && tt.errorContains == "decode" {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			const rockTokenID = "5VsPJ2EG7jjo3k2LPzQVriENKKQkNUTzujEzuaj4Aisf"
			priceData, ok := jupiterResp[rockTokenID]

			if tt.shouldError && tt.errorContains == "not found" {
				assert.False(t, ok, "Expected token not to be found")
				return
			}

			if tt.shouldError && tt.errorContains == "invalid" {
				require.True(t, ok)
				assert.True(t, priceData.USDPrice <= 0, "Price should be zero or negative")
				return
			}

			require.True(t, ok, "Expected ROCK token to be present")
			assert.Equal(t, tt.expectedPrice, priceData.USDPrice)
			assert.Greater(t, priceData.USDPrice, 0.0)

			// Verify we can convert to math.LegacyDec
			priceStr := fmt.Sprintf("%.18f", priceData.USDPrice)
			priceDec, err := math.LegacyNewDecFromStr(priceStr)
			require.NoError(t, err)
			assert.False(t, priceDec.IsNil())
			assert.True(t, priceDec.IsPositive())
		})
	}
}

func TestJupiterAPI_Integration(t *testing.T) {
	tests := []struct {
		name          string
		response      string
		statusCode    int
		shouldError   bool
		errorContains string
	}{
		{
			name: "Successful API call",
			response: `{
				"5VsPJ2EG7jjo3k2LPzQVriENKKQkNUTzujEzuaj4Aisf": {
					"usdPrice": 0.567890,
					"blockId": 378315808,
					"decimals": 6,
					"priceChange24h": 5.5
				}
			}`,
			statusCode:  http.StatusOK,
			shouldError: false,
		},
		{
			name:          "Malformed JSON response",
			response:      `{broken json}`,
			statusCode:    http.StatusOK,
			shouldError:   true,
			errorContains: "decode",
		},
		{
			name: "Missing token in response",
			response: `{
				"SomeOtherToken": {
					"usdPrice": 1.23,
					"blockId": 123,
					"decimals": 9,
					"priceChange24h": 0
				}
			}`,
			statusCode:  http.StatusOK,
			shouldError: false, // Decode succeeds, but token lookup fails
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			// Make request to test server
			resp, err := http.Get(server.URL)
			require.NoError(t, err)
			defer resp.Body.Close()

			var jupiterResp JupiterPriceResponse
			err = json.NewDecoder(resp.Body).Decode(&jupiterResp)

			if tt.shouldError && tt.errorContains == "decode" {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			const rockTokenID = "5VsPJ2EG7jjo3k2LPzQVriENKKQkNUTzujEzuaj4Aisf"
			priceData, ok := jupiterResp[rockTokenID]

			if !ok {
				// Token not found - this is expected for some tests
				return
			}

			assert.Greater(t, priceData.USDPrice, 0.0)

			// Verify conversion to math.LegacyDec works
			priceStr := fmt.Sprintf("%.18f", priceData.USDPrice)
			priceDec, err := math.LegacyNewDecFromStr(priceStr)
			require.NoError(t, err)
			assert.True(t, priceDec.IsPositive())
		})
	}
}

func TestJupiterAPI_RealEndpoint(t *testing.T) {
	// This test makes a real call to Jupiter API
	// Skip in CI or when network is unavailable
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	const jupiterURL = "https://lite-api.jup.ag/price/v3?ids=5VsPJ2EG7jjo3k2LPzQVriENKKQkNUTzujEzuaj4Aisf"

	resp, err := http.Get(jupiterURL)
	require.NoError(t, err, "Failed to fetch from Jupiter API")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected 200 OK from Jupiter API")

	var jupiterResp JupiterPriceResponse
	err = json.NewDecoder(resp.Body).Decode(&jupiterResp)
	require.NoError(t, err, "Failed to decode Jupiter response")

	const rockTokenID = "5VsPJ2EG7jjo3k2LPzQVriENKKQkNUTzujEzuaj4Aisf"
	priceData, ok := jupiterResp[rockTokenID]
	require.True(t, ok, "ROCK token not found in Jupiter response")
	require.Greater(t, priceData.USDPrice, 0.0, "Price should be positive")

	// Verify we can parse the price
	priceStr := fmt.Sprintf("%.18f", priceData.USDPrice)
	priceDec, err := math.LegacyNewDecFromStr(priceStr)
	require.NoError(t, err, "Failed to parse price as decimal")
	assert.True(t, priceDec.IsPositive(), "Price should be positive")

	t.Logf("Successfully fetched ROCK price from Jupiter: $%.6f", priceData.USDPrice)
	t.Logf("Block ID: %d, Decimals: %d, 24h Change: %.2f%%",
		priceData.BlockID, priceData.Decimals, priceData.PriceChange24h)
}

func TestJupiterAPI_ErrorHandling(t *testing.T) {
	t.Run("Missing usdPrice field", func(t *testing.T) {
		jsonResp := `{
			"5VsPJ2EG7jjo3k2LPzQVriENKKQkNUTzujEzuaj4Aisf": {
				"blockId": 378315808,
				"decimals": 6
			}
		}`

		var jupiterResp JupiterPriceResponse
		err := json.Unmarshal([]byte(jsonResp), &jupiterResp)
		require.NoError(t, err)

		const rockTokenID = "5VsPJ2EG7jjo3k2LPzQVriENKKQkNUTzujEzuaj4Aisf"
		priceData, ok := jupiterResp[rockTokenID]
		require.True(t, ok)
		assert.Equal(t, 0.0, priceData.USDPrice, "Missing price field should default to 0")
	})

	t.Run("Null response", func(t *testing.T) {
		jsonResp := `null`

		var jupiterResp JupiterPriceResponse
		err := json.Unmarshal([]byte(jsonResp), &jupiterResp)
		require.NoError(t, err)
		assert.Nil(t, jupiterResp)
	})

	t.Run("Float precision preservation", func(t *testing.T) {
		// Test that we can handle very small and very precise prices
		testPrices := []float64{
			0.000001,
			0.123456789012345,
			123.456789012345,
			0.02171898832948882, // Actual ROCK price from API
		}

		for _, price := range testPrices {
			priceStr := fmt.Sprintf("%.18f", price)
			priceDec, err := math.LegacyNewDecFromStr(priceStr)
			require.NoError(t, err, "Should parse price: %f", price)
			assert.True(t, priceDec.IsPositive(), "Price should be positive: %f", price)
		}
	})
}

func TestJupiterAPI_PriceConversion(t *testing.T) {
	t.Run("Convert float to LegacyDec correctly", func(t *testing.T) {
		testCases := []struct {
			floatPrice    float64
			expectedValid bool
		}{
			{0.123456, true},
			{1.0, true},
			{999999.99, true},
			{0.000001, true},
			{0, false}, // Zero price should be invalid
			{-1.23, false}, // Negative price should be invalid
		}

		for _, tc := range testCases {
			priceStr := fmt.Sprintf("%.18f", tc.floatPrice)
			priceDec, err := math.LegacyNewDecFromStr(priceStr)

			if tc.expectedValid {
				require.NoError(t, err, "Should parse valid price: %f", tc.floatPrice)
				if tc.floatPrice > 0 {
					assert.True(t, priceDec.IsPositive(), "Price should be positive: %f", tc.floatPrice)
				}
			} else {
				// For invalid prices (0 or negative), parsing succeeds but value check should fail
				require.NoError(t, err)
				if tc.floatPrice <= 0 {
					assert.False(t, priceDec.IsPositive() || priceDec.IsNil(),
						"Zero/negative price should not be positive: %f", tc.floatPrice)
				}
			}
		}
	})
}
