package types

import (
	"context"

	policykeeper "github.com/Zenrock-Foundation/zrchain/v5/x/policy/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

type PolicyKeeper interface {
	policykeeper.ExportedKeeper
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
