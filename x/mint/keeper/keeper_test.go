package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	ubermock "go.uber.org/mock/gomock"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/mint"
	"github.com/Zenrock-Foundation/zrchain/v6/x/mint/keeper"
	minttestutil "github.com/Zenrock-Foundation/zrchain/v6/x/mint/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/mint/types"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	zenextypes "github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	zentptypes "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	mintKeeper    keeper.Keeper
	ctx           sdk.Context
	msgServer     types.MsgServer
	stakingKeeper *minttestutil.MockStakingKeeper
	bankKeeper    *minttestutil.MockBankKeeper
	accountKeeper *minttestutil.MockAccountKeeper
	zentpKeeper   *minttestutil.MockZentpKeeper
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
	ctrl := ubermock.NewController(s.T())
	accountKeeper := minttestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := minttestutil.NewMockBankKeeper(ctrl)
	stakingKeeper := minttestutil.NewMockStakingKeeper(ctrl)
	zentpKeeper := minttestutil.NewMockZentpKeeper(ctrl)
	accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(sdk.AccAddress{}).AnyTimes()

	// Assign the mock keepers to the suite fields
	s.accountKeeper = accountKeeper
	s.bankKeeper = bankKeeper
	s.stakingKeeper = stakingKeeper
	s.zentpKeeper = zentpKeeper
	s.mintKeeper = keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		stakingKeeper,
		accountKeeper,
		bankKeeper,
		zentpKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
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
				s.accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(moduleAddr).AnyTimes()
				s.bankKeeper.EXPECT().
					GetBalance(s.ctx, moduleAddr, "urock").
					Return(sdk.NewCoin("urock", math.NewInt(1000))).AnyTimes()
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
					Return(sdk.NewCoin("urock", math.NewInt(100))).AnyTimes()
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

func (s *IntegrationTestSuite) TestClaimTxFees() {
	// Setup test parameters
	params := types.DefaultParams()
	err := s.mintKeeper.Params.Set(s.ctx, params)
	s.Require().NoError(err)

	// Setup expected fees amount
	expectedFees := sdk.NewCoin(params.MintDenom, math.NewInt(500000))

	// Mock getting the fee collector address
	feeCollectorAddr := sdk.AccAddress{}
	s.accountKeeper.EXPECT().
		GetModuleAddress(authtypes.FeeCollectorName).
		Return(feeCollectorAddr)

	// Mock the GetBalance call with the fee collector address
	s.bankKeeper.EXPECT().
		GetBalance(s.ctx, feeCollectorAddr, params.MintDenom).
		Return(expectedFees)

	// Mock the SendCoinsFromModuleToModule call
	s.bankKeeper.EXPECT().
		SendCoinsFromModuleToModule(
			s.ctx,
			authtypes.FeeCollectorName,
			types.ModuleName,
			sdk.NewCoins(expectedFees),
		).
		Return(nil)

	// Call the function being tested
	actualFees, err := s.mintKeeper.ClaimTxFees(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(expectedFees, actualFees)
}

func (s *IntegrationTestSuite) TestBaseDistribution() {
	// Setup context and parameters
	ctx := s.ctx
	totalRewards := sdk.NewCoin("urock", math.NewInt(1000))

	// Set up parameters for the test
	params := types.Params{
		BurnRate:              math.LegacyNewDecWithPrec(20, 2), // 20%
		ProtocolWalletRate:    math.LegacyNewDecWithPrec(30, 2), // 30%
		MintDenom:             "urock",
		ProtocolWalletAddress: "zen1fhln2vnudxddpymqy82vzqhnlsfh4stjd683ze",
	}

	// Set the parameters in the keeper
	err := s.mintKeeper.Params.Set(ctx, params)
	s.Require().NoError(err)

	// Calculate the burn amount
	burnAmount := math.LegacyNewDecFromInt(totalRewards.Amount).Mul(params.BurnRate).TruncateInt()

	// Set up expectation for burning coins
	s.bankKeeper.EXPECT().BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(params.MintDenom, burnAmount))).Return(nil)

	// Calculate the protocol wallet portion
	protocolWalletPortion := math.LegacyNewDecFromInt(totalRewards.Amount).Mul(params.ProtocolWalletRate).TruncateInt()

	// Set up expectation for sending coins to the protocol wallet
	protocolWalletAddr, err := sdk.AccAddressFromBech32(params.ProtocolWalletAddress)
	s.Require().NoError(err)
	s.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(ctx, types.ModuleName, protocolWalletAddr, sdk.NewCoins(sdk.NewCoin(params.MintDenom, protocolWalletPortion))).Return(nil)

	// Call the function being tested
	remainingRewards, err := s.mintKeeper.ZenexFeeProcessing(ctx)
	s.Require().NoError(err)

	// Check the remaining rewards after distribution
	expectedRemaining := sdk.NewCoin("urock", math.NewInt(500)) // 1000 - 200 (burn) - 300 (protocol wallet)
	s.Require().Equal(expectedRemaining, remainingRewards)
}

