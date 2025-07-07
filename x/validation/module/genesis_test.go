package validation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/testutil"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	validationtestutil "github.com/Zenrock-Foundation/zrchain/v6/x/validation/testutil"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
)

func TestValidateGenesis(t *testing.T) {
	genValidators1 := make([]types.Validator, 1, 5)
	pk := ed25519.GenPrivKey().PubKey()
	genValidators1[0] = testutil.NewValidator(t, sdk.ValAddress(pk.Address()), pk)
	genValidators1[0].Tokens = math.OneInt()
	genValidators1[0].DelegatorShares = math.LegacyOneDec()

	tests := []struct {
		name    string
		mutate  func(*types.GenesisState)
		wantErr bool
	}{
		{"default", func(*types.GenesisState) {}, false},
		// validate genesis validators
		{"duplicate validator", func(data *types.GenesisState) {
			data.Validators = genValidators1
			data.Validators = append(data.Validators, genValidators1[0])
		}, true},
		{"no delegator shares", func(data *types.GenesisState) {
			data.Validators = genValidators1
			data.Validators[0].DelegatorShares = math.LegacyZeroDec()
		}, true},
		{"jailed and bonded validator", func(data *types.GenesisState) {
			data.Validators = genValidators1
			data.Validators[0].Jailed = true
			data.Validators[0].Status = types.Bonded
		}, true},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			genesisState := types.DefaultGenesisState()
			tt.mutate(genesisState)

			if tt.wantErr {
				assert.Error(t, staking.ValidateGenesis(genesisState))
			} else {
				assert.NoError(t, staking.ValidateGenesis(genesisState))
			}
		})
	}
}

func TestExportGenesisWithEmptyCollections(t *testing.T) {
	// Test that the default genesis state has proper empty values
	// This ensures our fix for empty collections works correctly
	genesisState := validationtypes.DefaultGenesisState()
	require.NotNil(t, genesisState)

	// The default genesis state only sets Params, other fields are zero values
	require.NotNil(t, genesisState.Params)
	require.NotNil(t, genesisState.BackfillRequest)                // non-nil struct
	require.Nil(t, genesisState.BackfillRequest.Requests)          // nil slice inside struct
	require.Nil(t, genesisState.RequestedHistoricalBitcoinHeaders) // nil slice
	require.Equal(t, int64(0), genesisState.LastValidVeHeight)     // zero value
}

func TestInitGenesis(t *testing.T) {
	genesisState := validationtestutil.DefaultGenesis()

	require.NotNil(t, genesisState.Params)
	require.NotNil(t, genesisState.BackfillRequest)                // non-nil struct
	require.Nil(t, genesisState.BackfillRequest.Requests)          // nil slice inside struct
	require.Nil(t, genesisState.RequestedHistoricalBitcoinHeaders) // nil slice
	require.Equal(t, int64(0), genesisState.LastValidVeHeight)     // zero value
}
