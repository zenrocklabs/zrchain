package zentp

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/Zenrock-Foundation/zrchain/v5/testutil/sample"
	zentpsimulation "github.com/Zenrock-Foundation/zrchain/v5/x/zentp/simulation"
	"github.com/Zenrock-Foundation/zrchain/v5/x/zentp/types"
)

// avoid unused import issue
var (
	_ = zentpsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgBurn = "op_weight_msg_burn"
	// TODO: Determine the simulation weight value
	defaultWeightMsgBurn int = 100

	opWeightMsgBurnRock = "op_weight_msg_burn_rock"
	// TODO: Determine the simulation weight value
	defaultWeightMsgBurnRock int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	zentpGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&zentpGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgBurn int
	simState.AppParams.GetOrGenerate(opWeightMsgBurn, &weightMsgBurn, nil,
		func(_ *rand.Rand) {
			weightMsgBurn = defaultWeightMsgBurn
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgBurn,
		zentpsimulation.SimulateMsgBurn(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgBurnRock int
	simState.AppParams.GetOrGenerate(opWeightMsgBurnRock, &weightMsgBurnRock, nil,
		func(_ *rand.Rand) {
			weightMsgBurnRock = defaultWeightMsgBurnRock
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgBurnRock,
		zentpsimulation.SimulateMsgBurnRock(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgBurn,
			defaultWeightMsgBurn,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				zentpsimulation.SimulateMsgBurn(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgBurnRock,
			defaultWeightMsgBurnRock,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				zentpsimulation.SimulateMsgBurnRock(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
