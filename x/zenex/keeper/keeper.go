package keeper

import (
	"fmt"

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
		accountKeeper    types.AccountKeeper

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
	accountKeeper types.AccountKeeper,
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
		accountKeeper:    accountKeeper,

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

func (k Keeper) GetPair(ctx sdk.Context, pair types.TradePair) (*types.SwapPair, math.LegacyDec, error) {
	var pairType types.SwapPair

	var price math.LegacyDec

	assetPrices, err := k.validationKeeper.GetAssetPrices(ctx)
	if err != nil {
		return nil, math.LegacyDec{}, err
	}

	if assetPrices[validationtypes.Asset_ROCK].IsZero() || assetPrices[validationtypes.Asset_BTC].IsZero() {
		return nil, math.LegacyDec{}, fmt.Errorf("price is zero, check sidecar consensus, got: ROCK=%s, BTC=%s", assetPrices[validationtypes.Asset_ROCK].String(), assetPrices[validationtypes.Asset_BTC].String())
	}

	switch pair {
	case types.TradePair_TRADE_PAIR_ROCK_BTC:
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
	case types.TradePair_TRADE_PAIR_BTC_ROCK:
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
		return nil, math.LegacyDec{}, fmt.Errorf("unknown pair: %s", pair)
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

func (k Keeper) GetPrice(ctx sdk.Context, pair types.TradePair) (math.LegacyDec, error) {
	switch pair {
	case types.TradePair_TRADE_PAIR_ROCK_BTC:
		return k.validationKeeper.GetRockBtcPrice(ctx)
	case types.TradePair_TRADE_PAIR_BTC_ROCK:
		return k.validationKeeper.GetBtcRockPrice(ctx)
	default:
		return math.LegacyDec{}, fmt.Errorf("unknown pair: %s", pair)
	}
}

func (k Keeper) GetAmountOut(ctx sdk.Context, pair types.TradePair, amountIn uint64, price math.LegacyDec) (uint64, error) {
	switch pair {
	case types.TradePair_TRADE_PAIR_ROCK_BTC:
		// returns BTC amount in satoshis to transfer
		amountInDec := math.LegacyNewDecFromInt(math.NewIntFromUint64(amountIn))
		satoshisDec := amountInDec.Mul(price.Abs())
		satoshis := satoshisDec.TruncateInt().Uint64()
		if k.GetParams(ctx).MinimumSatoshis > satoshis {
			return 0, fmt.Errorf("calculated satoshis %d is less than the minimum satoshis %d", satoshis, k.GetParams(ctx).MinimumSatoshis)
		}
		return satoshis, nil
	case types.TradePair_TRADE_PAIR_BTC_ROCK:
		if k.GetParams(ctx).MinimumSatoshis > amountIn {
			return 0, fmt.Errorf("%d satoshis in is less than the minimum satoshis %d", amountIn, k.GetParams(ctx).MinimumSatoshis)
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

	rockAddress, err := k.GetRockAddress(ctx, senderKey.Id)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.MustAccAddressFromBech32(rockAddress), types.ZenexCollectorName, sdk.NewCoins(sdk.NewCoin(params.BondDenom, math.NewIntFromUint64(amount))))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetRockAddress(ctx sdk.Context, rockKeyId uint64) (string, error) {
	rockKey, err := k.treasuryKeeper.GetKey(ctx, rockKeyId)
	if err != nil {
		return "", err
	}
	rockAddress, err := treasurytypes.NativeAddress(rockKey, "zen")
	if err != nil {
		return "", err
	}
	return rockAddress, nil
}

func (k Keeper) CheckRedeemableAsset(ctx sdk.Context, amountOut uint64, price math.LegacyDec) error {
	availableRockBalance := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(types.ZenexCollectorName), params.BondDenom).Amount.Uint64()

	if amountOut > availableRockBalance {
		return fmt.Errorf("amount %d is greater than the available rock balance %d", amountOut, availableRockBalance)
	}

	return nil
}
