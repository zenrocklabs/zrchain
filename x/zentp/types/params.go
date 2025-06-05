package types

import (
	fmt "fmt"

	"cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var DefaultSolanaParams = &Solana{
	SignerKeyId:       10,
	ProgramId:         "DXREJumiQhNejXa1b5EFPUxtSYdyJXBdiHeu6uX1ribA",
	NonceAuthorityKey: 11,
	NonceAccountKey:   12,
	MintAddress:       "StVNdHNSFK3uVTL5apWHysgze4M8zrsqwjEAH1JM87i",
	FeeWallet:         "FzqGcRG98v1KhKxatX2Abb2z1aJ2rViQwBK5GHByKCAd",
	Fee:               50000000,
	Btl:               20,
	// MultisigKeyAddress: "8cmZY2id22vxpXs2H3YYQNARuPHNuYwa7jipW1q1v9Fy",
}

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		Solana:    DefaultSolanaParams,
		BridgeFee: math.LegacyNewDecWithPrec(0, 0),
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateBridgeFee(p.BridgeFee); err != nil {
		return err
	}

	return nil
}

func validateBridgeFee(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("bridge fee must have a value: %s", v)
	}
	if v.IsNegative() || v.IsZero() {
		return fmt.Errorf("bridge fee must be non-negative: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("bridge fee too large: %s", v)
	}

	return nil
}
