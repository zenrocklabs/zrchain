package keeper

import (
	"github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
)

var _ types.QueryServer = Keeper{}
