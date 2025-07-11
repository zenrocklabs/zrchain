package main

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solrock/generated/rock_spl_token"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
)

func TestSolanaEventPerformance_Manual(t *testing.T) {
	t.Skip("Skipping performance tests in CI")
	const numRuns = 3

	testCases := []struct {
		name                         string
		solanaEventScanTxLimit       int
		solanaEventFetchBatchSize    int
		solanaEventFetchMinBatchSize int
		solanaMaxConcurrentRPCCalls  int
		solanaSleepInterval          time.Duration
		solanaFallbackSleepInterval  time.Duration
		solanaEventFetchRetrySleep   time.Duration
		solanaRPCTimeout             time.Duration
		solanaBatchTimeout           time.Duration
	}{
		// {
		// 	name:                         "LongerTimeouts",
		// 	solanaEventScanTxLimit:       440,
		// 	solanaEventFetchBatchSize:    10,
		// 	solanaEventFetchMinBatchSize: 2,
		// 	solanaMaxConcurrentRPCCalls:  20,
		// 	solanaSleepInterval:          50 * time.Millisecond,
		// 	solanaFallbackSleepInterval:  10 * time.Millisecond,
		// 	solanaEventFetchRetrySleep:   100 * time.Millisecond,
		// 	solanaRPCTimeout:             25 * time.Second,
		// 	solanaBatchTimeout:           35 * time.Second,
		// },
		// {
		// 	name:                         "StandardTimeouts",
		// 	solanaEventScanTxLimit:       440,
		// 	solanaEventFetchBatchSize:    10,
		// 	solanaEventFetchMinBatchSize: 2,
		// 	solanaMaxConcurrentRPCCalls:  20,
		// 	solanaSleepInterval:          50 * time.Millisecond,
		// 	solanaFallbackSleepInterval:  10 * time.Millisecond,
		// 	solanaEventFetchRetrySleep:   100 * time.Millisecond,
		// 	solanaRPCTimeout:             20 * time.Second,
		// 	solanaBatchTimeout:           30 * time.Second,
		// },
		// {
		// 	name:                         "ShorterTimeouts",
		// 	solanaEventScanTxLimit:       440,
		// 	solanaEventFetchBatchSize:    10,
		// 	solanaEventFetchMinBatchSize: 2,
		// 	solanaMaxConcurrentRPCCalls:  20,
		// 	solanaSleepInterval:          50 * time.Millisecond,
		// 	solanaFallbackSleepInterval:  10 * time.Millisecond,
		// 	solanaEventFetchRetrySleep:   100 * time.Millisecond,
		// 	solanaRPCTimeout:             15 * time.Second,
		// 	solanaBatchTimeout:           25 * time.Second,
		// },
		{
			name:                         "EvenShorterTimeouts",
			solanaEventScanTxLimit:       440,
			solanaEventFetchBatchSize:    10,
			solanaEventFetchMinBatchSize: 2,
			solanaMaxConcurrentRPCCalls:  20,
			solanaSleepInterval:          50 * time.Millisecond,
			solanaFallbackSleepInterval:  10 * time.Millisecond,
			solanaEventFetchRetrySleep:   100 * time.Millisecond,
			solanaRPCTimeout:             10 * time.Second,
			solanaBatchTimeout:           20 * time.Second,
		},
	}

	cfg := LoadConfig("", "")
	solanaClient := rpc.New(cfg.SolanaRPC[cfg.Network])

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save original values
			originalSolanaEventScanTxLimit := sidecartypes.SolanaEventScanTxLimit
			originalSolanaEventFetchBatchSize := sidecartypes.SolanaEventFetchBatchSize
			originalSolanaEventFetchMinBatchSize := sidecartypes.SolanaEventFetchMinBatchSize
			originalSolanaMaxConcurrentRPCCalls := sidecartypes.SolanaMaxConcurrentRPCCalls
			originalSolanaSleepInterval := sidecartypes.SolanaSleepInterval
			originalSolanaFallbackSleepInterval := sidecartypes.SolanaFallbackSleepInterval
			originalSolanaEventFetchRetrySleep := sidecartypes.SolanaEventFetchRetrySleep
			originalSolanaRPCTimeout := sidecartypes.SolanaRPCTimeout
			originalSolanaBatchTimeout := sidecartypes.SolanaBatchTimeout

			// Defer restoration
			defer func() {
				sidecartypes.SolanaEventScanTxLimit = originalSolanaEventScanTxLimit
				sidecartypes.SolanaEventFetchBatchSize = originalSolanaEventFetchBatchSize
				sidecartypes.SolanaEventFetchMinBatchSize = originalSolanaEventFetchMinBatchSize
				sidecartypes.SolanaMaxConcurrentRPCCalls = originalSolanaMaxConcurrentRPCCalls
				sidecartypes.SolanaSleepInterval = originalSolanaSleepInterval
				sidecartypes.SolanaFallbackSleepInterval = originalSolanaFallbackSleepInterval
				sidecartypes.SolanaEventFetchRetrySleep = originalSolanaEventFetchRetrySleep
				sidecartypes.SolanaRPCTimeout = originalSolanaRPCTimeout
				sidecartypes.SolanaBatchTimeout = originalSolanaBatchTimeout
			}()

			// Apply test case parameters
			sidecartypes.SolanaEventScanTxLimit = tc.solanaEventScanTxLimit
			sidecartypes.SolanaEventFetchBatchSize = tc.solanaEventFetchBatchSize
			sidecartypes.SolanaEventFetchMinBatchSize = tc.solanaEventFetchMinBatchSize
			sidecartypes.SolanaMaxConcurrentRPCCalls = tc.solanaMaxConcurrentRPCCalls
			sidecartypes.SolanaSleepInterval = tc.solanaSleepInterval
			sidecartypes.SolanaFallbackSleepInterval = tc.solanaFallbackSleepInterval
			sidecartypes.SolanaEventFetchRetrySleep = tc.solanaEventFetchRetrySleep
			sidecartypes.SolanaRPCTimeout = tc.solanaRPCTimeout
			sidecartypes.SolanaBatchTimeout = tc.solanaBatchTimeout

			oracle := NewOracle(cfg, nil, nil, solanaClient, nil, false, true)

			programID := sidecartypes.SolRockProgramID[oracle.Config.Network]
			program, err := solana.PublicKeyFromBase58(programID)
			require.NoError(t, err)

			limit := 440
			signatures, err := solanaClient.GetSignaturesForAddressWithOpts(context.Background(), program, &rpc.GetSignaturesForAddressOpts{
				Limit: &limit,
			})
			require.NoError(t, err)
			require.Len(t, signatures, limit, "failed to fetch 440 signatures for performance test")

			// Create the processor function
			processor := func(txResult *rpc.GetTransactionResult, program solana.PublicKey, sig solana.Signature, debugMode bool) ([]any, error) {
				return oracle.processMintTransaction(txResult, program, sig, debugMode,
					func(tx *rpc.GetTransactionResult, prog solana.PublicKey) ([]any, error) {
						events, err := rock_spl_token.DecodeEvents(tx, prog)
						if err != nil {
							return nil, err
						}
						var interfaceEvents []any
						for _, event := range events {
							interfaceEvents = append(interfaceEvents, event)
						}
						return interfaceEvents, nil
					},
					func(data any) (solana.PublicKey, uint64, uint64, solana.PublicKey, bool) {
						eventData, ok := data.(*rock_spl_token.TokensMintedWithFeeEventData)
						if !ok {
							return solana.PublicKey{}, 0, 0, solana.PublicKey{}, false
						}
						return eventData.Recipient, eventData.Value, eventData.Fee, eventData.Mint, true
					},
					"SolRockMint",
				)
			}

			totalStartTime := time.Now()
			var runTimes []time.Duration

			for i := 0; i < numRuns; i++ {
				runStartTime := time.Now()
				_, _, err := oracle.processSignatures(context.Background(), signatures, program, "SolRockMint", processor)
				runDuration := time.Since(runStartTime)
				runTimes = append(runTimes, runDuration)

				t.Logf("--- %s | Run %d/%d: %s ---", tc.name, i+1, numRuns, runDuration.String())

				if err != nil {
					// The oracle is designed to be resilient to transient errors.
					// We only want to fail the test on unrecoverable errors.
					if !strings.Contains(err.Error(), "context deadline exceeded") && !strings.Contains(err.Error(), "unexpected end of JSON input") {
						require.NoError(t, err)
					}
				}
			}
			totalElapsedTime := time.Since(totalStartTime)

			// Calculate average
			var totalRunTime time.Duration
			for _, rt := range runTimes {
				totalRunTime += rt
			}
			avgRunTime := totalRunTime / time.Duration(len(runTimes))

			t.Logf("--- %s | Summary: Total=%s, Average=%s, Runs=%d ---", tc.name, totalElapsedTime.String(), avgRunTime.String(), numRuns)
		})
	}
}
