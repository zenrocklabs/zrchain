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
