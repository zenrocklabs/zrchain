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
	bankKeeper := zenextestutil.NewMockBankKeeper(ctrl)
	accountKeeper := zenextestutil.NewMockAccountKeeper(ctrl)

	// Assign the mock keepers to the suite fields
	s.treasuryKeeper = treasuryKeeper
	s.identityKeeper = identityKeeper
	s.validationKeeper = validationKeeper
	s.bankKeeper = bankKeeper
	s.accountKeeper = accountKeeper
	s.zenexKeeper = keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		testCtx.Ctx.Logger(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		identityKeeper,
		treasuryKeeper,
		validationKeeper,
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
		pair        string
		rockPrice   math.LegacyDec
		btcPrice    math.LegacyDec
		expected    math.LegacyDec
		expectedErr error
	}{
		{
			name:      "rockbtc",
			pair:      "rockbtc",
			rockPrice: zenextestutil.SampleRockUSDPrice, // $0.025 USD
			btcPrice:  zenextestutil.SampleBtcUSDPrice,  // $110,000 USD
			expected:  zenextestutil.SampleRockBtcPrice, // 4,400,000 ROCK per BTC
		},
		{
			name:      "btcrock",
			pair:      "btcrock",
			rockPrice: zenextestutil.SampleRockUSDPrice, // $0.025 USD
			btcPrice:  zenextestutil.SampleBtcUSDPrice,  // $110,000 USD
			expected:  zenextestutil.SampleBtcRockPrice, // 0.000000227 BTC per ROCK
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
			amountIn: 10000000000,                      // 10 billion urock (10 ROCK)
			price:    zenextestutil.SampleRockBtcPrice, // 0.00000022727 satoshis per urock
			expected: 2272,                             // 10,000,000,000 * 0.00000022727 = 2272.7 satoshis (truncated to 2272)
		},
		{
			name:        "rockbtc",
			pair:        "rockbtc",
			amountIn:    1000000,                                                                // 1 million urock (1 ROCK)
			price:       zenextestutil.SampleRockBtcPrice,                                       // 0.00000022727 satoshis per urock
			expected:    0,                                                                      // 1,000,000 * 0.00000022727 = 0.22727 satoshis (truncated to 0)
			expectedErr: fmt.Errorf("amount 1000000 in is less than the minimum satoshis 1000"), // 0.22727 < 1000 minimum
		},
		{
			name:        "rockbtc with less than minimum satoshis",
			pair:        "rockbtc",
			amountIn:    1,                                                                // 1 urock
			price:       math.LegacyNewDecFromInt(math.NewInt(1)),                         // 1 satoshi per urock (very small price)
			expectedErr: fmt.Errorf("amount 1 in is less than the minimum satoshis 1000"), // 1 * 1 = 1 satoshi < 1000 minimum
		},
		{
			name:     "btcrock",
			pair:     "btcrock",
			amountIn: 10000,                            // 10,000 satoshis
			price:    zenextestutil.SampleBtcRockPrice, // 4,400,000 urock per satoshi
			expected: 44000000000,                      // 10,000 * 4,400,000 = 44 billion urock
		},
		{
			name:        "btcrock with 999 satoshis",
			pair:        "btcrock",
			amountIn:    999, // 999 satoshis (above 1000 minimum)
			price:       zenextestutil.SampleBtcRockPrice,
			expectedErr: fmt.Errorf("amount 999 in is less than the minimum satoshis 1000"),
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
