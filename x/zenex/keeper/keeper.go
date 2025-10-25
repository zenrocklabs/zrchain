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
		zenbtcKeeper     types.ZenbtcKeeper

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
	zenbtcKeeper types.ZenbtcKeeper,
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
		zenbtcKeeper:     zenbtcKeeper,
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

// Returns the pair and trade pair price
// including the asset prices, base and quote token
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

// Returns all swaps from the swaps store
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

// Returns the price for a given trade pair
// based on the asset prices from the validation keeper
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

// Calculates the amount out for a given pair and amount in
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

// Escrows rock from the sender key to the zenex collector from a swap
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

// Returns the rock address for the given rock key id
// as a recipient address for zenex swaps
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

// Checks if we have enough rock balance to swap for BTC
func (k Keeper) CheckRedeemableAsset(ctx sdk.Context, amountOut uint64, price math.LegacyDec) error {
	availableRockBalance := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(types.ZenexCollectorName), params.BondDenom).Amount.Uint64()

	if amountOut > availableRockBalance {
		return fmt.Errorf("amount %d is greater than the available rock balance %d", amountOut, availableRockBalance)
	}

	return nil
}

// Calculates the minimum rock balance required to swap for BTC
func (k Keeper) GetRequiredRockBalance(ctx sdk.Context) (uint64, error) {

	rockBtcPrice, err := k.validationKeeper.GetRockBtcPrice(ctx)
	if err != nil {
		return 0, err
	}

	if rockBtcPrice.IsZero() || rockBtcPrice.IsNegative() {
		return 0, fmt.Errorf("rock to btc price must be positive, got: %s", rockBtcPrice.String())
	}

	thresholdSatoshis := math.LegacyNewDecFromInt(math.NewIntFromUint64(k.GetParams(ctx).SwapThresholdSatoshis))

	requiredRockBalance := thresholdSatoshis.Quo(rockBtcPrice)

	if requiredRockBalance.IsZero() || requiredRockBalance.IsNegative() {
		return 0, fmt.Errorf("required rock balance is zero or negative, got: %s", requiredRockBalance.String())
	}

	return requiredRockBalance.TruncateInt().Uint64(), nil
}

// Returns the balance of the zenbtc rewards collector
func (k Keeper) GetZenBtcRewardsCollectorBalance(ctx sdk.Context) uint64 {
	balance := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(types.ZenBtcRewardsCollectorName), params.BondDenom)
	return balance.Amount.Uint64()
}

// swaps ROCK for BTC and funds the zenbtc reward address with BTC
// if the rock fee pool balance is greater than the required rock balance
// checks if it's over the minimum amount of satoshis
func (k Keeper) CreateRockBtcSwap(ctx sdk.Context, amountIn uint64) error {

	swapPair, price, err := k.GetPair(ctx, types.TradePair_TRADE_PAIR_ROCK_BTC)
	if err != nil {
		return err
	}

	amountOutRaw, err := k.GetAmountOut(ctx, types.TradePair_TRADE_PAIR_ROCK_BTC, amountIn, price)
	if err != nil {
		return err
	}

	zenbtcparams, err := k.zenbtcKeeper.GetParams(ctx)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ZenBtcRewardsCollectorName, types.ZenexCollectorName, sdk.NewCoins(sdk.NewCoin(params.BondDenom, math.NewIntFromUint64(amountIn))))
	if err != nil {
		return err
	}

	swapCount, err := k.SwapsCount.Get(ctx)
	if err != nil {
		return err
	}

	swapCount++
	swap := types.Swap{
		Creator: "",
		SwapId:  swapCount,
		Status:  types.SwapStatus_SWAP_STATUS_INITIATED,
		Pair:    types.TradePair_TRADE_PAIR_ROCK_BTC,
		Data: &types.SwapData{
			BaseToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_ROCK,
				PriceUSD:  swapPair.BaseToken.PriceUSD,
				Precision: 6,
			},
			QuoteToken: &validationtypes.AssetData{
				Asset:     validationtypes.Asset_BTC,
				PriceUSD:  swapPair.QuoteToken.PriceUSD,
				Precision: 8,
			},
			Price:     price,
			AmountIn:  amountIn,
			AmountOut: amountOutRaw, // TODO: consider reducing the amount out by 500 satoshis for btc fee
		},
		RockKeyId:      0,
		BtcKeyId:       zenbtcparams.RewardsDepositKeyID,
		ZenexPoolKeyId: k.GetParams(ctx).ZenexPoolKeyId,
		Workspace:      "",
		ZenbtcSwap:     true,
	}

	err = k.SwapsStore.Set(ctx, swapCount, swap)
	if err != nil {
		return err
	}

	err = k.SwapsCount.Set(ctx, swapCount)
	if err != nil {
		return err
	}

	return nil
}
