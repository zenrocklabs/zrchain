package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"testing"
	"time"

	"cosmossdk.io/math"
	sidecar "github.com/Zenrock-Foundation/zrchain/v5/sidecar"
	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const (
	// bufSize is the size of the buffer for the gRPC connection
	bufSize = 1024 * 1024
	// ethBlockHeight is the Ethereum block height to query in the test
	ethBlockHeight uint64 = 3419775
)

// OracleData represents the processed data from the sidecar service
type OracleData struct {
	EigenDelegationsMap        map[string]map[string]*big.Int
	ValidatorDelegations       []ValidatorDelegations
	EthBlockHeight             uint64
	EthGasLimit                uint64
	EthBaseFee                 uint64
	EthTipCap                  uint64
	SolanaLamportsPerSignature uint64
	EthBurnEvents              []api.BurnEvent
	Redemptions                []api.Redemption
	ROCKUSDPrice               math.LegacyDec
	BTCUSDPrice                math.LegacyDec
	ETHUSDPrice                math.LegacyDec
}

// ValidatorDelegations represents the delegations for a validator
type ValidatorDelegations struct {
	Validator string
	Stake     math.Int
}

// TestGetSidecarStateByEthHeight tests the GetSidecarStateByEthHeight method
// with a configurable height variable
func TestGetSidecarStateByEthHeight(t *testing.T) {
	// Skip this test in short mode as it requires external connections
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	// Initialize the oracle using the existing test function
	oracle := initTestOracle()

	// Setup a buffer connection for gRPC
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()

	// Create the actual service implementation
	// This will use the real implementation of GetSidecarStateByEthHeight
	service := sidecar.NewOracleService(oracle)

	// Register the service with our oracle
	api.RegisterSidecarServiceServer(s, service)

	// Start the server
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Errorf("Server exited with error: %v", err)
		}
	}()
	defer s.Stop()

	// Create a client connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	// Create a client
	client := api.NewSidecarServiceClient(conn)

	// Create a request with the configured height
	req := &api.SidecarStateByEthHeightRequest{
		EthBlockHeight: ethBlockHeight,
	}

	// Call the service method
	response, err := client.GetSidecarStateByEthHeight(ctx, req)
	require.NoError(t, err, "GetSidecarStateByEthHeight should not return an error")
	require.NotNil(t, response, "Response should not be nil")

	// Process the response into OracleData format (similar to processOracleResponse in abci_utils.go)
	oracleData, err := processOracleResponse(response)
	require.NoError(t, err, "Failed to process oracle response")

	// Log the processed OracleData instead of the raw response
	t.Logf("Processed OracleData: %+v", oracleData)
	t.Logf("EthBlockHeight: %d", oracleData.EthBlockHeight)
	t.Logf("EthGasLimit: %d", oracleData.EthGasLimit)
	t.Logf("EthBaseFee: %d", oracleData.EthBaseFee)
	t.Logf("EthTipCap: %d", oracleData.EthTipCap)
	t.Logf("ROCKUSDPrice: %s", oracleData.ROCKUSDPrice.String())
	t.Logf("BTCUSDPrice: %s", oracleData.BTCUSDPrice.String())
	t.Logf("ETHUSDPrice: %s", oracleData.ETHUSDPrice.String())

	if len(oracleData.ValidatorDelegations) > 0 {
		t.Logf("Validator Delegations:")
		for i, vd := range oracleData.ValidatorDelegations {
			t.Logf("  %d. Validator: %s, Stake: %s", i+1, vd.Validator, vd.Stake.String())
		}
	}

	// Verify the response contains the expected data
	if ethBlockHeight > 0 {
		// If we specified a specific height, check that it matches
		assert.Equal(t, ethBlockHeight, response.EthBlockHeight, "Response should contain the requested block height")
		t.Logf("Testing with specific block height: %d", ethBlockHeight)
	} else {
		// For height 0, we're testing with whatever state is returned
		// This could be the latest state or an empty state depending on the implementation
		t.Logf("Testing with height 0, received block height: %d", response.EthBlockHeight)
	}

	// Additional assertions based on expected state
	// Note: In a test environment, some of these values might be empty or zero
	// So we'll make the assertions conditional
	if response.EthBlockHeight > 0 {
		assert.Greater(t, response.EthGasLimit, uint64(0), "EthGasLimit should be greater than 0")
		assert.Greater(t, response.EthBaseFee, uint64(0), "EthBaseFee should be greater than 0")
	}

	// Log the response details
	t.Logf("Successfully retrieved state for Ethereum block height: %d", response.EthBlockHeight)
	t.Logf("Gas limit: %d, Base fee: %d", response.EthGasLimit, response.EthBaseFee)
	t.Logf("ROCK/USD Price: %s, BTC/USD Price: %s, ETH/USD Price: %s",
		response.ROCKUSDPrice, response.BTCUSDPrice, response.ETHUSDPrice)
}

// processOracleResponse processes the response from the sidecar service into OracleData format
// This is based on the processOracleResponse function in abci_utils.go
func processOracleResponse(resp *api.SidecarStateResponse) (*OracleData, error) {
	var delegations map[string]map[string]*big.Int

	if err := json.Unmarshal(resp.EigenDelegations, &delegations); err != nil {
		return nil, err
	}

	validatorDelegations, err := processDelegations(delegations)
	if err != nil {
		return nil, fmt.Errorf("error processing delegations: %w", err)
	}

	ROCKUSDPrice, err := math.LegacyNewDecFromStr(resp.ROCKUSDPrice)
	if err != nil {
		return nil, fmt.Errorf("error parsing rock price: %w", err)
	}

	BTCUSDPrice, err := math.LegacyNewDecFromStr(resp.BTCUSDPrice)
	if err != nil {
		return nil, fmt.Errorf("error parsing btc price: %w", err)
	}

	ETHUSDPrice, err := math.LegacyNewDecFromStr(resp.ETHUSDPrice)
	if err != nil {
		return nil, fmt.Errorf("error parsing eth price: %w", err)
	}

	return &OracleData{
		EigenDelegationsMap:        delegations,
		ValidatorDelegations:       validatorDelegations,
		EthBlockHeight:             resp.EthBlockHeight,
		EthGasLimit:                resp.EthGasLimit,
		EthBaseFee:                 resp.EthBaseFee,
		EthTipCap:                  resp.EthTipCap,
		SolanaLamportsPerSignature: resp.SolanaLamportsPerSignature,
		EthBurnEvents:              resp.EthBurnEvents,
		Redemptions:                resp.Redemptions,
		ROCKUSDPrice:               ROCKUSDPrice,
		BTCUSDPrice:                BTCUSDPrice,
		ETHUSDPrice:                ETHUSDPrice,
	}, nil
}

// processDelegations processes the delegations map into a slice of ValidatorDelegations
// This is based on the processDelegations function in abci_utils.go
func processDelegations(delegations map[string]map[string]*big.Int) ([]ValidatorDelegations, error) {
	validatorTotals := make(map[string]*big.Int)
	for validator, delegatorMap := range delegations {
		total := new(big.Int)
		for _, amount := range delegatorMap {
			total.Add(total, amount)
		}
		validatorTotals[validator] = total
	}

	validatorDelegations := make([]ValidatorDelegations, 0, len(validatorTotals))
	for validator, totalStake := range validatorTotals {
		validatorDelegations = append(validatorDelegations, ValidatorDelegations{
			Validator: validator,
			Stake:     math.NewIntFromBigInt(totalStake),
		})
	}

	return validatorDelegations, nil
}
