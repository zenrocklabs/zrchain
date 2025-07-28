package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	idTypes "github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	treasuryTypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
)

type (
	Keeper struct {
		cdc             codec.BinaryCodec
		storeService    store.KVStoreService
		memStoreService store.MemoryStoreService
		logger          log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority        string
		treasuryKeeper   types.TreasuryKeeper
		bankKeeper       types.BankKeeper
		accountKeeper    types.AccountKeeper
		identityKeeper   types.IdentityKeeper
		validationKeeper types.ValidationKeeper
		mintKeeper       types.MintKeeper
		mintStore        collections.Map[uint64, types.Bridge] // DEPRECATED
		MintStore        collections.Map[uint64, types.Bridge]
		MintCount        collections.Item[uint64]
		burnStore        collections.Map[uint64, types.Bridge] // DEPRECATED
		BurnStore        collections.Map[uint64, types.Bridge]
		BurnCount        collections.Item[uint64]
		ParamStore       collections.Item[types.Params]
		SolanaROCKSupply collections.Item[math.Int]
		ZentpFees        collections.Item[uint64]
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
	validationKeeper types.ValidationKeeper,
	mintKeeper types.MintKeeper,
	memStoreService store.MemoryStoreService,
	testMode bool,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}
	// ensure mint module account is set
	if !testMode {
		if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
			panic(fmt.Sprintf("the x/%s module account has not been set", types.ModuleName))
		}
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:              cdc,
		storeService:     storeService,
		memStoreService:  memStoreService,
		MintStore:        collections.NewMap(sb, types.MintsKey, types.MintsIndex, collections.Uint64Key, codec.CollValue[types.Bridge](cdc)),
		BurnStore:        collections.NewMap(sb, types.BurnsKey, types.BurnsIndex, collections.Uint64Key, codec.CollValue[types.Bridge](cdc)),
		mintStore:        collections.NewMap(sb, types.MintsKeyDeprecated, types.MintsIndexDeprecated, collections.Uint64Key, codec.CollValue[types.Bridge](cdc)),
		burnStore:        collections.NewMap(sb, types.BurnsKeyDeprecated, types.BurnsIndexDeprecated, collections.Uint64Key, codec.CollValue[types.Bridge](cdc)),
		MintCount:        collections.NewItem(sb, types.MintCountKey, types.MintCountIndex, collections.Uint64Value),
		BurnCount:        collections.NewItem(sb, types.BurnCountKey, types.BurnCountIndex, collections.Uint64Value),
		ParamStore:       collections.NewItem(sb, types.ParamsKey, types.ParamsIndex, codec.CollValue[types.Params](cdc)),
		SolanaROCKSupply: collections.NewItem(sb, types.SolanaROCKSupplyKey, types.SolanaROCKSupplyIndex, sdk.IntValue),
		authority:        authority,
		logger:           logger,
		treasuryKeeper:   treasuryKeeper,
		bankKeeper:       bankKeeper,
		accountKeeper:    accountKeeper,
		identityKeeper:   identityKeeper,
		validationKeeper: validationKeeper,
		mintKeeper:       mintKeeper,
		ZentpFees:        collections.NewItem(sb, types.ZentpFeesKey, types.ZentpFeesIndex, collections.Uint64Value),
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
	resp, err := k.identityKeeper.Workspaces(goCtx, &idTypes.QueryWorkspacesRequest{Creator: user, Owner: user})
	if err != nil {
		k.Logger().Error("failed to get workspaces for user: "+user, err.Error())
		return false
	}

	for _, ws := range resp.Workspaces {
		if key.WorkspaceAddr == ws.Address {
			return true
		}
	}

	return false
}

