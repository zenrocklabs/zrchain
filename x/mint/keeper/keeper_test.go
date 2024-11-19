package keeper_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/Zenrock-Foundation/zrchain/v5/x/mint"
	"github.com/Zenrock-Foundation/zrchain/v5/x/mint/keeper"
	minttestutil "github.com/Zenrock-Foundation/zrchain/v5/x/mint/testutil"
	"github.com/Zenrock-Foundation/zrchain/v5/x/mint/types"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	mintKeeper     keeper.Keeper
	ctx            sdk.Context
	msgServer      types.MsgServer
	stakingKeeper  *minttestutil.MockStakingKeeper
	bankKeeper     *minttestutil.MockBankKeeper
	treasuryKeeper *minttestutil.MockTreasuryKeeper
	accountKeeper  *minttestutil.MockAccountKeeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupTest() {
	encCfg := moduletestutil.MakeTestEncodingConfig(mint.AppModuleBasic{})
	key := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx

	// gomock initializations
	ctrl := gomock.NewController(s.T())
	accountKeeper := minttestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := minttestutil.NewMockBankKeeper(ctrl)
	stakingKeeper := minttestutil.NewMockStakingKeeper(ctrl)
	treasuryKeeper := minttestutil.NewMockTreasuryKeeper(ctrl)
	accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(sdk.AccAddress{})

	// Assign the mock keepers to the suite fields
	s.accountKeeper = accountKeeper
	s.bankKeeper = bankKeeper
	s.stakingKeeper = stakingKeeper
	s.treasuryKeeper = treasuryKeeper

	s.mintKeeper = keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		stakingKeeper,
		accountKeeper,
		bankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		treasuryKeeper,
	)

	s.Require().Equal(testCtx.Ctx.Logger().With("module", "x/"+types.ModuleName),
		s.mintKeeper.Logger(testCtx.Ctx))

	err := s.mintKeeper.Params.Set(s.ctx, types.DefaultParams())
	s.Require().NoError(err)
	s.Require().NoError(s.mintKeeper.Minter.Set(s.ctx, types.DefaultInitialMinter()))

	s.msgServer = keeper.NewMsgServerImpl(s.mintKeeper)
}

