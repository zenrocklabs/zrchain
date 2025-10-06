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
		Status:  types.SwapStatus_SWAP_STATUS_INITIATED,
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
		Status:  types.SwapStatus_SWAP_STATUS_INITIATED,
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
		SwapId:  4,
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
		SignReqId:      1,
		SourceTxHash:   "source_tx_hash",
	},
	{

		Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SwapId:  5,
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
		SignReqId:      1,
		SourceTxHash:   "source_tx_hash",
	},
	{
		Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SwapId:  6,
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
		SignReqId:      1,
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
		SignReqId:      1,
	},
	{
		Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SwapId:  8,
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
			Price:     SampleBtcRockPrice,
			AmountIn:  4400000,
			AmountOut: 2000,
		},
		RockKeyId:      1,
		BtcKeyId:       2,
		ZenexPoolKeyId: 3,
		Workspace:      "workspace14a2hpadpsy9h4auve2z8lw",
		SignReqId:      1,
	},
	{
		Creator: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SwapId:  9,
		Status:  types.SwapStatus_SWAP_STATUS_INITIATED,
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
		RockKeyId:      14,
		BtcKeyId:       15,
		ZenexPoolKeyId: 16,
		Workspace:      "workspace14a2hpadpsy9h4auve2z8lw",
	},
}

var Btc_tx = []struct {
	id             uint64
	DataForSigning []*types.InputHashes
	CacheId        []byte
	UnsignedPlusTx []byte
}{
	{
		DataForSigning: []*types.InputHashes{
			{
				KeyId: 15,
				Hash:  "349e10f8545d66c6c95a3aae132dee0e8612f1d5f40d19596d2c507802f7f2ee",
			},
		},
		CacheId:        []byte("600fa96d9b87050e4ada945cffa616ef2be0d84b22cfed4790e6fd2a0152d719"),
		UnsignedPlusTx: []byte("01000000000101fcf96c9e03fc23b11b633695effe7d220998e2327a0ddf260c669187049aa2a20100000000ffffffff02aa260000000000001600142ff47fdcce66e16b0e66c3d6611577115d8a5834723b0100000000001600144b01225b08241dbc749aaff8f30446508b9b4dd70808a08601000000000021022b9d2e93b7b69bc88774d19b6567a3c5c9031c5f05c197807f18cdcac81bda661600144b01225b08241dbc749aaff8f30446508b9b4dd708010000000000000020349e10f8545d66c6c95a3aae132dee0e8612f1d5f40d19596d2c507802f7f2ee0008746573746e6574340e4c4f43414c5a454e524f434b313000000000"),
	},
}
