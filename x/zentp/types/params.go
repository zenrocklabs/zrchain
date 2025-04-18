package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var DefaultSolanaParams = &Solana{
	SignerKeyId:       7,
	ProgramId:         "3jo4mdc6QbGRigia2jvmKShbmz3aWq4Y8bgUXfur5StT",
	NonceAuthorityKey: 8,
	NonceAccountKey:   9,
	MintAddress:       "9oBkgQUkq8jvzK98D7Uib6GYSZZmjnZ6QEGJRrAeKnDj",
	FeeWallet:         "FzqGcRG98v1KhKxatX2Abb2z1aJ2rViQwBK5GHByKCAd",
	Fee:               0,
	// MultisigKeyAddress: "8cmZY2id22vxpXs2H3YYQNARuPHNuYwa7jipW1q1v9Fy",
	Btl: 20,
}

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{Solana: DefaultSolanaParams}
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
	return nil
}
