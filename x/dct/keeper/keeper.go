package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"

	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	treasury "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/keeper"
	validation "github.com/Zenrock-Foundation/zrchain/v6/x/validation/keeper"
)

var (
	errAssetNotConfigured = errors.New("asset not configured")
)

// Keeper manages decentralised custody token state.
type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	logger       log.Logger

	validationKeeper *validation.Keeper
	treasuryKeeper   *treasury.Keeper

	authority string

	Schema collections.Schema

	Params collections.Item[dcttypes.Params]

	LockTransactions collections.Map[collections.Pair[string, string], dcttypes.LockTransaction]

	PendingMintTransactions        collections.Map[collections.Pair[string, uint64], dcttypes.PendingMintTransaction]
	PendingMintTransactionCount    collections.Map[string, uint64]
	FirstPendingStakeTransaction   collections.Map[string, uint64]
	FirstPendingSolMintTransaction collections.Map[string, uint64]
	FirstPendingEthMintTransaction collections.Map[string, uint64]

	BurnEvents            collections.Map[collections.Pair[string, uint64], dcttypes.BurnEvent]
	BurnEventCount        collections.Map[string, uint64]
	FirstPendingBurnEvent collections.Map[string, uint64]

	Redemptions                 collections.Map[collections.Pair[string, uint64], dcttypes.Redemption]
	FirstPendingRedemption      collections.Map[string, uint64]
	FirstRedemptionAwaitingSign collections.Map[string, uint64]

	Supply collections.Map[string, dcttypes.Supply]
}

// NewKeeper constructs a new DCT keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	validationKeeper *validation.Keeper,
	treasuryKeeper *treasury.Keeper,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:              cdc,
		storeService:     storeService,
		logger:           logger,
		validationKeeper: validationKeeper,
		treasuryKeeper:   treasuryKeeper,
		authority:        authority,
		Params:           collections.NewItem(sb, dcttypes.ParamsKey, dcttypes.ParamsIndex, codec.CollValue[dcttypes.Params](cdc)),
		LockTransactions: collections.NewMap(sb, dcttypes.LockTransactionsNewKey, dcttypes.LockTransactionsNewIndex, collections.PairKeyCodec(collections.StringKey, collections.StringKey), codec.CollValue[dcttypes.LockTransaction](cdc)),
		PendingMintTransactions: collections.NewMap(
			sb,
			dcttypes.PendingMintTransactionsMapKey,
			dcttypes.PendingMintTransactionsMapIndex,
			collections.PairKeyCodec(collections.StringKey, collections.Uint64Key),
			codec.CollValue[dcttypes.PendingMintTransaction](cdc),
		),
		PendingMintTransactionCount: collections.NewMap(
			sb,
			dcttypes.PendingMintTransactionCountKey,
			dcttypes.PendingMintTransactionCountIndex,
			collections.StringKey,
			collections.Uint64Value,
		),
		FirstPendingStakeTransaction: collections.NewMap(
			sb,
			dcttypes.FirstPendingStakeTransactionKey,
			dcttypes.FirstPendingStakeTransactionIndex,
			collections.StringKey,
			collections.Uint64Value,
		),
		FirstPendingSolMintTransaction: collections.NewMap(
			sb,
			dcttypes.FirstPendingSolMintTransactionKey,
			dcttypes.FirstPendingSolMintTransactionIndex,
			collections.StringKey,
			collections.Uint64Value,
		),
		FirstPendingEthMintTransaction: collections.NewMap(
			sb,
			dcttypes.FirstPendingEthMintTransactionKey,
			dcttypes.FirstPendingEthMintTransactionIndex,
			collections.StringKey,
			collections.Uint64Value,
		),
		BurnEvents: collections.NewMap(
			sb,
			dcttypes.BurnEventsKey,
			dcttypes.BurnEventsIndex,
			collections.PairKeyCodec(collections.StringKey, collections.Uint64Key),
			codec.CollValue[dcttypes.BurnEvent](cdc),
		),
		BurnEventCount: collections.NewMap(
			sb,
			dcttypes.BurnEventCountKey,
			dcttypes.BurnEventCountIndex,
			collections.StringKey,
			collections.Uint64Value,
		),
		FirstPendingBurnEvent: collections.NewMap(
			sb,
			dcttypes.FirstPendingBurnEventKey,
			dcttypes.FirstPendingBurnEventIndex,
			collections.StringKey,
			collections.Uint64Value,
		),
		Redemptions: collections.NewMap(
			sb,
			dcttypes.RedemptionsKey,
			dcttypes.RedemptionsIndex,
			collections.PairKeyCodec(collections.StringKey, collections.Uint64Key),
			codec.CollValue[dcttypes.Redemption](cdc),
		),
		FirstPendingRedemption: collections.NewMap(
			sb,
			dcttypes.FirstPendingRedemptionKey,
			dcttypes.FirstPendingRedemptionIndex,
			collections.StringKey,
			collections.Uint64Value,
		),
		FirstRedemptionAwaitingSign: collections.NewMap(
			sb,
			dcttypes.FirstRedemptionAwaitingSignKey,
			dcttypes.FirstRedemptionAwaitingSignIndex,
			collections.StringKey,
			collections.Uint64Value,
		),
		Supply: collections.NewMap(
			sb,
			dcttypes.SupplyKey,
			dcttypes.SupplyIndex,
			collections.StringKey,
			codec.CollValue[dcttypes.Supply](cdc),
		),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return &k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", dcttypes.ModuleName))
}

