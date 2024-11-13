package keeper

import (
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
)

var _ types.QueryServer = Keeper{}
