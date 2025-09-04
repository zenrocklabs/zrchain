package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
    

	"github.com/Zenrock-Foundation/zrchain/v6/x/bedrock/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

        // the address capable of executing a MsgUpdateParams message. Typically, this
        // should be the x/gov module account.
        authority string
        
		
        identityKeeper types.IdentityKeeper
        treasuryKeeper types.TreasuryKeeper
        validationKeeper types.ValidationKeeper
	}
)

func NewKeeper(
    cdc codec.BinaryCodec,
	storeService store.KVStoreService,
    logger log.Logger,
	authority string,
    
    identityKeeper types.IdentityKeeper,
    treasuryKeeper types.TreasuryKeeper,
    validationKeeper types.ValidationKeeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,
		
		identityKeeper: identityKeeper,
		treasuryKeeper: treasuryKeeper,
		validationKeeper: validationKeeper,
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


