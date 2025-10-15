package keeper_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	appparams "github.com/Zenrock-Foundation/zrchain/v6/app/params"
	"github.com/Zenrock-Foundation/zrchain/v6/x/mint"
	treasurytestutil "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/testutil"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
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
	zenbtcKeeper     *zenextestutil.MockZenbtcKeeper
	bankKeeper       *zenextestutil.MockBankKeeper
	accountKeeper    *zenextestutil.MockAccountKeeper
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
	zenbtcKeeper := zenextestutil.NewMockZenbtcKeeper(ctrl)
	bankKeeper := zenextestutil.NewMockBankKeeper(ctrl)
	accountKeeper := zenextestutil.NewMockAccountKeeper(ctrl)

	// Assign the mock keepers to the suite fields
	s.treasuryKeeper = treasuryKeeper
	s.identityKeeper = identityKeeper
	s.validationKeeper = validationKeeper
	s.bankKeeper = bankKeeper
	s.zenbtcKeeper = zenbtcKeeper
	s.accountKeeper = accountKeeper
	s.zenexKeeper = keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		testCtx.Ctx.Logger(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		identityKeeper,
		treasuryKeeper,
		validationKeeper,
		zenbtcKeeper,
		bankKeeper,
		accountKeeper,
	)

	s.Require().Equal(testCtx.Ctx.Logger().With("module", "x/"+types.ModuleName),
		s.zenexKeeper.Logger())

	s.msgServer = keeper.NewMsgServerImpl(s.zenexKeeper)
}

