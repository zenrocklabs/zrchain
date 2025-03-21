package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"cosmossdk.io/math"

	"github.com/Zenrock-Foundation/zrchain/v6/x/mint/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// Simulation parameter constants
const (
	Inflation                = "inflation"
	InflationRateChange      = "inflation_rate_change"
	InflationMax             = "inflation_max"
	InflationMin             = "inflation_min"
	GoalBonded               = "goal_bonded"
	StakingYield             = "staking_yield"
	AdditionalStakingRewards = "additional_staking_rewards"
	AdditionalMpcRewards     = "additional_mpc_rewards"
	AdditionalBurnRate       = "additional_burn_rate"
	ProtocolWalletRate       = "protocol_wallet_rate"
	ProtocolWalletAddress    = "protocol_wallet_address"
	BurnRate                 = "burn_rate"
)

// GenInflation randomized Inflation
func GenInflation(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(int64(r.Intn(99)), 2)
}

// GenInflationRateChange randomized InflationRateChange
func GenInflationRateChange(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(int64(r.Intn(99)), 2)
}

// GenInflationMax randomized InflationMax
func GenInflationMax(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(20, 2)
}

// GenInflationMin randomized InflationMin
func GenInflationMin(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(7, 2)
}

// GenGoalBonded randomized GoalBonded
func GenGoalBonded(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(67, 2)
}

// GenStakingYield randomized StakingYield
func GenStakingYield(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(10, 2)
}

// GenAdditionalStakingRewards randomized AdditionalStakingRewards
func GenAdditionalStakingRewards(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(10, 2)
}

// GenAdditionalMpcRewards randomized AdditionalMpcRewards
func GenAdditionalMpcRewards(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(10, 2)
}

// GenAdditionalBurnRate randomized AdditionalBurnRate
func GenAdditionalBurnRate(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(10, 2)
}

// GenProtocolWalletRate randomized ProtocolWalletRate
func GenProtocolWalletRate(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(10, 2)
}

// GenProtocolWalletAddress randomized ProtocolWalletAddress
func GenProtocolWalletAddress(r *rand.Rand) string {
	return "zen1qwnafe2s9eawhah5x6v4593v3tljdntl9zcqpn"
}

// GenBurnRate randomized BurnRate
func GenBurnRate(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(10, 2)
}

// RandomizedGenState generates a random GenesisState for mint
func RandomizedGenState(simState *module.SimulationState) {
	// minter
	var inflation math.LegacyDec
	simState.AppParams.GetOrGenerate(Inflation, &inflation, simState.Rand, func(r *rand.Rand) { inflation = GenInflation(r) })

	// params
	var inflationRateChange math.LegacyDec
	simState.AppParams.GetOrGenerate(InflationRateChange, &inflationRateChange, simState.Rand, func(r *rand.Rand) { inflationRateChange = GenInflationRateChange(r) })

	var inflationMax math.LegacyDec
	simState.AppParams.GetOrGenerate(InflationMax, &inflationMax, simState.Rand, func(r *rand.Rand) { inflationMax = GenInflationMax(r) })

	var inflationMin math.LegacyDec
	simState.AppParams.GetOrGenerate(InflationMin, &inflationMin, simState.Rand, func(r *rand.Rand) { inflationMin = GenInflationMin(r) })

	var goalBonded math.LegacyDec
	simState.AppParams.GetOrGenerate(GoalBonded, &goalBonded, simState.Rand, func(r *rand.Rand) { goalBonded = GenGoalBonded(r) })

	var stakingYield math.LegacyDec
	simState.AppParams.GetOrGenerate(StakingYield, &stakingYield, simState.Rand, func(r *rand.Rand) { stakingYield = GenStakingYield(r) })

	var additionalStakingRewards math.LegacyDec
	simState.AppParams.GetOrGenerate(AdditionalStakingRewards, &additionalStakingRewards, simState.Rand, func(r *rand.Rand) { additionalStakingRewards = GenAdditionalStakingRewards(r) })

	var additionalMpcRewards math.LegacyDec
	simState.AppParams.GetOrGenerate(AdditionalMpcRewards, &additionalMpcRewards, simState.Rand, func(r *rand.Rand) { additionalMpcRewards = GenAdditionalMpcRewards(r) })

	var additionalBurnRate math.LegacyDec
	simState.AppParams.GetOrGenerate(AdditionalBurnRate, &additionalBurnRate, simState.Rand, func(r *rand.Rand) { additionalBurnRate = GenAdditionalBurnRate(r) })

	var protocolWalletRate math.LegacyDec
	simState.AppParams.GetOrGenerate(ProtocolWalletRate, &protocolWalletRate, simState.Rand, func(r *rand.Rand) { protocolWalletRate = GenProtocolWalletRate(r) })

	var protocolWalletAddress string
	simState.AppParams.GetOrGenerate(ProtocolWalletAddress, &protocolWalletAddress, simState.Rand, func(r *rand.Rand) { protocolWalletAddress = GenProtocolWalletAddress(r) })

	var burnRate math.LegacyDec
	simState.AppParams.GetOrGenerate(BurnRate, &burnRate, simState.Rand, func(r *rand.Rand) { burnRate = GenBurnRate(r) })

	mintDenom := simState.BondDenom
	blocksPerYear := uint64(60 * 60 * 8766 / 5)
	params := types.NewParams(
		mintDenom,
		protocolWalletAddress,
		inflationRateChange,
		inflationMax,
		inflationMin,
		goalBonded,
		stakingYield,
		additionalStakingRewards,
		additionalMpcRewards,
		additionalBurnRate,
		protocolWalletRate,
		burnRate,
		blocksPerYear,
	)

	mintGenesis := types.NewGenesisState(types.InitialMinter(inflation), params)

	bz, err := json.MarshalIndent(&mintGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated minting parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(mintGenesis)
}
