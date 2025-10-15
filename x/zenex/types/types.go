package types

import (
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/btcsuite/btcd/chaincfg"
)

var ValidPairTypes = []string{TradePair_TRADE_PAIR_ROCK_BTC.String(), TradePair_TRADE_PAIR_BTC_ROCK.String()}

const (
	ZenexCollectorName = "zenex_collector"
)

// ChainFromWalletType returns the corresponding chain configuration parameters based on the provided wallet type.
func ChainFromWalletType(walletType treasurytypes.WalletType) *chaincfg.Params {
	switch walletType {
	case treasurytypes.WalletType_WALLET_TYPE_BTC_MAINNET:
		return &chaincfg.MainNetParams
	case treasurytypes.WalletType_WALLET_TYPE_BTC_TESTNET:
		return &chaincfg.TestNet3Params
	case treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET:
		return &chaincfg.RegressionNetParams
	case treasurytypes.WalletType_WALLET_TYPE_ZCASH_MAINNET:
		return &chaincfg.MainNetParams // ZCash uses similar params structure
	case treasurytypes.WalletType_WALLET_TYPE_ZCASH_TESTNET:
		return &chaincfg.TestNet3Params // ZCash testnet (also used for regtest)
	default:
		return nil
	}
}
