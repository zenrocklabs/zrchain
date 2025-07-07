package keeper

import (
	"context"
	"errors"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"

	"cosmossdk.io/collections"
	addresscodec "cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	"github.com/Zenrock-Foundation/zrchain/v6/shared"
	sidecar "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	zenbtctypes "github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

type Keeper struct {
	storeService          storetypes.KVStoreService
	cdc                   codec.BinaryCodec
	authKeeper            types.AccountKeeper
	bankKeeper            types.BankKeeper
	hooks                 types.StakingHooks
	authority             string
	treasuryKeeper        types.TreasuryKeeper
	zenBTCKeeper          shared.ZenBTCKeeper
	validatorAddressCodec addresscodec.Codec
	consensusAddressCodec addresscodec.Codec
	txDecoder             sdk.TxDecoder
	zrConfig              *params.ZRConfig
	sidecarClient         sidecarClient
	zentpKeeper           types.ZentpKeeper
	// AVSDelegations - keys: validator addr + delegator addr (operator) | value: delegation amount
	AVSDelegations collections.Map[collections.Pair[string, string], math.Int]
	// ValidatorDelegations - key: validator addr | value: total amount delegated to validator
	ValidatorDelegations collections.Map[string, math.Int]
	// AVSRewardsPool - key: address | value: total unclaimed rewards for that address
	AVSRewardsPool collections.Map[string, math.Int]
	// AssetPrices - key: asset type | value: asset price + precision
	AssetPrices collections.Map[types.Asset, math.LegacyDec]
	// LastValidVEHeight - value: height of last valid VE
	LastValidVEHeight collections.Item[int64]
	// SlashEvents - key: id number | value: slash event struct
	SlashEvents collections.Map[uint64, types.SlashEvent]
	// SlashEventCount - value: number of slash events
	SlashEventCount collections.Item[uint64]
	// Params related to Hybrid Validation (separate from vanilla x/staking SDK module due to interface constraints)
	HVParams collections.Item[types.HVParams]
	// ValidationInfos - key: block height | value: validation info for block
	ValidationInfos collections.Map[int64, types.ValidationInfo]
	// BitcoinMerkleRoots - key: block height | value: merkle root of Bitcoin block
	BtcBlockHeaders collections.Map[int64, sidecar.BTCBlockHeader]
	// LatestBtcHeaderHeight - value: height of the latest btc header stored
	LatestBtcHeaderHeight collections.Item[int64]
	// EthereumNonceRequested - key: key ID | value: bool (is requested)
	EthereumNonceRequested collections.Map[uint64, bool]
	// SolanaNonceRequested - key: key ID | value: bool (is requested)
	SolanaNonceRequested collections.Map[uint64, bool]
	// LastUsedEthereumNonce - map: key ID | value: last used Ethereum nonce data
	SolanaAccountsRequested      collections.Map[string, bool]
	SolanaZenTPAccountsRequested collections.Map[string, bool]
	LastUsedEthereumNonce        collections.Map[uint64, zenbtctypes.NonceData]
	LastUsedSolanaNonce          collections.Map[uint64, types.SolanaNonce]
	// RequestedHistoricalBitcoinHeaders - keys: block height
	RequestedHistoricalBitcoinHeaders collections.Item[zenbtctypes.RequestedBitcoinHeaders]
	// BackfillRequests - key: tx hash | value: bool (is requested)
	BackfillRequests collections.Item[types.BackfillRequests]
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	authority string,
	txDecoder sdk.TxDecoder,
	zrConfig *params.ZRConfig,
	treasuryKeeper types.TreasuryKeeper,
	zenBTCKeeper shared.ZenBTCKeeper,
	zentpKeeper types.ZentpKeeper,
	validatorAddressCodec addresscodec.Codec,
	consensusAddressCodec addresscodec.Codec,
) *Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	// ensure bonded and not bonded module accounts are set
	if addr := ak.GetModuleAddress(types.BondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	if addr := ak.GetModuleAddress(types.NotBondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	// ensure that authority is a valid AccAddress
	if _, err := ak.AddressCodec().StringToBytes(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	if validatorAddressCodec == nil || consensusAddressCodec == nil {
		panic("validator and/or consensus address codec are nil")
	}

	// The config isn't created when we export the genesis state, so we blank it out
	sidecarAddr := ""
	if zrConfig != nil {
		sidecarAddr = zrConfig.SidecarAddr
	}
	oracleClient, err := NewSidecarClient(sidecarAddr)
	if err != nil {
		panic("error creating sidecar client")
	}

	sb := collections.NewSchemaBuilder(storeService)

	return &Keeper{
		storeService:                      storeService,
		cdc:                               cdc,
		authKeeper:                        ak,
		bankKeeper:                        bk,
		hooks:                             nil,
		authority:                         authority,
		txDecoder:                         txDecoder,
		zrConfig:                          zrConfig,
		sidecarClient:                     oracleClient,
		treasuryKeeper:                    treasuryKeeper,
		zenBTCKeeper:                      zenBTCKeeper,
		zentpKeeper:                       zentpKeeper,
		validatorAddressCodec:             validatorAddressCodec,
		consensusAddressCodec:             consensusAddressCodec,
		AVSDelegations:                    collections.NewMap(sb, types.AVSDelegationsKey, types.AVSDelegationsIndex, collections.PairKeyCodec(collections.StringKey, collections.StringKey), sdk.IntValue),
		ValidatorDelegations:              collections.NewMap(sb, types.ValidatorDelegationsKey, types.ValidatorDelegationsIndex, collections.StringKey, sdk.IntValue),
		AVSRewardsPool:                    collections.NewMap(sb, types.AVSRewardsPoolKey, types.AVSRewardsPoolIndex, collections.StringKey, sdk.IntValue),
		AssetPrices:                       collections.NewMap(sb, types.AssetPricesKey, types.AssetPricesIndex, types.AssetKey{}, sdk.LegacyDecValue),
		SlashEvents:                       collections.NewMap(sb, types.SlashEventsKey, types.SlashEventsIndex, collections.Uint64Key, codec.CollValue[types.SlashEvent](cdc)),
		SlashEventCount:                   collections.NewItem(sb, types.SlashEventCountKey, types.SlashEventCountIndex, collections.Uint64Value),
		HVParams:                          collections.NewItem(sb, types.HVParamsKey, types.HVParamsIndex, codec.CollValue[types.HVParams](cdc)),
		ValidationInfos:                   collections.NewMap(sb, types.ValidationInfosKey, types.ValidationInfosIndex, collections.Int64Key, codec.CollValue[types.ValidationInfo](cdc)),
		BtcBlockHeaders:                   collections.NewMap(sb, types.BtcBlockHeadersKey, types.BtcBlockHeadersIndex, collections.Int64Key, codec.CollValue[sidecar.BTCBlockHeader](cdc)),
		LatestBtcHeaderHeight:             collections.NewItem(sb, types.LatestBtcHeaderHeightKey, types.LatestBtcHeaderHeightIndex, collections.Int64Value),
		EthereumNonceRequested:            collections.NewMap(sb, types.EthereumNonceRequestedKey, types.EthereumNonceRequestedIndex, collections.Uint64Key, collections.BoolValue),
		SolanaNonceRequested:              collections.NewMap(sb, types.SolanaNonceRequestedKey, types.SolanaNonceRequestedIndex, collections.Uint64Key, collections.BoolValue),
		SolanaAccountsRequested:           collections.NewMap(sb, types.SolanaAccountsRequestedKey, types.SolanaAccountsRequestedIndex, collections.StringKey, collections.BoolValue),
		SolanaZenTPAccountsRequested:      collections.NewMap(sb, types.SolanaZenTPAccountsRequestedKey, types.SolanaZenTPAccountsRequestedIndex, collections.StringKey, collections.BoolValue),
		LastUsedEthereumNonce:             collections.NewMap(sb, types.LastUsedEthereumNonceKey, types.LastUsedEthereumNonceIndex, collections.Uint64Key, codec.CollValue[zenbtctypes.NonceData](cdc)),
		LastUsedSolanaNonce:               collections.NewMap(sb, types.LastUsedSolanaNonceKey, types.LastUsedSolanaNonceIndex, collections.Uint64Key, codec.CollValue[types.SolanaNonce](cdc)),
		RequestedHistoricalBitcoinHeaders: collections.NewItem(sb, types.RequestedHistoricalBitcoinHeadersKey, types.RequestedHistoricalBitcoinHeadersIndex, codec.CollValue[zenbtctypes.RequestedBitcoinHeaders](cdc)),
		LastValidVEHeight:                 collections.NewItem(sb, types.LastValidVEHeightKey, types.LastValidVEHeightIndex, collections.Int64Value),
		BackfillRequests:                  collections.NewItem(sb, types.BackfillRequestsKey, types.BackfillRequestsIndex, codec.CollValue[types.BackfillRequests](cdc)),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	return sdk.UnwrapSDKContext(ctx).Logger().With("module", "x/"+types.ModuleName)
}

// Hooks gets the hooks for staking *Keeper {
func (k *Keeper) Hooks() types.StakingHooks {
	if k.hooks == nil {
		// return a no-op implementation if no hooks are set
		return types.MultiStakingHooks{}
	}
	return k.hooks
}

// SetHooks sets the validator hooks.  In contrast to other receivers, this method must take a pointer due to nature
// of the hooks interface and SDK start up sequence.
func (k *Keeper) SetHooks(sh types.StakingHooks) {
	if k.hooks != nil {
		panic("cannot set validator hooks twice")
	}
	k.hooks = sh
}

// GetLastTotalPower loads the last total validator power.
func (k Keeper) GetLastTotalPower(ctx context.Context) (math.Int, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.LastTotalPowerKey)
	if err != nil {
		return math.ZeroInt(), err
	}

	if bz == nil {
		return math.ZeroInt(), nil
	}

	ip := sdk.IntProto{}
	err = k.cdc.Unmarshal(bz, &ip)
	if err != nil {
		return math.ZeroInt(), err
	}

	return ip.Int, nil
}

// SetLastTotalPower sets the last total validator power.
func (k Keeper) SetLastTotalPower(ctx context.Context, power math.Int) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&sdk.IntProto{Int: power})
	if err != nil {
		return err
	}
	return store.Set(types.LastTotalPowerKey, bz)
}

// GetAuthority returns the x/staking module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// ValidatorAddressCodec returns the app validator address codec.
func (k Keeper) ValidatorAddressCodec() addresscodec.Codec {
	return k.validatorAddressCodec
}

// ConsensusAddressCodec returns the app consensus address codec.
func (k Keeper) ConsensusAddressCodec() addresscodec.Codec {
	return k.consensusAddressCodec
}

// SetValidatorUpdates sets the ABCI validator power updates for the current block.
func (k Keeper) SetValidatorUpdates(ctx context.Context, valUpdates []abci.ValidatorUpdate) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&types.ValidatorUpdates{Updates: valUpdates})
	if err != nil {
		return err
	}
	return store.Set(types.ValidatorUpdatesKey, bz)
}

