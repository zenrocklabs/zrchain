package identity

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"
	identitysimulation "github.com/Zenrock-Foundation/zrchain/v6/x/identity/simulation"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
)

// avoid unused import issue
var (
	_ = identitysimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgNewWorkspace = "op_weight_msg_new_workspace"
	// TODO: Determine the simulation weight value
	defaultWeightMsgNewWorkspace int = 100

	opWeightMsgAddWorkspaceOwner = "op_weight_msg_add_workspace_owner"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddWorkspaceOwner int = 100

	opWeightMsgAppendChildWorkspace = "op_weight_msg_append_child_workspace"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAppendChildWorkspace int = 100

	opWeightMsgNewChildWorkspace = "op_weight_msg_new_child_workspace"
	// TODO: Determine the simulation weight value
	defaultWeightMsgNewChildWorkspace int = 100

	opWeightMsgRemoveWorkspaceOwner = "op_weight_msg_remove_workspace_owner"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRemoveWorkspaceOwner int = 100

	opWeightMsgAddKeyringParty = "op_weight_msg_add_keyring_party"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddKeyringParty int = 100

	opWeightMsgUpdateKeyring = "op_weight_msg_update_keyring"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateKeyring int = 100

	opWeightMsgRemoveKeyringParty = "op_weight_msg_remove_keyring_party"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRemoveKeyringParty int = 100

	opWeightMsgNewKeyring = "op_weight_msg_new_keyring"
	// TODO: Determine the simulation weight value
	defaultWeightMsgNewKeyring int = 100

	opWeightMsgAddKeyringAdmin = "op_weight_msg_add_keyring_admin"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddKeyringAdmin int = 100

	opWeightMsgRemoveKeyringAdmin = "op_weight_msg_remove_keyring_admin"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRemoveKeyringAdmin int = 100

	opWeightMsgUpdateWorkspace = "op_weight_msg_update_workspace"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateWorkspace int = 100

	opWeightMsgDeactivateKeyring = "op_weight_msg_deactivate_keyring"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeactivateKeyring int = 100

	opWeightMsgNewZrSignWorkspace = "op_weight_msg_new_zr_sign_workspace"
	// TODO: Determine the simulation weight value
	defaultWeightMsgNewZrSignWorkspace int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	identityGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&identityGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgNewWorkspace int
	simState.AppParams.GetOrGenerate(opWeightMsgNewWorkspace, &weightMsgNewWorkspace, nil,
		func(_ *rand.Rand) {
			weightMsgNewWorkspace = defaultWeightMsgNewWorkspace
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgNewWorkspace,
		identitysimulation.SimulateMsgNewWorkspace(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAddWorkspaceOwner int
	simState.AppParams.GetOrGenerate(opWeightMsgAddWorkspaceOwner, &weightMsgAddWorkspaceOwner, nil,
		func(_ *rand.Rand) {
			weightMsgAddWorkspaceOwner = defaultWeightMsgAddWorkspaceOwner
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddWorkspaceOwner,
		identitysimulation.SimulateMsgAddWorkspaceOwner(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAppendChildWorkspace int
	simState.AppParams.GetOrGenerate(opWeightMsgAppendChildWorkspace, &weightMsgAppendChildWorkspace, nil,
		func(_ *rand.Rand) {
			weightMsgAppendChildWorkspace = defaultWeightMsgAppendChildWorkspace
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAppendChildWorkspace,
		identitysimulation.SimulateMsgAppendChildWorkspace(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgNewChildWorkspace int
	simState.AppParams.GetOrGenerate(opWeightMsgNewChildWorkspace, &weightMsgNewChildWorkspace, nil,
		func(_ *rand.Rand) {
			weightMsgNewChildWorkspace = defaultWeightMsgNewChildWorkspace
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgNewChildWorkspace,
		identitysimulation.SimulateMsgNewChildWorkspace(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgRemoveWorkspaceOwner int
	simState.AppParams.GetOrGenerate(opWeightMsgRemoveWorkspaceOwner, &weightMsgRemoveWorkspaceOwner, nil,
		func(_ *rand.Rand) {
			weightMsgRemoveWorkspaceOwner = defaultWeightMsgRemoveWorkspaceOwner
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRemoveWorkspaceOwner,
		identitysimulation.SimulateMsgRemoveWorkspaceOwner(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAddKeyringParty int
	simState.AppParams.GetOrGenerate(opWeightMsgAddKeyringParty, &weightMsgAddKeyringParty, nil,
		func(_ *rand.Rand) {
			weightMsgAddKeyringParty = defaultWeightMsgAddKeyringParty
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddKeyringParty,
		identitysimulation.SimulateMsgAddKeyringParty(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateKeyring int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateKeyring, &weightMsgUpdateKeyring, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateKeyring = defaultWeightMsgUpdateKeyring
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateKeyring,
		identitysimulation.SimulateMsgUpdateKeyring(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgRemoveKeyringParty int
	simState.AppParams.GetOrGenerate(opWeightMsgRemoveKeyringParty, &weightMsgRemoveKeyringParty, nil,
		func(_ *rand.Rand) {
			weightMsgRemoveKeyringParty = defaultWeightMsgRemoveKeyringParty
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRemoveKeyringParty,
		identitysimulation.SimulateMsgRemoveKeyringParty(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgNewKeyring int
	simState.AppParams.GetOrGenerate(opWeightMsgNewKeyring, &weightMsgNewKeyring, nil,
		func(_ *rand.Rand) {
			weightMsgNewKeyring = defaultWeightMsgNewKeyring
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgNewKeyring,
		identitysimulation.SimulateMsgNewKeyring(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAddKeyringAdmin int
	simState.AppParams.GetOrGenerate(opWeightMsgAddKeyringAdmin, &weightMsgAddKeyringAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgAddKeyringAdmin = defaultWeightMsgAddKeyringAdmin
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddKeyringAdmin,
		identitysimulation.SimulateMsgAddKeyringAdmin(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgRemoveKeyringAdmin int
	simState.AppParams.GetOrGenerate(opWeightMsgRemoveKeyringAdmin, &weightMsgRemoveKeyringAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgRemoveKeyringAdmin = defaultWeightMsgRemoveKeyringAdmin
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRemoveKeyringAdmin,
		identitysimulation.SimulateMsgRemoveKeyringAdmin(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateWorkspace int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateWorkspace, &weightMsgUpdateWorkspace, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateWorkspace = defaultWeightMsgUpdateWorkspace
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateWorkspace,
		identitysimulation.SimulateMsgUpdateWorkspace(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeactivateKeyring int
	simState.AppParams.GetOrGenerate(opWeightMsgDeactivateKeyring, &weightMsgDeactivateKeyring, nil,
		func(_ *rand.Rand) {
			weightMsgDeactivateKeyring = defaultWeightMsgDeactivateKeyring
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeactivateKeyring,
		identitysimulation.SimulateMsgDeactivateKeyring(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgNewZrSignWorkspace int
	simState.AppParams.GetOrGenerate(opWeightMsgNewZrSignWorkspace, &weightMsgNewZrSignWorkspace, nil,
		func(_ *rand.Rand) {
			weightMsgNewZrSignWorkspace = defaultWeightMsgNewZrSignWorkspace
		},
	)

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgNewWorkspace,
			defaultWeightMsgNewWorkspace,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				identitysimulation.SimulateMsgNewWorkspace(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddWorkspaceOwner,
			defaultWeightMsgAddWorkspaceOwner,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				identitysimulation.SimulateMsgAddWorkspaceOwner(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgAppendChildWorkspace,
			defaultWeightMsgAppendChildWorkspace,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				identitysimulation.SimulateMsgAppendChildWorkspace(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgNewChildWorkspace,
			defaultWeightMsgNewChildWorkspace,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				identitysimulation.SimulateMsgNewChildWorkspace(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgRemoveWorkspaceOwner,
			defaultWeightMsgRemoveWorkspaceOwner,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				identitysimulation.SimulateMsgRemoveWorkspaceOwner(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// simulation.NewWeightedProposalMsg(
		// 	opWeightMsgNewKeyring,
		// 	defaultWeightMsgNewKeyring,
		// 	func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		// 		identitysimulation.SimulateMsgNewKeyring(am.accountKeeper, am.bankKeeper, am.keeper)
		// 		return nil
		// 	},
		// ),
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddKeyringParty,
			defaultWeightMsgAddKeyringParty,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				identitysimulation.SimulateMsgAddKeyringParty(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateKeyring,
			defaultWeightMsgUpdateKeyring,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				identitysimulation.SimulateMsgUpdateKeyring(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgRemoveKeyringParty,
			defaultWeightMsgRemoveKeyringParty,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				identitysimulation.SimulateMsgRemoveKeyringParty(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgNewKeyring,
			defaultWeightMsgNewKeyring,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				identitysimulation.SimulateMsgNewKeyring(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddKeyringAdmin,
			defaultWeightMsgAddKeyringAdmin,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				identitysimulation.SimulateMsgAddKeyringAdmin(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgRemoveKeyringAdmin,
			defaultWeightMsgRemoveKeyringAdmin,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				identitysimulation.SimulateMsgRemoveKeyringAdmin(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateWorkspace,
			defaultWeightMsgUpdateWorkspace,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				identitysimulation.SimulateMsgUpdateWorkspace(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateWorkspace,
			defaultWeightMsgUpdateWorkspace,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				identitysimulation.SimulateMsgUpdateWorkspace(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateWorkspace,
			defaultWeightMsgUpdateWorkspace,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				identitysimulation.SimulateMsgUpdateWorkspace(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateWorkspace,
			defaultWeightMsgUpdateWorkspace,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				identitysimulation.SimulateMsgUpdateWorkspace(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeactivateKeyring,
			defaultWeightMsgDeactivateKeyring,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				identitysimulation.SimulateMsgDeactivateKeyring(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
