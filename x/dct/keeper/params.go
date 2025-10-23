package keeper

import (
	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

var (
	defaultZenZECSolana = &types.Solana{
		SignerKeyId:         7,
		ProgramId:           "7q3u7bL1nKzftYFbwUE8yuzhzbQxFwM7xyMh1cZFzenZ",
		NonceAccountKey:     9,
		NonceAuthorityKey:   8,
		MintAddress:         "ZC3hZPnfYg1y5SP62x9XyLJmMnt4zUpfcu8JZenZEC",
		FeeWallet:           "FvzecFee1YdzWm7Nzvy2Kc7CquBrEPNvkMRPQHrfeee",
		Fee:                 0,
		MultisigKeyAddress:  "7ZeCmxg9APm9wJgZtPy1X6Qt23afVqzVrqVZC9HexbBP",
		Btl:                 20,
		EventStoreProgramId: "7zecEvtStreProg11111111111111111111111111111",
	}

	defaultZenZECCfg = types.AssetParams{
		Asset:               types.Asset_ASSET_ZENZEC,
		DepositKeyringAddr:  "keyring1pfnq7r04rept47gaf5cpdew2",
		StakerKeyId:         6,
		EthMinterKeyId:      2,
		UnstakerKeyId:       4,
		CompleterKeyId:      7,
		RewardsDepositKeyId: 5,
		ChangeAddressKeyIds: []uint64{3},
		ProxyAddress:        "zen1trdxe6r48aqvhm026akay7tjnzuarf2rxuz0ah",
		EthTokenAddr:        "",
		ControllerAddr:      "",
		Solana:              defaultZenZECSolana,
	}
)

func DefaultParams() *types.Params {
	return &types.Params{Assets: []types.AssetParams{defaultZenZECCfg}}
}
