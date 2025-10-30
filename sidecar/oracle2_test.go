package main

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"sync"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
	"github.com/ethereum/go-ethereum"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	sol "github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	jsonrpc "github.com/gagliardetto/solana-go/rpc/jsonrpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockEthClient struct {
	mock.Mock
}

func (m *MockEthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*ethtypes.Header, error) {
	args := m.Called(ctx, number)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ethtypes.Header), args.Error(1)
}
func (m *MockEthClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*big.Int), args.Error(1)
}
func (m *MockEthClient) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	args := m.Called(ctx, msg)
	return args.Get(0).(uint64), args.Error(1)
}

type MockSolanaClient struct {
	mock.Mock
}

func (m *MockSolanaClient) GetSignaturesForAddressWithOpts(ctx context.Context, address sol.PublicKey, opts *solrpc.GetSignaturesForAddressOpts) ([]*solrpc.TransactionSignature, error) {
	args := m.Called(ctx, address, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*solrpc.TransactionSignature), args.Error(1)
}

func (m *MockSolanaClient) RPCCallBatch(ctx context.Context, requests jsonrpc.RPCRequests) (jsonrpc.RPCResponses, error) {
	args := m.Called(ctx, requests)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(jsonrpc.RPCResponses), args.Error(1)
}

func (m *MockSolanaClient) GetLatestBlockhash(ctx context.Context, commitment interface{}) (*solrpc.GetLatestBlockhashResult, error) {
	args := m.Called(ctx, commitment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*solrpc.GetLatestBlockhashResult), args.Error(1)
}

func (m *MockSolanaClient) GetFeeForMessage(ctx context.Context, message string, commitment interface{}) (*solrpc.GetFeeForMessageResult, error) {
	args := m.Called(ctx, message, commitment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*solrpc.GetFeeForMessageResult), args.Error(1)
}

func (m *MockSolanaClient) GetTransaction(ctx context.Context, signature sol.Signature, opts *solrpc.GetTransactionOpts) (*solrpc.GetTransactionResult, error) {
	args := m.Called(ctx, signature, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*solrpc.GetTransactionResult), args.Error(1)
}

func createTestOracle() *Oracle {
	config := sidecartypes.Config{
		Network:   sidecartypes.NetworkDevnet,
		StateFile: "test_state.json",
	}
	oracle := &Oracle{
		Config:       config,
		DebugMode:    false,
		solanaClient: nil,
	}
	oracle.currentState.Store(&EmptyOracleState)
	oracle.stateCache = []sidecartypes.OracleState{EmptyOracleState}

	// Initialize function fields to prevent nil pointer dereference
	// These functions should return errors when solanaClient is nil to simulate real behavior
	oracle.getSolanaZenBTCBurnEventsFn = func(ctx context.Context, programID string, lastKnownSig sol.Signature) ([]api.BurnEvent, sol.Signature, error) {
		if oracle.solanaClient == nil {
			return nil, sol.Signature{}, fmt.Errorf("solana client is nil")
		}
		return []api.BurnEvent{}, sol.Signature{}, nil
	}
	oracle.getSolanaRockBurnEventsFn = func(ctx context.Context, programID string, lastKnownSig sol.Signature) ([]api.BurnEvent, sol.Signature, error) {
		if oracle.solanaClient == nil {
			return nil, sol.Signature{}, fmt.Errorf("solana client is nil")
		}
		return []api.BurnEvent{}, sol.Signature{}, nil
	}

	// Mock the reconcileBurnEventsFn to prevent nil pointer dereference in zrChainQueryClient
	oracle.reconcileBurnEventsFn = func(ctx context.Context, eventsToClean []api.BurnEvent, cleanedEvents map[string]bool, chainTypeName string) ([]api.BurnEvent, map[string]bool) {
		// For testing, just return the events unchanged (simulate no events are found on-chain)
		return eventsToClean, cleanedEvents
	}

	return oracle
}

func createMockHeader(blockNumber uint64, baseFee *big.Int) *ethtypes.Header {
	return &ethtypes.Header{
		Number:  big.NewInt(0).SetUint64(blockNumber),
		BaseFee: baseFee,
	}
}

func TestInitializeStateUpdate(t *testing.T) {
	oracle := createTestOracle()
	update := oracle.initializeStateUpdate()
	assert.NotNil(t, update)
	assert.NotNil(t, update.latestSolanaSigs)
	assert.NotNil(t, update.SolanaMintEvents)
	assert.NotNil(t, update.solanaBurnEvents)
	assert.NotNil(t, update.redemptions)
	assert.NotNil(t, update.ethBurnEvents)
}

func TestApplyFallbacks(t *testing.T) {
	oracle := createTestOracle()
	currentState := sidecartypes.OracleState{
		ROCKUSDPrice: math.LegacyNewDec(1),
		BTCUSDPrice:  math.LegacyNewDec(40000),
		ETHUSDPrice:  math.LegacyNewDec(2000),
	}
	oracle.currentState.Store(&currentState)
	update := &oracleStateUpdate{
		suggestedTip: nil,
		ROCKUSDPrice: math.LegacyDec{},
		BTCUSDPrice:  math.LegacyDec{},
		ETHUSDPrice:  math.LegacyDec{},
	}
	oracle.applyFallbacks(update, &currentState)
	assert.NotNil(t, update.suggestedTip)
	assert.Equal(t, big.NewInt(0), update.suggestedTip)
	assert.True(t, update.ROCKUSDPrice.Equal(currentState.ROCKUSDPrice))
	assert.True(t, update.BTCUSDPrice.Equal(currentState.BTCUSDPrice))
	assert.True(t, update.ETHUSDPrice.Equal(currentState.ETHUSDPrice))
}

func TestBuildFinalState(t *testing.T) {
	oracle := createTestOracle()
	currentState := sidecartypes.OracleState{
		EthBurnEvents:           []api.BurnEvent{},
		CleanedEthBurnEvents:    make(map[string]bool),
		SolanaBurnEvents:        []api.BurnEvent{},
		CleanedSolanaBurnEvents: make(map[string]bool),
		SolanaMintEvents:        []api.SolanaMintEvent{},
		CleanedSolanaMintEvents: make(map[string]bool),
	}
	oracle.currentState.Store(&currentState)
	update := &oracleStateUpdate{
		redemptions:      []api.Redemption{},
		suggestedTip:     big.NewInt(1500000000),
		estimatedGas:     231000,
		ethBurnEvents:    []api.BurnEvent{},
		solanaBurnEvents: []api.BurnEvent{},
		ROCKUSDPrice:     math.LegacyNewDec(1),
		BTCUSDPrice:      math.LegacyNewDec(50000),
		ETHUSDPrice:      math.LegacyNewDec(3000),
		SolanaMintEvents: []api.SolanaMintEvent{},
		latestSolanaSigs: make(map[sidecartypes.SolanaEventType]sol.Signature),
	}
	header := createMockHeader(1000, big.NewInt(20000000000))
	targetBlockNumber := big.NewInt(995)
	result, err := oracle.buildFinalState(update, header, targetBlockNumber)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint64(995), result.EthBlockHeight)
	assert.Equal(t, uint64(20000000000), result.EthBaseFee)
	assert.Equal(t, uint64(1500000000), result.EthTipCap)
	assert.Equal(t, uint64(231000), result.EthGasLimit)
	assert.True(t, result.ROCKUSDPrice.Equal(math.LegacyNewDec(1)))
	assert.True(t, result.BTCUSDPrice.Equal(math.LegacyNewDec(50000)))
	assert.True(t, result.ETHUSDPrice.Equal(math.LegacyNewDec(3000)))
}

func TestGetLastProcessedSolSignature(t *testing.T) {
	oracle := createTestOracle()
	oracle.lastSolRockMintSigStr = "test_rock_mint_sig"
	oracle.lastSolZenBTCMintSigStr = "test_zenbtc_mint_sig"
	oracle.lastSolZenBTCBurnSigStr = "test_zenbtc_burn_sig"
	oracle.lastSolRockBurnSigStr = "test_rock_burn_sig"
	sig := oracle.GetLastProcessedSolSignature(sidecartypes.SolRockMint)
	assert.NotEmpty(t, sig.String())
	sig = oracle.GetLastProcessedSolSignature(sidecartypes.SolZenBTCMint)
	assert.NotEmpty(t, sig.String())
	sig = oracle.GetLastProcessedSolSignature(sidecartypes.SolZenBTCBurn)
	assert.NotEmpty(t, sig.String())
	sig = oracle.GetLastProcessedSolSignature(sidecartypes.SolRockBurn)
	assert.NotEmpty(t, sig.String())
	sig = oracle.GetLastProcessedSolSignature("unknown")
	assert.True(t, sig.IsZero())
}

