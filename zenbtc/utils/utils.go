package utils

import (
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
)

// ChainFromString returns the corresponding chain configuration parameters based on the provided chain name.
func ChainFromString(chainName string) *chaincfg.Params {
	chainName = strings.ToLower(chainName)
	switch chainName {
	case "mainnet":
		return &chaincfg.MainNetParams
	case "testnet", "testnet3":
		return &chaincfg.TestNet3Params
	case "regtest", "regnet":
		return &chaincfg.RegressionNetParams
	case "testnet4":
		return &chaincfg.TestNet3Params //TestNet4Params not available yet (22/7/24)

	default:
		return nil
	}
}
