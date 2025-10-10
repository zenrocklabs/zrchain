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
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"

	treasury "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/keeper"
	validation "github.com/Zenrock-Foundation/zrchain/v6/x/validation/keeper"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		validationKeeper *validation.Keeper
		treasuryKeeper   *treasury.Keeper

		authority string

		Schema collections.Schema
		Params collections.Item[types.Params]
		// LockTransactions - key: hash of lock transaction rawTx + vout | value: lock transaction data
		LockTransactions collections.Map[string, types.LockTransaction]
		// PendingMintTransactionsMap - key: pending zenBTC mint transaction id | value: pending zenBTC mint transaction
		PendingMintTransactionsMap collections.Map[uint64, types.PendingMintTransaction]
		// FirstPendingEthMintTransaction - value: lowest key of pending Ethereum mint transaction
		FirstPendingEthMintTransaction collections.Item[uint64]
		// FirstPendingSolMintTransaction - value: lowest key of pending Solana mint transaction
		FirstPendingSolMintTransaction collections.Item[uint64]
		// PendingMintTransactionCount - value: count of pending zenBTC mint transactions
		PendingMintTransactionCount collections.Item[uint64]
		// BurnEvents - key: burn event index | value: burn event data
		BurnEvents collections.Map[uint64, types.BurnEvent]
		// FirstPendingBurnEvent - value: lowest key of pending burn event
		FirstPendingBurnEvent collections.Item[uint64]
		// BurnEventCount - value: count of burn events
		BurnEventCount collections.Item[uint64]
		// Redemptions - key: redemption index | value: redemption data
		Redemptions collections.Map[uint64, types.Redemption]
		// FirstPendingRedemption - value: lowest key of pending redemption
		FirstPendingRedemption collections.Item[uint64]
		// FirstRedemptionAwaitingSign - value: lowest key of pending redemption awaiting signature
		FirstRedemptionAwaitingSign collections.Item[uint64]
		// Supply - value: zenBTC supply data
		Supply collections.Item[types.Supply]
		// FirstPendingStakeTransaction - value: lowest key of pending stake transaction
		FirstPendingStakeTransaction collections.Item[uint64]
		// DEPRECATED
		LockTransactionStore collections.Map[collections.Pair[string, uint64], types.LockTransaction]
		// DEPRECATED
		PendingMintTransactions collections.Item[types.PendingMintTransactions]
	}
)

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
		cdc:                            cdc,
		storeService:                   storeService,
		logger:                         logger,
		validationKeeper:               validationKeeper,
		treasuryKeeper:                 treasuryKeeper,
		authority:                      authority,
		Params:                         collections.NewItem(sb, types.ParamsKey, types.ParamsIndex, codec.CollValue[types.Params](cdc)),
		LockTransactions:               collections.NewMap(sb, types.LockTransactionsNewKey, types.LockTransactionsNewIndex, collections.StringKey, codec.CollValue[types.LockTransaction](cdc)),
		LockTransactionStore:           collections.NewMap(sb, types.LockTransactionsKey, types.LockTransactionsIndex, collections.PairKeyCodec(collections.StringKey, collections.Uint64Key), codec.CollValue[types.LockTransaction](cdc)),
		PendingMintTransactions:        collections.NewItem(sb, types.PendingMintTransactionsKey, types.PendingMintTransactionsIndex, codec.CollValue[types.PendingMintTransactions](cdc)),
		PendingMintTransactionsMap:     collections.NewMap(sb, types.PendingMintTransactionsMapKey, types.PendingMintTransactionsMapIndex, collections.Uint64Key, codec.CollValue[types.PendingMintTransaction](cdc)),
		FirstPendingStakeTransaction:   collections.NewItem(sb, types.FirstPendingStakeTransactionKey, types.FirstPendingStakeTransactionIndex, collections.Uint64Value),
		FirstPendingEthMintTransaction: collections.NewItem(sb, types.FirstPendingEthMintTransactionKey, types.FirstPendingEthMintTransactionIndex, collections.Uint64Value),
		FirstPendingSolMintTransaction: collections.NewItem(sb, types.FirstPendingSolMintTransactionKey, types.FirstPendingSolMintTransactionIndex, collections.Uint64Value),
		PendingMintTransactionCount:    collections.NewItem(sb, types.PendingMintTransactionCountKey, types.PendingMintTransactionCountIndex, collections.Uint64Value),
		BurnEvents:                     collections.NewMap(sb, types.BurnEventsKey, types.BurnEventsIndex, collections.Uint64Key, codec.CollValue[types.BurnEvent](cdc)),
		FirstPendingBurnEvent:          collections.NewItem(sb, types.FirstPendingBurnEventKey, types.FirstPendingBurnEventIndex, collections.Uint64Value),
		BurnEventCount:                 collections.NewItem(sb, types.BurnEventCountKey, types.BurnEventCountIndex, collections.Uint64Value),
		Redemptions:                    collections.NewMap(sb, types.RedemptionsKey, types.RedemptionsIndex, collections.Uint64Key, codec.CollValue[types.Redemption](cdc)),
		FirstPendingRedemption:         collections.NewItem(sb, types.FirstPendingRedemptionKey, types.FirstPendingRedemptionIndex, collections.Uint64Value),
		FirstRedemptionAwaitingSign:    collections.NewItem(sb, types.FirstRedemptionAwaitingSignKey, types.FirstRedemptionAwaitingSignIndex, collections.Uint64Value),
		Supply:                         collections.NewItem(sb, types.SupplyKey, types.SupplyIndex, codec.CollValue[types.Supply](cdc)),
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
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetZenBTCExchangeRate returns the current exchange rate between BTC and zenBTC
// Returns the number of BTC represented by 1 zenBTC
func (k Keeper) GetExchangeRate(ctx context.Context) (math.LegacyDec, error) {
	supply, err := k.Supply.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return math.LegacyNewDec(0), err
		}
		return math.LegacyNewDec(1), nil // Initial exchange rate of 1:1
	}

	totalZenBTC := supply.MintedZenBTC + supply.PendingZenBTC
	if totalZenBTC == 0 {
		return math.LegacyNewDec(1), nil // If no mints/deposits yet, use 1:1 rate
	}

	return math.LegacyNewDecFromInt(math.NewIntFromUint64(supply.CustodiedBTC)).Quo(math.LegacyNewDecFromInt(math.NewIntFromUint64(totalZenBTC))), nil
}