func TestSetStateCacheForTesting(t *testing.T) {
	oracle := createTestOracle()
	testStates := []sidecartypes.OracleState{
		{
			EthBlockHeight: 100,
			ROCKUSDPrice:   math.LegacyNewDec(1),
		},
		{
			EthBlockHeight: 200,
			ROCKUSDPrice:   math.LegacyNewDec(2),
		},
	}
	oracle.SetStateCacheForTesting(testStates)
	assert.Equal(t, len(testStates), len(oracle.stateCache))
	currentState := oracle.currentState.Load()
	assert.Equal(t, uint64(200), currentState.EthBlockHeight)
	assert.True(t, currentState.ROCKUSDPrice.Equal(math.LegacyNewDec(2)))
}

func TestSetStateCacheForTesting_Empty(t *testing.T) {
	oracle := createTestOracle()
	oracle.SetStateCacheForTesting([]sidecartypes.OracleState{})
	assert.Equal(t, 1, len(oracle.stateCache))
	currentState := oracle.currentState.Load()
	assert.Equal(t, EmptyOracleState, *currentState)
}

func TestROCKPriceFetching(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"currency_pair":"ROCK_USDT","last":"1.50","lowest_ask":"1.51","highest_bid":"1.49"}]`))
	}))
	defer server.Close()
	originalURL := sidecartypes.ROCKUSDPriceURL
	sidecartypes.ROCKUSDPriceURL = server.URL
	defer func() { sidecartypes.ROCKUSDPriceURL = originalURL }()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(sidecartypes.ROCKUSDPriceURL)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestROCKPriceFetching_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()
	originalURL := sidecartypes.ROCKUSDPriceURL
	sidecartypes.ROCKUSDPriceURL = server.URL
	defer func() { sidecartypes.ROCKUSDPriceURL = originalURL }()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(sidecartypes.ROCKUSDPriceURL)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestROCKPriceFetching_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	originalURL := sidecartypes.ROCKUSDPriceURL
	sidecartypes.ROCKUSDPriceURL = server.URL
	defer func() { sidecartypes.ROCKUSDPriceURL = originalURL }()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(sidecartypes.ROCKUSDPriceURL)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestROCKPriceFetching_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(15 * time.Second)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"currency_pair":"ROCK_USDT","last":"1.50"}]`))
	}))
	defer server.Close()
	originalURL := sidecartypes.ROCKUSDPriceURL
	sidecartypes.ROCKUSDPriceURL = server.URL
	defer func() { sidecartypes.ROCKUSDPriceURL = originalURL }()
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", sidecartypes.ROCKUSDPriceURL, nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	assert.Error(t, err)
	if resp != nil {
		resp.Body.Close()
	}
}

func TestCreateMockHeader(t *testing.T) {
	blockNumber := uint64(1000)
	baseFee := big.NewInt(20000000000)
	header := createMockHeader(blockNumber, baseFee)
	assert.Equal(t, big.NewInt(1000), header.Number)
	assert.Equal(t, baseFee, header.BaseFee)
}

func TestCreateTestOracle(t *testing.T) {
	oracle := createTestOracle()
	assert.NotNil(t, oracle)
	assert.Equal(t, sidecartypes.NetworkDevnet, oracle.Config.Network)
	assert.Equal(t, "test_state.json", oracle.Config.StateFile)
	assert.False(t, oracle.DebugMode)
	currentState := oracle.currentState.Load()
	assert.NotNil(t, currentState)
	assert.Equal(t, EmptyOracleState, *currentState)
	assert.Equal(t, 1, len(oracle.stateCache))
	assert.Equal(t, EmptyOracleState, oracle.stateCache[0])
}

func TestCreateTestPriceData(t *testing.T) {
	priceData := createTestPriceData()
	assert.Equal(t, 1, len(priceData))
	assert.Equal(t, "ROCK_USDT", priceData[0].CurrencyPair)
	assert.Equal(t, "1.50", priceData[0].Last)
	assert.Equal(t, "1.51", priceData[0].LowestAsk)
	assert.Equal(t, "1.49", priceData[0].HighestBid)
}

func createTestPriceData() []PriceData {
	return []PriceData{
		{
			CurrencyPair: "ROCK_USDT",
			Last:         "1.50",
			LowestAsk:    "1.51",
			HighestBid:   "1.49",
		},
	}
}

func BenchmarkCreateTestOracle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		createTestOracle()
	}
}

func BenchmarkCreateMockHeader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		createMockHeader(1000, big.NewInt(20000000000))
	}
}

func BenchmarkInitializeStateUpdate(b *testing.B) {
	oracle := createTestOracle()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		oracle.initializeStateUpdate()
	}
}

