package types

// DONTCOVER

import sdkerrors "cosmossdk.io/errors"

// x/dct module sentinel errors.
var (
	ErrInvalidSigner                   = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrUnknownAsset                    = sdkerrors.Register(ModuleName, 1101, "unknown asset")
	ErrMissingAssetConfig              = sdkerrors.Register(ModuleName, 1102, "asset configuration missing")
	ErrMissingSolanaData               = sdkerrors.Register(ModuleName, 1103, "solana configuration missing for asset")
	ErrDuplicatePendingMintTransaction = sdkerrors.Register(ModuleName, 1104, "duplicate pending mint transaction")
	ErrDuplicateRedemption             = sdkerrors.Register(ModuleName, 1105, "duplicate redemption")
)
