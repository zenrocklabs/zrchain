package main_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"testing"

	sidecar "github.com/Zenrock-Foundation/zrchain/v6/sidecar"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
)

// TestGetSolanaNonceAccount tests fetching the specific nonce account that's failing
func TestGetSolanaNonceAccount(t *testing.T) {
	t.Skip() // Don't run on CI
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

// TestVerifyZenBTCAccounts verifies all accounts needed for zenBTC wrap transaction
func TestVerifyZenBTCAccounts(t *testing.T) {
	t.Skip() // Don't run on CI
	// Configuration from params
	mintAddress := "nRy9PYAWC6vYQcwGk44S8kqRaNs56noWzJjfg1mY4fx"
	feeWallet := "FzqGcRG98v1KhKxatX2Abb2z1aJ2rViQwBK5GHByKCAd"
	multisigKeyAddress := "BDEF27kbv73TRrhN7uSDPq4ndQhMDnnZvLGdu4qe889k"
	programID := "BoPDAvu4Q3JjFzQKjHru6BmJLbEYMBxuKqy8pHVTA7A3"
	signerAddress := "4GCX9fgq9gzBH282tVMnwebAW6gW8QX4N7JzGD1djYWT" // Your signer

	// Example recipient - replace with actual recipient from your transaction
	recipient := "FzqGcRG98v1KhKxatX2Abb2z1aJ2rViQwBK5GHByKCAd" // Using fee wallet as example

	// Parse keys
	mintKey, err := solana.PublicKeyFromBase58(mintAddress)
	require.NoError(t, err)

	feeKey, err := solana.PublicKeyFromBase58(feeWallet)
	require.NoError(t, err)

	multisigKey, err := solana.PublicKeyFromBase58(multisigKeyAddress)
	require.NoError(t, err)

	recipientKey, err := solana.PublicKeyFromBase58(recipient)
	require.NoError(t, err)

	programKey, err := solana.PublicKeyFromBase58(programID)
	require.NoError(t, err)

	signerKey, err := solana.PublicKeyFromBase58(signerAddress)
	require.NoError(t, err)

	// Derive ATAs
	feeWalletAta, _, err := solana.FindAssociatedTokenAddress(feeKey, mintKey)
	require.NoError(t, err)
	fmt.Printf("Fee Wallet ATA: %s\n", feeWalletAta.String())

	receiverAta, _, err := solana.FindAssociatedTokenAddress(recipientKey, mintKey)
	require.NoError(t, err)
	fmt.Printf("Receiver ATA: %s\n", receiverAta.String())

	// Derive global_config PDA
	globalConfigPDA, _, err := solana.FindProgramAddress(
		[][]byte{[]byte("global_config")},
		programKey,
	)
	require.NoError(t, err)
	fmt.Printf("Global Config PDA: %s\n", globalConfigPDA.String())

	// Connect to devnet
	client := rpc.New(rpc.DevNet_RPC)
	ctx := context.Background()

	// Check mint account
	t.Run("Mint account exists", func(t *testing.T) {
		info, err := client.GetAccountInfo(ctx, mintKey)
		require.NoError(t, err)
		require.NotNil(t, info.Value)
		fmt.Printf("✓ Mint account exists: %s (owner: %s)\n", mintKey, info.Value.Owner)
	})

	// Check fee wallet account
	t.Run("Fee wallet exists", func(t *testing.T) {
		info, err := client.GetAccountInfo(ctx, feeKey)
		require.NoError(t, err)
		require.NotNil(t, info.Value)
		fmt.Printf("✓ Fee wallet exists: %s (balance: %d lamports)\n", feeKey, info.Value.Lamports)
	})

	// Check multisig account
	t.Run("Multisig account exists", func(t *testing.T) {
		info, err := client.GetAccountInfo(ctx, multisigKey)
		require.NoError(t, err)
		require.NotNil(t, info.Value)
		fmt.Printf("✓ Multisig account exists: %s (owner: %s, size: %d bytes)\n",
			multisigKey, info.Value.Owner, len(info.Value.Data.GetBinary()))
	})

	// Check program account
	t.Run("Program account exists and is executable", func(t *testing.T) {
		info, err := client.GetAccountInfo(ctx, programKey)
		require.NoError(t, err)
		require.NotNil(t, info.Value)
		require.True(t, info.Value.Executable, "Program should be executable")
		fmt.Printf("✓ Program account exists: %s (executable: %v)\n", programKey, info.Value.Executable)
	})

	// Check fee wallet ATA - THIS IS CRITICAL
	t.Run("Fee wallet ATA", func(t *testing.T) {
		info, err := client.GetAccountInfo(ctx, feeWalletAta)
		if err != nil || info.Value == nil {
			fmt.Printf("Fee wallet ATA does NOT exist: %s\n", feeWalletAta)
			fmt.Printf("\nTo fix, run:\n")
			fmt.Printf("spl-token create-account %s --owner %s --url devnet\n\n", mintKey, feeKey)
			t.Errorf("Fee wallet ATA does not exist")
		} else {
			fmt.Printf("✓ Fee wallet ATA exists: %s\n", feeWalletAta)
		}
	})

	// Check receiver ATA
	t.Run("Receiver ATA", func(t *testing.T) {
		info, err := client.GetAccountInfo(ctx, receiverAta)
		if err != nil || info.Value == nil {
			fmt.Printf("✗ Receiver ATA does NOT exist: %s\n", receiverAta)
			fmt.Printf("  Note: This should be created by fundReceiver flag or already exist\n")
		} else {
			fmt.Printf("✓ Receiver ATA exists: %s\n", receiverAta)
		}
	})

	// Check global_config PDA - CRITICAL
	t.Run("Global Config PDA", func(t *testing.T) {
		info, err := client.GetAccountInfo(ctx, globalConfigPDA)
		if err != nil || info.Value == nil {
			fmt.Printf("\n✗✗✗ CRITICAL: Global Config PDA does NOT exist ✗✗✗\n")
			fmt.Printf("Global Config PDA: %s\n", globalConfigPDA)
			fmt.Printf("The program must be initialized first!\n")
			t.Errorf("Global Config PDA does not exist - program not initialized")
		} else {
			fmt.Printf("✓ Global Config PDA exists: %s (owner: %s, size: %d bytes)\n",
				globalConfigPDA, info.Value.Owner, len(info.Value.Data.GetBinary()))
		}
	})

	// Check signer has balance
	t.Run("Signer account", func(t *testing.T) {
		info, err := client.GetAccountInfo(ctx, signerKey)
		if err != nil || info.Value == nil {
			fmt.Printf("✗ Signer account does NOT exist: %s\n", signerKey)
		} else {
			fmt.Printf("✓ Signer account exists: %s (balance: %d lamports)\n", signerKey, info.Value.Lamports)
		}
	})
}
