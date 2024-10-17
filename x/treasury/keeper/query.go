package keeper

import (
	"github.com/Zenrock-Foundation/zrchain/v4/x/treasury/types"
)

var _ types.QueryServer = Keeper{}
