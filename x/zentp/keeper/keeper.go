package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	treasuryTypes "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/Zenrock-Foundation/zrchain/v5/x/zentp/types"
)

type (
	Keeper struct {
		cdc             codec.BinaryCodec
		storeService    store.KVStoreService
		memStoreService store.MemoryStoreService
		logger          log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority      string
		treasuryKeeper types.TreasuryKeeper
		bankKeeper     types.BankKeeper
		accountKeeper  types.AccountKeeper
		identityKeeper types.IdentityKeeper
		mintStore      collections.Map[uint64, types.BridgeRock]
		MintCount      collections.Item[uint64]
		burnStore      collections.Map[uint64, types.BridgeRock]
		BurnCount      collections.Item[uint64]
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
	memStoreService store.MemoryStoreService,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}
	// ensure mint module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("the x/%s module account has not been set", types.ModuleName))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:             cdc,
		storeService:    storeService,
		memStoreService: memStoreService,
		mintStore:       collections.NewMap(sb, types.MintsKey, types.MintsIndex, collections.Uint64Key, codec.CollValue[types.BridgeRock](cdc)),
		burnStore:       collections.NewMap(sb, types.BurnsKey, types.BurnsIndex, collections.Uint64Key, codec.CollValue[types.BridgeRock](cdc)),
		MintCount:       collections.NewItem(sb, types.MintCountKey, types.MintCountIndex, collections.Uint64Value),
		BurnCount:       collections.NewItem(sb, types.BurnCountKey, types.BurnCountIndex, collections.Uint64Value),
		authority:       authority,
		logger:          logger,
		treasuryKeeper:  treasuryKeeper,
		bankKeeper:      bankKeeper,
		accountKeeper:   accountKeeper,
		identityKeeper:  identityKeeper,
	}

	return k
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

func (k Keeper) GetMints(goCtx context.Context, address string, chainID string) ([]*types.BridgeRock, error) {
	mints, _, err := query.CollectionFilteredPaginate[uint64, types.BridgeRock, collections.Map[uint64, types.BridgeRock], *types.BridgeRock](
		goCtx,
		k.mintStore,
		nil,
		func(key uint64, value types.BridgeRock) (bool, error) {
			return value.SourceAddress == address &&
				value.DestinationChain == chainID, nil
		},
		func(key uint64, value types.BridgeRock) (*types.BridgeRock, error) {
			return &value, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return mints, nil
}

func (k Keeper) GetBurns(goCtx context.Context, address string, chainID string) ([]*types.BridgeRock, error) {
	burns, _, err := query.CollectionFilteredPaginate[uint64, types.BridgeRock, collections.Map[uint64, types.BridgeRock], *types.BridgeRock](
		goCtx,
		k.burnStore,
		nil,
		func(key uint64, value types.BridgeRock) (bool, error) {
			return value.RecipientAddress == address &&
				value.SourceChain == chainID, nil
		},
		func(key uint64, value types.BridgeRock) (*types.BridgeRock, error) {
			return &value, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return burns, nil
}

func (k Keeper) GetSignerKeyID(ctx context.Context) uint64 {
	params := k.GetParams(ctx)

	return params.ZrchainRelayerKeyId
}
