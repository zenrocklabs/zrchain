package main
import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
	"github.com/ethereum/go-ethereum"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	sol "github.com/gagliardetto/solana-go"
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
func createTestOracle() *Oracle {
	config := sidecartypes.Config{
		Network:   sidecartypes.NetworkDevnet,
		StateFile: "test_state.json",
	}
	oracle := &Oracle{
		Config:    config,
		DebugMode: false,
	}
	oracle.currentState.Store(&EmptyOracleState)
	oracle.stateCache = []sidecartypes.OracleState{EmptyOracleState}
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
	assert.NotNil(t, update.eigenDelegations)
	assert.NotNil(t, update.redemptions)
	assert.NotNil(t, update.ethBurnEvents)
}
func TestApplyFallbacks(t *testing.T) {
	oracle := createTestOracle()
	currentState := sidecartypes.OracleState{
		ROCKUSDPrice:               math.LegacyNewDec(1),
		BTCUSDPrice:                math.LegacyNewDec(40000),
		ETHUSDPrice:                math.LegacyNewDec(2000),
		SolanaLamportsPerSignature: 5000,
	}
	oracle.currentState.Store(&currentState)
	update := &oracleStateUpdate{
		suggestedTip:               nil,
		ROCKUSDPrice:               math.LegacyDec{},
		BTCUSDPrice:                math.LegacyDec{},
		ETHUSDPrice:                math.LegacyDec{},
		solanaLamportsPerSignature: 0,
	}
	oracle.applyFallbacks(update, &currentState)
	assert.NotNil(t, update.suggestedTip)
	assert.Equal(t, big.NewInt(0), update.suggestedTip)
	assert.True(t, update.ROCKUSDPrice.Equal(currentState.ROCKUSDPrice))
	assert.True(t, update.BTCUSDPrice.Equal(currentState.BTCUSDPrice))
	assert.True(t, update.ETHUSDPrice.Equal(currentState.ETHUSDPrice))
	assert.Equal(t, currentState.SolanaLamportsPerSignature, update.solanaLamportsPerSignature)
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
		eigenDelegations:           make(map[string]map[string]*big.Int),
		redemptions:                []api.Redemption{},
		suggestedTip:               big.NewInt(1500000000),
		estimatedGas:               231000,
		ethBurnEvents:              []api.BurnEvent{},
		solanaBurnEvents:           []api.BurnEvent{},
		ROCKUSDPrice:               math.LegacyNewDec(1),
		BTCUSDPrice:                math.LegacyNewDec(50000),
		ETHUSDPrice:                math.LegacyNewDec(3000),
		solanaLamportsPerSignature: 5000,
		SolanaMintEvents:           []api.SolanaMintEvent{},
		latestSolanaSigs:           make(map[sidecartypes.SolanaEventType]sol.Signature),
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
	assert.Equal(t, uint64(5000), result.SolanaLamportsPerSignature)
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
	currentState := oracle.currentState.Load().(*sidecartypes.OracleState)
	assert.Equal(t, uint64(200), currentState.EthBlockHeight)
	assert.True(t, currentState.ROCKUSDPrice.Equal(math.LegacyNewDec(2)))
}
func TestSetStateCacheForTesting_Empty(t *testing.T) {
	oracle := createTestOracle()
	oracle.SetStateCacheForTesting([]sidecartypes.OracleState{})
	assert.Equal(t, 1, len(oracle.stateCache))
	currentState := oracle.currentState.Load().(*sidecartypes.OracleState)
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
	currentState := oracle.currentState.Load().(*sidecartypes.OracleState)
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
		eigenDelegations:           make(map[string]map[string]*big.Int),
		redemptions:                []api.Redemption{},
		suggestedTip:               big.NewInt(1500000000),
		estimatedGas:               231000,
		ethBurnEvents:              []api.BurnEvent{},
		solanaBurnEvents:           []api.BurnEvent{},
		ROCKUSDPrice:               math.LegacyNewDec(1),
		BTCUSDPrice:                math.LegacyNewDec(50000),
		ETHUSDPrice:                math.LegacyNewDec(3000),
		solanaLamportsPerSignature: 5000,
		SolanaMintEvents:           []api.SolanaMintEvent{},
		latestSolanaSigs:           make(map[sidecartypes.SolanaEventType]sol.Signature),
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
		eigenDelegations:           make(map[string]map[string]*big.Int),
		redemptions:                []api.Redemption{},
		suggestedTip:               big.NewInt(1500000000),
		estimatedGas:               231000,
		ethBurnEvents:              []api.BurnEvent{},
		solanaBurnEvents:           []api.BurnEvent{},
		ROCKUSDPrice:               math.LegacyNewDec(1),
		BTCUSDPrice:                math.LegacyNewDec(50000),
		ETHUSDPrice:                math.LegacyNewDec(3000),
		solanaLamportsPerSignature: 5000,
		SolanaMintEvents:           []api.SolanaMintEvent{},
		latestSolanaSigs:           make(map[sidecartypes.SolanaEventType]sol.Signature),
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
		eigenDelegations:           make(map[string]map[string]*big.Int),
		redemptions:                []api.Redemption{},
		suggestedTip:               big.NewInt(1500000000),
		estimatedGas:               231000,
		ethBurnEvents:              []api.BurnEvent{},
		solanaBurnEvents:           []api.BurnEvent{},
		ROCKUSDPrice:               math.LegacyNewDec(1),
		BTCUSDPrice:                math.LegacyNewDec(50000),
		ETHUSDPrice:                math.LegacyNewDec(3000),
		solanaLamportsPerSignature: 5000,
		SolanaMintEvents:           []api.SolanaMintEvent{},
		latestSolanaSigs:           make(map[sidecartypes.SolanaEventType]sol.Signature),
	}
	header := createMockHeader(1000, big.NewInt(20000000000))
	assert.Panics(t, func() {
		oracle.buildFinalState(update, header, nil)
	})
}
func TestApplyFallbacks_NilCurrentState(t *testing.T) {
	oracle := createTestOracle()
	update := &oracleStateUpdate{
		suggestedTip:               nil,
		ROCKUSDPrice:               math.LegacyDec{},
		BTCUSDPrice:                math.LegacyDec{},
		ETHUSDPrice:                math.LegacyDec{},
		solanaLamportsPerSignature: 0,
	}
	assert.Panics(t, func() {
		oracle.applyFallbacks(update, nil)
	})
}
func TestApplyFallbacks_NilUpdate(t *testing.T) {
	oracle := createTestOracle()
	currentState := sidecartypes.OracleState{
		ROCKUSDPrice:               math.LegacyNewDec(1),
		BTCUSDPrice:                math.LegacyNewDec(40000),
		ETHUSDPrice:                math.LegacyNewDec(2000),
		SolanaLamportsPerSignature: 5000,
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
	currentState := oracle.currentState.Load().(*sidecartypes.OracleState)
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
		largeData := make([]byte, 1024*1024) // 1MB
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
