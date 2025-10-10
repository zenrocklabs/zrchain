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
	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	zenbtctypes "github.com/Zenrock-Foundation/zrchain/v6/x/zenbtc/types"
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
	dctKeeper             shared.DCTKeeper
	validatorAddressCodec addresscodec.Codec
	consensusAddressCodec addresscodec.Codec
	txDecoder             sdk.TxDecoder
	zrConfig              *params.ZRConfig
	sidecarClient         sidecarClient
	zentpKeeper           types.ZentpKeeper
	slashingKeeper        types.SlashingKeeper
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
	SolanaNonceRequested         collections.Map[uint64, bool]
	SolanaAccountsRequested      collections.Map[string, bool]
	SolanaZenTPAccountsRequested collections.Map[string, bool]
	SolanaDCTAccountsRequested   collections.Map[collections.Pair[string, string], bool]
	// LastUsedEthereumNonce - map: key ID | value: last used Ethereum nonce data
	LastUsedEthereumNonce collections.Map[uint64, zenbtctypes.NonceData]
	LastUsedSolanaNonce   collections.Map[uint64, types.SolanaNonce]
	// RequestedHistoricalBitcoinHeaders - keys: block height
	RequestedHistoricalBitcoinHeaders collections.Item[zenbtctypes.RequestedBitcoinHeaders]
	// BackfillRequests - key: tx hash | value: bool (is requested)
	BackfillRequests collections.Item[types.BackfillRequests]
	// ValidatorMismatchCounts - key: validator hex address | value: mismatch count data
	ValidatorMismatchCounts collections.Map[string, types.ValidatorMismatchCount]
	// LastCompletedZentpMintID - value: id of last completed zentp mint
	LastCompletedZentpMintID collections.Item[uint64]
	// BedrockValidatorSet - key: validator address | value: is in bedrock validator set
	BedrockValidatorSet collections.Map[string, bool]
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
	dctKeeper shared.DCTKeeper,
	zentpKeeper types.ZentpKeeper,
	slashingKeeper types.SlashingKeeper,
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
		dctKeeper:                         dctKeeper,
		zentpKeeper:                       zentpKeeper,
		slashingKeeper:                    slashingKeeper,
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
		SolanaDCTAccountsRequested:        collections.NewMap(sb, types.SolanaDCTAccountsRequestedKey, types.SolanaDCTAccountsRequestedIndex, collections.PairKeyCodec(collections.StringKey, collections.StringKey), collections.BoolValue),
		LastUsedEthereumNonce:             collections.NewMap(sb, types.LastUsedEthereumNonceKey, types.LastUsedEthereumNonceIndex, collections.Uint64Key, codec.CollValue[zenbtctypes.NonceData](cdc)),
		LastUsedSolanaNonce:               collections.NewMap(sb, types.LastUsedSolanaNonceKey, types.LastUsedSolanaNonceIndex, collections.Uint64Key, codec.CollValue[types.SolanaNonce](cdc)),
		RequestedHistoricalBitcoinHeaders: collections.NewItem(sb, types.RequestedHistoricalBitcoinHeadersKey, types.RequestedHistoricalBitcoinHeadersIndex, codec.CollValue[zenbtctypes.RequestedBitcoinHeaders](cdc)),
		LastValidVEHeight:                 collections.NewItem(sb, types.LastValidVEHeightKey, types.LastValidVEHeightIndex, collections.Int64Value),
		BackfillRequests:                  collections.NewItem(sb, types.BackfillRequestsKey, types.BackfillRequestsIndex, codec.CollValue[types.BackfillRequests](cdc)),
		ValidatorMismatchCounts:           collections.NewMap(sb, types.ValidatorMismatchCounts, types.ValidatorMismatchCountsIndex, collections.StringKey, codec.CollValue[types.ValidatorMismatchCount](cdc)),
		LastCompletedZentpMintID:          collections.NewItem(sb, types.LastCompletedZentpMintIDKey, types.LastCompletedZentpMintIDIndex, collections.Uint64Value),
		BedrockValidatorSet:               collections.NewMap(sb, types.BedrockValidatorSetKey, types.BedrockValidatorSetIndex, collections.StringKey, collections.BoolValue),
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

