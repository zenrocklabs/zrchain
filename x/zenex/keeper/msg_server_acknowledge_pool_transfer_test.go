package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"

	appparams "github.com/Zenrock-Foundation/zrchain/v6/app/params"
	treasurytestutil "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/testutil"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	zenextestutil "github.com/Zenrock-Foundation/zrchain/v6/x/zenex/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) TestMsgAcknowledgePoolTransfer() {

	tests := []struct {
		name      string
		input     *types.MsgAcknowledgePoolTransfer
		expErr    bool
		expErrMsg string
		wantSwap  types.Swap
	}{
		{
			name: "Pass: Happy Path rockbtc",
			input: &types.MsgAcknowledgePoolTransfer{
				Creator:      "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
				SwapId:       4,
				SourceTxHash: "source_tx_hash",
				Status:       types.SwapStatus_SWAP_STATUS_COMPLETED,
			},
			expErr: false,
			wantSwap: types.Swap{
				Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				SwapId:  4,
				Status:  types.SwapStatus_SWAP_STATUS_COMPLETED,
				Pair:    "rockbtc",
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
				SourceTxHash:   "source_tx_hash",
			},
		},
		{
			name: "Pass: Happy Path btcrock",
			input: &types.MsgAcknowledgePoolTransfer{
				Creator:      "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
				SwapId:       5,
				SourceTxHash: "source_tx_hash",
				Status:       types.SwapStatus_SWAP_STATUS_COMPLETED,
			},
			expErr: false,
			wantSwap: types.Swap{
				Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				SwapId:  5,
				Status:  types.SwapStatus_SWAP_STATUS_COMPLETED,
				Pair:    "btcrock",
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
				SourceTxHash:   "source_tx_hash",
			},
		},
		{
			name: "Pass: swap status rejected",
			input: &types.MsgAcknowledgePoolTransfer{
				Creator:      "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
				SwapId:       6,
				SourceTxHash: "source_tx_hash",
				Status:       types.SwapStatus_SWAP_STATUS_REJECTED,
				RejectReason: "reject_reason",
			},
			expErr: false,
			wantSwap: types.Swap{
				Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				SwapId:  6,
				Status:  types.SwapStatus_SWAP_STATUS_REJECTED,
				Pair:    "btcrock",
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
				RejectReason:   "reject_reason",
			},
		},
		{
			name: "FAIL: msg status not swap transfer completed/rejected",
			input: &types.MsgAcknowledgePoolTransfer{
				Creator:      "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
				SwapId:       6,
				SourceTxHash: "source_tx_hash",
				Status:       types.SwapStatus_SWAP_STATUS_SWAP_TRANSFER_REQUESTED,
			},
			expErr:    true,
			expErrMsg: fmt.Sprintf("msg status is not completed or rejected: %s", types.SwapStatus_SWAP_STATUS_SWAP_TRANSFER_REQUESTED.String()),
		},
		{
			name: "FAIL: swap status is already completed/rejected",
			input: &types.MsgAcknowledgePoolTransfer{
				Creator:      "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
				SwapId:       7,
				SourceTxHash: "source_tx_hash",
				Status:       types.SwapStatus_SWAP_STATUS_COMPLETED,
			},
			expErr:    true,
			expErrMsg: fmt.Sprintf("swap status is not swap transfer requested: %s", types.SwapStatus_SWAP_STATUS_COMPLETED.String()),
		},
		{
			name: "FAIL: msg creator is not the btc proxy address",
			input: &types.MsgAcknowledgePoolTransfer{
				Creator:      "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				SwapId:       7,
				SourceTxHash: "source_tx_hash",
				Status:       types.SwapStatus_SWAP_STATUS_COMPLETED,
			},
			expErr:    true,
			expErrMsg: fmt.Sprintf("message creator is not the btc proxy address: %s", "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.zenexKeeper.SwapsCount.Set(s.ctx, 0)
			s.Require().NoError(err)

			params := types.DefaultParams()
			s.zenexKeeper.SetParams(s.ctx, params)

			// Set up the swap in the store first
			if tt.wantSwap != (types.Swap{}) {
				for _, swap := range zenextestutil.SampleSwap {
					err = s.zenexKeeper.SwapsStore.Set(s.ctx, swap.SwapId, swap)
					s.Require().NoError(err)
				}
			}

			if !tt.expErr && tt.wantSwap.Pair == "btcrock" && tt.input.Status == types.SwapStatus_SWAP_STATUS_COMPLETED {
				s.treasuryKeeper.EXPECT().GetKey(s.ctx, tt.wantSwap.RockKeyId).Return(&treasurytestutil.DefaultKeys[tt.wantSwap.RockKeyId-1], nil)
				senderAddress, err := treasurytypes.NativeAddress(&treasurytestutil.DefaultKeys[tt.wantSwap.RockKeyId-1], "zen")
				if err != nil {
					s.T().Fatalf("failed to convert sender key to zenrock address: %v", err)
				}
				s.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(s.ctx, types.ZenexCollectorName, sdk.MustAccAddressFromBech32(senderAddress), sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, math.NewIntFromUint64(tt.wantSwap.Data.AmountOut)))).Return(nil)
			}

			resp, err := s.msgServer.AcknowledgePoolTransfer(s.ctx, tt.input)

			if tt.expErr {
				s.Require().Error(err)
				s.Require().Equal(tt.expErrMsg, err.Error())
				s.Require().Nil(resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				swap, err := s.zenexKeeper.SwapsStore.Get(s.ctx, tt.input.SwapId)
				s.Require().NoError(err)
				s.Require().Equal(tt.wantSwap, swap)
				s.Require().NotNil(resp)
			}
		})
	}
}
