package keeper

import (
	"fmt"
	"strings"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string

		identityKeeper   types.IdentityKeeper
		treasuryKeeper   types.TreasuryKeeper
		validationKeeper types.ValidationKeeper

		Schema     collections.Schema
		SwapsCount collections.Item[uint64]
		SwapsStore collections.Map[uint64, types.Swap]
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

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,

		identityKeeper:   identityKeeper,
		treasuryKeeper:   treasuryKeeper,
		validationKeeper: validationKeeper,

		SwapsCount: collections.NewItem(sb, types.SwapsCountKey, types.SwapsCountIndex, collections.Uint64Value),
		SwapsStore: collections.NewMap(sb, types.SwapsKey, types.SwapsIndex, collections.Uint64Key, codec.CollValue[types.Swap](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

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

func (k Keeper) GetPair(ctx sdk.Context, pair string) (*types.SwapPair, error) {
	var pairType types.SwapPair
	typeStr := strings.ToLower(pair)

	assetPrices, err := k.validationKeeper.GetAssetPrices(ctx)
	if err != nil {
		return nil, err
	}

	switch {
	case strings.Contains(typeStr, "rockbtc"):
		pairType = types.SwapPair{
			BaseToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_ROCK,
				Precision: 6,
				PriceUSD:  assetPrices[validationtypes.Asset_ROCK],
			},
			QuoteToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_BTC,
				Precision: 8,
				PriceUSD:  assetPrices[validationtypes.Asset_BTC],
			},
		}
	case strings.Contains(typeStr, "btcrock"):
		pairType = types.SwapPair{
			BaseToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_BTC,
				Precision: 8,
				PriceUSD:  assetPrices[validationtypes.Asset_BTC],
			},
			QuoteToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_ROCK,
				Precision: 6,
				PriceUSD:  assetPrices[validationtypes.Asset_ROCK],
			},
		}
	default:
		return nil, fmt.Errorf("unknown key type: %s", pair)
	}

	return &pairType, nil
}

func (k Keeper) GetSwaps(ctx sdk.Context) ([]types.Swap, error) {
	swapStore, err := k.SwapsStore.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	swaps, err := swapStore.Values()
	if err != nil {
		return nil, err
	}
	return swaps, nil
}