// GetValidatorUpdates returns the ABCI validator power updates within the current block.
func (k Keeper) GetValidatorUpdates(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.ValidatorUpdatesKey)
	if err != nil {
		return nil, err
	}

	var valUpdates types.ValidatorUpdates
	err = k.cdc.Unmarshal(bz, &valUpdates)
	if err != nil {
		return nil, err
	}

	return valUpdates.Updates, nil
}

// SetSolanaZenBTCRequestedAccount sets the requested state for a Solana account (owner address) for ZenBTC.
func (k Keeper) SetSolanaZenBTCRequestedAccount(ctx context.Context, ownerAddress string, state bool) error {
	return k.SolanaAccountsRequested.Set(ctx, ownerAddress, state)
}

// SetSolanaZenTPRequestedAccount sets the requested state for a Solana account (owner address) for ZenTP.
func (k Keeper) SetSolanaZenTPRequestedAccount(ctx context.Context, ownerAddress string, state bool) error {
	return k.SolanaZenTPAccountsRequested.Set(ctx, ownerAddress, state)
}

func (k Keeper) SetSolanaRequestedNonce(ctx context.Context, keyID uint64, state bool) error {
	return k.SolanaNonceRequested.Set(ctx, keyID, state)
}

// SetSidecarClient sets the sidecar client for the keeper.
func (k *Keeper) SetSidecarClient(client sidecarClient) {
	k.sidecarClient = client
}

