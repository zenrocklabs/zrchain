package keeper_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/mint"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/keeper"
	zenextestutil "github.com/Zenrock-Foundation/zrchain/v6/x/zenex/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	ubermock "go.uber.org/mock/gomock"
)

func TestKeeper(t *testing.T) {
	require.True(t, true)
}

type IntegrationTestSuite struct {
	suite.Suite

	zenexKeeper      keeper.Keeper
	ctx              sdk.Context
	msgServer        types.MsgServer
	identityKeeper   *zenextestutil.MockIdentityKeeper
	treasuryKeeper   *zenextestutil.MockTreasuryKeeper
	validationKeeper *zenextestutil.MockValidationKeeper
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
	treasuryKeeper := zenextestutil.NewMockTreasuryKeeper(ctrl)
	identityKeeper := zenextestutil.NewMockIdentityKeeper(ctrl)
	validationKeeper := zenextestutil.NewMockValidationKeeper(ctrl)

	// Assign the mock keepers to the suite fields
	s.treasuryKeeper = treasuryKeeper
	s.identityKeeper = identityKeeper
	s.validationKeeper = validationKeeper
	s.zenexKeeper = keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		testCtx.Ctx.Logger(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		identityKeeper,
		treasuryKeeper,
		validationKeeper,
	)

	s.Require().Equal(testCtx.Ctx.Logger().With("module", "x/"+types.ModuleName),
		s.zenexKeeper.Logger())

	s.msgServer = keeper.NewMsgServerImpl(s.zenexKeeper)
}