func (k Keeper) SetPendingMintTransaction(ctx context.Context, pendingMintTransaction types.PendingMintTransaction) error {
	return k.PendingMintTransactionsMap.Set(ctx, pendingMintTransaction.Id, pendingMintTransaction)
}

func (k Keeper) WalkPendingMintTransactions(ctx context.Context, fn func(id uint64, pendingMintTransaction types.PendingMintTransaction) (stop bool, err error)) error {
	return k.PendingMintTransactionsMap.Walk(ctx, nil, fn)
}

func (k Keeper) HasRedemption(ctx context.Context, id uint64) (bool, error) {
	return k.Redemptions.Has(ctx, id)
}

func (k Keeper) SetRedemption(ctx context.Context, id uint64, redemption types.Redemption) error {
	return k.Redemptions.Set(ctx, id, redemption)
}

func (k Keeper) WalkRedemptions(ctx context.Context, fn func(id uint64, redemption types.Redemption) (stop bool, err error)) error {
	return k.Redemptions.Walk(ctx, nil, fn)
}

func (k Keeper) WalkRedemptionsDescending(ctx context.Context, fn func(id uint64, redemption types.Redemption) (stop bool, err error)) error {
	rng := new(collections.Range[uint64]).Descending()
	return k.Redemptions.Walk(ctx, rng, fn)
}

func (k Keeper) GetSupply(ctx context.Context) (types.Supply, error) {
	return k.Supply.Get(ctx)
}

func (k Keeper) SetSupply(ctx context.Context, supply types.Supply) error {
	return k.Supply.Set(ctx, supply)
}

func (k Keeper) GetBurnEvent(ctx context.Context, id uint64) (types.BurnEvent, error) {
	return k.BurnEvents.Get(ctx, id)
}

func (k Keeper) SetBurnEvent(ctx context.Context, id uint64, burnEvent types.BurnEvent) error {
	return k.BurnEvents.Set(ctx, id, burnEvent)
}

func (k Keeper) WalkBurnEvents(ctx context.Context, fn func(id uint64, burnEvent types.BurnEvent) (stop bool, err error)) error {
	return k.BurnEvents.Walk(ctx, nil, fn)
}

func (k Keeper) GetRedemption(ctx context.Context, id uint64) (types.Redemption, error) {
	return k.Redemptions.Get(ctx, id)
}

func (k Keeper) GetPendingMintTransaction(ctx context.Context, id uint64) (types.PendingMintTransaction, error) {
	return k.PendingMintTransactionsMap.Get(ctx, id)
}

func (k Keeper) HasPendingMintTransaction(ctx context.Context, id uint64) (bool, error) {
	return k.PendingMintTransactionsMap.Has(ctx, id)
}

// GetFirstPendingMintTransaction returns the ID of the first pending mint transaction
func (k Keeper) GetFirstPendingEthMintTransaction(ctx context.Context) (uint64, error) {
	return k.FirstPendingEthMintTransaction.Get(ctx)
}