// GetAuthority returns the governance authority for the module.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// GetParams returns the global module parameters.
func (k Keeper) GetParams(ctx context.Context) (dcttypes.Params, error) {
	return k.Params.Get(ctx)
}

// SetParams updates the global module parameters. Mostly used by governance proposals.
func (k Keeper) SetParams(ctx context.Context, params dcttypes.Params) error {
	return k.Params.Set(ctx, params)
}

// ListSupportedAssets returns all assets that have configuration.
// Always includes ASSET_ZENZEC to ensure nonce account retrieval works.
func (k Keeper) ListSupportedAssets(ctx context.Context) ([]dcttypes.Asset, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	assetsMap := make(map[dcttypes.Asset]bool)
	for _, ap := range params.Assets {
		if ap.Asset == dcttypes.Asset_ASSET_UNSPECIFIED {
			continue
		}
		assetsMap[ap.Asset] = true
	}

	// Always ensure zenZEC is included, even if params are empty or misconfigured
	assetsMap[dcttypes.Asset_ASSET_ZENZEC] = true

	assets := make([]dcttypes.Asset, 0, len(assetsMap))
	for asset := range assetsMap {
		assets = append(assets, asset)
	}
	return assets, nil
}

// GetAssetParams returns the parameter set associated with the given asset.
func (k Keeper) GetAssetParams(ctx context.Context, asset dcttypes.Asset) (dcttypes.AssetParams, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return dcttypes.AssetParams{}, err
	}
	for _, ap := range params.Assets {
		if ap.Asset == asset {
			return ap, nil
		}
	}
	return dcttypes.AssetParams{}, fmt.Errorf("%w: %s", errAssetNotConfigured, asset.String())
}

// GetSolanaParams returns the Solana configuration for the asset.
func (k Keeper) GetSolanaParams(ctx context.Context, asset dcttypes.Asset) (*dcttypes.Solana, error) {
	ap, err := k.GetAssetParams(ctx, asset)
	if err != nil {
		return nil, err
	}
	if ap.Solana == nil {
		return nil, dcttypes.ErrMissingSolanaData
	}
	return ap.Solana, nil
}

func (k Keeper) getAssetKey(asset dcttypes.Asset) (string, error) {
	if asset == dcttypes.Asset_ASSET_UNSPECIFIED {
		return "", dcttypes.ErrUnknownAsset
	}
	return asset.String(), nil
}

func (k Keeper) assetFromKey(assetKey string) (dcttypes.Asset, error) {
	if value, ok := dcttypes.Asset_value[assetKey]; ok {
		return dcttypes.Asset(value), nil
	}
	return dcttypes.Asset_ASSET_UNSPECIFIED, dcttypes.ErrUnknownAsset
}

// --- Parameter helpers ---

func (k Keeper) GetStakerKeyID(ctx context.Context, asset dcttypes.Asset) (uint64, error) {
	ap, err := k.GetAssetParams(ctx, asset)
	if err != nil {
		return 0, err
	}
	return ap.StakerKeyId, nil
}

func (k Keeper) GetRewardsDepositKeyID(ctx context.Context, asset dcttypes.Asset) (uint64, error) {
	ap, err := k.GetAssetParams(ctx, asset)
	if err != nil {
		return 0, err
	}
	return ap.RewardsDepositKeyId, nil
}

