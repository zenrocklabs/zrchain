package types

import (
	"fmt"

	"cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyBtcproxyaddress  = []byte("Btcproxyaddress")
	KeyMinimumBtcAmount = []byte("MinimumBtcAmount")
	// TODO: Determine the default value
	DefaultBtcproxyaddress string = "btcproxyaddress"
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	btcProxyAddress string,
	minimumBtcAmount math.LegacyDec,
) Params {
	return Params{
		BtcProxyAddress:  btcProxyAddress,
		MinimumBtcAmount: minimumBtcAmount,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultBtcproxyaddress,
		math.LegacyNewDecFromIntWithPrec(math.NewInt(5000), 8), // 5000 satoshis = 0.00005 BTC
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBtcproxyaddress, &p.BtcProxyAddress, validateBtcproxyaddress),
		paramtypes.NewParamSetPair(KeyMinimumBtcAmount, &p.MinimumBtcAmount, validateMinimumBtcAmount),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateBtcproxyaddress(p.BtcProxyAddress); err != nil {
		return err
	}
	if err := validateMinimumBtcAmount(p.MinimumBtcAmount); err != nil {
		return err
	}

	return nil
}

// validateBtcproxyaddress validates the Btcproxyaddress param
func validateBtcproxyaddress(v interface{}) error {
	btcproxyaddress, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// TODO implement validation
	_ = btcproxyaddress

	return nil
}

// validateMinimumBtcAmount validates the MinimumBtcAmount param
func validateMinimumBtcAmount(v interface{}) error {
	minimumBtcAmount, ok := v.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	if minimumBtcAmount.IsNegative() {
		return fmt.Errorf("minimum BTC amount cannot be negative")
	}

	return nil
}