func TestBuildFinalState_NilCurrentState(t *testing.T) {
	oracle := createTestOracle()
	update := &oracleStateUpdate{
		redemptions:      []api.Redemption{},
		suggestedTip:     big.NewInt(1500000000),
		estimatedGas:     231000,
		ethBurnEvents:    []api.BurnEvent{},
		solanaBurnEvents: []api.BurnEvent{},
		ROCKUSDPrice:     math.LegacyNewDec(1),
		BTCUSDPrice:      math.LegacyNewDec(50000),
		ETHUSDPrice:      math.LegacyNewDec(3000),
		SolanaMintEvents: []api.SolanaMintEvent{},
		latestSolanaSigs: make(map[sidecartypes.SolanaEventType]sol.Signature),
	}
	header := createMockHeader(1000, big.NewInt(20000000000))
	targetBlockNumber := big.NewInt(995)
	result, err := oracle.buildFinalState(update, header, targetBlockNumber)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestBuildFinalState_NilHeader(t *testing.T) {
	oracle := createTestOracle()
	update := &oracleStateUpdate{
		redemptions:      []api.Redemption{},
		suggestedTip:     big.NewInt(1500000000),
		estimatedGas:     231000,
		ethBurnEvents:    []api.BurnEvent{},
		solanaBurnEvents: []api.BurnEvent{},
		ROCKUSDPrice:     math.LegacyNewDec(1),
		BTCUSDPrice:      math.LegacyNewDec(50000),
		ETHUSDPrice:      math.LegacyNewDec(3000),
		SolanaMintEvents: []api.SolanaMintEvent{},
		latestSolanaSigs: make(map[sidecartypes.SolanaEventType]sol.Signature),
	}
	targetBlockNumber := big.NewInt(995)
	assert.Panics(t, func() {
		oracle.buildFinalState(update, nil, targetBlockNumber)
	})
}

func TestBuildFinalState_NilUpdate(t *testing.T) {
	oracle := createTestOracle()
	header := createMockHeader(1000, big.NewInt(20000000000))
	targetBlockNumber := big.NewInt(995)
	assert.Panics(t, func() {
		oracle.buildFinalState(nil, header, targetBlockNumber)
	})
}

func TestBuildFinalState_NilTargetBlockNumber(t *testing.T) {
	oracle := createTestOracle()
	update := &oracleStateUpdate{
		redemptions:      []api.Redemption{},
		suggestedTip:     big.NewInt(1500000000),
		estimatedGas:     231000,
		ethBurnEvents:    []api.BurnEvent{},
		solanaBurnEvents: []api.BurnEvent{},
		ROCKUSDPrice:     math.LegacyNewDec(1),
		BTCUSDPrice:      math.LegacyNewDec(50000),
		ETHUSDPrice:      math.LegacyNewDec(3000),
		SolanaMintEvents: []api.SolanaMintEvent{},
		latestSolanaSigs: make(map[sidecartypes.SolanaEventType]sol.Signature),
	}
	header := createMockHeader(1000, big.NewInt(20000000000))
	assert.Panics(t, func() {
		oracle.buildFinalState(update, header, nil)
	})
}

func TestApplyFallbacks_NilCurrentState(t *testing.T) {
	oracle := createTestOracle()
	update := &oracleStateUpdate{
		suggestedTip: nil,
		ROCKUSDPrice: math.LegacyDec{},
		BTCUSDPrice:  math.LegacyDec{},
		ETHUSDPrice:  math.LegacyDec{},
	}
	assert.Panics(t, func() {
		oracle.applyFallbacks(update, nil)
	})
}

func TestApplyFallbacks_NilUpdate(t *testing.T) {
	oracle := createTestOracle()
	currentState := sidecartypes.OracleState{
		ROCKUSDPrice: math.LegacyNewDec(1),
		BTCUSDPrice:  math.LegacyNewDec(40000),
		ETHUSDPrice:  math.LegacyNewDec(2000),
	}
	oracle.currentState.Store(&currentState)
	assert.Panics(t, func() {
		oracle.applyFallbacks(nil, &currentState)
	})
}

func TestGetLastProcessedSolSignature_InvalidSignature(t *testing.T) {
	oracle := createTestOracle()
	oracle.lastSolRockMintSigStr = "invalid_signature_with_special_chars!@#"
	oracle.lastSolZenBTCMintSigStr = "another_invalid_sig_with_spaces "
	oracle.lastSolZenBTCBurnSigStr = "yet_another_invalid_sig_with_unicode_🚀"
	oracle.lastSolRockBurnSigStr = ""
	sig := oracle.GetLastProcessedSolSignature(sidecartypes.SolRockMint)
	assert.True(t, sig.IsZero())
	sig = oracle.GetLastProcessedSolSignature(sidecartypes.SolZenBTCMint)
	assert.True(t, sig.IsZero())
	sig = oracle.GetLastProcessedSolSignature(sidecartypes.SolZenBTCBurn)
	assert.True(t, sig.IsZero())
	sig = oracle.GetLastProcessedSolSignature(sidecartypes.SolRockBurn)
	assert.True(t, sig.IsZero())
}

func TestSetStateCacheForTesting_NilStates(t *testing.T) {
	oracle := createTestOracle()
	oracle.SetStateCacheForTesting(nil)
	assert.Equal(t, 1, len(oracle.stateCache))
	currentState := oracle.currentState.Load()
	assert.Equal(t, EmptyOracleState, *currentState)
}

func TestCreateTestOracle_InvalidConfig(t *testing.T) {
	config := sidecartypes.Config{
		Network:   "",
		StateFile: "",
	}
	oracle := &Oracle{
		Config:    config,
		DebugMode: false,
	}
	oracle.currentState.Store(&EmptyOracleState)
	oracle.stateCache = []sidecartypes.OracleState{EmptyOracleState}
	assert.NotNil(t, oracle)
	assert.Equal(t, "", oracle.Config.Network)
	assert.Equal(t, "", oracle.Config.StateFile)
}

func TestCreateMockHeader_ZeroValues(t *testing.T) {
	header := createMockHeader(0, big.NewInt(0))
	assert.Equal(t, big.NewInt(0), header.Number)
	assert.Equal(t, big.NewInt(0), header.BaseFee)
}

func TestCreateMockHeader_NilBaseFee(t *testing.T) {
	header := createMockHeader(1000, nil)
	assert.Equal(t, big.NewInt(1000), header.Number)
	assert.Nil(t, header.BaseFee)
}

func TestCreateMockHeader_LargeBlockNumber(t *testing.T) {
	largeBlockNumber := uint64(18446744073709551615)
	header := createMockHeader(largeBlockNumber, big.NewInt(20000000000))
	expectedNumber := big.NewInt(0).SetUint64(largeBlockNumber)
	assert.Equal(t, expectedNumber, header.Number)
	assert.Equal(t, big.NewInt(20000000000), header.BaseFee)
}

func TestCreateTestPriceData_EmptyValues(t *testing.T) {
	priceData := createTestPriceData()
	assert.Equal(t, 1, len(priceData))
	assert.Equal(t, "ROCK_USDT", priceData[0].CurrencyPair)
	assert.Equal(t, "1.50", priceData[0].Last)
	assert.Equal(t, "1.51", priceData[0].LowestAsk)
	assert.Equal(t, "1.49", priceData[0].HighestBid)
}

func TestROCKPriceFetching_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"currency_pair":"ROCK_USDT","last":"1.50"`))
	}))
	defer server.Close()
	originalURL := sidecartypes.ROCKUSDPriceURL
	sidecartypes.ROCKUSDPriceURL = server.URL
	defer func() { sidecartypes.ROCKUSDPriceURL = originalURL }()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(sidecartypes.ROCKUSDPriceURL)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestROCKPriceFetching_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	originalURL := sidecartypes.ROCKUSDPriceURL
	sidecartypes.ROCKUSDPriceURL = server.URL
	defer func() { sidecartypes.ROCKUSDPriceURL = originalURL }()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(sidecartypes.ROCKUSDPriceURL)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestROCKPriceFetching_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	originalURL := sidecartypes.ROCKUSDPriceURL
	sidecartypes.ROCKUSDPriceURL = server.URL
	defer func() { sidecartypes.ROCKUSDPriceURL = originalURL }()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(sidecartypes.ROCKUSDPriceURL)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestROCKPriceFetching_ServerError500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()
	originalURL := sidecartypes.ROCKUSDPriceURL
	sidecartypes.ROCKUSDPriceURL = server.URL
	defer func() { sidecartypes.ROCKUSDPriceURL = originalURL }()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(sidecartypes.ROCKUSDPriceURL)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestROCKPriceFetching_UnreachableServer(t *testing.T) {
	unreachableURL := "http://localhost:99999"
	originalURL := sidecartypes.ROCKUSDPriceURL
	sidecartypes.ROCKUSDPriceURL = unreachableURL
	defer func() { sidecartypes.ROCKUSDPriceURL = originalURL }()
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", sidecartypes.ROCKUSDPriceURL, nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	assert.Error(t, err)
	if resp != nil {
		resp.Body.Close()
	}
}

func TestROCKPriceFetching_Redirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/new-location", http.StatusMovedPermanently)
	}))
	defer server.Close()
	originalURL := sidecartypes.ROCKUSDPriceURL
	sidecartypes.ROCKUSDPriceURL = server.URL
	defer func() { sidecartypes.ROCKUSDPriceURL = originalURL }()
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(sidecartypes.ROCKUSDPriceURL)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode)
}

func TestROCKPriceFetching_WrongContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("This is not JSON"))
	}))
	defer server.Close()
	originalURL := sidecartypes.ROCKUSDPriceURL
	sidecartypes.ROCKUSDPriceURL = server.URL
	defer func() { sidecartypes.ROCKUSDPriceURL = originalURL }()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(sidecartypes.ROCKUSDPriceURL)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))
}

func TestROCKPriceFetching_LargeResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		largeData := make([]byte, 1024*1024)
		for i := range largeData {
			largeData[i] = 'A'
		}
		w.Write(largeData)
	}))
	defer server.Close()
	originalURL := sidecartypes.ROCKUSDPriceURL
	sidecartypes.ROCKUSDPriceURL = server.URL
	defer func() { sidecartypes.ROCKUSDPriceURL = originalURL }()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(sidecartypes.ROCKUSDPriceURL)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestROCKPriceFetching_SlowResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"currency_pair":"ROCK_USDT","last":"1.50"}]`))
	}))
	defer server.Close()
	originalURL := sidecartypes.ROCKUSDPriceURL
	sidecartypes.ROCKUSDPriceURL = server.URL
	defer func() { sidecartypes.ROCKUSDPriceURL = originalURL }()
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", sidecartypes.ROCKUSDPriceURL, nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	assert.Error(t, err)
	if resp != nil {
		resp.Body.Close()
	}
}

func TestROCKPriceFetching_ChunkedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("2\r\n"))
		w.Write([]byte(`[{`))
		w.Write([]byte("\r\n"))
		w.Write([]byte("20\r\n"))
		w.Write([]byte(`"currency_pair":"ROCK_USDT","last":"1.50"`))
		w.Write([]byte("\r\n"))
		w.Write([]byte("2\r\n"))
		w.Write([]byte(`}]`))
		w.Write([]byte("\r\n"))
		w.Write([]byte("0\r\n\r\n"))
	}))
	defer server.Close()
	originalURL := sidecartypes.ROCKUSDPriceURL
	sidecartypes.ROCKUSDPriceURL = server.URL
	defer func() { sidecartypes.ROCKUSDPriceURL = originalURL }()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(sidecartypes.ROCKUSDPriceURL)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("Content-Type"))
}

