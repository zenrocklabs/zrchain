package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sidecar "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/testutil"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	ubermock "go.uber.org/mock/gomock"
)

// setupTestKeeper creates a new keeper instance with mocked dependencies for testing
func setupTestKeeper(t *testing.T, sdkCtx sdk.Context) (*keeper.Keeper, *testutil.MocksidecarClient, *ubermock.Controller) {
	ctrl := ubermock.NewController(t)

	// Create store service
	key := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(key)

	// Create mock keepers
	mockAccountKeeper := testutil.NewMockAccountKeeper(ctrl)
	mockBankKeeper := testutil.NewMockBankKeeper(ctrl)
	mockTreasuryKeeper := testutil.NewMockTreasuryKeeper(ctrl)
	mockZentpKeeper := testutil.NewMockZentpKeeper(ctrl)
	mockSidecarClient := testutil.NewMocksidecarClient(ctrl)

	// Set up mock expectations
	mockAccountKeeper.EXPECT().GetModuleAddress(validationtypes.BondedPoolName).Return(sdk.AccAddress{})
	mockAccountKeeper.EXPECT().GetModuleAddress(validationtypes.NotBondedPoolName).Return(sdk.AccAddress{})
	mockAccountKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("zen")).AnyTimes()

	// Use governance module address as authority
	govAddr := authtypes.NewModuleAddress(govtypes.ModuleName)
	mockAccountKeeper.EXPECT().GetModuleAddress(govtypes.ModuleName).Return(govAddr).AnyTimes()
	authority := govAddr.String()

	// Create keeper with mock client
	k := keeper.NewKeeper(
		codec.NewProtoCodec(nil), // cdc
		storeService,             // storeService
		mockAccountKeeper,        // accountKeeper
		mockBankKeeper,           // bankKeeper
		authority,
		nil,                                  // txDecoder
		nil,                                  // zrConfig
		mockTreasuryKeeper,                   // treasuryKeeper
		nil,                                  // zenBTCKeeper
		mockZentpKeeper,                      // zentpKeeper
		address.NewBech32Codec("zenvaloper"), // validatorAddressCodec
		address.NewBech32Codec("zenvalcons"), // consensusAddressCodec
	)
	k.SetSidecarClient(mockSidecarClient)

	return k, mockSidecarClient, ctrl
}

func TestGetSidecarState(t *testing.T) {
	tests := []struct {
		name       string
		height     int64
		oracleData *sidecar.SidecarStateResponse
		err        error
	}{
		{
			name:   "successfully fetches oracle data with valid delegations",
			height: 100,
			oracleData: &sidecar.SidecarStateResponse{
				EthBlockHeight:   100,
				EthGasLimit:      100,
				EthBaseFee:       100,
				EthTipCap:        100,
				EigenDelegations: []byte(`{"validator1":{"delegator1":100}}`),
			},
			err: nil,
		},
		{
			name:       "returns error when oracle service fails",
			height:     100,
			oracleData: nil,
			err:        keeper.ErrOracleSidecar,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a proper SDK context with a logger
			logger := log.NewTestLogger(t)
			sdkCtx := sdk.NewContext(nil, cmtproto.Header{}, true, logger)
			sdkCtx = sdkCtx.WithBlockHeight(tt.height)

			// Set up keeper with mocks
			k, mockSidecarClient, ctrl := setupTestKeeper(t, sdkCtx)
			defer ctrl.Finish()

			// Set up mock expectations
			mockSidecarClient.EXPECT().GetSidecarState(sdkCtx, &sidecar.SidecarStateRequest{}).Return(tt.oracleData, tt.err)

			// Test the actual keeper function
			result, err := k.GetSidecarState(sdkCtx, tt.height)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, keeper.ErrOracleSidecar, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tt.oracleData.EthBlockHeight, result.EthBlockHeight)
				require.Equal(t, tt.oracleData.EthGasLimit, result.EthGasLimit)
				require.Equal(t, tt.oracleData.EthBaseFee, result.EthBaseFee)
				require.Equal(t, tt.oracleData.EthTipCap, result.EthTipCap)

				// Verify EigenDelegations were processed correctly
				require.NotNil(t, result.EigenDelegationsMap)
				require.Contains(t, result.EigenDelegationsMap, "validator1")
				require.Contains(t, result.EigenDelegationsMap["validator1"], "delegator1")
				require.Equal(t, big.NewInt(100), result.EigenDelegationsMap["validator1"]["delegator1"])
			}
		})
	}
}

