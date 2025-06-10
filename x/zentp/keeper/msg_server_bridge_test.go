package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/suite"

	minttypes "github.com/Zenrock-Foundation/zrchain/v6/x/mint/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestBridge() {
	// Setup test parameters
	params := types.DefaultParams()
	err := s.zentpKeeper.ParamStore.Set(s.ctx, params)
	s.Require().NoError(err)

	// Setup Solana ROCK supply for invariant check
	err = s.zentpKeeper.SetSolanaROCKSupply(s.ctx, math.NewIntFromUint64(100_000_000_000_000)) // 100M ROCK
	s.Require().NoError(err)

	// Mock bank keeper GetSupply for invariant check
	s.bankKeeper.EXPECT().GetSupply(s.ctx, "urock").Return(
		sdk.NewCoin("urock", math.NewIntFromUint64(200_000_000_000_000)), // 200M ROCK
	).AnyTimes()

	// Mock for new check in CheckROCKSupplyCap
	zentpModuleAddr := authtypes.NewModuleAddress(types.ModuleName)
	s.accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(zentpModuleAddr).AnyTimes()
	s.bankKeeper.EXPECT().GetBalance(s.ctx, zentpModuleAddr, "urock").Return(
		sdk.NewCoin("urock", math.ZeroInt()), // Assume module has zero balance
	).AnyTimes()

	// Mock getting the mint params
	s.mintKeeper.EXPECT().GetParams(s.ctx).Return(minttypes.DefaultParams(), nil)

	// Create test message
	msg := &types.MsgBridge{
		Creator:          "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		DestinationChain: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		Amount:           1000,
		Denom:            "urock",
		RecipientAddress: "1BbzosnmC3EVe7XcMgHYd6fUtcfdzUvfeaVZxaZ2QsE",
	}

	// Mock bank keeper GetBalance
	s.bankKeeper.EXPECT().GetBalance(
		s.ctx,
		sdk.MustAccAddressFromBech32(msg.Creator),
		msg.Denom,
	).Return(sdk.NewCoin(msg.Denom, math.NewIntFromUint64(msg.Amount+100000000+params.Solana.Fee*2)))

	// Calculate total amount including bridge fee and Solana fee
	baseAmountInt := math.NewIntFromUint64(msg.Amount)
	bridgeFeeAmount := math.LegacyNewDecFromInt(baseAmountInt).Mul(params.BridgeFee).TruncateInt()
	totalAmountInt := baseAmountInt.Add(bridgeFeeAmount).Add(math.NewIntFromUint64(params.Solana.Fee))

	// Mock bank keeper SendCoinsFromAccountToModule
	s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(
		s.ctx,
		sdk.MustAccAddressFromBech32(msg.Creator),
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin("urock", totalAmountInt)),
	).Return(nil)

	// Mock validation keeper SetSolanaRequestedNonce
	s.validationKeeper.EXPECT().SetSolanaRequestedNonce(
		s.ctx,
		params.Solana.NonceAccountKey, // Default nonce account key
		true,
	).Return(nil)

	// Mock validation keeper SetSolanaRequestedAccount
	s.validationKeeper.EXPECT().SetSolanaZenTPRequestedAccount(
		s.ctx,
		msg.RecipientAddress,
		true,
	).Return(nil)

	// Call the Bridge handler
	response, err := s.msgServer.Bridge(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().Equal(uint64(1), response.Id)
}

