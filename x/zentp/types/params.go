package types

import (
	fmt "fmt"

	"cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var DefaultSolanaParams = &Solana{
	SignerKeyId:       10,
	ProgramId:         "AgoRvPWg2R7nkKhxvipvms79FmxQr75r2GwNSpPtxcLg",
	NonceAuthorityKey: 11,
	NonceAccountKey:   12,
	MintAddress:       "4oUDGAy46CmemmozTt6kWT5E3rqkLp2rCvAumpMWqR5T",
	FeeWallet:         "5aLz81F9uugwKBmvUY3DcXB1B7G2Yf7tB9zacdJBhZbh",
	Fee:               0,
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
		BridgeFee: math.LegacyNewDecWithPrec(5, 3),
		FlatFee:   200000000,
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