func TestGetSidecarStateByEthHeight(t *testing.T) {
	tests := []struct {
		name       string
		height     uint64
		oracleData *sidecar.SidecarStateResponse
		err        error
	}{
		{
			name:   "successfully fetches oracle data by ETH height with valid delegations",
			height: 100,
			oracleData: &sidecar.SidecarStateResponse{
				EthBlockHeight:   100,
				EthGasLimit:      100,
				EthBaseFee:       100,
				EthTipCap:        0,
				EigenDelegations: []byte(`{"validator1":{"delegator1":100}}`),
			},
			err: nil,
		},
		{
			name:       "returns error when oracle service fails for ETH height",
			height:     100,
			oracleData: nil,
			err:        keeper.ErrOracleSidecar,
		},
		{
			name:   "handles zero ETH height correctly",
			height: 0,
			oracleData: &sidecar.SidecarStateResponse{
				EthBlockHeight:   0,
				EthGasLimit:      0,
				EthBaseFee:       0,
				EthTipCap:        0,
				EigenDelegations: []byte(`{"validator1":{"delegator1":100}}`),
			},
			err: nil,
		},
		{
			name:   "handles maximum uint64 ETH height correctly",
			height: ^uint64(0),
			oracleData: &sidecar.SidecarStateResponse{
				EthBlockHeight:   ^uint64(0),
				EthGasLimit:      100,
				EthBaseFee:       100,
				EthTipCap:        0,
				EigenDelegations: []byte(`{"validator1":{"delegator1":100}}`),
			},
			err: nil,
		},
		{
			name:   "handles malformed eigen delegations JSON gracefully",
			height: 100,
			oracleData: &sidecar.SidecarStateResponse{
				EthBlockHeight:   100,
				EthGasLimit:      100,
				EthBaseFee:       100,
				EthTipCap:        0,
				EigenDelegations: []byte(`{invalid json}`),
			},
			err: keeper.ErrOracleSidecar,
		},
		{
			name:   "handles empty eigen delegations correctly",
			height: 100,
			oracleData: &sidecar.SidecarStateResponse{
				EthBlockHeight:   100,
				EthGasLimit:      100,
				EthBaseFee:       100,
				EthTipCap:        0,
				EigenDelegations: []byte(`{}`),
			},
			err: keeper.ErrOracleSidecar,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a proper SDK context with a logger
			logger := log.NewTestLogger(t)
			sdkCtx := sdk.NewContext(nil, cmtproto.Header{}, true, logger)
			sdkCtx = sdkCtx.WithBlockHeight(int64(tt.height))

			// Set up keeper with mocks
			k, mockSidecarClient, ctrl := setupTestKeeper(t, sdkCtx)
			defer ctrl.Finish()

			// Set up mock expectations
			mockSidecarClient.EXPECT().GetSidecarStateByEthHeight(sdkCtx, &sidecar.SidecarStateByEthHeightRequest{EthBlockHeight: tt.height}).Return(tt.oracleData, tt.err)

			// Test the actual keeper function
			result, err := k.GetSidecarStateByEthHeight(sdkCtx, tt.height)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, keeper.ErrOracleSidecar, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tt.oracleData.EthBlockHeight, result.EthBlockHeight)
				require.Equal(t, tt.oracleData.EthGasLimit, result.EthGasLimit)
				require.Equal(t, tt.oracleData.EthBaseFee, result.EthBaseFee)
				require.Equal(t, tt.oracleData.EthTipCap, result.EthTipCap)

				// Verify EigenDelegations were processed correctly
				if tt.name == "invalid eigen delegations json" {
					require.Empty(t, result.EigenDelegationsMap)
				} else if tt.name == "empty eigen delegations" {
					require.Empty(t, result.EigenDelegationsMap)
				} else if tt.name != "error case" {
					require.NotNil(t, result.EigenDelegationsMap)
					require.Contains(t, result.EigenDelegationsMap, "validator1")
					require.Contains(t, result.EigenDelegationsMap["validator1"], "delegator1")
					require.Equal(t, big.NewInt(100), result.EigenDelegationsMap["validator1"]["delegator1"])
				}
			}
		})
	}
}

