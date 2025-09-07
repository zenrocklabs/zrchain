package keeper_test

import (
	"cosmossdk.io/math"

	identitytestutil "github.com/Zenrock-Foundation/zrchain/v6/x/identity/testutil"
	treasurytestutil "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/testutil"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	zenextestutil "github.com/Zenrock-Foundation/zrchain/v6/x/zenex/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
)

func (s *IntegrationTestSuite) TestMsgSwap() {

	tests := []struct {
		name      string
		input     *types.MsgSwap
		expErr    bool
		expErrMsg string
		want      *types.MsgSwapResponse
		wantSwap  *types.Swap
	}{
		{
			name: "Pass: Happy Path",
			input: &types.MsgSwap{
				Creator:      "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				Pair:         "rockbtc",
				Workspace:    "workspace14a2hpadpsy9h4auve2z8lw",
				AmountIn:     100000,
				Yield:        false,
				SenderKey:    1,
				RecipientKey: 2,
			},
			expErr: false,
			want: &types.MsgSwapResponse{
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
				SenderKeyId:    1,
				RecipientKeyId: 2,
				Workspace:      "workspace14a2hpadpsy9h4auve2z8lw",
				ZenbtcYield:    false,
			},
		},
		{
			name: "FAIL: Invalid pair",
			input: &types.MsgSwap{
				Creator:      "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				Pair:         "wrongpair",
				Workspace:    "workspace14a2hpadpsy9h4auve2z8lw",
				AmountIn:     100000,
				Yield:        false,
				SenderKey:    1,
				RecipientKey: 2,
			},
			expErr:    true,
			expErrMsg: "invalid keytype wrongpair, valid types [rockbtc btcrock]: invalid request",
		},
		{
			name: "Pass: Proxy address",
			input: &types.MsgSwap{
				Creator:      "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
				Pair:         "rockbtc",
				Workspace:    "workspace14a2hpadpsy9h4auve2z8lw",
				AmountIn:     100000,
				Yield:        false,
				SenderKey:    1,
				RecipientKey: 2,
			},
			expErr: false,
			want: &types.MsgSwapResponse{
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
				SenderKeyId:    1,
				RecipientKeyId: 2,
				Workspace:      "workspace14a2hpadpsy9h4auve2z8lw",
				ZenbtcYield:    false,
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.zenexKeeper.SwapsCount.Set(s.ctx, 0)
			s.Require().NoError(err)

			params := types.DefaultParams()
			params.BtcProxyAddress = tt.input.Creator
			s.zenexKeeper.SetParams(s.ctx, params)

			if !tt.expErr {
				s.identityKeeper.EXPECT().GetWorkspace(s.ctx, tt.input.Workspace).Return(&identitytestutil.DefaultWsWithAlice, nil)
				s.treasuryKeeper.EXPECT().GetKey(s.ctx, tt.input.SenderKey).Return(&treasurytestutil.DefaultKeys[0], nil)
				s.treasuryKeeper.EXPECT().GetKey(s.ctx, tt.input.RecipientKey).Return(&treasurytestutil.DefaultKeys[1], nil)
				s.validationKeeper.EXPECT().GetAssetPrices(s.ctx).Return(map[validationtypes.Asset]math.LegacyDec{
					validationtypes.Asset_ROCK: zenextestutil.SampleRockBtcPrice,
					validationtypes.Asset_BTC:  zenextestutil.SampleBtcRockPrice,
				}, nil)
				s.validationKeeper.EXPECT().GetRockBtcPrice(s.ctx).Return(zenextestutil.SampleRockBtcPrice, nil)
			}

			swapId, err := s.msgServer.Swap(s.ctx, tt.input)

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
