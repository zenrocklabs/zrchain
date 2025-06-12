package keeper_test

import (
	"cosmossdk.io/math"

	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) TestBurn() {
	// Set up initial balance
	maccbalance := sdk.NewCoin(params.BondDenom, math.NewIntFromUint64(10000))
	moduleAddr := s.accountKeeper.GetModuleAddress(types.ModuleName)

	// Happy Path
	msg := &types.MsgBurn{
		Authority:     s.zentpKeeper.GetAuthority(),
		ModuleAccount: types.ModuleName,
		Amount:        1000,
		Denom:         params.BondDenom,
	}

	// Wrong Authority
	msgWrongAuthority := &types.MsgBurn{
		Authority:     "wrongAuthority",
		ModuleAccount: types.ModuleName,
		Amount:        1000,
		Denom:         params.BondDenom,
	}

	// Wrong Module Account
	msgWrongModuleAccount := &types.MsgBurn{
		Authority:     s.zentpKeeper.GetAuthority(),
		ModuleAccount: "wrongModuleAccount",
		Amount:        1000,
		Denom:         params.BondDenom,
	}

	// Insufficient Balance
	msgInsufficientBalance := &types.MsgBurn{
		Authority:     s.zentpKeeper.GetAuthority(),
		ModuleAccount: types.ModuleName,
		Amount:        100000000000,
		Denom:         params.BondDenom,
	}

	wrongDenom := &types.MsgBurn{
		Authority:     s.zentpKeeper.GetAuthority(),
		ModuleAccount: types.ModuleName,
		Amount:        100000000000,
		Denom:         "wrongDenom",
	}

	// Test cases
	testCases := []struct {
		name          string
		msg           *types.MsgBurn
		setupMocks    func()
		expectedError bool
	}{
		{
			name: "Happy Path",
			msg:  msg,
			setupMocks: func() {
				s.accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(moduleAddr).AnyTimes()
				s.accountKeeper.EXPECT().HasAccount(s.ctx, moduleAddr).Return(true).AnyTimes()
				s.bankKeeper.EXPECT().GetBalance(s.ctx, moduleAddr, params.BondDenom).Return(maccbalance).AnyTimes()
				s.bankKeeper.EXPECT().BurnCoins(s.ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(params.BondDenom, math.NewIntFromUint64(msg.Amount)))).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "Wrong Authority",
			msg:  msgWrongAuthority,
			setupMocks: func() {
				s.accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(moduleAddr).AnyTimes()
				s.accountKeeper.EXPECT().HasAccount(s.ctx, moduleAddr).Return(true).AnyTimes()
			},
			expectedError: true,
		},
		{
			name: "Wrong Module Account",
			msg:  msgWrongModuleAccount,
			setupMocks: func() {
				s.accountKeeper.EXPECT().GetModuleAddress("wrongModuleAccount").Return(sdk.AccAddress{}).AnyTimes()
				s.accountKeeper.EXPECT().HasAccount(s.ctx, sdk.AccAddress{}).Return(false).AnyTimes()
			},
			expectedError: true,
		},
		{
			name: "Insufficient Balance",
			msg:  msgInsufficientBalance,
			setupMocks: func() {
				s.accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(moduleAddr).AnyTimes()
				s.accountKeeper.EXPECT().HasAccount(s.ctx, moduleAddr).Return(true).AnyTimes()
				s.bankKeeper.EXPECT().GetBalance(s.ctx, moduleAddr, params.BondDenom).Return(maccbalance).AnyTimes()
			},
			expectedError: true,
		},
		{
			name: "Wrong Denom",
			msg:  wrongDenom,
			setupMocks: func() {
				s.accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(moduleAddr).AnyTimes()
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.setupMocks()
			response, err := s.msgServer.Burn(s.ctx, tc.msg)
			if tc.expectedError {
				s.Require().Error(err)
				s.Require().Nil(response)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(response)
			}
		})
	}
}
