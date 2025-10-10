package zenbtc

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"

	"github.com/zenrocklabs/zenbtc/x/zenbtc/keeper"
	zenbtcsimulation "github.com/zenrocklabs/zenbtc/x/zenbtc/simulation"

	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

// avoid unused import issue
var (
	_ = zenbtcsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgVerifyDepositBlockInclusion = "op_weight_msg_verify_block_inclusion"
	// TODO: Determine the simulation weight value
	defaultWeightMsgVerifyDepositBlockInclusion int = 100

	opWeightMsgSubmitUnlockTransaction = "op_weight_msg_submit_solana_unlock_transaction"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSubmitUnlockTransaction int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	zenbtcGenesis := types.GenesisState{
		Params: *keeper.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&zenbtcGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgVerifyDepositBlockInclusion int
	simState.AppParams.GetOrGenerate(opWeightMsgVerifyDepositBlockInclusion, &weightMsgVerifyDepositBlockInclusion, nil,
		func(_ *rand.Rand) {
			weightMsgVerifyDepositBlockInclusion = defaultWeightMsgVerifyDepositBlockInclusion
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgVerifyDepositBlockInclusion,
		zenbtcsimulation.SimulateMsgVerifyDepositBlockInclusion(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgVerifyDepositBlockInclusion,
			defaultWeightMsgVerifyDepositBlockInclusion,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				zenbtcsimulation.SimulateMsgVerifyDepositBlockInclusion(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
