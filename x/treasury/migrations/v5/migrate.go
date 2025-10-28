package v5

import (
	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
)

func UpdateZenbtcKeys(ctx sdk.Context, keyCol collections.Map[uint64, types.Key], codec codec.BinaryCodec) error {
	ctx.Logger().With("module", types.ModuleName).Info("updating zenbtc keys")

	err := keyCol.Walk(ctx, nil, func(key uint64, value types.Key) (stop bool, err error) {
		if value.ZenbtcMetadata != nil && value.ZenbtcMetadata.Asset.String() != dcttypes.Asset_ASSET_ZENZEC.String() {
			key, err := keyCol.Get(ctx, value.Id)
			if err != nil {
				return false, err
			}
			key.ZenbtcMetadata.Asset = dcttypes.Asset_ASSET_ZENBTC
			err = keyCol.Set(ctx, value.Id, key)
			if err != nil {
				return false, err
			}
			return false, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}
