package types

import treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"

func Caip2ToSolanaNetwork(caip string) treasurytypes.SolanaNetworkType {
	switch caip {
	case "solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp":
		return treasurytypes.SolanaNetworkType_MAINNET
	case "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1":
		return treasurytypes.SolanaNetworkType_DEVNET
	case "solana:4uhcVJyU9pJkvQyS88uRDiswHXSCkY3z":
		return treasurytypes.SolanaNetworkType_TESTNET
	case "solana:HK8b7Skns2TX3FvXQxm2mPQbY2nVY8GD":
		return treasurytypes.SolanaNetworkType_REGNET
	default:
		return treasurytypes.SolanaNetworkType_UNDEFINED
	}

}

const (
	ZentpCollectorName = "zentp_collector"
)