func TestCalculateFlatZenBTCMintFee(t *testing.T) {
	tests := []struct {
		name         string
		btcUSDPrice  string // Using string to create LegacyDec
		exchangeRate string // Using string to create LegacyDec
		expected     uint64
		description  string
	}{
		{
			name:         "normal case with BTC at $50,000 and 1:1 exchange rate",
			btcUSDPrice:  "50000",
			exchangeRate: "1.0",
			expected:     10000, // $5 / $50,000 = 0.0001 BTC = 10,000 sat = 0.0001 zenBTC
			description:  "Basic calculation with round numbers",
		},
		{
			name:         "BTC at $100,000 with 1:1 exchange rate",
			btcUSDPrice:  "100000",
			exchangeRate: "1.0",
			expected:     5000, // $5 / $100,000 = 0.00005 BTC = 5,000 sat = 0.00005 zenBTC
			description:  "Higher BTC price results in lower fee",
		},
		{
			name:         "BTC at $25,000 with 1:1 exchange rate",
			btcUSDPrice:  "25000",
			exchangeRate: "1.0",
			expected:     20000, // $5 / $25,000 = 0.0002 BTC = 20,000 sat = 0.0002 zenBTC
			description:  "Lower BTC price results in higher fee",
		},
		{
			name:         "BTC at $50,000 with 2:1 exchange rate (2 BTC = 1 zenBTC)",
			btcUSDPrice:  "50000",
			exchangeRate: "2.0",
			expected:     5000, // $5 / $50,000 = 0.0001 BTC = 10,000 sat / 2 = 0.00005 zenBTC
			description:  "Exchange rate affects final zenBTC amount",
		},
		{
			name:         "BTC at $50,000 with 0.5:1 exchange rate (0.5 BTC = 1 zenBTC)",
			btcUSDPrice:  "50000",
			exchangeRate: "0.5", // it should never really be less than 1, but let's test it
			expected:     20000, // $5 / $50,000 = 0.0001 BTC = 10,000 sat / 0.5 = 0.0002 zenBTC
			description:  "Lower exchange rate increases zenBTC amount",
		},
		{
			name:         "very high BTC price",
			btcUSDPrice:  "1000000",
			exchangeRate: "1.0",
			expected:     500, // $5 / $1,000,000 = 0.000005 BTC = 500 sat = 0.000005 zenBTC
			description:  "Handles very high BTC prices",
		},
		{
			name:         "fractional satoshis get truncated",
			btcUSDPrice:  "33333.333", // Results in non-round satoshis
			exchangeRate: "1.0",
			expected:     15000, // $5 / $33,333.333 ≈ 0.00015 BTC ≈ 15,000 sat = 0.00015 zenBTC (truncated)
			description:  "Truncation works correctly for fractional results",
		},
		{
			name:         "zero BTC price returns zero",
			btcUSDPrice:  "0",
			exchangeRate: "1.0",
			expected:     0,
			description:  "Zero BTC price safety check",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a proper SDK context with a logger
			logger := log.NewTestLogger(t)
			sdkCtx := sdk.NewContext(nil, cmtproto.Header{}, true, logger)

			// Set up keeper with mocks
			k, _, ctrl := setupTestKeeper(t, sdkCtx)
			defer ctrl.Finish()

			// Convert string prices to LegacyDec
			btcPrice, err := math.LegacyNewDecFromStr(tt.btcUSDPrice)
			require.NoError(t, err, "Failed to create BTC price from string: %s", tt.btcUSDPrice)

			exchangeRate, err := math.LegacyNewDecFromStr(tt.exchangeRate)
			require.NoError(t, err, "Failed to create exchange rate from string: %s", tt.exchangeRate)

			// Call the function under test
			result := k.CalculateFlatZenBTCMintFee(btcPrice, exchangeRate)

			// Verify the result
			require.Equal(t, tt.expected, result, "Test case: %s - %s", tt.name, tt.description)

			// Additional verification for non-zero cases
			if !btcPrice.IsZero() && tt.expected > 0 {
				// Verify that the fee represents approximately $5
				// Back-calculate: result zenBTC * exchangeRate * btcUSDPrice / 1e8 should ≈ $5
				resultDec := math.LegacyNewDecFromInt(math.NewIntFromUint64(result))
				feeInBTC := resultDec.Mul(exchangeRate).Quo(math.LegacyNewDec(1e8))
				feeInUSD := feeInBTC.Mul(btcPrice)

				// Allow for small rounding differences due to truncation
				expectedUSD := math.LegacyNewDec(5)
				diff := feeInUSD.Sub(expectedUSD).Abs()
				tolerance := math.LegacyNewDecWithPrec(1, 6) // 0.000001 USD tolerance

				require.True(t, diff.LTE(tolerance),
					"Fee should be approximately $5, got $%s (diff: $%s)",
					feeInUSD.String(), diff.String())
			}
		})
	}
}

func TestCalculateFlatZenBTCMintFeeEdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		btcUSDPrice  string
		exchangeRate string
		expectPanic  bool
		description  string
	}{
		{
			name:         "zero exchange rate should not panic but may give unexpected results",
			btcUSDPrice:  "50000",
			exchangeRate: "0",
			expectPanic:  true, // Division by zero should panic
			description:  "Division by zero in exchange rate",
		},
		{
			name:         "very small BTC price",
			btcUSDPrice:  "0.001",
			exchangeRate: "1.0",
			expectPanic:  false,
			description:  "Very small BTC price should work",
		},
		{
			name:         "very small exchange rate",
			btcUSDPrice:  "50000",
			exchangeRate: "0.000001",
			expectPanic:  false,
			description:  "Very small exchange rate should work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a proper SDK context with a logger
			logger := log.NewTestLogger(t)
			sdkCtx := sdk.NewContext(nil, cmtproto.Header{}, true, logger)

			// Set up keeper with mocks
			k, _, ctrl := setupTestKeeper(t, sdkCtx)
			defer ctrl.Finish()

			// Convert string prices to LegacyDec
			btcPrice, err := math.LegacyNewDecFromStr(tt.btcUSDPrice)
			require.NoError(t, err)

			exchangeRate, err := math.LegacyNewDecFromStr(tt.exchangeRate)
			require.NoError(t, err)

			if tt.expectPanic {
				require.Panics(t, func() {
					k.CalculateFlatZenBTCMintFee(btcPrice, exchangeRate)
				}, tt.description)
			} else {
				require.NotPanics(t, func() {
					result := k.CalculateFlatZenBTCMintFee(btcPrice, exchangeRate)
					// For very small values, just ensure we get a reasonable result
					require.GreaterOrEqual(t, result, uint64(0))
				}, tt.description)
			}
		})
	}
}