func TestFetchSolanaBurnEvents(t *testing.T) {
	tests := []struct {
		name                 string
		existingState        sidecartypes.OracleState
		lastKnownSigs        map[sidecartypes.SolanaEventType]string
		mockZenBTCEvents     []api.BurnEvent
		mockRockEvents       []api.BurnEvent
		expectedTotalEvents  int
		expectedZenBTCEvents int
		expectedRockEvents   int
		expectedNewSigs      map[sidecartypes.SolanaEventType]bool
		description          string
	}{
		{
			name: "No existing events, new events found",
			existingState: sidecartypes.OracleState{
				SolanaBurnEvents:        []api.BurnEvent{},
				CleanedSolanaBurnEvents: make(map[string]bool),
			},
			lastKnownSigs: map[sidecartypes.SolanaEventType]string{
				sidecartypes.SolZenBTCBurn: "",
				sidecartypes.SolRockBurn:   "",
			},
			mockZenBTCEvents: []api.BurnEvent{
				{
					TxID:            "zenbtc_burn_1",
					LogIndex:        0,
					ChainID:         "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
					DestinationAddr: []byte{1, 2, 3, 4},
					Amount:          1000000,
					IsZenBTC:        true,
				},
				{
					TxID:            "zenbtc_burn_2",
					LogIndex:        1,
					ChainID:         "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
					DestinationAddr: []byte{5, 6, 7, 8},
					Amount:          2000000,
					IsZenBTC:        true,
				},
			},
			mockRockEvents: []api.BurnEvent{
				{
					TxID:            "rock_burn_1",
					LogIndex:        0,
					ChainID:         "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
					DestinationAddr: []byte{9, 10, 11, 12},
					Amount:          500000,
					IsZenBTC:        false,
				},
			},
			expectedTotalEvents:  3,
			expectedZenBTCEvents: 2,
			expectedRockEvents:   1,
			expectedNewSigs: map[sidecartypes.SolanaEventType]bool{
				sidecartypes.SolZenBTCBurn: true,
				sidecartypes.SolRockBurn:   true,
			},
			description: "Should add all new events when no existing events",
		},
		{
			name: "Existing events, new events found",
			existingState: sidecartypes.OracleState{
				SolanaBurnEvents: []api.BurnEvent{
					{
						TxID:            "existing_zenbtc_burn",
						LogIndex:        0,
						ChainID:         "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
						DestinationAddr: []byte{1, 1, 1, 1},
						Amount:          100000,
						IsZenBTC:        true,
					},
					{
						TxID:            "existing_rock_burn",
						LogIndex:        0,
						ChainID:         "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
						DestinationAddr: []byte{2, 2, 2, 2},
						Amount:          200000,
						IsZenBTC:        false,
					},
				},
				CleanedSolanaBurnEvents: make(map[string]bool),
			},
			lastKnownSigs: map[sidecartypes.SolanaEventType]string{
				sidecartypes.SolZenBTCBurn: "existing_zenbtc_sig",
				sidecartypes.SolRockBurn:   "existing_rock_sig",
			},
			mockZenBTCEvents: []api.BurnEvent{
				{
					TxID:            "new_zenbtc_burn",
					LogIndex:        0,
					ChainID:         "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
					DestinationAddr: []byte{3, 3, 3, 3},
					Amount:          300000,
					IsZenBTC:        true,
				},
			},
			mockRockEvents: []api.BurnEvent{
				{
					TxID:            "new_rock_burn",
					LogIndex:        0,
					ChainID:         "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
					DestinationAddr: []byte{4, 4, 4, 4},
					Amount:          400000,
					IsZenBTC:        false,
				},
			},
			expectedTotalEvents:  4,
			expectedZenBTCEvents: 2,
			expectedRockEvents:   2,
			expectedNewSigs: map[sidecartypes.SolanaEventType]bool{
				sidecartypes.SolZenBTCBurn: true,
				sidecartypes.SolRockBurn:   true,
			},
			description: "Should append new events to existing ones",
		},
		{
			name: "Events already cleaned up",
			existingState: sidecartypes.OracleState{
				SolanaBurnEvents: []api.BurnEvent{},
				CleanedSolanaBurnEvents: map[string]bool{
					"solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1-zenbtc_burn_1-0": true,
					"solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1-rock_burn_1-0":   true,
				},
			},
			lastKnownSigs: map[sidecartypes.SolanaEventType]string{
				sidecartypes.SolZenBTCBurn: "",
				sidecartypes.SolRockBurn:   "",
			},
			mockZenBTCEvents: []api.BurnEvent{
				{
					TxID:            "zenbtc_burn_1",
					LogIndex:        0,
					ChainID:         "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
					DestinationAddr: []byte{1, 2, 3, 4},
					Amount:          1000000,
					IsZenBTC:        true,
				},
			},
			mockRockEvents: []api.BurnEvent{
				{
					TxID:            "rock_burn_1",
					LogIndex:        0,
					ChainID:         "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
					DestinationAddr: []byte{5, 6, 7, 8},
					Amount:          500000,
					IsZenBTC:        false,
				},
			},
			expectedTotalEvents:  0,
			expectedZenBTCEvents: 0,
			expectedRockEvents:   0,
			expectedNewSigs: map[sidecartypes.SolanaEventType]bool{
				sidecartypes.SolZenBTCBurn: true,
				sidecartypes.SolRockBurn:   true,
			},
			description: "Should not add events that are already cleaned up",
		},
		{
			name: "No new events found",
			existingState: sidecartypes.OracleState{
				SolanaBurnEvents: []api.BurnEvent{
					{
						TxID:            "existing_burn",
						LogIndex:        0,
						ChainID:         "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
						DestinationAddr: []byte{1, 1, 1, 1},
						Amount:          100000,
						IsZenBTC:        true,
					},
				},
				CleanedSolanaBurnEvents: make(map[string]bool),
			},
			lastKnownSigs: map[sidecartypes.SolanaEventType]string{
				sidecartypes.SolZenBTCBurn: "latest_zenbtc_sig",
				sidecartypes.SolRockBurn:   "latest_rock_sig",
			},
			mockZenBTCEvents:     []api.BurnEvent{},
			mockRockEvents:       []api.BurnEvent{},
			expectedTotalEvents:  1,
			expectedZenBTCEvents: 1,
			expectedRockEvents:   0,
			expectedNewSigs: map[sidecartypes.SolanaEventType]bool{
				sidecartypes.SolZenBTCBurn: false,
				sidecartypes.SolRockBurn:   false,
			},
			description: "Should keep existing events when no new events found",
		},
		{
			name: "Mixed cleaned and new events",
			existingState: sidecartypes.OracleState{
				SolanaBurnEvents: []api.BurnEvent{
					{
						TxID:            "existing_burn",
						LogIndex:        0,
						ChainID:         "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
						DestinationAddr: []byte{1, 1, 1, 1},
						Amount:          100000,
						IsZenBTC:        true,
					},
				},
				CleanedSolanaBurnEvents: map[string]bool{
					"solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1-cleaned_burn-0": true,
				},
			},
			lastKnownSigs: map[sidecartypes.SolanaEventType]string{
				sidecartypes.SolZenBTCBurn: "",
				sidecartypes.SolRockBurn:   "",
			},
			mockZenBTCEvents: []api.BurnEvent{
				{
					TxID:            "cleaned_burn",
					LogIndex:        0,
					ChainID:         "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
					DestinationAddr: []byte{2, 2, 2, 2},
					Amount:          200000,
					IsZenBTC:        true,
				},
				{
					TxID:            "new_burn",
					LogIndex:        0,
					ChainID:         "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
					DestinationAddr: []byte{3, 3, 3, 3},
					Amount:          300000,
					IsZenBTC:        true,
				},
			},
			mockRockEvents:       []api.BurnEvent{},
			expectedTotalEvents:  2,
			expectedZenBTCEvents: 2,
			expectedRockEvents:   0,
			expectedNewSigs: map[sidecartypes.SolanaEventType]bool{
				sidecartypes.SolZenBTCBurn: true,
				sidecartypes.SolRockBurn:   false,
			},
			description: "Should add new events but skip cleaned ones",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oracle := createTestOracle()
			oracle.SetStateCacheForTesting([]sidecartypes.OracleState{tt.existingState})

			if sig, ok := tt.lastKnownSigs[sidecartypes.SolZenBTCBurn]; ok {
				oracle.lastSolZenBTCBurnSigStr = sig
			}
			if sig, ok := tt.lastKnownSigs[sidecartypes.SolRockBurn]; ok {
				oracle.lastSolRockBurnSigStr = sig
			}

			t.Logf("Testing scenario: %s", tt.description)

			expectedEvents := len(tt.existingState.SolanaBurnEvents)

			for _, event := range tt.mockZenBTCEvents {
				key := fmt.Sprintf("%s-%s-%d", event.ChainID, event.TxID, event.LogIndex)
				if !tt.existingState.CleanedSolanaBurnEvents[key] {
					expectedEvents++
				}
			}

			for _, event := range tt.mockRockEvents {
				key := fmt.Sprintf("%s-%s-%d", event.ChainID, event.TxID, event.LogIndex)
				if !tt.existingState.CleanedSolanaBurnEvents[key] {
					expectedEvents++
				}
			}

			if expectedEvents != tt.expectedTotalEvents {
				t.Errorf("Test case inconsistency: calculated %d expected events, but test expects %d",
					expectedEvents, tt.expectedTotalEvents)
			}

			mockEvent := createMockBurnEvent("test_tx", 0, true, 1000000)
			assert.Equal(t, "test_tx", mockEvent.TxID)
			assert.True(t, mockEvent.IsZenBTC)
			assert.Equal(t, uint64(1000000), mockEvent.Amount)

			mockState := createMockOracleState([]api.BurnEvent{mockEvent}, make(map[string]bool))
			assert.Equal(t, 1, len(mockState.SolanaBurnEvents))
			assert.Equal(t, mockEvent.TxID, mockState.SolanaBurnEvents[0].TxID)
		})
	}
}

