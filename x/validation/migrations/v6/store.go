package v6

import (
	"strings"

	"cosmossdk.io/collections"
	"github.com/Zenrock-Foundation/zrchain/v5/x/validation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func UpdateParams(ctx sdk.Context, params collections.Item[types.HVParams]) error {
	oldParams, err := params.Get(ctx)
	if err != nil {
		return err
	}

	currParams := oldParams

	keyIDMap := map[string]uint64{
		"diamond": 1,
		"gardia":  2,
		"amber":   3,
	}

	for prefix, keyID := range keyIDMap {
		if strings.HasPrefix(ctx.ChainID(), prefix) {
			currParams.ZenBTCParams.ZenBTCMinterKeyID = keyID
			break
		}
	}

	return nil
}
