package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/cosmos/gogoproto/proto"
)

type SignMethod interface {
	proto.Message

	VerifyConfig(ctx sdk.Context) error
	GetConfigId() string
	GetParticipantId() string
	IsActive() bool
	SetActive(active bool)
}
