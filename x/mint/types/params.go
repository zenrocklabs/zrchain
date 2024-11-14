package types

import (
	"errors"
	"fmt"
	"strings"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParams returns Params instance with the given values.
func NewParams(mintDenom, protocolWalletAddress string, inflationRateChange, inflationMax, inflationMin, goalBonded, stakingYield, burnRate, protocolWalletRate, retentionRate, additionalStakingRewards, additionalMpcRewards, additionalBurnRate math.LegacyDec, blocksPerYear uint64) Params {
	return Params{
		MintDenom:                mintDenom,
		InflationRateChange:      inflationRateChange,
		InflationMax:             inflationMax,
		InflationMin:             inflationMin,
		GoalBonded:               goalBonded,
		BlocksPerYear:            blocksPerYear,
		StakingYield:             stakingYield,
		BurnRate:                 burnRate,
		ProtocolWalletRate:       protocolWalletRate,
		RetentionRate:            retentionRate,
		AdditionalStakingRewards: additionalStakingRewards,
		AdditionalMpcRewards:     additionalMpcRewards,
		AdditionalBurnRate:       additionalBurnRate,
		ProtocolWalletAddress:    protocolWalletAddress,
	}
}

// DefaultParams returns default x/mint module parameters.
func DefaultParams() Params {
	return Params{
		MintDenom:                sdk.DefaultBondDenom,
		InflationRateChange:      math.LegacyNewDecWithPrec(13, 2),
		InflationMax:             math.LegacyNewDecWithPrec(20, 2),
		InflationMin:             math.LegacyNewDecWithPrec(7, 2),
		GoalBonded:               math.LegacyNewDecWithPrec(67, 2),
		BlocksPerYear:            uint64(60 * 60 * 8766 / 5), // assuming 5 second block times
		StakingYield:             math.LegacyNewDecWithPrec(10, 2),
		BurnRate:                 math.LegacyNewDecWithPrec(0, 2),
		ProtocolWalletRate:       math.LegacyNewDecWithPrec(0, 2),
		RetentionRate:            math.LegacyNewDecWithPrec(0, 2),
		AdditionalStakingRewards: math.LegacyNewDecWithPrec(0, 2),
		AdditionalMpcRewards:     math.LegacyNewDecWithPrec(0, 2),
		AdditionalBurnRate:       math.LegacyNewDecWithPrec(0, 2),
		ProtocolWalletAddress:    "",
	}
}

// Validate does the sanity check on the params.
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateInflationRateChange(p.InflationRateChange); err != nil {
		return err
	}
	if err := validateInflationMax(p.InflationMax); err != nil {
		return err
	}
	if err := validateInflationMin(p.InflationMin); err != nil {
		return err
	}
	if err := validateGoalBonded(p.GoalBonded); err != nil {
		return err
	}
	if err := validateBlocksPerYear(p.BlocksPerYear); err != nil {
		return err
	}
	if err := validateStakingYield(p.StakingYield); err != nil {
		return err
	}
	if err := validateBurnRate(p.BurnRate); err != nil {
		return err
	}
	if err := validateProtocolWalletRate(p.ProtocolWalletRate); err != nil {
		return err
	}
	if err := validateRetentionRate(p.RetentionRate); err != nil {
		return err
	}
	if err := validateAdditionalStakingRewards(p.AdditionalStakingRewards); err != nil {
		return err
	}
	if err := validateAdditionalMpcRewards(p.AdditionalMpcRewards); err != nil {
		return err
	}
	if err := validateAdditionalBurnRate(p.AdditionalBurnRate); err != nil {
		return err
	}

	if p.InflationMax.LT(p.InflationMin) {
		return fmt.Errorf(
			"max inflation (%s) must be greater than or equal to min inflation (%s)",
			p.InflationMax, p.InflationMin,
		)
	}

	return nil
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateInflationRateChange(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("inflation rate change cannot be nil: %s", v)
	}
	if v.IsNegative() {
		return fmt.Errorf("inflation rate change cannot be negative: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("inflation rate change too large: %s", v)
	}

	return nil
}

func validateInflationMax(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("max inflation cannot be nil: %s", v)
	}
	if v.IsNegative() {
		return fmt.Errorf("max inflation cannot be negative: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("max inflation too large: %s", v)
	}

	return nil
}

func validateInflationMin(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("min inflation cannot be nil: %s", v)
	}
	if v.IsNegative() {
		return fmt.Errorf("min inflation cannot be negative: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("min inflation too large: %s", v)
	}

	return nil
}

func validateGoalBonded(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("goal bonded cannot be nil: %s", v)
	}
	if v.IsNegative() || v.IsZero() {
		return fmt.Errorf("goal bonded must be positive: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("goal bonded too large: %s", v)
	}

	return nil
}

func validateBlocksPerYear(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("blocks per year must be positive: %d", v)
	}

	return nil
}

func validateStakingYield(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("staking yield must be positive: %s", v)
	}
	if v.IsNegative() || v.IsZero() {
		return fmt.Errorf("staking yield must be positive: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("staking yield too large: %s", v)
	}

	return nil
}

func validateBurnRate(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("burn rate must be positive: %s", v)
	}
	if v.IsNegative() || v.IsZero() {
		return fmt.Errorf("burn rate must be positive: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("burn rate too large: %s", v)
	}

	return nil
}

func validateProtocolWalletRate(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("protocol wallet rate must be positive: %s", v)
	}
	if v.IsNegative() || v.IsZero() {
		return fmt.Errorf("protocol wallet rate must be positive: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("protocol wallet rate too large: %s", v)
	}

	return nil
}

func validateRetentionRate(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("retention rate must be positive: %s", v)
	}
	if v.IsNegative() || v.IsZero() {
		return fmt.Errorf("retention rate must be positive: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("retention rate too large: %s", v)
	}

	return nil
}

func validateAdditionalStakingRewards(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("additional staking rewards must be positive: %s", v)
	}
	if v.IsNegative() || v.IsZero() {
		return fmt.Errorf("additional staking rewards must be positive: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("additional staking rewards too large: %s", v)
	}

	return nil
}

func validateAdditionalMpcRewards(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("additional mpc rewards must be positive: %s", v)
	}
	if v.IsNegative() || v.IsZero() {
		return fmt.Errorf("additional mpc rewards must be positive: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("additional mpc rewards too large: %s", v)
	}

	return nil
}

func validateAdditionalBurnRate(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("additional burn rate must be positive: %s", v)
	}
	if v.IsNegative() || v.IsZero() {
		return fmt.Errorf("additional burn rate must be positive: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("additional burn rate too large: %s", v)
	}

	return nil
}
