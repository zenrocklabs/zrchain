package types

import treasurytypes "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"

func Caip2ToSolananNetwork(caip string) treasurytypes.SolanaNetworkType {
	switch caip {
	case "solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp:7S3P4HxJpyyigGzodYwHtCxZyUQe9JiBMHyRWXArAaKv":
		return treasurytypes.SolanaNetworkType_MAINNET
	case "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1:DYw8jCTfwHNRJhhmFcbXvVDTqWMEVFBX6ZKUmG5CNSKK":
		return treasurytypes.SolanaNetworkType_DEVNET
	case "solana:4uhcVJyU9pJkvQyS88uRDiswHXSCkY3z:6LmSRCiu3z6NCSpF19oz1pHXkYkN4jWbj9K1nVELpDkT":
		return treasurytypes.SolanaNetworkType_TESTNET
	default:
		return treasurytypes.SolanaNetworkType_UNDEFINED
	}

}