func (k Keeper) GetMints(goCtx context.Context, address string, chainID string) ([]*types.Bridge, error) {
	mints, _, err := query.CollectionFilteredPaginate(
		goCtx,
		k.MintStore,
		nil,
		func(key uint64, value types.Bridge) (bool, error) {
			return value.SourceAddress == address &&
				value.DestinationChain == chainID, nil
		},
		func(key uint64, value types.Bridge) (*types.Bridge, error) {
			return &value, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return mints, nil
}

func (k Keeper) GetBurns(goCtx context.Context, address, chainID, txHash string) ([]*types.Bridge, error) {
	burns, _, err := query.CollectionFilteredPaginate(
		goCtx,
		k.BurnStore,
		nil,
		func(key uint64, value types.Bridge) (bool, error) {
			addressFilter := address == "" || address == value.RecipientAddress
			chainFilter := chainID == "" || chainID == value.SourceChain
			txHashFilter := txHash == "" || txHash == value.TxHash
			return addressFilter && chainFilter && txHashFilter, nil
		},
		func(key uint64, value types.Bridge) (*types.Bridge, error) {
			return &value, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return burns, nil
}

func (k Keeper) GetSignerKeyID(ctx context.Context) uint64 {
	return k.GetSolanaParams(ctx).SignerKeyId
}

func (k Keeper) GetNonceAccountKey(ctx context.Context) uint64 {
	return k.GetSolanaParams(ctx).NonceAccountKey
}

func (k Keeper) GetNonceAuthorityKey(ctx context.Context) uint64 {
	return k.GetSolanaParams(ctx).NonceAuthorityKey
}

func (k Keeper) GetMintsWithStatus(goCtx context.Context, status types.BridgeStatus) ([]*types.Bridge, error) {

	mints, _, err := query.CollectionFilteredPaginate(
		goCtx,
		k.MintStore,
		nil,
		func(key uint64, value types.Bridge) (bool, error) {
			return value.State == status, nil
		},
		func(key uint64, value types.Bridge) (*types.Bridge, error) {
			return &value, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return mints, nil
}

func (k Keeper) GetMintsWithStatusPending(goCtx context.Context) ([]*types.Bridge, error) {
	lastCompletedZentpMint, err := k.validationKeeper.GetLastCompletedZentpMintID(goCtx)
	if err != nil {
		return nil, err
	}

	startKey := lastCompletedZentpMint
	queryRange := &collections.Range[uint64]{}
	pendingMints := []*types.Bridge{}

	k.MintStore.Walk(goCtx, queryRange.StartExclusive(startKey), func(key uint64, value types.Bridge) (bool, error) {
		if value.State == types.BridgeStatus_BRIDGE_STATUS_PENDING {
			pendingMints = append(pendingMints, &value)
		}
		return true, nil
	})

	return pendingMints, nil
}

func (k Keeper) UpdateMint(ctx context.Context, id uint64, mint *types.Bridge) error {
	return k.MintStore.Set(ctx, id, *mint)
}

func (k Keeper) UpdateBurn(ctx context.Context, id uint64, burn *types.Bridge) error {
	return k.BurnStore.Set(ctx, id, *burn)
}

func (k Keeper) GetSolanaROCKSupply(ctx context.Context) (math.Int, error) {
	supply, err := k.SolanaROCKSupply.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return math.ZeroInt(), nil
		}
		return math.Int{}, err
	}
	return supply, nil
}

func (k Keeper) SetSolanaROCKSupply(ctx context.Context, supply math.Int) error {
	return k.SolanaROCKSupply.Set(ctx, supply)
}

func (k Keeper) GetTotalROCKSupply(ctx sdk.Context) (math.Int, error) {
	solanaSupply, err := k.GetSolanaROCKSupply(ctx)
	if err != nil {
		return math.Int{}, fmt.Errorf("failed to get solana rock supply: %w", err)
	}
	zrchainSupply := k.bankKeeper.GetSupply(ctx, params.BondDenom).Amount
	return zrchainSupply.Add(solanaSupply), nil
}

func (k Keeper) AddBurn(ctx context.Context, burn *types.Bridge) error {
	burnID, err := k.BurnCount.Get(ctx)
	if err != nil {
		return err
	}

	if err := k.BurnCount.Set(ctx, burnID+1); err != nil {
		return err
	}

	return k.BurnStore.Set(ctx, burnID, *burn)
}

func (k Keeper) GetBridgeFeeParams(ctx context.Context) (sdk.AccAddress, math.LegacyDec, error) {
	mintParams, err := k.mintKeeper.GetParams(ctx)
	if err != nil {
		return nil, math.LegacyDec{}, err
	}
	protocolWalletAddress := sdk.MustAccAddressFromBech32(mintParams.ProtocolWalletAddress)

	params, err := k.ParamStore.Get(ctx)
	if err != nil {
		return nil, math.LegacyDec{}, err
	}

	bridgeFee := params.BridgeFee

	return protocolWalletAddress, bridgeFee, nil
}

func (k Keeper) GetBridgeFeeAmount(ctx context.Context, amount uint64, bridgeFee math.LegacyDec) (sdk.Coins, error) {
	amountInt := math.NewIntFromUint64(amount)
	bridgeFeeAmount := math.LegacyNewDecFromInt(amountInt).Mul(bridgeFee).TruncateInt()

	bridgeFeeCoins := sdk.NewCoins(sdk.NewCoin(params.BondDenom, bridgeFeeAmount))
	bridgeAmount := sdk.NewCoins(sdk.NewCoin(params.BondDenom, amountInt.Sub(bridgeFeeAmount)))

	if bridgeFeeCoins.AmountOf(params.BondDenom).Add(bridgeAmount.AmountOf(params.BondDenom)).GT(amountInt) {
		return nil, fmt.Errorf("bridge fee %s and bridge amount %s cannot exceed original amount: %s", bridgeFeeCoins.String(), bridgeAmount.String(), amountInt.String())
	}

	return bridgeFeeCoins, nil
}

func (k Keeper) AddFeeToBridgeAmount(ctx context.Context, amount uint64) (uint64, error) {
	_, bridgeFee, err := k.GetBridgeFeeParams(ctx)
	if err != nil {
		return 0, err
	}

	amountInt := math.NewIntFromUint64(amount)
	bridgeFeeAmount := math.LegacyNewDecFromInt(amountInt).Mul(bridgeFee).TruncateInt()
	totalAmount := bridgeFeeAmount.Add(amountInt)

	if !totalAmount.IsUint64() {
		return 0, fmt.Errorf("total amount %s exceeds max uint64", totalAmount.String())
	}

	return totalAmount.Uint64(), nil
}

func (k Keeper) UpdateZentpFees(ctx context.Context, fees uint64) error {
	zentpFees, err := k.ZentpFees.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return k.ZentpFees.Set(ctx, fees)
		}
		return err
	}

	return k.ZentpFees.Set(ctx, zentpFees+fees)
}
