package v8

import (
	"cosmossdk.io/collections"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	types "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func UpdateBtcBlockHeaders(ctx sdk.Context, btcBlockHeaders collections.Map[int64, api.BTCBlockHeader], validationInfos collections.Map[int64, types.ValidationInfo]) error {
	btcBlockHeaders.Walk(ctx, nil, func(key int64, value api.BTCBlockHeader) (stop bool, err error) {
		value.BlockHeight = key
		return false, btcBlockHeaders.Set(ctx, key, value)
	})

	validationInfos.Walk(ctx, nil, func(key int64, value types.ValidationInfo) (stop bool, err error) {
		value.BlockHeight = uint64(key)
		return false, validationInfos.Set(ctx, key, value)
	})
	return nil
}
