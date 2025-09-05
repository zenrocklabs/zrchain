package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyBtcproxyaddress = []byte("Btcproxyaddress")
	KeyMinimumSatoshis = []byte("MinimumSatoshis")
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
	MinimumSatoshis uint64,
) Params {
	return Params{
		BtcProxyAddress: btcProxyAddress,
		MinimumSatoshis: MinimumSatoshis,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultBtcproxyaddress,
		5000, // 5000 satoshis = 0.00005 BTC
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBtcproxyaddress, &p.BtcProxyAddress, validateBtcproxyaddress),
		paramtypes.NewParamSetPair(KeyMinimumSatoshis, &p.MinimumSatoshis, validateMinimumSatoshis),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateBtcproxyaddress(p.BtcProxyAddress); err != nil {
		return err
	}
	if err := validateMinimumSatoshis(p.MinimumSatoshis); err != nil {
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

// validateMinimumSatoshis validates the minimum satoshis param
func validateMinimumSatoshis(v interface{}) error {
	minimumSatoshis, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	if minimumSatoshis < 5000 {
		return fmt.Errorf("minimum satoshis cannot be smaller than 5000")
	}

	return nil
}
