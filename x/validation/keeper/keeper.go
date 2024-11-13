package keeper

import (
	"context"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"

	"cosmossdk.io/collections"
	addresscodec "cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v5/app/params"
	sidecar "github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
	treasury "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/keeper"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v5/x/validation/types"
)

type Keeper struct {
	storeService          storetypes.KVStoreService
	cdc                   codec.BinaryCodec
	authKeeper            types.AccountKeeper
	bankKeeper            types.BankKeeper
	hooks                 types.StakingHooks
	authority             string
	treasuryKeeper        *treasury.Keeper
	validatorAddressCodec addresscodec.Codec
	consensusAddressCodec addresscodec.Codec
	txDecoder             sdk.TxDecoder
	zrConfig              *params.ZRConfig
	sidecarClient         sidecarClient

	// AVSDelegations - keys: validator addr + delegator addr (operator) | value: delegation amount
	AVSDelegations collections.Map[collections.Pair[string, string], math.Int]
	// ValidatorDelegations - key: validator addr | value: total amount delegated to validator
	ValidatorDelegations collections.Map[string, math.Int]
	// AVSRewardsPool - key: address | value: total unclaimed rewards for that address
	AVSRewardsPool collections.Map[string, math.Int]
	// AssetPrices - key: asset symbol | value: asset price struct
	AssetPrices collections.Map[string, types.AssetPrice]
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
	// ConfirmedUnlockTxs - key: chain + tx hash | value: withdrawal metadata
	ConfirmedUnlockTxs collections.Map[collections.Pair[string, string], types.WithdrawalInfo]
	// EthereumNonceRequested - key: bool (is requested)
	EthereumNonceRequested collections.Item[bool]
	// PendingMintTransactions - key: pending zenBTC mint transaction
	PendingMintTransactions collections.Item[treasurytypes.PendingMintTransactions]
	// VoteExtensionRejected - key: bool (is rejected)
	VoteExtensionRejected collections.Item[bool]
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
	treasuryKeeper *treasury.Keeper,
	validatorAddressCodec addresscodec.Codec,
	consensusAddressCodec addresscodec.Codec,
) *Keeper {
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

	oracleClient, err := NewSidecarClient(zrConfig.SidecarAddr)
	if err != nil {
		panic("error creating sidecar client")
	}

	sb := collections.NewSchemaBuilder(storeService)

	return &Keeper{
		storeService:          storeService,
		cdc:                   cdc,
		authKeeper:            ak,
		bankKeeper:            bk,
		hooks:                 nil,
		authority:             authority,
		txDecoder:             txDecoder,
		zrConfig:              zrConfig,
		sidecarClient:         oracleClient,
		treasuryKeeper:        treasuryKeeper,
		validatorAddressCodec: validatorAddressCodec,
		consensusAddressCodec: consensusAddressCodec,

		AVSDelegations:          collections.NewMap(sb, types.AVSDelegationsKey, types.AVSDelegationsIndex, collections.PairKeyCodec(collections.StringKey, collections.StringKey), sdk.IntValue),
		ValidatorDelegations:    collections.NewMap(sb, types.ValidatorDelegationsKey, types.ValidatorDelegationsIndex, collections.StringKey, sdk.IntValue),
		AVSRewardsPool:          collections.NewMap(sb, types.AVSRewardsPoolKey, types.AVSRewardsPoolIndex, collections.StringKey, sdk.IntValue),
		AssetPrices:             collections.NewMap(sb, types.AssetPricesKey, types.AssetPricesIndex, collections.StringKey, codec.CollValue[types.AssetPrice](cdc)),
		SlashEvents:             collections.NewMap(sb, types.SlashEventsKey, types.SlashEventsIndex, collections.Uint64Key, codec.CollValue[types.SlashEvent](cdc)),
		SlashEventCount:         collections.NewItem(sb, types.SlashEventCountKey, types.SlashEventCountIndex, collections.Uint64Value),
		HVParams:                collections.NewItem(sb, types.HVParamsKey, types.HVParamsIndex, codec.CollValue[types.HVParams](cdc)),
		ValidationInfos:         collections.NewMap(sb, types.ValidationInfosKey, types.ValidationInfosIndex, collections.Int64Key, codec.CollValue[types.ValidationInfo](cdc)),
		BtcBlockHeaders:         collections.NewMap(sb, types.BtcBlockHeadersKey, types.BtcBlockHeadersIndex, collections.Int64Key, codec.CollValue[sidecar.BTCBlockHeader](cdc)),
		ConfirmedUnlockTxs:      collections.NewMap(sb, types.ConfirmedUnlockTxsKey, types.ConfirmedUnlockTxsIndex, collections.PairKeyCodec(collections.StringKey, collections.StringKey), codec.CollValue[types.WithdrawalInfo](cdc)),
		EthereumNonceRequested:  collections.NewItem(sb, types.EthereumNonceRequestedKey, types.EthereumNonceRequestedIndex, collections.BoolValue),
		PendingMintTransactions: collections.NewItem(sb, types.PendingMintTransactionsKey, types.PendingMintTransactionsIndex, codec.CollValue[treasurytypes.PendingMintTransactions](cdc)),
		VoteExtensionRejected:   collections.NewItem(sb, types.VoteExtensionRejectedKey, types.VoteExtensionRejectedIndex, collections.BoolValue),
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