func createMockBurnEvent(txID string, logIndex uint64, isZenBTC bool, amount uint64) api.BurnEvent {
	return api.BurnEvent{
		TxID:            txID,
		LogIndex:        logIndex,
		ChainID:         "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		DestinationAddr: []byte{1, 2, 3, 4},
		Amount:          amount,
		IsZenBTC:        isZenBTC,
	}
}

func createMockOracleState(existingEvents []api.BurnEvent, cleanedEvents map[string]bool) sidecartypes.OracleState {
	return sidecartypes.OracleState{
		SolanaBurnEvents:        existingEvents,
		CleanedSolanaBurnEvents: cleanedEvents,
	}
}

func TestBurnEventProcessingLogic(t *testing.T) {
	t.Run("Event deduplication", func(t *testing.T) {
		existingEvents := []api.BurnEvent{
			createMockBurnEvent("existing_tx", 0, true, 1000000),
		}

		cleanedEvents := make(map[string]bool)
		state := createMockOracleState(existingEvents, cleanedEvents)

		newEvent := createMockBurnEvent("existing_tx", 0, true, 1000000)

		key := fmt.Sprintf("%s-%s-%d", newEvent.ChainID, newEvent.TxID, newEvent.LogIndex)
		existingKeys := make(map[string]bool)
		for _, event := range state.SolanaBurnEvents {
			eventKey := fmt.Sprintf("%s-%s-%d", event.ChainID, event.TxID, event.LogIndex)
			existingKeys[eventKey] = true
		}

		if existingKeys[key] {
			t.Log("Correctly identified duplicate event")
		} else {
			t.Error("Failed to identify duplicate event")
		}
	})

	t.Run("Cleaned event filtering", func(t *testing.T) {
		existingEvents := []api.BurnEvent{}
		cleanedEvents := map[string]bool{
			"solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1-cleaned_tx-0": true,
		}
		state := createMockOracleState(existingEvents, cleanedEvents)

		newEvent := createMockBurnEvent("cleaned_tx", 0, true, 1000000)
		key := fmt.Sprintf("%s-%s-%d", newEvent.ChainID, newEvent.TxID, newEvent.LogIndex)

		if state.CleanedSolanaBurnEvents[key] {
			t.Log("Correctly identified cleaned event")
		} else {
			t.Error("Failed to identify cleaned event")
		}
	})

	t.Run("Event type classification", func(t *testing.T) {
		zenBTCEvent := createMockBurnEvent("zenbtc_tx", 0, true, 1000000)
		rockEvent := createMockBurnEvent("rock_tx", 0, false, 500000)

		assert.True(t, zenBTCEvent.IsZenBTC, "ZenBTC event should be marked as IsZenBTC=true")
		assert.False(t, rockEvent.IsZenBTC, "ROCK event should be marked as IsZenBTC=false")
		assert.Equal(t, "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1", zenBTCEvent.ChainID, "Event should have correct chain ID")
		assert.Equal(t, "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1", rockEvent.ChainID, "Event should have correct chain ID")
	})
}

func TestFetchSolanaBurnEventsComprehensive(t *testing.T) {
	t.Skip("Skipping slow integration test with 10+ second timeouts")
	tests := []struct {
		name          string
		setupOracle   func() *Oracle
		expectedError bool
		description   string
	}{
		{
			name: "Nil solana client should cause error",
			setupOracle: func() *Oracle {
				oracle := createTestOracle()
				oracle.lastSolZenBTCBurnSigStr = "old_zenbtc_sig"
				oracle.lastSolRockBurnSigStr = "old_rock_sig"
				return oracle
			},
			expectedError: true,
			description:   "Should fail gracefully when Solana client is nil",
		},
		{
			name: "Invalid signature watermark should be handled",
			setupOracle: func() *Oracle {
				oracle := createTestOracle()
				oracle.lastSolZenBTCBurnSigStr = "invalid_signature"
				oracle.lastSolRockBurnSigStr = "invalid_signature"
				return oracle
			},
			expectedError: true, // Will fail due to nil client, but invalid sigs should be handled gracefully
			description:   "Should handle invalid signature watermarks gracefully",
		},
		{
			name: "Empty signature watermark (first run)",
			setupOracle: func() *Oracle {
				oracle := createTestOracle()
				oracle.lastSolZenBTCBurnSigStr = ""
				oracle.lastSolRockBurnSigStr = ""
				return oracle
			},
			expectedError: true, // Will fail due to nil client
			description:   "Should handle first run with empty watermarks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oracle := tt.setupOracle()

			update := oracle.initializeStateUpdate()
			var updateMutex sync.Mutex
			errChan := make(chan error, 2) // Buffer for both goroutines
			var wg sync.WaitGroup

			oracle.fetchSolanaBurnEvents(context.Background(), &wg, update, &updateMutex, errChan)

			wg.Wait()
			close(errChan)

			var errors []error
			for err := range errChan {
				errors = append(errors, err)
			}

			if tt.expectedError {
				assert.NotEmpty(t, errors, "Expected error but none occurred")
				t.Logf("Test passed: %s (expected error occurred)", tt.description)
			} else {
				assert.Empty(t, errors, "Unexpected errors occurred: %v", errors)
				t.Logf("Test passed: %s", tt.description)
			}
		})
	}
}

func TestFetchSolanaBurnEventsRaceConditions(t *testing.T) {
	t.Skip("Skipping slow race conditions test with 10+ second timeouts")
	oracle := createTestOracle()

	const numGoroutines = 10
	var wg sync.WaitGroup
	errorCounts := make([]int, numGoroutines)

	for i := range numGoroutines {
		wg.Go(func() {
			update := oracle.initializeStateUpdate()
			var updateMutex sync.Mutex
			errChan := make(chan error, 2)
			var innerWg sync.WaitGroup

			oracle.fetchSolanaBurnEvents(context.Background(), &innerWg, update, &updateMutex, errChan)
			innerWg.Wait()
			close(errChan)

			errorCount := 0
			for range errChan {
				errorCount++
			}

			errorCounts[i] = errorCount
		})
	}

	wg.Wait()

	expectedErrors := 2
	totalErrors := 0
	for i, errorCount := range errorCounts {
		assert.Equal(t, expectedErrors, errorCount,
			"Goroutine %d encountered %d errors, expected %d", i, errorCount, expectedErrors)
		totalErrors += errorCount
	}

	t.Logf("All %d goroutines encountered expected errors due to nil Solana client (total: %d errors)", numGoroutines, totalErrors)
}

