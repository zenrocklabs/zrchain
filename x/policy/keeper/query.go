package keeper

import (
	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
)

var _ types.QueryServer = Keeper{}