func (s *IntegrationTestSuite) TestGetPrice() {

	tests := []struct {
		name        string
		pair        types.TradePair
		rockPrice   math.LegacyDec
		btcPrice    math.LegacyDec
		expected    math.LegacyDec
		expectedErr error
	}{
		{
			name:      "rockbtc",
			pair:      types.TradePair_TRADE_PAIR_ROCK_BTC,
			rockPrice: zenextestutil.SampleRockUSDPrice, // $0.025 USD
			btcPrice:  zenextestutil.SampleBtcUSDPrice,  // $110,000 USD
			expected:  zenextestutil.SampleRockBtcPrice, // 4,400,000 ROCK per BTC
		},
		{
			name:      "btcrock",
			pair:      types.TradePair_TRADE_PAIR_BTC_ROCK,
			rockPrice: zenextestutil.SampleRockUSDPrice, // $0.025 USD
			btcPrice:  zenextestutil.SampleBtcUSDPrice,  // $110,000 USD
			expected:  zenextestutil.SampleBtcRockPrice, // 0.000000227 BTC per ROCK
		},
		{
			name:        "unknown pair",
			pair:        types.TradePair_TRADE_PAIR_UNSPECIFIED,
			expected:    math.LegacyDec{},
			expectedErr: fmt.Errorf("unknown pair: TRADE_PAIR_UNSPECIFIED"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Set up expectations based on the pair
			if tt.pair == types.TradePair_TRADE_PAIR_ROCK_BTC {
				s.validationKeeper.EXPECT().GetRockBtcPrice(s.ctx).Return(tt.expected, nil)
			} else if tt.pair == types.TradePair_TRADE_PAIR_BTC_ROCK {
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
		pair        types.TradePair
		amountIn    uint64
		price       math.LegacyDec
		expected    uint64
		expectedErr error
	}{
		{
			name:     "happy path rockbtc",
			pair:     types.TradePair_TRADE_PAIR_ROCK_BTC,
			amountIn: 10000000000,                      // 10 billion urock (10 ROCK)
			price:    zenextestutil.SampleRockBtcPrice, // 0.00000022727 satoshis per urock
			expected: 2272,                             // 10,000,000,000 * 0.00000022727 = 2272.7 satoshis (truncated to 2272)
		},
		{
			name:        "rockbtc",
			pair:        types.TradePair_TRADE_PAIR_ROCK_BTC,
			amountIn:    100000000,                                                                   // 1 million urock (1 ROCK)
			price:       zenextestutil.SampleRockBtcPrice,                                            // 0.00000022727 satoshis per urock
			expected:    22,                                                                          // 1,000,000 * 0.00000022727 = 0.22727 satoshis (truncated to 0)
			expectedErr: fmt.Errorf("calculated satoshis 22 is less than the minimum satoshis 1000"), // 0.22727 < 1000 minimum
		},
		{
			name:        "rockbtc with less than minimum satoshis",
			pair:        types.TradePair_TRADE_PAIR_ROCK_BTC,
			amountIn:    10000,                                                                      // 10,000 urock (0.01 ROCK)
			price:       zenextestutil.SampleRockBtcPrice,                                           // 0.00000022727 satoshis per urock
			expectedErr: fmt.Errorf("calculated satoshis 0 is less than the minimum satoshis 1000"), // 10000 * 0.00000022727 = 0.0022727 satoshis (truncated to 0) < 1000 minimum
		},
		{
			name:     "btcrock",
			pair:     types.TradePair_TRADE_PAIR_BTC_ROCK,
			amountIn: 10000,                            // 10,000 satoshis
			price:    zenextestutil.SampleBtcRockPrice, // 4,400,000 urock per satoshi
			expected: 44000000000,                      // 10,000 * 4,400,000 = 44 billion urock
		},
		{
			name:        "btcrock with 999 satoshis",
			pair:        types.TradePair_TRADE_PAIR_BTC_ROCK,
			amountIn:    999, // 999 satoshis (above 1000 minimum)
			price:       zenextestutil.SampleBtcRockPrice,
			expectedErr: fmt.Errorf("999 satoshis in is less than the minimum satoshis 1000"),
		},
		{
			name:        "unknown pair",
			pair:        types.TradePair_TRADE_PAIR_UNSPECIFIED,
			expected:    0,
			expectedErr: fmt.Errorf("unknown pair: TRADE_PAIR_UNSPECIFIED"),
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
		pair          types.TradePair
		rockPrice     math.LegacyDec
		btcPrice      math.LegacyDec
		expected      types.SwapPair
		expectedPrice math.LegacyDec
		expectedErr   error
	}{
		{
			name:      "rockbtc",
			pair:      types.TradePair_TRADE_PAIR_ROCK_BTC,
			rockPrice: zenextestutil.SampleRockUSDPrice,
			btcPrice:  zenextestutil.SampleBtcUSDPrice,
			expected: types.SwapPair{
				BaseToken: &validationtypes.AssetData{
					Asset:     validationtypes.Asset_ROCK,
					Precision: 6,
					PriceUSD:  zenextestutil.SampleRockUSDPrice,
				},
				QuoteToken: &validationtypes.AssetData{
					Asset:     validationtypes.Asset_BTC,
					Precision: 8,
					PriceUSD:  zenextestutil.SampleBtcUSDPrice,
				},
			},
			expectedPrice: zenextestutil.SampleRockBtcPrice,
			expectedErr:   nil,
		},
		{
			name:      "btcrock",
			pair:      types.TradePair_TRADE_PAIR_BTC_ROCK,
			rockPrice: zenextestutil.SampleRockUSDPrice,
			btcPrice:  zenextestutil.SampleBtcUSDPrice,
			expected: types.SwapPair{
				BaseToken: &validationtypes.AssetData{
					Asset:     validationtypes.Asset_BTC,
					Precision: 8,
					PriceUSD:  zenextestutil.SampleBtcUSDPrice,
				},
				QuoteToken: &validationtypes.AssetData{
					Asset:     validationtypes.Asset_ROCK,
					Precision: 6,
					PriceUSD:  zenextestutil.SampleRockUSDPrice,
				},
			},
			expectedPrice: zenextestutil.SampleBtcRockPrice,
			expectedErr:   nil,
		},
		{
			name:          "unknown pair",
			pair:          types.TradePair_TRADE_PAIR_UNSPECIFIED,
			rockPrice:     zenextestutil.SampleRockUSDPrice,
			btcPrice:      zenextestutil.SampleBtcUSDPrice,
			expected:      types.SwapPair{},
			expectedPrice: math.LegacyDec{},
			expectedErr:   fmt.Errorf("unknown pair: TRADE_PAIR_UNSPECIFIED"),
		},
		{
			name:      "price is zero",
			pair:      types.TradePair_TRADE_PAIR_ROCK_BTC,
			rockPrice: math.LegacyNewDecFromInt(math.NewInt(0)),
			btcPrice:  math.LegacyNewDecFromInt(math.NewInt(0)),
			expected: types.SwapPair{
				BaseToken: &validationtypes.AssetData{
					Asset:     validationtypes.Asset_ROCK,
					Precision: 6,
					PriceUSD:  math.LegacyNewDecFromInt(math.NewInt(0)),
				},
				QuoteToken: &validationtypes.AssetData{
					Asset:     validationtypes.Asset_BTC,
					Precision: 8,
					PriceUSD:  math.LegacyNewDecFromInt(math.NewInt(0)),
				},
			},
			expectedPrice: math.LegacyDec{},
			expectedErr:   fmt.Errorf("price is zero, check sidecar consensus, got: ROCK=0.000000000000000000, BTC=0.000000000000000000"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// GetAssetPrices is always called first
			s.validationKeeper.EXPECT().GetAssetPrices(s.ctx).Return(map[validationtypes.Asset]math.LegacyDec{
				validationtypes.Asset_ROCK: tt.rockPrice,
				validationtypes.Asset_BTC:  tt.btcPrice,
			}, nil)

			if tt.expectedErr == nil {
				if tt.pair == types.TradePair_TRADE_PAIR_ROCK_BTC {
					s.validationKeeper.EXPECT().GetRockBtcPrice(s.ctx).Return(zenextestutil.SampleRockBtcPrice, nil)
				} else if tt.pair == types.TradePair_TRADE_PAIR_BTC_ROCK {
					s.validationKeeper.EXPECT().GetBtcRockPrice(s.ctx).Return(zenextestutil.SampleBtcRockPrice, nil)
				}
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

func (s *IntegrationTestSuite) TestEscrowRock() {
	tests := []struct {
		name        string
		senderKey   treasurytypes.Key
		amountIn    uint64
		errExpected bool
		expectedErr error
	}{
		{
			name:        "happy path",
			senderKey:   treasurytestutil.DefaultKeys[0],
			amountIn:    100000,
			errExpected: false,
			expectedErr: nil,
		},
		{
			name:        "fail: wrong key type",
			senderKey:   treasurytestutil.DefaultKeys[2],
			amountIn:    100000,
			errExpected: true,
			expectedErr: types.ErrWrongKeyType,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if !tt.errExpected {
				s.treasuryKeeper.EXPECT().GetKey(s.ctx, tt.senderKey.Id).Return(&tt.senderKey, nil).AnyTimes()
				senderAddress, err := treasurytypes.NativeAddress(&tt.senderKey, "zen")
				if err != nil {
					s.T().Fatalf("failed to convert sender key to zenrock address: %v", err)
				}
				s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(s.ctx, sdk.MustAccAddressFromBech32(senderAddress), types.ZenexCollectorName, sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, math.NewIntFromUint64(tt.amountIn)))).Return(nil)
			}

			err := s.zenexKeeper.EscrowRock(s.ctx, tt.senderKey, tt.amountIn)

			if tt.expectedErr != nil {
				s.Require().Error(err)
				s.Require().Equal(tt.expectedErr.Error(), err.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestCheckRedeemableAsset() {
	tests := []struct {
		name             string
		amountOut        uint64
		zenexRockBalance uint64
		errExpected      bool
		expectedErr      error
	}{
		{
			name:             "happy path",
			amountOut:        2000,
			zenexRockBalance: 1000000000000,
			errExpected:      false,
			expectedErr:      nil,
		},
		{
			name:             "fail: not enough rock balance in pool",
			amountOut:        440000000000,
			zenexRockBalance: 100000,
			errExpected:      true,
			expectedErr:      fmt.Errorf("amount 440000000000 is greater than the available rock balance 100000"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {

			s.accountKeeper.EXPECT().GetModuleAddress(types.ZenexCollectorName).Return(sdk.MustAccAddressFromBech32("zen1234wz2aaavp089ttnrj9jwjqraaqxkkadq0k03")).AnyTimes()
			s.bankKeeper.EXPECT().GetBalance(s.ctx, sdk.MustAccAddressFromBech32("zen1234wz2aaavp089ttnrj9jwjqraaqxkkadq0k03"), appparams.BondDenom).Return(sdk.NewCoin(appparams.BondDenom, math.NewIntFromUint64(tt.zenexRockBalance))).Times(1)

			err := s.zenexKeeper.CheckRedeemableAsset(s.ctx, tt.amountOut, zenextestutil.SampleBtcRockPrice)

			if tt.expectedErr != nil {
				s.Require().Error(err)
				s.Require().Equal(tt.expectedErr.Error(), err.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetSwapThreshold() {
	tests := []struct {
		name          string
		swapThreshold uint64
		rockPrice     math.LegacyDec
		btcPrice      math.LegacyDec
		expected      uint64
		expectedErr   error
	}{
		{
			name:          "happy path",
			swapThreshold: 100000,
			rockPrice:     zenextestutil.SampleRockBtcPrice,
			btcPrice:      zenextestutil.SampleBtcRockPrice,
			expected:      440005280063,
			expectedErr:   nil,
		},
		{
			name:          "happy path with default swap threshold",
			swapThreshold: 6100,
			rockPrice:     zenextestutil.SampleRockBtcPrice,
			btcPrice:      zenextestutil.SampleBtcRockPrice,
			expected:      26840322083,
			expectedErr:   nil,
		},
		{
			name:          "fail: rock price is zero",
			swapThreshold: 100000,
			rockPrice:     math.LegacyNewDecFromInt(math.NewInt(0)),
			btcPrice:      zenextestutil.SampleBtcRockPrice,
			expected:      0,
			expectedErr:   fmt.Errorf("rock to btc price must be positive, got: 0.000000000000000000"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()

			params := types.DefaultParams()
			params.SwapThresholdSatoshis = tt.swapThreshold
			s.zenexKeeper.SetParams(s.ctx, params)

			s.validationKeeper.EXPECT().GetRockBtcPrice(s.ctx).Return(tt.rockPrice, nil).AnyTimes()

			swapThreshold, err := s.zenexKeeper.GetRequiredRockBalance(s.ctx)

			if tt.expectedErr != nil {
				s.Require().Error(err)
				s.Require().Equal(tt.expectedErr.Error(), err.Error())
				s.Require().Equal(tt.expected, swapThreshold)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tt.expected, swapThreshold)
			}
		})
	}
}