func (s *IntegrationTestSuite) TestBridgeFailureScenarios() {
	// Setup test parameters
	params := types.DefaultParams()
	err := s.zentpKeeper.ParamStore.Set(s.ctx, params)
	s.Require().NoError(err)

	// Setup Solana ROCK supply for invariant check (for most tests)
	err = s.zentpKeeper.SetSolanaROCKSupply(s.ctx, math.NewIntFromUint64(100_000_000_000_000)) // 100M ROCK
	s.Require().NoError(err)

	// Mock bank keeper GetSupply for invariant check (for most tests)
	s.bankKeeper.EXPECT().GetSupply(s.ctx, "urock").Return(
		sdk.NewCoin("urock", math.NewIntFromUint64(200_000_000_000_000)), // 200M ROCK
	).AnyTimes()

	// Mock for new check in CheckROCKSupplyCap
	zentpModuleAddr := authtypes.NewModuleAddress(types.ModuleName)
	s.accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(zentpModuleAddr).AnyTimes()
	s.bankKeeper.EXPECT().GetBalance(s.ctx, zentpModuleAddr, "urock").Return(
		sdk.NewCoin("urock", math.ZeroInt()), // Assume module has zero balance for most tests
	).AnyTimes()

	// Mock getting the mint params
	s.mintKeeper.EXPECT().GetParams(s.ctx).Return(minttypes.DefaultParams(), nil).AnyTimes()

	// Create base test message
	baseMsg := &types.MsgBridge{
		Creator:          "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		DestinationChain: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		Amount:           1000,
		Denom:            "urock",
		RecipientAddress: "1BbzosnmC3EVe7XcMgHYd6fUtcfdzUvfeaVZxaZ2QsE",
	}

	testCases := []struct {
		name          string
		modifyMsg     func(*types.MsgBridge)
		setupMocks    func()
		expectedError string
	}{
		{
			name: "Invalid Destination Chain",
			modifyMsg: func(msg *types.MsgBridge) {
				msg.DestinationChain = "invalid:chain"
			},
			setupMocks:    func() {},
			expectedError: "invalid key type",
		},
		{
			name: "Invalid Recipient Address",
			modifyMsg: func(msg *types.MsgBridge) {
				msg.RecipientAddress = "invalid-address"
			},
			setupMocks:    func() {},
			expectedError: "invalid recipient address",
		},
		{
			name: "invalid denom",
			modifyMsg: func(msg *types.MsgBridge) {
				msg.Denom = "noturock"
			},
			setupMocks:    func() {},
			expectedError: "invalid denomination",
		},
		{
			name: "Insufficient Balance",
			modifyMsg: func(msg *types.MsgBridge) {
				msg.Amount = 1000000000000
			},
			setupMocks: func() {
				// Create a test message for mocking
				testMsg := &types.MsgBridge{
					Creator:          "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
					DestinationChain: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
					Amount:           1000000000000,
					Denom:            "urock",
					RecipientAddress: "1BbzosnmC3EVe7XcMgHYd6fUtcfdzUvfeaVZxaZ2QsE",
				}

				// First check balance
				s.bankKeeper.EXPECT().GetBalance(
					s.ctx,
					sdk.MustAccAddressFromBech32(testMsg.Creator),
					testMsg.Denom,
				).Return(sdk.NewCoin(testMsg.Denom, math.NewIntFromUint64(100))).AnyTimes() // Less than required amount

				// Mock SendCoinsFromAccountToModule even though it shouldn't be called
				s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(
					s.ctx,
					sdk.MustAccAddressFromBech32(testMsg.Creator),
					types.ModuleName,
					sdk.NewCoins(sdk.NewCoin(testMsg.Denom, math.NewIntFromUint64(testMsg.Amount+10))),
				).Return(nil).AnyTimes()

				// Mock validation keeper calls
				s.validationKeeper.EXPECT().SetSolanaRequestedNonce(
					s.ctx,
					params.Solana.NonceAccountKey,
					true,
				).Return(nil).AnyTimes()

				s.validationKeeper.EXPECT().SetSolanaZenTPRequestedAccount(
					s.ctx,
					testMsg.RecipientAddress,
					true,
				).Return(nil).AnyTimes()
			},
			expectedError: "not enough balance",
		},
		{
			name: "Supply Cap Exceeded",
			modifyMsg: func(msg *types.MsgBridge) {
				msg.Amount = 200_000_000_000_000 // 200M ROCK - will exceed cap with existing supplies
			},
			setupMocks: func() {
				// To hit the supply cap check, we need to pass the available bridging supply check first
				// Setup: 200M zrchain + 100M solana + 200M new = 500M total
				// But let's increase solana supply to make total > 1B cap
				// Setup: 200M zrchain + 800M solana + 200M new = 1.2B > 1B cap

				// Override solana supply to be much larger
				err := s.zentpKeeper.SetSolanaROCKSupply(s.ctx, math.NewIntFromUint64(800_000_000_000_000)) // 800M ROCK
				s.Require().NoError(err)

				// Create a test message for mocking balance check
				testMsg := &types.MsgBridge{
					Creator:          "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
					DestinationChain: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
					Amount:           200_000_000_000_000,
					Denom:            "urock",
					RecipientAddress: "1BbzosnmC3EVe7XcMgHYd6fUtcfdzUvfeaVZxaZ2QsE",
				}

				// Calculate total amount including bridge fee and Solana fee for balance check
				baseAmountInt := math.NewIntFromUint64(testMsg.Amount)
				bridgeFeeAmount := math.LegacyNewDecFromInt(baseAmountInt).Mul(params.BridgeFee).TruncateInt()
				totalAmountInt := baseAmountInt.Add(bridgeFeeAmount).Add(math.NewIntFromUint64(params.Solana.Fee))

				// Mock sufficient balance for this amount
				s.bankKeeper.EXPECT().GetBalance(
					s.ctx,
					sdk.MustAccAddressFromBech32(testMsg.Creator),
					testMsg.Denom,
				).Return(sdk.NewCoin(testMsg.Denom, totalAmountInt.Add(math.NewIntFromUint64(1000000)))).AnyTimes()
			},
			expectedError: "total ROCK supply including pending would exceed cap",
		},
		{
			name: "Bridge amount exceeds available zrchain supply",
			modifyMsg: func(msg *types.MsgBridge) {
				msg.Amount = 250_000_000_000_000 // 250M ROCK
			},
			setupMocks: func() {
				// Reset Solana supply back to 100M (from previous test that set it to 800M)
				err := s.zentpKeeper.SetSolanaROCKSupply(s.ctx, math.NewIntFromUint64(100_000_000_000_000)) // 100M ROCK
				s.Require().NoError(err)

				// With global setup: zrchain supply 200M, module balance 0M, available = 200M
				// Bridge amount 250M > 200M should trigger the "exceeds available supply" error

				// Create a test message for mocking balance check
				testMsg := &types.MsgBridge{
					Creator:          "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
					DestinationChain: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
					Amount:           250_000_000_000_000,
					Denom:            "urock",
					RecipientAddress: "1BbzosnmC3EVe7XcMgHYd6fUtcfdzUvfeaVZxaZ2QsE",
				}

				// Calculate total amount including bridge fee and Solana fee for balance check
				baseAmountInt := math.NewIntFromUint64(testMsg.Amount)
				bridgeFeeAmount := math.LegacyNewDecFromInt(baseAmountInt).Mul(params.BridgeFee).TruncateInt()
				totalAmountInt := baseAmountInt.Add(bridgeFeeAmount).Add(math.NewIntFromUint64(params.Solana.Fee))

				// Mock sufficient balance so the balance check passes and we hit the supply check
				s.bankKeeper.EXPECT().GetBalance(
					s.ctx,
					sdk.MustAccAddressFromBech32(testMsg.Creator),
					testMsg.Denom,
				).Return(sdk.NewCoin(testMsg.Denom, totalAmountInt.Add(math.NewIntFromUint64(1000000)))).AnyTimes()
			},
			expectedError: "bridge amount 250000000000000 exceeds available zrchain rock supply for bridging 200000000000000",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Create a copy of the base message
			msg := *baseMsg
			if tc.modifyMsg != nil {
				tc.modifyMsg(&msg)
			}

			// Setup mocks
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			// Call the Bridge handler
			response, err := s.msgServer.Bridge(s.ctx, &msg)
			s.Require().Error(err)
			s.Require().Nil(response)
			s.Require().Contains(err.Error(), tc.expectedError)
		})
	}
}
