package testutil

import (
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
)

var DefaultMints = []types.Bridge{
	{
		Creator:          "zenrock1vpg325p5a7v73jaqyk02muh6n22p852h60anxx",
		Amount:           100,
		BlockHeight:      100,
		Denom:            "urock",
		DestinationChain: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		Id:               1,
		RecipientAddress: "8e7ekBeWmMdU6sJqnCwhm3P2bHBpNwZZ6RNiWJyrMyYz",
		SourceAddress:    "zenrock1vpg325p5a7v73jaqyk02muh6n22p852h60anxx",
		SourceChain:      "cosmos:zenrock",
		State:            types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		TxId:             7,
	},
	{
		Creator:          "zen1xac8tywg8quhwwghr8w80r26j0whae4pjw28l4",
		Amount:           4000000,
		BlockHeight:      1060,
		Denom:            "urock",
		DestinationChain: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		Id:               2,
		RecipientAddress: "8e7ekBeWmMdU6sJqnCwhm3P2bHBpNwZZ6RNiWJyrMyYz",
		SourceAddress:    "zen1xac8tywg8quhwwghr8w80r26j0whae4pjw28l4",
		SourceChain:      "cosmos:zenrock",
		State:            types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		TxId:             8,
	},
	{
		Creator:          "zenrock1vpg325p5a7v73jaqyk02muh6n22p852h60anxx",
		Amount:           100,
		BlockHeight:      100,
		Denom:            "urock",
		DestinationChain: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		Id:               3,
		RecipientAddress: "HCKJKQ8xFfqm5BxTA6LaUCWYXNShYXF4eAwC7NhBYXkj",
		SourceAddress:    "zenrock1vpg325p5a7v73jaqyk02muh6n22p852h60anxx",
		SourceChain:      "cosmos:zenrock",
		State:            types.BridgeStatus_BRIDGE_STATUS_PENDING,
		TxId:             100,
	},
}

var DefaultBurns = []types.Bridge{
	{
		Creator:          "zenrock1vpg325p5a7v73jaqyk02muh6n22p852h60anxx",
		Amount:           400000,
		BlockHeight:      1096,
		Denom:            "urock",
		DestinationChain: "cosmos:zenrock",
		Id:               4,
		RecipientAddress: "zenrock1vpg325p5a7v73jaqyk02muh6n22p852h60anxx",
		SourceAddress:    "8e7ekBeWmMdU6sJqnCwhm3P2bHBpNwZZ6RNiWJyrMyYz",
		SourceChain:      "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		State:            types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		TxHash:           "4JGeiMjKPmXYuPiQQPxmsDKsxSTJKoTqdzDU47NeomPjUjT7SgwjDzm3UBxErnC19DNkFRKhAHVdZ1J1J8236gtw",
	},
	{
		Creator:          "zen1xac8tywg8quhwwghr8w80r26j0whae4pjw28l4",
		Amount:           500000,
		BlockHeight:      1116,
		Denom:            "urock",
		DestinationChain: "cosmos:zenrock",
		Id:               5,
		RecipientAddress: "zen1xac8tywg8quhwwghr8w80r26j0whae4pjw28l4",
		SourceAddress:    "8e7ekBeWmMdU6sJqnCwhm3P2bHBpNwZZ6RNiWJyrMyYz",
		SourceChain:      "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		State:            types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		TxHash:           "51JqTggEpYGSkg1yQRH5cYbVG2vXYdEp5pXdg2zZDmX7jrhYRyD2vswhrPQ9hK5FoTxvUMWQd4UWopWEioGrf7bT",
	},
	{
		Creator:          "zen1xac8tywg8quhwwghr8w80r26j0whae4pjw28l4",
		Amount:           600000,
		BlockHeight:      1116,
		Denom:            "urock",
		DestinationChain: "cosmos:zenrock",
		Id:               6,
		RecipientAddress: "zen1xac8tywg8quhwwghr8w80r26j0whae4pjw28l4",
		SourceAddress:    "8e7ekBeWmMdU6sJqnCwhm3P2bHBpNwZZ6RNiWJyrMyYz",
		SourceChain:      "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		State:            types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		TxHash:           "3pske6HU6z9NGz2zPsQ9uJZMxJPihNo8FpJqqzWNWTb2Bp4i5ob58hy3mdX2rzxpvNz5R6mfaGvXT97DQKM3kYBB",
	},
	{
		Creator:          "zenrock1vpg325p5a7v73jaqyk02muh6n22p852h60anxx",
		Amount:           700000,
		BlockHeight:      1116,
		Denom:            "urock",
		DestinationChain: "cosmos:zenrock",
		Id:               7,
		RecipientAddress: "zenrock1vpg325p5a7v73jaqyk02muh6n22p852h60anxx",
		SourceAddress:    "8e7ekBeWmMdU6sJqnCwhm3P2bHBpNwZZ6RNiWJyrMyYz",
		SourceChain:      "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		State:            types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		TxHash:           "qbBwRp9wq6sjFzNYjjRdrDT47HqiBXMFEt6DA8H2fK8r84qcdLweAEgxx51x1sFFWESAkY43StUQVeDsmnKosRX",
	},
	{
		Creator:          "zen1xac8tywg8quhwwghr8w80r26j0whae4pjw28l4",
		Amount:           800000,
		BlockHeight:      1116,
		Denom:            "urock",
		DestinationChain: "cosmos:zenrock",
		Id:               8,
		RecipientAddress: "zen1xac8tywg8quhwwghr8w80r26j0whae4pjw28l4",
		SourceAddress:    "8e7ekBeWmMdU6sJqnCwhm3P2bHBpNwZZ6RNiWJyrMyYz",
		SourceChain:      "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		State:            types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		TxHash:           "5aSgL1nJcBMsgbCtmkJ79bWpDQtBfy72u9uoa4bkvWSYDuNyAtuweyoyqksWbiSVPccsAqHjorez1sGb3f9LXpgp",
	},
	{
		Creator:          "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		Amount:           50,
		BlockHeight:      1166,
		Denom:            "urock",
		DestinationChain: "cosmos:zenrock",
		Id:               9,
		RecipientAddress: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SourceAddress:    "8e7ekBeWmMdU6sJqnCwhm3P2bHBpNwZZ6RNiWJyrMyYz",
		SourceChain:      "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		State:            types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		TxHash:           "2fzLENYbCtriyWgDjmn345Jhz1TyX5LSTg5srfadRf4mUqXRTGWbcJDVULPgoZrJ8HurmeKNZrbH2gGYNiFmbFoQ",
	},
}