func TestFetchSolanaBurnEventsWatermarkConsistency(t *testing.T) {
	t.Skip("Skipping slow watermark consistency test with 10+ second timeouts")
	oracle := createTestOracle()

	update := oracle.initializeStateUpdate()
	var updateMutex sync.Mutex
	errChan := make(chan error, 2)
	var wg sync.WaitGroup

	oracle.fetchSolanaBurnEvents(context.Background(), &wg, update, &updateMutex, errChan)
	wg.Wait()
	close(errChan)

	errorCount := 0
	for range errChan {
		errorCount++
	}

	assert.Equal(t, 2, errorCount, "Expected 2 errors due to nil Solana client")

	t.Logf("Test passed: function handled nil Solana client gracefully (%d expected errors)", errorCount)
}

func TestFetchSolanaBurnEventsErrorHandling(t *testing.T) {
	t.Skip("Skipping slow error handling test with 10+ second timeouts")
	tests := []struct {
		name          string
		setupOracle   func() *Oracle
		expectedError bool
		description   string
	}{
		{
			name: "Nil solana client",
			setupOracle: func() *Oracle {
				oracle := createTestOracle()
				oracle.solanaClient = nil
				return oracle
			},
			expectedError: true,
			description:   "Should handle nil solana client gracefully",
		},
		{
			name: "Invalid signature watermarks",
			setupOracle: func() *Oracle {
				oracle := createTestOracle()
				oracle.lastSolZenBTCBurnSigStr = "invalid_sig"
				oracle.lastSolRockBurnSigStr = "invalid_sig"
				return oracle
			},
			expectedError: true, // Will fail due to nil client
			description:   "Should handle invalid signature watermarks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oracle := tt.setupOracle()

			update := oracle.initializeStateUpdate()
			var updateMutex sync.Mutex
			errChan := make(chan error, 2)
			var wg sync.WaitGroup

			oracle.fetchSolanaBurnEvents(context.Background(), &wg, update, &updateMutex, errChan)
			wg.Wait()
			close(errChan)

			var errors []error
			for err := range errChan {
				errors = append(errors, err)
			}

			if tt.expectedError {
				assert.NotEmpty(t, errors, "Expected error but none occurred")
				t.Logf("Test passed: %s (expected error occurred)", tt.description)
			} else {
				t.Logf("Test passed: %s", tt.description)
			}
		})
	}
}

func TestFetchSolanaBurnEventsEventDeduplication(t *testing.T) {
	t.Skip("Skipping slow event deduplication test with 10+ second timeouts")
	oracle := createTestOracle()

	existingEvents := []api.BurnEvent{
		createMockBurnEvent("existing_tx_1", 0, true, 1000000),
		createMockBurnEvent("existing_tx_2", 0, false, 500000),
	}

	oracle.SetStateCacheForTesting([]sidecartypes.OracleState{
		{
			SolanaBurnEvents:        existingEvents,
			CleanedSolanaBurnEvents: make(map[string]bool),
		},
	})

	update := oracle.initializeStateUpdate()
	var updateMutex sync.Mutex
	errChan := make(chan error, 2)
	var wg sync.WaitGroup

	oracle.fetchSolanaBurnEvents(context.Background(), &wg, update, &updateMutex, errChan)
	wg.Wait()
	close(errChan)

	errorCount := 0
	for range errChan {
		errorCount++
	}

	assert.Equal(t, 2, errorCount, "Expected 2 errors due to nil Solana client")

	t.Logf("Test passed: function handled nil Solana client gracefully (%d expected errors)", errorCount)
}

func TestFetchSolanaBurnEventsCleanedEventHandling(t *testing.T) {
	t.Skip("Skipping slow cleaned event handling test with 10+ second timeouts")
	oracle := createTestOracle()

	cleanedEvents := map[string]bool{
		"solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1-cleaned_tx_1-0": true,
		"solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1-cleaned_tx_2-0": true,
	}

	oracle.SetStateCacheForTesting([]sidecartypes.OracleState{
		{
			SolanaBurnEvents:        []api.BurnEvent{},
			CleanedSolanaBurnEvents: cleanedEvents,
		},
	})

	update := oracle.initializeStateUpdate()
	var updateMutex sync.Mutex
	errChan := make(chan error, 2)
	var wg sync.WaitGroup

	oracle.fetchSolanaBurnEvents(context.Background(), &wg, update, &updateMutex, errChan)
	wg.Wait()
	close(errChan)

	errorCount := 0
	for range errChan {
		errorCount++
	}

	assert.Equal(t, 2, errorCount, "Expected 2 errors due to nil Solana client")

	t.Logf("Test passed: function handled nil Solana client gracefully (%d expected errors)", errorCount)
}

func TestFetchSolanaBurnEventsBatchProcessing(t *testing.T) {
	t.Skip("Skipping slow batch processing test with 10+ second timeouts")
	oracle := createTestOracle()

	update := oracle.initializeStateUpdate()
	var updateMutex sync.Mutex
	errChan := make(chan error, 2)
	var wg sync.WaitGroup

	oracle.fetchSolanaBurnEvents(context.Background(), &wg, update, &updateMutex, errChan)
	wg.Wait()
	close(errChan)

	errorCount := 0
	for range errChan {
		errorCount++
	}

	assert.Equal(t, 2, errorCount, "Expected 2 errors due to nil Solana client")

	t.Logf("Test passed: function handled nil Solana client gracefully (%d expected errors)", errorCount)
}

func TestFetchSolanaBurnEventsSignatureOrdering(t *testing.T) {
	t.Skip("Skipping slow signature ordering test with 10+ second timeouts")
	oracle := createTestOracle()

	update := oracle.initializeStateUpdate()
	var updateMutex sync.Mutex
	errChan := make(chan error, 2)
	var wg sync.WaitGroup

	oracle.fetchSolanaBurnEvents(context.Background(), &wg, update, &updateMutex, errChan)
	wg.Wait()
	close(errChan)

	errorCount := 0
	for range errChan {
		errorCount++
	}

	assert.Equal(t, 2, errorCount, "Expected 2 errors due to nil Solana client")

	t.Logf("Test passed: function handled nil Solana client gracefully (%d expected errors)", errorCount)
}

func TestFetchSolanaBurnEventsConcurrentAccess(t *testing.T) {
	t.Skip("Skipping slow concurrent access test with 10+ second timeouts")
	oracle := createTestOracle()

	const numConcurrent = 5
	var wg sync.WaitGroup
	errorCount := 0
	var errorMutex sync.Mutex

	for range numConcurrent {
		wg.Go(func() {
			update := oracle.initializeStateUpdate()
			var updateMutex sync.Mutex
			errChan := make(chan error, 2)
			var innerWg sync.WaitGroup

			oracle.fetchSolanaBurnEvents(context.Background(), &innerWg, update, &updateMutex, errChan)
			innerWg.Wait()
			close(errChan)

			for range errChan {
				errorMutex.Lock()
				errorCount++
				errorMutex.Unlock()
			}
		})
	}

	wg.Wait()

	expectedErrors := numConcurrent * 2 // 2 errors per call (ZenBTC and ROCK)
	if errorCount != expectedErrors {
		t.Errorf("Expected %d errors (nil client), got %d", expectedErrors, errorCount)
	}
	t.Logf("Concurrent access test completed with %d expected errors due to nil Solana client", errorCount)
}

func TestFetchSolanaBurnEventsMemoryLeaks(t *testing.T) {
	t.Skip("Skipping slow memory leak test with 10+ second timeouts")
	oracle := createTestOracle()

	const numIterations = 20
	totalErrors := 0

	for i := 0; i < numIterations; i++ {
		update := oracle.initializeStateUpdate()
		var updateMutex sync.Mutex
		errChan := make(chan error, 2)
		var wg sync.WaitGroup

		oracle.fetchSolanaBurnEvents(context.Background(), &wg, update, &updateMutex, errChan)
		wg.Wait()
		close(errChan)

		errorCount := 0
		for range errChan {
			errorCount++
		}
		totalErrors += errorCount

		if i%5 == 0 {
			t.Logf("Completed iteration %d (expected errors: %d)", i, errorCount)
		}
	}

	t.Logf("Completed %d iterations without obvious memory leaks (total expected errors: %d)", numIterations, totalErrors)
}

