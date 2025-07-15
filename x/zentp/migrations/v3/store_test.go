package v3_test

import (
	"context"
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	appparams "github.com/Zenrock-Foundation/zrchain/v6/app/params"
	validation "github.com/Zenrock-Foundation/zrchain/v6/x/validation/module"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"
	ubermock "go.uber.org/mock/gomock"

	storetypes "cosmossdk.io/store/types"

	"errors"

	v3 "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/migrations/v3"
	zentptestutil "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var (
	mint1 = &types.Bridge{
		Id:                 1,
		Denom:              "urock",
		Creator:            "test-creator",
		SourceAddress:      "test-creator",
		SourceChain:        "test-source-chain",
		DestinationChain:   "test-destination-chain",
		Amount:             1000,
		RecipientAddress:   "test-recipient-address",
		TxId:               123,
		TxHash:             "test-tx-hash",
		State:              types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		BlockHeight:        100,
		AwaitingEventSince: 0,
	}

	mint2 = &types.Bridge{
		Id:                 2,
		Denom:              "urock",
		Creator:            "test-creator",
		SourceAddress:      "test-creator",
		SourceChain:        "test-source-chain",
		DestinationChain:   "test-destination-chain",
		Amount:             2000,
		RecipientAddress:   "test-recipient-address",
		TxId:               456,
		TxHash:             "test-tx-hash-2",
		State:              types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		BlockHeight:        100,
		AwaitingEventSince: 0,
	}

	burn1 = &types.Bridge{
		Id:                 1,
		Denom:              "urock",
		Creator:            "test-creator",
		SourceAddress:      "test-creator",
		SourceChain:        "test-source-chain",
		DestinationChain:   "test-destination-chain",
		Amount:             1000,
		RecipientAddress:   "test-recipient-address",
		TxId:               123,
		TxHash:             "test-tx-hash",
		State:              types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		BlockHeight:        100,
		AwaitingEventSince: 0,
	}

	burn2 = &types.Bridge{
		Id:                 2,
		Denom:              "urock",
		Creator:            "test-creator",
		SourceAddress:      "test-creator",
		SourceChain:        "test-source-chain",
		DestinationChain:   "test-destination-chain",
		Amount:             1000,
		RecipientAddress:   "test-recipient-address",
		TxId:               123,
		TxHash:             "test-tx-hash",
		State:              types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		BlockHeight:        100,
		AwaitingEventSince: 0,
	}
)

func TestUpdateMintStore(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(validation.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)

	kvStoreService := runtime.NewKVStoreService(storeKey)
	sb := collections.NewSchemaBuilder(kvStoreService)
	oldMintsCol := collections.NewMap(sb, types.MintsKey, types.MintsIndex, collections.Uint64Key, codec.CollValue[types.Bridge](cdc))
	err := oldMintsCol.Set(ctx, 1, *mint1)
	require.NoError(t, err)
	err = oldMintsCol.Set(ctx, 2, *mint2)
	require.NoError(t, err)
	newMintsCol := collections.NewMap(sb, types.MintsKey, types.MintsIndex, collections.Uint64Key, codec.CollValue[types.Bridge](cdc))

	require.NoError(t, err)
	require.NoError(t, v3.UpdateMintStore(ctx, oldMintsCol, newMintsCol))

	mint, err := newMintsCol.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, mint1.Id, mint.Id)

	mintStore, err := newMintsCol.Iterate(ctx, nil)
	require.NoError(t, err)
	mints, err := mintStore.Values()
	require.NoError(t, err)
	require.Equal(t, 2, len(mints))

	mint, err = newMintsCol.Get(ctx, 2)
	require.NoError(t, err)
	require.Equal(t, mint2.Id, mint.Id)

	require.Equal(t, 2, len(mints))
}

func TestUpdateBurnStore(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(validation.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)

	kvStoreService := runtime.NewKVStoreService(storeKey)
	sb := collections.NewSchemaBuilder(kvStoreService)
	oldBurnsCol := collections.NewMap(sb, types.BurnsKey, types.BurnsIndex, collections.Uint64Key, codec.CollValue[types.Bridge](cdc))
	err := oldBurnsCol.Set(ctx, 1, *burn1)
	newBurnsCol := collections.NewMap(sb, types.BurnsKey, types.BurnsIndex, collections.Uint64Key, codec.CollValue[types.Bridge](cdc))

	require.NoError(t, err)
	require.NoError(t, v3.UpdateBurnStore(ctx, oldBurnsCol, newBurnsCol))

	burn, err := newBurnsCol.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, burn1.Id, burn.Id)
}