func (k Keeper) GetEthMinterKeyID(ctx context.Context, asset dcttypes.Asset) (uint64, error) {
	ap, err := k.GetAssetParams(ctx, asset)
	if err != nil {
		return 0, err
	}
	return ap.EthMinterKeyId, nil
}

func (k Keeper) GetUnstakerKeyID(ctx context.Context, asset dcttypes.Asset) (uint64, error) {
	ap, err := k.GetAssetParams(ctx, asset)
	if err != nil {
		return 0, err
	}
	return ap.UnstakerKeyId, nil
}

func (k Keeper) GetCompleterKeyID(ctx context.Context, asset dcttypes.Asset) (uint64, error) {
	ap, err := k.GetAssetParams(ctx, asset)
	if err != nil {
		return 0, err
	}
	return ap.CompleterKeyId, nil
}

func (k Keeper) GetChangeAddressKeyIDs(ctx context.Context, asset dcttypes.Asset) ([]uint64, error) {
	ap, err := k.GetAssetParams(ctx, asset)
	if err != nil {
		return nil, err
	}
	return ap.ChangeAddressKeyIds, nil
}

func (k Keeper) GetControllerAddr(ctx context.Context, asset dcttypes.Asset) (string, error) {
	ap, err := k.GetAssetParams(ctx, asset)
	if err != nil {
		return "", err
	}
	return ap.ControllerAddr, nil
}

func (k Keeper) GetEthTokenAddr(ctx context.Context, asset dcttypes.Asset) (string, error) {
	ap, err := k.GetAssetParams(ctx, asset)
	if err != nil {
		return "", err
	}
	return ap.EthTokenAddr, nil
}

func (k Keeper) GetDepositKeyringAddr(ctx context.Context, asset dcttypes.Asset) (string, error) {
	ap, err := k.GetAssetParams(ctx, asset)
	if err != nil {
		return "", err
	}
	return ap.DepositKeyringAddr, nil
}

func (k Keeper) GetProxyAddress(ctx context.Context, asset dcttypes.Asset) (string, error) {
	ap, err := k.GetAssetParams(ctx, asset)
	if err != nil {
		return "", err
	}
	return ap.ProxyAddress, nil
}

func (k Keeper) GetBitcoinProxyAddress(ctx context.Context, asset dcttypes.Asset) (string, error) {
	// For Bitcoin-fork assets the general proxy address doubles as the asset-specific proxy.
	return k.GetProxyAddress(ctx, asset)
}

// --- Lock transactions ---

func (k Keeper) lockKey(asset dcttypes.Asset, key string) (collections.Pair[string, string], error) {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return collections.Pair[string, string]{}, err
	}
	return collections.Join(assetKey, key), nil
}

func (k Keeper) pendingMintKey(asset dcttypes.Asset, id uint64) (collections.Pair[string, uint64], error) {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return collections.Pair[string, uint64]{}, err
	}
	return collections.Join(assetKey, id), nil
}

// --- Pending mint transactions ---

func (k Keeper) SetPendingMintTransaction(ctx context.Context, tx dcttypes.PendingMintTransaction) error {
	key, err := k.pendingMintKey(tx.Asset, tx.Id)
	if err != nil {
		return err
	}

	return k.PendingMintTransactions.Set(ctx, key, tx)
}

func (k Keeper) GetPendingMintTransaction(ctx context.Context, asset dcttypes.Asset, id uint64) (dcttypes.PendingMintTransaction, error) {
	key, err := k.pendingMintKey(asset, id)
	if err != nil {
		return dcttypes.PendingMintTransaction{}, err
	}
	return k.PendingMintTransactions.Get(ctx, key)
}

func (k Keeper) HasPendingMintTransaction(ctx context.Context, asset dcttypes.Asset, id uint64) (bool, error) {
	key, err := k.pendingMintKey(asset, id)
	if err != nil {
		return false, err
	}
	return k.PendingMintTransactions.Has(ctx, key)
}

func (k Keeper) WalkPendingMintTransactions(ctx context.Context, asset dcttypes.Asset, fn func(id uint64, tx dcttypes.PendingMintTransaction) (bool, error)) error {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return err
	}

	rng := collections.NewPrefixedPairRange[string, uint64](assetKey)
	return k.PendingMintTransactions.Walk(ctx, rng, func(key collections.Pair[string, uint64], value dcttypes.PendingMintTransaction) (bool, error) {
		return fn(key.K2(), value)
	})
}

func (k Keeper) GetPendingMintTransactionCount(ctx context.Context, asset dcttypes.Asset) (uint64, error) {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return 0, err
	}
	return k.PendingMintTransactionCount.Get(ctx, assetKey)
}

