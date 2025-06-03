package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/log"
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
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// setupTestKeeper creates a new keeper instance with mocked dependencies for testing
func setupTestKeeper(t *testing.T, sdkCtx sdk.Context) (*keeper.Keeper, *testutil.MocksidecarClient, *gomock.Controller) {
	ctrl := gomock.NewController(t)

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
