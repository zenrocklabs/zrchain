package keeper_test

import (
	"cosmossdk.io/math"

	appparams "github.com/Zenrock-Foundation/zrchain/v6/app/params"
	identitytestutil "github.com/Zenrock-Foundation/zrchain/v6/x/identity/testutil"
	treasurytestutil "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/testutil"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	zenextestutil "github.com/Zenrock-Foundation/zrchain/v6/x/zenex/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) TestMsgSwapRequest() {
	s.T().Skip("zenex module is currently disabled")

	tests := []struct {
		name           string
		input          *types.MsgSwapRequest
		assetPrice     math.LegacyDec
		assetPriceBtc  math.LegacyDec
		assetPriceRock math.LegacyDec
		expErr         bool
		expErrMsg      string
		want           *types.MsgSwapRequestResponse
		wantSwap       types.Swap
	}{
		{
			name: "Pass: Happy Path: rockbtc",
			input: &types.MsgSwapRequest{
				Creator:   "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				Pair:      types.TradePair_TRADE_PAIR_ROCK_BTC,
				Workspace: types.DefaultParams().ZenexWorkspaceAddress,
				AmountIn:  10000000000,
				RockKeyId: 1,
				BtcKeyId:  4,
			},
			assetPrice:     zenextestutil.SampleRockBtcPrice,
			assetPriceBtc:  zenextestutil.SampleBtcUSDPrice,
			assetPriceRock: zenextestutil.SampleRockUSDPrice,
			expErr:         false,
			want: &types.MsgSwapRequestResponse{
				SwapId: 1,
			},
			wantSwap: types.Swap{
				Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				SwapId:  1,
				Status:  types.SwapStatus_SWAP_STATUS_INITIATED,
				Pair:    types.TradePair_TRADE_PAIR_ROCK_BTC,
				Data: &types.SwapData{
					BaseToken: &validationtypes.AssetData{
						Asset:     validationtypes.Asset_ROCK,
						PriceUSD:  zenextestutil.SampleRockUSDPrice, // 0.025 USD per ROCK
						Precision: 6,
					},
					QuoteToken: &validationtypes.AssetData{
						Asset:     validationtypes.Asset_BTC,
						PriceUSD:  zenextestutil.SampleBtcUSDPrice, // 110,000 USD per BTC
						Precision: 8,
					},
					Price:     zenextestutil.SampleRockBtcPrice, // 0.00000022727 satoshis per urock
					AmountIn:  10000000000,                      // 10,000,000,000 urock
					AmountOut: 2272,                             // 10,000,000,000 * 0.00000022727 = 2272.7 satoshis
				},
				RockKeyId:      1,
				BtcKeyId:       4,
				ZenexPoolKeyId: 16,
				Workspace:      "workspace14a2hpadpsy9h4auve2z8lw",
			},
		},
		{
			name: "FAIL: Invalid pair",
			input: &types.MsgSwapRequest{
				Creator:   "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				Pair:      types.TradePair_TRADE_PAIR_UNSPECIFIED,
				Workspace: types.DefaultParams().ZenexWorkspaceAddress,
				AmountIn:  100000,
				RockKeyId: 1,
				BtcKeyId:  4,
			},
			assetPrice:     zenextestutil.SampleRockBtcPrice,
			assetPriceBtc:  zenextestutil.SampleBtcUSDPrice,
			assetPriceRock: zenextestutil.SampleRockUSDPrice,
			expErr:         true,
			expErrMsg:      "pair is unspecified",
		},
		{
			name: "Fail: Proxy address",
			input: &types.MsgSwapRequest{
				Creator:   types.DefaultParams().BtcProxyAddress,
				Pair:      types.TradePair_TRADE_PAIR_ROCK_BTC,
				Workspace: types.DefaultParams().ZenexWorkspaceAddress,
				AmountIn:  10000000000, // 10,000,000,000 urock
				RockKeyId: 1,
				BtcKeyId:  4,
			},
			assetPrice:     zenextestutil.SampleRockBtcPrice,
			assetPriceBtc:  zenextestutil.SampleBtcUSDPrice,
			assetPriceRock: zenextestutil.SampleRockUSDPrice,
			expErr:         true,
			expErrMsg:      "message creator zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq is not the owner of the workspace",
		},
		{
			name: "FAIL: Invalid creator",
			input: &types.MsgSwapRequest{
				Creator:   "zen10kmgv5gzygnecf46x092ecfe5xcvvv9rdaxmts",
				Pair:      types.TradePair_TRADE_PAIR_ROCK_BTC,
				Workspace: types.DefaultParams().ZenexWorkspaceAddress,
				AmountIn:  100000,
				RockKeyId: 1,
				BtcKeyId:  4,
			},
			assetPrice:     zenextestutil.SampleRockBtcPrice,
			assetPriceBtc:  zenextestutil.SampleBtcUSDPrice,
			assetPriceRock: zenextestutil.SampleRockUSDPrice,
			expErr:         true,
			expErrMsg:      "message creator zen10kmgv5gzygnecf46x092ecfe5xcvvv9rdaxmts is not the owner of the workspace",
		},
		{
			name: "Pass: Happy Path btcrock",
			input: &types.MsgSwapRequest{
				Creator:   "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				Pair:      types.TradePair_TRADE_PAIR_BTC_ROCK,
				Workspace: types.DefaultParams().ZenexWorkspaceAddress,
				AmountIn:  2000,
				RockKeyId: 1,
				BtcKeyId:  4,
			},
			assetPrice:     zenextestutil.SampleBtcRockPrice,
			assetPriceBtc:  zenextestutil.SampleBtcUSDPrice,
			assetPriceRock: zenextestutil.SampleRockUSDPrice,
			expErr:         false,
			want: &types.MsgSwapRequestResponse{
				SwapId: 1,
			},
			wantSwap: types.Swap{
				Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				SwapId:  1,
				Status:  types.SwapStatus_SWAP_STATUS_INITIATED,
				Pair:    types.TradePair_TRADE_PAIR_BTC_ROCK,
				Data: &types.SwapData{
					BaseToken: &validationtypes.AssetData{
						Asset:     validationtypes.Asset_BTC,
						PriceUSD:  zenextestutil.SampleBtcUSDPrice,
						Precision: 8,
					},
					QuoteToken: &validationtypes.AssetData{
						Asset:     validationtypes.Asset_ROCK,
						PriceUSD:  zenextestutil.SampleRockUSDPrice,
						Precision: 6,
					},
					Price:     zenextestutil.SampleBtcRockPrice, // 4,400,000 urock per satoshi
					AmountIn:  2000,                             // 2,000 satoshis
					AmountOut: 8800000000,                       // 2,000 * 4,400,000 = 8.8 billion urock
				},
				RockKeyId:      1,
				BtcKeyId:       4,
				ZenexPoolKeyId: 16,
				Workspace:      "workspace14a2hpadpsy9h4auve2z8lw",
			},
		},
		{
			name: "FAIL: Not enough rock balance in pool",
			input: &types.MsgSwapRequest{
				Creator:   "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				Pair:      types.TradePair_TRADE_PAIR_BTC_ROCK,
				Workspace: types.DefaultParams().ZenexWorkspaceAddress,
				AmountIn:  3000000, // waay too much
				RockKeyId: 1,
				BtcKeyId:  4,
			},
			assetPrice:     zenextestutil.SampleBtcRockPrice,
			assetPriceBtc:  zenextestutil.SampleBtcUSDPrice,
			assetPriceRock: zenextestutil.SampleRockUSDPrice,
			expErr:         true,
			expErrMsg:      "amount 13200000000000 is greater than the available rock balance 10000000000000",
		},
		{
			name: "FAIL: Not enough rock balance in pool",
			input: &types.MsgSwapRequest{
				Creator:   "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				Pair:      types.TradePair_TRADE_PAIR_BTC_ROCK,
				Workspace: types.DefaultParams().ZenexWorkspaceAddress,
				AmountIn:  3000000, // waay too much
				RockKeyId: 1,
				BtcKeyId:  4,
			},
			assetPrice:     zenextestutil.SampleBtcRockPrice,
			assetPriceBtc:  zenextestutil.SampleBtcUSDPrice,
			assetPriceRock: zenextestutil.SampleRockUSDPrice,
			expErr:         true,
			expErrMsg:      "amount 13200000000000 is greater than the available rock balance 10000000000000",
		},
		{
			name: "FAIL: Rock key is not an ECDSA SECP256K1 key",
			input: &types.MsgSwapRequest{
				Creator:   "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				Pair:      types.TradePair_TRADE_PAIR_BTC_ROCK,
				Workspace: types.DefaultParams().ZenexWorkspaceAddress,
				AmountIn:  2000,
				RockKeyId: 3,
				BtcKeyId:  4,
			},
			assetPrice:     zenextestutil.SampleBtcRockPrice,
			assetPriceBtc:  zenextestutil.SampleBtcUSDPrice,
			assetPriceRock: zenextestutil.SampleRockUSDPrice,
			expErr:         true,
			expErrMsg:      "rock key 3 or btc key 4 is not an ECDSA SECP256K1 or BITCOIN SECP256K1 key",
		},
		{
			name: "FAIL: Btc key is not an BITCOIN SECP256K1 key",
			input: &types.MsgSwapRequest{
				Creator:   "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				Pair:      types.TradePair_TRADE_PAIR_BTC_ROCK,
				Workspace: types.DefaultParams().ZenexWorkspaceAddress,
				AmountIn:  2000,
				RockKeyId: 1,
				BtcKeyId:  3,
			},
			assetPrice:     zenextestutil.SampleBtcRockPrice,
			assetPriceBtc:  zenextestutil.SampleBtcUSDPrice,
			assetPriceRock: zenextestutil.SampleRockUSDPrice,
			expErr:         true,
			expErrMsg:      "rock key 1 or btc key 3 is not an ECDSA SECP256K1 or BITCOIN SECP256K1 key",
		},
		{
			name: "FAIL: Price is zero",
			input: &types.MsgSwapRequest{
				Creator:   "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				Pair:      types.TradePair_TRADE_PAIR_BTC_ROCK,
				Workspace: types.DefaultParams().ZenexWorkspaceAddress,
				AmountIn:  2000,
				RockKeyId: 1,
				BtcKeyId:  4,
			},
			assetPriceBtc:  math.LegacyNewDecFromInt(math.NewInt(0)),
			assetPriceRock: math.LegacyNewDecFromInt(math.NewInt(0)),
			assetPrice:     math.LegacyNewDecFromInt(math.NewInt(0)),
			expErr:         true,
			expErrMsg:      "price is zero, check sidecar consensus, got: ROCK=0.000000000000000000, BTC=0.000000000000000000",
		},
		{
			name: "FAIL: key is not in the workspace",
			input: &types.MsgSwapRequest{
				Creator:   "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				Pair:      types.TradePair_TRADE_PAIR_BTC_ROCK,
				Workspace: types.DefaultParams().ZenexWorkspaceAddress,
				AmountIn:  2000,
				RockKeyId: 5,
				BtcKeyId:  4,
			},
			assetPrice:     zenextestutil.SampleBtcRockPrice,
			assetPriceBtc:  zenextestutil.SampleBtcUSDPrice,
			assetPriceRock: zenextestutil.SampleRockUSDPrice,
			expErr:         true,
			expErrMsg:      "rock key 5 or btc key 4 is not in the workspace workspace14a2hpadpsy9h4auve2z8lw",
		},
		{
			name: "FAIL: Invalid workspace",
			input: &types.MsgSwapRequest{
				Creator:   "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				Pair:      types.TradePair_TRADE_PAIR_ROCK_BTC,
				Workspace: "workspace1mphgzyhncnzyggfxmv4nmh",
				AmountIn:  100000,
				RockKeyId: 1,
				BtcKeyId:  4,
			},
			assetPrice:     zenextestutil.SampleRockBtcPrice,
			assetPriceBtc:  zenextestutil.SampleBtcUSDPrice,
			assetPriceRock: zenextestutil.SampleRockUSDPrice,
			expErr:         true,
			expErrMsg:      "workspace1mphgzyhncnzyggfxmv4nmh is not the zenex workspace address",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			err := s.zenexKeeper.SwapsCount.Set(s.ctx, 0)
			s.Require().NoError(err)

			params := types.DefaultParams()
			s.zenexKeeper.SetParams(s.ctx, params)

			s.identityKeeper.EXPECT().GetWorkspace(s.ctx, tt.input.Workspace).Return(&identitytestutil.DefaultWsWithAlice, nil).AnyTimes()
			s.treasuryKeeper.EXPECT().GetKey(s.ctx, tt.input.RockKeyId).Return(&treasurytestutil.DefaultKeys[tt.input.RockKeyId-1], nil).AnyTimes()
			s.treasuryKeeper.EXPECT().GetKey(s.ctx, tt.input.BtcKeyId).Return(&treasurytestutil.DefaultKeys[tt.input.BtcKeyId-1], nil).AnyTimes()

			// Set up mocks needed for both success and error cases
			if tt.input.Pair == types.TradePair_TRADE_PAIR_BTC_ROCK {
				s.validationKeeper.EXPECT().GetBtcRockPrice(s.ctx).Return(tt.assetPrice, nil).AnyTimes()
				s.accountKeeper.EXPECT().GetModuleAddress(types.ZenexCollectorName).Return(sdk.MustAccAddressFromBech32("zen1234wz2aaavp089ttnrj9jwjqraaqxkkadq0k03")).AnyTimes()
				s.bankKeeper.EXPECT().GetBalance(s.ctx, sdk.MustAccAddressFromBech32("zen1234wz2aaavp089ttnrj9jwjqraaqxkkadq0k03"), appparams.BondDenom).Return(sdk.NewCoin(appparams.BondDenom, math.NewIntFromUint64(10000000000000))).AnyTimes()
			}

			// Set up mocks needed for both success and error cases
			if tt.input.Pair == types.TradePair_TRADE_PAIR_BTC_ROCK {
				s.validationKeeper.EXPECT().GetBtcRockPrice(s.ctx).Return(tt.assetPrice, nil).AnyTimes()
				s.accountKeeper.EXPECT().GetModuleAddress(types.ZenexCollectorName).Return(sdk.MustAccAddressFromBech32("zen1234wz2aaavp089ttnrj9jwjqraaqxkkadq0k03")).AnyTimes()
				s.bankKeeper.EXPECT().GetBalance(s.ctx, sdk.MustAccAddressFromBech32("zen1234wz2aaavp089ttnrj9jwjqraaqxkkadq0k03"), appparams.BondDenom).Return(sdk.NewCoin(appparams.BondDenom, math.NewIntFromUint64(10000000000000))).AnyTimes()
			}

			s.validationKeeper.EXPECT().GetAssetPrices(s.ctx).Return(map[validationtypes.Asset]math.LegacyDec{
				validationtypes.Asset_ROCK: tt.assetPriceRock,
				validationtypes.Asset_BTC:  tt.assetPriceBtc,
			}, nil).AnyTimes()

			if !tt.expErr {
				// Set up price expectations based on the pair
				if tt.input.Pair == types.TradePair_TRADE_PAIR_ROCK_BTC {
					s.validationKeeper.EXPECT().GetRockBtcPrice(s.ctx).Return(tt.assetPrice, nil).AnyTimes()
				}
				senderAddress, err := treasurytypes.NativeAddress(&treasurytestutil.DefaultKeys[tt.input.RockKeyId-1], "zen")
				if err != nil {
					s.T().Fatalf("failed to convert sender key to zenrock address: %v", err)
				}
				s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(s.ctx, sdk.MustAccAddressFromBech32(senderAddress), types.ZenexCollectorName, sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, math.NewIntFromUint64(tt.input.AmountIn)))).Return(nil).AnyTimes()
			}

			response, err := s.msgServer.SwapRequest(s.ctx, tt.input)

			if tt.expErr {
				s.Require().Error(err)
				s.Require().Equal(tt.expErrMsg, err.Error())
				s.Require().Nil(response)
				s.Require().Nil(response)
			} else {
				swap, err := s.zenexKeeper.SwapsStore.Get(s.ctx, tt.want.SwapId)
				s.Require().NoError(err)
				s.Require().NotNil(swap)
				s.Require().Equal(tt.want.SwapId, response.SwapId)
				s.Require().Equal(tt.want.SwapId, response.SwapId)
				s.Require().Equal(tt.wantSwap, swap)
			}
		})
	}
}
