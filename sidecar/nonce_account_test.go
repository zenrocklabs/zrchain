package main_test

import (
	"context"
	"encoding/hex"
	"log"
	"testing"

	sidecar "github.com/Zenrock-Foundation/zrchain/v6/sidecar"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	"github.com/stretchr/testify/require"
)

// TestGetSolanaNonceAccount tests fetching the specific nonce account that's failing
func TestGetSolanaNonceAccount(t *testing.T) {
	// Initialize oracle with Solana client
	oracle := initTestOracle()
	service := sidecar.NewOracleService(oracle)
	require.NotNil(t, service)

	// The specific nonce account from your error
	nonceAccountPubKey := "14K6YBmWeL2VwQ84heNvs4r5wwRvC6rida6uLz2TTeSt"

	req := &api.SolanaAccountInfoRequest{
		PubKey: nonceAccountPubKey,
	}

	log.Printf("\n=== Testing Nonce Account: %s ===\n", nonceAccountPubKey)

	// Test 1: Get account info via gRPC service
	resp, err := service.GetSolanaAccountInfo(context.Background(), req)

	if err != nil {
		log.Printf("ERROR: Failed to get account info: %v\n", err)
		t.Fatalf("GetSolanaAccountInfo failed: %v", err)
	}

	require.NotNil(t, resp)
	log.Printf("Response received, account data length: %d bytes\n", len(resp.Account))

	if len(resp.Account) == 0 {
		t.Fatalf("ISSUE FOUND: Account data is empty. Check sidecar logs above for the reason (not found, nil value, etc.)")
	}

	log.Printf("SUCCESS: Received account data, length: %d bytes\n", len(resp.Account))
	log.Printf("Account data (hex): %s\n", hex.EncodeToString(resp.Account))
}