func (s *IntegrationTestSuite) TestAliasFunctions() {

	stakingTokenSupply := math.NewIntFromUint64(100000000000)
	s.stakingKeeper.EXPECT().StakingTokenSupply(s.ctx).Return(stakingTokenSupply, nil)
	tokenSupply, err := s.mintKeeper.StakingTokenSupply(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(tokenSupply, stakingTokenSupply)

	bondedRatio := math.LegacyNewDecWithPrec(15, 2)
	s.stakingKeeper.EXPECT().BondedRatio(s.ctx).Return(bondedRatio, nil)
	ratio, err := s.mintKeeper.BondedRatio(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(ratio, bondedRatio)

	coins := sdk.NewCoins(sdk.NewCoin("urock", math.NewInt(1000000)))
	s.bankKeeper.EXPECT().MintCoins(s.ctx, types.ModuleName, coins).Return(nil)
	s.Require().Equal(s.mintKeeper.MintCoins(s.ctx, sdk.NewCoins()), nil)
	s.Require().Nil(s.mintKeeper.MintCoins(s.ctx, coins))

	fees := sdk.NewCoins(sdk.NewCoin("urock", math.NewInt(1000)))
	s.bankKeeper.EXPECT().SendCoinsFromModuleToModule(s.ctx, types.ModuleName, authtypes.FeeCollectorName, fees).Return(nil)
	s.Require().Nil(s.mintKeeper.AddCollectedFees(s.ctx, fees))
}

func (s *IntegrationTestSuite) TestClaimKeyringFees() {
	// Setup test parameters
	params := types.DefaultParams()
	err := s.mintKeeper.Params.Set(s.ctx, params)
	s.Require().NoError(err)

	// Setup expected keyring rewards
	expectedRewards := sdk.NewCoin(params.MintDenom, math.NewInt(1000000))

	// Mock getting the module account address
	moduleAddr := sdk.AccAddress{}
	s.accountKeeper.EXPECT().
		GetModuleAddress(treasurytypes.KeyringCollectorName).
		Return(moduleAddr)

	// Mock the GetBalance call with the module account address
	s.bankKeeper.EXPECT().
		GetBalance(s.ctx, moduleAddr, params.MintDenom).
		Return(expectedRewards)

	// Mock the SendCoinsFromModuleToModule call
	s.bankKeeper.EXPECT().
		SendCoinsFromModuleToModule(
			s.ctx,
			treasurytypes.KeyringCollectorName,
			types.ModuleName,
			sdk.NewCoins(expectedRewards),
		).
		Return(nil)

	// Call the function being tested
	actualRewards, err := s.mintKeeper.ClaimKeyringFees(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(expectedRewards, actualRewards)
}

// func (s *IntegrationTestSuite) TestTopUpKeyringRewards() {
// 	// Setup test parameters
// 	params := types.DefaultParams()
// 	err := s.mintKeeper.Params.Set(s.ctx, params)
// 	s.Require().NoError(err)

// 	// Setup test amount
// 	topUpAmount := sdk.NewCoin(params.MintDenom, math.NewInt(1000000))

// 	// Convert the protocol wallet address string to AccAddress before using it
// 	protocolAddr, err := sdk.AccAddressFromBech32(params.ProtocolWalletAddress)
// 	s.Require().NoError(err)

// 	// Mock the SendCoinsFromAccountToModule call with the correct address
// 	s.bankKeeper.EXPECT().
// 		SendCoinsFromAccountToModule(
// 			s.ctx,
// 			protocolAddr,
// 			types.ModuleName,
// 			sdk.NewCoins(topUpAmount),
// 		).
// 		Return(nil)

// 	// Call the function being tested
// 	err = s.mintKeeper.TopUpTotalRewards(s.ctx, topUpAmount)
// 	s.Require().NoError(err)
// }

func (s *IntegrationTestSuite) TestCheckModuleBalance() {
	testCases := []struct {
		name        string
		setupMocks  func()
		reward      sdk.Coin
		expectError bool
	}{
		{
			name: "sufficient balance",
			setupMocks: func() {
				moduleAddr := sdk.AccAddress{}
				s.accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(moduleAddr)
				s.bankKeeper.EXPECT().
					GetBalance(s.ctx, moduleAddr, "urock").
					Return(sdk.NewCoin("urock", math.NewInt(1000)))
			},
			reward:      sdk.NewCoin("urock", math.NewInt(500)),
			expectError: false,
		},
		{
			name: "insufficient balance",
			setupMocks: func() {
				moduleAddr := sdk.AccAddress{}
				s.accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(moduleAddr)
				s.bankKeeper.EXPECT().
					GetBalance(s.ctx, moduleAddr, "urock").
					Return(sdk.NewCoin("urock", math.NewInt(100)))
			},
			reward:      sdk.NewCoin("urock", math.NewInt(500)),
			expectError: true,
		},
		{
			name:        "zero reward amount",
			setupMocks:  func() {},
			reward:      sdk.NewCoin("urock", math.NewInt(0)),
			expectError: true,
		},
		{
			name: "module address not found",
			setupMocks: func() {
				s.accountKeeper.EXPECT().
					GetModuleAddress(types.ModuleName).
					Return(nil)
			},
			reward:      sdk.NewCoin("urock", math.NewInt(500)),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Setup the mocks
			tc.setupMocks()

			// Execute the method
			err := s.mintKeeper.CheckModuleBalance(s.ctx, tc.reward)

			// Check results
			if tc.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestTotalBondedTokens() {
	testCases := []struct {
		name          string
		setupMocks    func()
		expectedValue math.Int
		expectError   bool
	}{
		{
			name: "successful query",
			setupMocks: func() {
				s.stakingKeeper.EXPECT().
					TotalBondedTokens(s.ctx).
					Return(math.NewInt(1000000), nil)
			},
			expectedValue: math.NewInt(1000000),
			expectError:   false,
		},
		{
			name: "staking keeper returns error",
			setupMocks: func() {
				s.stakingKeeper.EXPECT().
					TotalBondedTokens(s.ctx).
					Return(math.Int{}, fmt.Errorf("failed to get total bonded tokens"))
			},
			expectedValue: math.Int{},
			expectError:   true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Setup the mocks
			tc.setupMocks()

			// Execute the method
			result, err := s.mintKeeper.TotalBondedTokens(s.ctx)

			// Check results
			if tc.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedValue, result)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestNextStakingReward() {
	testCases := []struct {
		name string

		setupMocks     func()
		totalBonded    math.Int
		expectedReward sdk.Coin
		expectError    bool
	}{
		{
			name: "successful calculation",
			setupMocks: func() {
				params := types.Params{
					MintDenom:     "urock",
					BlocksPerYear: 6311520,
					StakingYield:  math.LegacyNewDecWithPrec(15, 2),
				}
				err := s.mintKeeper.Params.Set(s.ctx, params)
				s.Require().NoError(err)
			},
			totalBonded:    math.NewInt(1000000000),
			expectedReward: sdk.NewCoin("urock", math.NewInt(23)),
			expectError:    false,
		},
		{
			name: "params get error",
			setupMocks: func() {
				err := s.mintKeeper.Params.Remove(s.ctx)
				s.Require().NoError(err)
			},
			totalBonded:    math.NewInt(1000000000),
			expectedReward: sdk.Coin{},
			expectError:    true,
		},
		{
			name: "zero bonded tokens",
			setupMocks: func() {
				params := types.Params{
					MintDenom:     "urock",
					BlocksPerYear: 6311520,
					StakingYield:  math.LegacyNewDecWithPrec(15, 2),
				}
				err := s.mintKeeper.Params.Set(s.ctx, params)
				s.Require().NoError(err)
			},
			totalBonded:    math.NewInt(0),
			expectedReward: sdk.NewCoin("urock", math.NewInt(0)),
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Setup the mocks
			tc.setupMocks()

			// Execute the method
			reward, err := s.mintKeeper.NextStakingReward(s.ctx, tc.totalBonded)

			// Check results
			if tc.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedReward, reward)
			}
		})
	}
}
