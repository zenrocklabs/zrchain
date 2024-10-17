package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AdditionalSignature interface {
	Verify(ctx sdk.Context, config SignMethod, act Action) string
	GetConfigId() string
}