func (k Keeper) SetPendingMintTransactionCount(ctx context.Context, asset dcttypes.Asset, count uint64) error {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return err
	}
	return k.PendingMintTransactionCount.Set(ctx, assetKey, count)
}

func (k Keeper) GetFirstPendingStakeTransaction(ctx context.Context, asset dcttypes.Asset) (uint64, error) {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return 0, err
	}
	return k.FirstPendingStakeTransaction.Get(ctx, assetKey)
}

func (k Keeper) SetFirstPendingStakeTransaction(ctx context.Context, asset dcttypes.Asset, id uint64) error {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return err
	}
	return k.FirstPendingStakeTransaction.Set(ctx, assetKey, id)
}

func (k Keeper) GetFirstPendingSolMintTransaction(ctx context.Context, asset dcttypes.Asset) (uint64, error) {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return 0, err
	}
	return k.FirstPendingSolMintTransaction.Get(ctx, assetKey)
}

func (k Keeper) SetFirstPendingSolMintTransaction(ctx context.Context, asset dcttypes.Asset, id uint64) error {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return err
	}
	return k.FirstPendingSolMintTransaction.Set(ctx, assetKey, id)
}

func (k Keeper) GetFirstPendingEthMintTransaction(ctx context.Context, asset dcttypes.Asset) (uint64, error) {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return 0, err
	}
	return k.FirstPendingEthMintTransaction.Get(ctx, assetKey)
}

func (k Keeper) SetFirstPendingEthMintTransaction(ctx context.Context, asset dcttypes.Asset, id uint64) error {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return err
	}
	return k.FirstPendingEthMintTransaction.Set(ctx, assetKey, id)
}

// --- Supply ---

func (k Keeper) GetSupply(ctx context.Context, asset dcttypes.Asset) (dcttypes.Supply, error) {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return dcttypes.Supply{}, err
	}
	return k.Supply.Get(ctx, assetKey)
}

func (k Keeper) SetSupply(ctx context.Context, supply dcttypes.Supply) error {
	assetKey, err := k.getAssetKey(supply.Asset)
	if err != nil {
		return err
	}
	return k.Supply.Set(ctx, assetKey, supply)
}

func (k Keeper) GetExchangeRate(ctx context.Context, asset dcttypes.Asset) (math.LegacyDec, error) {
	supply, err := k.GetSupply(ctx, asset)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return math.LegacyNewDec(1), nil
		}
		return math.LegacyNewDec(0), err
	}

	totalMinted := supply.MintedAmount + supply.PendingAmount
	if totalMinted == 0 {
		return math.LegacyNewDec(1), nil
	}

	return math.LegacyNewDecFromInt(math.NewIntFromUint64(supply.CustodiedAmount)).Quo(math.LegacyNewDecFromInt(math.NewIntFromUint64(totalMinted))), nil
}

// --- Burn events ---

func (k Keeper) burnKey(asset dcttypes.Asset, id uint64) (collections.Pair[string, uint64], error) {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return collections.Pair[string, uint64]{}, err
	}
	return collections.Join(assetKey, id), nil
}

func (k Keeper) SetBurnEvent(ctx context.Context, asset dcttypes.Asset, id uint64, burnEvent dcttypes.BurnEvent) error {
	burnEvent.Asset = asset
	key, err := k.burnKey(asset, id)
	if err != nil {
		return err
	}
	return k.BurnEvents.Set(ctx, key, burnEvent)
}

func (k Keeper) GetBurnEvent(ctx context.Context, asset dcttypes.Asset, id uint64) (dcttypes.BurnEvent, error) {
	key, err := k.burnKey(asset, id)
	if err != nil {
		return dcttypes.BurnEvent{}, err
	}
	return k.BurnEvents.Get(ctx, key)
}

func (k Keeper) WalkBurnEvents(ctx context.Context, asset dcttypes.Asset, fn func(id uint64, burnEvent dcttypes.BurnEvent) (bool, error)) error {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return err
	}
	rng := collections.NewPrefixedPairRange[string, uint64](assetKey)
	return k.BurnEvents.Walk(ctx, rng, func(key collections.Pair[string, uint64], value dcttypes.BurnEvent) (bool, error) {
		return fn(key.K2(), value)
	})
}

func (k Keeper) GetBurnEventCount(ctx context.Context, asset dcttypes.Asset) (uint64, error) {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return 0, err
	}
	return k.BurnEventCount.Get(ctx, assetKey)
}

