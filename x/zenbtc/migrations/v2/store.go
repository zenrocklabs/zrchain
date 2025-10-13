package v2

import (
	"strings"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenbtc/types"
)

func UpdateParams(ctx sdk.Context, params collections.Item[types.Params]) error {
	paramsMap := map[string]types.Params{
		"zenrock": { // local
			DepositKeyringAddr:  "keyring1hpyh7xqr2w7h4eas5y8twnsg",
			StakerKeyID:         1,
			EthMinterKeyID:      2,
			UnstakerKeyID:       3,
			CompleterKeyID:      4,
			RewardsDepositKeyID: 5,
			ChangeAddressKeyIDs: []uint64{6},
			BitcoinProxyAddress: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
			EthTokenAddr:        "0x7692E9a796001FeE9023853f490A692bAB2E4834",
			ControllerAddr:      "0x2844bd31B68AE5a0335c672e6251e99324441B73",
		},
		"amber": { // devnet
			DepositKeyringAddr:  "keyring1k6vc6vhp6e6l3rxalue9v4ux",
			StakerKeyID:         1,
			EthMinterKeyID:      2,
			UnstakerKeyID:       3,
			CompleterKeyID:      4,
			RewardsDepositKeyID: 5,
			ChangeAddressKeyIDs: []uint64{6},
			BitcoinProxyAddress: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
			EthTokenAddr:        "0x7692E9a796001FeE9023853f490A692bAB2E4834",
			ControllerAddr:      "0x2844bd31B68AE5a0335c672e6251e99324441B73",
		},
		"gardia": { // testnet
			DepositKeyringAddr:  "keyring1k6vc6vhp6e6l3rxalue9v4ux",
			StakerKeyID:         1,
			EthMinterKeyID:      2,
			UnstakerKeyID:       3,
			CompleterKeyID:      4,
			RewardsDepositKeyID: 5,
			ChangeAddressKeyIDs: []uint64{6},
			BitcoinProxyAddress: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
			EthTokenAddr:        "0xfA32a2D7546f8C7c229F94E693422A786DaE5E18",
			ControllerAddr:      "0xaCE3634AAd9bCC48ef6A194f360F7ACe51F7d9f1",
		},
		"diamond": { // mainnet
			DepositKeyringAddr:  "keyring1k6vc6vhp6e6l3rxalue9v4ux",
			StakerKeyID:         24,
			EthMinterKeyID:      17,
			UnstakerKeyID:       19,
			CompleterKeyID:      28,
			RewardsDepositKeyID: 20,
			ChangeAddressKeyIDs: []uint64{18},
			BitcoinProxyAddress: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
			EthTokenAddr:        "",
			ControllerAddr:      "",
		},
	}

	chainID := ctx.ChainID()
	if chainID == "" {
		chainID = "zenrock"
	}

	newParams := types.Params{}

	for prefix, paramSet := range paramsMap {
		if strings.HasPrefix(chainID, prefix) {
			newParams = paramSet
			break
		}
	}

	if err := params.Set(ctx, newParams); err != nil {
		return err
	}

	return nil
}
