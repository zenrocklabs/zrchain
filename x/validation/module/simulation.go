package validation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	validationsimulation "github.com/Zenrock-Foundation/zrchain/v5/x/validation/simulation"

	"github.com/Zenrock-Foundation/zrchain/v5/testutil/sample"
	"github.com/Zenrock-Foundation/zrchain/v5/x/validation/types"
)

// avoid unused import issue
var (
	// _ = validationsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	// _ = simulation.MsgEntryKind
)

const (
// this line is used by starport scaffolding # simapp/module/const
)

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the staking module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	validationsimulation.RandomizedGenState(simState)
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return validationsimulation.ProposalMsgs()
}

// RegisterStoreDecoder registers a decoder for staking module's types
func (am AppModule) RegisterStoreDecoder(sdr simtypes.StoreDecoderRegistry) {
	sdr[types.StoreKey] = validationsimulation.NewDecodeStore(am.cdc)
}

// WeightedOperations returns the all the staking module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
