package keeper

import (
	"github.com/Zenrock-Foundation/zrchain/v4/x/identity/types"
)

var _ types.QueryServer = Keeper{}