package treasury

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"

	treasurysimulation "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/simulation"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
)

// avoid unused import issue
var (
	_ = treasurysimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgNewKeyRequest = "op_weight_msg_new_key_request"
	// TODO: Determine the simulation weight value
	defaultWeightMsgNewKeyRequest int = 100

	opWeightMsgFulfilKeyRequest = "op_weight_msg_fulfil_key_request"
	// TODO: Determine the simulation weight value
	defaultWeightMsgFulfilKeyRequest int = 100

	opWeightMsgNewSignatureRequest = "op_weight_msg_new_signature_request"
	// TODO: Determine the simulation weight value
	defaultWeightMsgNewSignatureRequest int = 100

	opWeightMsgFulfilSignatureRequest = "op_weight_msg_fulfil_signature_request"
	// TODO: Determine the simulation weight value
	defaultWeightMsgFulfilSignatureRequest int = 100

	opWeightMsgNewSignTransactionRequest = "op_weight_msg_new_sign_transaction_request"
	// TODO: Determine the simulation weight value
	defaultWeightMsgNewSignTransactionRequest int = 100

	opWeightMsgTransferFromKeyring = "op_weight_msg_transfer_from_keyring"
	// TODO: Determine the simulation weight value
	defaultWeightMsgTransferFromKeyring int = 100

	opWeightMsgNewICATransactionRequest = "op_weight_msg_new_ica_transaction_request"
	// TODO: Determine the simulation weight value
	defaultWeightMsgNewICATransactionRequest int = 100

	opWeightMsgFulfilICATransactionRequest = "op_weight_msg_fulfil_ica_transaction_request"
	// TODO: Determine the simulation weight value
	defaultWeightMsgFulfilICATransactionRequest int = 100

	opWeightMsgNewZrSignSignatureRequest = "op_weight_msg_new_zr_sign_signature_request"
	// TODO: Determine the simulation weight value
	defaultWeightMsgNewZrSignSignatureRequest int = 100

	opWeightMsgUpdateKeyPolicy = "op_weight_msg_update_key_policy"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateKeyPolicy int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	treasuryGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&treasuryGenesis)
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

	var weightMsgNewKeyRequest int
	simState.AppParams.GetOrGenerate(opWeightMsgNewKeyRequest, &weightMsgNewKeyRequest, nil,
		func(_ *rand.Rand) {
			weightMsgNewKeyRequest = defaultWeightMsgNewKeyRequest
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgNewKeyRequest,
		treasurysimulation.SimulateMsgNewKeyRequest(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgFulfilKeyRequest int
	simState.AppParams.GetOrGenerate(opWeightMsgFulfilKeyRequest, &weightMsgFulfilKeyRequest, nil,
		func(_ *rand.Rand) {
			weightMsgFulfilKeyRequest = defaultWeightMsgFulfilKeyRequest
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgFulfilKeyRequest,
		treasurysimulation.SimulateMsgFulfilKeyRequest(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgNewSignatureRequest int
	simState.AppParams.GetOrGenerate(opWeightMsgNewSignatureRequest, &weightMsgNewSignatureRequest, nil,
		func(_ *rand.Rand) {
			weightMsgNewSignatureRequest = defaultWeightMsgNewSignatureRequest
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgNewSignatureRequest,
		treasurysimulation.SimulateMsgNewSignatureRequest(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgFulfilSignatureRequest int
	simState.AppParams.GetOrGenerate(opWeightMsgFulfilSignatureRequest, &weightMsgFulfilSignatureRequest, nil,
		func(_ *rand.Rand) {
			weightMsgFulfilSignatureRequest = defaultWeightMsgFulfilSignatureRequest
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgFulfilSignatureRequest,
		treasurysimulation.SimulateMsgFulfilSignatureRequest(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgNewSignTransactionRequest int
	simState.AppParams.GetOrGenerate(opWeightMsgNewSignTransactionRequest, &weightMsgNewSignTransactionRequest, nil,
		func(_ *rand.Rand) {
			weightMsgNewSignTransactionRequest = defaultWeightMsgNewSignTransactionRequest
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgNewSignTransactionRequest,
		treasurysimulation.SimulateMsgNewSignTransactionRequest(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgTransferFromKeyring int
	simState.AppParams.GetOrGenerate(opWeightMsgTransferFromKeyring, &weightMsgTransferFromKeyring, nil,
		func(_ *rand.Rand) {
			weightMsgTransferFromKeyring = defaultWeightMsgTransferFromKeyring
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgTransferFromKeyring,
		treasurysimulation.SimulateMsgTransferFromKeyring(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgNewICATransactionRequest int
	simState.AppParams.GetOrGenerate(opWeightMsgNewICATransactionRequest, &weightMsgNewICATransactionRequest, nil,
		func(_ *rand.Rand) {
			weightMsgNewICATransactionRequest = defaultWeightMsgNewICATransactionRequest
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgNewICATransactionRequest,
		treasurysimulation.SimulateMsgNewICATransactionRequest(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgFulfilICATransactionRequest int
	simState.AppParams.GetOrGenerate(opWeightMsgFulfilICATransactionRequest, &weightMsgFulfilICATransactionRequest, nil,
		func(_ *rand.Rand) {
			weightMsgFulfilICATransactionRequest = defaultWeightMsgFulfilICATransactionRequest
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgFulfilICATransactionRequest,
		treasurysimulation.SimulateMsgFulfilICATransactionRequest(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgNewZrSignSignatureRequest int
	simState.AppParams.GetOrGenerate(opWeightMsgNewZrSignSignatureRequest, &weightMsgNewZrSignSignatureRequest, nil,
		func(_ *rand.Rand) {
			weightMsgNewZrSignSignatureRequest = defaultWeightMsgNewZrSignSignatureRequest
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgNewZrSignSignatureRequest,
		treasurysimulation.SimulateMsgNewZrSignSignatureRequest(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateKeyPolicy int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateKeyPolicy, &weightMsgUpdateKeyPolicy, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateKeyPolicy = defaultWeightMsgUpdateKeyPolicy
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateKeyPolicy,
		treasurysimulation.SimulateMsgUpdateKeyPolicy(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgNewKeyRequest,
			defaultWeightMsgNewKeyRequest,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				treasurysimulation.SimulateMsgNewKeyRequest(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgFulfilKeyRequest,
			defaultWeightMsgFulfilKeyRequest,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				treasurysimulation.SimulateMsgFulfilKeyRequest(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgNewSignatureRequest,
			defaultWeightMsgNewSignatureRequest,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				treasurysimulation.SimulateMsgNewSignatureRequest(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgFulfilSignatureRequest,
			defaultWeightMsgFulfilSignatureRequest,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				treasurysimulation.SimulateMsgFulfilSignatureRequest(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgNewSignTransactionRequest,
			defaultWeightMsgNewSignTransactionRequest,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				treasurysimulation.SimulateMsgNewSignTransactionRequest(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgTransferFromKeyring,
			defaultWeightMsgTransferFromKeyring,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				treasurysimulation.SimulateMsgTransferFromKeyring(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgNewICATransactionRequest,
			defaultWeightMsgNewICATransactionRequest,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				treasurysimulation.SimulateMsgNewICATransactionRequest(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgFulfilICATransactionRequest,
			defaultWeightMsgFulfilICATransactionRequest,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				treasurysimulation.SimulateMsgFulfilICATransactionRequest(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgNewZrSignSignatureRequest,
			defaultWeightMsgNewZrSignSignatureRequest,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				treasurysimulation.SimulateMsgNewZrSignSignatureRequest(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateKeyPolicy,
			defaultWeightMsgUpdateKeyPolicy,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				treasurysimulation.SimulateMsgUpdateKeyPolicy(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
