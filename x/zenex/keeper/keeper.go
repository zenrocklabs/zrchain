package keeper

import (
	"fmt"
	"strings"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
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
		bankKeeper       types.BankKeeper

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
	bankKeeper types.BankKeeper,
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
		bankKeeper:       bankKeeper,

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

func (k Keeper) GetPair(ctx sdk.Context, pair string) (*types.SwapPair, math.LegacyDec, error) {
	var pairType types.SwapPair
	typeStr := strings.ToLower(pair)

	var price math.LegacyDec

	assetPrices, err := k.validationKeeper.GetAssetPrices(ctx)
	if err != nil {
		return nil, math.LegacyDec{}, err
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
		price, err = k.validationKeeper.GetRockBtcPrice(ctx)
		if err != nil {
			return nil, math.LegacyDec{}, err
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
		price, err = k.validationKeeper.GetBtcRockPrice(ctx)
		if err != nil {
			return nil, math.LegacyDec{}, err
		}
	default:
		return nil, math.LegacyDec{}, fmt.Errorf("unknown key type: %s", pair)
	}

	if price.IsZero() || price.IsNegative() {
		return nil, math.LegacyDec{}, fmt.Errorf("price must be positive, got: %s", price.String())
	}

	return &pairType, price, nil
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

func (k Keeper) GetPrice(ctx sdk.Context, pair string) (math.LegacyDec, error) {
	switch pair {
	case "rockbtc":
		return k.validationKeeper.GetRockBtcPrice(ctx)
	case "btcrock":
		return k.validationKeeper.GetBtcRockPrice(ctx)
	default:
		return math.LegacyDec{}, fmt.Errorf("unknown pair: %s", pair)
	}
}

func (k Keeper) GetAmountOut(ctx sdk.Context, pair string, amountIn uint64, price math.LegacyDec) (uint64, error) {
	switch pair {
	case "rockbtc":
		// returns BTC amount in satoshis to transfer
		amountInDec := math.LegacyNewDecFromInt(math.NewIntFromUint64(amountIn))
		satoshisDec := amountInDec.Mul(price.Abs())
		satoshis := satoshisDec.TruncateInt().Uint64()
		if k.GetParams(ctx).MinimumSatoshis > satoshis {
			return 0, types.ErrMinimumSatoshis
		}
		return satoshis, nil
	case "btcrock":
		if k.GetParams(ctx).MinimumSatoshis > amountIn {
			return 0, types.ErrMinimumSatoshis
		}
		// returns ROCK amount in urock to transfer
		amountInDec := math.LegacyNewDecFromInt(math.NewIntFromUint64(amountIn))
		urockDec := amountInDec.Mul(price.Abs())
		return urockDec.TruncateInt().Uint64(), nil
	default:
		return 0, fmt.Errorf("unknown pair: %s", pair)
	}
}

func (k Keeper) EscrowRock(ctx sdk.Context, senderKey treasurytypes.Key, amount uint64) error {

	if senderKey.Type != treasurytypes.KeyType_KEY_TYPE_ECDSA_SECP256K1 {
		return types.ErrWrongKeyType
	}

	senderAddress, err := treasurytypes.NativeAddress(&senderKey, "zen")
	if err != nil {
		return fmt.Errorf("failed to convert sender key to zenrock address: %w", err)
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.MustAccAddressFromBech32(senderAddress), types.ZenexCollectorName, sdk.NewCoins(sdk.NewCoin(params.BondDenom, math.NewIntFromUint64(amount))))
	if err != nil {
		return err
	}

	return nil
}
