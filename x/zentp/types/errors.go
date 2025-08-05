package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/zentp module sentinel errors
var (
	ErrInvalidSigner        = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample               = sdkerrors.Register(ModuleName, 1101, "sample error")
	ErrInvalidModuleAccount = sdkerrors.Register(ModuleName, 1102, "invalid module account")
	ErrInsufficientBalance  = sdkerrors.Register(ModuleName, 1103, "insufficient balance")
	ErrDctAlreadyExists     = sdkerrors.Register(ModuleName, 1104, "DCT already exists")
)
