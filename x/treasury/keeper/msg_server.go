package keeper

import (
	policy "github.com/Zenrock-Foundation/zrchain/v4/x/policy/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/treasury/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = (*msgServer)(nil)

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	s := &msgServer{Keeper: keeper}

	policy.RegisterActionHandler(
		&keeper.policyKeeper,
		"/zrchain.treasury.MsgNewKeyRequest",
		s.NewKeyRequestActionHandler,
	)
	policy.RegisterPolicyGeneratorHandler(
		&keeper.policyKeeper,
		"/zrchain.treasury.MsgNewKeyRequest",
		s.NewKeyRequestPolicyGenerator,
	)

	policy.RegisterActionHandler(
		&keeper.policyKeeper,
		"/zrchain.treasury.MsgNewSignatureRequest",
		s.NewSignatureRequestActionHandler,
	)
	policy.RegisterPolicyGeneratorHandler(
		&keeper.policyKeeper,
		"/zrchain.treasury.MsgNewSignatureRequest",
		s.NewSignatureRequestPolicyGenerator,
	)

	policy.RegisterActionHandler(
		&keeper.policyKeeper,
		"/zrchain.treasury.MsgNewSignTransactionRequest",
		s.NewSignTransactionRequestActionHandler,
	)
	policy.RegisterPolicyGeneratorHandler(
		&keeper.policyKeeper,
		"/zrchain.treasury.MsgNewSignTransactionRequest",
		s.NewSignTransactionRequestPolicyGenerator,
	)

	policy.RegisterActionHandler(
		&keeper.policyKeeper,
		"/zrchain.treasury.MsgUpdateKeyPolicy",
		s.UpdateKeyPolicyActionHandler,
	)
	policy.RegisterPolicyGeneratorHandler(
		&keeper.policyKeeper,
		"/zrchain.treasury.MsgUpdateKeyPolicy",
		s.UpdateKeyPolicyPolicyGenerator,
	)

	return s
}
