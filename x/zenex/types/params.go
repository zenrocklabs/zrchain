package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyBtcproxyaddress       = []byte("Btcproxyaddress")
	KeyMinimumSatoshis       = []byte("MinimumSatoshis")
	KeyZenexBtcPoolKeyId     = []byte("ZenexBtcPoolKeyId")
	KeyZenexWorkspaceAddress = []byte("ZenexWorkspaceAddress")
	KeySwapThresholdSatoshis = []byte("SwapThresholdSatoshis")
	// TODO: Determine the default value
	DefaultBtcproxyaddress       string = "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq"
	DefaultMinimumSatoshis       uint64 = 1000
	DefaultZenexBtcPoolKeyId     uint64 = 16
	DefaultZenexWorkspaceAddress string = "workspace14a2hpadpsy9h4auve2z8lw"
	DefaultSwapThresholdSatoshis uint64 = 6100
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	btcProxyAddress string,
	MinimumSatoshis uint64,
	ZenexBtcPoolKeyId uint64,
	ZenexWorkspaceAddress string,
	SwapThresholdSatoshis uint64,
) Params {
	return Params{
		BtcProxyAddress:       btcProxyAddress,
		MinimumSatoshis:       MinimumSatoshis,
		ZenexPoolKeyId:        ZenexBtcPoolKeyId,
		ZenexWorkspaceAddress: ZenexWorkspaceAddress,
		SwapThresholdSatoshis: SwapThresholdSatoshis,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultBtcproxyaddress,
		1000, // 1000 satoshis = 0.0001 BTC
		DefaultZenexBtcPoolKeyId,
		DefaultZenexWorkspaceAddress,
		DefaultSwapThresholdSatoshis, // 100000 satoshis = 0.001 BTC
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBtcproxyaddress, &p.BtcProxyAddress, validateBtcproxyaddress),
		paramtypes.NewParamSetPair(KeyMinimumSatoshis, &p.MinimumSatoshis, validateMinimumSatoshis),
		paramtypes.NewParamSetPair(KeyZenexBtcPoolKeyId, &p.ZenexPoolKeyId, validateZenexBtcPoolKeyId),
		paramtypes.NewParamSetPair(KeyZenexWorkspaceAddress, &p.ZenexWorkspaceAddress, validateZenexWorkspaceAddress),
		paramtypes.NewParamSetPair(KeySwapThresholdSatoshis, &p.SwapThresholdSatoshis, validateSwapThresholdSatoshis),
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
	if err := validateZenexBtcPoolKeyId(p.ZenexPoolKeyId); err != nil {
		return err
	}
	if err := validateZenexWorkspaceAddress(p.ZenexWorkspaceAddress); err != nil {
		return err
	}
	if err := validateSwapThresholdSatoshis(p.SwapThresholdSatoshis); err != nil {
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

	if minimumSatoshis < 1000 {
		return fmt.Errorf("minimum satoshis cannot be smaller than 1000")
	}

	return nil
}

// validateBtcChangeAddressKeyId validates the BtcChangeAddressKeyId param
func validateZenexBtcPoolKeyId(v interface{}) error {
	_, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	return nil
}

// validateZenexWorkspaceAddress validates the ZenexWorkspaceAddress param
func validateZenexWorkspaceAddress(v interface{}) error {
	_, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	return nil
}

// validateSwapThresholdSatoshis validates the SwapThresholdSatoshis param
func validateSwapThresholdSatoshis(v interface{}) error {
	swapThresholdSatoshis, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	if swapThresholdSatoshis < 1000 {
		return fmt.Errorf("swap threshold satoshis cannot be smaller than 1000")
	}

	return nil
}