func (s *IntegrationTestSuite) TestCalculateTopUp() {
	testCases := []struct {
		name                string
		setupMocks          func()
		stakingRewards      sdk.Coin
		totalRewardRest     sdk.Coin
		expectedTopUpAmount sdk.Coin
		expectError         bool
	}{
		{
			name: "successful top-up calculation",
			setupMocks: func() {
				// No specific mocks needed for this case
			},
			stakingRewards:      sdk.NewCoin("urock", math.NewInt(100)),
			totalRewardRest:     sdk.NewCoin("urock", math.NewInt(50)),
			expectedTopUpAmount: sdk.NewCoin("urock", math.NewInt(50)),
			expectError:         false,
		},
		{
			name: "insufficient staking rewards",
			setupMocks: func() {
				// No specific mocks needed for this case
			},
			stakingRewards:      sdk.NewCoin("urock", math.NewInt(50)),
			totalRewardRest:     sdk.NewCoin("urock", math.NewInt(100)),
			expectedTopUpAmount: sdk.Coin{},
			expectError:         true,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Setup the mocks
			tc.setupMocks()

			// Execute the method
			topUpAmount, err := s.mintKeeper.CalculateTopUp(s.ctx, tc.stakingRewards, tc.totalRewardRest)

			// Check results
			if tc.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedTopUpAmount, topUpAmount)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestCalculateExcess() {
	testCases := []struct {
		name                    string
		setupMocks              func()
		totalBlockStakingReward sdk.Coin
		totalRewardsRest        sdk.Coin
		expectedExcess          sdk.Coin
		expectError             bool
	}{
		{
			name: "successful excess calculation",
			setupMocks: func() {
				// No specific mocks needed for this case
			},
			totalBlockStakingReward: sdk.NewCoin("urock", math.NewInt(50)),
			totalRewardsRest:        sdk.NewCoin("urock", math.NewInt(100)),
			expectedExcess:          sdk.NewCoin("urock", math.NewInt(50)),
			expectError:             false,
		},
		{
			name: "zero excess",
			setupMocks: func() {
				// No specific mocks needed for this case
			},
			totalBlockStakingReward: sdk.NewCoin("urock", math.NewInt(100)),
			totalRewardsRest:        sdk.NewCoin("urock", math.NewInt(100)),
			expectedExcess:          sdk.Coin{},
			expectError:             true,
		},
		{
			name: "negative excess",

			setupMocks: func() {
				// No specific mocks needed for this case
			},
			totalBlockStakingReward: sdk.NewCoin("urock", math.NewInt(150)),
			totalRewardsRest:        sdk.NewCoin("urock", math.NewInt(100)),
			expectedExcess:          sdk.Coin{},
			expectError:             true,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Setup the mocks
			tc.setupMocks()

			// Execute the method
			excess, err := s.mintKeeper.CalculateExcess(s.ctx, tc.totalBlockStakingReward, tc.totalRewardsRest)

			// Check results
			if tc.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedExcess, excess)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestAdditionalBurn() {
	// Setup context and parameters
	ctx := s.ctx
	totalExcess := sdk.NewCoin("urock", math.NewInt(100))

	// Set up parameters for the test
	params := types.Params{
		MintDenom:          "urock",
		AdditionalBurnRate: math.LegacyNewDecWithPrec(10, 2), // 10%
	}

	// Set the parameters in the keeper
	err := s.mintKeeper.Params.Set(ctx, params)
	s.Require().NoError(err)

	// Calculate the burn amount
	burnAmount := math.LegacyNewDecFromInt(totalExcess.Amount).Mul(params.AdditionalBurnRate).TruncateInt()

	// Set up expectation for burning coins
	s.bankKeeper.EXPECT().BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(params.MintDenom, burnAmount))).Return(nil)

	// Call the function being tested
	err = s.mintKeeper.AdditionalBurn(ctx, totalExcess)
	s.Require().NoError(err)

	// Check that the excess amount is reduced correctly
	expectedExcess := totalExcess.Amount.Sub(burnAmount)
	s.Require().Equal(expectedExcess, totalExcess.Amount.Sub(burnAmount))
}

func (s *IntegrationTestSuite) TestAdditionalMpcRewards() {
	// Setup context and parameters
	ctx := s.ctx
	totalExcess := sdk.NewCoin("urock", math.NewInt(100))

	// Set up parameters for the test
	params := types.Params{
		MintDenom:             "urock",
		ProtocolWalletAddress: "zen1fhln2vnudxddpymqy82vzqhnlsfh4stjd683ze",
		AdditionalMpcRewards:  math.LegacyNewDecWithPrec(20, 2), // 20%
	}

	// Set the parameters in the keeper
	err := s.mintKeeper.Params.Set(ctx, params)
	s.Require().NoError(err)

	// Calculate the MPC rewards amount
	mpcRewards := math.LegacyNewDecFromInt(totalExcess.Amount).Mul(params.AdditionalMpcRewards).TruncateInt()

	// Convert string address to AccAddress
	protocolAddr, err := sdk.AccAddressFromBech32(params.ProtocolWalletAddress)
	s.Require().NoError(err)

	// Set up expectation for sending coins to the protocol wallet
	s.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(ctx, types.ModuleName, protocolAddr, sdk.NewCoins(sdk.NewCoin(params.MintDenom, mpcRewards))).Return(nil)

	// Call the function being tested
	err = s.mintKeeper.AdditionalMpcRewards(ctx, totalExcess)
	s.Require().NoError(err)

	// Check that the excess amount is reduced correctly
	expectedExcess := totalExcess.Amount.Sub(mpcRewards)
	s.Require().Equal(expectedExcess, totalExcess.Amount.Sub(mpcRewards))
}

func (s *IntegrationTestSuite) TestAdditionalStakingRewards() {
	// Setup context and parameters
	ctx := s.ctx
	totalExcess := sdk.NewCoin("urock", math.NewInt(100))

	// Set up parameters for the test
	params := types.Params{
		MintDenom:                "urock",
		AdditionalStakingRewards: math.LegacyNewDecWithPrec(15, 2), // 15%
	}

	// Set the parameters in the keeper
	err := s.mintKeeper.Params.Set(ctx, params)
	s.Require().NoError(err)

	// Calculate the staking rewards amount
	stakingRewards := math.LegacyNewDecFromInt(totalExcess.Amount).Mul(params.AdditionalStakingRewards).TruncateInt()

	// Set up test cases
	testCases := []struct {
		name          string
		setupMocks    func()
		excess        sdk.Coin
		expectedError bool
	}{
		{
			name: "successful staking rewards distribution",
			setupMocks: func() {
				// Set up expectation for sending coins to the staking module
				s.bankKeeper.EXPECT().SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, sdk.NewCoins(sdk.NewCoin(params.MintDenom, stakingRewards))).Return(nil)
			},
			excess:        totalExcess,
			expectedError: false,
		},
		{
			name: "params get error",
			setupMocks: func() {
				err := s.mintKeeper.Params.Remove(ctx)
				s.Require().NoError(err)
				// No expectation for sending coins since it should error out
			},
			excess:        totalExcess,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Setup the mocks
			tc.setupMocks()

			// Call the function being tested
			err := s.mintKeeper.AdditionalStakingRewards(ctx, tc.excess)

			// Check results
			if tc.expectedError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestExcessDistribution() {
	// Setup context and parameters
	ctx := s.ctx
	excessAmount := sdk.NewCoin("urock", math.NewInt(1000))

	// Set up parameters for the test
	params := types.Params{
		MintDenom:                "urock",
		AdditionalBurnRate:       math.LegacyNewDecWithPrec(10, 2), // 10%
		AdditionalMpcRewards:     math.LegacyNewDecWithPrec(20, 2), // 20%
		AdditionalStakingRewards: math.LegacyNewDecWithPrec(15, 2), // 15%
		ProtocolWalletAddress:    "zen1fhln2vnudxddpymqy82vzqhnlsfh4stjd683ze",
	}

	// Set the parameters in the keeper
	err := s.mintKeeper.Params.Set(ctx, params)
	s.Require().NoError(err)

	// Calculate the burn amount
	burnAmount := math.LegacyNewDecFromInt(excessAmount.Amount).Mul(params.AdditionalBurnRate).TruncateInt()

	// Set up expectation for burning coins
	s.bankKeeper.EXPECT().BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(params.MintDenom, burnAmount))).Return(nil)

	// Calculate the MPC rewards amount
	mpcRewards := math.LegacyNewDecFromInt(excessAmount.Amount).Mul(params.AdditionalMpcRewards).TruncateInt()

	// Convert string address to AccAddress
	protocolWalletAddr, err := sdk.AccAddressFromBech32(params.ProtocolWalletAddress)
	s.Require().NoError(err)

	// Set up expectation for sending coins to the protocol wallet
	s.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(ctx, types.ModuleName, protocolWalletAddr, sdk.NewCoins(sdk.NewCoin(params.MintDenom, mpcRewards))).Return(nil)

	// Calculate the staking rewards amount
	stakingRewards := math.LegacyNewDecFromInt(excessAmount.Amount).Mul(params.AdditionalStakingRewards).TruncateInt()

	// Set up expectation for sending coins to the staking module
	s.bankKeeper.EXPECT().SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, sdk.NewCoins(sdk.NewCoin(params.MintDenom, stakingRewards))).Return(nil)

	// Call the function being tested
	err = s.mintKeeper.ExcessDistribution(ctx, excessAmount)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestGetMintModuleBalance() {
	// Setup context
	ctx := s.ctx

	// Set up parameters for the test
	params := types.Params{
		MintDenom: "urock",
	}

	// Set the parameters in the keeper
	err := s.mintKeeper.Params.Set(ctx, params)
	s.Require().NoError(err)

	// Define the expected balance for the mint module
	expectedBalance := sdk.NewCoin("urock", math.NewInt(500))

	// Mock the module address
	mintModuleAddr := sdk.AccAddress{} // This should be the expected address

	// Set up expectation for getting the module address
	s.accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(mintModuleAddr)

	// Set up expectation for getting the balance
	s.bankKeeper.EXPECT().GetBalance(ctx, mintModuleAddr, params.MintDenom).Return(expectedBalance)

	// Call the function being tested
	actualBalance, err := s.mintKeeper.GetMintModuleBalance(ctx)
	s.Require().NoError(err)

	// Check that the actual balance matches the expected balance
	s.Require().Equal(expectedBalance, actualBalance)
}

func (s *IntegrationTestSuite) TestZentpFeesDistribution() {
	tests := []struct {
		name             string
		zentpRockBalance uint64
		errExpected      bool
		expectedErr      error
	}{
		{
			name:             "happy path",
			zentpRockBalance: 2000,
			errExpected:      false,
			expectedErr:      nil,
		},
		{
			name:             "balance is zero",
			zentpRockBalance: 0,
			errExpected:      false,
			expectedErr:      nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {

			s.SetupTest()

			s.accountKeeper.EXPECT().GetModuleAddress(zentptypes.ZentpCollectorName).Return(sdk.MustAccAddressFromBech32("zen1234wz2aaavp089ttnrj9jwjqraaqxkkadq0k03")).AnyTimes()
			s.bankKeeper.EXPECT().GetBalance(s.ctx, sdk.MustAccAddressFromBech32("zen1234wz2aaavp089ttnrj9jwjqraaqxkkadq0k03"), types.DefaultParams().MintDenom).Return(sdk.NewCoin(types.DefaultParams().MintDenom, math.NewIntFromUint64(tt.zentpRockBalance))).AnyTimes()
			s.bankKeeper.EXPECT().SendCoinsFromModuleToModule(s.ctx, zentptypes.ZentpCollectorName, zenextypes.ZenexFeeCollectorName, sdk.NewCoins(sdk.NewCoin(types.DefaultParams().MintDenom, math.NewIntFromUint64(tt.zentpRockBalance)))).Return(nil).AnyTimes()
			s.zentpKeeper.EXPECT().UpdateZentpFees(s.ctx, tt.zentpRockBalance).Return(nil).AnyTimes()

			err := s.mintKeeper.DistributeZentpFeesToZenexFeeCollector(s.ctx)

			if tt.errExpected {
				s.Require().Error(err)
				s.Require().Equal(tt.expectedErr.Error(), err.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
