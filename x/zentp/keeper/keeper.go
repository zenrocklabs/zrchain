package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	treasuryTypes "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v5/x/zentp/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority      string
		treasuryKeeper types.TreasuryKeeper
		bankKeeper     types.BankKeeper
		accountKeeper  types.AccountKeeper
		identityKeeper types.IdentityKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	treasuryKeeper types.TreasuryKeeper,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	identityKeeper types.IdentityKeeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}
	// ensure mint module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("the x/%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		cdc:            cdc,
		storeService:   storeService,
		authority:      authority,
		logger:         logger,
		treasuryKeeper: treasuryKeeper,
		bankKeeper:     bankKeeper,
		accountKeeper:  accountKeeper,
		identityKeeper: identityKeeper,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) UserOwnsKey(goCtx context.Context, user string, key *treasuryTypes.Key) bool {
	userWS, err := k.identityKeeper.GetWorkspaces(goCtx, user)
	if err != nil {
		k.Logger().Error("failed to get workspaces for user: "+user, err.Error())
		return false
	}

	for _, ws := range userWS {
		if key.WorkspaceAddr == ws.Address {
			return true
		}
	}

	return false
}
