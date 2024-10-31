package policy

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/Zenrock-Foundation/zrchain/v5/testutil/sample"
	policysimulation "github.com/Zenrock-Foundation/zrchain/v5/x/policy/simulation"
	"github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
)

// avoid unused import issue
var (
	_ = policysimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgNewPolicy = "op_weight_msg_new_policy"
	// TODO: Determine the simulation weight value
	defaultWeightMsgNewPolicy int = 100

	opWeightMsgRevokeAction = "op_weight_msg_revoke_action"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRevokeAction int = 100

	opWeightMsgApproveAction = "op_weight_msg_approve_action"
	// TODO: Determine the simulation weight value
	defaultWeightMsgApproveAction int = 100

	opWeightMsgAddSignMethod = "op_weight_msg_add_sign_method"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddSignMethod int = 100

	opWeightMsgRemoveSignMethod = "op_weight_msg_remove_sign_method"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRemoveSignMethod int = 100

	opWeightMsgAddMultiGrant = "op_weight_msg_add_multi_grant"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddMultiGrant int = 100

	opWeightMsgRemoveMultiGrant = "op_weight_msg_remove_multi_grant"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRemoveMultiGrant int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	policyGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&policyGenesis)
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

	var weightMsgNewPolicy int
	simState.AppParams.GetOrGenerate(opWeightMsgNewPolicy, &weightMsgNewPolicy, nil,
		func(_ *rand.Rand) {
			weightMsgNewPolicy = defaultWeightMsgNewPolicy
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgNewPolicy,
		policysimulation.SimulateMsgNewPolicy(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgRevokeAction int
	simState.AppParams.GetOrGenerate(opWeightMsgRevokeAction, &weightMsgRevokeAction, nil,
		func(_ *rand.Rand) {
			weightMsgRevokeAction = defaultWeightMsgRevokeAction
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRevokeAction,
		policysimulation.SimulateMsgRevokeAction(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgApproveAction int
	simState.AppParams.GetOrGenerate(opWeightMsgApproveAction, &weightMsgApproveAction, nil,
		func(_ *rand.Rand) {
			weightMsgApproveAction = defaultWeightMsgApproveAction
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgApproveAction,
		policysimulation.SimulateMsgApproveAction(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAddSignMethod int
	simState.AppParams.GetOrGenerate(opWeightMsgAddSignMethod, &weightMsgAddSignMethod, nil,
		func(_ *rand.Rand) {
			weightMsgAddSignMethod = defaultWeightMsgAddSignMethod
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddSignMethod,
		policysimulation.SimulateMsgAddSignMethod(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgRemoveSignMethod int
	simState.AppParams.GetOrGenerate(opWeightMsgRemoveSignMethod, &weightMsgRemoveSignMethod, nil,
		func(_ *rand.Rand) {
			weightMsgRemoveSignMethod = defaultWeightMsgRemoveSignMethod
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRemoveSignMethod,
		policysimulation.SimulateMsgRemoveSignMethod(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAddMultiGrant int
	simState.AppParams.GetOrGenerate(opWeightMsgAddMultiGrant, &weightMsgAddMultiGrant, nil,
		func(_ *rand.Rand) {
			weightMsgAddMultiGrant = defaultWeightMsgAddMultiGrant
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddMultiGrant,
		policysimulation.SimulateMsgAddMultiGrant(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weighMsgRemoveMultiGrant int
	simState.AppParams.GetOrGenerate(opWeightMsgRemoveMultiGrant, &weighMsgRemoveMultiGrant, nil,
		func(_ *rand.Rand) {
			weighMsgRemoveMultiGrant = defaultWeightMsgRemoveMultiGrant
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weighMsgRemoveMultiGrant,
		policysimulation.SimulateMsgRemoveMultiGrant(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgNewPolicy,
			defaultWeightMsgNewPolicy,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				policysimulation.SimulateMsgNewPolicy(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgRevokeAction,
			defaultWeightMsgRevokeAction,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				policysimulation.SimulateMsgRevokeAction(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgApproveAction,
			defaultWeightMsgApproveAction,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				policysimulation.SimulateMsgApproveAction(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddSignMethod,
			defaultWeightMsgAddSignMethod,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				policysimulation.SimulateMsgAddSignMethod(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgRemoveSignMethod,
			defaultWeightMsgRemoveSignMethod,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				policysimulation.SimulateMsgRemoveSignMethod(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddMultiGrant,
			defaultWeightMsgAddMultiGrant,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				policysimulation.SimulateMsgAddMultiGrant(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgRemoveMultiGrant,
			defaultWeightMsgRemoveMultiGrant,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				policysimulation.SimulateMsgRemoveMultiGrant(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
