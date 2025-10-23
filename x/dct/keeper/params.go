package keeper

import (
	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

var (
	defaultZenZECSolana = &types.Solana{
		SignerKeyId:         17,
		ProgramId:           "Ney4F3GUe5BhPP1cFEmXRVWtzSUGhG4DrRkfERRYedh",
		NonceAccountKey:     19,
		NonceAuthorityKey:   18,
		MintAddress:         "4q9DEzEHLqNG637jsGMYSg8E56SbotNcGeH3GjtaYYJT",
		FeeWallet:           "4GCX9fgq9gzBH282tVMnwebAW6gW8QX4N7JzGD1djYWT",
		Fee:                 0,
		MultisigKeyAddress:  "F35sSPrvaKioWbr4dhrPti2KqqA3hkrJjN9nxLoUgLZ6",
		Btl:                 20,
		EventStoreProgramId: "HbUWCsvZzkQtHakTX6QovPKcaCeP1Pf34W9Bpw6j18J8",
	}

	defaultZenZECCfg = types.AssetParams{
		Asset:               types.Asset_ASSET_ZENZEC,
		DepositKeyringAddr:  "keyring1k6vc6vhp6e6l3rxalue9v4ux",
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
