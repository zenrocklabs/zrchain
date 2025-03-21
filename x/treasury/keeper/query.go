package keeper

import (
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
)

var _ types.QueryServer = Keeper{}
