package params

import (
	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// Name defines the application name of the Zenrock network.
	Name = "zen"

	// BondDenom defines the native staking token denomination.
	// NOTE: it is used by IBC, and must not change to avoid token migration in all IBC chains.
	BondDenom     = "urock"
	BaseDenomUnit = 10

	// DisplayDenom defines the name, symbol, and display value of the native token.
	DisplayDenom = "ROCK"

	// DefaultGasLimit - set to the same value as cosmos-sdk flags.DefaultGasLimit
	// this value is currently only used in tests.
	DefaultGasLimit = 200000
	MinGasFee       = "0.0001" + BondDenom

	AccountAddressPrefix = "zen"
)

var (
	ErrUnknownAddress = sdkerrors.Register("params", 9, "unknown address")
)

type ZRConfig struct {
	IsValidator bool
	SidecarAddr string
}

func init() {
	SetAddressPrefixes()
	RegisterDenoms()
}

// RegisterDenoms registers the base and display denominations to the SDK.
func RegisterDenoms() {
	if err := sdk.RegisterDenom(DisplayDenom, math.LegacyOneDec()); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(BondDenom, math.LegacyNewDecWithPrec(1, BaseDenomUnit)); err != nil {
		panic(err)
	}
}

func SetAddressPrefixes() {
	// Set prefixes
	accountPubKeyPrefix := AccountAddressPrefix + "pub"
	validatorAddressPrefix := AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := AccountAddressPrefix + "valconspub"

	// Set config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)

	config.SetAddressVerifier(AddressVerifier)
}

func AddressVerifier(bytes []byte) error {
	if len(bytes) == 0 {
		return sdkerrors.Wrap(ErrUnknownAddress, "addresses cannot be empty")
	}

	if len(bytes) > address.MaxAddrLen {
		return sdkerrors.Wrapf(ErrUnknownAddress, "address max length is %d, got %d", address.MaxAddrLen, len(bytes))
	}

	if len(bytes) != 20 && len(bytes) != 32 && len(bytes) != 11 {
		return sdkerrors.Wrapf(ErrUnknownAddress, "address length must be 11, 20, or 32 bytes, got %d", len(bytes))
	}

	return nil
}
