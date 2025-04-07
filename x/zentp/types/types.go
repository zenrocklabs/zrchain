package types

import treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"

func Caip2ToSolananNetwork(caip string) treasurytypes.SolanaNetworkType {
	switch caip {
	case "solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp":
		return treasurytypes.SolanaNetworkType_MAINNET
	case "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1":
		return treasurytypes.SolanaNetworkType_DEVNET
	case "solana:4uhcVJyU9pJkvQyS88uRDiswHXSCkY3z":
		return treasurytypes.SolanaNetworkType_TESTNET
	default:
		return treasurytypes.SolanaNetworkType_UNDEFINED
	}

}
