package keeper

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/stretchr/testify/require"
)

func TestWalletTypeFromChainName(t *testing.T) {
	tests := []struct {
		name           string
		chainName      string
		expectedWallet treasurytypes.WalletType
	}{
		// Bitcoin chains
		{
			name:           "Bitcoin mainnet",
			chainName:      "mainnet",
			expectedWallet: treasurytypes.WalletType_WALLET_TYPE_BTC_MAINNET,
		},
		{
			name:           "Bitcoin testnet",
			chainName:      "testnet",
			expectedWallet: treasurytypes.WalletType_WALLET_TYPE_BTC_TESTNET,
		},
		{
			name:           "Bitcoin testnet3",
			chainName:      "testnet3",
			expectedWallet: treasurytypes.WalletType_WALLET_TYPE_BTC_TESTNET,
		},
		{
			name:           "Bitcoin regtest",
			chainName:      "regtest",
			expectedWallet: treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET,
		},
		{
			name:           "Bitcoin regnet",
			chainName:      "regnet",
			expectedWallet: treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET,
		},
		// Zcash chains - hyphenated format
		{
			name:           "Zcash mainnet (hyphenated)",
			chainName:      "zcash-mainnet",
			expectedWallet: treasurytypes.WalletType_WALLET_TYPE_ZCASH_MAINNET,
		},
		{
			name:           "Zcash testnet (hyphenated)",
			chainName:      "zcash-testnet",
			expectedWallet: treasurytypes.WalletType_WALLET_TYPE_ZCASH_TESTNET,
		},
		{
			name:           "Zcash regtest (hyphenated)",
			chainName:      "zcash-regtest",
			expectedWallet: treasurytypes.WalletType_WALLET_TYPE_ZCASH_TESTNET,
		},
		// Zcash chains - non-hyphenated format
		{
			name:           "Zcash mainnet (non-hyphenated)",
			chainName:      "zcashmainnet",
			expectedWallet: treasurytypes.WalletType_WALLET_TYPE_ZCASH_MAINNET,
		},
		{
			name:           "Zcash testnet (non-hyphenated)",
			chainName:      "zcashtestnet",
			expectedWallet: treasurytypes.WalletType_WALLET_TYPE_ZCASH_TESTNET,
		},
		{
			name:           "Zcash regtest (non-hyphenated)",
			chainName:      "zcashregtest",
			expectedWallet: treasurytypes.WalletType_WALLET_TYPE_ZCASH_TESTNET,
		},
		// Unknown chain
		{
			name:           "Unknown chain",
			chainName:      "unknown-chain",
			expectedWallet: treasurytypes.WalletType_WALLET_TYPE_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &types.MsgVerifyDepositBlockInclusion{
				ChainName: tt.chainName,
			}

			result := WalletTypeFromChainName(msg)
			require.Equal(t, tt.expectedWallet, result,
				"Chain name %s should map to wallet type %s, got %s",
				tt.chainName, tt.expectedWallet, result)
		})
	}
}
