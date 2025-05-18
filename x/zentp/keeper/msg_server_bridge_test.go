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
		SourceAddress:    "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
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
	).Return(sdk.NewCoin(msg.Denom, math.NewIntFromUint64(msg.Amount*20)))

	burnAmount := math.LegacyNewDecFromInt(math.NewIntFromUint64(msg.Amount)).Mul(params.BridgeFee).TruncateInt()

	// Mock bank keeper SendCoinsFromAccountToModule
	s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(
		s.ctx,
		sdk.MustAccAddressFromBech32(msg.SourceAddress),
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin("urock", math.NewIntFromUint64(msg.Amount).Add(burnAmount))),
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
		SourceAddress:    "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		DestinationChain: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		Amount:           1000,
		Denom:            "urock",
		RecipientAddress: "1BbzosnmC3EVe7XcMgHYd6fUtcfdzUvfeaVZxaZ2QsE",
	}

	testCases := []struct {
		name          string
		modifyMsg     func(*types.MsgBridge)
		expectedError string
	}{
		{
			name: "Invalid Destination Chain",
			modifyMsg: func(msg *types.MsgBridge) {
				msg.DestinationChain = "invalid:chain"
			},
			expectedError: "invalid key type",
		},
		{
			name: "Invalid Recipient Address",
			modifyMsg: func(msg *types.MsgBridge) {
				msg.RecipientAddress = "invalid-address"
			},
			expectedError: "invalid recipient address",
		},
		{
			name: "Insufficient Balance",
			modifyMsg: func(msg *types.MsgBridge) {
				msg.Amount = 1000000000000
				// First check balance
				s.bankKeeper.EXPECT().GetBalance(
					s.ctx,
					sdk.MustAccAddressFromBech32(msg.Creator),
					msg.Denom,
				).Return(sdk.NewCoin(msg.Denom, math.NewIntFromUint64(100))).AnyTimes() // Less than required amount

				// Mock SendCoinsFromAccountToModule even though it shouldn't be called
				s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(
					s.ctx,
					sdk.MustAccAddressFromBech32(msg.SourceAddress),
					types.ModuleName,
					sdk.NewCoins(sdk.NewCoin(msg.Denom, math.NewIntFromUint64(msg.Amount+10))),
				).Return(nil).AnyTimes()

				// Mock validation keeper calls
				s.validationKeeper.EXPECT().SetSolanaRequestedNonce(
					s.ctx,
					params.Solana.NonceAccountKey,
					true,
				).Return(nil).AnyTimes()

				s.validationKeeper.EXPECT().SetSolanaZenTPRequestedAccount(
					s.ctx,
					msg.RecipientAddress,
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

			// Call the Bridge handler
			response, err := s.msgServer.Bridge(s.ctx, &msg)
			s.Require().Error(err)
			s.Require().Nil(response)
			s.Require().Contains(err.Error(), tc.expectedError)
		})
	}
}
