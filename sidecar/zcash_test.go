package main_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/bitcoin"
	sidecar "github.com/Zenrock-Foundation/zrchain/v6/sidecar"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	"github.com/stretchr/testify/require"
)

// TestGetZcashBlockHeaderByHeight tests fetching a specific Zcash block header by height
func TestGetZcashBlockHeaderByHeight(t *testing.T) {
	// t.Skip("Skipping test on CI")

	oracle := initTestOracle()
	service := sidecar.NewOracleService(oracle)
	require.NotNil(t, service)

	// Test with a known block height
	testHeight := int64(3623412)
	expectedBlockHash := "0023f6c253e30d63b74b5d47879aa6a56cad0d2858bedc6b4c6d5460ee34a0be"

	req := &api.BitcoinBlockHeaderByHeightRequest{
		BlockHeight: testHeight,
		ChainName:   "zcash",
	}
	resp, err := service.GetZcashBlockHeaderByHeight(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.BlockHeader)
	require.Equal(t, testHeight, resp.BlockHeight)
	require.Equal(t, expectedBlockHash, resp.BlockHeader.BlockHash)

	fmt.Printf("\n=== ZCash Block Header at Height %d ===\n", testHeight)
	fmt.Printf("Block Hash: %s\n", resp.BlockHeader.BlockHash)
	fmt.Printf("Previous Block: %s\n", resp.BlockHeader.PrevBlock)
	fmt.Printf("Merkle Root: %s\n", resp.BlockHeader.MerkleRoot)
	fmt.Printf("Block Commitments: %s\n", resp.BlockHeader.BlockCommitments)
	fmt.Printf("Timestamp: %d\n", resp.BlockHeader.TimeStamp)
	fmt.Printf("Version: %d\n", resp.BlockHeader.Version)
	fmt.Printf("Bits: %d\n", resp.BlockHeader.Bits)
	fmt.Printf("Nonce (int64): %d\n", resp.BlockHeader.Nonce)
	fmt.Printf("NonceHex: %s\n", resp.BlockHeader.NonceHex)
	fmt.Printf("========================================\n\n")
}

// TestVerifyZcashBlockHeader tests that the Zcash block header hash verification works correctly
func TestVerifyZcashBlockHeader(t *testing.T) {
	// t.Skip("Skipping test on CI")

	oracle := initTestOracle()
	service := sidecar.NewOracleService(oracle)
	require.NotNil(t, service)

	// Get a known Zcash block header
	testHeight := int64(3623412)
	expectedBlockHash := "0023f6c253e30d63b74b5d47879aa6a56cad0d2858bedc6b4c6d5460ee34a0be"

	req := &api.BitcoinBlockHeaderByHeightRequest{
		BlockHeight: testHeight,
		ChainName:   "zcash",
	}
	resp, err := service.GetZcashBlockHeaderByHeight(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.BlockHeader)

	// Verify the block header using CheckBlockHeader
	err = bitcoin.CheckBlockHeader(resp.BlockHeader)
	require.NoError(t, err, "Block header verification should pass")

	// Verify the block hash matches expected
	require.Equal(t, expectedBlockHash, resp.BlockHeader.BlockHash)

	fmt.Printf("\n=== Zcash Block Header Verification ===\n")
	fmt.Printf("Block Height: %d\n", testHeight)
	fmt.Printf("Block Hash: %s\n", resp.BlockHeader.BlockHash)
	fmt.Printf("Verification: ✓ PASSED\n")
	fmt.Printf("=====================================\n\n")
}

// TestVerifyZcashDepositBlockInclusion tests the full deposit verification flow
func TestVerifyZcashDepositBlockInclusion(t *testing.T) {
	// t.Skip("Skipping test on CI")

	oracle := initTestOracle()
	service := sidecar.NewOracleService(oracle)
	require.NotNil(t, service)

	// Use the coinbase transaction from block 3623412
	testHeight := int64(3623412)
	testTxID := "dc7cfca537414cbf45b568f7afad4ccdd402b8455fbdf9b5fed26e96b3fe4a83"
	testIndex := 0 // Coinbase is always first transaction

	// Get the block header
	headerReq := &api.BitcoinBlockHeaderByHeightRequest{
		BlockHeight: testHeight,
		ChainName:   "zcash",
	}
	headerResp, err := service.GetZcashBlockHeaderByHeight(context.Background(), headerReq)
	require.NoError(t, err)
	require.NotNil(t, headerResp)

	// For this test, we'll create a simple merkle proof
	// In a real scenario, this would come from BuildProof function
	// For a single-transaction block (coinbase only), proof would be empty
	proof := []string{} // Empty proof for coinbase in single-tx block

	// Get the raw transaction (coinbase transaction from block 3623412)
	rawTx := "050000800a27a726f04dec4d00000000f4493700010000000000000000000000000000000000000000000000000000000000000000ffffffff0603f449370106ffffffff0240597307000000001976a91402b7b5b3afa00d56eb1a2e76b8db889c935c3f1088ac20bcbe000000000017a9147a86d6c7eb12ce0aa309d7391a6f338eba3c242b87000000"

	// Test block header verification first
	err = bitcoin.CheckBlockHeader(headerResp.BlockHeader)
	require.NoError(t, err, "Block header should be valid")

	// Test merkle proof verification
	outputs, calculatedTxID, err := bitcoin.VerifyBTCLockTransaction(
		rawTx,
		"zcashtestnet",
		testIndex,
		proof,
		headerResp.BlockHeader,
		[]string{}, // no ignore addresses for this test
	)

	require.NoError(t, err, "Deposit verification should succeed")
	require.NotNil(t, outputs, "Should return transaction outputs")
	require.Equal(t, testTxID, calculatedTxID, "Calculated TXID should match expected")

	fmt.Printf("\n=== Zcash Deposit Verification Test ===\n")
	fmt.Printf("Block Height: %d\n", testHeight)
	fmt.Printf("Transaction ID: %s\n", testTxID)
	fmt.Printf("Calculated TXID: %s\n", calculatedTxID)
	fmt.Printf("Number of outputs: %d\n", len(outputs))
	fmt.Printf("Block Header Hash: %s\n", headerResp.BlockHeader.BlockHash)
	fmt.Printf("Block Header Verification: ✓ PASSED\n")
	fmt.Printf("Merkle Proof Verification: ✓ PASSED\n")
	fmt.Printf("======================================\n\n")
}
