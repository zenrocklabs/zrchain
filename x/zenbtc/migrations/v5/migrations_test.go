package v5_test

import (
	"testing"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"
	validation "github.com/Zenrock-Foundation/zrchain/v6/x/validation/module"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"

	v5 "github.com/zenrocklabs/zenbtc/x/zenbtc/migrations/v5"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

// Define expected parameters based on v5/store.go logic
var localParams = types.Params{ // local
	DepositKeyringAddr:  "keyring1k6vc6vhp6e6l3rxalue9v4ux",
	StakerKeyID:         1,
	EthMinterKeyID:      2,
	UnstakerKeyID:       3,
	CompleterKeyID:      4,
	RewardsDepositKeyID: 5,
	ChangeAddressKeyIDs: []uint64{6},
	BitcoinProxyAddress: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
	EthTokenAddr:        "0x7692E9a796001FeE9023853f490A692bAB2E4834",
	ControllerAddr:      "0x2844bd31B68AE5a0335c672e6251e99324441B73",
	Solana: &types.Solana{
		SignerKeyId:        7,
		ProgramId:          "3jo4mdc6QbGRigia2jvmKShbmz3aWq4Y8bgUXfur5StT",
		NonceAuthorityKey:  8,
		NonceAccountKey:    9,
		MintAddress:        "9oBkgQUkq8jvzK98D7Uib6GYSZZmjnZ6QEGJRrAeKnDj",
		FeeWallet:          "FzqGcRG98v1KhKxatX2Abb2z1aJ2rViQwBK5GHByKCAd",
		Fee:                0,
		MultisigKeyAddress: "8cmZY2id22vxpXs2H3YYQNARuPHNuYwa7jipW1q1v9Fy",
		Btl:                20,
	},
}
var amberParams = types.Params{ // devnet
	DepositKeyringAddr:  "keyring1k6vc6vhp6e6l3rxalue9v4ux",
	StakerKeyID:         1,
	EthMinterKeyID:      2,
	UnstakerKeyID:       3,
	CompleterKeyID:      4,
	RewardsDepositKeyID: 5,
	ChangeAddressKeyIDs: []uint64{6},
	BitcoinProxyAddress: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
	EthTokenAddr:        "0x7692E9a796001FeE9023853f490A692bAB2E4834",
	ControllerAddr:      "0x2844bd31B68AE5a0335c672e6251e99324441B73",
	Solana: &types.Solana{
		SignerKeyId:        7,
		ProgramId:          "3jo4mdc6QbGRigia2jvmKShbmz3aWq4Y8bgUXfur5StT",
		NonceAuthorityKey:  8,
		NonceAccountKey:    9,
		MintAddress:        "9oBkgQUkq8jvzK98D7Uib6GYSZZmjnZ6QEGJRrAeKnDj",
		FeeWallet:          "FzqGcRG98v1KhKxatX2Abb2z1aJ2rViQwBK5GHByKCAd",
		Fee:                0,
		MultisigKeyAddress: "8cmZY2id22vxpXs2H3YYQNARuPHNuYwa7jipW1q1v9Fy",
		Btl:                20,
	},
}
var gardiaParams = types.Params{ // testnet
	DepositKeyringAddr:  "keyring1k6vc6vhp6e6l3rxalue9v4ux",
	StakerKeyID:         1,
	EthMinterKeyID:      2,
	UnstakerKeyID:       3,
	CompleterKeyID:      4,
	RewardsDepositKeyID: 5,
	ChangeAddressKeyIDs: []uint64{6},
	BitcoinProxyAddress: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
	EthTokenAddr:        "0xfA32a2D7546f8C7c229F94E693422A786DaE5E18",
	ControllerAddr:      "0xaCE3634AAd9bCC48ef6A194f360F7ACe51F7d9f1",
	Solana: &types.Solana{
		SignerKeyId:        7,
		ProgramId:          "3jo4mdc6QbGRigia2jvmKShbmz3aWq4Y8bgUXfur5StT",
		NonceAuthorityKey:  8,
		NonceAccountKey:    9,
		MintAddress:        "9oBkgQUkq8jvzK98D7Uib6GYSZZmjnZ6QEGJRrAeKnDj",
		FeeWallet:          "FzqGcRG98v1KhKxatX2Abb2z1aJ2rViQwBK5GHByKCAd",
		Fee:                0,
		MultisigKeyAddress: "8cmZY2id22vxpXs2H3YYQNARuPHNuYwa7jipW1q1v9Fy",
		Btl:                20,
	},
}
var diamondParams = types.Params{ // mainnet
	DepositKeyringAddr:  "keyring1k6vc6vhp6e6l3rxalue9v4ux",
	StakerKeyID:         24,
	EthMinterKeyID:      17,
	UnstakerKeyID:       19,
	CompleterKeyID:      28,
	RewardsDepositKeyID: 20,
	ChangeAddressKeyIDs: []uint64{18},
	BitcoinProxyAddress: "zen1mgl98jt30nemuqtt5asldk49ju9lnx0pfke79q",
	EthTokenAddr:        "0x2fE9754d5D28bac0ea8971C0Ca59428b8644C776",
	ControllerAddr:      "0xa87bE298115bE701A12F34F9B4585586dF052008",
	Solana: &types.Solana{
		SignerKeyId:        7,
		ProgramId:          "3jo4mdc6QbGRigia2jvmKShbmz3aWq4Y8bgUXfur5StT",
		NonceAuthorityKey:  8,
		NonceAccountKey:    9,
		MintAddress:        "9oBkgQUkq8jvzK98D7Uib6GYSZZmjnZ6QEGJRrAeKnDj",
		FeeWallet:          "FzqGcRG98v1KhKxatX2Abb2z1aJ2rViQwBK5GHByKCAd",
		Fee:                0,
		MultisigKeyAddress: "8cmZY2id22vxpXs2H3YYQNARuPHNuYwa7jipW1q1v9Fy",
		Btl:                20,
	},
}
var emptyParams = types.Params{}

func TestMigrate(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(validation.AppModuleBasic{})
	cdc := encCfg.Codec
	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	kvStoreService := runtime.NewKVStoreService(storeKey)

	testCases := []struct {
		name           string
		chainID        string
		initialParams  types.Params // Can be used to test migration from specific states if needed
		expectedParams types.Params
	}{
		{
			name:           "empty chainID",
			chainID:        "", // Defaults to "zenrock" in UpdateParams
			initialParams:  types.Params{},
			expectedParams: localParams,
		},
		{
			name:           "zenrock chainID",
			chainID:        "zenrock",
			initialParams:  types.Params{},
			expectedParams: localParams,
		},
		{
			name:           "amber chainID",
			chainID:        "amber-1",
			initialParams:  types.Params{},
			expectedParams: amberParams,
		},
		{
			name:           "gardia chainID",
			chainID:        "gardia-3",
			initialParams:  types.Params{},
			expectedParams: gardiaParams,
		},
		{
			name:           "diamond chainID",
			chainID:        "diamond-1",
			initialParams:  types.Params{},
			expectedParams: diamondParams,
		},
		{
			name:           "unknown chainID prefix",
			chainID:        "unknown-0",
			initialParams:  types.Params{},
			expectedParams: emptyParams, // Expect empty params if no prefix matches
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testutil.DefaultContext(storeKey, tKey)
			// Set the ChainID for the context
			ctx = ctx.WithChainID(tc.chainID)

			sb := collections.NewSchemaBuilder(kvStoreService)
			paramsCol := collections.NewItem(sb, types.ParamsKey, types.ParamsIndex, codec.CollValue[types.Params](cdc))

			// Set initial state if needed (currently just empty)
			err := paramsCol.Set(ctx, tc.initialParams)
			require.NoError(t, err)

			// Run the migration
			err = v5.UpdateParams(ctx, paramsCol)
			require.NoError(t, err)

			// Verify the updated params
			finalParams, err := paramsCol.Get(ctx)
			require.NoError(t, err)
			require.Equal(t, tc.expectedParams, finalParams)

			// Optional: Double check raw bytes if necessary (less ideal than checking via collection)
			store := kvStoreService.OpenKVStore(ctx)
			var res types.Params
			bz, err := store.Get(types.ParamsKey)
			require.NoError(t, err)
			require.NoError(t, cdc.Unmarshal(bz, &res))
			require.Equal(t, tc.expectedParams, res)
		})
	}
}
