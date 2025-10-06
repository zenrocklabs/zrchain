package keeper

import (
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
)

var _ types.QueryServer = Keeper{}
