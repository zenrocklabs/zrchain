package main

import (
	"context"
	"math/big"
	"testing"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solzenbtc"
	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solzenbtc/generated/zenbtc_spl_token"
	"github.com/Zenrock-Foundation/zrchain/v6/go-client"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/eventstore"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

// TestSolanaCounterParity compares the Solana program counters with the on-chain counters.
//
// This test requires live network access. Set RUN_COUNTER_CHECK=1 to run it manually.
func TestSolanaCounterParity(t *testing.T) {
	t.Skip("skipping on CI")

	cfg := LoadConfig("", "")
	rpcURL := cfg.SolanaRPC[cfg.Network]
	if rpcURL == "" {
		t.Skipf("no Solana RPC configured for network %s", cfg.Network)
	}

	ctx := context.Background()
	solClient := solanarpc.New(rpcURL)

	programID := sidecartypes.ZenZECSolanaProgramID[cfg.Network]
	if programID == "" {
		t.Skip("no zenZEC program configured for this network")
	}

	solMintCounter := fetchSolanaMintCounter(t, ctx, solClient, programID)
	expectedChain := new(big.Int).Set(solMintCounter)
	if expectedChain.Sign() > 0 {
		expectedChain.Sub(expectedChain, big.NewInt(1))
	}

	zrConn, err := client.NewClientConn(cfg.ZRChainRPC, false)
	require.NoError(t, err)
	defer zrConn.Close()

	validationClient := client.NewValidationQueryClient(zrConn)
	chainResp, err := validationClient.SolanaCounters(ctx, dcttypes.Asset_ASSET_ZENZEC.String())
	require.NoError(t, err)

	counters, ok := chainResp.Counters[dcttypes.Asset_ASSET_ZENZEC.String()]
	require.True(t, ok, "chain did not return counters for asset")

	chainMint := new(big.Int).SetUint64(counters.MintCounter)
	chainRedemption := new(big.Int).SetUint64(counters.RedemptionCounter)

	t.Logf("Solana program %s mint_counter(next)=%s, expected chain=%s, chain stored=%s (redemptions=%s)",
		programID, solMintCounter.String(), expectedChain.String(), chainMint.String(), chainRedemption.String())

	require.Zero(t, expectedChain.Cmp(chainMint), "chain mint counter mismatch")

	eventStoreProgram := sidecartypes.ZenZECEventStoreProgramID[cfg.Network]
	if eventStoreProgram != "" {
		maxEventID := fetchMaxEventStoreID(t, ctx, solClient, eventStoreProgram)
		if maxEventID != nil {
			t.Logf("EventStore %s highest wrap event id=%s", eventStoreProgram, maxEventID.String())
		}
	}
}

func fetchSolanaMintCounter(t *testing.T, ctx context.Context, client *solanarpc.Client, programIDStr string) *big.Int {
	programID := solana.MustPublicKeyFromBase58(programIDStr)
	globalConfigPDA, err := solzenbtc.GetGlobalConfigPDA(programID)
	require.NoError(t, err)

	resp, err := client.GetAccountInfo(ctx, globalConfigPDA)
	require.NoError(t, err)
	require.NotNil(t, resp.Value)

	data := resp.Value.Data.GetBinary()
	require.NotEmpty(t, data)

	decoder := bin.NewBorshDecoder(data)
	account := new(zenbtc_spl_token.GlobalConfig)
	require.NoError(t, account.UnmarshalWithDecoder(decoder))

	return account.MintCounter.BigInt()
}

func fetchMaxEventStoreID(t *testing.T, ctx context.Context, client *solanarpc.Client, programIDStr string) *big.Int {
	programID := solana.MustPublicKeyFromBase58(programIDStr)
	esClient := eventstore.NewClient(client, &programID)

	events, err := esClient.GetZenbtcWrapEvents(ctx)
	require.NoError(t, err)
	if len(events) == 0 {
		return nil
	}

	maxID := big.NewInt(0)
	for _, evt := range events {
		val := eventStoreIDToBigInt(evt.ID)
		if val.Cmp(maxID) > 0 {
			maxID = val
		}
	}
	return maxID
}
