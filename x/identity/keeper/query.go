package keeper

import (
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
)

var _ types.QueryServer = Keeper{}
