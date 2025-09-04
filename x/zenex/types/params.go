package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)


var (
	KeyBtcproxyaddress = []byte("Btcproxyaddress")
	// TODO: Determine the default value
	DefaultBtcproxyaddress string = "btcproxyaddress"
)


// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	btcproxyaddress string,
) Params {
	return Params{
        Btcproxyaddress: btcproxyaddress,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
        DefaultBtcproxyaddress,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBtcproxyaddress, &p.Btcproxyaddress, validateBtcproxyaddress),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
   	if err := validateBtcproxyaddress(p.Btcproxyaddress); err != nil {
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
