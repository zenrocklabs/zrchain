package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTestGenesis(t *testing.T) {
	genesis := TestGenesis()

	// Verify basic structure
	require.NotNil(t, genesis)
	require.NotNil(t, genesis.Params)
	require.NotNil(t, genesis.Validators)
	require.Len(t, genesis.Validators, 1)
	require.NotNil(t, genesis.Delegations)
	require.Len(t, genesis.Delegations, 1)
	require.NotNil(t, genesis.LastValidatorPowers)
	require.Len(t, genesis.LastValidatorPowers, 1)

	// Verify validator data
	validator := genesis.Validators[0]
	require.Equal(t, "zenvaloper126hek6zagmp3jqf97x7pq7c0j9jqs0ndvcepy6", validator.OperatorAddress)
	require.Equal(t, "zenrock", validator.Description.Moniker)
	require.Equal(t, int64(125000000000000), validator.TokensNative.Int64())

	// Verify delegation data
	delegation := genesis.Delegations[0]
	require.Equal(t, "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq", delegation.DelegatorAddress)
	require.Equal(t, "zenvaloper126hek6zagmp3jqf97x7pq7c0j9jqs0ndvcepy6", delegation.ValidatorAddress)

	// Verify validation infos (should have 3 entries)
	require.Len(t, genesis.ValidationInfos, 3)

	// Verify other fields
	require.True(t, genesis.Exported)
	require.NotNil(t, genesis.HVParams)
	require.Len(t, genesis.AssetPrices, 3)
	require.Len(t, genesis.AvsRewardsPool, 1)
	require.Len(t, genesis.SolanaNonceRequested, 2)
	require.Equal(t, uint64(12), genesis.SolanaNonceRequested[0])
}
