package keeper

import (
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	policykeeper "github.com/Zenrock-Foundation/zrchain/v5/x/policy/keeper"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = (*msgServer)(nil)

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	s := &msgServer{Keeper: keeper}

	policykeeper.RegisterActionHandler(
		&keeper.policyKeeper,
		"/zrchain.identity.MsgAddWorkspaceOwner",
		s.AddOwnerActionHandler,
	)
	policykeeper.RegisterPolicyGeneratorHandler(
		&keeper.policyKeeper,
		"/zrchain.identity.MsgAddWorkspaceOwner",
		s.AddOwnerPolicyGenerator,
	)

	policykeeper.RegisterActionHandler(
		&keeper.policyKeeper,
		"/zrchain.identity.MsgRemoveWorkspaceOwner",
		s.RemoveOwnerActionHandler,
	)
	policykeeper.RegisterPolicyGeneratorHandler(
		&keeper.policyKeeper,
		"/zrchain.identity.MsgRemoveWorkspaceOwner",
		s.RemoveOwnerPolicyGenerator,
	)

	policykeeper.RegisterActionHandler(
		&keeper.policyKeeper,
		"/zrchain.identity.MsgAppendChildWorkspace",
		s.AppendChildWorkspaceActionHandler,
	)
	policykeeper.RegisterPolicyGeneratorHandler(
		&keeper.policyKeeper,
		"/zrchain.identity.MsgAppendChildWorkspace",
		s.AppendChildWorkspacePolicyGenerator,
	)

	policykeeper.RegisterActionHandler(
		&keeper.policyKeeper,
		"/zrchain.identity.MsgNewChildWorkspace",
		s.NewChildWorkspaceActionHandler,
	)
	policykeeper.RegisterPolicyGeneratorHandler(
		&keeper.policyKeeper,
		"/zrchain.identity.MsgNewChildWorkspace",
		s.NewChildWorkspacePolicyGenerator,
	)

	policykeeper.RegisterActionHandler(
		&keeper.policyKeeper,
		"/zrchain.identity.MsgUpdateWorkspace",
		s.UpdateWorkspaceActionHandler,
	)
	policykeeper.RegisterPolicyGeneratorHandler(
		&keeper.policyKeeper,
		"/zrchain.identity.MsgUpdateWorkspace",
		s.UpdateWorkspacePolicyGenerator,
	)

	return s
}
