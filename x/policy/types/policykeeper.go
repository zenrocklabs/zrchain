package types

import (
	"github.com/Zenrock-Foundation/zrchain/v5/policy"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PolicyKeeper interface {
	PolicyForAction(ctx sdk.Context, act *Action) (policy.Policy, error)
	SetAction(ctx sdk.Context, action *Action) error
	ActionHandler(actionType string) (func(sdk.Context, *Action) (any, error), bool)
	RegisterActionHandler(actionType string, f func(sdk.Context, *Action) (any, error))
	GeneratorHandler(url string) (func(sdk.Context, *cdctypes.Any) (policy.Policy, error), bool)
	RegisterPolicyGeneratorHandler(reqType string, f func(sdk.Context, *cdctypes.Any) (policy.Policy, error))
	Unpack(policyPb *Policy) (policy.Policy, error)
	Codec() codec.BinaryCodec
}
