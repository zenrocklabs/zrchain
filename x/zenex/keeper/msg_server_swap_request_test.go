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

	tests := []struct {
		name      string
		input     *types.MsgSwapRequest
		expErr    bool
		expErrMsg string
		want      *types.MsgSwapRequestResponse
		wantSwap  *types.Swap
	}{
		{
			name: "Pass: Happy Path: rockbtc",
			input: &types.MsgSwapRequest{
				Creator:   "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				Pair:      "rockbtc",
				Workspace: "workspace14a2hpadpsy9h4auve2z8lw",
				AmountIn:  100000,
				RockKeyId: 1,
				BtcKeyId:  2,
			},
			expErr: false,
			want: &types.MsgSwapRequestResponse{
				SwapId: 1,
			},
			wantSwap: &types.Swap{
				SwapId: 1,
				Status: types.SwapStatus_SWAP_STATUS_REQUESTED,
				Pair:   "rockbtc",
				Data: &types.SwapData{
					BaseToken: &validationtypes.AssetData{
						Asset:     validationtypes.Asset_ROCK,
						PriceUSD:  zenextestutil.SampleRockBtcPrice,
						Precision: 6,
					},
					QuoteToken: &validationtypes.AssetData{
						Asset:     validationtypes.Asset_BTC,
						PriceUSD:  zenextestutil.SampleBtcRockPrice,
						Precision: 8,
					},
					Price:     math.LegacyNewDec(100000),
					AmountIn:  100000,
					AmountOut: 100000,
				},
				RockKeyId:      1,
				BtcKeyId:       2,
				ZenexPoolKeyId: 3,
				Workspace:      "workspace14a2hpadpsy9h4auve2z8lw",
			},
		},
		{
			name: "FAIL: Invalid pair",
			input: &types.MsgSwapRequest{
				Creator:   "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				Pair:      "wrongpair",
				Workspace: "workspace14a2hpadpsy9h4auve2z8lw",
				AmountIn:  100000,
				RockKeyId: 1,
				BtcKeyId:  2,
			},
			expErr:    true,
			expErrMsg: "invalid keytype wrongpair, valid types [rockbtc btcrock]: invalid request",
		},
		{
			name: "Pass: Proxy address",
			input: &types.MsgSwapRequest{
				Creator:   "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
				Pair:      "rockbtc",
				Workspace: "workspace14a2hpadpsy9h4auve2z8lw",
				AmountIn:  100000,
				RockKeyId: 1,
				BtcKeyId:  2,
			},
			expErr: false,
			want: &types.MsgSwapRequestResponse{
				SwapId: 1,
			},
			wantSwap: &types.Swap{
				SwapId: 1,
				Status: types.SwapStatus_SWAP_STATUS_REQUESTED,
				Pair:   "rockbtc",
				Data: &types.SwapData{
					BaseToken: &validationtypes.AssetData{
						Asset:     validationtypes.Asset_ROCK,
						PriceUSD:  zenextestutil.SampleRockBtcPrice,
						Precision: 6,
					},
					QuoteToken: &validationtypes.AssetData{
						Asset:     validationtypes.Asset_BTC,
						PriceUSD:  zenextestutil.SampleBtcRockPrice,
						Precision: 8,
					},
					Price:     math.LegacyNewDec(100000),
					AmountIn:  100000,
					AmountOut: 100000,
				},
				RockKeyId:      1,
				BtcKeyId:       2,
				ZenexPoolKeyId: 3,
				Workspace:      "workspace14a2hpadpsy9h4auve2z8lw",
			},
		},
		{
			name: "FAIL: Invalid creator",
			input: &types.MsgSwapRequest{
				Creator:   "zen10kmgv5gzygnecf46x092ecfe5xcvvv9rdaxmts",
				Pair:      "rockbtc",
				Workspace: "workspace14a2hpadpsy9h4auve2z8lw",
				AmountIn:  100000,
				RockKeyId: 1,
				BtcKeyId:  2,
			},
			expErr:    true,
			expErrMsg: "message creator is not the owner of the workspace or the btc proxy address",
		},
		{
			name: "Pass: Happy Path btcrock",
			input: &types.MsgSwapRequest{
				Creator:   "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				Pair:      "btcrock",
				Workspace: "workspace14a2hpadpsy9h4auve2z8lw",
				AmountIn:  100000,
				RockKeyId: 1,
				BtcKeyId:  2,
			},
			expErr: false,
			want: &types.MsgSwapRequestResponse{
				SwapId: 1,
			},
			wantSwap: &types.Swap{
				SwapId: 1,
				Status: types.SwapStatus_SWAP_STATUS_REQUESTED,
				Pair:   "btcrock",
				Data: &types.SwapData{
					BaseToken: &validationtypes.AssetData{
						Asset:     validationtypes.Asset_BTC,
						PriceUSD:  zenextestutil.SampleBtcRockPrice,
						Precision: 8,
					},
					QuoteToken: &validationtypes.AssetData{
						Asset:     validationtypes.Asset_ROCK,
						PriceUSD:  zenextestutil.SampleRockBtcPrice,
						Precision: 6,
					},
					Price:     math.LegacyNewDec(100000),
					AmountIn:  100000,
					AmountOut: 100000,
				},
				RockKeyId:      1,
				BtcKeyId:       2,
				ZenexPoolKeyId: 3,
				Workspace:      "workspace14a2hpadpsy9h4auve2z8lw",
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.zenexKeeper.SwapsCount.Set(s.ctx, 0)
			s.Require().NoError(err)

			params := types.DefaultParams()
			s.zenexKeeper.SetParams(s.ctx, params)

			s.identityKeeper.EXPECT().GetWorkspace(s.ctx, tt.input.Workspace).Return(&identitytestutil.DefaultWsWithAlice, nil).AnyTimes()
			s.treasuryKeeper.EXPECT().GetKey(s.ctx, tt.input.RockKeyId).Return(&treasurytestutil.DefaultKeys[tt.input.RockKeyId-1], nil).AnyTimes()
			s.treasuryKeeper.EXPECT().GetKey(s.ctx, tt.input.BtcKeyId).Return(&treasurytestutil.DefaultKeys[tt.input.BtcKeyId-1], nil).AnyTimes()

			if !tt.expErr {
				s.validationKeeper.EXPECT().GetAssetPrices(s.ctx).Return(map[validationtypes.Asset]math.LegacyDec{
					validationtypes.Asset_ROCK: zenextestutil.SampleRockBtcPrice,
					validationtypes.Asset_BTC:  zenextestutil.SampleBtcRockPrice,
				}, nil).AnyTimes()
				// Set up price expectations based on the pair
				if tt.input.Pair == "rockbtc" {
					s.validationKeeper.EXPECT().GetRockBtcPrice(s.ctx).Return(zenextestutil.SampleRockBtcPrice, nil).AnyTimes()
				} else if tt.input.Pair == "btcrock" {
					s.validationKeeper.EXPECT().GetBtcRockPrice(s.ctx).Return(zenextestutil.SampleBtcRockPrice, nil).AnyTimes()
				}
				senderAddress, err := treasurytypes.NativeAddress(&treasurytestutil.DefaultKeys[tt.input.RockKeyId-1], "zen")
				if err != nil {
					s.T().Fatalf("failed to convert sender key to zenrock address: %v", err)
				}
				s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(s.ctx, sdk.MustAccAddressFromBech32(senderAddress), types.ZenexCollectorName, sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, math.NewIntFromUint64(tt.input.AmountIn)))).Return(nil).AnyTimes()
			}

			swapId, err := s.msgServer.SwapRequest(s.ctx, tt.input)

			if tt.expErr {
				s.Require().Error(err)
				s.Require().Equal(tt.expErrMsg, err.Error())
				s.Require().Nil(swapId)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(swapId)
				s.Require().Equal(tt.want.SwapId, swapId.SwapId)
			}
		})
	}
}
