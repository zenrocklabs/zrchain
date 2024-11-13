package keeper

import (
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
)

var _ types.QueryServer = Keeper{}