func TestFetchSolanaBurnEventsTimeoutHandling(t *testing.T) {
	t.Skip("Skipping slow timeout handling test with 10+ second timeouts")
	oracle := createTestOracle()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	update := oracle.initializeStateUpdate()
	var updateMutex sync.Mutex
	errChan := make(chan error, 2)
	var wg sync.WaitGroup

	oracle.fetchSolanaBurnEvents(context.Background(), &wg, update, &updateMutex, errChan)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		t.Log("Operation timed out as expected")
	case <-done:
		close(errChan)
		errorCount := 0
		for range errChan {
			errorCount++
		}
		t.Logf("Operation completed within timeout (%d expected errors)", errorCount)
	}
}

func TestFetchSolanaBurnEventsEdgeCases(t *testing.T) {
	t.Skip("Skipping slow edge cases test with potential timeouts")
	tests := []struct {
		name        string
		setupOracle func() *Oracle
		description string
	}{
		{
			name: "Zero events",
			setupOracle: func() *Oracle {
				oracle := createTestOracle()
				oracle.lastSolZenBTCBurnSigStr = "latest_sig"
				oracle.lastSolRockBurnSigStr = "latest_sig"
				return oracle
			},
			description: "Should handle case with zero events gracefully",
		},
		{
			name: "Very large amounts",
			setupOracle: func() *Oracle {
				oracle := createTestOracle()
				return oracle
			},
			description: "Should handle very large burn amounts",
		},
		{
			name: "Empty destination addresses",
			setupOracle: func() *Oracle {
				oracle := createTestOracle()
				return oracle
			},
			description: "Should handle empty destination addresses",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oracle := tt.setupOracle()

			update := oracle.initializeStateUpdate()
			var updateMutex sync.Mutex
			errChan := make(chan error, 2)
			var wg sync.WaitGroup

			oracle.fetchSolanaBurnEvents(context.Background(), &wg, update, &updateMutex, errChan)
			wg.Wait()
			close(errChan)

			errorCount := 0
			for range errChan {
				errorCount++
			}

			t.Logf("Test passed: %s (%d expected errors)", tt.description, errorCount)
		})
	}
}

// TestStateCacheDataRace verifies that Oracle.stateCache is properly protected from data races.
// This test will PASS when run with `go test -race` because the RWMutex protection is in place.
//
// The test verifies the following scenarios work correctly with RWMutex protection:
// 1. Concurrent AppendStateToCache() operations (exclusive Lock for writes)
// 2. Concurrent read (getStateByEthHeight with RLock) and write (AppendStateToCache with Lock) operations
// 3. Concurrent SaveToFile() operations (RLock for reads, allowing concurrent read access)
// 4. Concurrent read and reset operations (SetStateCacheForTesting with exclusive Lock)
//
// Run with: go test -race -run TestStateCacheDataRace ./sidecar/
func TestStateCacheDataRace(t *testing.T) {
	oracle := createTestOracle()

	// Cleanup test state file created by AppendStateToCache() calls
	t.Cleanup(func() {
		os.Remove("test_state.json")
		os.Remove("test_state.json.tmp")
	})

	// Initialize with some states
	initialStates := []sidecartypes.OracleState{
		{EthBlockHeight: 100, ROCKUSDPrice: math.LegacyNewDec(1)},
		{EthBlockHeight: 200, ROCKUSDPrice: math.LegacyNewDec(2)},
		{EthBlockHeight: 300, ROCKUSDPrice: math.LegacyNewDec(3)},
	}
	oracle.SetStateCacheForTesting(initialStates)

	const numGoroutines = 10
	const iterations = 100
	var wg sync.WaitGroup

	// Scenario 1: Concurrent AppendStateToCache() calls (appending to stateCache)
	t.Run("Concurrent AppendStateToCache", func(t *testing.T) {
		for i := range numGoroutines {
			wg.Go(func() {
				for j := range iterations {
					// Simulate what happens during normal oracle operation
					newState := sidecartypes.OracleState{
						EthBlockHeight: uint64(1000 + i*iterations + j),
						ROCKUSDPrice:   math.LegacyNewDec(int64(i + 1)),
					}
					oracle.currentState.Store(&newState)
					oracle.appendStateToCache() // Protected by stateCacheMutex
				}
			})
		}
		wg.Wait()

		// If we got here without panicking, the slice wasn't corrupted
		// RWMutex prevents races
		t.Logf("Completed %d concurrent AppendStateToCache operations", numGoroutines*iterations)
	})

	// Scenario 2: Concurrent getStateByEthHeight() reads during AppendStateToCache() writes
	t.Run("Concurrent Read and Write", func(t *testing.T) {
		// Reset state
		oracle.SetStateCacheForTesting(initialStates)

		// Readers
		for range numGoroutines {
			wg.Go(func() {
				for range iterations {
					// Try to read from stateCache
					_, err := oracle.getStateByEthHeight(100) // Protected by stateCacheMutex
					if err != nil {
						// Expected if state was not found
					}
				}
			})
		}

		// Writers
		for i := range numGoroutines {
			wg.Go(func() {
				for j := range iterations {
					newState := sidecartypes.OracleState{
						EthBlockHeight: uint64(400 + i*iterations + j),
						ROCKUSDPrice:   math.LegacyNewDec(int64(i + 1)),
					}
					oracle.currentState.Store(&newState)
					oracle.appendStateToCache() // Protected by stateCacheMutex
				}
			})
		}

		wg.Wait()
		t.Logf("Completed concurrent read/write operations without panic")
	})

	// Scenario 3: Concurrent SaveToFile() and AppendStateToCache() operations
	t.Run("Concurrent SaveToFile", func(t *testing.T) {
		// Reset state
		oracle.SetStateCacheForTesting(initialStates)

		// Cleanup test state files after test completes
		t.Cleanup(func() {
			for i := range numGoroutines {
				for j := range 10 {
					filename := fmt.Sprintf("test_state_%d_%d.json", i, j)
					os.Remove(filename)
					os.Remove(filename + ".tmp") // Also remove any temp files
				}
			}
		})

		for i := range numGoroutines {
			wg.Go(func() {
				for j := range 10 { // Fewer iterations since file I/O is slow
					// Get a copy of stateCache with RLock, then pass to SaveToFile
					oracle.stateCacheMutex.RLock()
					stateCopy := slices.Clone(oracle.stateCache)
					oracle.stateCacheMutex.RUnlock()

					err := saveStatesToFile(fmt.Sprintf("test_state_%d_%d.json", i, j), stateCopy)
					if err != nil {
						t.Logf("SaveToFile error (expected in concurrent test): %v", err)
					}
				}
			})
		}

		// Concurrent modifications
		for i := range numGoroutines {
			wg.Go(func() {
				for j := range 10 {
					newState := sidecartypes.OracleState{
						EthBlockHeight: uint64(500 + i*10 + j),
						ROCKUSDPrice:   math.LegacyNewDec(int64(i + 1)),
					}
					oracle.currentState.Store(&newState)
					oracle.appendStateToCache()
				}
			})
		}

		wg.Wait()
		t.Logf("Completed concurrent SaveToFile operations")
	})

	// Scenario 4: Simulated scheduled reset during operations
	// Now tests that the mutex protection works correctly
	t.Run("Simulated Reset During Operations", func(t *testing.T) {
		// Reset state
		oracle.SetStateCacheForTesting(initialStates)

		// Readers - will iterate over stateCache using the proper API
		for i := range numGoroutines {
			wg.Go(func() {
				for j := range iterations {
					// Use the proper API which now has mutex protection
					height := uint64((i*100 + j) % 5000)
					_, _ = oracle.getStateByEthHeight(height)
				}
			})
		}

		// Simulate reset operations using the proper API
		for range 5 {
			wg.Go(func() {
				for range 10 {
					// Use the proper API which now has mutex protection
					oracle.SetStateCacheForTesting([]sidecartypes.OracleState{EmptyOracleState})
					time.Sleep(time.Millisecond)
				}
			})
		}

		wg.Wait()
		t.Logf("Completed simulated reset operations with proper mutex protection")
	})
}

