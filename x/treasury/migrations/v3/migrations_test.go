package v3_test

import (
	"testing"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"
	v3 "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/migrations/v3"
	treasury "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/test-go/testify/require"
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

	require.NoError(t, v3.RejectBadTestnetRequests(ctx, signRequestStore, cdc))

	migratedReq1, err := signRequestStore.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED, migratedReq1.Status)
	require.Equal(t, "data for signing is empty", migratedReq1.RejectReason)

	migratedReq2, err := signRequestStore.Get(ctx, 2)
	require.NoError(t, err)
	require.Equal(t, types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING, migratedReq2.Status)
}
