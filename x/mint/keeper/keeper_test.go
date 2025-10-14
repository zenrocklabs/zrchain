package keeper_test

import (
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
	zenexKeeper   *minttestutil.MockZenexKeeper
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
	zenexKeeper := minttestutil.NewMockZenexKeeper(ctrl)
	accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(sdk.MustAccAddressFromBech32("zen1m3h30wlvsf8llruxtpukdvsy0km2kum8ju4et3")).AnyTimes()

	// Assign the mock keepers to the suite fields
	s.accountKeeper = accountKeeper
	s.bankKeeper = bankKeeper
	s.stakingKeeper = stakingKeeper
	s.zentpKeeper = zentpKeeper
	s.zenexKeeper = zenexKeeper
	s.mintKeeper = keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		stakingKeeper,
		accountKeeper,
		bankKeeper,
		zentpKeeper,
		zenexKeeper,
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
	tests := []struct {
		name            string
		keyringBalance  uint64
		expectedRewards uint64
		expectError     bool
	}{
		{
			name:            "successful claim",
			keyringBalance:  1000000,
			expectedRewards: 1000000,
			expectError:     false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()

			// Setup test parameters
			params := types.DefaultParams()
			s.Require().NoError(s.mintKeeper.Params.Set(s.ctx, params))

			s.accountKeeper.EXPECT().GetModuleAddress(treasurytypes.KeyringCollectorName).Return(sdk.MustAccAddressFromBech32("zen1a2q8gs3x4c0q9a99hapl8p9tyhxhhwvetrvure")).AnyTimes()
			s.accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(sdk.MustAccAddressFromBech32("zen1m3h30wlvsf8llruxtpukdvsy0km2kum8ju4et3")).AnyTimes()

			s.bankKeeper.EXPECT().GetBalance(s.ctx, sdk.MustAccAddressFromBech32("zen1a2q8gs3x4c0q9a99hapl8p9tyhxhhwvetrvure"), params.MintDenom).Return(sdk.NewCoin(params.MintDenom, math.NewIntFromUint64(tt.keyringBalance)))
			s.bankKeeper.EXPECT().SendCoinsFromModuleToModule(s.ctx, treasurytypes.KeyringCollectorName, types.ModuleName, sdk.NewCoins(sdk.NewCoin(params.MintDenom, math.NewIntFromUint64(tt.keyringBalance)))).Return(nil)

			// Call the function being tested
			actualRewards, err := s.mintKeeper.ClaimKeyringFees(s.ctx)
			if tt.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(sdk.NewCoin(params.MintDenom, math.NewIntFromUint64(tt.expectedRewards)), actualRewards)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestCheckModuleBalance() {

	testCases := []struct {
		name        string
		reward      sdk.Coin
		expectError bool
	}{
		{
			name:        "sufficient balance",
			reward:      sdk.NewCoin("urock", math.NewInt(500)),
			expectError: false,
		},
		{
			name:        "insufficient balance",
			reward:      sdk.NewCoin("urock", math.NewInt(500)),
			expectError: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()

			s.accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(sdk.MustAccAddressFromBech32("zen1m3h30wlvsf8llruxtpukdvsy0km2kum8ju4et3")).AnyTimes()
			s.bankKeeper.EXPECT().GetBalance(s.ctx, sdk.MustAccAddressFromBech32("zen1m3h30wlvsf8llruxtpukdvsy0km2kum8ju4et3"), tc.reward.Denom).Return(tc.reward).AnyTimes()

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
		expectedValue math.Int
		expectError   bool
	}{
		{
			name:          "successful query",
			expectedValue: math.NewInt(1000000),
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			s.stakingKeeper.EXPECT().TotalBondedTokens(s.ctx).Return(tc.expectedValue, nil)

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
	type args struct {
		totalBonded  math.Int
		stakingYield math.LegacyDec
	}

	testCases := []struct {
		name           string
		args           args
		expectedReward sdk.Coin
		expectError    bool
	}{
		{
			name: "successful calculation",
			args: args{
				totalBonded:  math.NewInt(1000000000),
				stakingYield: math.LegacyNewDecWithPrec(15, 2),
			},
			expectedReward: sdk.NewCoin("urock", math.NewInt(23)),
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			s.mintKeeper.Params.Set(s.ctx, types.Params{
				MintDenom:     "urock",
				BlocksPerYear: 6311520,
				StakingYield:  tc.args.stakingYield,
			})

			// Execute the method
			reward, err := s.mintKeeper.NextStakingReward(s.ctx, tc.args.totalBonded)

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

	testCases := []struct {
		name         string
		expectedFees sdk.Coin
		expectError  bool
	}{
		{
			name:         "successful claim",
			expectedFees: sdk.NewCoin("urock", math.NewInt(500000)),
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()

			// Setup test parameters
			params := types.DefaultParams()
			err := s.mintKeeper.Params.Set(s.ctx, params)
			s.Require().NoError(err)

			s.accountKeeper.EXPECT().GetModuleAddress(authtypes.FeeCollectorName).Return(sdk.MustAccAddressFromBech32("zen1m3h30wlvsf8llruxtpukdvsy0km2kum8ju4et3")).AnyTimes()
			s.bankKeeper.EXPECT().GetBalance(s.ctx, sdk.MustAccAddressFromBech32("zen1m3h30wlvsf8llruxtpukdvsy0km2kum8ju4et3"), params.MintDenom).Return(tc.expectedFees).AnyTimes()
			s.bankKeeper.EXPECT().SendCoinsFromModuleToModule(s.ctx, authtypes.FeeCollectorName, types.ModuleName, sdk.NewCoins(tc.expectedFees)).Return(nil)

			// Call the function being tested
			actualFees, err := s.mintKeeper.ClaimTxFees(s.ctx)
			if tc.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedFees, actualFees)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestZenexFeeProcessing() {
	testCases := []struct {
		name            string
		rewards         sdk.Coin
		leftoverBalance sdk.Coin
		expectedFees    sdk.Coin
		expectError     bool
	}{
		{
			name:            "successful processing",
			rewards:         sdk.NewCoin("urock", math.NewInt(1000)),
			expectedFees:    sdk.NewCoin("urock", math.NewInt(350)),
			leftoverBalance: sdk.NewCoin("urock", math.NewInt(300)),
			expectError:     false,
		},
		{
			name:            "low rewards",
			rewards:         sdk.NewCoin("urock", math.NewInt(1)),
			expectedFees:    sdk.NewCoin("urock", math.NewInt(0)),
			leftoverBalance: sdk.NewCoin("urock", math.NewInt(1)),
			expectError:     false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()

			params := types.DefaultParams()
			params.MintDenom = "urock" // Override to use urock instead of stake
			err := s.mintKeeper.Params.Set(s.ctx, params)
			s.Require().NoError(err)

			// need to change the address for the actual zenexfee collector address
			s.accountKeeper.EXPECT().GetModuleAddress(zenextypes.ZenexFeeCollectorName).Return(sdk.MustAccAddressFromBech32("zen1m3h30wlvsf8llruxtpukdvsy0km2kum8ju4et3")).AnyTimes()
			s.bankKeeper.EXPECT().GetBalance(s.ctx, sdk.MustAccAddressFromBech32("zen1m3h30wlvsf8llruxtpukdvsy0km2kum8ju4et3"), "urock").Return(tc.rewards).AnyTimes()
			s.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(s.ctx, zenextypes.ZenexFeeCollectorName, sdk.MustAccAddressFromBech32("zen1fhln2vnudxddpymqy82vzqhnlsfh4stjd683ze"), sdk.NewCoins(tc.expectedFees)).Return(nil)
			s.accountKeeper.EXPECT().GetModuleAddress(zenextypes.ZenBtcRewardsCollectorName).Return(sdk.MustAccAddressFromBech32("zen1fhln2vnudxddpymqy82vzqhnlsfh4stjd683ze")).AnyTimes()
			s.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(s.ctx, zenextypes.ZenexFeeCollectorName, sdk.MustAccAddressFromBech32("zen1fhln2vnudxddpymqy82vzqhnlsfh4stjd683ze"), sdk.NewCoins(tc.expectedFees)).Return(nil)

			// Call the function being tested
			remainingRewards, err := s.mintKeeper.ZenexFeeProcessing(s.ctx)
			if tc.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.leftoverBalance, remainingRewards)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestCalculateTopUp() {
	testCases := []struct {
		name                string
		stakingRewards      sdk.Coin
		totalRewardRest     sdk.Coin
		expectedTopUpAmount sdk.Coin
		expectError         bool
	}{
		{
			name:                "successful top-up calculation",
			stakingRewards:      sdk.NewCoin("urock", math.NewInt(100)),
			totalRewardRest:     sdk.NewCoin("urock", math.NewInt(50)),
			expectedTopUpAmount: sdk.NewCoin("urock", math.NewInt(50)),
			expectError:         false,
		},
		{
			name:                "insufficient staking rewards",
			stakingRewards:      sdk.NewCoin("urock", math.NewInt(50)),
			totalRewardRest:     sdk.NewCoin("urock", math.NewInt(100)),
			expectedTopUpAmount: sdk.Coin{},
			expectError:         true,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

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
			name:                    "successful excess calculation",
			totalBlockStakingReward: sdk.NewCoin("urock", math.NewInt(50)),
			totalRewardsRest:        sdk.NewCoin("urock", math.NewInt(100)),
			expectedExcess:          sdk.NewCoin("urock", math.NewInt(50)),
			expectError:             false,
		},
		{
			name:                    "zero excess",
			totalBlockStakingReward: sdk.NewCoin("urock", math.NewInt(100)),
			totalRewardsRest:        sdk.NewCoin("urock", math.NewInt(100)),
			expectedExcess:          sdk.Coin{},
			expectError:             true,
		},
		{
			name:                    "negative excess",
			totalBlockStakingReward: sdk.NewCoin("urock", math.NewInt(150)),
			totalRewardsRest:        sdk.NewCoin("urock", math.NewInt(100)),
			expectedExcess:          sdk.Coin{},
			expectError:             true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

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

	testCases := []struct {
		name               string
		totalExcess        sdk.Coin
		additionalBurnRate math.LegacyDec
		expectedBurnAmount math.Int
		expectError        bool
	}{
		{
			name:               "successful burn calculation",
			totalExcess:        sdk.NewCoin("urock", math.NewInt(100)),
			additionalBurnRate: math.LegacyNewDecWithPrec(10, 2),
			expectedBurnAmount: math.NewInt(10),
			expectError:        false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			s.SetupTest()

			params := types.Params{
				MintDenom:          "urock",
				AdditionalBurnRate: tc.additionalBurnRate,
			}
			err := s.mintKeeper.Params.Set(s.ctx, params)
			s.Require().NoError(err)

			burnAmount := math.LegacyNewDecFromInt(tc.totalExcess.Amount).Mul(tc.additionalBurnRate).TruncateInt()

			s.bankKeeper.EXPECT().BurnCoins(s.ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(params.MintDenom, burnAmount))).Return(nil)

			// Execute the method
			err = s.mintKeeper.AdditionalBurn(s.ctx, tc.totalExcess)

			// Check results
			if tc.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedBurnAmount, burnAmount)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestAdditionalMpcRewards() {

	testCases := []struct {
		name                 string
		totalExcess          sdk.Coin
		additionalMpcRewards math.LegacyDec
		expectedMpcRewards   math.Int
		expectError          bool
	}{
		{
			name:                 "successful mpc rewards calculation",
			totalExcess:          sdk.NewCoin("urock", math.NewInt(100)),
			additionalMpcRewards: math.LegacyNewDecWithPrec(20, 2),
			expectedMpcRewards:   math.NewInt(20),
			expectError:          false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			s.SetupTest()

			params := types.Params{
				MintDenom:             "urock",
				AdditionalMpcRewards:  tc.additionalMpcRewards,
				ProtocolWalletAddress: "zen1fhln2vnudxddpymqy82vzqhnlsfh4stjd683ze",
			}
			err := s.mintKeeper.Params.Set(s.ctx, params)
			s.Require().NoError(err)

			mpcRewards := math.LegacyNewDecFromInt(tc.totalExcess.Amount).Mul(tc.additionalMpcRewards).TruncateInt()

			s.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(s.ctx, types.ModuleName, sdk.MustAccAddressFromBech32(params.ProtocolWalletAddress), sdk.NewCoins(sdk.NewCoin(params.MintDenom, mpcRewards))).Return(nil)

			// Execute the method
			err = s.mintKeeper.AdditionalMpcRewards(s.ctx, tc.totalExcess)

			// Check results
			if tc.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedMpcRewards, mpcRewards)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestAdditionalStakingRewards() {

	testCases := []struct {
		name                     string
		totalExcess              sdk.Coin
		AdditionalStakingRewards math.LegacyDec
		expectedStakingRewards   math.Int
		expectError              bool
	}{
		{
			name:                     "successful mpc rewards calculation",
			totalExcess:              sdk.NewCoin("urock", math.NewInt(100)),
			AdditionalStakingRewards: math.LegacyNewDecWithPrec(20, 2),
			expectedStakingRewards:   math.NewInt(20),
			expectError:              false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			s.SetupTest()

			params := types.Params{
				MintDenom:                "urock",
				AdditionalStakingRewards: tc.AdditionalStakingRewards,
				ProtocolWalletAddress:    "zen1fhln2vnudxddpymqy82vzqhnlsfh4stjd683ze",
			}
			err := s.mintKeeper.Params.Set(s.ctx, params)
			s.Require().NoError(err)

			stakingRewards := math.LegacyNewDecFromInt(tc.totalExcess.Amount).Mul(tc.AdditionalStakingRewards).TruncateInt()

			s.bankKeeper.EXPECT().SendCoinsFromModuleToModule(s.ctx, types.ModuleName, authtypes.FeeCollectorName, sdk.NewCoins(sdk.NewCoin(params.MintDenom, stakingRewards))).Return(nil)

			// Execute the method
			err = s.mintKeeper.AdditionalStakingRewards(s.ctx, tc.totalExcess)

			// Check results
			if tc.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedStakingRewards, stakingRewards)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestExcessDistribution() {

	testCases := []struct {
		name        string
		totalExcess sdk.Coin
		expectError bool
	}{
		{
			name:        "successful staking rewards calculation",
			totalExcess: sdk.NewCoin("urock", math.NewInt(100)),
			expectError: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			s.SetupTest()

			params := types.Params{
				MintDenom:                "urock",
				AdditionalBurnRate:       math.LegacyNewDecWithPrec(10, 2), // 10%
				AdditionalMpcRewards:     math.LegacyNewDecWithPrec(20, 2), // 20%
				AdditionalStakingRewards: math.LegacyNewDecWithPrec(15, 2), // 15%
				ProtocolWalletAddress:    "zen1fhln2vnudxddpymqy82vzqhnlsfh4stjd683ze",
			}
			err := s.mintKeeper.Params.Set(s.ctx, params)
			s.Require().NoError(err)

			burnAmount := math.LegacyNewDecFromInt(tc.totalExcess.Amount).Mul(params.AdditionalBurnRate).TruncateInt()

			s.bankKeeper.EXPECT().BurnCoins(s.ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(params.MintDenom, burnAmount))).Return(nil)

			mpcRewards := math.LegacyNewDecFromInt(tc.totalExcess.Amount).Mul(params.AdditionalMpcRewards).TruncateInt()

			s.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(s.ctx, types.ModuleName, sdk.MustAccAddressFromBech32(params.ProtocolWalletAddress), sdk.NewCoins(sdk.NewCoin(params.MintDenom, mpcRewards))).Return(nil)

			stakingRewards := math.LegacyNewDecFromInt(tc.totalExcess.Amount).Mul(params.AdditionalStakingRewards).TruncateInt()

			s.bankKeeper.EXPECT().SendCoinsFromModuleToModule(s.ctx, types.ModuleName, authtypes.FeeCollectorName, sdk.NewCoins(sdk.NewCoin(params.MintDenom, stakingRewards))).Return(nil)

			// Execute the method
			err = s.mintKeeper.ExcessDistribution(s.ctx, tc.totalExcess)

			// Check results
			if tc.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetMintModuleBalance() {

	testCases := []struct {
		name        string
		expectError bool
	}{
		{
			name:        "successful balance calculation",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			s.SetupTest()

			params := types.Params{
				MintDenom: "urock",
			}
			err := s.mintKeeper.Params.Set(s.ctx, params)
			s.Require().NoError(err)

			s.accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(sdk.MustAccAddressFromBech32("zen1m3h30wlvsf8llruxtpukdvsy0km2kum8ju4et3")).AnyTimes()
			s.bankKeeper.EXPECT().GetBalance(s.ctx, sdk.MustAccAddressFromBech32("zen1m3h30wlvsf8llruxtpukdvsy0km2kum8ju4et3"), "urock").Return(sdk.NewCoin("urock", math.NewInt(10))).AnyTimes()

			// Execute the method
			_, err = s.mintKeeper.GetMintModuleBalance(s.ctx)

			// Check results
			if tc.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
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

func (s *IntegrationTestSuite) TestCheckZenBtcSwapThreshold() {
	tests := []struct {
		name                string
		requiredRockBalance uint64
		rockFeePoolBalance  uint64
		errExpected         bool
		expectedErr         error
	}{
		{
			name:                "happy path",
			requiredRockBalance: 1000000,
			rockFeePoolBalance:  1200000,
			errExpected:         false,
			expectedErr:         nil,
		},
		{
			name:                "balance is lower than required",
			requiredRockBalance: 1000000,
			rockFeePoolBalance:  800000,
			errExpected:         false,
			expectedErr:         nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()

			s.zenexKeeper.EXPECT().GetRequiredRockBalance(s.ctx).Return(tt.requiredRockBalance, nil).AnyTimes()
			s.zenexKeeper.EXPECT().GetRockFeePoolBalance(s.ctx).Return(tt.rockFeePoolBalance).AnyTimes()

			if !tt.errExpected {
				s.zenexKeeper.EXPECT().CreateRockBtcSwap(s.ctx, tt.rockFeePoolBalance).Return(nil).AnyTimes()
			}

			err := s.mintKeeper.CheckZenBtcSwapThreshold(s.ctx)
			if tt.errExpected {
				s.Require().Error(err)
				s.Require().Equal(tt.expectedErr.Error(), err.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