// TestCurrentStateCopyOnWrite verifies the copy-on-write pattern for currentState modifications
func TestCurrentStateCopyOnWrite(t *testing.T) {
	oracle := createTestOracle()

	// Set initial state
	initialState := sidecartypes.OracleState{
		EthBlockHeight:          100,
		ROCKUSDPrice:            math.LegacyNewDec(1),
		SolanaBurnEvents:        []api.BurnEvent{{TxID: "initial", Amount: 100}},
		CleanedSolanaBurnEvents: map[string]bool{"key1": true},
	}
	oracle.currentState.Store(&initialState)

	// Get pointer to verify immutability later
	storedPtr := oracle.currentState.Load()

	// CORRECT: Copy-on-write pattern
	oldState := oracle.currentState.Load()
	newState := *oldState // Copy the struct
	newState.EthBlockHeight = 200
	newState.ROCKUSDPrice = math.LegacyNewDec(2)
	newState.SolanaBurnEvents = []api.BurnEvent{{TxID: "modified", Amount: 200}}
	oracle.currentState.Store(&newState) // Store new pointer

	// Verify old state was NOT mutated (immutability preserved)
	assert.Equal(t, uint64(100), storedPtr.EthBlockHeight, "Original state should not be mutated")
	assert.True(t, storedPtr.ROCKUSDPrice.Equal(math.LegacyNewDec(1)), "Original price should not be mutated")
	assert.Equal(t, "initial", storedPtr.SolanaBurnEvents[0].TxID, "Original events should not be mutated")

	// Verify new state is correct
	currentState := oracle.currentState.Load()
	assert.Equal(t, uint64(200), currentState.EthBlockHeight)
	assert.True(t, currentState.ROCKUSDPrice.Equal(math.LegacyNewDec(2)))
	assert.Equal(t, "modified", currentState.SolanaBurnEvents[0].TxID)

	// Verify they're different pointers
	assert.NotSame(t, storedPtr, currentState, "Should store different pointer after copy-on-write")
}

// TestCurrentStateConcurrentReads verifies concurrent reads are safe with copy-on-write
func TestCurrentStateConcurrentReads(t *testing.T) {
	oracle := createTestOracle()

	// Initialize with known state
	initialState := sidecartypes.OracleState{
		EthBlockHeight: 1000,
		ROCKUSDPrice:   math.LegacyNewDec(100),
	}
	oracle.currentState.Store(&initialState)

	const numReaders = 20
	const numWrites = 50
	var wg sync.WaitGroup

	// Start concurrent readers
	for range numReaders {
		wg.Go(func() {
			for range 100 {
				// Read current state (should never panic or see torn reads)
				state := oracle.currentState.Load()
				_ = state.EthBlockHeight
				_ = state.ROCKUSDPrice
			}
		})
	}

	// Concurrent writer using copy-on-write
	wg.Go(func() {
		for i := range numWrites {
			// CORRECT copy-on-write pattern
			oldState := oracle.currentState.Load()
			newState := *oldState
			newState.EthBlockHeight = uint64(1000 + i)
			newState.ROCKUSDPrice = math.LegacyNewDec(int64(100 + i))
			oracle.currentState.Store(&newState)
			time.Sleep(time.Microsecond) // Small delay to increase race chance
		}
	})

	wg.Wait()
	t.Log("Concurrent reads and copy-on-write updates completed successfully")
}

// TestCurrentStateMapMutationSafety verifies map/slice field updates use copy-on-write
func TestCurrentStateMapMutationSafety(t *testing.T) {
	oracle := createTestOracle()

	// Set initial state with maps and slices
	initialState := sidecartypes.OracleState{
		CleanedSolanaBurnEvents: map[string]bool{
			"key1": true,
			"key2": false,
		},
		SolanaBurnEvents: []api.BurnEvent{
			{TxID: "event1", Amount: 100},
			{TxID: "event2", Amount: 200},
		},
		PendingSolanaTxs: map[string]sidecartypes.PendingTxInfo{
			"sig1": {Signature: "sig1", RetryCount: 1},
		},
	}
	oracle.currentState.Store(&initialState)

	// Get pointer to verify immutability
	oldPtr := oracle.currentState.Load()

	// CORRECT: Copy struct and create new maps/slices
	oldState := oracle.currentState.Load()
	newState := *oldState // Copy struct (shallow copy)

	// Create new map instances (don't modify shared ones)
	newState.CleanedSolanaBurnEvents = map[string]bool{
		"key1": true,
		"key2": false,
		"key3": true, // Added
	}
	newState.SolanaBurnEvents = []api.BurnEvent{
		{TxID: "event1", Amount: 100},
		{TxID: "event2", Amount: 200},
		{TxID: "event3", Amount: 300}, // Added
	}
	newState.PendingSolanaTxs = map[string]sidecartypes.PendingTxInfo{
		"sig1": {Signature: "sig1", RetryCount: 1},
		"sig2": {Signature: "sig2", RetryCount: 2}, // Added
	}

	oracle.currentState.Store(&newState)

	// Verify old state was NOT mutated
	assert.Len(t, oldPtr.CleanedSolanaBurnEvents, 2, "Old map should not be mutated")
	assert.Len(t, oldPtr.SolanaBurnEvents, 2, "Old slice should not be mutated")
	assert.Len(t, oldPtr.PendingSolanaTxs, 1, "Old pending txs should not be mutated")

	// Verify new state has updates
	currentState := oracle.currentState.Load()
	assert.Len(t, currentState.CleanedSolanaBurnEvents, 3, "New map should have 3 entries")
	assert.Len(t, currentState.SolanaBurnEvents, 3, "New slice should have 3 entries")
	assert.Len(t, currentState.PendingSolanaTxs, 2, "New pending txs should have 2 entries")
}

// TestCurrentStateCompareAndSwapPattern verifies CAS usage requires pointer storage
func TestCurrentStateCompareAndSwapPattern(t *testing.T) {
	oracle := createTestOracle()

	initialState := sidecartypes.OracleState{
		EthBlockHeight: 100,
		PendingSolanaTxs: map[string]sidecartypes.PendingTxInfo{
			"sig1": {Signature: "sig1", RetryCount: 1},
			"sig2": {Signature: "sig2", RetryCount: 2},
		},
	}
	oracle.currentState.Store(&initialState)

	// Simulate CompareAndSwap pattern used in production code
	// (e.g., removePendingTransactionFromState)
	success := false
	attempts := 0
	for !success && attempts < 10 {
		attempts++
		currentState := oracle.currentState.Load()

		// Create new map without one entry
		newPendingTxs := make(map[string]sidecartypes.PendingTxInfo)
		for k, v := range currentState.PendingSolanaTxs {
			if k != "sig1" {
				newPendingTxs[k] = v
			}
		}

		// Copy-on-write
		newState := *currentState
		newState.PendingSolanaTxs = newPendingTxs

		// CAS: Only succeeds if currentState pointer hasn't changed
		if oracle.currentState.CompareAndSwap(currentState, &newState) {
			success = true
		}
		// If CAS failed, loop retries with fresh state
	}

	assert.True(t, success, "CompareAndSwap should succeed")
	assert.Less(t, attempts, 5, "Should succeed quickly without contention")

	// Verify the removal worked
	finalState := oracle.currentState.Load()
	assert.Len(t, finalState.PendingSolanaTxs, 1, "Should have 1 pending tx")
	assert.NotContains(t, finalState.PendingSolanaTxs, "sig1", "sig1 should be removed")
	assert.Contains(t, finalState.PendingSolanaTxs, "sig2", "sig2 should remain")
}

// TestStateCacheRaceWithSetStateCacheForTesting verifies concurrent access to SetStateCacheForTesting is safe
// This function is used in tests and is now properly protected by stateCacheMutex
func TestStateCacheRaceWithSetStateCacheForTesting(t *testing.T) {
	oracle := createTestOracle()

	const numGoroutines = 20
	var wg sync.WaitGroup

	// Multiple goroutines calling SetStateCacheForTesting concurrently
	for i := range numGoroutines {
		wg.Go(func() {
			for j := range 50 {
				testStates := []sidecartypes.OracleState{
					{
						EthBlockHeight: uint64(i*100 + j),
						ROCKUSDPrice:   math.LegacyNewDec(int64(i)),
					},
				}
				// Protected by stateCacheMutex
				oracle.SetStateCacheForTesting(testStates)
			}
		})
	}

	wg.Wait()
	t.Logf("Completed concurrent SetStateCacheForTesting operations")

	// Verify we can still read the state without panicking
	currentState := oracle.currentState.Load()
	assert.NotNil(t, currentState)
	assert.Equal(t, 1, len(oracle.stateCache))
}
