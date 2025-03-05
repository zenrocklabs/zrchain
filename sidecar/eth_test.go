package main_test

import (
	"context"
	"testing"

	sidecar "github.com/Zenrock-Foundation/zrchain/v5/sidecar"
	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func Test_QueryEthNonceAndBlockHeight(t *testing.T) {
	// Initialize the test oracle
	oracle := initTestOracle()
	require.NotNil(t, oracle, "Oracle should be initialized")

	// Create oracle service for testing
	service := sidecar.NewOracleService(oracle)

	// Create the test request with the specified Ethereum address
	req := &api.LatestEthereumNonceForAccountRequest{
		Address: "0x75F1068e904815398045878A41e4324317c93aE4",
	}

	// Query the nonce
	ctx := context.Background()
	res, err := service.GetLatestEthereumNonceForAccount(ctx, req)

	// Verify the results
	require.NoError(t, err, "Should get nonce without error")
	require.NotNil(t, res, "Response should not be nil")

	// Get the latest block header to retrieve the block height
	latestHeader, err := oracle.EthClient.HeaderByNumber(ctx, nil)
	require.NoError(t, err, "Should get latest block header without error")

	// Log the nonce value and block height for verification
	t.Logf("Current nonce for address %s: %d", req.Address, res.Nonce)
	t.Logf("Latest block height: %d", latestHeader.Number.Uint64())

	// Optionally verify the address directly through the client
	nonce, err := oracle.EthClient.NonceAt(ctx, common.HexToAddress(req.Address), nil)
	require.NoError(t, err, "Should get nonce directly without error")
	require.Equal(t, nonce, res.Nonce, "Nonces should match")
}