func (s *IntegrationTestSuite) TestGetPrice() {

	tests := []struct {
		name        string
		pair        string
		rockPrice   math.LegacyDec
		btcPrice    math.LegacyDec
		expected    math.LegacyDec
		expectedErr error
	}{
		{
			name:      "rockbtc",
			pair:      "rockbtc",
			rockPrice: math.LegacyNewDecFromInt(math.NewInt(25)).Quo(math.LegacyNewDecFromInt(math.NewInt(1000))), // $0.025 USD
			btcPrice:  math.LegacyNewDecFromInt(math.NewInt(110000)),                                              // $110,000 USD
			expected:  zenextestutil.SampleRockBtcPrice,                                                           // 4,400,000 ROCK per BTC
		},
		{
			name:      "btcrock",
			pair:      "btcrock",
			rockPrice: math.LegacyNewDecFromInt(math.NewInt(25)).Quo(math.LegacyNewDecFromInt(math.NewInt(1000))), // $0.025 USD
			btcPrice:  math.LegacyNewDecFromInt(math.NewInt(110000)),                                              // $110,000 USD
			expected:  zenextestutil.SampleBtcRockPrice,                                                           // 0.000000227 BTC per ROCK
		},
		{
			name:        "unknown pair",
			pair:        "unknown",
			expected:    math.LegacyDec{},
			expectedErr: fmt.Errorf("unknown pair: unknown"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Set up expectations based on the pair
			if tt.pair == "rockbtc" {
				s.validationKeeper.EXPECT().GetRockBtcPrice(s.ctx).Return(tt.expected, nil)
			} else if tt.pair == "btcrock" {
				s.validationKeeper.EXPECT().GetBtcRockPrice(s.ctx).Return(tt.expected, nil)
			}

			price, err := s.zenexKeeper.GetPrice(s.ctx, tt.pair)
			if tt.expectedErr != nil {
				s.Require().Error(err)
				s.Require().Equal(tt.expectedErr.Error(), err.Error())
				s.Require().Equal(math.LegacyDec{}, price)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tt.expected, price)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetAmountOut() {

	tests := []struct {
		name        string
		pair        string
		amountIn    uint64
		price       math.LegacyDec
		expected    uint64
		expectedErr error
	}{
		{
			name:     "rockbtc",
			pair:     "rockbtc",
			amountIn: 1000000,                          // 1 million urock (1 ROCK)
			price:    zenextestutil.SampleRockBtcPrice, // 4,400,000 satoshis per urock
			expected: 4400000000000,                    // 1,000,000 * 4,400,000 = 4.4 trillion satoshis
		},
		{
			name:        "rockbtc with less than minimum satoshis",
			pair:        "rockbtc",
			amountIn:    1,                                        // 1 urock
			price:       math.LegacyNewDecFromInt(math.NewInt(1)), // 1 satoshi per urock (very small price)
			expectedErr: types.ErrMinimumSatoshis,                 // 1 * 1 = 1 satoshi < 5000 minimum
		},
		{
			name:     "btcrock",
			pair:     "btcrock",
			amountIn: 10000,                            // 10,000 satoshis
			price:    zenextestutil.SampleBtcRockPrice, // 2,272,727 urock per satoshi
			expected: 22727270000,                      // 10,000 * 2,272,727 = 22.7 billion urock
		},
		{
			name:        "btcrock with less than minimum satoshis",
			pair:        "btcrock",
			amountIn:    4999, // 4,999 satoshis (less than 5000 minimum)
			price:       zenextestutil.SampleBtcRockPrice,
			expectedErr: types.ErrMinimumSatoshis,
		},
		{
			name:        "unknown pair",
			pair:        "unknown",
			expected:    0,
			expectedErr: fmt.Errorf("unknown pair: unknown"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.zenexKeeper.SetParams(s.ctx, types.DefaultParams())
			price, err := s.zenexKeeper.GetAmountOut(s.ctx, tt.pair, tt.amountIn, tt.price)

			if tt.expectedErr != nil {
				s.Require().Error(err)
				s.Require().Equal(tt.expectedErr.Error(), err.Error())
				s.Require().Equal(uint64(0), price)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tt.expected, price)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetSwaps() {
	swaps, err := s.zenexKeeper.GetSwaps(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(0, len(swaps))
}

func (s *IntegrationTestSuite) TestGetPair() {

	tests := []struct {
		name          string
		pair          string
		expected      types.SwapPair
		expectedPrice math.LegacyDec
		expectedErr   error
	}{
		{
			name: "rockbtc",
			pair: "rockbtc",
			expected: types.SwapPair{
				BaseToken: &validationtypes.AssetData{
					Asset:     validationtypes.Asset_ROCK,
					Precision: 6,
					PriceUSD:  zenextestutil.SampleRockBtcPrice,
				},
				QuoteToken: &validationtypes.AssetData{
					Asset:     validationtypes.Asset_BTC,
					Precision: 8,
					PriceUSD:  zenextestutil.SampleBtcRockPrice,
				},
			},
			expectedPrice: zenextestutil.SampleRockBtcPrice,
			expectedErr:   nil,
		},
		{
			name: "btcrock",
			pair: "btcrock",
			expected: types.SwapPair{
				BaseToken: &validationtypes.AssetData{
					Asset:     validationtypes.Asset_BTC,
					Precision: 8,
					PriceUSD:  zenextestutil.SampleBtcRockPrice,
				},
				QuoteToken: &validationtypes.AssetData{
					Asset:     validationtypes.Asset_ROCK,
					Precision: 6,
					PriceUSD:  zenextestutil.SampleRockBtcPrice,
				},
			},
			expectedPrice: zenextestutil.SampleBtcRockPrice,
			expectedErr:   nil,
		},
		{
			name:          "unknown pair",
			pair:          "unknown",
			expected:      types.SwapPair{},
			expectedPrice: math.LegacyDec{},
			expectedErr:   fmt.Errorf("unknown key type: unknown"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// GetAssetPrices is always called first
			s.validationKeeper.EXPECT().GetAssetPrices(s.ctx).Return(map[validationtypes.Asset]math.LegacyDec{
				validationtypes.Asset_ROCK: zenextestutil.SampleRockBtcPrice,
				validationtypes.Asset_BTC:  zenextestutil.SampleBtcRockPrice,
			}, nil)

			// Set up expectations based on the pair
			if tt.pair == "rockbtc" {
				s.validationKeeper.EXPECT().GetRockBtcPrice(s.ctx).Return(zenextestutil.SampleRockBtcPrice, nil)
			} else if tt.pair == "btcrock" {
				s.validationKeeper.EXPECT().GetBtcRockPrice(s.ctx).Return(zenextestutil.SampleBtcRockPrice, nil)
			}

			pair, price, err := s.zenexKeeper.GetPair(s.ctx, tt.pair)

			if tt.expectedErr != nil {
				s.Require().Error(err)
				s.Require().Equal(tt.expectedErr.Error(), err.Error())
				s.Require().Nil(pair)
				s.Require().Equal(tt.expectedPrice, price)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(pair)
				s.Require().Equal(tt.expected, *pair)
				s.Require().Equal(tt.expectedPrice, price)
			}
		})
	}
}
