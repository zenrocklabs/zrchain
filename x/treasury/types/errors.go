package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/treasury module sentinel errors
var (
	ErrInvalidSigner        = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample               = sdkerrors.Register(ModuleName, 1101, "sample error")
	ErrInvalidPacketTimeout = sdkerrors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion       = sdkerrors.Register(ModuleName, 1501, "invalid version")
	ErrNotFound             = sdkerrors.Register(ModuleName, 1503, "not found")
	ErrInternal             = sdkerrors.Register(ModuleName, 1504, "internal error")
	ErrInvalidArgument      = sdkerrors.Register(ModuleName, 1505, "invalid argument")
	ErrInvalidCommission    = sdkerrors.Register(ModuleName, 1506, "commission must be between 0 and 100")
)