func TestSendZentpFeesToMintModule(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(validation.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)

	kvStoreService := runtime.NewKVStoreService(storeKey)
	sb := collections.NewSchemaBuilder(kvStoreService)
	oldMintsCol := collections.NewMap(sb, types.MintsKey, types.MintsIndex, collections.Uint64Key, codec.CollValue[types.Bridge](cdc))
	err := oldMintsCol.Set(ctx, 1, *mint1)
	require.NoError(t, err)
	err = oldMintsCol.Set(ctx, 2, *mint2)
	require.NoError(t, err)

	ctrl := ubermock.NewController(t)
	defer ctrl.Finish()

	accountKeeper := zentptestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := zentptestutil.NewMockBankKeeper(ctrl)

	zentpModuleAddr := authtypes.NewModuleAddress(types.ModuleName)
	accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(zentpModuleAddr).AnyTimes()

	currentBalance := sdk.NewCoin(appparams.BondDenom, math.NewIntFromUint64(10000)) // 10k urock

	bankKeeper.EXPECT().GetBalance(ctx, zentpModuleAddr, appparams.BondDenom).DoAndReturn(
		func(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
			return currentBalance
		},
	).AnyTimes()

	bankKeeper.EXPECT().SendCoinsFromModuleToModule(
		ctx,
		types.ModuleName,
		types.ZentpCollectorName,
		ubermock.Any(),
	).DoAndReturn(
		func(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error {
			// Update the tracked balance by subtracting the sent amount
			currentBalance = currentBalance.Sub(amt[0])
			return nil
		},
	).Times(1)

	getPendingMints := func(context.Context) ([]*types.Bridge, error) {
		mintStore, err := oldMintsCol.Iterate(ctx, nil)
		if err != nil {
			return nil, err
		}
		mints, err := mintStore.Values()
		if err != nil {
			return nil, err
		}

		var pendingMints []*types.Bridge
		for _, mint := range mints {
			if mint.State == types.BridgeStatus_BRIDGE_STATUS_PENDING {
				pendingMints = append(pendingMints, &mint)
			}
		}
		return pendingMints, nil
	}

	getBridgeFeeParams := func(context.Context) (sdk.AccAddress, math.LegacyDec, error) {
		return zentpModuleAddr, math.LegacyNewDecWithPrec(5, 3), nil // 0.5% bridge fee
	}

	err = v3.SendZentpFeesToMintModule(ctx, getPendingMints, getBridgeFeeParams, bankKeeper, accountKeeper)

	require.NoError(t, err)

	// Final balance should be 0
	require.True(t, currentBalance.IsZero(), "Final balance should be empty after sending all coins")
}

func TestSendZentpFeesToMintModuleWithPendingMints(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(validation.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)

	kvStoreService := runtime.NewKVStoreService(storeKey)
	sb := collections.NewSchemaBuilder(kvStoreService)
	oldMintsCol := collections.NewMap(sb, types.MintsKey, types.MintsIndex, collections.Uint64Key, codec.CollValue[types.Bridge](cdc))

	// Create a pending mint for fee calculation
	pendingMint := types.Bridge{
		Id:                 3,
		Denom:              "urock",
		Creator:            "test-creator",
		SourceAddress:      "test-creator",
		SourceChain:        "test-source-chain",
		DestinationChain:   "test-destination-chain",
		Amount:             10000, // 10k urock
		RecipientAddress:   "test-recipient-address",
		TxId:               789,
		TxHash:             "test-tx-hash-3",
		State:              types.BridgeStatus_BRIDGE_STATUS_PENDING,
		BlockHeight:        100,
		AwaitingEventSince: 0,
	}

	err := oldMintsCol.Set(ctx, 3, pendingMint)
	require.NoError(t, err)

	ctrl := ubermock.NewController(t)
	defer ctrl.Finish()

	accountKeeper := zentptestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := zentptestutil.NewMockBankKeeper(ctrl)

	zentpModuleAddr := authtypes.NewModuleAddress(types.ModuleName)
	accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(zentpModuleAddr).AnyTimes()

	currentBalance := sdk.NewCoin(appparams.BondDenom, math.NewIntFromUint64(10050)) // 10050 urock

	bankKeeper.EXPECT().GetBalance(ctx, zentpModuleAddr, appparams.BondDenom).DoAndReturn(
		func(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
			return currentBalance
		},
	).AnyTimes()

	bankKeeper.EXPECT().SendCoinsFromModuleToModule(
		ctx,
		types.ModuleName,
		types.ZentpCollectorName,
		ubermock.Any(),
	).DoAndReturn(
		func(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error {
			currentBalance = currentBalance.Sub(amt[0])
			return nil
		},
	).Times(1)

	getPendingMints := func(context.Context) ([]*types.Bridge, error) {
		mintStore, err := oldMintsCol.Iterate(ctx, nil)
		if err != nil {
			return nil, err
		}
		mints, err := mintStore.Values()
		if err != nil {
			return nil, err
		}

		var pendingMints []*types.Bridge
		for _, mint := range mints {
			if mint.State == types.BridgeStatus_BRIDGE_STATUS_PENDING {
				pendingMints = append(pendingMints, &mint)
			}
		}
		return pendingMints, nil
	}

	getBridgeFeeParams := func(context.Context) (sdk.AccAddress, math.LegacyDec, error) {
		return zentpModuleAddr, math.LegacyNewDecWithPrec(5, 3), nil
	}

	err = v3.SendZentpFeesToMintModule(ctx, getPendingMints, getBridgeFeeParams, bankKeeper, accountKeeper)
	require.NoError(t, err)

	expectedFinalBalance := sdk.NewCoin(appparams.BondDenom, math.NewIntFromUint64(10000))

	require.Equal(t, expectedFinalBalance, currentBalance, "Final balance should match expected pending mint amount")
}

func TestSendZentpFeesToMintModuleEdgeCases(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(validation.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)

	kvStoreService := runtime.NewKVStoreService(storeKey)
	sb := collections.NewSchemaBuilder(kvStoreService)
	oldMintsCol := collections.NewMap(sb, types.MintsKey, types.MintsIndex, collections.Uint64Key, codec.CollValue[types.Bridge](cdc))

	pendingMint1 := types.Bridge{
		Id:     1,
		Denom:  "urock",
		Amount: 1000,
		State:  types.BridgeStatus_BRIDGE_STATUS_PENDING,
	}
	pendingMint2 := types.Bridge{
		Id:     2,
		Denom:  "urock",
		Amount: 2000,
		State:  types.BridgeStatus_BRIDGE_STATUS_PENDING,
	}

	err := oldMintsCol.Set(ctx, 1, pendingMint1)
	require.NoError(t, err)
	err = oldMintsCol.Set(ctx, 2, pendingMint2)
	require.NoError(t, err)

	ctrl := ubermock.NewController(t)
	defer ctrl.Finish()

	accountKeeper := zentptestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := zentptestutil.NewMockBankKeeper(ctrl)

	zentpModuleAddr := authtypes.NewModuleAddress(types.ModuleName)
	accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(zentpModuleAddr).AnyTimes()

	currentBalance := sdk.NewCoin(appparams.BondDenom, math.NewIntFromUint64(4030))

	bankKeeper.EXPECT().GetBalance(ctx, zentpModuleAddr, appparams.BondDenom).DoAndReturn(
		func(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
			return currentBalance
		},
	).AnyTimes()

	bankKeeper.EXPECT().SendCoinsFromModuleToModule(
		ctx,
		types.ModuleName,
		types.ZentpCollectorName,
		ubermock.Any(),
	).DoAndReturn(
		func(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error {
			currentBalance = currentBalance.Sub(amt[0])
			return nil
		},
	).Times(1)

	getPendingMints := func(context.Context) ([]*types.Bridge, error) {
		mintStore, err := oldMintsCol.Iterate(ctx, nil)
		if err != nil {
			return nil, err
		}
		mints, err := mintStore.Values()
		if err != nil {
			return nil, err
		}

		var pendingMints []*types.Bridge
		for _, mint := range mints {
			if mint.State == types.BridgeStatus_BRIDGE_STATUS_PENDING {
				pendingMints = append(pendingMints, &mint)
			}
		}
		return pendingMints, nil
	}

	getBridgeFeeParams := func(context.Context) (sdk.AccAddress, math.LegacyDec, error) {
		return zentpModuleAddr, math.LegacyNewDecWithPrec(1, 2), nil // 1% bridge fee
	}

	err = v3.SendZentpFeesToMintModule(ctx, getPendingMints, getBridgeFeeParams, bankKeeper, accountKeeper)
	require.NoError(t, err)

	// Verify calculations:
	// Pending mint 1: 1000 urock * 1% = 10 urock fee
	// Pending mint 2: 2000 urock * 1% = 20 urock fee
	// Total pending fees: 30 urock
	// Amount sent: 4030 - 3000 = 1030 urock
	// Final balance: 3000 urock
	expectedFinalBalance := sdk.NewCoin(appparams.BondDenom, math.NewIntFromUint64(3000))
	require.Equal(t, expectedFinalBalance, currentBalance, "Final balance should match total pending amounts without fees")
}

func TestZentpFees(t *testing.T) {

	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)

	kvStoreService := runtime.NewKVStoreService(storeKey)
	sb := collections.NewSchemaBuilder(kvStoreService)
	zentpFees := collections.NewItem(sb, types.ZentpFeesKey, types.ZentpFeesIndex, collections.Uint64Value)

	sampledFees := uint64(0)

	fees, err := zentpFees.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			err = zentpFees.Set(ctx, sampledFees)
			require.NoError(t, err)
		} else {
			require.NoError(t, err)
		}
	}

	fees, err = zentpFees.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, sampledFees, fees)
}
