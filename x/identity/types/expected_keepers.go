package types

import (
	"context"

	policy "github.com/Zenrock-Foundation/zrchain/v5/policy"
	policytypes "github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

type PolicyKeeper interface {
	AddAction(ctx sdk.Context, creator string, msg sdk.Msg, policyID, btl uint64, policyData map[string][]byte, wsOwners []string) (*policytypes.Action, error)
	GetPolicyParticipants(ctx context.Context, policyId uint64) (map[string]struct{}, error)
	PolicyMembersAreOwners(ctx context.Context, policyId uint64, wsOwners []string) error
	GetPolicy(ctx sdk.Context, policyId uint64) (*policytypes.Policy, error)
	SetAction(ctx sdk.Context, action *policytypes.Action) error
	PolicyForAction(ctx sdk.Context, act *policytypes.Action) (policy.Policy, error)
	ActionHandler(actionType string) (func(sdk.Context, *policytypes.Action) (any, error), bool)
	GeneratorHandler(reqType string) (func(sdk.Context, *cdctypes.Any) (policy.Policy, error), bool)
	Unpack(policyPb *policytypes.Policy) (policy.Policy, error)
	Codec() codec.BinaryCodec
	RegisterActionHandler(actionType string, f func(sdk.Context, *policytypes.Action) (any, error))
	RegisterPolicyGeneratorHandler(reqType string, f func(sdk.Context, *cdctypes.Any) (policy.Policy, error))
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