func (k *Keeper) SetBackfillRequests(ctx context.Context, requests types.BackfillRequests) error {
	return k.BackfillRequests.Set(ctx, requests)
}

func (k Keeper) GetAssetPrices(ctx context.Context) ([]*types.AssetData, error) {
	storeIterator, err := k.AssetPrices.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer storeIterator.Close()

	keys, err := storeIterator.Keys()
	if err != nil {
		return nil, err
	}

	var assetPrices []*types.AssetData
	for _, key := range keys {
		value, err := k.AssetPrices.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		assetPrices = append(assetPrices, &types.AssetData{
			Asset:    key,
			PriceUSD: value,
		})
	}

	return assetPrices, nil
}

func (k Keeper) GetLastValidVeHeight(ctx context.Context) (int64, error) {
	lastValidVeHeight, err := k.LastValidVEHeight.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// Return 0 when the collection is empty
			return 0, nil
		}
		return 0, err
	}
	return lastValidVeHeight, nil
}

func (k Keeper) GetSlashEvents(ctx context.Context) ([]types.SlashEvent, uint64, error) {
	storeIterator, err := k.SlashEvents.Iterate(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer storeIterator.Close()

	keys, err := storeIterator.Keys()
	if err != nil {
		return nil, 0, err
	}

	var slashEvents []types.SlashEvent
	var slashEventCount uint64
	for _, key := range keys {
		value, err := k.SlashEvents.Get(ctx, key)
		if err != nil {
			return nil, 0, err
		}
		slashEvents = append(slashEvents, value)
		slashEventCount++
	}

	return slashEvents, slashEventCount, nil
}

func (k Keeper) GetValidationInfos(ctx context.Context) ([]types.ValidationInfo, error) {
	storeIterator, err := k.ValidationInfos.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer storeIterator.Close()

	keys, err := storeIterator.Keys()
	if err != nil {
		return nil, err
	}

	var validationInfos []types.ValidationInfo
	for _, key := range keys {
		value, err := k.ValidationInfos.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		validationInfos = append(validationInfos, value)
	}
	return validationInfos, nil
}

func (k Keeper) GetBtcBlockHeaders(ctx context.Context) ([]sidecar.BTCBlockHeader, error) {
	storeIterator, err := k.BtcBlockHeaders.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer storeIterator.Close()

	keys, err := storeIterator.Keys()
	if err != nil {
		return nil, err
	}

	var btcBlockHeaders []sidecar.BTCBlockHeader
	for _, key := range keys {
		value, err := k.BtcBlockHeaders.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		btcBlockHeaders = append(btcBlockHeaders, value)
	}
	return btcBlockHeaders, nil
}

func (k Keeper) GetLastUsedSolanaNonce(ctx context.Context) ([]types.SolanaNonce, error) {
	storeIterator, err := k.LastUsedSolanaNonce.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer storeIterator.Close()

	keys, err := storeIterator.Keys()
	if err != nil {
		return nil, err
	}

	var lastUsedSolanaNonce []types.SolanaNonce
	for _, key := range keys {
		value, err := k.LastUsedSolanaNonce.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		lastUsedSolanaNonce = append(lastUsedSolanaNonce, value)
	}
	return lastUsedSolanaNonce, nil
}

func (k Keeper) GetBackfillRequests(ctx context.Context) (types.BackfillRequests, error) {
	backfillRequest, err := k.BackfillRequests.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// Return empty BackfillRequests when the collection is empty
			return types.BackfillRequests{}, nil
		}
		return types.BackfillRequests{}, err
	}

	return backfillRequest, nil
}

