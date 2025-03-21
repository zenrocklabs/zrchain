package keeper

import (
	"testing"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	identitykeeper "github.com/Zenrock-Foundation/zrchain/v6/x/identity/keeper"
	policykeeper "github.com/Zenrock-Foundation/zrchain/v6/x/policy/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
)

func TreasuryKeeper(t testing.TB, policyKeeper *policykeeper.Keeper, identityKeeper *identitykeeper.Keeper, bankKeeper types.BankKeeper, db dbm.DB, stateStore storetypes.CommitMultiStore) (keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(registry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	k := keeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(storeKey),
		log.NewNopLogger(),
		authority.String(),
		bankKeeper,
		*identityKeeper,
		*policyKeeper,

		nil,
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.ParamStore.Set(ctx, types.DefaultParams())

	return k, ctx
}
