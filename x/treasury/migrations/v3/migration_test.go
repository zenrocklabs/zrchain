package v3_test

import (
	"testing"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"
	v3 "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/migrations/v3"
	treasury "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"
)

func TestMigrate(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(treasury.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)

	kvStoreService := runtime.NewKVStoreService(storeKey)
	sb := collections.NewSchemaBuilder(kvStoreService)

	signRequestStore := collections.NewMap(sb, types.SignRequestsKey, types.SignRequestsIndex, collections.Uint64Key, codec.CollValue[types.SignRequest](cdc))
	keyRequestStore := collections.NewMap(sb, types.KeyRequestsKey, types.KeyRequestsIndex, collections.Uint64Key, codec.CollValue[types.KeyRequest](cdc))

	// Setup sign requests
	req1 := types.SignRequest{
		Id:             1,
		DataForSigning: [][]byte{},
		Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
		Creator:        "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		KeyId:          1,
	}

	req2 := types.SignRequest{
		Id:             2,
		DataForSigning: [][]byte{[]byte("to_be_signed_over")},
		Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
		Creator:        "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		KeyId:          1,
	}

	err := signRequestStore.Set(ctx, req1.Id, req1)
	require.NoError(t, err)
	err = signRequestStore.Set(ctx, req2.Id, req2)
	require.NoError(t, err)

	// Setup key requests
	keyReq1 := types.KeyRequest{
		Id:            1,
		Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
		Creator:       "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		KeyringAddr:   "keyring1",
		WorkspaceAddr: "workspace1",
	}

	keyReq2 := types.KeyRequest{
		Id:            2,
		Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PARTIAL,
		Creator:       "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		KeyringAddr:   "keyring2",
		WorkspaceAddr: "workspace2",
	}

	keyReq3 := types.KeyRequest{
		Id:            3,
		Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
		Creator:       "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		KeyringAddr:   "keyring3",
		WorkspaceAddr: "workspace3",
	}

	err = keyRequestStore.Set(ctx, keyReq1.Id, keyReq1)
	require.NoError(t, err)
	err = keyRequestStore.Set(ctx, keyReq2.Id, keyReq2)
	require.NoError(t, err)
	err = keyRequestStore.Set(ctx, keyReq3.Id, keyReq3)
	require.NoError(t, err)

	// Run migration
	require.NoError(t, v3.RejectBadTestnetRequests(ctx, signRequestStore, keyRequestStore, cdc))

	// Check sign requests after migration
	migratedReq1, err := signRequestStore.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED, migratedReq1.Status)
	require.Equal(t, "rejected by migration", migratedReq1.RejectReason)

	migratedReq2, err := signRequestStore.Get(ctx, 2)
	require.NoError(t, err)
	require.Equal(t, types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED, migratedReq2.Status)
	require.Equal(t, "rejected by migration", migratedReq2.RejectReason)

	// Check key requests after migration
	migratedKeyReq1, err := keyRequestStore.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, types.KeyRequestStatus_KEY_REQUEST_STATUS_REJECTED, migratedKeyReq1.Status)
	require.Equal(t, "rejected by migration", migratedKeyReq1.RejectReason)

	migratedKeyReq2, err := keyRequestStore.Get(ctx, 2)
	require.NoError(t, err)
	require.Equal(t, types.KeyRequestStatus_KEY_REQUEST_STATUS_REJECTED, migratedKeyReq2.Status)
	require.Equal(t, "rejected by migration", migratedKeyReq2.RejectReason)

	// Check that fulfilled key requests are not affected
	migratedKeyReq3, err := keyRequestStore.Get(ctx, 3)
	require.NoError(t, err)
	require.Equal(t, types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED, migratedKeyReq3.Status)
	require.Empty(t, migratedKeyReq3.RejectReason)
}
