package v6

import (
	"strings"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v5/x/validation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func UpdateParams(ctx sdk.Context, params collections.Item[types.HVParams]) error {
	oldParams, err := params.Get(ctx)
	if err != nil {
		return err
	}

	currParams := oldParams

	paramsMap := map[string]types.HVParams{
		"zenrock": {
			ZenBTCParams: &types.ZenBTCParams{
				ZenBTCEthBatcherAddr:      "0x912D79F8d489d0d007aBE0E26fD5d2f06BA4A2AA",
				ZenBTCDepositKeyringAddr:  "keyring1hpyh7xqr2w7h4eas5y8twnsg",
				ZenBTCWithdrawerKeyID:     1,
				ZenBTCMinterKeyID:         2,
				ZenBTCChangeAddressKeyIDs: []uint64{3},
				ZenBTCUnstakerKeyID:       4,
				ZenBTCRewardsDepositKeyID: 5,
				BitcoinProxyCreatorID:     "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				StakeableAssets: []*types.AssetData{
					{Asset: types.Asset_ROCK, Precision: 6, PriceUSD: math.LegacyZeroDec()},
					{Asset: types.Asset_zenBTC, Precision: 8, PriceUSD: math.LegacyZeroDec()},
					{Asset: types.Asset_stETH, Precision: 18, PriceUSD: math.LegacyZeroDec()},
				},
			},
		},
		"amber": {
			ZenBTCParams: &types.ZenBTCParams{
				ZenBTCEthBatcherAddr:      "0x912D79F8d489d0d007aBE0E26fD5d2f06BA4A2AA",
				ZenBTCDepositKeyringAddr:  "keyring1hpyh7xqr2w7h4eas5y8twnsg",
				ZenBTCWithdrawerKeyID:     29,
				ZenBTCMinterKeyID:         30,
				ZenBTCChangeAddressKeyIDs: []uint64{31},
				ZenBTCUnstakerKeyID:       32,
				ZenBTCRewardsDepositKeyID: 33,
				BitcoinProxyCreatorID:     "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				StakeableAssets: []*types.AssetData{
					{Asset: types.Asset_ROCK, Precision: 6, PriceUSD: math.LegacyZeroDec()},
					{Asset: types.Asset_zenBTC, Precision: 8, PriceUSD: math.LegacyZeroDec()},
					{Asset: types.Asset_stETH, Precision: 18, PriceUSD: math.LegacyZeroDec()},
				},
			},
		},
		"gardia": {
			ZenBTCParams: &types.ZenBTCParams{
				ZenBTCEthBatcherAddr:      "0xbd903A8D04d98bCA97eD091C87e7A00b7b8F3629",
				ZenBTCDepositKeyringAddr:  "keyring1w887ucurq2nmnj5mq5uaju6a",
				ZenBTCWithdrawerKeyID:     1272,
				ZenBTCMinterKeyID:         1273,
				ZenBTCChangeAddressKeyIDs: []uint64{1274},
				ZenBTCUnstakerKeyID:       1275,
				ZenBTCRewardsDepositKeyID: 1276,
				BitcoinProxyCreatorID:     "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				StakeableAssets: []*types.AssetData{
					{Asset: types.Asset_ROCK, Precision: 6, PriceUSD: math.LegacyZeroDec()},
					{Asset: types.Asset_zenBTC, Precision: 8, PriceUSD: math.LegacyZeroDec()},
					{Asset: types.Asset_stETH, Precision: 18, PriceUSD: math.LegacyZeroDec()},
				},
			},
		},
		"diamond": {
			ZenBTCParams: &types.ZenBTCParams{
				ZenBTCEthBatcherAddr:      "",
				ZenBTCDepositKeyringAddr:  "keyring1k6vc6vhp6e6l3rxalue9v4ux",
				ZenBTCWithdrawerKeyID:     16,
				ZenBTCMinterKeyID:         17,
				ZenBTCChangeAddressKeyIDs: []uint64{18},
				ZenBTCUnstakerKeyID:       19,
				ZenBTCRewardsDepositKeyID: 20,
				BitcoinProxyCreatorID:     "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				StakeableAssets: []*types.AssetData{
					{Asset: types.Asset_ROCK, Precision: 6, PriceUSD: math.LegacyZeroDec()},
					{Asset: types.Asset_zenBTC, Precision: 8, PriceUSD: math.LegacyZeroDec()},
				},
			},
		},
	}

	chainID := ctx.ChainID()
	if chainID == "" {
		chainID = "zenrock"
	}

	for prefix, params := range paramsMap {
		if strings.HasPrefix(chainID, prefix) {
			currParams = params
			break
		}
	}

	if err := params.Set(ctx, currParams); err != nil {
		return err
	}

	return nil
}