// SetFirstPendingEthMintTransaction sets the ID of the first pending Ethereum mint transaction
func (k Keeper) SetFirstPendingEthMintTransaction(ctx context.Context, id uint64) error {
	return k.FirstPendingEthMintTransaction.Set(ctx, id)
}

// GetFirstPendingSolMintTransaction returns the ID of the first pending Solana mint transaction
func (k Keeper) GetFirstPendingSolMintTransaction(ctx context.Context) (uint64, error) {
	return k.FirstPendingSolMintTransaction.Get(ctx)
}

// SetFirstPendingSolMintTransaction sets the ID of the first pending Solana mint transaction
func (k Keeper) SetFirstPendingSolMintTransaction(ctx context.Context, id uint64) error {
	return k.FirstPendingSolMintTransaction.Set(ctx, id)
}

// GetFirstPendingBurnEvent returns the ID of the first pending burn event
func (k Keeper) GetFirstPendingBurnEvent(ctx context.Context) (uint64, error) {
	return k.FirstPendingBurnEvent.Get(ctx)
}

// SetFirstPendingBurnEvent sets the ID of the first pending burn event
func (k Keeper) SetFirstPendingBurnEvent(ctx context.Context, id uint64) error {
	return k.FirstPendingBurnEvent.Set(ctx, id)
}

// GetFirstPendingRedemption returns the ID of the first pending redemption
func (k Keeper) GetFirstPendingRedemption(ctx context.Context) (uint64, error) {
	return k.FirstPendingRedemption.Get(ctx)
}

// SetFirstPendingRedemption sets the ID of the first pending redemption
func (k Keeper) SetFirstPendingRedemption(ctx context.Context, id uint64) error {
	return k.FirstPendingRedemption.Set(ctx, id)
}

// GetFirstPendingStakeTransaction returns the ID of the first pending stake transaction
func (k Keeper) GetFirstPendingStakeTransaction(ctx context.Context) (uint64, error) {
	return k.FirstPendingStakeTransaction.Get(ctx)
}

// SetFirstPendingStakeTransaction sets the ID of the first pending stake transaction
func (k Keeper) SetFirstPendingStakeTransaction(ctx context.Context, id uint64) error {
	return k.FirstPendingStakeTransaction.Set(ctx, id)
}

// GetFirstRedemptionAwaitingSign returns the ID of the first pending redemption awaiting signature
func (k Keeper) GetFirstRedemptionAwaitingSign(ctx context.Context) (uint64, error) {
	return k.FirstRedemptionAwaitingSign.Get(ctx)
}

// SetFirstRedemptionAwaitingSign sets the ID of the first pending redemption awaiting signature
func (k Keeper) SetFirstRedemptionAwaitingSign(ctx context.Context, id uint64) error {
	return k.FirstRedemptionAwaitingSign.Set(ctx, id)
}

func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) SetAuthority(authority string) {
	k.authority = authority
}

func (k Keeper) GetParams(ctx context.Context) (types.Params, error) {
	return k.Params.Get(ctx)
}

func (k Keeper) GetLockTransactionsMap(ctx context.Context) (map[string]types.LockTransaction, error) {
	lockTxMap := make(map[string]types.LockTransaction)
	err := k.LockTransactions.Walk(ctx, nil, func(key string, value types.LockTransaction) (bool, error) {
		lockTxMap[key] = value
		return false, nil
	})
	return lockTxMap, err
}

func (k Keeper) GetPendingMintTransactionsMap(ctx context.Context) (map[uint64]types.PendingMintTransaction, error) {
	pendingMintTxMap := make(map[uint64]types.PendingMintTransaction)
	err := k.PendingMintTransactionsMap.Walk(ctx, nil, func(key uint64, value types.PendingMintTransaction) (bool, error) {
		pendingMintTxMap[key] = value
		return false, nil
	})
	return pendingMintTxMap, err
}

func (k Keeper) GetBurnEventsMap(ctx context.Context) (map[uint64]types.BurnEvent, error) {
	burnEventMap := make(map[uint64]types.BurnEvent)
	err := k.BurnEvents.Walk(ctx, nil, func(key uint64, value types.BurnEvent) (bool, error) {
		burnEventMap[key] = value
		return false, nil
	})
	return burnEventMap, err
}

func (k Keeper) GetRedemptionsMap(ctx context.Context) (map[uint64]types.Redemption, error) {
	redemptionMap := make(map[uint64]types.Redemption)
	err := k.Redemptions.Walk(ctx, nil, func(key uint64, value types.Redemption) (bool, error) {
		redemptionMap[key] = value
		return false, nil
	})
	return redemptionMap, err
}
