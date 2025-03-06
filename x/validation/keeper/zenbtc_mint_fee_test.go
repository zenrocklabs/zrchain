package keeper

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/assert"
)

func TestCalculateZenBTCMintFee(t *testing.T) {
	keeper := Keeper{}

	tests := []struct {
		name         string
		ethBaseFee   uint64
		ethTipCap    uint64
		ethGasLimit  uint64
		btcUSDPrice  sdkmath.LegacyDec
		ethUSDPrice  sdkmath.LegacyDec
		exchangeRate sdkmath.LegacyDec
		expected     uint64
	}{
		{
			name:         "zero BTC price returns zero fee to avoid division by zero",
			ethBaseFee:   30_000_000_000, // 30 gwei
			ethTipCap:    2_000_000_000,  // 2 gwei
			ethGasLimit:  285_000,
			btcUSDPrice:  sdkmath.LegacyNewDec(0),
			ethUSDPrice:  sdkmath.LegacyNewDec(2000),
			exchangeRate: sdkmath.LegacyNewDec(1),
			expected:     0,
		},
		{
			name:         "typical Ethereum mainnet values",
			ethBaseFee:   30_000_000_000, // 30 gwei
			ethTipCap:    2_000_000_000,  // 2 gwei
			ethGasLimit:  285_000,
			btcUSDPrice:  sdkmath.LegacyNewDec(90_000_00), // $90k BTC in cents
			ethUSDPrice:  sdkmath.LegacyNewDec(3_000_00),  // $3k ETH in cents
			exchangeRate: sdkmath.LegacyNewDec(1),
			expected:     30399, // 0.00030399 BTC in fees (1:1 exchange rate)
		},
		{
			name:         "typical Ethereum Holesky testnet values",
			ethBaseFee:   1816605,
			ethTipCap:    1000000,
			ethGasLimit:  239646,
			btcUSDPrice:  sdkmath.LegacyNewDec(90_000_00), // $90k BTC in cents
			ethUSDPrice:  sdkmath.LegacyNewDec(2_000_00),  // $2k ETH in cents
			exchangeRate: sdkmath.LegacyNewDec(1),
			expected:     1, // 0.00000001 BTC in fees (1:1 exchange rate)
		},
		{
			name:         "high gas price scenario",
			ethBaseFee:   100_000_000_000, // 100 gwei
			ethTipCap:    5_000_000_000,   // 5 gwei
			ethGasLimit:  285_000,
			btcUSDPrice:  sdkmath.LegacyNewDec(90_000_00),
			ethUSDPrice:  sdkmath.LegacyNewDec(3_000_00),
			exchangeRate: sdkmath.LegacyNewDec(1),
			expected:     99749, // 0.00099749 BTC in fees (1:1 exchange rate)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := keeper.CalculateZenBTCMintFee(
				tt.ethBaseFee,
				tt.ethTipCap,
				tt.ethGasLimit,
				tt.btcUSDPrice,
				tt.ethUSDPrice,
				tt.exchangeRate,
			)
			assert.Equal(t, tt.expected, result)
		})
	}
}
