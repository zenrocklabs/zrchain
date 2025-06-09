package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/suite"

	minttypes "github.com/Zenrock-Foundation/zrchain/v6/x/mint/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestBridge() {
	// Setup test parameters
	params := types.DefaultParams()
	err := s.zentpKeeper.ParamStore.Set(s.ctx, params)
	s.Require().NoError(err)

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

	// Mock getting the mint params
	s.mintKeeper.EXPECT().GetParams(s.ctx).Return(minttypes.DefaultParams(), nil)

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
