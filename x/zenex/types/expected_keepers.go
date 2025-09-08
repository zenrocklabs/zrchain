package types

import (
	"context"

	"cosmossdk.io/math"
	identitytypes "github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type IdentityKeeper interface {
	GetWorkspace(ctx sdk.Context, id string) (*identitytypes.Workspace, error)
}

type TreasuryKeeper interface {
	GetKey(ctx sdk.Context, keyID uint64) (*treasurytypes.Key, error)
	MakeSignTransactionRequest(goCtx context.Context, msg *treasurytypes.MsgNewSignTransactionRequest) (*treasurytypes.MsgNewSignTransactionRequestResponse, error)
}

type ValidationKeeper interface {
	GetRockBtcPrice(ctx context.Context) (math.LegacyDec, error)
	GetBtcRockPrice(ctx context.Context) (math.LegacyDec, error)
	GetAssetPrices(ctx context.Context) (map[validationtypes.Asset]math.LegacyDec, error)
}

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(context.Context, sdk.AccAddress, string, sdk.Coins) error
	SendCoinsFromModuleToAccount(context.Context, string, sdk.AccAddress, sdk.Coins) error
	// Methods imported from bank should be defined here
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
