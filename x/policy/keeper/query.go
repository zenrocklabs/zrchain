package keeper

import (
	"github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
)

var _ types.QueryServer = Keeper{}