func (k Keeper) SetSolanaDCTRequestedAccount(ctx context.Context, asset dcttypes.Asset, ownerAddress string, state bool) error {
	key, err := k.dctAccountKey(asset, ownerAddress)
	if err != nil {
		return err
	}
	return k.SolanaDCTAccountsRequested.Set(ctx, key, state)
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

func (k Keeper) dctAccountKey(asset dcttypes.Asset, ownerAddress string) (collections.Pair[string, string], error) {
	if asset == dcttypes.Asset_ASSET_UNSPECIFIED {
		return collections.Pair[string, string]{}, dcttypes.ErrUnknownAsset
	}
	if ownerAddress == "" {
		return collections.Pair[string, string]{}, fmt.Errorf("owner address must be provided")
	}
	return collections.Join(asset.String(), ownerAddress), nil
}

func (k Keeper) GetAssetPrices(ctx context.Context) (map[types.Asset]math.LegacyDec, error) {
	assetPrices := make(map[types.Asset]math.LegacyDec)

	// Fetch only known assets to avoid decoding legacy/corrupt entries that may be stored under the same prefix.
	// TODO: Gracefully handle the case where the asset is not found.
	knownAssets := []types.Asset{types.Asset_ROCK, types.Asset_BTC, types.Asset_ETH}
	for _, asset := range knownAssets {
		price, err := k.AssetPrices.Get(ctx, asset)
		if err != nil {
			if errors.Is(err, collections.ErrNotFound) {
				k.Logger(ctx).Info("Asset not found", "key", asset.String())
				continue
			}
			k.Logger(ctx).Error("Failed to read asset price", "asset", asset.String(), "error", err)
			continue
		}

		assetPrices[asset] = price
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

func (k Keeper) GetLastCompletedZentpMintID(ctx context.Context) (uint64, error) {
	lastCompletedZentpMintID, err := k.LastCompletedZentpMintID.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			if err := k.LastCompletedZentpMintID.Set(ctx, 0); err != nil {
				return 0, err
			}
			return 0, nil
		}
		return 0, err
	}
	return lastCompletedZentpMintID, nil
}

func (k Keeper) GetSlashEvents(ctx context.Context) (map[uint64]types.SlashEvent, error) {
	slashEvents := make(map[uint64]types.SlashEvent)

	err := k.SlashEvents.Walk(ctx, nil, func(key uint64, value types.SlashEvent) (stop bool, err error) {
		slashEvents[key] = value
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return slashEvents, nil
}

func (k Keeper) GetValidationInfos(ctx context.Context) (map[int64]types.ValidationInfo, error) {
	validationInfos := make(map[int64]types.ValidationInfo)

	// Limit iteration window to avoid OOM: walk from a recent height only.
	// TODO: only used for module exports, remove this once we have a better way to export the data
	const maxHeightsToScan = int64(10)
	currentHeight := sdk.UnwrapSDKContext(ctx).BlockHeight()
	startHeight := currentHeight - maxHeightsToScan
	if startHeight < 0 {
		startHeight = 0
	}
	queryRange := &collections.Range[int64]{}

	err := k.ValidationInfos.Walk(ctx, queryRange.StartInclusive(startHeight), func(key int64, value types.ValidationInfo) (stop bool, err error) {
		validationInfos[key] = value
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return validationInfos, nil
}

func (k Keeper) GetBtcBlockHeaders(ctx context.Context) (map[int64]sidecar.BTCBlockHeader, error) {
	btcBlockHeaders := make(map[int64]sidecar.BTCBlockHeader)

	err := k.BtcBlockHeaders.Walk(ctx, nil, func(key int64, value sidecar.BTCBlockHeader) (stop bool, err error) {
		btcBlockHeaders[key] = value
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return btcBlockHeaders, nil
}

func (k Keeper) GetLastUsedSolanaNonce(ctx context.Context) (map[uint64]types.SolanaNonce, error) {
	lastUsedSolanaNonce := make(map[uint64]types.SolanaNonce)

	err := k.LastUsedSolanaNonce.Walk(ctx, nil, func(key uint64, value types.SolanaNonce) (stop bool, err error) {
		lastUsedSolanaNonce[key] = value
		return false, nil
	})
	if err != nil {
		return nil, err
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

func (k Keeper) GetLastUsedEthereumNonce(ctx context.Context) (map[uint64]zenbtctypes.NonceData, error) {
	lastUsedEthereumNonce := make(map[uint64]zenbtctypes.NonceData)

	err := k.LastUsedEthereumNonce.Walk(ctx, nil, func(key uint64, value zenbtctypes.NonceData) (stop bool, err error) {
		lastUsedEthereumNonce[key] = value
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return lastUsedEthereumNonce, nil
}

func (k Keeper) GetRequestedHistoricalBitcoinHeaders(ctx context.Context) (zenbtctypes.RequestedBitcoinHeaders, error) {
	requestedHistoricalBitcoinHeaders, err := k.RequestedHistoricalBitcoinHeaders.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return zenbtctypes.RequestedBitcoinHeaders{}, nil
		}
		return zenbtctypes.RequestedBitcoinHeaders{}, err
	}

	return requestedHistoricalBitcoinHeaders, nil
}

func (k Keeper) GetAvsRewardsPool(ctx context.Context) (map[string]math.Int, error) {
	avsRewardsPool := make(map[string]math.Int)

	err := k.AVSRewardsPool.Walk(ctx, nil, func(key string, value math.Int) (stop bool, err error) {
		avsRewardsPool[key] = value
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return avsRewardsPool, nil
}

func (k Keeper) GetEthereumNonceRequested(ctx context.Context) (map[uint64]bool, error) {
	ethereumNonceRequested := make(map[uint64]bool)

	err := k.EthereumNonceRequested.Walk(ctx, nil, func(key uint64, value bool) (stop bool, err error) {
		ethereumNonceRequested[key] = value
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return ethereumNonceRequested, nil
}

func (k Keeper) GetSolanaNonceRequested(ctx context.Context) (map[uint64]bool, error) {
	solanaNonceRequested := make(map[uint64]bool)

	err := k.SolanaNonceRequested.Walk(ctx, nil, func(key uint64, value bool) (stop bool, err error) {
		solanaNonceRequested[key] = value
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return solanaNonceRequested, nil
}

func (k Keeper) GetSolanaAccountsRequested(ctx context.Context) (map[string]bool, error) {
	solanaAccountsRequested := make(map[string]bool)

	err := k.SolanaAccountsRequested.Walk(ctx, nil, func(key string, value bool) (stop bool, err error) {
		solanaAccountsRequested[key] = value
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return solanaAccountsRequested, nil
}

func (k Keeper) GetSolanaZenTPAccountsRequested(ctx context.Context) (map[string]bool, error) {
	solanaZenTPAccountsRequested := make(map[string]bool)

	err := k.SolanaZenTPAccountsRequested.Walk(ctx, nil, func(key string, value bool) (stop bool, err error) {
		solanaZenTPAccountsRequested[key] = value
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return solanaZenTPAccountsRequested, nil
}

func (k Keeper) GetSolanaDCTAccountsRequested(ctx context.Context) (map[string]bool, error) {
	requested := make(map[string]bool)
	err := k.SolanaDCTAccountsRequested.Walk(ctx, nil, func(key collections.Pair[string, string], value bool) (bool, error) {
		compositeKey := fmt.Sprintf("%s:%s", key.K1(), key.K2())
		requested[compositeKey] = value
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return requested, nil
}

func (k Keeper) GetValidatorMismatchCounts(ctx context.Context) (map[string]types.ValidatorMismatchCount, error) {
	validatorMismatchCounts := make(map[string]types.ValidatorMismatchCount)

	err := k.ValidatorMismatchCounts.Walk(ctx, nil, func(key string, value types.ValidatorMismatchCount) (stop bool, err error) {
		validatorMismatchCounts[key] = value
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return validatorMismatchCounts, nil
}

func (k Keeper) GetRockBtcPrice(ctx context.Context) (math.LegacyDec, error) {
	rockPrice, err := k.AssetPrices.Get(ctx, types.Asset_ROCK)
	if err != nil {
		return math.LegacyDec{}, err
	}
	btcPrice, err := k.AssetPrices.Get(ctx, types.Asset_BTC)
	if err != nil {
		return math.LegacyDec{}, err
	}

	return rockPrice.Quo(btcPrice), nil
}

func (k Keeper) GetBtcRockPrice(ctx context.Context) (math.LegacyDec, error) {
	rockPrice, err := k.AssetPrices.Get(ctx, types.Asset_ROCK)
	if err != nil {
		return math.LegacyDec{}, err
	}
	btcPrice, err := k.AssetPrices.Get(ctx, types.Asset_BTC)
	if err != nil {
		return math.LegacyDec{}, err
	}

	return btcPrice.Quo(rockPrice), nil
}

func (k Keeper) GetAssets(ctx context.Context) ([]types.Asset, error) {

	assets := []types.Asset{types.Asset_ROCK, types.Asset_BTC, types.Asset_ETH}

	return assets, nil
}
