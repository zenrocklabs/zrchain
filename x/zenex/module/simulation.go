package zenex

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"
	zenexsimulation "github.com/Zenrock-Foundation/zrchain/v6/x/zenex/simulation"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
)

// avoid unused import issue
var (
	_ = zenexsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgSwap = "op_weight_msg_swap"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSwap int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	zenexGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&zenexGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgSwap int
	simState.AppParams.GetOrGenerate(opWeightMsgSwap, &weightMsgSwap, nil,
		func(_ *rand.Rand) {
			weightMsgSwap = defaultWeightMsgSwap
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSwap,
		zenexsimulation.SimulateMsgSwap(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgSwap,
			defaultWeightMsgSwap,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				zenexsimulation.SimulateMsgSwap(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
