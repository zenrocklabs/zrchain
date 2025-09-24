package keeper_test

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	treasurytestutil "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/testutil"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	zenextestutil "github.com/Zenrock-Foundation/zrchain/v6/x/zenex/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	"go.uber.org/mock/gomock"
)

func (s *IntegrationTestSuite) TestMsgZenexTransferRequest() {

	tests := []struct {
		name      string
		input     *types.MsgZenexTransferRequest
		expErr    bool
		expErrMsg string
		want      *types.MsgZenexTransferRequestResponse
		wantSwap  types.Swap
	}{
		{
			name: "Pass: Happy Path rockbtc",
			input: &types.MsgZenexTransferRequest{
				Creator:        types.DefaultParams().BtcProxyAddress,
				SwapId:         zenextestutil.SampleSwap[0].SwapId,
				UnsignedPlusTx: zenextestutil.Btc_tx[0].UnsignedPlusTx,
				CacheId:        zenextestutil.Btc_tx[0].CacheId,
				DataForSigning: zenextestutil.Btc_tx[0].DataForSigning,
				WalletType:     treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET,
			},
			expErr: false,
			want: &types.MsgZenexTransferRequestResponse{
				SignReqId: 1,
			},
			wantSwap: types.Swap{
				Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				SwapId:  1,
				Status:  types.SwapStatus_SWAP_STATUS_REQUESTED,
				Pair:    types.TradePair_TRADE_PAIR_ROCK_BTC,
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
				SignReqId:      1,
			},
		},
		{
			name: "Pass: Happy Path btcrock",
			input: &types.MsgZenexTransferRequest{
				Creator:        types.DefaultParams().BtcProxyAddress,
				SwapId:         zenextestutil.SampleSwap[1].SwapId,
				UnsignedPlusTx: zenextestutil.Btc_tx[0].UnsignedPlusTx,
				CacheId:        zenextestutil.Btc_tx[0].CacheId,
				DataForSigning: zenextestutil.Btc_tx[0].DataForSigning,
				WalletType:     treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET,
			},
			expErr: false,
			want: &types.MsgZenexTransferRequestResponse{
				SignReqId: 1,
			},
			wantSwap: types.Swap{
				Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				SwapId:  2,
				Status:  types.SwapStatus_SWAP_STATUS_REQUESTED,
				Pair:    types.TradePair_TRADE_PAIR_BTC_ROCK,
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
				SignReqId:      1,
			},
		},
		{
			name: "FAIL: invalid msg sender",
			input: &types.MsgZenexTransferRequest{
				Creator:        "zen1qwnafe2s9eawhah5x6v4593v3tljdntl9zcqpn",
				SwapId:         zenextestutil.SampleSwap[1].SwapId,
				UnsignedPlusTx: zenextestutil.Btc_tx[0].UnsignedPlusTx,
				CacheId:        zenextestutil.Btc_tx[0].CacheId,
				DataForSigning: zenextestutil.Btc_tx[0].DataForSigning,
				WalletType:     treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET,
			},
			expErr:    true,
			expErrMsg: "message creator is not the btc proxy address",
		},
		{
			name: "FAIL: swap id not found",
			input: &types.MsgZenexTransferRequest{
				Creator:        types.DefaultParams().BtcProxyAddress,
				SwapId:         1000,
				UnsignedPlusTx: zenextestutil.Btc_tx[0].UnsignedPlusTx,
				CacheId:        zenextestutil.Btc_tx[0].CacheId,
				DataForSigning: zenextestutil.Btc_tx[0].DataForSigning,
				WalletType:     treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET,
			},
			expErr:    true,
			expErrMsg: collections.ErrNotFound.Error(),
		},
		{
			name: "FAIL: swap status not requested",
			input: &types.MsgZenexTransferRequest{
				Creator:        "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				SwapId:         zenextestutil.SampleSwap[2].SwapId,
				UnsignedPlusTx: zenextestutil.Btc_tx[0].UnsignedPlusTx,
				CacheId:        zenextestutil.Btc_tx[0].CacheId,
				DataForSigning: zenextestutil.Btc_tx[0].DataForSigning,
				WalletType:     treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET,
			},
			expErr:    true,
			expErrMsg: "swap status is not requested",
		},
		{
			name: "FAIL: invalid wallet type",
			input: &types.MsgZenexTransferRequest{
				Creator:        "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				SwapId:         zenextestutil.SampleSwap[0].SwapId,
				UnsignedPlusTx: zenextestutil.Btc_tx[0].UnsignedPlusTx,
				CacheId:        zenextestutil.Btc_tx[0].CacheId,
				DataForSigning: zenextestutil.Btc_tx[0].DataForSigning,
				WalletType:     treasurytypes.WalletType_WALLET_TYPE_NATIVE,
			},
			expErr:    true,
			expErrMsg: "invalid wallet type: WALLET_TYPE_NATIVE",
		},
		{
			name: "Pass: Reject Swap",
			input: &types.MsgZenexTransferRequest{
				Creator:        types.DefaultParams().BtcProxyAddress,
				SwapId:         zenextestutil.SampleSwap[0].SwapId,
				UnsignedPlusTx: nil,
				CacheId:        nil,
				DataForSigning: nil,
				WalletType:     treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET,
				RejectReason:   "some_reject_reason",
			},
			expErr: false,
			want: &types.MsgZenexTransferRequestResponse{
				SignReqId: 0,
			},
			wantSwap: types.Swap{
				Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				SwapId:  1,
				Status:  types.SwapStatus_SWAP_STATUS_REJECTED,
				Pair:    types.TradePair_TRADE_PAIR_ROCK_BTC,
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
				RejectReason:   "some_reject_reason",
				ZenexPoolKeyId: 3,
				Workspace:      "workspace14a2hpadpsy9h4auve2z8lw",
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			err := s.zenexKeeper.SwapsCount.Set(s.ctx, 0)
			s.Require().NoError(err)

			// Set up the swap in the store first
			if tt.wantSwap != (types.Swap{}) {
				for _, swap := range zenextestutil.SampleSwap {
					err = s.zenexKeeper.SwapsStore.Set(s.ctx, swap.SwapId, swap)
					s.Require().NoError(err)
				}
			} else {
				for _, swap := range zenextestutil.SampleSwap {
					if swap.SwapId == tt.input.SwapId {
						err = s.zenexKeeper.SwapsStore.Set(s.ctx, swap.SwapId, swap)
						s.Require().NoError(err)
						break
					}
				}
			}

			params := types.DefaultParams()
			s.zenexKeeper.SetParams(s.ctx, params)

			if !tt.expErr {
				swap, err := s.zenexKeeper.SwapsStore.Get(s.ctx, tt.input.SwapId)
				s.Require().NoError(err)

				s.treasuryKeeper.EXPECT().HandleSignatureRequest(s.ctx, gomock.Any()).Return(&treasurytypes.MsgNewSignatureRequestResponse{
					SigReqId: tt.want.SignReqId,
				}, nil).AnyTimes()

				if swap.BtcKeyId > 0 {
					s.treasuryKeeper.EXPECT().GetKey(s.ctx, swap.BtcKeyId).Return(&treasurytestutil.DefaultKeys[swap.BtcKeyId], nil).AnyTimes()
				}
			}

			zenexTransferResult, err := s.msgServer.ZenexTransferRequest(s.ctx, tt.input)

			if tt.expErr {
				s.Require().Error(err)
				if tt.expErrMsg == collections.ErrNotFound.Error() {
					s.Require().ErrorIs(err, collections.ErrNotFound)
				} else {
					s.Require().Equal(tt.expErrMsg, err.Error())
				}
				s.Require().Nil(zenexTransferResult)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(zenexTransferResult)
				s.Require().Equal(tt.want.SignReqId, zenexTransferResult.SignReqId)

				swapResult, err := s.zenexKeeper.SwapsStore.Get(s.ctx, tt.input.SwapId)
				s.Require().NoError(err)
				s.Require().Equal(tt.wantSwap, swapResult)
			}
		})
	}
}
