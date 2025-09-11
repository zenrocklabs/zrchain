package testutil

import (
	math "cosmossdk.io/math"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
)

var (
	// ROCK price: 0.025 USD per ROCK (1,000,000 urock)
	// BTC price: 110,000 USD per BTC (100,000,000 satoshis)
	// ROCK/BTC price = 0.025 / 110,000 = 0.00000022727 satoshis per urock
	SampleRockBtcPrice, _ = math.LegacyNewDecFromStr("0.00000022727")
	// BTC/ROCK price = 110,000 / 0.025 = 4,400,000 urock per satoshi
	SampleBtcRockPrice, _ = math.LegacyNewDecFromStr("4400000")
	// ROCK price: 0.025 USD per ROCK
	SampleRockUSDPrice, _ = math.LegacyNewDecFromStr("0.025")
	// BTC price: 110,000 USD per BTC
	SampleBtcUSDPrice, _ = math.LegacyNewDecFromStr("110000")
)

var SampleSwap = []types.Swap{
	{
		Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SwapId:  1,
		Status:  types.SwapStatus_SWAP_STATUS_REQUESTED,
		Pair:    types.TradePair_TRADE_PAIR_ROCK_BTC,
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
		Pair:    types.TradePair_TRADE_PAIR_BTC_ROCK,
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
		Pair:    types.TradePair_TRADE_PAIR_BTC_ROCK,
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
		Pair:    types.TradePair_TRADE_PAIR_ROCK_BTC,
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
		Pair:    types.TradePair_TRADE_PAIR_BTC_ROCK,
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
		Pair:    types.TradePair_TRADE_PAIR_BTC_ROCK,
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
		Pair:    types.TradePair_TRADE_PAIR_BTC_ROCK,
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
