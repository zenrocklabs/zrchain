package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/zenex module sentinel errors
var (
	ErrInvalidSigner     = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample            = sdkerrors.Register(ModuleName, 1101, "sample error")
	ErrMinimumSatoshis   = sdkerrors.Register(ModuleName, 1102, "amount in is less than the minimum satoshis")
	ErrInsufficientFunds = sdkerrors.Register(ModuleName, 1103, "insufficient funds")
	ErrWrongKeyType      = sdkerrors.Register(ModuleName, 1104, "sender key is not an ECDSA SECP256K1 key")
	ErrInvalidWalletType = sdkerrors.Register(ModuleName, 1105, "invalid wallet type: %s")
)