func (k Keeper) SetBurnEventCount(ctx context.Context, asset dcttypes.Asset, count uint64) error {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return err
	}
	return k.BurnEventCount.Set(ctx, assetKey, count)
}

func (k Keeper) GetFirstPendingBurnEvent(ctx context.Context, asset dcttypes.Asset) (uint64, error) {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return 0, err
	}
	return k.FirstPendingBurnEvent.Get(ctx, assetKey)
}

func (k Keeper) SetFirstPendingBurnEvent(ctx context.Context, asset dcttypes.Asset, id uint64) error {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return err
	}
	return k.FirstPendingBurnEvent.Set(ctx, assetKey, id)
}

func (k Keeper) CreateBurnEvent(ctx context.Context, asset dcttypes.Asset, burnEvent *dcttypes.BurnEvent) (uint64, error) {
	count, err := k.GetBurnEventCount(ctx, asset)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return 0, err
		}
		count = 0
	}
	nextID := count + 1
	burnEvent.Id = nextID
	burnEvent.Asset = asset
	if err := k.SetBurnEvent(ctx, asset, nextID, *burnEvent); err != nil {
		return 0, err
	}
	if err := k.SetBurnEventCount(ctx, asset, nextID); err != nil {
		return 0, err
	}
	return nextID, nil
}

// --- Redemptions ---

func (k Keeper) redemptionKey(asset dcttypes.Asset, id uint64) (collections.Pair[string, uint64], error) {
	return k.burnKey(asset, id)
}

func (k Keeper) SetRedemption(ctx context.Context, asset dcttypes.Asset, id uint64, redemption dcttypes.Redemption) error {
	redemption.Data.Asset = asset
	key, err := k.redemptionKey(asset, id)
	if err != nil {
		return err
	}
	return k.Redemptions.Set(ctx, key, redemption)
}

func (k Keeper) GetRedemption(ctx context.Context, asset dcttypes.Asset, id uint64) (dcttypes.Redemption, error) {
	key, err := k.redemptionKey(asset, id)
	if err != nil {
		return dcttypes.Redemption{}, err
	}
	return k.Redemptions.Get(ctx, key)
}

func (k Keeper) HasRedemption(ctx context.Context, asset dcttypes.Asset, id uint64) (bool, error) {
	key, err := k.redemptionKey(asset, id)
	if err != nil {
		return false, err
	}
	return k.Redemptions.Has(ctx, key)
}

func (k Keeper) WalkRedemptions(ctx context.Context, asset dcttypes.Asset, fn func(id uint64, redemption dcttypes.Redemption) (bool, error)) error {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return err
	}
	rng := collections.NewPrefixedPairRange[string, uint64](assetKey)
	return k.Redemptions.Walk(ctx, rng, func(key collections.Pair[string, uint64], value dcttypes.Redemption) (bool, error) {
		return fn(key.K2(), value)
	})
}

func (k Keeper) WalkRedemptionsDescending(ctx context.Context, asset dcttypes.Asset, fn func(id uint64, redemption dcttypes.Redemption) (bool, error)) error {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return err
	}
	rng := collections.NewPrefixedPairRange[string, uint64](assetKey).Descending()
	return k.Redemptions.Walk(ctx, rng, func(key collections.Pair[string, uint64], value dcttypes.Redemption) (bool, error) {
		return fn(key.K2(), value)
	})
}

func (k Keeper) GetFirstPendingRedemption(ctx context.Context, asset dcttypes.Asset) (uint64, error) {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return 0, err
	}
	return k.FirstPendingRedemption.Get(ctx, assetKey)
}

func (k Keeper) SetFirstPendingRedemption(ctx context.Context, asset dcttypes.Asset, id uint64) error {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return err
	}
	return k.FirstPendingRedemption.Set(ctx, assetKey, id)
}

func (k Keeper) GetFirstRedemptionAwaitingSign(ctx context.Context, asset dcttypes.Asset) (uint64, error) {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return 0, err
	}
	return k.FirstRedemptionAwaitingSign.Get(ctx, assetKey)
}

func (k Keeper) SetFirstRedemptionAwaitingSign(ctx context.Context, asset dcttypes.Asset, id uint64) error {
	assetKey, err := k.getAssetKey(asset)
	if err != nil {
		return err
	}
	return k.FirstRedemptionAwaitingSign.Set(ctx, assetKey, id)
}
