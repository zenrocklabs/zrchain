package testutil

import (
	math "cosmossdk.io/math"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
)

var (
	SampleRockBtcPrice = math.LegacyNewDecFromInt(math.NewInt(4400000))
	SampleBtcRockPrice = math.LegacyNewDecFromInt(math.NewInt(2272727))
)

var SampleSwap = []types.Swap{
	{
		Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SwapId:  1,
		Status:  types.SwapStatus_SWAP_STATUS_REQUESTED,
		Pair:    "rockbtc",
		Data: &types.SwapData{
			BaseToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_ROCK,
				PriceUSD:  SampleRockBtcPrice,
				Precision: 6,
			},
			QuoteToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_BTC,
				PriceUSD:  SampleBtcRockPrice,
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
	{
		Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SwapId:  2,
		Status:  types.SwapStatus_SWAP_STATUS_REQUESTED,
		Pair:    "btcrock",
		Data: &types.SwapData{
			BaseToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_BTC,
				PriceUSD:  SampleBtcRockPrice,
				Precision: 8,
			},
			QuoteToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_ROCK,
				PriceUSD:  SampleRockBtcPrice,
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
	{
		Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SwapId:  3,
		Status:  types.SwapStatus_SWAP_STATUS_SWAP_TRANSFER_REQUESTED,
		Pair:    "btcrock",
		Data: &types.SwapData{
			BaseToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_BTC,
				PriceUSD:  SampleBtcRockPrice,
				Precision: 8,
			},
			QuoteToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_ROCK,
				PriceUSD:  SampleRockBtcPrice,
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
	{
		Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SwapId:  4,
		Status:  types.SwapStatus_SWAP_STATUS_SWAP_TRANSFER_REQUESTED,
		Pair:    "rockbtc",
		Data: &types.SwapData{
			BaseToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_ROCK,
				PriceUSD:  SampleRockBtcPrice,
				Precision: 6,
			},
			QuoteToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_BTC,
				PriceUSD:  SampleBtcRockPrice,
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
	{

		Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SwapId:  5,
		Status:  types.SwapStatus_SWAP_STATUS_SWAP_TRANSFER_REQUESTED,
		Pair:    "btcrock",
		Data: &types.SwapData{
			BaseToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_BTC,
				PriceUSD:  SampleBtcRockPrice,
				Precision: 8,
			},
			QuoteToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_ROCK,
				PriceUSD:  SampleRockBtcPrice,
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
	{
		Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SwapId:  6,
		Status:  types.SwapStatus_SWAP_STATUS_SWAP_TRANSFER_REQUESTED,
		Pair:    "btcrock",
		Data: &types.SwapData{
			BaseToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_BTC,
				PriceUSD:  SampleBtcRockPrice,
				Precision: 8,
			},
			QuoteToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_ROCK,
				PriceUSD:  SampleRockBtcPrice,
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
	{
		Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SwapId:  7,
		Status:  types.SwapStatus_SWAP_STATUS_COMPLETED,
		Pair:    "btcrock",
		Data: &types.SwapData{
			BaseToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_BTC,
				PriceUSD:  SampleBtcRockPrice,
				Precision: 8,
			},
			QuoteToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_ROCK,
				PriceUSD:  SampleRockBtcPrice,
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
}
