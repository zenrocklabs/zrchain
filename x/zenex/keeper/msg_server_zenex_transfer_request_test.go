package keeper_test

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	treasurytestutil "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/testutil"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	zenextestutil "github.com/Zenrock-Foundation/zrchain/v6/x/zenex/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
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
				Creator:    types.DefaultParams().BtcProxyAddress,
				SwapId:     zenextestutil.SampleSwap[0].SwapId,
				UnsignedTx: []byte("unsigned_tx"),
				WalletType: treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET,
			},
			expErr: false,
			want: &types.MsgZenexTransferRequestResponse{
				SignTxId: 1,
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
				SignTxId:       1,
			},
		},
		{
			name: "Pass: Happy Path btcrock",
			input: &types.MsgZenexTransferRequest{
				Creator:    types.DefaultParams().BtcProxyAddress,
				SwapId:     zenextestutil.SampleSwap[1].SwapId,
				UnsignedTx: []byte("unsigned_tx"),
				WalletType: treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET,
			},
			expErr: false,
			want: &types.MsgZenexTransferRequestResponse{
				SignTxId: 1,
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
				SignTxId:       1,
			},
		},
		{
			name: "FAIL: invalid msg sender",
			input: &types.MsgZenexTransferRequest{
				Creator:    "zen1qwnafe2s9eawhah5x6v4593v3tljdntl9zcqpn",
				SwapId:     zenextestutil.SampleSwap[1].SwapId,
				UnsignedTx: []byte("unsigned_tx"),
				WalletType: treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET,
			},
			expErr:    true,
			expErrMsg: "message creator is not the btc proxy address",
		},
		{
			name: "FAIL: swap id not found",
			input: &types.MsgZenexTransferRequest{
				Creator:    types.DefaultParams().BtcProxyAddress,
				SwapId:     1000,
				UnsignedTx: []byte("unsigned_tx"),
				WalletType: treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET,
			},
			expErr:    true,
			expErrMsg: collections.ErrNotFound.Error(),
		},
		{
			name: "FAIL: swap status not requested",
			input: &types.MsgZenexTransferRequest{
				Creator:    "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				SwapId:     zenextestutil.SampleSwap[2].SwapId,
				UnsignedTx: []byte("unsigned_tx"),
				WalletType: treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET,
			},
			expErr:    true,
			expErrMsg: "swap status is not requested",
		},
		{
			name: "FAIL: invalid wallet type",
			input: &types.MsgZenexTransferRequest{
				Creator:    "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				SwapId:     zenextestutil.SampleSwap[0].SwapId,
				UnsignedTx: []byte("unsigned_tx"),
				WalletType: treasurytypes.WalletType_WALLET_TYPE_NATIVE,
			},
			expErr:    true,
			expErrMsg: "invalid wallet type: WALLET_TYPE_NATIVE",
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

			var expectedBitcoinTx *treasurytypes.MsgNewSignTransactionRequest
			if !tt.expErr {
				switch tt.wantSwap.Pair {
				case types.TradePair_TRADE_PAIR_ROCK_BTC:
					expectedBitcoinTx = treasurytypes.NewMsgNewSignTransactionRequest(
						tt.wantSwap.Creator,            // txCreator from swap
						[]uint64{tt.wantSwap.BtcKeyId}, // senderKeyId from swap for "rockbtc" pair
						tt.input.WalletType,
						tt.input.UnsignedTx,
						nil,
						treasurytypes.DefaultParams().DefaultBtl,
					)
				case types.TradePair_TRADE_PAIR_BTC_ROCK:
					expectedBitcoinTx = treasurytypes.NewMsgNewSignTransactionRequest(
						params.BtcProxyAddress,
						[]uint64{tt.wantSwap.ZenexPoolKeyId},
						tt.input.WalletType,
						tt.input.UnsignedTx,
						nil,
						treasurytypes.DefaultParams().DefaultBtl,
					)
				}

				s.treasuryKeeper.EXPECT().MakeSignTransactionRequest(s.ctx, expectedBitcoinTx).Return(&treasurytypes.MsgNewSignTransactionRequestResponse{
					Id: tt.want.SignTxId,
				}, nil)
			}

			// Only set up GetKey expectation if we have a valid BtcKeyId
			if tt.wantSwap.BtcKeyId > 0 {
				s.treasuryKeeper.EXPECT().GetKey(s.ctx, tt.wantSwap.BtcKeyId).Return(&treasurytestutil.DefaultKeys[tt.wantSwap.BtcKeyId-1], nil).AnyTimes()
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
				s.Require().Equal(tt.want.SignTxId, zenexTransferResult.SignTxId)

				swapResult, err := s.zenexKeeper.SwapsStore.Get(s.ctx, tt.input.SwapId)
				s.Require().NoError(err)
				s.Require().Equal(tt.wantSwap, swapResult)
			}
		})
	}
}