func (k Keeper) GetLastUsedEthereumNonce(ctx context.Context) ([]zenbtctypes.NonceData, error) {
	storeIterator, err := k.LastUsedEthereumNonce.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer storeIterator.Close()

	keys, err := storeIterator.Keys()
	if err != nil {
		return nil, err
	}

	var lastUsedEthereumNonce []zenbtctypes.NonceData
	for _, key := range keys {
		value, err := k.LastUsedEthereumNonce.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		lastUsedEthereumNonce = append(lastUsedEthereumNonce, value)
	}
	return lastUsedEthereumNonce, nil
}

func (k Keeper) GetRequestedHistoricalBitcoinHeaders(ctx context.Context) ([]zenbtctypes.RequestedBitcoinHeaders, error) {
	requestedHistoricalBitcoinHeaders, err := k.RequestedHistoricalBitcoinHeaders.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// Return empty RequestedBitcoinHeaders when the collection is empty
			return []zenbtctypes.RequestedBitcoinHeaders{{}}, nil
		}
		return nil, err
	}

	return []zenbtctypes.RequestedBitcoinHeaders{requestedHistoricalBitcoinHeaders}, nil
}

func (k Keeper) GetAvsRewardsPool(ctx context.Context) ([]string, error) {
	avsRewardsPool := []string{}

	storeIterator, err := k.AVSRewardsPool.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	keys, err := storeIterator.Keys()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		avsRewardsPool = append(avsRewardsPool, key)
	}

	return avsRewardsPool, nil
}

func (k Keeper) GetEthereumNonceRequested(ctx context.Context) ([]uint64, error) {
	ethereumNonceRequested := []uint64{}

	err := k.EthereumNonceRequested.Walk(ctx, nil, func(key uint64, value bool) (stop bool, err error) {
		ethereumNonceRequested = append(ethereumNonceRequested, key)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return ethereumNonceRequested, nil
}

func (k Keeper) GetSolanaNonceRequested(ctx context.Context) ([]uint64, error) {
	solanaNonceRequested := []uint64{}

	err := k.SolanaNonceRequested.Walk(ctx, nil, func(key uint64, value bool) (stop bool, err error) {
		solanaNonceRequested = append(solanaNonceRequested, key)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return solanaNonceRequested, nil
}

func (k Keeper) GetSolanaAccountsRequested(ctx context.Context) ([]string, error) {
	solanaAccountsRequested := []string{}

	err := k.SolanaAccountsRequested.Walk(ctx, nil, func(key string, value bool) (stop bool, err error) {
		solanaAccountsRequested = append(solanaAccountsRequested, key)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return solanaAccountsRequested, nil
}

func (k Keeper) GetSolanaZenTPAccountsRequested(ctx context.Context) ([]string, error) {
	solanaZenTPAccountsRequested := []string{}

	err := k.SolanaZenTPAccountsRequested.Walk(ctx, nil, func(key string, value bool) (stop bool, err error) {
		solanaZenTPAccountsRequested = append(solanaZenTPAccountsRequested, key)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return solanaZenTPAccountsRequested, nil
}
